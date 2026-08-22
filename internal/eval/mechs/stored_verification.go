package mechs

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
)

type StoredVerification struct{}

func (StoredVerification) ID() string { return idStoredVerification }

func (StoredVerification) Run(ctx context.Context, taskID string, svc *domain.Service) (bool, string, string, error) {
	gateOK, reason, err := svc.CheckVerificationGate(ctx, taskID)
	if err != nil {
		return false, "", "", err
	}
	if gateOK {
		return true, "verification gate satisfied", "", nil
	}
	hasRow, err := svc.HasVerificationOutcome(ctx, taskID)
	if err != nil {
		return false, "", "", err
	}
	if hasRow {
		return true, "stored verification outcome present", "", nil
	}
	summary := strings.TrimSpace(reason)
	if summary == "" {
		summary = "verification gate not satisfied"
	}
	return false, summary, "", nil
}
