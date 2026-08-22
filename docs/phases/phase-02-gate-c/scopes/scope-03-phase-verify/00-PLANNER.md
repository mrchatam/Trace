# P02 / S03 / 00-PLANNER — Phase 02 VERIFY

## Metadata
- id: P02-S03-00
- todo_ids: [P02-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Write, Grep]
- verification: automated

## Objective
Lock Phase 02 VERIFY commands + Gate C decision artifact checks + **DR-HANDOFF** for Phase 03 (progressive planner) — or explicit No-Go / no-successor. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 3
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- S01 Gate C evidence + S02 REVIEW (if any)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Depends-on
- S01 done (`GATE-C-NOTES.md` + `docs/verification/gate-c-x0/` metrics exist; S01-02 review closed).
- S02 done **or** skipped with reason after No-Go.
- VERIFY must re-check Gate C evidence artifacts — **not** Phase 01 dry-run alone.
- **S02 hardening surfaces (S02-02 APPROVE high):** re-check GC-01 (`TestWhyTaskIncludesDiscoveryPlanChange`, `TestTaskContextIncludesDiscoveryPlanChange`) and GC-02 (`TestFixtureReadmeHasNoGTUUIDOracle`; pins hash `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22`). Mode-B packs remain historical; do not require pack rewrite for VERIFY pass. GC-03/04 stay deferred unless promoted. Residual (do not fail VERIFY): global attach of all `discovery_causes_plan_change` edges to every task Expand — see [../scope-02-slice-hardening/REVIEW-NOTES.md](../scope-02-slice-hardening/REVIEW-NOTES.md).

## Handoff policy (phase-locked)
| Gate C outcome | Handoff |
|----------------|---------|
| Go / iterate-with-continue | Scaffold `phase-03-progressive-planner` (name may refine) + `P03-00` board row |
| No-Go | Record explicit stop / `no successor` (user override may reopen later) |

**Phase 01 dry-run ≠ Gate C pass** — VERIFY must re-check Gate C evidence, not trust dry-run alone.

## Planner work
- Thicken `01-verify.md` with evidence table + spawn rules + handoff checklist.
- Ensure `02-scope-review.md` owns handoff completion (mirror Phase 01 S05 pattern).
- Sync SCOPE-TODOS + board.

## Exit criteria
- [x] `01-verify.md` has evidence table + spawn rules + handoff checklist
- [x] `02-scope-review.md` owns handoff completion
- [x] Board synced

## Minimal todos
- [x] Thicken 01/02
- [x] SCOPE-TODOS + board
