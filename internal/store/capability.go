package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Capability kind vocabulary (canonical; domain rejects unknown/empty).
const (
	CapabilityKindSkill = "SKILL"
	CapabilityKindRule  = "RULE"
	CapabilityKindMCP   = "MCP"
	CapabilityKindTool  = "TOOL"
	CapabilityKindHook  = "HOOK"
	CapabilityKindAgent = "AGENT"
)

// Capability status vocabulary (canonical; empty on create → UNKNOWN at domain).
const (
	CapabilityStatusAvailable   = "AVAILABLE"
	CapabilityStatusUnavailable = "UNAVAILABLE"
	CapabilityStatusUnknown     = "UNKNOWN"
)

// Capability is a catalog entry (skill/rule/MCP/tool/hook).
type Capability struct {
	ID        string
	Kind      string
	Slug      string
	Title     string
	Status    string
	Body      string
	CreatedAt string
	UpdatedAt string
}

// TaskCapabilityRequirement links a task to a required capability.
type TaskCapabilityRequirement struct {
	ID           string
	TaskID       string
	CapabilityID string
	CreatedAt    string
}

// CapabilityListFilter optionally filters ListCapabilities.
type CapabilityListFilter struct {
	Kind   string // empty = any
	Status string // empty = any
}

// UpsertCapability inserts or replaces a capability by id. Empty ID allocates a UUID.
func (s *Store) UpsertCapability(c Capability) (Capability, error) {
	now := nowRFC3339()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	if strings.TrimSpace(c.Kind) == "" {
		return Capability{}, fmt.Errorf("store: upsert capability: kind required")
	}
	if strings.TrimSpace(c.Slug) == "" {
		return Capability{}, fmt.Errorf("store: upsert capability: slug required")
	}
	if c.Status == "" {
		c.Status = CapabilityStatusUnknown
	}
	if c.CreatedAt == "" {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO capabilities(id, kind, slug, title, status, body, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			kind = excluded.kind,
			slug = excluded.slug,
			title = excluded.title,
			status = excluded.status,
			body = excluded.body,
			updated_at = excluded.updated_at
	`, c.ID, c.Kind, c.Slug, c.Title, c.Status, c.Body, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return Capability{}, fmt.Errorf("store: upsert capability: %w", err)
	}
	return s.GetCapability(c.ID)
}

// GetCapability loads a capability by id.
func (s *Store) GetCapability(id string) (Capability, error) {
	if id == "" {
		return Capability{}, fmt.Errorf("store: get capability: id required")
	}
	var c Capability
	err := s.db.QueryRow(`
		SELECT id, kind, slug, title, status, body, created_at, updated_at
		FROM capabilities WHERE id = ?
	`, id).Scan(
		&c.ID, &c.Kind, &c.Slug, &c.Title, &c.Status, &c.Body, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Capability{}, fmt.Errorf("store: capability %q: %w", id, err)
	}
	if err != nil {
		return Capability{}, fmt.Errorf("store: get capability: %w", err)
	}
	return c, nil
}

// GetCapabilityBySlug loads a capability by unique slug.
func (s *Store) GetCapabilityBySlug(slug string) (Capability, error) {
	if slug == "" {
		return Capability{}, fmt.Errorf("store: get capability by slug: slug required")
	}
	var c Capability
	err := s.db.QueryRow(`
		SELECT id, kind, slug, title, status, body, created_at, updated_at
		FROM capabilities WHERE slug = ?
	`, slug).Scan(
		&c.ID, &c.Kind, &c.Slug, &c.Title, &c.Status, &c.Body, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Capability{}, fmt.Errorf("store: capability slug %q: %w", slug, err)
	}
	if err != nil {
		return Capability{}, fmt.Errorf("store: get capability by slug: %w", err)
	}
	return c, nil
}

// ListCapabilities returns capabilities matching the optional filter, ordered by slug.
func (s *Store) ListCapabilities(f CapabilityListFilter) ([]Capability, error) {
	q := `
		SELECT id, kind, slug, title, status, body, created_at, updated_at
		FROM capabilities
		WHERE 1=1
	`
	var args []any
	if f.Kind != "" {
		q += ` AND kind = ?`
		args = append(args, f.Kind)
	}
	if f.Status != "" {
		q += ` AND status = ?`
		args = append(args, f.Status)
	}
	q += ` ORDER BY slug ASC, id ASC`

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list capabilities: %w", err)
	}
	defer rows.Close()
	var out []Capability
	for rows.Next() {
		var c Capability
		if err := rows.Scan(
			&c.ID, &c.Kind, &c.Slug, &c.Title, &c.Status, &c.Body, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan capability: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// InsertTaskCapabilityRequirement inserts a requirement row. Empty ID allocates a UUID.
// UNIQUE(task_id, capability_id) conflicts return the existing row via Get.
func (s *Store) InsertTaskCapabilityRequirement(r TaskCapabilityRequirement) (TaskCapabilityRequirement, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.TaskID == "" {
		return TaskCapabilityRequirement{}, fmt.Errorf("store: insert task capability requirement: task_id required")
	}
	if r.CapabilityID == "" {
		return TaskCapabilityRequirement{}, fmt.Errorf("store: insert task capability requirement: capability_id required")
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}

	_, err := s.db.Exec(`
		INSERT INTO task_capability_requirements(id, task_id, capability_id, created_at)
		VALUES (?, ?, ?, ?)
	`, r.ID, r.TaskID, r.CapabilityID, r.CreatedAt)
	if err != nil {
		// UNIQUE conflict — return existing
		existing, getErr := s.GetTaskCapabilityRequirement(r.TaskID, r.CapabilityID)
		if getErr == nil {
			return existing, nil
		}
		return TaskCapabilityRequirement{}, fmt.Errorf("store: insert task capability requirement: %w", err)
	}
	return s.GetTaskCapabilityRequirement(r.TaskID, r.CapabilityID)
}

// GetTaskCapabilityRequirement loads a requirement by task+capability.
func (s *Store) GetTaskCapabilityRequirement(taskID, capabilityID string) (TaskCapabilityRequirement, error) {
	if taskID == "" || capabilityID == "" {
		return TaskCapabilityRequirement{}, fmt.Errorf("store: get task capability requirement: task_id and capability_id required")
	}
	var r TaskCapabilityRequirement
	err := s.db.QueryRow(`
		SELECT id, task_id, capability_id, created_at
		FROM task_capability_requirements
		WHERE task_id = ? AND capability_id = ?
	`, taskID, capabilityID).Scan(&r.ID, &r.TaskID, &r.CapabilityID, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return TaskCapabilityRequirement{}, fmt.Errorf("store: task capability requirement %s/%s: %w", taskID, capabilityID, err)
	}
	if err != nil {
		return TaskCapabilityRequirement{}, fmt.Errorf("store: get task capability requirement: %w", err)
	}
	return r, nil
}

// DeleteTaskCapabilityRequirement removes a requirement. Missing row is a no-op.
func (s *Store) DeleteTaskCapabilityRequirement(taskID, capabilityID string) error {
	if taskID == "" || capabilityID == "" {
		return fmt.Errorf("store: delete task capability requirement: task_id and capability_id required")
	}
	_, err := s.db.Exec(`
		DELETE FROM task_capability_requirements
		WHERE task_id = ? AND capability_id = ?
	`, taskID, capabilityID)
	if err != nil {
		return fmt.Errorf("store: delete task capability requirement: %w", err)
	}
	return nil
}

// ListTaskCapabilityRequirementsByTaskID returns requirement rows for a task, ordered by created_at.
func (s *Store) ListTaskCapabilityRequirementsByTaskID(taskID string) ([]TaskCapabilityRequirement, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list task capability requirements: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, capability_id, created_at
		FROM task_capability_requirements
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list task capability requirements: %w", err)
	}
	defer rows.Close()
	var out []TaskCapabilityRequirement
	for rows.Next() {
		var r TaskCapabilityRequirement
		if err := rows.Scan(&r.ID, &r.TaskID, &r.CapabilityID, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan task capability requirement: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ListCapabilitiesRequiredByTaskID joins requirements → capabilities, ordered by slug.
func (s *Store) ListCapabilitiesRequiredByTaskID(taskID string) ([]Capability, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list capabilities required by task: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT c.id, c.kind, c.slug, c.title, c.status, c.body, c.created_at, c.updated_at
		FROM task_capability_requirements r
		JOIN capabilities c ON c.id = r.capability_id
		WHERE r.task_id = ?
		ORDER BY c.slug ASC, c.id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list capabilities required by task: %w", err)
	}
	defer rows.Close()
	var out []Capability
	for rows.Next() {
		var c Capability
		if err := rows.Scan(
			&c.ID, &c.Kind, &c.Slug, &c.Title, &c.Status, &c.Body, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan required capability: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
