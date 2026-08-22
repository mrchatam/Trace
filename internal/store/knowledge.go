package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// EngineeringKnowledge is a persisted project-specific knowledge row.
type EngineeringKnowledge struct {
	ID              string
	Title           string
	BodyJSON        string
	Topic           string
	EvidenceIDsJSON string
	Confidence      float64
	Status          string
	SourceType      string
	CreatedAt       string
	UpdatedAt       string
}

const (
	KnowledgeStatusActive     = "active"
	KnowledgeStatusSuperseded = "superseded"
)

func requireBodyJSONObject(raw string) (string, error) {
	if raw == "" {
		return "{}", nil
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("store: engineering knowledge: body_json must be valid JSON")
	}
	if _, ok := v.(map[string]any); !ok {
		return "", fmt.Errorf("store: engineering knowledge: body_json must be a JSON object")
	}
	if len(raw) > 8192 {
		return "", fmt.Errorf("store: engineering knowledge: body_json exceeds 8192 bytes")
	}
	return raw, nil
}

// UpsertEngineeringKnowledge inserts or replaces an engineering_knowledge row by id.
func (s *Store) UpsertEngineeringKnowledge(row EngineeringKnowledge) (EngineeringKnowledge, error) {
	now := nowRFC3339()
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	bodyJSON, err := requireBodyJSONObject(row.BodyJSON)
	if err != nil {
		return EngineeringKnowledge{}, err
	}
	evJSON, err := requireEvidenceIDsJSONArray(row.EvidenceIDsJSON)
	if err != nil {
		return EngineeringKnowledge{}, err
	}
	if row.Status == "" {
		row.Status = KnowledgeStatusActive
	}
	if row.CreatedAt == "" {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	row.BodyJSON = bodyJSON
	row.EvidenceIDsJSON = evJSON

	_, err = s.db.Exec(`
		INSERT INTO engineering_knowledge(
			id, title, body_json, topic, evidence_ids_json, confidence, status,
			source_type, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body_json = excluded.body_json,
			topic = excluded.topic,
			evidence_ids_json = excluded.evidence_ids_json,
			confidence = excluded.confidence,
			status = excluded.status,
			source_type = excluded.source_type,
			updated_at = excluded.updated_at
	`, row.ID, row.Title, row.BodyJSON, row.Topic, row.EvidenceIDsJSON, row.Confidence,
		row.Status, row.SourceType, row.CreatedAt, row.UpdatedAt)
	if err != nil {
		return EngineeringKnowledge{}, fmt.Errorf("store: upsert engineering knowledge: %w", err)
	}
	return s.GetEngineeringKnowledge(row.ID)
}

// GetEngineeringKnowledge loads one engineering_knowledge row by id.
func (s *Store) GetEngineeringKnowledge(id string) (EngineeringKnowledge, error) {
	if id == "" {
		return EngineeringKnowledge{}, fmt.Errorf("store: get engineering knowledge: id required")
	}
	var row EngineeringKnowledge
	err := s.db.QueryRow(`
		SELECT id, title, body_json, topic, evidence_ids_json, confidence, status,
			source_type, created_at, updated_at
		FROM engineering_knowledge WHERE id = ?
	`, id).Scan(
		&row.ID, &row.Title, &row.BodyJSON, &row.Topic, &row.EvidenceIDsJSON,
		&row.Confidence, &row.Status, &row.SourceType, &row.CreatedAt, &row.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return EngineeringKnowledge{}, fmt.Errorf("store: engineering knowledge %q: %w", id, err)
	}
	if err != nil {
		return EngineeringKnowledge{}, fmt.Errorf("store: get engineering knowledge: %w", err)
	}
	return row, nil
}

func scanEngineeringKnowledge(rows *sql.Rows) ([]EngineeringKnowledge, error) {
	var out []EngineeringKnowledge
	for rows.Next() {
		var row EngineeringKnowledge
		if err := rows.Scan(
			&row.ID, &row.Title, &row.BodyJSON, &row.Topic, &row.EvidenceIDsJSON,
			&row.Confidence, &row.Status, &row.SourceType, &row.CreatedAt, &row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan engineering knowledge: %w", err)
		}
		out = append(out, row)
	}
	if out == nil {
		out = []EngineeringKnowledge{}
	}
	return out, rows.Err()
}

const engineeringKnowledgeSelect = `
		SELECT id, title, body_json, topic, evidence_ids_json, confidence, status,
			source_type, created_at, updated_at
		FROM engineering_knowledge`

// ListEngineeringKnowledge returns knowledge rows with optional topic/status filters.
func (s *Store) ListEngineeringKnowledge(topic, status string, limit int) ([]EngineeringKnowledge, error) {
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	topic = strings.TrimSpace(topic)
	status = strings.TrimSpace(status)

	query := engineeringKnowledgeSelect + ` WHERE 1=1`
	args := []any{}
	if topic != "" {
		query += ` AND topic = ?`
		args = append(args, topic)
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY updated_at DESC, id ASC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list engineering knowledge: %w", err)
	}
	defer rows.Close()
	return scanEngineeringKnowledge(rows)
}

// ListAllEngineeringKnowledge returns every engineering_knowledge row.
func (s *Store) ListAllEngineeringKnowledge() ([]EngineeringKnowledge, error) {
	rows, err := s.db.Query(engineeringKnowledgeSelect + `
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all engineering knowledge: %w", err)
	}
	defer rows.Close()
	return scanEngineeringKnowledge(rows)
}

// ListAllChangePatterns returns every change_patterns row (for seed export).
func (s *Store) ListAllChangePatterns() ([]ChangePattern, error) {
	rows, err := s.db.Query(`
		SELECT change_kind, outcome_kind, count_positive, count_negative, last_seen
		FROM change_patterns
		ORDER BY change_kind ASC, outcome_kind ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all change patterns: %w", err)
	}
	defer rows.Close()
	var out []ChangePattern
	for rows.Next() {
		var p ChangePattern
		if err := rows.Scan(&p.ChangeKind, &p.OutcomeKind, &p.CountPositive, &p.CountNegative, &p.LastSeen); err != nil {
			return nil, fmt.Errorf("store: scan change pattern: %w", err)
		}
		out = append(out, p)
	}
	if out == nil {
		out = []ChangePattern{}
	}
	return out, rows.Err()
}
