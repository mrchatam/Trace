package store

import (
	"database/sql"
	"fmt"
)

// ListPromotionCandidateDiscoveries returns BLOCKING discoveries that lack a live
// discovery_mentions_task → tasks row, ordered by created_at, id.
// Callers should pass limit+1 when they need truncation honesty.
// limit <= 0 returns nil.
func (s *Store) ListPromotionCandidateDiscoveries(limit int) ([]Discovery, error) {
	if limit <= 0 {
		return nil, nil
	}
	const q = `
		SELECT d.id, d.title, d.body, d.source_type, d.confidence, d.status, d.severity,
		       d.created_at, d.updated_at, d.last_verified_at
		FROM discoveries d
		WHERE d.severity = ?
		  AND NOT EXISTS (
			SELECT 1
			FROM entity_links el
			INNER JOIN tasks t ON t.id = el.to_id
			WHERE el.from_type = 'discovery'
			  AND el.from_id = d.id
			  AND el.rel = 'discovery_mentions_task'
			  AND el.to_type = 'task'
		  )
		ORDER BY d.created_at ASC, d.id ASC
		LIMIT ?
	`
	rows, err := s.db.Query(q, SeverityBlocking, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list promotion candidate discoveries: %w", err)
	}
	defer rows.Close()
	var out []Discovery
	for rows.Next() {
		var d Discovery
		var lastVerified sql.NullString
		if err := rows.Scan(
			&d.ID, &d.Title, &d.Body, &d.SourceType, &d.Confidence, &d.Status, &d.Severity,
			&d.CreatedAt, &d.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan promotion candidate: %w", err)
		}
		d.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Discovery{}
	}
	return out, nil
}
