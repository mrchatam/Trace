package domain

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// ReviewInput creates a Review. Result starts open ("").
type ReviewInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

// ReviewResultOptions controls SetReviewResult (actor+reason required).
type ReviewResultOptions struct {
	Actor  string
	Reason string
}

// CreateReview persists a review with empty result and appends entity.created.
func (s *Service) CreateReview(ctx context.Context, in ReviewInput) (store.Review, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Review{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	r, err := s.store.UpsertReview(store.Review{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		Result:         store.ReviewResultOpen,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Review{}, err
	}
	if err := s.appendCreated(EntityReview, r.ID, r.Title); err != nil {
		return store.Review{}, err
	}
	return r, nil
}

// SetReviewResult sets PASS|FAIL|UNCERTAIN on a review. Actor+Reason required.
// TransitionTask does not accept narrative PASS — only this API writes results.
func (s *Service) SetReviewResult(ctx context.Context, reviewID, result string, opts ReviewResultOptions) error {
	_ = ctx
	if reviewID == "" {
		return &ErrValidation{Msg: "reviewID is required"}
	}
	if strings.TrimSpace(opts.Actor) == "" || strings.TrimSpace(opts.Reason) == "" {
		return &ErrValidation{Msg: "actor and reason are required"}
	}
	switch result {
	case store.ReviewResultPass, store.ReviewResultFail, store.ReviewResultUncertain:
	default:
		return &ErrValidation{Msg: "result must be PASS, FAIL, or UNCERTAIN"}
	}
	r, err := s.store.GetReview(reviewID)
	if err != nil {
		return err
	}
	r.Result = result
	if _, err := s.store.UpsertReview(r); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{
		"result": result,
		"actor":  opts.Actor,
		"reason": opts.Reason,
	})
	_, err = s.store.AppendEvent(store.Event{
		Type:        EventReviewResult,
		EntityType:  EntityReview,
		EntityID:    reviewID,
		PayloadJSON: string(payload),
	})
	return err
}

// LinkReviewTask inserts entity_links rel=review_judges_task (from=review, to=task).
func (s *Service) LinkReviewTask(ctx context.Context, reviewID, taskID string, meta LinkMeta) error {
	_ = ctx
	if reviewID == "" || taskID == "" {
		return &ErrValidation{Msg: "reviewID and taskID are required"}
	}
	if _, err := s.store.GetReview(reviewID); err != nil {
		return err
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityReview,
		FromID:     reviewID,
		Rel:        RelReviewJudgesTask,
		ToType:     EntityTask,
		ToID:       taskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityReview, reviewID, RelReviewJudgesTask, EntityTask, taskID, meta)
}

// LinkReviewScope inserts entity_links rel=review_judges_scope (from=review, to=plan_scope).
// Recording only — does not mutate plan_scopes.status or task DONE gates.
func (s *Service) LinkReviewScope(ctx context.Context, reviewID, scopeID string, meta LinkMeta) error {
	_ = ctx
	if reviewID == "" || scopeID == "" {
		return &ErrValidation{Msg: "reviewID and scopeID are required"}
	}
	if _, err := s.store.GetReview(reviewID); err != nil {
		return err
	}
	if _, err := s.store.GetPlanScope(scopeID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityReview,
		FromID:     reviewID,
		Rel:        RelReviewJudgesScope,
		ToType:     EntityPlanScope,
		ToID:       scopeID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityReview, reviewID, RelReviewJudgesScope, EntityPlanScope, scopeID, meta)
}

// GetReview loads a review by id.
func (s *Service) GetReview(ctx context.Context, id string) (store.Review, error) {
	_ = ctx
	return s.store.GetReview(id)
}

// ListReviews returns all reviews ordered by created_at (G19 thin store wrap).
func (s *Service) ListReviews(ctx context.Context) ([]store.Review, error) {
	_ = ctx
	return s.store.ListReviews()
}

// ListReviewsByTaskID returns reviews linked via review_judges_task to taskID,
// ordered by review created_at (stable).
func (s *Service) ListReviewsByTaskID(ctx context.Context, taskID string) ([]store.Review, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "taskID is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return nil, err
	}
	links, err := s.store.ListLinksTo(EntityTask, taskID)
	if err != nil {
		return nil, err
	}
	out := make([]store.Review, 0)
	seen := map[string]bool{}
	for _, l := range links {
		if l.Rel != RelReviewJudgesTask || l.FromType != EntityReview {
			continue
		}
		if seen[l.FromID] {
			continue
		}
		r, err := s.store.GetReview(l.FromID)
		if err != nil {
			return nil, err
		}
		seen[l.FromID] = true
		out = append(out, r)
	}
	// Stable order by created_at, then id (ListLinksTo is by link created_at).
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].CreatedAt != out[j].CreatedAt {
			return out[i].CreatedAt < out[j].CreatedAt
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}
