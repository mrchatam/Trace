# P04 / S03 / 00-PLANNER — Phase 04 VERIFY

## Metadata
- id: P04-S03-00
- todo_ids: [P04-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 04 VERIFY commands + evidence table + spawn rights + **DR-HANDOFF** Phase 05 (`phase-05-decision-impact`). Thicken `01-verify.md` / `02-scope-review.md`. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 5
- Pattern: Phase 03 VERIFY [`../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/`](../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/)
- Depends: S01 + S02 done (`P04-S02-02` APPROVE high); Gate G prelim = `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** + **`schema-gate-g.json`** / temp **`metrics-gate-g.json`** (locked P04-S02-00; planted escapes=1/caught=2/attempts=3 + OPEN `POLICY_EXCEPTION`); keep Paths A/B/C; S01 APIs (`review_judges_scope` + residuals / `CountOpenResidualsByScope`) must remain green in VERIFY suite

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Planner work
- [x] Lock Gate G prelim bar path (named test / metrics from S02) + carry-forward honesty/p0x/x0/replan/Gate C/Gate E commands.
- [x] Thicken 01-verify (evidence table + spawn 01a/b/c + Phase 05 checklist) + 02 (owns DR-HANDOFF completion).
- [x] Sync SCOPE-TODOS + board; mark done.

## Locked defaults (respect — P04-00)
| Item | Value |
|------|-------|
| Product Go | Forbidden on VERIFY |
| Gate G bar | `CGO_ENABLED=0 go test ./evals/honesty/... -run TestHonestyEscapeRateGateGPrelim` (+ schema/metrics evidence); exact VERIFY command list finalized in S03-00 |
| Phase 05 folder | **`phase-05-decision-impact`** |
| DR-HANDOFF | S03-01 starts Phase 05 scaffold; S03-02 owns completion |
| Dry-run ≠ Gate C | Explicit in VERIFY-NOTES |
| Carry-forward | Honesty A/B/C; p0x 7/7; Gate E `TestPlantedDiscoveryReplan`; Gate C `dry_run:false` |

## Out of scope
- Implementing Phase 05 features; reopening Gate C without contradicting evidence; inventing Gate G without harness
