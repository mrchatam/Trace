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
	"github.com/mrchatam/Trace/internal/testrun"
)

func TestChangesCapturePromotesMeaningfulCommit(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	writeGoFile(t, dir, "main.go", "package main\nfunc main() {}\n")
	git("add", "-A")
	git("commit", "-m", "initial go")
	head := git("rev-parse", "HEAD")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "index", "main.go"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "capture"})
	})
	var rows []changeCaptureRow
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if len(rows) != 1 {
		t.Fatalf("want one promoted change, got %d: %q", len(rows), out)
	}
	if rows[0].GitCommit != head || rows[0].TaskID != domain.VCSCaptureTaskID {
		t.Fatalf("row: %+v", rows[0])
	}
	if rows[0].Status != store.ChangeStatusRecorded {
		t.Fatalf("status: %+v", rows[0])
	}

	out2 := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "capture"})
	})
	var rows2 []changeCaptureRow
	if err := json.Unmarshal([]byte(out2), &rows2); err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 || rows2[0].ID != rows[0].ID {
		t.Fatalf("second capture idempotent: first=%+v second=%q", rows[0], out2)
	}
}

// TestChangesRecordCLIThenTestSelection proves the CLI-drivable loop:
// `changes record` persists a task-scoped change (OPEN without --commit,
// RECORDED with it), `changes list --task` sees it, and its paths seed
// `test run` selection (validating test selected without --paths).
func TestChangesRecordCLIThenTestSelection(t *testing.T) {
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
	_ = st.Close()

	// OPEN without --commit.
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "record", "--task", task.ID,
			"--reason", "first cut", "pkg/foo.go"})
	})
	var rec struct {
		OK     bool   `json:"ok"`
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatalf("json: %v out=%q", err, out)
	}
	if !rec.OK || rec.ID == "" || rec.Status != "OPEN" {
		t.Fatalf("record: %+v out=%q", rec, out)
	}

	// Visible via changes list --task.
	listOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "list", "--task", task.ID})
	})
	if !strings.Contains(listOut, rec.ID) {
		t.Fatalf("changes list missing recorded change: %q", listOut)
	}

	// OPEN changes deliberately do NOT seed selection (provenance honesty:
	// only RECORDED/COMPARED evidence does). Record with a commit SHA.
	out2 := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "record", "--task", task.ID,
			"--commit", "abc1234", "--reason", "committed cut", "pkg/foo.go"})
	})
	var rec2 struct {
		OK     bool   `json:"ok"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(out2), &rec2); err != nil {
		t.Fatalf("json: %v out=%q", err, out2)
	}
	if !rec2.OK || rec2.Status != "RECORDED" {
		t.Fatalf("record with commit: %+v out=%q", rec2, out2)
	}

	// Selection derives from the recorded change paths (library seam used by the
	// CLI): the validating test is selected without explicit --paths.
	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	targets, err := testrun.SelectTestTargets(context.Background(), st2, domain.New(st2), task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Name != "TestFoo" {
		t.Fatalf("targets=%+v want TestFoo seeded by recorded change paths", targets)
	}
}

func TestChangesCaptureSkipsDocsOnlyUnlessAll(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "-A")
	git("commit", "-m", "docs only")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "index"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "capture"})
	})
	var rows []changeCaptureRow
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("docs-only should skip: %q", out)
	}

	outAll := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "capture", "--all"})
	})
	var rowsAll []changeCaptureRow
	if err := json.Unmarshal([]byte(outAll), &rowsAll); err != nil {
		t.Fatal(err)
	}
	if len(rowsAll) != 1 {
		t.Fatalf("docs-only with --all: %q", outAll)
	}
}

func TestChangesCompare(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	writeGoFile(t, dir, "main.go", "package main\nfunc A() {}\n")
	git("add", "-A")
	git("commit", "-m", "first")
	c1 := git("rev-parse", "HEAD")

	writeGoFile(t, dir, "main.go", "package main\nfunc B() {}\n")
	writeGoFile(t, dir, "extra.go", "package main\n")
	git("add", "-A")
	git("commit", "-m", "second")
	c2 := git("rev-parse", "HEAD")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "compare", "--from", c1, "--to", c2})
	})
	var result struct {
		From     string   `json:"from"`
		To       string   `json:"to"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if result.From != c1 || result.To != c2 {
		t.Fatalf("range: %+v", result)
	}
	if len(result.Added) != 1 || result.Added[0] != "extra.go" {
		t.Fatalf("added: %+v", result.Added)
	}
	if len(result.Modified) != 1 || result.Modified[0] != "main.go" {
		t.Fatalf("modified: %+v", result.Modified)
	}
	if len(result.Removed) != 0 {
		t.Fatalf("removed: %+v", result.Removed)
	}
}

func TestCLIChangesList(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	task := mustChangeTaskCLI(t, svc)

	older, err := st.UpsertChange(store.Change{
		TaskID:     task.ID,
		Reason:     "older change",
		Status:     store.ChangeStatusRecorded,
		SourceType: domain.DefaultSourceType,
		CreatedAt:  "2026-01-01T00:00:00Z",
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := st.InsertChangePath(store.ChangePath{ChangeID: older.ID, Path: "a.go", Status: "M"}); err != nil {
		st.Close()
		t.Fatal(err)
	}

	newer, err := st.UpsertChange(store.Change{
		TaskID:     task.ID,
		Reason:     "newer change",
		Status:     store.ChangeStatusRecorded,
		SourceType: domain.DefaultSourceType,
		CreatedAt:  "2026-02-01T00:00:00Z",
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := st.InsertChangePath(store.ChangePath{ChangeID: newer.ID, Path: "b.go", Status: "A"}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "list", "--task", task.ID, "--limit", "10"})
	})
	var rows []changeListRow
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &rows); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if len(rows) != 2 {
		t.Fatalf("want 2 rows, got %d: %q", len(rows), out)
	}
	if rows[0].ID != newer.ID || rows[1].ID != older.ID {
		t.Fatalf("newest-first order: %+v", rows)
	}
	for _, r := range rows {
		if r.TaskID != task.ID {
			t.Fatalf("task_id: %+v", r)
		}
	}

	emptyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "list", "--task", "00000000-0000-0000-0000-000000000001"})
	})
	var emptyRows []changeListRow
	if err := json.Unmarshal([]byte(strings.TrimSpace(emptyOut)), &emptyRows); err != nil {
		t.Fatal(err)
	}
	if len(emptyRows) != 0 {
		t.Fatalf("unknown task should return [], got %q", emptyOut)
	}
}

func TestCLIChangesShow(t *testing.T) {
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
	task := mustChangeTaskCLI(t, svc)
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "show paths",
		Paths: []domain.ChangePathInput{
			{Path: "internal/a.go", Status: "M", SymbolID: "sym-1"},
			{Path: "internal/b.go", Status: "A"},
		},
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "show", chg.ID})
	})
	var resp changeShowResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &resp); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if !resp.OK || resp.Change.ID != chg.ID {
		t.Fatalf("change: %+v", resp)
	}
	if len(resp.Paths) != 2 {
		t.Fatalf("paths: %+v", resp.Paths)
	}
	for _, p := range resp.Paths {
		if p.Path == "" {
			t.Fatalf("path required: %+v", p)
		}
	}
	body := strings.ToLower(out)
	if strings.Contains(body, "package ") || strings.Contains(body, "func ") {
		t.Fatalf("show must not include file content: %q", out)
	}

	if code := run([]string{"-C", dir, "changes", "show", "00000000-0000-0000-0000-000000000001"}); code == exitOK {
		t.Fatal("missing change id should fail")
	}
}

func mustChangeTaskCLI(t *testing.T, svc *domain.Service) store.Task {
	t.Helper()
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "changes CLI task"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	return task
}

func TestChangesSimilarCLI(t *testing.T) {
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
	task := mustChangeTaskCLI(t, svc)
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "internal one",
		Paths:  []domain.ChangePathInput{{Path: "internal/a.go", Status: "M"}},
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "cmd other",
		Paths:  []domain.ChangePathInput{{Path: "cmd/main.go", Status: "M"}},
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "similar", "--kind", "seg:internal"})
	})
	var result domain.SimilarChangesResult
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &result); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if len(result.Changes) != 1 || result.Changes[0].ChangeKind != "seg:internal" {
		t.Fatalf("similar by kind: %+v", result)
	}

	outPath := captureStdout(t, func() int {
		return run([]string{"-C", dir, "changes", "similar", "--path", "internal/"})
	})
	var byPath domain.SimilarChangesResult
	if err := json.Unmarshal([]byte(strings.TrimSpace(outPath)), &byPath); err != nil {
		t.Fatalf("json path: %v body=%q", err, outPath)
	}
	if len(byPath.Changes) != 1 {
		t.Fatalf("similar by path: %+v", byPath)
	}

	if code := run([]string{"-C", dir, "changes", "similar"}); code == exitOK {
		t.Fatal("similar without path/kind should fail")
	}
}
