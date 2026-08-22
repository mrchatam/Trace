package replan_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func openReplan(t *testing.T) (*planner.Service, *domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return planner.New(st), domain.New(st), st
}

// TestPlantedDiscoveryReplan is the named honesty-style library demo for discovery replan.
func TestPlantedDiscoveryReplan(t *testing.T) {
	psvc, dsvc, st := openReplan(t)
	ctx := context.Background()

	// Plant Goal → CreateCoarsePlan → SetCurrentScope → DeepPlan
	goal, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "Ship progressive planner"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: goal.ID,
		Phases: []planner.PhaseInput{{
			Title: "Foundation",
			Scopes: []planner.ScopeInput{
				{Title: "Coarse planner"},
				{Title: "Discovery replan"},
			},
		}},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, goal.ID, scopeID); err != nil {
		t.Fatalf("SetCurrentScope: %v", err)
	}
	initial, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{"coarse plan works"},
		WorkItems:    []planner.WorkItem{{Title: "ship S01"}},
	})
	if err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}

	// PLAN_AFFECTING → ApplyDiscoveryReplan → deep plan superseded + count
	discPA, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Lookahead must be fail-closed",
		Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery PLAN_AFFECTING: %v", err)
	}
	resPA, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID:     discPA.ID,
		ScopeID:         scopeID,
		PlanChangeTitle: "Tighten lookahead",
		ExitCriteria:    []string{"lookahead fail-closed"},
		WorkItems:       []planner.WorkItem{{Title: "fix DeepPlan gate"}},
	})
	if err != nil {
		t.Fatalf("ApplyDiscoveryReplan PLAN_AFFECTING: %v", err)
	}
	if !resPA.AutoReplanApplied || resPA.AutoReplanCount != 1 || resPA.RevisionID == "" {
		t.Fatalf("PLAN_AFFECTING result: %+v", resPA)
	}
	active, err := st.GetActiveScopeDeepPlan(scopeID)
	if err != nil || active.ID != resPA.RevisionID || active.ID == initial.RevisionID {
		t.Fatalf("supersede mismatch: active=%+v initial=%s res=%s err=%v", active, initial.RevisionID, resPA.RevisionID, err)
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, discPA.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelDiscoveryCausesPlanChange {
		t.Fatalf("causal link: %+v err=%v", links, err)
	}

	// INFO → no supersede
	beforeINFO := active.ID
	scBeforeINFO, _ := st.GetPlanScope(scopeID)
	discINFO, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "docs typo", Severity: domain.SeverityINFO,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery INFO: %v", err)
	}
	resINFO, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID:  discINFO.ID,
		ScopeID:      scopeID,
		ExitCriteria: []string{"must-not-apply"},
	})
	if err != nil {
		t.Fatalf("ApplyDiscoveryReplan INFO: %v", err)
	}
	if resINFO.AutoReplanApplied || resINFO.Reason != "severity_info" {
		t.Fatalf("INFO result: %+v", resINFO)
	}
	afterINFO, _ := st.GetActiveScopeDeepPlan(scopeID)
	if afterINFO.ID != beforeINFO {
		t.Fatalf("INFO superseded unexpectedly: %s → %s", beforeINFO, afterINFO.ID)
	}
	scAfterINFO, _ := st.GetPlanScope(scopeID)
	if scAfterINFO.AutoReplanCount != scBeforeINFO.AutoReplanCount {
		t.Fatalf("INFO incremented count: %d → %d", scBeforeINFO.AutoReplanCount, scAfterINFO.AutoReplanCount)
	}

	// Exhaust budget (already count=1; need 4 more to reach 5, then 6th fails)
	for i := 0; i < 4; i++ {
		d, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
			Title: "churn", Severity: domain.SeverityPlanAffecting,
		})
		if err != nil {
			t.Fatalf("churn CreateDiscovery %d: %v", i, err)
		}
		_, err = psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
			DiscoveryID: d.ID, ScopeID: scopeID, ExitCriteria: []string{"c"},
		})
		if err != nil {
			t.Fatalf("churn replan %d: %v", i, err)
		}
	}
	scFull, _ := st.GetPlanScope(scopeID)
	if scFull.AutoReplanCount != planner.DefaultMaxAutoReplans {
		t.Fatalf("want count=%d got %d", planner.DefaultMaxAutoReplans, scFull.AutoReplanCount)
	}
	over, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "over", Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("over CreateDiscovery: %v", err)
	}
	_, err = psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID: over.ID, ScopeID: scopeID, ExitCriteria: []string{"x"},
	})
	if err == nil || !errors.Is(err, planner.ErrReplanBudgetExceeded) {
		t.Fatalf("want ErrReplanBudgetExceeded, got %v", err)
	}

	// Ack → succeed once
	if err := psvc.AckReplan(ctx, scopeID); err != nil {
		t.Fatalf("AckReplan: %v", err)
	}
	again, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "post-ack", Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("post-ack CreateDiscovery: %v", err)
	}
	resAgain, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID: again.ID, ScopeID: scopeID, ExitCriteria: []string{"ok"},
	})
	if err != nil {
		t.Fatalf("post-ack Apply: %v", err)
	}
	if !resAgain.AutoReplanApplied || resAgain.AutoReplanCount != 1 {
		t.Fatalf("post-ack result: %+v", resAgain)
	}
}
