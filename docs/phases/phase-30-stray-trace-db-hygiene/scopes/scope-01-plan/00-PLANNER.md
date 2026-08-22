# P30-S01-00 — Scope planner (plan)

## Metadata
- id: P30-S01-00
- todo_ids: [P30-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective

Against **live** `INVESTIGATION.md`, lock defaults and thicken `01-plan.md` so S01-01 writes `PLAN.md`. **No product code.**

## Depends

S00 review **PASS** (`P30-S00-02`, 2026-08-21): verdict **agent hygiene**; INTAKE confirmed; independent stub repro 0-byte; no Trace dual-store bug.

## Locked plan track (S00 confirmed — hygiene only)

| ID | Change | Risk |
|----|--------|------|
| T1 | Install/AGENTS/docs: only `.trace/trace.db`; never open/create root `trace.db` | low |
| T2 | Optional warn once on store open if `<root>/trace.db` exists (stderr; non-fatal) | low |
| T3 | Scaffold/gitignore recommend ignoring root `trace.db` | low |
| T4 | Tests: warn fires when stub present; Trace still uses `.trace/` | low |

Do **not** auto-delete the stray file by default. Do **not** redesign store path.

## Overturn path

**Closed for S01.** S00 did **not** find a Trace bug. Do not invent a path-fix task. (If later evidence appears, that is a new phase/scope — not this PLAN.)

## Hard constraints (carry forward)

- Canonical store remains `<root>/.trace/trace.db` — **locked**; no store-path change in S01/S02
- Cite [`../scope-00-investigate/INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md) recommendations § as non-binding input to PLAN.md
- No silent delete without documented flag
- No HTTP/daemon scope creep

## Exit criteria

- [x] `01-plan.md` runnable with PLAN.md template + acceptance hooks
- [x] Board Notes cite INVESTIGATION verdict used
- [x] Next: **P30-S01-01**

## Next

`P30-S01-01`
