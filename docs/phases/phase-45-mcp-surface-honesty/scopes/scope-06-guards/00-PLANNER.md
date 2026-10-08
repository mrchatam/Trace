# Phase 45 / S06 / Planner — Guards

## Metadata
- id: P45-S06-00
- todo_ids: [P45-S06-00]
- role: planner
- skills: [testing-and-validation, api-and-interface-design]
- mcps: [user-trace]
- verification: automated

## Objective

Plan the permanent guards that keep the phase's fixes fixed: the key-shape sweep and the protocol conformance suite (the permanent replacement for the ephemeral `/tmp/mcp-probe/` Python harness used during the audit). **No product code in this row** (test-only code is the deliverable of S06-01).

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md) — SDK facts (esp. fact 3: schema failures are `IsError`, not Go errors)
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)

## Locked defaults

| Item | Value |
|------|-------|
| Key shape | Sweep **real tool calls**, not struct reflection — reflection misses `map[string]any` payloads |
| Sweep set | Every tool callable against a seeded fixture (21 calls), plus a self-keeper proving the guard can fail |
| Conformance | SDK in-memory transports (`sdkmcp.NewInMemoryTransports`), real `tools/list`, real `tools/call` |
| Tool set | **Set comparison** — `tools/list` results are sorted by the SDK, so order is not asserted |
| Unknown params / bad enums | Assert `CallToolResult.IsError == true` with Go `err == nil` (SDK fact 3) |
| Session lifecycle | `serverSession.Close()` — **never `Wait()`, it hangs** |

## Why these two guards

Every defect class in this phase had a "third recurrence" flavour (#161: *"third recurrence of this defect class"*). One-off tests per fix would let the next recurrence land in a tool nobody patched. A key sweep that walks **every real response** catches the class wherever it appears next, and a protocol suite that talks real JSON-RPC catches wire-level regressions (method negotiation, param validation, enum rejection) that handler-level tests cannot see.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S06-01**

## Next

`P45-S06-01`