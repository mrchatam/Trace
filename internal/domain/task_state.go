package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// legalTransitions is the locked work_state graph.
var legalTransitions = map[string]map[string]bool{
	store.WorkStatePending: {
		store.WorkStateInProgress: true,
		store.WorkStateBlocked:    true,
		store.WorkStateSkipped:    true,
	},
	store.WorkStateInProgress: {
		store.WorkStateAwaitingReview: true,
		store.WorkStateBlocked:        true,
		store.WorkStateFailed:         true,
		store.WorkStateDone:           true,
		store.WorkStatePending:        true,
	},
	store.WorkStateAwaitingReview: {
		store.WorkStateDone:       true,
		store.WorkStateInProgress: true,
		store.WorkStateFailed:     true,
		store.WorkStateBlocked:    true,
	},
	store.WorkStateBlocked: {
		store.WorkStatePending:    true,
		store.WorkStateInProgress: true,
		store.WorkStateSkipped:    true,
		store.WorkStateFailed:     true,
	},
	store.WorkStateFailed: {
		store.WorkStatePending:    true,
		store.WorkStateInProgress: true,
		store.WorkStateSkipped:    true,
	},
	store.WorkStateDone: {
		store.WorkStateStale:   true,
		store.WorkStatePending: true,
	},
	store.WorkStateSkipped: {
		store.WorkStatePending: true,
	},
	store.WorkStateStale: {
		store.WorkStatePending: true,
	},
}

// TransitionTask moves a task along the locked work_state graph.
// Actor and Reason must be non-empty. Check order: legal edge → DF-24 missing
// caps → →DONE policy (hatch OR (no linked FAIL ∧ PASS ∧ AllowOperatorDone)).
// Leaving DONE invalidates linked PASS reviews (DF-18). EvidenceIDs alone never
// authorize DONE. Actor string is never authorization.
func (s *Service) TransitionTask(ctx context.Context, taskID, toWorkState string, opts TransitionOptions) error {
	if taskID == "" {
		return &ErrValidation{Msg: "taskID is required"}
	}
	if strings.TrimSpace(opts.Actor) == "" || strings.TrimSpace(opts.Reason) == "" {
		return &ErrValidation{Msg: "actor and reason are required"}
	}
	task, err := s.store.GetTask(taskID)
	if err != nil {
		return err
	}
	from := task.WorkState
	if from == "" {
		from = store.WorkStatePending
	}
	allowed := legalTransitions[from]
	if allowed == nil || !allowed[toWorkState] {
		return &ErrInvalidTransition{From: from, To: toWorkState, Reason: "edge not in legal graph"}
	}

	// DF-24: fail-closed when any required capability is missing (unless override).
	if !opts.AllowMissingCapabilities {
		missing, err := s.MissingCapabilities(ctx, taskID)
		if err != nil {
			return err
		}
		if len(missing) > 0 {
			return &ErrInvalidTransition{
				From: from,
				To:   toWorkState,
				Reason: "missing required capabilities; declare AVAILABLE or pass AllowMissingCapabilities " +
					"(--allow-missing-caps / allow_missing_caps)",
			}
		}
	}

	var passReviewID string
	if toWorkState == store.WorkStateDone {
		if !opts.AllowDoneWithoutReview {
			// DF-43: any linked review_judges_task with result=FAIL blocks DONE
			// even when a sibling PASS exists. UNCERTAIN/empty do not block.
			hasFail, err := s.hasLinkedFailReview(taskID)
			if err != nil {
				return err
			}
			if hasFail {
				return &ErrInvalidTransition{
					From: from,
					To:   toWorkState,
					Reason: "DONE blocked by linked Review with result=FAIL; clear FAIL " +
						"(e.g. SetReviewResult → UNCERTAIN) before PASS + AllowOperatorDone " +
						"(--as-operator), or use AllowDoneWithoutReview (--allow-done)",
				}
			}
			passReviewID, err = s.findPassReviewID(taskID)
			if err != nil {
				return err
			}
			if passReviewID == "" {
				return &ErrInvalidTransition{
					From: from,
					To:   toWorkState,
					Reason: "DONE requires linked Review with result=PASS and AllowOperatorDone " +
						"(--as-operator), or AllowDoneWithoutReview (--allow-done)",
				}
			}
			if !opts.AllowOperatorDone {
				return &ErrInvalidTransition{
					From: from,
					To:   toWorkState,
					Reason: "DONE requires AllowOperatorDone (--as-operator) with linked Review PASS, " +
						"or AllowDoneWithoutReview (--allow-done)",
				}
			}
		}
	}

	// DF-18: leaving DONE invalidates sticky PASS reviews before state change.
	if from == store.WorkStateDone {
		if err := s.invalidatePassReviewsOnReopen(ctx, taskID, toWorkState, opts.Actor); err != nil {
			return err
		}
	}

	evidenceIDs := opts.EvidenceIDs
	if evidenceIDs == nil {
		evidenceIDs = []string{}
	}
	payloadMap := map[string]any{
		"actor":        opts.Actor,
		"from":         from,
		"to":           toWorkState,
		"reason":       opts.Reason,
		"evidence_ids": evidenceIDs,
	}
	if passReviewID != "" {
		payloadMap["review_id"] = passReviewID
	}
	if opts.AllowDoneWithoutReview {
		payloadMap["allow_done_without_review"] = true
	}
	if opts.AllowOperatorDone {
		payloadMap["allow_operator_done"] = true
	}
	if opts.AllowMissingCapabilities {
		payloadMap["allow_missing_capabilities"] = true
	}
	payload, _ := json.Marshal(payloadMap)

	task.WorkState = toWorkState
	if _, err := s.store.UpsertTask(task); err != nil {
		return err
	}
	_, err = s.store.AppendEvent(store.Event{
		Type:        EventTaskTransition,
		EntityType:  EntityTask,
		EntityID:    taskID,
		PayloadJSON: string(payload),
	})
	return err
}

// invalidatePassReviewsOnReopen sets linked review_judges_task PASS → UNCERTAIN (DF-18).
func (s *Service) invalidatePassReviewsOnReopen(ctx context.Context, taskID, toWorkState, actor string) error {
	links, err := s.store.ListLinksTo(EntityTask, taskID)
	if err != nil {
		return err
	}
	reason := "invalidated on reopen (DONE→" + toWorkState + ")"
	act := strings.TrimSpace(actor)
	if act == "" {
		act = "system"
	}
	for _, l := range links {
		if l.Rel != RelReviewJudgesTask || l.FromType != EntityReview {
			continue
		}
		r, err := s.store.GetReview(l.FromID)
		if err != nil {
			return err
		}
		if r.Result != store.ReviewResultPass {
			continue
		}
		if err := s.SetReviewResult(ctx, r.ID, store.ReviewResultUncertain, ReviewResultOptions{
			Actor:  act,
			Reason: reason,
		}); err != nil {
			return err
		}
	}
	return nil
}

// findPassReviewID returns the id of a review linked via review_judges_task
// with result=PASS, or "" if none.
func (s *Service) findPassReviewID(taskID string) (string, error) {
	links, err := s.store.ListLinksTo(EntityTask, taskID)
	if err != nil {
		return "", err
	}
	for _, l := range links {
		if l.Rel != RelReviewJudgesTask || l.FromType != EntityReview {
			continue
		}
		r, err := s.store.GetReview(l.FromID)
		if err != nil {
			return "", err
		}
		if r.Result == store.ReviewResultPass {
			return r.ID, nil
		}
	}
	return "", nil
}

// hasLinkedFailReview reports whether any review linked via review_judges_task
// currently has result=FAIL (DF-43). Scope-only links are ignored.
func (s *Service) hasLinkedFailReview(taskID string) (bool, error) {
	links, err := s.store.ListLinksTo(EntityTask, taskID)
	if err != nil {
		return false, err
	}
	for _, l := range links {
		if l.Rel != RelReviewJudgesTask || l.FromType != EntityReview {
			continue
		}
		r, err := s.store.GetReview(l.FromID)
		if err != nil {
			return false, err
		}
		if r.Result == store.ReviewResultFail {
			return true, nil
		}
	}
	return false, nil
}

// MarkStale sets provenance status=STALE on a supported entity and appends an event.
func (s *Service) MarkStale(ctx context.Context, entityType, entityID, reason string) error {
	_ = ctx
	if entityType == "" || entityID == "" {
		return &ErrValidation{Msg: "entityType and entityID are required"}
	}
	if strings.TrimSpace(reason) == "" {
		return &ErrValidation{Msg: "reason is required"}
	}

	switch entityType {
	case EntityGoal:
		g, err := s.store.GetGoal(entityID)
		if err != nil {
			return err
		}
		g.Status = store.StatusStale
		if _, err := s.store.UpsertGoal(g); err != nil {
			return err
		}
	case EntityTask:
		t, err := s.store.GetTask(entityID)
		if err != nil {
			return err
		}
		t.Status = store.StatusStale
		if _, err := s.store.UpsertTask(t); err != nil {
			return err
		}
	case EntityDecision:
		d, err := s.store.GetDecision(entityID)
		if err != nil {
			return err
		}
		d.Status = store.StatusStale
		if _, err := s.store.UpsertDecision(d); err != nil {
			return err
		}
	case EntityAssumption:
		a, err := s.store.GetAssumption(entityID)
		if err != nil {
			return err
		}
		a.Status = store.StatusStale
		if _, err := s.store.UpsertAssumption(a); err != nil {
			return err
		}
	case EntityDiscovery:
		d, err := s.store.GetDiscovery(entityID)
		if err != nil {
			return err
		}
		d.Status = store.StatusStale
		if _, err := s.store.UpsertDiscovery(d); err != nil {
			return err
		}
	case EntityPlanChange:
		p, err := s.store.GetPlanChange(entityID)
		if err != nil {
			return err
		}
		p.Status = store.StatusStale
		if _, err := s.store.UpsertPlanChange(p); err != nil {
			return err
		}
	case EntityClaim:
		c, err := s.store.GetClaim(entityID)
		if err != nil {
			return err
		}
		c.Status = store.StatusStale
		if _, err := s.store.UpsertClaim(c); err != nil {
			return err
		}
	case EntityEvidence:
		e, err := s.store.GetEvidence(entityID)
		if err != nil {
			return err
		}
		e.Status = store.StatusStale
		if _, err := s.store.UpsertEvidence(e); err != nil {
			return err
		}
	default:
		return &ErrValidation{Msg: "unsupported entityType for MarkStale: " + entityType}
	}

	payload, _ := json.Marshal(map[string]string{
		"reason": reason,
		"status": store.StatusStale,
	})
	_, err := s.store.AppendEvent(store.Event{
		Type:        "entity.stale",
		EntityType:  entityType,
		EntityID:    entityID,
		PayloadJSON: string(payload),
	})
	return err
}
