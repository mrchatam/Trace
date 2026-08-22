package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const (
	ExperimentStatusPlanned   = "planned"
	ExperimentStatusRunning   = "running"
	ExperimentStatusCompleted = "completed"
)

// Experiment is a thin §16 record — metadata + optional outcome link, no runner.
type Experiment struct {
	ID                string
	TaskID            string
	Label             string
	HypothesisSummary string
	Status            string
	OutcomeResultID   string
	CreatedAt         string
	UpdatedAt         string
}

// UpsertExperiment inserts or replaces an experiment by id. Empty ID allocates a UUID.
func (s *Store) UpsertExperiment(e Experiment) (Experiment, error) {
	now := nowRFC3339()
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.TaskID == "" {
		return Experiment{}, fmt.Errorf("store: upsert experiment: task_id required")
	}
	if e.Status == "" {
		e.Status = ExperimentStatusPlanned
	}
	if e.CreatedAt == "" {
		e.CreatedAt = now
	}
	e.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO experiments(
			id, task_id, label, hypothesis_summary, status, outcome_result_id,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			task_id = excluded.task_id,
			label = excluded.label,
			hypothesis_summary = excluded.hypothesis_summary,
			status = excluded.status,
			outcome_result_id = excluded.outcome_result_id,
			updated_at = excluded.updated_at
	`, e.ID, e.TaskID, e.Label, e.HypothesisSummary, e.Status, e.OutcomeResultID,
		e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return Experiment{}, fmt.Errorf("store: upsert experiment: %w", err)
	}
	return s.GetExperiment(e.ID)
}

// GetExperiment loads an experiment by id.
func (s *Store) GetExperiment(id string) (Experiment, error) {
	if id == "" {
		return Experiment{}, fmt.Errorf("store: get experiment: id required")
	}
	var e Experiment
	err := s.db.QueryRow(`
		SELECT id, task_id, label, hypothesis_summary, status, outcome_result_id,
			created_at, updated_at
		FROM experiments WHERE id = ?
	`, id).Scan(
		&e.ID, &e.TaskID, &e.Label, &e.HypothesisSummary, &e.Status, &e.OutcomeResultID,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Experiment{}, fmt.Errorf("store: experiment %q: %w", id, err)
	}
	if err != nil {
		return Experiment{}, fmt.Errorf("store: get experiment: %w", err)
	}
	return e, nil
}

// ListExperimentsByTaskID returns experiments for a task, oldest first.
func (s *Store) ListExperimentsByTaskID(taskID string) ([]Experiment, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list experiments: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, label, hypothesis_summary, status, outcome_result_id,
			created_at, updated_at
		FROM experiments
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list experiments: %w", err)
	}
	defer rows.Close()
	var out []Experiment
	for rows.Next() {
		var e Experiment
		if err := rows.Scan(
			&e.ID, &e.TaskID, &e.Label, &e.HypothesisSummary, &e.Status, &e.OutcomeResultID,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan experiment: %w", err)
		}
		out = append(out, e)
	}
	if out == nil {
		out = []Experiment{}
	}
	return out, rows.Err()
}
