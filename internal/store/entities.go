package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Status values for semantic entities.
const (
	StatusActive     = "ACTIVE"
	StatusStale      = "STALE"
	StatusSuperseded = "SUPERSEDED"
)

// Goal is a durable intent entity.
type Goal struct {
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

// Work-state values for tasks (separate from provenance Status).
const (
	WorkStatePending        = "PENDING"
	WorkStateInProgress     = "IN_PROGRESS"
	WorkStateAwaitingReview = "AWAITING_REVIEW"
	WorkStateBlocked        = "BLOCKED"
	WorkStateFailed         = "FAILED"
	WorkStateDone           = "DONE"
	WorkStateStale          = "STALE"
	WorkStateSkipped        = "SKIPPED"
)

// Task is a durable work entity (optional goal_id stub FK).
// WorkState is the task machine; Status is provenance (ACTIVE|STALE|SUPERSEDED).
type Task struct {
	ID             string
	GoalID         *string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	WorkState      string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Event is an append-only history row (DR-EVT). No update/delete API.
type Event struct {
	ID          string
	TS          string
	Type        string
	EntityType  string
	EntityID    string
	PayloadJSON string
}

// UpsertGoal inserts or replaces a goal by id. Empty ID allocates a UUID.
func (s *Store) UpsertGoal(g Goal) (Goal, error) {
	now := nowRFC3339()
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	if g.Status == "" {
		g.Status = StatusActive
	}
	if g.CreatedAt == "" {
		g.CreatedAt = now
	}
	g.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO goals(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, g.ID, g.Title, g.Body, g.SourceType, g.Confidence, g.Status, g.CreatedAt, g.UpdatedAt, nullStr(g.LastVerifiedAt))
	if err != nil {
		return Goal{}, fmt.Errorf("store: upsert goal: %w", err)
	}
	out, err := s.GetGoal(g.ID)
	if err != nil {
		return Goal{}, err
	}
	if err := s.SyncEntityFTS("goal", out.ID); err != nil {
		return Goal{}, err
	}
	return out, nil
}

// GetGoal loads a goal by id.
func (s *Store) GetGoal(id string) (Goal, error) {
	var g Goal
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM goals WHERE id = ?
	`, id).Scan(
		&g.ID, &g.Title, &g.Body, &g.SourceType, &g.Confidence, &g.Status,
		&g.CreatedAt, &g.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Goal{}, fmt.Errorf("store: goal %q: %w", id, err)
	}
	if err != nil {
		return Goal{}, fmt.Errorf("store: get goal: %w", err)
	}
	g.LastVerifiedAt = nullStrPtr(lastVerified)
	return g, nil
}

// UpsertTask inserts or replaces a task by id. Empty ID allocates a UUID.
// Empty WorkState defaults to PENDING.
func (s *Store) UpsertTask(t Task) (Task, error) {
	now := nowRFC3339()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	if t.Status == "" {
		t.Status = StatusActive
	}
	if t.WorkState == "" {
		t.WorkState = WorkStatePending
	}
	if t.CreatedAt == "" {
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO tasks(id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			goal_id = excluded.goal_id,
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			work_state = excluded.work_state,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, t.ID, nullStr(t.GoalID), t.Title, t.Body, t.SourceType, t.Confidence, t.Status, t.WorkState, t.CreatedAt, t.UpdatedAt, nullStr(t.LastVerifiedAt))
	if err != nil {
		return Task{}, fmt.Errorf("store: upsert task: %w", err)
	}
	out, err := s.GetTask(t.ID)
	if err != nil {
		return Task{}, err
	}
	if err := s.SyncEntityFTS("task", out.ID); err != nil {
		return Task{}, err
	}
	return out, nil
}

// UpsertTaskFromSeed upserts a task from seed import. On conflict, work_state is preserved.
func (s *Store) UpsertTaskFromSeed(t Task) (Task, error) {
	now := nowRFC3339()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	if t.Status == "" {
		t.Status = StatusActive
	}
	if t.WorkState == "" {
		t.WorkState = WorkStatePending
	}
	if t.CreatedAt == "" {
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO tasks(id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			goal_id = excluded.goal_id,
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, t.ID, nullStr(t.GoalID), t.Title, t.Body, t.SourceType, t.Confidence, t.Status, t.WorkState, t.CreatedAt, t.UpdatedAt, nullStr(t.LastVerifiedAt))
	if err != nil {
		return Task{}, fmt.Errorf("store: upsert task from seed: %w", err)
	}
	out, err := s.GetTask(t.ID)
	if err != nil {
		return Task{}, err
	}
	if err := s.SyncEntityFTS("task", out.ID); err != nil {
		return Task{}, err
	}
	return out, nil
}

// GetTask loads a task by id.
func (s *Store) GetTask(id string) (Task, error) {
	var t Task
	var goalID, lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at
		FROM tasks WHERE id = ?
	`, id).Scan(
		&t.ID, &goalID, &t.Title, &t.Body, &t.SourceType, &t.Confidence, &t.Status, &t.WorkState,
		&t.CreatedAt, &t.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Task{}, fmt.Errorf("store: task %q: %w", id, err)
	}
	if err != nil {
		return Task{}, fmt.Errorf("store: get task: %w", err)
	}
	t.GoalID = nullStrPtr(goalID)
	t.LastVerifiedAt = nullStrPtr(lastVerified)
	return t, nil
}

// AppendEvent inserts an event. Empty ID/TS/PayloadJSON get defaults.
// There is no Update/Delete for events.
func (s *Store) AppendEvent(e Event) (Event, error) {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.TS == "" {
		e.TS = nowRFC3339()
	}
	if e.PayloadJSON == "" {
		e.PayloadJSON = "{}"
	}
	if e.Type == "" || e.EntityType == "" || e.EntityID == "" {
		return Event{}, fmt.Errorf("store: append event: type, entity_type, and entity_id are required")
	}

	_, err := s.db.Exec(`
		INSERT INTO events(id, ts, type, entity_type, entity_id, payload_json)
		VALUES (?, ?, ?, ?, ?, ?)
	`, e.ID, e.TS, e.Type, e.EntityType, e.EntityID, e.PayloadJSON)
	if err != nil {
		return Event{}, fmt.Errorf("store: append event: %w", err)
	}
	return e, nil
}

// ListEventsByEntity returns events for an entity ordered by ts ascending.
func (s *Store) ListEventsByEntity(entityType, entityID string) ([]Event, error) {
	rows, err := s.db.Query(`
		SELECT id, ts, type, entity_type, entity_id, payload_json
		FROM events
		WHERE entity_type = ? AND entity_id = ?
		ORDER BY ts ASC, id ASC
	`, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("store: list events by entity: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

// ListRecentEvents returns the most recent events (newest first), limited.
func (s *Store) ListRecentEvents(limit int) ([]Event, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.Query(`
		SELECT id, ts, type, entity_type, entity_id, payload_json
		FROM events
		ORDER BY ts DESC, id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list recent events: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

func scanEvents(rows *sql.Rows) ([]Event, error) {
	var out []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.TS, &e.Type, &e.EntityType, &e.EntityID, &e.PayloadJSON); err != nil {
			return nil, fmt.Errorf("store: scan event: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func nullStr(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullStrPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

// NormalizePath forces repo-relative forward-slash paths.
func NormalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	return strings.TrimPrefix(path, "./")
}

// EnsureRFC3339 returns t formatted as RFC3339 UTC, or now if zero.
func EnsureRFC3339(t time.Time) string {
	if t.IsZero() {
		return nowRFC3339()
	}
	return t.UTC().Format(time.RFC3339)
}
