# Phase 45 / S02 / Review — Write-path validation

## Row `P45-S02-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** Implemented in one session by the coding agent. Evidence, not a second opinion.

## Honesty invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| H1 | `trace_add` cannot lose data without the caller asking (#150) | **yes** | `TestAddRejectsExistingIDWithoutReplace`; refusal message contains the remedy (`tools_write.go:104`) |
| H2 | The caller can always tell which happened | **yes** | `created` field on every `trace_add` response (`tools_write.go:205`); `TestAddReplaceTrueReportsCreatedFalse` |
| H3 | No dangling rows from `trace_add` (#151) | **yes** | `TestAddRejectsNonUUIDID` + `TestAddRejectsUnresolvableGoalID`; validated **before** insert |
| H4 | An invalid `work_state` is an error, not an empty list (#154) | **yes** | `TestTasksRejectsInvalidWorkState` |
| H5 | The plural filter names the offending element | **yes** | `work_states[%d]` (`tools_parity.go:95`) |
| H6 | Self-links are impossible through **any** adapter (#162) | **yes** | `TestLinkRejectsSelfReference` (MCP) **and** `TestLinkBlocksDomainRejectsSelfLink` (domain) |
| H7 | The self-link check runs before any lookup or write | **yes** | `toolLink` compares ids first; `refuseSelfLink` at the top of both domain link functions |
| H8 | `errors.Is(err, sql.ErrNoRows)` still true after the error refactor | **yes** | `NotFoundError.Is` (`errors.go:29`) matches that target explicitly; full store + domain + httpapi suites green |
| H9 | No FK migration, no cascade | **yes** | No migration file added; `replace=true` writes the entity row only |

## Commands run

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1 -run 'TestAdd|TestTasks|TestLink'
CGO_ENABLED=1 go test ./internal/store/ ./internal/domain/ ./internal/httpapi/ ./internal/mcp/ -count=1
```

All green.

## Residuals (accepted, stated)

- **`replace=true` does not cascade.** Overwriting a goal's title leaves its tasks pointing at a goal whose title changed. That is the *documented* behaviour for this phase; a real merge/delete operation is a different, destructive feature and was left out on purpose.
- **No FK on `tasks.goal_id`.** Validation covers the MCP write path; a direct `store`/`domain` caller that bypasses `trace_add` can still create a dangling row. Closed by migration in a later phase, not here.
- **Unknown `work_state` is rejected, not ignored.** An agent that was previously getting an empty list now gets an error. That is the intended fail-closed change, and it is a behaviour change for any caller passing a value outside the vocabulary.

## Verdict

**PASS.** #150, #151, #154, #162 closed by behaviour + permanent tests.