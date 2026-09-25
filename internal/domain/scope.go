package domain

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// ScopeInput creates or upserts a thin graph scope record.
type ScopeInput struct {
	ID    string
	Slug  string
	Title string
	Kind  string // feature | layer | business
}

// CreateScope upserts a thin scope (id/slug) and emits entity.created on first insert.
func (s *Service) CreateScope(ctx context.Context, in ScopeInput) (store.Scope, error) {
	_ = ctx
	slug := strings.TrimSpace(in.Slug)
	title := strings.TrimSpace(in.Title)
	kind := strings.TrimSpace(in.Kind)
	if slug == "" || title == "" {
		return store.Scope{}, &ErrValidation{Msg: "scope slug and title are required"}
	}
	switch kind {
	case store.ScopeKindFeature, store.ScopeKindLayer, store.ScopeKindBusiness:
	default:
		return store.Scope{}, &ErrValidation{Msg: "scope kind must be feature, layer, or business"}
	}

	id := strings.TrimSpace(in.ID)
	if id == "" {
		// Prefer update-by-slug when slug already exists.
		if existing, err := s.store.GetScopeBySlug(slug); err == nil {
			id = existing.ID
		} else if !errors.Is(err, sql.ErrNoRows) && !isNotFoundErr(err) {
			return store.Scope{}, err
		} else {
			id = uuid.NewString()
		}
	}

	existed, err := scopeExists(s.store, id)
	if err != nil {
		return store.Scope{}, err
	}

	sc, err := s.store.UpsertScope(store.Scope{
		ID: id, Slug: slug, Title: title, Kind: kind,
	})
	if err != nil {
		return store.Scope{}, err
	}
	if !existed {
		if err := s.appendCreated(EntityScope, sc.ID, sc.Title); err != nil {
			return store.Scope{}, err
		}
	}
	return sc, nil
}

func scopeExists(st *store.Store, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	_, err := st.GetScope(id)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) || isNotFoundErr(err) {
		return false, nil
	}
	return false, err
}

func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "sql: no rows")
}

// LinkScopeMember links any supported entity → scope (rel=scope_member).
func (s *Service) LinkScopeMember(ctx context.Context, fromID, scopeID string, meta LinkMeta) error {
	_ = ctx
	fromID = strings.TrimSpace(fromID)
	scopeID = strings.TrimSpace(scopeID)
	if fromID == "" || scopeID == "" {
		return &ErrValidation{Msg: "fromID and scopeID are required"}
	}
	fromType, err := s.resolveLinkableEntityType(fromID)
	if err != nil {
		return err
	}
	if _, err := s.store.GetScope(scopeID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   fromType,
		FromID:     fromID,
		Rel:        RelScopeMember,
		ToType:     EntityScope,
		ToID:       scopeID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(fromType, fromID, RelScopeMember, EntityScope, scopeID, meta)
}

// LinkAPIContract links task → task (rel=api_contract).
func (s *Service) LinkAPIContract(ctx context.Context, fromTaskID, toTaskID string, meta LinkMeta) error {
	_ = ctx
	fromTaskID = strings.TrimSpace(fromTaskID)
	toTaskID = strings.TrimSpace(toTaskID)
	if fromTaskID == "" || toTaskID == "" {
		return &ErrValidation{Msg: "from and to task ids are required"}
	}
	if _, err := s.store.GetTask(fromTaskID); err != nil {
		return err
	}
	if _, err := s.store.GetTask(toTaskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityTask,
		FromID:     fromTaskID,
		Rel:        RelAPIContract,
		ToType:     EntityTask,
		ToID:       toTaskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityTask, fromTaskID, RelAPIContract, EntityTask, toTaskID, meta)
}

// LinkImplements links task or decision → task (rel=implements). Distinct from change_implements_decision.
func (s *Service) LinkImplements(ctx context.Context, fromID, toTaskID string, meta LinkMeta) error {
	_ = ctx
	fromID = strings.TrimSpace(fromID)
	toTaskID = strings.TrimSpace(toTaskID)
	if fromID == "" || toTaskID == "" {
		return &ErrValidation{Msg: "fromID and toTaskID are required"}
	}
	fromType, err := s.resolveTaskOrDecision(fromID)
	if err != nil {
		return err
	}
	if _, err := s.store.GetTask(toTaskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   fromType,
		FromID:     fromID,
		Rel:        RelImplements,
		ToType:     EntityTask,
		ToID:       toTaskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(fromType, fromID, RelImplements, EntityTask, toTaskID, meta)
}

// LinkBlocks links task → task (rel=blocks). Distinct from uncertainty_blocks_task.
func (s *Service) LinkBlocks(ctx context.Context, fromTaskID, toTaskID string, meta LinkMeta) error {
	_ = ctx
	fromTaskID = strings.TrimSpace(fromTaskID)
	toTaskID = strings.TrimSpace(toTaskID)
	if fromTaskID == "" || toTaskID == "" {
		return &ErrValidation{Msg: "from and to task ids are required"}
	}
	if _, err := s.store.GetTask(fromTaskID); err != nil {
		return err
	}
	if _, err := s.store.GetTask(toTaskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityTask,
		FromID:     fromTaskID,
		Rel:        RelBlocks,
		ToType:     EntityTask,
		ToID:       toTaskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityTask, fromTaskID, RelBlocks, EntityTask, toTaskID, meta)
}

func (s *Service) resolveTaskOrDecision(id string) (string, error) {
	if _, err := s.store.GetTask(id); err == nil {
		return EntityTask, nil
	} else if !isNotFoundErr(err) {
		return "", err
	}
	if _, err := s.store.GetDecision(id); err == nil {
		return EntityDecision, nil
	} else if !isNotFoundErr(err) {
		return "", err
	}
	return "", &ErrValidation{Msg: "implements from must be task or decision"}
}

// resolveLinkableEntityType probes store for a known entity id (scope_member from side).
func (s *Service) resolveLinkableEntityType(id string) (string, error) {
	probes := []struct {
		typ string
		get func(string) error
	}{
		{EntityTask, func(id string) error { _, err := s.store.GetTask(id); return err }},
		{EntityGoal, func(id string) error { _, err := s.store.GetGoal(id); return err }},
		{EntityDecision, func(id string) error { _, err := s.store.GetDecision(id); return err }},
		{EntityAssumption, func(id string) error { _, err := s.store.GetAssumption(id); return err }},
		{EntityDiscovery, func(id string) error { _, err := s.store.GetDiscovery(id); return err }},
		{EntityPlanChange, func(id string) error { _, err := s.store.GetPlanChange(id); return err }},
		{EntityClaim, func(id string) error { _, err := s.store.GetClaim(id); return err }},
		{EntityEvidence, func(id string) error { _, err := s.store.GetEvidence(id); return err }},
		{EntityReview, func(id string) error { _, err := s.store.GetReview(id); return err }},
		{EntityChange, func(id string) error { _, err := s.store.GetChange(id); return err }},
		{EntityRegression, func(id string) error { _, err := s.store.GetRegression(id); return err }},
		{EntityScope, func(id string) error { _, err := s.store.GetScope(id); return err }},
	}
	for _, p := range probes {
		err := p.get(id)
		if err == nil {
			return p.typ, nil
		}
		if !isNotFoundErr(err) {
			return "", err
		}
	}
	return "", &ErrValidation{Msg: "entity not found for scope_member: " + id}
}
