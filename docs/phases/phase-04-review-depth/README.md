# Phase 04 — Review depth & evidence policies

## Goal

Deepen the **review / evidence** layer: scope-level review policies, richer verification rules, and residual tracking so honesty tests can report escape rates. Close toward **Gate G preliminary** (`A_PROJECT_PLAN` Phase 4).

**Depends on:** Phase 03 progressive planner **complete** (VERIFY PASS / Gate E mini-eval green; DR-HANDOFF closed on `P03-S03-02`). **Phase 04 complete** (2026-08-16) — S01 + S02 APPROVE; S03 VERIFY PASS / Gate G prelim green; DR-HANDOFF closed on `P04-S03-02`. Next runnable: **`P05-00`**.

## Prior phase outcomes (live — carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Gate E | **Green** — `evals/replan` `TestPlantedDiscoveryReplan` (PLAN_AFFECTING+; churn N=5 + ack) |
| Gate C | **Go** — mean G1 0.800 > B0 0.000; `docs/verification/gate-c-x0/` (`dry_run:false`, N=3) |
| Dry-run ≠ Gate C | Phase 01 dry-run remains regression-only |
| Honesty / P0-X / X0 / replan | Keep green (`evals/{honesty,p0x,x0,replan}`) |
| Planner | `internal/planner` + mig 006/007; `ApplyDiscoveryReplan` / `AckReplan` |
| Review surface (Phase 01) | Claim→Evidence→Review→DONE; task-level `review_judges_task`; `AllowDoneWithoutReview` escape hatch exists |
| Scope review (Phase 04 S01) | mig 008 `review_residuals`; `LinkReviewScope` / `review_judges_scope`; `CountOpenResidualsByScope`; codes incl. `POLICY_EXCEPTION` |
| Residuals (carry) | DPC-global attach; non-tx Apply; UNIQUE re-link; MCP no severity; GC-03/04 deferred |
| Daemon / HTTP / embeddings | Still forbidden as primary |

## Live inventory → Phase 04 gaps (P04-00)

| Need (`A_PROJECT_PLAN` Phase 4) | Live today | Gap owner |
|---------------------------------|------------|-----------|
| Task Review PASS→DONE | `CreateReview` / `SetReviewResult` / `LinkReviewTask`; mig 005; thin `trace review` | Keep; do not weaken |
| Claim↔Evidence | `CreateClaim` / `CreateEvidence` / `LinkClaimEvidence` | Keep |
| Scope-level review | **Shipped (S01 APPROVE)** — `review_judges_scope` via `LinkReviewScope`; mig 008 | Keep |
| Richer evidence / residual tracking | **Shipped (S01 APPROVE)** — `review_residuals` + Add/List/CountOpen/SetStatus | Keep |
| Escape-rate report / Gate G prelim | **Done (S02 APPROVE)** — `TestHonestyEscapeRateGateGPrelim`; `schema-gate-g.json` v1 / temp `metrics-gate-g.json`; escapes=1/caught=2/attempts=3 | **S02** |
| Phase VERIFY + Phase 05 handoff | **Done (S03-02 APPROVE)** — Gate G green; DR-HANDOFF complete (`phase-05-decision-impact` + `P05-00`) | **S03** |

## Locked phase defaults (do not weaken)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Review depth & evidence policies (`A_PROJECT_PLAN` Phase 4) |
| Validation | Gate G preliminary — `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** + **`schema-gate-g.json`** / temp **`metrics-gate-g.json`** (locked `P04-S02-00`) |
| Carry-forward bars | Honesty Paths A/B/C; p0x 7/7; Gate E replan; Gate C artifacts intact; x0 green |
| MCP | Optional; CLI path primary |
| Daemon / HTTP / embeddings | Forbidden as primary |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| VerifiedFact | **Out** this phase (Phase 01 deferred; cost risk) — residuals only, not promotion engine |
| AllowDoneWithoutReview | Keep as explicit escape; honesty Paths A/B/C never set it true |
| Phase 05 handoff | Folder **`phase-05-decision-impact`**; S03-01 starts scaffold; S03-02 owns completion |
| DR-HANDOFF | Closing VERIFY scaffolds Phase 05 (Decision impact & simulation) or records `no successor` |

## Scope run order (locked)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Scope review layer / richer evidence + residual hooks | P04-S01-00/01/02 | `scopes/scope-01-scope-review-layer/` |
| S02 | Honesty escape-rate / Gate G prelim harness | P04-S02-00/01/02 | `scopes/scope-02-honesty-escape-rate/` |
| S03 | Phase VERIFY + Phase 05 handoff | P04-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

## Cross-scope blast radius

- S01 review-depth surface thickens S02 escape-rate harness (hooks S02 can count).
- S02 must not weaken honesty Paths A/B/C or Gate E / Gate C evidence integrity.
- S03 scaffolds Phase 05 (`phase-05-decision-impact`) or records explicit `no successor`.

## Out of scope (this phase)

- Full multi-agent orchestration / UI
- Daemon / always-on HTTP / embeddings
- Reopening Gate C / inventing Gate E without `TestPlantedDiscoveryReplan`
- Commercial A1 validation
- VerifiedFact promotion engine
- Impact engine / plan simulate (Phase 05+)
- Multi-model commercial review cost explosion

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 4–5
- [`docs/init/I_BENCHMARK_PLAN.md`](../../init/I_BENCHMARK_PLAN.md) H5 / Gate G
- [`docs/init/H_VERIFICATION_STRATEGY.md`](../../init/H_VERIFICATION_STRATEGY.md) scope review layer
- [`docs/init/PROJECT_MODEL_SNAPSHOT.md`](../../init/PROJECT_MODEL_SNAPSHOT.md) Gate G
- Phase 03 VERIFY: [`../phase-03-progressive-planner/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-03-progressive-planner/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
- Protocol: [`docs/rules/agent-loop-protocol.md`](../../rules/agent-loop-protocol.md)
