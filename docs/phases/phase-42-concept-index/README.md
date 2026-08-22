# Phase 42+ — Concept & index

Human-promoted successor to Phase 41 close. **Complete — closed at `P42-S02-02` (2026-08-22).** G6+G7 delivered; **G1–G9 remediation complete.**

## Goal

- **G6:** Non-semantic concept retrieval — graph-label channel without vector (G-004b)
- **G7:** Index freshness & language coverage — analyzer/lang policy + index honesty (G-005)

## Evidence basis

- Phase 38 [REMEDIATION-PLAN.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §2 G6/G7
- Phase 41 [DR-HANDOFF.md](../phase-41-layers-intent/DR-HANDOFF.md) + [VERIFY-NOTES.md](../phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md)
- G8 progressive layers shipped (Phase 41 S00)
- G9 intent pipeline shipped (Phase 41 S01)

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | G6/G7 in/out, M-001 charter, rejects preserved |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P42-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P42-S02-02` — successor **no successor** |

Board: [`docs/TODO/phase-42.md`](../../TODO/phase-42.md).

## Scope sequence

```
S00 non-semantic concept retrieval (G6) → graph-label channel under DR-NOSSEM
 → S01 index freshness & langs (G7) → lang policy + index honesty
 → S02 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus |
|-------|-------|----------------|
| S00 | G6 | Graph-label retrieval channel; law review gate |
| S01 | G7 | Lang expansion policy; optional local watch/hook path |
| S02 | VERIFY | Blocks + successor Phase 43+ or `no successor` |

## In scope

(per [INTAKE.md](INTAKE.md) — planner thickens at P42-00)

- **S00 G6:** Non-semantic concept retrieval via graph labels/summaries (no vector)
- **S01 G7:** Index freshness ergonomics; language coverage policy

## Out of scope

- G-004a vector / embeddings
- Always-on daemon / public bind defaults
- Product dual-index default
- Rewriting Phase 41 delivery history

## Moat charter (M-001)

Concept retrieval and index changes merge into task loop + progressive caps — never replace moat or enable full-graph dump defaults.
