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

func TestVerifyInvariantsCLI(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	if code := run([]string{"-C", dir, "verify", "invariants", "--task", "00000000-0000-4000-8000-000000000001"}); code != exitFail {
		t.Fatalf("unknown task exit: %d want %d", code, exitFail)
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

	testdata := filepath.Join("..", "..", "internal", "analyzers", "testdata")
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
		TaskID:    task.ID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "internal/lib/lib.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "verify", "invariants", "--task", task.ID})
	})
	var payload struct {
		Passed     bool `json:"passed"`
		Violations []struct {
			FromLayer    string `json:"from_layer"`
			ToLayer      string `json:"to_layer"`
			ImporterPath string `json:"importer_path"`
		} `json:"violations"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if payload.Passed || len(payload.Violations) == 0 {
		t.Fatalf("want fail with violations, got %#v out=%q", payload, out)
	}
	if payload.Violations[0].FromLayer != "internal" || payload.Violations[0].ToLayer != "cmd" {
		t.Fatalf("layers: %+v", payload.Violations[0])
	}
}

func TestVerifyInvariantsRequiresTask(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "verify", "invariants"}); code != exitUsage {
		t.Fatalf("usage: %d", code)
	}
}

func TestHelpIncludesVerifyInvariants(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "verify invariants") {
		t.Fatalf("help missing verify invariants: %q", out)
	}
}
