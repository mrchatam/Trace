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
	libPath, _ := seedValidatesGraph(t, st)
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
	// The stub's output carries no per-file report, so the whole-run verdict
	// is recorded once under a synthetic suite name — never under TestFoo.
	if len(outcomes) != 1 || !strings.HasPrefix(outcomes[0].TestName, "suite:") {
		t.Fatalf("outcomes=%+v want single suite outcome", outcomes)
	}
	if outcomes[0].TestStatus != store.TestStatusPass {
		t.Fatalf("status=%q want pass", outcomes[0].TestStatus)
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

// runConfiguredRunner seeds a second validating test so the selection has two
// real targets, then returns a stubRunner for custom-config runs.
func seedSecondValidatesGraph(t *testing.T, st *store.Store) (string, string) {
	t.Helper()
	lib, err := st.UpsertFile("pkg/bar.go", "hlib2", nil)
	if err != nil {
		t.Fatal(err)
	}
	testFile, err := st.UpsertFile("pkg/bar_test.go", "htest2", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("pkg/bar.go", []store.Symbol{
		{Name: "Bar", Kind: "function", StartLine: 1, EndLine: 5},
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("pkg/bar_test.go", []store.Symbol{
		{Name: "TestBar", Kind: "test", StartLine: 3, EndLine: 10},
	}); err != nil {
		t.Fatal(err)
	}
	prod, err := st.ListSymbolsByPath("pkg/bar.go")
	if err != nil || len(prod) != 1 {
		t.Fatalf("prod symbols: %v %v", prod, err)
	}
	tests, err := st.ListSymbolsByPath("pkg/bar_test.go")
	if err != nil || len(tests) != 1 {
		t.Fatalf("test symbols: %v %v", tests, err)
	}
	fromID := tests[0].ID
	toID := prod[0].ID
	if err := st.ReplaceFileEdges("pkg/bar_test.go", []store.CodeEdge{
		{
			FromFileID: testFile.ID, FromSymbolID: &fromID,
			ToFileID: lib.ID, ToSymbolID: &toID,
			Rel: store.RelValidates, Provenance: store.ImportProvenanceInferred,
		},
	}); err != nil {
		t.Fatal(err)
	}
	return "pkg/bar.go", "TestBar"
}

// seedVitestGraph seeds two vitest-style targets (src/foo.test.ts,
// src/bar.test.ts) with validates edges so SelectTestTargets returns real
// custom-runner targets whose names are file basenames. Returns the seed
// change paths.
func seedVitestGraph(t *testing.T, st *store.Store) ([]string, string, string) {
	t.Helper()
	lib, err := st.UpsertFile("src/core.ts", "hcore", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("src/core.ts", []store.Symbol{
		{Name: "core", Kind: "function", StartLine: 1, EndLine: 5},
	}); err != nil {
		t.Fatal(err)
	}
	prod, err := st.ListSymbolsByPath("src/core.ts")
	if err != nil || len(prod) != 1 {
		t.Fatalf("prod symbols: %v %v", prod, err)
	}
	coreSym := prod[0].ID
	coreFile, err := st.GetFileByPath("src/core.ts")
	if err != nil {
		t.Fatal(err)
	}
	// Symbol names mirror the real indexed-repo shape: vitest file symbols are
	// named after the file basename, so targets carry basename names while the
	// runner reports repo-relative paths.
	for _, tc := range []struct{ path, symbol string }{
		{"src/foo.test.ts", "foo.test.ts"},
		{"src/bar.test.ts", "bar.test.ts"},
	} {
		tf, err := st.UpsertFile(tc.path, "h"+tc.symbol, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := st.ReplaceFileSymbols(tc.path, []store.Symbol{
			{Name: tc.symbol, Kind: "test", StartLine: 1, EndLine: 4},
		}); err != nil {
			t.Fatal(err)
		}
		syms, err := st.ListSymbolsByPath(tc.path)
		if err != nil || len(syms) != 1 {
			t.Fatalf("symbols %s: %v %v", tc.path, syms, err)
		}
		fromID := syms[0].ID
		if err := st.ReplaceFileEdges(tc.path, []store.CodeEdge{
			{
				FromFileID: tf.ID, FromSymbolID: &fromID,
				ToFileID: lib.ID, ToSymbolID: &coreSym,
				Rel: store.RelValidates, Provenance: store.ImportProvenanceInferred,
			},
		}); err != nil {
			t.Fatal(err)
		}
		_ = coreFile
	}
	return []string{"src/core.ts"}, "foo.test.ts", "bar.test.ts"
}

// TestTestRunCustomAttributionFromReportedResults is work item 2's core
// regression: with a custom runner whose command fails, only targets the run
// actually reported are fail, and the command is invoked exactly once.
func TestTestRunCustomAttributionFromReportedResults(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	// The command exits 1 but the run reports one passing and one failing file:
	// the failure belongs only to the file that failed, never to every target.
	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"runner","args":["run"]}`)
	stub := &stubRunner{
		exit: 1,
		out: "\u2713 src/foo.test.ts (2 tests) 12ms\n" +
			"\u00d7 src/bar.test.ts (3 tests | 1 failed) 9ms\n",
	}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("invocations=%d want 1 (run once, not once per target)", len(stub.calls))
	}
	byName := map[string]string{}
	for _, o := range outcomes {
		byName[o.TestName] = o.TestStatus
	}
	if got := byName[fooName]; got != store.TestStatusPass {
		t.Fatalf("%s status=%q want pass (reported passing)", fooName, got)
	}
	if got := byName[barName]; got != store.TestStatusFail {
		t.Fatalf("%s status=%q want fail (reported failing)", barName, got)
	}
	// Only the two selected targets appear — a reported file with no selected
	// target must NOT be recorded under a target's name.
	if len(outcomes) != 2 {
		t.Fatalf("outcomes=%d want exactly the 2 selected targets", len(outcomes))
	}
}

// TestTestRunCustomUnreportedTargetIsSkip asserts no outcome is invented for a
// selected target the run did not report: it is recorded as skip with a
// distinct source_type, so the miss is surfaced, not dropped or faked.
func TestTestRunCustomUnreportedTargetIsSkip(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"runner","args":["run"]}`)
	// Green run that reports only foo — bar was never reported.
	stub := &stubRunner{exit: 0, out: "\u2713 src/foo.test.ts (2 tests) 8ms\n"}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("invocations=%d want 1", len(stub.calls))
	}
	var skipRow *store.OutcomeResult
	for i := range outcomes {
		if outcomes[i].TestName == barName {
			skipRow = &outcomes[i]
		}
	}
	if skipRow == nil {
		t.Fatalf("unreported target must be recorded, outcomes=%+v", outcomes)
	}
	if skipRow.TestStatus != store.TestStatusSkip {
		t.Fatalf("unreported target status=%q want skip (no invented pass)", skipRow.TestStatus)
	}
	if skipRow.SourceType != sourceTypeUnreported {
		t.Fatalf("unreported target source_type=%q want %q", skipRow.SourceType, sourceTypeUnreported)
	}
	for _, o := range outcomes {
		if o.TestName == fooName && o.TestStatus != store.TestStatusPass {
			t.Fatalf("reported target must keep its own status, got %q", o.TestStatus)
		}
	}
}

// TestTestRunCustomAlwaysFailNotAllFail is the deterministic scenario from the
// bug report: a command that runs no tests at all (exit 1, no per-file output)
// must NOT mark every selected target fail.
func TestTestRunCustomAlwaysFailNotAllFail(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"sh","args":["-c","exit 1"]}`)
	stub := &stubRunner{exit: 1, out: ""}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("invocations=%d want 1", len(stub.calls))
	}
	for _, o := range outcomes {
		if o.TestName == fooName || o.TestName == barName {
			t.Fatalf("whole-run verdict must not be mis-attributed: %s=%s", o.TestName, o.TestStatus)
		}
	}
	if len(outcomes) != 1 || !strings.HasPrefix(outcomes[0].TestName, "suite:") {
		t.Fatalf("want single suite-level outcome, got %+v", outcomes)
	}
	if outcomes[0].TestStatus != store.TestStatusFail {
		t.Fatalf("suite outcome status=%q want fail", outcomes[0].TestStatus)
	}
}

// TestTestRunCustomResultsFilePreferred proves trace/test-runner-results.json
// wins over stdout parsing for attribution.
func TestTestRunCustomResultsFilePreferred(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"runner","args":["run"]}`)
	writeConfig(t, st.ProjectRoot(), filepath.Join("trace", "test-runner-results.json"),
		`{"files":[{"name":"src/foo.test.ts","status":"fail"}]}`)
	t.Cleanup(func() {
		_ = os.Remove(filepath.Join(st.ProjectRoot(), "trace", "test-runner-results.json"))
	})
	stub := &stubRunner{exit: 1, out: "\u2713 src/foo.test.ts (2 tests) 8ms\n"}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	byName := map[string]string{}
	for _, o := range outcomes {
		byName[o.TestName] = o.TestStatus
	}
	if byName[fooName] != store.TestStatusFail {
		t.Fatalf("%s=%q want fail from results file (overrides stdout)", fooName, byName[fooName])
	}
	if byName[barName] != store.TestStatusSkip {
		t.Fatalf("unreported bar=%q want skip", byName[barName])
	}
}

// TestTestRunCustomVitestReporterFormat covers the real vitest default-reporter
// line shape including the ✓/× glyphs and failed-count grouping.
func TestTestRunCustomVitestReporterFormat(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json",
		`{"command":"pnpm","args":["exec","vitest","run","--reporter=dot"]}`)
	out := "\n Test Files  1 failed | 1 passed (2)\n\n" +
		"\u2713 src/foo.test.ts (2 tests) 12ms\n" +
		"\u00d7 src/bar.test.ts (3 tests | 1 failed) 9ms\n"
	stub := &stubRunner{exit: 1, out: out}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("invocations=%d want 1", len(stub.calls))
	}
	byName := map[string]string{}
	for _, o := range outcomes {
		byName[o.TestName] = o.TestStatus
	}
	if got := byName[fooName]; got != store.TestStatusPass {
		t.Fatalf("%s=%q want pass (vitest \u2713 line)", fooName, got)
	}
	if got := byName[barName]; got != store.TestStatusFail {
		t.Fatalf("%s=%q want fail (vitest \u00d7 line with failed count)", barName, got)
	}
}

// TestTestRunCustomNoPlaceholderRunsOnce proves the no-{path} config invokes
// the command once total (was: once per target with identical args).
func TestTestRunCustomNoPlaceholderRunsOnce(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()
	task := mustTask(t, dom)
	seedPaths, fooName, barName := seedVitestGraph(t, st)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: seedPaths[0]}},
	}); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, st.ProjectRoot(), "test-runner.json", `{"command":"runner","args":["run","all"]}`)
	stub := &stubRunner{exit: 0, out: "\u2713 src/foo.test.ts (1 test) 4ms\n\u2713 src/bar.test.ts (1 test) 4ms\n"}
	outcomes, err := RunRelevantTests(ctx, st, dom, task.ID, Options{Runner: stub})
	if err != nil {
		t.Fatalf("RunRelevantTests: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("invocations=%d want 1 total", len(stub.calls))
	}
	got := stub.calls[0]
	if len(got.Args) != 2 || got.Args[0] != "run" || got.Args[1] != "all" {
		t.Fatalf("args=%v want passthrough [run all]", got.Args)
	}
	byName := map[string]string{}
	for _, o := range outcomes {
		byName[o.TestName] = o.TestStatus
	}
	if byName[fooName] != store.TestStatusPass || byName[barName] != store.TestStatusPass {
		t.Fatalf("both reported targets must be attributed: %v", byName)
	}
}

// TestStatusForTargetSuffixMatchIsDeterministic pins the path/name resolution:
// exact keys first, then a unique, sorted path-suffix match.
func TestStatusForTargetSuffixMatchIsDeterministic(t *testing.T) {
	reported := map[string]testFileResult{
		"src/foo.test.ts": {Name: "src/foo.test.ts", Status: store.TestStatusFail},
	}
	target := TestTarget{Name: "foo.test.ts", Path: "pkg/foo_test.go"}
	status, ok, err := statusForTarget(reported, target)
	if err != nil || !ok || status != store.TestStatusFail {
		t.Fatalf("suffix match: status=%q ok=%v err=%v", status, ok, err)
	}
	if _, ok, _ := statusForTarget(map[string]testFileResult{}, target); ok {
		t.Fatal("no report must not match")
	}
	// Ambiguous suffixes (two dirs, same basename) resolve to nothing — never
	// guessed: only an exact key or a unique suffix matches.
	amb := map[string]testFileResult{
		"a/foo.test.ts": {Status: store.TestStatusPass},
		"b/foo.test.ts": {Status: store.TestStatusFail},
	}
	if _, ok, _ := statusForTarget(amb, target); ok {
		t.Fatal("ambiguous suffix must not match")
	}
	// But when only one dir reports the basename, the unique suffix matches.
	uni := map[string]testFileResult{
		"b/foo.test.ts": {Status: store.TestStatusFail},
	}
	if status, ok, _ := statusForTarget(uni, target); !ok || status != store.TestStatusFail {
		t.Fatalf("unique suffix: status=%q ok=%v", status, ok)
	}
	// Same-basename reports under both path forms are unambiguous too.
	dup := map[string]testFileResult{
		"b/foo.test.ts":   {Status: store.TestStatusPass},
		"./b/foo.test.ts": {Status: store.TestStatusPass},
	}
	if status, ok, _ := statusForTarget(dup, target); !ok || status != store.TestStatusPass {
		t.Fatalf("same-file dup: status=%q ok=%v", status, ok)
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

// TestJSTSFallbackTargetsFromSiblingTestDir pins the non-Go fallback: a
// changed TS source with no incoming validates and no impact-walk tests still
// selects the indexed test files in its sibling test/ directory (the JS/TS
// analogue of the Go package fallback). Regression: packageFallbackTargets
// skipped every non-.go path, so a TypeScript repo failed closed with no
// explanation while a Go repo was rescued by the same code path.
func TestJSTSFallbackTargetsFromSiblingTestDir(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()

	if _, err := st.UpsertFile("src/util.ts", "hutil", nil); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("src/util.ts", []store.Symbol{
		{Name: "util", Kind: "function", StartLine: 1, EndLine: 2},
	}); err != nil {
		t.Fatal(err)
	}
	testFile, err := st.UpsertFile("src/test/core-util.test.ts", "htest", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("src/test/core-util.test.ts", []store.Symbol{
		{Name: "core-util.test.ts", Kind: "test", StartLine: 1, EndLine: 3},
	}); err != nil {
		t.Fatal(err)
	}
	// A non-test file in the same dir must never become a target.
	if _, err := st.UpsertFile("src/test/helpers.ts", "hh", nil); err != nil {
		t.Fatal(err)
	}

	task := mustTask(t, dom)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "src/util.ts"}},
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
	tg := targets[0]
	if tg.Name != "core-util.test.ts" || tg.Path != "src/test/core-util.test.ts" || tg.Package != "" {
		t.Fatalf("fallback target mismatch: %+v (test file %s)", tg, testFile.Path)
	}
}

// TestSelectTestTargetsDiagnosesZeroTargets pins the truthful failure: when
// every stage yields nothing, the error names the stage and the paths that are
// missing from the index, instead of a bare "no relevant tests selected".
func TestSelectTestTargetsDiagnosesZeroTargets(t *testing.T) {
	st, dom := openTestrun(t)
	ctx := context.Background()

	// Indexed source file with no incoming validates and no impact-walk tests.
	if _, err := st.UpsertFile("src/core.ts", "hcore", nil); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileSymbols("src/core.ts", []store.Symbol{
		{Name: "core", Kind: "function", StartLine: 1, EndLine: 2},
	}); err != nil {
		t.Fatal(err)
	}

	task := mustTask(t, dom)
	if _, err := dom.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{
			{Path: "src/core.ts"},
			{Path: "src/ghost.ts"},  // source path, never indexed
			{Path: "docs/notes.md"}, // non-source path
		},
	}); err != nil {
		t.Fatal(err)
	}

	_, err := SelectTestTargets(ctx, st, dom, task.ID, nil)
	if err == nil {
		t.Fatal("expected zero-target diagnosis error")
	}
	msg := err.Error()
	for _, want := range []string{
		"no relevant tests selected",
		"1 changed path(s) not in the index",
		`"src/ghost.ts"`,
		"ignored non-source path(s)",
		`"docs/notes.md"`,
		"no incoming validates edges",
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("diagnosis missing %q: %s", want, msg)
		}
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
