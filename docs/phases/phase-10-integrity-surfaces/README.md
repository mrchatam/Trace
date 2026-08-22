# Phase 10 — Integrity surfaces (post-dogfood)

**Status:** complete (2026-08-16). DR-HANDOFF = **`no successor`** (historical). **Phase 11** later reopened forward from post-P10 residuals — [`../phase-11-residual-hardening/`](../phase-11-residual-hardening/).

Reopened the board after Phase 09 (`no successor`) using **deduped** dogfood + bug-hunt findings.

## Why this phase exists

Phase 09 closed DF-01…05 (why/context after review, `trace tasks`, install wire). Parallel dogfood + an adversarial hunt then exposed **integrity** gaps: retrieval pollution, sticky PASS / operator DONE, index ghosts, MCP tool parity, capability gating. Findings SoT: [`experiments/DOGFOOD-FINDINGS.md`](../../../experiments/DOGFOOD-FINDINGS.md) (canonical **DF-17+** after 2026-08-16 collision map).

## Scope order (locked at P10-00)

| Scope | Focus |
|-------|--------|
| S01 | **Retrieval / why fidelity** — DF-19 global DPC attach; DF-23 `plan-change` vocab; DF-25 capability Exact; DF-27 decision trust labeling; DF-29 IncludeWhy errors |
| S02 | **MCP parity + install freshness** — DF-21 thin `trace_tasks` + capability tools (G19); DF-22 reload/pin after install; DF-32 JSON shape consistency |
| S03 | **Index GC** — DF-20 ghost files/symbols after rename (incremental delete-on-missing) |
| S04 | **Operator + capability gates** — DF-17 actor/operator DONE; DF-18 sticky PASS invalidate; DF-24 transition vs missing caps; DF-26 `--allow-done` UX; DF-31 capability missing usage |
| S05 | Phase VERIFY + DR-HANDOFF — successor decided at VERIFY (default assumption: **`no successor`** unless Notes promote) |

## Out of scope (unless promoted)

- Daemon / HTTP primary surface / embeddings  
- Full MCP surface for plan / impact / index (CLI stays primary; **thin tasks + capability only** this phase)  
- Handoff SoT product (DF-28) + planner depth leak (DF-35) — deferred design  
- DF-30 plan-show empty / DF-33 seed `from_id` aliases / DF-34…36 experiment method notes  
- Rewriting Phase 00–09 `done` history  

## Parallel track (not board-blocking)

Continue ladder / natural dogfood under `experiments/` — score into `RESULTS.md`; feed new DF-* **forward** only (no ID reuse).
