package compat_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "modernc.org/sqlite"
)

// MetricsCompat is the Phase 08 checklist artifact shape (schema-compat.json v1).
type MetricsCompat struct {
	SchemaVersion               int      `json:"schema_version"`
	Gate                        string   `json:"gate"`
	Suite                       string   `json:"suite"`
	DryRun                      bool     `json:"dry_run"`
	NamedTest                   string   `json:"named_test"`
	LanguageAdapterAPIVersion   int      `json:"language_adapter_api_version"`
	LanguageAdapterAPIVersionOK bool     `json:"language_adapter_api_version_ok"`
	PathLocalBindOK             bool     `json:"path_local_bind_ok"`
	TraceLockOK                 bool     `json:"trace_lock_ok"`
	MigrateStatusOK             bool     `json:"migrate_status_ok"`
	BackupRestoreOK             bool     `json:"backup_restore_ok"`
	NoBlobColumnsOK             bool     `json:"no_blob_columns_ok"`
	LocalAuthFailClosedOK       bool     `json:"local_auth_fail_closed_ok"`
	G19OK                       bool     `json:"g19_ok"`
	NoDaemonHTTPPrimaryOK       bool     `json:"no_daemon_http_primary_ok"`
	No011MigOK                  bool     `json:"no_011_mig_ok"`
	S01Hooks                    []string `json:"s01_hooks"`
	S02Hooks                    []string `json:"s02_hooks"`
	S03Hooks                    []string `json:"s03_hooks"`
	TraceVersion                string   `json:"trace_version,omitempty"`
}

func findModuleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	root := findModuleRoot()
	if root == "" {
		t.Fatal("go.mod not found above evals/compat")
	}
	return root
}

// TestCompatibilitySecurityChecklist is the Phase 08 planted compat+security gate.
func TestCompatibilitySecurityChecklist(t *testing.T) {
	root := moduleRoot(t)

	apiVer := analyzers.LanguageAdapterAPIVersion
	apiOK := apiVer == 1
	if !apiOK {
		t.Fatalf("LanguageAdapterAPIVersion=%d want 1", apiVer)
	}

	pathLocalOK := checkPathLocalBind(t)
	lockOK := checkTraceLock(t)
	migrateOK, no011OK := checkMigrateStatus(t, root)
	backupOK, noBlobOK := checkBackupRestore(t)
	authOK := checkLocalAuthFailClosed(t)
	g19OK := checkG19(t, root)
	noDaemonOK := checkNoDaemonHTTPPrimary(t, root)

	if !pathLocalOK || !lockOK || !migrateOK || !backupOK || !noBlobOK || !authOK || !g19OK || !noDaemonOK || !no011OK {
		t.Fatalf("checklist flags failed: path_local=%v lock=%v migrate=%v backup=%v no_blob=%v auth=%v g19=%v no_daemon=%v no_011=%v",
			pathLocalOK, lockOK, migrateOK, backupOK, noBlobOK, authOK, g19OK, noDaemonOK, no011OK)
	}

	metrics := MetricsCompat{
		SchemaVersion:               1,
		Gate:                        "compat-security",
		Suite:                       "compat",
		DryRun:                      false,
		NamedTest:                   "TestCompatibilitySecurityChecklist",
		LanguageAdapterAPIVersion:   apiVer,
		LanguageAdapterAPIVersionOK: apiOK,
		PathLocalBindOK:             pathLocalOK,
		TraceLockOK:                 lockOK,
		MigrateStatusOK:             migrateOK,
		BackupRestoreOK:             backupOK,
		NoBlobColumnsOK:             noBlobOK,
		LocalAuthFailClosedOK:       authOK,
		G19OK:                       g19OK,
		NoDaemonHTTPPrimaryOK:       noDaemonOK,
		No011MigOK:                  no011OK,
		S01Hooks: []string{
			"LanguageAdapterAPIVersion",
			"LanguageAdapter",
			"builtinAdapters (static)",
			"TestLanguageAdapterAPIVersion",
			"TestBuiltinLanguageAdaptersContributionPath",
		},
		S02Hooks: []string{
			"store.Open path-local Abs→<root>/.trace/",
			"exclusive .trace/trace.lock",
			"store.ErrLocked",
			"TestProjectBindPathLocalIsolation",
			"TestConcurrentStoreOpenFailClosed",
		},
		S03Hooks: []string{
			"MigrationStatus / migrate status",
			"BackupTo VACUUM INTO + Restore Abs rebind",
			"HasBlobLikeColumns false",
			"access.token + TRACE_ACCESS_TOKEN → ErrUnauthorized",
			"011_import_edge_provenance + 012_import_provenance_enum + 013_capability_tool_decisions + 014_capability_tool_decision_enum + 015_deliberation_state + 016_cognitive_artifacts + 017_changes_effects + 018_outcome_results_baselines + 019_regressions_reflections + 020_baselines_promotion + 021_experiments + 022_code_relationships + 023_graph_sync + 024_impact_compare + 025_engineering_knowledge + 026_eval_rules + 027_harness_agents present; no 028+",
			"no MCP auth/backup tools",
		},
		TraceVersion: "0.0.0-dev",
	}

	outDir := t.TempDir()
	metricsPath := filepath.Join(outDir, "metrics-compat.json")
	writeCompatMetrics(t, metricsPath, metrics)
	sch := loadCompatSchema(t, root)
	validateCompatMetricsFile(t, sch, metricsPath)

	raw, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["dry_run"] != false {
		t.Fatalf("dry_run must be false for checklist; got %#v", decoded["dry_run"])
	}
	if v, ok := decoded["language_adapter_api_version"].(float64); !ok || int(v) != 1 {
		t.Fatalf("language_adapter_api_version must be 1; got %#v", decoded["language_adapter_api_version"])
	}
	for _, key := range []string{
		"language_adapter_api_version_ok",
		"path_local_bind_ok",
		"trace_lock_ok",
		"migrate_status_ok",
		"backup_restore_ok",
		"no_blob_columns_ok",
		"local_auth_fail_closed_ok",
		"g19_ok",
		"no_daemon_http_primary_ok",
		"no_011_mig_ok",
	} {
		if decoded[key] != true {
			t.Fatalf("%s must be true; got %#v", key, decoded[key])
		}
	}
	t.Logf("checklist green: api_v=%d path_local=%v lock=%v migrate=%v backup=%v no_blob=%v auth=%v g19=%v no_daemon=%v no_011=%v",
		apiVer, pathLocalOK, lockOK, migrateOK, backupOK, noBlobOK, authOK, g19OK, noDaemonOK, no011OK)
}

func checkPathLocalBind(t *testing.T) bool {
	t.Helper()
	a := t.TempDir()
	b := t.TempDir()
	parent := t.TempDir()
	child := filepath.Join(parent, "child")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}

	sa, err := store.Open(a)
	if err != nil {
		t.Fatalf("Open(A): %v", err)
	}
	defer sa.Close()
	sb, err := store.Open(b)
	if err != nil {
		t.Fatalf("Open(B): %v", err)
	}
	defer sb.Close()
	sc, err := store.Open(child)
	if err != nil {
		t.Fatalf("Open(child): %v", err)
	}
	defer sc.Close()

	wantA := filepath.Join(a, ".trace", "trace.db")
	wantB := filepath.Join(b, ".trace", "trace.db")
	wantC := filepath.Join(child, ".trace", "trace.db")
	if sa.DBPath() != wantA || sb.DBPath() != wantB || sc.DBPath() != wantC {
		t.Errorf("DBPath mismatch: A=%q B=%q C=%q", sa.DBPath(), sb.DBPath(), sc.DBPath())
		return false
	}
	if sa.DBPath() == sb.DBPath() {
		t.Error("expected distinct DB paths for A and B")
		return false
	}
	// No walk-up / shared parent: parent must not gain .trace from child Open.
	if _, err := os.Stat(filepath.Join(parent, ".trace")); !os.IsNotExist(err) {
		t.Errorf("parent .trace must not exist after child Open; err=%v", err)
		return false
	}
	return true
}

func checkTraceLock(t *testing.T) bool {
	t.Helper()
	root := t.TempDir()
	s1, err := store.Open(root)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	defer s1.Close()

	var secondErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, secondErr = store.Open(root)
	}()
	wg.Wait()
	if !errors.Is(secondErr, store.ErrLocked) {
		t.Errorf("second Open: want ErrLocked, got %v", secondErr)
		return false
	}
	return true
}

func checkMigrateStatus(t *testing.T, moduleRoot string) (migrateOK, no011OK bool) {
	t.Helper()
	s, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	st, err := s.MigrationStatus()
	if err != nil {
		t.Fatalf("MigrationStatus: %v", err)
	}
	// Phase 26 S03 landed 028_deliberation_consecutive_empty; ceiling is 28 (no 029+).
	migrateOK = st.EmbedExpected == 28 && st.MaxApplied == 28 && st.PendingCount == 0 && len(st.AppliedVersions) == 28
	if !migrateOK {
		t.Errorf("MigrationStatus: embed=%d max=%d pending=%d applied=%v", st.EmbedExpected, st.MaxApplied, st.PendingCount, st.AppliedVersions)
	}
	for _, v := range st.AppliedVersions {
		if v >= 29 {
			t.Errorf("unexpected applied migration version %d (029_* forbidden)", v)
			migrateOK = false
		}
	}

	schemaDir := filepath.Join(moduleRoot, "internal", "store", "schema")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		t.Fatalf("ReadDir schema: %v", err)
	}
	saw011 := false
	saw012 := false
	saw013 := false
	saw014 := false
	saw015 := false
	saw016 := false
	saw017 := false
	saw018 := false
	saw019 := false
	saw020 := false
	saw021 := false
	saw022 := false
	saw023 := false
	saw024 := false
	saw025 := false
	saw026 := false
	saw027 := false
	saw028 := false
	no011OK = true
	for _, e := range entries {
		name := e.Name()
		if name == "011_import_edge_provenance.sql" {
			saw011 = true
			continue
		}
		if name == "012_import_provenance_enum.sql" {
			saw012 = true
			continue
		}
		if name == "013_capability_tool_decisions.sql" {
			saw013 = true
			continue
		}
		if name == "014_capability_tool_decision_enum.sql" {
			saw014 = true
			continue
		}
		if name == "015_deliberation_state.sql" {
			saw015 = true
			continue
		}
		if name == "016_cognitive_artifacts.sql" {
			saw016 = true
			continue
		}
		if name == "017_changes_effects.sql" {
			saw017 = true
			continue
		}
		if name == "018_outcome_results_baselines.sql" {
			saw018 = true
			continue
		}
		if name == "019_regressions_reflections.sql" {
			saw019 = true
			continue
		}
		if name == "020_baselines_promotion.sql" {
			saw020 = true
			continue
		}
		if name == "021_experiments.sql" {
			saw021 = true
			continue
		}
		if name == "022_code_relationships.sql" {
			saw022 = true
			continue
		}
		if name == "023_graph_sync.sql" {
			saw023 = true
			continue
		}
		if name == "024_impact_compare.sql" {
			saw024 = true
			continue
		}
		if name == "025_engineering_knowledge.sql" {
			saw025 = true
			continue
		}
		if name == "026_eval_rules.sql" {
			saw026 = true
			continue
		}
		if name == "027_harness_agents.sql" {
			saw027 = true
			continue
		}
		if name == "028_deliberation_consecutive_empty.sql" {
			saw028 = true
			continue
		}
		if strings.HasPrefix(name, "029_") || strings.Contains(name, "029_") {
			t.Errorf("forbidden migration file present: %s", name)
			no011OK = false
		}
	}
	if !saw011 {
		t.Errorf("expected migration file 011_import_edge_provenance.sql")
		no011OK = false
	}
	if !saw012 {
		t.Errorf("expected migration file 012_import_provenance_enum.sql")
		no011OK = false
	}
	if !saw013 {
		t.Errorf("expected migration file 013_capability_tool_decisions.sql")
		no011OK = false
	}
	if !saw014 {
		t.Errorf("expected migration file 014_capability_tool_decision_enum.sql")
		no011OK = false
	}
	if !saw015 {
		t.Errorf("expected migration file 015_deliberation_state.sql")
		no011OK = false
	}
	if !saw016 {
		t.Errorf("expected migration file 016_cognitive_artifacts.sql")
		no011OK = false
	}
	if !saw017 {
		t.Errorf("expected migration file 017_changes_effects.sql")
		no011OK = false
	}
	if !saw018 {
		t.Errorf("expected migration file 018_outcome_results_baselines.sql")
		no011OK = false
	}
	if !saw019 {
		t.Errorf("expected migration file 019_regressions_reflections.sql")
		no011OK = false
	}
	if !saw020 {
		t.Errorf("expected migration file 020_baselines_promotion.sql")
		no011OK = false
	}
	if !saw021 {
		t.Errorf("expected migration file 021_experiments.sql")
		no011OK = false
	}
	if !saw022 {
		t.Errorf("expected migration file 022_code_relationships.sql")
		no011OK = false
	}
	if !saw023 {
		t.Errorf("expected migration file 023_graph_sync.sql")
		no011OK = false
	}
	if !saw024 {
		t.Errorf("expected migration file 024_impact_compare.sql")
		no011OK = false
	}
	if !saw025 {
		t.Errorf("expected migration file 025_engineering_knowledge.sql")
		no011OK = false
	}
	if !saw026 {
		t.Errorf("expected migration file 026_eval_rules.sql")
		no011OK = false
	}
	if !saw027 {
		t.Errorf("expected migration file 027_harness_agents.sql")
		no011OK = false
	}
	if !saw028 {
		t.Errorf("expected migration file 028_deliberation_consecutive_empty.sql")
		no011OK = false
	}
	if st.EmbedExpected != 28 {
		t.Errorf("EmbedExpected=%d want 28", st.EmbedExpected)
		no011OK = false
	}
	// JSON field no_011_mig_ok retained for schema-compat v1; means ceiling OK (011–028 present, no 029+).
	return migrateOK, no011OK
}

func checkBackupRestore(t *testing.T) (backupOK, noBlobOK bool) {
	t.Helper()
	rootA := t.TempDir()
	sa, err := store.Open(rootA)
	if err != nil {
		t.Fatalf("Open A: %v", err)
	}
	goal, err := sa.UpsertGoal(store.Goal{
		Title:      "compat-checklist-marker",
		Body:       "round-trip",
		SourceType: "USER_ASSERTED",
		Confidence: 1,
		Status:     store.StatusActive,
	})
	if err != nil {
		sa.Close()
		t.Fatalf("UpsertGoal: %v", err)
	}
	if err := store.SetAccessToken(rootA, "secret-must-not-copy"); err != nil {
		sa.Close()
		t.Fatalf("SetAccessToken: %v", err)
	}
	bak := filepath.Join(t.TempDir(), "snap.db")
	if err := sa.BackupTo(bak); err != nil {
		sa.Close()
		t.Fatalf("BackupTo: %v", err)
	}
	if err := sa.Close(); err != nil {
		t.Fatalf("Close A: %v", err)
	}

	rootB := t.TempDir()
	if err := store.Restore(rootB, bak, false); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if _, err := os.Stat(filepath.Join(rootB, ".trace", "access.token")); !os.IsNotExist(err) {
		t.Errorf("access.token must be absent after default restore; err=%v", err)
		return false, false
	}

	sb, err := store.Open(rootB)
	if err != nil {
		t.Fatalf("Open B: %v", err)
	}

	got, err := sb.GetGoal(goal.ID)
	if err != nil || got.Title != "compat-checklist-marker" {
		t.Errorf("GetGoal on B: title=%q err=%v", got.Title, err)
		return false, false
	}
	absB, err := filepath.Abs(rootB)
	if err != nil {
		t.Fatal(err)
	}
	if sb.ProjectRoot() != absB {
		t.Errorf("ProjectRoot: got %q want %q", sb.ProjectRoot(), absB)
		return false, false
	}
	pid := sb.ProjectID()
	dbPath := sb.DBPath()

	bad, where, err := sb.HasBlobLikeColumns()
	if err != nil {
		t.Fatalf("HasBlobLikeColumns: %v", err)
	}
	noBlobOK = !bad
	if bad {
		t.Errorf("unexpected BLOB-like column: %s", where)
	}

	// Lock fail-closed on Backup while Open held.
	bak2 := filepath.Join(t.TempDir(), "locked.db")
	if err := store.Backup(rootB, bak2, store.BackupOptions{}); !errors.Is(err, store.ErrLocked) {
		t.Errorf("Backup while locked: want ErrLocked, got %v", err)
		return false, noBlobOK
	}
	if err := sb.Close(); err != nil {
		t.Fatalf("Close B: %v", err)
	}

	// Abs rebind: projects.root_path must match destination Abs after Restore.
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sqlite open restored db: %v", err)
	}
	defer db.Close()
	var rootPath string
	if err := db.QueryRow(`SELECT root_path FROM projects WHERE id = ?`, pid).Scan(&rootPath); err != nil {
		t.Errorf("root_path query: %v", err)
		return false, noBlobOK
	}
	if rootPath != absB {
		t.Errorf("rebind: root_path=%q want %q", rootPath, absB)
		return false, noBlobOK
	}
	return true, noBlobOK
}

func checkLocalAuthFailClosed(t *testing.T) bool {
	t.Helper()
	root := t.TempDir()
	s, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open cold: %v", err)
	}
	s.Close()

	if err := store.SetAccessToken(root, "correct-token"); err != nil {
		t.Fatalf("SetAccessToken: %v", err)
	}
	t.Setenv(store.AccessTokenEnv, "")
	if _, err := store.Open(root); !errors.Is(err, store.ErrUnauthorized) {
		t.Errorf("Open without env: want ErrUnauthorized, got %v", err)
		return false
	}
	t.Setenv(store.AccessTokenEnv, "wrong-token")
	if _, err := store.Open(root); !errors.Is(err, store.ErrUnauthorized) {
		t.Errorf("Open wrong env: want ErrUnauthorized, got %v", err)
		return false
	}
	t.Setenv(store.AccessTokenEnv, "correct-token")
	s2, err := store.Open(root)
	if err != nil {
		t.Errorf("Open matching token: %v", err)
		return false
	}
	s2.Close()
	return true
}

func checkG19(t *testing.T, moduleRoot string) bool {
	t.Helper()
	// Library packages under internal/ (non-test) must not import cmd/trace or cmd/trace-mcp.
	forbidden := []string{
		"github.com/mrchatam/Trace/cmd/trace",
		"github.com/mrchatam/Trace/cmd/trace-mcp",
	}
	internalRoot := filepath.Join(moduleRoot, "internal")
	ok := true
	err := filepath.WalkDir(internalRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range f.Imports {
			pathLit := strings.Trim(imp.Path.Value, `"`)
			for _, bad := range forbidden {
				if pathLit == bad {
					t.Errorf("G19 violation: %s imports %s", path, pathLit)
					ok = false
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk internal: %v", err)
	}
	return ok
}

func checkNoDaemonHTTPPrimary(t *testing.T, moduleRoot string) bool {
	t.Helper()
	ok := true

	// cmd/trace and cmd/trace-mcp must not import net/http or start listeners.
	for _, cmdPkg := range []string{"cmd/trace", "cmd/trace-mcp"} {
		dir := filepath.Join(moduleRoot, cmdPkg)
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ReadDir %s: %v", dir, err)
		}
		for _, e := range entries {
			name := e.Name()
			if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			path := filepath.Join(dir, name)
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			for _, imp := range f.Imports {
				pathLit := strings.Trim(imp.Path.Value, `"`)
				if pathLit == "net/http" || pathLit == "net/http/httptest" {
					t.Errorf("daemon/HTTP primary: %s imports %s", path, pathLit)
					ok = false
				}
			}
		}
	}

	// MCP must remain the locked six tools — no auth/backup tools.
	serverPath := filepath.Join(moduleRoot, "internal", "mcp", "server.go")
	raw, err := os.ReadFile(serverPath)
	if err != nil {
		t.Fatalf("read mcp server: %v", err)
	}
	body := string(raw)
	forbiddenTools := []string{"trace_backup", "trace_restore", "trace_auth", "trace_migrate"}
	for _, tool := range forbiddenTools {
		if strings.Contains(body, `"`+tool+`"`) || strings.Contains(body, "Name:        \""+tool+"\"") {
			t.Errorf("forbidden MCP tool registered: %s", tool)
			ok = false
		}
	}
	for _, required := range []string{"trace_why", "trace_context", "trace_add", "trace_link", "trace_transition", "trace_review"} {
		if !strings.Contains(body, required) {
			t.Errorf("expected MCP tool missing: %s", required)
			ok = false
		}
	}

	// Spot-check: no ListenAndServe in product cmd/internal (stdio MCP only).
	cmd := exec.Command("rg", "-n", `ListenAndServe|http\.ListenAndServe`, "cmd", "internal")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		t.Errorf("ListenAndServe found (daemon/HTTP primary risk):\n%s", out)
		ok = false
	} else if err != nil {
		// rg exit 1 = no matches (good); other errors: fall back to walk
		if ee, okExit := err.(*exec.ExitError); !okExit || ee.ExitCode() != 1 {
			// Fallback without rg
			_ = filepath.WalkDir(moduleRoot, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				rel, _ := filepath.Rel(moduleRoot, path)
				if !strings.HasPrefix(rel, "cmd"+string(os.PathSeparator)) && !strings.HasPrefix(rel, "internal"+string(os.PathSeparator)) {
					return nil
				}
				if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
					return nil
				}
				b, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if strings.Contains(string(b), "ListenAndServe") {
					t.Errorf("ListenAndServe in %s", rel)
					ok = false
				}
				return nil
			})
		}
	}
	return ok
}

func loadCompatSchema(t *testing.T, root string) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(root, "evals", "compat", "schema-compat.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema-compat.json: %v", err)
	}
	return sch
}

func writeCompatMetrics(t *testing.T, path string, m MetricsCompat) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validateCompatMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read metrics %s: %v", path, err)
	}
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal metrics: %v", err)
	}
	if err := sch.Validate(doc); err != nil {
		t.Fatalf("schema validation failed for %s: %v\n%s", path, err, raw)
	}
}
