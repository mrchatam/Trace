package mechs

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
)

type StoredEvaluation struct{}

func (StoredEvaluation) ID() string { return idStoredEvaluation }

func (StoredEvaluation) Run(ctx context.Context, taskID string, svc *domain.Service) (bool, string, string, error) {
	ok, err := svc.HasComputedEvaluation(ctx, taskID)
	if err != nil {
		return false, "", "", err
	}
	if ok {
		return true, "computed evaluation outcome present", "", nil
	}
	return false, "no computed evaluation outcome", "", nil
}
