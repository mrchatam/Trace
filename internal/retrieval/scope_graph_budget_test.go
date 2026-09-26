package retrieval_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestProjectGraphScopeIDPopulatedUnderBudget(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	// Create a scope
	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "test-scope", Title: "Test Scope", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create 5 tasks linked to the scope
	taskIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		tk, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task " + string(rune('A'+i))})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.LinkScopeMember(ctx, tk.ID, sc.ID, domain.LinkMeta{}); err != nil {
			t.Fatal(err)
		}
		taskIDs[i] = tk.ID
	}

	eng := retrieval.New(st)

	// Test various MaxNodes values: 3, 5, 6, 100
	for _, maxNodes := range []int{3, 5, 6, 100} {
		g, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{MaxNodes: maxNodes})
		if err != nil {
			t.Fatalf("MaxNodes=%d: %v", maxNodes, err)
		}
		hasTaskWithScope := false
		hasScopeNode := false
		for _, n := range g.Nodes {
			if n.Kind == "task" {
				if n.ScopeID != sc.ID {
					t.Fatalf("MaxNodes=%d: task %s has ScopeID=%q, want %q", maxNodes, n.ID, n.ScopeID, sc.ID)
				}
				hasTaskWithScope = true
			}
			if n.Kind == "scope" && n.ID == sc.ID {
				hasScopeNode = true
			}
		}
		if !hasTaskWithScope {
			t.Fatalf("MaxNodes=%d: expected at least one task node with scope_id", maxNodes)
		}
		// For MaxNodes=3 and 5, the scope node (kind=scope, order 12) should be excluded
		// because tasks (kind=task, order 1) come first in kind order.
		if maxNodes == 3 || maxNodes == 5 {
			if hasScopeNode {
				t.Fatalf("MaxNodes=%d: scope node should be excluded by budget, but was included", maxNodes)
			}
		}
	}
}

func TestProjectGraphScopedModeTruncationHonesty(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	// Create a scope
	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create 4 tasks linked to the scope
	taskIDs := make([]string, 4)
	for i := 0; i < 4; i++ {
		tk, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task " + string(rune('A'+i))})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.LinkScopeMember(ctx, tk.ID, sc.ID, domain.LinkMeta{}); err != nil {
			t.Fatal(err)
		}
		taskIDs[i] = tk.ID
	}

	eng := retrieval.New(st)

	// MaxNodes = 5 (scope + 4 members = 5 nodes, exact fit)
	filtered5, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{
		MaxNodes: 5, Scope: "auth", Depth: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if filtered5.Truncated {
		t.Fatalf("MaxNodes=5: expected Truncated=false for exact fit, got true")
	}
	if filtered5.TotalEntities != 5 {
		t.Fatalf("MaxNodes=5: expected TotalEntities=5, got %d", filtered5.TotalEntities)
	}
	if len(filtered5.Nodes) != 5 {
		t.Fatalf("MaxNodes=5: expected 5 nodes, got %d", len(filtered5.Nodes))
	}

	// MaxNodes = 3 (scope + 4 members, budget 3 -> scope + 2 members)
	filtered3, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{
		MaxNodes: 3, Scope: "auth", Depth: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !filtered3.Truncated {
		t.Fatalf("MaxNodes=3: expected Truncated=true when nodes dropped, got false")
	}
	if filtered3.TotalEntities <= len(filtered3.Nodes) {
		t.Fatalf("MaxNodes=3: expected TotalEntities > len(Nodes) (%d > %d), got %d <= %d",
			filtered3.TotalEntities, len(filtered3.Nodes), filtered3.TotalEntities, len(filtered3.Nodes))
	}
	if len(filtered3.Nodes) != 3 {
		t.Fatalf("MaxNodes=3: expected 3 nodes, got %d", len(filtered3.Nodes))
	}
}

func TestProjectGraphScopedModeDepthValidation(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "test", Title: "Test", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, tk.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	eng := retrieval.New(st)

	// Depth 3 should be capped to 2 internally (not an error at Engine level)
	_, err = eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{
		MaxNodes: 10, Scope: "test", Depth: 3,
	})
	if err != nil {
		t.Fatalf("depth=3 should be capped, not error: %v", err)
	}

	// Depth 0 should default to 1
	g1, err := eng.ProjectGraph(ctx, retrieval.ProjectGraphOpts{
		MaxNodes: 10, Scope: "test", Depth: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if g1.MaxNodes != 10 {
		t.Fatalf("MaxNodes not preserved: %d", g1.MaxNodes)
	}
}
