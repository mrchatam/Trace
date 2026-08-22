# P22-S05-03 — Implement: evidence queries

## Metadata
- id: P22-S05-03
- todo_ids: [P22-S05-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Agents can ask **what tests verify X**, **what previously failed**, **what approaches worked**, and **about regressions**. Closes **C17, C31, C32, C33, C34**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- S01: `ListValidatesForSymbol`, `ListValidatesForFile` (`internal/store/file_graph.go`)
- S03: `outcome_results` test/evaluation kinds + status vocab
- S04: `ListRegressionsByChangeID`, `ListAllRegressions`, `ListImprovementsByTaskID`/`ByChangeID`

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `trace test run` | `trace tests verifying` |
| `trace outcomes compare\|improvements` | `failed\|worked` |
| Domain regression/improvement APIs | `trace regressions` CLI |
| MCP **12** (after S05-01) | `trace_regressions` |
| Compat **24** | **025+** |

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **None**; compat stays **24** |
| Domain helpers | New `internal/domain/queries.go` (or `evidence_query.go`): thin wrappers over store |
| Tests verifying (C31) | `ListTestsValidatingSymbol(ctx, symbolID)` → reverse `validates` edges to test symbols/files; `ListTestsValidatingFile(ctx, path)` resolves path → file id → `ListValidatesForFile`; include `edge_provenance`, test symbol/file ids + paths |
| Outcomes failed (C32) | `ListFailedOutcomes(ctx, opts)` — rows where `kind=test` AND `test_status IN (fail,error)` OR `kind=evaluation` with failing signal in `summary`/scores (lock: **test fail/error only** for v1); optional `--task`; newest first; limit 32/64 |
| Outcomes worked (C33) | `ListWorkedApproaches(ctx, opts)` — **union** of: (a) `improvements` for task/change scope, (b) `outcome_results` test pass rows for task; dedupe by id; newest first; limit 32/64; compact JSON |
| Regressions (C17/C34) | `ListRegressions(ctx, opts)` — filter optional `task_id`, `change_id` (via `ListRegressionsByChangeID`); else `ListAllRegressions` truncated newest-first; include `attribution`, `dimension`, `summary`, associated change ids when linked |
| CLI | `trace tests verifying --symbol <uuid> \| --file <path>` (xor); `trace outcomes failed [--task <id>] [--limit N]`; `trace outcomes worked [--task <id>] [--limit N]`; `trace regressions list [--task <id>] [--change <id>] [--limit N]` |
| CLI gating | `cli:tests` / `cli:outcomes` / new `cli:regressions` in `BuiltinCLICapabilitySpecs` |
| MCP | **`trace_regressions`** only — action=`list`, optional `task_id`, `change_id`, `limit`; catalog **13** |
| Tests/outcomes MCP | **CLI-only** this row (planner lock — avoid catalog explosion) |
| Fail-closed | `--symbol` and `--file` both missing/extra on verifying → usage error; unknown symbol/file → validation error |

## Requirements

1. Domain query helpers + unit tests in `internal/domain/*_test.go`.
2. Extend `cmd/trace/test.go` → subcommand router `trace tests verifying` (or add `tests.go` if cleaner).
3. Extend `cmd/trace/outcomes.go` — `failed`, `worked` subcommands.
4. New `cmd/trace/regressions.go` + root wiring.
5. `internal/mcp/tools_regressions.go` + register; update `RegisteredToolNames`.
6. Checklist C17, C31–C34 **unboxed** until S05-04 review.

## Touch files

- `internal/domain/queries.go`, `queries_test.go` (new)
- `cmd/trace/test.go` or `tests.go`, `outcomes.go`, `regressions.go` (+ tests)
- `cmd/trace/root.go`, `help.go`
- `internal/mcp/tools_regressions.go`, `server.go`, `mcp_test.go`
- `internal/domain/capability.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestTestsVerifyingQuery` | validates edge fixture → verifying returns test symbol |
| `TestOutcomesFailedAndWorked` | seeded fail pass + improvement → failed/worked lists correct |
| `TestRegressionsListQueryable` | regression + associated change filter |
| `TestMCPRegressionsRegistered` | trace_regressions list works |
| `TestToolNamesRegistered` | exactly **13** tools |
| `TestListRegressionsByChangeID` | keeper (S04) |

```bash
go test ./internal/domain/... -count=1 -run 'TestTestsVerifying|TestOutcomesFailed|TestRegressionsList|TestListRegressionsByChange'
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestTestsVerifying|TestOutcomes|TestRegressions|TestMCPRegressions|TestToolNamesRegistered'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C17, C31–C34 true
- [ ] MCP catalog **13**; compat **24**
- [ ] Named tests PASS
- [ ] Board Notes

## Minimal todos

- [ ] Domain query helpers + tests
- [ ] CLI tests/outcomes/regressions
- [ ] MCP trace_regressions
- [ ] Capability specs + help
- [ ] Board notes

## Residual risks (carry to S05-04)

- File-path verifying when path not indexed — must fail-closed with clear error
- Project-wide failed/worked without `--task` may be noisy — bounded + newest-first required
- C33 “worked” vs S06 “successful approaches” overlap — improvements + pass rows sufficient for S05; S06 may extend
