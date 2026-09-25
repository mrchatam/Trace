package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestGraphInferCLI_DryRunAndInsert(t *testing.T) {
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
	if _, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	}); err != nil {
		t.Fatal(err)
	}
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateTask(ctx, domain.TaskInput{Title: "login auth", GoalID: &g.ID}); err != nil {
		t.Fatal(err)
	}
	st.Close()

	dryOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "graph", "infer", "--dry-run"})
	})
	var dry map[string]any
	if err := json.Unmarshal([]byte(dryOut), &dry); err != nil {
		t.Fatalf("dry JSON: %v\n%s", err, dryOut)
	}
	if dry["dry_run"] != true {
		t.Fatalf("dry_run=%v", dry["dry_run"])
	}
	if dry["inserted"] != float64(0) {
		t.Fatalf("dry inserted=%v", dry["inserted"])
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	links, err := st2.ListLinksByRel(domain.RelScopeMember)
	st2.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 0 {
		t.Fatalf("dry-run wrote links: %+v", links)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "graph", "infer"})
	})
	var rep map[string]any
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("infer JSON: %v\n%s", err, out)
	}
	if rep["inserted"] != float64(1) {
		t.Fatalf("inserted=%v want 1; %s", rep["inserted"], out)
	}

	again := captureStdout(t, func() int {
		return run([]string{"-C", dir, "graph", "infer"})
	})
	var rep2 map[string]any
	if err := json.Unmarshal([]byte(again), &rep2); err != nil {
		t.Fatal(err)
	}
	if rep2["inserted"] != float64(0) {
		t.Fatalf("idempotent inserted=%v", rep2["inserted"])
	}
}

func TestGraphHelpMentionsOptIn(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"graph", "help"})
	})
	if !strings.Contains(out, "infer") || !strings.Contains(out, "does not run on trace index") {
		t.Fatalf("graph help missing opt-in note:\n%s", out)
	}
}
