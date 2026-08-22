package mechs

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
)

// StoredTest mechanism (file storedtest.go — Go excludes *_test.go from product builds).
type StoredTest struct{}

func (StoredTest) ID() string { return idStoredTest }

func (StoredTest) Run(ctx context.Context, taskID string, svc *domain.Service) (bool, string, string, error) {
	ok, err := svc.HasTestOutcomeSinceLatestChange(ctx, taskID)
	if err != nil {
		return false, "", "", err
	}
	if ok {
		return true, "stored test outcome since latest change", "", nil
	}
	return false, "no stored test outcome since latest change", "", nil
}
