# P26-S02-01 — P25-A implementer (discovery→task)

## Metadata
- id: P26-S02-01
- todo_ids: [P26-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: []
- verification: automated
- hooks: []

## Objective

Close **INT-01 + INT-06 (FM-10)**: BLOCKING discoveries can become linked tasks **only** when the agent/human supplies promotion intent (`loop apply` `spawned_tasks[]` or explicit `trace add task`). Implement PLAN.md **S02-T01–T07**. No product auto-spawn on discovery import/add.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [00-PLANNER.md](00-PLANNER.md)
- [PLAN.md](../scope-01-planning/PLAN.md) — S02 table + architecture decisions
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) — INT-01 / INT-06
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) §4 (human gate)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended run: execute these locks; do not re-debate them.

## Locked defaults (do not re-litigate)

| Item | Value |
|------|-------|
| Promotion trigger | **Primary:** `writes.spawned_tasks[]` with optional `discovery_id`. **Secondary:** `trace add task --from-discovery <id>`. **Do not** add a `promote_discoveries[]` envelope field (keep `trace.loop.apply.v1`). |
| Apply schema | Stay `trace.loop.apply.v1`. Optional `discovery_id` on `ApplySpawnedTask` is backward compatible. |
| Apply result (D1) | Keep `new_spawned_tasks` count. Add **`SpawnedTaskIDs []string`** JSON **`spawned_task_ids`** listing task UUIDs created or re-used this apply. CLI `trace loop apply` already JSON-encodes `ApplyResult` — that is D2. |
| BLOCKING filter | Helper **fail-closed** if discovery missing or severity ≠ `BLOCKING`. INFO / PLAN_AFFECTING must not promote. |
| Idempotency | If discovery already has `discovery_mentions_task` → return that task (no second row). Else `ImportSeedTask` with **task ID = discovery UUID** (PLAN architecture decision). Re-promote is a no-op (same ID, preserve `work_state`). |
| Link | After spawn, `LinkDiscoveryMentionsTask` (same rel as MCP/CLI). Mentions-only without this path still does **not** create tasks. |
| Human gate | No spawn from `ImportSeedDiscovery`, standalone `trace add discovery`, MCP `trace_add` kind=discovery, or seed import. No background goroutine / daemon / HTTP. |
| Saturation (S03) | **Do not** change `out.Saturated` / P19 zero-write rule. Discoveries-only apply still saturates until S03-T01. Promotion that inserts tasks increments `NewSpawnedTasks` as today. |
| MCP (D3 / T04) | Change **`internal/mcp/server.go` `trace_add` Description only** (plus test). Discovery first in the kind list; one sentence: after BLOCKING discovery, `trace_add` kind=task **or** `loop apply` `spawned_tasks` with `discovery_id` before product edits. Do not add a new MCP tool. Do not auto-spawn from `tools_write.go` discovery case. |
| Harness (D4) | Append **one sentence** to `GapPassPrompt` in `internal/install/gappass.go`: BLOCKING gaps → `trace add task --from-discovery` or `loop apply` `spawned_tasks` (`discovery_id`). Do **not** wholesale rewrite P25-C. Do **not** wire `ParentOrchestratorRule` (S04). |
| CLI (T03) | `trace add task --from-discovery <id>` calls the same domain helper; `--title` optional (default discovery title); `--goal-id` required if helper cannot inherit a goal from context. Help text on `add` / `loop` mentions the path. Standalone `trace add discovery` never spawns. |
| `loop next` (T05) | Optional `promotion_candidates[]` on `NextPacket` (JSON `promotion_candidates`): BLOCKING discoveries with **no** `discovery_mentions_task` target. Empty array when none. Does not spawn. |
| Schema | No SQLite migration expected. If you add a column, **stop and bump schema version** in Notes — do not ship unversioned schema. |
| Out of scope | S03 saturation/reset/STOP-reason; S04 `ParentOrchestratorRule`; new MCP promote tool; `promote_discoveries[]`; changing `promotion_blocked` (advisory / unrelated). |

## Live paths (verified 2026-08-20)

| Path | Role |
|------|------|
| `internal/domain/` (new helper, e.g. `promote.go`) | S02-T01, T07 |
| `internal/domain/seed_import.go` `ImportSeedTask` L308–334 | Upsert; preserve `work_state` on conflict |
| `internal/loop/apply.go` `ApplySpawnedTask` L165–171; spawn loop L479–497; `ApplyResult` L184–190 | T02; add `discovery_id` + result IDs |
| `internal/loop/apply.go` discoveries loop L409–428 | Must remain import-only (no auto-promote) |
| `internal/loop/next.go` `NextPacket` L51–73; task list L190–195 | T05 |
| `cmd/trace/add.go` task L68–83; discovery L104–114 | T03 |
| `cmd/trace/loop.go` `cmdLoopApply` L228–275 (stdout `Encode(res)`) | D2 — IDs appear when `ApplyResult` has them |
| `internal/mcp/server.go` L77–78 | T04 description |
| `internal/mcp/tools_write.go` L121–129, L190–191 | Discovery create + mentions-task; no auto-spawn |
| `internal/install/gappass.go` L8–11 | D4 one-line append |
| `internal/loop/apply.go` L639–649 `buildPromotionBlocked` | **Do not change** |

## Pre-conditions

- Read PLAN.md S02 + AUDIT INT-01/06.
- Baseline: `go test ./internal/...` green before edits.

## Implementation order (T01→T07)

1. **S02-T01** Domain helper `PromoteBlockingDiscovery(ctx, discoveryID, goalID) (taskID string, inserted bool, err error)` (+ unit tests). Fail-closed on non-BLOCKING. Idempotent.
2. **S02-T02** Apply: optional `discovery_id` on spawned tasks → helper; fill `spawned_task_ids`. Discoveries-only apply: **zero** new tasks. Integration test: BLOCKING discovery in store → apply with `discovery_id` → `ListTasksByGoalID` +1 and link exists.
3. **S02-T03** CLI `--from-discovery` + help. Test or CLI test in `cmd/trace`.
4. **S02-T04** MCP `trace_add` Description reorder + string assertion in `internal/mcp/*_test.go`.
5. **S02-T05** `promotion_candidates[]` on next packet + unit test (fixture BLOCKING unlinked vs linked).
6. **S02-T07** Seed import: discoveries land without spawn; promote after import uses same idempotency key; existing seed task UUIDs unchanged.
7. **S02-T06** End-to-end: BLOCKING discovery → promote via apply (or add) → task in `ListTasksByGoalID`; MCP description contains promotion wording. `go test ./internal/loop/... ./internal/mcp/... ./internal/domain/... ./internal/install/...` (and `./cmd/trace/...` if CLI tests added).

## Exit criteria

- [ ] D1: apply result includes `spawned_task_ids` when promotions occur (empty slice OK when none)
- [ ] D2: `trace loop apply` stdout JSON includes those IDs (same struct)
- [ ] D3: MCP `trace_add` description lists discovery first + promotion sentence
- [ ] D4: `GapPassPrompt` mentions `--from-discovery` or `spawned_tasks` / `discovery_id`
- [ ] D5: test BLOCKING discovery → linked task; discoveries-only apply does **not** create tasks
- [ ] T05: next packet `promotion_candidates` present (may be `[]`)
- [ ] T07: import then promote = single task row; no auto-spawn on import
- [ ] Human gate: no silent spawn
- [ ] No daemon/HTTP; no unversioned schema change; no S03/S04 edits
- [ ] `go test ./internal/...` PASS
- [ ] Board Notes: files + test commands (status + notes only)

## Minimal todos

- [ ] T01 helper + tests
- [ ] T02 apply `discovery_id` + `spawned_task_ids` + tests
- [ ] T03 `--from-discovery` + help
- [ ] T04 MCP description + test
- [ ] T05 `promotion_candidates` + test
- [ ] T07 seed import parity test
- [ ] T06 e2e + `go test ./internal/...`
- [ ] Own row `done` with evidence

## Do not

- Silent autonomous spawn
- `promote_discoveries[]` or apply schema bump
- Touch deliberation reset / saturation / STOP reason (S03)
- Wire `ParentOrchestratorRule` (S04)
- Auto-promote inside MCP `trace_add` discovery handler
