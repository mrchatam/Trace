package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// UncertaintyInput creates an Uncertainty (title = question text).
type UncertaintyInput struct {
	ID             string
	Title          string
	Body           string
	Severity       string // empty → INFO; INFO|BLOCKING; PLAN_AFFECTING fail closed
	Status         string // empty → OPEN; create must be OPEN
	Kind           string // '' | risk | gap | unknown
	SourceType     string
	Confidence     float64
	TaskID         string // required when severity=BLOCKING
	GoalID         string // optional → uncertainty_affects_goal
	LastVerifiedAt *string
}

// HypothesisInput creates a Hypothesis (title = statement).
type HypothesisInput struct {
	ID             string
	Title          string
	Body           string
	Status         string // empty → OPEN
	SourceType     string
	Confidence     float64
	EvidenceIDs    []string // optional hypothesis_supported_by
	UncertaintyID  string   // optional hypothesis_addresses_uncertainty
	LastVerifiedAt *string
}

// InvalidateAssumptionInput is the cognitive invalidation path (no DELETE).
type InvalidateAssumptionInput struct {
	Status         string // required STALE | SUPERSEDED
	Reason         string // required non-empty
	SupersededBy   string // optional assumption id when SUPERSEDED
	EmitDiscovery  bool
	DiscoveryTitle string
	DiscoveryBody  string
	TaskIDs        []string // optional RelDiscoveryMentionsTask when discovery emitted
}

// ReconsiderationInput records an append-only decision reconsideration.
type ReconsiderationInput struct {
	Trigger      string
	Status       string // empty → FIRED
	Reason       string
	RelatedType  string
	RelatedID    string
	ReconsiderAt string // empty → now
}

func normalizeUncertaintySeverity(severity string) (string, error) {
	s := strings.TrimSpace(severity)
	if s == "" {
		return store.UncertaintySeverityINFO, nil
	}
	switch s {
	case store.UncertaintySeverityINFO, store.UncertaintySeverityBlocking:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "uncertainty severity must be INFO or BLOCKING"}
	}
}

func normalizeUncertaintyKind(kind string) (string, error) {
	k := strings.TrimSpace(kind)
	switch k {
	case store.UncertaintyKindNone, store.UncertaintyKindRisk, store.UncertaintyKindGap, store.UncertaintyKindUnknown:
		return k, nil
	default:
		return "", &ErrValidation{Msg: "uncertainty kind must be empty, risk, gap, or unknown"}
	}
}

func normalizeUncertaintyCreateStatus(status string) (string, error) {
	s := strings.TrimSpace(status)
	if s == "" {
		return store.UncertaintyStatusOpen, nil
	}
	if s != store.UncertaintyStatusOpen {
		return "", &ErrValidation{Msg: "uncertainty status on create must be OPEN"}
	}
	return s, nil
}

func normalizeHypothesisCreateStatus(status string) (string, error) {
	s := strings.TrimSpace(status)
	if s == "" {
		return store.HypothesisStatusOpen, nil
	}
	if s != store.HypothesisStatusOpen {
		return "", &ErrValidation{Msg: "hypothesis status on create must be OPEN"}
	}
	return s, nil
}

func normalizeReconsiderTrigger(trigger string) (string, error) {
	t := strings.TrimSpace(trigger)
	switch t {
	case store.ReconsiderTriggerContradictedEffect, store.ReconsiderTriggerNewEvidence, store.ReconsiderTriggerInvalidatedAssumption:
		return t, nil
	default:
		return "", &ErrValidation{Msg: "reconsideration trigger must be contradicted_effect, new_evidence, or invalidated_assumption"}
	}
}

func normalizeReconsiderStatus(status string) (string, error) {
	s := strings.TrimSpace(status)
	if s == "" {
		return store.ReconsiderStatusFired, nil
	}
	switch s {
	case store.ReconsiderStatusOpen, store.ReconsiderStatusFired:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "reconsideration status must be OPEN or FIRED"}
	}
}

func (s *Service) insertTypedLink(fromType, fromID, rel, toType, toID string) error {
	meta := LinkMeta{}.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   fromType,
		FromID:     fromID,
		Rel:        rel,
		ToType:     toType,
		ToID:       toID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(fromType, fromID, rel, toType, toID, meta)
}

func (s *Service) appendNamed(eventType, entityType, entityID string, payload map[string]string) error {
	raw, _ := json.Marshal(payload)
	_, err := s.store.AppendEvent(store.Event{
		Type:        eventType,
		EntityType:  entityType,
		EntityID:    entityID,
		PayloadJSON: string(raw),
	})
	return err
}

// CreateUncertainty persists an uncertainty. BLOCKING requires TaskID and inserts uncertainty_blocks_task.
func (s *Service) CreateUncertainty(ctx context.Context, in UncertaintyInput) (store.Uncertainty, error) {
	_ = ctx
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return store.Uncertainty{}, &ErrValidation{Msg: "title is required"}
	}
	sev, err := normalizeUncertaintySeverity(in.Severity)
	if err != nil {
		return store.Uncertainty{}, err
	}
	status, err := normalizeUncertaintyCreateStatus(in.Status)
	if err != nil {
		return store.Uncertainty{}, err
	}
	kind, err := normalizeUncertaintyKind(in.Kind)
	if err != nil {
		return store.Uncertainty{}, err
	}
	taskID := strings.TrimSpace(in.TaskID)
	goalID := strings.TrimSpace(in.GoalID)
	if sev == store.UncertaintySeverityBlocking {
		if taskID == "" {
			return store.Uncertainty{}, &ErrValidation{Msg: "BLOCKING uncertainty requires TaskID"}
		}
		if _, err := s.store.GetTask(taskID); err != nil {
			return store.Uncertainty{}, err
		}
	}
	if goalID != "" {
		if _, err := s.store.GetGoal(goalID); err != nil {
			return store.Uncertainty{}, err
		}
	}

	src := in.SourceType
	if src == "" {
		src = DefaultSourceType
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	u, err := s.store.UpsertUncertainty(store.Uncertainty{
		ID:             id,
		Title:          title,
		Body:           in.Body,
		Severity:       sev,
		Status:         status,
		Kind:           kind,
		Confidence:     in.Confidence,
		SourceType:     src,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Uncertainty{}, err
	}
	if sev == store.UncertaintySeverityBlocking {
		if err := s.insertTypedLink(EntityUncertainty, u.ID, RelUncertaintyBlocksTask, EntityTask, taskID); err != nil {
			return store.Uncertainty{}, err
		}
	}
	if goalID != "" {
		if err := s.insertTypedLink(EntityUncertainty, u.ID, RelUncertaintyAffectsGoal, EntityGoal, goalID); err != nil {
			return store.Uncertainty{}, err
		}
	}
	if err := s.appendCreated(EntityUncertainty, u.ID, u.Title); err != nil {
		return store.Uncertainty{}, err
	}
	return u, nil
}

// ResolveUncertainty transitions OPEN → RESOLVED. Resolution is required. Terminal afterwards.
func (s *Service) ResolveUncertainty(ctx context.Context, id, resolution string) (store.Uncertainty, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	resolution = strings.TrimSpace(resolution)
	if id == "" {
		return store.Uncertainty{}, &ErrValidation{Msg: "uncertainty id is required"}
	}
	if resolution == "" {
		return store.Uncertainty{}, &ErrValidation{Msg: "resolution is required"}
	}
	u, err := s.store.GetUncertainty(id)
	if err != nil {
		return store.Uncertainty{}, err
	}
	if u.Status != store.UncertaintyStatusOpen {
		return store.Uncertainty{}, &ErrValidation{Msg: "uncertainty resolve only from OPEN"}
	}
	u.Status = store.UncertaintyStatusResolved
	u.Resolution = resolution
	out, err := s.store.UpsertUncertainty(u)
	if err != nil {
		return store.Uncertainty{}, err
	}
	if err := s.appendNamed(EventUncertaintyResolved, EntityUncertainty, out.ID, map[string]string{
		"resolution": resolution,
		"status":     out.Status,
	}); err != nil {
		return store.Uncertainty{}, err
	}
	return out, nil
}

// SupersedeUncertainty transitions OPEN → SUPERSEDED. Reason is required. Terminal afterwards.
func (s *Service) SupersedeUncertainty(ctx context.Context, id, reason string) (store.Uncertainty, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	reason = strings.TrimSpace(reason)
	if id == "" {
		return store.Uncertainty{}, &ErrValidation{Msg: "uncertainty id is required"}
	}
	if reason == "" {
		return store.Uncertainty{}, &ErrValidation{Msg: "reason is required"}
	}
	u, err := s.store.GetUncertainty(id)
	if err != nil {
		return store.Uncertainty{}, err
	}
	if u.Status != store.UncertaintyStatusOpen {
		return store.Uncertainty{}, &ErrValidation{Msg: "uncertainty supersede only from OPEN"}
	}
	u.Status = store.UncertaintyStatusSuperseded
	u.Resolution = reason
	out, err := s.store.UpsertUncertainty(u)
	if err != nil {
		return store.Uncertainty{}, err
	}
	if err := s.appendNamed(EventUncertaintySuperseded, EntityUncertainty, out.ID, map[string]string{
		"reason": reason,
		"status": out.Status,
	}); err != nil {
		return store.Uncertainty{}, err
	}
	return out, nil
}

// CountBlockingUncertainties returns OPEN BLOCKING uncertainties linked to the task.
func (s *Service) CountBlockingUncertainties(ctx context.Context, taskID string) (int, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return 0, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.CountOpenBlockingUncertaintiesByTaskID(taskID)
}

// GetUncertainty loads an uncertainty by id.
func (s *Service) GetUncertainty(ctx context.Context, id string) (store.Uncertainty, error) {
	_ = ctx
	return s.store.GetUncertainty(id)
}

// CreateHypothesis persists a hypothesis and optional evidence/uncertainty links.
func (s *Service) CreateHypothesis(ctx context.Context, in HypothesisInput) (store.Hypothesis, error) {
	_ = ctx
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return store.Hypothesis{}, &ErrValidation{Msg: "title is required"}
	}
	status, err := normalizeHypothesisCreateStatus(in.Status)
	if err != nil {
		return store.Hypothesis{}, err
	}
	for _, eid := range in.EvidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return store.Hypothesis{}, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return store.Hypothesis{}, err
		}
	}
	uncID := strings.TrimSpace(in.UncertaintyID)
	if uncID != "" {
		if _, err := s.store.GetUncertainty(uncID); err != nil {
			return store.Hypothesis{}, err
		}
	}

	src := in.SourceType
	if src == "" {
		src = DefaultSourceType
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	h, err := s.store.UpsertHypothesis(store.Hypothesis{
		ID:             id,
		Title:          title,
		Body:           in.Body,
		Status:         status,
		Confidence:     in.Confidence,
		SourceType:     src,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Hypothesis{}, err
	}
	for _, eid := range in.EvidenceIDs {
		eid = strings.TrimSpace(eid)
		if err := s.insertTypedLink(EntityHypothesis, h.ID, RelHypothesisSupportedBy, EntityEvidence, eid); err != nil {
			return store.Hypothesis{}, err
		}
	}
	if uncID != "" {
		if err := s.insertTypedLink(EntityHypothesis, h.ID, RelHypothesisAddressesUncertainty, EntityUncertainty, uncID); err != nil {
			return store.Hypothesis{}, err
		}
	}
	if err := s.appendCreated(EntityHypothesis, h.ID, h.Title); err != nil {
		return store.Hypothesis{}, err
	}
	return h, nil
}

// ConfirmHypothesis transitions OPEN → CONFIRMED. Reason is required.
func (s *Service) ConfirmHypothesis(ctx context.Context, id, reason string) (store.Hypothesis, error) {
	return s.transitionHypothesis(ctx, id, store.HypothesisStatusConfirmed, reason, EventHypothesisConfirmed)
}

// RejectHypothesis transitions OPEN → REJECTED. Reason is required.
func (s *Service) RejectHypothesis(ctx context.Context, id, reason string) (store.Hypothesis, error) {
	return s.transitionHypothesis(ctx, id, store.HypothesisStatusRejected, reason, EventHypothesisRejected)
}

// SupersedeHypothesis transitions OPEN → SUPERSEDED. Reason is required.
func (s *Service) SupersedeHypothesis(ctx context.Context, id, reason string) (store.Hypothesis, error) {
	return s.transitionHypothesis(ctx, id, store.HypothesisStatusSuperseded, reason, EventHypothesisSuperseded)
}

func (s *Service) transitionHypothesis(ctx context.Context, id, toStatus, reason, eventType string) (store.Hypothesis, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	reason = strings.TrimSpace(reason)
	if id == "" {
		return store.Hypothesis{}, &ErrValidation{Msg: "hypothesis id is required"}
	}
	if reason == "" {
		return store.Hypothesis{}, &ErrValidation{Msg: "reason is required"}
	}
	h, err := s.store.GetHypothesis(id)
	if err != nil {
		return store.Hypothesis{}, err
	}
	if h.Status != store.HypothesisStatusOpen {
		return store.Hypothesis{}, &ErrValidation{Msg: "hypothesis transition only from OPEN"}
	}
	h.Status = toStatus
	out, err := s.store.UpsertHypothesis(h)
	if err != nil {
		return store.Hypothesis{}, err
	}
	if err := s.appendNamed(eventType, EntityHypothesis, out.ID, map[string]string{
		"reason": reason,
		"status": out.Status,
	}); err != nil {
		return store.Hypothesis{}, err
	}
	return out, nil
}

// GetHypothesis loads a hypothesis by id.
func (s *Service) GetHypothesis(ctx context.Context, id string) (store.Hypothesis, error) {
	_ = ctx
	return s.store.GetHypothesis(id)
}

// InvalidateAssumption sets STALE or SUPERSEDED on the existing row (no DELETE).
func (s *Service) InvalidateAssumption(ctx context.Context, assumptionID string, in InvalidateAssumptionInput) (store.Assumption, *store.Discovery, error) {
	assumptionID = strings.TrimSpace(assumptionID)
	if assumptionID == "" {
		return store.Assumption{}, nil, &ErrValidation{Msg: "assumptionID is required"}
	}
	reason := strings.TrimSpace(in.Reason)
	if reason == "" {
		return store.Assumption{}, nil, &ErrValidation{Msg: "reason is required"}
	}
	target := strings.ToUpper(strings.TrimSpace(in.Status))
	if target != store.StatusStale && target != store.StatusSuperseded {
		return store.Assumption{}, nil, &ErrValidation{Msg: "status must be STALE or SUPERSEDED"}
	}

	a, err := s.store.GetAssumption(assumptionID)
	if err != nil {
		return store.Assumption{}, nil, err
	}
	switch a.Status {
	case store.StatusActive:
		// ACTIVE → STALE or SUPERSEDED
	case store.StatusStale:
		if target != store.StatusSuperseded {
			return store.Assumption{}, nil, &ErrValidation{Msg: "STALE assumption may only move to SUPERSEDED"}
		}
	case store.StatusSuperseded:
		return store.Assumption{}, nil, &ErrValidation{Msg: "SUPERSEDED assumption is terminal"}
	default:
		return store.Assumption{}, nil, &ErrValidation{Msg: "assumption status cannot be invalidated"}
	}

	supersededBy := strings.TrimSpace(in.SupersededBy)
	if supersededBy != "" {
		if target != store.StatusSuperseded {
			return store.Assumption{}, nil, &ErrValidation{Msg: "SupersededBy is only valid when status is SUPERSEDED"}
		}
		if _, err := s.store.GetAssumption(supersededBy); err != nil {
			return store.Assumption{}, nil, err
		}
	}

	if in.EmitDiscovery && strings.TrimSpace(in.DiscoveryTitle) == "" {
		return store.Assumption{}, nil, &ErrValidation{Msg: "DiscoveryTitle is required when EmitDiscovery"}
	}

	a.Status = target
	out, err := s.store.UpsertAssumption(a)
	if err != nil {
		return store.Assumption{}, nil, err
	}
	payload := map[string]string{
		"reason": reason,
		"status": out.Status,
	}
	if supersededBy != "" {
		payload["superseded_by"] = supersededBy
	}
	if err := s.appendNamed(EventAssumptionInvalidated, EntityAssumption, out.ID, payload); err != nil {
		return store.Assumption{}, nil, err
	}

	links, err := s.store.ListLinksFrom(EntityAssumption, out.ID)
	if err != nil {
		return store.Assumption{}, nil, err
	}
	for _, l := range links {
		if l.Rel != RelAssumptionSupportsDecision || l.ToType != EntityDecision {
			continue
		}
		if _, err := s.AddImpactFinding(ctx, l.ToID, ImpactFindingInput{
			ImpactClass: ImpactClassCaution,
			Uncertainty: UncertaintyUNKNOWN,
			Kind:        FindingKindInvalidatedAssumption,
			Body:        reason,
			RelatedType: EntityAssumption,
			RelatedID:   out.ID,
		}); err != nil {
			return store.Assumption{}, nil, err
		}
		if _, err := s.RecordDecisionReconsideration(ctx, l.ToID, ReconsiderationInput{
			Trigger:     store.ReconsiderTriggerInvalidatedAssumption,
			Status:      store.ReconsiderStatusFired,
			Reason:      reason,
			RelatedType: EntityAssumption,
			RelatedID:   out.ID,
		}); err != nil {
			return store.Assumption{}, nil, err
		}
	}

	var disc *store.Discovery
	if in.EmitDiscovery {
		d, err := s.CreateDiscovery(ctx, DiscoveryInput{
			Title:    strings.TrimSpace(in.DiscoveryTitle),
			Body:     in.DiscoveryBody,
			Severity: SeverityPlanAffecting,
		})
		if err != nil {
			return store.Assumption{}, nil, err
		}
		if err := s.insertTypedLink(EntityDiscovery, d.ID, RelDiscoveryInvalidatesAssumption, EntityAssumption, out.ID); err != nil {
			return store.Assumption{}, nil, err
		}
		for _, tid := range in.TaskIDs {
			tid = strings.TrimSpace(tid)
			if tid == "" {
				continue
			}
			if err := s.LinkDiscoveryMentionsTask(ctx, d.ID, tid, LinkMeta{}); err != nil {
				return store.Assumption{}, nil, err
			}
		}
		disc = &d
	}
	return out, disc, nil
}

// RecordDecisionReconsideration appends a reconsideration child row. Does not delete Decision or alternatives.
func (s *Service) RecordDecisionReconsideration(ctx context.Context, decisionID string, in ReconsiderationInput) (store.DecisionReconsideration, error) {
	_ = ctx
	decisionID = strings.TrimSpace(decisionID)
	if decisionID == "" {
		return store.DecisionReconsideration{}, &ErrValidation{Msg: "decisionID is required"}
	}
	trigger, err := normalizeReconsiderTrigger(in.Trigger)
	if err != nil {
		return store.DecisionReconsideration{}, err
	}
	status, err := normalizeReconsiderStatus(in.Status)
	if err != nil {
		return store.DecisionReconsideration{}, err
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return store.DecisionReconsideration{}, err
	}
	row, err := s.store.InsertDecisionReconsideration(store.DecisionReconsideration{
		DecisionID:   decisionID,
		Trigger:      trigger,
		Status:       status,
		Reason:       in.Reason,
		RelatedType:  strings.TrimSpace(in.RelatedType),
		RelatedID:    strings.TrimSpace(in.RelatedID),
		ReconsiderAt: strings.TrimSpace(in.ReconsiderAt),
	})
	if err != nil {
		return store.DecisionReconsideration{}, err
	}
	if err := s.appendNamed(EventDecisionReconsider, EntityDecision, decisionID, map[string]string{
		"trigger":       row.Trigger,
		"status":        row.Status,
		"reason":        row.Reason,
		"related_type":  row.RelatedType,
		"related_id":    row.RelatedID,
		"reconsider_at": row.ReconsiderAt,
	}); err != nil {
		return store.DecisionReconsideration{}, err
	}
	return row, nil
}
