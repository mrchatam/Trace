package eval_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/eval"
	"github.com/mrchatam/Trace/internal/store"
)

func evalTestdata(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("testdata", name)
}

func writeEvalRulesFile(t *testing.T, root, srcFixture string) {
	t.Helper()
	src, err := os.ReadFile(evalTestdata(t, srcFixture))
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "trace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "eval-rules.json"), src, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestProjectEvalRulesLoaded(t *testing.T) {
	svc, st := openEvalDomain(t)
	root := st.ProjectRoot()
	writeEvalRulesFile(t, root, "eval-rules-default.json")

	ctx := context.Background()
	load, err := eval.LoadRules(ctx, root, st)
	if err != nil {
		t.Fatalf("LoadRules: %v", err)
	}
	if !load.Loaded {
		t.Fatal("expected loaded=true")
	}
	if load.Path != "trace/eval-rules.json" {
		t.Fatalf("path: %q", load.Path)
	}
	if len(load.Mechanisms) != 4 {
		t.Fatalf("mechanisms: %v", load.Mechanisms)
	}

	row, err := st.GetEvalRuleSet(store.EvalRuleSetDefaultID)
	if err != nil {
		t.Fatalf("GetEvalRuleSet: %v", err)
	}
	if row.SourcePath != "trace/eval-rules.json" {
		t.Fatalf("source_path: %q", row.SourcePath)
	}
	var body eval.RulesFile
	if err := json.Unmarshal([]byte(row.BodyJSON), &body); err != nil {
		t.Fatal(err)
	}
	if body.Version != 1 || len(body.Mechanisms) != 4 {
		t.Fatalf("cached body: %+v", body)
	}
	if load.CachedAt == "" {
		t.Fatal("expected cached_at")
	}
	_ = svc
}

func TestMissingEvalRulesUsesBuiltins(t *testing.T) {
	svc, st := openEvalDomain(t)
	root := st.ProjectRoot()

	ctx := context.Background()
	load, err := eval.LoadRules(ctx, root, st)
	if err != nil {
		t.Fatalf("LoadRules: %v", err)
	}
	if load.Loaded {
		t.Fatal("expected loaded=false when file missing")
	}
	want := []string{
		eval.MechanismStoredTest,
		eval.MechanismStoredVerification,
		eval.MechanismStoredEvaluation,
		eval.MechanismArchitecturalInvariant,
	}
	if len(load.Mechanisms) != len(want) {
		t.Fatalf("mechanisms: got %v want %v", load.Mechanisms, want)
	}
	for i, id := range want {
		if load.Mechanisms[i] != id {
			t.Fatalf("mechanism[%d]: got %q want %q", i, load.Mechanisms[i], id)
		}
	}

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	results, err := eval.RunAll(ctx, eval.EvalInput{TaskID: task.ID, Service: svc}, eval.RunOptions{
		Root:  root,
		Store: st,
	})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("want 4 mechanism results, got %d: %+v", len(results), results)
	}
}

func indexInvariantViolationFixture(t *testing.T, svc *domain.Service, st *store.Store, taskID string) {
	t.Helper()
	ctx := context.Background()
	root := st.ProjectRoot()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join("..", "analyzers", "testdata")
	cmdMain, err := os.ReadFile(filepath.Join(testdata, "architecture_cmd_main.go"))
	if err != nil {
		t.Fatal(err)
	}
	libSrc, err := os.ReadFile(filepath.Join(testdata, "architecture_internal_imports_cmd.go"))
	if err != nil {
		t.Fatal(err)
	}
	if err := analyzers.IndexFile(ctx, st, "cmd/app/main.go", cmdMain, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile cmd: %v", err)
	}
	if err := analyzers.IndexFile(ctx, st, "internal/lib/lib.go", libSrc, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile internal: %v", err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "internal/lib/lib.go"}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestProjectEvalRulesOverrideDefaultInvariant(t *testing.T) {
	svc, st := openEvalDomain(t)
	root := st.ProjectRoot()
	writeEvalRulesFile(t, root, "eval-rules-disable-invariant.json")

	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	indexInvariantViolationFixture(t, svc, st, task.ID)

	check, err := svc.CheckArchitecturalInvariants(ctx, task.ID)
	if err != nil {
		t.Fatalf("domain check: %v", err)
	}
	if check.Passed {
		t.Fatal("domain must still detect violation")
	}

	results, err := eval.RunAll(ctx, eval.EvalInput{TaskID: task.ID, Service: svc}, eval.RunOptions{
		Root:  root,
		Store: st,
	})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	var arch eval.EvalResult
	for _, r := range results {
		if r.MechanismID == eval.MechanismArchitecturalInvariant {
			arch = r
			break
		}
	}
	if arch.MechanismID == "" {
		t.Fatalf("missing architectural_invariant in %+v", results)
	}
	if !arch.Passed {
		t.Fatalf("disabled invariant must pass mechanism: %+v", arch)
	}
}

func TestLoadRulesInvalidJSONFailClosed(t *testing.T) {
	_, st := openEvalDomain(t)
	root := st.ProjectRoot()
	dir := filepath.Join(root, "trace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "eval-rules.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := eval.LoadRules(context.Background(), root, st)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadRulesWrongVersionFailClosed(t *testing.T) {
	_, st := openEvalDomain(t)
	root := st.ProjectRoot()
	dir := filepath.Join(root, "trace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	raw := []byte(`{"version":2,"mechanisms":["stored_test"],"invariants":[]}`)
	if err := os.WriteFile(filepath.Join(dir, "eval-rules.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := eval.LoadRules(context.Background(), root, st)
	if err == nil {
		t.Fatal("expected error for wrong version")
	}
}
