package domain_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
)

func requireGitCompare(t *testing.T) func(dir string, args ...string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
	return func(dir string, args ...string) string {
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

func writeCompareFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitAllCompare(t *testing.T, git func(string, ...string) string, dir, msg string) string {
	t.Helper()
	git(dir, "add", "-A")
	git(dir, "commit", "-m", msg)
	return git(dir, "rev-parse", "HEAD")
}

func openDomainGit(t *testing.T) (*domain.Service, *store.Store, func(string, ...string) string, string) {
	t.Helper()
	dir := t.TempDir()
	git := requireGitCompare(t)
	git(dir, "init")
	git(dir, "config", "user.email", "trace@test.local")
	git(dir, "config", "user.name", "Trace Test")

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st), st, git, dir
}

func TestCompareStatesPathDeltaNoBlob(t *testing.T) {
	svc, st, git, dir := openDomainGit(t)
	ctx := context.Background()

	writeCompareFile(t, dir, "keep.go", "package main\n")
	writeCompareFile(t, dir, "remove.go", "package main\nvar removed = 1\n")
	writeCompareFile(t, dir, "old.go", "package main\n// v1\n")
	c1 := commitAllCompare(t, git, dir, "baseline")

	writeCompareFile(t, dir, "add.go", "package extra\nconst added = true\n")
	os.Remove(filepath.Join(dir, "remove.go"))
	writeCompareFile(t, dir, "old.go", "package main\n// v2\n")
	c2 := commitAllCompare(t, git, dir, "delta")

	repo, err := gitcli.OpenWithStore(dir, st)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
	_ = repo.Close()

	got, err := svc.CompareStates(ctx, c1, c2)
	if err != nil {
		t.Fatalf("CompareStates: %v", err)
	}
	if got.From != c1 || got.To != c2 {
		t.Fatalf("range: %+v", got)
	}
	if !sliceEqual(got.Added, []string{"add.go"}) {
		t.Fatalf("added: %v", got.Added)
	}
	if !sliceEqual(got.Removed, []string{"remove.go"}) {
		t.Fatalf("removed: %v", got.Removed)
	}
	if !sliceEqual(got.Modified, []string{"old.go"}) {
		t.Fatalf("modified: %v", got.Modified)
	}
	if len(got.ChangeIDs) != 0 {
		t.Fatalf("change_ids: %v", got.ChangeIDs)
	}

	bad, where, err := st.HasBlobLikeColumns()
	if err != nil || bad {
		t.Fatalf("HasBlobLikeColumns: bad=%v where=%s err=%v", bad, where, err)
	}
}

func TestCompareStatesLinksChangeWhenPresent(t *testing.T) {
	svc, st, git, dir := openDomainGit(t)
	ctx := context.Background()

	writeCompareFile(t, dir, "pkg/main.go", "package pkg\nfunc A() {}\n")
	c1 := commitAllCompare(t, git, dir, "first")

	writeCompareFile(t, dir, "pkg/main.go", "package pkg\nfunc B() {}\n")
	c2 := commitAllCompare(t, git, dir, "second")

	repo, err := gitcli.OpenWithStore(dir, st)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
	_ = repo.Close()

	change, err := svc.PromoteVCSCommitToChange(ctx, c2)
	if err != nil {
		t.Fatal(err)
	}
	if change.ID == "" {
		t.Fatal("expected promoted change")
	}

	got, err := svc.CompareStates(ctx, c1, c2)
	if err != nil {
		t.Fatalf("CompareStates: %v", err)
	}
	if len(got.ChangeIDs) != 1 || got.ChangeIDs[0] != change.ID {
		t.Fatalf("change_ids: %+v want [%s]", got.ChangeIDs, change.ID)
	}
}

func TestCompareStatesUnknownOIDFailClosed(t *testing.T) {
	svc, _, git, dir := openDomainGit(t)
	ctx := context.Background()

	writeCompareFile(t, dir, "a.go", "package main\n")
	good := commitAllCompare(t, git, dir, "one")

	_, err := svc.CompareStates(ctx, good, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if err == nil {
		t.Fatal("expected error for unknown to OID")
	}

	_, err = svc.CompareStates(ctx, "badoid", good)
	if err == nil {
		t.Fatal("expected error for malformed from OID")
	}

	_, err = svc.CompareStates(ctx, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", good)
	if err == nil {
		t.Fatal("expected error for unknown from OID")
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
