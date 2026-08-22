package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInstallCursorPrintSnippet(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"install", "cursor"})
	})
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, out)
	}
	servers, ok := doc["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("missing mcpServers: %s", out)
	}
	trace, ok := servers["trace"].(map[string]any)
	if !ok {
		t.Fatalf("missing mcpServers.trace: %s", out)
	}
	if typ, _ := trace["type"].(string); typ != "stdio" {
		t.Fatalf("type: want stdio got %#v", trace["type"])
	}
	if cmd, _ := trace["command"].(string); cmd != "trace-mcp" {
		t.Fatalf("command: want trace-mcp got %#v", trace["command"])
	}
	args, ok := trace["args"].([]any)
	if !ok || len(args) != 2 {
		t.Fatalf("args: want [-C, ${workspaceFolder}] got %#v", trace["args"])
	}
	if args[0] != "-C" || args[1] != "${workspaceFolder}" {
		t.Fatalf("args: want [-C, ${workspaceFolder}] got %#v", args)
	}
	if !strings.Contains(out, "\n") {
		t.Fatalf("expected pretty JSON (multiline), got %q", out)
	}
}

// TestInstallCursorPrintReloadTip — DF-50/DF-22: print-only success emits the same
// stderr reload tip as --write; stdout stays pretty mcpServers.trace JSON only.
func TestInstallCursorPrintReloadTip(t *testing.T) {
	out, errOut := captureStdoutStderr(t, func() int {
		return run([]string{"install", "cursor"})
	})
	if !strings.Contains(errOut, "reload") || !strings.Contains(errOut, "trace-mcp") {
		t.Fatalf("print stderr should include DF-22/50 reload tip, got %q", errOut)
	}
	if !strings.Contains(errOut, installCursorReloadTip) {
		t.Fatalf("print tip must match shared helper text %q, got %q", installCursorReloadTip, errOut)
	}
	if strings.Contains(out, "reload") || strings.Contains(out, "install: tip") {
		t.Fatalf("tip must not appear on stdout, got %q", out)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, out)
	}
	servers, ok := doc["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("missing mcpServers: %s", out)
	}
	if _, ok := servers["trace"].(map[string]any); !ok {
		t.Fatalf("missing mcpServers.trace: %s", out)
	}
}

func TestInstallCursorPrintBin(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"install", "cursor", "--bin", "/opt/trace-mcp"})
	})
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout JSON: %v", err)
	}
	trace := doc["mcpServers"].(map[string]any)["trace"].(map[string]any)
	if cmd, _ := trace["command"].(string); cmd != "/opt/trace-mcp" {
		t.Fatalf("command: want /opt/trace-mcp got %#v", trace["command"])
	}
}

func TestInstallCursorWriteMergeBackup(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, "mcp.json")
	prior := []byte(`{
  "mcpServers": {
    "other": {
      "type": "stdio",
      "command": "other-mcp",
      "args": []
    },
    "trace": {
      "type": "stdio",
      "command": "old-trace-mcp",
      "args": ["-C", "/old"]
    }
  }
}
`)
	if err := os.WriteFile(mcpPath, prior, 0o644); err != nil {
		t.Fatal(err)
	}

	stderr := captureStderr(t, func() int {
		return run([]string{"install", "cursor", "--write", "--mcp-json", mcpPath, "--bin", "/tmp/trace-mcp"})
	})
	if !strings.Contains(stderr, mcpPath+".bak.") {
		t.Fatalf("stderr should mention backup path, got %q", stderr)
	}
	if !strings.Contains(stderr, "reload") || !strings.Contains(stderr, "trace-mcp") {
		t.Fatalf("stderr should include DF-22 reload tip, got %q", stderr)
	}
	if !strings.Contains(stderr, installCursorReloadTip) {
		t.Fatalf("write tip must match shared helper text %q, got %q", installCursorReloadTip, stderr)
	}

	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("written JSON: %v", err)
	}
	servers := doc["mcpServers"].(map[string]any)
	if _, ok := servers["other"]; !ok {
		t.Fatal("sibling server other must be preserved")
	}
	trace := servers["trace"].(map[string]any)
	if cmd, _ := trace["command"].(string); cmd != "/tmp/trace-mcp" {
		t.Fatalf("upserted command: %#v", trace["command"])
	}
	args := trace["args"].([]any)
	if args[0] != "-C" || args[1] != "${workspaceFolder}" {
		t.Fatalf("upserted args: %#v", args)
	}

	matches, err := filepath.Glob(mcpPath + ".bak.*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("want one backup, got %v", matches)
	}
	bak, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(bak) != string(prior) {
		t.Fatalf("backup content mismatch")
	}
	// UTC stamp shape: .bak.YYYYMMDDTHHMMSSZ
	suffix := strings.TrimPrefix(filepath.Base(matches[0]), filepath.Base(mcpPath)+".bak.")
	if _, err := time.Parse("20060102T150405Z", suffix); err != nil {
		t.Fatalf("backup UTC suffix %q: %v", suffix, err)
	}
}

func TestInstallCursorWriteCreateMissing(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, "nested", "mcp.json")
	stderr := captureStderr(t, func() int {
		return run([]string{"install", "cursor", "--write", "--mcp-json", mcpPath})
	})
	if !strings.Contains(stderr, "reload") {
		t.Fatalf("new-file --write should still print DF-22 tip, got %q", stderr)
	}
	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	if _, ok := doc["mcpServers"].(map[string]any)["trace"]; !ok {
		t.Fatalf("missing trace: %s", raw)
	}
	matches, _ := filepath.Glob(mcpPath + ".bak.*")
	if len(matches) != 0 {
		t.Fatalf("no backup for new file, got %v", matches)
	}
}

func TestInstallCursorWriteInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, "mcp.json")
	if err := os.WriteFile(mcpPath, []byte("{not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	code := run([]string{"install", "cursor", "--write", "--mcp-json", mcpPath})
	if code != exitFail {
		t.Fatalf("want exitFail got %d", code)
	}
	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "{not-json" {
		t.Fatalf("must not write on invalid JSON: %q", raw)
	}
	matches, _ := filepath.Glob(mcpPath + ".bak.*")
	if len(matches) != 0 {
		t.Fatalf("must not backup on fail-closed: %v", matches)
	}
}

func TestHelpIncludesCursorHook(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	for _, want := range []string{"cursor-hook", "TRACE_TASK_ID", "loop gate"} {
		if !strings.Contains(out, want) {
			t.Fatalf("help should mention %q, got excerpt:\n%s", want, out)
		}
	}
}

func TestInstallUsage(t *testing.T) {
	if code := run([]string{"install"}); code != exitUsage {
		t.Fatalf("install alone: want usage got %d", code)
	}
	if code := run([]string{"install", "vscode"}); code != exitUsage {
		t.Fatalf("unknown target: want usage got %d", code)
	}
}

// DF-68: install -C vs process cwd. Helpers follow cli_test.go Chdir+Cleanup.

func chdirTemp(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func absTemp(t *testing.T, dir string) string {
	t.Helper()
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func TestInstallClaudeDashCRefuseCitesProjectRoot(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	absProject := absTemp(t, projectDir)
	absCwd := absTemp(t, cwdDir)
	chdirTemp(t, cwdDir)

	code, _, stderr := runCapture(t, []string{"-C", projectDir, "install", "claude"})
	if code != exitFail {
		t.Fatalf("want refuse without marker, got exit %d stderr %q", code, stderr)
	}
	if !strings.Contains(stderr, "under "+absProject) {
		t.Fatalf("refuse should cite Abs(projectDir) after under, want %q got %q", absProject, stderr)
	}
	if strings.Contains(stderr, absCwd) {
		t.Fatalf("refuse must not cite cwdDir %q as the marker root, got %q", absCwd, stderr)
	}
}

func TestInstallClaudeDashCIgnoresCwdMarker(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(cwdDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	absProject := absTemp(t, projectDir)
	absCwd := absTemp(t, cwdDir)
	chdirTemp(t, cwdDir)

	code, _, stderr := runCapture(t, []string{"-C", projectDir, "install", "claude"})
	if code != exitFail {
		t.Fatalf("cwd marker must not authorize -C project, got exit %d stderr %q stdout skipped", code, stderr)
	}
	if !strings.Contains(stderr, absProject) {
		t.Fatalf("refuse should cite projectDir %q, got %q", absProject, stderr)
	}
	if strings.Contains(stderr, absCwd) {
		t.Fatalf("refuse must not cite cwdDir %q, got %q", absCwd, stderr)
	}
}

func TestInstallClaudeDashCWriteUsesProjectRoot(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	chdirTemp(t, cwdDir)

	code, _, stderr := runCapture(t, []string{"-C", projectDir, "install", "claude", "--write"})
	if code != exitOK {
		t.Fatalf("write under -C projectDir: want exitOK got %d stderr %q", code, stderr)
	}
	want := filepath.Join(projectDir, ".claude", "trace-mcp.json")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected write at %s: %v", want, err)
	}
	decoy := filepath.Join(cwdDir, ".claude", "trace-mcp.json")
	if _, err := os.Stat(decoy); !os.IsNotExist(err) {
		t.Fatalf("cwdDir must stay untouched, stat %s: %v", decoy, err)
	}
}

func TestInstallDetectDashCClaudeReasonCitesRoot(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	absProject := absTemp(t, projectDir)
	absCwd := absTemp(t, cwdDir)
	chdirTemp(t, cwdDir)

	code, stdout, stderr := runCapture(t, []string{"-C", projectDir, "install", "detect"})
	if code != exitOK {
		t.Fatalf("detect: want exitOK got %d stderr %q", code, stderr)
	}
	claude := detectInfoByID(t, stdout, "claude")
	if !claude.Detected {
		t.Fatalf("claude should be detected under projectDir, got %+v", claude)
	}
	if !strings.Contains(claude.Reason, absProject) {
		t.Fatalf("claude reason should contain Abs(projectDir) %q, got %q", absProject, claude.Reason)
	}
	if strings.Contains(claude.Reason, "under .") {
		t.Fatalf("claude reason must not cite under . as the root, got %q", claude.Reason)
	}
	if strings.Contains(claude.Reason, absCwd) {
		t.Fatalf("claude reason must not cite cwdDir %q, got %q", absCwd, claude.Reason)
	}
}

func TestInstallDetectDashCCursorHomeUnchanged(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	absProject := absTemp(t, projectDir)
	chdirTemp(t, cwdDir)

	code, stdout, stderr := runCapture(t, []string{"-C", projectDir, "install", "detect"})
	if code != exitOK {
		t.Fatalf("detect: want exitOK got %d stderr %q", code, stderr)
	}
	cursor := detectInfoByID(t, stdout, "cursor")
	if strings.Contains(cursor.Reason, absProject) {
		t.Fatalf("cursor reason must not contain projectDir %q (STABLE home), got %q", absProject, cursor.Reason)
	}
}

func TestInstallUninstallClaudeDashCUsesProjectRoot(t *testing.T) {
	cwdDir := t.TempDir()
	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	chdirTemp(t, cwdDir)

	code, _, stderr := runCapture(t, []string{"-C", projectDir, "install", "claude", "--write"})
	if code != exitOK {
		t.Fatalf("write under -C: want exitOK got %d stderr %q", code, stderr)
	}
	projectFile := filepath.Join(projectDir, ".claude", "trace-mcp.json")
	if _, err := os.Stat(projectFile); err != nil {
		t.Fatalf("expected project write: %v", err)
	}

	decoyDir := filepath.Join(cwdDir, ".claude")
	if err := os.MkdirAll(decoyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	decoy := filepath.Join(decoyDir, "trace-mcp.json")
	if err := os.WriteFile(decoy, []byte(`{"decoy":true}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, stderr = runCapture(t, []string{"-C", projectDir, "install", "uninstall", "claude"})
	if code != exitOK {
		t.Fatalf("uninstall under -C: want exitOK got %d stderr %q", code, stderr)
	}
	if _, err := os.Stat(projectFile); !os.IsNotExist(err) {
		t.Fatalf("projectDir trace-mcp.json should be removed, stat: %v", err)
	}
	raw, err := os.ReadFile(decoy)
	if err != nil {
		t.Fatalf("cwd decoy must remain: %v", err)
	}
	if string(raw) != `{"decoy":true}`+"\n" {
		t.Fatalf("cwd decoy must be untouched, got %q", raw)
	}
}

func detectInfoByID(t *testing.T, stdout, id string) struct {
	ID       string `json:"id"`
	Tier     string `json:"tier"`
	Detected bool   `json:"detected"`
	Reason   string `json:"reason"`
} {
	t.Helper()
	var infos []struct {
		ID       string `json:"id"`
		Tier     string `json:"tier"`
		Detected bool   `json:"detected"`
		Reason   string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(stdout), &infos); err != nil {
		t.Fatalf("detect JSON: %v\n%s", err, stdout)
	}
	for _, info := range infos {
		if info.ID == id {
			return info
		}
	}
	t.Fatalf("detect missing id %q in %s", id, stdout)
	return struct {
		ID       string `json:"id"`
		Tier     string `json:"tier"`
		Detected bool   `json:"detected"`
		Reason   string `json:"reason"`
	}{}
}

// captureStderr runs fn while redirecting os.Stderr to a pipe; requires exitOK.
func captureStderr(t *testing.T, fn func() int) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	code := fn()
	_ = w.Close()
	os.Stderr = old
	if code != exitOK {
		t.Fatalf("command exit %d", code)
	}
	buf := make([]byte, 1<<20)
	n, _ := r.Read(buf)
	_ = r.Close()
	return string(buf[:n])
}

// captureStdoutStderr redirects both stdout and stderr; requires exitOK.
func captureStdoutStderr(t *testing.T, fn func() int) (stdout, stderr string) {
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
	code := fn()
	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	if code != exitOK {
		t.Fatalf("command exit %d", code)
	}
	buf := make([]byte, 1<<20)
	nOut, _ := rOut.Read(buf)
	_ = rOut.Close()
	stdout = string(buf[:nOut])
	bufErr := make([]byte, 1<<20)
	nErr, _ := rErr.Read(bufErr)
	_ = rErr.Close()
	stderr = string(bufErr[:nErr])
	return stdout, stderr
}
