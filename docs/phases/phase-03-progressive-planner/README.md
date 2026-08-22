# Phase 03 — Progressive planner (minimal)

## Goal

Deliver a **minimal progressive planner**: coarse goal→phase→scope planning, deep-plan of the current scope, and discovery→PlanChange propagation with churn controls. Prove with a **replan demo** on a planted discovery; close with Gate E mini-eval readiness (`A_PROJECT_PLAN` Phase 3).

**Depends on:** Phase 02 Gate C **Go** (VERIFY PASS; DR-HANDOFF closed on `P02-S03-02`). Phase 02 is **complete**. **Phase 03 complete** (2026-08-16): S01–S03 all `done`; Gate E green (`evals/replan`/`TestPlantedDiscoveryReplan`); DR-HANDOFF closed → Phase 04. Next board row: **`P04-00`**.

## Prior phase outcomes (live — carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Gate C | **Go** — mean G1 `understanding_accuracy` 0.800 > B0 0.000; evidence `GATE-C-NOTES.md` + `docs/verification/gate-c-x0/` (`dry_run:false`, N=3) |
| Dry-run ≠ Gate C | Phase 01 dry-run remains regression-only — do not reopen Gate C without contradicting evidence |
| Slice harden | GC-01/02 done; GC-03/04 deferred; fixture hash `15fe50a1…` |
| Honesty / P0-X / X0 | `evals/honesty` Paths A/B/C; `evals/p0x` 7/7; `evals/x0` — keep green |
| Causal surface | `internal/domain` CreateGoal/Task/Discovery/PlanChange + LinkGoalTask / LinkDiscoveryPlanChange (`discovery_causes_plan_change`) |
| Retrieval | Why + TaskContext/Expand; GC-01 attaches DPC endpoints (global-attach residual) |
| CLI | `add`/`link`/`transition`/`seed`/`why`/`context`/`review`/`plan` (create-coarse\|set-current\|deep\|show); S02 adds discovery `--severity` + `plan apply-discovery`/`ack-replan` |
| Residual | Global attach of all `discovery_causes_plan_change` on task Expand (non-blocking; prefer scoping if measured) |
| Daemon / HTTP / embeddings | Still forbidden as primary |

## Live inventory → Phase 03 gaps (P03-00)

| Need (`A_PROJECT_PLAN` Phase 3) | Live today | Gap owner |
|---------------------------------|------------|-----------|
| Coarse goal→phase→scope | **S01 shipped** — `internal/planner` + mig `006_plan_hierarchy` | **S01 done** |
| Deep-plan current scope (+ one lookahead) | Live: current + one lookahead fail-closed; `SupersedeDeepPlan` | **S01 done** |
| Discovery→PlanChange + churn controls | Mig 007 + `ApplyDiscoveryReplan` + N=5 churn; S02 APPROVE | **S02 done** |
| Replan demo on planted discovery | `evals/replan` `TestPlantedDiscoveryReplan` PASS | **S02 done** |
| Gate E mini-eval on fixture | **VERIFY PASS + S03-02 APPROVE:** `evals/replan` `TestPlantedDiscoveryReplan` + severity/churn green; handoff closed | **S03 done** |

## Locked phase defaults (do not weaken)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Progressive planner minimal (`A_PROJECT_PLAN` Phase 3) |
| Validation | Gate E mini-eval = `evals/replan` `TestPlantedDiscoveryReplan` (PLAN_AFFECTING+; N=5 churn); locked by S03-00 |
| Phase 04 handoff | Folder **`phase-04-review-depth`**; S03-01 starts scaffold; S03-02 owns completion |
| Carry-forward bars | Honesty Paths A/B/C; p0x 7/7; Gate C artifacts intact; x0 packages green |
| MCP | Optional; **CLI path primary** |
| Daemon / HTTP / embeddings | Forbidden as primary |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| Churn law | G16 — plan-affecting updates outside active task need acknowledgment; scopes have replan budget (`J_BRAINSTORMING` ADOPT) |
| DR-HANDOFF | Closing VERIFY scaffolds Phase 04 — Review depth & evidence policies (`A_PROJECT_PLAN` Phase 4) — or records `no successor` |

## Scope run order (locked)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Coarse goal→phase→scope planner surface | P03-S01-00/01/02 | `scopes/scope-01-coarse-planner/` |
| S02 | Discovery→PlanChange replan + churn demo | P03-S02-00/01/02 | `scopes/scope-02-discovery-replan/` |
| S03 | Phase VERIFY + Gate E bar + Phase 04 handoff | P03-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

## Cross-scope blast radius

- S01 planner surface thickens S02 replan APIs (do not invent a second planning stack).
- S02 must not weaken honesty / P0-X / Gate C evidence integrity; DPC-global residual stays unless measurement demands scoping.
- S03 scaffolds Phase 04 (or records explicit `no successor`).

## Out of scope (this phase)

- Full multi-agent orchestration / UI
- Daemon / always-on HTTP
- Embeddings / env graph / impact engine (later phases)
- Rewriting Mode-B Gate C packs to invent q3 pass
- Declaring commercial packaging settled
- Auto-generating entire backlog from goal via LLM (`J_BRAINSTORMING` REJECT)

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 3–4
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-HANDOFF, DR-AGENT
- [`docs/init/G_PROJECT_LAWS.md`](../../init/G_PROJECT_LAWS.md) G16 churn
- [`docs/init/J_BRAINSTORMING_OUTCOMES.md`](../../init/J_BRAINSTORMING_OUTCOMES.md) planning ADOPTs
- Phase 02 VERIFY: [`../phase-02-gate-c/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-02-gate-c/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
- Protocol: [`docs/rules/agent-loop-protocol.md`](../../rules/agent-loop-protocol.md)
