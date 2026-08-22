package deliberation

// SelectNext is a pure first-match policy. Same inputs always yield the same phase.
// EXECUTE is selected when execute_pending is set and upstream gates are clear
// (blocking_uncertainty_count=0, no open regression).
// When Stopped, the persisted StopReason is returned (hop_budget_exceeded only if
// HopCount >= HopBudget or StopReason is empty).
func SelectNext(state State, inputs PolicyInputs) (Phase, ReasonCode, bool) {
	if state.HopCount >= HopBudget {
		return PhaseStop, ReasonHopBudgetExceeded, true
	}
	if state.Stopped {
		reason := ReasonCode(state.StopReason)
		if reason == "" {
			reason = ReasonHopBudgetExceeded
		}
		return PhaseStop, reason, true
	}
	if inputs.P19Saturated {
		return PhaseStop, ReasonP19Saturated, true
	}
	if inputs.BlockingUncertaintyCount > 0 {
		return PhaseInvestigate, ReasonBlockingUncertainty, false
	}
	if inputs.OpenRegression {
		return PhaseInvestigate, ReasonOpenRegression, false
	}
	if !inputs.PlanExists {
		return PhasePlan, ReasonPlanMissing, false
	}
	if inputs.OpenDecisionAlternatives > 0 && !inputs.PlanCritiqued {
		return PhaseExplore, ReasonExploreAlternatives, false
	}
	if inputs.PlanExists && !inputs.PlanCritiqued {
		return PhaseCritique, ReasonPlanUncritiqued, false
	}
	if inputs.ExecutePending && inputs.BlockingUncertaintyCount == 0 && !inputs.OpenRegression {
		return PhaseExecute, ReasonExecutePending, false
	}
	if inputs.TestPending {
		return PhaseTest, ReasonTestPending, false
	}
	if inputs.VerificationIncomplete {
		return PhaseVerify, ReasonVerificationIncomplete, false
	}
	if inputs.EvaluationPending {
		return PhaseEvaluate, ReasonEvaluationPending, false
	}
	if inputs.ReflectPending {
		return PhaseReflect, ReasonReflectPending, false
	}
	if inputs.ReplanNeeded {
		return PhaseReplan, ReasonReplanNeeded, false
	}
	return PhaseOrient, ReasonContinueOrient, false
}

// ApplyTransition computes the next state and event payload without I/O.
// Hop count increments for non-terminal hops; STOP from hop budget or an
// already-stopped row does not increment past HopBudget.
func ApplyTransition(state State, inputs PolicyInputs) (State, TransitionPayload) {
	from := state.CurrentPhase
	if from == "" {
		from = PhaseOrient
	}
	to, reason, stop := SelectNext(state, inputs)

	hop := state.HopCount
	if !state.Stopped && state.HopCount < HopBudget {
		hop = state.HopCount + 1
	}

	next := state
	if next.CurrentPhase == "" {
		next.CurrentPhase = PhaseOrient
	}
	next.LastPhase = from
	next.CurrentPhase = to
	next.HopCount = hop
	next.PlanCritiqued = inputs.PlanCritiqued
	next.Stopped = stop
	if stop {
		next.StopReason = string(reason)
	}

	payload := TransitionPayload{
		TaskID:       next.TaskID,
		GoalID:       next.GoalID,
		FromPhase:    from,
		ToPhase:      to,
		ReasonCode:   reason,
		HopCount:     hop,
		PolicyInputs: inputs,
	}
	return next, payload
}
