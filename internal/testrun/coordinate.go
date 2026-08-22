package testrun

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// CoordinateTestRun returns a domain.CoordinateTestRunFunc wired to RunRelevantTests.
func CoordinateTestRun(st *store.Store, dom *domain.Service, extra Options) domain.CoordinateTestRunFunc {
	return func(ctx context.Context, taskID string, paths []string) ([]store.OutcomeResult, error) {
		opts := extra
		if len(paths) > 0 {
			opts.Paths = paths
		}
		return RunRelevantTests(ctx, st, dom, taskID, opts)
	}
}
