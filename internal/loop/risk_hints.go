package loop

import (
	"context"
	"fmt"
	"sort"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxRiskHintsCap      = 4
	manyPathsThreshold   = 8
	highChurnMinChanges  = 3
	riskHintBlocking     = "blocking_uncertainty"
	riskHintVerification = "missing_verification"
	riskHintManyPaths    = "many_paths"
	riskHintHighChurn    = "high_churn_path"
)

type RiskHint struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Detail   string `json:"detail"`
}

type RiskHintsSection struct {
	Freshness string     `json:"freshness"`
	Items     []RiskHint `json:"items"`
}

func buildRiskHintsSection(ctx context.Context, dom *domain.Service, st *store.Store, taskID, freshness string) (RiskHintsSection, error) {
	empty := RiskHintsSection{Freshness: freshness, Items: []RiskHint{}}
	if dom == nil || st == nil || taskID == "" {
		return empty, nil
	}

	items := make([]RiskHint, 0, maxRiskHintsCap)

	blocking, err := dom.CountBlockingUncertainties(ctx, taskID)
	if err != nil {
		return RiskHintsSection{}, fmt.Errorf("risk hints: blocking count: %w", err)
	}
	if blocking > 0 {
		items = append(items, RiskHint{
			Code:     riskHintBlocking,
			Severity: "high",
			Detail:   fmt.Sprintf("%d blocking uncertaint(y/ies) open", blocking),
		})
	}

	verifyDebt, err := dom.HasVerificationDebt(ctx, taskID)
	if err != nil {
		return RiskHintsSection{}, fmt.Errorf("risk hints: verification debt: %w", err)
	}
	if verifyDebt && len(items) < maxRiskHintsCap {
		items = append(items, RiskHint{
			Code:     riskHintVerification,
			Severity: "medium",
			Detail:   "implementation recorded without satisfactory verification",
		})
	}

	if len(items) < maxRiskHintsCap {
		pathCount, err := latestChangePathCount(st, taskID)
		if err != nil {
			return RiskHintsSection{}, fmt.Errorf("risk hints: latest change paths: %w", err)
		}
		if pathCount > manyPathsThreshold {
			items = append(items, RiskHint{
				Code:     riskHintManyPaths,
				Severity: "medium",
				Detail:   fmt.Sprintf("%d paths on latest change", pathCount),
			})
		}
	}

	if len(items) < maxRiskHintsCap {
		churnPaths, err := st.ListHighChurnPaths(taskID, highChurnMinChanges)
		if err != nil {
			return RiskHintsSection{}, fmt.Errorf("risk hints: high churn paths: %w", err)
		}
		if len(churnPaths) > 0 {
			items = append(items, RiskHint{
				Code:     riskHintHighChurn,
				Severity: "low",
				Detail:   fmt.Sprintf("path %s touched %d times", churnPaths[0], highChurnMinChanges),
			})
		}
	}

	if len(items) > maxRiskHintsCap {
		items = items[:maxRiskHintsCap]
	}
	return RiskHintsSection{Freshness: freshness, Items: items}, nil
}

func latestChangePathCount(st *store.Store, taskID string) (int, error) {
	changes, err := st.ListChangesByTaskID(taskID)
	if err != nil {
		return 0, err
	}
	if len(changes) == 0 {
		return 0, nil
	}
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].CreatedAt == changes[j].CreatedAt {
			return changes[i].ID > changes[j].ID
		}
		return changes[i].CreatedAt > changes[j].CreatedAt
	})
	return st.CountChangePaths(changes[0].ID)
}
