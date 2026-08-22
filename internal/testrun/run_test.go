package testrun

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func openTestrun(t *testing.T) (*store.Store, *domain.Service) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/testrun\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st, domain.New(st)
}

type stubRunner struct {
	calls []RunSpec
	exit  int
	out   string
	err   error
}

func (s *stubRunner) Run(ctx context.Context, spec RunSpec) (int, string, error) {
	s.calls = append(s.calls, spec)
	return s.exit, s.out, s.err
}

func seedValidatesGraph(t *testing.T, st *store.Store) (libPath, testName string) {
	t.Helper()
	lib, err := st.UpsertFile("pkg/foo.go", "hlib", nil)
	if err != nil {
		t.Fatal(err)
	}
	testFile, err := st.UpsertFile("pkg/foo_test.go", "htest", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("pkg/foo.go", []store.Symbol{
		{Name: "Foo", Kind: "function", StartLine: 1, EndLine: 5},
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("pkg/foo_test.go", []store.Symbol{
		{Name: "TestFoo", Kind: "test", StartLine: 3, EndLine: 10},
	}); err != nil {
		t.Fatal(err)
	}
	prod, err := st.ListSymbolsByPath("pkg/foo.go")
	if err != nil || len(prod) != 1 {
		t.Fatalf("prod symbols: %v %v", prod, err)
	}
	tests, err := st.ListSymbolsByPath("pkg/foo_test.go")
	if err != nil || len(tests) != 1 {
		t.Fatalf("test symbols: %v %v", tests, err)
	}
	fromID := tests[0].ID
	toID := prod[0].ID
	if err := st.ReplaceFileEdges("pkg/foo_test.go", []store.CodeEdge{
		{
			FromFileID: testFile.ID, FromSymbolID: &fromID,
			ToFileID: lib.ID, ToSymbolID: &toID,
			Rel: store.RelValidates, Provenance: store.ImportProvenanceInferred,
		},
	}); err != nil {
		t.Fatal(err)
	}
	return "pkg/foo.go", "TestFoo"
}

func mustTask(t *testing.T, dom *domain.Service) store.Task {
	t.Helper()
	ctx := context.Background()
	g, err := dom.CreateGoal(ctx, domain.GoalInput{Title: "test goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dom.CreateTask(ctx, domain.TaskInput{Title: "work", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func TestTestRunRecordsOutcome(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, testName := seedValidatesGraph(t, st)

	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}

	stub := &stubRunner{exit: 0, out: "ok"}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(outcomes) != 1 {
		t.Fatalf("outcomes=%d want 1", len(outcomes))
	}
	if outcomes[0].Kind != store.OutcomeKindTest || outcomes[0].TestName != testName {
		t.Fatalf("row: %+v", outcomes[0])
	}
	if outcomes[0].TestStatus != store.TestStatusPass {
		t.Fatalf("status=%q want pass", outcomes[0].TestStatus)
	}
	rows, err := st.ListOutcomeResultsByTaskKind(task.ID, store.OutcomeKindTest)
	if err != nil || len(rows) != 1 {
		t.Fatalf("stored rows: %d err=%v", len(rows), err)
	}
}

func TestTestRunSelectsValidatingTests(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, testName := seedValidatesGraph(t, st)

	if err := os.MkdirAll(filepath.Join(st.ProjectRoot(), "other"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpsertFile("other/unrelated_test.go", "hu", nil); err != nil {
		t.Fatal(err)
	}

	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}

	stub := &stubRunner{exit: 0, out: "ok"}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatal(err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("calls=%d want 1 validating test only", len(stub.calls))
	}
	if stub.calls[0].Command != "go" {
		t.Fatalf("command=%q want go", stub.calls[0].Command)
	}
	foundRun := false
	for i, a := range stub.calls[0].Args {
		if a == "-run" && i+1 < len(stub.calls[0].Args) && stub.calls[0].Args[i+1] == "^"+testName+"$" {
			foundRun = true
		}
	}
	if !foundRun {
		t.Fatalf("args=%v want -run ^TestFoo$", stub.calls[0].Args)
	}
}

func TestTestRunFailClosedWithoutCommand(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	dom := domain.New(st)
	task := mustTask(t, dom)

	_, err = RunRelevantTests(context.Background(), st, dom, task.ID, Options{Runner: &stubRunner{}})
	if err == nil {
		t.Fatal("expected error without go.mod/config")
	}
	rows, err := st.ListOutcomeResultsByTaskKind(task.ID, store.OutcomeKindTest)
	if err != nil || len(rows) != 0 {
		t.Fatalf("must not record fake pass: rows=%d err=%v", len(rows), err)
	}
}

func TestTestRunUsesStubRunner(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, _ := seedValidatesGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}
	stub := &stubRunner{exit: 1, out: "fail output"}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatal(err)
	}
	if len(stub.calls) == 0 {
		t.Fatal("stub runner must be invoked")
	}
}

func TestTestRunClearsTestPending(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, _ := seedValidatesGraph(t, st)

	g, err := st.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if g.GoalID == nil {
		t.Fatal("task missing goal")
	}
	goalID := *g.GoalID
	if _, err := st.UpsertDeliberationState(store.DeliberationState{
		TaskID: task.ID, GoalID: goalID, PlanCritiqued: true,
	}); err != nil {
		t.Fatal(err)
	}
	psvc := planner.New(st)

	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}

	before, err := loop.BuildPolicyInputs(ctx, dom, psvc, task.ID, goalID, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if !before.TestPending {
		t.Fatalf("want test_pending before run: %#v", before)
	}

	stub := &stubRunner{exit: 0, out: "ok"}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatal(err)
	}

	after, err := loop.BuildPolicyInputs(ctx, dom, psvc, task.ID, goalID, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if after.TestPending {
		t.Fatalf("test_pending must clear after recorded outcome: %#v", after)
	}
}

func TestSelectTargetsFromImpactWalk(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	libPath, testName := seedValidatesGraph(t, st)

	targets, err := SelectTestTargets(ctx, st, dom, "unused", []string{libPath})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Name != testName {
		t.Fatalf("targets=%+v want TestFoo", targets)
	}
}

func TestPackageFallbackTargets(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "internal/new/feature.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	targets, err := SelectTestTargets(ctx, st, dom, task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("targets=%+v", targets)
	}
	if targets[0].Name != "package:example.com/testrun/internal/new" {
		t.Fatalf("name=%q", targets[0].Name)
	}
}

func TestImpactWalkStillWorks(t *testing.T) {
	st, _ := openTestrun(t)
	eng := retrieval.New(st)
	libPath, _ := seedValidatesGraph(t, st)
	f, err := st.GetFileByPath(libPath)
	if err != nil {
		t.Fatal(err)
	}
	res, err := eng.ImpactWalk(context.Background(), []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: f.ID},
	}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.AffectedTests) != 1 {
		t.Fatalf("affected_tests=%+v", res.AffectedTests)
	}
}
