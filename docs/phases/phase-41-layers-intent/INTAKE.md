# INTAKE — Layers & intent

**P41-00 complete** (2026-08-22). Scopes locked; next runnable **P41-S00-00**.

## Trigger

Phase 40 delivered G5+G2 entry wave. [REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) ranks **G8 + G9** as Phase 41+ entry themes.

## Human locks

| Lock | Value |
|------|-------|
| Phase 40 | **CLOSED** — G5+G2 delivered; do not reopen implement rows |
| Phase 41+ mode | **Implement** G8+G9 entry |
| M-001 moat | Progressive task packet, loop+gate+review, harness, `trace_why`, plan tree, Laws 6–7 — **non-negotiable** |
| G-004a vector | **Out** — permanent defer |

## Entry themes

| Theme | GAP ids | Scope | Deliverable (P41-00 locked) |
|-------|---------|-------|----------------------------|
| **G8** Progressive layers L2–L3 | G-003 | S00 | **Ship** L2–L3 in compiler (opt-in `max_layer`; default L0–L1 unchanged) |
| **G9** Intent pipeline | G-009 | S01 | **Implement** bounded rule-based intent in `internal/retrieval/`; revise §3 (semantic deferred) |

## In scope (P41 locked)

- G8: Compiler layer expansion — L2/L3 admission, budget trim, CLI/MCP `--max-layer`
- G9: `ExtractIntent` rule-based stage wired into Search; §3 doc honest

## Out of scope (P41 locked)

- G-004a semantic/vector channel
- G6–G7 primary implement (secondary queue — see below)
- Hosted SaaS / daemon defaults
- Product dual-index default
- Rewriting Phase 40 artifacts
- LLM intent extraction

## Secondary queue (not P41 implement rows)

| Rank | Theme | GAP ids | Phase sketch | P41-00 decision |
|------|-------|---------|--------------|-----------------|
| 6 | **G6** Non-semantic concept retrieval | G-004b | Phase 42+ default | **Forward only** — no S03 row |
| 7 | **G7** Index freshness & langs | G-005 | Phase 42+ default | **Forward only** — no S04 row |

Human may promote G6 or G7 in a future phase planner row.

## Resolved decisions (P41-00)

| # | Question | Resolution |
|---|----------|------------|
| 1 | G9 implement vs doc-revise? | **Implement bounded rule-based** — deterministic `ExtractIntent`; doc-revise §3 only as S01 fallback |
| 2 | G6/G7 as secondary scopes in P41? | **Document only** in DR-HANDOFF — no S03/S04 board rows |
| 3 | G8 ship vs spec-revise? | **Ship** — opt-in L2/L3 via `max_layer`; spec-revise only if S00-01 blocked |

## Evidence pointers (read-only)

| Artifact | Path |
|----------|------|
| REMEDIATION-PLAN §2 G8/G9 | [§2](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) |
| GAP-REGISTRY G-003/G-009 | [GAP-REGISTRY.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md) |
| Phase 40 VERIFY-NOTES | [VERIFY-NOTES.md](../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md) |
| Layer design vs shipped | `experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md` |
| Intent pipeline gap | `experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md` |
| RETRIEVAL_AND_CONTEXT §3–§4 | `docs/RETRIEVAL_AND_CONTEXT.md` |
| G2 shipped | `internal/compiler/explore.go` |
| G5 shipped | `web/src/components/GraphOrientPanel.tsx` |

## Next runnable

**P41-S00-00** (G8 scope planner — prompts thickened at P41-00)
