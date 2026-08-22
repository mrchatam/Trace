# P22-S09-07 — Implement: CLI + MCP `trace agents`

## Metadata
- id: P22-S09-07
- todo_ids: [P22-S09-07]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Expose harness agent catalog and recommendations via CLI and stdio MCP so agents can query understood profiles without Trace acting as harness. **E02, E04** surface; strengthens **C39**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-agent-schema-routing.md](01-agent-schema-routing.md) — routing library
- Live: `internal/mcp/tools_loop.go` (thin mirror pattern), `cmd/trace/loop.go`, `internal/mcp/server.go`
- S08 pattern: G19 — handlers call library/domain only

## Live baseline (after S09-05)

| Present | Absent |
|---------|--------|
| Routing in loop next | CLI `trace agents` |
| Bundled catalog in DB after install | MCP `trace_agents` |
| MCP catalog **14** | Catalog **15** |

## Locked defaults

| Item | Value |
|------|-------|
| CLI | `trace agents list|describe|recommend` |
| MCP tool | **`trace_agents`** actions: `list`, `describe`, `recommend` (mirror CLI) |
| G19 | Handlers call `internal/agents` + store/domain — no duplicate routing logic |
| Annotations | All actions: `ReadOnlyHint=true`, `DestructiveHint=false` |
| Assert | `assertMCPToolAllowed(ctx, st, "trace_agents")` before every action |
| Catalog | **15** tools — append `trace_agents` after `trace_loop` in registration order |
| Specs | Add `mcp:trace_agents` to `BuiltinMCPCapabilitySpecs` |
| CLI specs | Add `cli:agents` to `BuiltinCLICapabilitySpecs` |
| Schema | Stays **27** — no new SQL this row |
| Hosted MCP | **Out** |

## CLI commands

### `trace agents list`

- JSON array: slug, title, subagent_type, deliberation_phases, requirement slugs
- Empty catalog → `[]`

### `trace agents describe <slug>`

- Full profile + requirements + registry metadata
- Unknown slug → exit non-zero, stderr message

### `trace agents recommend`

- Flags: `--task <id>` **or** `--phase CRITIQUE` (one required)
- Optional: `--goal-keywords "perf latency"` for keyword injection
- Output: same shape as `harness_recommendations[]` items (max 4)
- Uses same `RecommendAgents` as loop next

## MCP `trace_agents`

```go
type AgentsInput struct {
    Project  string `json:"project,omitempty"`
    Action   string `json:"action"` // list|describe|recommend
    Slug     string `json:"slug,omitempty"`
    TaskID   string `json:"task_id,omitempty"`
    Phase    string `json:"phase,omitempty"`
    Keywords string `json:"keywords,omitempty"`
}
```

- `recommend`: resolve task from DB when `task_id` set (title + deliberation phase); else use `phase` + `keywords`
- Return JSON text result (consistent with other trace MCP tools)

## Requirements

1. **`cmd/trace/agents.go`** (new) — subcommand wiring + tests
2. **`internal/mcp/tools_agents.go`** (new) — register in `server.go`
3. Update `RegisteredToolNames`, `cmd/trace-mcp/main.go` `-h`, `TestToolNamesRegistered` → **15** tools
4. Update help strings in `cmd/trace/help.go`
5. **No** spawn, no hosted server, no new migrations

## Touch files

- `cmd/trace/agents.go`, `cmd/trace/agents_test.go` (new)
- `cmd/trace/main.go` — register subcommand
- `cmd/trace/help.go`
- `internal/mcp/tools_agents.go` (new)
- `internal/mcp/server.go`, `internal/mcp/mcp_test.go`
- `internal/domain/capability.go` — specs
- `cmd/trace-mcp/main.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestCLIAgentsList` | list after install agents |
| `TestCLIAgentsDescribe` | describe known slug |
| `TestCLIAgentsRecommend` | recommend --phase CRITIQUE |
| `TestMCPAgentsRecommend` | MCP parity with CLI |
| `TestMCPAgentsList` | MCP list action |
| `TestToolNamesRegistered` | **15** tools; `trace_agents` last |

```bash
go test ./cmd/trace/... ./internal/mcp/... ./internal/domain/... -count=1 -run 'TestCLIAgents|TestMCPAgents|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] **E02** + **E04** query surface complete (CLI + MCP)
- [ ] Catalog **15**; specs + `-h` aligned
- [ ] Named tests PASS; compat **27**
- [ ] Board Notes: rebuild hint (`CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp`)

## Minimal todos

- [ ] CLI commands + help
- [ ] MCP tool + register + specs
- [ ] Named tests
- [ ] Board notes

## Residual risks (carry to S09-08)

- Recommend parity: MCP task_id path must load same inputs as loop next
- DENIED capability gating on `trace_agents` — mirror other read tools
- S08-07 VERIFY expects live `./bin/trace-mcp -h` with 15 tools — rebuild required
