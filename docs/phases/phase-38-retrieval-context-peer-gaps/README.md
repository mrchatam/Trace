# Phase 38 — Retrieval & context peer-gap investigation

Human-promoted **2026-08-22**. **Complete — closed at `P38-S07-02`.** Investigate → saturate → plan only. No implement.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | Hypotheses H1–H11 (verify, not implement) |
| [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) | Saturation gate + no-implement law |
| [`PEER-FIXTURES.md`](PEER-FIXTURES.md) | Codegraph, UA, Graphify, tools |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | `P38-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P38-S07-02` — successor Phase 39 |

Board: [`docs/TODO/phase-38.md`](../../TODO/phase-38.md).

## Scope sequence (investigation loops → plan)

```
S00 investigation index (hypothesis register + peer map)
 → S01 Trace live audit
 → S02 Codegraph peer deep-dive
 → S03 UA + Graphify peer deep-dive
 → S04 cross-matrix → GAP-REGISTRY.md
 → S05 saturation gate → SATURATION-NOTES.md  (exit loops here or spawn back)
 → S06 remediation PLAN only → REMEDIATION-PLAN.md
 → S07 VERIFY + DR-HANDOFF (successor = implementation phase TBD)
```

| Scope | Artifact | Investigate focus |
|-------|----------|-------------------|
| S00 | `INVESTIGATION-INDEX.md` | H* register, methods, spawn rules |
| S01 | `TRACE-AUDIT.md` | context, search, compiler, MCP, index, GUI |
| S02 | `PEER-CG.md` | explore tool, watch, benchmarks |
| S03 | `PEER-UA-GF.md` | query context, graph build, graph.html |
| S04 | `GAP-REGISTRY.md` | matrix Trace vs peers; moat row |
| S05 | `SATURATION-NOTES.md` | confident exit or spawn list |
| S06 | `REMEDIATION-PLAN.md` | ranked future phases — **no code** |
| S07 | `VERIFY-NOTES.md` + CLOSED | |

## Hard constraints

- **No product code** in P38 (evidence files OK)
- Saturation **before** REMEDIATION-PLAN
- Spawn investigation rows freely; do not skip to planning early
