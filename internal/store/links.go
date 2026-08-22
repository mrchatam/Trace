package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// EntityLink is a typed causal/structural edge between two entities.
type EntityLink struct {
	ID         string
	FromType   string
	FromID     string
	Rel        string
	ToType     string
	ToID       string
	SourceType string
	Confidence float64
	CreatedAt  string
}

// InsertLinkOrIgnore inserts a link or no-ops on duplicate endpoints.
// Returns inserted=true when a new row was written.
func (s *Store) InsertLinkOrIgnore(l EntityLink) (inserted bool, link EntityLink, err error) {
	if l.FromType == "" || l.FromID == "" || l.Rel == "" || l.ToType == "" || l.ToID == "" {
		return false, EntityLink{}, fmt.Errorf("store: insert link: from_type, from_id, rel, to_type, to_id are required")
	}
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	if l.CreatedAt == "" {
		l.CreatedAt = nowRFC3339()
	}

	res, err := s.db.Exec(`
		INSERT INTO entity_links(id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(from_type, from_id, rel, to_type, to_id) DO NOTHING
	`, l.ID, l.FromType, l.FromID, l.Rel, l.ToType, l.ToID, l.SourceType, l.Confidence, l.CreatedAt)
	if err != nil {
		return false, EntityLink{}, fmt.Errorf("store: insert link or ignore: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, EntityLink{}, fmt.Errorf("store: insert link or ignore rows: %w", err)
	}
	if n > 0 {
		out, err := s.GetLink(l.ID)
		return true, out, err
	}
	existing, err := s.getLinkByEndpoints(l.FromType, l.FromID, l.Rel, l.ToType, l.ToID)
	return false, existing, err
}

func (s *Store) getLinkByEndpoints(fromType, fromID, rel, toType, toID string) (EntityLink, error) {
	var l EntityLink
	err := s.db.QueryRow(`
		SELECT id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at
		FROM entity_links
		WHERE from_type = ? AND from_id = ? AND rel = ? AND to_type = ? AND to_id = ?
	`, fromType, fromID, rel, toType, toID).Scan(
		&l.ID, &l.FromType, &l.FromID, &l.Rel, &l.ToType, &l.ToID,
		&l.SourceType, &l.Confidence, &l.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return EntityLink{}, fmt.Errorf("store: link endpoints: %w", err)
	}
	if err != nil {
		return EntityLink{}, fmt.Errorf("store: get link by endpoints: %w", err)
	}
	return l, nil
}

// InsertLink inserts a link. Empty ID allocates a UUID. Empty CreatedAt uses now.
// Duplicate UNIQUE(from_type, from_id, rel, to_type, to_id) returns an error.
func (s *Store) InsertLink(l EntityLink) (EntityLink, error) {
	if l.FromType == "" || l.FromID == "" || l.Rel == "" || l.ToType == "" || l.ToID == "" {
		return EntityLink{}, fmt.Errorf("store: insert link: from_type, from_id, rel, to_type, to_id are required")
	}
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	if l.CreatedAt == "" {
		l.CreatedAt = nowRFC3339()
	}

	_, err := s.db.Exec(`
		INSERT INTO entity_links(id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, l.ID, l.FromType, l.FromID, l.Rel, l.ToType, l.ToID, l.SourceType, l.Confidence, l.CreatedAt)
	if err != nil {
		return EntityLink{}, fmt.Errorf("store: insert link: %w", err)
	}
	return s.GetLink(l.ID)
}

// GetLink loads a link by id.
func (s *Store) GetLink(id string) (EntityLink, error) {
	var l EntityLink
	err := s.db.QueryRow(`
		SELECT id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at
		FROM entity_links WHERE id = ?
	`, id).Scan(
		&l.ID, &l.FromType, &l.FromID, &l.Rel, &l.ToType, &l.ToID,
		&l.SourceType, &l.Confidence, &l.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return EntityLink{}, fmt.Errorf("store: link %q: %w", id, err)
	}
	if err != nil {
		return EntityLink{}, fmt.Errorf("store: get link: %w", err)
	}
	return l, nil
}

// ListLinksFrom returns links originating from an entity, ordered by created_at.
func (s *Store) ListLinksFrom(fromType, fromID string) ([]EntityLink, error) {
	rows, err := s.db.Query(`
		SELECT id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at
		FROM entity_links
		WHERE from_type = ? AND from_id = ?
		ORDER BY created_at ASC, id ASC
	`, fromType, fromID)
	if err != nil {
		return nil, fmt.Errorf("store: list links from: %w", err)
	}
	defer rows.Close()
	return scanLinks(rows)
}

// ListLinksTo returns links targeting an entity, ordered by created_at.
func (s *Store) ListLinksTo(toType, toID string) ([]EntityLink, error) {
	rows, err := s.db.Query(`
		SELECT id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at
		FROM entity_links
		WHERE to_type = ? AND to_id = ?
		ORDER BY created_at ASC, id ASC
	`, toType, toID)
	if err != nil {
		return nil, fmt.Errorf("store: list links to: %w", err)
	}
	defer rows.Close()
	return scanLinks(rows)
}

// ListLinksByRel returns all entity_links with the given rel, ordered by created_at.
func (s *Store) ListLinksByRel(rel string) ([]EntityLink, error) {
	if rel == "" {
		return nil, fmt.Errorf("store: list links by rel: rel required")
	}
	rows, err := s.db.Query(`
		SELECT id, from_type, from_id, rel, to_type, to_id, source_type, confidence, created_at
		FROM entity_links
		WHERE rel = ?
		ORDER BY created_at ASC, id ASC
	`, rel)
	if err != nil {
		return nil, fmt.Errorf("store: list links by rel: %w", err)
	}
	defer rows.Close()
	return scanLinks(rows)
}

func scanLinks(rows *sql.Rows) ([]EntityLink, error) {
	var out []EntityLink
	for rows.Next() {
		var l EntityLink
		if err := rows.Scan(
			&l.ID, &l.FromType, &l.FromID, &l.Rel, &l.ToType, &l.ToID,
			&l.SourceType, &l.Confidence, &l.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan link: %w", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
