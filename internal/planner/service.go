package planner

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Event type names for thin plan history.
const (
	EventCoarseCreated  = "plan.coarse_created"
	EventCurrentSet     = "plan.current_set"
	EventDeepPlanned    = "plan.deep_planned"
	EventDeepSuperseded = "plan.deep_superseded"
)

const (
	entityGoal  = "goal"
	entityScope = "plan_scope"
)

// Service orchestrates progressive coarse planning on a store.
type Service struct {
	store *store.Store
}

// New constructs a planner Service. st must be non-nil and already opened.
// Goal existence is checked via store.GetGoal (no domain import required).
func New(st *store.Store) *Service {
	if st == nil {
		panic("planner: New: store is nil")
	}
	return &Service{store: st}
}

// CreateCoarsePlan persists phases/scopes under an existing goal with stable ord.
// Does not deep-plan or create Tasks. Initializes goal_plan_state if missing.
func (s *Service) CreateCoarsePlan(ctx context.Context, in CoarsePlanInput) (CoarsePlan, error) {
	_ = ctx
	goalID := strings.TrimSpace(in.GoalID)
	if goalID == "" {
		return CoarsePlan{}, &ErrValidation{Msg: "goal_id is required"}
	}
	if _, err := s.store.GetGoal(goalID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CoarsePlan{}, fmt.Errorf("planner: goal %q: %w", goalID, ErrNotFound)
		}
		return CoarsePlan{}, fmt.Errorf("planner: get goal: %w", err)
	}
	if len(in.Phases) == 0 {
		return CoarsePlan{}, &ErrValidation{Msg: "at least one phase is required"}
	}

	out := CoarsePlan{GoalID: goalID}
	for pi, ph := range in.Phases {
		title := strings.TrimSpace(ph.Title)
		if title == "" {
			return CoarsePlan{}, &ErrValidation{Msg: "phase title is required"}
		}
		if len(ph.Scopes) == 0 {
			return CoarsePlan{}, &ErrValidation{Msg: "each phase requires at least one scope"}
		}
		stored, err := s.store.InsertPlanPhase(store.PlanPhase{
			GoalID: goalID,
			Title:  title,
			Body:   ph.Body,
			Ord:    pi,
		})
		if err != nil {
			return CoarsePlan{}, err
		}
		pv := PhaseView{
			ID: stored.ID, Title: stored.Title, Body: stored.Body,
			Ord: stored.Ord, Status: stored.Status,
			Scopes: []ScopeView{},
		}
		for si, sc := range ph.Scopes {
			stitle := strings.TrimSpace(sc.Title)
			if stitle == "" {
				return CoarsePlan{}, &ErrValidation{Msg: "scope title is required"}
			}
			ss, err := s.store.InsertPlanScope(store.PlanScope{
				PhaseID: stored.ID,
				Title:   stitle,
				Body:    sc.Body,
				Ord:     si,
			})
			if err != nil {
				return CoarsePlan{}, err
			}
			pv.Scopes = append(pv.Scopes, ScopeView{
				ID: ss.ID, PhaseID: ss.PhaseID, Title: ss.Title, Body: ss.Body,
				Ord: ss.Ord, Status: ss.Status, AutoReplanCount: ss.AutoReplanCount,
			})
		}
		out.Phases = append(out.Phases, pv)
	}

	if _, err := s.store.EnsureGoalPlanState(goalID); err != nil {
		return CoarsePlan{}, err
	}

	actor := in.Actor
	if actor == "" {
		actor = "planner"
	}
	payload, _ := json.Marshal(map[string]any{
		"goal_id":     goalID,
		"phase_count": len(out.Phases),
		"scope_count": countScopes(out),
		"actor":       actor,
	})
	_, _ = s.store.AppendEvent(store.Event{
		Type: EventCoarseCreated, EntityType: entityGoal, EntityID: goalID,
		PayloadJSON: string(payload),
	})
	return out, nil
}

func countScopes(cp CoarsePlan) int {
	n := 0
	for _, p := range cp.Phases {
		n += len(p.Scopes)
	}
	return n
}

// SetCurrentScope updates goal_plan_state after validating scope belongs to goal.
func (s *Service) SetCurrentScope(ctx context.Context, goalID, scopeID string) error {
	_ = ctx
	goalID = strings.TrimSpace(goalID)
	scopeID = strings.TrimSpace(scopeID)
	if goalID == "" || scopeID == "" {
		return &ErrValidation{Msg: "goal_id and scope_id are required"}
	}
	owner, err := s.store.GoalIDForScope(scopeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return err
	}
	if owner != goalID {
		return &ErrValidation{Msg: "scope does not belong to goal"}
	}
	cur := scopeID
	if _, err := s.store.UpsertGoalPlanState(store.GoalPlanState{
		GoalID: goalID, CurrentScopeID: &cur,
	}); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{
		"goal_id": goalID, "scope_id": scopeID, "actor": "planner",
	})
	_, _ = s.store.AppendEvent(store.Event{
		Type: EventCurrentSet, EntityType: entityGoal, EntityID: goalID,
		PayloadJSON: string(payload),
	})
	return nil
}

// GetCurrentScope returns the current scope ref, or ErrNotFound if unset.
func (s *Service) GetCurrentScope(ctx context.Context, goalID string) (ScopeRef, error) {
	_ = ctx
	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return ScopeRef{}, &ErrValidation{Msg: "goal_id is required"}
	}
	st, err := s.store.GetGoalPlanState(goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ScopeRef{}, fmt.Errorf("planner: current scope unset: %w", ErrNotFound)
		}
		return ScopeRef{}, err
	}
	if st.CurrentScopeID == nil || *st.CurrentScopeID == "" {
		return ScopeRef{}, fmt.Errorf("planner: current scope unset: %w", ErrNotFound)
	}
	return s.scopeRef(goalID, *st.CurrentScopeID)
}

// ListScopes returns ACTIVE scopes for a goal in progressive order.
func (s *Service) ListScopes(ctx context.Context, goalID string) ([]ScopeRef, error) {
	_ = ctx
	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return nil, &ErrValidation{Msg: "goal_id is required"}
	}
	scopes, err := s.store.ListPlanScopesByGoal(goalID)
	if err != nil {
		return nil, err
	}
	phases, err := s.store.ListPlanPhasesByGoal(goalID)
	if err != nil {
		return nil, err
	}
	phaseOrd := map[string]int{}
	for _, p := range phases {
		phaseOrd[p.ID] = p.Ord
	}
	out := make([]ScopeRef, 0, len(scopes))
	for _, sc := range scopes {
		out = append(out, ScopeRef{
			ID: sc.ID, PhaseID: sc.PhaseID, GoalID: goalID,
			Title: sc.Title, Body: sc.Body, Ord: sc.Ord,
			PhaseOrd: phaseOrd[sc.PhaseID],
		})
	}
	return out, nil
}

func (s *Service) scopeRef(goalID, scopeID string) (ScopeRef, error) {
	sc, err := s.store.GetPlanScope(scopeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ScopeRef{}, fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return ScopeRef{}, err
	}
	ph, err := s.store.GetPlanPhase(sc.PhaseID)
	if err != nil {
		return ScopeRef{}, err
	}
	if ph.GoalID != goalID {
		return ScopeRef{}, &ErrValidation{Msg: "scope does not belong to goal"}
	}
	return ScopeRef{
		ID: sc.ID, PhaseID: sc.PhaseID, GoalID: goalID,
		Title: sc.Title, Body: sc.Body, Ord: sc.Ord, PhaseOrd: ph.Ord,
	}, nil
}

// DeepPlan writes an ACTIVE deep plan for the current scope only (+ shallow lookahead).
func (s *Service) DeepPlan(ctx context.Context, in DeepPlanInput) (DeepPlanResult, error) {
	_ = ctx
	scopeID := strings.TrimSpace(in.ScopeID)
	if scopeID == "" {
		return DeepPlanResult{}, &ErrValidation{Msg: "scope_id is required"}
	}
	goalID, err := s.store.GoalIDForScope(scopeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeepPlanResult{}, fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return DeepPlanResult{}, err
	}
	st, err := s.store.GetGoalPlanState(goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeepPlanResult{}, fmt.Errorf("planner: current scope unset: %w", ErrNotFound)
		}
		return DeepPlanResult{}, err
	}
	cur := ""
	if st.CurrentScopeID != nil {
		cur = *st.CurrentScopeID
	}
	if cur == "" || cur != scopeID {
		return DeepPlanResult{}, &ErrNotCurrent{ScopeID: scopeID, CurrentID: cur}
	}

	lookaheadID, lookaheadSummary, err := s.resolveLookahead(goalID, scopeID, in.LookaheadSummary)
	if err != nil {
		return DeepPlanResult{}, err
	}

	doc := DeepPlanDocument{
		ScopeID:          scopeID,
		ExitCriteria:     nonNilStrings(in.ExitCriteria),
		Constraints:      nonNilStrings(in.Constraints),
		WorkItems:        nonNilWorkItems(in.WorkItems),
		LookaheadScopeID: lookaheadID,
		LookaheadSummary: lookaheadSummary,
	}
	if in.LookaheadSummary != "" && lookaheadID != "" {
		if err := s.store.UpdatePlanScopeBody(lookaheadID, in.LookaheadSummary); err != nil {
			return DeepPlanResult{}, err
		}
	}

	n, err := s.store.SupersedeActiveScopeDeepPlans(scopeID)
	if err != nil {
		return DeepPlanResult{}, err
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		return DeepPlanResult{}, fmt.Errorf("planner: marshal deep plan: %w", err)
	}
	rev, err := s.store.InsertScopeDeepPlan(store.ScopeDeepPlan{
		ScopeID: scopeID, ContentJSON: string(raw), Status: store.StatusActive,
	})
	if err != nil {
		return DeepPlanResult{}, err
	}

	actor := in.Actor
	if actor == "" {
		actor = "planner"
	}
	evType := EventDeepPlanned
	if n > 0 {
		evType = EventDeepSuperseded
	}
	payload, _ := json.Marshal(map[string]any{
		"goal_id": goalID, "scope_id": scopeID, "revision_id": rev.ID,
		"superseded": n, "lookahead_scope_id": lookaheadID, "actor": actor,
	})
	_, _ = s.store.AppendEvent(store.Event{
		Type: evType, EntityType: entityScope, EntityID: scopeID,
		PayloadJSON: string(payload),
	})

	return DeepPlanResult{
		RevisionID: rev.ID, Document: doc, SupersededCount: n,
		LookaheadScopeID: lookaheadID,
	}, nil
}

func (s *Service) resolveLookahead(goalID, currentScopeID, summaryOverride string) (id, summary string, err error) {
	scopes, err := s.store.ListPlanScopesByGoal(goalID)
	if err != nil {
		return "", "", err
	}
	idx := -1
	for i, sc := range scopes {
		if sc.ID == currentScopeID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", "", fmt.Errorf("planner: current scope missing from list: %w", ErrNotFound)
	}
	if idx+1 >= len(scopes) {
		return "", "", nil
	}
	next := scopes[idx+1]
	sum := summaryOverride
	if sum == "" {
		sum = next.Body
	}
	return next.ID, sum, nil
}

// SupersedeDeepPlan marks ACTIVE deep plan SUPERSEDED and writes a new ACTIVE revision.
// Does not require scope==current (S02 replan hook). Does not enforce churn budget.
func (s *Service) SupersedeDeepPlan(ctx context.Context, in SupersedeInput) (DeepPlanResult, error) {
	_ = ctx
	scopeID := strings.TrimSpace(in.ScopeID)
	if scopeID == "" {
		return DeepPlanResult{}, &ErrValidation{Msg: "scope_id is required"}
	}
	if _, err := s.store.GetPlanScope(scopeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeepPlanResult{}, fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return DeepPlanResult{}, err
	}

	doc := DeepPlanDocument{
		ScopeID:          scopeID,
		ExitCriteria:     nonNilStrings(in.ExitCriteria),
		Constraints:      nonNilStrings(in.Constraints),
		WorkItems:        nonNilWorkItems(in.WorkItems),
		LookaheadScopeID: in.LookaheadScopeID,
		LookaheadSummary: in.LookaheadSummary,
	}
	n, err := s.store.SupersedeActiveScopeDeepPlans(scopeID)
	if err != nil {
		return DeepPlanResult{}, err
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		return DeepPlanResult{}, fmt.Errorf("planner: marshal deep plan: %w", err)
	}
	rev, err := s.store.InsertScopeDeepPlan(store.ScopeDeepPlan{
		ScopeID: scopeID, ContentJSON: string(raw), Status: store.StatusActive,
	})
	if err != nil {
		return DeepPlanResult{}, err
	}

	actor := in.Actor
	if actor == "" {
		actor = "planner"
	}
	payload, _ := json.Marshal(map[string]any{
		"scope_id": scopeID, "revision_id": rev.ID, "superseded": n, "actor": actor,
	})
	_, _ = s.store.AppendEvent(store.Event{
		Type: EventDeepSuperseded, EntityType: entityScope, EntityID: scopeID,
		PayloadJSON: string(payload),
	})

	return DeepPlanResult{
		RevisionID: rev.ID, Document: doc, SupersededCount: n,
		LookaheadScopeID: in.LookaheadScopeID,
	}, nil
}

// GetPlan returns phases/scopes, current pointer, current ACTIVE deep plan, and lookahead.
func (s *Service) GetPlan(ctx context.Context, goalID string) (PlanView, error) {
	_ = ctx
	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return PlanView{}, &ErrValidation{Msg: "goal_id is required"}
	}
	if _, err := s.store.GetGoal(goalID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PlanView{}, fmt.Errorf("planner: goal %q: %w", goalID, ErrNotFound)
		}
		return PlanView{}, err
	}

	phases, err := s.store.ListPlanPhasesByGoal(goalID)
	if err != nil {
		return PlanView{}, err
	}
	view := PlanView{GoalID: goalID, Phases: []PhaseView{}}
	for _, ph := range phases {
		pv := PhaseView{
			ID: ph.ID, Title: ph.Title, Body: ph.Body, Ord: ph.Ord, Status: ph.Status,
			Scopes: []ScopeView{},
		}
		scopes, err := s.store.ListPlanScopesByPhase(ph.ID)
		if err != nil {
			return PlanView{}, err
		}
		for _, sc := range scopes {
			pv.Scopes = append(pv.Scopes, ScopeView{
				ID: sc.ID, PhaseID: sc.PhaseID, Title: sc.Title, Body: sc.Body,
				Ord: sc.Ord, Status: sc.Status, AutoReplanCount: sc.AutoReplanCount,
			})
		}
		view.Phases = append(view.Phases, pv)
	}

	st, err := s.store.GetGoalPlanState(goalID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return PlanView{}, err
	}
	if err == nil {
		view.CurrentScopeID = st.CurrentScopeID
	}

	if view.CurrentScopeID != nil && *view.CurrentScopeID != "" {
		cur := *view.CurrentScopeID
		if dp, err := s.store.GetActiveScopeDeepPlan(cur); err == nil {
			var doc DeepPlanDocument
			if json.Unmarshal([]byte(dp.ContentJSON), &doc) == nil {
				view.CurrentDeepPlan = &doc
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return PlanView{}, err
		}
		laID, laSum, err := s.resolveLookahead(goalID, cur, "")
		if err != nil {
			return PlanView{}, err
		}
		view.LookaheadScopeID = laID
		view.LookaheadSummary = laSum
		if laID != "" {
			if sc, err := s.store.GetPlanScope(laID); err == nil && sc.Body != "" {
				view.LookaheadSummary = sc.Body
			}
		}
	}
	return view, nil
}

func nonNilStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

func nonNilWorkItems(in []WorkItem) []WorkItem {
	if in == nil {
		return []WorkItem{}
	}
	return in
}
