package eval_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/eval"
	"github.com/mrchatam/Trace/internal/store"
)

func openEvalDomain(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st), st
}

func TestEvalRegistryMultipleMechanisms(t *testing.T) {
	ids := eval.DefaultRegistry().ListMechanismIDs()
	if len(ids) < 4 {
		t.Fatalf("want >=4 mechanisms, got %d: %v", len(ids), ids)
	}
	want := map[string]bool{
		eval.MechanismStoredTest:             false,
		eval.MechanismStoredVerification:     false,
		eval.MechanismStoredEvaluation:       false,
		eval.MechanismArchitecturalInvariant: false,
	}
	for _, id := range ids {
		if _, ok := want[id]; ok {
			want[id] = true
		}
	}
	for id, seen := range want {
		if !seen {
			t.Fatalf("missing built-in mechanism %q in registry: %v", id, ids)
		}
	}
}

type fakeEcho struct {
	msg string
}

func (f fakeEcho) ID() string { return "fake_echo" }

func (f fakeEcho) Run(ctx context.Context, in eval.EvalInput) (eval.EvalResult, error) {
	_ = ctx
	_ = in
	return eval.EvalResult{
		MechanismID: f.ID(),
		Passed:      true,
		Summary:     f.msg,
	}, nil
}

func TestAddMechanismWithoutSchemaChange(t *testing.T) {
	eval.Register(fakeEcho{msg: "echo-ok"})
	t.Cleanup(func() {
		// fake_echo persists in default registry for process lifetime — acceptable for extensibility proof.
	})

	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "eval extensibility"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	results, err := eval.RunAll(ctx, eval.EvalInput{TaskID: task.ID, Service: svc}, eval.RunOptions{
		MechanismIDs: []string{"fake_echo"},
	})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	if len(results) != 1 || results[0].MechanismID != "fake_echo" || !results[0].Passed || results[0].Summary != "echo-ok" {
		t.Fatalf("fake_echo result: %+v err=%v", results, err)
	}

	dir := schemaDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	sqlCount := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlCount++
		}
	}
	if sqlCount != 28 {
		t.Fatalf("schema sql file count: got %d want 28", sqlCount)
	}

	// outcome_results kind CHECK unchanged (018 migration).
	body, err := os.ReadFile(filepath.Join(dir, "018_outcome_results_baselines.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "CHECK (kind IN ('test', 'verification', 'evaluation'))") {
		t.Fatal("outcome_results kind CHECK changed")
	}
}

func TestRunAllInvalidTaskIDFailClosed(t *testing.T) {
	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	_, err := eval.RunAll(ctx, eval.EvalInput{TaskID: " ", Service: svc}, eval.RunOptions{})
	if err == nil {
		t.Fatal("expected error for empty task_id")
	}
}

func TestRunAllContinuesOnMechanismFailure(t *testing.T) {
	eval.Register(fakeEcho{msg: "always-ok"})
	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	results, err := eval.RunAll(ctx, eval.EvalInput{TaskID: task.ID, Service: svc}, eval.RunOptions{})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	if len(results) < 4 {
		t.Fatalf("want >=4 results, got %d", len(results))
	}
	for _, r := range results {
		if r.MechanismID == eval.MechanismStoredTest && r.Passed {
			t.Fatal("stored_test should fail without change+test row")
		}
		if r.MechanismID == "fake_echo" && !r.Passed {
			t.Fatalf("fake_echo should pass: %+v", r)
		}
	}
}

func TestRunAllFiltersUnknownMechanismIDs(t *testing.T) {
	svc, _ := openEvalDomain(t)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	results, err := eval.RunAll(ctx, eval.EvalInput{TaskID: task.ID, Service: svc}, eval.RunOptions{
		MechanismIDs: []string{"nonexistent_mech", eval.MechanismStoredTest},
	})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	if len(results) != 1 || results[0].MechanismID != eval.MechanismStoredTest {
		t.Fatalf("filtered results: %+v", results)
	}
}

func schemaDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "internal", "store", "schema")
}
