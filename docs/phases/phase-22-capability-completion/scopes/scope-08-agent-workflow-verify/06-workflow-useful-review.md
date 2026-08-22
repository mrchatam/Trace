# P22-S08-06 — Review: workflow usefulness

## Metadata
- id: P22-S08-06
- todo_ids: [P22-S08-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C39** (and C38 both halves still hold).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [05-workflow-useful.md](05-workflow-useful.md)
- Completion bar [README.md](../../README.md) §8–9

## Review checklist

### C39 — useful in normal agent workflow

- [ ] `TestHelpIncludesSearchTestVerify` PASS — search, test run, verify, changes, knowledge, loop, git-hook discoverable
- [ ] CONTRIBUTING **Agent workflow** paragraph present and accurate (index → loop → test → apply → search/why)
- [ ] Fresh `./bin/trace-mcp -h` lists **14** tools including `trace_loop`, `trace_regressions`
- [ ] `TestToolNamesRegistered` = `RegisteredToolNames()` = `-h` order
- [ ] Rebuild evidence in S08-05 Notes (mtime or command output)

### C38 hold (both halves)

- [ ] C38 checklist still `[x]` from S08-02
- [ ] `TestMCPLoop*` still PASS
- [ ] CLI test/verify (S03) + MCP loop (S08-01) both evidenced

### Landmines

- [ ] No hosted MCP / daemon docs added
- [ ] Do **not** require `trace_agents` yet — S09 scope
- [ ] Schema **26** at this row (027 comes S09)

## Spawn policy

If C39 unmet (missing help, stale `-h`, no CONTRIBUTING workflow): spawn **`P22-S08-06a`** + **`P22-S08-06b`**. Do not close with residuals.

## Re-run commands

```bash
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestHelpIncludesSearchTestVerify|TestToolNamesRegistered|TestHelpIncludesLoopNext|TestMCPLoop'
./bin/trace-mcp -h  # manual evidence: 14 tool names
grep -n "Agent workflow" CONTRIBUTING.md
```

## Exit criteria

- [ ] C39 closed (checklist `[x]`) or spawned
- [ ] C38 still closed
- [ ] Confidence **high**
- [ ] Board Notes: checklist boxed + `-h` evidence
