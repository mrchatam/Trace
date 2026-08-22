package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func mustEvalWithRegression(t *testing.T, svc *domain.Service, taskID string) store.OutcomeResult {
	t.Helper()
	ctx := context.Background()
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     taskID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func mustContradictedEffect(t *testing.T, svc *domain.Service, taskID string, withHyp bool) store.Effect {
	t.Helper()
	ctx := context.Background()
	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "green"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	in := domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     "red",
		Comparison: store.EffectComparisonContradicted,
	}
	if withHyp {
		in.CreateHypothesis = true
		in.HypothesisTitle = "flaky fixture"
	}
	eff, _, err := svc.RecordActualEffect(ctx, c.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	return eff
}

func hypothesisExplainingEffect(t *testing.T, st *store.Store, effectID string) string {
	t.Helper()
	links, err := st.ListLinksTo(domain.EntityEffect, effectID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == domain.RelHypothesisExplainsEffect && l.FromType == domain.EntityHypothesis {
			return l.FromID
		}
	}
	t.Fatal("missing hypothesis_explains_effect")
	return ""
}

func TestRecordRegressionFromEvaluationDefaultsCorrelated(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)

	open, err := svc.HasOpenRegression(ctx, task.ID)
	if err != nil || open {
		t.Fatalf("RecordEvaluationOutcome must not auto-insert regression: open=%v err=%v", open, err)
	}

	row, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID,
		TaskID:    task.ID,
		Summary:   "correctness dropped",
		Actor:     "agent",
	})
	if err != nil {
		t.Fatalf("RecordRegressionFromEvaluation: %v", err)
	}
	if row.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("attribution=%s want correlated", row.Attribution)
	}
	if row.Attribution == store.RegressionAttributionCaused {
		t.Fatal("create must not set caused")
	}
	if row.Dimension != "overall" {
		t.Fatalf("dimension=%s want overall", row.Dimension)
	}
	again, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID,
		TaskID:    task.ID,
	})
	if err != nil || again.ID != row.ID {
		t.Fatalf("re-record must return existing: %+v err=%v", again, err)
	}
	links, err := st.ListLinksFrom(domain.EntityRegression, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel == domain.RelRegressionFromEvaluation && l.ToID == out.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing regression_from_evaluation: %+v", links)
	}
}

func TestRegressionAutoLinkedFromContradictedEffect(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	eff := mustContradictedEffect(t, svc, task.ID, false)

	row, err := svc.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
		EffectID: eff.ID,
		TaskID:   task.ID,
		Summary:  "effect contradicted",
	})
	if err != nil {
		t.Fatalf("RecordRegressionFromContradictedEffect: %v", err)
	}
	if row.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("attribution=%s want correlated", row.Attribution)
	}
	if row.Attribution == store.RegressionAttributionCaused {
		t.Fatal("auto-link must not set caused")
	}

	chg, err := st.GetChange(eff.ChangeID)
	if err != nil {
		t.Fatal(err)
	}
	links, err := st.ListLinksFrom(domain.EntityRegression, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel == domain.RelRegressionAssociatedChange && l.ToType == domain.EntityChange && l.ToID == chg.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing regression_associated_change: %+v", links)
	}

	again, err := svc.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
		EffectID: eff.ID, TaskID: task.ID,
	})
	if err != nil || again.ID != row.ID {
		t.Fatalf("re-record idempotent: %+v err=%v", again, err)
	}
}

func TestListRegressionsByChangeID(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	regEval, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.AssociateRegressionWithChange(ctx, regEval.ID, chg.ID); err != nil {
		t.Fatalf("AssociateRegressionWithChange: %v", err)
	}

	eff := mustContradictedEffect(t, svc, task.ID, false)
	regEff, err := svc.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
		EffectID: eff.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	effChg, err := st.GetChange(eff.ChangeID)
	if err != nil {
		t.Fatal(err)
	}
	if effChg.ID == chg.ID {
		t.Fatal("test setup: effect change must differ from manual change")
	}

	byChange, err := svc.ListRegressionsByChangeID(ctx, chg.ID)
	if err != nil {
		t.Fatalf("ListRegressionsByChangeID: %v", err)
	}
	if len(byChange) != 1 || byChange[0].ID != regEval.ID {
		t.Fatalf("ListRegressionsByChangeID: %+v want [%s]", byChange, regEval.ID)
	}

	alias, err := svc.RegressionsForChange(ctx, effChg.ID)
	if err != nil {
		t.Fatalf("RegressionsForChange: %v", err)
	}
	if len(alias) != 1 || alias[0].ID != regEff.ID {
		t.Fatalf("RegressionsForChange: %+v want [%s]", alias, regEff.ID)
	}

	_, err = svc.ListRegressionsByChangeID(ctx, "missing-change")
	if err == nil {
		t.Fatal("missing change must fail closed")
	}
}

func TestRegressionLinkedToChangeCausedWithEvidence(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "deadbeef",
		Paths:     []domain.ChangePathInput{{Path: "b.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.AssociateRegressionWithChange(ctx, row.ID, chg.ID); err != nil {
		t.Fatalf("AssociateRegressionWithChange: %v", err)
	}

	h, err := svc.CreateHypothesis(ctx, domain.HypothesisInput{Title: "bad refactor"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.LinkHypothesisToRegression(ctx, h.ID, row.ID); err != nil {
		t.Fatalf("LinkHypothesisToRegression: %v", err)
	}
	if _, err := svc.ConfirmHypothesis(ctx, h.ID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	ev := mustEvidence(t, svc, "git bisect")
	caused, err := svc.SetRegressionAttributionCaused(ctx, row.ID, []string{ev.ID})
	if err != nil {
		t.Fatalf("SetRegressionAttributionCaused: %v", err)
	}
	if caused.Attribution != store.RegressionAttributionCaused {
		t.Fatalf("attribution=%s want caused", caused.Attribution)
	}

	links, err := st.ListLinksFrom(domain.EntityRegression, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	sawChange, sawEvidence := false, false
	for _, l := range links {
		if l.Rel == domain.RelRegressionAssociatedChange && l.ToID == chg.ID {
			sawChange = true
		}
		if l.Rel == domain.RelRegressionSupportedBy && l.ToID == ev.ID {
			sawEvidence = true
		}
	}
	if !sawChange {
		t.Fatalf("caused row must retain change link: %+v", links)
	}
	if !sawEvidence {
		t.Fatalf("caused row must have evidence link: %+v", links)
	}
}

func TestRecordRegressionFromContradictedEffectDefaultsCorrelated(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	eff := mustContradictedEffect(t, svc, task.ID, false)

	row, err := svc.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
		EffectID: eff.ID,
		TaskID:   task.ID,
		Summary:  "effect contradicted",
	})
	if err != nil {
		t.Fatalf("RecordRegressionFromContradictedEffect: %v", err)
	}
	if row.Attribution != store.RegressionAttributionCorrelated || row.Dimension != "correctness" {
		t.Fatalf("row: %+v", row)
	}
	if row.SourceKind != store.RegressionSourceContradictedEffect {
		t.Fatalf("source_kind=%s", row.SourceKind)
	}
}

func TestCorrelationAndContradictionNeverAutoSetCaused(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)

	evalRow, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if evalRow.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("eval attribution=%s", evalRow.Attribution)
	}

	eff := mustContradictedEffect(t, svc, task.ID, true)
	hypID := hypothesisExplainingEffect(t, st, eff.ID)
	effRow, err := svc.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
		EffectID: eff.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if effRow.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("contradiction must stay correlated, got %s", effRow.Attribution)
	}
	if _, err := svc.ConfirmHypothesis(ctx, hypID, "looks right"); err != nil {
		t.Fatalf("ConfirmHypothesis: %v", err)
	}
	got, err := svc.GetRegression(ctx, effRow.ID)
	if err != nil || got.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("ConfirmHypothesis must not set caused/hypothesized: %+v err=%v", got, err)
	}

	_, err = st.UpsertRegression(store.Regression{
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    "other-outcome",
		Dimension:   "overall",
		Attribution: store.RegressionAttributionCaused,
	})
	if err == nil {
		t.Fatal("create-time caused must fail closed")
	}
	_, err = st.UpsertRegression(store.Regression{
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    "other-outcome-2",
		Dimension:   "overall",
		Attribution: store.RegressionAttributionHypothesized,
	})
	if err == nil {
		t.Fatal("create-time hypothesized must fail closed")
	}
}

func TestLinkHypothesisUpgradesToHypothesizedNotCaused(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := svc.CreateHypothesis(ctx, domain.HypothesisInput{Title: "cache miss"})
	if err != nil {
		t.Fatal(err)
	}
	up, err := svc.LinkHypothesisToRegression(ctx, h.ID, row.ID)
	if err != nil {
		t.Fatalf("LinkHypothesisToRegression: %v", err)
	}
	if up.Attribution != store.RegressionAttributionHypothesized {
		t.Fatalf("want hypothesized, got %s", up.Attribution)
	}
	if up.Attribution == store.RegressionAttributionCaused {
		t.Fatal("link must not set caused")
	}
	if _, err := svc.ConfirmHypothesis(ctx, h.ID, "confirmed"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetRegression(ctx, row.ID)
	if err != nil || got.Attribution != store.RegressionAttributionHypothesized {
		t.Fatalf("ConfirmHypothesis must not set caused: %+v err=%v", got, err)
	}
}

func TestSetAttributionCausedFailClosedWithoutEvidence(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, _ := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	h, _ := svc.CreateHypothesis(ctx, domain.HypothesisInput{Title: "h"})
	if _, err := svc.LinkHypothesisToRegression(ctx, h.ID, row.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ConfirmHypothesis(ctx, h.ID, "ok"); err != nil {
		t.Fatal(err)
	}

	_, err := svc.SetRegressionAttributionCaused(ctx, row.ID, nil)
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("empty evidence: want ErrValidation, got %v", err)
	}
	_, err = svc.SetRegressionAttributionCaused(ctx, row.ID, []string{"missing-ev"})
	if err == nil {
		t.Fatal("missing evidence row must fail closed")
	}
	got, _ := svc.GetRegression(ctx, row.ID)
	if got.Attribution != store.RegressionAttributionHypothesized {
		t.Fatalf("still hypothesized: %s", got.Attribution)
	}
}

func TestSetAttributionCausedFailClosedFromCorrelated(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, _ := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	ev := mustEvidence(t, svc, "trace")
	_, err := svc.SetRegressionAttributionCaused(ctx, row.ID, []string{ev.ID})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("from correlated: want ErrValidation, got %v", err)
	}
	got, _ := svc.GetRegression(ctx, row.ID)
	if got.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("must stay correlated: %s", got.Attribution)
	}
}

func TestSetAttributionCausedRequiresConfirmedHypothesisAndEvidence(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, _ := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	h, _ := svc.CreateHypothesis(ctx, domain.HypothesisInput{Title: "root cause"})
	if _, err := svc.LinkHypothesisToRegression(ctx, h.ID, row.ID); err != nil {
		t.Fatal(err)
	}
	ev := mustEvidence(t, svc, "profile")

	_, err := svc.SetRegressionAttributionCaused(ctx, row.ID, []string{ev.ID})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("OPEN hypothesis insufficient: got %v", err)
	}

	if _, err := svc.ConfirmHypothesis(ctx, h.ID, "yes"); err != nil {
		t.Fatal(err)
	}
	caused, err := svc.SetRegressionAttributionCaused(ctx, row.ID, []string{ev.ID})
	if err != nil {
		t.Fatalf("caused: %v", err)
	}
	if caused.Attribution != store.RegressionAttributionCaused {
		t.Fatalf("want caused, got %s", caused.Attribution)
	}
	links, err := st.ListLinksFrom(domain.EntityRegression, caused.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel == domain.RelRegressionSupportedBy && l.ToID == ev.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing regression_supported_by: %+v", links)
	}
}

func TestHasOpenRegressionFeedsApplyDeliberationTransition(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	out := mustEvalWithRegression(t, svc, task.ID)
	if _, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	}); err != nil {
		t.Fatal(err)
	}

	open, err := svc.HasOpenRegression(ctx, task.ID)
	if err != nil || !open {
		t.Fatalf("HasOpenRegression=%v err=%v want true", open, err)
	}
	if _, err := st.GetDeliberationState(task.ID); err == nil {
		t.Fatal("RecordRegression must not auto-hop / write deliberation_state")
	}

	inputs := deliberation.PolicyInputs{
		OpenRegression: open,
		PlanExists:     true,
		PlanCritiqued:  true,
	}
	next, ev, err := svc.ApplyDeliberationTransition(ctx, task.ID, g.ID, inputs)
	if err != nil {
		t.Fatalf("ApplyDeliberationTransition: %v", err)
	}
	if next.CurrentPhase != deliberation.PhaseInvestigate {
		t.Fatalf("phase %s want INVESTIGATE", next.CurrentPhase)
	}
	if next.CurrentPhase == deliberation.PhaseExecute {
		t.Fatal("EXECUTE forbidden while open regression")
	}
	var payload deliberation.TransitionPayload
	if err := json.Unmarshal([]byte(ev.PayloadJSON), &payload); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if payload.ReasonCode != deliberation.ReasonOpenRegression {
		t.Fatalf("reason %s want open_regression", payload.ReasonCode)
	}
	if !payload.PolicyInputs.OpenRegression {
		t.Fatal("policy open_regression must be true")
	}
}

func TestResolveRegressionClearsHasOpenRegression(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	row, _ := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	open, _ := svc.HasOpenRegression(ctx, task.ID)
	if !open {
		t.Fatal("expected open")
	}
	if _, err := svc.ResolveRegression(ctx, row.ID, ""); err == nil {
		t.Fatal("reason required")
	}
	resolved, err := svc.ResolveRegression(ctx, row.ID, "fixed")
	if err != nil {
		t.Fatalf("ResolveRegression: %v", err)
	}
	if resolved.Status != store.RegressionStatusResolved {
		t.Fatalf("status=%s", resolved.Status)
	}
	open, err = svc.HasOpenRegression(ctx, task.ID)
	if err != nil || open {
		t.Fatalf("after resolve open=%v err=%v", open, err)
	}
	n, err := svc.CountOpenRegressionsByTaskID(ctx, task.ID)
	if err != nil || n != 0 {
		t.Fatalf("count=%d err=%v", n, err)
	}
	listed, err := svc.ListOpenRegressions(ctx, task.ID)
	if err != nil || len(listed) != 0 {
		t.Fatalf("list=%+v err=%v", listed, err)
	}
}

func TestReflectionPersistsStructuredFieldsQueryable(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "CGO free"})
	if err != nil {
		t.Fatal(err)
	}
	row, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID:                   task.ID,
		Summary:                  "learned path coupling",
		InvalidatedAssumptionIDs: []string{a.ID},
		NewDependencies: []domain.DependencyRef{
			{Kind: "path", Ref: "./internal/store/regressions.go"},
			{Kind: "symbol", Ref: "Store.HasOpenRegression"},
		},
		UsefulTests:      []string{"TestHasOpenRegressionFeedsApplyDeliberationTransition"},
		BroadenTestsNote: "add soak later",
		Actor:            "agent",
	})
	if err != nil {
		t.Fatalf("CreateReflection: %v", err)
	}

	got, err := svc.GetReflection(ctx, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	var assump []string
	if err := json.Unmarshal([]byte(got.InvalidatedAssumptionsJSON), &assump); err != nil || len(assump) != 1 || assump[0] != a.ID {
		t.Fatalf("invalidated_assumptions_json: %s err=%v", got.InvalidatedAssumptionsJSON, err)
	}
	var deps []domain.DependencyRef
	if err := json.Unmarshal([]byte(got.NewDependenciesJSON), &deps); err != nil || len(deps) != 2 {
		t.Fatalf("new_dependencies_json: %s err=%v", got.NewDependenciesJSON, err)
	}
	if deps[0].Kind != "path" || deps[0].Ref != "internal/store/regressions.go" {
		t.Fatalf("path not normalized: %+v", deps[0])
	}
	var tests []string
	if err := json.Unmarshal([]byte(got.UsefulTestsJSON), &tests); err != nil || len(tests) != 1 {
		t.Fatalf("useful_tests_json: %s err=%v", got.UsefulTestsJSON, err)
	}
	listed, err := svc.ListReflectionsByTaskID(ctx, task.ID)
	if err != nil || len(listed) != 1 || listed[0].ID != row.ID {
		t.Fatalf("ListReflectionsByTaskID: %+v err=%v", listed, err)
	}
	links, err := st.ListLinksFrom(domain.EntityReflection, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel == domain.RelReflectionInvalidatesAssumption && l.ToID == a.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing reflection_invalidates_assumption: %+v", links)
	}
}

func TestReflectionEssayOnlyFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	_, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID:  task.ID,
		Summary: "we learned a lot in prose",
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("summary-only: want ErrValidation, got %v", err)
	}
	_, err = svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID:           task.ID,
		BroadenTestsNote: "try more tests",
	})
	if !errors.As(err, &ve) {
		t.Fatalf("note-only: want ErrValidation, got %v", err)
	}
}

func TestObservedRelationshipLinkWithConfidenceNoEvidence(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	link, err := svc.RecordObservedRelationship(ctx, domain.RelInput{
		FromType:   domain.EntityChange,
		FromID:     c.ID,
		ToType:     domain.EntityTask,
		ToID:       task.ID,
		Confidence: 0.7,
	})
	if err != nil {
		t.Fatalf("RecordObservedRelationship: %v", err)
	}
	if link.Rel != domain.RelObservedRelationship || link.Confidence != 0.7 {
		t.Fatalf("link: %+v", link)
	}
	from, err := st.ListLinksFrom(domain.EntityChange, c.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range from {
		if l.Rel == domain.RelRelationshipSupportedBy {
			t.Fatalf("observed must not require evidence links: %+v", from)
		}
	}
	_, err = svc.RecordObservedRelationship(ctx, domain.RelInput{
		FromType:   domain.EntityChange,
		FromID:     c.ID,
		ToType:     domain.EntityTask,
		ToID:       task.ID,
		Confidence: math.NaN(),
	})
	if err == nil {
		t.Fatal("NaN confidence must fail closed")
	}
}

func TestCausalRelationshipFailClosedWithoutEvidence(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	reg, _ := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.RecordCausalRelationship(ctx, domain.RelInput{
		FromType:   domain.EntityChange,
		FromID:     c.ID,
		ToType:     domain.EntityRegression,
		ToID:       reg.ID,
		Confidence: 0.9,
	}, nil)
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("no evidence: want ErrValidation, got %v", err)
	}

	ev := mustEvidence(t, svc, "git blame")
	link, err := svc.RecordCausalRelationship(ctx, domain.RelInput{
		FromType:   domain.EntityChange,
		FromID:     c.ID,
		ToType:     domain.EntityRegression,
		ToID:       reg.ID,
		Confidence: 0.9,
	}, []string{ev.ID})
	if err != nil {
		t.Fatalf("with evidence: %v", err)
	}
	if link.Rel != domain.RelCausedBy {
		t.Fatalf("rel=%s", link.Rel)
	}
	got, err := svc.GetRegression(ctx, reg.ID)
	if err != nil || got.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("caused_by must not set attribution=caused: %+v err=%v", got, err)
	}
	from, err := st.ListLinksFrom(domain.EntityChange, c.ID)
	if err != nil {
		t.Fatal(err)
	}
	sawSupport := false
	for _, l := range from {
		if l.Rel == domain.RelRelationshipSupportedBy && l.ToID == ev.ID {
			sawSupport = true
		}
	}
	if !sawSupport {
		t.Fatalf("missing relationship_supported_by: %+v", from)
	}
}

func TestUnknownAttributionFailClosed(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	_, err := st.UpsertRegression(store.Regression{
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    "o1",
		Dimension:   "overall",
		Attribution: "inferred",
	})
	if err == nil {
		t.Fatal("unknown attribution must fail closed")
	}

	out := mustEvalWithRegression(t, svc, task.ID)
	row, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if row.Attribution != store.RegressionAttributionCorrelated &&
		row.Attribution != store.RegressionAttributionHypothesized &&
		row.Attribution != store.RegressionAttributionCaused {
		t.Fatalf("unknown persisted attribution %q", row.Attribution)
	}
	if row.Attribution != store.RegressionAttributionCorrelated {
		t.Fatalf("create must be correlated, got %s", row.Attribution)
	}
}
