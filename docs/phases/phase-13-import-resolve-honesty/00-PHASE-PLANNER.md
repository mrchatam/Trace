# P13 / 00-PHASE-PLANNER — Import resolve & honesty residuals (thin) phase scaffold

## Metadata
- id: P13-00
- todo_ids: [P13-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **thin Phase 13** after Phase 12 closed with `no successor`, driven by human-approved promotion of post–P12 findings **DF-60…67**. Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + `PROJECT_DOCS_INDEX.md` + DOGFOOD-FINDINGS schedule. Optionally scaffold `experiments/ab-import-resolve/` (prepare + rubric only). **Do not** implement product Go in this row. **Do not** board deferred research impact/install/supersession unless already a DF.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 5–7, 12–13, 18–19; DR-INCREMENTAL
- [phase README](./README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md) — DF-60…67
- [experiments/POST-P12-DOGFOOD.md](../../../experiments/POST-P12-DOGFOOD.md)
- [experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md](../../../experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md)
- Phase 12 closeout: [DR-HANDOFF.md](../phase-12-peer-honesty-surfaces/DR-HANDOFF.md) — historical `no successor` (intact)
- Design: [RETRIEVAL_AND_CONTEXT.md](../../RETRIEVAL_AND_CONTEXT.md), Laws 5–7 / 18

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings/Neo4j | Forbidden as primary / SoT |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 12 `done` prompts; Phase 13 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 5 | Edge provenance honesty — resolve paths so live trees surface enum; do not fake precise calls |
| Laws 6–7 | Budgeted packets; no silent truncation / false-fresh |
| Law 18 | Causal STALE ≠ index honesty banner |
| Phase 12 product | Named honesty tests stay green; do not regress |

## Problem statement (post–P12)

Phase 12 named locks hold, but live agents hit:

1. **Import resolve dead** — Expand `GetFileByPath(NormalizePath(imp))` strips `./` only; no importer-dir join / extension normalize → subdir `./util` never neighbors → **`edge_provenance` invisible** (**DF-60**, high; already clear).
2. **Packet honesty quiet edges** — stale_paths ≤8 silent (**DF-61**); trim drops file → null honesty (**DF-62**); post-cap `items_total` undercount (**DF-63**); context FTS≠import hops (**DF-65**).
3. **Enum soft + residuals** — no provenance CHECK; empty omitempty / garbage (**DF-64**); INFERRED unreachable live (**DF-66**); symbol-entity staleness out (**DF-67**).

## Assumptions locked (grill deferred — human already cut scope)

| # | Assumption |
|---|------------|
| A1 | **Board DF-60…67** in four scopes (S01–S03 + VERIFY); no clarifying experiment required to board (DF-60 clear) |
| A2 | Scope planners pick exact APIs/migrations; phase planner locks **homes** + order + non-goals |
| A3 | S01 fix is **resolve-time** (importer dir + `./` + extensions) — not re-keying all analyzer strings to project-relative paths (planner may choose helper strategy) |
| A4 | S02 keeps false-fresh preference; Law 18 causal STALE untouched |
| A5 | DF-66/67 are product residuals in **S03** — scope planner may thin-fix or document wontfix with evidence; do not silently drop |
| A6 | Thin MCP parity only if library already exposes fields (G19); no new MCP tool menu |
| A7 | VERIFY default DR-HANDOFF = **`no successor`** |
| A8 | Optional `ab-import-resolve` is acceptance/dogfood hook — **not** a board blocker |

## Canonical open DF list (2026-08-17) → Phase 13 homes

| ID | Sev | Status | Phase home |
|----|-----|--------|------------|
| DF-60 | **high** | scheduled → P13 | **S01** |
| DF-61 | med | scheduled → P13 | **S02** |
| DF-62 | med | scheduled → P13 | **S02** |
| DF-63 | med | scheduled → P13 | **S02** |
| DF-64 | med | scheduled → P13 | **S03** |
| DF-65 | med | scheduled → P13 | **S02** |
| DF-66 | low | scheduled → P13 | **S03** |
| DF-67 | low | scheduled → P13 | **S03** |

**Count: 8.** IDs DF-52…59 remain unused gap. Next free: **DF-68**.

## Scope order (locked)
1. **S01 import-path-resolve** — DF-60  
2. **S02 packet-honesty-residuals** — DF-61, DF-62, DF-63, DF-65  
3. **S03 provenance-schema** — DF-64, DF-66, DF-67  
4. **S04 VERIFY** — S01–S03 regressions + carry-forward gates + DR-HANDOFF (+ optional ab-import-resolve prepare)

## Non-goals
- Product Go in this planner row  
- Boarding research ranks 4+ / impact / install / supersession (unless already a DF)  
- MCP/daemon/HTTP as P0; embeddings / Neo4j SoT  
- Full-rebuild indexer; rewriting Phase 00–12 `done` history  
- Requiring new clarifying experiments before boarding  
- Claiming Mode-B Gate C / new Gate C from this phase  

## Acceptance (this row)
- Phase folder + README + DR-HANDOFF stub  
- Scope stubs S01–S04 with 00/01/02 + SCOPE-TODOS  
- Board Phase 13 section; **P13-00 done**; next **P13-S01-00**  
- AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS scheduled → Phase 13  
- Optional: `experiments/ab-import-resolve/` prepare + rubric (subdir `./` imports)

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S04 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 13 section + P13-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS schedule note  
- [x] Optional ab-import-resolve scaffold (prepare + rubric)  
- [ ] Product Go — **not** this row  

## Next
**P13-S01-00** (scope planner for import-path-resolve / DF-60).
