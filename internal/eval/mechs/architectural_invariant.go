package mechs

import (
	"context"
	"encoding/json"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/eval/rulectx"
)

type ArchitecturalInvariant struct{}

func (ArchitecturalInvariant) ID() string { return idArchitecturalInvariant }

func (ArchitecturalInvariant) Run(ctx context.Context, taskID string, svc *domain.Service) (bool, string, string, error) {
	check, err := svc.CheckArchitecturalInvariants(ctx, taskID)
	if err != nil {
		return false, "", "", err
	}
	if rules := rulectx.From(ctx); rules != nil {
		filtered := check.Violations[:0]
		for _, v := range check.Violations {
			if rules.IsInvariantEnabled(v.Rule) {
				filtered = append(filtered, v)
			}
		}
		check.Violations = filtered
		check.Passed = len(filtered) == 0
	}
	details, _ := json.Marshal(check)
	if check.Passed {
		return true, "architectural invariants satisfied", string(details), nil
	}
	return false, "architectural invariant violations detected", string(details), nil
}
