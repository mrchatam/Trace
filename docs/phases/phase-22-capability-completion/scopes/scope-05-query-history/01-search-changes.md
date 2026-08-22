# P22-S05-01 — Implement: search + changes history

## Metadata
- id: P22-S05-01
- todo_ids: [P22-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Wire FTS and change history to **CLI + stdio MCP**. Closes **C29, C30, C37**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/retrieval/search.go`, `internal/store/fts.go`, `internal/store/changes.go` (`ListAllChanges`, `GetChange`, `ListChangePaths`)
- S02: `CompareStates` already on `trace changes compare` — **do not reimplement**
- G19: query logic in domain/retrieval; CLI/MCP encode JSON only

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `Engine.Search` + `SearchFTS` (compiler-internal) | `trace search` CLI |
| `trace changes capture\|compare` | `list\|show` subcommands |
| `ListAllChanges` (unbounded ASC) | bounded newest-first list helper |
| MCP catalog **10** | `trace_search`, `trace_changes` |
| Compat/schema **24** | **025+** |

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **None**; compat stays **24** |
| Search CLI | `trace search <query> [--limit N]` — default 32, cap 64; stdout JSON `{ok:true,hits:[],count}`; hits = `retrieval.Hit` fields (entity_type, entity_id, title, excerpt, path, reason_code, score, distance) |
| Search errors | Empty query → exit 0 + `{ok:true,hits:[],count:0}`; FTS error → stderr + non-zero |
| Changes list | `trace changes list [--task <uuid>] [--limit N]` — **newest first** (`created_at DESC, id DESC`); default 32, cap 64; rows: id, task_id, git_commit, status, reason, source_type, created_at (no path blobs inline) |
| Changes show | `trace changes show <change-id>` — `{ok:true,change:{...},paths:[{path,status,symbol_id}]}`; missing id fail-closed |
| List helper | Add `domain.ListChanges` or store `ListChangesRecent(limit, taskID)` — **do not** load unbounded `ListAllChanges` in CLI |
| Capability | `cli:search` AUTO_ALLOW; extend `cli:changes` gating for `list`/`show` |
| MCP `trace_search` | Args: `query` (required), `limit` (optional); mirrors search CLI JSON |
| MCP `trace_changes` | Args: `action` = `list\|show\|compare`; `list`: optional `task_id`, `limit`; `show`: `change_id`; `compare`: `from`, `to` — delegate to existing domain compare |
| MCP catalog | **12** tools after this row; append to `RegisteredToolNames` in registration order after `trace_version` |
| C37 | Search hits include indexed historical entities (changes, regressions, outcome_results, reflections in FTS) — `TestCLISearchUsesFTS` seeds at least one non-task entity |

## Requirements

1. **`cmd/trace/search.go`** — new root subcommand; wire `retrieval.New(store).Search`.
2. **Extend `cmd/trace/changes.go`** — add `list`, `show`; update usage strings + `help.go`.
3. **`internal/mcp/tools_search.go`** (new) + extend `server.go` registerTools.
4. **`internal/mcp/tools_changes.go`** (new) or combine with search file if tiny — action dispatch.
5. **`internal/domain/capability.go`** — add `cli:search` to `BuiltinCLICapabilitySpecs`; add `trace_search`, `trace_changes` to `BuiltinMCPCapabilitySpecs`.
6. **`cmd/trace-mcp/main.go`** help text — list new tools if it enumerates catalog.
7. Named tests below; checklist C29/C30/C37 **unboxed** (reviewer boxes S05-02).

## Touch files

- `cmd/trace/search.go`, `search_test.go` (new)
- `cmd/trace/changes.go`, `changes_test.go` (extend)
- `cmd/trace/root.go`, `help.go`
- `internal/domain/changes.go` or `internal/store/changes.go` (bounded list query)
- `internal/mcp/server.go`, `tools_search.go`, `tools_changes.go` (new), `mcp_test.go`
- `internal/domain/capability.go`
- `cmd/trace-mcp/main.go` (help only if needed)

## Named tests

| Test | Proves |
|------|--------|
| `TestCLISearchUsesFTS` | Seeds task + change/regression token; search returns `reason_code=fts_match` hit |
| `TestCLIChangesList` | Two changes; list returns newest-first bounded rows |
| `TestCLIChangesShow` | show returns paths without file content |
| `TestMCPSearchRegistered` | `trace_search` in catalog + CallTool returns hits |
| `TestMCPChangesList` | `trace_changes` action=list returns JSON array |
| `TestToolNamesRegistered` | exactly **12** names, stable order |
| `TestChangesCompare` | keeper (S02) — compare still PASS |

```bash
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestCLISearch|TestCLIChanges|TestMCPSearch|TestMCPChanges|TestToolNamesRegistered|TestChangesCompare'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C29, C30, C37 true (evidence via named tests)
- [ ] MCP catalog **12**; compat **24** unchanged
- [ ] Checklist caps **unboxed** until S05-02
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Bounded changes list + show CLI
- [ ] Search CLI via retrieval.Engine
- [ ] MCP trace_search + trace_changes
- [ ] Capability specs + tests
- [ ] Board status + notes

## Residual risks (carry to S05-02)

- `ListAllChanges` ASC vs list DESC — ensure store query is explicit, not sort-in-CLI only
- MCP compare must reuse domain `CompareStates`, not duplicate git diff
- Empty-task filter on list must return [] not error
