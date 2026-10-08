# Phase 45 / S06 / Implement — Guards

## Row `P45-S06-01` — **done** 2026-10-04.

## 1. Key-shape sweep — `internal/mcp/keyshape_test.go` (new)

- `walkJSON` / `assertLowerKeys`: recursively walk decoded JSON; **fail on any object key containing an uppercase letter**. Catches the #161 class (PascalCase Go fields) wherever it appears — including `map[string]any` payloads that struct reflection would miss.
- `TestNoToolResponseKeyContainsUppercase`: seeds `seedKeyShapeFixture` (plan tree, 2 tasks, change, regression, capability, indexed file `pkg/keyshape.go`, scope, agent catalog) and sweeps **21 real tool calls**. This closed the audit's previously-unverified `scope-member`, `api-contract` and real-regression response paths.
- `TestKeyShapeGuardCatchesPascalCase`: self-keeper — feeds the guard a deliberately PascalCase document and proves it fails.

## 2. Protocol conformance — `internal/mcp/protocol_conformance_test.go` (new)

Permanent replacement for the ephemeral `/tmp/mcp-probe/` Python harness from the audit. Uses `sdkmcp.NewInMemoryTransports()` + `srv.MCP().Connect(ctx, transport, nil)` / `client.Connect(ctx, t, nil)` via `connectProbe`:

| Test | Asserts |
|------|---------|
| `TestProtocolNegotiation` | initialize handshake answers with the server identity |
| `TestProtocolListsTheLockedToolSet` | `tools/list` returns **exactly** `RegisteredToolNames()` as a set (SDK sorts; order not asserted) |
| `TestProtocolStandardMethodsRespond` | `ping` (1 return value), `tools/list`, a real `tools/call` succeed |
| `TestProtocolUnknownMethodIsMethodNotFound` | unknown method → JSON-RPC method-not-found |
| `TestProtocolRejectsUnknownParams` | unknown parameter → `CallToolResult.IsError == true`, Go `err == nil` (SDK fact 3) |
| `TestProtocolRejectsOutOfEnumValue` | out-of-enum value → `IsError`, `err == nil` |
| `TestProtocolConcurrentReadersShareOneProject` | concurrent sessions over one project |

Helpers: `resultText`, `containsAny`. Session teardown is `serverSession.Close()` — `Wait()` hangs (verified the hard way during the audit harness).

## 3. Module hygiene

`go mod tidy` run: `github.com/google/jsonschema-go` promoted from `// indirect` to a direct dependency (it is now load-bearing in `schema.go`, not just a transitive SDK dep). `go.mod`/`go.sum` updated accordingly.

## 4. Acceptance

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1 -run 'TestProtocol|TestNoToolResponseKey|TestKeyShape'
```
All green; both files are permanent suite members, not skip-gated.

## Next

`P45-S06-02`