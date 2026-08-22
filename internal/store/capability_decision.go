package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Capability tool-decision statuses (canonical; domain normalizes uppercase).
const (
	ToolDecisionAutoAllowed = "AUTO_ALLOWED"
	ToolDecisionPending     = "PENDING"
	ToolDecisionAllowed     = "ALLOWED"
	ToolDecisionDenied      = "DENIED"
)

// CapabilityToolDecision is a durable allow/deny/auto-allow audit row.
type CapabilityToolDecision struct {
	ID        string
	Slug      string
	Decision  string
	Reason    string
	Actor     string
	CreatedAt string
	UpdatedAt string
}

// UpsertCapabilityToolDecision inserts or updates a decision by slug.
// Empty ID allocates a UUID; CreatedAt preserved on conflict.
func (s *Store) UpsertCapabilityToolDecision(d CapabilityToolDecision) (CapabilityToolDecision, error) {
	now := nowRFC3339()
	if strings.TrimSpace(d.Slug) == "" {
		return CapabilityToolDecision{}, fmt.Errorf("store: upsert capability tool decision: slug required")
	}
	d.Decision = strings.TrimSpace(d.Decision)
	if d.Decision == "" {
		return CapabilityToolDecision{}, fmt.Errorf("store: upsert capability tool decision: decision required")
	}
	switch d.Decision {
	case ToolDecisionAutoAllowed, ToolDecisionPending, ToolDecisionAllowed, ToolDecisionDenied:
	default:
		return CapabilityToolDecision{}, fmt.Errorf("store: upsert capability tool decision: invalid decision %q", d.Decision)
	}
	if d.ID == "" {
		// Prefer existing id when updating by slug.
		if existing, err := s.GetCapabilityToolDecisionBySlug(d.Slug); err == nil {
			d.ID = existing.ID
			d.CreatedAt = existing.CreatedAt
		} else if err != nil && !isNoRows(err) {
			return CapabilityToolDecision{}, err
		} else {
			d.ID = uuid.NewString()
		}
	}
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO capability_tool_decisions(id, slug, decision, reason, actor, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			decision = excluded.decision,
			reason = excluded.reason,
			actor = excluded.actor,
			updated_at = excluded.updated_at
	`, d.ID, d.Slug, d.Decision, d.Reason, d.Actor, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return CapabilityToolDecision{}, fmt.Errorf("store: upsert capability tool decision: %w", err)
	}
	return s.GetCapabilityToolDecisionBySlug(d.Slug)
}

// GetCapabilityToolDecisionBySlug loads a decision by unique slug.
func (s *Store) GetCapabilityToolDecisionBySlug(slug string) (CapabilityToolDecision, error) {
	if slug == "" {
		return CapabilityToolDecision{}, fmt.Errorf("store: get capability tool decision: slug required")
	}
	var d CapabilityToolDecision
	err := s.db.QueryRow(`
		SELECT id, slug, decision, reason, actor, created_at, updated_at
		FROM capability_tool_decisions WHERE slug = ?
	`, slug).Scan(
		&d.ID, &d.Slug, &d.Decision, &d.Reason, &d.Actor, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return CapabilityToolDecision{}, fmt.Errorf("store: capability tool decision slug %q: %w", slug, err)
	}
	if err != nil {
		return CapabilityToolDecision{}, fmt.Errorf("store: get capability tool decision: %w", err)
	}
	return d, nil
}

// ListCapabilityToolDecisions returns all decisions ordered by slug.
func (s *Store) ListCapabilityToolDecisions() ([]CapabilityToolDecision, error) {
	rows, err := s.db.Query(`
		SELECT id, slug, decision, reason, actor, created_at, updated_at
		FROM capability_tool_decisions
		ORDER BY slug ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list capability tool decisions: %w", err)
	}
	defer rows.Close()
	var out []CapabilityToolDecision
	for rows.Next() {
		var d CapabilityToolDecision
		if err := rows.Scan(
			&d.ID, &d.Slug, &d.Decision, &d.Reason, &d.Actor, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan capability tool decision: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
