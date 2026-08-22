package eval

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
)

// RunAll executes selected mechanisms in stable id order.
// Individual mechanism failures set Passed=false and continue; aggregate error
// only when TaskID is invalid or Service is nil.
func RunAll(ctx context.Context, in EvalInput, opts RunOptions) ([]EvalResult, error) {
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return nil, &domain.ErrValidation{Msg: "task_id is required"}
	}
	if in.Service == nil {
		return nil, &domain.ErrValidation{Msg: "service is required"}
	}

	rules := opts.Rules
	if rules == nil && strings.TrimSpace(opts.Root) != "" && opts.Store != nil {
		load, err := LoadRules(ctx, opts.Root, opts.Store)
		if err != nil {
			return nil, err
		}
		rules = &load.Rules
	}
	in.Rules = rules

	reg := DefaultRegistry()
	targets := selectMechanismIDs(reg, opts, rules)
	now := time.Now().UTC().Format(time.RFC3339)

	results := make([]EvalResult, 0, len(targets))
	for _, id := range targets {
		m, ok := reg.mechanism(id)
		if !ok {
			continue
		}
		runCtx := ctx
		if rules != nil {
			runCtx = WithRules(runCtx, rules)
		}
		res, err := m.Run(runCtx, in)
		if err != nil {
			res = EvalResult{
				MechanismID: id,
				Passed:      false,
				Summary:     err.Error(),
				RecordedAt:  now,
			}
		}
		if res.MechanismID == "" {
			res.MechanismID = id
		}
		if res.RecordedAt == "" {
			res.RecordedAt = now
		}
		results = append(results, res)
	}
	return results, nil
}

func selectMechanismIDs(reg *Registry, opts RunOptions, rules *RulesFile) []string {
	registered := map[string]struct{}{}
	for _, id := range reg.ListMechanismIDs() {
		registered[id] = struct{}{}
	}

	if rules != nil && len(rules.Mechanisms) > 0 {
		targets := rules.FilterRegisteredMechanisms(registered)
		if len(opts.MechanismIDs) == 0 {
			return targets
		}
		allowed := mechanismAllowSet(opts.MechanismIDs)
		out := make([]string, 0, len(targets))
		for _, id := range targets {
			if _, ok := allowed[id]; ok {
				out = append(out, id)
			}
		}
		return out
	}

	filter := opts.MechanismIDs
	all := reg.ListMechanismIDs()
	if len(filter) == 0 {
		return all
	}
	allowed := mechanismAllowSet(filter)
	var out []string
	for _, id := range all {
		if _, ok := allowed[id]; ok {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

func mechanismAllowSet(filter []string) map[string]struct{} {
	allowed := map[string]struct{}{}
	for _, id := range filter {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		allowed[id] = struct{}{}
	}
	return allowed
}
