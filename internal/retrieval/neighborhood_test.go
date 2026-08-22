package retrieval_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestNeighborhoodRequiresBudget(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	_ = svc.LinkGoalTask(ctx, g.ID, task.ID, domain.LinkMeta{})

	eng := retrieval.New(st)
	_, err = eng.Neighborhood(ctx, retrieval.NeighborhoodOpts{Center: task.ID, MaxNodes: 0})
	if err == nil {
		t.Fatal("max_nodes 0 must fail")
	}
	out, err := eng.Neighborhood(ctx, retrieval.NeighborhoodOpts{Center: task.ID, MaxNodes: 10, Depth: 2})
	if err != nil {
		t.Fatal(err)
	}
	if out.Center != task.ID || len(out.Nodes) < 1 {
		t.Fatalf("%+v", out)
	}
}

func TestNeighborhoodEdgesJSONNotNull(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Orphan"})
	if err != nil {
		t.Fatal(err)
	}

	eng := retrieval.New(st)
	out, err := eng.Neighborhood(ctx, retrieval.NeighborhoodOpts{Center: dec.ID, MaxNodes: 10, Depth: 2})
	if err != nil {
		t.Fatal(err)
	}
	if out.Edges == nil {
		t.Fatal("Edges must be non-nil slice")
	}
	if len(out.Edges) != 0 {
		t.Fatalf("expected no edges for orphan decision, got %d", len(out.Edges))
	}
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"edges":null`) {
		t.Fatalf("edges must JSON-marshal as [] not null: %s", b)
	}
}
