# P37-S02-02 — Review implement

## Metadata
- id: P37-S02-02
- todo_ids: [P37-S02-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter, security-and-hardening]
- verification: automated
- hooks: []

## Objective

Independent review of P37-S02-01 deliverables vs [PLAN.md](../scope-01-plan/PLAN.md) + [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) + [RESIDUALS.md](../scope-00-triage/RESIDUALS.md). Re-run S02 acceptance tests + Phase 36 regression subset. Small inline fixes OK; spawn `02a`/`02b` for structural gaps. **No new features.** Fresh subagent — not S02-01 session.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-plan/PLAN.md) — §2 accept, §4 touch-list, §5 tests, §7 non-goals, §8 waves + risk notes
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [RESIDUALS.md](../scope-00-triage/RESIDUALS.md)
- [01-implement.md](01-implement.md)
- Implementer board Notes on P37-S02-01

## Session start

Fresh subagent. Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Confidence bar | **High** (or **Medium** with explicit residual risks listed) |
| Spawn threshold | Blocker/high without pending follow-up → spawn implement+review pair immediately below this row |
| Accept set | R1–R6, R8, R11 only — no scope creep into R7/R9/R8-full/R10 |
| R1 invariant | **Advisory only** — `bootstrap_recommended`; **NEVER** `PlanExists=true` |
| R7 invariant | `LoadEnforceMode` missing-config → `EnforceOff` unchanged |
| Law 19 | HTTP/MCP/GUI thin adapters — business logic in library |
| Live dogfood | Feet-seller fixture read-only in S02 — R9/R10 are S03 |

## Preflight / Plan

1. Read implementer Notes + changed files (`git diff` / file list)
2. Map each accepted R* to repo evidence
3. Run S02 acceptance tests + P36 regression subset
4. Spot-read OpenAPI + MCP schema updates for R2/R3/R5
5. Verify wave order A→D reflected in deliverables

## Checklist — accept set (PLAN §2)

### Wave A — R5 `advisories[]` + goal-structure warning — ACCEPT

- [ ] `StatusResult` in `internal/loop/apply.go` includes `advisories[]` field
- [ ] Code `goal_structure_warning` emitted when task count > 15 without `PlanExists`
- [ ] Wired from existing `GoalStructureWarning` (`internal/planner/advisory.go:13–41`)
- [ ] `violations[]` and `schema_version` unchanged — orthogonal channel
- [ ] Existing show/MCP stderr warning **not** removed (P36 partial preserved)

### Wave A — R1 bootstrap_recommended advisory — ACCEPT

- [ ] When `!PlanExists(goal)` and ≥**1** linked plan-change (`goalLinkedPlanChangeIDs`), `advisories[]` contains `bootstrap_recommended`
- [ ] Message recommends `trace plan bootstrap --goal <id>` or MCP `trace_plan action=bootstrap`
- [ ] **`TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists`** passes — `plan_exists` stays **false**
- [ ] **No** edits to `internal/loop/policy.go` — PlanExists from store read only
- [ ] Active-work `plan_missing` block **not** weakened

### Wave B — R3 MCP `trace_loop action=gate` — ACCEPT

- [ ] `internal/mcp/tools_loop.go` accepts `action=gate`
- [ ] JSON envelope matches CLI `trace loop gate` (`cmd/trace/loop.go:155–179`)
- [ ] Uses `loop.EvaluateGate` — **no** forked gate logic in MCP
- [ ] Blocked edit path includes violations in result
- [ ] `LoopInput` jsonschema updated; locked tool name list still **16** tools

### Wave B — R2 HTTP POST plan bootstrap — ACCEPT

- [ ] `POST /v1/plans/bootstrap` registered in `internal/httpapi/server.go`
- [ ] Handler in `handlers_p1.go` (or equivalent) calls `planner.Service.BootstrapFromPlanChanges`
- [ ] **No** business logic in handler beyond parse/validate/delegate
- [ ] `api/openapi.yaml` documents path + request/response schemas
- [ ] `TestHTTPPlanBootstrap_CreatesPlannerRows` passes

### Wave B — R4 bootstrap help refinement — ACCEPT

- [ ] `printPlanHelp` and/or bootstrap stderr states bootstrap yields **minimal** plan
- [ ] Mentions human refinement via `create-coarse` / `deep`
- [ ] No LLM generation claims; bootstrap limits honest
- [ ] `TestPlanHelp_MentionsRefinement` passes

### Wave C — R6 enforce nudge test — ACCEPT

- [ ] `TestWarnIfTraceDirWithoutConfig` in `internal/config/enforce_test.go`
- [ ] Temp `.trace/` without valid config → stderr nudge substring
- [ ] **`LoadEnforceMode` default unchanged** — missing config still `EnforceOff` (R7)

### Wave C — R11 critique workflow doc — ACCEPT

- [ ] Doc update cites post-bootstrap critique via **`trace loop apply`** + plan_changes
- [ ] References `TestGreenfield_MCPPlanBootstrap_EditGatePasses` as canonical
- [ ] **No** new MCP `critique-seed` tool added

### Wave D — R8 Overview surface — ACCEPT

- [ ] `web/src/screens/Overview.tsx` shows plan-gap/advisory copy when gate violation or status `advisories[]` present
- [ ] Consumes HTTP only (`GET /v1/loop/gate`, loop status) — Law 19
- [ ] `TaskDetail.tsx:205–211` bootstrap paragraph unchanged (unless copy-only with reason in Notes)
- [ ] Vitest or manual checklist evidence for advisory visibility

## Checklist — MCP / OpenAPI (R2, R3, R5)

### R5 + R1 — status advisories

- [ ] MCP `trace_loop action=status` returns `advisories[]` from library (not recomputed in MCP)
- [ ] OpenAPI loop status schema updated if status response is documented there

### R3 — MCP gate

- [ ] `trace_loop` tool description/schema lists `gate` in action enum
- [ ] Gate params align with CLI (`task_id`, `--for` equivalent if exposed)
- [ ] HTTP `GET /v1/loop/gate` unchanged as precedent — no duplicate logic added to HTTP layer

### R2 — HTTP bootstrap

- [ ] OpenAPI `POST /v1/plans/bootstrap` present with goal_id body and success response
- [ ] OpenAPI diff reviewable; no undocumented POST plan routes beyond accept set

### Tool catalog

- [ ] `RegisteredToolNames()` still returns **16** tools — no new tools
- [ ] `TestRegisteredToolNames_IncludesTracePlan` still passes

## Checklist — acceptance tests (PLAN §5)

Run:

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./internal/httpapi/... ./cmd/trace/... ./internal/config/... -count=1 \
  -run 'LoopStatus_.*Advisory|BootstrapAdvisoryNeverSetsPlanExists|MCPLoopGate|HTTPPlanBootstrap|PlanHelp_MentionsRefinement|WarnIfTraceDirWithoutConfig|Greenfield_MCPPlanBootstrap|FeetSellerExport_GateHonesty|ActiveWork_PlanMissing|TerminalPlanGapAdvisory|PlanBootstrap_Idempotent|GoalStructureWarning_OverThreshold|RegisteredToolNames'
```

| Test | Residual | Verify |
|------|----------|--------|
| `TestLoopStatus_IncludesGoalStructureAdvisory` | R5 | `goal_structure_warning` in JSON |
| `TestLoopStatus_BootstrapRecommendedAdvisory` | R1 | `bootstrap_recommended` in JSON |
| `TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists` | R1 guard | `plan_exists: false` |
| `TestMCPLoopGate_MatchesCLI` | R3 | Envelope parity |
| `TestHTTPPlanBootstrap_CreatesPlannerRows` | R2 | 200 + planner rows |
| `TestPlanHelp_MentionsRefinement` | R4 | create-coarse/deep mention |
| `TestWarnIfTraceDirWithoutConfig` | R6 | stderr nudge |
| Overview test or checklist | R8 | Banner visible |
| Block 0 regression | R11 | Greenfield MCP test green |

## Checklist — Phase 36 regression subset (PLAN §5 — must stay green)

- [ ] `TestGreenfield_MCPPlanBootstrap_EditGatePasses`
- [ ] `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`
- [ ] `TestActiveWork_PlanMissingStillBlocksEdit`
- [ ] `TestEvaluateGate_Done_TerminalPlanGapAdvisory`
- [ ] `TestPlanBootstrap_Idempotent`
- [ ] `TestGoalStructureWarning_OverThresholdNoPlan`
- [ ] `TestRegisteredToolNames_IncludesTracePlan`

## Checklist — preserve invariants (PLAN §7 non-goals)

- [ ] **No** silent PlanExists bridge in `policy.go` or deliberation
- [ ] **No** weaken active-work `plan_missing` edit block
- [ ] **No** default enforce flip (R7) — `EnforceOff` when config missing
- [ ] **No** new MCP critique-seed tool
- [ ] **No** business-logic fork in `web/` or HTTP handlers
- [ ] **No** full plan tree GUI (R8-full)
- [ ] **No** feet-seller task history rewrite
- [ ] Phase 36 MCP plan/bootstrap/terminal advisory guarantees intact

## Checklist — Law 19 boundaries

- [ ] Advisory assembly in `internal/planner/` + `internal/loop/` — not in MCP/HTTP/GUI
- [ ] HTTP bootstrap handler: parse → service call → JSON — no planner heuristic inline
- [ ] MCP gate: delegate to `loop.EvaluateGate` — same as CLI
- [ ] Overview: maps API JSON to copy/banners — no independent gate evaluation

## Checklist — touch-list completeness (PLAN §4)

| File | Expected | Residual |
|------|----------|----------|
| `internal/planner/advisory.go` | Edited — bootstrap helper | R1 |
| `internal/planner/advisory_test.go` | Edited | R1, R5 |
| `internal/loop/apply.go` | Edited — `advisories[]` | R5, R1 |
| `internal/loop/apply_test.go` or `cmd/trace/loop_test.go` | Edited | R1, R5 |
| `internal/mcp/tools_loop.go` | Edited — gate action | R3 |
| `internal/mcp/mcp_test.go` | Edited | R3 |
| `internal/httpapi/handlers_p1.go` | Edited — POST bootstrap | R2 |
| `internal/httpapi/server.go` | Edited — route | R2 |
| `api/openapi.yaml` | Edited | R2 |
| `internal/httpapi/*_test.go` | Edited/Created | R2 |
| `cmd/trace/plan.go` | Edited — help | R4 |
| `cmd/trace/plan_test.go` or `help_test.go` | Edited | R4 |
| `internal/config/enforce_test.go` | Edited | R6 |
| docs cross-ref | Edited | R11 |
| `web/src/screens/Overview.tsx` | Edited | R8 |
| `internal/loop/policy.go` | **Not** mutated | — |
| `LoadEnforceMode` default | **Not** changed | R7 |

## Risk review (PLAN §8)

| Risk | Check |
|------|-------|
| R1 accidentally sets PlanExists | Guard test present + `policy.go` diff clean |
| `advisories[]` breaks status consumers | `schema_version` + `violations[]` unchanged |
| HTTP handler duplicates CLI | Single `BootstrapFromPlanChanges` call path |
| Overview duplicates TaskDetail | Uses API advisory codes only; TaskDetail unchanged |

## Reviewer loop

```text
1. Compare deliverables to checklists above
2. Classify findings: blocker | high | medium | low | nit
3. blocker/high: small fix OR spawn 02a implement + 02b review immediately below this row
4. Re-run tests after fixes
5. UNTIL no open blocker/high without pending follow-up
   AND confidence medium or high with evidence
```

## Todo updates

Set own row `done` + Notes: PASS/FAIL, confidence, findings count, spawn IDs if any, test command output summary.

## Exit criteria

- [ ] All 8 accept residuals (R1–R6, R8, R11) evidenced in diff
- [ ] S02 acceptance tests verified passing
- [ ] Phase 36 regression subset green
- [ ] R1 guard confirmed; R7 default unchanged
- [ ] OpenAPI + MCP schema updates verified for R2/R3/R5
- [ ] Law 19 — no business-logic fork
- [ ] No open blocker/high without spawned follow-up row
- [ ] Confidence **high** (or **medium** with residual risks in Notes)
- [ ] Next: **P37-S03-00** on PASS

## Next

`P37-S03-00` (on PASS)
