# P26-S03-01 — P25-B implementer (recalibration + reset)

## Metadata
- id: P26-S03-01
- todo_ids: [P26-S03-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: []
- verification: automated
- hooks: []

## Objective

Close **INT-02 + INT-05 + INT-09**: consecutive empty-apply saturation (≥2), deliberation reset CLI, and unified STOP reason across gate/status/export. Implement PLAN.md **S03-T01–T07**. Do not revert S02 promotion (`spawned_task_ids`, `discovery_id`, etc.).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [00-PLANNER.md](00-PLANNER.md)
- [PLAN.md](../scope-01-planning/PLAN.md) — S03 table + architecture decisions
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) — INT-02 / INT-05 / INT-09

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended run: execute these locks; do not re-debate them.

## Locked defaults (do not re-litigate)

| Item | Value |
|------|-------|
| Threshold (D1) | Named const **`SaturationEmptyThreshold = 2`** in `internal/deliberation` (next to `HopBudget`). First pure-empty apply must **not** sticky-STOP greenfield. |
| Empty increment | Pure empty: `NewPlanChanges==0 && NewSpawnedTasks==0` **and** no `writes.discoveries[]` imported this apply → increment `consecutive_empty_applies`. |
| Discoveries-only | **Non-saturating:** discoveries imported with zero plan/spawn → do **not** increment counter; do **not** set `P19Saturated` / `out.Saturated` for transition. |
| Non-empty clears counter | `NewPlanChanges>0` **or** `NewSpawnedTasks>0` → set `consecutive_empty_applies = 0`. |
| Saturation fire | `consecutive_empty_applies >= SaturationEmptyThreshold` → `P19Saturated=true` → STOP with reason **`p19_saturated`**. |
| `MaxIterationsReached` | Immediate saturate/STOP (bypass consecutive counter; unchanged). |
| Persistence (T02) | New column **`consecutive_empty_applies`** on `deliberation_state` + migration **`028_*.sql`**. Do **not** derive from loop-step history (reset cannot clear audit steps without re-STOP). Bump embed expectations (store tests + CLI). |
| HopBudget (T07) | Stay **`HopBudget = 12`**. Do not raise. |
| Reset API (D2 / T03) | Domain `ResetDeliberationState(taskID)` + CLI **`trace loop reset --task <id>`**. Clear `Stopped`, `StopReason`, `HopCount`, `consecutive_empty_applies`; set `CurrentPhase=EXECUTE`; **preserve** `PlanCritiqued`. No MCP reset. |
| Reset vs gate | Gate uses **fresh** `SelectNext` (`gate.go`). Reset must prevent immediate re-STOP from saturation. Edit allow still needs `SelectNext` → EXECUTE (`ExecutePending` / `plan_critiqued`). G1 with `plan_critiqued=false` may still block edit until critique — that is OK; assert no sticky `p19_saturated` re-STOP after reset alone. |
| INT-09 (D3) | When `Stopped=true`, `SelectNext` returns **persisted** `StopReason` (canonical saturation string `p19_saturated`). `hop_budget_exceeded` **only** when `HopCount >= HopBudget` (or empty `StopReason` fallback). Status JSON adds **`stop_reason`** alongside `why_selected`. Gate `reason_code` == export `stop_reason` for same task. |
| Replay / status / next | Replay `Saturated`, `p19SaturatedFromLastStep`, status `Saturated`, and next/gate policy inputs must use the **consecutive** rule (or persisted counter), not single last-step zero-write. |
| Schema | Migration 028 required for the new column. Update `internal/store/deliberation_test.go` embed/max (**27→28**) and **`cmd/trace/cli_test.go` `TestMigrateBackupAuthCLI`** (stale expect embed/max **14**; live was 27 — set to **28** with this migration). |
| Out of scope | S02 promotion (`PromoteBlockingDiscovery`, `discovery_id`, MCP description, `GapPassPrompt`); S04 `ParentOrchestratorRule`; daemon/HTTP; MCP reset tool. |

## Live paths (verified 2026-08-20, post-S02)

| Path | Role |
|------|------|
| `internal/deliberation/types.go` L8–9 `HopBudget`; L35–37 reason consts; `State` L71–81 | T01 const; T02 field; T07 |
| `internal/deliberation/select.go` L7–11 sticky STOP + P19 branches; L49–87 `ApplyTransition` | T01, T05 |
| `internal/deliberation/select_test.go` (hop_budget / stopped cases) | T01, T05, T07 |
| `internal/loop/policy.go` L103–109 `p19SaturatedFromLastStep` | T01 — replace last-step-only rule |
| `internal/loop/apply.go` L400–407 replay Saturated; L418+ discoveries loop; L526–538 saturate + transition; L621 status saturated; L679–698 `buildStatusDeliberation` (no `stop_reason` today) | T01, T04, T05 |
| `internal/loop/gate.go` L152–159, L188–192 | T04, T05 — reason from SelectNext |
| `internal/loop/deliberation_packet.go` L160–169 `StatusDeliberation`; L182–203 next packet dual labels | T05 |
| `internal/loop/next.go` L235 `p19SaturatedFromLastStep` | T01 |
| `internal/domain/deliberation.go` L14–72 transition (forward-only); store mappers L84–109 | T02, T03 |
| `internal/store/deliberation.go` L8–19 struct; L28–60 upsert; schema `015_deliberation_state.sql` | T02 + 028 |
| `internal/store/schema/` (embed max **27** today) | Add **028**; bump tests |
| `cmd/trace/loop.go` L36–56 dispatch (no `reset` yet) | T03 CLI |
| `internal/domain/seed_export.go` L710–719 export `stop_reason` | T06 (already source-correct) |
| `internal/domain/seed_import.go` L936+ `ImportSeedDeliberationState` | Preserve/import new column if exported; reset remains product path after import |
| S02 surfaces | `ApplyResult.spawned_task_ids`, `spawned_tasks[].discovery_id` — **do not revert** |

## Pre-conditions

- Read PLAN.md S03 + AUDIT INT-02/05/09.
- Baseline: `go test ./internal/...` green before edits.
- Do not edit S02 promotion or S04 installer wiring.

## Implementation order (T01→T07)

1. **S03-T01** Add `SaturationEmptyThreshold = 2`. Change apply + policy so first pure-empty does not STOP; second consecutive pure-empty does. Discoveries-only does not increment. Unit tests: 1 empty → not stopped; 2 empty → `p19_saturated`; discoveries-only → not increment.
2. **S03-T02** Migration **028**: column `consecutive_empty_applies INTEGER NOT NULL DEFAULT 0` on `deliberation_state`. Wire store/domain State + upsert/scan/list/export/import as needed. Counter survives round-trip; clears on non-empty apply. Document version bump in Notes.
3. **S03-T03** Domain `ResetDeliberationState` + CLI `trace loop reset --task <id>` (help + dispatch). Clears stopped/hop/counter; phase EXECUTE; preserve plan_critiqued. Unit + CLI/domain tests.
4. **S03-T04** Integration: saturate → reset → discoveries-only / first empty does **not** immediate re-STOP. Status may hint recovery when saturated (optional one-liner; not a new subsystem).
5. **S03-T05** Reorder `SelectNext`: stopped → persisted `StopReason`; hop budget only when hop ≥ budget. Add `stop_reason` on status deliberation JSON. Align `why_selected` / gate `reason_code` with persisted reason when stopped.
6. **S03-T06** Test: after saturation STOP, export `stop_reason` == gate `reason_code` (`p19_saturated`).
7. **S03-T07** Keep `HopBudget = 12`; update `select_test.go` for new branch order. Fix embed expectations (**28**) including **`TestMigrateBackupAuthCLI`**.

## Exit criteria

- [ ] D1: `SaturationEmptyThreshold = 2`; first pure-empty no sticky STOP; second consecutive empties STOP with `p19_saturated`
- [ ] D2: `trace loop reset --task <id>` clears stopped, hop_count, consecutive empty; persisted phase EXECUTE
- [ ] D3: gate `reason_code` == export `stop_reason` (same constant) after saturation STOP; status exposes `stop_reason`
- [ ] D4: tests for first/second empty, discoveries-only non-increment, reset, reason alignment
- [ ] Migration 028 + embed tests updated (incl. `TestMigrateBackupAuthCLI` → 28)
- [ ] HopBudget remains 12; no S02/S04 edits; no daemon/HTTP
- [ ] `go test ./internal/...` PASS; `go test ./cmd/trace/...` PASS if CLI/migrate tests touched
- [ ] Board Notes: files + test commands (status + notes only)

## Minimal todos

- [ ] T01 saturation threshold + tests
- [ ] T02 column + migration 028 + store/domain wire
- [ ] T03 reset domain + CLI + tests
- [ ] T04 reset ↔ gap-pass / no immediate re-STOP
- [ ] T05 SelectNext + status `stop_reason`
- [ ] T06 export/gate alignment test
- [ ] T07 HopBudget docs/tests + embed 28 / CLI migrate test
- [ ] Own row `done` with evidence

## Do not

- Revert or rewrite S02 promotion (`spawned_task_ids`, `discovery_id`, MCP/GapPassPrompt)
- Wire `ParentOrchestratorRule` (S04)
- Raise `SaturationEmptyThreshold` above 2 or `HopBudget` above 12
- Derive consecutive-empty solely from loop-step history
- Add MCP reset or HTTP/daemon
- Ship column without migration 028
