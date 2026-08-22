package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// ResidualInput creates a structured residual on a review.
type ResidualInput struct {
	Code     string
	Body     string
	Severity string // empty → INFO
	Status   string // empty → OPEN
}

// ResidualStatusOptions controls SetResidualStatus (actor+reason required).
type ResidualStatusOptions struct {
	Actor  string
	Reason string
}

// NormalizeResidualSeverity returns a valid residual severity.
// Empty defaults to INFO; unknown values fail closed.
func NormalizeResidualSeverity(severity string) (string, error) {
	s := strings.TrimSpace(severity)
	if s == "" {
		return ResidualSeverityINFO, nil
	}
	switch s {
	case ResidualSeverityINFO, ResidualSeverityWARN, ResidualSeverityBlocking:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "residual severity must be INFO, WARN, or BLOCKING"}
	}
}

// NormalizeResidualStatus returns a valid residual status.
// Empty defaults to OPEN; unknown values fail closed.
func NormalizeResidualStatus(status string) (string, error) {
	s := strings.TrimSpace(status)
	if s == "" {
		return ResidualStatusOpen, nil
	}
	switch s {
	case ResidualStatusOpen, ResidualStatusAcked, ResidualStatusResolved:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "residual status must be OPEN, ACKED, or RESOLVED"}
	}
}

// AddResidual attaches a structured residual to an existing review.
func (s *Service) AddResidual(ctx context.Context, reviewID string, in ResidualInput) (store.ReviewResidual, error) {
	_ = ctx
	if reviewID == "" {
		return store.ReviewResidual{}, &ErrValidation{Msg: "reviewID is required"}
	}
	code := strings.TrimSpace(in.Code)
	if code == "" {
		return store.ReviewResidual{}, &ErrValidation{Msg: "residual code is required"}
	}
	sev, err := NormalizeResidualSeverity(in.Severity)
	if err != nil {
		return store.ReviewResidual{}, err
	}
	st, err := NormalizeResidualStatus(in.Status)
	if err != nil {
		return store.ReviewResidual{}, err
	}
	if _, err := s.store.GetReview(reviewID); err != nil {
		return store.ReviewResidual{}, err
	}
	return s.store.InsertReviewResidual(store.ReviewResidual{
		ReviewID: reviewID,
		Code:     code,
		Body:     in.Body,
		Severity: sev,
		Status:   st,
	})
}

// SetResidualStatus updates residual status. Actor+Reason required.
func (s *Service) SetResidualStatus(ctx context.Context, residualID, status string, opts ResidualStatusOptions) error {
	_ = ctx
	if residualID == "" {
		return &ErrValidation{Msg: "residualID is required"}
	}
	if strings.TrimSpace(opts.Actor) == "" || strings.TrimSpace(opts.Reason) == "" {
		return &ErrValidation{Msg: "actor and reason are required"}
	}
	raw := strings.TrimSpace(status)
	if raw == "" {
		return &ErrValidation{Msg: "residual status must be OPEN, ACKED, or RESOLVED"}
	}
	st, err := NormalizeResidualStatus(raw)
	if err != nil {
		return err
	}
	r, err := s.store.GetReviewResidual(residualID)
	if err != nil {
		return err
	}
	if err := s.store.UpdateReviewResidualStatus(residualID, st); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{
		"status": st,
		"actor":  opts.Actor,
		"reason": opts.Reason,
		"code":   r.Code,
		"review": r.ReviewID,
	})
	_, err = s.store.AppendEvent(store.Event{
		Type:        "residual.status",
		EntityType:  EntityReview,
		EntityID:    r.ReviewID,
		PayloadJSON: string(payload),
	})
	return err
}

// ListResidualsByReview returns residuals for a review.
func (s *Service) ListResidualsByReview(ctx context.Context, reviewID string) ([]store.ReviewResidual, error) {
	_ = ctx
	if reviewID == "" {
		return nil, &ErrValidation{Msg: "reviewID is required"}
	}
	return s.store.ListReviewResidualsByReviewID(reviewID)
}

// ListResidualsByScope returns residuals on reviews linked via review_judges_scope.
func (s *Service) ListResidualsByScope(ctx context.Context, scopeID string) ([]store.ReviewResidual, error) {
	_ = ctx
	if scopeID == "" {
		return nil, &ErrValidation{Msg: "scopeID is required"}
	}
	return s.store.ListReviewResidualsByScopeID(scopeID)
}

// CountOpenResidualsByScope counts OPEN residuals on reviews judging the scope.
func (s *Service) CountOpenResidualsByScope(ctx context.Context, scopeID string) (int, error) {
	_ = ctx
	if scopeID == "" {
		return 0, &ErrValidation{Msg: "scopeID is required"}
	}
	return s.store.CountOpenReviewResidualsByScopeID(scopeID)
}
