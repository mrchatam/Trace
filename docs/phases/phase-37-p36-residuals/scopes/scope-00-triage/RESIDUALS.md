# RESIDUALS — Phase 37 triage (R1–R11)

**Author:** P37-S00-01 · **Verified:** 2026-08-22 against live repo at `/home/ali/Desktop/Trace`.

---

## 1. Summary

| Decision | Count | IDs |
|----------|------:|-----|
| **Accept** | 9 | R1–R6, R8, R10, R11 |
| **Defer** | 2 | R7, R9 |
| **Reject** | 2 patterns | Silent `PlanExists` bridge (R1 alternative); new MCP `critique-seed` tool (R11 alternative) |

**Scope split:** 8 accepts land in **S02** (R1–R6, R8, R11). **R10** accepts in **S03** (live browser verify). **R9** defers to **S03** dogfood (documented refinement path; no mandatory history rewrite).

### Partial P36 S02 shipment (live re-verified 2026-08-22)

| ID | Shipped in P36 S02 | Still open (P37) |
|----|-------------------|------------------|
| R4 | `trace plan bootstrap` + MCP `trace_plan bootstrap` (`cmd/trace/plan.go:62–63`, `tools_plan.go:164+`) | Help omits explicit human-refinement note (PLAN §2.2; `plan.go:51–72` lists subcommands only) |
| R5 | `GoalStructureWarning` + unit test; wired to `plan show` stderr + MCP show field (`advisory.go:13–41`, `plan.go:297–298`, `tools_plan.go:144–156`) | `StatusResult` has no `advisories[]` (`apply.go:196–209`) |
| R6 | `WarnIfTraceDirWithoutConfig` + `init` call (`enforce.go:43–57`, `cmd/trace/init.go:50`) | No dedicated unit test (`enforce_test.go` covers `LoadEnforceMode` only) |
| R7 | Stderr nudge helper (P36 §2.6 document-only accept) | Default enforce still `EnforceOff` when config missing (`enforce.go:25–29`) |
| R8 | TaskDetail bootstrap advisory paragraph (`TaskDetail.tsx:205–211`) | Overview GateStrip only — no plan-gap copy (`Overview.tsx:72–88`) |
| R1, R2, R3 | — | Full gap (see table §2) |
| R9–R11 | MCP greenfield path (R11 Block 0) | Feet refinement quality (R9); live GUI (R10); CLI verify critique doc gap (R11) |

**Live spot-checks (no overrides to S00-00 locks):**

- `StatusResult` — no `advisories` field (`internal/loop/apply.go:196–209`)
- `trace_loop` — actions `next|apply|status` only (`tools_loop.go:39–47`); no `gate`
- HTTP plans — `GET /v1/plans` only (`server.go:282`, `handlers_p1.go:115`; `openapi.yaml:760` GET only)
- `policy.go` — `PlanExists` from store read only (`policy.go:45–48`, `:91`); no plan-change-count advisory
- `goalLinkedPlanChangeIDs` exists in `advisory.go:53–82` (used by bootstrap, reusable for R1)

---

## 2. Triage table

| ID | Item | P36 source | Decision | Effort | Risk | Rationale | S02? | Test sketch | Live cite |
|----|------|------------|----------|--------|------|-----------|------|-------------|-----------|
| R1 | PlanExists advisory bridge — recommend bootstrap when goal-linked plan-changes exist and `!PlanExists`; **never** silent satisfy | PLAN §2.4 defer; VERIFY-NOTES § Residuals #2 | **accept** | M | low | Feet-seller had 11 plan-changes with 0 planner rows; operators need discoverable nudge without faking progressive plan or weakening `plan_missing` on active work | Y (Wave A) | Unit: goal with ≥1 linked plan-change, no plan → `trace loop status` JSON includes `advisories[]` entry `bootstrap_recommended`; must not flip `deliberation.policy_inputs.plan_exists` | `policy.go:45–48`; `advisory.go:53–82`; no bridge in `apply.go:196–209` |
| R2 | HTTP POST plan routes — GUI/API parity for bootstrap (min); create-coarse/deep if S01 expands | PLAN touch-list defer; VERIFY-NOTES § Residuals #2 | **accept** | M | low–med | Law 19: thin handlers over `planner.Service`; mirror CLI. Enables browser bootstrap without shell | Y (Wave B) | Integration: `POST /v1/plans/bootstrap` (or agreed path) with goal_id → 200 + planner rows; OpenAPI documents route | `server.go:282` GET only; `handlers_p1.go:115`; `openapi.yaml:760` |
| R3 | MCP `trace_loop action=gate` — gate check without shell | defer in VERIFY-NOTES #2 | **accept** | S | low | Parity with CLI `loop gate` and existing HTTP `GET /v1/loop/gate`; library call already canonical | Y (Wave B) | MCP test: `trace_loop action=gate task_id=…` returns same envelope shape as CLI; blocked edit → violations present | `tools_loop.go:39–47`; CLI `loop.go:122–176`; HTTP `handlers_loop.go:104` |
| R4 | Bootstrap help — human-refinement note per PLAN §2.2 | S02-02 low; VERIFY-NOTES #1 | **accept** | S | low | Bootstrap recovers minimal plan; operators should expect `create-coarse`/`deep` for quality — honest help only | Y (Wave B) | `go test ./cmd/trace -run TestHelp` or snapshot help text contains refinement sentence | `plan.go:51–72` subcommand list; bootstrap stderr `:350` lacks refinement note |
| R5 | `loop status advisories[]` — wire `GoalStructureWarning` into status JSON | S02-02 low / PLAN §2.7 partial; VERIFY-NOTES #1 | **accept** | S | low | Helper + tests shipped; gap is status channel only. Blocks R1/R8 data contract | Y (Wave A) | Extend `TestLoopStatus*` or loop package test: >15 tasks, no plan → `advisories[]` contains `goal_structure_warning` code | `advisory.go:13–41`; `apply.go:196–209` no field |
| R6 | `WarnIfTraceDirWithoutConfig` unit test | nit; VERIFY-NOTES #1 | **accept** | S | low | Helper shipped; test locks stderr nudge when `.trace/` exists without valid config | Y (Wave C) | `TestWarnIfTraceDirWithoutConfig`: temp dir with `.trace/` no config → stderr contains nudge substring | `enforce.go:43–57`; `enforce_test.go` — no Warn test |
| R7 | Enforce default `warn` when `.trace/` without config | PLAN §2.6 defer | **defer** | — | med | Product decision: P36 explicitly preserved `EnforceOff` fail-open; stderr nudge sufficient until human locks flip | N | N/A — re-defer §6 | `enforce.go:25–29` → `EnforceOff`; nudge `:43–57` |
| R8 | Goal/plan surface UX beyond TaskDetail — Overview minimal plan-gap / advisories | defer in DR-HANDOFF | **accept** (minimal) | M | low | Law 19: consume API (`GET /v1/loop/gate` violations or status `advisories[]`); no planner logic in `web/`. Full plan tree screen → §6 re-defer | Y (Wave D) | Vitest or manual: Overview shows bootstrap/advisory copy when gate violation or advisory present; TaskDetail unchanged | `TaskDetail.tsx:205–211`; `Overview.tsx:72–88` GateStrip only |
| R9 | Feet-seller planner quality — post-bootstrap refinement via create-coarse/deep | VERIFY Block 6 note | **defer** → S03 | S | low | Bootstrap path works; deep refinement is human dogfood + doc, not S02 code. No history rewrite | N (S03) | S03: document path on fixture; optional quality spot-check | VERIFY-NOTES Block 6; `plan_uncritiqued` post-bootstrap expected |
| R10 | Live GUI browser verify — terminal + plan surfaces | VERIFY Block 4 deferred | **accept** → S03 | S | low | Pre-bootstrap API + GateStrip sufficient for P36; cheap browser spot-check closes verify gap | N (S03) | Pin screenshot or `docs/verification/` note: TaskDetail advisory + Overview surface after S02 | VERIFY-NOTES Block 4–5 |
| R11 | CLI greenfield critique path — agent discoverability | VERIFY Block 1 partial | **accept** (doc) | S | low | Block 0 MCP test seeds critique via `loop apply`; gap is locked CLI verify script + workflow docs — **not** new MCP tool | Y (Wave C) | Doc review + optional script pointer; existing `TestGreenfield_MCPPlanBootstrap_EditGatePasses` remains canonical | VERIFY-NOTES #4; `mcp_test.go` Block 0 pattern |

### Rejects (documented alternatives)

| Pattern | Related ID | Why rejected |
|---------|------------|--------------|
| Silent `PlanExists=true` from plan-changes alone | R1 | Violates DESIGN-LOCKS + Phase 36 guarantee; would fake progressive plan and weaken honesty |
| New MCP `trace_plan action=critique-seed` (or similar) | R11 | Duplicates existing `trace loop apply` with plan_changes envelope; prefer doc/workflow path |
| Business-logic fork in `web/` | R8 | Law 19 — adapters consume library/API only |
| Feet-seller task history rewrite | R9 | Out of scope per INTAKE |

---

## 3. R1 — PlanExists advisory bridge (locked signal)

**Decision:** Accept advisory-only. **Reject** silent satisfy (`PlanExists=true` without real planner rows).

| Field | Lock |
|-------|------|
| Trigger | `!PlanExists(goal)` **and** `len(goalLinkedPlanChanges) ≥ 1` |
| Threshold N | **1** (any goal-linked plan-change without progressive plan warrants nudge) |
| Counting | Reuse `goalLinkedPlanChangeIDs` in `internal/planner/advisory.go:53–82` (same path as `bootstrap.go:107`) |
| Surface | `trace.loop.status.v1` → `advisories[]` entry |
| Advisory code | `bootstrap_recommended` (stable snake_case) |
| Message shape | Recommend `trace plan bootstrap --goal <id>` or MCP `trace_plan action=bootstrap` |
| Must not | Set `PlanExists`, write planner rows, or change `policy.go` deliberation inputs (`policy.go:45–48`, `:91`) |

**Orthogonal advisory:** R5 `goal_structure_warning` fires on **task count > 15** without plan (`GoalStructureWarningThreshold` in `advisory.go:13–14`). Both may appear in `advisories[]` simultaneously.

**Adapter target (S02):** extend `loop.Status` / `attachStatusViolations` path in `internal/loop/apply.go` — add advisory assembly calling planner helpers; no changes to gate deny semantics for active `plan_missing`.

---

## 4. R7 — Enforce default `warn` (re-defer)

P36 §2.6 shipped **document-only** acceptance: preserve `LoadEnforceMode` → `EnforceOff` when `.trace/config.json` is missing or invalid (`internal/config/enforce.go:25–29`). Stderr nudge via `WarnIfTraceDirWithoutConfig` (`:43–57`) runs from `trace init` (`cmd/trace/init.go:50`).

**P37 decision:** Re-defer flipping the default to `warn`. Operators opt in via `trace install` / explicit config. Changing fail-open default affects all harness projects without config — requires explicit product decision (§6).

---

## 5. S02 wave order (dependency order)

```text
Wave A — status schema (blocks R1, R8 data):
  R5  advisories[] field on StatusResult + GoalStructureWarning wire
  R1  bootstrap_recommended advisory (same channel)

Wave B — adapters (parallel after A):
  R3  MCP trace_loop action=gate → loop.EvaluateGate (mirror cmd/trace/loop.go:155–162)
  R2  HTTP POST plan routes (bootstrap minimum; create-coarse/deep if S01 expands)
      → handlers in internal/httpapi/handlers_p1.go; OpenAPI api/openapi.yaml
  R4  bootstrap help text (cmd/trace/plan.go)

Wave C — tests + docs:
  R6  WarnIfTraceDirWithoutConfig unit test (internal/config/enforce_test.go)
  R11 agent-workflow / verify-script doc — critique = loop apply with plan_changes

Wave D — GUI (after A):
  R8  Overview plan-gap / advisories surface (web/src/screens/Overview.tsx)

Out of S02:
  R7  defer (product default flip)
  R9, R10 → S03 verify
```

### Law 19 adapter targets (R2 / R3)

| Residual | Adapter file | Library call |
|----------|--------------|--------------|
| R2 POST bootstrap (min) | `internal/httpapi/handlers_p1.go` (new handler) | `planner.Service.BootstrapFromPlanChanges` — mirror `cmd/trace/plan.go:351` |
| R2 POST create-coarse (if S01 expands) | same pattern | `planner.Service.CreateCoarsePlan` — mirror `cmd/trace/plan.go:88` |
| R3 MCP gate | `internal/mcp/tools_loop.go` switch + helper | `loop.EvaluateGate` — mirror `cmd/trace/loop.go:155–162` |
| R3 HTTP (exists) | `internal/httpapi/handlers_loop.go:104` | precedent only — do not duplicate |

---

## 6. Re-defer registry

| ID | Item | Owner | Trigger | Notes |
|----|------|-------|---------|-------|
| R7 | Enforce default `warn` when config missing | Human / product | Explicit S01+ decision to change `LoadEnforceMode` missing-config behavior; until then default **off** | P36 stderr nudge remains; R6 test does not imply default flip |
| R8 | Full plan screen / plan tree GUI | Phase 38 planner (or human) | S03 VERIFY shows Overview minimal surface insufficient for operator needs | S02 delivers Overview copy/advisories only; `GET /v1/plans` read-only today |
| R9 | Feet-seller deep refinement quality | Human dogfood | Post-bootstrap `create-coarse`/`deep` path documented in S03; live quality acceptable for verify | Fixture read-only; no task history rewrite |
| — | Silent PlanExists bridge | — | **Rejected** permanently for Phase 37 | See §3 |
| — | MCP critique-seed tool | — | **Rejected** — doc path per R11 | Extend agent-loop / verify docs |

---

## Next

**P37-S00-02** — independent review against DESIGN-LOCKS, INTAKE, and live cites.
