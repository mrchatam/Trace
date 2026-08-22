# S01 — Index partial-path GC — scope todos

**Depends-on:** P11-00 done. Owns **DF-40** (DF-20 residual: rename + partial `index <new-path>` ghosts).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL content-hash orphan GC locks |
| 2 | 01-index-partial-path-gc | implement | pending — ready for **P11-S01-01** |
| 3 | 02-scope-review | review | pending |

## Locked summary (from 00-PLANNER)
- Partial argv after successful index: delete other DB paths with **same `content_hash`** that are **missing on disk**
- Full-tree set-diff + missing-argv single delete **unchanged** (P10 DF-20)
- No project-wide GC on argv; no mig; no MCP index; no full-rebuild
- Required test: `TestIndexPartialArgvGCAfterRename` (+ keep DF-20/isolation)

## Reminders
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Residual OK: rename+edit (hash diverges) may need full-tree index
- Next after APPROVE: **P11-S02-00**
