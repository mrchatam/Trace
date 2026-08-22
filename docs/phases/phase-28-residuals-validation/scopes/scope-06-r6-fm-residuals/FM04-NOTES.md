# FM-04 / FR-P28-03 — notes (P28-S06-05)

**Date:** 2026-08-20  
**Status:** implemented (Option A intact; Option B not shipped)

## What shipped

1. **`ParentOrchestratorRule` (INT-04 / FM-04)** — parent must set `TRACE_TASK_ID` before edits and before delegation; parent owns gap pass / discoveries / decisions (no graph-only offload to workers while parent lacks task); worker inheritance via prompt + explicit env; Multitask product limit documented (no env-inheritance guarantee; Option B deferred).
2. **`AgentsEnforcementBlock`** — now includes `ParentOrchestratorRule` (closes AGENTS vs Cursor asymmetry for FM-04).
3. **Harness** — `PROTOCOL.md` FM-04 paragraph + Multitask step; `PROMPT-G1-BUILD.md` Multitask section; `ENFORCEMENT.md` Multitask/worker inheritance note under harness install.

## Multitask limits (product-unfixable)

Cursor Multitask does not guarantee workers inherit parent `TRACE_TASK_ID`. Trace cannot detect “parent orchestrator” without Option B (FR-P28-D4 deferred). Acceptance is therefore:

- Rules + install text require parent-must-set-task and worker prompt+env inheritance.
- Option A hook still denies empty-task under `enforce=strict` **per process**.
- Scripted coverage: install text assertions + existing Option A deny/allow tests.

## Evidence

```bash
GOPROXY=direct go test ./internal/install/... -count=1 -run 'CursorLoopGate|HookDrift|ParentOrchestrator|InstallCursorIncludes|InstallAgentsMDEnforcement'
GOPROXY=direct go test ./internal/install/... -count=1
```

Primary coverage: `TestParentOrchestratorRuleNonEmpty` (graph-offload ban, Worker inheritance, Multitask, Option B note, AGENTS block); `TestInstallAgentsMDEnforcementBlock` / `TestInstallCursorIncludesLoopGateRule` include Parent orchestrator (+ Worker inheritance on Cursor rules).
