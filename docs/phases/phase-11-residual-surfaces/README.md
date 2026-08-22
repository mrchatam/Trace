# Phase 11 — Residual surfaces (post–P10 open findings)

**Status:** **complete** (2026-08-16) — S01–S07 APPROVE; S08 VERIFY PASS + review APPROVE; **DR-HANDOFF = `no successor`** historically (thin Phase 12 later human-scheduled forward).

Reopened the board after Phase 10 (`no successor`) to schedule **all** still-open / ops / deferred dogfood findings — severity-agnostic (user promotion).

## Why this phase exists

Phase 10 closed DF-17…32 in binary (VERIFY APPROVE). Post–P10 dogfood, MCP parity, and adversarial hunts left **18 canonical open DFs** (partial-index ghosts, PASS+FAIL DONE, store lock, capability slug upsert / hatch vs caps, retrieval polish, MCP reload UX, seed/plan/review polish, deferred handoff design). Findings SoT: [`experiments/DOGFOOD-FINDINGS.md`](../../../experiments/DOGFOOD-FINDINGS.md).

## Scope order (locked at P11-00)

| Scope | Focus |
|-------|--------|
| S01 | **Index partial-path GC** — DF-40 (rename + `index <new-path>` ghosts; DF-20 residual) |
| S02 | **Review PASS+FAIL / operator identity** — DF-43 (sibling FAIL ignored); DF-44 (`--as-operator` ≠identity) |
| S03 | **Store lock / concurrency** — DF-47 (CLI↔CLI / CLI↔MCP exclusive lock UX) |
| S04 | **Capability upsert + hatch vs caps** — DF-41 (slug upsert); DF-51 (`--allow-done` vs missing-caps) |
| S05 | **Retrieval: why symbol, depth leak, trust, DPC attribution** — DF-49, DF-35, DF-48, DF-42 |
| S06 | **MCP / install reload UX** — DF-22, DF-37, DF-50 |
| S07 | **Seed / plan / review show polish** — DF-28, DF-30, DF-33, DF-45, DF-46 |
| S08 | Phase VERIFY + DR-HANDOFF |

## Out of scope (unless promoted)

- Daemon / HTTP primary surface / embeddings  
- Full MCP surface for plan / impact / index (CLI stays primary unless a scope planner promotes thin tools)  
- Rewriting Phase 00–10 `done` history  
- Experiment method notes only: DF-06/07/13/34/36 (not product board scopes)  
- Full-rebuild-on-any-change indexer architecture  

## Parallel track (not board-blocking)

Continue ladder / natural dogfood under `experiments/` — score into `RESULTS.md`; feed new DF-* **forward** only (no ID reuse; next free **DF-52**).
