package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Change lifecycle (not provenance ACTIVE/STALE).
const (
	ChangeStatusOpen       = "OPEN"
	ChangeStatusRecorded   = "RECORDED"
	ChangeStatusCompared   = "COMPARED"
	ChangeStatusSuperseded = "SUPERSEDED"
)

// Effect comparison vocabulary. Empty until actual is recorded.
const (
	EffectComparisonNone               = ""
	EffectComparisonSupported          = "supported"
	EffectComparisonPartiallySupported = "partially_supported"
	EffectComparisonContradicted       = "contradicted"
)

// Change is a first-class Git SHA + path-ref object (no blobs).
type Change struct {
	ID             string
	TaskID         string
	GitCommit      string
	ParentChangeID string
	Actor          string
	Reason         string
	Status         string
	SourceType     string
	Confidence     float64
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
	RowID          int64 // SQLite rowid when loaded; 0 if unset.
}

// ChangePath is a repo-relative path touched by a change. No file contents.
type ChangePath struct {
	ChangeID string
	Path     string
	Status   string
	SymbolID string
}

// Effect is one expected-vs-actual dimension on a change.
type Effect struct {
	ID         string
	ChangeID   string
	Dimension  string
	Expected   string
	Actual     string
	Comparison string
	Confidence float64
	SourceType string
	CreatedAt  string
	UpdatedAt  string
}

// UpsertChange inserts or replaces a change by id. Empty ID allocates a UUID.
func (s *Store) UpsertChange(c Change) (Change, error) {
	now := nowRFC3339()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	if c.TaskID == "" {
		return Change{}, fmt.Errorf("store: upsert change: task_id required")
	}
	if c.Status == "" {
		c.Status = ChangeStatusOpen
	}
	if c.CreatedAt == "" {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO changes(
			id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			task_id = excluded.task_id,
			git_commit = excluded.git_commit,
			parent_change_id = excluded.parent_change_id,
			actor = excluded.actor,
			reason = excluded.reason,
			status = excluded.status,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, c.ID, c.TaskID, c.GitCommit, c.ParentChangeID, c.Actor, c.Reason, c.Status,
		c.SourceType, c.Confidence, c.CreatedAt, c.UpdatedAt, nullStr(c.LastVerifiedAt))
	if err != nil {
		return Change{}, fmt.Errorf("store: upsert change: %w", err)
	}
	out, err := s.GetChange(c.ID)
	if err != nil {
		return Change{}, err
	}
	if err := s.SyncEntityFTS("change", out.ID); err != nil {
		return Change{}, err
	}
	return out, nil
}

// GetChangeByGitCommit loads a change by git_commit OID. Missing returns sql.ErrNoRows.
func (s *Store) GetChangeByGitCommit(oid string) (Change, error) {
	if oid == "" {
		return Change{}, fmt.Errorf("store: get change by git commit: oid required")
	}
	var c Change
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at
		FROM changes WHERE git_commit = ?
	`, oid).Scan(
		&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
		&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Change{}, err
	}
	if err != nil {
		return Change{}, fmt.Errorf("store: get change by git commit: %w", err)
	}
	c.LastVerifiedAt = nullStrPtr(lastVerified)
	return c, nil
}

// GetChange loads a change by id.
func (s *Store) GetChange(id string) (Change, error) {
	if id == "" {
		return Change{}, fmt.Errorf("store: get change: id required")
	}
	var c Change
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at
		FROM changes WHERE id = ?
	`, id).Scan(
		&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
		&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Change{}, fmt.Errorf("store: change %q: %w", id, err)
	}
	if err != nil {
		return Change{}, fmt.Errorf("store: get change: %w", err)
	}
	c.LastVerifiedAt = nullStrPtr(lastVerified)
	return c, nil
}

// ListChangesRecent returns up to limit changes newest-first (created_at DESC, id DESC).
// limit defaults to 32 and is capped at 64. Empty taskID lists across all tasks.
func (s *Store) ListChangesRecent(limit int, taskID string) ([]Change, error) {
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
		rows, err = s.db.Query(`
			SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
				source_type, confidence, created_at, updated_at, last_verified_at
			FROM changes
			WHERE task_id = ?
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, taskID, limit)
	} else {
		rows, err = s.db.Query(`
			SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
				source_type, confidence, created_at, updated_at, last_verified_at
			FROM changes
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("store: list changes recent: %w", err)
	}
	defer rows.Close()
	var out []Change
	for rows.Next() {
		var c Change
		var lastVerified sql.NullString
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
			&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan change: %w", err)
		}
		c.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, c)
	}
	if out == nil {
		out = []Change{}
	}
	return out, rows.Err()
}

// ListAllChanges returns every change row.
func (s *Store) ListAllChanges() ([]Change, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at
		FROM changes
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all changes: %w", err)
	}
	defer rows.Close()
	var out []Change
	for rows.Next() {
		var c Change
		var lastVerified sql.NullString
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
			&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan change: %w", err)
		}
		c.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, c)
	}
	if out == nil {
		out = []Change{}
	}
	return out, rows.Err()
}

// ListAllEffects returns every effect row.
func (s *Store) ListAllEffects() ([]Effect, error) {
	rows, err := s.db.Query(`
		SELECT id, change_id, dimension, expected, actual, comparison,
			confidence, source_type, created_at, updated_at
		FROM effects
		ORDER BY change_id ASC, dimension ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all effects: %w", err)
	}
	defer rows.Close()
	var out []Effect
	for rows.Next() {
		var e Effect
		if err := rows.Scan(
			&e.ID, &e.ChangeID, &e.Dimension, &e.Expected, &e.Actual, &e.Comparison,
			&e.Confidence, &e.SourceType, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan effect: %w", err)
		}
		out = append(out, e)
	}
	if out == nil {
		out = []Effect{}
	}
	return out, rows.Err()
}

// UpsertChangePath inserts or replaces a path row keyed by (change_id, path).
func (s *Store) UpsertChangePath(p ChangePath) (ChangePath, error) {
	if p.ChangeID == "" {
		return ChangePath{}, fmt.Errorf("store: upsert change path: change_id required")
	}
	if p.Path == "" {
		return ChangePath{}, fmt.Errorf("store: upsert change path: path required")
	}
	_, err := s.db.Exec(`
		INSERT INTO change_paths(change_id, path, status, symbol_id)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(change_id, path) DO UPDATE SET
			status = excluded.status,
			symbol_id = excluded.symbol_id
	`, p.ChangeID, p.Path, p.Status, p.SymbolID)
	if err != nil {
		return ChangePath{}, fmt.Errorf("store: upsert change path: %w", err)
	}
	return s.GetChangePath(p.ChangeID, p.Path)
}

// ListChangesByTaskID returns changes for a task, oldest first.
func (s *Store) ListChangesByTaskID(taskID string) ([]Change, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list changes: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at
		FROM changes
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list changes: %w", err)
	}
	defer rows.Close()
	var out []Change
	for rows.Next() {
		var c Change
		var lastVerified sql.NullString
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
			&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan change: %w", err)
		}
		c.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, c)
	}
	if out == nil {
		out = []Change{}
	}
	return out, rows.Err()
}

// InsertChangePath inserts one path row. Empty path fails closed.
func (s *Store) InsertChangePath(p ChangePath) (ChangePath, error) {
	if p.ChangeID == "" {
		return ChangePath{}, fmt.Errorf("store: insert change path: change_id required")
	}
	if p.Path == "" {
		return ChangePath{}, fmt.Errorf("store: insert change path: path required")
	}
	_, err := s.db.Exec(`
		INSERT INTO change_paths(change_id, path, status, symbol_id)
		VALUES (?, ?, ?, ?)
	`, p.ChangeID, p.Path, p.Status, p.SymbolID)
	if err != nil {
		return ChangePath{}, fmt.Errorf("store: insert change path: %w", err)
	}
	return p, nil
}

// ListHighChurnPaths returns paths touched on >= minChanges distinct changes for taskID.
func (s *Store) ListHighChurnPaths(taskID string, minChanges int) ([]string, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list high churn paths: task_id required")
	}
	if minChanges < 1 {
		return nil, fmt.Errorf("store: list high churn paths: min_changes must be >= 1")
	}
	rows, err := s.db.Query(`
		SELECT cp.path
		FROM change_paths cp
		JOIN changes c ON c.id = cp.change_id
		WHERE c.task_id = ?
		GROUP BY cp.path
		HAVING COUNT(DISTINCT cp.change_id) >= ?
		ORDER BY cp.path ASC
	`, taskID, minChanges)
	if err != nil {
		return nil, fmt.Errorf("store: list high churn paths: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, fmt.Errorf("store: scan high churn path: %w", err)
		}
		out = append(out, path)
	}
	if out == nil {
		out = []string{}
	}
	return out, rows.Err()
}

// CountChangePaths returns the number of paths recorded for a change.
func (s *Store) CountChangePaths(changeID string) (int, error) {
	if changeID == "" {
		return 0, fmt.Errorf("store: count change paths: change_id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM change_paths WHERE change_id = ?
	`, changeID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count change paths: %w", err)
	}
	return n, nil
}

// ListChangePaths returns paths for a change, ordered by path.
func (s *Store) ListChangePaths(changeID string) ([]ChangePath, error) {
	if changeID == "" {
		return nil, fmt.Errorf("store: list change paths: change_id required")
	}
	rows, err := s.db.Query(`
		SELECT change_id, path, status, symbol_id
		FROM change_paths
		WHERE change_id = ?
		ORDER BY path ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("store: list change paths: %w", err)
	}
	defer rows.Close()
	var out []ChangePath
	for rows.Next() {
		var p ChangePath
		if err := rows.Scan(&p.ChangeID, &p.Path, &p.Status, &p.SymbolID); err != nil {
			return nil, fmt.Errorf("store: scan change path: %w", err)
		}
		out = append(out, p)
	}
	if out == nil {
		out = []ChangePath{}
	}
	return out, rows.Err()
}

// GetChangePath loads one path row. Missing is an error.
func (s *Store) GetChangePath(changeID, path string) (ChangePath, error) {
	if changeID == "" || path == "" {
		return ChangePath{}, fmt.Errorf("store: get change path: change_id and path required")
	}
	var p ChangePath
	err := s.db.QueryRow(`
		SELECT change_id, path, status, symbol_id
		FROM change_paths
		WHERE change_id = ? AND path = ?
	`, changeID, path).Scan(&p.ChangeID, &p.Path, &p.Status, &p.SymbolID)
	if err == sql.ErrNoRows {
		return ChangePath{}, fmt.Errorf("store: change path %q on %q: %w", path, changeID, err)
	}
	if err != nil {
		return ChangePath{}, fmt.Errorf("store: get change path: %w", err)
	}
	return p, nil
}

// UpsertEffect inserts or replaces an effect by id. Empty ID allocates a UUID.
func (s *Store) UpsertEffect(e Effect) (Effect, error) {
	now := nowRFC3339()
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.ChangeID == "" {
		return Effect{}, fmt.Errorf("store: upsert effect: change_id required")
	}
	if e.Dimension == "" {
		return Effect{}, fmt.Errorf("store: upsert effect: dimension required")
	}
	if e.CreatedAt == "" {
		e.CreatedAt = now
	}
	e.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO effects(
			id, change_id, dimension, expected, actual, comparison,
			confidence, source_type, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			change_id = excluded.change_id,
			dimension = excluded.dimension,
			expected = excluded.expected,
			actual = excluded.actual,
			comparison = excluded.comparison,
			confidence = excluded.confidence,
			source_type = excluded.source_type,
			updated_at = excluded.updated_at
	`, e.ID, e.ChangeID, e.Dimension, e.Expected, e.Actual, e.Comparison,
		e.Confidence, e.SourceType, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return Effect{}, fmt.Errorf("store: upsert effect: %w", err)
	}
	return s.GetEffect(e.ID)
}

// GetEffect loads an effect by id.
func (s *Store) GetEffect(id string) (Effect, error) {
	if id == "" {
		return Effect{}, fmt.Errorf("store: get effect: id required")
	}
	var e Effect
	err := s.db.QueryRow(`
		SELECT id, change_id, dimension, expected, actual, comparison,
			confidence, source_type, created_at, updated_at
		FROM effects WHERE id = ?
	`, id).Scan(
		&e.ID, &e.ChangeID, &e.Dimension, &e.Expected, &e.Actual, &e.Comparison,
		&e.Confidence, &e.SourceType, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Effect{}, fmt.Errorf("store: effect %q: %w", id, err)
	}
	if err != nil {
		return Effect{}, fmt.Errorf("store: get effect: %w", err)
	}
	return e, nil
}

// GetEffectByChangeDimension loads the unique effect for (change_id, dimension).
func (s *Store) GetEffectByChangeDimension(changeID, dimension string) (Effect, error) {
	if changeID == "" || dimension == "" {
		return Effect{}, fmt.Errorf("store: get effect by dimension: change_id and dimension required")
	}
	var e Effect
	err := s.db.QueryRow(`
		SELECT id, change_id, dimension, expected, actual, comparison,
			confidence, source_type, created_at, updated_at
		FROM effects WHERE change_id = ? AND dimension = ?
	`, changeID, dimension).Scan(
		&e.ID, &e.ChangeID, &e.Dimension, &e.Expected, &e.Actual, &e.Comparison,
		&e.Confidence, &e.SourceType, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Effect{}, fmt.Errorf("store: effect dimension %q on %q: %w", dimension, changeID, err)
	}
	if err != nil {
		return Effect{}, fmt.Errorf("store: get effect by dimension: %w", err)
	}
	return e, nil
}

// ListEffectsByChangeID returns effects for a change, ordered by dimension.
func (s *Store) ListEffectsByChangeID(changeID string) ([]Effect, error) {
	if changeID == "" {
		return nil, fmt.Errorf("store: list effects: change_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, change_id, dimension, expected, actual, comparison,
			confidence, source_type, created_at, updated_at
		FROM effects
		WHERE change_id = ?
		ORDER BY dimension ASC, id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("store: list effects: %w", err)
	}
	defer rows.Close()
	var out []Effect
	for rows.Next() {
		var e Effect
		if err := rows.Scan(
			&e.ID, &e.ChangeID, &e.Dimension, &e.Expected, &e.Actual, &e.Comparison,
			&e.Confidence, &e.SourceType, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan effect: %w", err)
		}
		out = append(out, e)
	}
	if out == nil {
		out = []Effect{}
	}
	return out, rows.Err()
}
