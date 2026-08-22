package impact_test

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func openImpact(t *testing.T) *domain.Service {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st)
}

// plantLinkedDecision creates an isolated decision ↔ task pair via decision_affects_task.
func plantLinkedDecision(t *testing.T, svc *domain.Service, ctx context.Context, title string) string {
	t.Helper()
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: title})
	if err != nil {
		t.Fatalf("CreateDecision %q: %v", title, err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: title + " task"})
	if err != nil {
		t.Fatalf("CreateTask %q: %v", title, err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDecisionTask %q: %v", title, err)
	}
	return dec.ID
}

// MetricsGateF is the Gate F prelim artifact shape (schema-gate-f.json v1).
type MetricsGateF struct {
	SchemaVersion  int      `json:"schema_version"`
	Gate           string   `json:"gate"`
	Suite          string   `json:"suite"`
	Prelim         bool     `json:"prelim"`
	DryRun         bool     `json:"dry_run"`
	TruePositives  int      `json:"true_positives"`
	FalsePositives int      `json:"false_positives"`
	FalseNegatives int      `json:"false_negatives"`
	TrueNegatives  int      `json:"true_negatives"`
	Precision      float64  `json:"precision"`
	Recall         float64  `json:"recall"`
	Probes         int      `json:"probes"`
	NamedTest      string   `json:"named_test"`
	Mig            string   `json:"mig"`
	S01Hooks       []string `json:"s01_hooks"`
	ProbeIDs       []string `json:"probe_ids,omitempty"`
	TraceVersion   string   `json:"trace_version,omitempty"`
}

// TestPlantedImpactConflictsGateFPrelim is the Gate F preliminary precision/recall harness.
// Plants Pos-1..3 conflict probes + Neg-1 clean SAFE; scores ImpactReport GT assertions.
func TestPlantedImpactConflictsGateFPrelim(t *testing.T) {
	svc := openImpact(t)
	ctx := context.Background()

	var tp, fn, fp, tn int
	probeIDs := []string{"pos-1-unknown", "pos-2-rollup", "pos-3-empty-findings", "neg-1-clean-safe"}

	// Part 1 — Pos-1: UNKNOWN uncertainty conflict
	{
		decID := plantLinkedDecision(t, svc, ctx, "Gate F Pos-1 UNKNOWN")
		if _, err := svc.AddImpactFinding(ctx, decID, domain.ImpactFindingInput{
			ImpactClass: domain.ImpactClassCaution,
			Uncertainty: domain.UncertaintyUNKNOWN,
			Kind:        domain.FindingKindUnresolved,
			Body:        "unresolved impact class with UNKNOWN uncertainty",
		}); err != nil {
			t.Fatalf("Pos-1 AddImpactFinding: %v", err)
		}
		rep, err := svc.ImpactReport(ctx, decID)
		if err != nil {
			t.Fatalf("Pos-1 ImpactReport: %v", err)
		}
		ok := rep.HasUnknown && rep.Incomplete
		hasUnkFinding := false
		for _, f := range rep.Findings {
			if f.Uncertainty == domain.UncertaintyUNKNOWN {
				hasUnkFinding = true
				break
			}
		}
		ok = ok && hasUnkFinding
		// Do NOT assert OverallClass == UNKNOWN (rollup is severity band, not uncertainty).
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-1 FN: want HasUnknown+Incomplete+UNKNOWN finding; got OverallClass=%q HasUnknown=%v Incomplete=%v findings=%d",
				rep.OverallClass, rep.HasUnknown, rep.Incomplete, len(rep.Findings))
		}
	}

	// Part 2 — Pos-2: SAFE + DESTRUCTIVE rollup
	{
		decID := plantLinkedDecision(t, svc, ctx, "Gate F Pos-2 rollup")
		if _, err := svc.AddImpactFinding(ctx, decID, domain.ImpactFindingInput{
			ImpactClass: domain.ImpactClassSAFE,
			Uncertainty: domain.UncertaintyKNOWN,
			Kind:        domain.FindingKindAffectedWork,
			Body:        "known SAFE impact",
		}); err != nil {
			t.Fatalf("Pos-2 SAFE finding: %v", err)
		}
		if _, err := svc.AddImpactFinding(ctx, decID, domain.ImpactFindingInput{
			ImpactClass: domain.ImpactClassDestructive,
			Uncertainty: domain.UncertaintyKNOWN,
			Kind:        domain.FindingKindAffectedWork,
			Body:        "known DESTRUCTIVE impact",
		}); err != nil {
			t.Fatalf("Pos-2 DESTRUCTIVE finding: %v", err)
		}
		rep, err := svc.ImpactReport(ctx, decID)
		if err != nil {
			t.Fatalf("Pos-2 ImpactReport: %v", err)
		}
		ok := rep.OverallClass == domain.ImpactClassDestructive && !rep.HasUnknown
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-2 FN: want OverallClass=DESTRUCTIVE HasUnknown=false; got OverallClass=%q HasUnknown=%v Incomplete=%v",
				rep.OverallClass, rep.HasUnknown, rep.Incomplete)
		}
	}

	// Part 3 — Pos-3: linked tasks, empty findings (incomplete conflict)
	{
		decID := plantLinkedDecision(t, svc, ctx, "Gate F Pos-3 empty findings")
		rep, err := svc.ImpactReport(ctx, decID)
		if err != nil {
			t.Fatalf("Pos-3 ImpactReport: %v", err)
		}
		ok := rep.HasUnknown && rep.Incomplete && rep.OverallClass == ""
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-3 FN: want HasUnknown+Incomplete+OverallClass=\"\"; got OverallClass=%q HasUnknown=%v Incomplete=%v findings=%d",
				rep.OverallClass, rep.HasUnknown, rep.Incomplete, len(rep.Findings))
		}
	}

	// Part 4 — Neg-1: clean known SAFE (no false alarm)
	{
		decID := plantLinkedDecision(t, svc, ctx, "Gate F Neg-1 clean SAFE")
		if _, err := svc.AddImpactFinding(ctx, decID, domain.ImpactFindingInput{
			ImpactClass: domain.ImpactClassSAFE,
			Uncertainty: domain.UncertaintyKNOWN,
			Kind:        domain.FindingKindAffectedWork,
			Body:        "clean known SAFE affected work",
		}); err != nil {
			t.Fatalf("Neg-1 AddImpactFinding: %v", err)
		}
		rep, err := svc.ImpactReport(ctx, decID)
		if err != nil {
			t.Fatalf("Neg-1 ImpactReport: %v", err)
		}
		ok := !rep.HasUnknown && !rep.Incomplete && rep.OverallClass == domain.ImpactClassSAFE
		if ok {
			tn++
		} else {
			fp++
			t.Errorf("Neg-1 FP: want HasUnknown=false Incomplete=false OverallClass=SAFE; got OverallClass=%q HasUnknown=%v Incomplete=%v",
				rep.OverallClass, rep.HasUnknown, rep.Incomplete)
		}
	}

	probes := 4
	if tp != 3 || fn != 0 || fp != 0 || tn != 1 {
		t.Fatalf("tallies want TP=3 FN=0 FP=0 TN=1; got TP=%d FN=%d FP=%d TN=%d", tp, fn, fp, tn)
	}
	precision := float64(tp) / float64(tp+fp)
	recall := float64(tp) / float64(tp+fn)
	if math.Abs(precision-1.0) > 1e-9 || math.Abs(recall-1.0) > 1e-9 {
		t.Fatalf("precision/recall want 1.0/1.0; got precision=%v recall=%v", precision, recall)
	}

	// Part 5 — Metrics write + schema validate
	metricsDir := t.TempDir()
	metricsPath := filepath.Join(metricsDir, "metrics-gate-f.json")
	m := MetricsGateF{
		SchemaVersion:  1,
		Gate:           "F",
		Suite:          "impact",
		Prelim:         true,
		DryRun:         false,
		TruePositives:  tp,
		FalsePositives: fp,
		FalseNegatives: fn,
		TrueNegatives:  tn,
		Precision:      precision,
		Recall:         recall,
		Probes:         probes,
		NamedTest:      "TestPlantedImpactConflictsGateFPrelim",
		Mig:            "009_decision_impact",
		S01Hooks: []string{
			"AddImpactFinding",
			"LinkDecisionTask",
			"ImpactReport",
			"decision_affects_task",
		},
		ProbeIDs:     probeIDs,
		TraceVersion: "0.0.0-dev",
	}
	writeGateFMetrics(t, metricsPath, m)
	validateGateFMetricsFile(t, loadGateFSchema(t), metricsPath)
}

func findImpactModuleRoot() string {
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

func impactModuleRoot(t *testing.T) string {
	t.Helper()
	root := findImpactModuleRoot()
	if root == "" {
		t.Fatal("go.mod not found above evals/impact")
	}
	return root
}

func loadGateFSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(impactModuleRoot(t), "evals", "impact", "schema-gate-f.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema-gate-f.json: %v", err)
	}
	return sch
}

func writeGateFMetrics(t *testing.T, path string, m MetricsGateF) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validateGateFMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) {
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
