package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// seedCoordinateCycle prepares the reachable part of the chain: task with a
// goal, a RECORDED change at commit abc1234, a passing test outcome, and a
// verified verification with evidence. Returns the task and evidence ids.
func seedCoordinateCycle(t *testing.T, svc *domain.Service, commit string) (string, string) {
	t.Helper()
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: commit,
		Paths:     []domain.ChangePathInput{{Path: "src/agent-runtime.ts"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "agent-runtime.test.ts", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	ev := mustEvidence(t, svc, "proof")
	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}
	return task.ID, c.ID
}

// TestCoordinateVerificationCreatesBaselineAndRecordsEvaluation is the
// regression for the unreachable evaluation stage: a task with a change, test
// outcomes, and a verification — but NO baseline — must record a
// kind=evaluation outcome via CoordinateVerification, creating the
// commit-scoped baseline on demand. Fails on the pre-fix build (evaluation was
// never recorded; baselines stayed empty).
func TestCoordinateVerificationCreatesBaselineAndRecordsEvaluation(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	taskID, _ := seedCoordinateCycle(t, svc, "abc1234")

	res, err := svc.CoordinateVerification(ctx, taskID, domain.CoordinateOptions{
		ScoresJSON: `{"correctness":0.95}`,
		Actor:      "cli",
		SourceType: "CLI",
	})
	if err != nil {
		t.Fatalf("CoordinateVerification: %v", err)
	}
	if !res.EvaluationRecorded {
		t.Fatalf("evaluation must be recorded, got %+v (stop_reason=%q)", res, res.StopReason)
	}
	if res.StopReason == "no_active_baseline" {
		t.Fatal("no_active_baseline must not occur when the change has a commit")
	}

	baselines, err := st.ListAllBaselines()
	if err != nil {
		t.Fatal(err)
	}
	if len(baselines) != 1 {
		t.Fatalf("baselines=%d want 1 (created on demand)", len(baselines))
	}
	b := baselines[0]
	if b.GitCommit != "abc1234" || b.Status != store.BaselineStatusActive {
		t.Fatalf("baseline: %+v", b)
	}

	rows, err := st.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("evaluation rows=%d want 1", len(rows))
	}
	o := rows[0]
	if o.BaselineID != b.ID {
		t.Fatalf("evaluation baseline_id=%q want %q", o.BaselineID, b.ID)
	}
	// comparison_json must be the computed structure, not a boolean pass, and
	// the on-demand baseline snapshots the first evaluation's own scores (an
	// initialization pass — zero deltas against its own snapshot).
	if o.ScoresJSON != `{"correctness":0.95}` {
		t.Fatalf("scores_json=%q", o.ScoresJSON)
	}
	if b.ScoresJSON != `{"correctness":0.95}` {
		t.Fatalf("baseline snapshot scores=%q want the first evaluation's scores", b.ScoresJSON)
	}

	// The gate must now see a computed evaluation.
	ok, err := svc.HasComputedEvaluation(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("HasComputedEvaluation must be true after the coordinated evaluation")
	}

	// Idempotence: a second coordinate run must not create another baseline.
	if _, err := svc.CoordinateVerification(ctx, taskID, domain.CoordinateOptions{
		ScoresJSON: `{"correctness":0.94}`,
		Actor:      "cli",
		SourceType: "CLI",
	}); err != nil {
		t.Fatalf("second CoordinateVerification: %v", err)
	}
	baselines, err = st.ListAllBaselines()
	if err != nil {
		t.Fatal(err)
	}
	if len(baselines) != 1 {
		t.Fatalf("baselines after re-run=%d want 1 (reuse, never duplicate)", len(baselines))
	}
}

// TestCoordinateVerificationReactivatesSupersededBaseline covers the supersede
// path: a baseline for the commit that was superseded (e.g. deliberately via
// SupersedeBaseline) is reactivated via PromoteBaseline — status active again,
// supersedes_id linking the prior active row — instead of inserting a
// duplicate row for the same commit.
func TestCoordinateVerificationReactivatesSupersededBaseline(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	// Commit A gets a baseline with scores (simulating a prior evaluation era),
	// then it is deliberately superseded (operator action via the public API).
	first, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "aaaaaaa",
		ScoresJSON: `{"correctness":0.90}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.SupersedeBaseline(ctx, first.ID); err != nil {
		t.Fatal(err)
	}
	sup, err := st.GetBaseline(first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if sup.Status != store.BaselineStatusSuperseded {
		t.Fatalf("precondition: first baseline status=%q want superseded", sup.Status)
	}

	taskID, _ := seedCoordinateCycle(t, svc, "aaaaaaa")
	res, err := svc.CoordinateVerification(ctx, taskID, domain.CoordinateOptions{
		ScoresJSON: `{"correctness":0.93}`,
		Actor:      "cli",
		SourceType: "CLI",
	})
	if err != nil {
		t.Fatalf("CoordinateVerification: %v", err)
	}
	if !res.EvaluationRecorded {
		t.Fatalf("evaluation must be recorded against the reactivated baseline: %+v", res)
	}

	baselines, err := st.ListAllBaselines()
	if err != nil {
		t.Fatal(err)
	}
	if len(baselines) != 1 {
		t.Fatalf("baselines=%d want 1 (reactivated, not duplicated)", len(baselines))
	}
	reactivated, err := st.GetBaseline(first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reactivated.Status != store.BaselineStatusActive {
		t.Fatalf("superseded baseline must be reactivated, status=%q", reactivated.Status)
	}
	// compactJSONObject normalizes 0.90 → 0.9; the point is: unchanged.
	if reactivated.ScoresJSON != `{"correctness":0.9}` {
		t.Fatalf("baseline scores must never be mutated, scores=%q", reactivated.ScoresJSON)
	}

	rows, err := st.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].BaselineID != first.ID {
		t.Fatalf("evaluation must anchor to the reactivated baseline: %+v", rows)
	}
}

// TestCoordinateVerificationWithoutCommitFailsClosed pins the fail-closed
// contract: with no change commit to anchor to, evaluation is not recorded and
// stop_reason reports the missing baseline (never a fabricated one).
func TestCoordinateVerificationWithoutCommitFailsClosed(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	ev := mustEvidence(t, svc, "proof")
	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}

	res, err := svc.CoordinateVerification(ctx, task.ID, domain.CoordinateOptions{
		ScoresJSON: `{"correctness":0.95}`,
		Actor:      "cli",
		SourceType: "CLI",
	})
	if err != nil {
		t.Fatalf("CoordinateVerification: %v", err)
	}
	if res.EvaluationRecorded {
		t.Fatal("evaluation must not be recorded without a change commit")
	}
	if res.StopReason != "no_active_baseline" {
		t.Fatalf("stop_reason=%q want no_active_baseline", res.StopReason)
	}
	baselines, err := st.ListAllBaselines()
	if err != nil {
		t.Fatal(err)
	}
	if len(baselines) != 0 {
		t.Fatalf("no baseline may be fabricated without a commit: %d", len(baselines))
	}
}
