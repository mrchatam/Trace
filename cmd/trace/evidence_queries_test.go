package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func seedTestsVerifyingCLI(t *testing.T, dir string) string {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	root := st.ProjectRoot()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		st.Close()
		t.Fatal(err)
	}
	fooGo, err := os.ReadFile(filepath.Join("..", "..", "internal", "analyzers", "testdata", "foo.go"))
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	fooTest, err := os.ReadFile(filepath.Join("..", "..", "internal", "analyzers", "testdata", "foo_test.go"))
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if err := analyzers.IndexFile(ctx, st, "pkg/foo.go", fooGo, analyzers.IndexOptions{}); err != nil {
		st.Close()
		t.Fatalf("index foo.go: %v", err)
	}
	if err := analyzers.IndexFile(ctx, st, "pkg/foo_test.go", fooTest, analyzers.IndexOptions{}); err != nil {
		st.Close()
		t.Fatalf("index foo_test.go: %v", err)
	}
	prod, err := st.ListSymbolsByPath("pkg/foo.go")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	var fooID string
	for _, s := range prod {
		if s.Name == "Foo" {
			fooID = s.ID
			break
		}
	}
	if fooID == "" {
		st.Close()
		t.Fatal("Foo symbol missing")
	}
	st.Close()
	return fooID
}

func TestTestsVerifyingCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	fooID := seedTestsVerifyingCLI(t, dir)

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tests", "verifying", "--symbol", fooID})
	})
	var payload struct {
		OK    bool `json:"ok"`
		Count int  `json:"count"`
		Tests []struct {
			TestSymbolName string `json:"test_symbol_name"`
		} `json:"tests"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if !payload.OK || payload.Count != 1 || len(payload.Tests) != 1 || payload.Tests[0].TestSymbolName != "TestFoo" {
		t.Fatalf("payload: %#v out=%q", payload, out)
	}

	if code := run([]string{"-C", dir, "tests", "verifying"}); code != exitUsage {
		t.Fatalf("xor usage: %d", code)
	}
	if code := run([]string{"-C", dir, "tests", "verifying", "--file", "missing.go"}); code != exitFail {
		t.Fatalf("missing file: %d", code)
	}
}

func TestOutcomesFailedCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestFail", TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "outcomes", "failed", "--task", task.ID})
	})
	var payload struct {
		Count  int `json:"count"`
		Failed []struct {
			TestStatus string `json:"test_status"`
		} `json:"failed"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Count != 1 || payload.Failed[0].TestStatus != store.TestStatusFail {
		t.Fatalf("payload: %#v", payload)
	}
}

func TestOutcomesWorkedCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestPass", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "outcomes", "worked", "--task", task.ID})
	})
	var payload struct {
		Count  int `json:"count"`
		Worked []struct {
			Kind string `json:"kind"`
		} `json:"worked"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Count < 1 {
		t.Fatalf("payload: %#v", payload)
	}
	foundPass := false
	for _, w := range payload.Worked {
		if w.Kind == "test_pass" {
			foundPass = true
		}
	}
	if !foundPass {
		t.Fatalf("worked: %#v", payload.Worked)
	}
}

func TestRegressionsListCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	outEval := mustEvalWithRegressionDomain(t, svc, task.ID)
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: outEval.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "regressions", "list", "--task", task.ID})
	})
	var payload struct {
		Count       int `json:"count"`
		Regressions []struct {
			ID string `json:"id"`
		} `json:"regressions"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if payload.Count < 1 {
		t.Fatalf("payload: %#v", payload)
	}
	found := false
	for _, r := range payload.Regressions {
		if r.ID == reg.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("regression %s not in list: %+v", reg.ID, payload.Regressions)
	}
}

func TestHelpIncludesEvidenceQueries(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	for _, needle := range []string{
		"tests verifying", "outcomes failed", "outcomes worked", "regressions list",
	} {
		if !strings.Contains(out, needle) {
			t.Fatalf("help missing %q: %q", needle, out)
		}
	}
}

// mustEvalWithRegressionDomain duplicates regressions_test helper for CLI package tests.
func mustEvalWithRegressionDomain(t *testing.T, svc *domain.Service, taskID string) store.OutcomeResult {
	t.Helper()
	ctx := context.Background()
	bl, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"latency_ms": 10}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: bl.ID, ScoresJSON: `{"latency_ms": 50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
