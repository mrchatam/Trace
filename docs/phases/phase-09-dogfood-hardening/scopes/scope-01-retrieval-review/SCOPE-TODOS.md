# S01 — Retrieval review (DF-01) — scope todos

**Depends-on:** P09-00 done. Live gap: `lookupEntity` omits `"review"`; Expand via `review_judges_task` hard-fails why/context.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | finalize locks; thicken 01+02 |
| 2 | 01-retrieval-review | implement | ExactLookup `review` + regression; carry-forward green |
| 3 | 02-scope-review | review | independent DF-01 check; APPROVE or spawn |

## Locked reminders
- Full Hit (title + result/body excerpt); ReasonCode from caller; map `review_judges_*` rels
- Prefer fix over fail-soft skip
- No new mig; reuse `GetReview`
- Carry-forward: honesty Gate G + A/B/C, p0x, x0, `./...`
