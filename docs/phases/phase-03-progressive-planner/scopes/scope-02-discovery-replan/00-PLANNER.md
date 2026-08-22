# P03 / S02 / 00-PLANNER — Discovery→PlanChange replan

## Metadata
- id: P03-S02-00
- todo_ids: [P03-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-discovery-replan.md` for **discovery→PlanChange propagation with churn controls** and a planted-discovery replan demo. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 3
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) G16
- [docs/init/J_BRAINSTORMING_OUTCOMES.md](../../../../init/J_BRAINSTORMING_OUTCOMES.md) — churn budget; severity tiers; supersede not delete
- Live: `discovery_causes_plan_change` (GC-01); S02 DPC-global residual; S01 planner surface (after S01 done)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (replan gaps)
| Item | Today | S02 need |
|------|-------|----------|
| Link | `LinkDiscoveryPlanChange` + CLI `link discovery-plan-change` | Keep; wire into replan flow |
| Retrieval | Why/TaskContext attach DPC (global) | Prefer **scoping** only if measurement demands; do not weaken honesty/Gate C |
| Severity | Discovery has no INFO / PLAN_AFFECTING / BLOCKING | Only PLAN_AFFECTING+ auto-opens replan (`J` ADOPT) |
| Churn | `plan_scopes.auto_replan_count` (S01 column + read via `ScopeView`; **no mutator/enforce yet**) | Max N auto-replans per scope without human ack (G16 / DR-CHURN default N=5); fail closed past budget; **add store increment/ack helper** (S01 only DEFAULT 0 + preserve on body update) |
| Demo | Fixture has static discovery+plan_change UUIDs | **Planted discovery** replan demo (library test and/or `evals/` harness) |
| S01 hooks | **Shipped (P03-S01-02 APPROVE)** — consume; do not fork | `internal/planner`: `SupersedeDeepPlan`, `GetCurrentScope`/`ListScopes`, `GetPlan`; current cursor via `goal_plan_state`; column `auto_replan_count` present — S02 owns increment/check/ack |

## Phase defaults already locked (respect)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Carry-forward | Keep honesty / p0x / Gate C artifacts intact |
| Residual | Prefer scoping DPC attach if measurement demands (do not weaken without evidence) |
| MCP | Optional; CLI primary |

## Depends note (from S01 surface — do not invent a second planner)

S01 shipped (`P03-S01-02` APPROVE). S02 **must** call `internal/planner.SupersedeDeepPlan` (and read current/list scopes) to apply PLAN_AFFECTING discovery→PlanChange updates. Keep `LinkDiscoveryPlanChange` for causality. Churn: add store mutator + enforce budget on `plan_scopes.auto_replan_count` (fail closed; human ack resets/allows) — S01 does not increment. No parallel plan stack under `internal/domain`. Live CLI note: `trace plan create-coarse` uses ordered `--phase`/`--scope` argv (not `--from` JSON).

## Planner work
- Lock replan demo fixture/harness path + churn controls + CLI/library surface against **live S01**.
- Thicken `01-discovery-replan.md`; light-update upcoming S03 VERIFY Depends (Gate E inputs).
- Sync SCOPE-TODOS.md + board Notes.

## Exit criteria
- [x] `01-discovery-replan.md` runnable alone
- [x] Churn budget + severity policy locked
- [x] Light S03 Depends note updated
- [x] No product Go in this row

## Minimal todos
- [x] Inventory discovery↔plan_change + S01 planner hooks (live)
- [x] Thicken 01 + light S03 Depends
- [x] Mark P03-S02-00 done
