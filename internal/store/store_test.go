package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func openTempStore(t *testing.T) (*Store, string) {
	t.Helper()
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s, root
}

func TestOpenExistingMissingReturnsErrNotInitialized(t *testing.T) {
	root := t.TempDir()
	s, err := OpenExisting(root)
	if s != nil {
		t.Fatal("expected nil store")
	}
	if !errors.Is(err, ErrNotInitialized) {
		t.Fatalf("want ErrNotInitialized, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, traceDirName)); !os.IsNotExist(err) {
		t.Fatalf("OpenExisting must not mkdir .trace/: %v", err)
	}
}

func TestOpenExistingEmptyTraceDir(t *testing.T) {
	root := t.TempDir()
	traceDir := filepath.Join(root, traceDirName)
	if err := os.Mkdir(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := OpenExisting(root)
	if s != nil {
		t.Fatal("expected nil store")
	}
	if !errors.Is(err, ErrNotInitialized) {
		t.Fatalf("want ErrNotInitialized, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(traceDir, dbFileName)); !os.IsNotExist(err) {
		t.Fatalf("OpenExisting must not create trace.db: %v", err)
	}
}

func TestOpenCreatesDBAndMigratesIdempotent(t *testing.T) {
	root := t.TempDir()

	s1, err := Open(root)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	dbPath := s1.DBPath()
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected db at %s: %v", dbPath, err)
	}
	if s1.ProjectID() == "" {
		t.Fatal("expected project id")
	}
	pid1 := s1.ProjectID()
	if err := s1.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	s2, err := Open(root)
	if err != nil {
		t.Fatalf("second Open (idempotent migrate): %v", err)
	}
	defer s2.Close()
	if s2.ProjectID() != pid1 {
		t.Fatalf("project id changed on re-open: %q vs %q", pid1, s2.ProjectID())
	}

	for _, want := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24} {
		var ver int
		if err := s2.db.QueryRow(`SELECT version FROM schema_migrations WHERE version = ?`, want).Scan(&ver); err != nil {
			t.Fatalf("migration v%d missing: %v", want, err)
		}
		if ver != want {
			t.Fatalf("want version %d, got %d", want, ver)
		}
	}

	required := []string{
		"schema_migrations", "projects",
		"goals", "decisions", "assumptions", "tasks", "discoveries",
		"claims", "evidence", "reviews", "plan_changes",
		"events", "files", "symbols", "imports",
		"vcs_commits", "vcs_commit_paths", "vcs_meta",
		"entity_links",
		"fts_docs",
		"plan_phases", "plan_scopes", "scope_deep_plans", "goal_plan_state",
		"review_residuals",
		"decision_impact_findings", "decision_alternatives",
		"capabilities", "task_capability_requirements",
		"capability_tool_decisions",
		"deliberation_state",
		"uncertainties", "hypotheses", "decision_reconsiderations",
		"changes", "change_paths", "effects",
		"outcome_results", "baselines",
		"regressions", "reflections",
		"experiments",
		"code_edges",
		"graph_sync_state",
		"impact_predictions",
		"improvements",
	}
	for _, table := range required {
		var name string
		err := s2.db.QueryRow(
			`SELECT name FROM sqlite_master WHERE name=?`, table,
		).Scan(&name)
		if err != nil {
			t.Errorf("missing table %s: %v", table, err)
		}
	}
}

func TestGoalTaskEventRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	goal, err := s.UpsertGoal(Goal{
		Title:      "Ship store",
		Body:       "SQLite under .trace",
		SourceType: "USER_ASSERTED",
		Confidence: 0.9,
		Status:     StatusActive,
	})
	if err != nil {
		t.Fatalf("UpsertGoal: %v", err)
	}
	gotGoal, err := s.GetGoal(goal.ID)
	if err != nil {
		t.Fatalf("GetGoal: %v", err)
	}
	if gotGoal.Title != "Ship store" || gotGoal.SourceType != "USER_ASSERTED" || gotGoal.Status != StatusActive {
		t.Fatalf("goal provenance/roundtrip mismatch: %+v", gotGoal)
	}
	if gotGoal.Confidence != 0.9 {
		t.Fatalf("confidence: got %v", gotGoal.Confidence)
	}

	goalID := goal.ID
	task, err := s.UpsertTask(Task{
		GoalID:     &goalID,
		Title:      "Write migrations",
		Body:       "v1 schema",
		SourceType: "AGENT_PROPOSED",
		Confidence: 0.7,
		Status:     StatusActive,
	})
	if err != nil {
		t.Fatalf("UpsertTask: %v", err)
	}
	gotTask, err := s.GetTask(task.ID)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if gotTask.Title != "Write migrations" || gotTask.SourceType != "AGENT_PROPOSED" {
		t.Fatalf("task roundtrip mismatch: %+v", gotTask)
	}
	if gotTask.GoalID == nil || *gotTask.GoalID != goal.ID {
		t.Fatalf("task goal_id: %+v", gotTask.GoalID)
	}
	if gotTask.WorkState != WorkStatePending {
		t.Fatalf("task work_state default: got %q", gotTask.WorkState)
	}

	ev, err := s.AppendEvent(Event{
		Type:        "task.created",
		EntityType:  "task",
		EntityID:    task.ID,
		PayloadJSON: `{"title":"Write migrations"}`,
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	listed, err := s.ListEventsByEntity("task", task.ID)
	if err != nil {
		t.Fatalf("ListEventsByEntity: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != ev.ID {
		t.Fatalf("events visible: %+v", listed)
	}
	recent, err := s.ListRecentEvents(10)
	if err != nil || len(recent) < 1 {
		t.Fatalf("ListRecentEvents: %v %#v", err, recent)
	}
}

func TestReplaceFileSymbolsIncrementalIsolation(t *testing.T) {
	s, _ := openTempStore(t)

	if _, err := s.UpsertFile("a.go", "hash-a", nil); err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	if _, err := s.UpsertFile("b.go", "hash-b", nil); err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}

	if err := s.ReplaceFileSymbols("a.go", []Symbol{
		{Name: "A1", Kind: "func", StartLine: 1, EndLine: 5},
		{Name: "A2", Kind: "type", StartLine: 7, EndLine: 10},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols a: %v", err)
	}
	if err := s.ReplaceFileSymbols("b.go", []Symbol{
		{Name: "B1", Kind: "func", StartLine: 1, EndLine: 3},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols b: %v", err)
	}

	// Replace only A — B must remain.
	if err := s.ReplaceFileSymbols("a.go", []Symbol{
		{Name: "A3", Kind: "func", StartLine: 1, EndLine: 2},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols a again: %v", err)
	}

	symsA, err := s.ListSymbolsByPath("a.go")
	if err != nil {
		t.Fatalf("ListSymbolsByPath a: %v", err)
	}
	if len(symsA) != 1 || symsA[0].Name != "A3" {
		t.Fatalf("symbols for a: %+v", symsA)
	}

	symsB, err := s.ListSymbolsByPath("b.go")
	if err != nil {
		t.Fatalf("ListSymbolsByPath b: %v", err)
	}
	if len(symsB) != 1 || symsB[0].Name != "B1" {
		t.Fatalf("symbols for b were disturbed: %+v", symsB)
	}

	// Imports same isolation pattern.
	symName := "Fmt"
	if err := s.ReplaceFileImports("a.go", []Import{{ImportedPath: "fmt"}}); err != nil {
		t.Fatalf("imports a: %v", err)
	}
	if err := s.ReplaceFileImports("b.go", []Import{{ImportedPath: "os", Symbol: &symName}}); err != nil {
		t.Fatalf("imports b: %v", err)
	}
	if err := s.ReplaceFileImports("a.go", []Import{{ImportedPath: "io"}}); err != nil {
		t.Fatalf("imports a replace: %v", err)
	}
	impsB, err := s.ListImportsByPath("b.go")
	if err != nil || len(impsB) != 1 || impsB[0].ImportedPath != "os" {
		t.Fatalf("imports for b disturbed: %v %+v", err, impsB)
	}

	// File upsert updates hash without dropping symbols of other files.
	oid := "abc123"
	if _, err := s.UpsertFile("a.go", "hash-a2", &oid); err != nil {
		t.Fatalf("re-upsert a: %v", err)
	}
	symsB, err = s.ListSymbolsByPath("b.go")
	if err != nil || len(symsB) != 1 {
		t.Fatalf("b symbols after a upsert: %v %+v", err, symsB)
	}
}

func TestImportProvenanceRoundTrip(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("a.go", "ha", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}

	// Empty Provenance defaults to EXTRACTED.
	if err := s.ReplaceFileImports("a.go", []Import{
		{ImportedPath: "fmt"},
	}); err != nil {
		t.Fatalf("Replace empty provenance: %v", err)
	}
	imps, err := s.ListImportsByPath("a.go")
	if err != nil || len(imps) != 1 {
		t.Fatalf("list: %v %+v", err, imps)
	}
	if imps[0].Provenance != ImportProvenanceExtracted {
		t.Fatalf("default provenance: got %q want %q", imps[0].Provenance, ImportProvenanceExtracted)
	}

	// Explicit round-trip of all three enum values.
	if err := s.ReplaceFileImports("a.go", []Import{
		{ImportedPath: "extracted/pkg", Provenance: ImportProvenanceExtracted},
		{ImportedPath: "inferred/pkg", Provenance: ImportProvenanceInferred},
		{ImportedPath: "ambiguous/pkg", Provenance: ImportProvenanceAmbiguous},
	}); err != nil {
		t.Fatalf("Replace three: %v", err)
	}
	imps, err = s.ListImportsByPath("a.go")
	if err != nil {
		t.Fatalf("list three: %v", err)
	}
	got := map[string]string{}
	for _, im := range imps {
		got[im.ImportedPath] = im.Provenance
	}
	want := map[string]string{
		"extracted/pkg": ImportProvenanceExtracted,
		"inferred/pkg":  ImportProvenanceInferred,
		"ambiguous/pkg": ImportProvenanceAmbiguous,
	}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("path %s: got %q want %q", k, got[k], v)
		}
	}
}

func TestListImportEdges(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("a.go", "ha", nil); err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	if _, err := s.UpsertFile("b.go", "hb", nil); err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	if err := s.ReplaceFileImports("b.go", []Import{
		{ImportedPath: "a.go", Provenance: ImportProvenanceInferred},
	}); err != nil {
		t.Fatalf("ReplaceFileImports: %v", err)
	}
	edges, err := s.ListImportEdges()
	if err != nil {
		t.Fatalf("ListImportEdges: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("edges: %+v", edges)
	}
	if edges[0].ImporterPath != "b.go" || edges[0].ImportedPath != "a.go" {
		t.Fatalf("edge paths: %+v", edges[0])
	}
	if edges[0].Provenance != ImportProvenanceInferred {
		t.Fatalf("provenance: %q", edges[0].Provenance)
	}
}

// DF-64: garbage provenance must error on write (no silent coerce).
func TestReplaceFileImportsRejectsGarbageProvenance(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("a.go", "ha", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	err := s.ReplaceFileImports("a.go", []Import{
		{ImportedPath: "fmt", Provenance: "MADE_UP"},
	})
	if err == nil {
		t.Fatal("expected error for garbage provenance")
	}
	if !strings.Contains(err.Error(), "invalid import provenance") {
		t.Fatalf("error should name invalid provenance: %v", err)
	}
	imps, listErr := s.ListImportsByPath("a.go")
	if listErr != nil {
		t.Fatalf("list after reject: %v", listErr)
	}
	if len(imps) != 0 {
		t.Fatalf("reject must not leave partial inserts; got %+v", imps)
	}
}

// DF-64: empty write → stored EXTRACTED; read normalize maps empty → EXTRACTED.
func TestImportProvenanceEmptyWriteAndReadNormalize(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("a.go", "ha", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := s.ReplaceFileImports("a.go", []Import{
		{ImportedPath: "fmt"}, // empty Provenance
	}); err != nil {
		t.Fatalf("Replace empty: %v", err)
	}
	var raw string
	if err := s.db.QueryRow(`SELECT provenance FROM imports WHERE imported_path = ?`, "fmt").Scan(&raw); err != nil {
		t.Fatalf("raw select: %v", err)
	}
	if raw != ImportProvenanceExtracted {
		t.Fatalf("stored provenance: got %q want EXTRACTED", raw)
	}
	imps, err := s.ListImportsByPath("a.go")
	if err != nil || len(imps) != 1 {
		t.Fatalf("list: %v %+v", err, imps)
	}
	if imps[0].Provenance != ImportProvenanceExtracted {
		t.Fatalf("list provenance: got %q want EXTRACTED", imps[0].Provenance)
	}
	if got := normalizeImportProvenance(""); got != ImportProvenanceExtracted {
		t.Fatalf("normalize empty: got %q want EXTRACTED", got)
	}
	// CHECK rejects garbage at SQL layer after mig 012.
	_, err = s.db.Exec(`
		INSERT INTO imports(id, file_id, imported_path, symbol, provenance)
		VALUES ('x', (SELECT id FROM files WHERE path = 'a.go'), 'bad', NULL, 'MADE_UP')
	`)
	if err == nil {
		t.Fatal("expected CHECK failure for MADE_UP provenance")
	}
}

func TestGetSymbolByID(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("pkg/x.go", "hx", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := s.ReplaceFileSymbols("pkg/x.go", []Symbol{
		{Name: "ByID", Kind: "func", StartLine: 2, EndLine: 5},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols: %v", err)
	}
	listed, err := s.ListSymbolsByPath("pkg/x.go")
	if err != nil || len(listed) != 1 {
		t.Fatalf("list: %v %+v", err, listed)
	}
	sym, path, err := s.GetSymbolByID(listed[0].ID)
	if err != nil {
		t.Fatalf("GetSymbolByID: %v", err)
	}
	if sym.Name != "ByID" || sym.Kind != "func" || path != "pkg/x.go" {
		t.Fatalf("got %+v path=%q", sym, path)
	}
	if _, _, err := s.GetSymbolByID("00000000-0000-0000-0000-000000000000"); err == nil {
		t.Fatal("expected miss")
	}
}

func TestSetFileLanguage(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertFile("src/a.ts", "hash", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := s.SetFileLanguage("src/a.ts", "typescript"); err != nil {
		t.Fatalf("SetFileLanguage: %v", err)
	}
	f, err := s.GetFileByPath("src/a.ts")
	if err != nil {
		t.Fatalf("GetFileByPath: %v", err)
	}
	if f.Language == nil || *f.Language != "typescript" {
		t.Fatalf("language: got %v want typescript", f.Language)
	}
	if err := s.SetFileLanguage("missing.ts", "typescript"); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestListFilePathsAndDeleteFileByPath(t *testing.T) {
	s, _ := openTempStore(t)

	if _, err := s.UpsertFile("keep.js", "hash-keep", nil); err != nil {
		t.Fatalf("UpsertFile keep: %v", err)
	}
	if _, err := s.UpsertFile("gone.js", "hash-gone", nil); err != nil {
		t.Fatalf("UpsertFile gone: %v", err)
	}
	if err := s.ReplaceFileSymbols("gone.js", []Symbol{
		{Name: "GhostFn", Kind: "function", StartLine: 1, EndLine: 2},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols: %v", err)
	}
	symName := "named"
	if err := s.ReplaceFileImports("gone.js", []Import{
		{ImportedPath: "./x", Symbol: &symName},
	}); err != nil {
		t.Fatalf("ReplaceFileImports: %v", err)
	}

	paths, err := s.ListFilePaths()
	if err != nil {
		t.Fatalf("ListFilePaths: %v", err)
	}
	if len(paths) != 2 || paths[0] != "gone.js" || paths[1] != "keep.js" {
		t.Fatalf("ListFilePaths order/content: %v", paths)
	}

	hits, err := s.SearchFTS("GhostFn", 10)
	if err != nil {
		t.Fatalf("SearchFTS before delete: %v", err)
	}
	foundSym := false
	for _, h := range hits {
		if h.EntityType == "symbol" && h.SymbolName == "GhostFn" {
			foundSym = true
			break
		}
	}
	if !foundSym {
		t.Fatalf("expected GhostFn in FTS before delete, got %+v", hits)
	}

	if err := s.DeleteFileByPath("gone.js"); err != nil {
		t.Fatalf("DeleteFileByPath: %v", err)
	}
	if _, err := s.GetFileByPath("gone.js"); err == nil {
		t.Fatal("expected GetFileByPath error after delete")
	}
	if _, err := s.ListSymbolsByPath("gone.js"); err == nil {
		t.Fatal("expected ListSymbolsByPath error after delete")
	}
	if _, err := s.ListImportsByPath("gone.js"); err == nil {
		t.Fatal("expected ListImportsByPath error after delete")
	}

	afterHits, err := s.SearchFTS("GhostFn", 10)
	if err != nil {
		t.Fatalf("SearchFTS after delete: %v", err)
	}
	for _, h := range afterHits {
		if h.EntityType == "symbol" && h.SymbolName == "GhostFn" {
			t.Fatalf("GhostFn still in FTS after delete: %+v", h)
		}
		if h.EntityType == "file" && h.Path == "gone.js" {
			t.Fatalf("gone.js still in FTS after delete: %+v", h)
		}
	}

	if _, err := s.GetFileByPath("keep.js"); err != nil {
		t.Fatalf("sibling keep.js missing: %v", err)
	}
	if err := s.DeleteFileByPath("gone.js"); err != nil {
		t.Fatalf("idempotent DeleteFileByPath: %v", err)
	}
	if err := s.DeleteFileByPath("never-existed.js"); err != nil {
		t.Fatalf("delete missing path: %v", err)
	}

	paths2, err := s.ListFilePaths()
	if err != nil {
		t.Fatalf("ListFilePaths after: %v", err)
	}
	if len(paths2) != 1 || paths2[0] != "keep.js" {
		t.Fatalf("ListFilePaths after delete: %v", paths2)
	}
}

func TestListFilePathsByContentHash(t *testing.T) {
	s, _ := openTempStore(t)

	if _, err := s.UpsertFile("a.js", "hash-same", nil); err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	if _, err := s.UpsertFile("c.js", "hash-same", nil); err != nil {
		t.Fatalf("UpsertFile c: %v", err)
	}
	if _, err := s.UpsertFile("b.js", "hash-other", nil); err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}

	got, err := s.ListFilePathsByContentHash("hash-same")
	if err != nil {
		t.Fatalf("ListFilePathsByContentHash: %v", err)
	}
	if len(got) != 2 || got[0] != "a.js" || got[1] != "c.js" {
		t.Fatalf("expected ordered a.js,c.js; got %v", got)
	}

	empty, err := s.ListFilePathsByContentHash("no-such-hash")
	if err != nil {
		t.Fatalf("empty hash list: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty for unknown hash, got %v", empty)
	}
	nilHash, err := s.ListFilePathsByContentHash("")
	if err != nil {
		t.Fatalf("empty string hash: %v", err)
	}
	if len(nilHash) != 0 {
		t.Fatalf("expected nil/empty for blank hash, got %v", nilHash)
	}
}

func TestNoSourceContentColumns(t *testing.T) {
	s, _ := openTempStore(t)

	forbiddenSubstr := []string{
		"content_blob", "source_blob", "file_body", "source_text",
		"file_content", "raw_source", "blob",
	}
	// Allow content_hash (hash only, not body).
	rows, err := s.db.Query(`
		SELECT m.name AS table_name, p.name AS col_name
		FROM sqlite_master m
		JOIN pragma_table_info(m.name) p
		WHERE m.type = 'table'
	`)
	if err != nil {
		t.Fatalf("pragma: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table, col string
		if err := rows.Scan(&table, &col); err != nil {
			t.Fatal(err)
		}
		lower := strings.ToLower(col)
		if lower == "content_hash" || lower == "payload_json" || lower == "body" {
			// body on causal entities / payload_json on events are allowed; not source files.
			continue
		}
		for _, bad := range forbiddenSubstr {
			if strings.Contains(lower, bad) {
				t.Errorf("forbidden source-content-ish column %s.%s", table, col)
			}
		}
	}

	// Explicit: files table must not have a content/body/blob column.
	cols, err := fileTableColumns(s.conn)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cols {
		switch strings.ToLower(c) {
		case "content", "body", "blob", "source", "data", "bytes", "raw":
			t.Errorf("files has forbidden column %q", c)
		}
	}
}

func fileTableColumns(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`PRAGMA table_info(files)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cols []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}
	return cols, rows.Err()
}

func TestTraceDirPermissionsAndPath(t *testing.T) {
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	info, err := os.Stat(filepath.Join(root, ".trace"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal(".trace not a dir")
	}
	want := filepath.Join(root, ".trace", "trace.db")
	if s.DBPath() != want {
		t.Fatalf("DBPath: got %s want %s", s.DBPath(), want)
	}
}

func TestCausalEntitiesAndLinks(t *testing.T) {
	s, _ := openTempStore(t)

	d, err := s.UpsertDecision(Decision{Title: "Use SQLite", Body: "local-first", SourceType: "USER_ASSERTED", Confidence: 0.95})
	if err != nil {
		t.Fatalf("UpsertDecision: %v", err)
	}
	got, err := s.GetDecision(d.ID)
	if err != nil || got.Title != "Use SQLite" || got.Status != StatusActive {
		t.Fatalf("decision roundtrip: %+v err=%v", got, err)
	}

	task, err := s.UpsertTask(Task{Title: "Wire store", WorkState: WorkStateInProgress})
	if err != nil {
		t.Fatalf("UpsertTask: %v", err)
	}
	if task.WorkState != WorkStateInProgress {
		t.Fatalf("work_state: %q", task.WorkState)
	}

	link, err := s.InsertLink(EntityLink{
		FromType:   "decision",
		FromID:     d.ID,
		Rel:        "decision_affects_task",
		ToType:     "task",
		ToID:       task.ID,
		SourceType: "USER_ASSERTED",
		Confidence: 1,
	})
	if err != nil {
		t.Fatalf("InsertLink: %v", err)
	}
	from, err := s.ListLinksFrom("decision", d.ID)
	if err != nil || len(from) != 1 || from[0].ID != link.ID {
		t.Fatalf("ListLinksFrom: %+v err=%v", from, err)
	}
	byRel, err := s.ListLinksByRel("decision_affects_task")
	if err != nil || len(byRel) != 1 || byRel[0].ID != link.ID {
		t.Fatalf("ListLinksByRel: %+v err=%v", byRel, err)
	}

	var ws string
	if err := s.db.QueryRow(`SELECT work_state FROM tasks WHERE id = ?`, task.ID).Scan(&ws); err != nil {
		t.Fatalf("work_state column: %v", err)
	}
	if ws != WorkStateInProgress {
		t.Fatalf("column work_state=%q", ws)
	}
}

func TestReviewUpsertGetResult(t *testing.T) {
	s, _ := openTempStore(t)

	r, err := s.UpsertReview(Review{
		Title:      "Independent review",
		Body:       "checks DONE gate",
		SourceType: "USER_ASSERTED",
		Confidence: 1,
		Result:     ReviewResultOpen,
	})
	if err != nil {
		t.Fatalf("UpsertReview: %v", err)
	}
	if r.ID == "" || r.Result != ReviewResultOpen {
		t.Fatalf("want open result, got %+v", r)
	}

	got, err := s.GetReview(r.ID)
	if err != nil {
		t.Fatalf("GetReview: %v", err)
	}
	if got.Title != "Independent review" || got.Result != "" {
		t.Fatalf("GetReview: %+v", got)
	}

	got.Result = ReviewResultPass
	updated, err := s.UpsertReview(got)
	if err != nil {
		t.Fatalf("UpsertReview PASS: %v", err)
	}
	if updated.Result != ReviewResultPass {
		t.Fatalf("want PASS, got %q", updated.Result)
	}

	again, err := s.GetReview(r.ID)
	if err != nil || again.Result != ReviewResultPass {
		t.Fatalf("persisted result: %+v err=%v", again, err)
	}

	var col string
	if err := s.db.QueryRow(`SELECT result FROM reviews WHERE id = ?`, r.ID).Scan(&col); err != nil {
		t.Fatalf("result column missing (mig 005?): %v", err)
	}
	if col != ReviewResultPass {
		t.Fatalf("raw result=%q", col)
	}
}

func TestReviewResidualsCRUD(t *testing.T) {
	s, _ := openTempStore(t)

	rev, err := s.UpsertReview(Review{Title: "r", SourceType: "USER_ASSERTED"})
	if err != nil {
		t.Fatalf("UpsertReview: %v", err)
	}
	res, err := s.InsertReviewResidual(ReviewResidual{
		ReviewID: rev.ID,
		Code:     "OPEN_GAP",
		Body:     "gap",
	})
	if err != nil {
		t.Fatalf("InsertReviewResidual: %v", err)
	}
	if res.Severity != ResidualSeverityINFO || res.Status != ResidualStatusOpen {
		t.Fatalf("defaults: %+v", res)
	}
	got, err := s.GetReviewResidual(res.ID)
	if err != nil || got.Code != "OPEN_GAP" {
		t.Fatalf("GetReviewResidual: %+v err=%v", got, err)
	}
	list, err := s.ListReviewResidualsByReviewID(rev.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("ListReviewResidualsByReviewID: %+v err=%v", list, err)
	}
	if err := s.UpdateReviewResidualStatus(res.ID, ResidualStatusResolved); err != nil {
		t.Fatalf("UpdateReviewResidualStatus: %v", err)
	}
	again, _ := s.GetReviewResidual(res.ID)
	if again.Status != ResidualStatusResolved {
		t.Fatalf("status=%q", again.Status)
	}
}

func TestDecisionImpactFindingsAndAlternativesCRUD(t *testing.T) {
	s, _ := openTempStore(t)

	d, err := s.UpsertDecision(Decision{Title: "d", SourceType: "USER_ASSERTED"})
	if err != nil {
		t.Fatalf("UpsertDecision: %v", err)
	}
	f, err := s.InsertDecisionImpactFinding(DecisionImpactFinding{
		DecisionID:  d.ID,
		ImpactClass: ImpactClassHigh,
		Kind:        FindingKindAffectedWork,
		Body:        "blast",
	})
	if err != nil {
		t.Fatalf("InsertDecisionImpactFinding: %v", err)
	}
	if f.Uncertainty != UncertaintyUNKNOWN {
		t.Fatalf("default uncertainty: %+v", f)
	}
	flist, err := s.ListDecisionImpactFindingsByDecisionID(d.ID)
	if err != nil || len(flist) != 1 || flist[0].ID != f.ID {
		t.Fatalf("ListDecisionImpactFindingsByDecisionID: %+v err=%v", flist, err)
	}

	a, err := s.InsertDecisionAlternative(DecisionAlternative{
		DecisionID: d.ID, Title: "keep", Body: "status quo",
	})
	if err != nil {
		t.Fatalf("InsertDecisionAlternative: %v", err)
	}
	if a.IsRecommended {
		t.Fatalf("default recommended false: %+v", a)
	}
	if err := s.ClearRecommendedAlternatives(d.ID); err != nil {
		t.Fatalf("ClearRecommendedAlternatives: %v", err)
	}
	if err := s.UpdateDecisionAlternativeRecommended(a.ID, true); err != nil {
		t.Fatalf("UpdateDecisionAlternativeRecommended: %v", err)
	}
	got, err := s.GetDecisionAlternative(a.ID)
	if err != nil || !got.IsRecommended {
		t.Fatalf("GetDecisionAlternative recommended: %+v err=%v", got, err)
	}
	alist, err := s.ListDecisionAlternativesByDecisionID(d.ID)
	if err != nil || len(alist) != 1 {
		t.Fatalf("ListDecisionAlternativesByDecisionID: %+v err=%v", alist, err)
	}
}
