package eval

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	resultsDefaultLimit = 32
	resultsMaxLimit     = 64
)

// ResultRow is one stored eval-relevant outcome with a derived mechanism_id.
type ResultRow struct {
	ID             string `json:"id"`
	TaskID         string `json:"task_id"`
	MechanismID    string `json:"mechanism_id"`
	Kind           string `json:"kind"`
	Passed         *bool  `json:"passed,omitempty"`
	Summary        string `json:"summary,omitempty"`
	ScoresJSON     string `json:"scores_json,omitempty"`
	ComparisonJSON string `json:"comparison_json,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// ListResults returns stored test, verification, and evaluation outcomes for a task.
// Stored outcomes only — use RunAll for live mechanism checks.
func ListResults(ctx context.Context, svc *domain.Service, taskID string) ([]ResultRow, error) {
	return ListResultsWithLimit(ctx, svc, taskID, resultsDefaultLimit)
}

// ListResultsWithLimit is ListResults with an explicit row cap (default 32, max 64).
func ListResultsWithLimit(ctx context.Context, svc *domain.Service, taskID string, limit int) ([]ResultRow, error) {
	if svc == nil {
		return nil, &domain.ErrValidation{Msg: "service is required"}
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &domain.ErrValidation{Msg: "task_id is required"}
	}
	if limit <= 0 {
		limit = resultsDefaultLimit
	}
	if limit > resultsMaxLimit {
		limit = resultsMaxLimit
	}
	rows, err := svc.ListEvalOutcomeResults(ctx, taskID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]ResultRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapOutcomeToResultRow(row))
	}
	return out, nil
}

func mapOutcomeToResultRow(o store.OutcomeResult) ResultRow {
	row := ResultRow{
		ID:             o.ID,
		TaskID:         o.TaskID,
		Kind:           o.Kind,
		Summary:        o.Summary,
		ScoresJSON:     o.ScoresJSON,
		ComparisonJSON: o.ComparisonJSON,
		CreatedAt:      o.CreatedAt,
	}
	switch o.Kind {
	case store.OutcomeKindTest:
		row.MechanismID = MechanismStoredTest
		row.Passed = testPassedPtr(o.TestStatus)
	case store.OutcomeKindVerification:
		row.MechanismID = MechanismStoredVerification
		row.Passed = verificationPassedPtr(o.VerificationStatus)
	case store.OutcomeKindEvaluation:
		row.MechanismID = MechanismStoredEvaluation
	}
	return row
}

func testPassedPtr(status string) *bool {
	switch status {
	case store.TestStatusPass:
		v := true
		return &v
	case store.TestStatusFail, store.TestStatusError:
		v := false
		return &v
	default:
		return nil
	}
}

func verificationPassedPtr(status string) *bool {
	switch status {
	case store.VerificationStatusVerified:
		v := true
		return &v
	case store.VerificationStatusFailed:
		v := false
		return &v
	default:
		return nil
	}
}
