package planner

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// GoalStructureWarningThreshold is the task count above which goals without PlanExists emit advisory.
const GoalStructureWarningThreshold = 15

// StatusAdvisory is a planner-side status advisory (code + message).
type StatusAdvisory struct {
	Code    string
	Message string
}

const (
	AdvisoryCodeGoalStructureWarning = "goal_structure_warning"
	AdvisoryCodeBootstrapRecommended = "bootstrap_recommended"
)

// GoalStructureWarning returns a non-empty advisory when the goal exceeds the threshold without PlanExists.
func (s *Service) GoalStructureWarning(ctx context.Context, goalID string) (string, error) {
	_ = ctx
	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return "", &ErrValidation{Msg: "goal_id is required"}
	}
	view, err := s.GetPlan(ctx, goalID)
	if err != nil {
		return "", err
	}
	planExists := view.CurrentScopeID != nil && *view.CurrentScopeID != "" && view.CurrentDeepPlan != nil
	if planExists {
		return "", nil
	}
	tasks, err := s.store.ListTasksByGoalID(goalID)
	if err != nil {
		return "", err
	}
	if len(tasks) <= GoalStructureWarningThreshold {
		return "", nil
	}
	return fmt.Sprintf(
		"goal %s has %d tasks but no progressive plan (threshold %d); run trace plan create-coarse or trace plan bootstrap --goal %s",
		goalID, len(tasks), GoalStructureWarningThreshold, goalID,
	), nil
}

// StatusAdvisories returns goal-level advisories for loop status (orthogonal codes may both appear).
func (s *Service) StatusAdvisories(ctx context.Context, goalID string) ([]StatusAdvisory, error) {
	var out []StatusAdvisory
	if warn, err := s.GoalStructureWarning(ctx, goalID); err != nil {
		return nil, err
	} else if warn != "" {
		out = append(out, StatusAdvisory{Code: AdvisoryCodeGoalStructureWarning, Message: warn})
	}
	if adv, ok, err := s.bootstrapRecommendedAdvisory(ctx, goalID); err != nil {
		return nil, err
	} else if ok {
		out = append(out, adv)
	}
	return out, nil
}

func (s *Service) bootstrapRecommendedAdvisory(ctx context.Context, goalID string) (StatusAdvisory, bool, error) {
	exists, err := s.PlanExists(ctx, goalID)
	if err != nil {
		return StatusAdvisory{}, false, err
	}
	if exists {
		return StatusAdvisory{}, false, nil
	}
	linked := goalLinkedPlanChangeIDs(s.store, goalID)
	if len(linked) < 1 {
		return StatusAdvisory{}, false, nil
	}
	return StatusAdvisory{
		Code: AdvisoryCodeBootstrapRecommended,
		Message: fmt.Sprintf(
			"goal %s has %d linked plan-change(s) but no progressive plan; run trace plan bootstrap --goal %s or MCP trace_plan action=bootstrap",
			goalID, len(linked), goalID,
		),
	}, true, nil
}

// PlanExists reports whether the goal has current scope and active deep plan.
func (s *Service) PlanExists(ctx context.Context, goalID string) (bool, error) {
	view, err := s.GetPlan(ctx, goalID)
	if err != nil {
		return false, err
	}
	return view.CurrentScopeID != nil && *view.CurrentScopeID != "" && view.CurrentDeepPlan != nil, nil
}

func goalLinkedPlanChangeIDs(st *store.Store, goalID string) map[string]struct{} {
	out := map[string]struct{}{}
	tasks, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		return out
	}
	for _, task := range tasks {
		links, err := st.ListLinksTo(domain.EntityTask, task.ID)
		if err != nil {
			continue
		}
		for _, l := range links {
			if l.Rel != domain.RelDiscoveryMentionsTask && l.Rel != "discovery-mentions-task" {
				continue
			}
			discID := l.FromID
			dlinks, err := st.ListLinksFrom(domain.EntityDiscovery, discID)
			if err != nil {
				continue
			}
			for _, dl := range dlinks {
				if dl.Rel != domain.RelDiscoveryCausesPlanChange && dl.Rel != "discovery_causes_plan_change" {
					continue
				}
				out[dl.ToID] = struct{}{}
			}
		}
	}
	return out
}

func sortPlanChanges(candidates []store.PlanChange) {
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.CreatedAt != b.CreatedAt {
			return a.CreatedAt > b.CreatedAt
		}
		if len(a.Title) != len(b.Title) {
			return len(a.Title) > len(b.Title)
		}
		return a.ID > b.ID
	})
}
