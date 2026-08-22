# P40-S01-02 — Review (G2 unified explore)

## Metadata
- id: P40-S01-02
- todo_ids: [P40-S01-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening, silent-failure-hunter]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Fresh independent review of S01-01 G2 implementation vs REMEDIATION-PLAN G2, GAP-REGISTRY G-007, M-001 moat, Laws 6–7/19.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G2-T1–T7 acceptance map + 16→17 touch-list
- [00-PLANNER.md](00-PLANNER.md) — locks + law spike desk-check (waived)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — desk-check table
- [REMEDIATION-PLAN G2](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-007](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Live pre-ship baseline (P40-S01-00): `server.go:229–237` (16 tools), `instructions.go:13–16` (compose-first, no explore)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b below this row |
| Tool count | **17** — verify `RegisteredToolNames()` length + `AddTool` count |
| Cap defaults | 4096/32/64 packet defaults unchanged unless explicit scoped explore caps |
| Pre-Phase-40 | Was 16 tools (`server.go:59–226` AddTool ×16), no `trace_explore` |
| Stale hygiene | Instructions + test must say `9/17` not `9/16` after ship |

## Review checklist

### A — G2 gap closure

- [ ] `trace_explore` registered with ReadOnlyHint
- [ ] `task_id` required (G2-T1 / G2-T1-MCP)
- [ ] Optional `query` merges via G1 path (G2-T3)
- [ ] Capped response with truncation honesty (G2-T4)
- [ ] Why/neighborhood slice included bounded (G2-T6)

### B — M-001 moat

- [ ] Task identity present in explore response (G2-T2)
- [ ] No query-only code path
- [ ] Write tools still registered and documented — not hidden behind explore
- [ ] Instructions moat-first lead unchanged (`tasks` → `context` → `loop` → `review`)

### C — Laws 6–7

- [ ] No full-graph dump (G2-T5, `TestNoDumpAPI`)
- [ ] Fail-open search behavior (G2-T7)
- [ ] No default cap inflation on core packet defaults

### D — Law 19

- [ ] Compose logic in library (`internal/retrieval/` or `internal/compiler/`)
- [ ] MCP handler thin — schema pass-through + marshal

### E — Tests

- [ ] G2-T1–T7 + MCP mirrors evidenced
- [ ] G1 regression green
- [ ] `go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... -count=1` passes

### F — Rejects

- [ ] Not CG 1-tool facade (16 write/read tools still available)
- [ ] No semantic/vector channel
- [ ] No claim h7-compose-desk-check equivalence

### G — Live verification commands

```bash
go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... -count=1 \
  -run 'Explore|G1|NoDump|ToolNames|ServerInstructions'

grep -c 'AddTool(s.mcp' internal/mcp/server.go          # expect 17
grep 'trace_explore' internal/mcp/server.go internal/mcp/tools_explore.go
grep -A5 'start here' internal/mcp/instructions.go      # moat lead intact
grep 'trace_explore' internal/mcp/instructions.go       # optional convenience mention
grep '9/17' internal/mcp/instructions.go                # stale hygiene updated

# Write tools still present (must not hide behind explore)
grep -E 'trace_add|trace_transition|trace_review|trace_loop' internal/mcp/server.go | wc -l  # expect 4+
```

### H — Law spike desk-check (waived at planner — verify at review)

| Check | Evidence |
|-------|----------|
| Task required | G2-T1 / G2-T1-MCP |
| Query optional via G1 | G2-T3 / G2-T3-MCP + `compiler.go:158–165` reuse |
| Caps honest | G2-T4 |
| Write surface visible | 17 tools; write handlers unchanged |
| Not CG facade | Task packet section in explore response; compose recipe still in instructions |
| Compose-first preserved | `TestServerInstructionsComposeRecipe` green; explore after compose block |

## Exit criteria

- [ ] APPROVE or spawn with evidence
- [ ] M-001 + Law 6–7/19 verified
- [ ] Board row → `done` with verdict + confidence in Notes

## Next

`P40-S02-00`
