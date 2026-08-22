package store

import (
	"database/sql"
	"fmt"
)

// ListGoals returns all goals ordered by created_at, then id.
func (s *Store) ListGoals() ([]Goal, error) {
	rows, err := s.db.Query(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM goals
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list goals: %w", err)
	}
	defer rows.Close()
	var out []Goal
	for rows.Next() {
		var g Goal
		var lastVerified sql.NullString
		if err := rows.Scan(
			&g.ID, &g.Title, &g.Body, &g.SourceType, &g.Confidence, &g.Status,
			&g.CreatedAt, &g.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan goal: %w", err)
		}
		g.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ListTasks returns all tasks ordered by created_at, then id.
func (s *Store) ListTasks() ([]Task, error) {
	rows, err := s.db.Query(`
		SELECT id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at
		FROM tasks
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListTasksByGoalID returns tasks with goal_id = goalID, ordered by created_at.
func (s *Store) ListTasksByGoalID(goalID string) ([]Task, error) {
	if goalID == "" {
		return nil, fmt.Errorf("store: list tasks by goal: goal_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at
		FROM tasks WHERE goal_id = ?
		ORDER BY created_at ASC, id ASC
	`, goalID)
	if err != nil {
		return nil, fmt.Errorf("store: list tasks by goal: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]Task, error) {
	var out []Task
	for rows.Next() {
		var t Task
		var gid, lastVerified sql.NullString
		if err := rows.Scan(
			&t.ID, &gid, &t.Title, &t.Body, &t.SourceType, &t.Confidence, &t.Status, &t.WorkState,
			&t.CreatedAt, &t.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan task: %w", err)
		}
		t.GoalID = nullStrPtr(gid)
		t.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, t)
	}
	return out, rows.Err()
}

type causalEntityRow struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

func (s *Store) listCausalEntities(table string) ([]causalEntityRow, error) {
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM %s
		ORDER BY created_at ASC, id ASC
	`, table))
	if err != nil {
		return nil, fmt.Errorf("store: list %s: %w", table, err)
	}
	defer rows.Close()
	var out []causalEntityRow
	for rows.Next() {
		var r causalEntityRow
		var lastVerified sql.NullString
		if err := rows.Scan(
			&r.ID, &r.Title, &r.Body, &r.SourceType, &r.Confidence, &r.Status,
			&r.CreatedAt, &r.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan %s: %w", table, err)
		}
		r.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, r)
	}
	return out, rows.Err()
}

// ListDecisions returns all decisions ordered by created_at, then id.
func (s *Store) ListDecisions() ([]Decision, error) {
	rows, err := s.listCausalEntities("decisions")
	if err != nil {
		return nil, err
	}
	out := make([]Decision, len(rows))
	for i, r := range rows {
		out[i] = Decision(r)
	}
	return out, nil
}

// ListAssumptions returns all assumptions ordered by created_at, then id.
func (s *Store) ListAssumptions() ([]Assumption, error) {
	rows, err := s.listCausalEntities("assumptions")
	if err != nil {
		return nil, err
	}
	out := make([]Assumption, len(rows))
	for i, r := range rows {
		out[i] = Assumption(r)
	}
	return out, nil
}

// ListDiscoveries returns all discoveries ordered by created_at, then id.
func (s *Store) ListDiscoveries() ([]Discovery, error) {
	rows, err := s.db.Query(`
		SELECT id, title, body, source_type, confidence, status, severity, created_at, updated_at, last_verified_at
		FROM discoveries
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list discoveries: %w", err)
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
			return nil, fmt.Errorf("store: scan discovery: %w", err)
		}
		d.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, d)
	}
	return out, rows.Err()
}

// ListPlanChanges returns all plan_changes ordered by created_at, then id.
func (s *Store) ListPlanChanges() ([]PlanChange, error) {
	rows, err := s.listCausalEntities("plan_changes")
	if err != nil {
		return nil, err
	}
	out := make([]PlanChange, len(rows))
	for i, r := range rows {
		out[i] = PlanChange(r)
	}
	return out, nil
}

// ListClaims returns all claims ordered by created_at, then id.
func (s *Store) ListClaims() ([]Claim, error) {
	rows, err := s.listCausalEntities("claims")
	if err != nil {
		return nil, err
	}
	out := make([]Claim, len(rows))
	for i, r := range rows {
		out[i] = Claim(r)
	}
	return out, nil
}

// ListEvidence returns all evidence ordered by created_at, then id.
func (s *Store) ListEvidence() ([]Evidence, error) {
	rows, err := s.listCausalEntities("evidence")
	if err != nil {
		return nil, err
	}
	out := make([]Evidence, len(rows))
	for i, r := range rows {
		out[i] = Evidence(r)
	}
	return out, nil
}

// GetFileByID loads a file row by id.
func (s *Store) GetFileByID(id string) (FileRecord, error) {
	if id == "" {
		return FileRecord{}, fmt.Errorf("store: get file by id: id required")
	}
	var f FileRecord
	var gitOID, language, lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, path, content_hash, git_oid, language, indexed_at, status, source_type, confidence, created_at, updated_at, last_verified_at
		FROM files WHERE id = ?
	`, id).Scan(
		&f.ID, &f.Path, &f.ContentHash, &gitOID, &language, &f.IndexedAt, &f.Status,
		&f.SourceType, &f.Confidence, &f.CreatedAt, &f.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return FileRecord{}, fmt.Errorf("store: file id %q: %w", id, err)
	}
	if err != nil {
		return FileRecord{}, fmt.Errorf("store: get file by id: %w", err)
	}
	f.GitOID = nullStrPtr(gitOID)
	f.Language = nullStrPtr(language)
	f.LastVerifiedAt = nullStrPtr(lastVerified)
	return f, nil
}
