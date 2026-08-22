# P01 / S03 / 00-PLANNER — Agent X0 harness (CLI)

## Metadata
- id: P01-S03-00
- todo_ids: [P01-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-x0-harness.md` for **Agent X0 harness (CLI)** against live CLI + `fixtures/x0` + I_BENCHMARK_PLAN X0. Lock B0/G1, metrics schema path, dry-run bar. No product domain redesign.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 1 validation gate
- Live: `evals/p0x`, `fixtures/x0`, `cmd/trace` why/context
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Planner work
- Confirm B0 vs G1 invocation shape (scripted CLI for G1; repo tools for B0).
- Lock `evals/x0` layout + `schema.json` (I_BENCHMARK_PLAN).
- Dry-run = metrics emit for both conditions; **not** Gate C scoring.
- Thicken `01-x0-harness.md`; note keep `evals/p0x` green.
- Sync SCOPE-TODOS.md.

## Depends-on
- S01–S02 done preferred for honesty coverage; **hard** dependency is working CLI context/why (Phase 00) + fixtures/x0. Do not start until S02 review done (board order).
- **S01 surface (P01-S01-00):** DONE gate is Review PASS or explicit `allow_done`; if X0 seed ever marks DONE, plan for that — default dry-run should avoid DONE.
- **S02 note (P01-S02-00):** H5 honesty lives in **`evals/honesty`** (fail-closed domain demo). Keep **`evals/x0`** as B0/G1 harness only — do not merge honesty into X0 metrics schema this phase; keep `evals/p0x` untouched.

## Exit criteria
- [x] `01-*` runnable without guessing (paths, conditions, metrics)
- [x] SCOPE-TODOS + TODO.md Notes updated
- [x] No product Go in this planner row

## Minimal todos
- [x] Inspect p0x harness patterns + I_BENCHMARK_PLAN X0
- [x] Thicken 01 + 02 prompts; light S04 Depends notes
- [x] Sync SCOPE-TODOS / board Notes
