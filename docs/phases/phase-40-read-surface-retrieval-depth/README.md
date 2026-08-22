# Phase 40+ — Read surface & retrieval depth

Human-promoted successor to Phase 39 close. **Complete — closed at `P40-S02-02` (2026-08-22).** G5+G2 delivered.

## Goal

- **G5:** Graph-first onboarding UX — GUI orient adapter (G-008)
- **G2:** Unified `trace_explore` — task-aware capped read (G-007)

## Evidence basis

- Phase 38 [REMEDIATION-PLAN.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §3 Phase 40+
- Phase 39 [DR-HANDOFF.md](../phase-39-context-orient-harness/DR-HANDOFF.md) + [VERIFY-NOTES.md](../phase-39-context-orient-harness/scopes/scope-03-verify/VERIFY-NOTES.md)
- G1 query merge shipped (Phase 39 S00)
- G3 compose-first recipe shipped (Phase 39 S01)

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | G5/G2 in/out, M-001 charter, P39 links |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P40-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P40-S02-02` — successor Phase 41+ |

Board: [`docs/TODO/phase-40.md`](../../TODO/phase-40.md).

## Scope sequence

```
S00 GUI graph orient (G5) → adapter over canonical API
 → S01 unified explore (G2) → product code after law spike
 → S02 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus |
|-------|-------|----------------|
| S00 | G5 | GUI `/` graph route orient; Law 19 adapter only |
| S01 | G2 | Task-aware capped `trace_explore` MCP tool |
| S02 | VERIFY | Blocks + successor Phase 41+ queue |

## In scope

(per [INTAKE.md](INTAKE.md) — locked at P40-00)

- **S00 G5:** Orient panel on `/` Graph + install hook narrative (Law 19 adapter)
- **S01 G2:** Library explore compose + MCP `trace_explore` (17th tool, task_id required)

## Out of scope

- G-004a vector / embeddings
- Always-on daemon / public bind defaults
- Product dual-index default
- Rewriting Phase 39 delivery history

## Moat charter (M-001)

Unified explore merges into task loop — never query-only replacement.
GUI orient is Law 19 adapter over canonical library/API.
