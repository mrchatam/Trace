package planner_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func openPlannerTest(t *testing.T) (*planner.Service, *domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return planner.New(st), domain.New(st), st
}

func TestPlanBootstrap_Idempotent(t *testing.T) {
	psvc, dsvc, st := openPlannerTest(t)
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "boot goal"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Recovery plan", Body: "exit when ready"}); err != nil {
		t.Fatal(err)
	}

	first, err := psvc.BootstrapFromPlanChanges(ctx, g.ID, "test")
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if first.AlreadyExists || first.ScopeID == "" {
		t.Fatalf("first bootstrap: %+v", first)
	}
	ok, err := psvc.PlanExists(ctx, g.ID)
	if err != nil || !ok {
		t.Fatalf("PlanExists after bootstrap: ok=%v err=%v", ok, err)
	}

	second, err := psvc.BootstrapFromPlanChanges(ctx, g.ID, "test")
	if err != nil {
		t.Fatalf("idempotent bootstrap: %v", err)
	}
	if !second.AlreadyExists {
		t.Fatalf("want already_exists on second bootstrap: %+v", second)
	}
	if second.ScopeID != first.ScopeID {
		t.Fatalf("scope changed on idempotent bootstrap: %s vs %s", second.ScopeID, first.ScopeID)
	}
	_ = st
}
