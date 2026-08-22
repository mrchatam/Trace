package planner_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
)

func TestGoalStructureWarning_OverThresholdNoPlan(t *testing.T) {
	psvc, dsvc, _ := openPlannerTest(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "mega goal"})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < planner.GoalStructureWarningThreshold+1; i++ {
		if _, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID}); err != nil {
			t.Fatal(err)
		}
	}
	warn, err := psvc.GoalStructureWarning(ctx, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if warn == "" {
		t.Fatal("expected goal structure warning")
	}

	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatal(err)
	}
	if _, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID: scopeID, ExitCriteria: []string{"done"},
	}); err != nil {
		t.Fatal(err)
	}
	warn2, err := psvc.GoalStructureWarning(ctx, g.ID)
	if err != nil || warn2 != "" {
		t.Fatalf("want no warning after plan exists: %q err=%v", warn2, err)
	}
}

func TestBootstrapRecommendedAdvisory_NoPlanWithLinkedPlanChange(t *testing.T) {
	psvc, dsvc, _ := openPlannerTest(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "bootstrap goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "disc"})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.LinkDiscoveryMentionsTask(ctx, disc.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	pc, err := dsvc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "pc"})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.LinkDiscoveryPlanChange(ctx, disc.ID, pc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	advs, err := psvc.StatusAdvisories(ctx, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, a := range advs {
		if a.Code == planner.AdvisoryCodeBootstrapRecommended {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("want bootstrap_recommended in %#v", advs)
	}
}

func TestBootstrapRecommendedAdvisory_SuppressedWhenPlanExists(t *testing.T) {
	psvc, dsvc, _ := openPlannerTest(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "planned goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "disc"})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.LinkDiscoveryMentionsTask(ctx, disc.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	pc, err := dsvc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "pc"})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.LinkDiscoveryPlanChange(ctx, disc.ID, pc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatal(err)
	}
	if _, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{ScopeID: scopeID, ExitCriteria: []string{"done"}}); err != nil {
		t.Fatal(err)
	}

	advs, err := psvc.StatusAdvisories(ctx, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range advs {
		if a.Code == planner.AdvisoryCodeBootstrapRecommended {
			t.Fatalf("bootstrap advisory must not appear when plan exists: %#v", advs)
		}
	}
}
