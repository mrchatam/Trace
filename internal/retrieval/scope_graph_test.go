package retrieval_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestProjectGraphScopeEdgesAndFilter(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Ship"})
	if err != nil {
		t.Fatal(err)
	}
	tFE, err := svc.CreateTask(ctx, domain.TaskInput{Title: "FE auth", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	tBE, err := svc.CreateTask(ctx, domain.TaskInput{Title: "BE auth", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	tOther, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Unrelated", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, tFE.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, tBE.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkAPIContract(ctx, tFE.ID, tBE.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	// plan_scopes alone must not appear as kind=scope
	ph, err := st.InsertPlanPhase(store.PlanPhase{GoalID: g.ID, Title: "P", Ord: 0})
	if err != nil {
		t.Fatal(err)
	}
	ps, err := st.InsertPlanScope(store.PlanScope{PhaseID: ph.ID, Title: "plan-only", Ord: 0})
	if err != nil {
		t.Fatal(err)
	}

	eng := retrieval.New(st)
	full, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: 100})
	if err != nil {
		t.Fatal(err)
	}
	hasScopeNode := false
	hasAPI := false
	hasMember := false
	hasPlanScopeAsScope := false
	for _, n := range full.Nodes {
		if n.Kind == "scope" && n.ID == sc.ID {
			hasScopeNode = true
		}
		if n.Kind == "scope" && n.ID == ps.ID {
			hasPlanScopeAsScope = true
		}
		if n.ID == tFE.ID && n.ScopeID != sc.ID {
			t.Fatalf("FE task scope_id=%q want %q", n.ScopeID, sc.ID)
		}
	}
	for _, e := range full.Edges {
		if e.Rel == domain.RelAPIContract && e.From == tFE.ID && e.To == tBE.ID {
			hasAPI = true
			if e.Provenance != "explicit" {
				t.Fatalf("provenance=%q", e.Provenance)
			}
		}
		if e.Rel == domain.RelScopeMember {
			hasMember = true
		}
	}
	if !hasScopeNode || !hasAPI || !hasMember {
		t.Fatalf("missing scope edges: scope=%v api=%v member=%v", hasScopeNode, hasAPI, hasMember)
	}
	if hasPlanScopeAsScope {
		t.Fatal("plan_scopes must not appear as kind=scope")
	}

	// A3: scope filter
	filtered, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{
		MaxNodes: 50, Scope: "auth", Depth: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, n := range filtered.Nodes {
		ids[n.ID] = true
	}
	if !ids[sc.ID] || !ids[tFE.ID] || !ids[tBE.ID] {
		t.Fatalf("filter missing members: %#v", ids)
	}
	if ids[tOther.ID] {
		t.Fatal("unrelated task must be absent from scope filter")
	}

	_, err = eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: 0, Scope: "auth"})
	if err == nil {
		t.Fatal("missing max_nodes must fail")
	}
	_, err = eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: 5001, Scope: "auth"})
	var bud *retrieval.ErrBudgetExceeded
	if !errors.As(err, &bud) {
		t.Fatalf("want ErrBudgetExceeded, got %v", err)
	}
}
