package install_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/install"
)

func cursorTarget(t *testing.T) install.Target {
	t.Helper()
	tgt, err := install.Lookup(install.TargetCursor)
	if err != nil {
		t.Fatal(err)
	}
	return tgt
}

func claudeTarget(t *testing.T) install.Target {
	t.Helper()
	tgt, err := install.Lookup(install.TargetClaude)
	if err != nil {
		t.Fatal(err)
	}
	return tgt
}

func cursorHookTarget(t *testing.T) install.Target {
	t.Helper()
	tgt, err := install.Lookup(install.TargetCursorHook)
	if err != nil {
		t.Fatal(err)
	}
	return tgt
}

func setupProjectWithCursor(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".cursor"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func setupProjectWithClaude(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestInstallCursorIncludesLoopGateRule(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")

	tgt := cursorTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Bin:         "/tmp/trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	rulesPath := filepath.Join(root, ".cursor", "rules", "trace-enforcement.mdc")
	raw, err := os.ReadFile(rulesPath)
	if err != nil {
		t.Fatalf("rules file: %v", err)
	}
	body := string(raw)
	if !strings.Contains(body, "loop gate") {
		t.Fatalf("missing loop gate: %s", body)
	}
	if !strings.Contains(body, "TRACE_TASK_ID") {
		t.Fatalf("missing TRACE_TASK_ID: %s", body)
	}
	if !strings.Contains(body, "mandatory gap pass") {
		t.Fatalf("missing mandatory gap pass: %s", body)
	}
	if !strings.Contains(body, "Parent orchestrator") {
		t.Error("cursor rules missing Parent orchestrator rule (INT-04 / P25-2)")
	}
	if !strings.Contains(body, "Worker inheritance") {
		t.Error("cursor rules missing Worker inheritance guidance (FM-04)")
	}
	if !(strings.Contains(body, "trace gap --task") || strings.Contains(body, "trace loop status")) {
		t.Fatalf("missing trace gap/loop status command guidance: %s", body)
	}
}

func TestInstallAgentsMDEnforcementBlock(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")

	tgt := cursorTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "# begin-trace-enforcement") || !strings.Contains(body, "# end-trace-enforcement") {
		t.Fatalf("missing enforcement markers: %s", body)
	}
	if !strings.Contains(body, "loop gate") {
		t.Fatalf("missing gate command: %s", body)
	}
	if !strings.Contains(body, "mandatory gap pass") {
		t.Fatalf("missing mandatory gap pass: %s", body)
	}
	if !strings.Contains(body, "TRACE_TASK_ID") {
		t.Fatalf("missing TRACE_TASK_ID: %s", body)
	}
	if !strings.Contains(body, "Parent orchestrator") {
		t.Error("AGENTS.md missing Parent orchestrator rule (INT-04 / FM-04)")
	}
	if !(strings.Contains(body, "trace gap --task") || strings.Contains(body, "trace loop status")) {
		t.Fatalf("missing trace gap/loop status command guidance: %s", body)
	}
}

func TestInstallAgentsMDMarkersIdempotent(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}
	tgt := cursorTarget(t)
	if err := tgt.Install(opts); err != nil {
		t.Fatalf("first Install: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tgt.Install(opts); err != nil {
		t.Fatalf("second Install: %v", err)
	}
	second, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(second), "# begin-trace-enforcement") != 1 {
		t.Fatalf("want single begin marker, got:\n%s", second)
	}
	if strings.Count(string(second), "# end-trace-enforcement") != 1 {
		t.Fatalf("want single end marker, got:\n%s", second)
	}
	if !strings.Contains(string(second), "loop gate") {
		t.Fatalf("block content missing after re-run")
	}
	_ = first
}

func TestInstallCursorRulesPrintOnlyStderr(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")

	var stdout, stderr strings.Builder
	tgt := cursorTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       false,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Out:         &stdout,
		ErrOut:      &stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal([]byte(stdout.String()), &doc); err != nil {
		t.Fatalf("stdout JSON: %v\n%s", err, stdout.String())
	}
	if _, ok := doc["mcpServers"].(map[string]any)["trace"]; !ok {
		t.Fatal("missing mcpServers.trace on stdout")
	}
	if !strings.Contains(stderr.String(), "loop gate") {
		t.Fatalf("stderr should mention loop gate, got %q", stderr.String())
	}
}

func TestInstallCursorUninstallRemovesRules(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}
	tgt := cursorTarget(t)
	if err := tgt.Install(opts); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := tgt.Uninstall(install.InstallOpts{ProjectRoot: root, MCPJSON: mcpPath}); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	rulesPath := filepath.Join(root, ".cursor", "rules", "trace-enforcement.mdc")
	if _, err := os.Stat(rulesPath); !os.IsNotExist(err) {
		t.Fatalf("rules file should be removed: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err == nil && strings.Contains(string(raw), "# begin-trace-enforcement") {
		t.Fatalf("AGENTS block should be stripped: %s", raw)
	}
	var doc map[string]any
	mcpRaw, _ := os.ReadFile(mcpPath)
	if len(mcpRaw) > 0 {
		if err := json.Unmarshal(mcpRaw, &doc); err == nil {
			if servers, ok := doc["mcpServers"].(map[string]any); ok {
				if _, ok := servers["trace"]; ok {
					t.Fatal("MCP trace should be removed")
				}
			}
		}
	}
}

func TestInstallClaudeIncludesLoopGateRule(t *testing.T) {
	root := setupProjectWithClaude(t)
	tgt := claudeTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	foundLoopGate := false
	foundGapPass := false
	for _, path := range []string{
		filepath.Join(root, "CLAUDE.md"),
		filepath.Join(root, ".claude", "trace-enforcement-rules.md"),
	} {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		body := string(raw)
		if strings.Contains(body, "loop gate") {
			foundLoopGate = true
		}
		if strings.Contains(body, "mandatory gap pass") &&
			strings.Contains(body, "TRACE_TASK_ID") &&
			(strings.Contains(body, "trace gap --task") || strings.Contains(body, "trace loop status")) {
			foundGapPass = true
		}
	}
	if !foundLoopGate {
		t.Fatal("no claude rules surface references loop gate")
	}
	if !foundGapPass {
		t.Fatal("no claude rules surface references mandatory gap pass prompt")
	}
}

func TestInstallClaudeRulesWithCLAUDEmd(t *testing.T) {
	root := setupProjectWithClaude(t)
	claudeMD := filepath.Join(root, "CLAUDE.md")
	if err := os.WriteFile(claudeMD, []byte("# Project\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tgt := claudeTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	raw, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "# begin-trace-enforcement") {
		t.Fatalf("marker block should be inside CLAUDE.md: %s", body)
	}
}

func TestInstallClaudeRulesFallbackFile(t *testing.T) {
	root := setupProjectWithClaude(t)
	tgt := claudeTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	path := filepath.Join(root, ".claude", "trace-enforcement-rules.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("fallback rules file missing: %v", err)
	}
}

func TestInstallCursorHookCallsGate(t *testing.T) {
	root := setupProjectWithCursor(t)
	tgt := cursorHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".cursor", "hooks", "trace-loop-gate.sh"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "loop gate") {
		t.Fatalf("missing loop gate in script: %s", body)
	}
	if !strings.Contains(body, "--for edit") {
		t.Fatalf("missing --for edit in script: %s", body)
	}
	// Option A: empty-task branch must gate on enforce=strict (not always-allow).
	if !strings.Contains(body, `"enforce"`) || !strings.Contains(body, "strict") {
		t.Fatalf("script missing Option A strict config check: %s", body)
	}
}

func TestCursorLoopGateFailClosedStrictNoTask(t *testing.T) {
	root := t.TempDir()
	scriptPath := writeLoopGateScript(t, root)
	if err := os.MkdirAll(filepath.Join(root, ".trace"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".trace", "config.json"), []byte(`{"enforce":"strict"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := runLoopGateHook(t, root, scriptPath, "", nil)
	assertPermissionJSON(t, out, "deny")
	if code == 0 {
		t.Fatalf("strict+no-task: want non-zero exit, got 0 stdout=%s", out)
	}
}

func TestCursorLoopGateAllowNonStrictNoTask(t *testing.T) {
	cases := []struct {
		name   string
		config string // empty = missing file
	}{
		{name: "missing", config: ""},
		{name: "off", config: `{"enforce":"off"}`},
		{name: "warn", config: `{"enforce":"warn"}`},
		{name: "invalid", config: `{"enforce":"nope"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			scriptPath := writeLoopGateScript(t, root)
			if tc.config != "" {
				if err := os.MkdirAll(filepath.Join(root, ".trace"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, ".trace", "config.json"), []byte(tc.config), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			out, code := runLoopGateHook(t, root, scriptPath, "", nil)
			assertPermissionJSON(t, out, "allow")
			if code != 0 {
				t.Fatalf("%s+no-task: want exit 0, got %d stdout=%s", tc.name, code, out)
			}
		})
	}
}

func TestCursorLoopGateScriptPreservesGatePathWhenTaskSet(t *testing.T) {
	body := install.CursorLoopGateHookScript()
	if !strings.Contains(body, "loop gate") || !strings.Contains(body, "--for edit") {
		t.Fatalf("task-set path must still call loop gate --for edit:\n%s", body)
	}
	// Empty-task allow must remain reachable for non-strict (default-off).
	if !strings.Contains(body, `{"permission":"allow"}`) {
		t.Fatal("script must still emit allow JSON for non-strict empty task")
	}
}

func writeLoopGateScript(t *testing.T, root string) string {
	t.Helper()
	scriptPath := filepath.Join(root, ".cursor", "hooks", "trace-loop-gate.sh")
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scriptPath, []byte(install.CursorLoopGateHookScript()), 0o755); err != nil {
		t.Fatal(err)
	}
	return scriptPath
}

func TestInstallCursorHookPreToolUseMatcher(t *testing.T) {
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
		t.Fatal(err)
	}
	hooks, _ := doc["hooks"].(map[string]any)
	pre, _ := hooks["preToolUse"].([]any)
	found := false
	for _, item := range pre {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		cmd, _ := m["command"].(string)
		if strings.Contains(cmd, "trace-loop-gate.sh") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("hooks.json missing trace preToolUse entry: %s", raw)
	}
}

func TestInstallDetectIncludesCursorHook(t *testing.T) {
	root := setupProjectWithCursor(t)
	infos := install.ListTargets(install.InstallOpts{ProjectRoot: root})
	var hook *install.DetectInfo
	for i := range infos {
		if infos[i].ID == install.TargetCursorHook {
			hook = &infos[i]
			break
		}
	}
	if hook == nil {
		t.Fatalf("detect missing cursor-hook: %+v", infos)
	}
	if hook.Tier != install.TierConditional {
		t.Fatalf("tier: want CONDITIONAL got %q", hook.Tier)
	}
}

func TestInstallDetectCursorHookConditional(t *testing.T) {
	root := t.TempDir()
	infos := install.ListTargets(install.InstallOpts{ProjectRoot: root})
	var hook *install.DetectInfo
	for i := range infos {
		if infos[i].ID == install.TargetCursorHook {
			hook = &infos[i]
			break
		}
	}
	if hook == nil {
		t.Fatal("detect missing cursor-hook")
	}
	if hook.Detected {
		t.Fatalf("want not detected without .cursor/, got reason %q", hook.Reason)
	}
}

func TestInstallCursorHookUninstallRemovesScript(t *testing.T) {
	root := setupProjectWithCursor(t)
	tgt := cursorHookTarget(t)
	opts := install.InstallOpts{Write: true, ProjectRoot: root, ErrOut: os.Stderr}
	if err := tgt.Install(opts); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := tgt.Uninstall(install.InstallOpts{ProjectRoot: root}); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	scriptPath := filepath.Join(root, ".cursor", "hooks", "trace-loop-gate.sh")
	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		t.Fatalf("script should be removed: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".cursor", "hooks.json"))
	if err == nil {
		if strings.Contains(string(raw), "trace-loop-gate.sh") {
			t.Fatalf("hooks.json entry should be removed: %s", raw)
		}
	}
}

func TestInstallEnforcementIdempotent(t *testing.T) {
	root := setupProjectWithCursor(t)
	mcpPath := filepath.Join(t.TempDir(), "mcp.json")
	cursorOpts := install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Bin:         "trace-mcp",
		ErrOut:      os.Stderr,
	}
	hookOpts := install.InstallOpts{Write: true, ProjectRoot: root, ErrOut: os.Stderr}

	cursorTgt := cursorTarget(t)
	hookTgt := cursorHookTarget(t)
	for i := 0; i < 2; i++ {
		if err := cursorTgt.Install(cursorOpts); err != nil {
			t.Fatalf("cursor Install pass %d: %v", i, err)
		}
		if err := hookTgt.Install(hookOpts); err != nil {
			t.Fatalf("cursor-hook Install pass %d: %v", i, err)
		}
	}
	raw, err := os.ReadFile(filepath.Join(root, ".cursor", "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(raw), "trace-loop-gate.sh") != 1 {
		t.Fatalf("want single hook entry after idempotent installs: %s", raw)
	}
}

// Tests for P25-S01-01 constants (INT-03, INT-04, INT-11).

func TestGapPassPromptNonEmpty(t *testing.T) {
	if install.GapPassPrompt == "" {
		t.Fatal("GapPassPrompt must be non-empty")
	}
	if !strings.Contains(install.GapPassPrompt, "gap") {
		t.Fatalf("GapPassPrompt must contain the word 'gap', got: %s", install.GapPassPrompt)
	}
	if !strings.Contains(install.GapPassPrompt, "TRACE_TASK_ID") {
		t.Fatalf("GapPassPrompt must reference TRACE_TASK_ID, got: %s", install.GapPassPrompt)
	}
	if !strings.Contains(install.GapPassPrompt, "promotion_candidates") {
		t.Fatalf("GapPassPrompt must mention promotion_candidates (post-import nudge), got: %s", install.GapPassPrompt)
	}
	// FM-02: write discoveries/decisions before seed export --strict --enforce.
	if !strings.Contains(install.GapPassPrompt, "Write-before-export") {
		t.Fatalf("GapPassPrompt must require Write-before-export, got: %s", install.GapPassPrompt)
	}
	if !strings.Contains(install.GapPassPrompt, "seed export") || !strings.Contains(install.GapPassPrompt, "--strict --enforce") {
		t.Fatalf("GapPassPrompt must mention seed export --strict --enforce, got: %s", install.GapPassPrompt)
	}
	if !strings.Contains(install.GapPassPrompt, "discovery") || !strings.Contains(install.GapPassPrompt, "decision") {
		t.Fatalf("GapPassPrompt must mention discovery and decision writes, got: %s", install.GapPassPrompt)
	}
	// FM-08 / INT-06: post-discovery → task/promotion before product edits.
	for _, want := range []string{
		"Post-discovery nudge",
		"before product edits",
		"--from-discovery",
		"spawned_tasks[].discovery_id",
		"discovery-only",
	} {
		if !strings.Contains(install.GapPassPrompt, want) {
			t.Fatalf("GapPassPrompt must mention %q (FM-08 post-discovery nudge), got: %s", want, install.GapPassPrompt)
		}
	}
}

func TestParentOrchestratorRuleNonEmpty(t *testing.T) {
	if install.ParentOrchestratorRule == "" {
		t.Fatal("ParentOrchestratorRule must be non-empty")
	}
	if !strings.Contains(install.ParentOrchestratorRule, "TRACE_TASK_ID") {
		t.Fatalf("ParentOrchestratorRule must reference TRACE_TASK_ID")
	}
	if !strings.Contains(install.ParentOrchestratorRule, "preToolUse") {
		t.Fatalf("ParentOrchestratorRule must reference preToolUse deny path")
	}
	if !strings.Contains(install.ParentOrchestratorRule, "failClosed") {
		t.Fatalf("ParentOrchestratorRule must contain failClosed semantics")
	}
	// FM-04: parent owns graph; workers get task via prompt+env; Multitask limits documented.
	if !strings.Contains(install.ParentOrchestratorRule, "offloading graph") &&
		!strings.Contains(install.ParentOrchestratorRule, "graph-only") {
		t.Fatalf("ParentOrchestratorRule must forbid offloading graph-only work to workers: %s", install.ParentOrchestratorRule)
	}
	if !strings.Contains(install.ParentOrchestratorRule, "Worker inheritance") {
		t.Fatalf("ParentOrchestratorRule must document Worker inheritance")
	}
	if !strings.Contains(install.ParentOrchestratorRule, "Multitask") {
		t.Fatalf("ParentOrchestratorRule must document Multitask limits")
	}
	if !strings.Contains(install.ParentOrchestratorRule, "Option B") {
		t.Fatalf("ParentOrchestratorRule must note Option B deferred")
	}
	block := install.AgentsEnforcementBlock()
	if !strings.Contains(block, "Parent orchestrator") {
		t.Fatal("AgentsEnforcementBlock must include ParentOrchestratorRule (FM-04)")
	}
}

func TestHookDriftNoteNonEmpty(t *testing.T) {
	if install.HookDriftNote == "" {
		t.Fatal("HookDriftNote must be non-empty")
	}
	if !strings.Contains(install.HookDriftNote, "preToolUse") {
		t.Fatalf("HookDriftNote must reference preToolUse schema field")
	}
	if !strings.Contains(install.HookDriftNote, "trace install") {
		t.Fatalf("HookDriftNote must reference trace install upgrade step")
	}
}

func TestInstallGitHookUnchanged(t *testing.T) {
	root := initGitRepo(t)
	tgt := gitHookTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	hooksDir, err := resolveHooksDirForTest(root)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(hooksDir, "post-commit"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if !strings.Contains(body, "# begin-trace") {
		t.Fatal("git-hook must use # begin-trace marker")
	}
	if strings.Contains(body, "# begin-trace-enforcement") {
		t.Fatal("git-hook must not use enforcement markers")
	}
	if !strings.Contains(body, "index") || !strings.Contains(body, "seed export") {
		t.Fatalf("git-hook must retain index/export commands: %s", body)
	}
}

func TestInstallCursorMCPUnchanged(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, "mcp.json")
	prior := []byte(`{
  "mcpServers": {
    "other": {"type": "stdio", "command": "other-mcp", "args": []}
  }
}
`)
	if err := os.WriteFile(mcpPath, prior, 0o644); err != nil {
		t.Fatal(err)
	}
	root := setupProjectWithCursor(t)
	tgt := cursorTarget(t)
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		MCPJSON:     mcpPath,
		Bin:         "/opt/trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	servers := doc["mcpServers"].(map[string]any)
	if _, ok := servers["other"]; !ok {
		t.Fatal("sibling other must be preserved")
	}
	trace := servers["trace"].(map[string]any)
	if cmd, _ := trace["command"].(string); cmd != "/opt/trace-mcp" {
		t.Fatalf("command: %#v", trace["command"])
	}
	args := trace["args"].([]any)
	if args[0] != "-C" || args[1] != "${workspaceFolder}" {
		t.Fatalf("args: %#v", args)
	}
}
