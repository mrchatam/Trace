# Phase 45 / S06 / Review — Guards

## Row `P45-S06-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** Implemented in one session by the coding agent. Evidence, not a second opinion.

## Guard-quality invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| G1 | The key guard can actually fail (not a tautology) | **yes** | `TestKeyShapeGuardCatchesPascalCase` feeds it a PascalCase doc and asserts failure |
| G2 | The key guard sees real responses, not only structs | **yes** | 21 real handler calls decoded and walked — covers `map[string]any` payloads |
| G3 | The previously-unverified response paths are covered | **yes** | `seedKeyShapeFixture` includes a scope, `api-contract` link, and a real regression; all swept |
| G4 | The wire protocol is tested over real JSON-RPC, not just handler calls | **yes** | In-memory transports; initialize/tools-list/tools-call/ping round-trips |
| G5 | The locked 17-tool set cannot silently grow or shrink | **yes** | `TestProtocolListsTheLockedToolSet` set-compares against `RegisteredToolNames()` |
| G6 | Schema enforcement is proven at the wire level (#158) | **yes** | Unknown-param and out-of-enum `tools/call` both yield `IsError` with `err == nil` |
| G7 | Concurrent sessions cannot corrupt project state | **yes** | `TestProtocolConcurrentReadersShareOneProject` |
| G8 | Guards are permanent suite members | **yes** | No skip/t.Mkip; run in every `go test ./internal/mcp/` |

## Commands run

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1 -run 'TestProtocol|TestNoToolResponseKey|TestKeyShape'
CGO_ENABLED=1 go test ./internal/mcp/ -count=1
CGO_ENABLED=1 go test ./... -count=1
```

All green.

## Residuals (accepted, stated)

- **The 21-call sweep is a fixed list.** A new tool that is never added to the sweep escapes the key guard (the enum/description drift guard has the same shape). Deriving the sweep from `RegisteredToolNames()` programmatically is the obvious follow-up; deferred to keep each call's fixture explicit and debuggable.
- **`Wait()` hangs** is an SDK behaviour worth an upstream note; the suite simply avoids it.

## Verdict

**PASS.** Guard suite permanent and green.