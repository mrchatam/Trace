package p0x

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
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/store"
)

// Stable UUIDs from fixtures/x0/seed/gt.json (must match README).
const (
	GoalID       = "11111111-1111-1111-1111-111111111111"
	TaskID       = "22222222-2222-2222-2222-222222222222"
	DecisionID   = "33333333-3333-3333-3333-333333333333"
	DiscoveryID  = "44444444-4444-4444-4444-444444444444"
	PlanChangeID = "55555555-5555-5555-5555-555555555555"
)

const (
	TSPath = "src/greeter.ts"
	PyPath = "src/math_util.py"
)

// Paths holds resolved locations for one harness run.
type Paths struct {
	ModuleRoot string
	Work       string
	SeedAbs    string
	TraceBin   string
	Metrics    string
}

// MetricsP0X is the success artifact schema.
type MetricsP0X struct {
	OK       bool             `json:"ok"`
	Criteria map[string]bool  `json:"criteria"`
	Timings  map[string]int64 `json:"timings_ms,omitempty"`
}

var sharedTraceBin string

func TestMain(m *testing.M) {
	root := findModuleRoot()
	code := 1
	if root != "" {
		tmp, err := os.MkdirTemp("", "trace-p0x-bin-*")
		if err == nil {
			bin := filepath.Join(tmp, "trace")
			cmd := exec.Command("go", "build", "-o", bin, "./cmd/trace")
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "p0x: go build ./cmd/trace failed: %v\n%s\n", err, out)
			} else {
				sharedTraceBin = bin
			}
		}
	}
	if sharedTraceBin != "" {
		code = m.Run()
	} else {
		fmt.Fprintln(os.Stderr, "p0x: skipping tests — trace binary build failed (need CGO_ENABLED=1)")
		// Still run so individual tests fail with a clear message.
		code = m.Run()
	}
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
		t.Fatal("go.mod not found above evals/p0x")
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

func setupWork(t *testing.T) Paths {
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
	return Paths{
		ModuleRoot: root,
		Work:       work,
		SeedAbs:    seedAbs,
		TraceBin:   bin,
		Metrics:    filepath.Join(work, "metrics-p0x.json"),
	}
}

func runTrace(t *testing.T, p Paths, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmdArgs := append([]string{"-C", p.Work}, args...)
	cmd := exec.Command(p.TraceBin, cmdArgs...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return outBuf.String(), errBuf.String(), ee.ExitCode()
		}
		t.Fatalf("run %v: %v\nstderr=%s", args, err, errBuf.String())
	}
	return outBuf.String(), errBuf.String(), 0
}

func mustTrace(t *testing.T, p Paths, args ...string) string {
	t.Helper()
	out, errOut, code := runTrace(t, p, args...)
	if code != 0 {
		t.Fatalf("trace %v exit %d\nstdout=%s\nstderr=%s", args, code, out, errOut)
	}
	return out
}

func decodeJSON(t *testing.T, raw string, dest any) {
	t.Helper()
	raw = strings.TrimSpace(raw)
	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		t.Fatalf("json decode: %v\nraw=%s", err, raw)
	}
}

func writeMetrics(t *testing.T, path string, m MetricsP0X) {
	t.Helper()
	if m.Criteria == nil {
		m.Criteria = map[string]bool{}
	}
	for i := 1; i <= 7; i++ {
		k := fmt.Sprintf("%d", i)
		if _, ok := m.Criteria[k]; !ok {
			t.Fatalf("metrics missing criteria[%q]", k)
		}
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func openStore(t *testing.T, work string) *store.Store {
	t.Helper()
	st, err := store.Open(work)
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	return st
}

func closeStore(st **store.Store) {
	if st == nil || *st == nil {
		return
	}
	_ = (*st).Close()
	*st = nil
}

// fingerprintSymbolsImports is a stable, order-independent snapshot for incremental asserts.
func fingerprintSymbolsImports(t *testing.T, st *store.Store, rel string) string {
	t.Helper()
	syms, err := st.ListSymbolsByPath(rel)
	if err != nil {
		t.Fatalf("ListSymbolsByPath %s: %v", rel, err)
	}
	imps, err := st.ListImportsByPath(rel)
	if err != nil {
		t.Fatalf("ListImportsByPath %s: %v", rel, err)
	}
	var lines []string
	for _, s := range syms {
		lines = append(lines, fmt.Sprintf("S:%s:%s:%d:%d", s.Name, s.Kind, s.StartLine, s.EndLine))
	}
	for _, im := range imps {
		sym := ""
		if im.Symbol != nil {
			sym = *im.Symbol
		}
		lines = append(lines, fmt.Sprintf("I:%s:%s", im.ImportedPath, sym))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func initSeedIndex(t *testing.T, p Paths) {
	t.Helper()
	mustTrace(t, p, "init")
	// ABS seed path — -C does not rewrite it.
	out := mustTrace(t, p, "seed", "import", p.SeedAbs)
	var summary map[string]any
	decodeJSON(t, out, &summary)
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}
	mustTrace(t, p, "index")
}

// TestP0XAllCriteria is the primary 7/7 gate (includes ≥5 understanding queries).
func TestP0XAllCriteria(t *testing.T) {
	start := time.Now()
	p := setupWork(t)
	timings := map[string]int64{}

	t0 := time.Now()
	initSeedIndex(t, p)
	timings["init_seed_index"] = time.Since(t0).Milliseconds()

	criteria := map[string]bool{
		"1": false, "2": false, "3": false, "4": false,
		"5": false, "6": false, "7": false,
	}

	st := openStore(t, p.Work)
	defer closeStore(&st)

	// --- Criterion 1: round-trip + ACTIVE provenance + transition ---
	t.Run("criterion-1-roundtrip", func(t *testing.T) {
		g, err := st.GetGoal(GoalID)
		if err != nil {
			t.Fatalf("GetGoal: %v", err)
		}
		if g.Status != store.StatusActive {
			t.Fatalf("goal status: %q", g.Status)
		}
		task, err := st.GetTask(TaskID)
		if err != nil {
			t.Fatalf("GetTask: %v", err)
		}
		if task.Status != store.StatusActive {
			t.Fatalf("task status: %q", task.Status)
		}
		if task.WorkState != store.WorkStateInProgress {
			t.Fatalf("task work_state want IN_PROGRESS got %q", task.WorkState)
		}
		if task.GoalID == nil || *task.GoalID != GoalID {
			t.Fatalf("task goal_id: %v", task.GoalID)
		}
		d, err := st.GetDecision(DecisionID)
		if err != nil || d.Status != store.StatusActive {
			t.Fatalf("GetDecision: %+v err=%v", d, err)
		}
		disc, err := st.GetDiscovery(DiscoveryID)
		if err != nil || disc.Status != store.StatusActive {
			t.Fatalf("GetDiscovery: %+v err=%v", disc, err)
		}
		criteria["1"] = true
	})

	// --- Criterion 2: TS + Py files with symbol and/or import ---
	t.Run("criterion-2-files-symbols", func(t *testing.T) {
		for _, path := range []string{TSPath, PyPath} {
			f, err := st.GetFileByPath(path)
			if err != nil {
				t.Fatalf("GetFileByPath %s: %v", path, err)
			}
			if f.Path != path {
				t.Fatalf("path: %q", f.Path)
			}
			syms, err := st.ListSymbolsByPath(path)
			if err != nil {
				t.Fatal(err)
			}
			imps, err := st.ListImportsByPath(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(syms) == 0 && len(imps) == 0 {
				t.Fatalf("%s: expected ≥1 symbol or import", path)
			}
		}
		criteria["2"] = true
	})

	// Release exclusive lock before CLI why/context (second store.Open).
	closeStore(&st)

	// --- Criterion 3: why causal chain + reason codes ---
	var whyTaskJSON string
	t.Run("criterion-3-why", func(t *testing.T) {
		whyTaskJSON = mustTrace(t, p, "why", "task", TaskID)
		var why map[string]any
		decodeJSON(t, whyTaskJSON, &why)
		steps, ok := why["steps"].([]any)
		if !ok || len(steps) == 0 {
			t.Fatalf("why missing steps: %v", why)
		}
		var sawGoalOrDecision bool
		for _, s := range steps {
			step := s.(map[string]any)
			rc, _ := step["reason_code"].(string)
			if strings.TrimSpace(rc) == "" {
				t.Fatalf("empty reason_code: %v", step)
			}
			id, _ := step["entity_id"].(string)
			rcStr := rc
			if id == GoalID || id == DecisionID ||
				rcStr == "goal_has_task" || rcStr == "decision_affects_task" {
				sawGoalOrDecision = true
			}
		}
		if !sawGoalOrDecision {
			t.Fatalf("why chain missing goal/decision neighbor:\n%s", whyTaskJSON)
		}
		criteria["3"] = true
	})

	// --- Criterion 4: context bounded ---
	var ctxJSON string
	t.Run("criterion-4-context", func(t *testing.T) {
		ctxJSON = mustTrace(t, p, "context", TaskID, "--format", "json")
		var pkt map[string]any
		decodeJSON(t, ctxJSON, &pkt)
		items, ok := pkt["items"].([]any)
		if !ok {
			t.Fatalf("missing items: %v", pkt)
		}
		if len(items) > 32 {
			t.Fatalf("items length %d > 32", len(items))
		}
		budget, ok := pkt["budget"].(map[string]any)
		if !ok {
			t.Fatalf("missing budget: %v", pkt)
		}
		tl, _ := budget["token_limit"].(float64)
		mi, _ := budget["max_items"].(float64)
		if int(tl) != 4096 {
			t.Fatalf("token_limit want 4096 got %v", budget["token_limit"])
		}
		if int(mi) > 32 {
			t.Fatalf("max_items %v > 32", budget["max_items"])
		}
		te, _ := budget["tokens_est"].(float64)
		if te > 4096 {
			t.Fatalf("tokens_est %v exceeds 4096", te)
		}
		foundTrust := false
		for _, it := range items {
			m := it.(map[string]any)
			if m["trust"] != nil && m["trust"] != "" {
				foundTrust = true
			}
			if excerpt, _ := m["excerpt"].(string); excerpt != "" {
				tr, _ := m["trust"].(string)
				if tr != "untrusted_data" && tr != "system" {
					t.Fatalf("excerpt item unexpected trust: %v", m)
				}
			}
		}
		if !foundTrust {
			t.Fatal("expected trust / untrusted_data labeling on context items")
		}
		criteria["4"] = true
	})

	st = openStore(t, p.Work)

	// --- Criterion 5: GT UUIDs + links ---
	t.Run("criterion-5-gt-match", func(t *testing.T) {
		for _, id := range []struct {
			kind string
			get  func() error
		}{
			{"goal", func() error { _, err := st.GetGoal(GoalID); return err }},
			{"task", func() error { _, err := st.GetTask(TaskID); return err }},
			{"decision", func() error { _, err := st.GetDecision(DecisionID); return err }},
			{"discovery", func() error { _, err := st.GetDiscovery(DiscoveryID); return err }},
			{"plan_change", func() error { _, err := st.GetPlanChange(PlanChangeID); return err }},
		} {
			if err := id.get(); err != nil {
				t.Fatalf("missing GT %s: %v", id.kind, err)
			}
		}
		links, err := st.ListLinksFrom("decision", DecisionID)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, l := range links {
			if l.Rel == "decision_affects_task" && l.ToID == TaskID {
				found = true
			}
		}
		if !found {
			t.Fatalf("decision_affects_task link missing: %+v", links)
		}
		dpc, err := st.ListLinksFrom("discovery", DiscoveryID)
		if err != nil {
			t.Fatal(err)
		}
		foundPC := false
		for _, l := range dpc {
			if l.Rel == "discovery_causes_plan_change" && l.ToID == PlanChangeID {
				foundPC = true
			}
		}
		if !foundPC {
			t.Fatalf("discovery_causes_plan_change missing: %+v", dpc)
		}
		if !strings.Contains(whyTaskJSON, GoalID) && !strings.Contains(whyTaskJSON, DecisionID) {
			t.Fatalf("why JSON missing GT ids:\n%s", whyTaskJSON)
		}
		criteria["5"] = true
	})

	// --- Criterion 6: ≥5 understanding queries ---
	queriesOK := 0
	t.Run("criterion-6-queries", func(t *testing.T) {
		t.Run("why-task", func(t *testing.T) {
			if whyTaskJSON == "" {
				closeStore(&st)
				whyTaskJSON = mustTrace(t, p, "why", "task", TaskID)
				st = openStore(t, p.Work)
			}
			if !strings.Contains(whyTaskJSON, "reason_code") {
				t.Fatal("missing reason_code")
			}
			if !strings.Contains(whyTaskJSON, GoalID) && !strings.Contains(whyTaskJSON, DecisionID) {
				t.Fatal("missing goal/decision neighbor")
			}
			queriesOK++
		})
		t.Run("why-decision", func(t *testing.T) {
			closeStore(&st)
			out := mustTrace(t, p, "why", "decision", DecisionID)
			st = openStore(t, p.Work)
			var why map[string]any
			decodeJSON(t, out, &why)
			steps, ok := why["steps"].([]any)
			if !ok || len(steps) == 0 {
				t.Fatalf("empty why decision: %v", why)
			}
			for _, s := range steps {
				step := s.(map[string]any)
				if strings.TrimSpace(fmt.Sprint(step["reason_code"])) == "" || step["reason_code"] == nil {
					t.Fatalf("empty reason_code: %v", step)
				}
			}
			queriesOK++
		})
		t.Run("decision-constraint", func(t *testing.T) {
			blob := whyTaskJSON + "\n" + ctxJSON
			if !strings.Contains(blob, DecisionID) &&
				!strings.Contains(blob, "decision_affects_task") &&
				!strings.Contains(strings.ToLower(blob), "typescript greeter") {
				t.Fatalf("decision constraint not visible in why/context")
			}
			queriesOK++
		})
		t.Run("import-or-symbol-neighbor", func(t *testing.T) {
			if st == nil {
				st = openStore(t, p.Work)
			}
			tsFP := fingerprintSymbolsImports(t, st, TSPath)
			pyFP := fingerprintSymbolsImports(t, st, PyPath)
			if tsFP == "" || pyFP == "" {
				t.Fatalf("empty fingerprint ts=%q py=%q", tsFP, pyFP)
			}
			if !strings.Contains(tsFP, "greet") && !strings.Contains(tsFP, "format") {
				t.Fatalf("TS fingerprint missing expected symbol/import:\n%s", tsFP)
			}
			if !strings.Contains(pyFP, "add") && !strings.Contains(pyFP, "math") {
				t.Fatalf("Py fingerprint missing expected symbol/import:\n%s", pyFP)
			}
			queriesOK++
		})
		t.Run("context-boundedness", func(t *testing.T) {
			if ctxJSON == "" {
				closeStore(&st)
				ctxJSON = mustTrace(t, p, "context", TaskID, "--format", "json")
				st = openStore(t, p.Work)
			}
			var pkt map[string]any
			decodeJSON(t, ctxJSON, &pkt)
			items := pkt["items"].([]any)
			if len(items) > 32 {
				t.Fatalf("items %d", len(items))
			}
			// Not a dump of the whole DB: item count tiny vs seeded entity set + files.
			if len(items) > 100 {
				t.Fatalf("looks like dump-all: %d items", len(items))
			}
			queriesOK++
		})
		if t.Failed() {
			return
		}
		if queriesOK < 5 {
			t.Fatalf("want ≥5 understanding queries, got %d", queriesOK)
		}
		criteria["6"] = true
	})

	// --- Criterion 7: one-file reindex isolation ---
	t.Run("criterion-7-incremental", func(t *testing.T) {
		t0 := time.Now()
		if st == nil {
			st = openStore(t, p.Work)
		}
		pyBefore := fingerprintSymbolsImports(t, st, PyPath)
		fPy, err := st.GetFileByPath(PyPath)
		if err != nil {
			t.Fatal(err)
		}
		hashBefore := fPy.ContentHash

		greeterPath := filepath.Join(p.Work, filepath.FromSlash(TSPath))
		orig, err := os.ReadFile(greeterPath)
		if err != nil {
			t.Fatal(err)
		}
		mutated := string(orig) + "\nexport function greetAgain(name: string): string { return greet(name); }\n"
		if err := os.WriteFile(greeterPath, []byte(mutated), 0o644); err != nil {
			t.Fatal(err)
		}
		// Close store before CLI reindex (same DB file).
		closeStore(&st)

		mustTrace(t, p, "index", TSPath)

		st2, err := store.Open(p.Work)
		if err != nil {
			t.Fatal(err)
		}
		defer st2.Close()

		pyAfter := fingerprintSymbolsImports(t, st2, PyPath)
		if pyAfter != pyBefore {
			t.Fatalf("Py symbols/imports changed after TS-only reindex\nbefore:\n%s\nafter:\n%s", pyBefore, pyAfter)
		}
		fPy2, err := st2.GetFileByPath(PyPath)
		if err != nil {
			t.Fatal(err)
		}
		if fPy2.ContentHash != hashBefore {
			t.Fatalf("Py content_hash changed: %s → %s", hashBefore, fPy2.ContentHash)
		}
		tsAfter := fingerprintSymbolsImports(t, st2, TSPath)
		if !strings.Contains(tsAfter, "greetAgain") {
			t.Fatalf("TS symbols did not update after mutation:\n%s", tsAfter)
		}
		timings["incremental"] = time.Since(t0).Milliseconds()
		criteria["7"] = true
	})

	for i := 1; i <= 7; i++ {
		k := fmt.Sprintf("%d", i)
		if !criteria[k] {
			t.Fatalf("criterion %s not proven", k)
		}
	}

	timings["total"] = time.Since(start).Milliseconds()
	writeMetrics(t, p.Metrics, MetricsP0X{
		OK:       true,
		Criteria: criteria,
		Timings:  timings,
	})
	raw, err := os.ReadFile(p.Metrics)
	if err != nil {
		t.Fatal(err)
	}
	var check MetricsP0X
	if err := json.Unmarshal(raw, &check); err != nil {
		t.Fatal(err)
	}
	if !check.OK || len(check.Criteria) != 7 {
		t.Fatalf("metrics schema: %+v", check)
	}
	for i := 1; i <= 7; i++ {
		if !check.Criteria[fmt.Sprintf("%d", i)] {
			t.Fatalf("metrics criteria[%d] false", i)
		}
	}
}

func TestSeedPathIsAbsoluteNotUnderC(t *testing.T) {
	p := setupWork(t)
	mustTrace(t, p, "init")
	// Relative path would fail if cwd ≠ work; we require abs.
	if !filepath.IsAbs(p.SeedAbs) {
		t.Fatalf("SeedAbs not absolute: %s", p.SeedAbs)
	}
	out := mustTrace(t, p, "seed", "import", p.SeedAbs)
	if !strings.Contains(out, `"ok": true`) && !strings.Contains(out, `"ok":true`) {
		t.Fatalf("seed import failed: %s", out)
	}
}
