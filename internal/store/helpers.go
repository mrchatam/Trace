package store

import (
	"database/sql"
	"fmt"
	"strings"
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

// Task list budgets (Laws 6–7): agent/HTTP surfaces must page; seed/graph may use ListAllTasks.
const (
	DefaultTaskListLimit     = 50
	HTTPDefaultTaskListLimit = 100
	MaxTaskListLimit         = 500
)

// TaskListFilter controls bounded task listing.
type TaskListFilter struct {
	GoalID      string
	WorkStates  []string // empty = any state
	Limit       int      // 0 → DefaultTaskListLimit; capped at MaxTaskListLimit
	Cursor      string   // opaque keyset cursor from a prior NextCursor
	IncludeBody bool     // list rows omit body unless true
}

// TaskListResult is a page of tasks.
type TaskListResult struct {
	Tasks      []Task
	NextCursor string
	Truncated  bool
}

// ListAllTasks returns every task ordered by created_at, then id (no limit).
// Prefer ListTasksFiltered for agent/HTTP surfaces.
func (s *Store) ListAllTasks() ([]Task, error) {
	rows, err := s.db.Query(`
		SELECT id, goal_id, title, body, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at
		FROM tasks
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListTasks is an alias for ListAllTasks for older internal callers.
// Deprecated: use ListTasksFiltered (agent/HTTP) or ListAllTasks (seed/graph).
func (s *Store) ListTasks() ([]Task, error) {
	return s.ListAllTasks()
}

// ListTasksByGoalID returns all tasks for a goal (unbounded). Prefer ListTasksFiltered with GoalID.
func (s *Store) ListTasksByGoalID(goalID string) ([]Task, error) {
	if goalID == "" {
		return nil, fmt.Errorf("store: list tasks by goal: goal_id required")
	}
	res, err := s.ListTasksFiltered(TaskListFilter{GoalID: goalID, Limit: -1, IncludeBody: true})
	if err != nil {
		return nil, err
	}
	return res.Tasks, nil
}

// ListTasksFiltered returns a bounded task page. Limit -1 means unbounded (internal).
func (s *Store) ListTasksFiltered(f TaskListFilter) (TaskListResult, error) {
	limit := f.Limit
	unbounded := limit < 0
	if !unbounded {
		if limit <= 0 {
			limit = DefaultTaskListLimit
		}
		if limit > MaxTaskListLimit {
			limit = MaxTaskListLimit
		}
	}

	bodyCol := "'' AS body"
	if f.IncludeBody {
		bodyCol = "body"
	}

	q := fmt.Sprintf(`
		SELECT id, goal_id, title, %s, source_type, confidence, status, work_state, created_at, updated_at, last_verified_at
		FROM tasks
		WHERE 1=1
	`, bodyCol)
	var args []any
	if f.GoalID != "" {
		q += ` AND goal_id = ?`
		args = append(args, f.GoalID)
	}
	if len(f.WorkStates) > 0 {
		q += ` AND work_state IN (` + placeholders(len(f.WorkStates)) + `)`
		for _, ws := range f.WorkStates {
			args = append(args, ws)
		}
	}
	if curCreated, curID, ok := decodeTaskCursor(f.Cursor); ok {
		q += ` AND (created_at > ? OR (created_at = ? AND id > ?))`
		args = append(args, curCreated, curCreated, curID)
	} else if strings.TrimSpace(f.Cursor) != "" {
		return TaskListResult{}, fmt.Errorf("store: list tasks: invalid cursor")
	}
	q += ` ORDER BY created_at ASC, id ASC`
	fetch := limit
	if !unbounded {
		fetch = limit + 1
		q += ` LIMIT ?`
		args = append(args, fetch)
	}

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return TaskListResult{}, fmt.Errorf("store: list tasks filtered: %w", err)
	}
	defer rows.Close()
	tasks, err := scanTasks(rows)
	if err != nil {
		return TaskListResult{}, err
	}
	out := TaskListResult{Tasks: tasks}
	if !unbounded && len(tasks) > limit {
		out.Truncated = true
		out.Tasks = tasks[:limit]
		last := out.Tasks[len(out.Tasks)-1]
		out.NextCursor = encodeTaskCursor(last.CreatedAt, last.ID)
	}
	if out.Tasks == nil {
		out.Tasks = []Task{}
	}
	return out, nil
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	s := "?"
	for i := 1; i < n; i++ {
		s += ",?"
	}
	return s
}

func encodeTaskCursor(createdAt, id string) string {
	return createdAt + "\x1f" + id
}

func decodeTaskCursor(cur string) (createdAt, id string, ok bool) {
	cur = strings.TrimSpace(cur)
	if cur == "" {
		return "", "", false
	}
	i := strings.IndexByte(cur, '\x1f')
	if i <= 0 || i >= len(cur)-1 {
		return "", "", false
	}
	return cur[:i], cur[i+1:], true
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

// CountInTable returns COUNT(*) for an allowlisted entity table (graph budgets).
func (s *Store) CountInTable(table string) (int, error) {
	switch table {
	case "goals", "tasks", "decisions", "assumptions", "discoveries", "plan_changes",
		"claims", "evidence", "reviews", "capabilities", "changes", "regressions":
	default:
		return 0, fmt.Errorf("store: count: table %q not allowlisted", table)
	}
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&n); err != nil {
		return 0, fmt.Errorf("store: count %s: %w", table, err)
	}
	return n, nil
}

// listCausalEntities lists causal rows. limit < 0 means unbounded; 0 means empty.
func (s *Store) listCausalEntities(table string, limit int) ([]causalEntityRow, error) {
	if limit == 0 {
		return nil, nil
	}
	q := fmt.Sprintf(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM %s
		ORDER BY created_at ASC, id ASC
	`, table)
	var rows *sql.Rows
	var err error
	if limit > 0 {
		q += ` LIMIT ?`
		rows, err = s.db.Query(q, limit)
	} else {
		rows, err = s.db.Query(q)
	}
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
	rows, err := s.listCausalEntities("decisions", -1)
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
	rows, err := s.listCausalEntities("assumptions", -1)
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
	rows, err := s.listCausalEntities("plan_changes", -1)
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
	rows, err := s.listCausalEntities("claims", -1)
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
	rows, err := s.listCausalEntities("evidence", -1)
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
