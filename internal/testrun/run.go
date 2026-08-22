package testrun

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const defaultRunTimeout = 5 * time.Minute

// Options controls test invocation.
type Options struct {
	Paths      []string
	Runner     Runner
	Timeout    time.Duration
	Actor      string
	SourceType string
}

// RunRelevantTests selects relevant tests, invokes the runner, and records kind=test outcomes.
func RunRelevantTests(
	ctx context.Context,
	st *store.Store,
	dom *domain.Service,
	taskID string,
	opts Options,
) ([]store.OutcomeResult, error) {
	if st == nil {
		return nil, fmt.Errorf("testrun: store is required")
	}
	if dom == nil {
		return nil, fmt.Errorf("testrun: domain service is required")
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &domain.ErrValidation{Msg: "task_id is required"}
	}

	baseSpec, ok, err := resolveDefaultRunner(st.ProjectRoot())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("testrun: no test command available")
	}

	customCfg, _ := loadRunnerConfig(st.ProjectRoot())
	useCustom := customCfg != nil

	targets, err := SelectTestTargets(ctx, st, dom, taskID, opts.Paths)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, errors.New("testrun: no relevant tests selected")
	}

	runner := opts.Runner
	if runner == nil {
		runner = ExecRunner{}
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultRunTimeout
	}
	actor := strings.TrimSpace(opts.Actor)
	if actor == "" {
		actor = "cli"
	}
	sourceType := strings.TrimSpace(opts.SourceType)
	if sourceType == "" {
		sourceType = "CLI"
	}

	var recorded []store.OutcomeResult
	for _, target := range targets {
		spec := baseSpec
		if !useCustom {
			spec = goTestSpec(st.ProjectRoot(), target)
		}
		runCtx, cancel := context.WithTimeout(ctx, timeout)
		exitCode, output, runErr := runner.Run(runCtx, spec)
		cancel()

		status := exitCodeToStatus(exitCode, runErr)
		out, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
			TaskID:     taskID,
			TestName:   target.Name,
			TestStatus: status,
			Summary:    output,
			Actor:      actor,
			SourceType: sourceType,
		})
		if err != nil {
			return recorded, err
		}
		recorded = append(recorded, out)
	}
	return recorded, nil
}

func goTestSpec(root string, target TestTarget) RunSpec {
	args := []string{"test"}
	if target.Package != "" {
		args = append(args, target.Package)
	} else {
		args = append(args, "./...")
	}
	if target.RunPattern != "" {
		args = append(args, "-run", target.RunPattern)
	}
	return RunSpec{Command: "go", Args: args, Cwd: root}
}

func exitCodeToStatus(exitCode int, runErr error) string {
	if runErr != nil {
		return store.TestStatusError
	}
	if exitCode == 0 {
		return store.TestStatusPass
	}
	return store.TestStatusFail
}
