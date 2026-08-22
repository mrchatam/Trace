package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PlanPhase is a coarse phase under a goal.
type PlanPhase struct {
	ID        string
	GoalID    string
	Title     string
	Body      string
	Ord       int
	Status    string
	CreatedAt string
	UpdatedAt string
}

// PlanScope is a coarse scope under a phase.
type PlanScope struct {
	ID              string
	PhaseID         string
	Title           string
	Body            string
	Ord             int
	Status          string
	AutoReplanCount int
	CreatedAt       string
	UpdatedAt       string
}

// ScopeDeepPlan is one deep-plan revision for a scope.
type ScopeDeepPlan struct {
	ID          string
	ScopeID     string
	ContentJSON string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

// GoalPlanState tracks the current scope pointer for a goal.
type GoalPlanState struct {
	GoalID         string
	CurrentScopeID *string
	UpdatedAt      string
}

// InsertPlanPhase inserts a phase row. Empty ID allocates a UUID.
func (s *Store) InsertPlanPhase(p PlanPhase) (PlanPhase, error) {
	now := nowRFC3339()
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.GoalID == "" {
		return PlanPhase{}, fmt.Errorf("store: insert plan phase: goal_id required")
	}
	if p.Title == "" {
		return PlanPhase{}, fmt.Errorf("store: insert plan phase: title required")
	}
	if p.Status == "" {
		p.Status = StatusActive
	}
	if p.CreatedAt == "" {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO plan_phases(id, goal_id, title, body, ord, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, p.ID, p.GoalID, p.Title, p.Body, p.Ord, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return PlanPhase{}, fmt.Errorf("store: insert plan phase: %w", err)
	}
	return s.GetPlanPhase(p.ID)
}

// UpsertPlanPhase inserts or updates a phase by id. created_at is preserved on conflict.
func (s *Store) UpsertPlanPhase(p PlanPhase) (PlanPhase, error) {
	now := nowRFC3339()
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.GoalID == "" {
		return PlanPhase{}, fmt.Errorf("store: upsert plan phase: goal_id required")
	}
	if p.Title == "" {
		return PlanPhase{}, fmt.Errorf("store: upsert plan phase: title required")
	}
	if p.Status == "" {
		p.Status = StatusActive
	}
	if p.CreatedAt == "" {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO plan_phases(id, goal_id, title, body, ord, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			goal_id = excluded.goal_id,
			title = excluded.title,
			body = excluded.body,
			ord = excluded.ord,
			status = excluded.status,
			updated_at = excluded.updated_at
	`, p.ID, p.GoalID, p.Title, p.Body, p.Ord, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return PlanPhase{}, fmt.Errorf("store: upsert plan phase: %w", err)
	}
	return s.GetPlanPhase(p.ID)
}

// GetPlanPhase loads a phase by id.
func (s *Store) GetPlanPhase(id string) (PlanPhase, error) {
	var p PlanPhase
	err := s.db.QueryRow(`
		SELECT id, goal_id, title, body, ord, status, created_at, updated_at
		FROM plan_phases WHERE id = ?
	`, id).Scan(&p.ID, &p.GoalID, &p.Title, &p.Body, &p.Ord, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return PlanPhase{}, fmt.Errorf("store: plan phase %q: %w", id, err)
	}
	if err != nil {
		return PlanPhase{}, fmt.Errorf("store: get plan phase: %w", err)
	}
	return p, nil
}

// ListPlanPhasesByGoal returns ACTIVE phases for a goal ordered by ord.
func (s *Store) ListPlanPhasesByGoal(goalID string) ([]PlanPhase, error) {
	if goalID == "" {
		return nil, fmt.Errorf("store: list plan phases: goal_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, goal_id, title, body, ord, status, created_at, updated_at
		FROM plan_phases
		WHERE goal_id = ? AND status = ?
		ORDER BY ord ASC, id ASC
	`, goalID, StatusActive)
	if err != nil {
		return nil, fmt.Errorf("store: list plan phases: %w", err)
	}
	defer rows.Close()

	var out []PlanPhase
	for rows.Next() {
		var p PlanPhase
		if err := rows.Scan(&p.ID, &p.GoalID, &p.Title, &p.Body, &p.Ord, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan plan phase: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// InsertPlanScope inserts a scope row. Empty ID allocates a UUID.
func (s *Store) InsertPlanScope(sc PlanScope) (PlanScope, error) {
	now := nowRFC3339()
	if sc.ID == "" {
		sc.ID = uuid.NewString()
	}
	if sc.PhaseID == "" {
		return PlanScope{}, fmt.Errorf("store: insert plan scope: phase_id required")
	}
	if sc.Title == "" {
		return PlanScope{}, fmt.Errorf("store: insert plan scope: title required")
	}
	if sc.Status == "" {
		sc.Status = StatusActive
	}
	if sc.CreatedAt == "" {
		sc.CreatedAt = now
	}
	sc.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO plan_scopes(id, phase_id, title, body, ord, status, auto_replan_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, sc.ID, sc.PhaseID, sc.Title, sc.Body, sc.Ord, sc.Status, sc.AutoReplanCount, sc.CreatedAt, sc.UpdatedAt)
	if err != nil {
		return PlanScope{}, fmt.Errorf("store: insert plan scope: %w", err)
	}
	return s.GetPlanScope(sc.ID)
}

// UpsertPlanScope inserts or updates a scope by id. created_at is preserved on conflict.
func (s *Store) UpsertPlanScope(sc PlanScope) (PlanScope, error) {
	now := nowRFC3339()
	if sc.ID == "" {
		sc.ID = uuid.NewString()
	}
	if sc.PhaseID == "" {
		return PlanScope{}, fmt.Errorf("store: upsert plan scope: phase_id required")
	}
	if sc.Title == "" {
		return PlanScope{}, fmt.Errorf("store: upsert plan scope: title required")
	}
	if sc.Status == "" {
		sc.Status = StatusActive
	}
	if sc.CreatedAt == "" {
		sc.CreatedAt = now
	}
	sc.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO plan_scopes(id, phase_id, title, body, ord, status, auto_replan_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			phase_id = excluded.phase_id,
			title = excluded.title,
			body = excluded.body,
			ord = excluded.ord,
			status = excluded.status,
			auto_replan_count = excluded.auto_replan_count,
			updated_at = excluded.updated_at
	`, sc.ID, sc.PhaseID, sc.Title, sc.Body, sc.Ord, sc.Status, sc.AutoReplanCount, sc.CreatedAt, sc.UpdatedAt)
	if err != nil {
		return PlanScope{}, fmt.Errorf("store: upsert plan scope: %w", err)
	}
	return s.GetPlanScope(sc.ID)
}

// GetPlanScope loads a scope by id.
func (s *Store) GetPlanScope(id string) (PlanScope, error) {
	var sc PlanScope
	err := s.db.QueryRow(`
		SELECT id, phase_id, title, body, ord, status, auto_replan_count, created_at, updated_at
		FROM plan_scopes WHERE id = ?
	`, id).Scan(
		&sc.ID, &sc.PhaseID, &sc.Title, &sc.Body, &sc.Ord, &sc.Status, &sc.AutoReplanCount,
		&sc.CreatedAt, &sc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return PlanScope{}, fmt.Errorf("store: plan scope %q: %w", id, err)
	}
	if err != nil {
		return PlanScope{}, fmt.Errorf("store: get plan scope: %w", err)
	}
	return sc, nil
}

// UpdatePlanScopeBody sets the coarse body on a scope.
func (s *Store) UpdatePlanScopeBody(scopeID, body string) error {
	if scopeID == "" {
		return fmt.Errorf("store: update plan scope body: scope_id required")
	}
	now := nowRFC3339()
	res, err := s.db.Exec(`UPDATE plan_scopes SET body = ?, updated_at = ? WHERE id = ?`, body, now, scopeID)
	if err != nil {
		return fmt.Errorf("store: update plan scope body: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: update plan scope body rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: plan scope %q: %w", scopeID, sql.ErrNoRows)
	}
	return nil
}

// IncrementAutoReplanCount atomically increments plan_scopes.auto_replan_count by 1
// and returns the new value. Fails if the scope is missing.
func (s *Store) IncrementAutoReplanCount(scopeID string) (int, error) {
	if scopeID == "" {
		return 0, fmt.Errorf("store: increment auto_replan_count: scope_id required")
	}
	now := nowRFC3339()
	res, err := s.db.Exec(`
		UPDATE plan_scopes SET auto_replan_count = auto_replan_count + 1, updated_at = ?
		WHERE id = ?
	`, now, scopeID)
	if err != nil {
		return 0, fmt.Errorf("store: increment auto_replan_count: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("store: increment auto_replan_count rows: %w", err)
	}
	if n == 0 {
		return 0, fmt.Errorf("store: plan scope %q: %w", scopeID, sql.ErrNoRows)
	}
	sc, err := s.GetPlanScope(scopeID)
	if err != nil {
		return 0, err
	}
	return sc.AutoReplanCount, nil
}

// AckAutoReplan resets plan_scopes.auto_replan_count to 0. Fails if the scope is missing.
func (s *Store) AckAutoReplan(scopeID string) error {
	if scopeID == "" {
		return fmt.Errorf("store: ack auto_replan: scope_id required")
	}
	now := nowRFC3339()
	res, err := s.db.Exec(`
		UPDATE plan_scopes SET auto_replan_count = 0, updated_at = ?
		WHERE id = ?
	`, now, scopeID)
	if err != nil {
		return fmt.Errorf("store: ack auto_replan: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: ack auto_replan rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: plan scope %q: %w", scopeID, sql.ErrNoRows)
	}
	return nil
}

// ListPlanScopesByPhase returns ACTIVE scopes for a phase ordered by ord.
func (s *Store) ListPlanScopesByPhase(phaseID string) ([]PlanScope, error) {
	if phaseID == "" {
		return nil, fmt.Errorf("store: list plan scopes: phase_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, phase_id, title, body, ord, status, auto_replan_count, created_at, updated_at
		FROM plan_scopes
		WHERE phase_id = ? AND status = ?
		ORDER BY ord ASC, id ASC
	`, phaseID, StatusActive)
	if err != nil {
		return nil, fmt.Errorf("store: list plan scopes: %w", err)
	}
	defer rows.Close()

	var out []PlanScope
	for rows.Next() {
		var sc PlanScope
		if err := rows.Scan(
			&sc.ID, &sc.PhaseID, &sc.Title, &sc.Body, &sc.Ord, &sc.Status, &sc.AutoReplanCount,
			&sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan plan scope: %w", err)
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

// ListPlanScopesByGoal returns ACTIVE scopes for a goal ordered by phase.ord, scope.ord.
func (s *Store) ListPlanScopesByGoal(goalID string) ([]PlanScope, error) {
	if goalID == "" {
		return nil, fmt.Errorf("store: list plan scopes by goal: goal_id required")
	}
	rows, err := s.db.Query(`
		SELECT s.id, s.phase_id, s.title, s.body, s.ord, s.status, s.auto_replan_count, s.created_at, s.updated_at
		FROM plan_scopes s
		JOIN plan_phases p ON p.id = s.phase_id
		WHERE p.goal_id = ? AND s.status = ? AND p.status = ?
		ORDER BY p.ord ASC, s.ord ASC, s.id ASC
	`, goalID, StatusActive, StatusActive)
	if err != nil {
		return nil, fmt.Errorf("store: list plan scopes by goal: %w", err)
	}
	defer rows.Close()

	var out []PlanScope
	for rows.Next() {
		var sc PlanScope
		if err := rows.Scan(
			&sc.ID, &sc.PhaseID, &sc.Title, &sc.Body, &sc.Ord, &sc.Status, &sc.AutoReplanCount,
			&sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan plan scope: %w", err)
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

// GoalIDForScope resolves the owning goal_id for a scope via its phase.
func (s *Store) GoalIDForScope(scopeID string) (string, error) {
	if scopeID == "" {
		return "", fmt.Errorf("store: goal id for scope: scope_id required")
	}
	var goalID string
	err := s.db.QueryRow(`
		SELECT p.goal_id
		FROM plan_scopes s
		JOIN plan_phases p ON p.id = s.phase_id
		WHERE s.id = ?
	`, scopeID).Scan(&goalID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("store: plan scope %q: %w", scopeID, err)
	}
	if err != nil {
		return "", fmt.Errorf("store: goal id for scope: %w", err)
	}
	return goalID, nil
}

// InsertScopeDeepPlan inserts a deep-plan revision. Empty ID allocates a UUID.
func (s *Store) InsertScopeDeepPlan(d ScopeDeepPlan) (ScopeDeepPlan, error) {
	now := nowRFC3339()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.ScopeID == "" {
		return ScopeDeepPlan{}, fmt.Errorf("store: insert scope deep plan: scope_id required")
	}
	if d.ContentJSON == "" {
		d.ContentJSON = "{}"
	}
	if d.Status == "" {
		d.Status = StatusActive
	}
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO scope_deep_plans(id, scope_id, content_json, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, d.ID, d.ScopeID, d.ContentJSON, d.Status, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return ScopeDeepPlan{}, fmt.Errorf("store: insert scope deep plan: %w", err)
	}
	return s.GetScopeDeepPlan(d.ID)
}

// UpsertScopeDeepPlan inserts or updates a deep plan by id. created_at is preserved on conflict.
func (s *Store) UpsertScopeDeepPlan(d ScopeDeepPlan) (ScopeDeepPlan, error) {
	now := nowRFC3339()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.ScopeID == "" {
		return ScopeDeepPlan{}, fmt.Errorf("store: upsert scope deep plan: scope_id required")
	}
	if d.ContentJSON == "" {
		d.ContentJSON = "{}"
	}
	if d.Status == "" {
		d.Status = StatusActive
	}
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO scope_deep_plans(id, scope_id, content_json, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			scope_id = excluded.scope_id,
			content_json = excluded.content_json,
			status = excluded.status,
			updated_at = excluded.updated_at
	`, d.ID, d.ScopeID, d.ContentJSON, d.Status, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return ScopeDeepPlan{}, fmt.Errorf("store: upsert scope deep plan: %w", err)
	}
	return s.GetScopeDeepPlan(d.ID)
}

// GetScopeDeepPlan loads a deep-plan revision by id.
func (s *Store) GetScopeDeepPlan(id string) (ScopeDeepPlan, error) {
	var d ScopeDeepPlan
	err := s.db.QueryRow(`
		SELECT id, scope_id, content_json, status, created_at, updated_at
		FROM scope_deep_plans WHERE id = ?
	`, id).Scan(&d.ID, &d.ScopeID, &d.ContentJSON, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return ScopeDeepPlan{}, fmt.Errorf("store: scope deep plan %q: %w", id, err)
	}
	if err != nil {
		return ScopeDeepPlan{}, fmt.Errorf("store: get scope deep plan: %w", err)
	}
	return d, nil
}

// GetActiveScopeDeepPlan returns the ACTIVE deep plan for a scope, if any.
func (s *Store) GetActiveScopeDeepPlan(scopeID string) (ScopeDeepPlan, error) {
	if scopeID == "" {
		return ScopeDeepPlan{}, fmt.Errorf("store: get active deep plan: scope_id required")
	}
	var d ScopeDeepPlan
	err := s.db.QueryRow(`
		SELECT id, scope_id, content_json, status, created_at, updated_at
		FROM scope_deep_plans
		WHERE scope_id = ? AND status = ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, scopeID, StatusActive).Scan(&d.ID, &d.ScopeID, &d.ContentJSON, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return ScopeDeepPlan{}, fmt.Errorf("store: active deep plan for scope %q: %w", scopeID, err)
	}
	if err != nil {
		return ScopeDeepPlan{}, fmt.Errorf("store: get active deep plan: %w", err)
	}
	return d, nil
}

// SupersedeActiveScopeDeepPlans marks all ACTIVE deep plans for a scope SUPERSEDED.
func (s *Store) SupersedeActiveScopeDeepPlans(scopeID string) (int64, error) {
	if scopeID == "" {
		return 0, fmt.Errorf("store: supersede deep plans: scope_id required")
	}
	now := nowRFC3339()
	res, err := s.db.Exec(`
		UPDATE scope_deep_plans SET status = ?, updated_at = ?
		WHERE scope_id = ? AND status = ?
	`, StatusSuperseded, now, scopeID, StatusActive)
	if err != nil {
		return 0, fmt.Errorf("store: supersede deep plans: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("store: supersede deep plans rows: %w", err)
	}
	return n, nil
}

// CountScopeDeepPlans returns how many deep-plan rows exist for a scope (any status).
func (s *Store) CountScopeDeepPlans(scopeID string) (int, error) {
	if scopeID == "" {
		return 0, fmt.Errorf("store: count deep plans: scope_id required")
	}
	var n int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM scope_deep_plans WHERE scope_id = ?`, scopeID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count deep plans: %w", err)
	}
	return n, nil
}

// UpsertGoalPlanState inserts or updates the current-scope pointer for a goal.
func (s *Store) UpsertGoalPlanState(st GoalPlanState) (GoalPlanState, error) {
	if st.GoalID == "" {
		return GoalPlanState{}, fmt.Errorf("store: upsert goal plan state: goal_id required")
	}
	now := nowRFC3339()
	st.UpdatedAt = now
	_, err := s.db.Exec(`
		INSERT INTO goal_plan_state(goal_id, current_scope_id, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(goal_id) DO UPDATE SET
			current_scope_id = excluded.current_scope_id,
			updated_at = excluded.updated_at
	`, st.GoalID, nullStr(st.CurrentScopeID), st.UpdatedAt)
	if err != nil {
		return GoalPlanState{}, fmt.Errorf("store: upsert goal plan state: %w", err)
	}
	return s.GetGoalPlanState(st.GoalID)
}

// EnsureGoalPlanState creates a goal_plan_state row with NULL current if missing.
func (s *Store) EnsureGoalPlanState(goalID string) (GoalPlanState, error) {
	if goalID == "" {
		return GoalPlanState{}, fmt.Errorf("store: ensure goal plan state: goal_id required")
	}
	existing, err := s.GetGoalPlanState(goalID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return GoalPlanState{}, err
	}
	return s.UpsertGoalPlanState(GoalPlanState{GoalID: goalID, CurrentScopeID: nil})
}

// ListAllPlanPhases returns every plan phase row (all statuses).
func (s *Store) ListAllPlanPhases() ([]PlanPhase, error) {
	rows, err := s.db.Query(`
		SELECT id, goal_id, title, body, ord, status, created_at, updated_at
		FROM plan_phases
		ORDER BY ord ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all plan phases: %w", err)
	}
	defer rows.Close()
	var out []PlanPhase
	for rows.Next() {
		var p PlanPhase
		if err := rows.Scan(&p.ID, &p.GoalID, &p.Title, &p.Body, &p.Ord, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan plan phase: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ListAllPlanScopes returns every plan scope row (all statuses).
func (s *Store) ListAllPlanScopes() ([]PlanScope, error) {
	rows, err := s.db.Query(`
		SELECT id, phase_id, title, body, ord, status, auto_replan_count, created_at, updated_at
		FROM plan_scopes
		ORDER BY ord ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all plan scopes: %w", err)
	}
	defer rows.Close()
	var out []PlanScope
	for rows.Next() {
		var sc PlanScope
		if err := rows.Scan(
			&sc.ID, &sc.PhaseID, &sc.Title, &sc.Body, &sc.Ord, &sc.Status, &sc.AutoReplanCount,
			&sc.CreatedAt, &sc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan plan scope: %w", err)
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

// ListAllScopeDeepPlans returns every deep-plan revision (ACTIVE and SUPERSEDED).
func (s *Store) ListAllScopeDeepPlans() ([]ScopeDeepPlan, error) {
	rows, err := s.db.Query(`
		SELECT id, scope_id, content_json, status, created_at, updated_at
		FROM scope_deep_plans
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all scope deep plans: %w", err)
	}
	defer rows.Close()
	var out []ScopeDeepPlan
	for rows.Next() {
		var d ScopeDeepPlan
		if err := rows.Scan(&d.ID, &d.ScopeID, &d.ContentJSON, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan scope deep plan: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ListAllGoalPlanStates returns every goal_plan_state row.
func (s *Store) ListAllGoalPlanStates() ([]GoalPlanState, error) {
	rows, err := s.db.Query(`
		SELECT goal_id, current_scope_id, updated_at FROM goal_plan_state
		ORDER BY goal_id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all goal plan states: %w", err)
	}
	defer rows.Close()
	var out []GoalPlanState
	for rows.Next() {
		var st GoalPlanState
		var cur sql.NullString
		if err := rows.Scan(&st.GoalID, &cur, &st.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan goal plan state: %w", err)
		}
		st.CurrentScopeID = nullStrPtr(cur)
		out = append(out, st)
	}
	return out, rows.Err()
}

// GetGoalPlanState loads the plan state for a goal.
func (s *Store) GetGoalPlanState(goalID string) (GoalPlanState, error) {
	if goalID == "" {
		return GoalPlanState{}, fmt.Errorf("store: get goal plan state: goal_id required")
	}
	var st GoalPlanState
	var cur sql.NullString
	err := s.db.QueryRow(`
		SELECT goal_id, current_scope_id, updated_at FROM goal_plan_state WHERE goal_id = ?
	`, goalID).Scan(&st.GoalID, &cur, &st.UpdatedAt)
	if err == sql.ErrNoRows {
		return GoalPlanState{}, fmt.Errorf("store: goal plan state %q: %w", goalID, err)
	}
	if err != nil {
		return GoalPlanState{}, fmt.Errorf("store: get goal plan state: %w", err)
	}
	st.CurrentScopeID = nullStrPtr(cur)
	return st, nil
}
