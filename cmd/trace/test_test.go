package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestTestRunCLI(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/cli\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
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
	prod, _ := st.ListSymbolsByPath("pkg/foo.go")
	tests, _ := st.ListSymbolsByPath("pkg/foo_test.go")
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

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "pkg/foo.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	pkgDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "foo.go"), []byte("package pkg\nfunc Foo() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "foo_test.go"), []byte("package pkg\nimport \"testing\"\nfunc TestFoo(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "test", "run", "--task", task.ID})
	})
	var body struct {
		OK       bool `json:"ok"`
		Count    int  `json:"count"`
		Outcomes []struct {
			TestName   string `json:"test_name"`
			TestStatus string `json:"test_status"`
		} `json:"outcomes"`
	}
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if !body.OK || body.Count < 1 {
		t.Fatalf("body: %+v out=%q", body, out)
	}
}

func TestTestRunCLIRequiresTask(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "test", "run"}); code != exitUsage {
		t.Fatalf("usage: %d", code)
	}
}

func TestHelpIncludesTestRun(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(out, "test run") || !strings.Contains(out, "--task") {
		t.Fatalf("help missing test run: %q", out)
	}
}
