package loop

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
)

func buildVerificationCycle(ctx context.Context, dom *domain.Service, inputs deliberation.PolicyInputs, taskID string) (*VerificationCycleStatus, error) {
	out := &VerificationCycleStatus{
		ExecutePending:         inputs.ExecutePending,
		TestPending:            inputs.TestPending,
		VerificationIncomplete: inputs.VerificationIncomplete,
		EvaluationPending:      inputs.EvaluationPending,
		ReflectPending:         inputs.ReflectPending,
	}
	out.IncompleteReason = verificationCycleIncompleteReason(inputs)
	if dom != nil && taskID != "" {
		reg, err := dom.DetectAnyTestRegression(ctx, taskID)
		if err != nil {
			return nil, err
		}
		out.RegressionDetected = reg.Detected
		out.TestName = reg.TestName
	}
	return out, nil
}

func verificationCycleIncompleteReason(inputs deliberation.PolicyInputs) string {
	var parts []string
	if inputs.TestPending {
		parts = append(parts, "test_pending")
	}
	if inputs.VerificationIncomplete {
		parts = append(parts, "verification_incomplete")
	}
	if inputs.EvaluationPending {
		parts = append(parts, "evaluation_pending")
	}
	return strings.Join(parts, "; ")
}
