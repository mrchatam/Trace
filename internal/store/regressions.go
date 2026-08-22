package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const listOpenRegressionsLimit = 32

// Regression source kinds (thin derivation only).
const (
	RegressionSourceEvaluation         = "evaluation"
	RegressionSourceContradictedEffect = "contradicted_effect"
)

// Regression attribution vocabulary. Create is always correlated.
const (
	RegressionAttributionCorrelated   = "correlated"
	RegressionAttributionHypothesized = "hypothesized"
	RegressionAttributionCaused       = "caused"
)

// Regression lifecycle (not provenance).
const (
	RegressionStatusOpen       = "OPEN"
	RegressionStatusResolved   = "RESOLVED"
	RegressionStatusSuperseded = "SUPERSEDED"
)

const emptyJSONArray = "[]"

// Regression is a thin derived row from an evaluation flag or contradicted effect.
type Regression struct {
	ID          string
	TaskID      string
	SourceKind  string
	SourceID    string
	Dimension   string
	Attribution string
	Status      string
	Summary     string
	Actor       string
	SourceType  string
	Confidence  float64
	CreatedAt   string
	UpdatedAt   string
}

// Reflection is a structured learning artifact (JSON arrays, no essay body).
type Reflection struct {
	ID                         string
	TaskID                     string
	Summary                    string
	InvalidatedAssumptionsJSON string
	NewDependenciesJSON        string
	UsefulTestsJSON            string
	BroadenTestsNote           string
	Actor                      string
	SourceType                 string
	Confidence                 float64
	CreatedAt                  string
	UpdatedAt                  string
}

func validateRegressionEnums(r Regression) error {
	switch r.SourceKind {
	case RegressionSourceEvaluation, RegressionSourceContradictedEffect:
	default:
		return fmt.Errorf("store: upsert regression: unknown source_kind %q", r.SourceKind)
	}
	switch r.Attribution {
	case RegressionAttributionCorrelated, RegressionAttributionHypothesized, RegressionAttributionCaused:
	default:
		return fmt.Errorf("store: upsert regression: unknown attribution %q", r.Attribution)
	}
	switch r.Status {
	case RegressionStatusOpen, RegressionStatusResolved, RegressionStatusSuperseded:
	default:
		return fmt.Errorf("store: upsert regression: unknown status %q", r.Status)
	}
	return nil
}

func (s *Store) regressionExists(id string) (bool, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM regressions WHERE id = ?`, id).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("store: regression exists: %w", err)
	}
	return n > 0, nil
}

// UpsertRegression inserts or replaces a regression by id. Empty ID allocates a UUID.
// New rows must have attribution=correlated (caused is never a create default).
func (s *Store) UpsertRegression(r Regression) (Regression, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.TaskID == "" {
		return Regression{}, fmt.Errorf("store: upsert regression: task_id required")
	}
	if r.SourceID == "" {
		return Regression{}, fmt.Errorf("store: upsert regression: source_id required")
	}
	if r.Attribution == "" {
		r.Attribution = RegressionAttributionCorrelated
	}
	if r.Status == "" {
		r.Status = RegressionStatusOpen
	}
	if err := validateRegressionEnums(r); err != nil {
		return Regression{}, err
	}
	exists, err := s.regressionExists(r.ID)
	if err != nil {
		return Regression{}, err
	}
	if !exists && r.Attribution != RegressionAttributionCorrelated {
		return Regression{}, fmt.Errorf("store: upsert regression: create attribution must be correlated")
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	_, err = s.db.Exec(`
		INSERT INTO regressions(
			id, task_id, source_kind, source_id, dimension, attribution, status,
			summary, actor, source_type, confidence, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			attribution = excluded.attribution,
			status = excluded.status,
			summary = excluded.summary,
			actor = excluded.actor,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			updated_at = excluded.updated_at
	`, r.ID, r.TaskID, r.SourceKind, r.SourceID, r.Dimension, r.Attribution, r.Status,
		r.Summary, r.Actor, r.SourceType, r.Confidence, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return Regression{}, fmt.Errorf("store: upsert regression: %w", err)
	}
	out, err := s.GetRegression(r.ID)
	if err != nil {
		return Regression{}, err
	}
	if err := s.SyncEntityFTS("regression", out.ID); err != nil {
		return Regression{}, err
	}
	return out, nil
}

// GetRegression loads a regression by id.
func (s *Store) GetRegression(id string) (Regression, error) {
	if id == "" {
		return Regression{}, fmt.Errorf("store: get regression: id required")
	}
	var r Regression
	err := s.db.QueryRow(`
		SELECT id, task_id, source_kind, source_id, dimension, attribution, status,
			summary, actor, source_type, confidence, created_at, updated_at
		FROM regressions WHERE id = ?
	`, id).Scan(
		&r.ID, &r.TaskID, &r.SourceKind, &r.SourceID, &r.Dimension, &r.Attribution, &r.Status,
		&r.Summary, &r.Actor, &r.SourceType, &r.Confidence, &r.CreatedAt, &r.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Regression{}, fmt.Errorf("store: regression %q: %w", id, err)
	}
	if err != nil {
		return Regression{}, fmt.Errorf("store: get regression: %w", err)
	}
	return r, nil
}

// GetRegressionBySource loads the unique row for (source_kind, source_id, dimension).
func (s *Store) GetRegressionBySource(sourceKind, sourceID, dimension string) (Regression, error) {
	if sourceKind == "" || sourceID == "" {
		return Regression{}, fmt.Errorf("store: get regression by source: source_kind and source_id required")
	}
	var r Regression
	err := s.db.QueryRow(`
		SELECT id, task_id, source_kind, source_id, dimension, attribution, status,
			summary, actor, source_type, confidence, created_at, updated_at
		FROM regressions WHERE source_kind = ? AND source_id = ? AND dimension = ?
	`, sourceKind, sourceID, dimension).Scan(
		&r.ID, &r.TaskID, &r.SourceKind, &r.SourceID, &r.Dimension, &r.Attribution, &r.Status,
		&r.Summary, &r.Actor, &r.SourceType, &r.Confidence, &r.CreatedAt, &r.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Regression{}, fmt.Errorf("store: regression source: %w", err)
	}
	if err != nil {
		return Regression{}, fmt.Errorf("store: get regression by source: %w", err)
	}
	return r, nil
}

func scanRegressions(rows *sql.Rows) ([]Regression, error) {
	var out []Regression
	for rows.Next() {
		var r Regression
		if err := rows.Scan(
			&r.ID, &r.TaskID, &r.SourceKind, &r.SourceID, &r.Dimension, &r.Attribution, &r.Status,
			&r.Summary, &r.Actor, &r.SourceType, &r.Confidence, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan regression: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []Regression{}
	}
	return out, rows.Err()
}

const regressionSelect = `
		SELECT id, task_id, source_kind, source_id, dimension, attribution, status,
			summary, actor, source_type, confidence, created_at, updated_at
		FROM regressions`

// ListAllRegressions returns every regression row.
func (s *Store) ListAllRegressions() ([]Regression, error) {
	rows, err := s.db.Query(regressionSelect + `
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all regressions: %w", err)
	}
	defer rows.Close()
	return scanRegressions(rows)
}

// ListAllReflections returns every reflection row.
func (s *Store) ListAllReflections() ([]Reflection, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, summary, invalidated_assumptions_json, new_dependencies_json,
			useful_tests_json, broaden_tests_note, actor, source_type, confidence,
			created_at, updated_at
		FROM reflections
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all reflections: %w", err)
	}
	defer rows.Close()
	var out []Reflection
	for rows.Next() {
		var r Reflection
		if err := rows.Scan(
			&r.ID, &r.TaskID, &r.Summary, &r.InvalidatedAssumptionsJSON, &r.NewDependenciesJSON,
			&r.UsefulTestsJSON, &r.BroadenTestsNote, &r.Actor, &r.SourceType, &r.Confidence,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan reflection: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []Reflection{}
	}
	return out, rows.Err()
}

// HasOpenRegression reports whether the task has any OPEN regression.
func (s *Store) HasOpenRegression(taskID string) (bool, error) {
	n, err := s.CountOpenRegressionsByTaskID(taskID)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// CountOpenRegressionsByTaskID counts OPEN regressions for a task.
func (s *Store) CountOpenRegressionsByTaskID(taskID string) (int, error) {
	if taskID == "" {
		return 0, fmt.Errorf("store: count open regressions: task_id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM regressions WHERE task_id = ? AND status = 'OPEN'
	`, taskID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count open regressions: %w", err)
	}
	return n, nil
}

// ListRegressionsByChangeID returns regressions linked to a change via regression_associated_change.
func (s *Store) ListRegressionsByChangeID(changeID string) ([]Regression, error) {
	if changeID == "" {
		return nil, fmt.Errorf("store: list regressions by change: change_id required")
	}
	rows, err := s.db.Query(`
		SELECT r.id, r.task_id, r.source_kind, r.source_id, r.dimension, r.attribution, r.status,
			r.summary, r.actor, r.source_type, r.confidence, r.created_at, r.updated_at
		FROM regressions r
		INNER JOIN entity_links l ON l.from_type = 'regression' AND l.from_id = r.id
		WHERE l.rel = 'regression_associated_change' AND l.to_type = 'change' AND l.to_id = ?
		ORDER BY r.created_at ASC, r.id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("store: list regressions by change: %w", err)
	}
	defer rows.Close()
	return scanRegressions(rows)
}

// ListRegressionsRecent returns up to limit regressions newest-first. Empty taskID lists all tasks.
func (s *Store) ListRegressionsRecent(limit int, taskID string) ([]Regression, error) {
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	taskID = strings.TrimSpace(taskID)
	var (
		rows *sql.Rows
		err  error
	)
	if taskID != "" {
		rows, err = s.db.Query(regressionSelect+`
			WHERE task_id = ?
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, taskID, limit)
	} else {
		rows, err = s.db.Query(regressionSelect+`
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("store: list regressions recent: %w", err)
	}
	defer rows.Close()
	return scanRegressions(rows)
}

// ListOpenRegressions returns OPEN regressions for a task, oldest first, max 32.
func (s *Store) ListOpenRegressions(taskID string) ([]Regression, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list open regressions: task_id required")
	}
	rows, err := s.db.Query(regressionSelect+`
		WHERE task_id = ? AND status = 'OPEN'
		ORDER BY created_at ASC, id ASC
		LIMIT ?
	`, taskID, listOpenRegressionsLimit)
	if err != nil {
		return nil, fmt.Errorf("store: list open regressions: %w", err)
	}
	defer rows.Close()
	return scanRegressions(rows)
}

func requireJSONArray(label, raw string) (string, error) {
	if raw == "" {
		return emptyJSONArray, nil
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("store: upsert reflection: %s must be valid JSON", label)
	}
	switch v.(type) {
	case []any:
		return raw, nil
	default:
		return "", fmt.Errorf("store: upsert reflection: %s must be a JSON array", label)
	}
}

// UpsertReflection inserts or replaces a reflection by id. Empty ID allocates a UUID.
func (s *Store) UpsertReflection(r Reflection) (Reflection, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.TaskID == "" {
		return Reflection{}, fmt.Errorf("store: upsert reflection: task_id required")
	}
	inv, err := requireJSONArray("invalidated_assumptions_json", r.InvalidatedAssumptionsJSON)
	if err != nil {
		return Reflection{}, err
	}
	deps, err := requireJSONArray("new_dependencies_json", r.NewDependenciesJSON)
	if err != nil {
		return Reflection{}, err
	}
	tests, err := requireJSONArray("useful_tests_json", r.UsefulTestsJSON)
	if err != nil {
		return Reflection{}, err
	}
	r.InvalidatedAssumptionsJSON = inv
	r.NewDependenciesJSON = deps
	r.UsefulTestsJSON = tests
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	_, err = s.db.Exec(`
		INSERT INTO reflections(
			id, task_id, summary, invalidated_assumptions_json, new_dependencies_json,
			useful_tests_json, broaden_tests_note, actor, source_type, confidence,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			summary = excluded.summary,
			invalidated_assumptions_json = excluded.invalidated_assumptions_json,
			new_dependencies_json = excluded.new_dependencies_json,
			useful_tests_json = excluded.useful_tests_json,
			broaden_tests_note = excluded.broaden_tests_note,
			actor = excluded.actor,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			updated_at = excluded.updated_at
	`, r.ID, r.TaskID, r.Summary, r.InvalidatedAssumptionsJSON, r.NewDependenciesJSON,
		r.UsefulTestsJSON, r.BroadenTestsNote, r.Actor, r.SourceType, r.Confidence,
		r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return Reflection{}, fmt.Errorf("store: upsert reflection: %w", err)
	}
	out, err := s.GetReflection(r.ID)
	if err != nil {
		return Reflection{}, err
	}
	if err := s.SyncEntityFTS("reflection", out.ID); err != nil {
		return Reflection{}, err
	}
	return out, nil
}

// GetReflection loads a reflection by id.
func (s *Store) GetReflection(id string) (Reflection, error) {
	if id == "" {
		return Reflection{}, fmt.Errorf("store: get reflection: id required")
	}
	var r Reflection
	err := s.db.QueryRow(`
		SELECT id, task_id, summary, invalidated_assumptions_json, new_dependencies_json,
			useful_tests_json, broaden_tests_note, actor, source_type, confidence,
			created_at, updated_at
		FROM reflections WHERE id = ?
	`, id).Scan(
		&r.ID, &r.TaskID, &r.Summary, &r.InvalidatedAssumptionsJSON, &r.NewDependenciesJSON,
		&r.UsefulTestsJSON, &r.BroadenTestsNote, &r.Actor, &r.SourceType, &r.Confidence,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Reflection{}, fmt.Errorf("store: reflection %q: %w", id, err)
	}
	if err != nil {
		return Reflection{}, fmt.Errorf("store: get reflection: %w", err)
	}
	return r, nil
}

// ListReflectionsByTaskID returns reflections for a task, oldest first.
func (s *Store) ListReflectionsByTaskID(taskID string) ([]Reflection, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list reflections: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, summary, invalidated_assumptions_json, new_dependencies_json,
			useful_tests_json, broaden_tests_note, actor, source_type, confidence,
			created_at, updated_at
		FROM reflections
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list reflections: %w", err)
	}
	defer rows.Close()
	var out []Reflection
	for rows.Next() {
		var r Reflection
		if err := rows.Scan(
			&r.ID, &r.TaskID, &r.Summary, &r.InvalidatedAssumptionsJSON, &r.NewDependenciesJSON,
			&r.UsefulTestsJSON, &r.BroadenTestsNote, &r.Actor, &r.SourceType, &r.Confidence,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan reflection: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []Reflection{}
	}
	return out, rows.Err()
}
