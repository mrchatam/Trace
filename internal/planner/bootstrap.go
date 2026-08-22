package planner

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// BootstrapResult reports progressive planner recovery from plan_changes.
type BootstrapResult struct {
	GoalID        string `json:"goal_id"`
	ScopeID       string `json:"scope_id,omitempty"`
	AlreadyExists bool   `json:"already_exists"`
	Note          string `json:"note,omitempty"`
}

// BootstrapFromPlanChanges creates minimal coarse + current + deep planner state
// from existing plan_changes so PlanExists becomes true. Idempotent when plan exists.
func (s *Service) BootstrapFromPlanChanges(ctx context.Context, goalID, actor string) (BootstrapResult, error) {
	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return BootstrapResult{}, &ErrValidation{Msg: "goal_id is required"}
	}
	if _, err := s.store.GetGoal(goalID); err != nil {
		return BootstrapResult{}, fmt.Errorf("planner: goal %q: %w", goalID, ErrNotFound)
	}

	view, err := s.GetPlan(ctx, goalID)
	if err != nil {
		return BootstrapResult{}, err
	}
	if view.CurrentScopeID != nil && *view.CurrentScopeID != "" && view.CurrentDeepPlan != nil {
		return BootstrapResult{
			GoalID:        goalID,
			ScopeID:       *view.CurrentScopeID,
			AlreadyExists: true,
			Note:          "progressive plan already exists; bootstrap skipped",
		}, nil
	}

	pc, err := s.pickPrimaryPlanChange(goalID)
	if err != nil {
		return BootstrapResult{}, err
	}

	phaseTitle := "Recovery"
	scopeTitle := "Bootstrap scope"
	exitText := "Bootstrap from plan-changes"
	if pc != nil {
		if t := strings.TrimSpace(pc.Title); t != "" {
			phaseTitle = truncateRunes(t, 80)
			scopeTitle = truncateRunes(t, 60)
		}
		if b := strings.TrimSpace(pc.Body); b != "" {
			exitText = truncateRunes(b, 200)
		}
	}

	if actor == "" {
		actor = "planner"
	}
	cp, err := s.CreateCoarsePlan(ctx, CoarsePlanInput{
		GoalID: goalID,
		Actor:  actor,
		Phases: []PhaseInput{{
			Title:  phaseTitle,
			Scopes: []ScopeInput{{Title: scopeTitle, Body: exitText}},
		}},
	})
	if err != nil {
		return BootstrapResult{}, err
	}
	if len(cp.Phases) == 0 || len(cp.Phases[0].Scopes) == 0 {
		return BootstrapResult{}, &ErrValidation{Msg: "bootstrap produced no scope"}
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := s.SetCurrentScope(ctx, goalID, scopeID); err != nil {
		return BootstrapResult{}, err
	}
	if _, err := s.DeepPlan(ctx, DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{exitText},
		WorkItems:    []WorkItem{{Title: "Refine bootstrap plan"}},
		Actor:        actor,
	}); err != nil {
		return BootstrapResult{}, err
	}

	note := "bootstrapped progressive plan from plan_changes"
	if pc != nil {
		note = fmt.Sprintf("bootstrapped from plan_change %q", pc.Title)
	}
	return BootstrapResult{GoalID: goalID, ScopeID: scopeID, Note: note}, nil
}

func (s *Service) pickPrimaryPlanChange(goalID string) (*store.PlanChange, error) {
	all, err := s.store.ListPlanChanges()
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, &ErrValidation{Msg: "no plan_changes to bootstrap from"}
	}

	linked := goalLinkedPlanChangeIDs(s.store, goalID)
	candidates := make([]store.PlanChange, 0, len(all))
	for _, pc := range all {
		if strings.EqualFold(pc.Status, "SUPERSEDED") {
			continue
		}
		if len(linked) > 0 {
			if _, ok := linked[pc.ID]; !ok {
				continue
			}
		}
		candidates = append(candidates, pc)
	}
	if len(candidates) == 0 {
		for _, pc := range all {
			if !strings.EqualFold(pc.Status, "SUPERSEDED") {
				candidates = append(candidates, pc)
			}
		}
	}
	if len(candidates) == 0 {
		return nil, &ErrValidation{Msg: "no usable plan_changes to bootstrap from"}
	}

	sortPlanChanges(candidates)
	pc := candidates[0]
	return &pc, nil
}

func truncateRunes(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max]) + "…"
}
