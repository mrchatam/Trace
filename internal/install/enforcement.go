package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	enforcementBeginMarker = "# begin-trace-enforcement"
	enforcementEndMarker   = "# end-trace-enforcement"

	cursorRulesRelPath  = ".cursor/rules/trace-enforcement.mdc"
	claudeFallbackRules = ".claude/trace-enforcement-rules.md"
	cursorHookScriptRel = ".cursor/hooks/trace-loop-gate.sh"
	cursorHooksJSONRel  = ".cursor/hooks.json"
)

// ParentOrchestratorRule documents the parent-orchestrator TRACE_TASK_ID ownership
// contract (INT-04 / FM-04). This is a document-level constant; the enforcement
// implementation is the existing CursorLoopGateHookScript preToolUse path.
//
// failClosed policy (Option A): When TRACE_TASK_ID is absent and
// .trace/config.json has enforce=strict, deny the edit. off/warn/missing
// preserve default-off allow. Option B (parent-orchestrator detection) is out.
const ParentOrchestratorRule = `Parent orchestrator TRACE_TASK_ID ownership (INT-04 / FM-04):

1. The parent orchestrator MUST set TRACE_TASK_ID to the active seed task UUID
   before any product-code edit and before delegating work to any subagent.
2. Parent owns the Trace graph for the active task: gap pass, discoveries, and
   decisions. Do not complete an edit path by offloading graph-only work to
   workers while the parent edits without TRACE_TASK_ID / loop gate.
3. Worker inheritance: before each worker, set TRACE_TASK_ID and
   TRACE_PROJECT_ROOT; put workspace path + task UUID in every worker prompt;
   workers must export the same env before product edits. Do not assume Cursor
   Multitask inherits parent env automatically.
4. preToolUse deny fires when a Write/edit is attempted without an active
   TRACE_TASK_ID under enforce=strict — enforced via CursorLoopGateHookScript
   (.cursor/hooks/trace-loop-gate.sh). Option A applies per process to that
   process's TRACE_TASK_ID.
5. failClosed (Option A): When TRACE_TASK_ID is absent AND .trace/config.json
   enforce=strict, deny the edit rather than allowing untracked work.
   off/warn/missing/invalid enforce still allow (default-off preserved).
6. Multitask limit: Trace cannot product-enforce worker env inheritance or
   detect "parent orchestrator" (Option B deferred). Rules + Option A hook are
   the harness choke points; orchestrators still verify board order.

Implementation path: see CursorLoopGateHookScript in internal/install/enforcement.go.`

// EnforcementRulesMarkdown returns the shared enforcement rules body (stderr hint / docs).
func EnforcementRulesMarkdown() string {
	return `## Trace enforcement (installed by trace install)

When ` + "`TRACE_TASK_ID`" + ` is set for the active seed task:

1. Before product code edits, run ` + "`trace loop gate --task \"$TRACE_TASK_ID\" --for edit`" + `.
2. New goal without coarse plan: bootstrap via ` + "`trace plan create-coarse`" + ` or MCP ` + "`trace_plan`" + ` before edit gate can pass (or ` + "`trace plan bootstrap --goal <id>`" + ` for recovery).
3. Before DONE, run ` + "`trace loop status --task \"$TRACE_TASK_ID\"`" + ` and resolve any ` + "`violations[]`" + `.
4. Use ` + "`--enforce`" + ` on transition DONE and ` + "`seed export --strict --enforce`" + ` only when the team opts in.

Environment: ` + "`TRACE_TASK_ID`" + ` (active task UUID), ` + "`TRACE_PROJECT_ROOT`" + ` (optional -C root).
See docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md. Local config: ` + "`.trace/config.json`" + ` (` + "`enforce`" + `, default off).
`
}

// AgentsEnforcementBlock returns the AGENTS.md / CLAUDE.md marker-delimited block.
func AgentsEnforcementBlock() string {
	return enforcementBeginMarker + `
## Trace enforcement (harness)

When a Trace seed task is active, set ` + "`TRACE_TASK_ID`" + ` to its UUID.

- **Before edits:** ` + "`trace loop gate --task \"$TRACE_TASK_ID\" --for edit`" + ` (exit 0 = proceed).
- **Coarse plan:** New goal without progressive plan — ` + "`trace plan create-coarse`" + ` + ` + "`set-current`" + ` + ` + "`deep`" + `, or MCP ` + "`trace_plan`" + `; recovery: ` + "`trace plan bootstrap --goal <id>`" + `.
- **Before DONE:** ` + "`trace loop status --task \"$TRACE_TASK_ID\"`" + ` — resolve non-empty ` + "`violations[]`" + `.
- **Opt-in strict:** ` + "`--enforce`" + ` on ` + "`trace transition … --to DONE`" + `; ` + "`trace seed export --strict --enforce`" + ` for CI.
- **Config:** ` + "`.trace/config.json`" + ` → ` + "`{ \"enforce\": \"off\"|\"warn\"|\"strict\" }`" + ` (default off).

Product design: docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md
` + "\n" + GapPassPrompt + "\n" + ParentOrchestratorRule + "\n" + enforcementEndMarker + "\n"
}

func cursorRulesMDCContent() string {
	return `---
description: Trace loop gate — pre-edit enforcement for active tasks
alwaysApply: true
---

## Trace enforcement (installed by trace install cursor --write)

When ` + "`TRACE_TASK_ID`" + ` is set for the active seed task:

1. Before product code edits, run ` + "`trace loop gate --task \"$TRACE_TASK_ID\" --for edit`" + `.
2. New goal without coarse plan: bootstrap via ` + "`trace plan create-coarse`" + ` or MCP ` + "`trace_plan`" + ` before edit gate can pass (or ` + "`trace plan bootstrap --goal <id>`" + ` for recovery).
3. Before DONE, run ` + "`trace loop status --task \"$TRACE_TASK_ID\"`" + ` and resolve any ` + "`violations[]`" + `.
4. Use ` + "`--enforce`" + ` on transition DONE and ` + "`seed export --strict --enforce`" + ` only when the team opts in.

See docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md. Local config: ` + "`.trace/config.json`" + ` (` + "`enforce`" + `, default off).
` + "\n" + GapPassPrompt + "\n" + ParentOrchestratorRule + "\n"
}

func claudeFallbackRulesContent() string {
	return `# Trace enforcement rules (installed by trace install claude --write)

When ` + "`TRACE_TASK_ID`" + ` is set for the active seed task:

1. Before product code edits, run ` + "`trace loop gate --task \"$TRACE_TASK_ID\" --for edit`" + `.
2. New goal without coarse plan: bootstrap via ` + "`trace plan create-coarse`" + ` or MCP ` + "`trace_plan`" + ` before edit gate can pass (or ` + "`trace plan bootstrap --goal <id>`" + ` for recovery).
3. Before DONE, run ` + "`trace loop status --task \"$TRACE_TASK_ID\"`" + ` and resolve any ` + "`violations[]`" + `.
4. Use ` + "`--enforce`" + ` on transition DONE and ` + "`seed export --strict --enforce`" + ` only when the team opts in.

Environment: ` + "`TRACE_TASK_ID`" + `, ` + "`TRACE_PROJECT_ROOT`" + ` (optional).
See docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md.
` + "\n" + GapPassPrompt + "\n" + ParentOrchestratorRule + "\n"
}

// CursorLoopGateHookScript is the preToolUse shell script body.
// Option A failClosed: empty TRACE_TASK_ID + enforce=strict → deny; otherwise
// empty task allows (default-off). Task set → existing loop gate path.
func CursorLoopGateHookScript() string {
	return `#!/usr/bin/env bash
set -euo pipefail
input="$(cat)"
task_id="${TRACE_TASK_ID:-}"
root="${TRACE_PROJECT_ROOT:-$PWD}"
if [[ -z "$task_id" ]]; then
  # Option A: failClosed only when enforce=strict (mirror LoadEnforceMode fail-open).
  cfg="$root/.trace/config.json"
  if [[ -f "$cfg" ]] && grep -q '"enforce"[[:space:]]*:[[:space:]]*"strict"' "$cfg" 2>/dev/null; then
    echo '{"permission":"deny","user_message":"TRACE_TASK_ID required when .trace/config.json enforce=strict. Set TRACE_TASK_ID to the active seed task UUID before editing.","agent_message":"Set TRACE_TASK_ID to the active seed task UUID before editing product code."}'
    exit 2
  fi
  echo '{"permission":"allow"}'
  exit 0
fi
if ! command -v trace >/dev/null 2>&1; then
  echo '{"permission":"allow"}'
  exit 0
fi
if trace -C "$root" loop gate --task "$task_id" --for edit >/dev/null 2>&1; then
  echo '{"permission":"allow"}'
  exit 0
fi
msg="Trace loop gate blocked edit for task ${task_id}. Run: trace loop gate --task ${task_id} --for edit"
# JSON string escape (quotes/backslashes); avoid bash %q which is not JSON.
msg_json=${msg//\\/\\\\}
msg_json=${msg_json//\"/\\\"}
printf '{"permission":"deny","user_message":"%s","agent_message":"Follow recommended_phase from gate JSON before editing product code."}\n' "$msg_json"
exit 2
`
}

// UpsertAgentsMD merges the enforcement block into AGENTS.md under root.
func UpsertAgentsMD(root string) error {
	path := filepath.Join(root, "AGENTS.md")
	var existing []byte
	if raw, err := os.ReadFile(path); err == nil {
		existing = raw
	} else if !os.IsNotExist(err) {
		return err
	}
	block := AgentsEnforcementBlock()
	updated := UpsertMarkedFragment(string(existing), enforcementBeginMarker, enforcementEndMarker, block)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(updated), 0o644)
}

// StripAgentsMD removes the enforcement block from AGENTS.md.
func StripAgentsMD(root string) error {
	path := filepath.Join(root, "AGENTS.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	updated := RemoveMarkedFragment(string(raw), enforcementBeginMarker, enforcementEndMarker)
	if updated == string(raw) {
		return nil
	}
	if strings.TrimSpace(updated) == "" {
		return os.Remove(path)
	}
	return os.WriteFile(path, []byte(updated), 0o644)
}

// UpsertClaudeRules writes enforcement rules into CLAUDE.md or fallback file.
func UpsertClaudeRules(root string) error {
	claudeMD := filepath.Join(root, "CLAUDE.md")
	if _, err := os.Stat(claudeMD); err == nil {
		raw, err := os.ReadFile(claudeMD)
		if err != nil {
			return err
		}
		block := AgentsEnforcementBlock()
		updated := UpsertMarkedFragment(string(raw), enforcementBeginMarker, enforcementEndMarker, block)
		return os.WriteFile(claudeMD, []byte(updated), 0o644)
	}
	claudeDir := filepath.Join(root, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(root, claudeFallbackRules)
	return os.WriteFile(path, []byte(claudeFallbackRulesContent()), 0o644)
}

// StripClaudeRules removes enforcement rules from CLAUDE.md or fallback file.
func StripClaudeRules(root string) error {
	claudeMD := filepath.Join(root, "CLAUDE.md")
	if raw, err := os.ReadFile(claudeMD); err == nil {
		updated := RemoveMarkedFragment(string(raw), enforcementBeginMarker, enforcementEndMarker)
		if updated != string(raw) {
			if strings.TrimSpace(updated) == "" {
				return os.Remove(claudeMD)
			}
			return os.WriteFile(claudeMD, []byte(updated), 0o644)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	path := filepath.Join(root, claudeFallbackRules)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// WriteCursorRulesMDC writes .cursor/rules/trace-enforcement.mdc.
func WriteCursorRulesMDC(root string) error {
	path := filepath.Join(root, cursorRulesRelPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(cursorRulesMDCContent()), 0o644)
}

// RemoveCursorRulesMDC deletes the Cursor rules file.
func RemoveCursorRulesMDC(root string) error {
	path := filepath.Join(root, cursorRulesRelPath)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func hasCursorDir(root string) bool {
	cursorDir := filepath.Join(root, ".cursor")
	st, err := os.Stat(cursorDir)
	return err == nil && st.IsDir()
}

func enforcementPrintHint(errOut interface{ Write([]byte) (int, error) }) {
	if errOut == nil {
		return
	}
	fmt.Fprint(errOut, EnforcementRulesMarkdown())
}
