# P03 / S03 / 00-PLANNER — Phase 03 VERIFY

## Metadata
- id: P03-S03-00
- todo_ids: [P03-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize Phase 03 VERIFY + review prompts: Gate E mini-eval bar, evidence table, spawn rules, and **DR-HANDOFF** to Phase 04 (Review depth & evidence policies). No product code.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 3–4
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) Gate E
- Prior VERIFY pattern: Phase 02 `scopes/scope-03-phase-verify/`
- [docs/TODO.md](../../../../TODO.md)
- Live S02: [../scope-02-discovery-replan/REVIEW-NOTES.md](../scope-02-discovery-replan/REVIEW-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (VERIFY / Gate E) — locked

| Item | Locked value |
|------|--------------|
| Gate E definition | Discovery propagation reduces downstream rework |
| Gate E path / harness | `evals/replan` **`TestPlantedDiscoveryReplan`** |
| Severity | Only `PLAN_AFFECTING`+ auto-replan; INFO does not |
| Churn | N=5 fail-closed + ack (`DefaultMaxAutoReplans=5`) |
| S01 | `internal/planner` + mig `006_plan_hierarchy.sql` |
| S02 | mig `007_discovery_severity.sql` + `ApplyDiscoveryReplan` |
| Regression bars | honesty A/B/C; p0x 7/7; x0; Gate C `dry_run:false` artifacts intact |
| Spawn | On fail: `01a` implement / `01b` review / optional `01c` re-VERIFY |
| DR-HANDOFF | S03-01 starts `phase-04-review-depth` + `P04-00`; S03-02 owns completion |

## Phase defaults already locked (respect)
| Item | Value |
|------|-------|
| VERIFY | Independent re-run; spawn 01a/01b/01c on fail |
| DR-HANDOFF | S03-01 starts Phase 04 scaffold; S03-02 owns completion |
| Carry-forward | Honesty / p0x / Gate C integrity |
| Product Go | Forbidden in VERIFY rows |

## Planner work
- Lock VERIFY commands + Gate E evidence path + evidence table columns.
- Thicken `01-verify.md` + `02-scope-review.md` (S03-02 owns handoff completion).
- Sync SCOPE-TODOS.md + board Notes.

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable
- [x] Gate E path + regression command list locked
- [x] No product Go in this row

## Minimal todos
- [x] Thicken VERIFY + review prompts from live S01/S02 artifacts
- [x] Mark P03-S03-00 done
