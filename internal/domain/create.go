package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// GoalInput creates a Goal.
type GoalInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

// DecisionInput creates a Decision.
type DecisionInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

// AssumptionInput creates an Assumption.
type AssumptionInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
	DecisionIDs    []string // optional assumption_supports_decision
	TaskIDs        []string // optional assumption_affects_task
}

// TaskInput creates a Task. GoalID optional.
type TaskInput struct {
	ID             string
	GoalID         *string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	WorkState      string
	LastVerifiedAt *string
}

// DiscoveryInput creates a Discovery.
type DiscoveryInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	Severity       string // optional; default INFO; INFO|PLAN_AFFECTING|BLOCKING
	LastVerifiedAt *string
}

// PlanChangeInput creates a PlanChange.
type PlanChangeInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

func applyProvenance(title, sourceType, status string) (string, string, error) {
	if strings.TrimSpace(title) == "" {
		return "", "", &ErrValidation{Msg: "title is required"}
	}
	if sourceType == "" {
		sourceType = DefaultSourceType
	}
	if status == "" {
		status = store.StatusActive
	}
	return sourceType, status, nil
}

func (s *Service) appendCreated(entityType, entityID, title string) error {
	payload, _ := json.Marshal(map[string]string{"title": title})
	_, err := s.store.AppendEvent(store.Event{
		Type:        EventEntityCreated,
		EntityType:  entityType,
		EntityID:    entityID,
		PayloadJSON: string(payload),
	})
	return err
}

// CreateGoal persists a goal and appends entity.created.
func (s *Service) CreateGoal(ctx context.Context, in GoalInput) (store.Goal, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Goal{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	g, err := s.store.UpsertGoal(store.Goal{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Goal{}, err
	}
	if err := s.appendCreated(EntityGoal, g.ID, g.Title); err != nil {
		return store.Goal{}, err
	}
	return g, nil
}

// CreateDecision persists a decision and appends entity.created.
func (s *Service) CreateDecision(ctx context.Context, in DecisionInput) (store.Decision, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Decision{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	d, err := s.store.UpsertDecision(store.Decision{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Decision{}, err
	}
	if err := s.appendCreated(EntityDecision, d.ID, d.Title); err != nil {
		return store.Decision{}, err
	}
	return d, nil
}

// CreateAssumption persists an assumption and appends entity.created.
func (s *Service) CreateAssumption(ctx context.Context, in AssumptionInput) (store.Assumption, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Assumption{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	for _, did := range in.DecisionIDs {
		did = strings.TrimSpace(did)
		if did == "" {
			return store.Assumption{}, &ErrValidation{Msg: "decision id is required"}
		}
		if _, err := s.store.GetDecision(did); err != nil {
			return store.Assumption{}, err
		}
	}
	for _, tid := range in.TaskIDs {
		tid = strings.TrimSpace(tid)
		if tid == "" {
			return store.Assumption{}, &ErrValidation{Msg: "task id is required"}
		}
		if _, err := s.store.GetTask(tid); err != nil {
			return store.Assumption{}, err
		}
	}
	a, err := s.store.UpsertAssumption(store.Assumption{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Assumption{}, err
	}
	for _, did := range in.DecisionIDs {
		if err := s.insertTypedLink(EntityAssumption, a.ID, RelAssumptionSupportsDecision, EntityDecision, strings.TrimSpace(did)); err != nil {
			return store.Assumption{}, err
		}
	}
	for _, tid := range in.TaskIDs {
		if err := s.insertTypedLink(EntityAssumption, a.ID, RelAssumptionAffectsTask, EntityTask, strings.TrimSpace(tid)); err != nil {
			return store.Assumption{}, err
		}
	}
	if err := s.appendCreated(EntityAssumption, a.ID, a.Title); err != nil {
		return store.Assumption{}, err
	}
	return a, nil
}

// CreateTask persists a task (default work_state PENDING) and appends entity.created.
func (s *Service) CreateTask(ctx context.Context, in TaskInput) (store.Task, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Task{}, err
	}
	ws := in.WorkState
	if ws == "" {
		ws = store.WorkStatePending
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	t, err := s.store.UpsertTask(store.Task{
		ID:             id,
		GoalID:         in.GoalID,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		WorkState:      ws,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Task{}, err
	}
	if err := s.appendCreated(EntityTask, t.ID, t.Title); err != nil {
		return store.Task{}, err
	}
	return t, nil
}

// CreateDiscovery persists a discovery and appends entity.created.
func (s *Service) CreateDiscovery(ctx context.Context, in DiscoveryInput) (store.Discovery, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Discovery{}, err
	}
	sev, err := NormalizeSeverity(in.Severity)
	if err != nil {
		return store.Discovery{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	d, err := s.store.UpsertDiscovery(store.Discovery{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		Severity:       sev,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Discovery{}, err
	}
	if err := s.appendCreated(EntityDiscovery, d.ID, d.Title); err != nil {
		return store.Discovery{}, err
	}
	return d, nil
}

// SetDiscoverySeverity updates a discovery's severity (validated enum).
func (s *Service) SetDiscoverySeverity(ctx context.Context, discoveryID, severity string) error {
	_ = ctx
	if strings.TrimSpace(discoveryID) == "" {
		return &ErrValidation{Msg: "discoveryID is required"}
	}
	sev, err := NormalizeSeverity(severity)
	if err != nil {
		return err
	}
	d, err := s.store.GetDiscovery(discoveryID)
	if err != nil {
		return err
	}
	d.Severity = sev
	_, err = s.store.UpsertDiscovery(d)
	return err
}

// CreatePlanChange persists a plan_change and appends entity.created.
func (s *Service) CreatePlanChange(ctx context.Context, in PlanChangeInput) (store.PlanChange, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.PlanChange{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	p, err := s.store.UpsertPlanChange(store.PlanChange{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.PlanChange{}, err
	}
	if err := s.appendCreated(EntityPlanChange, p.ID, p.Title); err != nil {
		return store.PlanChange{}, err
	}
	return p, nil
}

// GetGoal / GetTask / GetDecision / GetDiscovery / GetAssumption / GetPlanChange thin wrappers.

func (s *Service) GetGoal(ctx context.Context, id string) (store.Goal, error) {
	_ = ctx
	return s.store.GetGoal(id)
}

func (s *Service) GetTask(ctx context.Context, id string) (store.Task, error) {
	_ = ctx
	return s.store.GetTask(id)
}

func (s *Service) GetDecision(ctx context.Context, id string) (store.Decision, error) {
	_ = ctx
	return s.store.GetDecision(id)
}

func (s *Service) GetAssumption(ctx context.Context, id string) (store.Assumption, error) {
	_ = ctx
	return s.store.GetAssumption(id)
}

func (s *Service) GetDiscovery(ctx context.Context, id string) (store.Discovery, error) {
	_ = ctx
	return s.store.GetDiscovery(id)
}

func (s *Service) GetPlanChange(ctx context.Context, id string) (store.PlanChange, error) {
	_ = ctx
	return s.store.GetPlanChange(id)
}
