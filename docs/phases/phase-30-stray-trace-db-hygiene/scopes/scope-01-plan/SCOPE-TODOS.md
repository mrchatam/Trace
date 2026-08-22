# Scope 01 — board map

Plan-only scope. Serial: **S01-00 → S01-01**. No separate S01 review — S02-00 re-reads `PLAN.md`.

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 531 | P30-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (lock defaults; thicken 01-plan; no product code) |
| 532 | P30-S01-01 | [01-plan.md](01-plan.md) | Implementer → write `PLAN.md` |

Primary artifact: `PLAN.md` (written by **P30-S01-01**).

Planner lock (P30-S01-00, 2026-08-21): S00 verdict **agent hygiene** (INTAKE confirmed; no Trace dual-store bug). Track **T1–T4** only. Canonical store `.trace/trace.db` locked. Warn: stderr, once-per-open, non-fatal. Gitignore: `/trace.db`. No silent delete; no path redesign; no HTTP creep.

S01-01 done 2026-08-21 — artifact: [`PLAN.md`](PLAN.md). Next: **P30-S02-00**.
