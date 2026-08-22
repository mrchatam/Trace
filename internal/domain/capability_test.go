package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestUpsertCapabilityGetAndReject(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	if _, err := svc.UpsertCapability(ctx, domain.CapabilityInput{Slug: "x"}); err == nil {
		t.Fatal("empty kind must fail")
	}
	if _, err := svc.UpsertCapability(ctx, domain.CapabilityInput{Kind: "WIDGET", Slug: "x"}); err == nil {
		t.Fatal("unknown kind must fail")
	}
	if _, err := svc.UpsertCapability(ctx, domain.CapabilityInput{Kind: "SKILL", Slug: ""}); err == nil {
		t.Fatal("empty slug must fail")
	}
	if _, err := svc.UpsertCapability(ctx, domain.CapabilityInput{Kind: "SKILL", Slug: "s", Status: "READY"}); err == nil {
		t.Fatal("unknown status must fail")
	}

	agent, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: "AGENT", Slug: "agent:keeper-test", Title: "Keeper agent",
	})
	if err != nil {
		t.Fatalf("AGENT kind must succeed: %v", err)
	}
	if agent.Kind != domain.CapabilityKindAgent {
		t.Fatalf("want AGENT kind, got %q", agent.Kind)
	}

	c, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: "skill", Slug: "skill:tdd", Title: "TDD",
	})
	if err != nil {
		t.Fatalf("UpsertCapability: %v", err)
	}
	if c.Kind != domain.CapabilityKindSkill || c.Status != domain.CapabilityStatusUnknown {
		t.Fatalf("want SKILL/UNKNOWN, got %+v", c)
	}

	byID, err := svc.GetCapability(ctx, c.ID)
	if err != nil || byID.ID != c.ID {
		t.Fatalf("GetCapability: %+v err=%v", byID, err)
	}
	bySlug, err := svc.GetCapabilityBySlug(ctx, "skill:tdd")
	if err != nil || bySlug.ID != c.ID {
		t.Fatalf("GetCapabilityBySlug: %+v err=%v", bySlug, err)
	}

	// unique slug conflict on explicit different id still fails (DF-41)
	_, err = svc.UpsertCapability(ctx, domain.CapabilityInput{
		ID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Kind: domain.CapabilityKindTool, Slug: "skill:tdd", Title: "dup",
	})
	if err == nil {
		t.Fatal("duplicate slug with different id must fail")
	}
}

// TestUpsertCapabilityBySlugUpdatesExisting — DF-41: empty-ID re-declare by slug reuses id.
func TestUpsertCapabilityBySlugUpdatesExisting(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	first, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindSkill, Slug: "skill:upsert", Title: "One", Status: domain.CapabilityStatusUnknown,
	})
	if err != nil {
		t.Fatalf("first declare: %v", err)
	}

	second, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindSkill, Slug: "skill:upsert", Title: "Two", Status: domain.CapabilityStatusAvailable,
	})
	if err != nil {
		t.Fatalf("empty-ID re-declare: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("want same id %q, got %q", first.ID, second.ID)
	}
	if second.Status != domain.CapabilityStatusAvailable || second.Title != "Two" {
		t.Fatalf("want updated status/title, got %+v", second)
	}

	bySlug, err := svc.GetCapabilityBySlug(ctx, "skill:upsert")
	if err != nil || bySlug.ID != first.ID || bySlug.Status != domain.CapabilityStatusAvailable {
		t.Fatalf("GetCapabilityBySlug: %+v err=%v", bySlug, err)
	}
}

func TestListCapabilitiesFilter(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	_, _ = svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindSkill, Slug: "skill:a", Status: domain.CapabilityStatusAvailable,
	})
	_, _ = svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindMCP, Slug: "mcp:b", Status: domain.CapabilityStatusUnavailable,
	})
	all, err := svc.ListCapabilities(ctx, domain.ListCapabilitiesFilter{})
	if err != nil || len(all) != 2 {
		t.Fatalf("List all: %+v err=%v", all, err)
	}
	if all[0].Slug != "mcp:b" || all[1].Slug != "skill:a" {
		t.Fatalf("want stable slug order, got %+v", all)
	}
	skills, err := svc.ListCapabilities(ctx, domain.ListCapabilitiesFilter{Kind: "SKILL"})
	if err != nil || len(skills) != 1 || skills[0].Slug != "skill:a" {
		t.Fatalf("filter kind: %+v err=%v", skills, err)
	}
}

func TestRequireUnrequireAndMissing(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Need caps"})
	if err != nil {
		t.Fatal(err)
	}
	avail, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindTool, Slug: "tool:ok", Status: domain.CapabilityStatusAvailable,
	})
	if err != nil {
		t.Fatal(err)
	}
	unavail, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindTool, Slug: "tool:down", Status: domain.CapabilityStatusUnavailable,
	})
	if err != nil {
		t.Fatal(err)
	}
	unknown, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindRule, Slug: "rule:maybe",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.RequireCapability(ctx, "missing-task", avail.ID); err == nil {
		t.Fatal("require with missing task must fail")
	}
	if _, err := svc.RequireCapability(ctx, task.ID, "missing-cap"); err == nil {
		t.Fatal("require with missing capability must fail")
	}

	r1, err := svc.RequireCapability(ctx, task.ID, avail.ID)
	if err != nil {
		t.Fatalf("RequireCapability avail: %v", err)
	}
	r1b, err := svc.RequireCapability(ctx, task.ID, avail.ID)
	if err != nil || r1b.ID != r1.ID {
		t.Fatalf("Require idempotent: %+v vs %+v err=%v", r1, r1b, err)
	}
	if _, err := svc.RequireCapability(ctx, task.ID, unavail.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RequireCapability(ctx, task.ID, unknown.ID); err != nil {
		t.Fatal(err)
	}

	req, err := svc.ListRequiredCapabilities(ctx, task.ID)
	if err != nil || len(req) != 3 {
		t.Fatalf("ListRequired: %+v err=%v", req, err)
	}

	missing, err := svc.MissingCapabilities(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 2 {
		t.Fatalf("want 2 missing (UNAVAILABLE+UNKNOWN), got %+v", missing)
	}
	for _, m := range missing {
		if m.Status == domain.CapabilityStatusAvailable {
			t.Fatalf("AVAILABLE must not be missing: %+v", m)
		}
	}

	if err := svc.UnrequireCapability(ctx, task.ID, unavail.ID); err != nil {
		t.Fatalf("Unrequire: %v", err)
	}
	if err := svc.UnrequireCapability(ctx, task.ID, unavail.ID); err != nil {
		t.Fatalf("Unrequire no-op: %v", err)
	}
	req, _ = svc.ListRequiredCapabilities(ctx, task.ID)
	if len(req) != 2 {
		t.Fatalf("after unrequire want 2, got %+v", req)
	}
}

func TestBuiltinMCPCapabilitySpecs(t *testing.T) {
	specs := domain.BuiltinMCPCapabilitySpecs()
	want := []string{
		"mcp:trace_why", "mcp:trace_context", "mcp:trace_add",
		"mcp:trace_link", "mcp:trace_transition", "mcp:trace_review",
		"mcp:trace_tasks", "mcp:trace_capability", "mcp:trace_impact",
		"mcp:trace_version", "mcp:trace_search", "mcp:trace_changes", "mcp:trace_regressions",
		"mcp:trace_loop", "mcp:trace_agents", "mcp:trace_plan", "mcp:trace_explore",
	}
	if len(specs) != 17 {
		t.Fatalf("want 17 specs, got %d", len(specs))
	}
	for i, s := range specs {
		if s.Kind != domain.CapabilityKindMCP || s.Status != domain.CapabilityStatusAvailable {
			t.Fatalf("spec %d: %+v", i, s)
		}
		if s.Slug != want[i] || s.Title != want[i][len("mcp:"):] {
			t.Fatalf("spec %d slug/title: %+v want %s", i, s, want[i])
		}
	}
	// no auto-seed: Open leaves catalog empty
	_, st := openDomain(t)
	list, err := st.ListCapabilities(store.CapabilityListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("BuiltinMCP must not auto-seed; got %+v", list)
	}
}

func TestResolveCapabilityIDOrSlug(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	c, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindHook, Slug: "hook:pre", Status: domain.CapabilityStatusAvailable,
	})
	if err != nil {
		t.Fatal(err)
	}
	byID, err := svc.ResolveCapabilityIDOrSlug(ctx, c.ID)
	if err != nil || byID.ID != c.ID {
		t.Fatalf("by id: %+v err=%v", byID, err)
	}
	bySlug, err := svc.ResolveCapabilityIDOrSlug(ctx, "hook:pre")
	if err != nil || bySlug.ID != c.ID {
		t.Fatalf("by slug: %+v err=%v", bySlug, err)
	}
}
