# S01 — Retrieval / why fidelity — scope todos

**Depends-on:** P10-00 done. Owns DF-19, DF-23, DF-25, DF-27, DF-29.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16; 01+02 thickened |
| 2 | 01-retrieval-why-fidelity | implement | **done** — DF-19/23/25/27/29 |
| 3 | 02-scope-review | review | **done** — APPROVE high; [REVIEW-NOTES.md](REVIEW-NOTES.md); no spawns |

## Locked reminders (from P10-S01-00)
- DF-19: goal-scoped + pair-completion; multi-goal must not leak; single-goal unattributed fallback for x0
- DF-23: canonical `plan_change`; accept `plan-change` at why/Exact
- DF-25: Exact/Why `capability`; **`plan_scope` residual** (do not implement)
- DF-27: reword decision MD; keep `untrusted_data`; no system elevation (Law 4/9)
- DF-29: IncludeWhy fail-closed
- No new mig; no new MCP tools (S02); no daemon/HTTP/embeddings
- Carry-forward: honesty G, E/F/ablation/H/compat, p0x, x0, Gate C `dry_run:false`

## Next after S01
**P10-S02-00** (MCP parity + install) — inherits why/Exact alias + IncludeWhy behavior.
