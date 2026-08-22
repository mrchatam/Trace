package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyRunCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	if code := run([]string{"-C", dir, "verify", "run", "--task", "00000000-0000-4000-8000-000000000001"}); code != exitFail {
		t.Fatalf("unknown task exit: %d want %d", code, exitFail)
	}

	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "verify", "run", "--task", taskID})
	})
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if payload["ok"] != true || payload["task"] != taskID {
		t.Fatalf("payload: %#v", payload)
	}
}

func TestHelpIncludesVerifyRun(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "verify run") {
		t.Fatalf("help missing verify run: %q", out)
	}
}

func TestVerifyRunRequiresTask(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "verify", "run"}); code != exitUsage {
		t.Fatalf("usage: %d", code)
	}
}
