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

func TestEvalRulesCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "eval", "rules"})
	})
	var payload struct {
		Loaded     bool     `json:"loaded"`
		Path       string   `json:"path"`
		Mechanisms []string `json:"mechanisms"`
		Invariants []struct {
			ID      string `json:"id"`
			Enabled bool   `json:"enabled"`
		} `json:"invariants"`
		CachedAt string `json:"cached_at"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if payload.Loaded {
		t.Fatal("expected loaded=false without rules file")
	}
	if payload.Path != "trace/eval-rules.json" {
		t.Fatalf("path: %q", payload.Path)
	}
	if len(payload.Mechanisms) != 4 {
		t.Fatalf("mechanisms: %v", payload.Mechanisms)
	}
	if payload.CachedAt != "" {
		t.Fatalf("cached_at should be empty without file load: %q", payload.CachedAt)
	}

	fixture, err := os.ReadFile(filepath.Join("..", "..", "internal", "eval", "testdata", "eval-rules-default.json"))
	if err != nil {
		t.Fatal(err)
	}
	traceDir := filepath.Join(dir, "trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(traceDir, "eval-rules.json"), fixture, 0o644); err != nil {
		t.Fatal(err)
	}

	outLoaded := captureStdout(t, func() int {
		return run([]string{"-C", dir, "eval", "rules"})
	})
	var loaded struct {
		Loaded   bool   `json:"loaded"`
		CachedAt string `json:"cached_at"`
	}
	if err := json.Unmarshal([]byte(outLoaded), &loaded); err != nil {
		t.Fatalf("decode loaded: %v out=%q", err, outLoaded)
	}
	if !loaded.Loaded {
		t.Fatalf("expected loaded=true: %q", outLoaded)
	}
	if loaded.CachedAt == "" {
		t.Fatalf("expected cached_at when file loaded: %q", outLoaded)
	}
}

func TestEvalRulesUsage(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "eval"}); code != exitUsage {
		t.Fatalf("eval usage: %d", code)
	}
	if code := run([]string{"-C", dir, "eval", "rules", "--bogus"}); code != exitUsage {
		t.Fatalf("eval rules extra args: %d", code)
	}
}

func TestHelpIncludesEvalRules(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "eval rules") {
		t.Fatalf("help missing eval rules: %q", out)
	}
	if !strings.Contains(out, "eval results") {
		t.Fatalf("help missing eval results: %q", out)
	}
}

func TestEvalResultsCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	ctx := context.Background()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "cli eval"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "cli task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abcdef0",
		ScoresJSON: `{"correctness": 1.0}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     task.ID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.99}`,
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	taskID := task.ID
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "eval", "results", "--task", taskID})
	})
	var rows []struct {
		ID          string `json:"id"`
		MechanismID string `json:"mechanism_id"`
		Kind        string `json:"kind"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if len(rows) != 1 {
		t.Fatalf("rows: got %d want 1: %q", len(rows), out)
	}
	if rows[0].MechanismID != "stored_evaluation" {
		t.Fatalf("mechanism_id: %q", rows[0].MechanismID)
	}
	if rows[0].Kind != "evaluation" {
		t.Fatalf("kind: %q", rows[0].Kind)
	}
}

func TestEvalResultsUsage(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "eval", "results"}); code != exitUsage {
		t.Fatalf("eval results usage: %d", code)
	}
	if code := run([]string{"-C", dir, "eval", "results", "--task", ""}); code != exitUsage {
		t.Fatalf("eval results empty task: %d", code)
	}
}
