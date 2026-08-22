package install_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/install"
)

func initGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "trace@test.local")
	git(t, dir, "config", "user.name", "Trace Test")
	return dir
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func gitHookTarget(t *testing.T) install.Target {
	t.Helper()
	tgt, err := install.Lookup(install.TargetGitHook)
	if err != nil {
		t.Fatal(err)
	}
	return tgt
}

func TestInstallDetectIncludesGitHook(t *testing.T) {
	root := initGitRepo(t)
	infos := install.ListTargets(install.InstallOpts{ProjectRoot: root})
	var hook *install.DetectInfo
	for i := range infos {
		if infos[i].ID == install.TargetGitHook {
			hook = &infos[i]
			break
		}
	}
	if hook == nil {
		t.Fatalf("detect missing git-hook: %+v", infos)
	}
	if hook.Tier != install.TierConditional {
		t.Fatalf("tier: want CONDITIONAL got %q", hook.Tier)
	}
	if !hook.Detected {
		t.Fatalf("want detected in git repo, got reason %q", hook.Reason)
	}
}

func TestInstallGitHookWritesPostCommit(t *testing.T) {
	root := initGitRepo(t)
	tgt := gitHookTarget(t)

	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	hooksDir, err := resolveHooksDirForTest(root)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(hooksDir, "post-commit")
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("post-commit missing: %v", err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Fatal("post-commit must be executable")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "trace -C") {
		t.Fatalf("missing trace -C: %s", body)
	}
	if !strings.Contains(body, "index") {
		t.Fatalf("missing index: %s", body)
	}
	if !strings.Contains(body, installTraceBeginMarker()) {
		t.Fatal("missing begin marker")
	}
	if !strings.Contains(body, installTraceEndMarker()) {
		t.Fatal("missing end marker")
	}
}

func TestInstallGitHookHonorsCoreHooksPath(t *testing.T) {
	root := initGitRepo(t)
	git(t, root, "config", "core.hooksPath", ".husky")

	tgt := gitHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	custom := filepath.Join(root, ".husky", "post-commit")
	if _, err := os.Stat(custom); err != nil {
		t.Fatalf(".husky/post-commit missing: %v", err)
	}
	defaultHook := filepath.Join(root, ".git", "hooks", "post-commit")
	if _, err := os.Stat(defaultHook); err == nil {
		t.Fatal("expected post-commit under .husky, not .git/hooks")
	}
}

func TestInstallGitHookDoesNotWrapCommit(t *testing.T) {
	root := initGitRepo(t)
	tgt := gitHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	hooksDir, err := resolveHooksDirForTest(root)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(hooksDir, "post-commit"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	for _, forbidden := range []string{"git commit", "git add"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("hook must not contain %q:\n%s", forbidden, body)
		}
	}
}

func TestUninstallGitHookRemovesFragment(t *testing.T) {
	root := initGitRepo(t)
	hooksDir, err := resolveHooksDirForTest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(hooksDir, "post-commit")
	prior := "#!/bin/sh\n# user hook\necho keep-me\n"
	if err := os.WriteFile(path, []byte(prior), 0o755); err != nil {
		t.Fatal(err)
	}

	tgt := gitHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := tgt.Uninstall(install.InstallOpts{ProjectRoot: root}); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if strings.Contains(body, installTraceBeginMarker()) || strings.Contains(body, installTraceEndMarker()) {
		t.Fatalf("markers must be removed: %s", body)
	}
	if !strings.Contains(body, "keep-me") {
		t.Fatalf("sibling hook line must remain: %s", body)
	}
	if strings.Contains(body, "trace -C") {
		t.Fatalf("trace fragment must be gone: %s", body)
	}
}

func resolveHooksDirForTest(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--git-path", "hooks")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(out))
	if !filepath.IsAbs(p) {
		p = filepath.Join(root, p)
	}
	return filepath.Abs(p)
}

// installTraceBeginMarker / installTraceEndMarker mirror private constants for assertions.
func installTraceBeginMarker() string { return "# begin-trace" }
func installTraceEndMarker() string   { return "# end-trace" }
