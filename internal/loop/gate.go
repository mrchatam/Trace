package loop

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

type GateFor string

const (
	GateForOrient  GateFor = "orient"
	GateForEdit    GateFor = "edit"
	GateForExecute GateFor = "execute"
	GateForDone    GateFor = "done"
	GateForExport  GateFor = "export"
)

const (
	violationCodePrematureImplementation = "premature_implementation"
	violationCodeGateOrientFailed        = "gate_orient_failed"
	violationCodeGoalPlanGapAdvisory     = "goal_plan_gap_advisory"

	reasonCodeGoalPlanGapTerminalAdvisory = "goal_plan_gap_terminal_advisory"

	orientReasonTaskNotFound       = "task_not_found"
	orientReasonMissingGoalID      = "missing_goal_id"
	orientReasonMissingPlanContext = "missing_plan_context"

	executeReasonNotPending = "execute_not_pending"
)

type Violation struct {
	Code             string `json:"code"`
	For              string `json:"for"`
	Message          string `json:"message"`
	RecommendedPhase string `json:"recommended_phase,omitempty"`
	ReasonCode       string `json:"reason_code"`
}

// EvaluateGate is the single harness/product gate entrypoint.
func EvaluateGate(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (allowed bool, violations []Violation, err error) {
	switch gateFor {
	case GateForOrient:
		return evaluateOrient(ctx, plan, st, taskID, gateFor)
	case GateForEdit:
		return evaluateEdit(ctx, dom, plan, st, taskID, gateFor)
	case GateForExecute:
		return evaluateExecute(ctx, dom, plan, st, taskID, gateFor)
	case GateForDone, GateForExport:
		return evaluateDone(ctx, dom, plan, st, taskID, gateFor)
	default:
		return false, nil, fmt.Errorf("loop gate: unknown gate %q", gateFor)
	}
}

type gateContext struct {
	task    store.Task
	goalID  string
	inputs  deliberation.PolicyInputs
	dState  deliberation.State
	phase   deliberation.Phase
	reason  deliberation.ReasonCode
	stopped bool
}

func evaluateOrient(
	ctx context.Context,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (bool, []Violation, error) {
	if st == nil {
		return false, nil, fmt.Errorf("loop gate: store is required")
	}
	if plan == nil {
		return false, nil, fmt.Errorf("loop gate: planner is required")
	}

	task, err := st.GetTask(taskID)
	if err != nil {
		if isTaskNotFound(err) {
			return false, []Violation{orientViolation(gateFor, taskID, orientReasonTaskNotFound,
				fmt.Sprintf("orient failed: task %q not found", taskID))}, nil
		}
		return false, nil, fmt.Errorf("loop gate: load task: %w", err)
	}
	if task.GoalID == nil || *task.GoalID == "" {
		return false, []Violation{orientViolation(gateFor, taskID, orientReasonMissingGoalID,
			fmt.Sprintf("orient failed: task %q has no goal_id", taskID))}, nil
	}
	goalID := *task.GoalID

	planView, err := plan.GetPlan(ctx, goalID)
	if err != nil {
		return false, []Violation{orientViolation(gateFor, taskID, orientReasonMissingPlanContext,
			fmt.Sprintf("orient failed: missing plan context for goal %q", goalID))}, nil
	}
	if planView.CurrentScopeID == nil || *planView.CurrentScopeID == "" || planView.CurrentDeepPlan == nil {
		return false, []Violation{orientViolation(gateFor, taskID, orientReasonMissingPlanContext,
			fmt.Sprintf("orient failed: missing plan context for goal %q", goalID))}, nil
	}
	return true, nil, nil
}

func loadGateContext(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (gateContext, *Violation, error) {
	if dom == nil {
		return gateContext{}, nil, fmt.Errorf("loop gate: domain service is required")
	}
	if plan == nil {
		return gateContext{}, nil, fmt.Errorf("loop gate: planner is required")
	}
	if st == nil {
		return gateContext{}, nil, fmt.Errorf("loop gate: store is required")
	}

	task, err := st.GetTask(taskID)
	if err != nil {
		if isTaskNotFound(err) {
			v := orientViolation(gateFor, taskID, orientReasonTaskNotFound,
				fmt.Sprintf("gate failed: task %q not found", taskID))
			return gateContext{}, &v, nil
		}
		return gateContext{}, nil, fmt.Errorf("loop gate: load task: %w", err)
	}
	if task.GoalID == nil || *task.GoalID == "" {
		v := orientViolation(gateFor, taskID, orientReasonMissingGoalID,
			fmt.Sprintf("gate failed: task %q has no goal_id", taskID))
		return gateContext{}, &v, nil
	}
	goalID := *task.GoalID

	seed := ApplySeed{TaskID: taskID, GoalID: goalID}
	p19Sat := p19SaturatedFromLastStep(st, seed)
	inputs, err := BuildPolicyInputs(ctx, dom, plan, taskID, goalID, nil, p19Sat)
	if err != nil {
		return gateContext{}, nil, fmt.Errorf("loop gate: policy inputs: %w", err)
	}

	dState := loadDeliberationState(ctx, dom, taskID, goalID)
	phase, reason, stopped := deliberation.SelectNext(dState, inputs)

	return gateContext{
		task:    task,
		goalID:  goalID,
		inputs:  inputs,
		dState:  dState,
		phase:   phase,
		reason:  reason,
		stopped: stopped,
	}, nil, nil
}

func evaluateEdit(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (bool, []Violation, error) {
	gc, pre, err := loadGateContext(ctx, dom, plan, st, taskID, gateFor)
	if err != nil {
		return false, nil, err
	}
	if pre != nil {
		return false, []Violation{*pre}, nil
	}

	if gc.phase == deliberation.PhaseExecute && !gc.stopped {
		return true, nil, nil
	}
	if terminalPlanGapAdvisory(gc) {
		return true, []Violation{terminalPlanGapViolation(gateFor, gc.goalID)}, nil
	}
	return false, []Violation{prematureViolation(gateFor, taskID, gc.phase, gc.reason, gc.stopped,
		fmt.Sprintf("edit blocked: recommended phase %s (%s)", recommendedPhase(gc.phase, gc.stopped), gc.reason))}, nil
}

func evaluateExecute(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (bool, []Violation, error) {
	gc, pre, err := loadGateContext(ctx, dom, plan, st, taskID, gateFor)
	if err != nil {
		return false, nil, err
	}
	if pre != nil {
		return false, []Violation{*pre}, nil
	}

	allowed := gc.phase == deliberation.PhaseExecute &&
		gc.reason == deliberation.ReasonExecutePending &&
		gc.inputs.ExecutePending &&
		!gc.stopped
	if allowed {
		return true, nil, nil
	}

	reasonCode := string(gc.reason)
	if gc.phase == deliberation.PhaseExecute && !gc.inputs.ExecutePending {
		reasonCode = executeReasonNotPending
	}
	return false, []Violation{prematureViolationWithReason(gateFor, taskID, gc.phase, gc.stopped, reasonCode,
		fmt.Sprintf("execute blocked: recommended phase %s (%s)", recommendedPhase(gc.phase, gc.stopped), reasonCode))}, nil
}

func evaluateDone(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskID string,
	gateFor GateFor,
) (bool, []Violation, error) {
	gc, pre, err := loadGateContext(ctx, dom, plan, st, taskID, gateFor)
	if err != nil {
		return false, nil, err
	}
	if pre != nil {
		return false, []Violation{*pre}, nil
	}

	if gc.inputs.VerificationIncomplete {
		return false, []Violation{prematureViolationWithReason(gateFor, taskID, gc.phase, gc.stopped,
			string(deliberation.ReasonVerificationIncomplete),
			fmt.Sprintf("done blocked: verification incomplete (%s)", deliberation.ReasonVerificationIncomplete))}, nil
	}
	if gc.inputs.OpenRegression {
		return false, []Violation{prematureViolationWithReason(gateFor, taskID, gc.phase, gc.stopped,
			string(deliberation.ReasonOpenRegression),
			fmt.Sprintf("done blocked: open regression (%s)", deliberation.ReasonOpenRegression))}, nil
	}
	if gc.inputs.BlockingUncertaintyCount > 0 {
		return false, []Violation{prematureViolationWithReason(gateFor, taskID, deliberation.PhaseInvestigate, gc.stopped,
			string(deliberation.ReasonBlockingUncertainty),
			fmt.Sprintf("done blocked: recommended phase INVESTIGATE (%s)", deliberation.ReasonBlockingUncertainty))}, nil
	}

	if terminalPlanGapAdvisory(gc) {
		return true, []Violation{terminalPlanGapViolation(gateFor, gc.goalID)}, nil
	}

	if gc.phase == deliberation.PhaseOrient || gc.phase == deliberation.PhaseStop {
		return true, nil, nil
	}

	return false, []Violation{prematureViolation(gateFor, taskID, gc.phase, gc.reason, gc.stopped,
		fmt.Sprintf("done blocked: recommended phase %s (%s)", recommendedPhase(gc.phase, gc.stopped), gc.reason))}, nil
}

func orientViolation(gateFor GateFor, taskID, reasonCode, message string) Violation {
	return Violation{
		Code:       violationCodeGateOrientFailed,
		For:        string(gateFor),
		Message:    message,
		ReasonCode: reasonCode,
	}
}

func prematureViolation(
	gateFor GateFor,
	taskID string,
	phase deliberation.Phase,
	reason deliberation.ReasonCode,
	stopped bool,
	message string,
) Violation {
	return prematureViolationWithReason(gateFor, taskID, phase, stopped, string(reason), message)
}

func prematureViolationWithReason(
	gateFor GateFor,
	taskID string,
	phase deliberation.Phase,
	stopped bool,
	reasonCode string,
	message string,
) Violation {
	return Violation{
		Code:             violationCodePrematureImplementation,
		For:              string(gateFor),
		Message:          message,
		RecommendedPhase: recommendedPhase(phase, stopped),
		ReasonCode:       reasonCode,
	}
}

func recommendedPhase(phase deliberation.Phase, stopped bool) string {
	if stopped && phase == deliberation.PhaseStop {
		return string(deliberation.PhaseStop)
	}
	return string(phase)
}

func isTerminalWorkState(ws string) bool {
	return ws == store.WorkStateDone || ws == store.WorkStateSkipped
}

func terminalPlanGapAdvisory(gc gateContext) bool {
	return isTerminalWorkState(gc.task.WorkState) &&
		!gc.inputs.PlanExists &&
		gc.reason == deliberation.ReasonPlanMissing
}

func terminalPlanGapViolation(gateFor GateFor, goalID string) Violation {
	return Violation{
		Code:       violationCodeGoalPlanGapAdvisory,
		For:        string(gateFor),
		ReasonCode: reasonCodeGoalPlanGapTerminalAdvisory,
		Message: fmt.Sprintf(
			"goal %s lacks progressive plan (work already terminal); run trace plan bootstrap --goal %s or MCP trace_plan action=bootstrap",
			goalID, goalID,
		),
	}
}

func isTaskNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
