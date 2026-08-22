package gitcli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mrchatam/Trace/internal/vcs"
)

// runner executes git -C <root> … with context.
type runner struct {
	root string
	git  string
}

func newRunner(root string) *runner {
	return &runner{root: root, git: "git"}
}

func (r *runner) run(ctx context.Context, args ...string) (string, error) {
	out, err := r.runBytes(ctx, args...)
	return string(out), err
}

func (r *runner) runBytes(ctx context.Context, args ...string) ([]byte, error) {
	full := append([]string{"-C", r.root}, args...)
	cmd := exec.CommandContext(ctx, r.git, full...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		if errors.Is(err, exec.ErrNotFound) {
			return nil, &vcs.Error{Op: "git", Err: vcs.ErrGitMissing}
		}
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, fmt.Errorf("git %v: %s", args, msg)
		}
		return nil, fmt.Errorf("git %v: %w (%s)", args, err, msg)
	}
	return stdout.Bytes(), nil
}
