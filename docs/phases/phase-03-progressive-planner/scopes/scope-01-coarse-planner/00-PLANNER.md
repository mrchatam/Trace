# P03 / S01 / 00-PLANNER — Coarse progressive planner

## Metadata
- id: P03-S01-00
- todo_ids: [P03-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-coarse-planner.md` for **goal→phase→scope** coarse planning against live domain/store. Lock package paths, APIs, persistence, CLI surface, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 3
- [docs/init/J_BRAINSTORMING_OUTCOMES.md](../../../../init/J_BRAINSTORMING_OUTCOMES.md) — deep-plan current + one lookahead; no LLM auto-backlog
- Live: `internal/domain` (Goal/Task/Discovery/PlanChange), `cmd/trace` add/link
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (coarse-planner gaps)
| Item | Today (post–Phase 02) | S01 need |
|------|----------------------|----------|
| Hierarchy | Goal↔Task (`goal_id` / LinkGoalTask) | Coarse **phase** + **scope** under a goal |
| Planner API | Absent — no `internal/planner` | Library API to create/list coarse plan; deep-plan **current scope only** (+ one lookahead horizon) |
| Persistence | Goals/tasks/events/links in `.trace/trace.db` | Persist plan structure without source blobs; prefer supersede over delete |
| CLI | `trace add`/`link` only | Thin CLI (stdlib argv) if needed — G19 library-first |
| MCP | Optional six tools | **Not** required for S01 (CLI primary) |
| Bars | honesty / p0x / x0 / Gate C artifacts | Keep green |

## Phase defaults already locked (respect; refine paths only)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Honesty / p0x / x0 | Keep green |
| Daemon/HTTP/embeddings | Forbidden as primary |
| MCP | Optional; CLI primary |
| Progressive rule | Deep-plan current scope + one lookahead — not entire backlog (`J` ADOPT / REJECT) |

## Locked by this planner (2026-08-16)

| Item | Value |
|------|-------|
| Package | **`internal/planner`** (not domain Phase/Scope CRUD) |
| Persistence | Store mig **`006_plan_hierarchy.sql`**: `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` |
| Hierarchy | Goal (domain) → phases → scopes; deep plan = ACTIVE `scope_deep_plans` revision |
| Progressive | `DeepPlan` fail-closed unless scope is current; lookahead = next scope shallow only |
| S02 hooks | `SupersedeDeepPlan`, `GetCurrentScope`/`ListScopes`, `auto_replan_count` column |
| CLI | Thin `trace plan` (create-coarse / set-current / deep / show) |
| Out | CONFLICT edges; tasks.scope_id; MCP plan tools; churn budget enforcement (S02) |

## Planner work
- [x] Lock coarse planner surface (`internal/planner`) and persistence model.
- [x] Thicken `01-coarse-planner.md` exit criteria enough to run alone (files, APIs, tests, bars).
- [x] Light-update **upcoming** S02 stubs with expected discovery→replan hooks from the S01 surface.
- [x] Sync SCOPE-TODOS.md + board Notes.

## Exit criteria
- [x] `01-coarse-planner.md` runnable alone
- [x] Package path + coarse hierarchy model locked
- [x] Light S02 Depends note updated
- [x] No product Go in this row

## Minimal todos
- [x] Inventory live Goal/Task APIs vs phase/scope needs (confirm still current)
- [x] Thicken 01 + light S02 Depends
- [x] Mark P03-S01-00 done
