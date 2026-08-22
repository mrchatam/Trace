# P39-S00-02 — Review G1 context orient merge

## Metadata
- id: P39-S00-02
- todo_ids: [P39-S00-02]
- role: reviewer
- skills: [code-review-and-quality, context-engineering]
- mcps: [user-trace]
- verification: mixed

## Objective

Fresh independent review of S00-01 G1 implementation vs REMEDIATION-PLAN G1, GAP-REGISTRY G-001/G-002, and M-001 moat charter.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — acceptance tests T1–T6
- [00-PLANNER.md](00-PLANNER.md) — locks + rejects
- [REMEDIATION-PLAN G1](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-001/G-002](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high without trivial fix → spawn 02a/02b pair below this row |
| MCP tool count | Must remain **16** — verify `RegisteredToolNames()` in `server.go:227–235` unchanged length |
| Cap defaults | `MaxCandidateHits=64`, `DefaultMaxItems=32`, `DefaultTokenBudget=4096` unchanged (`packet.go:18–20`) |
| Live pre-state (P39-S00-00) | No `Query` on `ContextOptions`; no `--query`; no MCP `query` field |

## Review checklist

### A — G1 gap closure

- [ ] Optional `query` on compiler (`ContextOptions`), CLI (`--query`), MCP (`trace_context.query`)
- [ ] Query FTS hits merged into packet (not separate compose step for agents)
- [ ] Title FTS path still active when query provided (G-002 partial lift)

### B — M-001 moat

- [ ] `task_id` required; Layer 0 task core present
- [ ] No query-only code path
- [ ] Gates/loop/review tools untouched

### C — Laws 6–7

- [ ] No default cap increase
- [ ] `Budget.Truncated` / `CandidatesCapped` honest when over budget
- [ ] No full-graph dump API added

### D — Law 19

- [ ] Merge logic in `internal/compiler/` — CLI/MCP thin
- [ ] No business-logic fork in MCP beyond schema pass-through

### E — Tests

- [ ] T1–T6 + T1-MCP from implement prompt evidenced (run or read test bodies)
- [ ] Regression green: `TestTaskContextContinuesWhenSearchErrors`, `TestCandidateCapSetsTruncated`, `TestItemCapNeverExceeded`, `TestBudgetLoudTotals`, `TestNoDumpAPI`
- [ ] `go test ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1` passes

### G — Live verification commands

```bash
go test ./internal/compiler/... ./internal/mcp/... -count=1 -run 'G1|CandidateCap|SearchErrors|NoDump'
grep -c 'AddTool(s.mcp' internal/mcp/server.go   # expect 16
grep 'Query' internal/compiler/compiler.go internal/mcp/tools_context.go cmd/trace/context.go
```

### F — Rejects

- [ ] No `trace_explore`
- [ ] No semantic/vector channel
- [ ] No new MCP tool

## Findings template

Report: blocker | high | medium | low | nit — with file:line evidence.

## Exit criteria

- [ ] APPROVE (medium+ confidence) **or** spawn repair pair with full prompts
- [ ] Residual risks listed explicitly if medium confidence
- [ ] Board row → `done` with verdict + confidence in Notes

## Next

`P39-S01-00`
