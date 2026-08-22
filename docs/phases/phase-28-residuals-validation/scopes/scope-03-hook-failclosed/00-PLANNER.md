# P28-S03-00 — Scope planner (hook failClosed)

## Metadata
- id: P28-S03-00
- todo_ids: [P28-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, security-and-hardening]
- verification: automated
- hooks: []

## Objective

Plan hook failClosed hardening for INT-04/11 beyond install text (R2, R3, R8). Lock implementer defaults so strict enforce denies edits when `TRACE_TASK_ID` is absent.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R2/R3/R8
- [Phase 28 README](../../README.md)
- `internal/install/enforcement.go` — `CursorLoopGateHookScript()` L106–108 allow path
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) — FM-05

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Problem (R2, R3, R8)

Phase 25/26 wired **install text** for INT-04 (Parent orchestrator) but hook still **allows** edits when `TRACE_TASK_ID` is absent (`CursorLoopGateHookScript` L106–108). FM-05: strict config without deny permits post-STOP untracked edits.

INT-11: no automated hook drift check on Cursor schema changes.

## Locked defaults (S03-01 implements)

| Item | Value |
|------|-------|
| **Preferred (A)** | When `.trace/config.json` enforce=strict AND no TRACE_TASK_ID → deny (failClosed) |
| Option B | Deny only when parent orchestrator rule detected (harder — defer unless audit requires) |
| INT-11 | Add `internal/install/hook_drift_test.go` or golden JSON schema check |
| Scope files | `internal/install/cursorhook.go`, `enforcement.go`, hook script body |
| Out of scope | Rewriting Cursor Multitask product; daemon/HTTP |

## Planner gate

- [x] `RESIDUAL-AUDIT.md` R2/R3/R8 seeds present
- [x] `01-implement.md` + `02-review.md` runnable
- [x] Default-off enforce projects preserved (locked: off/warn/missing + empty task → allow)

## Exit criteria

- [x] Implementer prompt locked for fresh subagent
- [x] Board row P28-S03-00 Notes cite option A lock
- [x] Next runnable **P28-S03-01**

## Planner completion (2026-08-20)

**Locked Option A:** `.trace/config.json` `enforce=strict` AND no `TRACE_TASK_ID` → deny; Cursor hooks.json `failClosed` stays `false` (script owns policy). Option B deferred. INT-11 → `hook_drift_test.go`. Thickened `01-implement.md`, `02-review.md`, `SCOPE-TODOS.md`.

## Todo updates

Status + notes on **P28-S03-00** only.

## Next

`P28-S03-01`
