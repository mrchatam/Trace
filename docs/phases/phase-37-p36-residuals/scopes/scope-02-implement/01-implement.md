# P37-S02-01 — Implement

## Metadata
- id: P37-S02-01
- todo_ids: [P37-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, api-and-interface-design]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Implement the **8 accepted S02 residuals** (R1–R6, R8, R11) from [PLAN.md](../scope-01-plan/PLAN.md) in **wave order A→D**, preserving Phase 36 guarantees (MCP `trace_plan`, bootstrap, terminal advisory, active `plan_missing` block). **No** silent PlanExists bridge; **no** enforce default flip (R7).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [PLAN.md](../scope-01-plan/PLAN.md) — **SoT** for accept set, touch-list, waves, acceptance tests, P36 regression subset
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [RESIDUALS.md](../scope-00-triage/RESIDUALS.md)
- [00-PLANNER.md](00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22):
  - `internal/loop/apply.go:196–209` — `StatusResult` (no `advisories[]` yet)
  - `internal/planner/advisory.go:13–82` — `GoalStructureWarning`, `goalLinkedPlanChangeIDs`
  - `internal/loop/policy.go:45–48` — `PlanExists` from store read only (**no edit**)
  - `internal/mcp/tools_loop.go:17–47` — `trace_loop` actions `next|apply|status` only
  - `cmd/trace/loop.go:155–179` — CLI gate via `loop.EvaluateGate` + JSON envelope
  - `internal/httpapi/server.go:258,282` — `GET /v1/loop/gate`, `GET /v1/plans` only
  - `internal/httpapi/handlers_loop.go:104+` — HTTP gate precedent (Law 19)
  - `cmd/trace/plan.go:351+` — CLI `plan bootstrap` → `planner.Service.BootstrapFromPlanChanges`
  - `internal/config/enforce.go:43–57` — `WarnIfTraceDirWithoutConfig` (helper shipped; test gap)
  - `web/src/screens/Overview.tsx:72–88` — GateStrip + status violations only (no advisories banner)
  - `web/src/screens/TaskDetail.tsx:205–211` — bootstrap paragraph (P36 partial R8 — **unchanged**)
  - `internal/mcp/mcp_test.go:1210+` — `TestGreenfield_MCPPlanBootstrap_EditGatePasses` (R11 canonical)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| SoT | [PLAN.md](../scope-01-plan/PLAN.md) §2 accept, §4 touch-list, §5 tests, §7 non-goals, §8 waves |
| Scope | R1–R6, R8, R11 only — R7/R9/R8-full re-deferred; R10 is S03 |
| Wave order | **Strict A→B→C→D** — see §Implementation waves below |
| Law 19 | Policy in library; MCP/HTTP/GUI are thin adapters — single service call paths |
| R1 | **Advisory only** — `bootstrap_recommended` on `advisories[]`; **NEVER** set `PlanExists=true`; N=**1** linked plan-change |
| R5 | `advisories[]` on `StatusResult`; code `goal_structure_warning`; orthogonal to R1 |
| R3 | MCP `trace_loop action=gate` mirrors CLI JSON + exit semantics; **no** gate logic fork |
| R2 | `POST /v1/plans/bootstrap` → `planner.Service.BootstrapFromPlanChanges`; OpenAPI updated |
| R4 | Help/bootstrap stderr: bootstrap yields **minimal** plan; refinement via create-coarse/deep expected |
| R6 | Unit test only — **do not** flip `LoadEnforceMode` default (R7 stays `EnforceOff`) |
| R8 | Overview minimal surface — consume HTTP JSON only; TaskDetail paragraph unchanged |
| R11 | Doc path only — critique via `trace loop apply` + plan_changes; **no** new MCP tool |
| MCP tool count | **16** — `RegisteredToolNames()` unchanged count; extend `trace_loop` schema only |
| Dogfood fixture | `/home/ali/Desktop/feet seller telegram app` — **read-only** in S02 |
| Graph export | If Trace entities change: `trace seed export -o trace/graph.json` per CONTRIBUTING |

## Accept set (do not re-debate)

| ID | Wave | Summary | Partial P36? |
|----|------|---------|--------------|
| **R5** | A | Wire `GoalStructureWarning` → `advisories[]` code `goal_structure_warning` | Helper exists; status channel gap |
| **R1** | A | `bootstrap_recommended` when ≥1 linked plan-change, no plan | Full gap |
| **R3** | B | MCP `trace_loop action=gate` | Full gap |
| **R2** | B | HTTP `POST /v1/plans/bootstrap` + OpenAPI | Full gap |
| **R4** | B | Bootstrap help refinement note | Help gap only |
| **R6** | C | `TestWarnIfTraceDirWithoutConfig` | Helper exists; test gap |
| **R11** | C | Agent workflow doc (loop apply critique path) | Block 0 test exists |
| **R8** | D | Overview plan-gap / advisories banner | TaskDetail partial only |

## Implementation waves (strict — dependency order)

```text
Wave A — status schema (blocks R8 data contract):
  R5  advisories[] + GoalStructureWarning wire
  R1  bootstrap_recommended on same channel

Wave B — adapters (after A stable):
  R3  MCP trace_loop action=gate
  R2  HTTP POST /v1/plans/bootstrap (+ OpenAPI)
  R4  bootstrap help refinement note

Wave C — tests + docs:
  R6  WarnIfTraceDirWithoutConfig unit test
  R11 agent-workflow / verify doc path

Wave D — GUI (after A):
  R8  Overview minimal plan-gap / advisories surface
```

## Touch-list (ordered — library → MCP → HTTP → install → GUI)

| Step | File | Action | Residuals |
|------|------|--------|-----------|
| 1 | `internal/planner/advisory.go` | **Edit** — bootstrap advisory builder (reuse `goalLinkedPlanChangeIDs`) | R1 |
| 2 | `internal/planner/advisory_test.go` | **Edit** — tests for bootstrap + goal-structure codes | R1, R5 |
| 3 | `internal/loop/apply.go` | **Edit** — add `Advisories []Advisory` to `StatusResult`; assemble in status path | R5, R1 |
| 4 | `internal/loop/apply_test.go` or `cmd/trace/loop_test.go` | **Edit** — R1/R5 status JSON tests | R1, R5 |
| 5 | `internal/mcp/tools_loop.go` | **Edit** — add `gate` action + helper mirroring CLI | R3 |
| 6 | `internal/mcp/mcp_test.go` | **Edit** — `TestMCPLoopGate_MatchesCLI`; update `LoopInput` jsonschema | R3 |
| 7 | `internal/httpapi/handlers_p1.go` | **Edit** — `POST /v1/plans/bootstrap` handler | R2 |
| 8 | `internal/httpapi/server.go` | **Edit** — route registration | R2 |
| 9 | `api/openapi.yaml` | **Edit** — path + request/response schemas | R2 |
| 10 | `internal/httpapi/handlers_p1_test.go` (or integration test file) | **Edit/Create** — `TestHTTPPlanBootstrap_CreatesPlannerRows` | R2 |
| 11 | `cmd/trace/plan.go` | **Edit** — refinement sentence in `printPlanHelp` + bootstrap stderr | R4 |
| 12 | `cmd/trace/plan_test.go` or `help_test.go` | **Edit** — `TestPlanHelp_MentionsRefinement` | R4 |
| 13 | `internal/config/enforce_test.go` | **Edit** — `TestWarnIfTraceDirWithoutConfig` | R6 |
| 14 | `docs/rules/agent-loop-protocol.md` and/or `docs/phases/phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md` | **Edit** — R11 critique workflow cross-ref | R11 |
| 15 | `web/src/screens/Overview.tsx` | **Edit** — advisory/plan-gap banner from gate or status `advisories[]` | R8 |
| 16 | `web/src/screens/Overview.test.tsx` (if Vitest present) | **Edit/Create** — advisory copy visible | R8 |

**Explicit non-touch:**

- `internal/loop/policy.go` — no PlanExists bridge heuristic
- `LoadEnforceMode` default path — no R7 flip
- `web/src/screens/TaskDetail.tsx` — bootstrap paragraph unchanged unless copy-only tweak explicitly needed
- New MCP `critique-seed` tool
- Hosted SaaS routes; full plan tree GUI (R8-full)

## Implementation order (strict — align steps 1–16)

```text
1. internal/planner/advisory.go (+ advisory_test.go) — bootstrap advisory helper
2. internal/loop/apply.go (+ tests) — StatusResult.advisories[]
3. internal/mcp/tools_loop.go (+ mcp_test.go) — gate action
4. internal/httpapi/handlers_p1.go, server.go, openapi.yaml (+ tests) — POST bootstrap
5. cmd/trace/plan.go (+ help test) — R4 refinement note
6. internal/config/enforce_test.go — R6
7. docs verify/agent-loop cross-ref — R11
8. web/src/screens/Overview.tsx (+ test) — R8
9. Full test pass + graph export if entities changed
```

## Role work

### Wave A — R5 + R1 (`advisories[]`)

1. Define `Advisory` type (code + message; snake_case JSON matching `trace.loop.status.v1`).
2. In status assembly (`apply.go`), call:
   - `planner.GoalStructureWarning` → append `{code: "goal_structure_warning", ...}` when non-empty
   - New bootstrap helper → append `{code: "bootstrap_recommended", ...}` when `!PlanExists && len(goalLinkedPlanChangeIDs) ≥ 1`
3. **Preserve** `violations[]`, `schema_version`, deliberation `policy_inputs.plan_exists` — advisory channel is orthogonal.
4. **R1 guard:** `plan_exists` in status deliberation must remain **false** when only plan-changes exist.

### Wave B — R3 MCP gate

1. Extend `LoopInput.Action` jsonschema: `next|apply|status|gate`.
2. Add `gate` branch calling same `loop.EvaluateGate` path as `cmd/trace/loop.go:155–179`.
3. Mirror CLI JSON envelope on stdout; match exit/blocked semantics in MCP result.
4. HTTP `GET /v1/loop/gate` is precedent only — do not duplicate gate logic.

### Wave B — R2 HTTP bootstrap

1. Thin handler: parse `goal_id`, call `planner.Service.BootstrapFromPlanChanges`.
2. Register `POST /v1/plans/bootstrap` in `server.go`.
3. Update `api/openapi.yaml` with path, request body, 200 response — document alongside existing `GET /v1/plans`.

### Wave B — R4 help

1. Add sentence to `printPlanHelp` and/or bootstrap stderr: bootstrap produces **minimal** plan; human refinement via `create-coarse` / `deep` expected.
2. No LLM generation claims; do not hide bootstrap limits.

### Wave C — R6 test

1. Temp dir with `.trace/` but no valid config → capture stderr from `WarnIfTraceDirWithoutConfig`.
2. Assert nudge substring (enforce warn suggestion) present.

### Wave C — R11 doc

1. Document post-bootstrap critique path: **`trace loop apply`** with plan_changes envelope.
2. Cite `TestGreenfield_MCPPlanBootstrap_EditGatePasses` as canonical Block 0 pattern.
3. **Reject path:** no new MCP `critique-seed` tool.

### Wave D — R8 Overview

1. When active task gate shows plan-gap advisory **or** loop status includes non-empty `advisories[]`, show minimal banner/copy.
2. Consume `GET /v1/loop/gate` + loop status API only — Law 19.
3. Reuse `GateStrip` warn path where applicable; do not duplicate TaskDetail bootstrap paragraph logic.

## MCP / OpenAPI update expectations

| Residual | MCP | OpenAPI |
|----------|-----|---------|
| **R5** | `trace_loop action=status` JSON gains `advisories[]` (pass-through from library) | `LoopStatusResponse` / status schema gains `advisories` array if documented |
| **R1** | Same status channel — `bootstrap_recommended` entries | Same |
| **R3** | `trace_loop` tool schema: add `gate` to action enum; gate params mirror CLI `--for` | No new path (HTTP gate already exists) |
| **R2** | No new MCP tool | **Required:** `POST /v1/plans/bootstrap` path + schemas |
| **R4** | N/A | N/A |
| **R6** | N/A | N/A |
| **R11** | N/A | N/A |
| **R8** | N/A | Optional: ensure status/gate response schemas include `advisories` for GUI typing |

**MCP tool count:** remains **16** — extend `trace_loop` only; do not add tools.

## Test strategy (required)

### S02 acceptance tests (PLAN §5 — must pass)

| Name | File | Residual | Assert shape |
|------|------|----------|--------------|
| `TestLoopStatus_IncludesGoalStructureAdvisory` | `cmd/trace/loop_test.go` or `internal/loop/apply_test.go` | R5 | >15 tasks, no plan → `advisories[]` contains `goal_structure_warning` |
| `TestLoopStatus_BootstrapRecommendedAdvisory` | same | R1 | ≥1 linked plan-change, no plan → `advisories[]` contains `bootstrap_recommended` |
| `TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists` | same | R1 guard | Same setup — `deliberation.policy_inputs.plan_exists` still **false** |
| `TestMCPLoopGate_MatchesCLI` | `internal/mcp/mcp_test.go` | R3 | `trace_loop action=gate` envelope matches CLI; blocked edit → violations present |
| `TestHTTPPlanBootstrap_CreatesPlannerRows` | `internal/httpapi/*_test.go` | R2 | POST + goal_id → 200; planner rows exist |
| `TestPlanHelp_MentionsRefinement` | `cmd/trace/plan_test.go` or `help_test.go` | R4 | Help/bootstrap text mentions create-coarse/deep refinement |
| `TestWarnIfTraceDirWithoutConfig` | `internal/config/enforce_test.go` | R6 | `.trace/` without config → stderr nudge substring |
| `Overview.test.tsx` or manual checklist | `web/src/screens/Overview.test.tsx` | R8 | Advisory/plan-gap copy visible when fixture has advisory |
| Doc review + Block 0 regression | — | R11 | Workflow documented; greenfield MCP test still passes |

### Phase 36 regression subset (must stay green — re-run in S02)

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... ./internal/config/... ./internal/httpapi/... -count=1 \
  -run 'Greenfield_MCPPlanBootstrap|FeetSellerExport_GateHonesty|ActiveWork_PlanMissing|TerminalPlanGapAdvisory|PlanBootstrap_Idempotent|GoalStructureWarning_OverThreshold|RegisteredToolNames_IncludesTracePlan|MCPLoopGate|LoopStatus_.*Advisory|HTTPPlanBootstrap|PlanHelp_MentionsRefinement|WarnIfTraceDirWithoutConfig'
```

Individual tests (all must pass):

- `TestGreenfield_MCPPlanBootstrap_EditGatePasses`
- `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`
- `TestActiveWork_PlanMissingStillBlocksEdit`
- `TestEvaluateGate_Done_TerminalPlanGapAdvisory`
- `TestPlanBootstrap_Idempotent`
- `TestGoalStructureWarning_OverThresholdNoPlan`
- `TestRegisteredToolNames_IncludesTracePlan`

### Primary test commands

```bash
# Wave A–B core
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... -count=1

# HTTP + CLI
go test ./internal/httpapi/... ./cmd/trace/... -count=1

# Config + web (if Vitest configured)
go test ./internal/config/... -count=1
cd web && npm test -- --run Overview 2>/dev/null || true

# Full S02 sweep
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./internal/httpapi/... ./cmd/trace/... ./internal/config/... -count=1
```

## Preflight / Plan

Before coding, in Plan mode confirm:

1. Wave A completes before R8 (GUI needs `advisories[]` contract)
2. No ambiguity on R1 vs R5 orthogonality (both may appear)
3. HTTP handler delegates to same bootstrap service as CLI — no duplicate heuristic
4. OpenAPI diff is included in PR scope for R2

## Todo updates

Per board-rights: set own row `done` + Notes listing files touched + closed residual IDs. Do **not** rewrite future prompts.

## Exit criteria

- [ ] All 8 accept residuals (R1–R6, R8, R11) implemented per PLAN §2
- [ ] Wave order A→D respected in commit history or Notes
- [ ] All S02 acceptance tests pass (§5 above)
- [ ] Phase 36 regression subset green
- [ ] R1 guard: no `PlanExists` flip from advisory bridge; `policy.go` untouched
- [ ] R7 preserved: `LoadEnforceMode` default unchanged
- [ ] OpenAPI documents `POST /v1/plans/bootstrap`; MCP `trace_loop` schema includes `gate`
- [ ] Law 19: no business-logic fork in HTTP handlers or `web/`
- [ ] Board Notes: files + residual IDs + test evidence
- [ ] Next: **P37-S02-02**

## Minimal todos

- [ ] **T1** — Wave A: `advisory.go` bootstrap helper + tests
- [ ] **T2** — Wave A: `apply.go` `advisories[]` + R1/R5/R1-guard tests
- [ ] **T3** — Wave B: `tools_loop.go` gate action + `TestMCPLoopGate_MatchesCLI`
- [ ] **T4** — Wave B: HTTP POST bootstrap + OpenAPI + `TestHTTPPlanBootstrap_CreatesPlannerRows`
- [ ] **T5** — Wave B: R4 help refinement + test
- [ ] **T6** — Wave C: R6 enforce test + R11 doc cross-ref
- [ ] **T7** — Wave D: Overview advisories surface + test/checklist
- [ ] **T8** — P36 regression subset + full S02 test sweep
- [ ] **T9** — Graph export if entities changed; board row `done` with evidence

## Next

`P37-S02-02`
