package loop_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func evalGate(t *testing.T, st *store.Store, psvc *planner.Service, dsvc *domain.Service, taskID string, gateFor loop.GateFor) (bool, []loop.Violation) {
	t.Helper()
	ctx := context.Background()
	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, gateFor)
	if err != nil {
		t.Fatalf("EvaluateGate(%s): %v", gateFor, err)
	}
	return allowed, violations
}

func assertBlocked(t *testing.T, allowed bool, violations []loop.Violation, code, reasonCode string) {
	t.Helper()
	if allowed {
		t.Fatalf("allowed=true want false; violations=%+v", violations)
	}
	if len(violations) != 1 {
		t.Fatalf("violations len=%d want 1: %+v", len(violations), violations)
	}
	v := violations[0]
	if v.Code != code {
		t.Fatalf("code=%q want %q", v.Code, code)
	}
	if reasonCode != "" && v.ReasonCode != reasonCode {
		t.Fatalf("reason_code=%q want %q; violation=%+v", v.ReasonCode, reasonCode, v)
	}
}

func assertAllowed(t *testing.T, allowed bool, violations []loop.Violation) {
	t.Helper()
	if !allowed {
		t.Fatalf("allowed=false want true; violations=%+v", violations)
	}
	if len(violations) != 0 {
		t.Fatalf("violations=%+v want none", violations)
	}
}

func seedGoalTaskNoPlan(t *testing.T, psvc *planner.Service, dsvc *domain.Service) (goalID, taskID string) {
	t.Helper()
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "no plan goal"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "no plan task", GoalID: &g.ID})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	return g.ID, task.ID
}

func seedFullCycleClear(t *testing.T, psvc *planner.Service, dsvc *domain.Service, st *store.Store) (goalID, taskID string) {
	t.Helper()
	ctx := context.Background()
	goalID, taskID, _ = seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
	if _, err := dsvc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskID, TestName: "TestAll", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatalf("RecordTestOutcome: %v", err)
	}
	ev, err := dsvc.CreateEvidence(ctx, domain.EvidenceInput{Title: "verify proof"})
	if err != nil {
		t.Fatalf("CreateEvidence: %v", err)
	}
	if _, err := dsvc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             taskID,
		GoalID:             goalID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatalf("RecordVerificationOutcome: %v", err)
	}
	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatalf("CreateBaseline: %v", err)
	}
	if _, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.95}`,
	}); err != nil {
		t.Fatalf("RecordEvaluationOutcome: %v", err)
	}
	if _, err := dsvc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskID, Summary: "cycle complete", UsefulTests: []string{"TestAll"},
	}); err != nil {
		t.Fatalf("CreateReflection: %v", err)
	}
	return goalID, taskID
}

func TestEvaluateGate_Orient_UnknownTask(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	allowed, violations := evalGate(t, st, psvc, dsvc, uuid.NewString(), loop.GateForOrient)
	assertBlocked(t, allowed, violations, "gate_orient_failed", "task_not_found")
}

func TestEvaluateGate_Orient_MissingGoal(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "orphan"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	allowed, violations := evalGate(t, st, psvc, dsvc, task.ID, loop.GateForOrient)
	assertBlocked(t, allowed, violations, "gate_orient_failed", "missing_goal_id")
}

func TestEvaluateGate_Orient_MissingPlan(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID := seedGoalTaskNoPlan(t, psvc, dsvc)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForOrient)
	assertBlocked(t, allowed, violations, "gate_orient_failed", "missing_plan_context")
}

func TestEvaluateGate_Orient_OK(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForOrient)
	assertAllowed(t, allowed, violations)
}

func TestEvaluateGate_Edit_BlockingUncertainty(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "blocker", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "blocking_uncertainty")
	if violations[0].RecommendedPhase != string(deliberation.PhaseInvestigate) {
		t.Fatalf("recommended_phase=%q want INVESTIGATE", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Edit_OpenRegression(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.9}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.5}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskID,
	}); err != nil {
		t.Fatalf("RecordRegressionFromEvaluation: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "open_regression")
	if violations[0].RecommendedPhase != string(deliberation.PhaseInvestigate) {
		t.Fatalf("recommended_phase=%q want INVESTIGATE", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Edit_PlanMissing(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID := seedGoalTaskNoPlan(t, psvc, dsvc)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "plan_missing")
	if violations[0].RecommendedPhase != string(deliberation.PhasePlan) {
		t.Fatalf("recommended_phase=%q want PLAN", violations[0].RecommendedPhase)
	}
}

func TestActiveWork_PlanMissingStillBlocksEdit(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "active goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "active task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{Actor: "test", Reason: "start"}); err != nil {
		t.Fatal(err)
	}
	allowed, violations := evalGate(t, st, psvc, dsvc, task.ID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "plan_missing")
	if violations[0].RecommendedPhase != string(deliberation.PhasePlan) {
		t.Fatalf("recommended_phase=%q want PLAN", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Done_TerminalPlanGapAdvisory(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "terminal goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "terminal task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	seedTerminalDoneTask(t, st, task.ID)
	allowed, violations := evalGate(t, st, psvc, dsvc, task.ID, loop.GateForDone)
	if !allowed {
		t.Fatalf("terminal DONE want allowed=true; violations=%+v", violations)
	}
	if len(violations) != 1 || violations[0].ReasonCode != "goal_plan_gap_terminal_advisory" {
		t.Fatalf("violations=%+v want goal_plan_gap_terminal_advisory", violations)
	}
}

func TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap(t *testing.T) {
	st, psvc, dsvc := importFeetExportFixture(t)
	ctx := context.Background()
	goalID := "a1111111-1111-4111-8111-111111111111"
	taskID := "b2222222-2222-4222-8222-222222222222"
	seedTerminalDoneTask(t, st, taskID)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForDone)
	if !allowed {
		t.Fatalf("pre-bootstrap DONE gate want allowed; violations=%+v", violations)
	}
	if len(violations) != 1 || violations[0].ReasonCode != "goal_plan_gap_terminal_advisory" {
		t.Fatalf("pre-bootstrap violations=%+v", violations)
	}
	if _, err := psvc.BootstrapFromPlanChanges(ctx, goalID, "test"); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	ok, err := psvc.PlanExists(ctx, goalID)
	if err != nil || !ok {
		t.Fatalf("PlanExists after bootstrap: ok=%v err=%v", ok, err)
	}
	markPlanCritiqued(t, st, taskID, goalID)
	editAllowed, editViolations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	if !editAllowed {
		t.Fatalf("post-bootstrap edit gate want allowed; violations=%+v", editViolations)
	}
}

func TestEvaluateGate_Edit_PlanUncritiqued(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "plan_uncritiqued")
	if violations[0].RecommendedPhase != string(deliberation.PhaseCritique) {
		t.Fatalf("recommended_phase=%q want CRITIQUE", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Edit_ExecuteReady(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertAllowed(t, allowed, violations)
}

func TestEvaluateGate_Edit_HopBudgetStopped(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := st.UpsertDeliberationState(store.DeliberationState{
		TaskID:        taskID,
		GoalID:        goalID,
		CurrentPhase:  string(deliberation.PhaseExecute),
		HopCount:      deliberation.HopBudget,
		PlanCritiqued: true,
	}); err != nil {
		t.Fatalf("UpsertDeliberationState: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForEdit)
	assertBlocked(t, allowed, violations, "premature_implementation", "hop_budget_exceeded")
	if violations[0].RecommendedPhase != string(deliberation.PhaseStop) {
		t.Fatalf("recommended_phase=%q want STOP", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Execute_NotExecutePending(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForExecute)
	assertBlocked(t, allowed, violations, "premature_implementation", "test_pending")
}

func TestEvaluateGate_Execute_ExecutePendingClear(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForExecute)
	assertAllowed(t, allowed, violations)
}

func TestEvaluateGate_Done_VerificationDebt(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForDone)
	assertBlocked(t, allowed, violations, "premature_implementation", "verification_incomplete")
}

func TestEvaluateGate_Done_OpenRegression(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.9}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.5}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskID,
	}); err != nil {
		t.Fatalf("RecordRegressionFromEvaluation: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForDone)
	assertBlocked(t, allowed, violations, "premature_implementation", "open_regression")
}

func TestEvaluateGate_Done_DeliberationIncomplete(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
	if _, err := dsvc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskID, TestName: "TestGate", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatalf("RecordTestOutcome: %v", err)
	}
	ev, err := dsvc.CreateEvidence(ctx, domain.EvidenceInput{Title: "proof"})
	if err != nil {
		t.Fatalf("CreateEvidence: %v", err)
	}
	if _, err := dsvc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             taskID,
		GoalID:             goalID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatalf("RecordVerificationOutcome: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForDone)
	assertBlocked(t, allowed, violations, "premature_implementation", "evaluation_pending")
	if violations[0].RecommendedPhase != string(deliberation.PhaseEvaluate) {
		t.Fatalf("recommended_phase=%q want EVALUATE", violations[0].RecommendedPhase)
	}
}

func TestEvaluateGate_Done_Clean(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID := seedFullCycleClear(t, psvc, dsvc, st)
	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForDone)
	assertAllowed(t, allowed, violations)
}

func TestEvaluateGate_Export_SameAsDone(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}

	allowed, violations := evalGate(t, st, psvc, dsvc, taskID, loop.GateForExport)
	assertBlocked(t, allowed, violations, "premature_implementation", "verification_incomplete")

	st2, psvc2, dsvc2 := openLoopTestStore(t)
	_, cleanTaskID := seedFullCycleClear(t, psvc2, dsvc2, st2)
	allowedClean, violationsClean := evalGate(t, st2, psvc2, dsvc2, cleanTaskID, loop.GateForExport)
	assertAllowed(t, allowedClean, violationsClean)
}
