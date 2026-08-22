# Phase 36 / Scope 00 — Investigation (feet-seller `plan_missing`)

**Author:** P36-S00-01 (2026-08-22)  
**Fixture (read-only):** `/home/ali/Desktop/feet seller telegram app`  
**Trace binary:** `go build -o /tmp/trace ./cmd/trace` (repo `/home/ali/Desktop/Trace`)  
**Goal ID:** `353b12a4-57dd-4d68-8379-b2024e064733`  
**Step 1 task:** `33247e2d-aa10-4b25-b194-4b7afb5a6359`  
**Loop 112 task:** `99d8fb92-65ac-462c-82c4-21bcf198c09e`

---

## Repro

All commands run 2026-08-22 against the feet-seller fixture (DB read-only; no writes).

### Task inventory

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" tasks | jq 'length'
# → 123

/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" tasks | jq '[.[] | .work_state] | unique'
# → ["DONE"]

/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" tasks | jq '[.[] | .goal_id] | unique'
# → ["353b12a4-57dd-4d68-8379-b2024e064733"]
```

**123 tasks, all `work_state: DONE`, single `goal_id`.**

### Loop gate — Step 1 (`--for done`)

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" loop gate \
  --task 33247e2d-aa10-4b25-b194-4b7afb5a6359 --for done
```

```json
{"schema_version":"trace.loop.gate.v1","task_id":"33247e2d-aa10-4b25-b194-4b7afb5a6359","for":"done","allowed":false,"recommended_phase":"PLAN","reason_code":"plan_missing","violations":[{"code":"premature_implementation","for":"done","message":"done blocked: recommended phase PLAN (plan_missing)","recommended_phase":"PLAN","reason_code":"plan_missing"}]}
```

### Loop gate — Loop 112 (`--for done`)

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" loop gate \
  --task 99d8fb92-65ac-462c-82c4-21bcf198c09e --for done
```

```json
{"schema_version":"trace.loop.gate.v1","task_id":"99d8fb92-65ac-462c-82c4-21bcf198c09e","for":"done","allowed":false,"recommended_phase":"PLAN","reason_code":"plan_missing","violations":[{"code":"premature_implementation","for":"done","message":"done blocked: recommended phase PLAN (plan_missing)","recommended_phase":"PLAN","reason_code":"plan_missing"}]}
```

**Identical JSON** on Step 1 and Loop 112 (only `task_id` differs).

### Loop gate — edit (Step 1, supplemental)

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" loop gate \
  --task 33247e2d-aa10-4b25-b194-4b7afb5a6359 --for edit
```

```json
{"schema_version":"trace.loop.gate.v1","task_id":"33247e2d-aa10-4b25-b194-4b7afb5a6359","for":"edit","allowed":false,"recommended_phase":"PLAN","reason_code":"plan_missing","violations":[{"code":"premature_implementation","for":"edit","message":"edit blocked: recommended phase PLAN (plan_missing)","recommended_phase":"PLAN","reason_code":"plan_missing"}]}
```

### Plan show

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" plan show \
  --goal 353b12a4-57dd-4d68-8379-b2024e064733
```

Progressive planner header (stderr: `plan: showing plan for goal …`; stdout JSON excerpt):

```json
{
  "goal_id": "353b12a4-57dd-4d68-8379-b2024e064733",
  "current_scope_id": null,
  "phases": [],
  "current_deep_plan": null,
  "lookahead_scope_id": "",
  "lookahead_summary": "",
  "tasks": [ "... 123 task objects, all work_state DONE ..." ]
}
```

Full stdout lists **123 tasks** under one goal; planner fields empty.

### Seed export counts

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" seed export -o /tmp/feet-export.json
```

| Top-level array | Count |
|-----------------|------:|
| `plan_changes` | **11** |
| `plan_phases` | **0** |
| `plan_scopes` | **0** |
| `scope_deep_plans` | **0** |
| `goal_plan_state` | **0** (field absent → treated as 0) |
| `tasks` | 123 |
| `goals` | 1 |
| `discoveries` | 45 |
| `decisions` | 24 |
| `reviews` | omitted (default export per `internal/domain/seed_export.go:349`) |

Sample `plan_changes` titles (11 total): *Architecture plan: TMA adult marketplace Stars-first*, *Production-ready MVP plan*, *Grilling closed*, *Process pivot: local tests default*, *Loop 97: optional naming rename becomes active*, …

### Why identical gate on terminal tasks is goal-level inheritance

Gate evaluation loads the task, resolves **`goal_id`**, then builds policy inputs from **goal planner state**, not from task `work_state`:

1. `loadGateContext` reads `task.GoalID` and passes it to `BuildPolicyInputs` (`internal/loop/gate.go:135–153`).
2. `PlanExists` is set only when `plan.GetPlan(goalID)` has **both** non-empty `current_scope_id` **and** non-null `current_deep_plan` (`internal/loop/policy.go:45–49`).
3. With progressive planner empty for the sole goal, `PlanExists` is **false for every task** on that goal.
4. `SelectNext` returns `PhasePlan` / `ReasonPlanMissing` when `!inputs.PlanExists` (`internal/deliberation/select.go:28–29`).
5. `evaluateDone` has **no** short-circuit for tasks already `DONE` — it runs the same deliberation path (`internal/loop/gate.go:227–265`).

So Step 1, Loop 112, and any other task on goal `353b12a4-…` inherit the same **`plan_missing`** outcome. Task terminal state does not change gate semantics today.

---

## Verdict table

| Layer | Blame | Evidence |
|-------|-------|----------|
| **Trace product** | **Primary** | Two disconnected planning signals: 11 `plan-change` entities vs 0 progressive planner rows (`/tmp/feet-export.json` counts). `PlanExists` reads only progressive planner (`internal/loop/policy.go:45–49`); `PlanCritiqued` can be satisfied by plan-changes on apply path (`:60–62`) but that does **not** set `PlanExists`. MCP exposes `trace_add kind=plan-change` (`internal/mcp/tools_write.go:130–136`) but **no** plan CLI/MCP (`cmd/trace/plan.go:67`). Agents on MCP cannot populate what the gate checks. |
| **Agent misuse** | **Secondary** | Agents could have run CLI `trace plan create-coarse|set-current|deep` — documented only on CLI, not MCP/install. Feet-seller shows rich planning via plan-changes (11 titles) consistent with agent *intent* to plan, but wrong store for `PlanExists`. Not “no planning”; wrong **representation**. |
| **Harness** | **Secondary** | No `trace install` artifacts on fixture (see §5). `.trace/config.json` absent → `LoadEnforceMode` → `off` (`internal/config/enforce.go:21–27`). Transition to DONE does not run gate unless `--enforce` (`cmd/trace/transition.go:36–69`). MCP `trace_transition` calls `TransitionTask` directly with no gate (`internal/mcp/tools_write.go:224–231`). **127** reviews all `PASS` via `trace review list` — work completed via review + `as_operator` path without gate ever blocking. |
| **GUI** | **Tertiary** | `TaskDetail.tsx:76–86` always fetches done gate on mount; `GateStrip.tsx:41–56` shows red “Gate blocked” + `plan_missing` even when `work_state` is already DONE. Misleading UX on finished work; **not** root cause of empty progressive planner. |

**Answer #1:** Trace product design is the primary fault (dual planning model + MCP gap). Agent and harness share secondary blame (undiscoverable bootstrap, enforce/install off). GUI amplifies confusion but did not cause 123 DONE without `PlanExists`.

---

## Two planning systems

| System | Storage | Agent surface | Gate signal |
|--------|---------|---------------|-------------|
| **Causal graph** | `plan_changes` entities (+ links) | `trace add` / MCP `trace_add kind=plan-change` | **`PlanCritiqued`** only on apply path when writes include plan-changes (`internal/loop/policy.go:60–62`) — **not** `PlanExists` |
| **Progressive planner** | `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` | CLI `trace plan create-coarse\|set-current\|deep\|show` only | **`PlanExists`** requires `current_scope_id` + `current_deep_plan` (`internal/loop/policy.go:45–49`) |

**Live proof — 11 vs 0:** export counts above; `plan show` → `phases: []`, `current_scope_id: null`, `current_deep_plan: null`.

**Why rich plan-changes still yield `plan_missing`:** deliberation selects PLAN when `!PlanExists` before critique/execute (`internal/deliberation/select.go:28–29`). Feet-seller never ran `create-coarse` (0 planner rows). Agents recording MVP/PIR pivots as plan-changes satisfied **critique** semantics in isolation but never created progressive planner state the gate requires.

**Answer #2:** Confirmed — two systems, one gate input. MCP/agent path naturally fills (1); gate checks (2).

---

## Transition audit

How **123 tasks reached DONE** while `PlanExists` was never true:

### 1. Enforce off (default)

Fixture filesystem audit:

| Path | Status |
|------|--------|
| `AGENTS.md` | **Absent** |
| `.trace/config.json` | **Absent** |
| `.cursor/rules/trace-enforcement.mdc` | **Absent** |
| `.cursor/hooks/trace-loop-gate.sh` | **Absent** |

Missing `.trace/config.json` → `LoadEnforceMode` returns `EnforceOff` (`internal/config/enforce.go:21–27`). No hook pre-edit gate, no strict deny.

### 2. Transition gate opt-in

`trace transition --enforce` runs `GateForDone` only when flag set (`cmd/trace/transition.go:36–69`). Default DONE path: linked Review PASS + `--as-operator` (or `--allow-done` hatch) — **not** progressive plan.

### 3. MCP path

`trace_transition` mirrors CLI transition options but **never** calls `loop.EvaluateGate` (`internal/mcp/tools_write.go:207–231`). Agents using MCP could transition to DONE without gate.

### 4. Reviews

```bash
/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" review list | jq 'length'
# → 127

/tmp/trace -C "/home/ali/Desktop/feet seller telegram app" review list | jq '[.[] | select(.result == "PASS")] | length'
# → 127
```

~127 PASS reviews (INTAKE ~120) — consistent with loop implement+review workflow completing without progressive planner bootstrap.

### 5. Loop status (Loop 112)

`trace loop status --task 99d8fb92-…` shows `policy_inputs.plan_exists: false`, `why_selected: "plan_missing"`, edit violation present — gate **would** block edit/done **if invoked**, but workflow did not require it.

**Answer #3:** Gate logic is **correct per code** (`plan_missing` when progressive planner empty) but was **never blocking** during agent workflow (enforce off, no `--enforce`, no MCP gate, no install hooks). Transitions succeeded on domain rules (reviews + operator claim).

---

## MCP + install audit

### MCP tool inventory (`internal/mcp/server.go:52–214`, `RegisteredToolNames` `:216–224`)

Fifteen tools: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`, `trace_search`, `trace_changes`, `trace_regressions`, `trace_loop`, `trace_agents`.

| Check | Result |
|-------|--------|
| `trace_add` supports `plan-change` | **Yes** — `tools_write.go:130–136` |
| `trace_plan` / create-coarse / set-current / deep / show | **No** — `cmd/trace/plan.go:67`: *"No MCP plan tools."* |
| `trace_loop` actions | **`next\|apply\|status` only** — `tools_loop.go:39–47`; **no `gate`** |
| `trace_transition` gate | **No** — direct `TransitionTask` (`tools_write.go:224–231`) |

**Answer #4:** MCP cannot bootstrap progressive planner or run loop gate; primary agent surface is incomplete for PLAN phase satisfaction.

### Install audit

`internal/install/enforcement.go` documents:

- Pre-edit: `trace loop gate --task "$TRACE_TASK_ID" --for edit` (`:57`, `:73`, `:92`)
- Pre-DONE: `trace loop status` + resolve violations (`:58`, `:74`)
- Opt-in `--enforce` on transition/export (`:59`, `:75`)

**Grep `create-coarse` under `internal/install/`:** **no matches** — install contract does not document how to bootstrap `trace plan create-coarse`.

Feet-seller: no install artifacts (table above) → agents had no AGENTS.md enforcement block, no Cursor hook, enforce default off.

**Answer #5:** Harness never wired; install docs omit mandatory coarse-plan bootstrap even when install is run elsewhere.

---

## Mega-goal pattern

**123 tasks / 1 goal** — `plan show` lists all 123 under `353b12a4-…` with empty `phases`.

| Aspect | Observation |
|--------|-------------|
| Structure | Single mega-goal spans research steps, grilling, PIR loops 01–112 |
| Gate amplification | Goal-level `PlanExists` → identical `plan_missing` on **every** task detail view |
| Product guidance | No warning when task count ≫ N without coarse plan (candidate S01 fix) |

**Answer #6:** Mega-goal is a **product guidance gap**, not the sole root cause (empty progressive planner would block even a smaller goal), but it **amplifies** identical red gate on all 123 DONE tasks in GUI.

---

## Fundamental fix options

Ranked for **S01** (`DESIGN-LOCKS.md`); Law 19 — library policy first, adapters follow.

| Rank | Option | Rationale |
|------|--------|-----------|
| 1 | **MCP plan tools** (`trace_plan`: create-coarse, set-current, deep, show) | Closes primary agent surface gap; agents can satisfy `PlanExists` without undocumented CLI-only steps. |
| 2 | **Install contract** | First goal → mandatory create-coarse in AGENTS.md + cursor rules; document bootstrap in `internal/install/enforcement.go`. |
| 3 | **Bootstrap command** (`trace plan bootstrap --goal`) | Dogfood recovery for feet-seller-like repos (11 plan-changes, 0 planner). |
| 4 | **Terminal gate honesty** | Library: DONE/SKIPPED tasks — don’t emit actionable “done blocked” for goal-level `plan_missing` (`evaluateDone` / gate display). |
| 5 | **Goal structure warning** | Warn when >N tasks under goal without coarse plan. |
| 6 | **Enforce nudge** | Default `warn` when `.trace/` exists without config. |
| 7 | **PlanExists bridge** | Heuristic from plan-change density — **careful**; must not fake progressive plan. |

**Recommended S01 direction:** (1) + (2) + (3) for agent-complete bootstrap; (4) for honesty on terminal tasks; defer bridge (7) unless S01 rejects pure MCP path.

---

## GUI secondary

Adapter-only (Law 19); defers product fix to S02.

| Location | Behavior |
|----------|----------|
| `web/src/screens/TaskDetail.tsx:76–86` | `loadGate()` always calls `getLoopGate(taskId, 'done')` on mount — **including DONE tasks** |
| `web/src/screens/TaskDetail.tsx:198–210` | “DONE gate” panel + warn banner when gate blocked (advisory; API still decides transition) |
| `web/src/components/GateStrip.tsx:41–56` | `banner--error`, “Gate blocked”, displays `reason_code` (`plan_missing`) |

**Answer #7:** GUI shows identical red alarm on all finished tasks because it faithfully mirrors library gate JSON with no terminal-task nuance. Secondary UX issue once root planning model is fixed.

---

## Red-capable tests (S02 sketches)

Reference patterns: `cmd/trace/loop_test.go:1137–1167` (`createCurrentDeepPlanForLoopTest`), `internal/loop/gate_test.go:198–205` (`TestEvaluateGate_Edit_PlanMissing`).

### 1. `TestGreenfield_MCPPlanBootstrap_EditGatePasses`

**Setup:** Temp dir, `trace init`, add goal + task via MCP/CLI.  
**Act:** MCP `trace_plan create-coarse` + `set-current` + `deep` (or CLI equivalent until MCP lands).  
**Assert:** `trace loop gate --for edit` → `allowed: true`, `reason_code` not `plan_missing`; `policy_inputs.plan_exists: true` in status JSON.  
**Shape:** Mirror `createCurrentDeepPlanForLoopTest` + `TestEvaluateGate_Edit_ExecuteReady`.

### 2. `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`

**Setup:** Import/copy feet-seller export (`plan_changes: 11`, planner arrays 0) into temp store.  
**Assert (pre-fix):** `loop gate --for done` on a DONE task → `plan_missing` (documents current behavior).  
**Assert (post S02 terminal honesty):** DONE task gate returns honest terminal signal (e.g. allowed or distinct reason — not identical actionable “done blocked” alarm); after `plan bootstrap`, edit gate passes.  
**Shape:** Extend `gate_test.go` eval patterns; optional fixture JSON from `/tmp/feet-export.json`.

---

## Out of scope rejects

Per `SCOPE-TODOS.md` — documented rejects:

1. **Blaming Phase 35 pick logic** — Phase 35 fixed active-task selection (Loop 112 landing). Symptom here is gate **meaning** on terminal DONE tasks + goal-level `plan_missing`, not pick order.
2. **“Just create a plan” as the only fix** — Insufficient: does not explain misleading DONE gate on **finished** work, MCP gap, or 123 identical alarms; needs bootstrap + honesty + agent surface.
3. **Weakening `plan_missing` globally for active PLAN work** — Rejected; DESIGN-LOCKS preserves PLAN enforcement for active non-terminal work.
4. **Deleting feet-seller history** — Rejected; recovery/bootstrap path preferred over erasing 123 DONE tasks, 11 plan-changes, 127 reviews.

---

## Summary

| Item | Finding |
|------|---------|
| Repro | 123 DONE / 1 goal; both gates `plan_missing` identical; planner empty; export **11 / 0** |
| Root cause | Progressive planner never initialized; gate checks planner only; MCP/install cannot bootstrap |
| Transitions | 127 PASS reviews, enforce off, no `--enforce`, MCP no gate → DONE without `PlanExists` |
| Verdict | Trace **primary**; agent + harness **secondary**; GUI **tertiary** |
| S01 recommendation | MCP plan tools + install bootstrap docs + recovery command + terminal gate honesty |

**Next row:** P36-S00-02 (independent review vs DESIGN-LOCKS + INTAKE + live cites).
