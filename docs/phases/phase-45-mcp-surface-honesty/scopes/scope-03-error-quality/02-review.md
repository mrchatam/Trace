# Phase 45 / S03 / Review — Error quality

## Row `P45-S03-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** Implemented in one session by the coding agent. Evidence, not a second opinion.

## Honesty invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| H1 | No SQL/driver text crosses the MCP boundary (#155) | **yes** | `assertNoDriverText` (`errors_test.go:188`) sweeps every wrapped error; `TestToolErrorsNeverLeakDriverText` |
| H2 | The error names the entity type actually queried, not the function's fixed type (#155) | **yes** | `wrapToolErr` reads `nf.Entity` from the S02 typed error |
| H3 | `trace_link` names the caller's wrong `from`/`to`, not the lookup order | **yes** | `wrapLinkErr` (`errors.go:58–62`) compares `nf.ID` against the supplied ids; `TestLinkErrorsNameRelVocabulary` |
| H4 | The rel's expectations are shown, so the caller learns the rule | **yes** | Message includes `expects from=… to=…` from `linkRelShape` |
| H5 | Hyphen **and** underscore rel aliases resolve (`api-contract` / `api_contract`) | **yes** | Both spellings present in `linkRelShape`; `strings.ToLower`+`TrimSpace` on lookup |
| H6 | `loop.next` and `loop.gate` speak one schema version (#156) | **yes** | Both render `mcpGateSchemaVersion = "trace.loop.gate.v1"`; no `trace.loop.next.v1` exists |
| H7 | A blocked `next` is machine-readable, not prose | **yes** | `reason_code` ∈ {`goal_link_missing`, `seed_task_not_found`}; `TestLoopNextReturnsGateEnvelopeOnPrecondition`, `TestLoopNextUnknownTaskReturnsGateEnvelope` |
| H8 | A blocked `next` says what to do | **yes** | `Remedy` field added to `mcpGateEnvelope` (`:39`) |
| H9 | Preconditions are checked before packet assembly, not after | **yes** | `checkLoopNextPreconditions` is called at `tools_loop.go:140`, before assembly |
| H10 | Non-not-found diagnostics keep their detail | **yes** | `wrapToolErr` passes them through with `%w` under the tool prefix only |

## Commands run

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1 -run 'TestLinkErrors|TestToolErrors|TestLoopNext'
CGO_ENABLED=1 go test ./internal/mcp/ ./internal/loop/ ./internal/store/ ./internal/domain/ -count=1
```

All green. Existing `TestMCPLoopNext` (happy path) still passes; the duplicate was deleted.

## Residuals (accepted, stated)

- **`linkRelShape` is a hand-maintained table.** A new rel added to `internal/domain` will not automatically appear here; an unknown rel degrades to the generic `wrapToolErr` message, which is still correct and still leak-free — just less helpful. Closing the loop (deriving the table from the domain rel vocabulary) is a real follow-up, listed on the board.
- **`remedy` is prose.** It is a string, not an enum, so an agent cannot switch on it. `reason_code` is the machine-readable field; `remedy` is for the model reading the message.
- Error rewriting is scoped to the MCP boundary. In-process Go callers still see the underlying store error — deliberate, so the fix does not churn internal call sites.

## Verdict

**PASS.** #155 and #156 closed by behaviour + permanent tests.