package eval

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// Mechanism is one verification/evaluation check (C40/C43 additive contract).
type Mechanism interface {
	ID() string
	Run(ctx context.Context, in EvalInput) (EvalResult, error)
}

// EvalInput carries task context; mechanisms read store state via domain.Service.
type EvalInput struct {
	TaskID  string
	Service *domain.Service
	Rules   *RulesFile
}

// EvalResult is an in-memory mechanism outcome (no new outcome_results kind).
type EvalResult struct {
	MechanismID string `json:"mechanism_id"`
	Passed      bool   `json:"passed"`
	Summary     string `json:"summary"`
	DetailsJSON string `json:"details_json,omitempty"`
	RecordedAt  string `json:"recorded_at"`
}

// RunOptions selects which registered mechanisms to run.
type RunOptions struct {
	MechanismIDs []string
	Rules        *RulesFile
	Root         string
	Store        *store.Store
}
