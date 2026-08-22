package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func setupVerificationDebtFixture(t *testing.T, dir string) (goalID, taskID string) {
	t.Helper()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID = addGoalForLoopTest(t, dir, "Goal")
	taskID = addTaskForLoopTest(t, dir, goalID, "Debt task")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	if _, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	return goalID, taskID
}

func setupCleanFullCycleFixture(t *testing.T, dir string) (goalID, taskID string) {
	t.Helper()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID = addGoalForLoopTest(t, dir, "Goal")
	taskID = addTaskForLoopTest(t, dir, goalID, "Clean task")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskID, TestName: "TestAll", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	ev, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "verify proof"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             taskID,
		GoalID:             goalID,
		VerificationStatus: store.VerificationStatusVerified,
		EvidenceIDs:        []string{ev.ID},
	}); err != nil {
		t.Fatal(err)
	}
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.95}`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskID, Summary: "cycle complete", UsefulTests: []string{"TestAll"},
	}); err != nil {
		t.Fatal(err)
	}
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Architecture choice"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, taskID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	return goalID, taskID
}

func createReviewPassForTask(t *testing.T, dir, taskID string) {
	t.Helper()
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS", "--reason", "start"}); code != exitOK {
		t.Fatalf("IN_PROGRESS: %d", code)
	}
	revOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "create", "--title", "Gate", "--task", taskID})
	})
	var revRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(revOut)), &revRes); err != nil {
		t.Fatalf("review create: %v (%s)", err, revOut)
	}
	revID := revRes["id"].(string)
	if code := run([]string{"-C", dir, "review", "set", "--id", revID, "--result", "PASS", "--reason", "ok", "--actor", "rev"}); code != exitOK {
		t.Fatalf("review set: %d", code)
	}
}

func assertTaskWorkState(t *testing.T, dir, taskID, want string) {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task, err := st.GetTask(taskID)
	if err != nil {
		t.Fatal(err)
	}
	if task.WorkState != want {
		t.Fatalf("work_state want %q got %q", want, task.WorkState)
	}
}

func TestTransitionDoneEnforceBlocksVerificationDebt(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupVerificationDebtFixture(t, dir)
	createReviewPassForTask(t, dir, taskID)

	code, _, stderr := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "promote", "--as-operator", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "verification_incomplete") {
		t.Fatalf("stderr missing verification_incomplete: %q", stderr)
	}
	assertTaskWorkState(t, dir, taskID, store.WorkStateInProgress)
}

func TestTransitionDoneWithoutEnforceUnchanged(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupVerificationDebtFixture(t, dir)
	createReviewPassForTask(t, dir, taskID)

	code, _, _ := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "promote", "--as-operator",
	})
	if code != exitOK {
		t.Fatalf("want exitOK got %d", code)
	}
	assertTaskWorkState(t, dir, taskID, store.WorkStateDone)
}

func TestTransitionDoneEnforceAllowsClean(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupCleanFullCycleFixture(t, dir)
	createReviewPassForTask(t, dir, taskID)

	code, _, _ := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "promote", "--as-operator", "--enforce",
	})
	if code != exitOK {
		t.Fatalf("want exitOK got %d", code)
	}
	assertTaskWorkState(t, dir, taskID, store.WorkStateDone)
}

func TestTransitionDoneEnforcePreservesAllowDone(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Hatch"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatal(err)
	}
	taskID := taskRes["id"].(string)
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS", "--reason", "start"}); code != exitOK {
		t.Fatalf("IN_PROGRESS: %d", code)
	}
	code, _, stderr := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "escape", "--allow-done",
	})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "WARNING") || !strings.Contains(stderr, "allow-done") {
		t.Fatalf("expected allow-done WARNING, got %q", stderr)
	}
}

func TestTransitionDoneEnforceBlocksDespiteAllowDone(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupVerificationDebtFixture(t, dir)
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS", "--reason", "start"}); code != exitOK {
		t.Fatalf("IN_PROGRESS: %d", code)
	}
	code, _, _ := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "escape", "--allow-done", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d", code)
	}
	assertTaskWorkState(t, dir, taskID, store.WorkStateInProgress)
}

func TestTransitionDoneEnforceIgnoredForNonDone(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupVerificationDebtFixture(t, dir)
	code, _, _ := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS",
		"--reason", "start", "--enforce",
	})
	if code != exitOK {
		t.Fatalf("enforce on non-DONE want exitOK got %d", code)
	}
}

func TestTransitionDoneEnforcePreservesDomainReviewGate(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Review gate"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatal(err)
	}
	taskID := taskRes["id"].(string)
	createReviewPassForTask(t, dir, taskID)

	code, _, _ := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE", "--reason", "promote",
	})
	if code == exitOK {
		t.Fatal("DONE after PASS without --as-operator must fail")
	}
}

func TestTransitionDoneEnforceStderrHint(t *testing.T) {
	dir := t.TempDir()
	_, taskID := setupVerificationDebtFixture(t, dir)
	createReviewPassForTask(t, dir, taskID)

	code, stdout, stderr := runCapture(t, []string{
		"-C", dir, "transition", "--task", taskID, "--to", "DONE",
		"--reason", "promote", "--as-operator", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d", code)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Fatal("expected non-empty stderr on gate block")
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("blocked enforce must not emit success JSON on stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "verification") {
		t.Fatalf("stderr should contain violation hint: %q", stderr)
	}
}

func TestHelpIncludesTransitionEnforce(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "--enforce") {
		t.Fatalf("help missing --enforce: %q", out)
	}
	if !strings.Contains(out, "GateForDone") || !strings.Contains(out, "transition") {
		if !strings.Contains(out, "deliberation gate") {
			t.Fatalf("help missing transition enforce docs: %q", out)
		}
	}
}

func TestSeedExportStrictEnforceNoWriteOnViolation(t *testing.T) {
	dir := t.TempDir()
	setupVerificationDebtFixture(t, dir)
	outPath := filepath.Join(dir, "out.json")
	if _, err := os.Stat(outPath); err == nil {
		t.Fatal("out.json should not exist yet")
	}
	code, stdout, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(outPath); err == nil {
		t.Fatal("out.json must not be written on enforce violation")
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("enforce block must not write stdout payload: %q", stdout)
	}
}

func TestSeedExportStrictWithoutEnforceExitZero(t *testing.T) {
	dir := t.TempDir()
	setupVerificationDebtFixture(t, dir)
	outPath := filepath.Join(dir, "out.json")
	code, _, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict",
	})
	if code != exitOK {
		t.Fatalf("strict without enforce want exitOK got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "seed export strict:") {
		t.Fatalf("stderr missing violation line: %q", stderr)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("strict without enforce should still write file")
	}
}

func TestSeedExportStrictEnforceBlocksOpenRegression(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Reg task")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	evalOut, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "ffffffff-ffff-4fff-8fff-ffffffffffff",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []any{}, "plan_changes": []any{}, "spawned_tasks": []any{},
			"regressions": []map[string]any{{
				"source_kind": "evaluation", "source_id": evalOut.ID, "task_id": taskID,
			}},
		},
	}
	writeJSON(t, filepath.Join(dir, "reg.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "reg.json")}); code != exitOK {
		t.Fatalf("apply regression: %d", code)
	}

	outPath := filepath.Join(dir, "reg-out.json")
	code, stdout, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(outPath); err == nil {
		t.Fatal("out file must not be written on regression enforce block")
	}
	if !strings.Contains(stderr, "open_regression") && !strings.Contains(stderr, "regression") {
		t.Fatalf("stderr should mention regression: %q", stderr)
	}
}

func p26ExportSnippetPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "testdata", "p26-export-snippet.json")
}

func setupP26ThinGraphFixture(t *testing.T, dir string) {
	t.Helper()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := p26ExportSnippetPath(t)
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed import P26 snippet: %d", code)
	}
}

// TestSeedExportStrictBlockingOrphanDiscoverySingleHonestyViolation locks R4:
// a BLOCKING orphan discovery must produce exactly one honesty line for that ID
// (document CollectSeedDocumentHonestyViolations only — no store-backed duplicate).
func TestSeedExportStrictBlockingOrphanDiscoverySingleHonestyViolation(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	_ = addTaskForLoopTest(t, dir, goalID, "Task")

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	disc, err := svc.CreateDiscovery(context.Background(), domain.DiscoveryInput{
		Title:    "Blocking orphan",
		Severity: domain.SeverityBlocking,
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	discID := disc.ID
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	code, _, stderr := runCapture(t, []string{"-C", dir, "seed", "export", "--strict"})
	if code != exitOK {
		t.Fatalf("strict without enforce want exitOK got %d stderr=%q", code, stderr)
	}
	n := 0
	for _, line := range strings.Split(stderr, "\n") {
		if strings.Contains(line, discID) && strings.Contains(line, "discovery_mentions_task") {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("want exactly 1 honesty violation mentioning %s, got %d; stderr=%q", discID, n, stderr)
	}
	if strings.Contains(stderr, "BLOCKING discovery") {
		t.Fatalf("store-backed BLOCKING orphan message must be gone: %q", stderr)
	}
}

func TestSeedExportStrictEnforceBlocksP26ThinGraph(t *testing.T) {
	dir := t.TempDir()
	setupP26ThinGraphFixture(t, dir)
	outPath := filepath.Join(dir, "thin-out.json")
	code, stdout, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict", "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(outPath); err == nil {
		t.Fatal("thin graph enforce must not write output file")
	}
	if !strings.Contains(stderr, "graph honesty") {
		t.Fatalf("stderr missing graph honesty violation: %q", stderr)
	}
	if !strings.Contains(stderr, "discoveries=0") || !strings.Contains(stderr, "decisions=0") {
		t.Fatalf("stderr missing thin-graph counts: %q", stderr)
	}
}

func TestSeedExportPlainThinGraphEarlyWarnWrites(t *testing.T) {
	dir := t.TempDir()
	setupP26ThinGraphFixture(t, dir)
	outPath := filepath.Join(dir, "thin-early.json")
	code, _, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath,
	})
	if code != exitOK {
		t.Fatalf("plain export want exitOK got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "seed export: warn:") || !strings.Contains(stderr, "thin graph") {
		t.Fatalf("stderr missing early thin-graph warn: %q", stderr)
	}
	if !strings.Contains(stderr, "write discoveries/decisions before") || !strings.Contains(stderr, "--strict --enforce") {
		t.Fatalf("stderr missing write-before-export nudge: %q", stderr)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("plain export must still write file (warn-only)")
	}
}

func TestSeedExportStrictThinGraphWarnOnlyWrites(t *testing.T) {
	dir := t.TempDir()
	setupP26ThinGraphFixture(t, dir)
	outPath := filepath.Join(dir, "thin-warn.json")
	code, _, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict",
	})
	if code != exitOK {
		t.Fatalf("strict without enforce want exitOK got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "seed export strict:") || !strings.Contains(stderr, "graph honesty") {
		t.Fatalf("stderr missing honesty violation line: %q", stderr)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatal("strict without enforce should still write file")
	}
}

func TestSeedExportStrictCleanAllowsWrite(t *testing.T) {
	dir := t.TempDir()
	setupCleanFullCycleFixture(t, dir)
	outPath := filepath.Join(dir, "clean.json")
	code, _, _ := runCapture(t, []string{
		"-C", dir, "seed", "export", "-o", outPath, "--strict", "--enforce",
	})
	if code != exitOK {
		t.Fatalf("clean export want exitOK got %d", code)
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if doc["version"] != float64(1) {
		t.Fatalf("version want 1: %#v", doc["version"])
	}
}

func TestSeedExportWithoutStrictUnchanged(t *testing.T) {
	dir := t.TempDir()
	setupVerificationDebtFixture(t, dir)
	code, stdout, _ := runCapture(t, []string{"-C", dir, "seed", "export"})
	if code != exitOK {
		t.Fatalf("plain export want exitOK got %d", code)
	}
	if !strings.Contains(stdout, `"version"`) {
		t.Fatalf("expected seed JSON on stdout: %q", stdout)
	}
}

func TestSeedExportEnforceRequiresStrict(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	code, _, stderr := runCapture(t, []string{"-C", dir, "seed", "export", "--enforce"})
	if code != exitUsage {
		t.Fatalf("enforce without strict want exitUsage got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "--enforce requires --strict") {
		t.Fatalf("stderr missing requirement hint: %q", stderr)
	}
}

func TestSeedExportStrictTaskFilter(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	debtGoal := addGoalForLoopTest(t, dir, "Debt goal")
	debtID := addTaskForLoopTest(t, dir, debtGoal, "Debt")
	createCurrentDeepPlanForLoopTest(t, dir, debtGoal)
	markPlanCritiquedForLoopTest(t, dir, debtID, debtGoal)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	if _, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: debtID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	cleanGoal, cleanID := setupCleanFullCycleFixture(t, dir)

	code, _, _ := runCapture(t, []string{
		"-C", dir, "seed", "export", "--strict", "--task", debtID, "--enforce",
	})
	if code != exitGateBlocked {
		t.Fatalf("debt task filter want blocked got %d", code)
	}

	code, _, stderr := runCapture(t, []string{
		"-C", dir, "seed", "export", "--strict", "--task", cleanID, "--enforce",
	})
	if code != exitOK {
		t.Fatalf("clean task filter want exitOK got %d stderr=%q (cleanGoal=%s debtGoal=%s)", code, stderr, cleanGoal, debtGoal)
	}
}

func TestHelpIncludesSeedExportStrict(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "--strict") || !strings.Contains(out, "--enforce") {
		t.Fatalf("help missing seed export strict/enforce: %q", out)
	}
	if !strings.Contains(out, "seed export") {
		t.Fatalf("help missing seed export block: %q", out)
	}
}
