package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func mustGoalTask(t *testing.T, svc *domain.Service) (store.Goal, store.Task) {
	t.Helper()
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "ship verify"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "implement", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	return g, task
}

func mustEvidence(t *testing.T, svc *domain.Service, title string) store.Evidence {
	t.Helper()
	e, err := svc.CreateEvidence(context.Background(), domain.EvidenceInput{Title: title})
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func mustRecordedChange(t *testing.T, svc *domain.Service, taskID string) store.Change {
	t.Helper()
	c, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "internal/domain/outcomes.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestRecordTestOutcomeRequiresNameAndStatus(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	_, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestStatus: store.TestStatusPass,
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("missing name: want ErrValidation, got %v", err)
	}

	_, err = svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:   task.ID,
		TestName: "TestFoo",
	})
	if !errors.As(err, &ve) {
		t.Fatalf("missing status: want ErrValidation, got %v", err)
	}

	_, err = svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestFoo",
		TestStatus: "flaky",
	})
	if !errors.As(err, &ve) {
		t.Fatalf("unknown status: want ErrValidation, got %v", err)
	}

	got, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestFoo",
		TestStatus: store.TestStatusPass,
		Summary:    "ok",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != store.OutcomeKindTest || got.TestName != "TestFoo" || got.TestStatus != store.TestStatusPass {
		t.Fatalf("recorded: %+v", got)
	}
	if got.GoalID != "" || got.VerificationStatus != "" || got.BaselineID != "" {
		t.Fatalf("test kind pollution: %+v", got)
	}
}

func TestTestPassAloneCannotSatisfyVerificationGate(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	ok, _, err := svc.CheckTestGate(ctx, task.ID, "TestShip")
	if err != nil || !ok {
		t.Fatalf("test gate: ok=%v err=%v", ok, err)
	}
	vOK, reason, err := svc.CheckVerificationGate(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if vOK {
		t.Fatalf("test pass must not satisfy verification gate; reason=%q", reason)
	}
}

func TestVerificationRequiresGoalAndEvidenceIDs(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	ev := mustEvidence(t, svc, "trace.log excerpt")

	_, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("missing goal: want ErrValidation, got %v", err)
	}

	_, err = svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
	})
	if !errors.As(err, &ve) {
		t.Fatalf("missing evidence: want ErrValidation, got %v", err)
	}

	got, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
		Summary:            "goal met",
	})
	if err != nil {
		t.Fatal(err)
	}
	ok, _, err := svc.CheckVerificationGate(ctx, task.ID)
	if err != nil || !ok {
		t.Fatalf("verification gate: ok=%v err=%v outcome=%+v", ok, err, got)
	}
}

func TestVerificationMissingEvidenceFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)

	_, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{"missing-evidence-id"},
	})
	if err == nil {
		t.Fatal("missing evidence row must fail closed")
	}
	ok, _, err := svc.CheckVerificationGate(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("verification gate must stay false without evidence")
	}
}

func TestEvaluationComparesScoresToBaselineNotBoolean(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "ABCDEF012",
		ScoresJSON: `{"correctness": 0.98, "performance_p95_ms": 310}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.95, "performance_p95_ms": 280}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	var probe map[string]any
	if err := json.Unmarshal([]byte(out.ComparisonJSON), &probe); err != nil {
		t.Fatalf("comparison_json: %v", err)
	}
	if _, ok := probe["pass"]; ok {
		t.Fatalf("comparison_json must not be boolean PASS: %s", out.ComparisonJSON)
	}
	if probe["baseline_id"] != b.ID {
		t.Fatalf("baseline_id: %+v", probe)
	}
	dims, _ := probe["dimensions"].(map[string]any)
	if len(dims) < 2 {
		t.Fatalf("expected per-dimension deltas: %s", out.ComparisonJSON)
	}
	corr, _ := dims["correctness"].(map[string]any)
	if corr == nil {
		t.Fatalf("missing correctness dim: %s", out.ComparisonJSON)
	}
	if _, ok := corr["delta"]; !ok {
		t.Fatalf("missing delta: %+v", corr)
	}
	ok, _, err := svc.CheckEvaluationGate(ctx, task.ID, b.ID)
	if err != nil || !ok {
		t.Fatalf("evaluation gate: ok=%v err=%v", ok, err)
	}
}

func TestEvaluationRegressionFlagInComparisonJSON(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98, "performance_p95_ms": 310}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.90, "performance_p95_ms": 250}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	var cmp domain.ScoreComparison
	if err := json.Unmarshal([]byte(out.ComparisonJSON), &cmp); err != nil {
		t.Fatal(err)
	}
	corr, ok := cmp.Dimensions["correctness"]
	if !ok || !corr.Regression {
		t.Fatalf("correctness drop must regress: %+v", cmp.Dimensions)
	}
	lat, ok := cmp.Dimensions["performance_p95_ms"]
	if !ok || lat.Regression {
		t.Fatalf("lower latency must not regress: %+v", cmp.Dimensions)
	}
	if !cmp.OverallRegression {
		t.Fatal("overall_regression must be true when any dimension regresses")
	}
}

func TestBaselineStoresCommitOIDAndScoresJSONOnly(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	_, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "not-hex",
		ScoresJSON: `{"a":1}`,
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("bad OID: want ErrValidation, got %v", err)
	}
	_, err = svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abcdef0",
		ScoresJSON: `{}`,
	})
	if !errors.As(err, &ve) {
		t.Fatalf("empty scores: want ErrValidation, got %v", err)
	}
	_, err = svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abcdef0",
		ScoresJSON: `[{"a":1}]`,
	})
	if !errors.As(err, &ve) {
		t.Fatalf("array root: want ErrValidation, got %v", err)
	}

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "ABCDEF0",
		ScoresJSON: `{"correctness": 1}`,
		Label:      "B1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.GitCommit != "abcdef0" {
		t.Fatalf("OID lowercase: %q", b.GitCommit)
	}
	if !strings.Contains(b.ScoresJSON, "correctness") {
		t.Fatalf("scores_json: %s", b.ScoresJSON)
	}
	got, err := svc.GetBaseline(ctx, b.ID)
	if err != nil || got.Label != "B1" {
		t.Fatalf("GetBaseline: %+v err=%v", got, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityBaseline, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	var saw bool
	for _, e := range evs {
		if e.Type == domain.EventBaselineCreated {
			saw = true
		}
	}
	if !saw {
		t.Fatalf("missing baseline.created: %+v", evs)
	}
	bad, where, err := st.HasBlobLikeColumns()
	if err != nil || bad {
		t.Fatalf("HasBlobLikeColumns: bad=%v where=%s err=%v", bad, where, err)
	}
}

func TestVerificationDebtWhenImplementationWithoutVerification(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	signal, err := svc.HasImplementationSignal(ctx, task.ID)
	if err != nil || signal {
		t.Fatalf("no changes → no signal: %v %v", signal, err)
	}
	debt, err := svc.HasVerificationDebt(ctx, task.ID)
	if err != nil || debt {
		t.Fatalf("no implementation → no debt: %v %v", debt, err)
	}

	mustRecordedChange(t, svc, task.ID)
	signal, err = svc.HasImplementationSignal(ctx, task.ID)
	if err != nil || !signal {
		t.Fatalf("RECORDED change → signal: %v %v", signal, err)
	}
	debt, err = svc.HasVerificationDebt(ctx, task.ID)
	if err != nil || !debt {
		t.Fatalf("implementation without verification → debt: %v %v", debt, err)
	}
	items, err := svc.ListVerificationDebtSummary(ctx, task.ID)
	if err != nil || len(items) != 1 || items[0].Missing == "" {
		t.Fatalf("debt summary: %+v err=%v", items, err)
	}
}

func TestVerificationDebtClearsWhenVerifiedWithEvidence(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	mustRecordedChange(t, svc, task.ID)
	ev := mustEvidence(t, svc, "integration evidence")

	debt, err := svc.HasVerificationDebt(ctx, task.ID)
	if err != nil || !debt {
		t.Fatalf("want debt before verify: %v %v", debt, err)
	}
	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}
	debt, err = svc.HasVerificationDebt(ctx, task.ID)
	if err != nil || debt {
		t.Fatalf("verified with evidence clears debt: %v %v", debt, err)
	}
}

func TestPromotionGateRequiresStoredTestNotAgentClaim(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	ok, reason, err := svc.CheckTestGate(ctx, task.ID, "TestShip")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("agent claim without stored row must not pass gate; reason=%q", reason)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}
	ok, _, err = svc.CheckTestGate(ctx, task.ID, "TestShip")
	if err != nil || ok {
		t.Fatalf("fail status must not satisfy test gate: ok=%v err=%v", ok, err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	ok, _, err = svc.CheckTestGate(ctx, task.ID, "TestShip")
	if err != nil || !ok {
		t.Fatalf("stored pass must satisfy test gate: ok=%v err=%v", ok, err)
	}
}

func TestEvaluationMissingBaselineFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	_, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		ScoresJSON: `{"correctness": 1}`,
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("missing baseline_id: want ErrValidation, got %v", err)
	}
	_, err = svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: "missing-baseline",
		ScoresJSON: `{"correctness": 1}`,
	})
	if err == nil {
		t.Fatal("missing baseline row must fail closed")
	}
}

func TestPartialVerificationCountsAsDebt(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, task := mustGoalTask(t, svc)
	mustRecordedChange(t, svc, task.ID)
	ev := mustEvidence(t, svc, "partial evidence")

	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             task.ID,
		GoalID:             g.ID,
		VerificationStatus: store.VerificationStatusPartial,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}
	ok, _, err := svc.CheckVerificationGate(ctx, task.ID)
	if err != nil || ok {
		t.Fatalf("partial must not satisfy verification gate: ok=%v err=%v", ok, err)
	}
	debt, err := svc.HasVerificationDebt(ctx, task.ID)
	if err != nil || !debt {
		t.Fatalf("partial counts as debt: %v %v", debt, err)
	}
}

func TestPromoteBaselineSupersedesPrior(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	b100, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatal(err)
	}
	b101, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.99}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatal(err)
	}
	promoted, err := svc.PromoteBaseline(ctx, b101.ID)
	if err != nil {
		t.Fatal(err)
	}
	if promoted.Status != store.BaselineStatusActive || promoted.SupersedesID != b100.ID {
		t.Fatalf("promoted B101: %+v want active supersedes=%s", promoted, b100.ID)
	}
	got100, err := st.GetBaseline(b100.ID)
	if err != nil || got100.Status != store.BaselineStatusSuperseded {
		t.Fatalf("B100 superseded: %+v err=%v", got100, err)
	}
}

func TestPromoteBaselineIdempotent(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	b100, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatal(err)
	}
	b101, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.99}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.PromoteBaseline(ctx, b101.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.PromoteBaseline(ctx, b101.ID); err != nil {
		t.Fatal(err)
	}
	got100, err := st.GetBaseline(b100.ID)
	if err != nil || got100.Status != store.BaselineStatusSuperseded {
		t.Fatalf("B100 still superseded: %+v err=%v", got100, err)
	}
	got101, err := st.GetBaseline(b101.ID)
	if err != nil || got101.Status != store.BaselineStatusActive || got101.SupersedesID != b100.ID {
		t.Fatalf("B101 active chain: %+v err=%v", got101, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityBaseline, b101.ID)
	if err != nil {
		t.Fatal(err)
	}
	var promoted int
	for _, e := range evs {
		if e.Type == domain.EventBaselinePromoted {
			promoted++
		}
	}
	if promoted != 1 {
		t.Fatalf("expected one baseline.promoted event, got %d", promoted)
	}
}

func TestEvalRegressionBlocksPromotionGate(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98, "performance_p95_ms": 310}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.90, "performance_p95_ms": 250}`,
	}); err != nil {
		t.Fatal(err)
	}
	allowed, reason, err := svc.CheckPromotionGate(ctx, task.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if allowed || reason != "eval_regression" {
		t.Fatalf("gate blocked on regression: allowed=%v reason=%q", allowed, reason)
	}
}

func TestEvalRegressionGateClearsAfterResolve(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98, "performance_p95_ms": 310}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	regressed, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.90, "performance_p95_ms": 250}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	clean, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.99, "performance_p95_ms": 300}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		ID: regressed.ID, TaskID: task.ID, Kind: store.OutcomeKindEvaluation,
		BaselineID: b.ID, ScoresJSON: regressed.ScoresJSON, ComparisonJSON: regressed.ComparisonJSON,
		CreatedAt: "2026-01-01T00:00:00Z", UpdatedAt: "2026-01-01T00:00:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		ID: clean.ID, TaskID: task.ID, Kind: store.OutcomeKindEvaluation,
		BaselineID: b.ID, ScoresJSON: clean.ScoresJSON, ComparisonJSON: clean.ComparisonJSON,
		CreatedAt: "2026-01-02T00:00:00Z", UpdatedAt: "2026-01-02T00:00:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	allowed, reason, err := svc.CheckPromotionGate(ctx, task.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if !allowed || reason != "" {
		t.Fatalf("newer clean eval opens gate: allowed=%v reason=%q", allowed, reason)
	}
}

func TestPromotionGateIndependentOfTestPassAlone(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	ok, _, err := svc.CheckTestGate(ctx, task.ID, "TestShip")
	if err != nil || !ok {
		t.Fatalf("test gate passes: ok=%v err=%v", ok, err)
	}
	allowed, reason, err := svc.CheckPromotionGate(ctx, task.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if allowed || reason != "no_stored_evaluation" {
		t.Fatalf("test pass alone must not open promotion gate: allowed=%v reason=%q", allowed, reason)
	}
}

func TestBaselinePromotionRequiresStoredEvaluation(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	allowed, reason, err := svc.CheckPromotionGate(ctx, task.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if allowed || reason != "no_stored_evaluation" {
		t.Fatalf("no eval blocks gate: allowed=%v reason=%q", allowed, reason)
	}
}

func TestHasTestOutcomeSinceLatestChange(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	got, err := svc.HasTestOutcomeSinceLatestChange(ctx, task.ID)
	if err != nil || got {
		t.Fatalf("no change → false: got=%v err=%v", got, err)
	}

	mustRecordedChange(t, svc, task.ID)
	got, err = svc.HasTestOutcomeSinceLatestChange(ctx, task.ID)
	if err != nil || got {
		t.Fatalf("change without test → false: got=%v err=%v", got, err)
	}

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	got, err = svc.HasTestOutcomeSinceLatestChange(ctx, task.ID)
	if err != nil || !got {
		t.Fatalf("test since change → true: got=%v err=%v", got, err)
	}

	_, task2 := mustGoalTask(t, svc)
	if _, err := st.UpsertOutcomeResult(store.OutcomeResult{
		TaskID:     task2.ID,
		Kind:       store.OutcomeKindTest,
		TestName:   "TestStale",
		TestStatus: store.TestStatusPass,
		CreatedAt:  "2020-01-01T00:00:00Z",
	}); err != nil {
		t.Fatal(err)
	}
	mustRecordedChange(t, svc, task2.ID)
	got, err = svc.HasTestOutcomeSinceLatestChange(ctx, task2.ID)
	if err != nil || got {
		t.Fatalf("stale test before latest change → false: got=%v err=%v", got, err)
	}
}

func TestHasComputedEvaluation(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	got, err := svc.HasComputedEvaluation(ctx, task.ID)
	if err != nil || got {
		t.Fatalf("no eval → false: got=%v err=%v", got, err)
	}

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "ABCDEF012",
		ScoresJSON: `{"correctness": 0.98}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.95}`,
	}); err != nil {
		t.Fatal(err)
	}
	got, err = svc.HasComputedEvaluation(ctx, task.ID)
	if err != nil || !got {
		t.Fatalf("computed eval → true: got=%v err=%v", got, err)
	}
}

func TestHasReflectionSinceEvaluation(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	got, err := svc.HasReflectionSinceEvaluation(ctx, task.ID)
	if err != nil || got {
		t.Fatalf("no eval → false: got=%v err=%v", got, err)
	}

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "ABCDEF012",
		ScoresJSON: `{"correctness": 0.98}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.95}`,
	}); err != nil {
		t.Fatal(err)
	}
	got, err = svc.HasReflectionSinceEvaluation(ctx, task.ID)
	if err != nil || got {
		t.Fatalf("eval without reflection → false: got=%v err=%v", got, err)
	}

	if _, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID:      task.ID,
		Summary:     "cycle closed",
		UsefulTests: []string{"TestHasReflectionSinceEvaluation"},
	}); err != nil {
		t.Fatal(err)
	}
	got, err = svc.HasReflectionSinceEvaluation(ctx, task.ID)
	if err != nil || !got {
		t.Fatalf("reflection at/after eval → true: got=%v err=%v", got, err)
	}
}

func TestCompareIterationOutcomes(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestShip",
		TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := svc.CompareIterationOutcomes(ctx, task.ID, store.OutcomeKindTest)
	if err != nil {
		t.Fatalf("CompareIterationOutcomes: %v", err)
	}
	if got.Previous.TestStatus != store.TestStatusPass || got.Current.TestStatus != store.TestStatusFail {
		t.Fatalf("snapshots: previous=%+v current=%+v", got.Previous, got.Current)
	}
	if got.Delta.TestStatus == nil || got.Delta.TestStatus.From != store.TestStatusPass || got.Delta.TestStatus.To != store.TestStatusFail {
		t.Fatalf("test_status delta: %+v", got.Delta)
	}
}

func TestCompareIterationOutcomesEvaluation(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.90}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.90}`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.95}`,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := svc.CompareIterationOutcomes(ctx, task.ID, store.OutcomeKindEvaluation)
	if err != nil {
		t.Fatalf("CompareIterationOutcomes evaluation: %v", err)
	}
	dim, ok := got.Delta.Dimensions["correctness"]
	if !ok || dim.Delta == nil {
		t.Fatalf("correctness dimension delta: %+v", got.Delta.Dimensions)
	}
	if *dim.Delta < 0.04 || *dim.Delta > 0.06 {
		t.Fatalf("correctness delta=%v want ~0.05", *dim.Delta)
	}
}
