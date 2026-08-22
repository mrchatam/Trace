# P39-S01-02 — Review G3 harness orient

## Metadata
- id: P39-S01-02
- todo_ids: [P39-S01-02]
- role: reviewer
- skills: [code-review-and-quality, context-engineering, writing-for-agents]
- mcps: [user-trace]
- verification: mixed

## Objective

Fresh independent review of S01-01 G3 harness changes vs REMEDIATION-PLAN G3, GAP-REGISTRY G-006/G-010, and M-001 moat charter.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G3-A1–A6
- [00-PLANNER.md](00-PLANNER.md)
- [REMEDIATION-PLAN G3](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [INTAKE.md](../../INTAKE.md) — Q5 9/16 hybrid resolution

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high without trivial fix → spawn 02a/02b pair below this row |
| MCP tool count | Must remain **16** — `RegisteredToolNames()` + 16× `AddTool` |
| 9/16 fix | Docs + Instructions + trace_version — **not** tool reduction |
| G1 cross-check | Instructions reference optional `trace_context.query` (shipped S00) |
| Rejects | CG 1-tool facade; MP 44-tool copy; hidden write tools; `trace_explore` |

## Review checklist

### A — G3 gap closure

- [ ] MCP Instructions present and agent-facing (go-sdk `ServerOptions.Instructions` wired in `NewServer`)
- [ ] Ranked moat-first playbook orients agents without 16-tool roulette
- [ ] Lead sequence: `trace_tasks` → `trace_context`(optional `query`) → `trace_loop` → `trace_review` → `trace_plan`
- [ ] Gate path documented (`trace_loop action=gate`)
- [ ] 9/16 hygiene documented + Instructions mention `trace_version`/reload

### B — M-001 / rejects

- [ ] All 16 write/read tools still registered (`TestToolNamesRegistered`)
- [ ] No CG 1-tool-only facade or “use only trace_context” language
- [ ] No MP-style tool explosion or new MCP tool
- [ ] Task/write path not hidden; write tools remain in registration

### C — G2 compose-first (partial — full explore deferred)

- [ ] Instructions include compose-first read recipe (search → why → impact → capability)
- [ ] No `trace_explore` implement snuck in
- [ ] Codegraph mentioned as optional complement, not bundled MCP

### D — Law 19

- [ ] Instructions are adapter-layer string content; no domain logic fork in MCP handlers
- [ ] `internal/mcp/instructions.go` is presentation only

### E — Tests + docs

- [ ] G3-A1–A6 evidenced (run tests or read test bodies)
- [ ] CONTRIBUTING reload/orient section accurate vs live install paths
- [ ] Optional install stderr tweaks consistent with Instructions (not contradictory)

### F — G1 integration (S00 shipped)

- [ ] Orient recipe mentions optional `query` on `trace_context` (not “after G1” defer language)
- [ ] No regression to G1 compiler/MCP query wiring

## Live verification commands

```bash
go test ./internal/mcp/... ./internal/install/... -count=1 -run 'ToolNames|ServerInstructions'
grep -c 'AddTool(s.mcp' internal/mcp/server.go                    # expect 16
grep -n 'ServerOptions\|ServerInstructions\|Instructions:' internal/mcp/server.go internal/mcp/instructions.go
grep -E 'trace_tasks|trace_context|trace_loop|trace_review|trace_plan|trace_version' internal/mcp/instructions.go
grep -E 'trace_version|reload|9/16|stale' CONTRIBUTING.md internal/mcp/instructions.go
wc -l internal/mcp/instructions.go                                  # expect non-zero
```

## Exit criteria

- [ ] APPROVE (medium+ confidence) or spawn repair pair
- [ ] Board row → `done` with verdict + findings counts in Notes

## Next

`P39-S02-00`
