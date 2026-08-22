# P01 / S02 / 02 — Scope review (Honesty demo)

## Metadata
- id: P01-S02-02
- todo_ids: [P01-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S02 (Honesty demo). Findings by severity; small fixes or spawn `a`/`b` pairs. Forward-only. May thicken **upcoming** S03 if `evals/honesty` vs `evals/x0` package boundaries need clarity.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling [01-honesty.md](01-honesty.md) + board Notes
- S01 [REVIEW-NOTES.md](../scope-01-claim-review/REVIEW-NOTES.md) (promotion contract)
- Live: `evals/honesty`, `internal/domain/{review.go,task_state.go,claim.go}`

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus (checklist)

### Fail-closed honesty (H5 partial)
- [ ] Planted claim narrative present (Claim + weak Evidence), not a tautology that always PASS
- [ ] Path A: EvidenceIDs alone → DONE rejected; WorkState unchanged
- [ ] Path B: `SetReviewResult(FAIL)` + `LinkReviewTask` → DONE rejected
- [ ] Path C: second Review PASS → DONE succeeds (recovery)
- [ ] `AllowDoneWithoutReview` **not** used to greenwash Paths A–B
- [ ] Artifact is under **`evals/honesty`** (not only a rename of `TestDoneRequiresReviewPass`)

### S01 contract fidelity
- [ ] Uses live APIs: `CreateReview` / `SetReviewResult` / `LinkReviewTask` / `TransitionTask`
- [ ] Does not invent VerifiedFact / alternate DONE gate
- [ ] Does not weaken S01 policy in product code

### Laws / regression
- [ ] No MCP / daemon / HTTP creep
- [ ] `CGO_ENABLED=0 go test ./evals/honesty/...` PASS
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` PASS
- [ ] Notes point to test name (findable evidence)

### Cross-scope
- [ ] S03 still owns `evals/x0` only — honesty package not swallowed into X0 metrics unless Phase planner later says so
- [ ] Light-update upcoming S03 prompts only if package layout surprises the harness planner

## Exit criteria
- [ ] Findings recorded (REVIEW-NOTES.md preferred)
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated

## Minimal todos
- [ ] Compare 01 claims + Notes vs repo evidence
- [ ] Run honesty + p0x + `./...` verification commands
- [ ] Fix or spawn
- [ ] Re-verify → board update
