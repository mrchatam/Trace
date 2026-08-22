# P36-S00-01 — Investigate

## Metadata
- id: P36-S00-01
- todo_ids: [P36-S00-01]
- role: implementer
- skills: [diagnosing-bugs, systematic-debugging, planning-and-task-breakdown, test-driven-development]
- mcps: [user-trace (read-only audit), Shell via trace CLI]
- agents: []
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-00-investigate/INVESTIGATION.md`: live feet-seller repro, **Trace vs agent vs harness verdict**, two planning systems analysis, MCP/install/enforce audit, transition-without-gate path, mega-goal pattern, GUI secondary notes, red-capable test sketches for S02. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [00-PLANNER.md](00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live policy/gate: `internal/loop/policy.go`, `internal/loop/gate.go`, `internal/deliberation/select.go`
- Live MCP: `internal/mcp/server.go`, `internal/mcp/tools_write.go`, `internal/mcp/tools_loop.go`
- Live plan CLI: `cmd/trace/plan.go`
- Live transition/enforce: `cmd/trace/transition.go`, `internal/config/enforce.go`
- Live install: `cmd/trace/install.go`, `internal/install/enforcement.go`
- Live GUI adapters: `web/src/screens/TaskDetail.tsx`, `web/src/components/GateStrip.tsx`
- Fixture: `/home/ali/Desktop/feet seller telegram app`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). No material ambiguity — locked defaults below are authoritative.

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-36-gate-honesty-terminal-tasks/scopes/scope-00-investigate/INVESTIGATION.md` **only** |
| Product edits | **Forbidden** (no Go/CLI/web/SQLite writes in Trace repo or fixture) |
| Fixture | `/home/ali/Desktop/feet seller telegram app` — **read-only** (do not mutate `.trace/trace.db`) |
| Goal ID | `353b12a4-57dd-4d68-8379-b2024e064733` |
| Step 1 task | `33247e2d-aa10-4b25-b194-4b7afb5a6359` |
| Loop 112 task | `99d8fb92-65ac-462c-82c4-21bcf198c09e` |
| Trace binary | Build from repo: `go build -o /tmp/trace ./cmd/trace` (or `trace` on PATH) |
| Primary feedback loop | CLI `trace loop gate` JSON asserts |
| Secondary feedback loop | GUI TaskDetail DONE gate strip (screenshot optional) |
| Export proof path | `trace -C "<fixture>" seed export -o /tmp/feet-export.json` then count top-level arrays |
| Sequence | → `P36-S00-02`; S01 waits for S00-02 PASS |

## Preflight / Plan

1. Build `trace` from `/home/ali/Desktop/Trace`.
2. Run locked CLI commands (below) and capture stdout JSON verbatim in INVESTIGATION.md.
3. Read cited library files; quote **live** line numbers (planner verified 2026-08-22 — re-verify if lines drift).
4. Write INVESTIGATION.md sections in template order; update board row **P36-S00-01** status + Notes only.

## Role work

### A — Repro (section 1)

Run from Trace repo (fixture read-only; gate needs DB open — use full shell, not sandbox):

```bash
trace -C "/home/ali/Desktop/feet seller telegram app" tasks | jq 'length'
trace -C "/home/ali/Desktop/feet seller telegram app" loop gate \
  --task 33247e2d-aa10-4b25-b194-4b7afb5a6359 --for done
trace -C "/home/ali/Desktop/feet seller telegram app" loop gate \
  --task 99d8fb92-65ac-462c-82c4-21bcf198c09e --for done
trace -C "/home/ali/Desktop/feet seller telegram app" plan show \
  --goal 353b12a4-57dd-4d68-8379-b2024e064733
trace -C "/home/ali/Desktop/feet seller telegram app" seed export -o /tmp/feet-export.json
```

**Expected (planner spot-check — cite your run):**

| Check | Expected |
|-------|----------|
| tasks | 123, all `work_state: DONE`, single `goal_id` |
| both `--for done` gates | `allowed: false`, `reason_code: plan_missing`, `recommended_phase: PLAN` — **identical** on Step 1 and Loop 112 |
| `plan show` | `phases: []`, `current_scope_id: null`, `current_deep_plan: null` |
| export counts | `plan_changes: 11`; `plan_phases/plan_scopes/scope_deep_plans/goal_plan_state`: **0** |

Document why identical gate on terminal tasks is **goal-level inheritance** (`PlanExists` reads goal planner state, not task `work_state`).

### B — Verdict table (section 2) — **must answer #1**

Three-row table: **Trace product** | **Agent misuse** | **Harness**. For each: share of blame (primary/secondary/tertiary), evidence bullets, file:line or CLI JSON cites. Align with INTAKE verdict sketch but **prove** with live data — do not copy without evidence.

### C — Two planning systems (section 3) — **must answer #2**

| System | Storage | Agent surface | Gate signal |
|--------|---------|---------------|-------------|
| Causal graph | `plan_changes` entities | `trace add` / MCP `trace_add kind=plan-change` | **`PlanCritiqued`** only (`policy.go:60–62`) |
| Progressive planner | `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` | CLI `trace plan create-coarse\|set-current\|deep\|show` | **`PlanExists`** (`policy.go:45–49`) |

Prove **11 vs 0** from export (and optionally `plan show` JSON). Explain why agents recording rich plan-changes still get `plan_missing`: deliberation selects PLAN when `!PlanExists` (`deliberation/select.go:28–29`).

### D — Transition audit (section 4) — **must answer #3**

Explain how **123 tasks reached DONE** while `PlanExists` was never true:

1. **Enforce off:** no `.trace/config.json` → `config.LoadEnforceMode` → `off` (`internal/config/enforce.go`). Document missing file on fixture.
2. **Transition gate opt-in:** `--enforce` on `trace transition --to DONE` is **optional** (`cmd/trace/transition.go:36–69`); default DONE path uses Review PASS + `as_operator`, not `GateForDone`.
3. **MCP path:** `trace_transition` mirrors CLI (`internal/mcp/tools_write.go`); no mandatory gate before DONE unless harness adds it.
4. **Reviews:** ~120 PASS linked — cite export review count or sample.

State explicitly: gate was **correct per code** but **never blocking** during agent workflow.

### E — MCP + install audit (section 5) — **must answer #4 and #5**

**MCP inventory** (`internal/mcp/server.go`): list all registered tools. Confirm:

- `trace_add` supports `plan-change` (`tools_write.go:130–136`)
- **No** `trace_plan` / create-coarse / set-current / deep / show
- `trace_loop` actions = `next|apply|status` only — **not** `gate` (`tools_loop.go:39–47`)
- `cmd/trace/plan.go:67` — "No MCP plan tools."

**Install audit** on fixture (filesystem, read-only):

| Artifact | Expected on feet-seller |
|----------|-------------------------|
| `AGENTS.md` | **Absent** |
| `.trace/config.json` | **Absent** → enforce off |
| `.cursor/rules/trace-enforcement.mdc` | **Absent** |
| `.cursor/hooks/trace-loop-gate.sh` | **Absent** |

Read `internal/install/enforcement.go` — note install docs mention `loop gate --for edit` but **do not** document `trace plan create-coarse` bootstrap (grep confirms no match in `internal/install/`).

### F — Mega-goal pattern (section 6) — **must answer #6**

123 tasks / 1 goal — product guidance gap? Cite `plan show` listing all 123 tasks under one goal. Discuss whether Trace should warn when task count >> N without coarse plan (candidate fix for S01, not implemented).

### G — Fundamental fix options (section 7)

Rank options from DESIGN-LOCKS for S01: MCP plan tools, PlanExists bridge, bootstrap command, install contract, enforce nudge, goal structure warning, terminal gate honesty. Law 19: library policy first; adapters follow.

### H — GUI secondary (section 8) — **must answer #7**

Cite adapter path only:

- `TaskDetail.tsx:76–86` — always fetches `getLoopGate(taskId, 'done')` on mount (even for DONE tasks)
- `TaskDetail.tsx:198–210` — DONE gate panel + warn banner when blocked
- `GateStrip.tsx:41–56` — `banner--error` + "Gate blocked" + `plan_missing` copy

Note: misleading UX on **finished** work; not root cause. Defer product fix to S02.

### I — Red-capable tests for S02 (section 9) — **must answer #8**

Sketch **two** test scenarios S02 can implement (names + assert shape, no code here):

1. **Greenfield bootstrap:** temp dir → goal + MCP/CLI coarse plan → `loop gate --for edit` passes Plan phase → agent path complete.
2. **Legacy dogfood:** import or copy feet-seller export → gate still `plan_missing` until recovery/bootstrap → terminal DONE tasks get honest gate (not identical red alarm) once S02 fix lands.

Reference existing test patterns: `cmd/trace/loop_test.go` (create-coarse setup), `internal/loop/gate_test.go`.

### J — Out of scope rejects (section 10)

Address all four rejects in SCOPE-TODOS.md (Phase 35 pick, GUI-only fix, global weaken, delete history).

## INVESTIGATION.md template (required sections)

Use these as `##` headings in order (content per Role work above):

1. **Repro**
2. **Verdict table**
3. **Two planning systems**
4. **Transition audit**
5. **MCP + install audit**
6. **Mega-goal pattern**
7. **Fundamental fix options** (ranked for S01)
8. **GUI secondary**
9. **Red-capable tests** (S02 sketches)
10. **Out of scope rejects**

Every factual claim needs **file:line** or **CLI JSON** evidence.

## Todo updates

Set **P36-S00-01** row `done` with Notes listing INVESTIGATION.md path + key verdict one-liner. **Do not** edit other board rows.

## Exit criteria

- [ ] `INVESTIGATION.md` complete with file:line + CLI cites; all eight must-answer items covered
- [ ] No product code changes
- [ ] Board Notes on P36-S00-01 only
- [ ] Next: **P36-S00-02**

## Minimal todos

- [ ] Preflight: build trace; run locked CLI commands; capture JSON
- [ ] Export counts: prove 11 plan-changes vs 0 progressive planner rows
- [ ] Code audit: PlanExists vs PlanCritiqued vs evaluateDone vs MCP surface
- [ ] Fixture filesystem audit: missing install/config artifacts
- [ ] Write INVESTIGATION.md (10 sections); update board row

## Next

`P36-S00-02`
