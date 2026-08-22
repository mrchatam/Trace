# Phase 39 — Context orient & harness

Human-promoted successor to Phase 38 close. **Complete — closed at `P39-S03-02` (2026-08-22).** G1+G3+G4 delivered.

## Goal

- **G1:** Query+task orient merge (G-001, G-002)
- **G3:** MCP/harness orient playbook (G-006, G-010)
- **G4:** Dual-stack documentation — doc-only (G-011)

## Evidence basis

- Phase 38 [REMEDIATION-PLAN.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- Phase 38 [GAP-REGISTRY.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Phase 38 [SATURATION-NOTES.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-05-saturation-gate/SATURATION-NOTES.md)
- Phase 38 [DR-HANDOFF.md](../phase-38-retrieval-context-peer-gaps/DR-HANDOFF.md)

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | G1/G3/G4 in/out, M-001 charter, P38 links |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P39-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P39-S03-02` — successor Phase 40+ |

Board: [`docs/TODO/phase-39.md`](../../TODO/phase-39.md).

## Scope sequence

```
S00 context orient merge (G1) → product code
 → S01 harness orient (G3) → product + docs
 → S02 dual-stack docs (G4) → doc-only
 → S03 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus |
|-------|-------|----------------|
| S00 | G1 | Optional `query` on context/compiler; merge search hits into packet |
| S01 | G3 | MCP playbook, moat-first bootstrap, 9/16 hygiene |
| S02 | G4 | CONTRIBUTING/AGENTS Trace+Codegraph recipe |
| S03 | VERIFY | Blocks + successor Phase 40+ queue |

## In scope

(per [INTAKE.md](INTAKE.md) — planner thickens at P39-00)

- G1 query+task context merge preserving task UUID, gates, reason_codes, Laws 6–7
- G3 harness orient without copying MP 44-tool surface or hiding write tools
- G4 doc-only dual-stack recipe (`.trace/` vs `.codegraph/`, Law 19)

## Out of scope

- G2 unified `trace_explore` (Phase 40+)
- G-004a vector / embeddings
- Rewriting Phase 38 investigation history
- Product dual-index default or bundled MCP
- G5 GUI orient (secondary — Phase 39–40 if capacity)

## Moat charter (M-001)

All work merges peer patterns **into** the task loop + gates — never replaces moat.
