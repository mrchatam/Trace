package loop

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// RecommendHarnessInput mirrors trace agents recommend / trace_agents recommend.
type RecommendHarnessInput struct {
	TaskID   string
	Phase    string
	Keywords string
}

// RecommendHarness returns ranked agent suggestions (same shape as harness_recommendations[] items).
func RecommendHarness(ctx context.Context, dom *domain.Service, st *store.Store, psvc *planner.Service, in RecommendHarnessInput) ([]agents.Recommendation, error) {
	taskID := strings.TrimSpace(in.TaskID)
	phaseStr := strings.TrimSpace(in.Phase)

	if taskID != "" && phaseStr != "" {
		return nil, fmt.Errorf("agents recommend: specify task or phase, not both")
	}
	if taskID == "" && phaseStr == "" {
		return nil, fmt.Errorf("agents recommend: task or phase is required")
	}

	extraKW := parseKeywordString(in.Keywords)

	if taskID != "" {
		task, err := st.GetTask(taskID)
		if err != nil {
			return nil, fmt.Errorf("agents recommend: load task: %w", err)
		}
		var goalKeywords []string
		phase := deliberation.Phase("")
		if task.GoalID != nil && *task.GoalID != "" {
			goalID := *task.GoalID
			if g, err := st.GetGoal(goalID); err == nil {
				goalKeywords = goalKeywordsFromTitle(g.Title)
			}
			if psvc != nil {
				planView, err := psvc.GetPlan(ctx, goalID)
				if err == nil && planView.CurrentScopeID != nil && *planView.CurrentScopeID != "" && planView.CurrentDeepPlan != nil {
					dState := loadDeliberationState(ctx, dom, task.ID, goalID)
					seed := ApplySeed{TaskID: task.ID, GoalID: goalID}
					inputs, err := BuildPolicyInputs(ctx, dom, psvc, task.ID, goalID, nil, p19SaturatedFromLastStep(st, seed))
					if err == nil {
						phase, _, _ = deliberation.SelectNext(dState, inputs)
					}
				}
			}
		}
		if len(extraKW) > 0 {
			goalKeywords = append(goalKeywords, extraKW...)
		}
		if phase == "" {
			phase = deliberation.PhaseOrient
		}
		sec, err := buildHarnessRecommendationsSection(ctx, dom, st, phase, task.Title, goalKeywords)
		if err != nil {
			return nil, err
		}
		return sec.Items, nil
	}

	phase := deliberation.Phase(strings.ToUpper(phaseStr))
	sec, err := buildHarnessRecommendationsSection(ctx, dom, st, phase, "", extraKW)
	if err != nil {
		return nil, err
	}
	return sec.Items, nil
}

func parseKeywordString(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.ToLower(p))
	}
	return out
}
