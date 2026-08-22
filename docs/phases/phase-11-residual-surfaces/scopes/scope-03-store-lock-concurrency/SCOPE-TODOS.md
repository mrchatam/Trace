# S03 — Store lock / concurrency — scope todos

**Depends-on:** P11-S02-02 done. Owns DF-47. (S02 FINAL: DF-43 FAIL-blocks-DONE + DF-44 conscious-flag docs — no lock coupling.)

**S03-00 locks (2026-08-16):** Keep exclusive `.trace/trace.lock`; **short bounded retry** (~250–500ms) on Open; clearer **ErrLocked**/help serialize CLI↔MCP or worktrees; exit **2**; **no** multi-writer / no mig.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | done — FINAL locks |
| 2 | 01-store-lock-concurrency | implement | pending |
| 3 | 02-scope-review | review | pending |

## Reminders
- Exclusive flock retained — DF-47 is UX + brief-race recovery, not unlimited same-root concurrency
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P11-S04-00**
