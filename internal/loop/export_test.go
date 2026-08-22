package loop

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// BuildRiskHintsSectionForTest exposes buildRiskHintsSection to external tests.
func BuildRiskHintsSectionForTest(ctx context.Context, dom *domain.Service, st *store.Store, taskID, freshness string) (RiskHintsSection, error) {
	return buildRiskHintsSection(ctx, dom, st, taskID, freshness)
}
