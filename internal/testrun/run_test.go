package testrun

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

// writeConfig writes a test-runner.json at rel (relative to project root).
func writeConfig(t *testing.T, dir, rel, contents string) {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestTestRunUsesCustomRunnerConfig proves the configured command is invoked
// (regression: root-level test-runner.json was ignored entirely — RunRelevantTests
// fell back to go test even with a valid config present).
func TestTestRunUsesCustomRunnerConfig(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, testName := seedValidatesGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json",
		`{"command":"custom-runner","args":["run","all"]}`)

	stub := &stubRunner{exit: 0, out: "ok"}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(outcomes) != 1 || outcomes[0].TestName != testName {
		t.Fatalf("outcomes=%+v want 1 %q", outcomes, testName)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("calls=%d want 1", len(stub.calls))
	}
	got := stub.calls[0]
	if got.Command != "custom-runner" {
		t.Fatalf("command=%q want custom-runner (config must not be overridden by go test)", got.Command)
	}
	if len(got.Args) != 2 || got.Args[0] != "run" || got.Args[1] != "all" {
		t.Fatalf("args=%v want [run all] (raw config, no go-style rewriting)", got.Args)
	}
}

// TestTestRunCustomConfigAtProjectRoot proves the root-level convenience location.
func TestTestRunCustomConfigAtProjectRoot(t *testing.T) {
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

	writeConfig(t, st.ProjectRoot(), "test-runner.json",
		`{"command":"root-runner","args":["-x"]}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 || stub.calls[0].Command != "root-runner" {
		t.Fatalf("calls=%+v want root-runner", stub.calls)
	}
}

// TestTestRunCustomConfigCanonicalTraceDir proves the canonical P22 location
// (trace/test-runner.json) keeps working end to end.
func TestTestRunCustomConfigCanonicalTraceDir(t *testing.T) {
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

	writeConfig(t, st.ProjectRoot(), filepath.Join("trace", "test-runner.json"),
		`{"command":"trace-dir-runner"}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 || stub.calls[0].Command != "trace-dir-runner" {
		t.Fatalf("calls=%+v want trace-dir-runner", stub.calls)
	}
}

// TestTestRunCustomConfigWinsOverGoMod proves precedence config > go fallback.
func TestTestRunCustomConfigWinsOverGoMod(t *testing.T) {
	st, dom := openTestrun(t) // openTestrun writes go.mod
	ctx := context.Background()
	task := mustTask(t, dom)
	libPath, _ := seedValidatesGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: libPath}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"cfg-runner"}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 || stub.calls[0].Command != "cfg-runner" {
		t.Fatalf("calls=%+v want cfg-runner (config beats go.mod)", stub.calls)
	}
}

// TestTestRunCustomConfigNoTargetsErrors is the non-Go fail-closed contract
// (the exact reported scenario): no go.mod, root-level test-runner.json present,
// task with no changes and no seed paths. The config must be honored (error is
// about target selection, not a missing config), and no outcome may be recorded.
func TestTestRunCustomConfigNoTargetsErrors(t *testing.T) {
	dir := t.TempDir() // no go.mod — non-Go repo
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	dom := domain.New(st)
	task := mustTask(t, dom)
	writeConfig(t, dir, "test-runner.json", `{"command":"never-runner"}`)

	_, err = RunRelevantTests(context.Background(), st, dom, task.ID, Options{Runner: &stubRunner{}})
	if err == nil || !strings.Contains(err.Error(), "no relevant tests selected") {
		t.Fatalf("err=%v want no relevant tests selected", err)
	}
	rows, err := st.ListOutcomeResultsByTaskKind(task.ID, store.OutcomeKindTest)
	if err != nil || len(rows) != 0 {
		t.Fatalf("must not record outcomes when nothing selected: rows=%d err=%v", len(rows), err)
	}
}

// TestTestRunCustomConfigPathPlaceholder proves {path} substitution lets a
// custom runner execute only the selected test file (relevant tests, not the suite).
func TestTestRunCustomConfigPathPlaceholder(t *testing.T) {
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

	writeConfig(t, st.ProjectRoot(), "test-runner.json",
		`{"command":"vitest","args":["run","{path}","--reporter=dot"]}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("calls=%d want 1", len(stub.calls))
	}
	got := stub.calls[0]
	wantArgs := []string{"run", "pkg/foo_test.go", "--reporter=dot"}
	if got.Command != "vitest" || !reflect.DeepEqual(got.Args, wantArgs) {
		t.Fatalf("spec=%+v want vitest %v", got, wantArgs)
	}
}

// TestResolveDefaultRunnerCustomCwd proves custom cwd resolves relative to root.
func TestResolveDefaultRunnerCustomCwd(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "test-runner.json")
	if err := os.WriteFile(path, []byte(`{"command":"c","cwd":"sub"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec, source, err := resolveDefaultRunner(root)
	if err != nil || source != runnerSourceCustom {
		t.Fatalf("source=%v err=%v", source, err)
	}
	if spec.Cwd != filepath.Join(root, "sub") {
		t.Fatalf("cwd=%q want %q", spec.Cwd, filepath.Join(root, "sub"))
	}
}

// TestTestRunCustomConfigCommandGo is the edge where a custom config's command
// is literally `go` — the runnerSource-based classification must still treat it
// as custom (args passed through unrewritten, no package fallback semantics).
func TestTestRunCustomConfigCommandGo(t *testing.T) {
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

	writeConfig(t, st.ProjectRoot(), "test-runner.json",
		`{"command":"go","args":["test","./pkg/..."]}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("calls=%d want 1", len(stub.calls))
	}
	got := stub.calls[0]
	wantArgs := []string{"test", "./pkg/..."}
	if got.Command != "go" || !reflect.DeepEqual(got.Args, wantArgs) {
		t.Fatalf("spec=%+v want go %v (config args must pass through unrewritten)", got, wantArgs)
	}
}

// TestTestRunCustomConfigWhenTraceIsFile proves the loader falls through to the
// root-level config when trace/ exists as a regular file (ENOTDIR, not NotExist).
func TestTestRunCustomConfigWhenTraceIsFile(t *testing.T) {
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

	writeConfig(t, st.ProjectRoot(), "trace", "not a directory")
	t.Cleanup(func() { _ = os.Remove(filepath.Join(st.ProjectRoot(), "trace")) })
	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"root-runner"}`)

	stub := &stubRunner{exit: 0}
	if _, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub}); err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 || stub.calls[0].Command != "root-runner" {
		t.Fatalf("calls=%+v want root-runner (trace-as-file must not break root config)", stub.calls)
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
