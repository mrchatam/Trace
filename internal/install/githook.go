package install

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	traceBeginMarker = "# begin-trace"
	traceEndMarker   = "# end-trace"
)

var gitHookFiles = []string{"post-commit", "pre-push"}

type gitHookTarget struct{}

func (gitHookTarget) ID() string   { return TargetGitHook }
func (gitHookTarget) Tier() string { return TierConditional }

func (g gitHookTarget) Detect(opts InstallOpts) DetectResult {
	root := projectRoot(opts)
	ok, err := gitInsideWorkTree(root)
	if err != nil {
		return DetectResult{Detected: false, Reason: "git unavailable under " + root}
	}
	if ok {
		return DetectResult{Detected: true, Reason: "git work tree under " + root}
	}
	return DetectResult{Detected: false, Reason: "not a git work tree under " + root}
}

func (g gitHookTarget) Install(opts InstallOpts) error {
	root := projectRoot(opts)
	out := opts.Out
	if out == nil {
		out = io.Discard
	}
	errOut := opts.ErrOut
	if errOut == nil {
		errOut = io.Discard
	}

	fragment := traceHookFragment()

	if !opts.Write {
		if _, err := out.Write([]byte(fragment)); err != nil {
			return err
		}
		return nil
	}

	ok, err := gitInsideWorkTree(root)
	if err != nil {
		return fmt.Errorf("install: git-hook: %w", err)
	}
	if !ok {
		return fmt.Errorf("install: git-hook is CONDITIONAL — require git work tree under %s", root)
	}

	hooksDir, err := resolveHooksDir(root)
	if err != nil {
		return fmt.Errorf("install: git-hook: %w", err)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}

	for _, name := range gitHookFiles {
		path := filepath.Join(hooksDir, name)
		if err := upsertHookFile(path, fragment); err != nil {
			return err
		}
		fmt.Fprintf(errOut, "install: wrote %s\n", path)
	}
	return nil
}

func (g gitHookTarget) Uninstall(opts InstallOpts) error {
	root := projectRoot(opts)
	ok, err := gitInsideWorkTree(root)
	if err != nil || !ok {
		return nil // idempotent when not a git repo
	}
	hooksDir, err := resolveHooksDir(root)
	if err != nil {
		return nil
	}
	for _, name := range gitHookFiles {
		path := filepath.Join(hooksDir, name)
		if err := stripHookFile(path); err != nil {
			return err
		}
	}
	return nil
}

func traceHookFragment() string {
	return traceBeginMarker + "\n" +
		`ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
command -v trace >/dev/null 2>&1 || exit 0
PATHS="$(git diff-tree --no-commit-id --name-only -r HEAD 2>/dev/null || true)"
if [ -n "$PATHS" ]; then
  trace -C "$ROOT" index $PATHS || true
else
  trace -C "$ROOT" index || true
fi
trace -C "$ROOT" seed export -o trace/graph.json 2>/dev/null || true
` + traceEndMarker + "\n"
}

func upsertHookFile(path, fragment string) error {
	var existing []byte
	mode := os.FileMode(0o755)
	if fi, err := os.Stat(path); err == nil {
		mode = fi.Mode()
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		existing = raw
	} else if !os.IsNotExist(err) {
		return err
	}

	updated := UpsertMarkedFragment(string(existing), traceBeginMarker, traceEndMarker, fragment)
	if len(existing) == 0 && !strings.HasPrefix(updated, "#!") {
		updated = "#!/bin/sh\n" + updated
	}
	perm := mode & 0o777
	if perm&0o100 == 0 {
		perm |= 0o111
	}
	if err := os.WriteFile(path, []byte(updated), perm); err != nil {
		return err
	}
	return os.Chmod(path, perm)
}

func stripHookFile(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	mode := fi.Mode()
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	updated := RemoveMarkedFragment(string(raw), traceBeginMarker, traceEndMarker)
	if updated == string(raw) {
		return nil
	}
	if strings.TrimSpace(updated) == "" {
		return os.Remove(path)
	}
	if err := os.WriteFile(path, []byte(updated), mode&0o777); err != nil {
		return err
	}
	return os.Chmod(path, mode)
}

func gitInsideWorkTree(root string) (bool, error) {
	out, err := runGit(root, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if isGitNotRepo(err) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}

func resolveHooksDir(root string) (string, error) {
	out, err := runGit(root, "rev-parse", "--git-path", "hooks")
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(out)
	if p == "" {
		return "", fmt.Errorf("empty hooks path")
	}
	if !filepath.IsAbs(p) {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			return "", err
		}
		p = filepath.Join(absRoot, p)
	}
	return filepath.Abs(p)
}

func runGit(root string, args ...string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	full := append([]string{"-C", abs}, args...)
	cmd := exec.Command("git", full...)
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", fmt.Errorf("git not found in PATH")
		}
		return "", fmt.Errorf("git %v: %w", args, err)
	}
	return string(out), nil
}

func isGitNotRepo(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not a git repository") ||
		strings.Contains(msg, "fatal:")
}
