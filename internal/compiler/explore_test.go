package compiler_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

type failAgentQuerySearch struct {
	inner     compiler.ExploreEngine
	taskTitle string
}

func (f failAgentQuerySearch) Search(ctx context.Context, q string, opts retrieval.SearchOptions) ([]retrieval.Hit, error) {
	if strings.TrimSpace(q) != strings.TrimSpace(f.taskTitle) {
		return nil, fmt.Errorf("stub: agent query search failed")
	}
	return f.inner.Search(ctx, q, opts)
}

func (f failAgentQuerySearch) SearchGraphLabels(ctx context.Context, intent retrieval.Intent, opts retrieval.SearchOptions) ([]retrieval.Hit, error) {
	return f.inner.SearchGraphLabels(ctx, intent, opts)
}

func (f failAgentQuerySearch) Why(ctx context.Context, entityType, entityID string) (retrieval.WhyResult, error) {
	return f.inner.Why(ctx, entityType, entityID)
}

func (f failAgentQuerySearch) Neighborhood(ctx context.Context, opts retrieval.NeighborhoodOpts) (*retrieval.BoundedGraph, error) {
	return f.inner.Neighborhood(ctx, opts)
}

func openExploreFixtures(t *testing.T) (*compiler.Compiler, compiler.ExploreEngine, *domain.Service, context.Context) {
	t.Helper()
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	eng := retrieval.New(st)
	return compiler.New(st).WithRetrieval(eng), eng, domain.New(st), context.Background()
}

func TestExploreTaskRequired(t *testing.T) {
	comp, eng, _, ctx := openExploreFixtures(t)
	_, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{})
	if err == nil || !strings.Contains(err.Error(), "task_id is required") {
		t.Fatalf("expected task_id required, got %v", err)
	}
}

func TestExploreTaskMoatPreserved(t *testing.T) {
	comp, eng, svc, ctx := openExploreFixtures(t)
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Moat explore task", Body: "body"})
	if err != nil {
		t.Fatal(err)
	}

	out, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{TaskID: task.ID, Query: "any"})
	if err != nil {
		t.Fatal(err)
	}
	if out.TaskSummary.TaskID != task.ID {
		t.Fatalf("task_id=%q want %q", out.TaskSummary.TaskID, task.ID)
	}
	if out.TaskSummary.Title != task.Title {
		t.Fatalf("title=%q want %q", out.TaskSummary.Title, task.Title)
	}
	if out.TaskSummary.WorkState == "" {
		t.Fatal("expected work_state in task summary")
	}
}

func TestExploreQueryMerged(t *testing.T) {
	comp, eng, svc, ctx := openExploreFixtures(t)
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "G2 query merge task"})
	if err != nil {
		t.Fatal(err)
	}
	const token = "g2explorequerytokenXYZ"
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Unlinked " + token})
	if err != nil {
		t.Fatal(err)
	}

	without, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{TaskID: task.ID})
	if err != nil {
		t.Fatal(err)
	}
	with, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{TaskID: task.ID, Query: token})
	if err != nil {
		t.Fatal(err)
	}
	if with.PacketBudget.ItemsKept <= without.PacketBudget.ItemsKept {
		t.Fatalf("expected G1 query merge to increase kept items: without=%d with=%d",
			without.PacketBudget.ItemsKept, with.PacketBudget.ItemsKept)
	}
	found := false
	for _, h := range with.SearchHits {
		if h.EntityType == "decision" && h.EntityID == dec.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected query search hit for decision; hits=%+v", with.SearchHits)
	}
}

func TestExploreCappedHonest(t *testing.T) {
	comp, eng, svc, ctx := openExploreFixtures(t)
	const capToken = "g2explorecaphonest"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: capToken + " anchor"})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 80; i++ {
		dec, err := svc.CreateDecision(ctx, domain.DecisionInput{
			Title: fmt.Sprintf("%s linked unique%d", capToken, i),
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
			t.Fatal(err)
		}
	}

	out, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{
		TaskID: task.ID,
		Query:  capToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Truncated {
		t.Fatal("expected truncated=true for over-budget explore")
	}
	if !out.PacketBudget.CandidatesCapped && !out.PacketBudget.Truncated {
		t.Fatalf("expected packet cap honesty: budget=%+v", out.PacketBudget)
	}
	if len(out.SearchHits) > compiler.MaxExploreSearchLimit {
		t.Fatalf("search hits exceeded cap: %d", len(out.SearchHits))
	}
}

func TestExploreNoDump(t *testing.T) {
	comp, eng, svc, ctx := openExploreFixtures(t)
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "No dump task"})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 40; i++ {
		dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: fmt.Sprintf("linked node %d", i)})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
			t.Fatal(err)
		}
	}

	out, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{TaskID: task.ID})
	if err != nil {
		t.Fatal(err)
	}
	if out.Neighborhood == nil {
		t.Fatal("expected neighborhood slice")
	}
	if out.Neighborhood.MaxNodes != compiler.DefaultExploreMaxNodes {
		t.Fatalf("default max_nodes=%d want %d", out.Neighborhood.MaxNodes, compiler.DefaultExploreMaxNodes)
	}
	if len(out.Neighborhood.Nodes) > compiler.DefaultExploreMaxNodes {
		t.Fatalf("neighborhood nodes=%d exceeds default cap %d", len(out.Neighborhood.Nodes), compiler.DefaultExploreMaxNodes)
	}
	if len(out.SearchHits) > compiler.MaxExploreSearchLimit {
		t.Fatalf("search hits=%d exceeds cap %d", len(out.SearchHits), compiler.MaxExploreSearchLimit)
	}
}

func TestExploreWhyIncluded(t *testing.T) {
	comp, eng, svc, ctx := openExploreFixtures(t)
	const token = "g2whyincludetoken"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Why task"})
	if err != nil {
		t.Fatal(err)
	}
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Policy " + token})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	out, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{TaskID: task.ID, Query: token})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.SearchHits) == 0 {
		t.Fatal("expected search hits for why slice")
	}
	if len(out.WhySlices) == 0 {
		t.Fatal("expected why slices on top hits")
	}
	if out.Neighborhood == nil || len(out.Neighborhood.Nodes) < 1 {
		t.Fatal("expected bounded neighborhood on task")
	}
	if len(out.WhySlices) > compiler.DefaultExploreWhyTopN {
		t.Fatalf("why slices=%d exceed top-N=%d", len(out.WhySlices), compiler.DefaultExploreWhyTopN)
	}
}

func TestExploreFailOpenSearch(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	ctx := context.Background()

	taskTitle := "plain title no slash"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: taskTitle})
	if err != nil {
		t.Fatal(err)
	}
	inner := retrieval.New(st)
	comp := compiler.New(st).WithRetrieval(inner)
	stub := failAgentQuerySearch{inner: inner, taskTitle: taskTitle}

	out, err := compiler.Explore(ctx, comp, stub, compiler.ExploreOpts{
		TaskID: task.ID,
		Query:  "agent query fails open",
	})
	if err != nil {
		t.Fatalf("Explore: %v", err)
	}
	if out.TaskSummary.TaskID != task.ID {
		t.Fatalf("task moat lost: %+v", out.TaskSummary)
	}
	if len(out.SearchHits) != 0 {
		t.Fatalf("expected empty search hits on fail-open, got %d", len(out.SearchHits))
	}
}
