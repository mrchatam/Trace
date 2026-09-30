package loop_test

// Discoverability regressions for the deliberation loop, reported by a
// consumer project driving `trace loop` on every phase. See
// docs/research/LOOP-GATES-ENVELOPE-PROMPT.md (D1–D3) for the full report.
//
// These tests lock ergonomics only: every negative guard here proves the
// gates still refuse what they refused before the fixes.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/store"
)

// --- D1: the reflection error must name the keys a caller can reach for ---

func TestApplyReflectionEssayOnlyErrorNamesStructuredKeys(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Reflections: []loop.ApplyReflection{{
				Summary:          "closing reflection",
				BroadenTestsNote: "none",
			}},
		},
	}
	_, err := loop.Apply(ctx, st, psvc, env)
	if err == nil {
		t.Fatal("essay-only reflection must keep failing closed (CreateReflection contract)")
	}
	msg := err.Error()
	for _, key := range []string{"useful_tests", "new_dependencies", "invalidated_assumption_ids"} {
		if !strings.Contains(msg, key) {
			t.Fatalf("reflection error must name accepted key %q so a caller need not guess; got %q", key, msg)
		}
	}
}

func TestParseApplyEnvelopeRejectsUnknownReflectionKeys(t *testing.T) {
	// A caller's guessed key inside a reflection object must not be silently
	// dropped: silent-ignore is what turned 25 guesses into one key-less error.
	raw := []byte(`{"schema_version":"trace.loop.apply.v1","apply_id":"` + uuid.NewString() +
		`","seed":{"task_id":"` + uuid.NewString() + `","goal_id":"` + uuid.NewString() +
		`"},"writes":{"reflections":[{"invalidated_assumptions":["x"]}]}}`)
	if _, err := loop.ParseApplyEnvelope(raw); err == nil {
		t.Fatal("unknown key inside reflections[] must be rejected, not silently ignored")
	}
}

func TestApplyReflectionStructuredArraysAccepted(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	asm, err := dsvc.CreateAssumption(ctx, domain.AssumptionInput{Title: "stale assumption"})
	if err != nil {
		t.Fatalf("CreateAssumption: %v", err)
	}
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Reflections: []loop.ApplyReflection{{
				Summary:                  "structured reflection",
				InvalidatedAssumptionIDs: []string{asm.ID},
				NewDependencies:          []loop.ApplyDependency{{Kind: "path", Ref: "internal/loop"}},
				UsefulTests:              []string{"TestAll"},
			}},
		},
	}
	if _, err := loop.Apply(ctx, st, psvc, env); err != nil {
		t.Fatalf("structured reflection must apply: %v", err)
	}
	refs, err := dsvc.ListReflectionsByTaskID(ctx, taskID)
	if err != nil || len(refs) != 1 {
		t.Fatalf("ListReflectionsByTaskID: n=%d err=%v", len(refs), err)
	}
}

// --- D2: envelope validation must report every problem at once ---

func TestParseApplyEnvelopeReportsAllMissingFields(t *testing.T) {
	raw := []byte(`{"schema_version":"trace.loop.apply.v1","writes":{}}`)
	_, err := loop.ParseApplyEnvelope(raw)
	if err == nil {
		t.Fatal("want validation error for missing apply_id/seed")
	}
	msg := err.Error()
	for _, want := range []string{`"apply_id"`, `"seed"`} {
		if !strings.Contains(msg, want) {
			t.Fatalf("single-pass validation must name %s; got %q", want, msg)
		}
	}
}

func TestParseApplyEnvelopeReportsAllInvalidFields(t *testing.T) {
	raw := []byte(`{"schema_version":"trace.loop.apply.v0","apply_id":"nope","seed":{"task_id":"x","goal_id":"y"},` +
		`"writes":{"plan_changes":[{"title":"t"},{"id":"` + uuid.NewString() + `","discovery_id":"z"}]}}`)
	_, err := loop.ParseApplyEnvelope(raw)
	if err == nil {
		t.Fatal("want validation error")
	}
	msg := err.Error()
	for _, want := range []string{
		"schema_version",
		"apply_id",
		"seed.task_id",
		"seed.goal_id",
		"writes.plan_changes[0].id",
		"writes.plan_changes[1].discovery_id",
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("validation must report %s in one pass; got %q", want, msg)
		}
	}
}

// --- D3: gate refusals must carry the remedy at the point of failure ---
// (the negative guards below are the load-bearing assertions: the gates must
// keep refusing exactly what they refused before)

func TestEditGatePlanUncritiquedViolationCarriesRemedy(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	allowed, violations, err := loop.EvaluateGate(context.Background(), dsvc, psvc, st, taskID, loop.GateForEdit)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertBlocked(t, allowed, violations, "premature_implementation", "plan_uncritiqued")

	remedy := violations[0].Remedy
	if strings.TrimSpace(remedy) == "" {
		t.Fatal("plan_uncritiqued violation must carry a remedy")
	}
	for _, want := range []string{"writes.plan_changes", "TestGreenfield_MCPPlanBootstrap_EditGatePasses"} {
		if !strings.Contains(remedy, want) {
			t.Fatalf("remedy must mention %q; got %q", want, remedy)
		}
	}
}

func TestExecuteGateTestPendingViolationCarriesRemedy(t *testing.T) {
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

	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForExecute)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertBlocked(t, allowed, violations, "premature_implementation", "test_pending")

	remedy := violations[0].Remedy
	if !strings.Contains(remedy, "trace test run") {
		t.Fatalf("test_pending remedy must point at `trace test run`; got %q", remedy)
	}
}

func TestApplyReflectionDoesNotClearPlanUncritiquedButPlanChangesDo(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	// A well-formed reflection documents the cycle; it must NOT set plan_critiqued.
	reflEnv := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Reflections: []loop.ApplyReflection{{Summary: "critique notes", UsefulTests: []string{"TestAll"}}},
		},
	}
	if _, err := loop.Apply(ctx, st, psvc, reflEnv); err != nil {
		t.Fatalf("apply reflection: %v", err)
	}
	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForEdit)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertBlocked(t, allowed, violations, "premature_implementation", "plan_uncritiqued")

	// Only a loop apply with writes.plan_changes[] clears plan_uncritiqued.
	pcEnv := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			PlanChanges: []loop.ApplyPlanChange{{ID: uuid.NewString(), Title: "critique fix"}},
		},
	}
	if _, err := loop.Apply(ctx, st, psvc, pcEnv); err != nil {
		t.Fatalf("apply plan_changes: %v", err)
	}
	allowed, violations, err = loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForEdit)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertAllowed(t, allowed, violations)
}

// Empty applies saturate (P19) and the resulting STOP covers DONE. This locks
// the observed fail-open valve as documented behavior: an apply is never
// evidence of work, but saturation is the deliberate escape hatch that ends a
// stalled loop.
func TestEmptyAppliesSaturateAndStopCoversDone(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForDone)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertBlocked(t, allowed, violations, "premature_implementation", "execute_pending")

	for i := 0; i < 2; i++ {
		env := loop.ApplyEnvelope{
			SchemaVersion: loop.ApplySchemaVersion,
			ApplyID:       uuid.NewString(),
			Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		}
		res, err := loop.Apply(ctx, st, psvc, env)
		if err != nil {
			t.Fatalf("empty apply %d: %v", i+1, err)
		}
		if i == 0 && res.Saturated {
			t.Fatal("first empty apply must not saturate (threshold=2)")
		}
		if i == 1 && !res.Saturated {
			t.Fatal("second consecutive empty apply must saturate")
		}
	}

	allowed, violations, err = loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForDone)
	if err != nil {
		t.Fatalf("EvaluateGate: %v", err)
	}
	assertAllowed(t, allowed, violations)
}

// --- D4: the discovery-plan-change error must name both id spaces ---

func TestLinkDiscoveryPlanChangeScopeIDErrorNamesBothIDSpaces(t *testing.T) {
	dsvc, st := openDomainSvc(t)
	ctx := context.Background()

	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "probe discovery"})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	goal, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "probe goal"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	ph, err := st.InsertPlanPhase(store.PlanPhase{GoalID: goal.ID, Title: "Phase 1"})
	if err != nil {
		t.Fatalf("InsertPlanPhase: %v", err)
	}
	scope, err := st.InsertPlanScope(store.PlanScope{PhaseID: ph.ID, Title: "Scope 1"})
	if err != nil {
		t.Fatalf("InsertPlanScope: %v", err)
	}

	err = dsvc.LinkDiscoveryPlanChange(ctx, disc.ID, scope.ID, domain.LinkMeta{})
	if err == nil {
		t.Fatal("plan_scopes id is not a plan_changes id; link must fail closed")
	}
	msg := err.Error()
	for _, want := range []string{"plan_changes", "plan_scopes"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("error must name both id spaces (%s); got %q", want, msg)
		}
	}
}

func TestLinkDiscoveryPlanChangeUnknownTargetNamesPlanChanges(t *testing.T) {
	dsvc, _ := openDomainSvc(t)
	ctx := context.Background()

	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "probe discovery"})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	err = dsvc.LinkDiscoveryPlanChange(ctx, disc.ID, uuid.NewString(), domain.LinkMeta{})
	if err == nil {
		t.Fatal("unknown target must fail closed")
	}
	if msg := err.Error(); !strings.Contains(msg, "plan_changes") {
		t.Fatalf("error must name the plan_changes id space; got %q", msg)
	}
}

// openDomainSvc is a small local helper so this file does not depend on the
// domain package's internal test helpers.
func openDomainSvc(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, psvc, dsvc := openLoopTestStore(t)
	_ = psvc
	return dsvc, st
}

// firstGoalIDForTask resolves a task's goal_id from the store.
func firstGoalIDForTask(t *testing.T, st *store.Store, taskID string) (string, error) {
	t.Helper()
	task, err := st.GetTask(taskID)
	if err != nil {
		return "", err
	}
	if task.GoalID == nil {
		return "", fmt.Errorf("task %s has no goal_id", taskID)
	}
	return *task.GoalID, nil
}

// guard: ApplyResult keeps exposing the created plan_changes ids (D4 echo).
func TestApplyResultEncodesPlanChangeIDs(t *testing.T) {
	raw, err := json.Marshal(loop.ApplyResult{PlanChangeIDs: []string{uuid.NewString()}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["plan_change_ids"]; !ok {
		t.Fatalf("ApplyResult must expose plan_change_ids; got %s", raw)
	}
}
