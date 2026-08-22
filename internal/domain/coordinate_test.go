package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestRegressionDetectedSameSecondTimestamps(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	const sameSecond = "2026-08-18T12:00:00Z"
	// UUID lex order would put fail before pass; rowid preserves insertion order.
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		ID:         "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz",
		TaskID:     task.ID,
		Kind:       store.OutcomeKindTest,
		TestName:   "TestFoo",
		TestStatus: store.TestStatusPass,
		CreatedAt:  sameSecond,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		ID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		TaskID:     task.ID,
		Kind:       store.OutcomeKindTest,
		TestName:   "TestFoo",
		TestStatus: store.TestStatusFail,
		CreatedAt:  sameSecond,
	}); err != nil {
		t.Fatal(err)
	}

	sig, err := svc.DetectTestRegression(ctx, task.ID, "TestFoo")
	if err != nil {
		t.Fatal(err)
	}
	if !sig.Detected || sig.TestName != "TestFoo" {
		t.Fatalf("regression: %+v", sig)
	}
	if !sig.PriorPass || !sig.CurrentFail {
		t.Fatalf("prior/current flags: %+v", sig)
	}
}

func TestRegressionDetectedVsPriorPassingTest(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestFoo", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestFoo", TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}
	rows, err := st.ListOutcomeResultsByTaskKind(task.ID, store.OutcomeKindTest)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) < 2 {
		t.Fatalf("want >=2 test rows, got %d", len(rows))
	}

	sig, err := svc.DetectTestRegression(ctx, task.ID, "TestFoo")
	if err != nil {
		t.Fatal(err)
	}
	if !sig.Detected || sig.TestName != "TestFoo" {
		t.Fatalf("regression: %+v", sig)
	}
	if !sig.PriorPass || !sig.CurrentFail {
		t.Fatalf("prior/current flags: %+v", sig)
	}
}

func TestCoordinateVerificationOrder(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	mustRecordedChange(t, svc, task.ID)

	ev := mustEvidence(t, svc, "proof")
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abcdef0",
		ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}

	var order []string
	_, err = svc.CoordinateVerification(ctx, task.ID, domain.CoordinateOptions{
		EvidenceIDs: []string{ev.ID},
		ScoresJSON:  `{"correctness":0.95}`,
		Hooks: &domain.CoordinateTestHooks{
			RunTests: func(ctx context.Context, taskID string, paths []string) ([]store.OutcomeResult, error) {
				order = append(order, "test")
				return nil, nil
			},
			EnsureVerification: func(ctx context.Context) error {
				order = append(order, "verify")
				_, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
					TaskID:             task.ID,
					GoalID:             g.ID,
					VerificationStatus: store.VerificationStatusVerified,
					EvidenceIDs:        []string{ev.ID},
				})
				return err
			},
			RunEvaluation: func(ctx context.Context) (store.OutcomeResult, error) {
				order = append(order, "evaluate")
				return svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
					TaskID: task.ID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.95}`,
				})
			},
		},
	})
	if err != nil {
		t.Fatalf("CoordinateVerification: %v", err)
	}
	want := []string{"test", "verify", "evaluate"}
	if len(order) != len(want) {
		t.Fatalf("order=%v want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order=%v want %v", order, want)
		}
	}
}

func TestCoordinateVerificationStopsOnTestFailUnlessForceEval(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	mustRecordedChange(t, svc, task.ID)

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestFoo", TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}

	res, err := svc.CoordinateVerification(ctx, task.ID, domain.CoordinateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !res.StoppedEarly || res.StopReason != "test_failed" {
		t.Fatalf("want stop on fail: %+v", res)
	}
	if res.EvaluationRecorded {
		t.Fatal("eval must not run after test fail")
	}

	res2, err := svc.CoordinateVerification(ctx, task.ID, domain.CoordinateOptions{
		ForceEval:  true,
		ScoresJSON: `{"correctness":1}`,
		Hooks: &domain.CoordinateTestHooks{
			RunEvaluation: func(ctx context.Context) (store.OutcomeResult, error) {
				return store.OutcomeResult{}, nil
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res2.StoppedEarly {
		t.Fatalf("force-eval must not stop early: %+v", res2)
	}
}
