package agents_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func openAgentsTest(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st), st
}

func seedDefaultHarnessCatalog(t *testing.T, svc *domain.Service) {
	t.Helper()
	ctx := context.Background()
	catalog := []domain.HarnessAgentInput{
		{Slug: "agent:code-reviewer", Title: "Code Reviewer", SubagentType: "code-reviewer", RecommendSubagent: true, DeliberationPhases: `["CRITIQUE"]`},
		{Slug: "agent:nested-reviewer", Title: "Nested Reviewer", SubagentType: "nested-reviewer", RecommendSubagent: true, DeliberationPhases: `["CRITIQUE"]`},
		{Slug: "agent:performance-reviewer", Title: "Performance Reviewer", SubagentType: "performance-reviewer", DeliberationPhases: `["VERIFY"]`, TaskKeywords: `["perf","performance","latency"]`},
		{Slug: "agent:security-reviewer", Title: "Security Reviewer", SubagentType: "security-reviewer", DeliberationPhases: `["CRITIQUE","VERIFY"]`},
		{Slug: "agent:explore", Title: "Explore", SubagentType: "explore", DeliberationPhases: `["INVESTIGATE","ORIENT"]`},
		{Slug: "agent:generalPurpose", Title: "General Purpose", SubagentType: "generalPurpose"},
	}
	for _, in := range catalog {
		if _, err := svc.UpsertHarnessAgent(ctx, in); err != nil {
			t.Fatalf("UpsertHarnessAgent %s: %v", in.Slug, err)
		}
	}
}

func TestRecommendAgentForPhaseCritique(t *testing.T) {
	svc, st := openAgentsTest(t)
	seedDefaultHarnessCatalog(t, svc)
	ctx := context.Background()

	recs, err := agents.RecommendAgents(ctx, st, agents.RecommendInput{Phase: "CRITIQUE"})
	if err != nil {
		t.Fatalf("RecommendAgents: %v", err)
	}
	if len(recs) < 2 {
		t.Fatalf("want at least 2 recommendations, got %+v", recs)
	}
	if recs[0].AgentSlug != "agent:code-reviewer" {
		t.Fatalf("first recommendation want code-reviewer, got %+v", recs[0])
	}
	if recs[1].AgentSlug != "agent:nested-reviewer" {
		t.Fatalf("second recommendation want nested-reviewer, got %+v", recs[1])
	}
}

func TestRecommendPerformanceReviewerForPerfTask(t *testing.T) {
	svc, st := openAgentsTest(t)
	seedDefaultHarnessCatalog(t, svc)
	ctx := context.Background()

	recs, err := agents.RecommendAgents(ctx, st, agents.RecommendInput{
		Phase:     "VERIFY",
		TaskTitle: "Fix latency regression in benchmark suite",
		TaskTags:  []string{"perf"},
	})
	if err != nil {
		t.Fatalf("RecommendAgents: %v", err)
	}
	if len(recs) == 0 || recs[0].AgentSlug != "agent:performance-reviewer" {
		t.Fatalf("want performance-reviewer first, got %+v", recs)
	}
}

func TestRoutingDeterministic(t *testing.T) {
	svc, st := openAgentsTest(t)
	seedDefaultHarnessCatalog(t, svc)
	ctx := context.Background()

	in := agents.RecommendInput{
		Phase:        "VERIFY",
		TaskTitle:    "Audit auth injection owasp",
		TaskTags:     []string{"security"},
		GoalKeywords: []string{"xss"},
		HarnessCaps:  map[string]string{"harness:subagent": "AVAILABLE"},
	}
	a, err := agents.RecommendAgents(ctx, st, in)
	if err != nil {
		t.Fatalf("first RecommendAgents: %v", err)
	}
	b, err := agents.RecommendAgents(ctx, st, in)
	if err != nil {
		t.Fatalf("second RecommendAgents: %v", err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("deterministic mismatch:\na=%+v\nb=%+v", a, b)
	}
}

func TestRecommendEmptyCatalog(t *testing.T) {
	_, st := openAgentsTest(t)
	ctx := context.Background()
	recs, err := agents.RecommendAgents(ctx, st, agents.RecommendInput{Phase: "CRITIQUE"})
	if err != nil {
		t.Fatalf("RecommendAgents: %v", err)
	}
	if len(recs) != 0 {
		t.Fatalf("empty catalog want [], got %+v", recs)
	}
}

func TestRecommendMissingCapabilities(t *testing.T) {
	svc, st := openAgentsTest(t)
	ctx := context.Background()
	if _, err := svc.UpsertHarnessAgent(ctx, domain.HarnessAgentInput{
		Slug: "agent:code-reviewer", Title: "CR", SubagentType: "code-reviewer",
		Requirements: []string{"skill:code-review", "mcp:trace_review"},
	}); err != nil {
		t.Fatal(err)
	}
	recs, err := agents.RecommendAgents(ctx, st, agents.RecommendInput{Phase: "CRITIQUE"})
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 1 {
		t.Fatalf("want 1 rec, got %+v", recs)
	}
	if len(recs[0].MissingCapabilities) != 2 {
		t.Fatalf("want missing caps, got %+v", recs[0].MissingCapabilities)
	}
}
