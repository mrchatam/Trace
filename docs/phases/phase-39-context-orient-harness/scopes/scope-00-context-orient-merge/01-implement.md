# P39-S00-01 — Implement G1 context orient merge

## Metadata
- id: P39-S00-01
- todo_ids: [P39-S00-01]
- role: implementer
- skills: [incremental-implementation, context-engineering, test-driven-development]
- mcps: [user-trace]
- verification: automated

## Objective

Implement **G1**: optional agent `query` on context/compiler path merges relevant FTS hits into the task packet without abandoning moat (G-001, G-002).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G1](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-001/G-002](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Live anchors (P39-S00-00 verified 2026-08-22):
  - `internal/compiler/compiler.go:14–19` — `ContextOptions` (no `Query` yet)
  - `internal/compiler/compiler.go:77–177` — `compileAtDepth` candidate pipeline
  - `internal/compiler/compiler.go:146–154` — title FTS only (`Search` on `task.Title`, Limit 16)
  - `internal/compiler/compiler.go:156–176` — file-seed expand (query merge **before** this block)
  - `internal/compiler/packet.go:18–20` — `DefaultTokenBudget=4096`, `DefaultMaxItems=32`, `MaxCandidateHits=64`
  - `internal/retrieval/types.go:13–14` — `ReasonFTSMatch`, `ReasonDirectTaskScope`
  - `cmd/trace/context.go:18–68` — no `--query`; `flagsFirst` map at :22–24
  - `internal/mcp/tools_context.go:14–21` — `ContextInput` (no `query` field)
  - `internal/mcp/server.go:67–75`, `:227–235` — `trace_context` desc; `RegisteredToolNames()` = 16

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-001, G-002 |
| Verdict | **Accept** — merge query+task |
| Query field | `ContextOptions.Query string`; CLI `--query`; MCP `query` optional json field |
| Merge point | After title FTS (`compiler.go:148`), before file-seed expand — append query FTS hits to `candidates` |
| Search opts | `retrieval.SearchOptions{Limit: 16}` — same bound as title FTS |
| Search errors | Fail-open like title FTS (DF-87): query search error → skip query hits, still emit packet |
| Layer 0 | Unchanged — task, task_state, goal always present when applicable |
| Reason code | `ReasonFTSMatch` on query hits |
| Must preserve | Task scope, gates, progressive caps, `Budget.Truncated`, index/graph honesty |
| Must not | Query-only drift; full-graph dump; G2 explore; new MCP tool |
| MCP tools | **16** — extend `ContextInput` schema only |

## Touch-list (library → CLI → MCP → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `internal/compiler/compiler.go` | Add `Query string` to `ContextOptions`; in `compileAtDepth`, if trimmed query non-empty, `retr.Search(ctx, query, …)` and append to candidates |
| 2 | `internal/compiler/compiler_test.go` | Add G1 acceptance tests (see below) |
| 3 | `cmd/trace/context.go` | Add `--query` string flag; add `"query": true` to `flagsFirst` map; pass to `ContextOptions.Query`; update usage string |
| 4 | *(none)* | No `context_test.go` — compiler + MCP tests cover wiring (Law 19 thin adapters) |
| 5 | `internal/mcp/tools_context.go` | Add `Query string` with `json:"query,omitempty"` + jsonschema; pass to `ContextOptions.Query` in `toolContext` |
| 6 | `internal/mcp/mcp_test.go` | MCP test: `trace_context` with `query` returns merged hit (reuse `callContext` seam) |
| 7 | `internal/mcp/server.go` | Update `trace_context` description to mention optional `query`; **do not** add `AddTool` — count stays 16 |

**Explicit non-touch:**

- `internal/retrieval/search.go` — no semantic channel
- `MaxCandidateHits` / `DefaultMaxItems` defaults
- New `trace_explore` tool
- `web/` — out of scope

## Implementation order

```text
1. ContextOptions.Query + compileAtDepth merge (compiler.go)
2. compiler_test.go — acceptance tests T1–T6
3. cmd/trace/context.go — --query flag
4. internal/mcp/tools_context.go + mcp_test.go
5. server.go description tweak
6. go test ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1
7. trace seed export if graph entities changed
```

## Acceptance tests (must pass)

Add in `internal/compiler/compiler_test.go` (names suggested; adjust if colliding):

| ID | Suggested name | Assert |
|----|----------------|--------|
| T1 | `TestG1QueryHitMerged` | Create indexed decision with unique token; `ContextOptions{Query: token}` includes it; empty Query excludes it |
| T2 | `TestG1TaskMoatPreserved` | Layer 0 `task` + `task_state` present with Query set; MCP rejects empty `task_id` (existing + extend) |
| T3 | `TestG1TitleFTSStillRunsWithQuery` | Task title matches entity A; query matches entity B; both appear when both channels hit |
| T4 | `TestG1QueryExpandDedupe` | Same entity reachable via expand and query FTS → one packet item |
| T5 | `TestG1QueryCapHonesty` | Many query matches → `Budget.CandidatesCapped` / `Truncated` honest (mirror `TestCandidateCapSetsTruncated`) |
| T6 | `TestG1QuerySearchFailOpen` | Retriever stub: fail `Search` only when query ≠ `task.Title` (agent query path); packet still valid with Layer 0 |

MCP mirror in `internal/mcp/mcp_test.go`:

| ID | Suggested name | Assert |
|----|----------------|--------|
| T1-MCP | `TestMCPContextQueryMerged` | `callContext` with `Query` field returns hit absent without query |

## Regression tests (must stay green)

Run existing suite unchanged — do **not** weaken:

- `TestTaskContextContinuesWhenSearchErrors` — title FTS fail-open (DF-87)
- `TestCandidateCapSetsTruncated`, `TestItemCapNeverExceeded`, `TestBudgetLoudTotals`
- `TestNoDumpAPI` — no dump surface added
- `internal/mcp/mcp_test.go` tool-count assertion (`RegisteredToolNames()` = 16)

## Role work

1. Add `Query` to `ContextOptions` (no behavior change when empty — backward compatible).
2. In `compileAtDepth` after title FTS append (`compiler.go:~154`), **before** file-seed expand (`:~156`):
   ```go
   q := strings.TrimSpace(opts.Query)
   if q != "" {
       qfts, err := c.retr.Search(ctx, q, retrieval.SearchOptions{Limit: 16})
       if err == nil {
           candidates = append(candidates, qfts...)
       }
       // Fail-open like title FTS (DF-87): query search error → skip query hits
   }
   ```
   Add `"strings"` import if not present.
3. Re-run existing sort/dedupe/cap pipeline — no fork.
4. Wire CLI `--query` and MCP `query` as thin adapters (Law 19).
5. Self-check all six acceptance tests before marking row done.

## Exit criteria

- [ ] T1–T6 green (`go test ./internal/compiler/... ./internal/mcp/... -count=1`)
- [ ] No Law 6–7 regression (existing compiler cap tests green)
- [ ] Board row → `done` with Notes listing files + test command

## Next

`P39-S00-02`
