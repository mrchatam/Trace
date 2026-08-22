package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxImprovementSummaryBytes = 4096
	maxImprovementEvidenceIDs  = 32
)

// ImprovementInput records a first-class improvement row (C18).
type ImprovementInput struct {
	ID          string
	ChangeID    string
	TaskID      string
	Dimension   string
	Summary     string
	EvidenceIDs []string
	SourceType  string
	Confidence  float64
}

func marshalEvidenceIDs(ids []string) (string, error) {
	raw, err := json.Marshal(ids)
	if err != nil {
		return "", &ErrValidation{Msg: "evidence_ids must be a JSON array"}
	}
	return string(raw), nil
}

func validateImprovementEvidenceIDs(s *Service, ids []string) ([]string, error) {
	if len(ids) > maxImprovementEvidenceIDs {
		return nil, &ErrValidation{Msg: "evidence_ids exceeds 32 items"}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			return nil, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		if _, err := s.store.GetEvidence(id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

// RecordImprovement persists a queryable improvement row for a change.
func (s *Service) RecordImprovement(ctx context.Context, in ImprovementInput) (store.Improvement, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.Improvement{}, err
	}
	changeID := strings.TrimSpace(in.ChangeID)
	if changeID == "" {
		return store.Improvement{}, &ErrValidation{Msg: "change_id is required"}
	}
	chg, err := s.store.GetChange(changeID)
	if err != nil {
		return store.Improvement{}, err
	}
	if chg.TaskID != task.ID {
		return store.Improvement{}, &ErrValidation{Msg: "task_id must match change.task_id"}
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxImprovementSummaryBytes); err != nil {
		return store.Improvement{}, err
	}
	if strings.TrimSpace(in.Summary) == "" {
		return store.Improvement{}, &ErrValidation{Msg: "summary is required"}
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return store.Improvement{}, err
	}
	evIDs, err := validateImprovementEvidenceIDs(s, in.EvidenceIDs)
	if err != nil {
		return store.Improvement{}, err
	}
	evJSON, err := marshalEvidenceIDs(evIDs)
	if err != nil {
		return store.Improvement{}, err
	}

	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}

	row, err := s.store.UpsertImprovement(store.Improvement{
		ID:              id,
		ChangeID:        changeID,
		TaskID:          task.ID,
		Dimension:       strings.TrimSpace(in.Dimension),
		Summary:         in.Summary,
		EvidenceIDsJSON: evJSON,
		SourceType:      defaultSource(in.SourceType),
		Confidence:      in.Confidence,
	})
	if err != nil {
		return store.Improvement{}, err
	}
	for _, eid := range evIDs {
		if err := s.insertLink(EntityImprovement, row.ID, RelImprovementSupportedBy, EntityEvidence, eid, row.SourceType, row.Confidence); err != nil {
			return store.Improvement{}, err
		}
	}
	if err := s.appendCreated(EntityImprovement, row.ID, row.Summary); err != nil {
		return store.Improvement{}, err
	}
	return row, nil
}

// GetImprovement loads an improvement by id.
func (s *Service) GetImprovement(ctx context.Context, id string) (store.Improvement, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.Improvement{}, &ErrValidation{Msg: "improvement id is required"}
	}
	return s.store.GetImprovement(id)
}

// ListImprovementsByChangeID returns improvements for a change.
func (s *Service) ListImprovementsByChangeID(ctx context.Context, changeID string) ([]store.Improvement, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil, &ErrValidation{Msg: "change_id is required"}
	}
	if _, err := s.store.GetChange(changeID); err != nil {
		return nil, err
	}
	return s.store.ListImprovementsByChangeID(changeID)
}

// ListImprovementsByTaskID returns improvements for a task.
func (s *Service) ListImprovementsByTaskID(ctx context.Context, taskID string) ([]store.Improvement, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return nil, err
	}
	return s.store.ListImprovementsByTaskID(taskID)
}
