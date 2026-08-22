package retrieval_test

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestNeighborhoodIncludesGoalIDEdge(t *testing.T) {
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
	// goal_id only — no entity_links row for goal→task
	eng := retrieval.New(st)
	out, err := eng.Neighborhood(ctx, retrieval.NeighborhoodOpts{Center: task.ID, MaxNodes: 10, Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	hasGoal := false
	hasGoalEdge := false
	for _, n := range out.Nodes {
		if n.ID == g.ID {
			hasGoal = true
		}
	}
	for _, e := range out.Edges {
		if e.Rel == retrieval.ReasonGoalHasTask && e.From == g.ID && e.To == task.ID {
			hasGoalEdge = true
		}
	}
	if !hasGoal || !hasGoalEdge {
		t.Fatalf("expected goal via goal_id: nodes=%+v edges=%+v", out.Nodes, out.Edges)
	}
}

func TestProjectGraphIncludesGoalAndTask(t *testing.T) {
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
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "D"})
	if err != nil {
		t.Fatal(err)
	}
	_ = svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{})

	eng := retrieval.New(st)
	out, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: 500})
	if err != nil {
		t.Fatal(err)
	}
	if out.Mode != "project" || out.TotalEntities < 3 {
		t.Fatalf("mode/total: %+v", out)
	}
	if len(out.Nodes) < 3 {
		t.Fatalf("nodes=%d want >=3", len(out.Nodes))
	}
	hasGoalEdge := false
	hasDecisionEdge := false
	for _, e := range out.Edges {
		if e.Rel == retrieval.ReasonGoalHasTask {
			hasGoalEdge = true
		}
		if e.Rel == retrieval.ReasonDecisionAffectsTask {
			hasDecisionEdge = true
		}
	}
	if !hasGoalEdge || !hasDecisionEdge {
		t.Fatalf("edges=%+v", out.Edges)
	}
	if out.Center != g.ID {
		t.Fatalf("center=%q want goal %q", out.Center, g.ID)
	}
}

func TestProjectGraphTruncationHonesty(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if _, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: fmt.Sprintf("D%d", i)}); err != nil {
			t.Fatal(err)
		}
	}

	eng := retrieval.New(st)
	out, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Truncated || out.TotalEntities < 5 || len(out.Nodes) != 2 {
		t.Fatalf("truncation: %+v", out)
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
