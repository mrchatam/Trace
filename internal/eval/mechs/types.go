package mechs

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
)

const (
	idStoredTest             = "stored_test"
	idStoredVerification     = "stored_verification"
	idStoredEvaluation       = "stored_evaluation"
	idArchitecturalInvariant = "architectural_invariant"
)

// ID constants for registration (mirrors eval package).
const (
	StoredTestID             = idStoredTest
	StoredVerificationID     = idStoredVerification
	StoredEvaluationID       = idStoredEvaluation
	ArchitecturalInvariantID = idArchitecturalInvariant
)

// Builtin is a built-in mechanism implementation (domain-only; registered by eval).
type Builtin interface {
	ID() string
	Run(ctx context.Context, taskID string, svc *domain.Service) (passed bool, summary, detailsJSON string, err error)
}

// All returns the four locked built-in mechanisms.
func All() []Builtin {
	return []Builtin{
		StoredTest{},
		StoredVerification{},
		StoredEvaluation{},
		ArchitecturalInvariant{},
	}
}
