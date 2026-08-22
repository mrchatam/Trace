package loop_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
)

func emptyApplyEnv(taskID, goalID, applyID string) loop.ApplyEnvelope {
	return loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       applyID,
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
}

func TestApplyConsecutiveEmptySaturationThreshold(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	first, err := loop.Apply(ctx, st, psvc, emptyApplyEnv(taskID, goalID, uuid.NewString()))
	if err != nil {
		t.Fatal(err)
	}
	if first.Saturated {
		t.Fatalf("first pure empty must not saturate: %+v", first)
	}
	state, err := dsvc.GetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if state.Stopped || state.ConsecutiveEmptyApplies != 1 {
		t.Fatalf("after first empty: %+v", state)
	}

	second, err := loop.Apply(ctx, st, psvc, emptyApplyEnv(taskID, goalID, uuid.NewString()))
	if err != nil {
		t.Fatal(err)
	}
	if !second.Saturated {
		t.Fatalf("second pure empty must saturate: %+v", second)
	}
	state, err = dsvc.GetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Stopped || state.StopReason != string(deliberation.ReasonP19Saturated) {
		t.Fatalf("after second empty: %+v", state)
	}
	if state.ConsecutiveEmptyApplies != 2 {
		t.Fatalf("consecutive=%d want 2", state.ConsecutiveEmptyApplies)
	}
}

func TestApplyDiscoveriesOnlyDoesNotIncrementEmptyCounter(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	discID := uuid.NewString()
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Discoveries: []loop.ApplyDiscovery{{
				ID: discID, Title: "gap finding", Body: "body",
			}},
		},
	}
	out, err := loop.Apply(ctx, st, psvc, env)
	if err != nil {
		t.Fatal(err)
	}
	if out.Saturated || out.NewDiscoveries != 1 {
		t.Fatalf("discoveries-only: %+v", out)
	}
	state, err := dsvc.GetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if state.ConsecutiveEmptyApplies != 0 || state.Stopped {
		t.Fatalf("discoveries-only must not increment/stop: %+v", state)
	}
}

func TestResetClearsSaturationAndPreventsImmediateReStop(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	for i := 0; i < 2; i++ {
		if _, err := loop.Apply(ctx, st, psvc, emptyApplyEnv(taskID, goalID, uuid.NewString())); err != nil {
			t.Fatal(err)
		}
	}
	state, err := dsvc.GetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Stopped {
		t.Fatal("expected saturated STOP before reset")
	}

	reset, err := dsvc.ResetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if reset.Stopped || reset.StopReason != "" || reset.HopCount != 0 || reset.ConsecutiveEmptyApplies != 0 {
		t.Fatalf("reset state: %+v", reset)
	}
	if reset.CurrentPhase != deliberation.PhaseExecute {
		t.Fatalf("phase=%s want EXECUTE", reset.CurrentPhase)
	}

	after, err := loop.Apply(ctx, st, psvc, emptyApplyEnv(taskID, goalID, uuid.NewString()))
	if err != nil {
		t.Fatal(err)
	}
	if after.Saturated {
		t.Fatalf("first empty after reset must not re-STOP: %+v", after)
	}
	state, err = dsvc.GetDeliberationState(ctx, taskID)
	if err != nil {
		t.Fatal(err)
	}
	if state.Stopped || state.StopReason == string(deliberation.ReasonP19Saturated) {
		t.Fatalf("sticky p19_saturated after reset: %+v", state)
	}

	discEnv := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Discoveries: []loop.ApplyDiscovery{{
				ID: uuid.NewString(), Title: "post-reset finding",
			}},
		},
	}
	discOut, err := loop.Apply(ctx, st, psvc, discEnv)
	if err != nil {
		t.Fatal(err)
	}
	if discOut.Saturated {
		t.Fatalf("discoveries-only after reset must not saturate: %+v", discOut)
	}
}

func TestExportStopReasonMatchesGateAfterSaturation(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	for i := 0; i < 2; i++ {
		if _, err := loop.Apply(ctx, st, psvc, emptyApplyEnv(taskID, goalID, uuid.NewString())); err != nil {
			t.Fatal(err)
		}
	}

	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, taskID, loop.GateForEdit)
	if err != nil {
		t.Fatal(err)
	}
	if allowed || len(violations) == 0 {
		t.Fatalf("expected edit blocked after saturation, allowed=%v violations=%v", allowed, violations)
	}
	gateReason := violations[0].ReasonCode

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var exportReason string
	for _, ds := range doc.DeliberationStates {
		if ds.TaskID == taskID {
			exportReason = ds.StopReason
			break
		}
	}
	if exportReason == "" {
		t.Fatal("missing deliberation_states stop_reason in export")
	}
	if gateReason != exportReason {
		t.Fatalf("gate reason_code=%q export stop_reason=%q", gateReason, exportReason)
	}
	if gateReason != string(deliberation.ReasonP19Saturated) {
		t.Fatalf("want p19_saturated, got %q", gateReason)
	}

	status, err := loop.Status(ctx, st, psvc, loop.ApplySeed{TaskID: taskID, GoalID: goalID})
	if err != nil {
		t.Fatal(err)
	}
	if status.Deliberation == nil || status.Deliberation.StopReason != exportReason {
		t.Fatalf("status stop_reason=%v want %q", status.Deliberation, exportReason)
	}
	if string(status.Deliberation.WhySelected) != exportReason {
		t.Fatalf("why_selected=%q want %q", status.Deliberation.WhySelected, exportReason)
	}
}
