# S07 — Seed / plan / review show polish — scope todos

**Depends-on:** P11-S06-02 done. Owns DF-28, DF-30, DF-33, DF-45, DF-46. (**S06 FINAL forward note:** DF-50 = print+write identical stderr reload tip; DF-22 = help/README tip parity + keep `trace_version`; DF-37 = tip/docs only — no Cursor process kill. No seed/plan/review coupling to S06 product work.)

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16 |
| 2 | 01-seed-plan-review-polish | implement | pending — next |
| 3 | 02-scope-review | review | pending |

## FINAL locks (summary)
- **DF-33:** seed link `from`/`to` **or** `from_id`/`to_id`; clear alias-aware empty error
- **DF-30:** `plan show` → `phases` always `[]` (never null) + goal `tasks` rows
- **DF-46:** plan show JSON snake_case (DF-32 parity)
- **DF-45:** CLI `review get`/`show`/`list [--task]`; snake_case; no MCP
- **DF-28:** thin help/docs handoff SoT (task body + Trace-pull); no entity/mig
- Migration **none**

## Reminders
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P11-S08-00**
