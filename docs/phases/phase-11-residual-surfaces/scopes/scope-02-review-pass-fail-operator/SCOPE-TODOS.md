# S02 — Review PASS+FAIL / operator identity — scope todos

**Depends-on:** P11-S01-02 done. Owns DF-43, DF-44.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16 |
| 2 | 01-review-pass-fail-operator | implement | **done** — DF-43/44 shipped 2026-08-16 |
| 3 | 02-scope-review | review | pending — next |

## Locked reminders (FINAL)
- **DF-43:** any linked `review_judges_task` **FAIL** blocks →DONE even if sibling PASS; UNCERTAIN/empty do not block; hatch bypasses FAIL; honesty Path C must supersede FAIL before DONE
- **DF-44:** keep conscious `--as-operator` / `as_operator` (no OAuth); close via help/MCP flag≠identity wording; Actor≠auth
- Migration: **none**; G19 domain-owned gates
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P11-S03-00**
