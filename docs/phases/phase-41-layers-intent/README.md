# Phase 41+ — Layers & intent

**Complete — closed at `P41-S02-02` (2026-08-22).** G8+G9 delivered.

## Goal

- **G8:** Progressive layers L2–L3 — **ship** opt-in via `max_layer` (default L0–L1) (G-003)
- **G9:** Intent pipeline — **implement** bounded rule-based `ExtractIntent` (G-009)

## Evidence basis

- Phase 38 [REMEDIATION-PLAN.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §3 Phase 41+
- Phase 40 [DR-HANDOFF.md](../phase-40-read-surface-retrieval-depth/DR-HANDOFF.md) + [VERIFY-NOTES.md](../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md)
- G5 GUI orient shipped (Phase 40 S00)
- G2 unified `trace_explore` shipped (Phase 40 S01)

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | G8/G9 in/out, M-001 charter, secondary queue G6/G7 |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P41-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P41-S02-02` — successor Phase 42+ |

Board: [`docs/TODO/phase-41.md`](../../TODO/phase-41.md).

## Scope sequence

```
S00 progressive layers (G8) → compiler layer expansion or spec revise
 → S01 intent pipeline (G9) → implement or doc-revise aspirational §3
 → S02 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus |
|-------|-------|----------------|
| S00 | G8 | L2–L3 progressive layers in compiler or documented alternative |
| S01 | G9 | Intent extraction implement or RETRIEVAL_AND_CONTEXT §3 supersede |
| S02 | VERIFY | Blocks + successor Phase 42+ or `no successor` |

## In scope

(per [INTAKE.md](INTAKE.md) — planner thickens at P41-00)

- **S00 G8:** Ship L2–L3 in compiler or revise spec with documented alternative
- **S01 G9:** Implement intent extraction **or** mark doc aspirational + supersede

## Out of scope

- G-004a vector / embeddings
- G6/G7 primary implement (secondary queue — human may promote before G8/G9)
- Always-on daemon / public bind defaults
- Product dual-index default
- Rewriting Phase 40 delivery history

## Moat charter (M-001)

Layer expansion and intent pipeline merge into task loop + progressive caps — never replace moat or enable full-graph dump defaults.
