package install_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/install"
)

// INT-11: hooks.json shape + allow/deny permission JSON from the loop-gate script.

func TestHookDriftHooksJSONShape(t *testing.T) {
	root := setupProjectWithCursor(t)
	tgt := cursorHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(root, ".cursor", "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("hooks.json: %v", err)
	}
	hooks, _ := doc["hooks"].(map[string]any)
	pre, _ := hooks["preToolUse"].([]any)
	var entry map[string]any
	for _, item := range pre {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		cmd, _ := m["command"].(string)
		if strings.Contains(cmd, "trace-loop-gate.sh") {
			entry = m
			break
		}
	}
	if entry == nil {
		t.Fatalf("missing trace-loop-gate preToolUse entry: %s", raw)
	}
	for _, key := range []string{"command", "matcher", "failClosed"} {
		if _, ok := entry[key]; !ok {
			t.Fatalf("hooks.json entry missing %q: %#v", key, entry)
		}
	}
	if fc, ok := entry["failClosed"].(bool); !ok || fc {
		t.Fatalf("failClosed must be false (policy is script Option A), got %#v", entry["failClosed"])
	}
	if !strings.Contains(entry["command"].(string), "trace-loop-gate.sh") {
		t.Fatalf("command: %#v", entry["command"])
	}
	matcher, _ := entry["matcher"].(string)
	if matcher == "" {
		t.Fatal("matcher must be non-empty")
	}
}

func TestHookDriftAllowDenyPermissionJSON(t *testing.T) {
	root := setupProjectWithCursor(t)
	scriptPath := filepath.Join(root, ".cursor", "hooks", "trace-loop-gate.sh")
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scriptPath, []byte(install.CursorLoopGateHookScript()), 0o755); err != nil {
		t.Fatal(err)
	}

	// allow: no config, empty TRACE_TASK_ID
	allowOut, allowCode := runLoopGateHook(t, root, scriptPath, "", nil)
	assertPermissionJSON(t, allowOut, "allow")
	if allowCode != 0 {
		t.Fatalf("allow exit: want 0 got %d stdout=%s", allowCode, allowOut)
	}

	// deny: enforce=strict, empty TRACE_TASK_ID
	if err := os.MkdirAll(filepath.Join(root, ".trace"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".trace", "config.json"), []byte(`{"enforce":"strict"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	denyOut, denyCode := runLoopGateHook(t, root, scriptPath, "", nil)
	assertPermissionJSON(t, denyOut, "deny")
	if denyCode == 0 {
		t.Fatalf("deny exit: want non-zero got 0 stdout=%s", denyOut)
	}
}

func assertPermissionJSON(t *testing.T, stdout, wantPerm string) {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &payload); err != nil {
		t.Fatalf("permission JSON: %v\nstdout=%q", err, stdout)
	}
	if _, ok := payload["permission"]; !ok {
		t.Fatalf("missing permission key: %#v", payload)
	}
	got, _ := payload["permission"].(string)
	if got != wantPerm {
		t.Fatalf("permission: want %q got %q payload=%#v", wantPerm, got, payload)
	}
}

func runLoopGateHook(t *testing.T, root, scriptPath, taskID string, extraEnv []string) (stdout string, exitCode int) {
	t.Helper()
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = root
	cmd.Stdin = bytes.NewReader(nil)
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"TRACE_PROJECT_ROOT=" + root,
	}
	if taskID != "" {
		env = append(env, "TRACE_TASK_ID="+taskID)
	}
	env = append(env, extraEnv...)
	cmd.Env = env
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err == nil {
		return out.String(), 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return out.String(), ee.ExitCode()
	}
	t.Fatalf("run hook: %v stderr=%s", err, errBuf.String())
	return "", -1
}
