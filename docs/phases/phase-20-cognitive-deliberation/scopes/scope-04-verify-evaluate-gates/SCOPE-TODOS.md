# S04 — Verify / evaluate / gates — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P20-S04-00 | scope planner | done |
| 2 | P20-S04-01 | implementer | done |
| 3 | P20-S04-02 | reviewer | done — APPROVE high; next P20-S05-00 |

**Feeds:** S01 `verification_incomplete` input (`HasVerificationDebt`), S06 packet `verification_debt`, S05 regression from evaluation `comparison_json`.

**FINAL locks (S04-00):** `outcome_results` + `baselines`; migration **018**; compat ceiling **18**; kinds `test|verification|evaluation`; verification = goal_id + evidence; evaluation = scores vs baseline (not boolean); debt query for implementation-without-verification; DONE/Review unchanged; library-only.
