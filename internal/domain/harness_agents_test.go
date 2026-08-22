package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestCapabilityKindAgent(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	c, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: "AGENT", Slug: "agent:perf-test", Title: "Perf agent",
	})
	if err != nil {
		t.Fatalf("UpsertCapability AGENT: %v", err)
	}
	if c.Kind != domain.CapabilityKindAgent {
		t.Fatalf("want AGENT kind, got %q", c.Kind)
	}
	got, err := svc.GetCapabilityBySlug(ctx, "agent:perf-test")
	if err != nil || got.Kind != domain.CapabilityKindAgent {
		t.Fatalf("GetCapabilityBySlug: %+v err=%v", got, err)
	}
}

func TestHarnessAgentSeedExportRoundTrip(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	agentID := "11111111-1111-4111-8111-111111111111"
	if _, err := svc.UpsertHarnessAgent(ctx, domain.HarnessAgentInput{
		ID: agentID, Slug: "agent:seed-roundtrip", Title: "Seed RT",
		SubagentType: "code-reviewer", DeliberationPhases: `["CRITIQUE"]`,
		TaskKeywords: `["review"]`, RecommendSubagent: true,
		Requirements: []string{"skill:code-review"},
	}); err != nil {
		t.Fatalf("UpsertHarnessAgent: %v", err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}
	if len(doc.HarnessAgents) != 1 {
		t.Fatalf("export harness_agents: %+v", doc.HarnessAgents)
	}
	if doc.HarnessAgents[0].Slug != "agent:seed-roundtrip" || len(doc.HarnessAgents[0].Requirements) != 1 {
		t.Fatalf("export agent detail: %+v", doc.HarnessAgents[0])
	}

	st2, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	svc2 := domain.New(st2)
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("ImportSeedDocument: %v", err)
	}
	got, err := svc2.GetHarnessAgentBySlug(ctx, "agent:seed-roundtrip")
	if err != nil || got.ID != agentID || got.Title != "Seed RT" {
		t.Fatalf("imported agent: %+v err=%v", got, err)
	}
	reqs, err := svc2.ListHarnessAgentRequirements(ctx, agentID)
	if err != nil || len(reqs) != 1 || reqs[0].RequiredCapabilitySlug != "skill:code-review" {
		t.Fatalf("imported requirements: %+v err=%v", reqs, err)
	}
}
