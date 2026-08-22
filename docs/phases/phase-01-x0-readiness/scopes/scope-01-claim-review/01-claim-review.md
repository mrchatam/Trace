# P01 / S01 / 01 — Claim/Evidence/Review promotion

## Metadata
- id: P01-S01-01
- todo_ids: [P01-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Replace light Claim/Evidence stubs with a real **Claim→Evidence→Review→DONE** promotion path so Task **DONE** cannot be reached by implementer self-claim alone (G2/G3/G14 + H authority matrix). Domain + store only for policy; thin CLI adapter. Keep G19. Do **not** break `evals/p0x`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md) — authority matrix
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G2/G3/G14
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-CLAIM
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- Live priors: `internal/domain/{claim_stub.go,task_state.go,service.go}`; `reviews` in `store/schema/001_init.sql` (table only — no Upsert/Get Review); CLI `cmd/trace/transition.go` `--allow-done` / `--evidence`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Package | `internal/domain` on `*store.Store` only — **no** second DB |
| Store Review | Add `store.Review` + `UpsertReview` / `GetReview` mirroring Claim/Evidence |
| Migration | Additive embed `005_review_promotion.sql`: `ALTER TABLE reviews ADD COLUMN result TEXT NOT NULL DEFAULT ''` — values `PASS` \| `FAIL` \| `UNCERTAIN` \| `''` (open). Do **not** rewrite `001_init.sql` |
| Links | Keep `claim_has_evidence`. Add **`review_judges_task`** (from=`review`, to=`task`). Optional `review_cites_evidence` OK if documented in package comment |
| Claim/Evidence | Keep CreateClaim / CreateEvidence / LinkClaimEvidence; rename `claim_stub.go` → `claim.go` (or equivalent) — stubs become real path |
| Review API | `CreateReview`, `SetReviewResult(PASS\|FAIL\|UNCERTAIN)`, Get/List as needed for tests |
| DONE policy | `TransitionTask` → `DONE` **only if** task has linked Review with `result=PASS` **or** `AllowDoneWithoutReview==true`. **Non-empty EvidenceIDs alone is insufficient** |
| Escape hatch | Keep `AllowDoneWithoutReview` / CLI `--allow-done` / seed `allow_done` — **opt-in only**, never default/silent; for harness/seed/tests |
| Actor authority | No session/RBAC this scope. Control = API separation: only `SetReviewResult` writes PASS/FAIL; TransitionTask does not accept narrative “I claim PASS” |
| Events | `entity.created` on review create; `review.result` on SetReviewResult (payload: `result`, `actor`, `reason`); `task.transition` may include `review_id` when promoting |
| VerifiedFact | **Out** — stop at Review PASS → DONE eligibility |
| CLI | Thin G19: add `trace review` (create + set `--result`); keep `transition --allow-done`; no business logic in `cmd/trace` |
| CGO | Domain + new store APIs must pass `CGO_ENABLED=0` |
| Surface | Library + thin CLI only — **no** MCP / daemon / HTTP |
| Regression | `CGO_ENABLED=1 go test ./evals/p0x/...` and `./...` must stay green |
| Out | Multi-model adversarial review; VerifiedFact research suite; embeddings; full Gate G honesty (S02) |

### DONE promotion (locked)

```text
CreateClaim → CreateEvidence → LinkClaimEvidence (claim_has_evidence)
CreateReview → Link review_judges_task → Task
SetReviewResult(PASS)   # FAIL / UNCERTAIN / missing → DONE rejected
TransitionTask(..., DONE)  # succeeds when PASS review linked OR AllowDoneWithoutReview
```

Legal work_state graph unchanged from Phase 00 (`task_state.go`). Only the **DONE gate** changes.

### Minimum public API (names may vary slightly; behavior locked)

```text
# Existing (keep)
CreateClaim / CreateEvidence / LinkClaimEvidence / GetClaim / GetEvidence

# New
CreateReview(ctx, ReviewInput) (store.Review, error)
  // provenance defaults like other creates; result starts ''
SetReviewResult(ctx, reviewID, result, opts ReviewResultOptions) error
  // result in {PASS,FAIL,UNCERTAIN}; Actor+Reason required non-empty
LinkReviewTask(ctx, reviewID, taskID, meta LinkMeta) error
  // entity_links rel=review_judges_task
GetReview(ctx, id) (store.Review, error)

# TransitionTask DONE gate (replace stub)
→ DONE iff AllowDoneWithoutReview OR exists linked Review with result=PASS for this task
EvidenceIDs may still be recorded on the transition payload but do NOT alone authorize DONE
```

### Target tree

```text
internal/domain/
  claim.go              # was claim_stub.go — Claim/Evidence create+link
  review.go             # CreateReview / SetReviewResult / LinkReviewTask
  task_state.go         # DONE gate: PASS review | escape hatch
  service.go            # RelReviewJudgesTask + EventReviewResult consts
  domain_test.go        # policy tests (see Exit criteria)

internal/store/
  schema/005_review_promotion.sql
  entities_causal.go    # Review type + UpsertReview/GetReview (+ result column)

cmd/trace/
  review.go             # thin create + --result
  help.go / root.go     # wire subcommand
  transition.go         # keep --allow-done; usage note that Evidence alone ≠ DONE
```

### Tests (required)

- Reject `→ DONE` with only EvidenceIDs (no PASS review, flag false)
- Reject `→ DONE` with no review / Review FAIL / UNCERTAIN
- Accept `→ DONE` after Review PASS linked via `review_judges_task`
- Escape hatch: `AllowDoneWithoutReview=true` still works (explicit)
- Update/replace `TestDonePolicyStub` accordingly
- Store: Upsert/Get Review with `result` column after mig 005

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] Review create + SetReviewResult PASS/FAIL/UNCERTAIN path exists
- [ ] DONE requires linked Review PASS (or explicit AllowDoneWithoutReview); EvidenceIDs alone insufficient
- [ ] Tests cover reject/accept/escape-hatch cases above
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...` green
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` green
- [ ] Thin CLI `review` wired (G19); help updated
- [ ] TODO.md status + Notes updated (this row only)

## Minimal todos
- [ ] Store: Review type + Upsert/Get + mig `005_review_promotion.sql`
- [ ] Domain: review API + DONE gate replacement; rename claim stub
- [ ] Domain tests for DONE policy matrix
- [ ] Thin CLI `review` + help/transition notes
- [ ] Full test + p0x regression; board Notes
