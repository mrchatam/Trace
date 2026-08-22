# Phase 05 — Decision impact & simulation

## Goal

Add **decision impact** classes and alternatives, then a preliminary **Gate F** path (impact precision/recall on planted conflicts). Full `plan simulate` stays deferred (roadmap P13). Close toward `A_PROJECT_PLAN` Phase 5 / Gate F preliminary — **not** a commercial multi-model impact engine (DR-NOIMP remains in force for that).

**Depends on:** Phase 04 review depth **complete** (VERIFY PASS / Gate G prelim green; DR-HANDOFF closed on `P04-S03-02`). Phase planner **`P05-00` done**; S01 **APPROVE** (`P05-S01-02`); S02 **APPROVE** (`P05-S02-02`); S03 planner **`P05-S03-00` done**; VERIFY **`P05-S03-01` PASS** (Gate F prelim green; Phase 06 scaffold started); next runnable: **`P05-S03-02`** (owns DR-HANDOFF completion).

## Prior phase outcomes (live — carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Gate G prelim | **Green** — `evals/honesty` `TestHonestyEscapeRateGateGPrelim` (escapes=1/caught=2/attempts=3; `schema-gate-g.json` v1) |
| Gate E | **Green** — `evals/replan` `TestPlantedDiscoveryReplan` |
| Gate C | **Go** — mean G1 0.800 > B0 0.000; `docs/verification/gate-c-x0/` (`dry_run:false`, N=3) |
| Dry-run ≠ Gate C | Phase 01 dry-run remains regression-only (also ≠ Gate G / Gate F) |
| Honesty / P0-X / X0 / replan | Keep green |
| Scope review | mig 008 `review_residuals`; `LinkReviewScope` / `CountOpenResidualsByScope` |
| Decision surface (S01 live) | mig `009_decision_impact`; findings+alternatives; fail-closed `ImpactReport`; CLI `trace impact`; `decision_affects_task` only |
| Gate F (S02-00 locked) | **`evals/impact`** / **`TestPlantedImpactConflictsGateFPrelim`**; `schema-gate-f.json` v1 + temp `metrics-gate-f.json`; TP=3/FN=0/FP=0/TN=1; precision=1.0; recall=1.0 |
| Migrations live | `001`…`009_decision_impact` |
| Residuals (carry) | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; GC-03/04 deferred; s01_hooks schema looseness |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| VerifiedFact | Still **out** unless a Phase 05 scope planner explicitly promotes with Notes |
| DR-NOIMP | Still PROVISIONAL — Phase 05 = planted/manual impact classes + Gate F prelim, **not** full automated commercial engine |

## Locked phase defaults (do not weaken — P05-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Decision impact & simulation (`A_PROJECT_PLAN` Phase 5) |
| Scope order | **S01** impact classes/alternatives → **S02** Gate F prelim → **S03** VERIFY + Phase 06 handoff |
| Validation | Gate F preliminary = planted precision/recall under **`evals/impact`** — named **`TestPlantedImpactConflictsGateFPrelim`**; schema **`schema-gate-f.json`** v1 + temp **`metrics-gate-f.json`**; tallies TP=3/FN=0/FP=0/TN=1 |
| Impact classes (S01-00 locked) | `SAFE` \| `CAUTION` \| `HIGH` \| `DESTRUCTIVE` \| `REVERSAL`; uncertainty `KNOWN` \| `LIKELY` \| `POSSIBLE` \| `UNKNOWN`; kinds AFFECTED_WORK…UNRESOLVED |
| Package / mig | Extend **`internal/domain`** + store mig **`009_decision_impact`** — **no** second impact stack; **no** planner fork; **no** new entity_links rels |
| `plan simulate` | **Out** this phase (P13 / later) |
| Carry-forward bars | Honesty A/B/C; Gate G prelim; Gate E; p0x 7/7; x0; Gate C `dry_run:false` intact |
| MCP | Optional; CLI path primary |
| Daemon / HTTP / embeddings | Forbidden as primary |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| DR-HANDOFF | Closing VERIFY scaffolds Phase 06 = **`phase-06-environment-capability`** (or records `no successor`) |

## Scope run order (locked P05-00)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Impact classes / alternatives surface | P05-S01-00/01/02 | `scopes/scope-01-impact-classes/` |
| S02 | Gate F prelim (planted conflicts) | P05-S02-00/01/02 | `scopes/scope-02-gate-f-prelim/` |
| S03 | Phase VERIFY + Phase 06 handoff | P05-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

## Out of scope (this phase — until planners promote)

- Full commercial impact engine / false-confidence UX theater (DR-NOIMP)
- `plan simulate` / PlanVersion branches (roadmap P13)
- Daemon / always-on HTTP / embeddings
- Reopening Gate C / inventing Gate G without `TestHonestyEscapeRateGateGPrelim`
- Commercial multi-model Gate F (no evidence path — planted `evals/impact` only)
- Commercial A1 validation
- Multi-agent orchestration / UI
- VerifiedFact promotion engine

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 5
- [`docs/init/PROJECT_MODEL_SNAPSHOT.md`](../../init/PROJECT_MODEL_SNAPSHOT.md) Gate F
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-HANDOFF / DR-NOIMP
- [`docs/ROADMAP.md`](../../ROADMAP.md) P12 impact bands
- Phase 04 VERIFY: [`../phase-04-review-depth/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-04-review-depth/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
