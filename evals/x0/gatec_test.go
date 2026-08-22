package x0

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	fixtureContentHash = "15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22"
	gateCModelPin      = "recorded-operator-sim/v1"
)

// TestX0GateCRecordedMetrics grades committed Mode-B answer packs (N≥3/condition),
// emits schema-valid dry_run:false metrics, and refreshes docs/verification/gate-c-x0/.
// Does not require a network LLM. Phase 01 dry-run is a separate test and is not a Gate C pass.
func TestX0GateCRecordedMetrics(t *testing.T) {
	root := moduleRoot(t)
	sch := loadSchema(t)
	bankPath := filepath.Join(root, "evals", "x0", "queries.json")
	bank, err := LoadQueryBank(bankPath)
	if err != nil {
		t.Fatal(err)
	}

	env := setupWork(t)
	initSeedIndex(t, env)
	ver := traceVersion(t, env)

	type condResult struct {
		condition string
		packs     []*AnswerPack
		qualities []UnderstandingQuality
		tools     []string
		path      string
	}

	runCond := func(cond string, files []string) condResult {
		t.Helper()
		dir := filepath.Join(root, "evals", "x0", "testdata", "gate-c")
		var packs []*AnswerPack
		var quals []UnderstandingQuality
		var tools []string
		runs := make([]MetricsRun, 0, len(files))
		for _, name := range files {
			pack, err := LoadAnswerPack(filepath.Join(dir, name))
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if !strings.EqualFold(pack.Condition, cond) {
				t.Fatalf("%s: condition want %s got %s", name, cond, pack.Condition)
			}
			if pack.Model != gateCModelPin {
				t.Fatalf("%s: model pin want %s got %s", name, gateCModelPin, pack.Model)
			}
			q := GradePack(bank, pack)
			rawQ, err := QualityJSON(q)
			if err != nil {
				t.Fatal(err)
			}
			eff := &Efficiency{}
			if pack.LatencyMS != nil {
				eff.LatencyMS = pack.LatencyMS
			}
			if pack.Tokens != nil {
				eff.Tokens = pack.Tokens
			}
			if pack.LatencyMS == nil && pack.Tokens == nil {
				eff = nil
			}
			runs = append(runs, MetricsRun{
				RunID:      pack.RunID,
				TaskFamily: "understanding",
				OK:         true, // completed graded understanding attempt
				Efficiency: eff,
				Quality:    rawQ,
				Notes:      pack.Notes,
			})
			packs = append(packs, pack)
			quals = append(quals, q)
			if tools == nil {
				tools = append([]string{}, pack.ToolsUsed...)
			}
		}
		if len(runs) < 3 {
			t.Fatalf("%s: want ≥3 runs got %d", cond, len(runs))
		}

		m := MetricsX0{
			SchemaVersion: 1,
			Experiment:    "X0",
			Condition:     cond,
			Fixture:       fixtureLabel,
			DryRun:        false,
			TraceVersion:  ver,
			Runs:          runs,
			ToolsUsed:     tools,
			Model:         gateCModelPin,
			// Stable repo-relative seed pin for committed verification artifacts
			// (harness still imports via env.SeedAbs in the temp worktree).
			Seed:   "fixtures/x0/seed/gt.json",
			GitSHA: fixtureContentHash,
		}
		tmp := filepath.Join(env.MetricsDir, fmt.Sprintf("metrics-%s.json", strings.ToLower(cond)))
		writeMetrics(t, tmp, m)
		doc := validateMetricsFile(t, sch, tmp)
		if doc["dry_run"] != false {
			t.Fatalf("%s dry_run want false got %v", cond, doc["dry_run"])
		}
		runsArr, _ := doc["runs"].([]any)
		if len(runsArr) < 3 {
			t.Fatalf("%s runs want ≥3 got %d", cond, len(runsArr))
		}
		return condResult{condition: cond, packs: packs, qualities: quals, tools: tools, path: tmp}
	}

	b0 := runCond("B0", []string{"b0-run1.json", "b0-run2.json", "b0-run3.json"})
	g1 := runCond("G1", []string{"g1-run1.json", "g1-run2.json", "g1-run3.json"})

	if toolsContainTraceCLI(b0.tools, "why") || toolsContainTraceCLI(b0.tools, "context") {
		t.Fatalf("B0 must not list why/context: %v", b0.tools)
	}
	if !toolsContainTraceCLI(g1.tools, "why") || !toolsContainTraceCLI(g1.tools, "context") {
		t.Fatalf("G1 tools_used must include why and context: %v", g1.tools)
	}

	meanB0 := MeanAccuracy(b0.qualities)
	meanG1 := MeanAccuracy(g1.qualities)
	cmB0 := MeanCriticalMisses(b0.qualities)
	cmG1 := MeanCriticalMisses(g1.qualities)
	t.Logf("Gate C means: B0 accuracy=%.3f cm=%.3f; G1 accuracy=%.3f cm=%.3f", meanB0, cmB0, meanG1, cmG1)

	// Persist committed verification artifacts (refresh from graded packs).
	outDir := filepath.Join(root, "docs", "verification", "gate-c-x0")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	copyFileMust(t, b0.path, filepath.Join(outDir, "metrics-b0.json"))
	copyFileMust(t, g1.path, filepath.Join(outDir, "metrics-g1.json"))
	validateMetricsFile(t, sch, filepath.Join(outDir, "metrics-b0.json"))
	validateMetricsFile(t, sch, filepath.Join(outDir, "metrics-g1.json"))

	// Kill criteria check is reported in GATE-C-NOTES; test asserts honesty of comparison inputs.
	if meanG1 <= meanB0 {
		t.Logf("KILL: mean G1 (%.3f) ≤ mean B0 (%.3f) — thesis endangered; Notes must be No-Go or iterate", meanG1, meanB0)
	}
}

func copyFileMust(t *testing.T, src, dst string) {
	t.Helper()
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		t.Fatal(err)
	}
}
