# P01 / 00-PHASE-PLANNER — X0 readiness phase scaffold

## Metadata
- id: P01-00
- todo_ids: [P01-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 01** now that Phase 00 / P0-X is complete. Confirm scope order, write/refresh scope stub prompts (00/01/02), lock phase defaults from live repo + init plan, and register/sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 1
- [docs/init/I_BENCHMARK_PLAN.md](../../init/I_BENCHMARK_PLAN.md) Experiment X0
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-AGENT, DR-SURFACE
- [docs/TODO.md](../../TODO.md)
- Phase 00 VERIFY: [VERIFY-NOTES.md](../phase-00-foundation/scopes/scope-09-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Layout | `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` + `cmd/trace` |
| Store | `.trace/trace.db`, modernc sqlite, embed migrations |
| CLI | stdlib argv; G19 library-only |
| P0 | Closed (7/7) — do not reopen P0-X bar as Phase 01 exit |

## Planner work
1. Inspect live `internal/domain` Claim/Evidence stubs and DONE policy; note what S01 must replace.
2. Confirm S01→S05 order: review promotion → honesty → X0 harness → MCP → phase verify.
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only (this phase); never edit Phase 00 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 01 section; set this row `done` with Notes.
6. Update phase README if order/locks change.

## Locked phase defaults (do not weaken)
| Item | Value |
|------|-------|
| X0 conditions | B0 = agent+repo tools; G1 = agent+`trace` CLI context/why (still may read repo) |
| X0 corpus | Start with synthetic `fixtures/x0` (+ seed); no new OSS required this phase |
| MCP | Thin adapter over library; **after** S01–S03; no business logic in MCP |
| Daemon/HTTP | Still forbidden as primary surface |
| Embeddings | Still forbidden |
| Human verification | Honesty demo may mark `verification: mixed` if scripted claim is insufficient — prefer automated fail-closed demo first |
| Review policy | Every scope gets `02-review` before the next scope’s implement runs in board order |

## Cross-scope blast radius
- S01 changes DONE promotion → thickens S02 honesty expectations and CLI `transition` docs in **upcoming** S02/S03 prompts only.
- S03 harness must not break `evals/p0x` (keep both).
- S04 must not fork domain logic into MCP package.

## Exit criteria
- [x] Phase README accurate vs live Phase 00 outcomes
- [x] Scopes S01–S05 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 01 lists planner + all scope rows; `P01-00` no longer `blocked`/`pending`/`in_progress`
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P01-S01-00`)
- [x] **Handoff check noted:** when this phase later VERIFYs, scaffold Phase 02 (Gate C) planner+stubs+board **before** marking VERIFY done (DR-HANDOFF) — encoded in S05 stubs

## Minimal todos
- [x] Inventory domain Claim/Evidence/DONE vs Phase 01 goals
- [x] Write/refresh S01–S05 stub prompts
- [x] Sync TODO.md Phase 01 section
- [x] Update README + mark P01-00 done
