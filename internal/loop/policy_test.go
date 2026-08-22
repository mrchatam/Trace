package loop_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/store"
)

func markPlanCritiqued(t *testing.T, st *store.Store, taskID, goalID string) {
	t.Helper()
	if _, err := st.UpsertDeliberationState(store.DeliberationState{
		TaskID:        taskID,
		GoalID:        goalID,
		CurrentPhase:  "CRITIQUE",
		PlanCritiqued: true,
	}); err != nil {
		t.Fatalf("UpsertDeliberationState: %v", err)
	}
}

func TestBuildPolicyInputsSetsExecutePending(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	inputs, err := loop.BuildPolicyInputs(ctx, dsvc, psvc, taskID, goalID, nil, false)
	if err != nil {
		t.Fatalf("BuildPolicyInputs: %v", err)
	}
	if !inputs.PlanExists || !inputs.PlanCritiqued {
		t.Fatalf("plan flags: exists=%v critiqued=%v", inputs.PlanExists, inputs.PlanCritiqued)
	}
	if inputs.BlockingUncertaintyCount != 0 || inputs.OpenRegression {
		t.Fatalf("gates: blocking=%d regression=%v", inputs.BlockingUncertaintyCount, inputs.OpenRegression)
	}
	if !inputs.ExecutePending {
		t.Fatalf("execute_pending want true, got %#v", inputs)
	}
	if inputs.ReplanNeeded || inputs.OpenDecisionAlternatives != 0 {
		t.Fatalf("stubs must stay zero: replan=%v alternatives=%d", inputs.ReplanNeeded, inputs.OpenDecisionAlternatives)
	}

	// applyWrites still treats plan_changes as plan_critiqued before store read.
	st2, psvc2, dsvc2 := openLoopTestStore(t)
	goalID2, taskID2, _ := seedGoalTaskPlan(t, psvc2, dsvc2)
	forced, err := loop.BuildPolicyInputs(ctx, dsvc2, psvc2, taskID2, goalID2, &loop.ApplyWrites{
		PlanChanges: []loop.ApplyPlanChange{{ID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Title: "critique"}},
	}, false)
	if err != nil {
		t.Fatalf("BuildPolicyInputs applyWrites: %v", err)
	}
	if !forced.PlanCritiqued || !forced.ExecutePending {
		t.Fatalf("applyWrites path: %#v", forced)
	}
	_ = st2
}

func TestBuildPolicyInputsSetsTestPendingAfterChange(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}

	inputs, err := loop.BuildPolicyInputs(ctx, dsvc, psvc, taskID, goalID, nil, false)
	if err != nil {
		t.Fatalf("BuildPolicyInputs: %v", err)
	}
	if inputs.ExecutePending {
		t.Fatalf("implementation signal must clear execute_pending: %#v", inputs)
	}
	if !inputs.TestPending {
		t.Fatalf("test_pending want true after RECORDED change with no test: %#v", inputs)
	}
	if !inputs.VerificationIncomplete {
		t.Fatalf("verification debt must stay wired: %#v", inputs)
	}
}

func TestBuildPolicyInputsSetsEvaluationPending(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	ev, err := dsvc.CreateEvidence(ctx, domain.EvidenceInput{Title: "verify log"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             taskID,
		GoalID:             goalID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatalf("RecordVerificationOutcome: %v", err)
	}

	inputs, err := loop.BuildPolicyInputs(ctx, dsvc, psvc, taskID, goalID, nil, false)
	if err != nil {
		t.Fatalf("BuildPolicyInputs: %v", err)
	}
	if !inputs.EvaluationPending {
		t.Fatalf("evaluation_pending want true (verified, no eval): %#v", inputs)
	}
	if inputs.ReflectPending {
		t.Fatalf("reflect_pending want false without computed eval: %#v", inputs)
	}
}

func TestBuildPolicyInputsSetsReflectPending(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     taskID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness":0.95}`,
	}); err != nil {
		t.Fatalf("RecordEvaluationOutcome: %v", err)
	}

	inputs, err := loop.BuildPolicyInputs(ctx, dsvc, psvc, taskID, goalID, nil, false)
	if err != nil {
		t.Fatalf("BuildPolicyInputs: %v", err)
	}
	if !inputs.ReflectPending {
		t.Fatalf("reflect_pending want true (computed eval, no reflection): %#v", inputs)
	}
	if inputs.ReplanNeeded || inputs.OpenDecisionAlternatives != 0 {
		t.Fatalf("stubs must stay zero: %#v", inputs)
	}
}
