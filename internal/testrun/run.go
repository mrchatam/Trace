package testrun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
//
// Attribution: for a custom runner the configured command is invoked exactly
// once (args passed through verbatim) and its per-file results — from
// --test-runner-results.json when the runner wrote one, else parsed from
// stdout — are mapped onto the selected targets. A target the run did not
// report gets a kind=test outcome with status=skip (source_type
// TESTRUNNER_UNREPORTED); a failing command without parseable per-file output
// records the whole-run exit code ONLY on the synthetic "suite" target, never
// under a real test file's name. The go test fallback keeps per-target
// invocations because `go test -run` filters per package and exit codes are
// per-invocation.
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

	baseSpec, source, err := resolveDefaultRunner(st.ProjectRoot())
	if err != nil {
		return nil, err
	}
	if source == runnerSourceNone {
		return nil, errors.New("testrun: no test command available")
	}

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

	if source == runnerSourceCustom {
		// An explicit {path} placeholder keeps the documented per-target mode:
		// each invocation executes exactly the selected file, so its exit code
		// is that file's own verdict (honest per-target attribution).
		if substitutable, ok := substitutableSpecForAnyTarget(baseSpec, targets); ok {
			return runCustomPerTarget(ctx, dom, substitutable, runner, timeout, actor, sourceType, taskID, targets)
		}
		return runCustomOnce(ctx, dom, baseSpec, runner, timeout, actor, sourceType, taskID, targets)
	}
	return runGoPerTarget(ctx, st, dom, runner, timeout, actor, sourceType, taskID, targets)
}

// runCustomPerTarget keeps the documented {path} contract: one invocation per
// selected target with the placeholder replaced by that target's path. The
// exit code belongs to exactly that file, so recording it under the target's
// own name is honest attribution.
func runCustomPerTarget(
	ctx context.Context,
	dom *domain.Service,
	specFor func(TestTarget) RunSpec,
	runner Runner,
	timeout time.Duration,
	actor, sourceType, taskID string,
	targets []TestTarget,
) ([]store.OutcomeResult, error) {
	var recorded []store.OutcomeResult
	for _, target := range targets {
		spec := specFor(target)
		runCtx, cancel := context.WithTimeout(ctx, timeout)
		exitCode, output, runErr := runner.Run(runCtx, spec)
		cancel()

		out, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
			TaskID:     taskID,
			TestName:   target.Name,
			TestStatus: exitCodeToStatus(exitCode, runErr),
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

// specForTarget substitutes the {path} placeholder in a custom runner's args
// with the selected test file path (root-relative). Args without the
// placeholder pass through unchanged.
func specForTarget(base RunSpec, target TestTarget) RunSpec {
	if target.Path == "" {
		return base
	}
	var sub bool
	args := make([]string, 0, len(base.Args))
	for _, a := range base.Args {
		if strings.Contains(a, "{path}") {
			args = append(args, strings.ReplaceAll(a, "{path}", target.Path))
			sub = true
		} else {
			args = append(args, a)
		}
	}
	if !sub {
		return base
	}
	return RunSpec{Command: base.Command, Args: args, Cwd: base.Cwd}
}

// substitutableSpecForAnyTarget reports whether the config's {path} placeholder
// resolves for every selected target (all have a path) and returns the
// per-target spec function when so.
func substitutableSpecForAnyTarget(base RunSpec, targets []TestTarget) (func(TestTarget) RunSpec, bool) {
	if len(targets) == 0 {
		return nil, false
	}
	for _, t := range targets {
		if t.Path == "" {
			return nil, false
		}
	}
	substitutes := false
	for _, a := range base.Args {
		if strings.Contains(a, "{path}") {
			substitutes = true
			break
		}
	}
	if !substitutes {
		return nil, false
	}
	return func(t TestTarget) RunSpec { return specForTarget(base, t) }, true
}

// runCustomOnce invokes the configured command once and maps reported results
// onto the selected targets. Results are attributed only to targets the run
// actually reported; everything else is skipped honestly (no invented pass).
func runCustomOnce(
	ctx context.Context,
	dom *domain.Service,
	baseSpec RunSpec,
	runner Runner,
	timeout time.Duration,
	actor, sourceType, taskID string,
	targets []TestTarget,
) ([]store.OutcomeResult, error) {
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	exitCode, output, runErr := runner.Run(runCtx, baseSpec)
	cancel()

	results, parseNote := parseRunnerResults(baseSpec, output)
	reported := map[string]testFileResult{}
	for _, r := range results {
		reported[r.Name] = r
	}
	if len(results) == 0 {
		// No per-file report available: the whole-run verdict must not be
		// mis-attributed to files that were never individually executed
		// (regression: one failing command marked every selected target fail).
		// Record one outcome under a synthetic name so the run is honest and
		// the summary is preserved; SelectTestTargets names are left untouched.
		name := suiteOutcomeName(baseSpec)
		status := exitCodeToStatus(exitCode, runErr)
		out, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
			TaskID:     taskID,
			TestName:   name,
			TestStatus: status,
			Summary:    outcomeSummary(baseSpec, parseNote, output),
			Actor:      actor,
			SourceType: sourceType,
		})
		if err != nil {
			return nil, err
		}
		return []store.OutcomeResult{out}, nil
	}

	var recorded []store.OutcomeResult
	for _, target := range targets {
		status, ok, err := statusForTarget(reported, target)
		if err != nil {
			return recorded, err
		}
		if !ok {
			// The run did not report this target. Record skip explicitly —
			// surfacing unreported targets instead of silently dropping them.
			out, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
				TaskID:     taskID,
				TestName:   target.Name,
				TestStatus: store.TestStatusSkip,
				Summary:    unreportedSummary(baseSpec, target, parseNote),
				Actor:      actor,
				SourceType: sourceTypeUnreported,
			})
			if err != nil {
				return recorded, err
			}
			recorded = append(recorded, out)
			continue
		}
		out, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
			TaskID:     taskID,
			TestName:   target.Name,
			TestStatus: status,
			Summary:    outcomeSummary(baseSpec, parseNote, output),
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

// statusForTarget resolves one target against the reported per-file results.
// ok=false means the run did not report the target. Reporters print
// repo-relative paths while target names are basenames, so exact-name and
// exact-path matches win first and a unique path-suffix match ("src/foo.test.ts"
// vs target "foo.test.ts") is accepted; keys are walked in sorted order so the
// resolution is deterministic.
func statusForTarget(reported map[string]testFileResult, target TestTarget) (string, bool, error) {
	if r, ok := reported[target.Name]; ok {
		return r.Status, true, nil
	}
	if target.Path != "" {
		for _, key := range []string{target.Path, "./" + target.Path} {
			if r, ok := reported[key]; ok {
				return r.Status, true, nil
			}
		}
	}
	// Repo-relative report keys are matched with a leading ./ stripped so
	// "./src/foo.test.ts" and "src/foo.test.ts" are the same file (deduped).
	// Path-suffix fallback covers repo-relative reports vs basename targets,
	// but only when unique: two dirs with the same basename must not be guessed.
	normKeys := make([]string, 0, len(reported))
	byNorm := map[string]string{}
	for k := range reported {
		norm := strings.TrimPrefix(k, "./")
		if _, ok := byNorm[norm]; !ok {
			byNorm[norm] = k
			normKeys = append(normKeys, norm)
		}
	}
	sort.Strings(normKeys)
	var match *testFileResult
	for _, norm := range normKeys {
		if norm == target.Name || (target.Path != "" && norm == target.Path) {
			return reported[byNorm[norm]].Status, true, nil
		}
		if strings.HasSuffix(norm, "/"+target.Name) {
			if match != nil {
				return "", false, nil
			}
			r := reported[byNorm[norm]]
			match = &r
		}
	}
	if match != nil {
		return match.Status, true, nil
	}
	return "", false, nil
}

// runGoPerTarget keeps the go fallback semantics: `go test` exit codes are
// per-invocation, so each target is invoked (and recorded) individually.
func runGoPerTarget(
	ctx context.Context,
	st *store.Store,
	dom *domain.Service,
	runner Runner,
	timeout time.Duration,
	actor, sourceType, taskID string,
	targets []TestTarget,
) ([]store.OutcomeResult, error) {
	var recorded []store.OutcomeResult
	for _, target := range targets {
		spec := goTestSpec(st.ProjectRoot(), target)
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

// suiteOutcomeName is the synthetic outcome name for a custom-runner run with
// no parseable per-file results. Distinct from any real file so the whole-run
// verdict can never be mis-attributed under a target's name.
func suiteOutcomeName(spec RunSpec) string {
	return fmt.Sprintf("suite:%s %s", spec.Command, strings.Join(spec.Args, " "))
}

func outcomeSummary(spec RunSpec, note, output string) string {
	return clampSummary(
		fmt.Sprintf("custom runner: %s %s (%s)\n", spec.Command, strings.Join(spec.Args, " "), note),
		output,
	)
}

// clampSummary keeps header + head of output within the domain's fail-closed
// summary cap (RecordTestOutcome rejects summaries over 4096 bytes).
func clampSummary(header, output string) string {
	const suffix = "\n…truncated"
	budget := maxOutcomeSummaryBytes - len(header)
	if budget < 0 {
		budget = 0
	}
	if len(output) > budget {
		if budget <= len(suffix) {
			return header
		}
		output = output[:budget-len(suffix)] + suffix
	}
	return header + output
}

func unreportedSummary(spec RunSpec, target TestTarget, note string) string {
	return clampSummary(
		fmt.Sprintf("custom runner: %s %s (%s)\n", spec.Command, strings.Join(spec.Args, " "), note),
		fmt.Sprintf("target %q was selected but the run reported no result for it; not executed as an individual test", target.Name),
	)
}

// --- result extraction ------------------------------------------------------

// testFileResult is one per-file result reported by a runner.
type testFileResult struct {
	Name   string
	Status string
}

// sourceTypeUnreported marks skip outcomes for selected targets the run did
// not report. Distinct source_type keeps honest skips auditable (Law 2 — a
// skip is not a pass and must not inherit the invoking actor's claim).
const sourceTypeUnreported = "TESTRUNNER_UNREPORTED"

// vitestFileLine matches vitest dot/default reporter per-file summary lines:
// a status glyph, the file path (repo-relative, no leading ./), a parenthesized
// per-file test summary, and the duration. Examples:
//
//	✓ src/foo.test.ts (2 tests) 12ms
//	× src/bar.test.ts (3 tests | 1 failed) 9ms
//
// Capture 1 is the file path, capture 2 the paren content. Unknown glyphs do
// not match — a file whose status cannot be determined is treated as
// unreported rather than guessed.
var vitestFileLine = regexp.MustCompile(
	`^[\p{So}\p{Sm}\p{Sk}xX]\s+(\S+\.tsx?)(?:\s+\(([^)]*)\))?\s+\d+(?:\.\d+)?ms\s*$`)

// parseRunnerResults extracts per-file results from a custom runner run.
// Preference order: trace/test-runner-results.json (machine-readable,
// runner-written; none of the bundled runners write it today), then vitest
// dot/default reporter stdout. Second return is an honest parse note for the
// outcome summary.
func parseRunnerResults(spec RunSpec, output string) ([]testFileResult, string) {
	if spec.Cwd != "" {
		if results, note, ok := tryResultsFile(spec.Cwd); ok {
			return results, note
		}
	}
	if results, ok := parseVitestOutput(output); ok {
		return results, "parsed from stdout (vitest reporter)"
	}
	return nil, "no per-file report found (write trace/test-runner-results.json for per-target attribution)"
}

// TryResultsFile reads trace/test-runner-results.json relative to the runner
// cwd. ok=false covers both missing file and malformed content — malformed
// falls back to stdout parsing rather than failing the run.
func tryResultsFile(cwd string) ([]testFileResult, string, bool) {
	b, err := os.ReadFile(filepath.Join(cwd, "trace", "test-runner-results.json"))
	if err != nil {
		return nil, "", false
	}
	var doc struct {
		Files []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"files"`
	}
	if err := json.Unmarshal(b, &doc); err != nil || len(doc.Files) == 0 {
		return nil, "", false
	}
	out := make([]testFileResult, 0, len(doc.Files))
	for _, f := range doc.Files {
		name := strings.TrimSpace(f.Name)
		status, err := normalizeStatus(f.Status)
		if name == "" || err != nil {
			continue
		}
		out = append(out, testFileResult{Name: name, Status: status})
	}
	if len(out) == 0 {
		return nil, "", false
	}
	return out, "read from trace/test-runner-results.json", true
}

// parseVitestOutput maps vitest's per-file summary lines onto test targets.
// Dot and default reporters print one line per file: a status glyph, the file
// path (without the leading ./), an optional parenthesized per-file summary,
// and the duration. A file fails when its summary names failed tests (or the
// glyph is an explicit cross); everything else the regex matched passes.
func parseVitestOutput(output string) ([]testFileResult, bool) {
	var out []testFileResult
	seen := map[string]struct{}{}
	for _, line := range strings.Split(output, "\n") {
		m := vitestFileLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := strings.TrimPrefix(strings.TrimSpace(m[1]), "./")
		if name == "" {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		paren := strings.ToLower(strings.TrimSpace(m[2]))
		glyph := strings.TrimSpace(line)
		status := store.TestStatusPass
		if strings.Contains(paren, "failed") || strings.HasPrefix(glyph, "\u2717") || strings.HasPrefix(glyph, "\u274c") {
			status = store.TestStatusFail
		}
		out = append(out, testFileResult{Name: name, Status: status})
	}
	return out, len(out) > 0
}

// normalizeStatus accepts the vocabulary shared by test-runner-results.json
// and the runner-side status mapping (pass/fail/skip/error).
func normalizeStatus(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	switch s {
	case store.TestStatusPass, store.TestStatusFail, store.TestStatusSkip, store.TestStatusError:
		return s, nil
	default:
		return "", fmt.Errorf("unknown test status %q", raw)
	}
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
