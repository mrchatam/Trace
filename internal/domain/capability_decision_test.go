package domain_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
	_ "modernc.org/sqlite"
)

func TestCapabilityDecisionAutoAllowBuiltinMCP(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	slug := "mcp:trace_why"

	d, err := svc.ResolveToolDecision(ctx, slug)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if d.Decision != domain.ToolDecisionAutoAllowed {
		t.Fatalf("want AUTO_ALLOWED got %q", d.Decision)
	}
	if !d.Durable {
		t.Fatal("AUTO_ALLOWED must be durable")
	}
	row, err := st.GetCapabilityToolDecisionBySlug(slug)
	if err != nil {
		t.Fatalf("store row: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("store decision: %q", row.Decision)
	}
	if err := svc.AssertToolAllowed(ctx, slug); err != nil {
		t.Fatalf("Assert should pass: %v", err)
	}
}

func TestCapabilityDecisionUnknownPendingFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	slug := "tool:unknown-yolo"

	d, err := svc.ResolveToolDecision(ctx, slug)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if d.Decision != domain.ToolDecisionPending {
		t.Fatalf("want PENDING got %q", d.Decision)
	}
	err = svc.AssertToolAllowed(ctx, slug)
	if err == nil {
		t.Fatal("Assert must fail-closed on PENDING")
	}
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %T %v", err, err)
	}
}

func TestCapabilityDecisionHumanAllowPersists(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	slug := "tool:custom-allow"

	row, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: slug, Decision: "ALLOWED", Reason: "reviewed", Actor: "cli",
	})
	if err != nil {
		t.Fatalf("DecideTool: %v", err)
	}
	if row.Decision != store.ToolDecisionAllowed {
		t.Fatalf("row: %+v", row)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	st2, err := store.Open(root)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	svc2 := domain.New(st2)

	d, err := svc2.ResolveToolDecision(ctx, slug)
	if err != nil {
		t.Fatalf("Resolve after reopen: %v", err)
	}
	if d.Decision != domain.ToolDecisionAllowed {
		t.Fatalf("want ALLOWED survived reopen, got %q", d.Decision)
	}
	if err := svc2.AssertToolAllowed(ctx, slug); err != nil {
		t.Fatalf("Assert: %v", err)
	}
}

func TestCapabilityDecisionDenyBlocks(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	slug := "tool:custom-deny"

	if _, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: slug, Decision: "DENIED", Reason: "unsafe",
	}); err != nil {
		t.Fatalf("DecideTool: %v", err)
	}
	d, err := svc.ResolveToolDecision(ctx, slug)
	if err != nil || d.Decision != domain.ToolDecisionDenied {
		t.Fatalf("Resolve: %+v err=%v", d, err)
	}
	if err := svc.AssertToolAllowed(ctx, slug); err == nil {
		t.Fatal("Assert must fail on DENIED")
	}
}

func TestResolveYOLOBuiltinDoesNotAutoAllow(t *testing.T) {
	root := plantPre014ToolDecision(t, "mcp:trace_add", "YOLO")
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open (014 heal): %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	ctx := context.Background()

	d, err := svc.ResolveToolDecision(ctx, "mcp:trace_add")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if d.Decision != domain.ToolDecisionPending {
		t.Fatalf("want PENDING after YOLO heal, got %q", d.Decision)
	}
	if !d.Durable {
		t.Fatal("healed PENDING must stay durable")
	}
	row, err := st.GetCapabilityToolDecisionBySlug("mcp:trace_add")
	if err != nil {
		t.Fatalf("store row: %v", err)
	}
	if row.Decision != store.ToolDecisionPending {
		t.Fatalf("store must stay PENDING, got %q", row.Decision)
	}
	if row.Decision == store.ToolDecisionAutoAllowed {
		t.Fatal("must not upsert AUTO_ALLOWED over healed YOLO")
	}
	if err := svc.AssertToolAllowed(ctx, "mcp:trace_add"); err == nil {
		t.Fatal("Assert must fail-closed on PENDING")
	}
}

func TestDecideUnprefixedMCPNameCanonicalizes(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	row, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: "trace_why", Decision: "DENIED", Reason: "unprefixed deny",
	})
	if err != nil {
		t.Fatalf("DecideTool: %v", err)
	}
	if row.Slug != "mcp:trace_why" {
		t.Fatalf("persisted slug: got %q want mcp:trace_why", row.Slug)
	}
	if _, err := st.GetCapabilityToolDecisionBySlug("trace_why"); err == nil {
		t.Fatal("must not leave gating-inert trace_why row")
	}
	if err := svc.AssertToolAllowed(ctx, "mcp:trace_why"); err == nil {
		t.Fatal("Assert mcp:trace_why must fail DENIED")
	}
	d, err := svc.ResolveToolDecision(ctx, "mcp:trace_why")
	if err != nil || d.Decision != domain.ToolDecisionDenied {
		t.Fatalf("Resolve mcp:trace_why: %+v err=%v", d, err)
	}
	d2, err := svc.ResolveToolDecision(ctx, "trace_why")
	if err != nil {
		t.Fatalf("Resolve trace_why: %v", err)
	}
	if d2.Slug != "mcp:trace_why" {
		t.Fatalf("Resolve(trace_why) slug: got %q want mcp:trace_why", d2.Slug)
	}
	if d2.Decision != domain.ToolDecisionDenied {
		t.Fatalf("Resolve(trace_why) decision: got %q", d2.Decision)
	}
}

func TestCapabilityDecisionAutoAllowBuiltinCLI(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	d, err := svc.ResolveToolDecision(ctx, "cli:add")
	if err != nil {
		t.Fatalf("Resolve cli:add: %v", err)
	}
	if d.Decision != domain.ToolDecisionAutoAllowed {
		t.Fatalf("want AUTO_ALLOWED got %q", d.Decision)
	}
	if d.Slug != "cli:add" {
		t.Fatalf("slug: got %q want cli:add", d.Slug)
	}
	if !d.Durable {
		t.Fatal("AUTO_ALLOWED must be durable")
	}
	if d.Reason != "builtin CLI command" {
		t.Fatalf("reason: got %q want builtin CLI command", d.Reason)
	}
	row, err := st.GetCapabilityToolDecisionBySlug("cli:add")
	if err != nil {
		t.Fatalf("store row: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("store decision: %q", row.Decision)
	}
	if err := svc.AssertToolAllowed(ctx, "cli:add"); err != nil {
		t.Fatalf("Assert cli:add should pass: %v", err)
	}

	if _, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: "mcp:trace_add", Decision: "DENIED", Reason: "mcp deny must not gate CLI",
	}); err != nil {
		t.Fatalf("DecideTool mcp:trace_add: %v", err)
	}
	d2, err := svc.ResolveToolDecision(ctx, "cli:add")
	if err != nil {
		t.Fatalf("Resolve after MCP DENIED: %v", err)
	}
	if d2.Decision != domain.ToolDecisionAutoAllowed {
		t.Fatalf("cli:add must stay AUTO_ALLOWED after mcp:trace_add DENIED, got %q", d2.Decision)
	}
	if err := svc.AssertToolAllowed(ctx, "cli:add"); err != nil {
		t.Fatalf("Assert cli:add after MCP DENIED: %v", err)
	}
}

func TestUnprefixedAddDecideDoesNotGateCLI(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	row, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: "add", Decision: "DENIED", Reason: "unprefixed custom",
	})
	if err != nil {
		t.Fatalf("DecideTool add: %v", err)
	}
	if row.Slug != "add" {
		t.Fatalf("persisted slug: got %q want add (must not fold to cli:add)", row.Slug)
	}
	got, err := st.GetCapabilityToolDecisionBySlug("add")
	if err != nil {
		t.Fatalf("store add row: %v", err)
	}
	if got.Decision != store.ToolDecisionDenied {
		t.Fatalf("add decision: got %q", got.Decision)
	}
	if _, err := st.GetCapabilityToolDecisionBySlug("cli:add"); err == nil {
		t.Fatal("must not persist cli:add from unprefixed add")
	}
	if err := svc.AssertToolAllowed(ctx, "cli:add"); err != nil {
		t.Fatalf("Assert cli:add must still AUTO_ALLOW: %v", err)
	}
	d, err := svc.ResolveToolDecision(ctx, "cli:add")
	if err != nil || d.Decision != domain.ToolDecisionAutoAllowed {
		t.Fatalf("Resolve cli:add: %+v err=%v", d, err)
	}
}

func TestCanonicalizeCLIReindexFoldsToIndex(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	d, err := svc.ResolveToolDecision(ctx, "cli:reindex")
	if err != nil {
		t.Fatalf("Resolve cli:reindex: %v", err)
	}
	if d.Slug != "cli:index" {
		t.Fatalf("Resolve(cli:reindex) slug: got %q want cli:index", d.Slug)
	}
	if d.Decision != domain.ToolDecisionAutoAllowed {
		t.Fatalf("want AUTO_ALLOWED got %q", d.Decision)
	}
	if _, err := st.GetCapabilityToolDecisionBySlug("cli:reindex"); err == nil {
		t.Fatal("must not persist gating-inert cli:reindex row")
	}
	row, err := st.GetCapabilityToolDecisionBySlug("cli:index")
	if err != nil {
		t.Fatalf("canonical cli:index row: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("cli:index decision: got %q", row.Decision)
	}

	if _, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: "cli:index", Decision: "DENIED", Reason: "alias deny",
	}); err != nil {
		t.Fatalf("DecideTool cli:index: %v", err)
	}
	if err := svc.AssertToolAllowed(ctx, "cli:index"); err == nil {
		t.Fatal("Assert cli:index must fail DENIED")
	}
	if err := svc.AssertToolAllowed(ctx, "cli:reindex"); err == nil {
		t.Fatal("Assert cli:reindex must fail after cli:index DENIED")
	}
	d2, err := svc.ResolveToolDecision(ctx, "cli:reindex")
	if err != nil || d2.Slug != "cli:index" || d2.Decision != domain.ToolDecisionDenied {
		t.Fatalf("Resolve cli:reindex after DENIED: %+v err=%v", d2, err)
	}

	svc2, st2 := openDomain(t)
	dec, err := svc2.DecideTool(ctx, domain.DecideToolInput{
		Slug: "cli:reindex", Decision: "DENIED", Reason: "fold on decide",
	})
	if err != nil {
		t.Fatalf("DecideTool cli:reindex: %v", err)
	}
	if dec.Slug != "cli:index" {
		t.Fatalf("Decide(cli:reindex) persisted %q want cli:index", dec.Slug)
	}
	if _, err := st2.GetCapabilityToolDecisionBySlug("cli:reindex"); err == nil {
		t.Fatal("Decide cli:reindex must not leave cli:reindex row")
	}
}

func TestCanonicalizeCustomAndCLISlugsUnchanged(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	cases := []string{"tool:custom-allow", "cli:add", "trace_*", "mcp:trace_*", "trace_why_extra"}
	for _, slug := range cases {
		row, err := svc.DecideTool(ctx, domain.DecideToolInput{
			Slug: slug, Decision: "ALLOWED", Reason: "keep as-given",
		})
		if err != nil {
			t.Fatalf("DecideTool %q: %v", slug, err)
		}
		if row.Slug != slug {
			t.Fatalf("slug %q canonicalized to %q (must stay unchanged)", slug, row.Slug)
		}
		got, err := st.GetCapabilityToolDecisionBySlug(slug)
		if err != nil {
			t.Fatalf("store %q: %v", slug, err)
		}
		if got.Slug != slug {
			t.Fatalf("stored slug %q want %q", got.Slug, slug)
		}
	}
	for _, slug := range []string{"MCP:trace_why", "Trace_Why"} {
		row, err := svc.DecideTool(ctx, domain.DecideToolInput{
			Slug: slug, Decision: "DENIED", Reason: "case exact",
		})
		if err != nil {
			t.Fatalf("DecideTool %q: %v", slug, err)
		}
		if row.Slug != slug {
			t.Fatalf("case-mismatched %q must not map to builtin, got %q", slug, row.Slug)
		}
	}
}

func plantPre014ToolDecision(t *testing.T, slug, decision string) string {
	t.Helper()
	root := t.TempDir()
	traceDir := filepath.Join(root, ".trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", filepath.Join(traceDir, "trace.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatal(err)
	}
	applyStoreSchemaThrough(t, db, 13)
	now := "2026-01-01T00:00:00Z"
	_, err = db.Exec(`
		INSERT INTO capability_tool_decisions(id, slug, decision, reason, actor, created_at, updated_at)
		VALUES (?, ?, ?, '', 'test', ?, ?)
	`, uuid.NewString(), slug, decision, now, now)
	if err != nil {
		t.Fatalf("plant %s=%s: %v", slug, decision, err)
	}
	return root
}

func applyStoreSchemaThrough(t *testing.T, db *sql.DB, maxVersion int) {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	schemaDir := filepath.Join(filepath.Dir(thisFile), "..", "store", "schema")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		t.Fatalf("schema dir: %v", err)
	}
	type mig struct {
		ver  int
		name string
	}
	var migs []mig
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		parts := strings.SplitN(e.Name(), "_", 2)
		ver, err := strconv.Atoi(parts[0])
		if err != nil {
			t.Fatalf("bad migration name %q", e.Name())
		}
		migs = append(migs, mig{ver: ver, name: e.Name()})
	}
	sort.Slice(migs, func(i, j int) bool { return migs[i].ver < migs[j].ver })
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	for _, m := range migs {
		if m.ver > maxVersion {
			continue
		}
		body, err := os.ReadFile(filepath.Join(schemaDir, m.name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(body)); err != nil {
			t.Fatalf("apply %s: %v", m.name, err)
		}
		if _, err := db.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, m.ver); err != nil {
			t.Fatalf("record %d: %v", m.ver, err)
		}
	}
}
