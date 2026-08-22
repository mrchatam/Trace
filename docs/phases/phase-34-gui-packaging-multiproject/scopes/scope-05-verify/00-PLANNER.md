# P34-S05-00 — Scope planner (VERIFY)

## Metadata
- id: P34-S05-00
- todo_ids: [P34-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock VERIFY floor for Phase 34: consumer-like temp (`.trace/` only) + `trace gui` → **real SPA not stub**; second concurrent → **different port**; docs aligned. Evidence under `experiments/runs/`. Leave DR-HANDOFF OPEN until S05-02.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- Prior PASS reviews S00–S04

## Session start

Follow agent-loop-protocol. Do not invent successor phase here.

## Locked defaults

| Item | Value |
|------|-------|
| Floor | (1) consumer temp no `web/` → real SPA; (2) concurrent second gui → free port + correct URL; (3) docs no consumer `web/` requirement |
| S01-00 seed (prefer PLAN) | VERIFY = PLAN **T9**: combine T1 (real SPA markers `#root`/`/assets/`, not stub phrase) + T4/T5 concurrent auto-port + T8 docs; fail if stub shipped when full SPA intended |
| T8 (S04-02) | **PASS** 2026-08-21 — quickstart/`web`/AGENTS/help/embeddist free of “no auto free-port” / consumer two-artifact primary; positive embed + `7432`–`7441` + `--addr` pin. Residual: keep `docs/TODO.md` + `AGENTS.md` next-runnable current (fixed to **P34-S05-00** in S04-02). |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p34-s05-01-verify/evidence/` |
| Notes | `scopes/scope-05-verify/VERIFY-NOTES.md` required |
| DR-HANDOFF | Stays OPEN until **P34-S05-02** |
| Successor lean | Default **no successor** (S05-02 decides) |
| Product code | **Forbidden** on VERIFY row |

## Fail vs residual

**Fail VERIFY for:** stub SPA on consumer path when release embed should be full; no auto-port on default busy; docs still require consumer `web/`; public bind default; SPA copied into consumer `.trace/` as primary.

**Do not fail solely for:** cosmetic help nits; contributor Trace-checkout `web/` DX still documented if labeled.

## Planner gate

- [x] `01-verify.md` + `02-dr-handoff.md` + `SCOPE-TODOS.md` ready

## Exit criteria

- [x] Next **P34-S05-01**

## Todo updates

Status + notes on **P34-S05-00** only.

## Next

`P34-S05-01`
