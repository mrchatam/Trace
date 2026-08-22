# Scope 02 — board map

**S02 implement** — per [PLAN.md](../scope-01-plan/PLAN.md) accept set. Serial: **P37-S02-00 → P37-S02-01 → P37-S02-02**.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 640 | P37-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Thicken implement/review from PLAN.md |
| 641 | P37-S02-01 | [01-implement.md](01-implement.md) | Implementer | Code + tests per PLAN waves A–D |
| 642 | P37-S02-02 | [02-review.md](02-review.md) | Reviewer | Independent review + test re-run |

## Inputs (must exist)

- [`PLAN.md`](../scope-01-plan/PLAN.md) — accept set, touch-list, §5 tests, §8 waves

## Wave order (locked for S02-01)

| Wave | Residuals | Layer |
|------|-----------|-------|
| **A** | R5 + R1 | `advisories[]` on loop status |
| **B** | R3 + R2 + R4 | MCP gate, HTTP POST bootstrap, help |
| **C** | R6 + R11 | enforce test, agent workflow doc |
| **D** | R8 | Overview minimal advisories surface |

## Preserve (Phase 36)

- MCP `trace_plan` (16 tools)
- CLI + MCP bootstrap
- Terminal gate honesty (`goal_plan_gap_terminal_advisory`)
- Active-work `plan_missing` block
- P36 regression subset — see `01-implement.md` §Test strategy

## Hard constraints

- R1: advisory-only `bootstrap_recommended` — **never** `PlanExists=true`
- R7: do **not** change enforce default
- Law 19: HTTP/MCP/GUI thin adapters

## Out of scope

- RESIDUALS defer/reject rows (R7, R9, R8-full) unless PLAN explicitly reopened via spawn
- R10 live GUI verify — S03
