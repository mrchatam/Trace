# P28-S06-09 — FM-08 / FR-P28-05 implementer

## Metadata
- id: P28-S06-09
- todo_ids: [P28-S06-09]
- role: implementer
- skills: [incremental-implementation, context-engineering]
- mcps: [user-codegraph, user-trace]
- verification: mixed
- hooks: []

## Objective

**FR-P28-05 / FM-08:** Make agents prefer task / promotion path over discovery-only edits after `trace_add`. Reinforce INT-06 MCP ordering + post-discovery nudge; optional apply-path smoke in dogfood.

## References

- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-08; INT-06
- `internal/mcp/server.go` description ordering; `internal/install/gappass.go` promotion nudge
- [TEST-MATRIX.md](../scope-01-integration-tests/TEST-MATRIX.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start.

## Acceptance hint

Session evidence: discovery → task (or spawned_tasks) before product edits; MCP description regression stays green.

## Preflight

```bash
cd /home/ali/Desktop/Trace
GOPROXY=direct go test ./internal/mcp/... -count=1
grep -n 'trace_add\|promotion\|discovery' internal/install/gappass.go internal/mcp/server.go | head
```

## Suggested work

1. Reinforce MCP ordering / gap-pass nudge text if weak.
2. Optional dogfood or smoke showing discovery→task before edits.
3. Keep MCP description unit regression green.
4. `FM08-NOTES.md` for session evidence.

## Out of scope

- Auto-spawn without gate (D1); other FMs; daemon/HTTP

## Exit criteria

- [ ] Acceptance hint met; MCP tests green
- [ ] Next runnable **P28-S06-10**

## Todo updates

Status + notes on **P28-S06-09** only.

## Next

`P28-S06-10`
