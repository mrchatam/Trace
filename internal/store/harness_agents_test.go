package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHarnessAgentCatalogMigrate027(t *testing.T) {
	s, _ := openTempStore(t)
	st, err := s.MigrationStatus()
	if err != nil {
		t.Fatalf("MigrationStatus: %v", err)
	}
	if st.EmbedExpected != 28 {
		t.Fatalf("EmbedExpected: got %d want 28", st.EmbedExpected)
	}
	if st.MaxApplied != 28 {
		t.Fatalf("MaxApplied: got %d want 28", st.MaxApplied)
	}
	if len(st.AppliedVersions) != 28 {
		t.Fatalf("AppliedVersions len: got %d want 28", len(st.AppliedVersions))
	}

	root := moduleRootFromStoreTest(t)
	schemaDir := filepath.Join(root, "internal", "store", "schema")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		t.Fatalf("ReadDir schema: %v", err)
	}
	sqlCount := 0
	saw027 := false
	saw028 := false
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		sqlCount++
		if e.Name() == "027_harness_agents.sql" {
			saw027 = true
		}
		if e.Name() == "028_deliberation_consecutive_empty.sql" {
			saw028 = true
		}
	}
	if sqlCount != 28 {
		t.Fatalf("schema sql file count: got %d want 28", sqlCount)
	}
	if !saw027 {
		t.Fatal("missing 027_harness_agents.sql")
	}
	if !saw028 {
		t.Fatal("missing 028_deliberation_consecutive_empty.sql")
	}

	agent, err := s.UpsertHarnessAgent(HarnessAgent{
		Slug: "agent:test-migrate", Title: "Test", SubagentType: "generalPurpose",
	})
	if err != nil {
		t.Fatalf("UpsertHarnessAgent: %v", err)
	}
	got, err := s.GetHarnessAgentBySlug("agent:test-migrate")
	if err != nil || got.ID != agent.ID {
		t.Fatalf("GetHarnessAgentBySlug: %+v err=%v", got, err)
	}
	agent2, err := s.UpsertHarnessAgent(HarnessAgent{
		ID: agent.ID, Slug: "agent:test-migrate", Title: "Test again",
	})
	if err != nil || agent2.Title != "Test again" {
		t.Fatalf("idempotent upsert: %+v err=%v", agent2, err)
	}
}

func moduleRootFromStoreTest(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestHarnessAgentRequirementsCRUD(t *testing.T) {
	s, _ := openTempStore(t)
	agent, err := s.UpsertHarnessAgent(HarnessAgent{
		Slug: "agent:req-test", Title: "Req", SubagentType: "explore",
	})
	if err != nil {
		t.Fatalf("UpsertHarnessAgent: %v", err)
	}
	if _, err := s.InsertHarnessAgentRequirement(HarnessAgentRequirement{
		AgentID: agent.ID, RequiredCapabilitySlug: "skill:tdd",
	}); err != nil {
		t.Fatalf("InsertHarnessAgentRequirement: %v", err)
	}
	reqs, err := s.ListHarnessAgentRequirements(agent.ID)
	if err != nil || len(reqs) != 1 || reqs[0].RequiredCapabilitySlug != "skill:tdd" {
		t.Fatalf("ListHarnessAgentRequirements: %+v err=%v", reqs, err)
	}
	if err := s.DeleteHarnessAgentRequirementsForAgent(agent.ID); err != nil {
		t.Fatalf("DeleteHarnessAgentRequirementsForAgent: %v", err)
	}
	reqs, err = s.ListHarnessAgentRequirements(agent.ID)
	if err != nil || len(reqs) != 0 {
		t.Fatalf("after delete want empty: %+v err=%v", reqs, err)
	}
}

func TestHarnessAgentJSONValidation(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertHarnessAgent(HarnessAgent{
		Slug: "agent:bad-json", DeliberationPhases: "{not-array}",
	}); err == nil {
		t.Fatal("invalid deliberation_phases JSON must fail")
	}
	if _, err := s.UpsertHarnessAgent(HarnessAgent{
		Slug: "agent:bad-kw", TaskKeywords: "123",
	}); err == nil {
		t.Fatal("invalid task_keywords JSON must fail")
	}
}
