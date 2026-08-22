package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestTasksConflictsCLI(t *testing.T) {
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

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "cli goal"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "cli task A", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "cli task B", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{taskA.ID, taskB.ID} {
		if err := svc.TransitionTask(ctx, id, store.WorkStateInProgress, domain.TransitionOptions{
			Actor: "test", Reason: "start",
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskA.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/shared/pkg"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskB.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/shared/pkg/a/b.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks", "conflicts"})
	})
	var resp map[string]any
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("json: %v (%s)", err, out)
	}
	if resp["ok"] != true {
		t.Fatalf("ok=%v", resp["ok"])
	}
	conflicts, ok := resp["conflicts"].([]any)
	if !ok || len(conflicts) != 1 {
		t.Fatalf("conflicts=%#v", resp["conflicts"])
	}
	row, ok := conflicts[0].(map[string]any)
	if !ok {
		t.Fatalf("row type: %#v", conflicts[0])
	}
	if !strings.Contains(row["reason"].(string), "path_overlap") {
		t.Fatalf("reason=%v", row["reason"])
	}

	filtered := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks", "conflicts", "--task", taskA.ID})
	})
	var filteredResp map[string]any
	if err := json.Unmarshal([]byte(filtered), &filteredResp); err != nil {
		t.Fatalf("filtered json: %v (%s)", err, filtered)
	}
	frows, ok := filteredResp["conflicts"].([]any)
	if !ok || len(frows) != 1 {
		t.Fatalf("filtered conflicts=%#v", filteredResp["conflicts"])
	}
	frow := frows[0].(map[string]any)
	if frow["task_a"] != taskA.ID && frow["task_b"] != taskA.ID {
		t.Fatalf("filter miss: %+v", frow)
	}

	emptyDir := t.TempDir()
	if code := run([]string{"-C", emptyDir, "init"}); code != exitOK {
		t.Fatal(code)
	}
	emptyOut := captureStdout(t, func() int {
		return run([]string{"-C", emptyDir, "tasks", "conflicts"})
	})
	var emptyResp map[string]any
	if err := json.Unmarshal([]byte(emptyOut), &emptyResp); err != nil {
		t.Fatal(err)
	}
	emptyConflicts, ok := emptyResp["conflicts"].([]any)
	if !ok || len(emptyConflicts) != 0 {
		t.Fatalf("empty conflicts=%#v", emptyResp["conflicts"])
	}
}

func TestHelpIncludesTasksConflicts(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "tasks conflicts") {
		t.Fatalf("help missing tasks conflicts: %s", out)
	}
}
