package perf_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Locked regression ceilings — formula from P07-S03-01:
//
//	ceiling_ms = max(measured_ms * 5, 2000)
//	ceiling_db = measured_db_bytes * 3
//
// Re-locked 2026-08-22 (Phase 43): max observed across dev linux/amd64 and
// GHA ubuntu-latest (run 32570568221) after schema/index growth through P42.
//
//	smoke:    initial=114ms incr=125ms  db=962560
//	rung-1k:  initial=496ms incr=436ms  db=2375680
//	rung-10k: initial=6746ms incr=7913ms db=16560128
var (
	ceilingSmokeInitialMS       int64 = 2000     // max(114*5, 2000)
	ceilingSmokeIncrementalMS   int64 = 2000     // max(125*5, 2000)
	ceilingSmokeDBBytes         int64 = 2887680  // 962560*3
	ceilingRung1kInitialMS      int64 = 2480     // max(496*5, 2000)
	ceilingRung1kIncrementalMS  int64 = 2180     // max(436*5, 2000)
	ceilingRung1kDBBytes        int64 = 7127040  // 2375680*3
	ceilingRung10kInitialMS     int64 = 33730    // max(6746*5, 2000)
	ceilingRung10kIncrementalMS int64 = 39565    // max(7913*5, 2000)
	ceilingRung10kDBBytes       int64 = 49680384 // 16560128*3
)

var sharedTraceBin string

func TestMain(m *testing.M) {
	root := findModuleRoot()
	code := 1
	if root != "" {
		tmp, err := os.MkdirTemp("", "trace-perf-bin-*")
		if err == nil {
			bin := filepath.Join(tmp, "trace")
			cmd := exec.Command("go", "build", "-o", bin, "./cmd/trace")
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "perf: go build ./cmd/trace failed: %v\n%s\n", err, out)
			} else {
				sharedTraceBin = bin
			}
		}
	}
	if sharedTraceBin == "" {
		fmt.Fprintln(os.Stderr, "perf: skipping tests — trace binary build failed (need CGO_ENABLED=1)")
	}
	code = m.Run()
	os.Exit(code)
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
		t.Fatal("go.mod not found above evals/perf")
	}
	return root
}

func requireTraceBin(t *testing.T) string {
	t.Helper()
	if sharedTraceBin == "" {
		t.Fatal("trace binary not built; run with CGO_ENABLED=1")
	}
	return sharedTraceBin
}

// MetricsGateH is the Gate H artifact shape (schema-gate-h.json v1).
type MetricsGateH struct {
	SchemaVersion          int            `json:"schema_version"`
	Gate                   string         `json:"gate"`
	Suite                  string         `json:"suite"`
	DryRun                 bool           `json:"dry_run"`
	NamedTest              string         `json:"named_test"`
	Rungs                  []RungMetrics  `json:"rungs"`
	Thresholds             map[string]any `json:"thresholds"`
	T0SkipOK               bool           `json:"t0_skip_ok"`
	IncrementalIsolationOK bool           `json:"incremental_isolation_ok"`
	GoAdapterExercised     bool           `json:"go_adapter_exercised"`
	S01Hooks               []string       `json:"s01_hooks"`
	S02Hooks               []string       `json:"s02_hooks"`
	TraceVersion           string         `json:"trace_version,omitempty"`
}

type RungMetrics struct {
	ID                 string `json:"id"`
	ApproxLOC          int    `json:"approx_loc"`
	FileCount          int    `json:"file_count"`
	InitialIndexMS     int64  `json:"initial_index_ms"`
	IncrementalIndexMS int64  `json:"incremental_index_ms"`
	DBBytes            int64  `json:"db_bytes"`
}

type rungPlant struct {
	id           string
	targetLOC    int
	jsFiles      int
	linesPerFile int
	withGo       bool
}

// TestPlantedPerfLadderGateH is the Gate H planted performance ladder.
func TestPlantedPerfLadderGateH(t *testing.T) {
	bin := requireTraceBin(t)
	root := moduleRoot(t)

	rungs := []rungPlant{
		{id: "smoke", targetLOC: 120, jsFiles: 3, linesPerFile: 40, withGo: true},
		{id: "rung-1k", targetLOC: 1000, jsFiles: 10, linesPerFile: 100, withGo: false},
		{id: "rung-10k", targetLOC: 10000, jsFiles: 50, linesPerFile: 200, withGo: false},
	}

	var rungMetrics []RungMetrics
	t0OK := true
	isoOK := true
	goExercised := false

	for _, rp := range rungs {
		rm, gotT0, gotIso, gotGo := runRung(t, bin, rp)
		rungMetrics = append(rungMetrics, rm)
		t0OK = t0OK && gotT0
		isoOK = isoOK && gotIso
		goExercised = goExercised || gotGo
		t.Logf("rung %s: files=%d approx_loc=%d initial_ms=%d incr_ms=%d db_bytes=%d t0=%v iso=%v go=%v",
			rm.ID, rm.FileCount, rm.ApproxLOC, rm.InitialIndexMS, rm.IncrementalIndexMS, rm.DBBytes, gotT0, gotIso, gotGo)
	}

	if !t0OK {
		t.Fatal("t0_skip_ok failed: T0 plant path was indexed or control path missing")
	}
	if !isoOK {
		t.Fatal("incremental_isolation_ok failed: sibling symbols drifted after single-file reindex")
	}
	if !goExercised {
		t.Fatal("go_adapter_exercised expected true (smoke plants a .go file)")
	}

	// Assert locked ceilings (MaxInt64 until first measure locks them).
	assertCeilings(t, rungMetrics)

	thresholds := map[string]any{
		"formula_ms": "max(measured_ms * 5, 2000)",
		"formula_db": "measured_db_bytes * 3",
		"smoke": map[string]int64{
			"initial_index_ms":     ceilingSmokeInitialMS,
			"incremental_index_ms": ceilingSmokeIncrementalMS,
			"db_bytes":             ceilingSmokeDBBytes,
		},
		"rung-1k": map[string]int64{
			"initial_index_ms":     ceilingRung1kInitialMS,
			"incremental_index_ms": ceilingRung1kIncrementalMS,
			"db_bytes":             ceilingRung1kDBBytes,
		},
		"rung-10k": map[string]int64{
			"initial_index_ms":     ceilingRung10kInitialMS,
			"incremental_index_ms": ceilingRung10kIncrementalMS,
			"db_bytes":             ceilingRung10kDBBytes,
		},
	}

	metrics := MetricsGateH{
		SchemaVersion:          1,
		Gate:                   "H",
		Suite:                  "perf",
		DryRun:                 false,
		NamedTest:              "TestPlantedPerfLadderGateH",
		Rungs:                  rungMetrics,
		Thresholds:             thresholds,
		T0SkipOK:               t0OK,
		IncrementalIsolationOK: isoOK,
		GoAdapterExercised:     goExercised,
		S01Hooks: []string{
			"isT0SkipDir",
			"isT0SkipPath",
			"walkIndexable T0→lang→T0 file→gitignore",
			"file-local IndexFile isolation",
		},
		S02Hooks: []string{
			"DetectLanguage .go → LangGo",
			"extract_go.go / tree-sitter-go v0.25.0",
		},
		TraceVersion: "0.0.0-dev",
	}

	outDir := t.TempDir()
	metricsPath := filepath.Join(outDir, "metrics-gate-h.json")
	writeGateHMetrics(t, metricsPath, metrics)
	sch := loadGateHSchema(t, root)
	validateGateHMetricsFile(t, sch, metricsPath)

	raw, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["dry_run"] != false {
		t.Fatalf("dry_run must be false for Gate H; got %#v", decoded["dry_run"])
	}
}

func assertCeilings(t *testing.T, rungs []RungMetrics) {
	t.Helper()
	byID := map[string]RungMetrics{}
	for _, r := range rungs {
		byID[r.ID] = r
	}
	checks := []struct {
		id       string
		initCeil int64
		incrCeil int64
		dbCeil   int64
	}{
		{"smoke", ceilingSmokeInitialMS, ceilingSmokeIncrementalMS, ceilingSmokeDBBytes},
		{"rung-1k", ceilingRung1kInitialMS, ceilingRung1kIncrementalMS, ceilingRung1kDBBytes},
		{"rung-10k", ceilingRung10kInitialMS, ceilingRung10kIncrementalMS, ceilingRung10kDBBytes},
	}
	for _, c := range checks {
		r, ok := byID[c.id]
		if !ok {
			t.Fatalf("missing rung %s", c.id)
		}
		if r.InitialIndexMS > c.initCeil {
			t.Fatalf("%s initial_index_ms=%d exceeds ceiling %d", c.id, r.InitialIndexMS, c.initCeil)
		}
		if r.IncrementalIndexMS > c.incrCeil {
			t.Fatalf("%s incremental_index_ms=%d exceeds ceiling %d", c.id, r.IncrementalIndexMS, c.incrCeil)
		}
		if r.DBBytes > c.dbCeil {
			t.Fatalf("%s db_bytes=%d exceeds ceiling %d", c.id, r.DBBytes, c.dbCeil)
		}
	}
}

func runRung(t *testing.T, bin string, rp rungPlant) (RungMetrics, bool, bool, bool) {
	t.Helper()
	work := t.TempDir()
	approxLOC, indexable, goPath := plantRung(t, work, rp)

	runTrace(t, bin, work, "init")

	start := time.Now()
	runTrace(t, bin, work, "index")
	initialMS := time.Since(start).Milliseconds()

	st, err := store.Open(work)
	if err != nil {
		t.Fatalf("store.Open after initial index: %v", err)
	}
	t0OK := checkT0Skip(t, st)
	goOK := false
	if goPath != "" {
		if _, err := st.GetFileByPath(goPath); err != nil {
			t.Fatalf("expected Go file indexed at %s: %v", goPath, err)
		}
		syms, err := st.ListSymbolsByPath(goPath)
		if err != nil {
			t.Fatalf("ListSymbolsByPath %s: %v", goPath, err)
		}
		if len(syms) == 0 {
			t.Fatalf("expected symbols on planted Go file %s", goPath)
		}
		goOK = true
	}
	// Count indexed planted files (exclude T0 decoys).
	fileCount := 0
	for _, rel := range indexable {
		if _, err := st.GetFileByPath(rel); err != nil {
			t.Fatalf("expected indexed file %s: %v", rel, err)
		}
		fileCount++
	}
	symsB1, err := st.ListSymbolsByPath("iso/b.js")
	if err != nil || len(symsB1) == 0 {
		st.Close()
		t.Fatalf("expected symbols on iso/b.js: err=%v n=%d", err, len(symsB1))
	}
	st.Close()

	// Re-index unchanged tree = incremental_index_ms (before isolation mutation).
	start = time.Now()
	runTrace(t, bin, work, "index")
	incrMS := time.Since(start).Milliseconds()

	dbPath := filepath.Join(work, ".trace", "trace.db")
	fi, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("stat db: %v", err)
	}
	dbBytes := fi.Size()

	isoOK := checkIsolationAfter(t, bin, work, symsB1)

	return RungMetrics{
		ID:                 rp.id,
		ApproxLOC:          approxLOC,
		FileCount:          fileCount,
		InitialIndexMS:     initialMS,
		IncrementalIndexMS: incrMS,
		DBBytes:            dbBytes,
	}, t0OK, isoOK, goOK
}

func plantRung(t *testing.T, work string, rp rungPlant) (approxLOC int, indexable []string, goPath string) {
	t.Helper()
	approxLOC = 0
	for i := 0; i < rp.jsFiles; i++ {
		rel := filepath.ToSlash(filepath.Join("src", fmt.Sprintf("mod_%03d.js", i)))
		body := syntheticJS(i, rp.linesPerFile)
		writePlant(t, work, rel, body)
		approxLOC += strings.Count(body, "\n")
		indexable = append(indexable, rel)
	}
	// Isolation pair: sibling files a.js / b.js
	aBody := "export function alpha() { return 1 }\n"
	bBody := "export function beta() { return 2 }\n"
	writePlant(t, work, "iso/a.js", aBody)
	writePlant(t, work, "iso/b.js", bBody)
	approxLOC += 2
	indexable = append(indexable, "iso/a.js", "iso/b.js")

	if rp.withGo {
		goBody := "package plant\n\nfunc Helper() int { return 1 }\n\ntype Counter struct{}\n\nfunc (c Counter) Run() {}\n"
		writePlant(t, work, "pkg/plant.go", goBody)
		approxLOC += strings.Count(goBody, "\n")
		indexable = append(indexable, "pkg/plant.go")
		goPath = "pkg/plant.go"
	}

	// T0 decoys — must NOT be indexed.
	writePlant(t, work, "node_modules/pkg/hidden.js", "export function hidden() { return 0 }\n")
	writePlant(t, work, "vendor/lib/vendored.js", "export function vendored() { return 0 }\n")
	writePlant(t, work, "foo.min.js", "export function minned() { return 0 }\n")

	if approxLOC < rp.targetLOC/2 {
		t.Fatalf("plant %s approx_loc=%d too low vs target ~%d", rp.id, approxLOC, rp.targetLOC)
	}
	return approxLOC, indexable, goPath
}

func syntheticJS(seed, lines int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("// synthetic plant seed=%d\n", seed))
	for i := 0; i < lines-1; i++ {
		b.WriteString(fmt.Sprintf("export function f_%d_%d() { return %d }\n", seed, i, i))
	}
	return b.String()
}

func writePlant(t *testing.T, root, rel, body string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runTrace(t *testing.T, bin, work string, args ...string) {
	t.Helper()
	cmdArgs := append([]string{"-C", work}, args...)
	cmd := exec.Command(bin, cmdArgs...)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("trace %v: %v\n%s", args, err, out)
	}
}

func checkT0Skip(t *testing.T, st *store.Store) bool {
	t.Helper()
	forbidden := []string{
		"node_modules/pkg/hidden.js",
		"vendor/lib/vendored.js",
		"foo.min.js",
	}
	for _, p := range forbidden {
		if _, err := st.GetFileByPath(p); err == nil {
			t.Errorf("T0 path unexpectedly indexed: %s", p)
			return false
		}
	}
	return true
}

func checkIsolationAfter(t *testing.T, bin, work string, symsB1 []store.Symbol) bool {
	t.Helper()
	aPath := filepath.Join(work, "iso", "a.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 99 }\nexport function alpha2() { return 0 }\n"), 0o644); err != nil {
		t.Error(err)
		return false
	}
	runTrace(t, bin, work, "index", "iso/a.js")

	st2, err := store.Open(work)
	if err != nil {
		t.Errorf("reopen store: %v", err)
		return false
	}
	defer st2.Close()
	symsB2, err := st2.ListSymbolsByPath("iso/b.js")
	if err != nil {
		t.Errorf("ListSymbolsByPath after reindex: %v", err)
		return false
	}
	if len(symsB2) != len(symsB1) {
		t.Errorf("iso/b.js symbols changed after indexing only iso/a.js: before=%d after=%d", len(symsB1), len(symsB2))
		return false
	}
	for i := range symsB1 {
		if symsB1[i].Name != symsB2[i].Name || symsB1[i].Kind != symsB2[i].Kind {
			t.Errorf("iso/b.js symbol drift: %+v vs %+v", symsB1[i], symsB2[i])
			return false
		}
	}
	symsA, err := st2.ListSymbolsByPath("iso/a.js")
	if err != nil || len(symsA) < 2 {
		t.Errorf("expected updated iso/a.js symbols, got %d err=%v", len(symsA), err)
		return false
	}
	return true
}

func loadGateHSchema(t *testing.T, root string) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(root, "evals", "perf", "schema-gate-h.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema-gate-h.json: %v", err)
	}
	return sch
}

func writeGateHMetrics(t *testing.T, path string, m MetricsGateH) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validateGateHMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) {
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
