package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Scope kinds for graph cartography (Phase 44). Distinct from plan_scopes.
const (
	ScopeKindFeature  = "feature"
	ScopeKindLayer    = "layer"
	ScopeKindBusiness = "business"
)

// Scope is a thin graph cartography record (feature / layer / business).
type Scope struct {
	ID        string
	Slug      string
	Title     string
	Kind      string
	CreatedAt string
	UpdatedAt string
}

// UpsertScope inserts or updates a scope by id. created_at is preserved on conflict.
// Empty ID allocates a UUID. Slug must be unique.
func (s *Store) UpsertScope(sc Scope) (Scope, error) {
	now := nowRFC3339()
	if sc.ID == "" {
		sc.ID = uuid.NewString()
	}
	if sc.Slug == "" {
		return Scope{}, fmt.Errorf("store: upsert scope: slug required")
	}
	if sc.Title == "" {
		return Scope{}, fmt.Errorf("store: upsert scope: title required")
	}
	if err := validateScopeKind(sc.Kind); err != nil {
		return Scope{}, err
	}
	if sc.CreatedAt == "" {
		sc.CreatedAt = now
	}
	sc.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO scopes(id, slug, title, kind, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			slug = excluded.slug,
			title = excluded.title,
			kind = excluded.kind,
			updated_at = excluded.updated_at
	`, sc.ID, sc.Slug, sc.Title, sc.Kind, sc.CreatedAt, sc.UpdatedAt)
	if err != nil {
		return Scope{}, fmt.Errorf("store: upsert scope: %w", err)
	}
	return s.GetScope(sc.ID)
}

func validateScopeKind(kind string) error {
	switch kind {
	case ScopeKindFeature, ScopeKindLayer, ScopeKindBusiness:
		return nil
	default:
		return fmt.Errorf("store: upsert scope: kind must be feature, layer, or business")
	}
}

// GetScope loads a scope by id.
func (s *Store) GetScope(id string) (Scope, error) {
	var sc Scope
	err := s.db.QueryRow(`
		SELECT id, slug, title, kind, created_at, updated_at
		FROM scopes WHERE id = ?
	`, id).Scan(&sc.ID, &sc.Slug, &sc.Title, &sc.Kind, &sc.CreatedAt, &sc.UpdatedAt)
	if err == sql.ErrNoRows {
		return Scope{}, fmt.Errorf("store: scope %q: %w", id, err)
	}
	if err != nil {
		return Scope{}, fmt.Errorf("store: get scope: %w", err)
	}
	return sc, nil
}

// GetScopeBySlug loads a scope by slug.
func (s *Store) GetScopeBySlug(slug string) (Scope, error) {
	var sc Scope
	err := s.db.QueryRow(`
		SELECT id, slug, title, kind, created_at, updated_at
		FROM scopes WHERE slug = ?
	`, slug).Scan(&sc.ID, &sc.Slug, &sc.Title, &sc.Kind, &sc.CreatedAt, &sc.UpdatedAt)
	if err == sql.ErrNoRows {
		return Scope{}, fmt.Errorf("store: scope slug %q: %w", slug, err)
	}
	if err != nil {
		return Scope{}, fmt.Errorf("store: get scope by slug: %w", err)
	}
	return sc, nil
}

// ListScopes returns all scopes ordered by slug.
func (s *Store) ListScopes() ([]Scope, error) {
	rows, err := s.db.Query(`
		SELECT id, slug, title, kind, created_at, updated_at
		FROM scopes
		ORDER BY slug ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list scopes: %w", err)
	}
	defer rows.Close()
	var out []Scope
	for rows.Next() {
		var sc Scope
		if err := rows.Scan(&sc.ID, &sc.Slug, &sc.Title, &sc.Kind, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan scope: %w", err)
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}
