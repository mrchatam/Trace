# Phase 13 — Import resolve & honesty residuals (thin)

**Status:** **complete** (2026-08-17) — S01–S03 APPROVE; VERIFY APPROVE; **DR-HANDOFF = `no successor`** (historical at close). Human-scheduled forward after Phase 12 (`no successor`) to close post–P12 dogfood findings **DF-60…67**. Later **Phase 14** human-scheduled from goals-gap #1 (ranks 4–6); parallel dogfood / ranks 7+ stay off-board.

## Why this phase exists

Phase 12 shipped peer-honesty surfaces (edge provenance + packet honesty) and closed with **DR-HANDOFF = `no successor`**. Live adversarial + natural dogfood then showed: subdirectory relative imports never resolve → **`edge_provenance` invisible** on real trees (**DF-60**, high); packet honesty quiet edges (**DF-61…63**, **DF-65**); provenance enum soft (**DF-64**); low residuals **DF-66/67**. Findings SoT: [`experiments/DOGFOOD-FINDINGS.md`](../../../experiments/DOGFOOD-FINDINGS.md). DF-60 is already clear — no clarifying experiment required to board.

## Scope order (locked at P13-00)

| Scope | Focus |
|-------|--------|
| S01 | **Import path resolve** — DF-60 (join importer dir, normalize `./`, extensions; live `edge_provenance` on subdir trees) |
| S02 | **Packet / index honesty residuals** — DF-61 stale totals, DF-62 trim silent stale, DF-63 post-cap undercount, DF-65 context vs hops |
| S03 | **Provenance schema / enum** — DF-64 (+ product DF-66 / DF-67 as planner locks) |
| S04 | Phase VERIFY (+ optional [`experiments/ab-import-resolve/`](../../../experiments/ab-import-resolve/) dogfood hook) |

## Out of scope (unless promoted later)

- Deferred research impact / install / supersession (ranks 4+) unless already a DF  
- Daemon / HTTP primary surface / embeddings / Neo4j  
- MCP/daemon as P0 architecture; full MCP surface dump  
- Full-rebuild-on-any-change indexer  
- Rewriting Phase 00–12 `done` history  
- Clarifying experiments that block boarding (DF-60 already high/clear)

## Parallel track (not board-blocking)

Optional isolation dogfood [`experiments/ab-import-resolve/`](../../../experiments/ab-import-resolve/) — natural G1 for subdir `./` imports (fails today / passes after S01). Continue other ladder work under `experiments/`; feed new DF-* **forward** only (next free **DF-68**).
