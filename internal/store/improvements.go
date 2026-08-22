package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Improvement is a first-class C18 row tied to a change and task.
type Improvement struct {
	ID              string
	ChangeID        string
	TaskID          string
	Dimension       string
	Summary         string
	EvidenceIDsJSON string
	SourceType      string
	Confidence      float64
	CreatedAt       string
	UpdatedAt       string
}

func requireEvidenceIDsJSONArray(raw string) (string, error) {
	if raw == "" {
		return emptyJSONArray, nil
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("store: upsert improvement: evidence_ids_json must be valid JSON")
	}
	arr, ok := v.([]any)
	if !ok {
		return "", fmt.Errorf("store: upsert improvement: evidence_ids_json must be a JSON array")
	}
	if len(arr) > 32 {
		return "", fmt.Errorf("store: upsert improvement: evidence_ids_json exceeds 32 items")
	}
	return raw, nil
}

// UpsertImprovement inserts or replaces an improvement by id. Empty ID allocates a UUID.
func (s *Store) UpsertImprovement(row Improvement) (Improvement, error) {
	now := nowRFC3339()
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	if row.ChangeID == "" {
		return Improvement{}, fmt.Errorf("store: upsert improvement: change_id required")
	}
	if row.TaskID == "" {
		return Improvement{}, fmt.Errorf("store: upsert improvement: task_id required")
	}
	evJSON, err := requireEvidenceIDsJSONArray(row.EvidenceIDsJSON)
	if err != nil {
		return Improvement{}, err
	}
	row.EvidenceIDsJSON = evJSON
	if row.CreatedAt == "" {
		row.CreatedAt = now
	}
	row.UpdatedAt = now

	_, err = s.db.Exec(`
		INSERT INTO improvements(
			id, change_id, task_id, dimension, summary, evidence_ids_json,
			source_type, confidence, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			change_id = excluded.change_id,
			task_id = excluded.task_id,
			dimension = excluded.dimension,
			summary = excluded.summary,
			evidence_ids_json = excluded.evidence_ids_json,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			updated_at = excluded.updated_at
	`, row.ID, row.ChangeID, row.TaskID, row.Dimension, row.Summary, row.EvidenceIDsJSON,
		row.SourceType, row.Confidence, row.CreatedAt, row.UpdatedAt)
	if err != nil {
		return Improvement{}, fmt.Errorf("store: upsert improvement: %w", err)
	}
	return s.GetImprovement(row.ID)
}

// GetImprovement loads an improvement by id.
func (s *Store) GetImprovement(id string) (Improvement, error) {
	if id == "" {
		return Improvement{}, fmt.Errorf("store: get improvement: id required")
	}
	var row Improvement
	err := s.db.QueryRow(`
		SELECT id, change_id, task_id, dimension, summary, evidence_ids_json,
			source_type, confidence, created_at, updated_at
		FROM improvements WHERE id = ?
	`, id).Scan(
		&row.ID, &row.ChangeID, &row.TaskID, &row.Dimension, &row.Summary, &row.EvidenceIDsJSON,
		&row.SourceType, &row.Confidence, &row.CreatedAt, &row.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Improvement{}, fmt.Errorf("store: improvement %q: %w", id, err)
	}
	if err != nil {
		return Improvement{}, fmt.Errorf("store: get improvement: %w", err)
	}
	return row, nil
}

func scanImprovements(rows *sql.Rows) ([]Improvement, error) {
	var out []Improvement
	for rows.Next() {
		var row Improvement
		if err := rows.Scan(
			&row.ID, &row.ChangeID, &row.TaskID, &row.Dimension, &row.Summary, &row.EvidenceIDsJSON,
			&row.SourceType, &row.Confidence, &row.CreatedAt, &row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan improvement: %w", err)
		}
		out = append(out, row)
	}
	if out == nil {
		out = []Improvement{}
	}
	return out, rows.Err()
}

const improvementSelect = `
		SELECT id, change_id, task_id, dimension, summary, evidence_ids_json,
			source_type, confidence, created_at, updated_at
		FROM improvements`

// ListAllImprovements returns every improvement row.
func (s *Store) ListAllImprovements() ([]Improvement, error) {
	rows, err := s.db.Query(improvementSelect + `
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all improvements: %w", err)
	}
	defer rows.Close()
	return scanImprovements(rows)
}

// ListImprovementsByChangeID returns improvements for a change, oldest first.
func (s *Store) ListImprovementsByChangeID(changeID string) ([]Improvement, error) {
	if changeID == "" {
		return nil, fmt.Errorf("store: list improvements by change: change_id required")
	}
	rows, err := s.db.Query(improvementSelect+`
		WHERE change_id = ?
		ORDER BY created_at ASC, id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("store: list improvements by change: %w", err)
	}
	defer rows.Close()
	return scanImprovements(rows)
}

// ListImprovementsByTaskID returns improvements for a task, oldest first.
func (s *Store) ListImprovementsByTaskID(taskID string) ([]Improvement, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list improvements by task: task_id required")
	}
	rows, err := s.db.Query(improvementSelect+`
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list improvements by task: %w", err)
	}
	defer rows.Close()
	return scanImprovements(rows)
}
