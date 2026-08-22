package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// HarnessAgent is a catalog entry for a harness agent profile.
type HarnessAgent struct {
	ID                 string
	Slug               string
	Title              string
	Description        string
	SubagentType       string
	DeliberationPhases string // JSON array
	TaskKeywords       string // JSON array
	RecommendSubagent  bool
	RegistrySource     string
	RegistryVersion    string
	ExternalURL        string // empty when NULL
	CreatedAt          string
	UpdatedAt          string
}

// HarnessAgentRequirement links an agent to a required capability slug.
type HarnessAgentRequirement struct {
	ID                     string
	AgentID                string
	RequiredCapabilitySlug string
	CreatedAt              string
}

func scanHarnessAgent(row *sql.Row) (HarnessAgent, error) {
	var a HarnessAgent
	var recommend int
	var externalURL sql.NullString
	err := row.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Description, &a.SubagentType,
		&a.DeliberationPhases, &a.TaskKeywords, &recommend,
		&a.RegistrySource, &a.RegistryVersion, &externalURL,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return HarnessAgent{}, err
	}
	a.RecommendSubagent = recommend != 0
	if externalURL.Valid {
		a.ExternalURL = externalURL.String
	}
	return a, nil
}

const harnessAgentSelect = `
	SELECT id, slug, title, description, subagent_type,
		deliberation_phases, task_keywords, recommend_subagent,
		registry_source, registry_version, external_url,
		created_at, updated_at
	FROM harness_agents`

// UpsertHarnessAgent inserts or replaces a harness_agents row by id.
func (s *Store) UpsertHarnessAgent(a HarnessAgent) (HarnessAgent, error) {
	now := nowRFC3339()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	if strings.TrimSpace(a.Slug) == "" {
		return HarnessAgent{}, fmt.Errorf("store: upsert harness agent: slug required")
	}
	phases, err := requireJSONArray("deliberation_phases", a.DeliberationPhases)
	if err != nil {
		return HarnessAgent{}, err
	}
	keywords, err := requireJSONArray("task_keywords", a.TaskKeywords)
	if err != nil {
		return HarnessAgent{}, err
	}
	if a.RegistrySource == "" {
		a.RegistrySource = "bundled"
	}
	if a.CreatedAt == "" {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	a.DeliberationPhases = phases
	a.TaskKeywords = keywords

	recommend := 0
	if a.RecommendSubagent {
		recommend = 1
	}
	var externalURL any
	if strings.TrimSpace(a.ExternalURL) != "" {
		externalURL = a.ExternalURL
	}

	_, err = s.db.Exec(`
		INSERT INTO harness_agents(
			id, slug, title, description, subagent_type,
			deliberation_phases, task_keywords, recommend_subagent,
			registry_source, registry_version, external_url,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			slug = excluded.slug,
			title = excluded.title,
			description = excluded.description,
			subagent_type = excluded.subagent_type,
			deliberation_phases = excluded.deliberation_phases,
			task_keywords = excluded.task_keywords,
			recommend_subagent = excluded.recommend_subagent,
			registry_source = excluded.registry_source,
			registry_version = excluded.registry_version,
			external_url = excluded.external_url,
			updated_at = excluded.updated_at
	`, a.ID, a.Slug, a.Title, a.Description, a.SubagentType,
		a.DeliberationPhases, a.TaskKeywords, recommend,
		a.RegistrySource, a.RegistryVersion, externalURL,
		a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return HarnessAgent{}, fmt.Errorf("store: upsert harness agent: %w", err)
	}
	return s.GetHarnessAgent(a.ID)
}

// GetHarnessAgent loads one harness_agents row by id.
func (s *Store) GetHarnessAgent(id string) (HarnessAgent, error) {
	if id == "" {
		return HarnessAgent{}, fmt.Errorf("store: get harness agent: id required")
	}
	row := s.db.QueryRow(harnessAgentSelect+` WHERE id = ?`, id)
	a, err := scanHarnessAgent(row)
	if err == sql.ErrNoRows {
		return HarnessAgent{}, fmt.Errorf("store: harness agent %q: %w", id, err)
	}
	if err != nil {
		return HarnessAgent{}, fmt.Errorf("store: get harness agent: %w", err)
	}
	return a, nil
}

// GetHarnessAgentBySlug loads one harness_agents row by slug.
func (s *Store) GetHarnessAgentBySlug(slug string) (HarnessAgent, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return HarnessAgent{}, fmt.Errorf("store: get harness agent: slug required")
	}
	row := s.db.QueryRow(harnessAgentSelect+` WHERE slug = ?`, slug)
	a, err := scanHarnessAgent(row)
	if err == sql.ErrNoRows {
		return HarnessAgent{}, fmt.Errorf("store: harness agent slug %q: %w", slug, err)
	}
	if err != nil {
		return HarnessAgent{}, fmt.Errorf("store: get harness agent by slug: %w", err)
	}
	return a, nil
}

func scanHarnessAgents(rows *sql.Rows) ([]HarnessAgent, error) {
	var out []HarnessAgent
	for rows.Next() {
		var a HarnessAgent
		var recommend int
		var externalURL sql.NullString
		if err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Description, &a.SubagentType,
			&a.DeliberationPhases, &a.TaskKeywords, &recommend,
			&a.RegistrySource, &a.RegistryVersion, &externalURL,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan harness agent: %w", err)
		}
		a.RecommendSubagent = recommend != 0
		if externalURL.Valid {
			a.ExternalURL = externalURL.String
		}
		out = append(out, a)
	}
	if out == nil {
		out = []HarnessAgent{}
	}
	return out, rows.Err()
}

// ListHarnessAgents returns all harness agents ordered by slug.
func (s *Store) ListHarnessAgents() ([]HarnessAgent, error) {
	rows, err := s.db.Query(harnessAgentSelect + ` ORDER BY slug ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("store: list harness agents: %w", err)
	}
	defer rows.Close()
	return scanHarnessAgents(rows)
}

// ListAllHarnessAgents returns every harness agent (seed export).
func (s *Store) ListAllHarnessAgents() ([]HarnessAgent, error) {
	return s.ListHarnessAgents()
}

// ListHarnessAgentRequirements returns requirements for one agent ordered by slug.
func (s *Store) ListHarnessAgentRequirements(agentID string) ([]HarnessAgentRequirement, error) {
	if agentID == "" {
		return nil, fmt.Errorf("store: list harness agent requirements: agent_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, agent_id, required_capability_slug, created_at
		FROM harness_agent_requirements
		WHERE agent_id = ?
		ORDER BY required_capability_slug ASC, id ASC
	`, agentID)
	if err != nil {
		return nil, fmt.Errorf("store: list harness agent requirements: %w", err)
	}
	defer rows.Close()
	var out []HarnessAgentRequirement
	for rows.Next() {
		var r HarnessAgentRequirement
		if err := rows.Scan(&r.ID, &r.AgentID, &r.RequiredCapabilitySlug, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan harness agent requirement: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []HarnessAgentRequirement{}
	}
	return out, rows.Err()
}

// InsertHarnessAgentRequirement adds one requirement row.
func (s *Store) InsertHarnessAgentRequirement(r HarnessAgentRequirement) (HarnessAgentRequirement, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.AgentID == "" || strings.TrimSpace(r.RequiredCapabilitySlug) == "" {
		return HarnessAgentRequirement{}, fmt.Errorf("store: insert harness agent requirement: agent_id and slug required")
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	_, err := s.db.Exec(`
		INSERT INTO harness_agent_requirements(id, agent_id, required_capability_slug, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(agent_id, required_capability_slug) DO UPDATE SET
			id = excluded.id,
			created_at = excluded.created_at
	`, r.ID, r.AgentID, r.RequiredCapabilitySlug, r.CreatedAt)
	if err != nil {
		return HarnessAgentRequirement{}, fmt.Errorf("store: insert harness agent requirement: %w", err)
	}
	return r, nil
}

// DeleteHarnessAgentRequirementsForAgent removes all requirements for an agent.
func (s *Store) DeleteHarnessAgentRequirementsForAgent(agentID string) error {
	if agentID == "" {
		return fmt.Errorf("store: delete harness agent requirements: agent_id required")
	}
	_, err := s.db.Exec(`DELETE FROM harness_agent_requirements WHERE agent_id = ?`, agentID)
	if err != nil {
		return fmt.Errorf("store: delete harness agent requirements: %w", err)
	}
	return nil
}

// ListAllHarnessAgentRequirements returns every requirement row (seed export).
func (s *Store) ListAllHarnessAgentRequirements() ([]HarnessAgentRequirement, error) {
	rows, err := s.db.Query(`
		SELECT id, agent_id, required_capability_slug, created_at
		FROM harness_agent_requirements
		ORDER BY agent_id ASC, required_capability_slug ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all harness agent requirements: %w", err)
	}
	defer rows.Close()
	var out []HarnessAgentRequirement
	for rows.Next() {
		var r HarnessAgentRequirement
		if err := rows.Scan(&r.ID, &r.AgentID, &r.RequiredCapabilitySlug, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: scan harness agent requirement: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []HarnessAgentRequirement{}
	}
	return out, rows.Err()
}
