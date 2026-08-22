# P40-S01-01 — Implement (G2 unified explore)

## Metadata
- id: P40-S01-01
- todo_ids: [P40-S01-01]
- role: implementer
- skills: [backend-dev, api-and-interface-design, test-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Implement **G2**: unified task-aware capped `trace_explore` — library compose + MCP read tool (G-007). M-001: merges into task loop; never query-only replacement.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G2](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-007](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Live anchors (P40-S01-00 re-verified 2026-08-22):
  - `internal/compiler/compiler.go:158–165` — G1 query merge (reuse via `compiler.ContextOptions.Query`)
  - `internal/compiler/packet.go:12–20` — caps `DefaultTokenBudget=4096`, `DefaultMaxItems=32`, `MaxCandidateHits=64`
  - `internal/retrieval/search.go` — FTS (`SearchOptions.Limit`; MCP default 32, cap 64)
  - `internal/retrieval/why.go` — causal neighborhood per hit
  - `internal/retrieval/neighborhood.go:8–17,53+` — `BoundedGraph.Truncated`, `MaxNeighborhoodNodes=5000`
  - `internal/mcp/tools_context.go:14–22,60–64` — `ContextInput` with optional `query`; task_id required pattern
  - `internal/mcp/tools_search.go:14–18` — limit defaults for explore search slice
  - `internal/mcp/tools_why.go:15–20` — why slice on top hits
  - `internal/mcp/server.go:54–227` — **16** `AddTool` registrations; **no** `trace_explore`
  - `internal/mcp/server.go:229–237` — `RegisteredToolNames()` returns 16 names
  - `internal/mcp/instructions.go:5–16` — moat-first + compose-first recipe (no explore yet)
  - `internal/mcp/instructions.go:18–21` — stale hygiene `9/16` → update to `9/17`
  - `internal/mcp/mcp_test.go:443–464` — `TestToolNamesRegistered` exact 16-name order lock
  - `internal/mcp/mcp_test.go:1385–1388` — `TestRegisteredToolNames_IncludesTracePlan` count=16
  - `internal/mcp/mcp_test.go:514–520` — `TestServerInstructionsStaleHygiene` checks `9/16`
  - `internal/httpapi/handlers_retrieval.go` — context/search/why/graph mirrors (HTTP explore **out of scope**)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-007 |
| Verdict | **Accept** — add `trace_explore` |
| Required input | `task_id` (UUID); optional `query` string |
| Optional inputs | `limit`, `max_nodes`, `depth` (bounded — mirror existing defaults) |
| Response shape | JSON struct e.g. `ExploreResult{TaskSummary, PacketBudget, SearchHits, WhySlices[], Neighborhood?, Truncated}` — task packet section **required** |
| Merge strategy | `compiler.TaskContext`/`ExpandContext` (depth 1 default, with `Query`) → title/query FTS hits → why on top N hits (default 3) → optional neighborhood on task entity |
| Default caps | Search limit 32 (cap 64); why top-N=3; neighborhood max_nodes default 100 depth 2; honor `packet.go` 4096/32/64 for task section |
| MCP annotations | `ReadOnlyHint: true`; not OpenWorld |
| Tool count | **17** after add — update `RegisteredToolNames()` + `TestToolNamesRegistered` |
| Instructions | Add explore bullet **after** moat + compose-first blocks — “optional convenience” |
| Must not | Query-only path; hide write tools; semantic/vector; dump API |
| Graph export | If entities change: `trace seed export -o trace/graph.json` |

## Touch-list (library → MCP → docs → tests) — 16→17 migration

| Step | File | Action |
|------|------|--------|
| 1 | `internal/retrieval/explore.go` (new) | `Explore(ctx, ExploreOpts)` — task-aware capped compose; calls `compiler` for task packet + retrieval for search/why/neighborhood |
| 2 | `internal/retrieval/explore_test.go` (new) | G2-T1–T7 acceptance tests |
| 3 | `internal/mcp/tools_explore.go` (new) | `ExploreInput{task_id, query?, limit?, max_nodes?, depth?}`; `toolExplore` thin adapter (mirror `tools_context.go` openStore/VCS/assert pattern) |
| 4 | `internal/mcp/server.go:217–227` | 17th `AddTool` for `trace_explore` after `trace_plan`; `ReadOnlyHint: true` |
| 5 | `internal/mcp/server.go:229–237` | Append `"trace_explore"` to `RegisteredToolNames()`; update comment sixteen→seventeen |
| 6 | `internal/mcp/instructions.go` | Add explore bullet **after** compose-first block — “optional convenience, not moat replacement”; update stale hygiene `9/16`→`9/17` |
| 7 | `internal/mcp/mcp_test.go:443–464` | `TestToolNamesRegistered`: want slice + `len(names)==17` |
| 8 | `internal/mcp/mcp_test.go:1385–1388` | `TestRegisteredToolNames_IncludesTracePlan`: count 17 |
| 9 | `internal/mcp/mcp_test.go` (new) | `TestMCPExploreTaskRequired`, `TestMCPExploreQueryMerged`, optional `TestServerInstructionsExploreOptional` |
| 10 | `internal/mcp/mcp_test.go:514–520` | `TestServerInstructionsStaleHygiene`: token `9/17` (replace `9/16`) |
| 11 | `cmd/trace/explore.go` (optional) | CLI mirror if low cost — not blocking |

**Registration order lock (append only):**

```text
trace_why, trace_context, trace_add, trace_link, trace_transition, trace_review,
trace_tasks, trace_capability, trace_impact, trace_version, trace_search,
trace_changes, trace_regressions, trace_loop, trace_agents, trace_plan, trace_explore
```

**Explicit non-touch:**

- Existing 16 tool handler bodies (except `RegisteredToolNames` comment)
- Cap default increases in `packet.go`
- `web/`, `internal/httpapi/` — no HTTP explore route (S00 GUI done; G2 is MCP-first)
- G-004a semantic/vector channel
- Write tools (`trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_capability`, `trace_impact`, `trace_loop`, `trace_plan`) — must remain registered

## Implementation order

```text
1. Library Explore compose + unit tests G2-T1–T7
2. MCP tools_explore.go + server registration
3. Update RegisteredToolNames + mcp_test tool count
4. ServerInstructions addendum
5. Optional CLI explore
6. go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1
7. trace seed export if entity schema changes
```

## Acceptance tests (must pass)

| ID | Suggested name | Assert |
|----|----------------|--------|
| G2-T1 | `TestExploreTaskRequired` | Empty task_id → error (library + MCP) |
| G2-T2 | `TestExploreTaskMoatPreserved` | Response includes task/Layer-0 identity |
| G2-T3 | `TestExploreQueryMerged` | Optional query adds hits via G1 merge path |
| G2-T4 | `TestExploreCappedHonest` | Over-budget → `truncated` / cap flags set |
| G2-T5 | `TestExploreNoDump` | Cannot fetch unbounded full graph |
| G2-T6 | `TestExploreWhyIncluded` | Top hit includes why/neighborhood slice (bounded) |
| G2-T7 | `TestExploreFailOpenSearch` | Search error → partial result, task moat intact |

MCP mirrors in `internal/mcp/mcp_test.go`:

| ID | Suggested name | Assert |
|----|----------------|--------|
| G2-T1-MCP | `TestMCPExploreTaskRequired` | MCP rejects missing `task_id` |
| G2-T3-MCP | `TestMCPExploreQueryMerged` | MCP `query` field merges hits |

## Regression tests (must stay green)

- All G1 tests (`TestG1*`)
- `TestNoDumpAPI`
- `TestServerInstructionsMoatLead` — moat ordering unchanged
- Compiler cap tests unchanged

## Role work

1. Implement library explore as **composition** of existing retrieval/compiler — no parallel index.
2. Wire MCP as thin adapter (Law 19).
3. Update tool count tests 16 → 17 explicitly.
4. Update Instructions without demoting task loop lead.
5. Self-check G2-T1–T7 before marking row done.

## Exit criteria

- [ ] G2-T1–T7 + MCP mirrors green
- [ ] 17 tools registered; Instructions updated
- [ ] Board row → `done` with files + test command in Notes

## Next

`P40-S01-02`
