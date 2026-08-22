package install

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	cursorHookCommandRel = ".cursor/hooks/trace-loop-gate.sh"
	cursorHookMatcher    = "Write|StrReplace|ApplyPatch|EditNotebook"
)

// HookDriftNote documents the Cursor upgrade schema compatibility verification
// steps agents and operators must follow after any Cursor upgrade (INT-11).
//
// Schema field to watch: the preToolUse entry format in .cursor/hooks.json.
// Policy failClosed is implemented in CursorLoopGateHookScript (Option A:
// enforce=strict + empty TRACE_TASK_ID → deny). Cursor hooks.json keeps
// failClosed: false so hook crashes do not lock out default-off projects.
const HookDriftNote = `Hook drift verification (INT-11) — run after each Cursor upgrade:

1. Run ` + "`trace install --check`" + ` (or ` + "`trace install cursor-hook`" + ` without --write) to
   verify .cursor/hooks.json schema compatibility.
2. Schema field to watch: preToolUse entry format in .cursor/hooks.json
   (each entry: {"command": "...", "matcher": "...", "failClosed": false}).
   Cursor failClosed stays false; script-level Option A deny applies when
   .trace/config.json enforce=strict and TRACE_TASK_ID is absent.
3. Escalation: if schema drift is detected, re-run ` + "`trace install`" + ` (with --write)
   to regenerate hooks.json from the current template.`

type cursorHookTarget struct{}

func (cursorHookTarget) ID() string   { return TargetCursorHook }
func (cursorHookTarget) Tier() string { return TierConditional }

func (cursorHookTarget) Detect(opts InstallOpts) DetectResult {
	root := projectRoot(opts)
	if hasCursorDir(root) {
		return DetectResult{Detected: true, Reason: "found .cursor/ under " + root}
	}
	if hasTraceHookInHooksJSON(root) {
		return DetectResult{Detected: true, Reason: "found trace hook in .cursor/hooks.json under " + root}
	}
	return DetectResult{Detected: false, Reason: "no .cursor/ under " + root}
}

func (c cursorHookTarget) Install(opts InstallOpts) error {
	root := projectRoot(opts)
	errOut := opts.ErrOut
	if errOut == nil {
		errOut = io.Discard
	}

	if !opts.Write {
		fmt.Fprint(errOut, EnforcementRulesMarkdown())
		fmt.Fprintf(errOut, "install: cursor-hook preToolUse → trace loop gate --for edit (requires .cursor/ under project)\n")
		return nil
	}

	if !hasCursorDir(root) {
		return fmt.Errorf("install: cursor-hook is CONDITIONAL — require .cursor/ under %s", root)
	}

	scriptPath := filepath.Join(root, cursorHookScriptRel)
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(scriptPath, []byte(CursorLoopGateHookScript()), 0o755); err != nil {
		return err
	}
	if err := upsertCursorHooksJSON(root); err != nil {
		return err
	}
	if err := UpsertAgentsMD(root); err != nil {
		return err
	}
	fmt.Fprintf(errOut, "install: wrote %s\n", scriptPath)
	return nil
}

func (c cursorHookTarget) Uninstall(opts InstallOpts) error {
	root := projectRoot(opts)
	scriptPath := filepath.Join(root, cursorHookScriptRel)
	_ = os.Remove(scriptPath)
	if err := removeTraceHookFromHooksJSON(root); err != nil {
		return err
	}
	return StripAgentsMD(root)
}

func hasTraceHookInHooksJSON(root string) bool {
	path := filepath.Join(root, cursorHooksJSONRel)
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(raw), "trace-loop-gate.sh")
}

func upsertCursorHooksJSON(root string) error {
	path := filepath.Join(root, cursorHooksJSONRel)
	var doc map[string]any
	if raw, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(raw, &doc); err != nil {
			return fmt.Errorf("invalid JSON in %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if doc == nil {
		doc = map[string]any{}
	}
	hooks, _ := doc["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		doc["hooks"] = hooks
	}
	pre, _ := hooks["preToolUse"].([]any)
	if pre == nil {
		pre = []any{}
	}
	entry := map[string]any{
		"command":    cursorHookCommandRel,
		"matcher":    cursorHookMatcher,
		"failClosed": false,
	}
	filtered := make([]any, 0, len(pre)+1)
	for _, item := range pre {
		m, ok := item.(map[string]any)
		if !ok {
			filtered = append(filtered, item)
			continue
		}
		cmd, _ := m["command"].(string)
		if strings.Contains(cmd, "trace-loop-gate.sh") {
			continue
		}
		filtered = append(filtered, item)
	}
	filtered = append(filtered, entry)
	hooks["preToolUse"] = filtered
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func removeTraceHookFromHooksJSON(root string) error {
	path := filepath.Join(root, cursorHooksJSONRel)
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	hooks, _ := doc["hooks"].(map[string]any)
	if hooks == nil {
		return nil
	}
	pre, _ := hooks["preToolUse"].([]any)
	if pre == nil {
		return nil
	}
	filtered := make([]any, 0, len(pre))
	for _, item := range pre {
		m, ok := item.(map[string]any)
		if !ok {
			filtered = append(filtered, item)
			continue
		}
		cmd, _ := m["command"].(string)
		if strings.Contains(cmd, "trace-loop-gate.sh") {
			continue
		}
		filtered = append(filtered, item)
	}
	if len(filtered) == len(pre) {
		return nil
	}
	hooks["preToolUse"] = filtered
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}
