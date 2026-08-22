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

// BuildPolicyInputs assembles live S06 policy signals. applyWrites is used on the
// apply path to treat plan_changes as plan_critiqued before post-write store reads.
func BuildPolicyInputs(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	taskID, goalID string,
	applyWrites *ApplyWrites,
	p19Saturated bool,
) (deliberation.PolicyInputs, error) {
	if dom == nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: domain service is required")
	}
	if plan == nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: planner is required")
	}

	blocking, err := dom.CountBlockingUncertainties(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: blocking uncertainties: %w", err)
	}
	verifyDebt, err := dom.HasVerificationDebt(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: verification debt: %w", err)
	}
	openReg, err := dom.HasOpenRegression(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: open regression: %w", err)
	}

	planExists := false
	planView, err := plan.GetPlan(ctx, goalID)
	if err == nil {
		planExists = planView.CurrentScopeID != nil && *planView.CurrentScopeID != "" && planView.CurrentDeepPlan != nil
	}

	planCritiqued := false
	dState, err := dom.GetDeliberationState(ctx, taskID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: deliberation state: %w", err)
		}
	} else {
		planCritiqued = dState.PlanCritiqued
	}
	if applyWrites != nil && len(applyWrites.PlanChanges) > 0 {
		planCritiqued = true
	}

	implSignal, err := dom.HasImplementationSignal(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: implementation signal: %w", err)
	}
	testSince, err := dom.HasTestOutcomeSinceLatestChange(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: test since change: %w", err)
	}
	verifyGate, _, err := dom.CheckVerificationGate(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: verification gate: %w", err)
	}
	hasVerif, err := dom.HasVerificationOutcome(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: verification outcome: %w", err)
	}
	hasEval, err := dom.HasComputedEvaluation(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: computed evaluation: %w", err)
	}
	hasReflect, err := dom.HasReflectionSinceEvaluation(ctx, taskID)
	if err != nil {
		return deliberation.PolicyInputs{}, fmt.Errorf("loop policy: reflection since evaluation: %w", err)
	}

	return deliberation.PolicyInputs{
		BlockingUncertaintyCount: blocking,
		PlanExists:               planExists,
		PlanCritiqued:            planCritiqued,
		VerificationIncomplete:   verifyDebt,
		OpenRegression:           openReg,
		P19Saturated:             p19Saturated,
		ExecutePending:           planExists && planCritiqued && blocking == 0 && !openReg && !implSignal,
		TestPending:              implSignal && !testSince,
		EvaluationPending:        (verifyGate || hasVerif) && !hasEval,
		ReflectPending:           hasEval && !hasReflect,
	}, nil
}

func p19SaturatedFromLastStep(st *store.Store, seed ApplySeed) bool {
	step, ok := latestLoopStep(st, seed)
	if ok && step.MaxIterationsReached {
		return true
	}
	ds, err := st.GetDeliberationState(seed.TaskID)
	if err != nil {
		return false
	}
	return deliberation.SaturatedFromCounter(false, ds.ConsecutiveEmptyApplies)
}
