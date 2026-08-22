# Scope S02 — Honesty escape-rate / Gate G prelim

**Depends-on:** S01 done (`P04-S01-02` APPROVE). Live S01 surface: mig `008` `review_residuals`; `LinkReviewScope` / `review_judges_scope`; `AddResidual` / `ListResidualsBy*` / **`CountOpenResidualsByScope`**; codes incl. `POLICY_EXCEPTION`.

**Locks (P04-S02-00):** Extend **`evals/honesty`**; named test **`TestHonestyEscapeRateGateGPrelim`**; schema **`evals/honesty/schema-gate-g.json`** v1; temp artifact **`metrics-gate-g.json`**; escape formula escapes=1 / caught=2 / attempts=3; keep Paths A/B/C; fail closed if S01 surface missing.

**Out:** Phase VERIFY (S03); full Gate G production policy; VerifiedFact; product Go outside harness.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P04-S02-00 | planner | done | 2026-08-16: Gate G prelim locks finalized |
| P04-S02-01 | implement | done | Gate G prelim harness + schema + metrics |
| P04-S02-02 | review | done | APPROVE high — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P04-S02-00 planner
- [x] P04-S02-01 implement
- [x] P04-S02-02 review
