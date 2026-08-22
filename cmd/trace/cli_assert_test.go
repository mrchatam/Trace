package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestCLIAddSucceedsWhenMCPAddDenied(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "mcp:trace_add", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide mcp:trace_add: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", "CLI add after MCP DENIED"})
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("add json: %v (%q)", err, out)
	}
	id, _ := res["id"].(string)
	if res["ok"] != true || id == "" {
		t.Fatalf("add payload: %v", res)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	if _, err := st.GetGoal(id); err != nil {
		t.Fatalf("entity missing: %v", err)
	}
	row, err := st.GetCapabilityToolDecisionBySlug("cli:add")
	if err != nil {
		t.Fatalf("expected durable cli:add AUTO_ALLOWED: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("cli:add want AUTO_ALLOWED got %q", row.Decision)
	}
}

func TestCLIAddDeniedFailClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "cli:add", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide cli:add: %d", code)
	}

	code, _, stderr := runCapture(t, []string{"-C", dir, "add", "goal", "--title", "must not persist"})
	if code == exitOK {
		t.Fatal("trace add must fail when cli:add is DENIED")
	}
	if !strings.Contains(stderr, "DENIED") {
		t.Fatalf("stderr should mention DENIED, got %q", stderr)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	goals, err := st.ListGoals()
	if err != nil {
		t.Fatal(err)
	}
	if len(goals) != 0 {
		t.Fatalf("DENIED add must not create entity, got %+v", goals)
	}
}

func TestCLIWhySucceedsWhenMCPWhyDenied(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Why seed"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task json: %v", err)
	}
	taskID := taskRes["id"].(string)
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "mcp:trace_why", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide mcp:trace_why: %d", code)
	}

	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "task", taskID})
	})
	if strings.TrimSpace(whyOut) == "" {
		t.Fatal("why stdout empty")
	}
}

func TestCLIWhyDeniedFailClosed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Why gated"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task json: %v", err)
	}
	taskID := taskRes["id"].(string)
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "cli:why", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide cli:why: %d", code)
	}

	code, _, stderr := runCapture(t, []string{"-C", dir, "why", "task", taskID})
	if code == exitOK {
		t.Fatal("trace why must fail when cli:why is DENIED")
	}
	if !strings.Contains(stderr, "DENIED") {
		t.Fatalf("stderr should mention DENIED, got %q", stderr)
	}
}

func TestUngatedCapabilityDecideWhenCLIAddDenied(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "cli:add", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide cli:add DENIED: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "capability", "decide", "--slug", "cli:add", "--decision", "ALLOWED"})
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("decide json: %v (%q)", err, out)
	}
	if res["ok"] != true || res["decision"] != "ALLOWED" {
		t.Fatalf("escape hatch decide: %v", res)
	}
}

func TestUnprefixedAddDecideDoesNotGateCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "add", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide add: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", "unprefixed deny is not CLI"})
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("add json: %v (%q)", err, out)
	}
	if res["ok"] != true {
		t.Fatalf("trace add must succeed when only unprefixed add is DENIED: %v", res)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	row, err := st.GetCapabilityToolDecisionBySlug("add")
	if err != nil {
		t.Fatalf("unprefixed add row: %v", err)
	}
	if row.Decision != store.ToolDecisionDenied {
		t.Fatalf("add slug must stay DENIED, got %q", row.Decision)
	}
	if _, err := st.GetCapabilityToolDecisionBySlug("cli:add"); err != nil {
		t.Fatalf("cli:add AUTO_ALLOWED row: %v", err)
	}
}

func TestCLIIndexAliasDenied(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "cli:index", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("decide cli:index: %d", code)
	}

	code, _, stderr := runCapture(t, []string{"-C", dir, "index"})
	if code == exitOK {
		t.Fatal("trace index must fail when cli:index is DENIED")
	}
	if !strings.Contains(stderr, "DENIED") {
		t.Fatalf("index stderr should mention DENIED, got %q", stderr)
	}

	code2, _, stderr2 := runCapture(t, []string{"-C", dir, "reindex"})
	if code2 == exitOK {
		t.Fatal("trace reindex must fail when cli:index is DENIED")
	}
	if !strings.Contains(stderr2, "DENIED") {
		t.Fatalf("reindex stderr should mention DENIED, got %q", stderr2)
	}
}

func runCapture(t *testing.T, args []string) (code int, stdout, stderr string) {
	t.Helper()
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = wOut, wErr
	code = run(args)
	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	buf := make([]byte, 1<<20)
	n, _ := rOut.Read(buf)
	_ = rOut.Close()
	stdout = string(buf[:n])
	buf2 := make([]byte, 1<<20)
	n2, _ := rErr.Read(buf2)
	_ = rErr.Close()
	stderr = string(buf2[:n2])
	return
}
