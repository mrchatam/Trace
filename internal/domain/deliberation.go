package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/store"
)

// ApplyDeliberationTransition selects the next phase, upserts deliberation_state,
// and appends a deliberation.transition event on the seed task. Missing state
// starts at ORIENT / hop_count 0.
func (s *Service) ApplyDeliberationTransition(ctx context.Context, taskID, goalID string, inputs deliberation.PolicyInputs) (deliberation.State, store.Event, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	goalID = strings.TrimSpace(goalID)
	if taskID == "" {
		return deliberation.State{}, store.Event{}, &ErrValidation{Msg: "task_id is required"}
	}
	if goalID == "" {
		return deliberation.State{}, store.Event{}, &ErrValidation{Msg: "goal_id is required"}
	}
	task, err := s.store.GetTask(taskID)
	if err != nil {
		return deliberation.State{}, store.Event{}, err
	}
	if task.GoalID == nil || strings.TrimSpace(*task.GoalID) == "" {
		return deliberation.State{}, store.Event{}, &ErrValidation{Msg: "task has no goal_id"}
	}
	if strings.TrimSpace(*task.GoalID) != goalID {
		return deliberation.State{}, store.Event{}, &ErrValidation{Msg: "goal_id does not match task"}
	}

	cur, err := s.store.GetDeliberationState(taskID)
	var dState deliberation.State
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return deliberation.State{}, store.Event{}, err
		}
		dState = deliberation.InitialState(taskID, goalID)
	} else {
		dState = deliberationFromStore(cur)
		dState.GoalID = goalID
	}

	next, payload := deliberation.ApplyTransition(dState, inputs)
	if _, err := s.store.UpsertDeliberationState(storeFromDeliberation(next)); err != nil {
		return deliberation.State{}, store.Event{}, err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return deliberation.State{}, store.Event{}, err
	}
	ev, err := s.store.AppendEvent(store.Event{
		Type:        EventDeliberationTransition,
		EntityType:  EntityTask,
		EntityID:    taskID,
		PayloadJSON: string(raw),
	})
	if err != nil {
		return deliberation.State{}, store.Event{}, err
	}
	persisted, err := s.store.GetDeliberationState(taskID)
	if err != nil {
		return deliberation.State{}, store.Event{}, err
	}
	return deliberationFromStore(persisted), ev, nil
}

// ResetDeliberationState clears sticky STOP and hop/empty counters so the gate can
// recover. Preserves PlanCritiqued; sets CurrentPhase to EXECUTE.
func (s *Service) ResetDeliberationState(ctx context.Context, taskID string) (deliberation.State, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return deliberation.State{}, &ErrValidation{Msg: "task_id is required"}
	}
	cur, err := s.store.GetDeliberationState(taskID)
	if err != nil {
		return deliberation.State{}, err
	}
	next := deliberationFromStore(cur)
	next.Stopped = false
	next.StopReason = ""
	next.HopCount = 0
	next.ConsecutiveEmptyApplies = 0
	next.CurrentPhase = deliberation.PhaseExecute
	if _, err := s.store.UpsertDeliberationState(storeFromDeliberation(next)); err != nil {
		return deliberation.State{}, err
	}
	persisted, err := s.store.GetDeliberationState(taskID)
	if err != nil {
		return deliberation.State{}, err
	}
	return deliberationFromStore(persisted), nil
}

// GetDeliberationState loads controller state for a seed task.
func (s *Service) GetDeliberationState(ctx context.Context, taskID string) (deliberation.State, error) {
	_ = ctx
	st, err := s.store.GetDeliberationState(strings.TrimSpace(taskID))
	if err != nil {
		return deliberation.State{}, err
	}
	return deliberationFromStore(st), nil
}

func storeFromDeliberation(st deliberation.State) store.DeliberationState {
	return store.DeliberationState{
		TaskID:                  st.TaskID,
		GoalID:                  st.GoalID,
		CurrentPhase:            string(st.CurrentPhase),
		HopCount:                st.HopCount,
		LastPhase:               string(st.LastPhase),
		PlanCritiqued:           st.PlanCritiqued,
		Stopped:                 st.Stopped,
		StopReason:              st.StopReason,
		ConsecutiveEmptyApplies: st.ConsecutiveEmptyApplies,
	}
}

func deliberationFromStore(st store.DeliberationState) deliberation.State {
	return deliberation.State{
		TaskID:                  st.TaskID,
		GoalID:                  st.GoalID,
		CurrentPhase:            deliberation.Phase(st.CurrentPhase),
		HopCount:                st.HopCount,
		LastPhase:               deliberation.Phase(st.LastPhase),
		PlanCritiqued:           st.PlanCritiqued,
		Stopped:                 st.Stopped,
		StopReason:              st.StopReason,
		ConsecutiveEmptyApplies: st.ConsecutiveEmptyApplies,
		UpdatedAt:               st.UpdatedAt,
	}
}
