package domain_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestSeedExportIncludesEvalRulesPath(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	root := st.ProjectRoot()

	traceDir := filepath.Join(root, "trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rules := []byte(`{"version":1,"mechanisms":["stored_test"],"invariants":[]}`)
	if err := os.WriteFile(filepath.Join(traceDir, "eval-rules.json"), rules, 0o644); err != nil {
		t.Fatal(err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{ProjectRoot: root})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}
	if doc.EvalRulesPath != "trace/eval-rules.json" {
		t.Fatalf("EvalRulesPath: %q", doc.EvalRulesPath)
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed["eval_rules_path"]; !ok {
		t.Fatal("export missing eval_rules_path key")
	}
	for _, key := range []string{"eval_rule_sets", "mechanisms"} {
		if _, ok := parsed[key]; ok {
			t.Fatalf("seed must not embed rules body key %q", key)
		}
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	svc2 := domain.New(st2)
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("ImportSeedDocument: %v", err)
	}
	row, err := st2.GetEvalRuleSet(store.EvalRuleSetDefaultID)
	if err != nil {
		t.Fatalf("GetEvalRuleSet: %v", err)
	}
	if row.SourcePath != "trace/eval-rules.json" {
		t.Fatalf("import source_path: %q", row.SourcePath)
	}
	if row.BodyJSON != "{}" {
		t.Fatalf("import must store pointer-only body: %q", row.BodyJSON)
	}
	_ = svc
}

func TestSeedExportOmitsEvalRulesPathWhenMissing(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{ProjectRoot: st.ProjectRoot()})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}
	if doc.EvalRulesPath != "" {
		t.Fatalf("EvalRulesPath should be empty: %q", doc.EvalRulesPath)
	}
	_ = svc
}
