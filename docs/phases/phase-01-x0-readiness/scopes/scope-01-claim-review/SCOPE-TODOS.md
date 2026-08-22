# Scope S01 — Claim/Evidence/Review promotion

**Depends-on:** Phase 00 complete. First Phase 01 product scope.

**Delivered (2026-08-16):** mig `005` `reviews.result`; store Upsert/Get Review; domain CreateReview/SetReviewResult/LinkReviewTask (`review_judges_task`); DONE = PASS review **or** explicit escape hatch (EvidenceIDs alone insufficient); thin `trace review`; `claim.go`; p0x green.

**Planner locks (2026-08-16):** DONE iff linked Review PASS **or** explicit escape hatch; EvidenceIDs alone insufficient; thin `trace review` CLI; VerifiedFact out; keep p0x green.

- [x] P01-S01-00 planner — 2026-08-16: locked DONE/Review policy; thickened 01+02; light S02/S03 Depends; no product Go
- [x] P01-S01-01 implement — 2026-08-16: Claim→Evidence→Review→DONE path shipped; tests + p0x PASS
- [x] P01-S01-02 review — 2026-08-16: APPROVE high; phase README live outcomes refreshed; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md)
