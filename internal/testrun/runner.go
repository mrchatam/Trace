package testrun

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

const maxOutcomeSummaryBytes = 4096

// RunSpec is one subprocess invocation.
type RunSpec struct {
	Command string
	Args    []string
	Cwd     string
}

// Runner executes a test command. Stub in unit tests; ExecRunner in production.
type Runner interface {
	Run(ctx context.Context, spec RunSpec) (exitCode int, output string, err error)
}

// ExecRunner runs real subprocesses with context cancellation.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, spec RunSpec) (int, string, error) {
	cmd := exec.CommandContext(ctx, spec.Command, spec.Args...)
	if spec.Cwd != "" {
		cmd.Dir = spec.Cwd
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := truncateSummary(buf.String())
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), out, nil
		}
		return -1, out, err
	}
	return 0, out, nil
}

func truncateSummary(s string) string {
	if len(s) <= maxOutcomeSummaryBytes {
		return s
	}
	return s[:maxOutcomeSummaryBytes]
}
