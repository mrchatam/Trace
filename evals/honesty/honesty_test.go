package honesty_test

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func openHonesty(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st), st
}

// TestHonestyFailClosedPlantedClaim is the named H5 partial honesty demo.
// Paths A/B reject planted completion; Path C recovers by superseding FAIL
// then a PASS review + AllowOperatorDone (DF-43).
// AllowDoneWithoutReview is never set true.
func TestHonestyFailClosedPlantedClaim(t *testing.T) {
	svc, st := openHonesty(t)
	ctx := context.Background()

	// Setup: task → IN_PROGRESS
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Ship feature X"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "implementer", Reason: "start work",
	}); err != nil {
		t.Fatalf("→ IN_PROGRESS: %v", err)
	}

	// Planted false / incomplete claim (implementer narrative)
	claim, err := svc.CreateClaim(ctx, domain.ClaimInput{Title: "Feature X is complete"})
	if err != nil {
		t.Fatalf("CreateClaim: %v", err)
	}
	evidence, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "tests pass on my machine"})
	if err != nil {
		t.Fatalf("CreateEvidence: %v", err)
	}
	if err := svc.LinkClaimEvidence(ctx, claim.ID, evidence.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkClaimEvidence: %v", err)
	}

	// Path A — EvidenceIDs alone must NOT unlock DONE
	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "claim complete", EvidenceIDs: []string{evidence.ID},
		AllowDoneWithoutReview: false,
	})
	assertDoneRejected(t, err)
	assertWorkState(t, svc, ctx, task.ID, store.WorkStateInProgress)

	// Path B — Independent review FAIL → still no DONE
	revFail, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Honesty check: completion claim"})
	if err != nil {
		t.Fatalf("CreateReview FAIL path: %v", err)
	}
	if err := svc.LinkReviewTask(ctx, revFail.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewTask FAIL: %v", err)
	}
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultFail, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "claim unproven / incomplete",
	}); err != nil {
		t.Fatalf("SetReviewResult FAIL: %v", err)
	}
	gotFail, err := svc.GetReview(ctx, revFail.ID)
	if err != nil || gotFail.Result != store.ReviewResultFail {
		t.Fatalf("GetReview after FAIL: %+v err=%v", gotFail, err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "retry after fail review",
		AllowDoneWithoutReview: false,
	})
	assertDoneRejectedForFail(t, err)
	assertWorkState(t, svc, ctx, task.ID, store.WorkStateInProgress)

	// Path C — supersede FAIL → UNCERTAIN, then PASS + AllowOperatorDone → DONE (DF-43)
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultUncertain, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "superseded by later review",
	}); err != nil {
		t.Fatalf("supersede FAIL→UNCERTAIN: %v", err)
	}
	superseded, err := svc.GetReview(ctx, revFail.ID)
	if err != nil || superseded.Result != store.ReviewResultUncertain {
		t.Fatalf("FAIL review must be UNCERTAIN after supersede: %+v err=%v", superseded, err)
	}

	revPass, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Re-review after remediation"})
	if err != nil {
		t.Fatalf("CreateReview PASS path: %v", err)
	}
	if err := svc.LinkReviewTask(ctx, revPass.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewTask PASS: %v", err)
	}
	if err := svc.SetReviewResult(ctx, revPass.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "evidence now adequate",
	}); err != nil {
		t.Fatalf("SetReviewResult PASS: %v", err)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "promote after PASS",
		AllowDoneWithoutReview: false,
		AllowOperatorDone:      true, // DF-17 Path C: operator flag required
	}); err != nil {
		t.Fatalf("Path C DONE after supersede+PASS: %v", err)
	}
	assertWorkState(t, svc, ctx, task.ID, store.WorkStateDone)

	// Optional: task.transition payload may include review_id == revPass.ID
	evs, err := st.ListEventsByEntity(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatalf("ListEventsByEntity: %v", err)
	}
	var donePayload map[string]any
	for i := range evs {
		if evs[i].Type != domain.EventTaskTransition {
			continue
		}
		var p map[string]any
		if json.Unmarshal([]byte(evs[i].PayloadJSON), &p) != nil {
			continue
		}
		if p["to"] == store.WorkStateDone {
			donePayload = p
		}
	}
	if donePayload == nil {
		t.Fatal("expected DONE task.transition event")
	}
	if rid, _ := donePayload["review_id"].(string); rid != revPass.ID {
		t.Fatalf("DONE payload review_id want %s, got %#v", revPass.ID, donePayload["review_id"])
	}

	// Prior FAIL review remains linked but was superseded to UNCERTAIN (not auto-cleared by PASS)
	stillLinked, err := svc.GetReview(ctx, revFail.ID)
	if err != nil || stillLinked.Result != store.ReviewResultUncertain {
		t.Fatalf("superseded review must remain UNCERTAIN: %+v err=%v", stillLinked, err)
	}
}

func assertDoneRejected(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected DONE rejection")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want *ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "Review") || !strings.Contains(inv.Reason, "AllowDoneWithoutReview") {
		t.Fatalf("reason should mention Review PASS / AllowDoneWithoutReview: %q", inv.Reason)
	}
	if !strings.Contains(inv.Reason, "PASS") {
		t.Fatalf("reason should mention PASS: %q", inv.Reason)
	}
	if !strings.Contains(inv.Reason, "AllowOperatorDone") && !strings.Contains(inv.Reason, "as-operator") {
		t.Fatalf("reason should mention AllowOperatorDone/--as-operator: %q", inv.Reason)
	}
}

func assertDoneRejectedForFail(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected DONE rejection")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want *ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "FAIL") {
		t.Fatalf("reason should mention linked FAIL: %q", inv.Reason)
	}
}

func assertWorkState(t *testing.T, svc *domain.Service, ctx context.Context, taskID, want string) {
	t.Helper()
	got, err := svc.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if got.WorkState != want {
		t.Fatalf("WorkState want %q, got %q", want, got.WorkState)
	}
}

// MetricsGateG is the Gate G prelim artifact shape (schema-gate-g.json v1).
type MetricsGateG struct {
	SchemaVersion        int            `json:"schema_version"`
	Gate                 string         `json:"gate"`
	Suite                string         `json:"suite"`
	Prelim               bool           `json:"prelim"`
	DryRun               bool           `json:"dry_run"`
	Attempts             int            `json:"attempts"`
	Escapes              int            `json:"escapes"`
	Caught               int            `json:"caught"`
	EscapeRate           float64        `json:"escape_rate"`
	OpenResidualsTotal   int            `json:"open_residuals_total"`
	OpenResidualsByScope map[string]int `json:"open_residuals_by_scope,omitempty"`
	PathsABCTest         string         `json:"paths_abc_test"`
	NamedTest            string         `json:"named_test"`
	S01Hooks             []string       `json:"s01_hooks"`
	TraceVersion         string         `json:"trace_version,omitempty"`
}

// TestHonestyEscapeRateGateGPrelim is the Gate G preliminary escape-rate report.
// Paths A/B/C remain covered by TestHonestyFailClosedPlantedClaim (not folded here).
// Escape = DONE via AllowDoneWithoutReview without linked Review PASS.
// Path C-style PASS→DONE is excluded from attempts (remediation, not escape).
func TestHonestyEscapeRateGateGPrelim(t *testing.T) {
	svc, st := openHonesty(t)
	psvc := planner.New(st)
	ctx := context.Background()

	escapes, caught := runEscapeRatePlantedCases(t, svc, ctx)
	attempts := escapes + caught
	if escapes != 1 || caught != 2 || attempts != 3 {
		t.Fatalf("tallies want escapes=1 caught=2 attempts=3; got escapes=%d caught=%d attempts=%d",
			escapes, caught, attempts)
	}
	escapeRate := float64(escapes) / float64(attempts)
	if escapeRate <= 0.33 || escapeRate >= 0.34 {
		if math.Abs(escapeRate-1.0/3.0) > 1e-9 {
			t.Fatalf("escape_rate want ≈1/3 (in (0.33,0.34)); got %v", escapeRate)
		}
	}

	openTotal, byScope := runS01ResidualTally(t, svc, psvc, ctx)
	if openTotal != 1 {
		t.Fatalf("open_residuals_total want 1, got %d", openTotal)
	}

	metricsDir := t.TempDir()
	metricsPath := filepath.Join(metricsDir, "metrics-gate-g.json")
	m := MetricsGateG{
		SchemaVersion:        1,
		Gate:                 "G",
		Suite:                "honesty",
		Prelim:               true,
		DryRun:               false,
		Attempts:             attempts,
		Escapes:              escapes,
		Caught:               caught,
		EscapeRate:           escapeRate,
		OpenResidualsTotal:   openTotal,
		OpenResidualsByScope: byScope,
		PathsABCTest:         "TestHonestyFailClosedPlantedClaim",
		NamedTest:            "TestHonestyEscapeRateGateGPrelim",
		S01Hooks: []string{
			"review_judges_scope",
			"CountOpenResidualsByScope",
			"POLICY_EXCEPTION",
		},
		TraceVersion: "0.0.0-dev",
	}
	writeGateGMetrics(t, metricsPath, m)
	validateGateGMetricsFile(t, loadGateGSchema(t), metricsPath)
}

func runEscapeRatePlantedCases(t *testing.T, svc *domain.Service, ctx context.Context) (escapes, caught int) {
	t.Helper()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Gate G planted completion"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "implementer", Reason: "start work",
	}); err != nil {
		t.Fatalf("→ IN_PROGRESS: %v", err)
	}

	evidence, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "tests pass on my machine"})
	if err != nil {
		t.Fatalf("CreateEvidence: %v", err)
	}

	// Caught-1 (mirror Path A): EvidenceIDs-alone DONE → reject
	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "claim complete", EvidenceIDs: []string{evidence.ID},
		AllowDoneWithoutReview: false,
	})
	assertDoneRejected(t, err)
	assertWorkState(t, svc, ctx, task.ID, store.WorkStateInProgress)
	caught++

	// Caught-2 (mirror Path B): FAIL review → DONE reject
	revFail, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Gate G FAIL review"})
	if err != nil {
		t.Fatalf("CreateReview FAIL: %v", err)
	}
	if err := svc.LinkReviewTask(ctx, revFail.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewTask FAIL: %v", err)
	}
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultFail, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "claim unproven",
	}); err != nil {
		t.Fatalf("SetReviewResult FAIL: %v", err)
	}
	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "retry after fail review",
		AllowDoneWithoutReview: false,
	})
	assertDoneRejectedForFail(t, err)
	assertWorkState(t, svc, ctx, task.ID, store.WorkStateInProgress)
	caught++

	// Escape-1: NEW task; AllowDoneWithoutReview hatch → DONE (report-only escape)
	escTask, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Gate G hatch escape"})
	if err != nil {
		t.Fatalf("CreateTask escape: %v", err)
	}
	if err := svc.TransitionTask(ctx, escTask.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "implementer", Reason: "start escape fixture",
	}); err != nil {
		t.Fatalf("escape → IN_PROGRESS: %v", err)
	}
	if err := svc.TransitionTask(ctx, escTask.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "implementer", Reason: "document hatch as escape",
		AllowDoneWithoutReview: true,
	}); err != nil {
		t.Fatalf("Escape-1 AllowDoneWithoutReview MUST succeed: %v", err)
	}
	assertWorkState(t, svc, ctx, escTask.ID, store.WorkStateDone)
	escapes++

	return escapes, caught
}

func runS01ResidualTally(t *testing.T, svc *domain.Service, psvc *planner.Service, ctx context.Context) (openTotal int, byScope map[string]int) {
	t.Helper()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Gate G residual goal"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: goal.ID,
		Phases: []planner.PhaseInput{{
			Title:  "Phase 1",
			Scopes: []planner.ScopeInput{{Title: "Scope A"}},
		}},
		Actor: "gate-g-prelim",
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	if len(cp.Phases) < 1 || len(cp.Phases[0].Scopes) < 1 {
		t.Fatalf("CreateCoarsePlan must yield ≥1 phase/scope; got %+v", cp)
	}
	scopeID := cp.Phases[0].Scopes[0].ID

	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Scope residual review"})
	if err != nil {
		t.Fatalf("CreateReview: %v", err)
	}
	if err := svc.LinkReviewScope(ctx, rev.ID, scopeID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewScope (S01 fail-closed): %v", err)
	}
	res, err := svc.AddResidual(ctx, rev.ID, domain.ResidualInput{
		Code:     domain.ResidualCodePolicyException,
		Body:     "planted POLICY_EXCEPTION for Gate G prelim",
		Severity: domain.ResidualSeverityWARN,
	})
	if err != nil {
		t.Fatalf("AddResidual POLICY_EXCEPTION (S01 fail-closed): %v", err)
	}
	if res.Status != store.ResidualStatusOpen {
		t.Fatalf("residual status want OPEN, got %q", res.Status)
	}
	if res.Code != domain.ResidualCodePolicyException {
		t.Fatalf("residual code want %s, got %q", domain.ResidualCodePolicyException, res.Code)
	}

	n, err := svc.CountOpenResidualsByScope(ctx, scopeID)
	if err != nil {
		t.Fatalf("CountOpenResidualsByScope (S01 fail-closed): %v", err)
	}
	if n != 1 {
		t.Fatalf("CountOpenResidualsByScope want 1, got %d", n)
	}
	listed, err := svc.ListResidualsByScope(ctx, scopeID)
	if err != nil {
		t.Fatalf("ListResidualsByScope: %v", err)
	}
	found := false
	for _, r := range listed {
		if r.ID == res.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("ListResidualsByScope missing residual %s", res.ID)
	}

	byScope = map[string]int{scopeID: n}
	return n, byScope
}

func findHonestyModuleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func honestyModuleRoot(t *testing.T) string {
	t.Helper()
	root := findHonestyModuleRoot()
	if root == "" {
		t.Fatal("go.mod not found above evals/honesty")
	}
	return root
}

func loadGateGSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(honestyModuleRoot(t), "evals", "honesty", "schema-gate-g.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema-gate-g.json: %v", err)
	}
	return sch
}

func writeGateGMetrics(t *testing.T, path string, m MetricsGateG) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validateGateGMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read metrics %s: %v", path, err)
	}
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal metrics: %v", err)
	}
	if err := sch.Validate(doc); err != nil {
		t.Fatalf("schema validation failed for %s: %v\n%s", path, err, raw)
	}
}
