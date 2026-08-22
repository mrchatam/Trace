# P04 / 00-PHASE-PLANNER — Review depth phase scaffold

## Metadata
- id: P04-00
- todo_ids: [P04-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 04 — Review depth & evidence policies** now that Phase 03 Gate E VERIFYs green. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 4, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 4
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 03 VERIFY: [VERIFY-NOTES.md](../phase-03-progressive-planner/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 03)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate E | Green via `evals/replan` `TestPlantedDiscoveryReplan` — do not weaken |
| Gate C | **Go** — do not silently reopen |
| Honesty / p0x / x0 / replan | Keep green |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary |

## Planner work
1. Inspect live review/claim/evidence APIs vs review-depth needs (scope review layer; escape-rate report; residual tracking).
2. Confirm S01→S03 order (or revise with Notes): scope review layer → honesty escape-rate / Gate G prelim → phase verify/handoff.
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 03 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 04 section; set this row `done` with Notes.
6. Update phase README if order/locks change.

## Locked phase defaults (do not weaken)
| Item | Value |
|------|-------|
| Goal | Review depth & evidence policies (`A_PROJECT_PLAN` Phase 4) |
| Validation gate | Gate G preliminary — honesty suite escape-rate report (path locked by scope planners) |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Honesty Paths A/B/C; p0x 7/7; Gate E replan; Gate C artifacts intact |
| MCP | Optional; CLI path primary |

## Cross-scope blast radius
- S01 review-depth surface thickens S02 escape-rate harness.
- S02 must not weaken honesty Paths A/B/C or Gate E / Gate C evidence integrity.
- S03 scaffolds Phase 05 (or records explicit `no successor`).

## Exit criteria
- [x] Phase README accurate vs live Phase 03 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 04 lists planner + scope rows; `P04-00` done after Phase 03 S03-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P04-S01-00`)

## Minimal todos
- [x] Inventory live review APIs vs review-depth gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 04 section
- [x] Update README + mark P04-00 done

## Out of scope
- Product feature implementation
- Re-scoring Gate C / inventing Gate E without planted replan test
- Starting Phase 05 before S03 VERIFY
