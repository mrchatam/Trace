# Scope S01 — Scope review layer

**Depends-on:** Phase 03 complete; `P04-00` done.

**Out:** Honesty escape-rate / Gate G harness (S02); phase VERIFY (S03); VerifiedFact promotion.

**Locks (P04-S01-00):** Extend `internal/domain`+store mig **`008_scope_review.sql`**; `review_judges_scope` → `plan_scope`; structured `review_residuals` (severity INFO|WARN|BLOCKING; status OPEN|ACKED|RESOLVED); reuse CreateReview/SetReviewResult; no planner fork; task DONE unchanged; CLI `--scope` + `residual add|list`; VerifiedFact out.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P04-S01-00 | planner | done | 2026-08-16: locks above; 01/02 thickened; S02 Depends noted |
| P04-S01-01 | implement | pending | Runnable — see thickened 01-scope-review-layer.md |
| P04-S01-02 | review | pending | |

## Checklist

- [x] P04-S01-00 planner
- [ ] P04-S01-01 implement
- [ ] P04-S01-02 review
