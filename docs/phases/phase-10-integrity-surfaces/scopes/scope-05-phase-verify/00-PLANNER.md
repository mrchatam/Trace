# P10 / S05 / 00-PLANNER — Phase 10 VERIFY / integrity closeout

## Metadata
- id: P10-S05-00
- todo_ids: [P10-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 10 VERIFY evidence table: **S01–S04 named DF regressions** + **carry-forward gates** + `./...` (product pkgs). Decide **DR-HANDOFF** = **`no successor`** unless Notes promote. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Pattern: Phase 09 VERIFY [`../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/`](../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/)
- Sibling REVIEW-NOTES: S01–S04 under `../scope-0{1,2,3,4}-*/REVIEW-NOTES.md`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Depends-on (S01–S04 — landed)

| Scope | Board | DF locks imported |
|-------|-------|-------------------|
| S01 | P10-S01-02 APPROVE high | DF-19/23/25/27/29 named tests (supersede GC-01 global DPC names) |
| S02 | P10-S02-02 APPROVE high | DF-21/22/32; **nine** MCP tools; G19; BuiltinMCP×9 |
| S03 | P10-S03-02 APPROVE high | DF-20 rename GC + argv isolation + `TestIndexIncrementalIsolation` |
| S04 | P10-S04-02 APPROVE high | DF-17/18/24/26/31; honesty Path C operator-flag; Gate G hatch retained |

## Live residuals → DR-HANDOFF decision (2026-08-16)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 10 | DF-17…32 (deduped) → S01–S04 | Closed by S01–S04 APPROVE high |
| Explicit residual OK into VERIFY | `plan_scope` Exact still out; Mode-B packs historical | Forward notes only — **not** a successor phase |
| Parallel dogfood (not board-blocking) | ab-index / ab-operator-gate / Cursor MCP reload | Stay in `experiments/` — **not** Phase 11 unless promoted |
| Known `./...` nit | `similar projects/graphify` space-in-path FAIL | Pre-existing non-product — VERIFY records product pkgs PASS |

**DR-HANDOFF = `no successor`.** Reopen only with explicit promotion + scaffold (same posture as Phase 08/09 historical closes).

## Planner work
1. Lock VERIFY command set (S01–S04 named tests + carry-forward + product `./...`).
2. Thicken `01-verify.md` evidence table + spawn 01a/b/c + handoff **start**.
3. Thicken `02-scope-review.md` owns DR-HANDOFF **completion** (`no successor`).
4. SCOPE-TODOS + board sync; light stamp `DR-HANDOFF.md` ownership.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Phase gate | Phase 10 integrity-surfaces closeout — **not** a new planted eval gate |
| S01 DF-19 | `TestWhyTaskDPCGoalScoped` + `TestWhyTaskDPCMultiGoalNoForeignPollution` (+ compiler `TestTaskContextDPCGoalScoped` / `TestTaskContextMultiGoalOmitsForeignDPC`) |
| S01 DF-23/25/27/29 | `TestExactWhyPlanChangeAlias` / `TestExactWhyCapability` / `TestDecisionMarkdownTrustLabels` / `TestIncludeWhyFailClosed` |
| S02 MCP | **Nine** tools: prior six + `trace_tasks` + `trace_capability` + `trace_version`; `TestToolNamesRegistered` + `TestBuiltinMCPCapabilitySpecs` + `TestImportBoundaryMCPNoPlanImpactIndexTools` |
| S02 DF-21/22/32 | `TestTraceTasksParity` / `TestTraceCapabilityActions` / `TestTraceVersion` / install tip asserts / `TestCapabilityListMissingSnakeCase` |
| S03 DF-20 | `TestIndexGCAfterPathRename` + `TestIndexArgvMissingPathDeletesOnlyThatPath` + `TestIndexIncrementalIsolation` (+ store `TestListFilePathsAndDeleteFileByPath`) |
| S04 DF-17/18/24/26/31 | `TestOperatorDoneRequiresFlag` / `TestOperatorDoneHatchBypassesOperator` / `TestReopenInvalidatesPassReviews` / `TestMissingCapabilitiesBlockTransition` / `TestAllowDoneWarnsOnStderr` / `TestCapabilityMissingRequiresTaskHint` + MCP `TestTransitionAllowDoneEmitsWarning` / `TestCapabilityMissingRequiresTaskParam` |
| Honesty | Paths A/B/C + Gate G; Path C uses `AllowOperatorDone` (operator-flag supersession); Gate G hatch **retained** |
| Carry-forward | Honesty A/B/C+G; Gate E; Gate F; capability ablation; Gate H; compat checklist; p0x 7/7; x0; Gate C `dry_run:false` N=3 |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` — **product pkgs PASS**; known FAIL only `similar projects/graphify` space (non-product) |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S05-01 starts** Notes; **S05-02 owns completion**. Do **not** scaffold Phase 11 unless Notes explicitly promote. |

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked
- [x] SCOPE-TODOS + board Notes; next `P10-S05-01`

## Out of scope
- Running VERIFY (S05-01)
- Product Go / new MCP tools / daemon / mig `011_*`
- Scaffolding Phase 11 without explicit promotion
- Closing parallel dogfood experiments
- Claiming Phase 09 historical handoff was wrong
