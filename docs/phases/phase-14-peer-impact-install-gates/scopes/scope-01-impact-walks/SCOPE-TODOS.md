# S01 — Impact walks — scope todos

**Depends-on:** P14-00 done. Owns research **rank 6** — multi-seed impact BFS + contains asymmetry.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL 2026-08-17 |
| 2 | 01-impact-walks | implement | **pending** — next runnable **P14-S01-01** |
| 3 | 02-scope-review | review | pending |

## FINAL locks (summary — see 00-PLANNER)
- Home: `internal/retrieval` ImpactWalk + `trace impact walk`; Expand untouched
- Multi-seed one BFS; seeds excluded; depth 1..2 (default 2)
- Incoming imports only; contains-OUT from files; no sibling climb from symbols
- Loud cap 64: `blast_total` / `blast_kept` / `truncated`; `hop_risk=float64(hop)`
- No mig / MCP / `internal/impact`; Gate F planted harness kept

## Reminders
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- No S02 install work; no S05 supersession; no `plan simulate`
- Do not re-open DF-60…67
- Next after APPROVE: **P14-S02-00**
