# P06 / S03 / 00-PLANNER — Phase 06 VERIFY

## Metadata
- id: P06-S03-00
- todo_ids: [P06-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 06 VERIFY commands + evidence table + DR-HANDOFF duties for Phase 07. Thicken `01-verify.md` / `02-scope-review.md`. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Pattern: Phase 05 VERIFY [`../../../phase-05-decision-impact/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-05-decision-impact/scopes/scope-03-phase-verify/01-verify.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 7
- S02 APPROVE: [../scope-02-capability-selection/REVIEW-NOTES.md](../scope-02-capability-selection/REVIEW-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Phase defaults already locked (respect — P06-00 + S02-00 FINAL)
| Item | Value |
|------|-------|
| Ablation gate | **`evals/capability`** / **`TestPlantedCapabilitySelectionAblation`** + `schema-capability.json` v1 / temp `metrics-capability.json` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) — **FINAL P06-S02-00** |
| Ablation re-prove | `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` |
| Phase 07 folder | **`phase-07-performance-ladder`** |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; Gate F; p0x; x0; Gate C `dry_run:false` |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation |
| DR-HANDOFF | S03-01 starts scaffold + `P07-00`; S03-02 owns completion |

## Depends (light — from S02-00)
S02 ablation APPROVE (**P06-S02-02 APPROVE high** — 2026-08-16). Re-prove command above is authoritative; S03-00 owns full evidence table + carry-forward command list.

### Ablation re-prove (from P06-S02-00 — stamped into VERIFY)
| Item | Value |
|------|-------|
| Package | `evals/capability` |
| Named test | `TestPlantedCapabilitySelectionAblation` |
| Schema / metrics | `schema-capability.json` v1 + temp `metrics-capability.json` |
| Tallies | TP=3 FN=0 FP=0 TN=1; precision=1.0; recall=1.0 |
| Command | `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` |
| Evidence note | VERIFY must re-run the named test (harness bar); do not treat board Notes alone as ablation pass |

## Locked defaults (stamp — respect)

| Item | Value |
|------|-------|
| Product Go | Forbidden on VERIFY planner/verify/review (scaffold docs only) |
| Ablation bar | `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` (+ schema/metrics; TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; Gate F; p0x 7/7; x0; domain/store/planner/compiler; Gate C `dry_run:false`; full `./...` incl. `evals/capability` |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | Explicit in VERIFY-NOTES |
| Spawn on fail | `P06-S03-01a` / `01b` / `01c` immediately below VERIFY |
| Phase 07 folder | **`phase-07-performance-ladder`** |
| DR-HANDOFF | S03-01 starts Phase 07 scaffold + `P07-00`; S03-02 owns completion |

## Planner work
- [x] Lock ablation / capability gate re-prove commands from live S02
- [x] Lock carry-forward Gate F/G/E / honesty / p0x / x0 / Gate C commands
- [x] Thicken 01-verify (evidence table + spawn 01a/b/c + Phase 07 checklist)
- [x] Thicken 02 (owns DR-HANDOFF completion)
- [x] Sync SCOPE-TODOS + board

## Exit criteria
- [x] `01-verify.md` runnable alone with locked commands
- [x] Phase 07 folder name locked (**`phase-07-performance-ladder`**)
- [x] No product Go

## Out of scope
- Running VERIFY (S03-01); closing phase (S03-02); implementing Phase 07 features
