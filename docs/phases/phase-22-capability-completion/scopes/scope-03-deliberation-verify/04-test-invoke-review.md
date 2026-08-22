# P22-S03-04 — Review: test invoke

## Metadata
- id: P22-S03-04
- todo_ids: [P22-S03-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C12** (real invoke/record, not docs) and **C38-CLI** test half; no daemon; Law 2 honored.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. Grep `internal/testrun` + `cmd/trace/test.go` — must exist; `trace test run` in help.
2. Grep: no `net/http` ListenAndServe / daemon loop; no background watcher.
3. Every success path calls `RecordTestOutcome` — no stdout-only “pass”.
4. Fail-closed without runner config on unknown stack.
5. Test selection uses S01 `validates` / impact (grep `ListValidates` or `AffectedTests` in testrun).
6. C09 hold: after run, `test_pending` clears when re-querying policy (spot-check or test).
7. Schema **23**; MCP catalog still **10** (`TestToolNamesRegistered`).

## Spawn policy

If C12 or C38-CLI unmet: spawn **`P22-S03-04a` + `P22-S03-04b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/testrun/... ./internal/domain/... -count=1 -run 'TestTestRun|TestPromotionGateRequiresStoredTestNotAgentClaim'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestTestRun'
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
```

## Exit criteria

- [ ] C12 + C38-CLI (test) closed or spawned
- [ ] Confidence **high**
- [ ] Board Notes
