package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestPatternsRefreshAndList(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "patterns CLI"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: task.ID,
		Reason: "internal tweak",
		Paths:  []domain.ChangePathInput{{Path: "internal/x.go", Status: "M"}},
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.RecordImprovement(context.Background(), domain.ImprovementInput{
		ChangeID: chg.ID,
		TaskID:   task.ID,
		Summary:  "better",
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	outRefresh := captureStdout(t, func() int {
		return run([]string{"-C", dir, "patterns", "refresh"})
	})
	var refreshResp patternsRefreshResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(outRefresh)), &refreshResp); err != nil {
		t.Fatalf("refresh json: %v body=%q", err, outRefresh)
	}
	if !refreshResp.OK || refreshResp.PatternsUpdated < 1 {
		t.Fatalf("refresh: %+v", refreshResp)
	}

	outList := captureStdout(t, func() int {
		return run([]string{"-C", dir, "patterns", "list"})
	})
	var rows []store.ChangePattern
	if err := json.Unmarshal([]byte(strings.TrimSpace(outList)), &rows); err != nil {
		t.Fatalf("list json: %v body=%q", err, outList)
	}
	if len(rows) == 0 {
		t.Fatalf("patterns list empty: %q", outList)
	}
}

func TestPatternsListEmptyBeforeRefresh(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "patterns", "list"})
	})
	var rows []store.ChangePattern
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &rows); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if len(rows) != 0 {
		t.Fatalf("want [] before refresh, got %q", out)
	}
}
