package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/store"
)

func gitTestHelper(t *testing.T, dir string) func(args ...string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
	return func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Trace Test",
			"GIT_AUTHOR_EMAIL=trace@test.local",
			"GIT_COMMITTER_NAME=Trace Test",
			"GIT_COMMITTER_EMAIL=trace@test.local",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
}

func writeGoFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestGraphSyncStaleWhenHeadDiffers proves HEAD advanced past watermark → status stale + packet notice.
func TestGraphSyncStaleWhenHeadDiffers(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	writeGoFile(t, dir, "main.go", "package main\nfunc A() {}\n")
	git("add", "-A")
	git("commit", "-m", "first")
	head1 := git("rev-parse", "HEAD")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "index", "main.go"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	writeGoFile(t, dir, "main.go", "package main\nfunc B() {}\n")
	git("add", "-A")
	git("commit", "-m", "second")
	head2 := git("rev-parse", "HEAD")
	if head1 == head2 {
		t.Fatal("expected HEAD to advance")
	}

	statusJSON := captureStdout(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	var status indexStatusJSON
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, statusJSON)
	}
	if status.Head != head2 {
		t.Fatalf("head=%q want %q", status.Head, head2)
	}
	if status.LastIndexedCommit != head1 {
		t.Fatalf("last_indexed_commit=%q want %q", status.LastIndexedCommit, head1)
	}
	if !status.Stale {
		t.Fatalf("stale=false want true: %+v", status)
	}

	goalID := "11111111-1111-1111-1111-111111111111"
	taskID := "22222222-2222-2222-2222-222222222222"
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{"version":1,"goals":[{"id":"` + goalID + `","title":"G","body":""}],"tasks":[{"id":"` + taskID + `","title":"T","body":"","goal_id":"` + goalID + `"}]}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed import: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	c := compiler.New(st)
	pkt, err := c.TaskContext(t.Context(), taskID, compiler.ContextOptions{MaxItems: 8})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	if pkt.GraphSyncHonesty == nil || !pkt.GraphSyncHonesty.StaleCommit {
		t.Fatalf("expected graph_sync_honesty.stale_commit; got %+v", pkt.GraphSyncHonesty)
	}
	if pkt.GraphSyncHonesty.Head != head2 {
		t.Fatalf("packet head=%q want %q", pkt.GraphSyncHonesty.Head, head2)
	}
	if pkt.GraphSyncHonesty.LastIndexedCommit != head1 {
		t.Fatalf("packet last_indexed=%q want %q", pkt.GraphSyncHonesty.LastIndexedCommit, head1)
	}
}

// TestHookIndexUpdatesLastIndexedCommit proves cmdIndex sets watermark to HEAD.
func TestHookIndexUpdatesLastIndexedCommit(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	writeGoFile(t, dir, "pkg/hook.go", "package pkg\nfunc Hook() {}\n")
	git("add", "-A")
	git("commit", "-m", "init")
	head := git("rev-parse", "HEAD")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "index", "pkg/hook.go"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	statusJSON := captureStdout(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	var status indexStatusJSON
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, statusJSON)
	}
	if status.LastIndexedCommit != head {
		t.Fatalf("last_indexed_commit=%q want HEAD %q", status.LastIndexedCommit, head)
	}
	if status.Stale {
		t.Fatalf("stale=true want false after index: %+v", status)
	}
}
