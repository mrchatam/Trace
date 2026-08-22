package gitcli_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
}

func git(t *testing.T, dir string, args ...string) string {
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

func initRepo(t *testing.T) string {
	t.Helper()
	requireGit(t)
	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "trace@test.local")
	git(t, dir, "config", "user.name", "Trace Test")
	return dir
}

func writeFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

func commitAll(t *testing.T, dir, msg string) string {
	t.Helper()
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-m", msg)
	return git(t, dir, "rev-parse", "HEAD")
}

func TestOpenRejectsNonRepo(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	_, err := gitcli.Open(dir)
	if err == nil {
		t.Fatal("expected error for non-repo")
	}
	if !strings.Contains(err.Error(), "not a git") && !isNotRepo(err) {
		t.Fatalf("want ErrNotRepo, got %v", err)
	}
}

func isNotRepo(err error) bool {
	return err != nil && (err == vcs.ErrNotRepo || strings.Contains(err.Error(), "not a git"))
}

func TestShowFileMatchesGitShow(t *testing.T) {
	ctx := context.Background()
	dir := initRepo(t)
	writeFile(t, dir, "hello.txt", "hello world\n")
	oid := commitAll(t, dir, "add hello")

	repo, err := gitcli.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer repo.Close()

	got, err := repo.ShowFile(ctx, oid, "hello.txt")
	if err != nil {
		t.Fatalf("ShowFile: %v", err)
	}
	cmd := exec.Command("git", "-C", dir, "show", oid+":hello.txt")
	want, err := cmd.Output()
	if err != nil {
		t.Fatalf("git show: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("ShowFile mismatch:\n got %q\nwant %q", got, want)
	}

	head, err := repo.Head(ctx)
	if err != nil || head != oid {
		t.Fatalf("Head: got %q want %q err=%v", head, oid, err)
	}
}

func TestHistoryCommitsBetweenChangesLastChanged(t *testing.T) {
	ctx := context.Background()
	dir := initRepo(t)

	writeFile(t, dir, "a.txt", "one\n")
	c1 := commitAll(t, dir, "c1 add a")

	writeFile(t, dir, "a.txt", "two\n")
	writeFile(t, dir, "b.txt", "bee\n")
	c2 := commitAll(t, dir, "c2 modify a add b")

	writeFile(t, dir, "b.txt", "bee2\n")
	c3 := commitAll(t, dir, "c3 modify b")

	repo, err := gitcli.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer repo.Close()

	if _, err := repo.Refresh(ctx); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	hist, err := repo.History(ctx, "a.txt", 10)
	if err != nil {
		t.Fatalf("History: %v", err)
	}
	if len(hist) != 2 {
		t.Fatalf("History(a.txt) len=%d want 2: %+v", len(hist), hist)
	}
	if hist[0].OID != c2 || hist[1].OID != c1 {
		t.Fatalf("History order: %+v want %s then %s", hist, c2, c1)
	}

	last, err := repo.LastChanged(ctx, "a.txt")
	if err != nil || last.OID != c2 {
		t.Fatalf("LastChanged: %+v err=%v", last, err)
	}

	between, err := repo.CommitsBetween(ctx, c1, c3)
	if err != nil {
		t.Fatalf("CommitsBetween: %v", err)
	}
	if len(between) != 2 {
		t.Fatalf("CommitsBetween len=%d want 2: %+v", len(between), between)
	}
	if between[0].OID != c2 || between[1].OID != c3 {
		t.Fatalf("CommitsBetween order: %+v", between)
	}

	ch, err := repo.Changes(ctx, c2)
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	byPath := map[string]string{}
	for _, p := range ch {
		byPath[p.Path] = p.Status
	}
	if byPath["a.txt"] != "M" {
		t.Fatalf("a.txt status: %+v", ch)
	}
	if byPath["b.txt"] != "A" {
		t.Fatalf("b.txt status: %+v", ch)
	}

	body, err := repo.ShowFile(ctx, c3, "b.txt")
	if err != nil || string(body) != "bee2\n" {
		t.Fatalf("ShowFile b.txt: %q err=%v", body, err)
	}
}

func TestHistoryFallsBackWhenWatermarkBehindHEAD(t *testing.T) {
	ctx := context.Background()
	dir := initRepo(t)

	writeFile(t, dir, "a.txt", "one\n")
	c1 := commitAll(t, dir, "c1")
	writeFile(t, dir, "a.txt", "two\n")
	c2 := commitAll(t, dir, "c2")

	repo, err := gitcli.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer repo.Close()

	if _, err := repo.Refresh(ctx); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	writeFile(t, dir, "a.txt", "three\n")
	c3 := commitAll(t, dir, "c3")

	// Index still at c2; History must not return stale index-only results.
	hist, err := repo.History(ctx, "a.txt", 10)
	if err != nil {
		t.Fatalf("History: %v", err)
	}
	if len(hist) < 3 || hist[0].OID != c3 {
		t.Fatalf("History stale/incomplete after new commit without Refresh: %+v (want newest %s, also %s %s)", hist, c3, c2, c1)
	}

	last, err := repo.LastChanged(ctx, "a.txt")
	if err != nil || last.OID != c3 {
		t.Fatalf("LastChanged: %+v err=%v want %s", last, err, c3)
	}
}

func TestIncrementalRefresh(t *testing.T) {
	ctx := context.Background()
	dir := initRepo(t)

	writeFile(t, dir, "f.txt", "v1\n")
	commitAll(t, dir, "n1")
	writeFile(t, dir, "f.txt", "v2\n")
	commitAll(t, dir, "n2")
	writeFile(t, dir, "f.txt", "v3\n")
	commitAll(t, dir, "n3")

	st, err := store.Open(dir)
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	defer st.Close()

	repo, err := gitcli.OpenWithStore(dir, st)
	if err != nil {
		t.Fatalf("OpenWithStore: %v", err)
	}
	defer repo.Close()

	r1, err := repo.Refresh(ctx)
	if err != nil {
		t.Fatalf("Refresh1: %v", err)
	}
	if r1.NewCommits != 3 {
		t.Fatalf("Refresh1 NewCommits=%d want 3", r1.NewCommits)
	}
	n, err := st.CountIndexedCommits()
	if err != nil || n != 3 {
		t.Fatalf("count after r1: %d err=%v", n, err)
	}
	wm1, _ := st.GetMeta(store.MetaVCSWatermark)

	// Noop refresh: no new commits.
	r2, err := repo.Refresh(ctx)
	if err != nil {
		t.Fatalf("Refresh2: %v", err)
	}
	if r2.NewCommits != 0 {
		t.Fatalf("Refresh2 NewCommits=%d want 0", r2.NewCommits)
	}
	n, _ = st.CountIndexedCommits()
	if n != 3 {
		t.Fatalf("count after noop rewrite? got %d", n)
	}
	wm2, _ := st.GetMeta(store.MetaVCSWatermark)
	if wm2 != wm1 {
		t.Fatalf("watermark changed on noop: %q -> %q", wm1, wm2)
	}

	// Add K=2 new commits.
	writeFile(t, dir, "f.txt", "v4\n")
	commitAll(t, dir, "n4")
	writeFile(t, dir, "g.txt", "g\n")
	commitAll(t, dir, "n5")

	r3, err := repo.Refresh(ctx)
	if err != nil {
		t.Fatalf("Refresh3: %v", err)
	}
	if r3.NewCommits != 2 {
		t.Fatalf("Refresh3 NewCommits=%d want 2", r3.NewCommits)
	}
	n, _ = st.CountIndexedCommits()
	if n != 5 {
		t.Fatalf("count after r3: %d want 5", n)
	}

	r4, err := repo.Refresh(ctx)
	if err != nil {
		t.Fatalf("Refresh4: %v", err)
	}
	if r4.NewCommits != 0 {
		t.Fatalf("Refresh4 NewCommits=%d want 0", r4.NewCommits)
	}
	n, _ = st.CountIndexedCommits()
	if n != 5 {
		t.Fatalf("count after final noop: %d want 5 (no wipe)", n)
	}
}

func TestNoBlobColumnsInVCSIndex(t *testing.T) {
	dir := initRepo(t)
	writeFile(t, dir, "x.txt", "x\n")
	commitAll(t, dir, "x")

	st, err := store.Open(dir)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer st.Close()

	repo, err := gitcli.OpenWithStore(dir, st)
	if err != nil {
		t.Fatalf("OpenWithStore: %v", err)
	}
	defer repo.Close()
	if _, err := repo.Refresh(context.Background()); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	bad, where, err := st.HasBlobLikeColumns()
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if bad {
		t.Fatalf("forbidden column in schema: %s", where)
	}

	// Ensure indexed rows have no patch-like payload columns populated via subject only.
	c, err := st.GetIndexedCommit(git(t, dir, "rev-parse", "HEAD"))
	if err != nil {
		t.Fatalf("GetIndexedCommit: %v", err)
	}
	if strings.Contains(c.Subject, "diff --git") {
		t.Fatalf("subject looks like a patch: %q", c.Subject)
	}
}
