# P01-S02-02 — Scope review notes (Honesty demo)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Reviewer:** independent (fresh context; claims checked against repo + re-run tests)

## Claims vs evidence

| Claim (P01-S02-01 Notes / 01 prompt) | Evidence |
|--------------------------------------|----------|
| Package `evals/honesty` (named H5 artifact) | `evals/honesty/doc.go` + `honesty_test.go`; not under `internal/domain` or `evals/p0x`/`evals/x0` |
| `TestHonestyFailClosedPlantedClaim` Paths A/B/C | `evals/honesty/honesty_test.go` — planted Claim+weak Evidence; A EvidenceIDs-alone reject; B FAIL review reject; C second PASS→DONE |
| WorkState unchanged on A/B rejects | `assertWorkState` → `IN_PROGRESS` after A and B |
| Path C DONE + `review_id` == PASS review | Asserts DONE + event payload `review_id == revPass.ID`; FAIL review remains FAIL |
| No `AllowDoneWithoutReview` greenwash | All three `TransitionTask` calls set `AllowDoneWithoutReview: false`; grep: no `: true` in package |
| Live S01 APIs only | `CreateClaim`/`CreateEvidence`/`LinkClaimEvidence`/`CreateReview`/`SetReviewResult`/`LinkReviewTask`/`GetReview`/`TransitionTask` — no VerifiedFact / alternate gate |
| Does not replace `TestDoneRequiresReviewPass` | Still present in `internal/domain/domain_test.go`; honesty is narrative sibling |
| No MCP / daemon / HTTP / product policy weaken | Honesty is eval-only; S01 `task_state.go` DONE gate unchanged |
| S03 owns `evals/x0` only | S03 `00-PLANNER` / `01-x0-harness` / `SCOPE-TODOS` already say keep honesty≠x0 — no thicken needed |
| Tests green | Independent: `CGO_ENABLED=0 go test ./evals/honesty/...` PASS; `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/...` + `./...` PASS |

## Checklist (02-scope-review)

### Fail-closed honesty (H5 partial) — PASS
- Planted claim narrative (not tautology always-PASS)
- Path A / B reject + state unchanged; Path C recovery via second PASS
- Escape hatch unused in proof
- Artifact under `evals/honesty`, distinct from domain unit test

### S01 contract fidelity — PASS
- Live CreateReview / SetReviewResult / LinkReviewTask / TransitionTask
- No VerifiedFact / alternate DONE gate; product policy not weakened

### Laws / regression — PASS
- No MCP/daemon/HTTP creep
- Honesty CGO=0 + p0x + full suite CGO=1 green
- Notes cite `TestHonestyFailClosedPlantedClaim`

### Cross-scope — PASS
- `evals/honesty` sibling to future `evals/x0`; S03 prompts already clear

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| nit | `SCOPE-TODOS.md` board checkboxes lagged implement/review | **Fixed inline** — boxes synced |

**Blocker/high:** none  
**Spawns:** none

## Residuals (none material)

- Claim entity is not linked to the task via entity_links — locked scenario uses EvidenceIDs + independent Review path; acceptable for H5 partial.
- Explicit `AllowDoneWithoutReview: false` is redundant with Go zero value but documents anti-greenwash intent.

## Next

**P01-S03-00** (X0 harness scope planner)
