# Scope S01 — Coarse progressive planner

**Depends-on:** Phase 02 complete; `P03-00` done.

**Out:** Discovery replan/churn demo (S02); Gate E VERIFY (S03).

**Locks (P03-S01-00):** `internal/planner` + store mig `006_plan_hierarchy.sql` (phases/scopes/deep plans/current pointer). Deep-plan = current scope + one lookahead. Thin `trace plan`. S02 hooks: `SupersedeDeepPlan`, cursor, `auto_replan_count`.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P03-S01-00 | planner | done | 2026-08-16: locked `internal/planner` + mig 006 hierarchy; 01 thickened; S02 Depends hooks noted |
| P03-S01-01 | implement | done | 2026-08-16: planner + mig 006 + thin `trace plan`; bars green |
| P03-S01-02 | review | done | 2026-08-16: APPROVE high — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P03-S01-00 planner
- [x] P03-S01-01 implement
- [x] P03-S01-02 review
