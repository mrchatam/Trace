package planner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func setupDeepPlanScope(t *testing.T) (*planner.Service, *domain.Service, *store.Store, string, string) {
	t.Helper()
	psvc, st := openPlanner(t)
	dsvc := domain.New(st)
	ctx := context.Background()
	g := mustGoal(t, st, "replan-goal")
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{
			Title:  "P0",
			Scopes: []planner.ScopeInput{{Title: "S0"}, {Title: "S1"}},
		}},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatalf("SetCurrentScope: %v", err)
	}
	_, err = psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{"initial"},
		WorkItems:    []planner.WorkItem{{Title: "w0"}},
	})
	if err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}
	return psvc, dsvc, st, g.ID, scopeID
}

func TestApplyDiscoveryReplanINFONoSupersede(t *testing.T) {
	psvc, dsvc, st, _, scopeID := setupDeepPlanScope(t)
	ctx := context.Background()

	before, err := st.GetActiveScopeDeepPlan(scopeID)
	if err != nil {
		t.Fatalf("active before: %v", err)
	}
	sc0, _ := st.GetPlanScope(scopeID)

	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "nice to know", Severity: domain.SeverityINFO,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}

	res, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID:     disc.ID,
		ScopeID:         scopeID,
		PlanChangeTitle: "optional note",
		ExitCriteria:    []string{"should-not-apply"},
	})
	if err != nil {
		t.Fatalf("ApplyDiscoveryReplan INFO: %v", err)
	}
	if res.AutoReplanApplied || res.Reason != "severity_info" {
		t.Fatalf("want AutoReplanApplied=false reason=severity_info, got %+v", res)
	}
	if res.PlanChangeID == "" {
		t.Fatal("expected PlanChange when title provided")
	}
	after, err := st.GetActiveScopeDeepPlan(scopeID)
	if err != nil {
		t.Fatalf("active after: %v", err)
	}
	if after.ID != before.ID {
		t.Fatalf("INFO must not supersede: before=%s after=%s", before.ID, after.ID)
	}
	sc1, _ := st.GetPlanScope(scopeID)
	if sc1.AutoReplanCount != sc0.AutoReplanCount {
		t.Fatalf("count changed on INFO: %d → %d", sc0.AutoReplanCount, sc1.AutoReplanCount)
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, disc.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelDiscoveryCausesPlanChange {
		t.Fatalf("expected discovery→plan_change link: %+v err=%v", links, err)
	}
}

func TestApplyDiscoveryReplanPlanAffectingSupersedes(t *testing.T) {
	psvc, dsvc, st, _, scopeID := setupDeepPlanScope(t)
	ctx := context.Background()

	before, err := st.GetActiveScopeDeepPlan(scopeID)
	if err != nil {
		t.Fatalf("active before: %v", err)
	}

	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "api shape wrong", Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}

	res, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID:     disc.ID,
		ScopeID:         scopeID,
		PlanChangeTitle: "revise API",
		ExitCriteria:    []string{"revised"},
		WorkItems:       []planner.WorkItem{{Title: "fix contracts"}},
	})
	if err != nil {
		t.Fatalf("ApplyDiscoveryReplan: %v", err)
	}
	if !res.AutoReplanApplied || res.RevisionID == "" || res.AutoReplanCount != 1 {
		t.Fatalf("want applied + revision + count=1, got %+v", res)
	}
	after, err := st.GetActiveScopeDeepPlan(scopeID)
	if err != nil {
		t.Fatalf("active after: %v", err)
	}
	if after.ID == before.ID || after.ID != res.RevisionID {
		t.Fatalf("expected new ACTIVE revision: before=%s after=%s res=%s", before.ID, after.ID, res.RevisionID)
	}
	old, err := st.GetScopeDeepPlan(before.ID)
	if err != nil || old.Status != store.StatusSuperseded {
		t.Fatalf("old revision supersede: %+v err=%v", old, err)
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, disc.ID)
	if err != nil || len(links) != 1 {
		t.Fatalf("LinkDiscoveryPlanChange missing: %+v err=%v", links, err)
	}
}

func TestApplyDiscoveryReplanBlockingLikePlanAffecting(t *testing.T) {
	psvc, dsvc, _, _, scopeID := setupDeepPlanScope(t)
	ctx := context.Background()

	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "blocker", Severity: domain.SeverityBlocking,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	res, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID:  disc.ID,
		ScopeID:      scopeID,
		ExitCriteria: []string{"unblocked"},
	})
	if err != nil {
		t.Fatalf("ApplyDiscoveryReplan BLOCKING: %v", err)
	}
	if !res.AutoReplanApplied || res.AutoReplanCount != 1 {
		t.Fatalf("BLOCKING should auto-replan: %+v", res)
	}
}

func TestApplyDiscoveryReplanBudgetAndAck(t *testing.T) {
	psvc, dsvc, st, _, scopeID := setupDeepPlanScope(t)
	ctx := context.Background()
	const N = 5

	for i := 0; i < N; i++ {
		disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
			Title: "churn", Severity: domain.SeverityPlanAffecting,
		})
		if err != nil {
			t.Fatalf("CreateDiscovery %d: %v", i, err)
		}
		res, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
			DiscoveryID:    disc.ID,
			ScopeID:        scopeID,
			ExitCriteria:   []string{"c"},
			MaxAutoReplans: N,
		})
		if err != nil {
			t.Fatalf("replan %d: %v", i, err)
		}
		if res.AutoReplanCount != i+1 {
			t.Fatalf("count want %d got %d", i+1, res.AutoReplanCount)
		}
	}

	discFail, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "over budget", Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery fail: %v", err)
	}
	_, err = psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID: discFail.ID, ScopeID: scopeID, ExitCriteria: []string{"x"}, MaxAutoReplans: N,
	})
	if err == nil || !errors.Is(err, planner.ErrReplanBudgetExceeded) {
		t.Fatalf("want ErrReplanBudgetExceeded, got %v", err)
	}

	if err := psvc.AckReplan(ctx, scopeID); err != nil {
		t.Fatalf("AckReplan: %v", err)
	}
	sc, err := st.GetPlanScope(scopeID)
	if err != nil || sc.AutoReplanCount != 0 {
		t.Fatalf("after ack count=%d err=%v", sc.AutoReplanCount, err)
	}

	discOK, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "after ack", Severity: domain.SeverityPlanAffecting,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery after ack: %v", err)
	}
	res, err := psvc.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
		DiscoveryID: discOK.ID, ScopeID: scopeID, ExitCriteria: []string{"again"}, MaxAutoReplans: N,
	})
	if err != nil {
		t.Fatalf("replan after ack: %v", err)
	}
	if !res.AutoReplanApplied || res.AutoReplanCount != 1 {
		t.Fatalf("after ack want count=1 applied: %+v", res)
	}
}

func TestCreateDiscoverySeverityValidation(t *testing.T) {
	_, dsvc, _, _, _ := setupDeepPlanScope(t)
	ctx := context.Background()

	d, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "default"})
	if err != nil {
		t.Fatalf("default: %v", err)
	}
	if d.Severity != domain.SeverityINFO {
		t.Fatalf("default severity=%q", d.Severity)
	}
	_, err = dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "bad", Severity: "CRITICAL"})
	if err == nil {
		t.Fatal("expected reject garbage severity")
	}
}
