# P04 / S02 / 00-PLANNER — Honesty escape-rate / Gate G prelim

## Metadata
- id: P04-S02-00
- todo_ids: [P04-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-honesty-escape-rate.md` for **honesty suite escape-rate reporting** toward Gate G preliminary. Lock harness path, metrics schema, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 4 — Gate G
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) H5
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) Gate G
- Live: `evals/honesty` Paths A/B/C; S01 review-depth surface (APPROVE)
- Depends: S01 done (`P04-S01-02` APPROVE)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (escape-rate gaps)
| Item | Today | S02 need |
|------|-------|----------|
| Honesty demo | `TestHonestyFailClosedPlantedClaim` Paths A/B/C | Keep; do not remove or weaken |
| Escape-rate metric | Absent | Gate G prelim **report** — `TestHonestyEscapeRateGateGPrelim` + `schema-gate-g.json` / temp `metrics-gate-g.json` |
| S01 hooks (APPROVE) | `review_judges_scope`; `AddResidual` / `ListResidualsBy*` / **`CountOpenResidualsByScope`**; `POLICY_EXCEPTION` | Consume for OPEN residual tallies; fail closed if missing; **do not** weaken Paths A/B/C |
| Bars | p0x / x0 / replan / Gate C / Gate E | Keep intact |

## Depends-on (S01 surface — locked consume)
- Link/query **`review_judges_scope`** via `LinkReviewScope` (review → `plan_scope`)
- Residual severity INFO\|WARN\|BLOCKING; status OPEN\|ACKED\|RESOLVED
- **`CountOpenResidualsByScope`** for Gate G residual tallies
- Planted code **`POLICY_EXCEPTION`** for documented residual signal — never rewrite honesty A/B/C to require residuals
- Coarse scope via `planner.CreateCoarsePlan` (shared store with domain)

## Planner work
- [x] Lock Gate G prelim / escape-rate harness path (**extend `evals/honesty`**).
- [x] Thicken `01-honesty-escape-rate.md` to run alone.
- [x] Light-update S03 VERIFY Depends notes (named test / metrics path).
- [x] Sync SCOPE-TODOS + board; mark done.

## Locked defaults (final — do not re-debate in S02-01)
| Item | Value |
|------|-------|
| Keep Paths A/B/C | Do not remove or weaken `TestHonestyFailClosedPlantedClaim` |
| Harness | Extend **`evals/honesty`** (no new evals package) |
| Named test | **`TestHonestyEscapeRateGateGPrelim`** |
| Schema | **`evals/honesty/schema-gate-g.json`** `schema_version` **1** |
| Artifact | Temp **`metrics-gate-g.json`** (validate in-test) |
| Escape formula | escapes=1, caught=2, attempts=3; escape = DONE without Review PASS via hatch only; PASS→DONE excluded from attempts |
| S01 consume | `LinkReviewScope` + OPEN `POLICY_EXCEPTION` + `CountOpenResidualsByScope==1` |
| Bars | p0x / x0 / replan / Gate E / Gate C intact |
| Out | Multi-model commercial review cost explosion; daemon/HTTP; inventing Gate G without harness evidence; VerifiedFact |

## Out of scope
- Product Go beyond harness; full Gate G production policy; Phase 05 impact; VerifiedFact
