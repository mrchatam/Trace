# P22-S03-03 — Implement: run/record relevant tests

## Metadata
- id: P22-S03-03
- todo_ids: [P22-S03-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Agents can **explicitly run relevant tests** (Trace invokes project test command and records results). Closes **C12** and **C38-CLI** (test half). Supersedes P21 D-16 for explicit invoke only (D-22-03). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- S01: `internal/retrieval/impact_walk.go` (`AffectedTests`), `internal/store/file_graph.go` (`ListValidatesForFile`, `ListValidatesForSymbol`)
- Domain: `RecordTestOutcome`, `CheckTestGate` — Law 2: no pass claim without row
- Changes: `ListChangePaths` on latest task change (S02)
- **Depends on S03-01** cycle flags (test_pending clears when outcome recorded)

## Locked defaults

| Item | Value |
|------|-------|
| Package | `internal/testrun` (**new**) |
| CLI | `trace test run --task <id> [--paths path,...]` |
| Config | Optional `trace/test-runner.json`: `{"command":"...","args":[...],"cwd":"..."}` — if absent and `go.mod` at project root → `go test ./...` |
| Timeout | Context-bound (default **5m**); kill process on cancel |
| Record | Always `domain.RecordTestOutcome` per executed test target; map exit code → `test_status` pass/fail/error |
| Summary | Truncate stdout/stderr to `maxOutcomeSummaryBytes` (4096) |
| Select tests | 1) If `--paths`: impact walk seeds from paths → collect `AffectedTests` symbol/file paths; map to `validates` test files. 2) Else: paths from latest RECORDED/COMPARED change for task. 3) Fallback: unique Go packages of changed `.go` files → synthetic test name `package:<importpath>` |
| Fail-closed | Unknown stack (no go.mod, no config) → error, **no** fake pass row |
| Daemon | **Forbidden** — explicit CLI/hook/agent trigger only |
| Schema / compat | **23** — no migration |
| MCP | **Do not** add MCP tool this row (S08) |

## Requirements

1. `testrun.RunRelevantTests(ctx, st, dom, taskID, opts) ([]store.OutcomeResult, error)`.
2. Stub runner interface for unit tests (`Runner` with `Run(ctx, spec) (exitCode int, output string, err error)`).
3. Thin CLI + help + `root.go` dispatch; capability `cli:test` AUTO_ALLOW in `capability.go` pattern.
4. After successful record, `BuildPolicyInputs` should clear `test_pending` when re-read (integration assertion in test).

## Touch files

- `internal/testrun/runner.go`, `select.go`, `run.go`, `*_test.go` (**new**)
- `cmd/trace/test.go`, `test_test.go` (**new**)
- `cmd/trace/root.go`, `help.go`, `capability.go`
- `internal/domain/outcomes.go` (only if small helper needed)

## Named tests

| Test | Proves |
|------|--------|
| `TestTestRunRecordsOutcome` | stored `kind=test` row after invoke |
| `TestTestRunSelectsValidatingTests` | validates graph picks test file, not all packages |
| `TestTestRunFailClosedWithoutCommand` | no go.mod/config → error, zero outcomes |
| `TestTestRunUsesStubRunner` | unit test never shells real `go test` unless integration subtest |
| `TestPromotionGateRequiresStoredTestNotAgentClaim` | keeper (domain) |

Optional integration (tagged or subtest): tiny module under `internal/testrun/testdata/minimodule/` with passing `TestX`.

```bash
go test ./internal/testrun/... ./internal/domain/... -count=1 -run 'TestTestRun|TestRecordTestOutcome|TestPromotionGateRequiresStoredTestNotAgentClaim'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestTestRun'
```

## Exit criteria

- [ ] C12 true (explicit invoke + record)
- [ ] C38 CLI test half true (`trace test run`)
- [ ] Named tests PASS; compat **23**
- [ ] Board Notes

## Minimal todos

- [ ] testrun adapter + test selection
- [ ] CLI + capability gate
- [ ] Named tests
- [ ] Board notes
