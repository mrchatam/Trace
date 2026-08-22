# P03 / 00-PHASE-PLANNER — Progressive planner phase scaffold

## Metadata
- id: P03-00
- todo_ids: [P03-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 03 — Progressive planner (minimal)** now that Phase 02 Gate C VERIFYs green (Go). Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 3, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 3
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 02 VERIFY: [VERIFY-NOTES.md](../phase-02-gate-c/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 02)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate C | **Go** — do not silently reopen No-Go without contradicting evidence |
| Honesty / p0x / x0 | Keep green |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary |

## Planner work
1. Inspect live domain/retrieval/compiler vs progressive-planner needs (goal→phase→scope; discovery→PlanChange churn).
2. Confirm S01→S03 order: coarse planner → discovery replan demo → phase verify/handoff.
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 02 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 03 section; set this row `done` with Notes.
6. Update phase README if order/locks change.

## Locked phase defaults (do not weaken)
| Item | Value |
|------|-------|
| Goal | Minimal progressive planner + replan demo (`A_PROJECT_PLAN` Phase 3) |
| Validation gate | Gate E mini-eval on fixture (path locked by scope planner) |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Honesty Paths A/B/C; p0x 7/7; Gate C artifacts intact |
| MCP | Optional; CLI path primary |

## Cross-scope blast radius
- S01 planner surface thickens S02 replan APIs.
- S02 must not weaken honesty / P0-X / Gate C evidence integrity.
- S03 scaffolds Phase 04 (or records explicit `no successor`).

## Exit criteria
- [x] Phase README accurate vs live Phase 02 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 03 lists planner + scope rows; `P03-00` done after Phase 02 S03-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P03-S01-00`)

## Minimal todos
- [x] Inventory live graph APIs vs progressive-planner gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 03 section
- [x] Update README + mark P03-00 done
