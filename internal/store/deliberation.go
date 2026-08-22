package store

import (
	"database/sql"
	"fmt"
)

// DeliberationState is one controller row per seed task_id.
type DeliberationState struct {
	TaskID                  string
	GoalID                  string
	CurrentPhase            string
	HopCount                int
	LastPhase               string
	PlanCritiqued           bool
	Stopped                 bool
	StopReason              string
	ConsecutiveEmptyApplies int
	UpdatedAt               string
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// UpsertDeliberationState inserts or replaces the row keyed by task_id.
func (s *Store) UpsertDeliberationState(st DeliberationState) (DeliberationState, error) {
	if st.TaskID == "" {
		return DeliberationState{}, fmt.Errorf("store: upsert deliberation state: task_id required")
	}
	if st.GoalID == "" {
		return DeliberationState{}, fmt.Errorf("store: upsert deliberation state: goal_id required")
	}
	if st.CurrentPhase == "" {
		st.CurrentPhase = "ORIENT"
	}
	st.UpdatedAt = nowRFC3339()

	_, err := s.db.Exec(`
		INSERT INTO deliberation_state(
			task_id, goal_id, current_phase, hop_count, last_phase,
			plan_critiqued, stopped, stop_reason, consecutive_empty_applies, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(task_id) DO UPDATE SET
			goal_id = excluded.goal_id,
			current_phase = excluded.current_phase,
			hop_count = excluded.hop_count,
			last_phase = excluded.last_phase,
			plan_critiqued = excluded.plan_critiqued,
			stopped = excluded.stopped,
			stop_reason = excluded.stop_reason,
			consecutive_empty_applies = excluded.consecutive_empty_applies,
			updated_at = excluded.updated_at
	`, st.TaskID, st.GoalID, st.CurrentPhase, st.HopCount, st.LastPhase,
		boolToInt(st.PlanCritiqued), boolToInt(st.Stopped), st.StopReason,
		st.ConsecutiveEmptyApplies, st.UpdatedAt)
	if err != nil {
		return DeliberationState{}, fmt.Errorf("store: upsert deliberation state: %w", err)
	}
	return s.GetDeliberationState(st.TaskID)
}

// ListAllDeliberationStates returns every deliberation_state row.
func (s *Store) ListAllDeliberationStates() ([]DeliberationState, error) {
	rows, err := s.db.Query(`
		SELECT task_id, goal_id, current_phase, hop_count, last_phase,
			plan_critiqued, stopped, stop_reason, consecutive_empty_applies, updated_at
		FROM deliberation_state
		ORDER BY task_id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all deliberation states: %w", err)
	}
	defer rows.Close()
	var out []DeliberationState
	for rows.Next() {
		var st DeliberationState
		var planCritiqued, stopped int
		if err := rows.Scan(
			&st.TaskID, &st.GoalID, &st.CurrentPhase, &st.HopCount, &st.LastPhase,
			&planCritiqued, &stopped, &st.StopReason, &st.ConsecutiveEmptyApplies, &st.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan deliberation state: %w", err)
		}
		st.PlanCritiqued = planCritiqued != 0
		st.Stopped = stopped != 0
		out = append(out, st)
	}
	if out == nil {
		out = []DeliberationState{}
	}
	return out, rows.Err()
}

// GetDeliberationState loads the row for task_id.
func (s *Store) GetDeliberationState(taskID string) (DeliberationState, error) {
	if taskID == "" {
		return DeliberationState{}, fmt.Errorf("store: get deliberation state: task_id required")
	}
	var st DeliberationState
	var planCritiqued, stopped int
	err := s.db.QueryRow(`
		SELECT task_id, goal_id, current_phase, hop_count, last_phase,
			plan_critiqued, stopped, stop_reason, consecutive_empty_applies, updated_at
		FROM deliberation_state WHERE task_id = ?
	`, taskID).Scan(
		&st.TaskID, &st.GoalID, &st.CurrentPhase, &st.HopCount, &st.LastPhase,
		&planCritiqued, &stopped, &st.StopReason, &st.ConsecutiveEmptyApplies, &st.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return DeliberationState{}, fmt.Errorf("store: deliberation state %q: %w", taskID, err)
	}
	if err != nil {
		return DeliberationState{}, fmt.Errorf("store: get deliberation state: %w", err)
	}
	st.PlanCritiqued = planCritiqued != 0
	st.Stopped = stopped != 0
	return st, nil
}
