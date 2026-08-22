# P02 / 00-PHASE-PLANNER — Gate C phase scaffold

## Metadata
- id: P02-00
- todo_ids: [P02-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 02 — Gate C evaluation & slice hardening** now that Phase 01 X0 readiness VERIFYs green. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` / `I_BENCHMARK_PLAN`, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that. **Do not** claim Gate C pass from Phase 01 dry-run alone.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 2
- [docs/init/I_BENCHMARK_PLAN.md](../../init/I_BENCHMARK_PLAN.md) Experiment X0
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-AGENT, DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 01 VERIFY: [VERIFY-NOTES.md](../phase-01-x0-readiness/scopes/scope-05-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 01)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| X0 dry-run | `evals/x0` schema v1; B0/G1 temp metrics with `dry_run:true` |
| Honesty | `evals/honesty` Paths A/B/C — keep green |
| P0-X | `evals/p0x` 7/7 — regression-keep |
| MCP | Thin stdio adapter; not required for Gate C (DR-AGENT) |
| Daemon/HTTP/embeddings | Still forbidden as primary |

## Planner work
1. Inspect live `evals/x0` dry-run vs Gate C needs (agent runs, scoring, N≥3, kill criteria).
2. Confirm S01→S03 order: Gate C run → slice hardening → phase verify/handoff.
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 01 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 02 section; set this row `done` with Notes.
6. Update phase README if order/locks change.

## Locked phase defaults (do not weaken)
| Item | Value |
|------|-------|
| Gate C bar | Documented pass / fail / iterate with evidence table (not vibes) |
| Conditions | B0 = repo tools; G1 = `trace why`/`context` (+ repo OK) |
| Instrument | Prefer extending `evals/x0`; keep honesty + p0x separate |
| Kill criteria | G1 ≤ B0 within error **and** non-trivial seeding cost → thesis endangered (`I_BENCHMARK_PLAN`) |
| MCP | Optional; CLI path sufficient |
| Review policy | Every scope gets `02-review` before next scope implement |

## Cross-scope blast radius
- S01 scoring protocol thickens S02 issue backlog.
- S02 must not weaken honesty / P0-X bars.
- S03 scaffolds Phase 03 only on Gate C Go (or records explicit No-Go stop).

## Exit criteria
- [x] Phase README accurate vs live Phase 01 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 02 lists planner + scope rows; `P02-00` done after Phase 01 S05-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P02-S01-00`)
- [x] Explicit: Phase 01 dry-run ≠ Gate C pass

## Minimal todos
- [x] Inventory `evals/x0` dry-run vs Gate C gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 02 section
- [x] Update README + mark P02-00 done
