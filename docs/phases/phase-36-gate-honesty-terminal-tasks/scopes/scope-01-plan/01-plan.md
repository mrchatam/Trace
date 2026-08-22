# P36-S01-01 — Plan

## Metadata
- id: P36-S01-01
- todo_ids: [P36-S01-01]
- role: implementer
- skills: [planning-and-task-breakdown, api-and-interface-design, domain-modeling, documentation-and-adrs]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-01-plan/PLAN.md`: locked policy, touch-list, acceptance tests, S02 handoff. Pick a **fundamental** fix set from DESIGN-LOCKS — MCP plan / bootstrap / install are **first-class required decisions**, not optional footnotes. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) — S00-02 PASS (high confidence)
- [00-PLANNER.md](00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors:
  - `internal/loop/policy.go:45–62` — `PlanExists` vs `PlanCritiqued`
  - `internal/loop/gate.go:227–265` — `evaluateDone` (no terminal short-circuit)
  - `internal/deliberation/select.go:28–29` — `!PlanExists` → PLAN / `plan_missing`
  - `cmd/trace/plan.go:46–67` — CLI subcommands; *"No MCP plan tools."*
  - `internal/mcp/server.go:216–224` — 15 registered tools (no `trace_plan`)
  - `internal/mcp/tools_loop.go:39–47` — `trace_loop` = next|apply|status only
  - `internal/mcp/tools_write.go:130–136`, `:224–231` — plan-change add; transition without gate
  - `internal/install/enforcement.go:51–98` — AGENTS/cursor rules (no create-coarse bootstrap)
  - `internal/config/enforce.go:21–27` — missing config → enforce off
  - `web/src/screens/TaskDetail.tsx:76–86`, `:198–210`
  - `web/src/components/GateStrip.tsx:41–56`
  - Test patterns: `cmd/trace/loop_test.go:1137–1167`, `internal/loop/gate_test.go:198–205`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md` **only** |
| Product edits | **Forbidden** on this row |
| S00 verdict | **Accepted** — Trace product primary; agent+harness secondary; GUI tertiary |
| Dogfood fixture | `/home/ali/Desktop/feet seller telegram app` (read-only unless PLAN scopes recovery write) |
| Primary fix axis | Trace product gaps (MCP / bootstrap / install) — not "agents should have known" |
| GUI | Adapter only (Law 19); terminal gate honesty **secondary** to agent-complete path |
| Preserve | Active-work PLAN enforcement (`PlanExists` for non-terminal tasks); fail-closed deliberation order |
| Touch-list order | **library → MCP → install → HTTP → GUI** (Law 19) |
| S02 SoT | PLAN.md touch-list + acceptance tests govern implement row |

## S00 verdict lock (do not re-litigate)

Summarize in PLAN.md §1 from INVESTIGATION.md:

| Layer | Verdict | One-line evidence |
|-------|---------|-------------------|
| Trace product | **Primary** | 11 plan-changes vs 0 progressive planner rows; gate reads planner only; MCP cannot populate |
| Agent misuse | Secondary | Could use CLI `create-coarse` — undiscoverable from MCP/install |
| Harness | Secondary | No install; enforce off; 127 PASS reviews without gate block |
| GUI | Tertiary | Mirrors library JSON on DONE tasks; not root cause |

**Two planning systems (locked fact):**

| System | Storage | Agent surface | Gate signal |
|--------|---------|---------------|-------------|
| Causal graph | `plan_changes` | `trace_add kind=plan-change` | `PlanCritiqued` only (apply path) |
| Progressive planner | `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` | CLI `trace plan create-coarse\|set-current\|deep\|show` | **`PlanExists`** |

## Decision framework (PLAN.md must apply to each candidate)

For **every** DESIGN-LOCKS candidate below, PLAN.md must record **accept | defer | reject** with:

1. **Rationale** — tied to S00 evidence (not speculation)
2. **S02 scope** — if accept: concrete deliverable; if defer: trigger + owner scope
3. **Risk** — especially for bridge/heuristic options
4. **Acceptance hook** — which test(s) prove the choice

**Default recommendation from S00 (implementer may override with rationale):**

| Rank | Option | S00 default |
|------|--------|-------------|
| 1 | MCP plan tools | **accept** |
| 2 | Install contract | **accept** |
| 3 | Bootstrap command | **accept** |
| 4 | Terminal gate | **accept** |
| 5 | Goal structure warning | accept or defer |
| 6 | Enforce nudge | accept or defer |
| 7 | PlanExists bridge | **defer** (unless rejecting pure MCP path) |

---

## PLAN.md required sections (author in order)

### §1 Root-cause lock

- S00 verdict table (above) accepted
- Feet-seller repro facts locked: 123 DONE / 1 goal; export 11/0; identical `plan_missing` on Step1+Loop112
- Explicit reject: GUI-only patch; global `plan_missing` weaken; deleting history

### §2 Fix set decision — all candidates first-class

Each subsection **must exist** even if deferred/rejected.

#### §2.1 MCP plan tools (`trace_plan`)

**Problem:** Primary agent surface (MCP) cannot satisfy `PlanExists`. CLI-only `trace plan` subcommands exist; MCP has 15 tools, no plan tool (`cmd/trace/plan.go:67`).

**Decision questions for PLAN.md:**

- Accept `trace_plan` MCP tool mirroring CLI: `create-coarse`, `set-current`, `deep`, `show`?
- Law 19: thin adapter over `internal/planner.Service` — same as `cmd/trace/plan.go`
- Tool shape: single `trace_plan` with `action` param vs separate tools?
- Registration: update `RegisteredToolNames()` (currently 15 → 16)
- Error contract: match CLI stderr/stdout JSON semantics

**If accept — S02 touch:** `internal/mcp/tools_plan.go` (new), `internal/mcp/server.go`, `internal/mcp/mcp_test.go`

**Acceptance hook:** greenfield MCP bootstrap test (§7.1)

#### §2.2 Bootstrap command (`trace plan bootstrap --goal`)

**Problem:** Feet-seller-like repos have rich plan-changes (11) but 0 progressive planner rows. Manual `create-coarse` requires human/agent to infer structure from 123 tasks.

**Decision questions:**

- Accept new CLI subcommand `trace plan bootstrap --goal <id>`?
- Input sources: plan-change titles/bodies? seed import? task clustering heuristic?
- Output: minimal valid coarse plan (at least one phase/scope + set-current + deep) so `PlanExists` becomes true
- Idempotency: safe to re-run on partially bootstrapped goal?
- MCP exposure: bootstrap via `trace_plan action=bootstrap` or CLI-only recovery?

**If accept — S02 touch:** `cmd/trace/plan.go`, `internal/planner/` (bootstrap helper if needed)

**Acceptance hook:** feet-seller recovery test (§7.2 post-bootstrap edit gate passes)

#### §2.3 Install contract (mandatory coarse plan on new goal)

**Problem:** `internal/install/enforcement.go` documents loop gate hooks but **no** `create-coarse` bootstrap. Feet-seller had no install artifacts → enforce off → 123 DONE without gate.

**Decision questions:**

- Accept install contract amendment: **first task under a new goal → mandatory `trace plan create-coarse`** before edit gate can pass?
- Where documented: `AgentsEnforcementBlock()`, `cursorRulesMDCContent()`, `EnforcementRulesMarkdown()`
- AGENTS.md template: add bootstrap step to `# begin-trace-enforcement` block
- Cursor rules: same block in `.cursor/rules/trace-enforcement.mdc`
- Optional: `trace install` post-install hint when goal exists but planner empty
- **Not** changing default enforce mode to strict (preserve opt-in)

**If accept — S02 touch:** `internal/install/enforcement.go`, possibly `cmd/trace/install.go` help text

**Acceptance hook:** install output/docs review + greenfield workflow includes bootstrap step

#### §2.4 PlanExists bridge (plan-change → progressive planner heuristic)

**Problem:** Agents naturally write plan-changes; gate ignores them for `PlanExists`.

**Decision questions:**

- Accept bridge that sets `PlanExists` from plan-change density/titles without real progressive planner rows?
- **Hard constraint:** must NOT fake progressive plan — if bridge sets `PlanExists`, what store rows are created?
- Alternative: bridge only **recommends** bootstrap (advisory), does not satisfy gate
- S00 default: **defer** unless MCP+bootstrap rejected

**If accept — S02 touch:** `internal/loop/policy.go` (high risk — justify heavily)

**If defer — PLAN.md must state:** deferred to future phase; MCP+bootstrap is preferred path

#### §2.5 Terminal gate (DONE/SKIPPED honesty)

**Problem:** `evaluateDone` (`gate.go:227–265`) has no `work_state` short-circuit. All 123 DONE feet-seller tasks show actionable "done blocked: recommended phase PLAN (plan_missing)".

**Decision questions:**

- Accept library short-circuit for terminal tasks (`work_state` DONE or SKIPPED)?
- Semantics: `allowed=true` with advisory? distinct `reason_code` (e.g. `plan_missing_terminal_advisory`)? separate violation severity?
- **Must preserve:** active non-terminal tasks still fail-closed on `!PlanExists`
- Goal-level advisory: surface that **goal** lacks progressive plan even when task is terminal
- HTTP/GUI adapter follows library JSON — no independent weaken

**If accept — S02 touch:** `internal/loop/gate.go`, `internal/loop/gate_test.go`

**Acceptance hook:** feet-seller terminal honesty test (§7.2)

#### §2.6 Enforce nudge (`.trace/` without config)

**Problem:** Missing `.trace/config.json` → `EnforceOff` (`internal/config/enforce.go:21–27`). Agents completed work without ever seeing gate block.

**Decision questions:**

- Accept: when `.trace/` exists but config missing/invalid, default or nudge to `warn`?
- Or: document-only nudge in `trace init` / first MCP call / install message?
- **Must preserve:** default-off for repos without `.trace/`; strict remains opt-in

**If accept — S02 touch:** `internal/config/enforce.go`, possibly `cmd/trace/init.go`

#### §2.7 Goal structure warning (>N tasks, no coarse plan)

**Problem:** 123 tasks / 1 goal amplifies identical gate on every task detail.

**Decision questions:**

- Accept warning when task count under goal exceeds threshold N (e.g. 10? 20?) without coarse plan?
- Surface: CLI `plan show` stderr? MCP response field? loop status advisory?
- Threshold N value and rationale

**If accept — S02 touch:** `internal/planner/` or `internal/loop/` advisory helper

#### §2.8 Feet-seller recovery (S02 vs deferred)

**Decision questions:**

- **In S02:** run bootstrap on live fixture (with user approval) OR ship command + docs only?
- **Deferred:** document manual recovery steps in PLAN.md / VERIFY without mutating fixture
- Minimum: post-bootstrap, edit gate passes on a representative DONE task's goal
- **Reject:** deleting or rewriting 123 DONE tasks / 11 plan-changes / 127 reviews

**PLAN.md must pick one:** `S02 execute recovery` | `S02 ship tool + S03 verify` | `defer to post-36`

### §3 Agent workflow (MCP-first)

Step-by-step greenfield path PLAN.md must document (happy path after fix set):

```text
1. trace init (or MCP equivalent)
2. Create goal + first task (trace_add / CLI)
3. Bootstrap coarse plan:
     MCP trace_plan create-coarse → set-current → deep
     OR trace plan bootstrap --goal (recovery)
4. trace loop gate --for edit → allowed=true, plan_exists=true
5. Agent work: trace_add discoveries, trace_loop apply, …
6. trace transition --to DONE (with review)
7. Terminal task: gate shows honest advisory (not actionable block) if goal plan gap remains
```

Include failure modes: agent writes plan-change only → still `plan_missing` until bootstrap (documents why bridge defer is OK).

### §4 Preserve — active-work PLAN enforcement

PLAN.md must explicitly lock:

- Non-terminal tasks with `!PlanExists` → `plan_missing`, edit/done blocked (unchanged)
- Deliberation order fail-closed (`select.go:28–29`)
- No global weaken of PLAN phase for active work
- `PlanCritiqued` via plan-changes does **not** substitute for `PlanExists`

### §5 Touch list (ordered — Law 19)

PLAN.md must list files/packages in **this order** with per-item rationale:

| Order | Layer | Typical paths | Notes |
|------:|-------|---------------|-------|
| 1 | Library policy | `internal/loop/policy.go`, `internal/loop/gate.go`, `internal/deliberation/` | Terminal gate; optional goal-structure advisory |
| 2 | Planner library | `internal/planner/service.go`, bootstrap helper if any | Shared by CLI + MCP |
| 3 | MCP adapter | `internal/mcp/tools_plan.go`, `server.go`, tests | `trace_plan` registration |
| 4 | CLI adapter | `cmd/trace/plan.go` | `bootstrap` subcommand |
| 5 | Install contract | `internal/install/enforcement.go`, `cmd/trace/install.go` | AGENTS + cursor bootstrap docs |
| 6 | Config nudge | `internal/config/enforce.go` | If enforce nudge accepted |
| 7 | HTTP adapter | `internal/http/` plan routes if any exist | Mirror MCP; defer if no HTTP plan routes today |
| 8 | GUI adapter | `web/src/screens/TaskDetail.tsx`, `web/src/components/GateStrip.tsx` | Secondary; follow library JSON only |

**Explicit non-touch unless PLAN accepts:** weakening `PlanExists` check globally; GUI-only hide without library change.

### §6 Acceptance tests (S02 must implement)

PLAN.md must name these tests (adapt names as needed):

#### §6.1 Greenfield MCP workflow

**Name:** `TestGreenfield_MCPPlanBootstrap_EditGatePasses` (from INVESTIGATION.md)

**Setup:** Temp dir, `trace init`, goal + task via MCP/CLI.

**Act:** MCP `trace_plan` create-coarse + set-current + deep (or CLI until MCP lands).

**Assert:**

- `trace loop gate --for edit` → `allowed: true`, `reason_code` ≠ `plan_missing`
- `trace loop status` → `policy_inputs.plan_exists: true`

**Pattern:** `createCurrentDeepPlanForLoopTest` in `cmd/trace/loop_test.go:1137–1167`

#### §6.2 Feet-seller live / import honesty

**Name:** `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`

**Setup:** Import feet-seller export (`plan_changes: 11`, planner arrays 0) into temp store.

**Assert (pre terminal-gate fix):** DONE task `--for done` → documents current `plan_missing` behavior.

**Assert (post terminal-gate fix):** DONE task → honest terminal signal (not identical actionable "done blocked" alarm).

**Assert (post bootstrap):** after `plan bootstrap`, edit gate passes; `plan_exists: true`.

**Fixture source:** `/tmp/feet-export.json` from fixture export or checked-in anonymized subset.

#### §6.3 Active-work preservation

**Name:** `TestActiveWork_PlanMissingStillBlocksEdit` (extend `gate_test.go:198–205`)

**Assert:** non-terminal task on goal without planner → edit gate still `plan_missing`, `allowed: false`.

### §7 Explicit non-goals

PLAN.md must list (reject with rationale):

1. **GUI-only patch** — hiding `plan_missing` in GateStrip/TaskDetail without library + agent path fix
2. **Global `plan_missing` weaken** — for active non-terminal PLAN work
3. **Deleting feet-seller history** — 123 tasks, 11 plan-changes, 127 reviews
4. **Blaming Phase 35 pick logic** — selection fixed; this is gate meaning
5. **Hosted SaaS / multi-tenant** — out of phase scope
6. **Auto-LLM backlog generation** — per `cmd/trace/plan.go` help

### §8 S02 handoff

- Ordered touch-list with accept/defer per file
- Acceptance test names + file locations for S02-01
- Feet-seller recovery scope (execute vs verify-only)
- Residual risks for S02 reviewer

---

## Preflight / Plan

Before writing PLAN.md:

1. Re-read INVESTIGATION.md §Fundamental fix options + §Red-capable tests
2. Skim live files listed in References (confirm line anchors still valid)
3. Decide each §2 candidate using decision framework
4. Draft PLAN.md sections §1–§8 in order
5. Self-check against Exit criteria below

## Role work

1. Apply decision framework to all 8 DESIGN-LOCKS candidates (+ feet-seller recovery)
2. Write `PLAN.md` with sections §1–§8
3. Do **not** edit product code, tests, or GUI
4. Update board row P36-S01-01 status + Notes only

## Todo updates

Per board-rights: implementer sets own row `done` + Notes on P36-S01-01 only.

## Exit criteria

- [ ] `PLAN.md` exists at `scopes/scope-01-plan/PLAN.md` with locked fix choices
- [ ] All 8 DESIGN-LOCKS candidates addressed (accept/defer/reject + rationale) in §2
- [ ] MCP plan / bootstrap / install are **first-class** (dedicated §2.1–§2.3 subsections, not buried)
- [ ] Touch-list ordered library → MCP → install → HTTP → GUI in §5
- [ ] Acceptance tests §6.1 + §6.2 + §6.3 named with assert shapes
- [ ] Explicit non-goals §7 present
- [ ] Active-work PLAN enforcement preserved in §4
- [ ] **No product code** edited in this row
- [ ] Board Notes on P36-S01-01 only

## Minimal todos

- [ ] Lock S00 verdict in PLAN.md §1
- [ ] Decide §2.1 MCP plan tools
- [ ] Decide §2.2 Bootstrap command
- [ ] Decide §2.3 Install contract
- [ ] Decide §2.4 PlanExists bridge (default defer)
- [ ] Decide §2.5 Terminal gate
- [ ] Decide §2.6 Enforce nudge
- [ ] Decide §2.7 Goal structure warning
- [ ] Decide §2.8 Feet-seller recovery scope
- [ ] Write §3 MCP-first agent workflow
- [ ] Write §5 ordered touch-list
- [ ] Write §6 acceptance tests
- [ ] Write §7 non-goals + §8 S02 handoff
- [ ] Mark P36-S01-01 done with Notes

## Next

`P36-S02-00` (after PLAN.md complete)
