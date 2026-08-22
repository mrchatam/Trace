package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestLoopNextPacketShape(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	goalID := addGoalForLoopTest(t, dir, "Loop goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "alpha")
	scopeID := createCurrentDeepPlanForLoopTest(t, dir, goalID)
	if scopeID == "" {
		t.Fatal("scope id empty")
	}

	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("import { alpha } from './a.js'\nexport function beta() { return alpha() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatalf("loop next json: %v\n%s", err, out)
	}
	if pkt["schema_version"] != "trace.loop.next.v1" {
		t.Fatalf("schema_version: %v", pkt["schema_version"])
	}
	seed, ok := pkt["seed"].(map[string]any)
	if !ok || seed["task_id"] != taskID || seed["goal_id"] != goalID {
		t.Fatalf("seed: %#v", pkt["seed"])
	}
	if seed["work_state"] != "PENDING" {
		t.Fatalf("seed work_state: %#v", seed)
	}

	tasks, ok := pkt["tasks"].(map[string]any)
	if !ok || tasks["goal_id"] != goalID {
		t.Fatalf("tasks section: %#v", pkt["tasks"])
	}
	taskItems, ok := tasks["items"].([]any)
	if !ok || len(taskItems) != 1 {
		t.Fatalf("tasks items: %#v", tasks["items"])
	}

	plan, ok := pkt["plan"].(map[string]any)
	if !ok {
		t.Fatalf("plan section missing: %#v", pkt["plan"])
	}
	planSnap, ok := plan["snapshot"].(map[string]any)
	if !ok || planSnap["goal_id"] != goalID {
		t.Fatalf("plan snapshot: %#v", plan["snapshot"])
	}
	if planSnap["current_scope_id"] == nil {
		t.Fatalf("current_scope_id missing: %#v", planSnap)
	}
	if planSnap["current_deep_plan"] == nil {
		t.Fatalf("current_deep_plan missing: %#v", planSnap)
	}

	why, ok := pkt["why"].(map[string]any)
	if !ok {
		t.Fatalf("why section: %#v", pkt["why"])
	}
	whySnap, ok := why["snapshot"].(map[string]any)
	if !ok || whySnap["seed_id"] != taskID {
		t.Fatalf("why snapshot: %#v", why["snapshot"])
	}

	ctx, ok := pkt["context"].(map[string]any)
	if !ok {
		t.Fatalf("context section: %#v", pkt["context"])
	}
	ctxSnap, ok := ctx["snapshot"].(map[string]any)
	if !ok || ctxSnap["task_id"] != taskID {
		t.Fatalf("context snapshot: %#v", ctx["snapshot"])
	}

	related, ok := pkt["related"].(map[string]any)
	if !ok {
		t.Fatalf("related section: %#v", pkt["related"])
	}
	if related["available"] != true {
		t.Fatalf("related available want true: %#v", related)
	}
	seeds, ok := related["seeds"].([]any)
	if !ok || len(seeds) == 0 {
		t.Fatalf("related seeds: %#v", related["seeds"])
	}
	relatedSnap, ok := related["snapshot"].(map[string]any)
	if !ok {
		t.Fatalf("related snapshot: %#v", related["snapshot"])
	}
	if relatedSnap["depth"] != float64(2) {
		t.Fatalf("related depth: %#v", relatedSnap["depth"])
	}

	loopHints, ok := pkt["loop_hints"].(map[string]any)
	if !ok || loopHints["available"] != false {
		t.Fatalf("loop_hints: %#v", pkt["loop_hints"])
	}
	if !strings.Contains(loopHints["unavailable_reason"].(string), "iteration metadata") {
		t.Fatalf("loop_hints reason: %#v", loopHints)
	}
}

func TestLoopNextMissingTaskFailsClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	code, _, stderr := runCapture(t, []string{"-C", dir, "loop", "next", "--task", "00000000-0000-0000-0000-000000000000"})
	if code != exitFail {
		t.Fatalf("want exitFail got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "load seed task") {
		t.Fatalf("stderr: %q", stderr)
	}
}

func TestLoopNextMissingGoalIDFailsClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskID := addTaskForLoopTest(t, dir, "", "ungoled")
	code, _, stderr := runCapture(t, []string{"-C", dir, "loop", "next", "--task", taskID})
	if code != exitFail {
		t.Fatalf("want exitFail got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "has no goal_id") {
		t.Fatalf("stderr: %q", stderr)
	}
}

func TestLoopNextMissingPlanContextFailsClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal without plan")
	taskID := addTaskForLoopTest(t, dir, goalID, "task without plan")
	code, _, stderr := runCapture(t, []string{"-C", dir, "loop", "next", "--task", taskID})
	if code != exitFail {
		t.Fatalf("want exitFail got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "missing goal plan context") {
		t.Fatalf("stderr: %q", stderr)
	}
}

func TestHelpIncludesLoopNext(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "loop next --task <id>") {
		t.Fatalf("help missing loop next: %q", out)
	}
	if !strings.Contains(out, "loop apply [--in <path>]") {
		t.Fatalf("help missing loop apply: %q", out)
	}
	if !strings.Contains(out, "loop status --task <id>") {
		t.Fatalf("help missing loop status: %q", out)
	}
}

func TestHelpIncludesLoopGate(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "loop gate --task") {
		t.Fatalf("help missing loop gate: %q", out)
	}
	if !strings.Contains(out, "Exit 0 allowed, 1 blocked, 2 usage") {
		t.Fatalf("help missing gate exit-code hint: %q", out)
	}
}

func TestLoopGateAllowedExitZero(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	if env["allowed"] != true {
		t.Fatalf("allowed: %#v", env["allowed"])
	}
	violations, ok := env["violations"].([]any)
	if !ok || len(violations) != 0 {
		t.Fatalf("violations: %#v", env["violations"])
	}
}

func TestLoopGateBlockedExitOne(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "blocker", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if code != exitGateBlocked {
		t.Fatalf("want exitGateBlocked (1) got %d stdout=%q", code, stdout)
	}
	if code == exitFail {
		t.Fatalf("blocked must not use exitFail (2); got %d", code)
	}
	env := parseGateOutput(t, stdout)
	if env["allowed"] != false {
		t.Fatalf("allowed: %#v", env["allowed"])
	}
	violations := env["violations"].([]any)
	v0 := violations[0].(map[string]any)
	if v0["code"] != "premature_implementation" {
		t.Fatalf("violation code: %#v", v0["code"])
	}
}

func TestLoopGateJSONSchemaVersion(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if code != exitGateBlocked {
		t.Fatalf("want blocked policy path got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	if env["schema_version"] != gateSchemaVersion {
		t.Fatalf("schema_version: %#v", env["schema_version"])
	}
}

func TestLoopGateTopLevelLiftFromViolation(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	violations := env["violations"].([]any)
	v0 := violations[0].(map[string]any)
	if env["recommended_phase"] != v0["recommended_phase"] {
		t.Fatalf("top-level recommended_phase=%#v violation=%#v", env["recommended_phase"], v0["recommended_phase"])
	}
	if env["reason_code"] != v0["reason_code"] {
		t.Fatalf("top-level reason_code=%#v violation=%#v", env["reason_code"], v0["reason_code"])
	}
}

func TestLoopGateAllowedEmptyViolations(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	})
	env := parseGateOutput(t, stdout)
	if env["allowed"] != true {
		t.Fatalf("allowed: %#v", env["allowed"])
	}
	violations, ok := env["violations"].([]any)
	if !ok {
		t.Fatalf("violations not array: %#v", env["violations"])
	}
	if len(violations) != 0 {
		t.Fatalf("violations want empty array: %#v", violations)
	}
}

func TestLoopGateBlockedStderrHint(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	_, _, stderr := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if stderr == "" {
		t.Fatal("stderr empty for blocked gate")
	}
	if !strings.Contains(stderr, "edit blocked") && !strings.Contains(stderr, "plan") {
		t.Fatalf("stderr missing violation hint: %q", stderr)
	}
}

func TestLoopGateDefaultForEdit(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	codeDefault, stdoutDefault, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID})
	codeExplicit, stdoutExplicit, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	if codeDefault != codeExplicit {
		t.Fatalf("exit codes differ: default=%d explicit=%d", codeDefault, codeExplicit)
	}
	envDefault := parseGateOutput(t, stdoutDefault)
	envExplicit := parseGateOutput(t, stdoutExplicit)
	if envDefault["allowed"] != envExplicit["allowed"] {
		t.Fatalf("allowed differ: default=%#v explicit=%#v", envDefault["allowed"], envExplicit["allowed"])
	}
	if envDefault["for"] != "edit" {
		t.Fatalf("default for: %#v", envDefault["for"])
	}
}

func TestLoopGateInvalidForFailClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "not-a-gate"})
	if code != exitFail {
		t.Fatalf("want exitFail (2) got %d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if strings.Contains(stdout, gateSchemaVersion) {
		t.Fatalf("stdout must not emit gate JSON on usage error: %q", stdout)
	}
}

func TestLoopGateMissingTaskFlag(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "gate"})
	if code != exitFail {
		t.Fatalf("want exitFail (2) got %d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "usage:") {
		t.Fatalf("stderr missing usage hint: %q", stderr)
	}
}

func TestLoopGateUnknownTaskOrientBlocked(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskID := "00000000-0000-0000-0000-000000000099"
	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "orient"})
	if code != exitGateBlocked {
		t.Fatalf("want exitGateBlocked got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	violations := env["violations"].([]any)
	v0 := violations[0].(map[string]any)
	if v0["code"] != "gate_orient_failed" {
		t.Fatalf("code: %#v", v0["code"])
	}
	if env["reason_code"] != "task_not_found" {
		t.Fatalf("reason_code: %#v", env["reason_code"])
	}
}

func TestLoopGateOrientAllowed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "orient"})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	if env["allowed"] != true {
		t.Fatalf("allowed: %#v", env["allowed"])
	}
}

func TestLoopGateDoneBlockedVerificationDebt(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
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
	_ = st.Close()

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "done"})
	if code != exitGateBlocked {
		t.Fatalf("want blocked got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	if env["reason_code"] != "verification_incomplete" {
		t.Fatalf("reason_code: %#v", env["reason_code"])
	}
}

func TestLoopGateExecuteAllowedWhenPending(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	code, stdout, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "execute"})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stdout=%q", code, stdout)
	}
	env := parseGateOutput(t, stdout)
	if env["allowed"] != true {
		t.Fatalf("allowed: %#v", env["allowed"])
	}
}

func markPlanCritiquedForLoopTest(t *testing.T, dir, taskID, goalID string) {
	t.Helper()
	pcID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []any{}, "spawned_tasks": []any{},
			"plan_changes": []map[string]any{{"id": pcID, "title": "critique done"}},
		},
	}
	writeJSON(t, filepath.Join(dir, "critique.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "critique.json")}); code != exitOK {
		t.Fatalf("apply plan change: %d", code)
	}
}

func parseGateOutput(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("gate json: %v\n%s", err, stdout)
	}
	return env
}

func writeTraceConfig(t *testing.T, root, content string) {
	t.Helper()
	dir := filepath.Join(root, ".trace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func addBlockingUncertaintyForLoopTest(t *testing.T, dir, taskID string) {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	if _, err := svc.CreateUncertainty(context.Background(), domain.UncertaintyInput{
		Title: "blocker", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
}

func setupBlockedLoopStatusFixture(t *testing.T) (dir, taskID string) {
	t.Helper()
	dir = t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID = addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	addBlockingUncertaintyForLoopTest(t, dir, taskID)
	return dir, taskID
}

func parseStatusOutput(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var status map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, stdout)
	}
	return status
}

func TestLoopStatusIncludesViolationsWhenBlocked(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	violations, ok := status["violations"].([]any)
	if !ok || len(violations) < 1 {
		t.Fatalf("violations: %#v", status["violations"])
	}
	v0 := violations[0].(map[string]any)
	if v0["code"] != "premature_implementation" {
		t.Fatalf("violation code: %#v", v0["code"])
	}
	if v0["for"] != "edit" {
		t.Fatalf("violation for: %#v", v0["for"])
	}
}

func TestLoopStatusViolationsEmptyWhenClean(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)
	markPlanCritiquedForLoopTest(t, dir, taskID, goalID)

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	violations, ok := status["violations"].([]any)
	if !ok {
		t.Fatalf("violations not array: %#v", status["violations"])
	}
	if len(violations) != 0 {
		t.Fatalf("violations want empty: %#v", violations)
	}
}

func TestLoopStatusViolationsMatchGateEdit(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)

	statusOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, statusOut)
	statusViolations := status["violations"]

	_, gateOut, _ := runCapture(t, []string{"-C", dir, "loop", "gate", "--task", taskID, "--for", "edit"})
	gate := parseGateOutput(t, gateOut)
	gateViolations := gate["violations"]

	statusRaw, err := json.Marshal(statusViolations)
	if err != nil {
		t.Fatal(err)
	}
	gateRaw, err := json.Marshal(gateViolations)
	if err != nil {
		t.Fatal(err)
	}
	if string(statusRaw) != string(gateRaw) {
		t.Fatalf("violations mismatch:\nstatus=%s\ngate=%s", statusRaw, gateRaw)
	}
}

func TestLoopStatusViolationsAlwaysArray(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	violations, ok := status["violations"].([]any)
	if !ok {
		t.Fatalf("violations key missing or not array: %#v", status["violations"])
	}
	_ = violations
}

func TestLoopStatusSchemaVersionUnchanged(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	if status["schema_version"] != "trace.loop.status.v1" {
		t.Fatalf("schema_version: %#v", status["schema_version"])
	}
}

func TestLoopStatus_IncludesGoalStructureAdvisory(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Mega goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	for i := 0; i < 16; i++ {
		addTaskForLoopTest(t, dir, goalID, fmt.Sprintf("task-%d", i))
	}

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	assertAdvisoryCode(t, status, "goal_structure_warning")
}

func TestLoopStatus_BootstrapRecommendedAdvisory(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	applyLinkedPlanChangeForLoopTest(t, dir, taskID, goalID)

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	assertAdvisoryCode(t, status, "bootstrap_recommended")
}

func TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	applyLinkedPlanChangeForLoopTest(t, dir, taskID, goalID)

	stdout := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	status := parseStatusOutput(t, stdout)
	assertAdvisoryCode(t, status, "bootstrap_recommended")
	delib, ok := status["deliberation"].(map[string]any)
	if !ok {
		t.Fatalf("deliberation: %#v", status["deliberation"])
	}
	policy, ok := delib["policy_inputs"].(map[string]any)
	if !ok {
		t.Fatalf("policy_inputs: %#v", delib["policy_inputs"])
	}
	if policy["plan_exists"] != false {
		t.Fatalf("plan_exists must stay false: %#v", policy["plan_exists"])
	}
}

func assertAdvisoryCode(t *testing.T, status map[string]any, code string) {
	t.Helper()
	advs, ok := status["advisories"].([]any)
	if !ok {
		t.Fatalf("advisories not array: %#v", status["advisories"])
	}
	for _, raw := range advs {
		a, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("advisory entry: %#v", raw)
		}
		if a["code"] == code {
			return
		}
	}
	t.Fatalf("advisories missing code %q: %#v", code, advs)
}

func applyLinkedPlanChangeForLoopTest(t *testing.T, dir, taskID, goalID string) {
	t.Helper()
	discID := "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
	pcID := "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee"
	applyID := "ffffffff-ffff-4fff-8fff-ffffffffffff"
	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       applyID,
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []map[string]any{{
				"id": discID, "title": "disc",
				"links": []map[string]any{{
					"rel": "discovery_mentions_task", "to_type": "task", "to_id": taskID,
				}},
			}},
			"plan_changes": []map[string]any{{
				"id": pcID, "title": "linked pc", "discovery_id": discID,
			}},
			"spawned_tasks": []any{},
		},
	}
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	if code, _, stderr := runLoopApplyWithStdin(t, dir, string(raw)); code != exitOK {
		t.Fatalf("loop apply plan change: %d stderr=%q", code, stderr)
	}
}

func TestTraceConfigEnforceDefaultOff(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)

	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stderr=%q", code, stderr)
	}
	status := parseStatusOutput(t, stdout)
	violations, ok := status["violations"].([]any)
	if !ok || len(violations) == 0 {
		t.Fatalf("stdout violations: %#v", status["violations"])
	}
	if stderr != "" {
		t.Fatalf("stderr want empty got %q", stderr)
	}
}

func TestTraceConfigEnforceMalformedFailClosedOff(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)
	writeTraceConfig(t, dir, "{not json")

	code, _, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d", code)
	}
	if stderr != "" {
		t.Fatalf("stderr want empty (fail-closed off) got %q", stderr)
	}
}

func TestTraceConfigEnforceInvalidValueFailClosedOff(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)
	writeTraceConfig(t, dir, `{"enforce":"loud"}`)

	code, _, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d", code)
	}
	if stderr != "" {
		t.Fatalf("stderr want empty got %q", stderr)
	}
}

func TestTraceConfigEnforceWarnSurfacesStderr(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)
	writeTraceConfig(t, dir, `{"enforce":"warn"}`)

	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stderr=%q", code, stderr)
	}
	status := parseStatusOutput(t, stdout)
	violations := status["violations"].([]any)
	msg := violations[0].(map[string]any)["message"].(string)
	if !strings.Contains(stderr, msg) {
		t.Fatalf("stderr missing violation message %q: %q", msg, stderr)
	}
	if !strings.Contains(stderr, "loop status:") {
		t.Fatalf("stderr missing prefix: %q", stderr)
	}
}

func TestTraceConfigEnforceStrictSurfacesStderr(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)
	writeTraceConfig(t, dir, `{"enforce":"strict"}`)

	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d stderr=%q", code, stderr)
	}
	status := parseStatusOutput(t, stdout)
	violations := status["violations"].([]any)
	msg := violations[0].(map[string]any)["message"].(string)
	if !strings.Contains(stderr, msg) {
		t.Fatalf("stderr missing violation message %q: %q", msg, stderr)
	}
}

func TestTraceConfigEnforceOffNoStderrOnViolation(t *testing.T) {
	dir, taskID := setupBlockedLoopStatusFixture(t)
	writeTraceConfig(t, dir, `{"enforce":"off"}`)

	code, stdout, stderr := runCapture(t, []string{"-C", dir, "loop", "status", "--task", taskID})
	if code != exitOK {
		t.Fatalf("want exitOK got %d", code)
	}
	status := parseStatusOutput(t, stdout)
	violations, ok := status["violations"].([]any)
	if !ok || len(violations) == 0 {
		t.Fatalf("stdout violations: %#v", status["violations"])
	}
	if stderr != "" {
		t.Fatalf("stderr want empty got %q", stderr)
	}
}

func TestHelpIncludesTraceConfig(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, ".trace/config.json") {
		t.Fatalf("help missing config path: %q", out)
	}
	for _, mode := range []string{"off", "warn", "strict"} {
		if !strings.Contains(out, mode) {
			t.Fatalf("help missing enforce mode %q", mode)
		}
	}
}

func TestLoopApplyMalformedInputFailsClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed task")

	missingApplyID := `{"schema_version":"trace.loop.apply.v1","seed":{"task_id":"` + taskID + `","goal_id":"` + goalID + `"},"writes":{"discoveries":[],"plan_changes":[],"spawned_tasks":[]}}`
	code, _, stderr := runLoopApplyWithStdin(t, dir, missingApplyID)
	if code != exitFail || !strings.Contains(stderr, "missing required field \"apply_id\"") {
		t.Fatalf("missing apply_id expected fail-closed, got code=%d stderr=%q", code, stderr)
	}

	badSchema := `{"schema_version":"trace.loop.apply.v0","apply_id":"11111111-1111-4111-8111-111111111111","seed":{"task_id":"` + taskID + `","goal_id":"` + goalID + `"},"writes":{"discoveries":[],"plan_changes":[],"spawned_tasks":[]}}`
	code, _, stderr = runLoopApplyWithStdin(t, dir, badSchema)
	if code != exitFail || !strings.Contains(stderr, "schema_version must be") {
		t.Fatalf("schema mismatch expected fail-closed, got code=%d stderr=%q", code, stderr)
	}

	missingItemID := `{"schema_version":"trace.loop.apply.v1","apply_id":"22222222-2222-4222-8222-222222222222","seed":{"task_id":"` + taskID + `","goal_id":"` + goalID + `"},"writes":{"discoveries":[{"title":"missing id"}],"plan_changes":[],"spawned_tasks":[]}}`
	code, _, stderr = runLoopApplyWithStdin(t, dir, missingItemID)
	if code != exitFail || !strings.Contains(stderr, "writes.discoveries[0].id is required") {
		t.Fatalf("missing item id expected fail-closed, got code=%d stderr=%q", code, stderr)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	discoveries, err := st.ListDiscoveries()
	if err != nil {
		t.Fatal(err)
	}
	if len(discoveries) != 0 {
		t.Fatalf("malformed apply should not write discoveries: %+v", discoveries)
	}
}

func TestLoopApplyReplayAndStatusFlow(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed task")

	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "33333333-3333-4333-8333-333333333333",
		"seed": map[string]any{
			"task_id": taskID,
			"goal_id": goalID,
		},
		"writes": map[string]any{
			"discoveries": []map[string]any{
				{
					"id":       "44444444-4444-4444-8444-444444444444",
					"title":    "Found gap",
					"severity": "PLAN_AFFECTING",
					"links": []map[string]any{
						{"rel": domain.RelDiscoveryMentionsTask, "to_type": domain.EntityTask, "to_id": taskID},
					},
				},
			},
			"plan_changes": []map[string]any{
				{
					"id":           "55555555-5555-4555-8555-555555555555",
					"title":        "Adjust plan",
					"discovery_id": "44444444-4444-4444-8444-444444444444",
				},
			},
			"spawned_tasks": []map[string]any{
				{
					"id":    "66666666-6666-4666-8666-666666666666",
					"title": "Follow-up",
				},
			},
		},
	}
	path := filepath.Join(dir, "apply.json")
	writeJSON(t, path, env)
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "apply", "--in", path})
	})
	var first map[string]any
	if err := json.Unmarshal([]byte(out), &first); err != nil {
		t.Fatalf("first apply json: %v (%s)", err, out)
	}
	if first["replay"] != false || first["new_spawned_tasks"] != float64(1) || first["new_plan_changes"] != float64(1) {
		t.Fatalf("first apply result: %#v", first)
	}

	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "apply", "--in", path})
	})
	var replay map[string]any
	if err := json.Unmarshal([]byte(out), &replay); err != nil {
		t.Fatalf("replay apply json: %v (%s)", err, out)
	}
	if replay["replay"] != true {
		t.Fatalf("replay result: %#v", replay)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	tasks, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Fatalf("replay must keep task count stable, got %d", len(tasks))
	}
	planChanges, err := st.ListPlanChanges()
	if err != nil {
		t.Fatal(err)
	}
	if len(planChanges) != 1 {
		t.Fatalf("replay must keep plan_changes stable, got %d", len(planChanges))
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, "44444444-4444-4444-8444-444444444444")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("replay must keep discovery links stable, got %d", len(links))
	}
	_ = st.Close()

	statusOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	var status map[string]any
	if err := json.Unmarshal([]byte(statusOut), &status); err != nil {
		t.Fatalf("status json: %v (%s)", err, statusOut)
	}
	if status["schema_version"] != "trace.loop.status.v1" {
		t.Fatalf("status schema: %#v", status)
	}
	if status["saturated"] != false || status["reason"] == "tasks_and_plan_unchanged" {
		t.Fatalf("status after non-zero delta should be unsaturated: %#v", status)
	}
}

func TestLoopStatusInsufficientHistory(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed task")

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	var status map[string]any
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status json: %v (%s)", err, out)
	}
	if status["saturated"] != false || status["reason"] != "insufficient_history" {
		t.Fatalf("expected insufficient_history, got %#v", status)
	}
}

func TestLoopStatusSaturatedByZeroDeltaAndMaxIteration(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed task")

	zeroDelta := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "77777777-7777-4777-8777-777777777777",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries":   []any{},
			"plan_changes":  []any{},
			"spawned_tasks": []any{},
		},
	}
	pathA := filepath.Join(dir, "apply-zero.json")
	writeJSON(t, pathA, zeroDelta)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", pathA}); code != exitOK {
		t.Fatalf("apply zero-delta: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	var status map[string]any
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status json: %v (%s)", err, out)
	}
	if status["saturated"] != false {
		t.Fatalf("first zero-delta must not saturate: %#v", status)
	}

	zeroDelta2 := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "77777777-7777-4777-8777-777777777778",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries":   []any{},
			"plan_changes":  []any{},
			"spawned_tasks": []any{},
		},
	}
	pathSecond := filepath.Join(dir, "apply-zero-2.json")
	writeJSON(t, pathSecond, zeroDelta2)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", pathSecond}); code != exitOK {
		t.Fatalf("apply second zero-delta: %d", code)
	}

	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status json: %v (%s)", err, out)
	}
	if status["saturated"] != true || status["reason"] != "tasks_and_plan_unchanged" {
		t.Fatalf("second zero-delta should saturate: %#v", status)
	}

	maxed := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "88888888-8888-4888-8888-888888888888",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries":   []any{},
			"plan_changes":  []any{},
			"spawned_tasks": []any{map[string]any{"id": "99999999-9999-4999-8999-999999999999", "title": "new task"}},
			"stop":          map[string]any{"reason": "budget", "max_iterations_reached": true},
		},
	}
	pathMax := filepath.Join(dir, "apply-max.json")
	writeJSON(t, pathMax, maxed)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", pathMax}); code != exitOK {
		t.Fatalf("apply max-iterations: %d", code)
	}

	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status json: %v (%s)", err, out)
	}
	if status["saturated"] != true || status["reason"] != "max_iterations_reached" {
		t.Fatalf("max-iterations should saturate: %#v", status)
	}
}

func TestLoopResetCLIClearsStop(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed task")

	for i, applyID := range []string{
		"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1",
		"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2",
	} {
		env := map[string]any{
			"schema_version": "trace.loop.apply.v1",
			"apply_id":       applyID,
			"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
			"writes": map[string]any{
				"discoveries":   []any{},
				"plan_changes":  []any{},
				"spawned_tasks": []any{},
			},
		}
		path := filepath.Join(dir, fmt.Sprintf("empty-%d.json", i))
		writeJSON(t, path, env)
		if code := run([]string{"-C", dir, "loop", "apply", "--in", path}); code != exitOK {
			t.Fatalf("apply %d: %d", i, code)
		}
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "reset", "--task", taskID})
	})
	var state map[string]any
	if err := json.Unmarshal([]byte(out), &state); err != nil {
		t.Fatalf("reset json: %v (%s)", err, out)
	}
	if stopped, _ := state["Stopped"].(bool); stopped {
		t.Fatalf("reset should clear Stopped: %#v", state)
	}
	if hop, _ := state["HopCount"].(float64); hop != 0 {
		t.Fatalf("HopCount=%v want 0", hop)
	}
	if phase, _ := state["CurrentPhase"].(string); phase != "EXECUTE" {
		t.Fatalf("CurrentPhase=%v want EXECUTE", phase)
	}
	if consec, _ := state["ConsecutiveEmptyApplies"].(float64); consec != 0 {
		t.Fatalf("ConsecutiveEmptyApplies=%v want 0", consec)
	}

	help := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "help"})
	})
	if !strings.Contains(help, "reset --task") {
		t.Fatalf("loop help missing reset: %q", help)
	}
}

func addGoalForLoopTest(t *testing.T, dir, title string) string {
	t.Helper()
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", title})
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("goal json: %v (%s)", err, out)
	}
	return res["id"].(string)
}

func addTaskForLoopTest(t *testing.T, dir, goalID, title string) string {
	t.Helper()
	args := []string{"-C", dir, "add", "task", "--title", title}
	if goalID != "" {
		args = append(args, "--goal-id", goalID)
	}
	out := captureStdout(t, func() int {
		return run(args)
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("task json: %v (%s)", err, out)
	}
	return res["id"].(string)
}

func createCurrentDeepPlanForLoopTest(t *testing.T, dir, goalID string) string {
	t.Helper()
	out := captureStdout(t, func() int {
		return run([]string{
			"-C", dir, "plan", "create-coarse", "--goal", goalID,
			"--phase", "Phase 1", "--scope", "Scope 1",
		})
	})
	var plan map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &plan); err != nil {
		t.Fatalf("plan create json: %v (%s)", err, out)
	}
	phases, ok := plan["phases"].([]any)
	if !ok || len(phases) != 1 {
		t.Fatalf("phases: %#v", plan["phases"])
	}
	phase := phases[0].(map[string]any)
	scopes, ok := phase["scopes"].([]any)
	if !ok || len(scopes) != 1 {
		t.Fatalf("scopes: %#v", phase["scopes"])
	}
	scopeID := scopes[0].(map[string]any)["id"].(string)
	if code := run([]string{"-C", dir, "plan", "set-current", "--goal", goalID, "--scope", scopeID}); code != exitOK {
		t.Fatalf("plan set-current: %d", code)
	}
	if code := run([]string{
		"-C", dir, "plan", "deep", "--scope", scopeID, "--exit", "packet ready", "--work", "emit loop packet",
	}); code != exitOK {
		t.Fatalf("plan deep: %d", code)
	}
	return scopeID
}

func runLoopApplyWithStdin(t *testing.T, dir, body string) (code int, stdout, stderr string) {
	t.Helper()
	rIn, wIn, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := wIn.Write([]byte(body)); err != nil {
		t.Fatal(err)
	}
	_ = wIn.Close()
	oldIn := os.Stdin
	os.Stdin = rIn
	defer func() {
		os.Stdin = oldIn
		_ = rIn.Close()
	}()
	return runCapture(t, []string{"-C", dir, "loop", "apply"})
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoopNextExecuteWhenPendingLive(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	pcID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []any{}, "spawned_tasks": []any{},
			"plan_changes": []map[string]any{{"id": pcID, "title": "critique done"}},
		},
	}
	writeJSON(t, filepath.Join(dir, "critique.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "critique.json")}); code != exitOK {
		t.Fatalf("apply plan change: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	delib, ok := pkt["deliberation"].(map[string]any)
	if !ok {
		t.Fatalf("deliberation missing: %#v", pkt["deliberation"])
	}
	inputs, ok := delib["policy_inputs"].(map[string]any)
	if !ok {
		t.Fatalf("policy_inputs missing: %#v", delib)
	}
	if inputs["execute_pending"] != true {
		t.Fatalf("execute_pending want true: %#v", inputs)
	}
	if delib["phase"] != "EXECUTE" {
		t.Fatalf("phase want EXECUTE: %#v", delib)
	}
}

func TestLoopNextDeliberationSectionPresent(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatalf("json: %v", err)
	}
	delib, ok := pkt["deliberation"].(map[string]any)
	if !ok {
		t.Fatalf("deliberation section missing: %#v", pkt["deliberation"])
	}
	if delib["phase"] == nil || delib["why_selected"] == nil {
		t.Fatalf("deliberation phase/why missing: %#v", delib)
	}
	inputs, ok := delib["policy_inputs"].(map[string]any)
	if !ok {
		t.Fatalf("policy_inputs missing: %#v", delib)
	}
	if _, ok := inputs["blocking_uncertainty_count"]; !ok {
		t.Fatalf("policy_inputs incomplete: %#v", inputs)
	}
}

func TestLoopNextPolicyInputsLiveQueries(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "blocking?", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	delib := pkt["deliberation"].(map[string]any)
	inputs := delib["policy_inputs"].(map[string]any)
	if inputs["blocking_uncertainty_count"].(float64) != 1 {
		t.Fatalf("blocking count: %#v", inputs["blocking_uncertainty_count"])
	}
	if delib["phase"] != "INVESTIGATE" {
		t.Fatalf("phase: %#v", delib["phase"])
	}
}

func TestLoopNextInvestigateNoRetrievalStderr(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "blocking?", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	stdout, stderr := captureStdoutStderr(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	if strings.Contains(stderr, "retrieval: unknown entity type") {
		t.Fatalf("stderr must not contain unknown entity type: %q", stderr)
	}
	var pkt map[string]any
	if err := json.Unmarshal([]byte(stdout), &pkt); err != nil {
		t.Fatalf("json: %v\nstdout=%s", err, stdout)
	}
	if pkt["deliberation"].(map[string]any)["phase"] != "INVESTIGATE" {
		t.Fatalf("want INVESTIGATE: %#v", pkt["deliberation"])
	}
}

func TestLoopNextInvestigateEmphasizesUncertainties(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "gap?", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	if pkt["deliberation"].(map[string]any)["phase"] != "INVESTIGATE" {
		t.Fatalf("want INVESTIGATE: %#v", pkt["deliberation"])
	}
	items := pkt["open_uncertainties"].(map[string]any)["items"].([]any)
	if len(items) == 0 {
		t.Fatal("open_uncertainties empty")
	}
}

func TestLoopNextExecuteEmphasizesContextAndRelated(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "alpha")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("import { alpha } from './a.js'\nexport function beta() { return alpha() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	related := pkt["related"].(map[string]any)
	if related["available"] != true {
		t.Fatalf("related unavailable: %#v", related)
	}
	if related["snapshot"].(map[string]any)["depth"] != float64(2) {
		t.Fatalf("related depth: %#v", related["snapshot"])
	}

	profile := loop.PhaseContextProfile(deliberation.PhaseExecute)
	if profile.ContextMaxItems != 32 || !profile.ContextIncludeMD || profile.RelatedDepth != 2 {
		t.Fatalf("EXECUTE profile: %#v", profile)
	}
}

func TestLoopNextVerifySurfacesVerificationDebt(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

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
		TaskID:     taskID,
		TestName:   "TestVerifyDebtFixture",
		TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()
	// Mark plan critiqued via apply plan_change.
	pcID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []any{}, "spawned_tasks": []any{},
			"plan_changes": []map[string]any{{"id": pcID, "title": "critique done"}},
		},
	}
	writeJSON(t, filepath.Join(dir, "critique.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "critique.json")}); code != exitOK {
		t.Fatalf("apply plan change: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	if pkt["deliberation"].(map[string]any)["phase"] != "VERIFY" {
		t.Fatalf("phase: %#v", pkt["deliberation"])
	}
	debt := pkt["verification_debt"].(map[string]any)
	if debt["present"] != true {
		t.Fatalf("verification_debt: %#v", debt)
	}
}

func TestLoopApplyUncertaintyWriteAffectsNextSelectNext(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	uncID := "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes": map[string]any{
			"discoveries": []any{}, "plan_changes": []any{},
			"spawned_tasks": []map[string]any{{
				"id": "99999999-9999-4999-8999-999999999998", "title": "follow-on",
			}},
			"uncertainties": []map[string]any{{
				"id": uncID, "title": "Need API shape", "severity": "BLOCKING", "task_id": taskID,
			}},
		},
	}
	writeJSON(t, filepath.Join(dir, "unc.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "unc.json")}); code != exitOK {
		t.Fatalf("apply: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	delib := pkt["deliberation"].(map[string]any)
	if delib["phase"] != "INVESTIGATE" {
		t.Fatalf("phase: %#v", delib["phase"])
	}
	if delib["policy_inputs"].(map[string]any)["blocking_uncertainty_count"].(float64) < 1 {
		t.Fatalf("blocking count not increased: %#v", delib["policy_inputs"])
	}
}

func TestLoopApplyRegressionWriteAffectsPolicyInputs(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
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

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	if pkt["deliberation"].(map[string]any)["policy_inputs"].(map[string]any)["open_regression"] != true {
		t.Fatalf("open_regression not true: %#v", pkt["deliberation"])
	}
}

func TestLoopStatusDeliberationFields(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	var status map[string]any
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatal(err)
	}
	delib, ok := status["deliberation"].(map[string]any)
	if !ok {
		t.Fatalf("deliberation block missing: %#v", status)
	}
	for _, key := range []string{"phase", "recommended_phase", "why_selected", "hop_count", "blocked", "policy_inputs"} {
		if delib[key] == nil {
			t.Fatalf("missing deliberation.%s: %#v", key, delib)
		}
	}
}

func TestLoopStatusBlockedWhenBlockingUncertainty(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title: "blocker", Severity: store.UncertaintySeverityBlocking, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "status", "--task", taskID})
	})
	var status map[string]any
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatal(err)
	}
	if status["deliberation"].(map[string]any)["blocked"] != true {
		t.Fatalf("blocked want true: %#v", status["deliberation"])
	}
}

func TestLoopRecentChangesNoFileBytes(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths:    []domain.ChangePathInput{{Path: "pkg/main.go"}},
		Expected: []domain.ExpectedEffectInput{{Dimension: "latency", Expected: "lower"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.RecordActualEffect(ctx, chg.ID, domain.RecordActualEffectInput{
		Dimension: "latency", Actual: "same", Comparison: store.EffectComparisonPartiallySupported,
	}); err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	raw := []byte(out)
	for _, forbidden := range []string{"content", "patch", "blob"} {
		if strings.Contains(string(raw), `"`+forbidden+`"`) {
			t.Fatalf("packet contains forbidden key %q", forbidden)
		}
	}
	var pkt map[string]any
	if err := json.Unmarshal(raw, &pkt); err != nil {
		t.Fatal(err)
	}
	items := pkt["recent_changes"].(map[string]any)["items"].([]any)
	if len(items) == 0 {
		t.Fatal("recent_changes empty")
	}
	item := items[0].(map[string]any)
	if item["git_commit"] == nil || len(item["paths"].([]any)) == 0 {
		t.Fatalf("recent change shape: %#v", item)
	}
}

func TestWhyTaskIncludesDeliberationTransition(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	env := map[string]any{
		"schema_version": "trace.loop.apply.v1",
		"apply_id":       "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		"seed":           map[string]any{"task_id": taskID, "goal_id": goalID},
		"writes":         map[string]any{"discoveries": []any{}, "plan_changes": []any{}, "spawned_tasks": []any{}},
	}
	writeJSON(t, filepath.Join(dir, "hop.json"), env)
	if code := run([]string{"-C", dir, "loop", "apply", "--in", filepath.Join(dir, "hop.json")}); code != exitOK {
		t.Fatalf("apply: %d", code)
	}

	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "task", taskID})
	})
	var why map[string]any
	if err := json.Unmarshal([]byte(whyOut), &why); err != nil {
		t.Fatalf("why json: %v\n%s", err, whyOut)
	}
	steps, ok := why["steps"].([]any)
	if !ok {
		t.Fatalf("why missing steps: %v", why)
	}
	found := false
	var expectedReason string
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityTask, taskID)
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range evs {
		if ev.Type != domain.EventDeliberationTransition {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal([]byte(ev.PayloadJSON), &payload); err != nil {
			t.Fatal(err)
		}
		expectedReason, _ = payload["reason_code"].(string)
		break
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if expectedReason == "" {
		t.Fatal("missing deliberation.transition event")
	}

	for _, s := range steps {
		step := s.(map[string]any)
		if step["reason_code"] == retrieval.ReasonDeliberationTransition {
			if step["detail"] != expectedReason {
				t.Fatalf("transition detail: want %q got %#v", expectedReason, step)
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing deliberation_transition step: %#v", steps)
	}
}

func TestLoopNextHistoricalRelationshipsSection(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "pkg/main.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordObservedRelationship(ctx, domain.RelInput{
		FromType: domain.EntityChange, FromID: chg.ID,
		ToType: domain.EntityTask, ToID: taskID,
		Confidence: 0.8,
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(out), &pkt); err != nil {
		t.Fatal(err)
	}
	section, ok := pkt["historical_relationships"].(map[string]any)
	if !ok {
		t.Fatalf("missing historical_relationships: %#v", pkt)
	}
	items, ok := section["items"].([]any)
	if !ok || len(items) < 1 || len(items) > 8 {
		t.Fatalf("historical_relationships.items: %#v", section["items"])
	}
	item := items[0].(map[string]any)
	if item["rel"] != domain.RelObservedRelationship {
		t.Fatalf("item rel: %#v", item)
	}
}

func TestHistoricalRelationshipsObservedVsCaused(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := addGoalForLoopTest(t, dir, "Goal")
	taskID := addTaskForLoopTest(t, dir, goalID, "Seed")
	createCurrentDeepPlanForLoopTest(t, dir, goalID)

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "pkg/main.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := mustEvalWithRegressionForLoop(t, svc, taskID)
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskID,
	})
	if err != nil {
		t.Fatal(err)
	}
	ev := mustEvidenceForLoop(t, svc, "git blame")
	if _, err := svc.RecordCausalRelationship(ctx, domain.RelInput{
		FromType: domain.EntityChange, FromID: chg.ID,
		ToType: domain.EntityRegression, ToID: reg.ID,
		Confidence: 0.9,
	}, []string{ev.ID}); err != nil {
		t.Fatal(err)
	}
	chgBare, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "def5678",
		Paths: []domain.ChangePathInput{{Path: "pkg/other.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// caused_by without evidence — must not appear in packet
	if _, err := st.InsertLink(store.EntityLink{
		FromType: domain.EntityChange, FromID: chgBare.ID,
		Rel:    domain.RelCausedBy,
		ToType: domain.EntityTask, ToID: taskID,
		Confidence: 0.5,
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	outJSON := captureStdout(t, func() int {
		return run([]string{"-C", dir, "loop", "next", "--task", taskID})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(outJSON), &pkt); err != nil {
		t.Fatal(err)
	}
	items := pkt["historical_relationships"].(map[string]any)["items"].([]any)
	var bareTaskCaused bool
	var evidencedRegressionCaused bool
	for _, it := range items {
		item := it.(map[string]any)
		if item["rel"] == domain.RelCausedBy && item["to_type"] == domain.EntityTask {
			bareTaskCaused = true
		}
		if item["rel"] == domain.RelCausedBy && item["to_type"] == domain.EntityRegression && item["to_id"] == reg.ID {
			evidencedRegressionCaused = true
		}
	}
	if bareTaskCaused {
		t.Fatal("caused_by without evidence must be excluded")
	}
	if !evidencedRegressionCaused {
		t.Fatal("caused_by with evidence must be included")
	}
}

func mustEvalWithRegressionForLoop(t *testing.T, svc *domain.Service, taskID string) store.OutcomeResult {
	t.Helper()
	ctx := context.Background()
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func mustEvidenceForLoop(t *testing.T, svc *domain.Service, title string) store.Evidence {
	t.Helper()
	ev, err := svc.CreateEvidence(context.Background(), domain.EvidenceInput{Title: title, Body: "proof"})
	if err != nil {
		t.Fatal(err)
	}
	return ev
}
