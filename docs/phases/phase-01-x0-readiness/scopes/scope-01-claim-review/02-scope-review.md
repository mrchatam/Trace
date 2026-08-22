# P01 / S01 / 02 — Scope review (Claim/Evidence/Review promotion)

## Metadata
- id: P01-S01-02
- todo_ids: [P01-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S01 (Claim/Evidence/Review promotion). Findings by severity; small fixes or spawn `a`/`b` pairs with full prompts. Forward-only. May thicken **upcoming** S02/S03 prompts if DONE/Review surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling [01-claim-review.md](01-claim-review.md) + board Notes
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G2/G3/G14

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus (checklist)

### Policy / laws
- [ ] DONE requires linked Review `result=PASS` (or documented escape hatch) — **EvidenceIDs alone must not authorize DONE**
- [ ] `AllowDoneWithoutReview` / `--allow-done` is opt-in only; **not** default or silent
- [ ] Implementer cannot self-certify PASS via TransitionTask narrative (only `SetReviewResult`)
- [ ] No VerifiedFact / MCP / daemon creep

### Implementation vs locks
- [ ] Mig `005_review_promotion.sql` additive (`reviews.result`); 001 untouched
- [ ] Store `UpsertReview`/`GetReview`; domain CreateReview / SetReviewResult / LinkReviewTask (`review_judges_task`)
- [ ] Events: `entity.created` + `review.result`; task.transition still appends
- [ ] Claim/Evidence path retained (`claim_has_evidence`)
- [ ] CLI thin only (G19); `trace review` present if claimed in Notes

### Verification
- [ ] Tests: reject EvidenceIDs-only; reject FAIL/no-review; accept PASS; escape hatch
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...`
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` still PASS
- [ ] Exit criteria in 01 honestly met (compare Notes ↔ repo)

### Cross-scope
- [ ] S02 honesty can call SetReviewResult FAIL + DONE reject (+ optional PASS path)
- [ ] S03 / CLI: if seeding DONE, must use Review PASS or explicit `allow_done` — thicken upcoming prompts if surface drifted

## Exit criteria
- [ ] Findings recorded (REVIEW-NOTES.md preferred)
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated; upcoming stubs thickened if needed

## Minimal todos
- [ ] Compare claims vs evidence against checklist
- [ ] Fix or spawn
- [ ] Re-verify (`go test` + p0x)
- [ ] Board update + REVIEW-NOTES.md
