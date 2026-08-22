package domain_test

import (
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
)

func TestPrematureImplementation_Code(t *testing.T) {
	err := &domain.PrematureImplementation{
		TaskID:           "task-1",
		For:              "edit",
		RecommendedPhase: "INVESTIGATE",
		ReasonCode:       "blocking_uncertainty",
	}
	if got := err.Code(); got != "premature_implementation" {
		t.Fatalf("Code() = %q want premature_implementation", got)
	}
	if err.Error() == "" {
		t.Fatal("Error() must be non-empty")
	}
}
