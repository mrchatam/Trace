# P22-S08-02 — Review: MCP loop

## Metadata
- id: P22-S08-02
- todo_ids: [P22-S08-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C38-MCP** (G19, no hosted server) and loop half of **C39**.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [01-mcp-loop.md](01-mcp-loop.md)
- S03 keepers: `TestLoopNextDeliberationSectionPresent`, `TestLoopApply`

## Review checklist

### C38-MCP — verification loop participation

- [ ] `trace_loop` registered; catalog **14**; `TestToolNamesRegistered` PASS
- [ ] `TestMCPLoopNext`, `TestMCPLoopApply`, `TestMCPLoopStatus` PASS
- [ ] Grep MCP handlers — call `loop.BuildNextPacket` / `Apply` / `Status`; **no** duplicated domain SQL
- [ ] **`apply`**: `DestructiveHint=true`; **`next`/`status`**: read-only hints
- [ ] `assertMCPToolAllowed` on all actions; slug `mcp:trace_loop` in specs
- [ ] Invalid action / missing task_id / bad envelope → fail-closed (not silent empty OK)

### C39 loop half (hold until S08-05 for full C39)

- [ ] MCP exposes same three loop subcommands as CLI
- [ ] Do **not** box full C39 this row — S08-05 owns help/docs/catalog gaps

### Landmines

- [ ] No hosted MCP / HTTP / daemon imports
- [ ] Schema **26** sql; compat **26** PASS
- [ ] S03 keepers still green: `go test ./internal/loop/... -count=1 -run 'TestLoopApply|TestLoopNext'`
- [ ] MCP catalog not regressed (still has search/changes/regressions)

## Spawn policy

If unmet: spawn **`P22-S08-02a`** (fix) + **`P22-S08-02b`** (re-review + box C38). Do not close with residuals.

## Re-run commands

```bash
go test ./internal/mcp/... -count=1 -run 'TestMCPLoop|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs'
go test ./internal/loop/... -count=1 -run 'TestLoopApply|TestLoopNextDeliberationSectionPresent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 26
```

## Exit criteria

- [ ] C38 closed (checklist C38 `[x]` with test citations) or spawned
- [ ] Confidence **high** | **medium** (spawn if medium+unmet)
- [ ] Board Notes: findings + confidence + C38 boxed when closed
