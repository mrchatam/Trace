# P01-S01-02 — Scope review notes (Claim/Evidence/Review promotion)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Reviewer:** independent (fresh context; claims checked against repo)

## Claims vs evidence

| Claim (P01-S01-01 Notes / 01 prompt) | Evidence |
|--------------------------------------|----------|
| Mig `005_review_promotion.sql` additive `reviews.result` | `internal/store/schema/005_review_promotion.sql`; `001_init.sql` reviews table has **no** `result` column (untouched) |
| Store `UpsertReview` / `GetReview` | `internal/store/entities_causal.go`; `TestReviewUpsertGetResult` |
| Domain `CreateReview` / `SetReviewResult` / `LinkReviewTask` (`review_judges_task`) | `internal/domain/review.go`; consts in `service.go` |
| DONE iff PASS review **or** `AllowDoneWithoutReview`; EvidenceIDs alone rejected | `task_state.go` `TransitionTask` + `findPassReviewID`; `TestDonePolicyStub` + `TestDoneRequiresReviewPass` |
| Escape hatch opt-in only | `AllowDoneWithoutReview` default false; CLI `--allow-done` default false |
| Only `SetReviewResult` writes PASS/FAIL | `TransitionOptions` has no result/PASS field; narrative cannot self-certify |
| Events `entity.created` + `review.result`; `task.transition` may include `review_id` | CreateReview → `appendCreated`; SetReviewResult → `EventReviewResult`; DONE payload asserts `review_id` in test |
| Claim/Evidence retained; `claim_stub.go` → `claim.go` | `claim.go` present; `claim_stub.go` absent; `RelClaimHasEvidence` |
| Thin CLI `trace review create\|set` | `cmd/trace/review.go` + `root.go`/`help.go`; CLI test covers evidence-alone reject + PASS promote |
| No VerifiedFact / MCP / daemon | Grep: no product MCP/daemon/VerifiedFact in domain/store/cmd for this scope |
| Tests green | Independent: `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...` PASS; `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` PASS |

## Checklist (02-scope-review)

### Policy / laws — PASS
- DONE requires linked Review PASS or explicit escape hatch; EvidenceIDs alone insufficient
- `--allow-done` / `AllowDoneWithoutReview` opt-in only
- Implementer cannot self-certify PASS via TransitionTask
- No VerifiedFact / MCP / daemon creep

### Implementation vs locks — PASS
- Mig 005 additive; 001 untouched
- Store + domain APIs as locked
- Events present; claim_has_evidence retained
- CLI thin (G19)

### Verification — PASS
- Reject EvidenceIDs-only; reject FAIL/open/UNCERTAIN/no-review; accept PASS; escape hatch
- Domain/store CGO=0 + p0x + full suite CGO=1 green

### Cross-scope — PASS (after doc fix)
- S02 honesty prompts already call SetReviewResult FAIL + DONE reject (+ optional PASS)
- S03 harness notes Review PASS or explicit `allow_done` if DONE needed
- **Medium fixed inline:** Phase 01 README “Prior phase outcomes” still described stub DONE / no Review API — updated to live S01 surface so S02+ agents are not misled

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| medium | Phase README live outcomes stale (stub DONE / no Review API; mig 001–004; CLI missing `review`) | **Fixed inline** — README Prior outcomes updated |
| nit | `TestDonePolicyStub` name still says Stub after real policy | Accept residual; rename optional later |
| nit | S01 `SCOPE-TODOS.md` checkboxes lagged implement | Updated with this review |

**Blocker/high:** none  
**Spawns:** none

## Residuals (none material)

- CLI defaults empty `--actor` to `"cli"` (same pattern as `transition`); reason still required.
- Store `UpsertReview` does not enum-validate `result` (domain `SetReviewResult` does) — acceptable G19 boundary.

## Next

**P01-S02-00** (honesty demo scope planner)
