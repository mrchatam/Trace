package domain_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	seedGoalID   = "11111111-1111-1111-1111-111111111111"
	seedTaskID   = "22222222-2222-2222-2222-222222222222"
	seedDecID    = "33333333-3333-3333-3333-333333333333"
	seedBaseline = "44444444-4444-4444-4444-444444444444"
	seedOutcome  = "55555555-5555-5555-5555-555555555555"
	seedChange   = "66666666-6666-6666-6666-666666666666"
	seedEffect   = "77777777-7777-7777-7777-777777777777"
	seedUncert   = "88888888-8888-8888-8888-888888888888"
	seedHypo     = "99999999-9999-9999-9999-999999999999"
	seedRecon    = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
	seedRegr     = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1"
	seedReflect  = "cccccccc-cccc-cccc-cccc-cccccccccc01"
)

func seedP20Cognition(t *testing.T, st *store.Store) {
	t.Helper()
	ctx := context.Background()
	svc := domain.New(st)

	if _, _, err := svc.ImportSeedGoal(ctx, domain.SeedEntity{ID: seedGoalID, Title: "G", Body: ""}); err != nil {
		t.Fatalf("goal: %v", err)
	}
	if _, _, err := svc.ImportSeedTask(ctx, domain.SeedTask{
		ID: seedTaskID, Title: "T", Body: "", GoalID: seedGoalID,
	}, strPtr(seedGoalID)); err != nil {
		t.Fatalf("task: %v", err)
	}
	if _, _, err := svc.ImportSeedDecision(ctx, domain.SeedEntity{ID: seedDecID, Title: "D", Body: ""}); err != nil {
		t.Fatalf("decision: %v", err)
	}

	if _, err := st.UpsertBaseline(store.Baseline{
		ID: seedBaseline, GitCommit: "abc1234", ScoresJSON: `{"latency":1.0}`, Label: "main",
	}); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		ID: seedOutcome, TaskID: seedTaskID, Kind: store.OutcomeKindTest,
		TestName: "TestFoo", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatalf("outcome: %v", err)
	}
	if _, err := st.UpsertChange(store.Change{
		ID: seedChange, TaskID: seedTaskID, GitCommit: "abc1234", Status: store.ChangeStatusRecorded,
	}); err != nil {
		t.Fatalf("change: %v", err)
	}
	if _, err := st.InsertChangePath(store.ChangePath{
		ChangeID: seedChange, Path: "internal/foo.go", Status: "M",
	}); err != nil {
		t.Fatalf("change path: %v", err)
	}
	if _, err := st.UpsertEffect(store.Effect{
		ID: seedEffect, ChangeID: seedChange, Dimension: "latency", Expected: "fast", Actual: "fast",
		Comparison: store.EffectComparisonSupported,
	}); err != nil {
		t.Fatalf("effect: %v", err)
	}
	if _, err := st.UpsertUncertainty(store.Uncertainty{
		ID: seedUncert, Title: "Will it scale?", Severity: store.UncertaintySeverityINFO,
	}); err != nil {
		t.Fatalf("uncertainty: %v", err)
	}
	if _, err := st.UpsertHypothesis(store.Hypothesis{
		ID: seedHypo, Title: "SQLite is enough",
	}); err != nil {
		t.Fatalf("hypothesis: %v", err)
	}
	if _, err := st.UpsertDecisionReconsideration(store.DecisionReconsideration{
		ID: seedRecon, DecisionID: seedDecID, Trigger: store.ReconsiderTriggerNewEvidence,
		Status: store.ReconsiderStatusOpen, Reason: "new data",
	}); err != nil {
		t.Fatalf("reconsideration: %v", err)
	}
	if _, err := st.UpsertRegression(store.Regression{
		ID: seedRegr, TaskID: seedTaskID, SourceKind: store.RegressionSourceEvaluation,
		SourceID: seedOutcome, Dimension: "latency",
	}); err != nil {
		t.Fatalf("regression: %v", err)
	}
	if _, err := st.UpsertReflection(store.Reflection{
		ID: seedReflect, TaskID: seedTaskID, Summary: "learned",
		InvalidatedAssumptionsJSON: `[]`, NewDependenciesJSON: `[]`, UsefulTestsJSON: `["TestFoo"]`,
	}); err != nil {
		t.Fatalf("reflection: %v", err)
	}
	if _, err := st.UpsertDeliberationState(store.DeliberationState{
		TaskID: seedTaskID, GoalID: seedGoalID, CurrentPhase: string(deliberation.PhaseInvestigate),
		HopCount: 2, LastPhase: string(deliberation.PhaseOrient), PlanCritiqued: true,
	}); err != nil {
		t.Fatalf("deliberation state: %v", err)
	}
}

func strPtr(s string) *string { return &s }

func TestSeedExportIncludesP20Cognition(t *testing.T) {
	svc, st := openDomain(t)
	seedP20Cognition(t, st)

	doc, err := domain.BuildSeedDocument(context.Background(), st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}

	wantKeys := []string{
		"deliberation_states", "uncertainties", "hypotheses", "decision_reconsiderations",
		"changes", "effects", "outcome_results", "baselines", "regressions", "reflections",
	}
	for _, key := range wantKeys {
		v, ok := parsed[key]
		if !ok {
			t.Fatalf("export missing key %q", key)
		}
		var arr []json.RawMessage
		if err := json.Unmarshal(v, &arr); err != nil {
			t.Fatalf("key %q not array: %v", key, err)
		}
		if len(arr) < 1 {
			t.Fatalf("key %q has no rows", key)
		}
	}

	var changes []domain.SeedChange
	if err := json.Unmarshal(parsed["changes"], &changes); err != nil {
		t.Fatal(err)
	}
	if len(changes[0].Paths) < 1 {
		t.Fatalf("changes[0] missing nested paths: %+v", changes[0])
	}

	_ = svc // seeded via store helpers
}
