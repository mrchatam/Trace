# P22-S09-08 — Review: CLI + MCP trace agents

## Metadata
- id: P22-S09-08
- todo_ids: [P22-S09-08]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Independent review of S09-07. Confirm CLI/MCP parity, no harness execution, tool catalog **15**, and **S09 scope complete (E01–E04)**.

## Session start

**Fresh subagent** (not S09-07). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [07-cli-mcp-agents.md](07-cli-mcp-agents.md)
- [VERIFY.md](../../VERIFY.md) — E01–E04 evidence required before S08-07
- [README.md](../../README.md) — enhancement matrix + completion bar item 9–10
- S08-07 preconditions: S09-00…S09-08 all `done`

## Review checklist

### E02 + E04 — CLI/MCP surface

- [x] `trace agents list|describe|recommend` work on fixture with `install agents`
- [x] `TestCLIAgentsList`, `TestCLIAgentsDescribe`, `TestCLIAgentsRecommend` PASS
- [x] `TestMCPAgentsRecommend`, `TestMCPAgentsList` PASS
- [x] MCP `recommend` matches CLI for same inputs (shared `loop.RecommendHarness`)

### W-40 — MCP catalog

- [x] `trace_agents` registered; `TestToolNamesRegistered` → **15** tools, `trace_agents` last
- [x] `mcp:trace_agents` in `BuiltinMCPCapabilitySpecs`
- [x] `cli:agents` in `BuiltinCLICapabilitySpecs`
- [x] `./bin/trace-mcp -h` lists all **15** via `RegisteredToolNames()` (`TestToolNamesRegistered`; **FYI** committed `bin/trace-mcp` may predate S09-07 — rebuild with `CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp`)
- [x] Read-only annotations; `assertMCPToolAllowed` on all actions

### E01 + E03 — end-to-end (cross-row)

- [x] Loop next includes `harness_recommendations` (S09-05 keeper green)
- [x] Bundled catalog installed (S09-03 keeper green)
- [x] Fresh subagent hint when hook AVAILABLE (S09-05 keeper)

### Hard boundaries (full S09)

- [x] Grep S09 code — no Task/spawn/subprocess agent runner
- [x] No HTTP registry fetch
- [x] No hosted MCP / daemon
- [x] Schema **27**; compat **27**
- [x] Trace role remains **recommend only** (D-22-25)

### S09 scope completion gate

| Enhancement | Evidence |
|-------------|----------|
| E01 | `TestRecommendSubagentWhenAvailable` + loop packet |
| E02 | Routing tests + CLI/MCP recommend |
| E03 | default.json + install + schema 027 |
| E04 | README + registry columns; no network |

## Spawn policy

Gaps in E02/E04 MCP surface or cross-row E01/E03 → spawn **`P22-S09-08a`** implement + **`P22-S09-08b`** review. If S09 incomplete, **block S08-07 VERIFY** (already in S08-07 preconditions).

## Re-run commands

```bash
go test ./cmd/trace/... ./internal/mcp/... ./internal/loop/... ./internal/install/... ./internal/agents/... -count=1 -run 'TestCLIAgents|TestMCPAgents|TestToolNamesRegistered|TestLoopNextIncludesHarness|TestInstallAgents|TestRecommendSubagent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp && ./bin/trace-mcp -h | tr ',' '\n' | tail -3
ls internal/store/schema/*.sql | wc -l  # expect 27
```

## Exit criteria

- [x] `trace_agents` registered; tests PASS
- [x] **E01–E04** closed with evidence citations in Notes
- [x] S09 scope complete — next runnable **P22-S08-07** VERIFY (board order)
- [x] Confidence **high**
