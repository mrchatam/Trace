# P22-S08-01 — Implement: MCP `trace_loop`

## Metadata
- id: P22-S08-01
- todo_ids: [P22-S08-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

MCP parity for the deliberation loop so agents can **initiate/participate in verification loops** over stdio (**C38-MCP**) and use Trace in normal workflow (**C39** loop half). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `cmd/trace/loop.go`, `internal/loop/{next,apply}.go`, `internal/mcp/tools_changes.go` (thin mirror pattern)

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| CLI `trace loop next\|apply\|status` gated `cli:loop` | MCP `trace_loop` |
| `loop.BuildNextPacket`, `loop.Apply`, `loop.Status` | `internal/mcp/tools_loop.go` |
| MCP catalog **13** | `TestMCPLoop*` |
| C38 CLI half closed (S03) | Checklist C38 still `[ ]` until this row + review |

## Locked defaults

| Item | Value |
|------|-------|
| Tool name | **`trace_loop`** |
| Actions | **`next`** (requires `task_id`), **`apply`** (requires `envelope` JSON string or object), **`status`** (requires `task_id`, optional `goal_id`) |
| G19 | Handlers call `internal/loop` only — no duplicate packet assembly or apply writes in MCP layer |
| Annotations | `next`/`status`: `ReadOnlyHint=true`, `DestructiveHint=false`; **`apply`**: `ReadOnlyHint=false`, **`DestructiveHint=true`** |
| Assert | `assertMCPToolAllowed(ctx, st, "trace_loop")` before every action |
| Catalog | **14** tools — append `trace_loop` after `trace_regressions` in registration order |
| Specs | Add `mcp:trace_loop` to `BuiltinMCPCapabilitySpecs` |
| Schema | **No SQL** — stays **26**; compat **26** |
| Checklist | Unbox C38 only when named tests PASS; C39 partial until S08-05 |

## Requirements

1. **`LoopInput`** struct: `project`, `action`, `task_id`, `goal_id`, `envelope` (apply only — raw JSON matching `trace.loop.apply.v1`).
2. **`next`**: mirror `cmdLoopNext` — open store, gate, build retrieval/compiler deps, `loop.BuildNextPacket`, return JSON text result.
3. **`apply`**: mirror `cmdLoopApply` — parse envelope via `loop.ParseApplyEnvelope` + `ValidateApplyEnvelope`, then `loop.Apply`; return apply result JSON.
4. **`status`**: mirror `cmdLoopStatus` — `loop.Status`, return status JSON.
5. Register in `registerTools`; update `RegisteredToolNames`, help in `cmd/trace-mcp/main.go` (include all **14** names).
6. **Do not** add hosted MCP, HTTP, or daemon paths.

## Touch files

- `internal/mcp/tools_loop.go` (new)
- `internal/mcp/server.go` — register tool + annotations
- `internal/mcp/mcp_test.go` — named tests
- `internal/domain/capability.go` — `BuiltinMCPCapabilitySpecs`
- `cmd/trace-mcp/main.go` — `-h` tool list

## Named tests

| Test | Proves |
|------|--------|
| `TestMCPLoopNext` | `action=next` returns valid `trace.loop.next.v1` packet for seeded task |
| `TestMCPLoopApply` | `action=apply` persists loop-step; destructive path works |
| `TestMCPLoopStatus` | `action=status` returns `trace.loop.status.v1` |
| `TestToolNamesRegistered` | **14** tools; `trace_loop` last |

```bash
go test ./internal/mcp/... -count=1 -run 'TestMCPLoop|TestToolNamesRegistered'
go test ./internal/domain/... -count=1 -run TestBuiltinMCPCapabilitySpecs
```

## Exit criteria

- [ ] C38-MCP true (named tests + grep G19)
- [ ] Catalog **14**; specs + `-h` aligned
- [ ] Compat **26** unchanged
- [ ] Board Notes: test summary + rebuild hint (`CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp`)

## Minimal todos

- [ ] `tools_loop.go` + register + specs
- [ ] Named tests (fixture: init + goal + task + minimal plan context)
- [ ] Board status + notes

## Residual risks (carry to S08-02)

- Apply envelope validation parity with CLI (schema_version, apply_id, seed mismatch)
- DENIED capability gating — mirror CLI fail-closed
- Missing compiler/retrieval wiring vs CLI next packet richness
