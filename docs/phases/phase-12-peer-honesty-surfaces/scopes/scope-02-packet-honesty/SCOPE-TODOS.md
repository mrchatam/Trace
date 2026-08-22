# S02 — Packet honesty — scope todos (FINAL) — retry 2026-08-17

**Depends-on:** P12-S01-02 done (APPROVE). Owns research **ranks 2–3** — emission-time staleness banners + loud truncation/totals.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** (retry) — FINAL locks + live inventory 2026-08-17 |
| 2 | 01-packet-honesty | implement | pending — **gap finish** (named tests + residual fixes; types mostly shipped) |
| 3 | 02-scope-review | review | pending |

## FINAL summary (from P12-S02-00 retry)
- **No migration**; Packet `SchemaVersion` **`0.2`** — **already in tree**
- Budget: `items_total` / `items_kept` / `candidates_capped` + `truncated` when kept&lt;total **or** candidate cap — **already wired**
- MD: loud `items=kept/total` (+ candidates_capped) — **already present**
- `index_honesty.stale_paths` via emission-time sha256(disk) vs `files.content_hash` on kept `file` items; **false-fresh** — **helper + wire present**; prefer **sort-then-cap 8**
- Do **not** conflate with Law 18 causal STALE; preserve S01 `edge_provenance`
- **Remaining bar:** named tests `TestBudgetLoudTotals`, `TestCandidateCapSetsTruncated`, `TestIndexStaleBanner` (+ S01 `TestContextWhyTraceEdgeProvenance`)
- Skeletonization/dedup optional — not bar
- Implementer: **do not blindly re-implement** shipped types

## Reminders
- Carry-forward + S01 stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P12-S03-00**
