# Phase 37 S01 — PLAN

**Author:** P37-S01-01 (2026-08-22)  
**SoT inputs:** [RESIDUALS.md](../scope-00-triage/RESIDUALS.md) (S00-02 APPROVE), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [01-plan.md](./01-plan.md) (S01-00 locks)  
**Dogfood fixture (read-only in S02/S03):** `/home/ali/Desktop/feet seller telegram app`

---

## §1 Summary

| Metric | Value |
|--------|------:|
| **Accept (S02)** | 8 — R1–R6, R8, R11 |
| **Accept (S03)** | 1 — R10 live GUI verify |
| **Re-defer** | 2 — R7, R9 (+ R8-full plan screen in §3 registry) |
| **Reject** | Silent PlanExists bridge; MCP critique-seed tool |
| **Theme** | Close Phase 36 residuals — advisories channel, adapter parity, honest help, minimal GUI surface |

**Partial P36 S02 already shipped (live re-verified 2026-08-22):** R4 bootstrap CLI/MCP exists (help gap only); R5 `GoalStructureWarning` helper + show/MCP stderr (status channel gap); R6 `WarnIfTraceDirWithoutConfig` helper (unit test gap); R8 TaskDetail bootstrap paragraph (`TaskDetail.tsx:205–211`); R11 Block 0 MCP greenfield path (`mcp_test.go:1210+`).

**Full gaps (S02 implements):** R1 advisory bridge; R2 HTTP POST bootstrap; R3 MCP `trace_loop action=gate`; R5 `StatusResult.advisories[]`; R6 dedicated test; R8 Overview surface; R4 help refinement note; R11 agent-workflow doc.

---

## §2 Accepted residuals (S02 — lock behavior)

| ID | Item | Wave | Locked behavior | Must not |
|----|------|------|-----------------|----------|
| **R1** | PlanExists advisory bridge | A | When `!PlanExists(goal)` **and** `len(goalLinkedPlanChanges) ≥ 1` (N=**1**), emit `advisories[]` entry `bootstrap_recommended` on `trace.loop.status.v1`. Message recommends `trace plan bootstrap --goal <id>` or MCP `trace_plan action=bootstrap`. Reuse `goalLinkedPlanChangeIDs` (`internal/planner/advisory.go:53–82`). | Set `PlanExists`, mutate `policy.go` deliberation inputs, write planner rows, weaken active `plan_missing` |
| **R5** | `loop status advisories[]` | A | Add `advisories[]` to `StatusResult` (`internal/loop/apply.go:196–209`). Wire existing `GoalStructureWarning` → code `goal_structure_warning` when task count > 15 without plan (`advisory.go:13–41`). Orthogonal to R1 — both may appear. | Change gate deny semantics; remove show/MCP stderr warning |
| **R3** | MCP `trace_loop action=gate` | B | Extend `trace_loop` actions to `next\|apply\|status\|gate` (`tools_loop.go:39–47`). Mirror CLI `trace loop gate` JSON envelope via `loop.EvaluateGate` (`cmd/trace/loop.go:155–162`). HTTP `GET /v1/loop/gate` precedent only — do not duplicate. | Fork gate logic in MCP; change exit semantics vs CLI |
| **R2** | HTTP POST plan routes | B | **Minimum:** `POST /v1/plans/bootstrap` thin handler → `planner.Service.BootstrapFromPlanChanges` (mirror `cmd/trace/plan.go:351`). Register in `server.go`; document in `api/openapi.yaml`. Optional expansion: POST create-coarse/deep if S02-00 thickens. | Business logic in handlers; skip OpenAPI |
| **R4** | Bootstrap help refinement note | B (install layer) | `printPlanHelp` / bootstrap stderr: state bootstrap yields **minimal** plan; human refinement via `create-coarse` / `deep` expected (PLAN §2.2 honesty). | LLM generation claims; hide bootstrap limits |
| **R6** | `WarnIfTraceDirWithoutConfig` unit test | C (install layer) | `TestWarnIfTraceDirWithoutConfig` in `internal/config/enforce_test.go`: temp dir with `.trace/` no valid config → stderr contains nudge substring (`enforce.go:43–57`). | Flip default enforce to `warn` (R7) |
| **R8** | Overview plan-gap / advisories (minimal) | D | `web/src/screens/Overview.tsx`: show bootstrap/advisory copy when gate violation or status `advisories[]` present. Consume HTTP only (`GET /v1/loop/gate`, loop status API) — Law 19. TaskDetail bootstrap paragraph unchanged (`TaskDetail.tsx:205–211`). | Planner logic in `web/`; full plan tree screen |
| **R11** | CLI greenfield critique discoverability | C (install layer) | Document agent workflow: after bootstrap, seed critique via **`trace loop apply`** with plan_changes envelope (Block 0 pattern in `mcp_test.go:1210+`). Update verify/agent-loop docs — **not** new MCP tool. | New `trace_plan action=critique-seed`; duplicate Block 0 in code |

### S03 accept (not S02 code)

| ID | Item | Locked behavior |
|----|------|-----------------|
| **R10** | Live GUI browser verify | Pin screenshot or `docs/verification/phase-37-p36-residuals/` note: TaskDetail advisory + Overview surface after S02. Use cursor-ide-browser MCP. |

---

## §3 Re-deferred residuals

| ID | Item | Owner | Trigger | Notes |
|----|------|-------|---------|-------|
| **R7** | Enforce default `warn` when `.trace/` without config | Human / product | Explicit decision to change `LoadEnforceMode` missing-config path (`enforce.go:25–29`); until then **`EnforceOff`** | P36 stderr nudge + R6 test remain; R6 does **not** imply default flip |
| **R9** | Feet-seller deep refinement quality | Human dogfood | S03: document post-bootstrap `create-coarse`/`deep` path on fixture; quality spot-check | Fixture read-only; `plan_uncritiqued` post-bootstrap expected; no task history rewrite |
| **R8-full** | Full plan screen / plan tree GUI | Phase 38 planner (or human) | S03 VERIFY shows Overview minimal surface insufficient | S02 delivers Overview copy only; `GET /v1/plans` read-only |

### Permanent rejects (document in PLAN)

| Pattern | Why |
|---------|-----|
| Silent `PlanExists=true` from plan-changes alone | Violates DESIGN-LOCKS + P36 guarantee |
| MCP `critique-seed` tool | Duplicates `trace loop apply`; R11 doc path preferred |
| Business-logic fork in `web/` | Law 19 |
| Feet-seller task history rewrite | INTAKE out of scope |

---

## §4 Touch-list (ordered — library → MCP → HTTP → install → GUI)

Policy in library first; adapters thin-wrap shared services. **S02 implements in this order.**

| Order | Layer | Paths | Residuals | Action |
|------:|-------|-------|-----------|--------|
| **1** | Library — loop status | `internal/loop/apply.go`, `internal/loop/apply_test.go` (or `cmd/trace/loop_test.go`) | R5, R1 | Add `Advisories []Advisory` to `StatusResult`; assemble in `Status` / `attachStatusViolations` path calling `planner.GoalStructureWarning` + new bootstrap advisory helper |
| **2** | Library — planner advisory | `internal/planner/advisory.go`, `internal/planner/advisory_test.go` | R1, R5 | Add bootstrap advisory builder (reuse `goalLinkedPlanChangeIDs`); extend tests for both codes |
| **3** | Library — policy (cite only) | `internal/loop/policy.go` | — | **No change** — `PlanExists` from store read only (`:45–48`, `:91`) |
| **4** | MCP adapter | `internal/mcp/tools_loop.go`, `internal/mcp/mcp_test.go` | R3 | Add `gate` action + helper mirroring CLI; extend locked tool schema tests |
| **5** | HTTP adapter | `internal/httpapi/handlers_p1.go`, `internal/httpapi/server.go`, `api/openapi.yaml`, `internal/httpapi/handlers_p1_test.go` (or integration test file) | R2 | `POST /v1/plans/bootstrap` handler; route registration; OpenAPI path + schema |
| **6** | Install — CLI help | `cmd/trace/plan.go`, `cmd/trace/plan_test.go` or `help_test.go` | R4 | Refinement sentence in `printPlanHelp` + bootstrap stderr |
| **7** | Install — config test | `internal/config/enforce_test.go` | R6 | `TestWarnIfTraceDirWithoutConfig` |
| **8** | Install — agent docs | `docs/rules/agent-loop-protocol.md` and/or `docs/phases/phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md` cross-ref, verify script pointer under `experiments/` or `docs/verification/` | R11 | Document critique = loop apply with plan_changes; cite `TestGreenfield_MCPPlanBootstrap_EditGatePasses` |
| **9** | GUI adapter | `web/src/screens/Overview.tsx`, optional `web/src/components/GateStrip.tsx` (reuse warn path), Vitest if present | R8 | Minimal advisory/plan-gap banner when gate or status advisories present |

**Explicit non-touch:** `internal/loop/policy.go` bridge heuristic; default enforce flip; hosted SaaS routes; silent PlanExists.

---

## §5 Acceptance tests (S02 implements)

| ID | Test name / method | File | Assert shape |
|----|-------------------|------|--------------|
| **R5** | `TestLoopStatus_IncludesGoalStructureAdvisory` (or extend `TestLoopStatus*`) | `cmd/trace/loop_test.go` or `internal/loop/apply_test.go` | >15 tasks, no plan → status JSON `advisories[]` contains `goal_structure_warning` |
| **R1** | `TestLoopStatus_BootstrapRecommendedAdvisory` | same | Goal with ≥1 linked plan-change, no plan → `advisories[]` contains `bootstrap_recommended`; `deliberation.policy_inputs.plan_exists` still **false** |
| **R1 guard** | `TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists` | same | Same setup — must **not** flip `plan_exists` true |
| **R3** | `TestMCPLoopGate_MatchesCLI` | `internal/mcp/mcp_test.go` | `trace_loop action=gate task_id=…` envelope matches CLI; blocked edit → violations present |
| **R2** | `TestHTTPPlanBootstrap_CreatesPlannerRows` | `internal/httpapi/*_test.go` | `POST /v1/plans/bootstrap` + goal_id → 200; planner rows exist; OpenAPI documents route |
| **R4** | `TestPlanHelp_MentionsRefinement` | `cmd/trace/plan_test.go` or `help_test.go` | Help/bootstrap text contains create-coarse/deep refinement sentence |
| **R6** | `TestWarnIfTraceDirWithoutConfig` | `internal/config/enforce_test.go` | `.trace/` without config → stderr nudge substring |
| **R8** | `Overview.test.tsx` or manual checklist | `web/src/screens/Overview.test.tsx` | Advisory/plan-gap copy visible when fixture gate/status has advisory |
| **R11** | Doc review + Block 0 regression | — | PLAN cites workflow; `TestGreenfield_MCPPlanBootstrap_EditGatePasses` remains canonical |

### Phase 36 regression subset (S03 Block 0 — must stay green)

Include in VERIFY Block 0:

`TestGreenfield_MCPPlanBootstrap_EditGatePasses`, `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`, `TestActiveWork_PlanMissingStillBlocksEdit`, `TestEvaluateGate_Done_TerminalPlanGapAdvisory`, `TestPlanBootstrap_Idempotent`, `TestGoalStructureWarning_OverThresholdNoPlan`, `TestRegisteredToolNames_IncludesTracePlan`.

---

## §6 VERIFY block mapping (S03-01)

| Block | Content | Residuals / evidence |
|-------|---------|---------------------|
| **0** | Phase 36 acceptance subset still green | All P36 tests above — `go test` log in `experiments/runs/` |
| **1** | Per accepted residual — test name or JSON capture | R1, R2, R3, R4, R5, R6, R8 (API/Vitest), R11 (doc pointer) |
| **2** | Feet-seller spot-check | R8 Overview + R9 documented refinement path (read-only fixture) |
| **3** | Greenfield MCP path | R3 gate action; R11 workflow aligns with Block 0 MCP test |
| **4** | Re-defer registry + R10 browser | R7, R9, R8-full updated in VERIFY-NOTES; R10 screenshot/`docs/verification/` |
| **5** | DR-HANDOFF successor table | Phase closure + any Phase 38 stubs |

---

## §7 Non-goals

| # | Non-goal |
|---|----------|
| 1 | Silent PlanExists bridge — fake progressive plan |
| 2 | Weaken active-work `plan_missing` block |
| 3 | Hosted SaaS / multi-tenant / always-on daemon |
| 4 | Full plan tree GUI (R8-full → Phase 38) |
| 5 | Feet-seller task history rewrite |
| 6 | Default enforce `warn` without product lock (R7) |
| 7 | New MCP critique-seed tool (R11 reject) |
| 8 | Business-logic fork in `web/` (Law 19) |
| 9 | Phase 36 core reopen (MCP plan, bootstrap, terminal advisory) |

---

## §8 S02 handoff (wave + strict file order)

### Dependency waves (from RESIDUALS §5)

```text
Wave A — status schema (blocks R8 data contract):
  R5  advisories[] + GoalStructureWarning wire
  R1  bootstrap_recommended on same channel

Wave B — adapters (after A):
  R3  MCP trace_loop action=gate
  R2  HTTP POST /v1/plans/bootstrap (+ OpenAPI)
  R4  bootstrap help refinement note

Wave C — tests + docs:
  R6  WarnIfTraceDirWithoutConfig unit test
  R11 agent-workflow / verify doc path

Wave D — GUI (after A):
  R8  Overview minimal plan-gap / advisories surface
```

### S02 implement order (touch-list aligned)

```text
1. internal/planner/advisory.go (+ tests) — bootstrap advisory helper
2. internal/loop/apply.go (+ tests) — StatusResult.advisories[]
3. internal/mcp/tools_loop.go (+ mcp_test.go) — gate action
4. internal/httpapi/handlers_p1.go, server.go, openapi.yaml (+ tests) — POST bootstrap
5. cmd/trace/plan.go (+ help test) — R4 refinement note
6. internal/config/enforce_test.go — R6
7. docs verify/agent-loop cross-ref — R11
8. web/src/screens/Overview.tsx (+ test) — R8
9. Full test pass + export graph if entities change (CONTRIBUTING)
```

### Risk notes for S02 reviewer

| Risk | Mitigation |
|------|------------|
| R1 accidentally sets PlanExists | Dedicated guard test; no `policy.go` edits |
| advisories[] breaks status schema consumers | Keep `schema_version`; violations[] unchanged |
| HTTP handler duplicates CLI | Single `planner.Service` call path |
| Overview duplicates TaskDetail logic | Shared advisory codes from API JSON only |

---

## Next

**P37-S02-00** — scope planner thickens implement prompts from this PLAN.
