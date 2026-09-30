package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// CoordinateTestRunFunc invokes relevant tests for a task (wired from internal/testrun at CLI).
type CoordinateTestRunFunc func(ctx context.Context, taskID string, paths []string) ([]store.OutcomeResult, error)

// CoordinateTestHooks overrides cycle steps for unit tests (call order proofs).
type CoordinateTestHooks struct {
	RunTests           CoordinateTestRunFunc
	EnsureVerification func(ctx context.Context) error
	RunEvaluation      func(ctx context.Context) (store.OutcomeResult, error)
}

// CoordinateOptions controls verification-cycle coordination.
type CoordinateOptions struct {
	ForceEval   bool
	Paths       []string
	RunTests    CoordinateTestRunFunc
	EvidenceIDs []string
	ScoresJSON  string
	Actor       string
	SourceType  string
	Hooks       *CoordinateTestHooks
}

// TestRegressionSignal is the C13 regression comparison result.
type TestRegressionSignal struct {
	Detected    bool   `json:"regression_detected"`
	TestName    string `json:"test_name,omitempty"`
	PriorPass   bool   `json:"prior_pass,omitempty"`
	CurrentFail bool   `json:"current_fail,omitempty"`
}

// CoordinateResult is the JSON payload from CoordinateVerification.
type CoordinateResult struct {
	TaskID             string                `json:"task_id"`
	TestOutcomes       []store.OutcomeResult `json:"test_outcomes,omitempty"`
	VerificationFound  bool                  `json:"verification_found"`
	EvaluationRecorded bool                  `json:"evaluation_recorded"`
	StoppedEarly       bool                  `json:"stopped_early"`
	StopReason         string                `json:"stop_reason,omitempty"`
	Regression         TestRegressionSignal  `json:"regression"`
	Steps              []string              `json:"steps,omitempty"`
}

// DetectTestRegression compares the latest two kind=test rows for testName.
// Prior pass + current fail yields regression_detected.
func (s *Service) DetectTestRegression(ctx context.Context, taskID, testName string) (TestRegressionSignal, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	testName = strings.TrimSpace(testName)
	if taskID == "" {
		return TestRegressionSignal{}, &ErrValidation{Msg: "task_id is required"}
	}
	if testName == "" {
		return TestRegressionSignal{}, &ErrValidation{Msg: "test_name is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return TestRegressionSignal{}, err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindTest)
	if err != nil {
		return TestRegressionSignal{}, err
	}
	var matching []store.OutcomeResult
	for _, o := range rows {
		if o.TestName == testName {
			matching = append(matching, o)
		}
	}
	if len(matching) < 2 {
		return TestRegressionSignal{TestName: testName}, nil
	}
	prior := matching[len(matching)-2]
	current := matching[len(matching)-1]
	sig := TestRegressionSignal{
		TestName:    testName,
		PriorPass:   prior.TestStatus == store.TestStatusPass,
		CurrentFail: current.TestStatus == store.TestStatusFail,
	}
	sig.Detected = sig.PriorPass && sig.CurrentFail
	return sig, nil
}

// DetectAnyTestRegression scans test names on the task for the first C13 signal.
func (s *Service) DetectAnyTestRegression(ctx context.Context, taskID string) (TestRegressionSignal, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return TestRegressionSignal{}, &ErrValidation{Msg: "task_id is required"}
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindTest)
	if err != nil {
		return TestRegressionSignal{}, err
	}
	seen := map[string]struct{}{}
	for _, o := range rows {
		if o.TestName == "" {
			continue
		}
		if _, ok := seen[o.TestName]; ok {
			continue
		}
		seen[o.TestName] = struct{}{}
		sig, err := s.DetectTestRegression(ctx, taskID, o.TestName)
		if err != nil {
			return TestRegressionSignal{}, err
		}
		if sig.Detected {
			return sig, nil
		}
	}
	return TestRegressionSignal{}, nil
}

// CoordinateVerification runs test → verification → evaluation in order (C36).
func (s *Service) CoordinateVerification(ctx context.Context, taskID string, opts CoordinateOptions) (CoordinateResult, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return CoordinateResult{}, &ErrValidation{Msg: "task_id is required"}
	}
	task, err := s.requireTask(taskID)
	if err != nil {
		return CoordinateResult{}, err
	}
	goalID := ""
	if task.GoalID != nil {
		goalID = strings.TrimSpace(*task.GoalID)
	}

	out := CoordinateResult{TaskID: taskID}
	appendStep := func(step string) { out.Steps = append(out.Steps, step) }

	testPending, err := s.isTestPending(ctx, taskID)
	if err != nil {
		return CoordinateResult{}, err
	}

	if testPending {
		appendStep("test")
		runTests := opts.RunTests
		if opts.Hooks != nil && opts.Hooks.RunTests != nil {
			runTests = opts.Hooks.RunTests
		}
		if runTests == nil {
			return out, fmt.Errorf("coordinate: test run function is required when test_pending")
		}
		testOutcomes, err := runTests(ctx, taskID, opts.Paths)
		if err != nil {
			return out, fmt.Errorf("coordinate: test run: %w", err)
		}
		out.TestOutcomes = testOutcomes
	}

	failed, err := s.latestTestsFailedSinceChange(taskID)
	if err != nil {
		return CoordinateResult{}, err
	}
	if failed && !opts.ForceEval {
		out.StoppedEarly = true
		out.StopReason = "test_failed"
		appendStep("stop_test_failed")
		reg, err := s.DetectAnyTestRegression(ctx, taskID)
		if err != nil {
			return out, err
		}
		out.Regression = reg
		return out, nil
	}

	verifyDebt, err := s.HasVerificationDebt(ctx, taskID)
	if err != nil {
		return CoordinateResult{}, err
	}
	if verifyDebt {
		appendStep("verify")
		if opts.Hooks != nil && opts.Hooks.EnsureVerification != nil {
			if err := opts.Hooks.EnsureVerification(ctx); err != nil {
				return out, fmt.Errorf("coordinate: verification: %w", err)
			}
		} else if len(opts.EvidenceIDs) > 0 && goalID != "" {
			_, err := s.RecordVerificationOutcome(ctx, VerificationOutcomeInput{
				TaskID:             taskID,
				GoalID:             goalID,
				VerificationStatus: store.VerificationStatusVerified,
				EvidenceIDs:        opts.EvidenceIDs,
				Actor:              opts.Actor,
				SourceType:         opts.SourceType,
			})
			if err != nil {
				return out, fmt.Errorf("coordinate: record verification: %w", err)
			}
		}
	}
	ok, _, err := s.CheckVerificationGate(ctx, taskID)
	if err != nil {
		return CoordinateResult{}, err
	}
	out.VerificationFound = ok

	evalPending, err := s.isEvaluationPending(ctx, taskID)
	if err != nil {
		return CoordinateResult{}, err
	}
	if evalPending && strings.TrimSpace(opts.ScoresJSON) != "" {
		appendStep("evaluate")
		if opts.Hooks != nil && opts.Hooks.RunEvaluation != nil {
			if _, err := opts.Hooks.RunEvaluation(ctx); err != nil {
				return out, fmt.Errorf("coordinate: evaluation: %w", err)
			}
			out.EvaluationRecorded = true
		} else {
			baseline, baselineErr := s.ensureBaselineForTaskCommit(ctx, taskID, opts.ScoresJSON)
			if baselineErr != nil {
				if errors.Is(baselineErr, sql.ErrNoRows) {
					// No change commit to anchor to — soft stop (pre-existing
					// contract). Never fabricate a baseline without a commit.
					out.StopReason = "no_active_baseline"
				} else {
					return out, fmt.Errorf("coordinate: baseline: %w", baselineErr)
				}
			} else {
				evalOut, err := s.RecordEvaluationOutcome(ctx, EvaluationOutcomeInput{
					TaskID:     taskID,
					BaselineID: baseline.ID,
					ScoresJSON: opts.ScoresJSON,
					Actor:      opts.Actor,
					SourceType: opts.SourceType,
				})
				if err != nil {
					return out, fmt.Errorf("coordinate: record evaluation: %w", err)
				}
				_ = evalOut
				out.EvaluationRecorded = true
			}
		}
	}

	reg, err := s.DetectAnyTestRegression(ctx, taskID)
	if err != nil {
		return out, err
	}
	out.Regression = reg
	return out, nil
}

func (s *Service) isTestPending(ctx context.Context, taskID string) (bool, error) {
	impl, err := s.HasImplementationSignal(ctx, taskID)
	if err != nil {
		return false, err
	}
	if !impl {
		return false, nil
	}
	since, err := s.HasTestOutcomeSinceLatestChange(ctx, taskID)
	if err != nil {
		return false, err
	}
	return !since, nil
}

func (s *Service) isEvaluationPending(ctx context.Context, taskID string) (bool, error) {
	verifyGate, _, err := s.CheckVerificationGate(ctx, taskID)
	if err != nil {
		return false, err
	}
	hasVerif, err := s.HasVerificationOutcome(ctx, taskID)
	if err != nil {
		return false, err
	}
	hasEval, err := s.HasComputedEvaluation(ctx, taskID)
	if err != nil {
		return false, err
	}
	return (verifyGate || hasVerif) && !hasEval, nil
}

func (s *Service) latestTestsFailedSinceChange(taskID string) (bool, error) {
	changes, err := s.store.ListChangesByTaskID(taskID)
	if err != nil {
		return false, err
	}
	latest, ok := latestRecordedOrComparedChange(changes)
	if !ok {
		return false, nil
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindTest)
	if err != nil {
		return false, err
	}
	for _, o := range rows {
		if o.CreatedAt >= latest.CreatedAt && o.TestStatus == store.TestStatusFail {
			return true, nil
		}
	}
	return false, nil
}

// activeBaselineForTaskCommit returns the active baseline for the latest
// RECORDED/COMPARED change's git commit. It is a lookup only — it never writes.
func (s *Service) activeBaselineForTaskCommit(taskID string) (store.Baseline, error) {
	sha, err := s.latestTaskCommit(taskID)
	if err != nil {
		return store.Baseline{}, err
	}
	return s.store.GetActiveBaselineByCommitLabel(sha, "")
}

// ensureBaselineForTaskCommit returns the active baseline for the task's
// latest change commit, creating one on demand when none exists. A baseline is
// a snapshot of scores at a commit, created when work is evaluated — not
// imported from a seed file. The snapshot is taken from the evaluation's own
// scores, so the first evaluation at a commit is an initialization pass (every
// delta is zero; nothing can regress against a snapshot that did not exist
// yet), and later re-evaluations at the same commit compare against that
// snapshot — drift within the commit is what becomes visible.
//
// When an active baseline already exists for the commit it is returned
// unchanged; scores are never mutated after creation. When a superseded
// baseline exists for the commit it is promoted back to active via the
// existing status/supersedes_id chain (Law 11 — no delete) instead of
// inserting a duplicate row for the same commit. Missing change or commit
// surfaces as sql.ErrNoRows (soft stop upstream) — no baseline is fabricated.
func (s *Service) ensureBaselineForTaskCommit(ctx context.Context, taskID, scoresJSON string) (store.Baseline, error) {
	sha, err := s.latestTaskCommit(taskID)
	if err != nil {
		return store.Baseline{}, err
	}
	if b, err := s.store.GetActiveBaselineByCommitLabel(sha, ""); err == nil {
		return b, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return store.Baseline{}, err
	}
	// No active baseline for this commit: create one. If a superseded
	// baseline exists for the same commit, PromoteBaseline reactivates it and
	// links the supersedes chain; otherwise it inserts a fresh active row.
	rows, err := s.store.ListAllBaselines()
	if err != nil {
		return store.Baseline{}, err
	}
	for _, b := range rows {
		if b.GitCommit != sha {
			continue
		}
		return s.PromoteBaseline(ctx, b.ID)
	}
	return s.CreateBaseline(ctx, BaselineInput{
		GitCommit:  sha,
		ScoresJSON: scoresJSON,
		SourceType: "COORDINATED",
	})
}

// latestTaskCommit resolves the git commit of the latest RECORDED/COMPARED
// change for a task. Fails closed (sql.ErrNoRows) when there is no change or
// the change has no commit — an evaluation without a commit has nothing to
// anchor to.
func (s *Service) latestTaskCommit(taskID string) (string, error) {
	changes, err := s.store.ListChangesByTaskID(taskID)
	if err != nil {
		return "", err
	}
	latest, ok := latestRecordedOrComparedChange(changes)
	if !ok || strings.TrimSpace(latest.GitCommit) == "" {
		return "", sql.ErrNoRows
	}
	sha, err := normalizeGitCommit(latest.GitCommit)
	if err != nil {
		return "", err
	}
	return sha, nil
}
