package eval_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/eval"
	"github.com/mrchatam/Trace/internal/store"
)

func seedGoalTask(t *testing.T, svc *domain.Service) (store.Goal, store.Task) {
	t.Helper()
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "eval results goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "eval results task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	return g, task
}

func TestListEvaluationResultsForFutureAgents(t *testing.T) {
	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	_, task := seedGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abcdef0",
		ScoresJSON: `{"correctness": 0.98}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	recorded, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.95}`,
	})
	if err != nil {
		t.Fatal(err)
	}

	rows, err := eval.ListResults(ctx, svc, task.ID)
	if err != nil {
		t.Fatalf("ListResults: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows: got %d want 1: %+v", len(rows), rows)
	}
	got := rows[0]
	if got.ID != recorded.ID {
		t.Fatalf("id: got %q want %q", got.ID, recorded.ID)
	}
	if got.TaskID != task.ID {
		t.Fatalf("task_id: got %q want %q", got.TaskID, task.ID)
	}
	if got.MechanismID != eval.MechanismStoredEvaluation {
		t.Fatalf("mechanism_id: got %q want %q", got.MechanismID, eval.MechanismStoredEvaluation)
	}
	if strings.TrimSpace(got.ComparisonJSON) == "" || got.ComparisonJSON == "{}" {
		t.Fatalf("comparison_json must be non-empty: %q", got.ComparisonJSON)
	}
	if got.Passed != nil {
		t.Fatalf("evaluation passed must be nil, got %v", *got.Passed)
	}
}

func TestEvalResultsIncludeMechanismID(t *testing.T) {
	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	g, task := seedGoalTask(t, svc)

	ev, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "proof"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestEvalResults",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "1234567",
		ScoresJSON: `{"correctness": 1.0}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.99}`,
	}); err != nil {
		t.Fatal(err)
	}

	rows, err := eval.ListResults(ctx, svc, task.ID)
	if err != nil {
		t.Fatalf("ListResults: %v", err)
	}
	if len(rows) < 3 {
		t.Fatalf("want >=3 rows, got %d: %+v", len(rows), rows)
	}
	seen := map[string]bool{
		eval.MechanismStoredTest:         false,
		eval.MechanismStoredVerification: false,
		eval.MechanismStoredEvaluation:   false,
	}
	for _, row := range rows {
		if strings.TrimSpace(row.MechanismID) == "" {
			t.Fatalf("missing mechanism_id on row %+v", row)
		}
		if _, ok := seen[row.MechanismID]; ok {
			seen[row.MechanismID] = true
		}
		switch row.Kind {
		case store.OutcomeKindTest:
			if row.Passed == nil || !*row.Passed {
				t.Fatalf("test row passed: %+v", row)
			}
		case store.OutcomeKindVerification:
			if row.Passed == nil || !*row.Passed {
				t.Fatalf("verification row passed: %+v", row)
			}
		case store.OutcomeKindEvaluation:
			if strings.TrimSpace(row.ComparisonJSON) == "" {
				t.Fatalf("evaluation comparison_json empty: %+v", row)
			}
		}
	}
	for id, ok := range seen {
		if !ok {
			t.Fatalf("missing mechanism mapping for %q in %+v", id, rows)
		}
	}
}
