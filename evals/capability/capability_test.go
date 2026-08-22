package capability_test

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type harness struct {
	svc *domain.Service
	cmp *compiler.Compiler
}

func openCapability(t *testing.T) harness {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return harness{svc: domain.New(st), cmp: compiler.New(st)}
}

func plantTask(t *testing.T, svc *domain.Service, ctx context.Context, title string) string {
	t.Helper()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: title})
	if err != nil {
		t.Fatalf("CreateTask %q: %v", title, err)
	}
	return task.ID
}

func capsContainSlug(caps []store.Capability, slug string) bool {
	for _, c := range caps {
		if c.Slug == slug || c.ID == slug {
			return true
		}
	}
	return false
}

func refsContainSlug(refs []compiler.CapabilityRef, slug string) bool {
	for _, r := range refs {
		if r.Slug == slug || r.ID == slug {
			return true
		}
	}
	return false
}

func refSlugs(refs []compiler.CapabilityRef) map[string]struct{} {
	out := make(map[string]struct{}, len(refs))
	for _, r := range refs {
		out[r.Slug] = struct{}{}
	}
	return out
}

// MetricsCapability is the capability-selection ablation artifact (schema-capability.json v1).
type MetricsCapability struct {
	SchemaVersion  int      `json:"schema_version"`
	Gate           string   `json:"gate"`
	Suite          string   `json:"suite"`
	Ablation       bool     `json:"ablation"`
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

// TestPlantedCapabilitySelectionAblation is the capability-selection precision/recall harness.
// Plants Pos-1..3 selection/missing probes + Neg-1 clean AVAILABLE; scores MissingCapabilities + packet GT.
func TestPlantedCapabilitySelectionAblation(t *testing.T) {
	h := openCapability(t)
	ctx := context.Background()

	var tp, fn, fp, tn int
	probeIDs := []string{
		"pos-1-unavailable",
		"pos-2-unknown",
		"pos-3-selection-filter",
		"neg-1-clean-available",
	}

	// Part 1 — Pos-1: UNAVAILABLE missing warning
	{
		taskID := plantTask(t, h.svc, ctx, "Cap Pos-1 UNAVAILABLE")
		cap, err := h.svc.UpsertCapability(ctx, domain.CapabilityInput{
			Kind:   domain.CapabilityKindTool,
			Slug:   "tool:down",
			Title:  "Down",
			Status: domain.CapabilityStatusUnavailable,
		})
		if err != nil {
			t.Fatalf("Pos-1 UpsertCapability: %v", err)
		}
		if _, err := h.svc.RequireCapability(ctx, taskID, cap.ID); err != nil {
			t.Fatalf("Pos-1 RequireCapability: %v", err)
		}
		miss, err := h.svc.MissingCapabilities(ctx, taskID)
		if err != nil {
			t.Fatalf("Pos-1 MissingCapabilities: %v", err)
		}
		pkt, err := h.cmp.TaskContext(ctx, taskID, compiler.ContextOptions{})
		if err != nil {
			t.Fatalf("Pos-1 TaskContext: %v", err)
		}
		ok := capsContainSlug(miss, "tool:down") &&
			refsContainSlug(pkt.MissingCapabilities, "tool:down") &&
			refsContainSlug(pkt.RequiredCapabilities, "tool:down")
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-1 FN: want miss+packet required+missing contain tool:down; miss=%+v req=%+v missing=%+v",
				miss, pkt.RequiredCapabilities, pkt.MissingCapabilities)
		}
	}

	// Part 2 — Pos-2: UNKNOWN missing warning
	{
		taskID := plantTask(t, h.svc, ctx, "Cap Pos-2 UNKNOWN")
		cap, err := h.svc.UpsertCapability(ctx, domain.CapabilityInput{
			Kind:   domain.CapabilityKindSkill,
			Slug:   "skill:maybe",
			Title:  "Maybe",
			Status: domain.CapabilityStatusUnknown,
		})
		if err != nil {
			t.Fatalf("Pos-2 UpsertCapability: %v", err)
		}
		if _, err := h.svc.RequireCapability(ctx, taskID, cap.ID); err != nil {
			t.Fatalf("Pos-2 RequireCapability: %v", err)
		}
		miss, err := h.svc.MissingCapabilities(ctx, taskID)
		if err != nil {
			t.Fatalf("Pos-2 MissingCapabilities: %v", err)
		}
		pkt, err := h.cmp.TaskContext(ctx, taskID, compiler.ContextOptions{})
		if err != nil {
			t.Fatalf("Pos-2 TaskContext: %v", err)
		}
		ok := capsContainSlug(miss, "skill:maybe") &&
			refsContainSlug(pkt.MissingCapabilities, "skill:maybe") &&
			refsContainSlug(pkt.RequiredCapabilities, "skill:maybe")
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-2 FN: want miss+packet required+missing contain skill:maybe; miss=%+v req=%+v missing=%+v",
				miss, pkt.RequiredCapabilities, pkt.MissingCapabilities)
		}
	}

	// Part 3 — Pos-3: selection filter (no catalog dump)
	{
		taskID := plantTask(t, h.svc, ctx, "Cap Pos-3 selection-filter")
		// Seed builtin MCP specs into catalog (AVAILABLE) plus extra TOOL/SKILL entries.
		var catalogSlugs []string
		var requiredIDs []string
		var requiredSlugs []string
		for _, spec := range domain.BuiltinMCPCapabilitySpecs() {
			c, err := h.svc.UpsertCapability(ctx, domain.CapabilityInput{
				Kind:   spec.Kind,
				Slug:   spec.Slug,
				Title:  spec.Title,
				Status: domain.CapabilityStatusAvailable,
			})
			if err != nil {
				t.Fatalf("Pos-3 UpsertBuiltin %s: %v", spec.Slug, err)
			}
			catalogSlugs = append(catalogSlugs, c.Slug)
		}
		extras := []domain.CapabilityInput{
			{Kind: domain.CapabilityKindTool, Slug: "tool:alpha", Title: "Alpha", Status: domain.CapabilityStatusAvailable},
			{Kind: domain.CapabilityKindTool, Slug: "tool:beta", Title: "Beta", Status: domain.CapabilityStatusAvailable},
			{Kind: domain.CapabilityKindSkill, Slug: "skill:gamma", Title: "Gamma", Status: domain.CapabilityStatusAvailable},
		}
		for _, in := range extras {
			c, err := h.svc.UpsertCapability(ctx, in)
			if err != nil {
				t.Fatalf("Pos-3 UpsertCapability %s: %v", in.Slug, err)
			}
			catalogSlugs = append(catalogSlugs, c.Slug)
			if in.Slug == "tool:alpha" || in.Slug == "tool:beta" {
				requiredIDs = append(requiredIDs, c.ID)
				requiredSlugs = append(requiredSlugs, c.Slug)
			}
		}
		if len(catalogSlugs) < 3 {
			t.Fatalf("Pos-3 catalog want ≥3 AVAILABLE; got %d", len(catalogSlugs))
		}
		for _, id := range requiredIDs {
			if _, err := h.svc.RequireCapability(ctx, taskID, id); err != nil {
				t.Fatalf("Pos-3 RequireCapability: %v", err)
			}
		}
		miss, err := h.svc.MissingCapabilities(ctx, taskID)
		if err != nil {
			t.Fatalf("Pos-3 MissingCapabilities: %v", err)
		}
		pkt, err := h.cmp.TaskContext(ctx, taskID, compiler.ContextOptions{})
		if err != nil {
			t.Fatalf("Pos-3 TaskContext: %v", err)
		}
		reqSet := refSlugs(pkt.RequiredCapabilities)
		wantSet := map[string]struct{}{"tool:alpha": {}, "tool:beta": {}}
		onlyRequired := len(pkt.RequiredCapabilities) == 2
		for s := range reqSet {
			if _, ok := wantSet[s]; !ok {
				onlyRequired = false
			}
		}
		for s := range wantSet {
			if _, ok := reqSet[s]; !ok {
				onlyRequired = false
			}
		}
		// Non-required catalog slug must not appear in packet required.
		noDump := true
		for _, slug := range catalogSlugs {
			if slug == "tool:alpha" || slug == "tool:beta" {
				continue
			}
			if refsContainSlug(pkt.RequiredCapabilities, slug) {
				noDump = false
				break
			}
		}
		ok := onlyRequired &&
			len(pkt.MissingCapabilities) == 0 &&
			len(miss) == 0 &&
			noDump &&
			len(requiredSlugs) == 2
		if ok {
			tp++
		} else {
			fn++
			t.Errorf("Pos-3 FN: want exactly 2 required (tool:alpha,tool:beta), no missing, no catalog dump; miss=%d req=%+v missing=%+v",
				len(miss), pkt.RequiredCapabilities, pkt.MissingCapabilities)
		}
	}

	// Part 4 — Neg-1: clean AVAILABLE (no false missing)
	{
		taskID := plantTask(t, h.svc, ctx, "Cap Neg-1 clean AVAILABLE")
		cap, err := h.svc.UpsertCapability(ctx, domain.CapabilityInput{
			Kind:   domain.CapabilityKindTool,
			Slug:   "tool:ok",
			Title:  "OK",
			Status: domain.CapabilityStatusAvailable,
		})
		if err != nil {
			t.Fatalf("Neg-1 UpsertCapability: %v", err)
		}
		if _, err := h.svc.RequireCapability(ctx, taskID, cap.ID); err != nil {
			t.Fatalf("Neg-1 RequireCapability: %v", err)
		}
		miss, err := h.svc.MissingCapabilities(ctx, taskID)
		if err != nil {
			t.Fatalf("Neg-1 MissingCapabilities: %v", err)
		}
		pkt, err := h.cmp.TaskContext(ctx, taskID, compiler.ContextOptions{})
		if err != nil {
			t.Fatalf("Neg-1 TaskContext: %v", err)
		}
		reqOK := len(pkt.RequiredCapabilities) == 1 && refsContainSlug(pkt.RequiredCapabilities, "tool:ok")
		ok := len(miss) == 0 && len(pkt.MissingCapabilities) == 0 && reqOK
		if ok {
			tn++
		} else {
			fp++
			t.Errorf("Neg-1 FP: want miss=0 packet.missing=0 required=[tool:ok]; miss=%d req=%+v missing=%+v",
				len(miss), pkt.RequiredCapabilities, pkt.MissingCapabilities)
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
	metricsPath := filepath.Join(metricsDir, "metrics-capability.json")
	m := MetricsCapability{
		SchemaVersion:  1,
		Gate:           "capability-selection",
		Suite:          "capability",
		Ablation:       true,
		DryRun:         false,
		TruePositives:  tp,
		FalsePositives: fp,
		FalseNegatives: fn,
		TrueNegatives:  tn,
		Precision:      precision,
		Recall:         recall,
		Probes:         probes,
		NamedTest:      "TestPlantedCapabilitySelectionAblation",
		Mig:            "010_capability_surface",
		S01Hooks: []string{
			"UpsertCapability",
			"RequireCapability",
			"MissingCapabilities",
			"required_capabilities",
			"missing_capabilities",
		},
		ProbeIDs:     probeIDs,
		TraceVersion: "0.0.0-dev",
	}
	writeCapabilityMetrics(t, metricsPath, m)
	validateCapabilityMetricsFile(t, loadCapabilitySchema(t), metricsPath)
}

func findCapabilityModuleRoot() string {
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

func capabilityModuleRoot(t *testing.T) string {
	t.Helper()
	root := findCapabilityModuleRoot()
	if root == "" {
		t.Fatal("go.mod not found above evals/capability")
	}
	return root
}

func loadCapabilitySchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	schemaPath := filepath.Join(capabilityModuleRoot(t), "evals", "capability", "schema-capability.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatal(err)
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile("file://" + filepath.ToSlash(abs))
	if err != nil {
		t.Fatalf("compile schema-capability.json: %v", err)
	}
	return sch
}

func writeCapabilityMetrics(t *testing.T, path string, m MetricsCapability) {
	t.Helper()
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validateCapabilityMetricsFile(t *testing.T, sch *jsonschema.Schema, path string) {
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
