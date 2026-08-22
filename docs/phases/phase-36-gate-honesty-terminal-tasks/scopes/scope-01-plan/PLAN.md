# Phase 36 S01 — PLAN

**Author:** P36-S01-01 (2026-08-22)  
**SoT inputs:** [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) (S00-02 PASS, high confidence), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [INTAKE.md](../../INTAKE.md)  
**Dogfood fixture (read-only in S01):** `/home/ali/Desktop/feet seller telegram app`

---

## §1 Root-cause lock

### S00 verdict (accepted — do not re-litigate)

| Layer | Verdict | One-line evidence |
|-------|---------|-------------------|
| **Trace product** | **Primary** | 11 `plan_changes` vs 0 progressive planner rows; `PlanExists` reads planner only (`internal/loop/policy.go:45–49`); MCP cannot populate (`cmd/trace/plan.go:67`) |
| **Agent misuse** | **Secondary** | Agents could use CLI `create-coarse` — undiscoverable from MCP/install; rich plan-change titles show planning *intent* in wrong store |
| **Harness** | **Secondary** | No install artifacts; `.trace/config.json` absent → `EnforceOff` (`internal/config/enforce.go:21–27`); 127 PASS reviews without gate block |
| **GUI** | **Tertiary** | `TaskDetail.tsx:76–86` fetches done gate on DONE tasks; `GateStrip.tsx:41–56` red “Gate blocked” — mirrors library JSON, not root cause |

### Two planning systems (locked fact)

| System | Storage | Agent surface | Gate signal |
|--------|---------|---------------|-------------|
| **Causal graph** | `plan_changes` | `trace_add kind=plan-change` (`internal/mcp/tools_write.go:130–136`) | `PlanCritiqued` only on apply path (`internal/loop/policy.go:60–62`) — **not** `PlanExists` |
| **Progressive planner** | `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` | CLI `trace plan create-coarse\|set-current\|deep\|show` | **`PlanExists`** — requires `current_scope_id` **and** `current_deep_plan` (`internal/loop/policy.go:45–49`) |

### Feet-seller repro facts (locked)

| Fact | Value | Source |
|------|-------|--------|
| Tasks | 123, all `work_state: DONE` | INVESTIGATION.md §Repro |
| Goals | 1 (`353b12a4-57dd-4d68-8379-b2024e064733`) | INVESTIGATION.md §Repro |
| Export counts | **11** plan-changes / **0** planner rows | `/tmp/feet-export.json` |
| Gate JSON | Identical `plan_missing` on Step1 + Loop112 (`--for done` and `--for edit`) | INVESTIGATION.md §Repro |
| Transitions | 127/127 PASS reviews; enforce off; MCP `trace_transition` skips gate (`tools_write.go:224–231`) | INVESTIGATION.md §Transition audit |

Goal-level inheritance: `loadGateContext` resolves `task.GoalID` → `BuildPolicyInputs` → `PlanExists` from goal planner state, not task `work_state` (`internal/loop/gate.go:135–153`, `policy.go:45–49`). `evaluateDone` has no terminal short-circuit (`internal/loop/gate.go:227–265`).

### Explicit rejects

| Reject | Rationale |
|--------|-----------|
| GUI-only patch | Hides symptom; MCP/install gap remains; agents still cannot satisfy `PlanExists` |
| Global `plan_missing` weaken | DESIGN-LOCKS + INTAKE: preserve PLAN enforcement for active non-terminal work |
| Deleting feet-seller history | 123 DONE tasks, 11 plan-changes, 127 reviews — recovery/bootstrap preferred |
| Blaming Phase 35 pick logic | Selection fixed; this is gate **meaning** on terminal DONE + empty progressive planner |
| Hosted SaaS / multi-tenant | Out of phase scope (DESIGN-LOCKS) |

### Locked fix set summary (S02 implements)

**Accept:** MCP plan tools (§2.1), bootstrap command (§2.2), install contract (§2.3), terminal gate honesty (§2.5), goal structure warning (§2.7), enforce document nudge (§2.6).  
**Defer:** PlanExists bridge (§2.4).  
**Recovery:** Ship bootstrap in S02; live feet-seller verify in S03 (§2.8).

---

## §2 Fix set decisions

Each candidate: **accept | defer | reject** with S00 evidence, S02 scope, risk, acceptance hook.

### §2.1 MCP plan tools (`trace_plan`) — **ACCEPT**

**Problem:** Primary agent surface (MCP) cannot satisfy `PlanExists`. CLI `trace plan` subcommands exist; MCP registers 15 tools, none for plan (`internal/mcp/server.go:216–224`; help: *"No MCP plan tools."* at `cmd/trace/plan.go:67`).

**Decision:**

| Question | Lock |
|----------|------|
| Accept `trace_plan` mirroring CLI? | **Yes** — actions: `create-coarse`, `set-current`, `deep`, `show`, `bootstrap` (see §2.2) |
| Law 19 shape | Thin adapter over `internal/planner.Service` — same calls as `cmd/trace/plan.go` |
| Tool shape | **Single `trace_plan` with `action` param** — matches `trace_loop` pattern (`internal/mcp/tools_loop.go:39–47`) |
| Registration | Add to `RegisteredToolNames()` → **16 tools**; update `internal/mcp/mcp_test.go` locked name list |
| Error contract | Match CLI stderr messages + stdout JSON from planner service (snake_case plan view on `show`) |

**Rationale:** INVESTIGATION.md §MCP audit — agents on MCP naturally write plan-changes but cannot run `create-coarse`. This closes the primary product gap without changing gate semantics.

**S02 deliverables:** `internal/mcp/tools_plan.go` (new), `internal/mcp/server.go` (register handler), `internal/mcp/mcp_test.go` (registration + greenfield workflow).

**Risk:** Low — adapter-only; no policy change.

**Acceptance hook:** `TestGreenfield_MCPPlanBootstrap_EditGatePasses` (§6.1).

---

### §2.2 Bootstrap command (`trace plan bootstrap --goal`) — **ACCEPT**

**Problem:** Feet-seller-like repos have 11 plan-changes, 0 progressive planner rows. Manual `create-coarse` requires inferring structure from 123 tasks (`INVESTIGATION.md` §Fundamental fix options #3).

**Decision:**

| Question | Lock |
|----------|------|
| Accept CLI subcommand? | **Yes** — `trace plan bootstrap --goal <id>` |
| Input sources | Goal-linked `plan_changes` (titles + bodies via store/domain); **no** LLM generation (consistent with `cmd/trace/plan.go:67`) |
| Heuristic | Pick primary plan-change (prefer most recent non-superseded by title length/recency); derive **one phase** (title from plan-change or `"Recovery"`) + **one scope** (truncated title) + minimal `deep` (`--exit` from plan-change body excerpt or fixed `"Bootstrap from plan-changes"`) |
| Output | Valid progressive planner state: `create-coarse` → `set-current` → `deep` so `PlanExists` becomes true |
| Idempotency | If goal already has `current_scope_id` + `current_deep_plan`, **no-op success** with stderr note |
| MCP exposure | **`trace_plan action=bootstrap`** — recovery must be agent-complete, not CLI-only |

**Rationale:** INVESTIGATION.md export 11/0 + sample titles (*Architecture plan*, *Production-ready MVP plan*) prove agents planned in causal graph. Bootstrap converts existing intent into gate-checked store without faking rows silently (explicit command + audit trail).

**S02 deliverables:** `internal/planner/bootstrap.go` (helper), `cmd/trace/plan.go` (subcommand), wire through MCP `tools_plan.go`.

**Risk:** Medium — heuristic quality varies; bootstrap produces **minimal** coarse plan, not full reconstruction of 123-task history. Document in command help that human refinement via `create-coarse` / `deep` is expected.

**Acceptance hook:** `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap` post-bootstrap assert (§6.2).

---

### §2.3 Install contract (mandatory coarse plan on new goal) — **ACCEPT**

**Problem:** `internal/install/enforcement.go` documents pre-edit/pre-DONE gate hooks (`:57–58`, `:73–74`, `:92–93`) but **grep `create-coarse` under `internal/install/` → no matches**. Feet-seller had no AGENTS.md, cursor rules, or hooks → enforce off path (`INVESTIGATION.md` §Install audit).

**Decision:**

| Question | Lock |
|----------|------|
| Accept install contract amendment? | **Yes** — document mandatory progressive planner bootstrap for new goals |
| Contract text | Before first edit on a goal without `PlanExists`: run `trace plan create-coarse …` then `set-current` + `deep` (or `trace plan bootstrap --goal` for recovery) |
| Where documented | `AgentsEnforcementBlock()` (`:67–79`), `cursorRulesMDCContent()` (`:82–97`), `EnforcementRulesMarkdown()` (`:52–63`) |
| AGENTS block | Add numbered step **between** TRACE_TASK_ID setup and pre-edit gate: *"New goal without coarse plan: bootstrap via `trace plan create-coarse` or MCP `trace_plan` before edit gate can pass."* |
| Post-install hint | `cmd/trace/install.go` agents/cursor success stderr: mention bootstrap when store has goals but empty planner |
| Default enforce mode | **Unchanged** — opt-in `off\|warn\|strict`; do not flip default to strict |

**Rationale:** Harness secondary blame (INVESTIGATION.md verdict) — install never told agents *how* to satisfy `PlanExists`. Contract closes discoverability gap without mandating strict enforce.

**S02 deliverables:** `internal/install/enforcement.go`, `cmd/trace/install.go` (hint text), install tests if present.

**Risk:** Low — documentation + harness contract only.

**Acceptance hook:** Manual review of install output + greenfield workflow step 3 (§3); optional snapshot test on `AgentsEnforcementBlock()` containing `create-coarse`.

---

### §2.4 PlanExists bridge (plan-change → progressive planner heuristic) — **DEFER**

**Problem:** Agents naturally write plan-changes; gate ignores them for `PlanExists` (`internal/deliberation/select.go:28–29`).

**Decision:** **Defer** to a future phase. Do **not** auto-satisfy `PlanExists` from plan-change density/titles without explicit bootstrap command.

| Question | Lock |
|----------|------|
| Accept silent bridge? | **No** — would fake progressive plan or create ambiguous store rows |
| Advisory-only bridge? | **Defer** — optional future: loop status field recommending bootstrap when plan-changes ≥ N and `!PlanExists` |
| Preferred path | MCP `trace_plan` + explicit `bootstrap` (§2.1–§2.2) |

**Rationale:** S00 default + DESIGN-LOCKS: *"must not fake progressive plan."* MCP+bootstrap gives honest, auditable path. Bridge risks agents believing plan-changes alone satisfy gate (perpetuates dual-system confusion).

**S02 scope:** None (no `policy.go` bridge).

**Risk if accepted:** High — weakens planning contract; hard to test heuristic edge cases.

**Acceptance hook:** N/A — defer documented; §6.3 ensures active work still blocked without bootstrap.

**Trigger to revisit:** If S03 VERIFY shows agents still skip bootstrap after MCP+install ship; owner: future phase planner.

---

### §2.5 Terminal gate (DONE/SKIPPED honesty) — **ACCEPT**

**Problem:** `evaluateDone` (`internal/loop/gate.go:227–265`) evaluates deliberation for all tasks regardless of `work_state`. Feet-seller DONE tasks emit actionable `"done blocked: recommended phase PLAN (plan_missing)"` identical to active work.

**Decision:**

| Question | Lock |
|----------|------|
| Accept library short-circuit? | **Yes** — when `task.WorkState` is `DONE` or `SKIPPED` |
| `--for done` semantics | `allowed: true`; emit **advisory** violation (non-blocking): `reason_code: goal_plan_gap_terminal_advisory`, message explains goal lacks progressive plan, work already terminal |
| `--for edit` on terminal task | Same advisory pattern: `allowed: true` with advisory (editing finished work is rare; blocking with red alarm is misleading) |
| Active non-terminal | **Unchanged** — `!PlanExists` → `plan_missing`, `allowed: false` |
| Goal-level signal | Advisory message includes `goal_id` and hint: `trace plan bootstrap --goal` or MCP `trace_plan` |
| HTTP/GUI | Follow library JSON only — no independent weaken (Law 19) |

**Rationale:** DESIGN-LOCKS *"Must fix (honesty)"* — terminal tasks must not imply actionable done-blocked for finished work. Root fix is MCP+bootstrap; honesty fix stops 123 identical red alarms while goal plan gap persists.

**S02 deliverables:** `internal/loop/gate.go` (`evaluateDone`, optionally `evaluateEdit` terminal branch), `internal/loop/gate_test.go`.

**Risk:** Low if scoped strictly to terminal `work_state`; medium if mis-scoped (must not bypass verification/regression blocks on terminal tasks — keep existing checks at `gate.go:243–257` **before** terminal advisory path).

**Acceptance hook:** `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap` terminal assert (§6.2); extend `gate_test.go` with `TestEvaluateGate_Done_TerminalPlanGapAdvisory`.

---

### §2.6 Enforce nudge (`.trace/` without config) — **ACCEPT (document-only)**

**Problem:** Missing `.trace/config.json` → `LoadEnforceMode` returns `EnforceOff` (`internal/config/enforce.go:21–27`). Feet-seller completed 123 DONE without gate ever blocking (`INVESTIGATION.md` §Transition audit).

**Decision:**

| Question | Lock |
|----------|------|
| Change default enforce mode? | **No** — preserve `EnforceOff` when config missing |
| Accept document nudge? | **Yes** — one-time stderr hint when `.trace/` exists but config missing/invalid |
| Surfaces | `trace init` completion message; first MCP `openStore` per process (optional, once); `trace install` success output |
| Suggest action | *"Consider `.trace/config.json` with `\"enforce\": \"warn\"` after `trace install`"* |

**Rationale:** Secondary harness blame — agents never saw enforcement. Document nudge raises visibility without breaking opt-in model (DESIGN-LOCKS preserve strict as opt-in).

**S02 deliverables:** `cmd/trace/init.go` (hint), optional helper in `internal/config/enforce.go` (`WarnIfTraceDirWithoutConfig`), `cmd/trace/install.go`.

**Risk:** Low — stderr only; no behavior change to gate evaluation.

**Acceptance hook:** CLI test or snapshot on init output when `.trace/` present, no config.

---

### §2.7 Goal structure warning (>N tasks, no coarse plan) — **ACCEPT**

**Problem:** 123 tasks / 1 goal amplifies identical gate on every task detail (`INVESTIGATION.md` §Mega-goal pattern). Even with terminal honesty, operators need early warning before mega-goal forms.

**Decision:**

| Question | Lock |
|----------|------|
| Accept warning? | **Yes** |
| Threshold **N** | **15** — below feet-seller 123 but catches unbounded loop growth early; aligns with typical phase scope sizes |
| Condition | `task_count(goal) > N && !PlanExists(goal)` |
| Surfaces | `trace plan show --goal` stderr advisory; MCP `trace_plan action=show` adds `"goal_structure_warning"` field; `trace loop status` optional `advisories[]` entry |
| Not a gate bypass | Warning only — does not set `PlanExists` or weaken edit block |

**Rationale:** Product guidance gap (INVESTIGATION.md #6) — complements bootstrap/install without bridge heuristic.

**S02 deliverables:** `internal/planner/advisory.go` or helper on `GetPlan`; wire into `cmd/trace/plan.go` show + MCP show + optionally `internal/loop/next.go` status.

**Risk:** Low — advisory only.

**Acceptance hook:** Unit test `TestGoalStructureWarning_OverThresholdNoPlan`; feet-seller import shows warning on `plan show`.

---

### §2.8 Feet-seller recovery — **S02 ship tool + S03 verify**

**Decision:** **S02 ship tool + S03 verify** (not in-fixture write in S02 without explicit user approval).

| Question | Lock |
|----------|------|
| S02 scope | Implement `trace plan bootstrap` + MCP `trace_plan bootstrap`; export-based test in temp store (§6.2) |
| S03 scope | Live verify on dogfood fixture: run bootstrap on goal `353b12a4-…`, confirm edit gate passes, GUI shows honest advisory/OK |
| Minimum proof | Post-bootstrap: `loop gate --for edit` → `allowed: true`, `plan_exists: true` on representative task |
| Reject | Deleting/rewriting 123 DONE tasks, 11 plan-changes, 127 reviews |

**Rationale:** INVESTIGATION.md dogfood fixture is read-only evidence; recovery command is the product fix; live mutation belongs in VERIFY with human promotion (INTAKE desired outcome #3).

**S02 deliverables:** Bootstrap implementation + import test only.

**Acceptance hook:** §6.2 temp import; S03 live checklist (out of S02).

---

## §3 Agent workflow (MCP-first greenfield happy path)

After S02 fix set, the canonical agent path:

```text
1. trace init
   └─ stderr nudge if no .trace/config.json (§2.6)

2. Create goal + first task
   └─ trace_add kind=goal / kind=task (MCP) or CLI equivalent

3. Bootstrap coarse plan (MANDATORY before edit gate passes)
   └─ MCP trace_plan action=create-coarse --goal <id> --phase "Phase 1" --scope "Scope 1"
   └─ MCP trace_plan action=set-current --goal <id> --scope <scope_id>
   └─ MCP trace_plan action=deep --scope <scope_id> --exit "…" --work "…"
   OR recovery: trace_plan action=bootstrap --goal <id>

4. Verify gate
   └─ trace loop gate --task <id> --for edit → allowed=true, reason_code ≠ plan_missing
   └─ trace loop status → policy_inputs.plan_exists: true

5. Agent work loop
   └─ trace_add discoveries/decisions/plan-changes
   └─ trace_loop next|apply|status
   └─ plan-changes satisfy PlanCritiqued on apply — NOT PlanExists

6. Complete task
   └─ trace_review → trace_transition --to DONE

7. Terminal task with lingering goal plan gap (legacy repos only)
   └─ gate --for done → allowed=true + goal_plan_gap_terminal_advisory (§2.5)
   └─ GUI GateStrip: warn tone, not error block
```

### Failure modes (documented)

| Agent action | Gate outcome | Remedy |
|--------------|--------------|--------|
| Writes plan-change only, no `trace_plan` bootstrap | `plan_missing` on edit/done (active tasks) | Run `trace_plan create-coarse` chain or `bootstrap` |
| Skips install / enforce off | Gate not invoked by harness | `trace install agents`; set `enforce: warn` optionally |
| Mega-goal without coarse plan | Advisory at `plan show`; repeated per-task signal until bootstrap | Bootstrap early; heed N>15 warning (§2.7) |

Bridge defer (§2.4) is intentional: plan-changes alone **must not** satisfy `PlanExists` until explicit bootstrap.

---

## §4 Preserve — active-work PLAN enforcement

Locked invariants (unchanged by Phase 36):

| Invariant | Anchor |
|-----------|--------|
| Non-terminal + `!PlanExists` → edit/done blocked, `reason_code: plan_missing` | `internal/deliberation/select.go:28–29`, `gate_test.go:198–205` |
| Deliberation order fail-closed | PLAN before CRITIQUE before EXECUTE (`select.go:28–38`) |
| `PlanExists` requires progressive planner rows | `internal/loop/policy.go:45–49` |
| `PlanCritiqued` via plan-changes does **not** substitute for `PlanExists` | `policy.go:60–62` vs `:45–49` |
| Verification/regression blocks on done | `gate.go:243–257` — preserved even with terminal advisory |
| No global weaken of PLAN phase for active work | DESIGN-LOCKS, INTAKE |

S02 must extend `TestActiveWork_PlanMissingStillBlocksEdit` (§6.3) to guard regressions.

---

## §5 Touch list (ordered — Law 19)

Policy in library first; adapters thin-wrap shared services.

| Order | Layer | Paths | Accept | Rationale |
|------:|-------|-------|--------|-----------|
| 1 | Library — gate | `internal/loop/gate.go`, `internal/loop/gate_test.go` | §2.5 | Terminal honesty in canonical evaluator |
| 2 | Library — deliberation | `internal/deliberation/select.go` | — (no change) | Order preserved; cite only |
| 3 | Library — policy | `internal/loop/policy.go` | — (no change) | `PlanExists` definition unchanged; bridge deferred |
| 4 | Planner library | `internal/planner/bootstrap.go`, `internal/planner/advisory.go`, `internal/planner/service.go` (if needed) | §2.2, §2.7 | Shared bootstrap + goal warning for CLI/MCP |
| 5 | MCP adapter | `internal/mcp/tools_plan.go` (new), `internal/mcp/server.go`, `internal/mcp/mcp_test.go` | §2.1 | `trace_plan` registration (16 tools) |
| 6 | CLI adapter | `cmd/trace/plan.go` | §2.2 | `bootstrap` subcommand + show stderr advisory |
| 7 | Install contract | `internal/install/enforcement.go`, `cmd/trace/install.go` | §2.3 | Bootstrap step in AGENTS/cursor blocks |
| 8 | Config nudge | `internal/config/enforce.go`, `cmd/trace/init.go` | §2.6 | Document-only enforce hint |
| 9 | HTTP adapter | `internal/httpapi/handlers_p1.go`, `internal/httpapi/server.go` | **Defer mutation routes** | `GET /v1/plans` exists (`server.go:282`); no POST create-coarse/set-current/deep today — mirror MCP only if routes added in follow-up |
| 10 | GUI adapter | `web/src/components/GateStrip.tsx`, `web/src/screens/TaskDetail.tsx` | §2.5 secondary | Map `allowed: true` + advisory reason to warn banner; no hide-without-library-change |

**Explicit non-touch:**

- Weakening `PlanExists` check globally
- GUI-only suppression of `plan_missing` without library terminal semantics
- `internal/loop/policy.go` bridge heuristic (§2.4 defer)

---

## §6 Acceptance tests (S02 implements)

Reference patterns: `createCurrentDeepPlanForLoopTest` (`cmd/trace/loop_test.go:1137–1167`), `TestEvaluateGate_Edit_PlanMissing` (`internal/loop/gate_test.go:198–205`).

### §6.1 Greenfield MCP workflow

**Name:** `TestGreenfield_MCPPlanBootstrap_EditGatePasses`  
**File:** `internal/mcp/mcp_test.go` (preferred) or `cmd/trace/loop_test.go`

| Step | Action |
|------|--------|
| Setup | Temp dir, `trace init`, goal + task via MCP/CLI |
| Act | MCP `trace_plan` create-coarse + set-current + deep |
| Assert edit gate | `trace loop gate --for edit` → `allowed: true`, `reason_code` ≠ `plan_missing` |
| Assert status | `trace loop status` → `policy_inputs.plan_exists: true` |

### §6.2 Feet-seller export honesty + recovery

**Name:** `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`  
**File:** `internal/loop/gate_test.go` + helper in `cmd/trace/` or `internal/planner/bootstrap_test.go`

| Step | Action |
|------|--------|
| Setup | Import feet-seller export subset (`plan_changes: 11`, planner arrays empty) into temp store; seed one DONE task |
| Assert pre-terminal-fix (optional baseline) | Document current `plan_missing` on done — may guard with build tag or subtest `PrePhase36` |
| Assert post-terminal-fix | DONE task `--for done` → `allowed: true`, advisory `goal_plan_gap_terminal_advisory` (not actionable premature_implementation block) |
| Assert post-bootstrap | Run `plan bootstrap --goal`; edit gate → `allowed: true`, `plan_exists: true` |

**Fixture:** Anonymized JSON checked in under `internal/loop/testdata/feet-export-min.json` (derived from `/tmp/feet-export.json`) or generated in test from INVESTIGATION counts.

### §6.3 Active-work preservation

**Name:** `TestActiveWork_PlanMissingStillBlocksEdit`  
**File:** `internal/loop/gate_test.go` (extend `TestEvaluateGate_Edit_PlanMissing` pattern)

| Assert | Non-terminal task on goal without planner → edit gate `allowed: false`, `reason_code: plan_missing`, `recommended_phase: PLAN` |

### Additional S02 tests (recommended)

| Name | File | Covers |
|------|------|--------|
| `TestEvaluateGate_Done_TerminalPlanGapAdvisory` | `gate_test.go` | §2.5 |
| `TestPlanBootstrap_Idempotent` | `bootstrap_test.go` | §2.2 |
| `TestGoalStructureWarning_OverThresholdNoPlan` | `advisory_test.go` | §2.7 |
| `TestRegisteredToolNames_IncludesTracePlan` | `mcp_test.go` | §2.1 (16 tools) |

---

## §7 Explicit non-goals

| # | Non-goal | Rationale |
|---|----------|-----------|
| 1 | GUI-only patch | INTAKE: user requested root cause fix; GUI follows library |
| 2 | Global `plan_missing` weaken for active PLAN work | Would allow edit without progressive planner — violates DESIGN-LOCKS preserve |
| 3 | Deleting feet-seller history | 123 tasks, 11 plan-changes, 127 reviews — bootstrap recovers |
| 4 | Blaming Phase 35 pick logic | Loop 112 landing correct; gate meaning on terminal DONE is separate |
| 5 | Hosted SaaS / multi-tenant | DESIGN-LOCKS out of scope |
| 6 | Auto-LLM backlog generation | `cmd/trace/plan.go:67` — caller-supplied structure only |
| 7 | PlanExists bridge (§2.4) | Deferred — must not fake progressive plan |
| 8 | HTTP POST plan mutation routes | Defer unless OpenAPI contract phase demands; MCP is primary agent surface |
| 9 | MCP `trace_loop action=gate` | Out of scope — gate via CLI/hook/GUI HTTP sufficient for S02 |
| 10 | Default enforce strict | Opt-in model preserved |

---

## §8 S02 handoff

### Implementation order (strict)

```text
1. internal/planner/bootstrap.go + bootstrap_test.go
2. internal/planner/advisory.go (goal structure warning)
3. internal/loop/gate.go terminal honesty + gate_test.go
4. internal/mcp/tools_plan.go + server.go registration + mcp_test.go
5. cmd/trace/plan.go bootstrap subcommand + plan show advisory
6. internal/install/enforcement.go + install hint
7. internal/config/enforce.go + cmd/trace/init.go nudge
8. web GateStrip/TaskDetail advisory mapping (if library JSON stable)
9. TestLegacy + TestGreenfield + TestActiveWork
```

### File targets with accept/defer

| File | Action | Section |
|------|--------|---------|
| `internal/planner/bootstrap.go` | **Create** | §2.2 |
| `internal/planner/advisory.go` | **Create** | §2.7 |
| `internal/planner/bootstrap_test.go` | **Create** | §6.2 |
| `internal/loop/gate.go` | **Edit** | §2.5 |
| `internal/loop/gate_test.go` | **Edit** | §6.2, §6.3 |
| `internal/loop/testdata/feet-export-min.json` | **Create** | §6.2 fixture |
| `internal/mcp/tools_plan.go` | **Create** | §2.1 |
| `internal/mcp/server.go` | **Edit** | §2.1 register + `RegisteredToolNames` |
| `internal/mcp/mcp_test.go` | **Edit** | §6.1, 16-tool lock |
| `cmd/trace/plan.go` | **Edit** | §2.2 bootstrap + show advisory |
| `internal/install/enforcement.go` | **Edit** | §2.3 |
| `cmd/trace/install.go` | **Edit** | §2.3 hint |
| `internal/config/enforce.go` | **Edit** | §2.6 |
| `cmd/trace/init.go` | **Edit** | §2.6 |
| `web/src/components/GateStrip.tsx` | **Edit** | §2.5 GUI secondary |
| `web/src/screens/TaskDetail.tsx` | **Edit** (optional) | Advisory copy |
| `internal/httpapi/server.go` | **Defer** | No POST plan routes in S02 |

### Feet-seller recovery scope

- **S02:** Ship bootstrap; prove on temp import (§6.2).
- **S03:** Live bootstrap on fixture with user approval; GUI spot-check; export `trace/graph.json` if entity changes (CONTRIBUTING).

### Residual risks for S02 reviewer

| Risk | Mitigation |
|------|------------|
| Bootstrap heuristic too minimal for complex goals | Document refinement path; S03 human review on feet-seller |
| Terminal advisory mis-applied to active tasks | Gate tests §6.3 + explicit `work_state` check |
| MCP tool schema drift from CLI | Shared planner service calls; mirror `loop_test` JSON asserts |
| 16-tool registration break agents | Update locked `RegisteredToolNames` test list |
| GUI still red on advisory | GateStrip must treat `allowed: true` as warn even with violations[] |

### Recommended next row

**P36-S02-00** — scope planner: thicken implement prompt from this PLAN.md touch-list + acceptance tests; no product code in S02-00.

---

## Exit criteria checklist (P36-S01-01)

- [x] PLAN.md exists with locked fix choices
- [x] All 8 DESIGN-LOCKS candidates addressed in §2.1–§2.8
- [x] MCP plan / bootstrap / install first-class (§2.1–§2.3)
- [x] Touch-list ordered library → MCP → install → HTTP → GUI (§5)
- [x] Acceptance tests §6.1–§6.3 named with assert shapes
- [x] Explicit non-goals §7
- [x] Active-work PLAN enforcement preserved §4
- [x] No product code edited in S01-01
