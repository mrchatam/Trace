# P36-S02-01 — Implement

## Metadata
- id: P36-S02-01
- todo_ids: [P36-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, api-and-interface-design]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Implement the locked fix set from [PLAN.md](../scope-01-plan/PLAN.md): progressive planner bootstrap (library + CLI + MCP), install contract, terminal gate honesty, goal-structure advisory, enforce nudge, GUI adapter — **without** weakening active-work PLAN enforcement or adding the deferred PlanExists bridge.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [PLAN.md](../scope-01-plan/PLAN.md) — **SoT** for touch-list, order, acceptance tests
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md)
- [00-PLANNER.md](00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22):
  - `internal/loop/policy.go:45–49` — `PlanExists` (unchanged)
  - `internal/loop/gate.go:227–265` — `evaluateDone` (no terminal short-circuit today)
  - `internal/deliberation/select.go:28–29` — `!PlanExists` → PLAN / `plan_missing`
  - `cmd/trace/plan.go:46–67` — CLI subcommands; *"No MCP plan tools."*
  - `internal/mcp/server.go:216–224` — 15 registered tools (no `trace_plan`)
  - `internal/mcp/tools_loop.go:39–47` — `trace_loop` action param pattern
  - `internal/install/enforcement.go:51–98` — AGENTS/cursor (no bootstrap step)
  - `internal/config/enforce.go:21–27` — missing config → `EnforceOff`
  - `web/src/components/GateStrip.tsx:33–56` — already uses warn when `allowed && violations`
  - Test helpers: `createCurrentDeepPlanForLoopTest` (`cmd/trace/loop_test.go:1137–1167`), `TestEvaluateGate_Edit_PlanMissing` (`internal/loop/gate_test.go:198–205`)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| SoT | [PLAN.md](../scope-01-plan/PLAN.md) §2 fix set, §5 touch-list, §6 acceptance tests, §7 non-goals |
| Law 19 | Policy in library; MCP/CLI/HTTP/GUI are thin adapters over `internal/planner` + `internal/loop` |
| Implementation order | **Strict** — see §Implementation order below |
| MCP tools | **16** after S02 (`trace_plan` added); update `RegisteredToolNames()` + `TestToolNamesRegistered` |
| `trace_plan` shape | Single tool, `action` param — mirror `trace_loop` (`tools_loop.go:39–47`) |
| Bootstrap | `trace plan bootstrap --goal <id>` + MCP `trace_plan action=bootstrap`; idempotent when plan exists |
| Terminal gate | `goal_plan_gap_terminal_advisory` — `allowed: true` for DONE/SKIPPED when only plan gap remains |
| Goal warning threshold | **N = 15** tasks under goal without `PlanExists` |
| Enforce default | **Unchanged** — `EnforceOff` when config missing; stderr nudge only |
| PlanExists bridge | **DEFER** — no `policy.go` heuristic |
| HTTP POST plan routes | **DEFER** — `GET /v1/plans` only |
| Feet-seller fixture | **Read-only** in S02 tests; live bootstrap in S03 only |
| Dogfood path | `/home/ali/Desktop/feet seller telegram app` — do not mutate in S02 |

## Fix set (accept / defer — do not re-debate)

| § | Fix | Action |
|---|-----|--------|
| 2.1 | MCP `trace_plan` (create-coarse, set-current, deep, show, bootstrap) | **Accept** |
| 2.2 | `trace plan bootstrap --goal` | **Accept** |
| 2.3 | Install contract (mandatory coarse-plan bootstrap step) | **Accept** |
| 2.4 | PlanExists bridge from plan-changes | **Defer** — no code |
| 2.5 | Terminal gate honesty (`goal_plan_gap_terminal_advisory`) | **Accept** |
| 2.6 | Enforce document nudge (`.trace/` without config) | **Accept** |
| 2.7 | Goal structure warning (N=15) | **Accept** |
| 2.8 | Feet-seller recovery | **S02:** ship tool + temp import test; **S03:** live verify |

## Touch-list (ordered — implement in this sequence)

| Step | File | Action | PLAN § |
|------|------|--------|--------|
| 1 | `internal/planner/bootstrap.go` | **Create** | 2.2 |
| 2 | `internal/planner/bootstrap_test.go` | **Create** | 6.2 |
| 3 | `internal/planner/advisory.go` | **Create** | 2.7 |
| 4 | `internal/planner/advisory_test.go` | **Create** | 2.7 |
| 5 | `internal/planner/service.go` | **Edit** (if bootstrap/advisory need hooks) | 2.2, 2.7 |
| 6 | `internal/loop/gate.go` | **Edit** — terminal honesty in `evaluateDone` (+ edit path if needed) | 2.5 |
| 7 | `internal/loop/gate_test.go` | **Edit** | 6.2, 6.3 |
| 8 | `internal/loop/testdata/feet-export-min.json` | **Create** — anonymized subset (11 plan-changes, 0 planner) | 6.2 |
| 9 | `internal/mcp/tools_plan.go` | **Create** — thin adapter over planner service | 2.1 |
| 10 | `internal/mcp/server.go` | **Edit** — register `trace_plan`; `RegisteredToolNames` → 16 | 2.1 |
| 11 | `internal/mcp/mcp_test.go` | **Edit** — 16-tool lock + greenfield workflow | 6.1 |
| 12 | `cmd/trace/plan.go` | **Edit** — `bootstrap` subcommand; show stderr advisory | 2.2, 2.7 |
| 13 | `internal/install/enforcement.go` | **Edit** — bootstrap step in AGENTS/cursor blocks | 2.3 |
| 14 | `cmd/trace/install.go` | **Edit** — post-install hint when goals exist, planner empty | 2.3 |
| 15 | `internal/config/enforce.go` | **Edit** — `WarnIfTraceDirWithoutConfig` helper | 2.6 |
| 16 | `cmd/trace/init.go` | **Edit** — enforce nudge on init completion | 2.6 |
| 17 | `web/src/components/GateStrip.tsx` | **Edit** (if needed) — advisory copy for terminal plan gap | 2.5 |
| 18 | `web/src/screens/TaskDetail.tsx` | **Edit** (optional) — advisory tone on DONE gate panel | 2.5 |

**Explicit non-touch:**

- `internal/loop/policy.go` — no bridge; `PlanExists` definition unchanged
- `internal/deliberation/select.go` — cite only; deliberation order unchanged
- `internal/httpapi/server.go` — no POST plan mutation routes in S02
- Global weaken of `plan_missing` for active non-terminal work
- GUI-only suppression without library terminal semantics
- Live mutation of feet-seller dogfood fixture

## Implementation order (strict)

```text
1. internal/planner/bootstrap.go + bootstrap_test.go
2. internal/planner/advisory.go (+ advisory_test.go)
3. internal/loop/gate.go terminal honesty + gate_test.go (+ feet-export-min.json fixture)
4. internal/mcp/tools_plan.go + server.go registration + mcp_test.go
5. cmd/trace/plan.go bootstrap subcommand + plan show advisory
6. internal/install/enforcement.go + cmd/trace/install.go hint
7. internal/config/enforce.go + cmd/trace/init.go nudge
8. web GateStrip/TaskDetail advisory mapping (only if library JSON stable)
9. Run full acceptance test suite (§6 below)
```

## Role work

### 1. Bootstrap (`internal/planner/bootstrap.go`)

- Input: goal-linked `plan_changes` (titles + bodies via store/domain)
- Heuristic: pick primary plan-change (most recent non-superseded; prefer title length/recency)
- Output: one phase + one scope + minimal `deep` — valid `PlanExists` state
- Idempotent: if goal already has `current_scope_id` + `current_deep_plan`, no-op success + stderr note
- **No LLM generation** — caller-supplied / derived from existing plan-change text only
- Wire: `cmd/trace/plan.go` subcommand `bootstrap --goal <id>`; MCP `trace_plan action=bootstrap`

### 2. Goal structure warning (`internal/planner/advisory.go`)

- Condition: `task_count(goal) > 15 && !PlanExists(goal)`
- Surfaces: `trace plan show` stderr; MCP `trace_plan action=show` field `goal_structure_warning`; optional `trace loop status` `advisories[]`
- **Advisory only** — does not set `PlanExists` or weaken edit block

### 3. Terminal gate honesty (`internal/loop/gate.go`)

- When `task.WorkState` is `DONE` or `SKIPPED` and deliberation would return `plan_missing`:
  - `--for done` and `--for edit`: `allowed: true` with advisory violation
  - `reason_code: goal_plan_gap_terminal_advisory`
  - Message includes `goal_id` + hint: `trace plan bootstrap --goal` or MCP `trace_plan`
- **Preserve** verification/regression/uncertainty blocks at `gate.go:243–257` **before** terminal advisory path
- Active non-terminal: unchanged — `plan_missing`, `allowed: false`

### 4. MCP `trace_plan` (`internal/mcp/tools_plan.go`)

- Actions: `create-coarse`, `set-current`, `deep`, `show`, `bootstrap`
- Thin adapter — same planner service calls as `cmd/trace/plan.go`
- Error contract: match CLI stderr + stdout JSON (snake_case on `show`)
- Register in `server.go`; update `RegisteredToolNames()` to **16 tools**

### 5. Install contract (`internal/install/enforcement.go`)

- Add numbered step **between** TRACE_TASK_ID setup and pre-edit gate:
  - *"New goal without coarse plan: bootstrap via `trace plan create-coarse` or MCP `trace_plan` before edit gate can pass."*
- Update: `AgentsEnforcementBlock()`, `cursorRulesMDCContent()`, `EnforcementRulesMarkdown()`
- `cmd/trace/install.go`: success stderr mentions bootstrap when store has goals but empty planner

### 6. Enforce nudge (`internal/config/enforce.go`, `cmd/trace/init.go`)

- One-time stderr when `.trace/` exists but config missing/invalid
- Suggest: *"Consider `.trace/config.json` with `\"enforce\": \"warn\"` after `trace install`"*
- **No behavior change** to gate evaluation or default enforce mode

### 7. GUI adapter (Law 19)

- `GateStrip.tsx` already maps `allowed && violations` → warn banner (`:41–56`)
- Verify terminal advisory renders as **warn**, not error block
- `TaskDetail.tsx`: optional copy tweak for DONE gate panel — no independent gate logic

## Test strategy (required)

### Primary acceptance tests (PLAN §6 — must pass)

| Name | File | Assert shape |
|------|------|--------------|
| `TestGreenfield_MCPPlanBootstrap_EditGatePasses` | `internal/mcp/mcp_test.go` (preferred) or `cmd/trace/loop_test.go` | Temp dir → init → goal+task → MCP `trace_plan` create-coarse + set-current + deep → edit gate `allowed: true`, `reason_code` ≠ `plan_missing`; status `policy_inputs.plan_exists: true` |
| `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap` | `internal/loop/gate_test.go` + helper | Import `feet-export-min.json` (11 plan-changes, 0 planner); DONE task `--for done` → `allowed: true`, advisory `goal_plan_gap_terminal_advisory`; post-`plan bootstrap` → edit gate `allowed: true`, `plan_exists: true` |
| `TestActiveWork_PlanMissingStillBlocksEdit` | `internal/loop/gate_test.go` | Non-terminal task, no planner → edit `allowed: false`, `reason_code: plan_missing`, `recommended_phase: PLAN` |

### Recommended additional tests

| Name | File | Covers |
|------|------|--------|
| `TestEvaluateGate_Done_TerminalPlanGapAdvisory` | `gate_test.go` | §2.5 terminal DONE |
| `TestPlanBootstrap_Idempotent` | `bootstrap_test.go` | §2.2 idempotency |
| `TestGoalStructureWarning_OverThresholdNoPlan` | `advisory_test.go` | §2.7 N=15 |
| `TestRegisteredToolNames_IncludesTracePlan` | `mcp_test.go` | §2.1 — 16 tools, `trace_plan` in locked list |

### Greenfield fixture pattern

```text
t.TempDir()
→ trace init (CLI or test helper)
→ add goal + task (MCP trace_add or domain helper)
→ trace_plan create-coarse + set-current + deep (MCP)
→ loop gate --for edit / EvaluateGate
```

Mirror `createCurrentDeepPlanForLoopTest` (`loop_test.go:1137–1167`) but via MCP path for §6.1.

### Feet-seller export fixture

- Create `internal/loop/testdata/feet-export-min.json` — anonymized subset from INVESTIGATION export counts (11 plan-changes, 0 planner rows, ≥1 DONE task)
- **Read-only import** into temp store in tests — do not write to dogfood fixture
- Live bootstrap on `/home/ali/Desktop/feet seller telegram app` is **S03** scope

### Test commands

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... -count=1
```

## Preflight / Plan

Before coding, in Plan mode confirm:

1. All touch-list files exist or are clearly new-create paths
2. No ambiguity on terminal advisory vs active-work block semantics
3. Bootstrap heuristic boundaries (minimal recovery, not full 123-task reconstruction)

## Todo updates

Per board-rights: set own row `done` + Notes listing files touched. Do **not** rewrite future prompts.

## Exit criteria

- [ ] All three primary acceptance tests pass (§6.1–§6.3)
- [ ] Recommended tests pass or are explicitly deferred in Notes with reason
- [ ] MCP registers **16** tools including `trace_plan`
- [ ] Bootstrap CLI + MCP paths both satisfy `PlanExists` on temp import
- [ ] Terminal DONE/SKIPPED tasks emit `goal_plan_gap_terminal_advisory` (not actionable done-blocked)
- [ ] Active non-terminal `plan_missing` edit block preserved
- [ ] Install contract documents bootstrap step; no PlanExists bridge in `policy.go`
- [ ] `go test` green for touched packages
- [ ] Board Notes: files touched + test evidence
- [ ] Next: **P36-S02-02**

## Minimal todos

- [ ] **T1** — `bootstrap.go` + `bootstrap_test.go` (idempotent recovery from plan-changes)
- [ ] **T2** — `advisory.go` + test (N=15 goal warning)
- [ ] **T3** — `gate.go` terminal honesty + `gate_test.go` (terminal advisory + active-work preservation)
- [ ] **T4** — `feet-export-min.json` fixture + `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`
- [ ] **T5** — `tools_plan.go` + `server.go` registration (16 tools)
- [ ] **T6** — `mcp_test.go` — `TestGreenfield_MCPPlanBootstrap_EditGatePasses` + tool-name lock
- [ ] **T7** — `cmd/trace/plan.go` bootstrap subcommand + show advisory
- [ ] **T8** — `internal/install/enforcement.go` + `cmd/trace/install.go` bootstrap contract/hint
- [ ] **T9** — `internal/config/enforce.go` + `cmd/trace/init.go` enforce nudge
- [ ] **T10** — GUI adapter verify (GateStrip warn on advisory; TaskDetail optional)
- [ ] **T11** — Full test run; board row `done` with evidence

## Next

`P36-S02-02`
