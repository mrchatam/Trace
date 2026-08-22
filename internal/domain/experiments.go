package domain

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// ExperimentInput creates a thin experiment record (no runner side effects).
type ExperimentInput struct {
	TaskID            string
	Label             string
	HypothesisSummary string
}

var experimentStatusOrder = map[string]int{
	store.ExperimentStatusPlanned:   0,
	store.ExperimentStatusRunning:   1,
	store.ExperimentStatusCompleted: 2,
}

func normalizeExperimentStatus(status string) (string, error) {
	status = strings.TrimSpace(status)
	switch status {
	case store.ExperimentStatusPlanned, store.ExperimentStatusRunning, store.ExperimentStatusCompleted:
		return status, nil
	case "":
		return store.ExperimentStatusPlanned, nil
	default:
		return "", &ErrValidation{Msg: "invalid experiment status: " + status}
	}
}

func validateExperimentForwardTransition(from, to string) error {
	if from == to {
		return nil
	}
	fromOrd, okFrom := experimentStatusOrder[from]
	toOrd, okTo := experimentStatusOrder[to]
	if !okFrom || !okTo {
		return &ErrValidation{Msg: "invalid experiment status transition"}
	}
	if toOrd != fromOrd+1 {
		return &ErrValidation{Msg: "experiment status may only advance planned→running→completed"}
	}
	return nil
}

// CreateExperiment allocates an experiment with status planned and no outcome link.
func (s *Service) CreateExperiment(ctx context.Context, in ExperimentInput) (store.Experiment, error) {
	_ = ctx
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return store.Experiment{}, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return store.Experiment{}, err
	}
	e, err := s.store.UpsertExperiment(store.Experiment{
		ID:                uuid.NewString(),
		TaskID:            taskID,
		Label:             strings.TrimSpace(in.Label),
		HypothesisSummary: strings.TrimSpace(in.HypothesisSummary),
		Status:            store.ExperimentStatusPlanned,
		OutcomeResultID:   "",
	})
	if err != nil {
		return store.Experiment{}, err
	}
	return e, nil
}

// SetExperimentStatus validates enum and allows forward-only lifecycle transitions.
func (s *Service) SetExperimentStatus(ctx context.Context, id, status string) error {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return &ErrValidation{Msg: "experiment id is required"}
	}
	next, err := normalizeExperimentStatus(status)
	if err != nil {
		return err
	}
	cur, err := s.store.GetExperiment(id)
	if err != nil {
		return err
	}
	if err := validateExperimentForwardTransition(cur.Status, next); err != nil {
		return err
	}
	if cur.Status == next {
		return nil
	}
	cur.Status = next
	_, err = s.store.UpsertExperiment(cur)
	return err
}

// LinkExperimentOutcome sets outcome_result_id when the outcome row exists.
func (s *Service) LinkExperimentOutcome(ctx context.Context, experimentID, outcomeResultID string) error {
	_ = ctx
	experimentID = strings.TrimSpace(experimentID)
	outcomeResultID = strings.TrimSpace(outcomeResultID)
	if experimentID == "" {
		return &ErrValidation{Msg: "experiment id is required"}
	}
	if outcomeResultID == "" {
		return &ErrValidation{Msg: "outcome_result_id is required"}
	}
	e, err := s.store.GetExperiment(experimentID)
	if err != nil {
		return err
	}
	if _, err := s.store.GetOutcomeResult(outcomeResultID); err != nil {
		return err
	}
	e.OutcomeResultID = outcomeResultID
	_, err = s.store.UpsertExperiment(e)
	return err
}

// GetExperiment loads an experiment by id.
func (s *Service) GetExperiment(ctx context.Context, id string) (store.Experiment, error) {
	_ = ctx
	return s.store.GetExperiment(strings.TrimSpace(id))
}

// ListExperimentsByTaskID returns experiments for a task.
func (s *Service) ListExperimentsByTaskID(ctx context.Context, taskID string) ([]store.Experiment, error) {
	_ = ctx
	return s.store.ListExperimentsByTaskID(strings.TrimSpace(taskID))
}
