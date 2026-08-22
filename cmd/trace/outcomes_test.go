package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestOutcomesCompareCLI(t *testing.T) {
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
		TaskID: task.ID, TestName: "TestShip", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestShip", TestStatus: store.TestStatusFail,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "outcomes", "compare", "--task", task.ID, "--kind", "test"})
	})
	var payload struct {
		Previous struct {
			TestStatus string `json:"test_status"`
		} `json:"previous"`
		Current struct {
			TestStatus string `json:"test_status"`
		} `json:"current"`
		Delta struct {
			TestStatus *struct {
				From string `json:"from"`
				To   string `json:"to"`
			} `json:"test_status"`
		} `json:"delta"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if payload.Previous.TestStatus != store.TestStatusPass || payload.Current.TestStatus != store.TestStatusFail {
		t.Fatalf("payload: %#v out=%q", payload, out)
	}
	if payload.Delta.TestStatus == nil || payload.Delta.TestStatus.From != "pass" || payload.Delta.TestStatus.To != "fail" {
		t.Fatalf("delta: %#v", payload.Delta)
	}
}

func TestOutcomesCompareRequiresTaskAndKind(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "outcomes", "compare"}); code != exitUsage {
		t.Fatalf("usage: %d", code)
	}
	if code := run([]string{"-C", dir, "outcomes", "compare", "--task", "x"}); code != exitUsage {
		t.Fatalf("kind required: %d", code)
	}
}

func TestOutcomesImprovementsCLI(t *testing.T) {
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
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234", Paths: []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID, TaskID: task.ID, Summary: "faster builds",
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "outcomes", "improvements", "--change", chg.ID})
	})
	var payload struct {
		OK           bool `json:"ok"`
		Count        int  `json:"count"`
		Improvements []struct {
			Summary string `json:"Summary"`
		} `json:"improvements"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v out=%q", err, out)
	}
	if !payload.OK || payload.Count != 1 || len(payload.Improvements) != 1 {
		t.Fatalf("payload: %#v out=%q", payload, out)
	}
	if payload.Improvements[0].Summary != "faster builds" {
		t.Fatalf("summary=%q", payload.Improvements[0].Summary)
	}

	outTask := captureStdout(t, func() int {
		return run([]string{"-C", dir, "outcomes", "improvements", "--task", task.ID})
	})
	var payloadTask struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal([]byte(outTask), &payloadTask); err != nil {
		t.Fatalf("decode task: %v", err)
	}
	if payloadTask.Count != 1 {
		t.Fatalf("task list count=%d", payloadTask.Count)
	}
}

func TestHelpIncludesOutcomesCompare(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(out, "outcomes compare") {
		t.Fatalf("help missing outcomes compare: %q", out)
	}
	if !strings.Contains(out, "outcomes improvements") {
		t.Fatalf("help missing outcomes improvements: %q", out)
	}
}
