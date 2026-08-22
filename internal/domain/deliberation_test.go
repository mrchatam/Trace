package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
)

func TestApplyDeliberationTransitionPersistsEvent(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "deliberate"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	inputs := deliberation.PolicyInputs{BlockingUncertaintyCount: 1, PlanExists: true, PlanCritiqued: true}
	next, ev, err := svc.ApplyDeliberationTransition(ctx, task.ID, g.ID, inputs)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if next.CurrentPhase != deliberation.PhaseInvestigate || next.StopReason != "" {
		t.Fatalf("state: %+v", next)
	}
	if next.HopCount != 1 || next.LastPhase != deliberation.PhaseOrient {
		t.Fatalf("hop/last: %+v", next)
	}
	if ev.Type != deliberation.EventTransition || ev.EntityType != domain.EntityTask || ev.EntityID != task.ID {
		t.Fatalf("event meta: %+v", ev)
	}

	events, err := st.ListEventsByEntity(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range events {
		if e.Type != deliberation.EventTransition {
			continue
		}
		found = true
		var payload map[string]any
		if err := json.Unmarshal([]byte(e.PayloadJSON), &payload); err != nil {
			t.Fatalf("payload json: %v", err)
		}
		for _, key := range []string{"task_id", "goal_id", "from_phase", "to_phase", "reason_code", "hop_count", "policy_inputs"} {
			if _, ok := payload[key]; !ok {
				t.Errorf("missing payload field %s", key)
			}
		}
		if payload["task_id"] != task.ID || payload["goal_id"] != g.ID {
			t.Errorf("ids: %v", payload)
		}
		if payload["from_phase"] != "ORIENT" || payload["to_phase"] != "INVESTIGATE" {
			t.Errorf("phases: %v", payload)
		}
		if payload["reason_code"] != "blocking_uncertainty" {
			t.Errorf("reason: %v", payload["reason_code"])
		}
		if payload["to_phase"] == "EXECUTE" {
			t.Error("EXECUTE persisted with blocking uncertainty")
		}
		pi, ok := payload["policy_inputs"].(map[string]any)
		if !ok {
			t.Fatalf("policy_inputs: %T", payload["policy_inputs"])
		}
		for _, key := range []string{
			"blocking_uncertainty_count", "plan_exists", "plan_critiqued",
			"verification_incomplete", "open_regression", "p19_saturated",
		} {
			if _, ok := pi[key]; !ok {
				t.Errorf("missing policy_inputs.%s", key)
			}
		}
		if pi["blocking_uncertainty_count"] != float64(1) {
			t.Errorf("blocking_uncertainty_count: %v", pi["blocking_uncertainty_count"])
		}
	}
	if !found {
		t.Fatalf("no deliberation.transition in events: %+v", events)
	}

	loaded, err := svc.GetDeliberationState(ctx, task.ID)
	if err != nil || loaded.CurrentPhase != deliberation.PhaseInvestigate {
		t.Fatalf("Get: %+v err=%v", loaded, err)
	}
}

func TestResetDeliberationStateClearsStopPreservesCritique(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "reset goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = svc.ApplyDeliberationTransition(ctx, task.ID, g.ID, deliberation.PolicyInputs{
		PlanExists: true, PlanCritiqued: true, P19Saturated: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	before, err := svc.GetDeliberationState(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !before.Stopped || !before.PlanCritiqued {
		t.Fatalf("precondition: %+v", before)
	}

	after, err := svc.ResetDeliberationState(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Stopped || after.StopReason != "" || after.HopCount != 0 || after.ConsecutiveEmptyApplies != 0 {
		t.Fatalf("reset: %+v", after)
	}
	if after.CurrentPhase != deliberation.PhaseExecute {
		t.Fatalf("phase=%s want EXECUTE", after.CurrentPhase)
	}
	if !after.PlanCritiqued {
		t.Fatal("plan_critiqued must be preserved")
	}
}

func TestApplyDeliberationTransitionPlanMissing(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	next, _, err := svc.ApplyDeliberationTransition(ctx, task.ID, g.ID, deliberation.PolicyInputs{})
	if err != nil {
		t.Fatal(err)
	}
	if next.CurrentPhase != deliberation.PhasePlan {
		t.Fatalf("phase %s want PLAN", next.CurrentPhase)
	}
}

func TestApplyDeliberationTransitionRequiresIDs(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	if _, _, err := svc.ApplyDeliberationTransition(ctx, "", "g", deliberation.PolicyInputs{}); err == nil {
		t.Fatal("expected task_id validation")
	}
	if _, _, err := svc.ApplyDeliberationTransition(ctx, "t", "", deliberation.PolicyInputs{}); err == nil {
		t.Fatal("expected goal_id validation")
	}
}

func TestApplyDeliberationTransitionRequiresMatchingGoalID(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	gA, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "goal A"})
	if err != nil {
		t.Fatal(err)
	}
	gB, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "goal B"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &gA.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ApplyDeliberationTransition(ctx, task.ID, gB.ID, deliberation.PolicyInputs{PlanExists: true})
	if err == nil {
		t.Fatal("expected goal_id mismatch error")
	}
	var val *domain.ErrValidation
	if !errors.As(err, &val) || val.Msg != "goal_id does not match task" {
		t.Fatalf("want ErrValidation goal_id mismatch, got %v", err)
	}

	evs, err := st.ListEventsByEntity(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == domain.EventDeliberationTransition {
			t.Fatalf("must not persist transition on goal mismatch: %+v", e)
		}
	}
}
