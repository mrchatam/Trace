package x0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Stable UUIDs from fixtures/x0/seed/gt.json.
const (
	TaskID = "22222222-2222-2222-2222-222222222222"
)

const fixtureLabel = "fixtures/x0"

var sharedTraceBin string

func TestMain(m *testing.M) {
	root := findModuleRoot()
	if root != "" {
		tmp, err := os.MkdirTemp("", "trace-x0-bin-*")
		if err == nil {
			bin := filepath.Join(tmp, "trace")
			cmd := exec.Command("go", "build", "-o", bin, "./cmd/trace")
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "x0: go build ./cmd/trace failed: %v\n%s\n", err, out)
			} else {
				sharedTraceBin = bin
			}
		}
	}
	if sharedTraceBin == "" {
		fmt.Fprintln(os.Stderr, "x0: skipping build — run with CGO_ENABLED=1")
	}
	os.Exit(m.Run())
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
		t.Fatal("go.mod not found above evals/x0")
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

func copyFixture(t *testing.T, src, dest string) {
	t.Helper()
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() && (d.Name() == ".trace" || d.Name() == ".git") {
			return fs.SkipDir
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
	if err != nil {
		t.Fatalf("copy fixture: %v", err)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

type workEnv struct {
	ModuleRoot string
	Work       string
	SeedAbs    string
	TraceBin   string
	MetricsDir string
}

func setupWork(t *testing.T) workEnv {
	t.Helper()
	root := moduleRoot(t)
	bin := requireTraceBin(t)
	src := filepath.Join(root, "fixtures", "x0")
	work := t.TempDir()
	copyFixture(t, src, work)
	seedAbs, err := filepath.Abs(filepath.Join(work, "seed", "gt.json"))
	if err != nil {
		t.Fatal(err)
	}
	metricsDir := t.TempDir()
	return workEnv{
		ModuleRoot: root,
		Work:       work,
		SeedAbs:    seedAbs,
		TraceBin:   bin,
		MetricsDir: metricsDir,
	}
}

func runTrace(t *testing.T, env workEnv, args ...string) (stdout, stderr string, code int, latencyMS int64) {
	t.Helper()
	cmdArgs := append([]string{"-C", env.Work}, args...)
	cmd := exec.Command(env.TraceBin, cmdArgs...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	start := time.Now()
	err := cmd.Run()
	latencyMS = time.Since(start).Milliseconds()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return outBuf.String(), errBuf.String(), ee.ExitCode(), latencyMS
		}
		t.Fatalf("run %v: %v\nstderr=%s", args, err, errBuf.String())
	}
	return outBuf.String(), errBuf.String(), 0, latencyMS
}

func mustTrace(t *testing.T, env workEnv, args ...string) (stdout string, latencyMS int64) {
	t.Helper()
	out, errOut, code, ms := runTrace(t, env, args...)
	if code != 0 {
		t.Fatalf("trace %v exit %d\nstdout=%s\nstderr=%s", args, code, out, errOut)
	}
	return out, ms
}

func initSeedIndex(t *testing.T, env workEnv) {
	t.Helper()
	mustTrace(t, env, "init")
	out, _ := mustTrace(t, env, "seed", "import", env.SeedAbs)
	var summary map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &summary); err != nil {
		t.Fatalf("seed json: %v\n%s", err, out)
	}
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}
	mustTrace(t, env, "index")
}

func traceVersion(t *testing.T, env workEnv) string {
	t.Helper()
	out, _ := mustTrace(t, env, "version")
	v := strings.TrimSpace(out)
	if v == "" {
		t.Fatal("trace version empty")
	}
	return v
}

// MetricsX0 is the dry-run artifact shape (schema.json v1).
type MetricsX0 struct {
	SchemaVersion int          `json:"schema_version"`
	Experiment    string       `json:"experiment"`
	Condition     string       `json:"condition"`
	Fixture       string       `json:"fixture"`
	DryRun        bool         `json:"dry_run"`
	TraceVersion  string       `json:"trace_version"`
	Runs          []MetricsRun `json:"runs"`
	ToolsUsed     []string     `json:"tools_used"`
	Model         string       `json:"model,omitempty"`
	Seed          string       `json:"seed,omitempty"`
	GitSHA        string       `json:"git_sha,omitempty"`
}

type MetricsRun struct {
	RunID      string          `json:"run_id"`
	TaskFamily string          `json:"task_family"`
	OK         bool            `json:"ok"`
	Efficiency *Efficiency     `json:"efficiency,omitempty"`
	Quality    json.RawMessage `json:"quality,omitempty"`
	Notes      string          `json:"notes,omitempty"`
}

type Efficiency struct {
	LatencyMS *float64 `json:"latency_ms,omitempty"`
	Tokens    *float64 `json:"tokens,omitempty"`
}

func writeMetrics(t *testing.T, path string, m MetricsX0) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func loadSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(moduleRoot(t), "evals", "x0", "schema.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	return sch
}

func validateMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) map[string]any {
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
	var asMap map[string]any
	if err := json.Unmarshal(raw, &asMap); err != nil {
		t.Fatal(err)
	}
	return asMap
}

func toolsContainTraceCLI(tools []string, needle string) bool {
	for _, t := range tools {
		tl := strings.ToLower(t)
		if strings.Contains(tl, needle) {
			return true
		}
	}
	return false
}

func stubB0Agent(t *testing.T, env workEnv) {
	t.Helper()
	// B0: ordinary repo tools only — read fixture files under work; no why/context.
	greeter := filepath.Join(env.Work, "src", "greeter.ts")
	b, err := os.ReadFile(greeter)
	if err != nil {
		t.Fatalf("B0 stub read: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("B0 stub: greeter.ts empty")
	}
}

// TestX0DryRunMetricsB0AndG1 emits schema-valid dry-run metrics for B0 and G1.
// Does not claim Gate C closed or that G1 beats B0.
func TestX0DryRunMetricsB0AndG1(t *testing.T) {
	sch := loadSchema(t)
	metricsDir := t.TempDir()
	pathB0 := filepath.Join(metricsDir, "metrics-b0.json")
	pathG1 := filepath.Join(metricsDir, "metrics-g1.json")

	// --- B0: prep + stub agent (no why/context) ---
	envB0 := setupWork(t)
	initSeedIndex(t, envB0)
	ver := traceVersion(t, envB0)
	stubB0Agent(t, envB0)
	writeMetrics(t, pathB0, MetricsX0{
		SchemaVersion: 1,
		Experiment:    "X0",
		Condition:     "B0",
		Fixture:       fixtureLabel,
		DryRun:        true,
		TraceVersion:  ver,
		Runs: []MetricsRun{{
			RunID:      "b0-dry-1",
			TaskFamily: "understanding",
			OK:         true,
			Notes:      "stub agent; repo file read only",
		}},
		ToolsUsed: []string{"read_file"},
		Seed:      envB0.SeedAbs,
	})

	// --- G1: prep + live why + context ---
	envG1 := setupWork(t)
	initSeedIndex(t, envG1)
	verG1 := traceVersion(t, envG1)

	whyOut, whyMS := mustTrace(t, envG1, "why", "task", TaskID)
	if strings.TrimSpace(whyOut) == "" {
		t.Fatal("G1 why: empty stdout")
	}
	ctxOut, ctxMS := mustTrace(t, envG1, "context", TaskID, "--format", "json")
	if strings.TrimSpace(ctxOut) == "" {
		t.Fatal("G1 context: empty stdout")
	}
	totalMS := float64(whyMS + ctxMS)
	writeMetrics(t, pathG1, MetricsX0{
		SchemaVersion: 1,
		Experiment:    "X0",
		Condition:     "G1",
		Fixture:       fixtureLabel,
		DryRun:        true,
		TraceVersion:  verG1,
		Runs: []MetricsRun{{
			RunID:      "g1-dry-1",
			TaskFamily: "understanding",
			OK:         true,
			Efficiency: &Efficiency{LatencyMS: &totalMS},
			Notes:      "dry-run invoked live trace why + context",
		}},
		ToolsUsed: []string{"trace why", "trace context", "read_file"},
		Seed:      envG1.SeedAbs,
	})

	docB0 := validateMetricsFile(t, sch, pathB0)
	docG1 := validateMetricsFile(t, sch, pathG1)

	if docB0["dry_run"] != true {
		t.Fatalf("B0 dry_run want true got %v", docB0["dry_run"])
	}
	if docG1["dry_run"] != true {
		t.Fatalf("G1 dry_run want true got %v", docG1["dry_run"])
	}
	if docB0["condition"] != "B0" {
		t.Fatalf("B0 condition: %v", docB0["condition"])
	}
	if docG1["condition"] != "G1" {
		t.Fatalf("G1 condition: %v", docG1["condition"])
	}

	toolsB0, _ := toStringSlice(docB0["tools_used"])
	toolsG1, _ := toStringSlice(docG1["tools_used"])
	if toolsContainTraceCLI(toolsB0, "why") || toolsContainTraceCLI(toolsB0, "context") {
		t.Fatalf("B0 must not list why/context tools: %v", toolsB0)
	}
	if !toolsContainTraceCLI(toolsG1, "why") {
		t.Fatalf("G1 tools_used missing why: %v", toolsG1)
	}
	if !toolsContainTraceCLI(toolsG1, "context") {
		t.Fatalf("G1 tools_used missing context: %v", toolsG1)
	}

	// Dry-run understanding task_family required.
	assertUnderstandingRun(t, docB0)
	assertUnderstandingRun(t, docG1)

	if _, err := os.Stat(pathB0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(pathG1); err != nil {
		t.Fatal(err)
	}
}

func toStringSlice(v any) ([]string, bool) {
	arr, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		s, ok := x.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

func assertUnderstandingRun(t *testing.T, doc map[string]any) {
	t.Helper()
	runs, ok := doc["runs"].([]any)
	if !ok || len(runs) < 1 {
		t.Fatalf("want ≥1 run, got %v", doc["runs"])
	}
	r0, ok := runs[0].(map[string]any)
	if !ok {
		t.Fatalf("run[0] not object: %v", runs[0])
	}
	if r0["task_family"] != "understanding" {
		t.Fatalf("task_family want understanding got %v", r0["task_family"])
	}
	if r0["ok"] != true {
		t.Fatalf("run ok want true got %v", r0["ok"])
	}
}
