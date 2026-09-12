package store

import (
	"database/sql"
	"fmt"
	"strings"
)

// Review list budgets (Laws 6–7): agent/HTTP surfaces must page; seed/graph may use ListReviews.
const (
	DefaultReviewListLimit     = 50
	HTTPDefaultReviewListLimit = 100
	MaxReviewListLimit         = 500
)

// ReviewListFilter controls bounded review listing.
type ReviewListFilter struct {
	TaskID      string // optional: reviews linked via review_judges_task
	Limit       int    // 0 → DefaultReviewListLimit; capped at MaxReviewListLimit; -1 unbounded
	Cursor      string // opaque keyset cursor from a prior NextCursor
	IncludeBody bool   // list rows omit body unless true
}

// ReviewListResult is a page of reviews.
type ReviewListResult struct {
	Reviews    []Review
	NextCursor string
	Truncated  bool
}

// ListReviewsFiltered returns a bounded review page. Limit -1 means unbounded (internal).
func (s *Store) ListReviewsFiltered(f ReviewListFilter) (ReviewListResult, error) {
	limit := f.Limit
	unbounded := limit < 0
	if !unbounded {
		if limit <= 0 {
			limit = DefaultReviewListLimit
		}
		if limit > MaxReviewListLimit {
			limit = MaxReviewListLimit
		}
	}

	bodyCol := "'' AS body"
	if f.IncludeBody {
		bodyCol = "r.body"
	}

	q := fmt.Sprintf(`
		SELECT r.id, r.title, %s, r.source_type, r.confidence, r.status, r.result, r.created_at, r.updated_at, r.last_verified_at
		FROM reviews r
		WHERE 1=1
	`, bodyCol)
	var args []any
	if tid := strings.TrimSpace(f.TaskID); tid != "" {
		q += `
		AND EXISTS (
			SELECT 1 FROM entity_links el
			WHERE el.from_type = 'review' AND el.from_id = r.id
			  AND el.rel = 'review_judges_task' AND el.to_type = 'task' AND el.to_id = ?
		)`
		args = append(args, tid)
	}
	if curCreated, curID, ok := decodeTaskCursor(f.Cursor); ok {
		q += ` AND (r.created_at > ? OR (r.created_at = ? AND r.id > ?))`
		args = append(args, curCreated, curCreated, curID)
	} else if strings.TrimSpace(f.Cursor) != "" {
		return ReviewListResult{}, fmt.Errorf("store: list reviews: invalid cursor")
	}
	q += ` ORDER BY r.created_at ASC, r.id ASC`
	if !unbounded {
		q += ` LIMIT ?`
		args = append(args, limit+1)
	}

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return ReviewListResult{}, fmt.Errorf("store: list reviews filtered: %w", err)
	}
	defer rows.Close()
	reviews, err := scanReviews(rows)
	if err != nil {
		return ReviewListResult{}, err
	}
	out := ReviewListResult{Reviews: reviews}
	if !unbounded && len(reviews) > limit {
		out.Truncated = true
		out.Reviews = reviews[:limit]
		last := out.Reviews[len(out.Reviews)-1]
		out.NextCursor = encodeTaskCursor(last.CreatedAt, last.ID)
	}
	if out.Reviews == nil {
		out.Reviews = []Review{}
	}
	return out, nil
}

func scanReviews(rows *sql.Rows) ([]Review, error) {
	var out []Review
	for rows.Next() {
		var r Review
		var lastVerified sql.NullString
		if err := rows.Scan(
			&r.ID, &r.Title, &r.Body, &r.SourceType, &r.Confidence, &r.Status, &r.Result,
			&r.CreatedAt, &r.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan review: %w", err)
		}
		r.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, r)
	}
	return out, rows.Err()
}
