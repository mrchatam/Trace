# P11-S02-02 — REVIEW-NOTES (Review PASS+FAIL / operator identity)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-43** — linked FAIL+PASS + `AllowOperatorDone` rejects →DONE; reason mentions FAIL | **Pass** — `hasLinkedFailReview` before `findPassReviewID`; `TestSiblingFailBlocksDone` asserts reject + `FAIL` in reason |
| 2 | **DF-43** — PASS alone / PASS+UNCERTAIN still authorize with flag; hatch bypasses FAIL | **Pass** — `TestSiblingPassAloneAllowsDone`, `TestSiblingPassPlusUncertainAllowsDone`, `TestHatchBypassesSiblingFail` |
| 3 | **DF-43** — honesty Path C supersedes FAIL before DONE; A/B + Gate G green | **Pass** — Path C `SetReviewResult` FAIL→UNCERTAIN then PASS+`AllowOperatorDone`; Path B uses FAIL reject reason; Gate G hatch retained |
| 4 | **DF-44** — freestanding flag; help/MCP conscious flag≠identity; Actor≠auth; no OAuth | **Pass** — `--as-operator`/`as_operator` retained; help + transition usage + MCP desc/schema; `TestAsOperatorFlagIdentityDocs`, `TestAsOperatorSchemaIdentityDocs`, `TestOperatorDoneRequiresFlag`; no OAuth/OIDC |
| 5 | G19 — no domain fork in CLI/MCP adapters | **Pass** — FAIL scan only in `internal/domain/task_state.go`; adapters pass flags only |
| 6 | No forbidden architecture | **Pass** — no new mig (schema ends `010_capability_surface.sql`); no daemon/HTTP/embeddings/full-rebuild; hatch kept |
| 7 | Carry-forward + Gate C `dry_run:false` + P11-S01 DF-40 | **Pass** — locked CGO0/CGO1 suites + product `./cmd/... ./internal/... ./evals/...`; Gate C artifacts untouched (`dry_run:false`) |
| 8 | Board Notes accurate; planner row no product Go | **Pass** — P11-S02-00 Notes claim no product Go; P11-S02-01 Notes match live APIs/tests |

## Focus answers

- Check order matches locks: actor+reason → legal edge → DF-24 caps → →DONE hatch **or** (no linked FAIL ∧ PASS ∧ operator).
- Recovery is explicit `SetReviewResult` (no auto-clear FAIL when new PASS linked) — proven by FAIL+PASS reject then supersede success in `TestSiblingFailBlocksDone`.
- Scope-only links ignored (`RelReviewJudgesTask` only).

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Empty/unset review result not covered by a dedicated unit test (only UNCERTAIN + FAIL paths) | Residual OK — helper only matches `ReviewResultFail`; locks say empty does not block |
| low | `go test ./...` still fails setup on `similar projects/graphify` path space | Pre-existing non-product; product pkgs PASS |

## Residuals (explicit)

1. Empty-result linked review non-block is code-path covered, not a named empty-result test.
2. Full-module `./...` FAIL only on pre-existing `similar projects/graphify` space path.
3. DF-51 hatch-vs-caps remains S04 (not expanded here).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS
Named: TestSiblingFailBlocksDone / TestSiblingPassAloneAllowsDone / TestSiblingPassPlusUncertainAllowsDone / TestHatchBypassesSiblingFail / TestHonestyFailClosedPlantedClaim / TestAsOperatorFlagIdentityDocs / TestAsOperatorSchemaIdentityDocs / TestOperatorDoneRequiresFlag → PASS
```

## Next

**P11-S03-00**
