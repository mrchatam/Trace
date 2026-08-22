package deliberation_test

import (
	"encoding/json"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
)

func TestSelectNext(t *testing.T) {
	clear := deliberation.PolicyInputs{
		PlanExists:    true,
		PlanCritiqued: true,
	}
	tests := []struct {
		name       string
		state      deliberation.State
		inputs     deliberation.PolicyInputs
		wantPhase  deliberation.Phase
		wantReason deliberation.ReasonCode
		wantStop   bool
	}{
		{
			name:       "blocking_uncertainty_count=1 → INVESTIGATE / blocking_uncertainty (never EXECUTE)",
			state:      deliberation.State{HopCount: 3},
			inputs:     deliberation.PolicyInputs{BlockingUncertaintyCount: 1, PlanExists: true, PlanCritiqued: true},
			wantPhase:  deliberation.PhaseInvestigate,
			wantReason: deliberation.ReasonBlockingUncertainty,
			wantStop:   false,
		},
		{
			name:       "plan_exists=false, uncertainties clear → PLAN / plan_missing",
			state:      deliberation.State{HopCount: 1},
			inputs:     deliberation.PolicyInputs{},
			wantPhase:  deliberation.PhasePlan,
			wantReason: deliberation.ReasonPlanMissing,
			wantStop:   false,
		},
		{
			name:       "plan_exists=true, plan_critiqued=false → CRITIQUE / plan_uncritiqued",
			state:      deliberation.State{HopCount: 1},
			inputs:     deliberation.PolicyInputs{PlanExists: true, PlanCritiqued: false},
			wantPhase:  deliberation.PhaseCritique,
			wantReason: deliberation.ReasonPlanUncritiqued,
			wantStop:   false,
		},
		{
			name:       "open_regression=true → INVESTIGATE / open_regression",
			state:      deliberation.State{HopCount: 2},
			inputs:     deliberation.PolicyInputs{OpenRegression: true, PlanExists: true, PlanCritiqued: true},
			wantPhase:  deliberation.PhaseInvestigate,
			wantReason: deliberation.ReasonOpenRegression,
			wantStop:   false,
		},
		{
			name:       "hop_count=12 → STOP / hop_budget_exceeded",
			state:      deliberation.State{HopCount: 12},
			inputs:     clear,
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonHopBudgetExceeded,
			wantStop:   true,
		},
		{
			name:       "p19_saturated=true → STOP / p19_saturated",
			state:      deliberation.State{HopCount: 4},
			inputs:     deliberation.PolicyInputs{P19Saturated: true, PlanExists: true, PlanCritiqued: true},
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonP19Saturated,
			wantStop:   true,
		},
		{
			name:       "verification_incomplete → VERIFY / verification_incomplete",
			state:      deliberation.State{HopCount: 2},
			inputs:     deliberation.PolicyInputs{PlanExists: true, PlanCritiqued: true, VerificationIncomplete: true},
			wantPhase:  deliberation.PhaseVerify,
			wantReason: deliberation.ReasonVerificationIncomplete,
			wantStop:   false,
		},
		{
			name:       "default → ORIENT / continue_orient",
			state:      deliberation.State{HopCount: 0},
			inputs:     clear,
			wantPhase:  deliberation.PhaseOrient,
			wantReason: deliberation.ReasonContinueOrient,
			wantStop:   false,
		},
		{
			name:       "stopped → STOP / persisted p19_saturated",
			state:      deliberation.State{HopCount: 5, Stopped: true, StopReason: string(deliberation.ReasonP19Saturated)},
			inputs:     clear,
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonP19Saturated,
			wantStop:   true,
		},
		{
			name:       "stopped with empty StopReason → hop_budget_exceeded fallback",
			state:      deliberation.State{HopCount: 5, Stopped: true, StopReason: ""},
			inputs:     clear,
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonHopBudgetExceeded,
			wantStop:   true,
		},
		{
			name:       "hop_count=12 beats p19_saturated (first-match)",
			state:      deliberation.State{HopCount: 12},
			inputs:     deliberation.PolicyInputs{P19Saturated: true},
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonHopBudgetExceeded,
			wantStop:   true,
		},
		{
			name:       "hop_count=12 beats stopped persisted reason",
			state:      deliberation.State{HopCount: 12, Stopped: true, StopReason: string(deliberation.ReasonP19Saturated)},
			inputs:     clear,
			wantPhase:  deliberation.PhaseStop,
			wantReason: deliberation.ReasonHopBudgetExceeded,
			wantStop:   true,
		},
		{
			name:       "blocking_uncertainty beats open_regression and missing plan",
			state:      deliberation.State{HopCount: 1},
			inputs:     deliberation.PolicyInputs{BlockingUncertaintyCount: 2, OpenRegression: true, PlanExists: false},
			wantPhase:  deliberation.PhaseInvestigate,
			wantReason: deliberation.ReasonBlockingUncertainty,
			wantStop:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			phase, reason, stop := deliberation.SelectNext(tc.state, tc.inputs)
			if phase != tc.wantPhase || reason != tc.wantReason || stop != tc.wantStop {
				t.Fatalf("SelectNext = (%s, %s, %v) want (%s, %s, %v)",
					phase, reason, stop, tc.wantPhase, tc.wantReason, tc.wantStop)
			}
			if tc.inputs.BlockingUncertaintyCount > 0 && phase == deliberation.PhaseExecute {
				t.Fatalf("SelectNext returned EXECUTE with blocking_uncertainty_count=%d", tc.inputs.BlockingUncertaintyCount)
			}
		})
	}
}

func TestSelectNextNeverExecuteOnBlockingUncertainty(t *testing.T) {
	// EXECUTE-path simulation: plan ready, no verify debt, no regression — still not EXECUTE.
	phase, reason, stop := deliberation.SelectNext(
		deliberation.State{HopCount: 8, CurrentPhase: deliberation.PhaseExecute},
		deliberation.PolicyInputs{
			BlockingUncertaintyCount: 1,
			PlanExists:               true,
			PlanCritiqued:            true,
			ExecutePending:           true,
			VerificationIncomplete:   false,
			OpenRegression:           false,
			P19Saturated:             false,
		},
	)
	if phase == deliberation.PhaseExecute {
		t.Fatal("SelectNext returned EXECUTE when blocking_uncertainty_count > 0")
	}
	if phase != deliberation.PhaseInvestigate || reason != deliberation.ReasonBlockingUncertainty || stop {
		t.Fatalf("got (%s, %s, %v) want INVESTIGATE / blocking_uncertainty / false", phase, reason, stop)
	}
}

func cycleClearBase() deliberation.PolicyInputs {
	return deliberation.PolicyInputs{
		PlanExists:    true,
		PlanCritiqued: true,
	}
}

func TestSelectNextExecuteWhenPending(t *testing.T) {
	inputs := cycleClearBase()
	inputs.ExecutePending = true
	phase, reason, stop := deliberation.SelectNext(deliberation.State{HopCount: 3}, inputs)
	if phase != deliberation.PhaseExecute || reason != deliberation.ReasonExecutePending || stop {
		t.Fatalf("got (%s, %s, %v) want EXECUTE / execute_pending / false", phase, reason, stop)
	}
}

func TestSelectNextTestWhenExecuteDone(t *testing.T) {
	inputs := cycleClearBase()
	inputs.TestPending = true
	phase, reason, stop := deliberation.SelectNext(deliberation.State{HopCount: 4}, inputs)
	if phase != deliberation.PhaseTest || reason != deliberation.ReasonTestPending || stop {
		t.Fatalf("got (%s, %s, %v) want TEST / test_pending / false", phase, reason, stop)
	}
}

func TestSelectNextEvaluateWhenTestPass(t *testing.T) {
	inputs := cycleClearBase()
	inputs.EvaluationPending = true
	phase, reason, stop := deliberation.SelectNext(deliberation.State{HopCount: 5}, inputs)
	if phase != deliberation.PhaseEvaluate || reason != deliberation.ReasonEvaluationPending || stop {
		t.Fatalf("got (%s, %s, %v) want EVALUATE / evaluation_pending / false", phase, reason, stop)
	}
}

func TestSelectNextReflectWhenEvaluationDone(t *testing.T) {
	inputs := cycleClearBase()
	inputs.ReflectPending = true
	phase, reason, stop := deliberation.SelectNext(deliberation.State{HopCount: 6}, inputs)
	if phase != deliberation.PhaseReflect || reason != deliberation.ReasonReflectPending || stop {
		t.Fatalf("got (%s, %s, %v) want REFLECT / reflect_pending / false", phase, reason, stop)
	}
}

func TestSelectNextReplanWhenFlagged(t *testing.T) {
	inputs := cycleClearBase()
	inputs.ReplanNeeded = true
	phase, reason, stop := deliberation.SelectNext(deliberation.State{HopCount: 7}, inputs)
	if phase != deliberation.PhaseReplan || reason != deliberation.ReasonReplanNeeded || stop {
		t.Fatalf("got (%s, %s, %v) want REPLAN / replan_needed / false", phase, reason, stop)
	}
}

func TestSelectNextExploreWhenAlternatives(t *testing.T) {
	phase, reason, stop := deliberation.SelectNext(
		deliberation.State{HopCount: 2},
		deliberation.PolicyInputs{
			PlanExists:               true,
			PlanCritiqued:            false,
			OpenDecisionAlternatives: 2,
		},
	)
	if phase != deliberation.PhaseExplore || reason != deliberation.ReasonExploreAlternatives || stop {
		t.Fatalf("got (%s, %s, %v) want EXPLORE / explore_alternatives / false", phase, reason, stop)
	}
}

func TestSelectNextFullCycleOrdering(t *testing.T) {
	state := deliberation.State{HopCount: 1}
	base := cycleClearBase()

	cases := []struct {
		name       string
		inputs     deliberation.PolicyInputs
		wantPhase  deliberation.Phase
		wantReason deliberation.ReasonCode
	}{
		{
			name: "row 8 execute_pending only",
			inputs: func() deliberation.PolicyInputs {
				in := base
				in.ExecutePending = true
				return in
			}(),
			wantPhase:  deliberation.PhaseExecute,
			wantReason: deliberation.ReasonExecutePending,
		},
		{
			name: "row 9 test_pending only",
			inputs: func() deliberation.PolicyInputs {
				in := base
				in.TestPending = true
				return in
			}(),
			wantPhase:  deliberation.PhaseTest,
			wantReason: deliberation.ReasonTestPending,
		},
		{
			name: "row 10 verification_incomplete only",
			inputs: func() deliberation.PolicyInputs {
				in := base
				in.VerificationIncomplete = true
				return in
			}(),
			wantPhase:  deliberation.PhaseVerify,
			wantReason: deliberation.ReasonVerificationIncomplete,
		},
		{
			name: "row 11 evaluation_pending only",
			inputs: func() deliberation.PolicyInputs {
				in := base
				in.EvaluationPending = true
				return in
			}(),
			wantPhase:  deliberation.PhaseEvaluate,
			wantReason: deliberation.ReasonEvaluationPending,
		},
		{
			name: "row 12 reflect_pending only",
			inputs: func() deliberation.PolicyInputs {
				in := base
				in.ReflectPending = true
				return in
			}(),
			wantPhase:  deliberation.PhaseReflect,
			wantReason: deliberation.ReasonReflectPending,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			phase, reason, stop := deliberation.SelectNext(state, tc.inputs)
			if phase != tc.wantPhase || reason != tc.wantReason || stop {
				t.Fatalf("SelectNext = (%s, %s, %v) want (%s, %s, false)",
					phase, reason, stop, tc.wantPhase, tc.wantReason)
			}
		})
	}

	// Higher-priority row wins when multiple cycle flags set.
	multi := base
	multi.ExecutePending = true
	multi.TestPending = true
	multi.EvaluationPending = true
	phase, reason, _ := deliberation.SelectNext(state, multi)
	if phase != deliberation.PhaseExecute || reason != deliberation.ReasonExecutePending {
		t.Fatalf("first-match: got (%s, %s) want EXECUTE / execute_pending", phase, reason)
	}
}

func TestTransitionPayloadIncludesOptionalScores(t *testing.T) {
	_, payload := deliberation.ApplyTransition(
		deliberation.InitialState("t1", "g1"),
		deliberation.PolicyInputs{
			PlanExists:          true,
			PlanCritiqued:       true,
			PlanConfidence:      0.85,
			RequirementCoverage: 0.72,
		},
	)
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	pi := m["policy_inputs"].(map[string]any)
	if pi["plan_confidence"] != 0.85 {
		t.Fatalf("plan_confidence: %v want 0.85", pi["plan_confidence"])
	}
	if pi["requirement_coverage"] != 0.72 {
		t.Fatalf("requirement_coverage: %v want 0.72", pi["requirement_coverage"])
	}
}

func TestApplyTransitionStickyInvestigateConsumesHopBudget(t *testing.T) {
	state := deliberation.InitialState("t", "g")
	inputs := deliberation.PolicyInputs{BlockingUncertaintyCount: 1, PlanExists: true, PlanCritiqued: true}
	for i := 0; i < deliberation.HopBudget; i++ {
		var payload deliberation.TransitionPayload
		state, payload = deliberation.ApplyTransition(state, inputs)
		if state.CurrentPhase != deliberation.PhaseInvestigate || state.Stopped {
			t.Fatalf("i=%d state=%+v want INVESTIGATE until hop budget", i, state)
		}
		if state.HopCount != i+1 {
			t.Fatalf("hop_count=%d want %d", state.HopCount, i+1)
		}
		if payload.ToPhase != deliberation.PhaseInvestigate || payload.ReasonCode != deliberation.ReasonBlockingUncertainty {
			t.Fatalf("payload=%+v", payload)
		}
	}
	next, payload := deliberation.ApplyTransition(state, inputs)
	if next.CurrentPhase != deliberation.PhaseStop || !next.Stopped || payload.ReasonCode != deliberation.ReasonHopBudgetExceeded {
		t.Fatalf("after N: state=%+v payload=%+v", next, payload)
	}
	if next.HopCount != deliberation.HopBudget {
		t.Fatalf("hop_count after STOP: %d", next.HopCount)
	}
}

func TestApplyTransitionHopBudgetDoesNotIncrementPastN(t *testing.T) {
	state := deliberation.State{
		TaskID:       "task-1",
		GoalID:       "goal-1",
		CurrentPhase: deliberation.PhaseOrient,
		HopCount:     deliberation.HopBudget,
	}
	next, payload := deliberation.ApplyTransition(state, deliberation.PolicyInputs{PlanExists: true, PlanCritiqued: true})
	if next.CurrentPhase != deliberation.PhaseStop || !next.Stopped {
		t.Fatalf("state: %+v", next)
	}
	if next.HopCount != deliberation.HopBudget {
		t.Fatalf("hop_count after STOP: %d want %d", next.HopCount, deliberation.HopBudget)
	}
	if payload.HopCount != next.HopCount || payload.ToPhase != deliberation.PhaseStop || payload.ReasonCode != deliberation.ReasonHopBudgetExceeded {
		t.Fatalf("payload: %+v", payload)
	}
	if payload.FromPhase != deliberation.PhaseOrient {
		t.Fatalf("from_phase: %s", payload.FromPhase)
	}
}

func TestTransitionPayloadJSONRequiredFields(t *testing.T) {
	_, payload := deliberation.ApplyTransition(
		deliberation.InitialState("t1", "g1"),
		deliberation.PolicyInputs{PlanExists: true},
	)
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"task_id", "goal_id", "from_phase", "to_phase", "reason_code", "hop_count", "policy_inputs"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing %s", key)
		}
	}
	pi := m["policy_inputs"].(map[string]any)
	for _, key := range []string{
		"blocking_uncertainty_count", "plan_exists", "plan_critiqued",
		"verification_incomplete", "open_regression", "p19_saturated",
		"execute_pending", "test_pending", "evaluation_pending",
		"reflect_pending", "replan_needed", "open_decision_alternatives",
	} {
		if _, ok := pi[key]; !ok {
			t.Errorf("missing policy_inputs.%s", key)
		}
	}
	if m["to_phase"] != "CRITIQUE" || m["reason_code"] != "plan_uncritiqued" {
		t.Fatalf("payload: %s", raw)
	}
}
