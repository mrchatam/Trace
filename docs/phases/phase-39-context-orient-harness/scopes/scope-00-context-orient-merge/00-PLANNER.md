# P39-S00-00 — Scope planner (G1 context orient merge)

## Metadata
- id: P39-S00-00
- todo_ids: [P39-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S00 **G1** against live repo: optional agent `query` merges into task context packet (G-001, G-002). Thicken `01-implement.md` + `02-review.md` with file targets, acceptance tests, and rejects. **No product code in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P39-00 resolutions
- [REMEDIATION-PLAN §2 G1](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-001/G-002](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [TRACE-AUDIT H1/H2](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-01-trace-audit/TRACE-AUDIT.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P39-S00-00):
  - `internal/compiler/compiler.go:14–19` — `ContextOptions` (no `Query` yet)
  - `internal/compiler/compiler.go:146–154` — FTS on `task.Title` only (Limit 16)
  - `internal/compiler/packet.go:18–20` — caps: TokenBudget 4096, MaxItems 32, MaxCandidateHits 64
  - `internal/retrieval/types.go:13–14` — `ReasonFTSMatch`, `ReasonDirectTaskScope`
  - `cmd/trace/context.go:18–20` — no `--query` flag
  - `internal/mcp/tools_context.go:14–21` — `ContextInput` has no `query`
  - `internal/mcp/server.go:67–75` — `trace_context` description

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: INTAKE + REMEDIATION-PLAN are authority.

## Locked defaults (FINAL — P39-00)

| Item | Value |
|------|-------|
| GAP ids | G-001, G-002 |
| Verdict | **Accept** per REMEDIATION-PLAN G1 |
| Peer pattern | UA `buildChatContext(query)` merge; MP `wake_up()` identity preserved |
| Moat (M-001) | `task_id` **required**; Layer 0 task core always present; gates/reason_codes unchanged |
| Query semantics | Optional `query` **adds** FTS hits to candidate pool; title FTS **still runs** when query set |
| Reason codes | Query hits use existing `fts_match`; no new reason_code unless documented in `retrieval/doc.go` |
| Law 6–7 | Existing `MaxCandidateHits`, `trimToBudget`, `Budget.Truncated` honesty — no cap increase |
| Out | Unified `trace_explore` (G2); query-only packet; semantic/vector (G-004a) |
| MCP tool count | **16** unchanged — extend `trace_context` schema only |
| Graph export | If entities change: `trace seed export -o trace/graph.json` |

## Accept / reject (G1)

| Decision | Item |
|----------|------|
| **Accept** | Optional `query` on `ContextOptions`, CLI `--query`, MCP `trace_context.query` |
| **Accept** | Merge query FTS hits into `compileAtDepth` candidates (dedupe with title/expand hits) |
| **Accept** | Tests proving query hit appears in packet when indexed; task Layer 0 unchanged |
| **Reject** | Query-only orient (no task_id) |
| **Reject** | Replace title FTS with query-only |
| **Reject** | Raise MaxCandidateHits / DefaultMaxItems defaults |
| **Reject** | New MCP tool for orient merge |

## Must lock for S00-01 (delivered in thickened 01-implement)

1. Touch-list: compiler → CLI → MCP → tests (library first, Law 19).
2. Six acceptance tests (see `01-implement.md`).
3. Regression: existing compiler cap/honesty tests green.

## Exit criteria

- [x] `01-implement.md` + `02-review.md` runnable with file targets + tests
- [x] SCOPE-TODOS updated
- [x] Board row → `done` with Notes

## Next

`P39-S00-01`
