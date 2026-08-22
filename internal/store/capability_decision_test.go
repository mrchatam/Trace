package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

// applyEmbeddedThrough applies embedded schema files with version <= maxVersion.
func applyEmbeddedThrough(t *testing.T, db *sql.DB, maxVersion int) {
	t.Helper()
	if _, err := db.Exec(migrationTableSQL); err != nil {
		t.Fatalf("schema_migrations: %v", err)
	}
	migs, err := listEmbeddedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range migs {
		if m.version > maxVersion {
			continue
		}
		var applied int
		if err := db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, m.version).Scan(&applied); err != nil {
			t.Fatalf("check migration %d: %v", m.version, err)
		}
		if applied > 0 {
			continue
		}
		body, err := schemaFS.ReadFile("schema/" + m.name)
		if err != nil {
			t.Fatalf("read %s: %v", m.name, err)
		}
		if _, err := db.Exec(string(body)); err != nil {
			t.Fatalf("apply %s: %v", m.name, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, m.version); err != nil {
			t.Fatalf("record %d: %v", m.version, err)
		}
	}
}

func plantPre014Decisions(t *testing.T, rows []CapabilityToolDecision) string {
	t.Helper()
	root := t.TempDir()
	traceDir := filepath.Join(root, traceDirName)
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", filepath.Join(traceDir, dbFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatal(err)
	}
	applyEmbeddedThrough(t, db, 13)
	now := nowRFC3339()
	for _, r := range rows {
		id := r.ID
		if id == "" {
			id = uuid.NewString()
		}
		created := r.CreatedAt
		if created == "" {
			created = now
		}
		updated := r.UpdatedAt
		if updated == "" {
			updated = now
		}
		_, err := db.Exec(`
			INSERT INTO capability_tool_decisions(id, slug, decision, reason, actor, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, id, r.Slug, r.Decision, r.Reason, r.Actor, created, updated)
		if err != nil {
			t.Fatalf("plant %q=%q: %v", r.Slug, r.Decision, err)
		}
	}
	return root
}

func TestCapabilityToolDecisionCheckRejectsYOLO(t *testing.T) {
	s, _ := openTempStore(t)

	_, err := s.db.Exec(`
		INSERT INTO capability_tool_decisions(id, slug, decision, reason, actor, created_at, updated_at)
		VALUES ('yolo-raw', 'mcp:trace_add', 'YOLO', '', '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')
	`)
	if err == nil {
		t.Fatal("raw INSERT YOLO must fail CHECK after 014")
	}

	_, err = s.db.Exec(`
		INSERT INTO capability_tool_decisions(id, slug, decision, reason, actor, created_at, updated_at)
		VALUES ('empty-raw', 'mcp:trace_link', '', '', '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')
	`)
	if err == nil {
		t.Fatal("raw INSERT empty decision must fail CHECK after 014")
	}

	_, err = s.UpsertCapabilityToolDecision(CapabilityToolDecision{
		Slug: "mcp:trace_tasks", Decision: "YOLO",
	})
	if err == nil {
		t.Fatal("Upsert YOLO must error")
	}
	_, err = s.UpsertCapabilityToolDecision(CapabilityToolDecision{
		Slug: "mcp:trace_review", Decision: "   ",
	})
	if err == nil {
		t.Fatal("Upsert empty/whitespace decision must error")
	}

	for _, d := range []string{
		ToolDecisionAutoAllowed, ToolDecisionPending, ToolDecisionAllowed, ToolDecisionDenied,
	} {
		row, err := s.UpsertCapabilityToolDecision(CapabilityToolDecision{
			Slug: "enum-" + d, Decision: d, Reason: "ok",
		})
		if err != nil {
			t.Fatalf("upsert %s: %v", d, err)
		}
		if row.Decision != d {
			t.Fatalf("upsert %s: stored %q", d, row.Decision)
		}
	}
}

func TestCapabilityToolDecisionMigrateHealsYOLOToPending(t *testing.T) {
	root := plantPre014Decisions(t, []CapabilityToolDecision{{
		Slug: "mcp:trace_add", Decision: "YOLO", Reason: "hunt garbage", Actor: "test",
	}})
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open after 013 YOLO plant: %v", err)
	}
	defer s.Close()

	row, err := s.GetCapabilityToolDecisionBySlug("mcp:trace_add")
	if err != nil {
		t.Fatalf("row dropped (must heal, not DROP): %v", err)
	}
	if row.Decision == ToolDecisionAutoAllowed {
		t.Fatal("YOLO must not heal to AUTO_ALLOWED")
	}
	if row.Decision != ToolDecisionPending {
		t.Fatalf("YOLO heal: got %q want PENDING", row.Decision)
	}
}

func TestMigrateUnprefixedDeniedFoldsOverAutoAllowed(t *testing.T) {
	root := plantPre014Decisions(t, []CapabilityToolDecision{
		{Slug: "trace_why", Decision: ToolDecisionDenied, Reason: "unprefixed deny", Actor: "cli"},
		{Slug: "mcp:trace_why", Decision: ToolDecisionAutoAllowed, Reason: "builtin MCP capability", Actor: "system"},
	})
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open after footgun plant: %v", err)
	}
	defer s.Close()

	all, err := s.ListCapabilityToolDecisions()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatalf("want one canonical row after fold, got %d: %+v", len(all), all)
	}
	if all[0].Slug != "mcp:trace_why" {
		t.Fatalf("canonical slug: got %q", all[0].Slug)
	}
	if all[0].Decision != ToolDecisionDenied {
		t.Fatalf("fold priority DENIED > AUTO_ALLOWED: got %q", all[0].Decision)
	}
	if _, err := s.GetCapabilityToolDecisionBySlug("trace_why"); err == nil {
		t.Fatal("unprefixed row must be dropped after fold")
	}
}
