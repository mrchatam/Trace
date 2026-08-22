package planner_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func openPlanner(t *testing.T) (*planner.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return planner.New(st), st
}

func mustGoal(t *testing.T, st *store.Store, title string) store.Goal {
	t.Helper()
	g, err := st.UpsertGoal(store.Goal{Title: title, Body: "b"})
	if err != nil {
		t.Fatalf("UpsertGoal: %v", err)
	}
	return g
}

func TestCreateCoarsePlanPersistsOrdAndRejects(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()

	_, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{})
	if err == nil {
		t.Fatal("expected error for empty goal")
	}

	_, err = svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: "missing-goal",
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S"}}}},
	})
	if err == nil || !errors.Is(err, planner.ErrNotFound) {
		t.Fatalf("want ErrNotFound for missing goal, got %v", err)
	}

	g := mustGoal(t, st, "G")
	_, err = svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "", Scopes: []planner.ScopeInput{{Title: "S"}}}},
	})
	if err == nil {
		t.Fatal("expected empty phase title reject")
	}
	_, err = svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: ""}}}},
	})
	if err == nil {
		t.Fatal("expected empty scope title reject")
	}

	cp, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{
			{
				Title: "Phase A",
				Scopes: []planner.ScopeInput{
					{Title: "S0"},
					{Title: "S1"},
				},
			},
			{
				Title:  "Phase B",
				Scopes: []planner.ScopeInput{{Title: "S2"}},
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	if len(cp.Phases) != 2 {
		t.Fatalf("phases=%d", len(cp.Phases))
	}
	if cp.Phases[0].Ord != 0 || cp.Phases[1].Ord != 1 {
		t.Fatalf("phase ord: %+v", cp.Phases)
	}
	if len(cp.Phases[0].Scopes) != 2 || cp.Phases[0].Scopes[0].Ord != 0 || cp.Phases[0].Scopes[1].Ord != 1 {
		t.Fatalf("scope ord: %+v", cp.Phases[0].Scopes)
	}
	if cp.Phases[1].Scopes[0].Ord != 0 {
		t.Fatalf("phase B scope ord: %d", cp.Phases[1].Scopes[0].Ord)
	}

	// goal_plan_state initialized with NULL current
	stt, err := st.GetGoalPlanState(g.ID)
	if err != nil {
		t.Fatalf("GetGoalPlanState: %v", err)
	}
	if stt.CurrentScopeID != nil {
		t.Fatalf("expected nil current, got %v", *stt.CurrentScopeID)
	}
}

func TestSetCurrentScopeRoundTripAndRejectForeign(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()
	g1 := mustGoal(t, st, "G1")
	g2 := mustGoal(t, st, "G2")

	cp1, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g1.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "A"}, {Title: "B"}}}},
	})
	if err != nil {
		t.Fatalf("cp1: %v", err)
	}
	cp2, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g2.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "X"}}}},
	})
	if err != nil {
		t.Fatalf("cp2: %v", err)
	}

	scopeA := cp1.Phases[0].Scopes[0].ID
	scopeX := cp2.Phases[0].Scopes[0].ID

	if err := svc.SetCurrentScope(ctx, g1.ID, scopeX); err == nil {
		t.Fatal("expected reject foreign scope")
	}
	if err := svc.SetCurrentScope(ctx, g1.ID, scopeA); err != nil {
		t.Fatalf("SetCurrentScope: %v", err)
	}
	ref, err := svc.GetCurrentScope(ctx, g1.ID)
	if err != nil {
		t.Fatalf("GetCurrentScope: %v", err)
	}
	if ref.ID != scopeA {
		t.Fatalf("got %s want %s", ref.ID, scopeA)
	}
}

func TestDeepPlanCurrentOnlyAndLookahead(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()
	g := mustGoal(t, st, "G")
	cp, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{
			Title: "P",
			Scopes: []planner.ScopeInput{
				{Title: "S0"},
				{Title: "S1"},
				{Title: "S2"},
			},
		}},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	s0 := cp.Phases[0].Scopes[0].ID
	s1 := cp.Phases[0].Scopes[1].ID
	s2 := cp.Phases[0].Scopes[2].ID

	_, err = svc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID: s0, ExitCriteria: []string{"done"}, WorkItems: []planner.WorkItem{{Title: "w"}},
	})
	if err == nil {
		t.Fatal("DeepPlan before SetCurrent should fail")
	}

	if err := svc.SetCurrentScope(ctx, g.ID, s0); err != nil {
		t.Fatalf("SetCurrent: %v", err)
	}

	_, err = svc.DeepPlan(ctx, planner.DeepPlanInput{ScopeID: s1, ExitCriteria: []string{"x"}})
	var nc *planner.ErrNotCurrent
	if !errors.As(err, &nc) {
		t.Fatalf("want ErrNotCurrent, got %v", err)
	}

	res, err := svc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:          s0,
		ExitCriteria:     []string{"exit-a"},
		Constraints:      []string{"c1"},
		WorkItems:        []planner.WorkItem{{Title: "item1", Notes: "n"}},
		LookaheadSummary: "shallow next",
	})
	if err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}
	if res.LookaheadScopeID != s1 {
		t.Fatalf("lookahead want %s got %s", s1, res.LookaheadScopeID)
	}
	if res.Document.LookaheadScopeID != s1 || res.Document.LookaheadSummary != "shallow next" {
		t.Fatalf("doc lookahead: %+v", res.Document)
	}

	// No deep revisions for non-current scopes in that call
	for _, id := range []string{s1, s2} {
		n, err := st.CountScopeDeepPlans(id)
		if err != nil {
			t.Fatalf("count %s: %v", id, err)
		}
		if n != 0 {
			t.Fatalf("scope %s deep plans=%d want 0", id, n)
		}
	}
	n0, err := st.CountScopeDeepPlans(s0)
	if err != nil || n0 != 1 {
		t.Fatalf("s0 deep plans=%d err=%v", n0, err)
	}

	// lookahead body updated shallowly
	sc1, err := st.GetPlanScope(s1)
	if err != nil {
		t.Fatalf("GetPlanScope s1: %v", err)
	}
	if sc1.Body != "shallow next" {
		t.Fatalf("lookahead body=%q", sc1.Body)
	}

	// re-deep supersedes prior ACTIVE
	res2, err := svc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID: s0, ExitCriteria: []string{"exit-b"}, WorkItems: []planner.WorkItem{{Title: "item2"}},
	})
	if err != nil {
		t.Fatalf("re-DeepPlan: %v", err)
	}
	if res2.SupersededCount != 1 {
		t.Fatalf("superseded=%d", res2.SupersededCount)
	}
	active, err := st.GetActiveScopeDeepPlan(s0)
	if err != nil {
		t.Fatalf("GetActive: %v", err)
	}
	if active.ID != res2.RevisionID {
		t.Fatalf("active id mismatch")
	}
	n0, _ = st.CountScopeDeepPlans(s0)
	if n0 != 2 {
		t.Fatalf("want 2 revisions, got %d", n0)
	}
	old, err := st.GetScopeDeepPlan(res.RevisionID)
	if err != nil {
		t.Fatalf("get old: %v", err)
	}
	if old.Status != store.StatusSuperseded {
		t.Fatalf("old status=%s", old.Status)
	}
}

func TestSupersedeDeepPlanSmoke(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()
	g := mustGoal(t, st, "G")
	cp, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S0"}, {Title: "S1"}}}},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	s0 := cp.Phases[0].Scopes[0].ID
	s1 := cp.Phases[0].Scopes[1].ID
	if err := svc.SetCurrentScope(ctx, g.ID, s0); err != nil {
		t.Fatalf("SetCurrent: %v", err)
	}
	first, err := svc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID: s0, ExitCriteria: []string{"a"}, WorkItems: []planner.WorkItem{{Title: "w1"}},
	})
	if err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}

	// Supersede without being current (target s1 which is not current)
	res, err := svc.SupersedeDeepPlan(ctx, planner.SupersedeInput{
		ScopeID:          s1,
		ExitCriteria:     []string{"replan"},
		WorkItems:        []planner.WorkItem{{Title: "replan-item"}},
		LookaheadSummary: "later",
	})
	if err != nil {
		t.Fatalf("SupersedeDeepPlan on non-current: %v", err)
	}
	if res.RevisionID == "" {
		t.Fatal("empty revision")
	}
	active1, err := st.GetActiveScopeDeepPlan(s1)
	if err != nil {
		t.Fatalf("active s1: %v", err)
	}
	if active1.ID != res.RevisionID {
		t.Fatalf("active mismatch")
	}

	// Supersede current scope's plan
	res2, err := svc.SupersedeDeepPlan(ctx, planner.SupersedeInput{
		ScopeID:      s0,
		ExitCriteria: []string{"b"},
		WorkItems:    []planner.WorkItem{{Title: "w2"}},
	})
	if err != nil {
		t.Fatalf("SupersedeDeepPlan s0: %v", err)
	}
	if res2.SupersededCount < 1 {
		t.Fatalf("expected supersede of prior ACTIVE, got %d", res2.SupersededCount)
	}
	active0, err := st.GetActiveScopeDeepPlan(s0)
	if err != nil {
		t.Fatalf("active s0: %v", err)
	}
	if active0.ID == first.RevisionID {
		t.Fatal("old revision still ACTIVE")
	}
	var doc planner.DeepPlanDocument
	if err := json.Unmarshal([]byte(active0.ContentJSON), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(doc.ExitCriteria) != 1 || doc.ExitCriteria[0] != "b" {
		t.Fatalf("doc=%+v", doc)
	}
}

func TestGetPlanView(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()
	g := mustGoal(t, st, "G")
	cp, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S0"}, {Title: "S1"}}}},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	s0 := cp.Phases[0].Scopes[0].ID
	s1 := cp.Phases[0].Scopes[1].ID
	_ = svc.SetCurrentScope(ctx, g.ID, s0)
	_, err = svc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID: s0, ExitCriteria: []string{"e"}, LookaheadSummary: "next-shallow",
	})
	if err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}
	view, err := svc.GetPlan(ctx, g.ID)
	if err != nil {
		t.Fatalf("GetPlan: %v", err)
	}
	if view.CurrentScopeID == nil || *view.CurrentScopeID != s0 {
		t.Fatalf("current=%v", view.CurrentScopeID)
	}
	if view.CurrentDeepPlan == nil || view.CurrentDeepPlan.ExitCriteria[0] != "e" {
		t.Fatalf("deep=%v", view.CurrentDeepPlan)
	}
	if view.LookaheadScopeID != s1 || view.LookaheadSummary != "next-shallow" {
		t.Fatalf("lookahead=%s %q", view.LookaheadScopeID, view.LookaheadSummary)
	}
	_ = st
}

func TestListScopesOrdered(t *testing.T) {
	svc, st := openPlanner(t)
	ctx := context.Background()
	g := mustGoal(t, st, "G")
	_, err := svc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{
			{Title: "P0", Scopes: []planner.ScopeInput{{Title: "A"}, {Title: "B"}}},
			{Title: "P1", Scopes: []planner.ScopeInput{{Title: "C"}}},
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	refs, err := svc.ListScopes(ctx, g.ID)
	if err != nil {
		t.Fatalf("ListScopes: %v", err)
	}
	if len(refs) != 3 {
		t.Fatalf("len=%d", len(refs))
	}
	titles := []string{refs[0].Title, refs[1].Title, refs[2].Title}
	want := []string{"A", "B", "C"}
	for i := range want {
		if titles[i] != want[i] {
			t.Fatalf("order=%v want %v", titles, want)
		}
	}
}

// Ensure sql.ErrNoRows wrapping is detectable for missing goals.
func TestMissingGoalWrapped(t *testing.T) {
	svc, _ := openPlanner(t)
	_, err := svc.CreateCoarsePlan(context.Background(), planner.CoarsePlanInput{
		GoalID: "nope",
		Phases: []planner.PhaseInput{{Title: "P", Scopes: []planner.ScopeInput{{Title: "S"}}}},
	})
	if !errors.Is(err, planner.ErrNotFound) && !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("err=%v", err)
	}
}
