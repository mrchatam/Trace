# P01 / S05 / 00-PLANNER — Phase 01 VERIFY

## Metadata
- id: P01-S05-00
- todo_ids: [P01-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-verify.md` for **Phase 01 VERIFY** gate: honesty green, X0 dry-run B0+G1 metrics, MCP documented, p0x regression, **DR-HANDOFF** for Phase 02. No product features.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 2 Gate C
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md)
- Phase 00 VERIFY pattern: `../../phase-00-foundation/scopes/scope-09-phase-verify/`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Planner work
- Thicken `01-verify.md` with concrete commands + evidence file (`VERIFY-NOTES.md`).
- Encode DR-HANDOFF: scaffold Phase 02 folder + board rows **before** final review marks phase complete.
- Fix board-rights: VERIFY **may spawn** remediations; S05-02 owns handoff completion check.
- Sync SCOPE-TODOS.md.

## Depends-on
- S01–S04 all `done` (or S04 skipped with reason — not expected).

## Exit criteria
- [x] `01-*` runnable without guessing
- [x] DR-HANDOFF duties explicit in 01 + 02
- [x] SCOPE-TODOS + TODO.md Notes updated
- [x] No product Go

## Minimal todos
- [x] Mirror P00 VERIFY structure for P01 bar — honesty/x0/p0x commands, evidence table, spawn convention, law checks
- [x] Thicken 01 + 02 handoff language — S05-01 creates `phase-02-gate-c` + `P02-00`; S05-02 owns completion check
- [x] Sync todos — SCOPE-TODOS + board Notes; **P01-S05-01 ready**
