# Phase 37 — Phase 36 residuals closure

**Phase planner.** Row `P37-00`.

## Metadata
- id: P37-00
- todo_ids: [P37-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation]
- verification: automated

## Mission

Triage and close Phase 36 residuals (R1–R11). Read [`INTAKE.md`](INTAKE.md) + [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) + Phase 36 DR-HANDOFF.

## Gate

Phase 36 **done**. Proceed.

## Scope sequence

| Scope | Rows | Artifact |
|-------|------|----------|
| S00 | P37-S00-00 → 02 | `RESIDUALS.md` |
| S01 | P37-S01-00 → 01 | `PLAN.md` |
| S02 | P37-S02-00 → 02 | Implementation |
| S03 | P37-S03-00 → 02 | VERIFY + CLOSED handoff |

## Hard constraints

- Preserve Phase 36 guarantees
- No silent PlanExists bridge
- Law 19

## Planner gate (P37-00)

- [x] Phase folder + INTAKE R1–R11 + locks
- [x] S00–S03 prompt stubs + SCOPE-TODOS
- [x] Board `docs/TODO/phase-37.md` + index link
- [x] DR-HANDOFF OPEN
- [x] No product code in this row

## P37-00 outcome (2026-08-22)

Gate **PASS**. Verified INTAKE R1–R11 against P36 DR-HANDOFF + VERIFY-NOTES § Residuals + PLAN §2.4/§2.6. README thickened with residual code-anchor table. Protocol-thickened S00–S03 prompts (References, Session start, Locked defaults, Exit criteria) + fixed SCOPE-TODOS board IDs (635–645). Spot-checked live repo: no `advisories[]` on `StatusResult`; `GoalStructureWarning` exists but unwired; `trace_loop` no gate; HTTP `GET /v1/plans` only; `WarnIfTraceDirWithoutConfig` untested. DR-HANDOFF OPEN with scope checklist. No product code. Next: **P37-S00-00**.

## Next

`P37-S00-00`
