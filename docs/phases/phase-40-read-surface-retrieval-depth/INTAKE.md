# INTAKE — Read surface & retrieval depth

**Scaffolded 2026-08-22** at Phase 39 close (P39-S03-02). **P40-00 complete 2026-08-22** — scopes locked.

## Trigger

Phase 39 delivered G1+G3+G4 entry co-wave. [REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) ranks **G5 + G2** as Phase 40+ entry themes.

## Human locks

| Lock | Value |
|------|-------|
| Phase 39 | **CLOSED** — G1+G3+G4 delivered; do not reopen implement rows |
| Phase 40+ mode | **Implement** G5+G2 entry |
| M-001 moat | Progressive task packet, loop+gate+review, harness, `trace_why`, plan tree, Laws 6–7 — **non-negotiable** |
| G5 | **Law 19 adapter** — GUI orient over canonical library/API |
| G-004a vector | **Out** — permanent defer |

## Entry themes

| Theme | GAP ids | Scope | Deliverable sketch |
|-------|---------|-------|-------------------|
| **G5** Graph-first onboarding UX | G-008 | S00 | Enhance `/` Graph route — orient panel + install narrative |
| **G2** Unified `trace_explore` | G-007 | S01 | Task-aware capped read; 17th MCP tool after G1 |

## In scope (P40 default)

- G5: GUI graph orient adapter on existing Explore route (orient panel + install hook docs)
- G2: Unified `trace_explore` MCP tool — task-aware, capped, merges into task loop

## Out of scope (P40 default)

- G-004a semantic/vector channel
- G6–G9 primary implement (secondary queue — **no S03/S04 rows added at P40-00**)
- Hosted SaaS / daemon defaults
- Product dual-index default
- Rewriting Phase 39 artifacts

## Evidence pointers (read-only)

| Artifact | Path |
|----------|------|
| REMEDIATION-PLAN §2 G5/G2 | [§2](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) |
| GAP-REGISTRY G-007/G-008 | [GAP-REGISTRY.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md) |
| Phase 39 VERIFY-NOTES | [VERIFY-NOTES.md](../phase-39-context-orient-harness/scopes/scope-03-verify/VERIFY-NOTES.md) |
| PEER-CG explore UX | [PEER-CG.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md) |
| G1 shipped | `internal/compiler/compiler.go:158–165` — `ContextOptions.Query` |
| G5 live GUI | `web/src/App.tsx:21–22`, `web/src/screens/Graph.tsx` |
| G2 compose baseline | `internal/mcp/instructions.go:13–16` — compose-first recipe |

## Open questions — resolved at P40-00

| # | Question | Resolution (P40-00) |
|---|----------|---------------------|
| 1 | G6/G7 as secondary scopes in P40? | **Document only** in DR-HANDOFF forward note — no S03/S04 board rows |
| 2 | G2 law spike gate — live MCP spike required? | **Waived** — G1 + compose-first shipped Phase 39; desk-check checklist at S01-00; implement at S01-01 |
| 3 | G5 static sketch vs full GUI route? | **Enhance existing `/` Graph route** — ReactFlow + seed compose already shipped; add orient panel + install narrative (not static-only, not Graphify port) |

## Next runnable

**P40-S00-00** (scope planner G5) — thicken already done at P40-00; implement wave starts **P40-S00-01**.
