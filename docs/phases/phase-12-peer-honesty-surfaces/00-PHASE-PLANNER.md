# P12 / 00-PHASE-PLANNER — Peer honesty surfaces (thin) phase scaffold

## Metadata
- id: P12-00
- todo_ids: [P12-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **thin Phase 12** after Phase 11 closed with `no successor`, driven by human-approved cut of [`docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md`](../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md): **S01 edge provenance + S02 packet honesty only**, then standard VERIFY. Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + `PROJECT_DOCS_INDEX.md` + research schedule note. **Do not** implement product Go in this row. **Do not** board research S03–S06 (impact / install / supersession / full peer VERIFY).

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 5–7, 12–13, 18–19; DR-INCREMENTAL
- [phase README](./README.md)
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — ranks 1–3 accept; §D rejects; §E FUTURE cut
- [SIMILAR-PROJECTS-REVIEW-PROMPT.md](../../research/SIMILAR-PROJECTS-REVIEW-PROMPT.md)
- Phase 11 closeout: [DR-HANDOFF.md](../phase-11-residual-surfaces/DR-HANDOFF.md) — historical `no successor` (intact)
- Design: [RETRIEVAL_AND_CONTEXT.md](../../RETRIEVAL_AND_CONTEXT.md), [EVALUATION.md](../../EVALUATION.md) H1/H5/H6

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings/Neo4j | Forbidden as primary / SoT |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 11 `done` prompts; Phase 12 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 5 | Confidence / provenance on edges — do not fake precise calls |
| Laws 6–7 | Budgeted packets; no silent truncation |
| Law 18 | STALE honesty — complements emission-time staleness banners |
| Phase 11 product | 18 DFs stay fixed; do not regress |

## Problem statement

Phases 00–11 closed residual DF honesty gates, but peer methods still show:

1. **Structural edge provenance missing** — code-graph edges (imports / calls / usages) are not tagged `EXTRACTED|INFERRED|AMBIGUOUS` and Why/context do not surface hop confidence (research rank 1 / graphify).
2. **Packet under-honesty** — `trace context` / compiler packets lack emission-time staleness banners and loud truncation + exact totals (ranks 2–3 / codegraph + codebase-memory-mcp).

Thin Phase 12 closes those two surfaces only.

## Assumptions locked (grill deferred — human already cut scope)

| # | Assumption |
|---|------------|
| A1 | **Board only S01 + S02 + VERIFY (S03)** — research FUTURE S03–S05 (impact / install / supersession) and full S06 peer VERIFY **not** boarded |
| A2 | Scope planners pick exact APIs/migrations; phase planner locks **homes** + order + non-goals + research citations |
| A3 | Prefer false-fresh over false-stale for staleness banners (codegraph method); Law 18 STALE remains authoritative for causal provenance |
| A4 | Edge enum applies to **structural** code edges first; causal entity `confidence` fields already exist — do not conflate without S01 planner lock |
| A5 | Skeletonization / session dedup are **optional** S02 stretch — not acceptance bar |
| A6 | Thin MCP parity only if Why/context already expose fields via library (G19); no new MCP tool menu |
| A7 | VERIFY default DR-HANDOFF = **`no successor`** |

## Research citations (in-scope steals)

| Rank | Technique | Source | Trace home |
|-----:|-----------|--------|------------|
| 1 | Edge confidence `EXTRACTED\|INFERRED\|AMBIGUOUS` + Why/context display | graphify | **S01** |
| 2 | Emission-time staleness / pending-index honesty banners | codegraph | **S02** |
| 3 | Loud truncation + totals in context packets | codebase-memory-mcp (+ codegraph budgets) | **S02** |

## Explicit deferrals (not this phase)

| Research item | Why deferred |
|---------------|--------------|
| Ranks 4–20 (install matrix, impact BFS, surgical explore composition, supersession, RRF, …) | Human thin-slice; promote later |
| §D rejects (MCP/daemon/HTTP P0, Neo4j, embeddings default, Hybrid LSP rewrite, …) | Hard boundaries |

## Scope order (locked)
1. **S01 edge-provenance** — enum + store/analyzers + Why/context surfacing  
2. **S02 packet-honesty** — staleness banners + loud truncation/totals  
3. **S03 VERIFY** — S01/S02 regressions + carry-forward gates + DR-HANDOFF  

## Non-goals
- Product Go in this planner row  
- Boarding research S03–S06 / ranks 4+  
- MCP/daemon/HTTP as P0 architecture; embeddings / Neo4j SoT  
- Full-rebuild indexer; rewriting Phase 00–11 `done` history  
- Claiming Mode-B Gate C / new Gate C from this phase  

## Acceptance (this row)
- Phase folder + README + DR-HANDOFF stub  
- Scope stubs S01–S03 with 00/01/02 + SCOPE-TODOS  
- Board Phase 12 section; **P12-00 done**; next **P12-S01-00**  
- AGENTS.md + PROJECT_DOCS_INDEX + research doc human-schedule note  

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S03 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 12 section + P12-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + research schedule note  
- [ ] Product Go — **not** this row  

## Next
**P12-S01-00** (scope planner for edge-provenance).
