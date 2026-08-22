# P05 / S03 / 00-PLANNER — Phase 05 VERIFY

## Metadata
- id: P05-S03-00
- todo_ids: [P05-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 05 VERIFY commands + evidence table + DR-HANDOFF Phase 06 checklist. Thicken `01-verify.md` / `02-scope-review.md`. No product Go.

## References
- [phase README](../../README.md)
- Gate F path from S02 APPROVE (`evals/impact` named test)
- Phase 04 VERIFY pattern: [VERIFY-NOTES.md](../../../phase-04-review-depth/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
- DR-HANDOFF → Phase 06 folder **`phase-06-environment-capability`** (A_PROJECT_PLAN Phase 6)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Depends
S01 + S02 APPROVE (**P05-S02-02 APPROVE high** — 2026-08-16; see [../scope-02-gate-f-prelim/REVIEW-NOTES.md](../scope-02-gate-f-prelim/REVIEW-NOTES.md)).

### Gate F re-prove (from P05-S02-00 — stamp into VERIFY command table under S03-00; confirmed live by S02-02)
| Item | Value |
|------|-------|
| Package | `evals/impact` |
| Named test | `TestPlantedImpactConflictsGateFPrelim` |
| Schema / metrics | `schema-gate-f.json` v1 + temp `metrics-gate-f.json` |
| Tallies | TP=3 FN=0 FP=0 TN=1; precision=1.0; recall=1.0 |
| Command | `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` |
| Evidence note | VERIFY must re-run the named test (harness bar); do not treat board Notes alone as Gate F |

## Phase locks to encode (P05-00 + S02-00 stamp)
| Item | Value |
|------|-------|
| Gate F re-prove | `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; p0x 7/7; x0; domain/store/planner; Gate C `dry_run:false` |
| Dry-run ≠ Gate C / ≠ Gate F | Explicit in VERIFY-NOTES |
| Spawn on fail | 01a / 01b / 01c convention |
| DR-HANDOFF | S03-01 starts `phase-06-environment-capability` + `P06-00`; S03-02 owns completion (or `no successor`) |

## Planner work
- [x] Lock Gate F named re-prove command + carry-forward command table
- [x] Spawn 01a/b/c convention on fail
- [x] DR-HANDOFF Phase 06 folder name checklist
- [x] Thicken 01-verify + 02

## Locked defaults (stamp — respect)

| Item | Value |
|------|-------|
| Product Go | Forbidden on VERIFY planner/verify/review (scaffold docs only) |
| Gate F bar | `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` (+ schema/metrics; TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; p0x 7/7; x0; domain/store/planner; Gate C `dry_run:false`; full `./...` incl. `evals/impact` |
| Dry-run ≠ Gate C / ≠ Gate F | Explicit in VERIFY-NOTES |
| Spawn on fail | `P05-S03-01a` / `01b` / `01c` immediately below VERIFY |
| Phase 06 folder | **`phase-06-environment-capability`** |
| DR-HANDOFF | S03-01 starts Phase 06 scaffold + `P06-00`; S03-02 owns completion |

## Exit criteria
- [x] VERIFY prompt runnable alone (`01-verify.md` thickened)
- [x] `02-scope-review.md` owns DR-HANDOFF completion checklist
- [x] No product Go

## Out of scope
- Product features; starting Phase 06 before VERIFY pass; inventing commercial Gate F; running VERIFY itself (S03-01)
