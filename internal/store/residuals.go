package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Residual severity vocabulary (canonical; domain rejects unknown).
const (
	ResidualSeverityINFO     = "INFO"
	ResidualSeverityWARN     = "WARN"
	ResidualSeverityBlocking = "BLOCKING"
)

// Residual status vocabulary.
const (
	ResidualStatusOpen     = "OPEN"
	ResidualStatusAcked    = "ACKED"
	ResidualStatusResolved = "RESOLVED"
)

// ReviewResidual is a structured tracking hook on a review (not VerifiedFact).
type ReviewResidual struct {
	ID        string
	ReviewID  string
	Code      string
	Body      string
	Severity  string
	Status    string
	CreatedAt string
	UpdatedAt string
}

// InsertReviewResidual inserts a residual row. Empty ID allocates a UUID.
func (s *Store) InsertReviewResidual(r ReviewResidual) (ReviewResidual, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.ReviewID == "" {
		return ReviewResidual{}, fmt.Errorf("store: insert review residual: review_id required")
	}
	if r.Code == "" {
		return ReviewResidual{}, fmt.Errorf("store: insert review residual: code required")
	}
	if r.Severity == "" {
		r.Severity = ResidualSeverityINFO
	}
	if r.Status == "" {
		r.Status = ResidualStatusOpen
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO review_residuals(id, review_id, code, body, severity, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.ReviewID, r.Code, r.Body, r.Severity, r.Status, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return ReviewResidual{}, fmt.Errorf("store: insert review residual: %w", err)
	}
	return s.GetReviewResidual(r.ID)
}

// GetReviewResidual loads a residual by id.
func (s *Store) GetReviewResidual(id string) (ReviewResidual, error) {
	var r ReviewResidual
	err := s.db.QueryRow(`
		SELECT id, review_id, code, body, severity, status, created_at, updated_at
		FROM review_residuals WHERE id = ?
	`, id).Scan(
		&r.ID, &r.ReviewID, &r.Code, &r.Body, &r.Severity, &r.Status,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return ReviewResidual{}, fmt.Errorf("store: review residual %q: %w", id, err)
	}
	if err != nil {
		return ReviewResidual{}, fmt.Errorf("store: get review residual: %w", err)
	}
	return r, nil
}

// UpdateReviewResidualStatus sets status and bumps updated_at.
func (s *Store) UpdateReviewResidualStatus(id, status string) error {
	if id == "" {
		return fmt.Errorf("store: update review residual status: id required")
	}
	if status == "" {
		return fmt.Errorf("store: update review residual status: status required")
	}
	res, err := s.db.Exec(`
		UPDATE review_residuals SET status = ?, updated_at = ? WHERE id = ?
	`, status, nowRFC3339(), id)
	if err != nil {
		return fmt.Errorf("store: update review residual status: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: update review residual status: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: review residual %q: %w", id, sql.ErrNoRows)
	}
	return nil
}

// ListReviewResidualsByReviewID returns residuals for a review, ordered by created_at.
func (s *Store) ListReviewResidualsByReviewID(reviewID string) ([]ReviewResidual, error) {
	if reviewID == "" {
		return nil, fmt.Errorf("store: list review residuals: review_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, review_id, code, body, severity, status, created_at, updated_at
		FROM review_residuals
		WHERE review_id = ?
		ORDER BY created_at ASC, id ASC
	`, reviewID)
	if err != nil {
		return nil, fmt.Errorf("store: list review residuals by review: %w", err)
	}
	defer rows.Close()
	return scanReviewResiduals(rows)
}

// ListReviewResidualsByScopeID returns residuals on reviews linked via
// review_judges_scope to the given plan_scope id.
func (s *Store) ListReviewResidualsByScopeID(scopeID string) ([]ReviewResidual, error) {
	if scopeID == "" {
		return nil, fmt.Errorf("store: list review residuals by scope: scope_id required")
	}
	rows, err := s.db.Query(`
		SELECT rr.id, rr.review_id, rr.code, rr.body, rr.severity, rr.status, rr.created_at, rr.updated_at
		FROM review_residuals rr
		INNER JOIN entity_links el
			ON el.from_type = 'review'
			AND el.from_id = rr.review_id
			AND el.rel = 'review_judges_scope'
			AND el.to_type = 'plan_scope'
			AND el.to_id = ?
		ORDER BY rr.created_at ASC, rr.id ASC
	`, scopeID)
	if err != nil {
		return nil, fmt.Errorf("store: list review residuals by scope: %w", err)
	}
	defer rows.Close()
	return scanReviewResiduals(rows)
}

// CountOpenReviewResidualsByScopeID counts residuals with status=OPEN on reviews
// linked via review_judges_scope to the given plan_scope.
func (s *Store) CountOpenReviewResidualsByScopeID(scopeID string) (int, error) {
	if scopeID == "" {
		return 0, fmt.Errorf("store: count open residuals by scope: scope_id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(1)
		FROM review_residuals rr
		INNER JOIN entity_links el
			ON el.from_type = 'review'
			AND el.from_id = rr.review_id
			AND el.rel = 'review_judges_scope'
			AND el.to_type = 'plan_scope'
			AND el.to_id = ?
		WHERE rr.status = ?
	`, scopeID, ResidualStatusOpen).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count open residuals by scope: %w", err)
	}
	return n, nil
}

func scanReviewResiduals(rows *sql.Rows) ([]ReviewResidual, error) {
	var out []ReviewResidual
	for rows.Next() {
		var r ReviewResidual
		if err := rows.Scan(
			&r.ID, &r.ReviewID, &r.Code, &r.Body, &r.Severity, &r.Status,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan review residual: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
