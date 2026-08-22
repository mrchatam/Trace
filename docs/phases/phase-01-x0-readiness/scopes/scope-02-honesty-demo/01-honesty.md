# P01 / S02 / 01 — Honesty demo

## Metadata
- id: P01-S02-01
- todo_ids: [P01-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Deliver a **named H5 partial** honesty artifact: a planted false/incomplete completion claim is rejected by independent Review and **cannot** reach Task `DONE`. Deterministic only — no LLM reviewer. Prove fail-closed against the **live S01** promotion surface. Do **not** break `evals/p0x`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md) — authority matrix; anti-pattern “skipping honesty”
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — H5
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G2/G3/G14
- Sibling [00-PLANNER.md](00-PLANNER.md) locks
- S01 APPROVED: [../scope-01-claim-review/REVIEW-NOTES.md](../scope-01-claim-review/REVIEW-NOTES.md)
- Live APIs: `internal/domain/{review.go,task_state.go,claim.go,service.go}`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Depends | S01 Claim/Review promotion **done** (APPROVE high, 2026-08-16) |
| Package | **`evals/honesty`** (new) — named H5 demo; **not** folded into `internal/domain` unit tests; **not** into `evals/p0x` or `evals/x0` |
| Verification | **`automated`** — no human gate; do not mark `mixed` |
| Escape hatch | **Forbidden in the honesty proof** — never set `AllowDoneWithoutReview` / `--allow-done` to “pass” the demo |
| S01 surface | `CreateReview` / `SetReviewResult` / `LinkReviewTask` (`review_judges_task`) / `GetReview`; DONE via `TransitionTask` only when linked Review `result=PASS` |
| Claim path | Use `CreateClaim` / `CreateEvidence` / `LinkClaimEvidence` (`claim_has_evidence`) for the planted narrative |
| Out | Full Gate G research suite; LLM-as-reviewer as sole proof; VerifiedFact; MCP/daemon/HTTP; rewriting S01 policy |
| Regression | Keep `evals/p0x` + `./...` green |
| CGO | Honesty package uses domain+store only → must pass `CGO_ENABLED=0 go test ./evals/honesty/...`; full regression `CGO_ENABLED=1 go test ./...` |
| G19 | Optional CLI transcript only; library must not import `cmd/trace`. Prefer pure domain demo (no CLI required) |
| Duplicate policy | Do **not** delete/replace `TestDoneRequiresReviewPass` — that stays S01 unit coverage. S02 adds the **narrative H5 artifact** under `evals/honesty` |

### Scenario (locked — implement exactly)

One Go test (suggested name: `TestHonestyFailClosedPlantedClaim`) that opens a temp `.trace` DB via store + `domain.New` (same pattern as `internal/domain` tests — copy helper or thin local open).

```text
# Setup
CreateTask(title ≈ "Ship feature X") → TransitionTask IN_PROGRESS (actor=implementer)

# Planted false / incomplete claim (implementer narrative)
CreateClaim(title ≈ "Feature X is complete")
CreateEvidence(title ≈ "tests pass on my machine")   # weak / non-independent
LinkClaimEvidence(claim, evidence)

# Path A — EvidenceIDs alone must NOT unlock DONE
TransitionTask(DONE, EvidenceIDs=[evidence.ID], AllowDoneWithoutReview=false)
  → MUST error (*ErrInvalidTransition; reason mentions Review PASS / AllowDoneWithoutReview)
  → task.WorkState still IN_PROGRESS

# Path B — Independent review FAIL → still no DONE
CreateReview(title ≈ "Honesty check: completion claim")
LinkReviewTask(review, task)           # review_judges_task
SetReviewResult(FAIL, actor=reviewer, reason ≈ "claim unproven / incomplete")
GetReview → result == FAIL
TransitionTask(DONE, AllowDoneWithoutReview=false)
  → MUST error (same policy)
  → WorkState still IN_PROGRESS

# Path C — Optional recovery (required in this scope)
CreateReview(title ≈ "Re-review after remediation")   # second review; leave FAIL linked
LinkReviewTask(review2, task)
SetReviewResult(PASS, actor=reviewer, reason ≈ "evidence now adequate")
TransitionTask(DONE) → MUST succeed; WorkState == DONE
  → task.transition payload may include review_id == review2.ID
```

**Asserts (minimum):**
1. Path A reject + state unchanged  
2. Path B: `GetReview` is FAIL; DONE reject; state unchanged  
3. Path C: DONE after second PASS review  
4. Never call `AllowDoneWithoutReview: true` in Paths A–B (Path C also false)

### Target tree

```text
evals/honesty/
  doc.go              # package purpose + how to run
  honesty_test.go     # TestHonestyFailClosedPlantedClaim (+ helpers)
```

Optional (not required): CLI transcript under Notes only — if added, keep G19 thin and do not make CLI the sole proof.

### How to run (implementer Notes must cite)

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `evals/honesty` exists with automated fail-closed test covering Paths A–C above
- [ ] Notes cite test name + package path (artifact findable)
- [ ] `AllowDoneWithoutReview` unused in fail paths
- [ ] `CGO_ENABLED=0 go test ./evals/honesty/...` green
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` green
- [ ] TODO.md status + Notes updated

## Minimal todos
- [ ] Add `evals/honesty` + fail-closed scenario test (A/B/C)
- [ ] Self-check exit criteria / run commands above
- [ ] Board status + notes (test name in Notes)
