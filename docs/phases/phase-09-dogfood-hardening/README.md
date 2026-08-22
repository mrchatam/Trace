# Phase 09 — Dogfood hardening & agent UX

Reopens the board after Phase 08 (`no successor`) using live Cursor dogfood evidence.

## Why this phase exists

Six A/B rungs (D01/D02/D03/D05/D07/D10) showed Trace helping agents — and exposed product gaps (especially **why/context after review**). Findings SoT: [`experiments/DOGFOOD-FINDINGS.md`](../../experiments/DOGFOOD-FINDINGS.md).

## Scope order (locked at P09-00)

| Scope | Focus |
|-------|--------|
| S01 | **Retrieval completeness** — fix DF-01 (`review` in ExactLookup / why+context); regression from honesty-shaped fixture |
| S02 | **Discoverability** — DF-02/DF-04: list tasks / work status; seed-path ergonomics |
| S03 | **Install & editor wire** — DF-03/DF-05: Cursor MCP install path (CLI or documented one-shot) |
| S04 | Phase VERIFY + DR-HANDOFF — locked at P09-S04-00 as **`no successor`** (ladder gaps stay parallel `experiments/`) |

**Status (2026-08-16):** Phase **complete**. VERIFY PASS; DR-HANDOFF closed on `P09-S04-02` = **`no successor`** (historical). **Phase 10** later reopened forward from deduped DF-17+ — [`phase-10-integrity-surfaces/`](../phase-10-integrity-surfaces/). Evidence: [`scopes/scope-04-phase-verify/VERIFY-NOTES.md`](scopes/scope-04-phase-verify/VERIFY-NOTES.md), [`scopes/scope-04-phase-verify/REVIEW-NOTES.md`](scopes/scope-04-phase-verify/REVIEW-NOTES.md).

## Out of scope (unless promoted)

- Daemon / HTTP primary surface  
- Embeddings  
- Rewriting Phase 00–08 `done` history  
- Replacing `experiments/ab-*` with harness-only (dogfood continues in parallel)

## Parallel track (not board-blocking)

Continue ladder batch: D04 / D06 / D11 under `experiments/` — score into `experiments/RESULTS.md`; feed new DF-* rows forward.
