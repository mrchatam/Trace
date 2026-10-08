# Phase 45 / S02 / Planner — Write-path validation

## Metadata
- id: P45-S02-00
- todo_ids: [P45-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan fail-closed validation for the three write surfaces that currently accept garbage and report success: `trace_add` (#150, #151), `trace_tasks` filters (#154), `trace_link` self-references (#162). **No product code in this row.**

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md) — locked decision **D2**
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)
- Live targets:
  - `internal/mcp/tools_write.go:69` (`toolAdd`)
  - `internal/mcp/tools_parity.go` (`toolTasks`)
  - `internal/domain/scope.go:213` (`LinkBlocks`), `:161` (`LinkAPIContract`)
  - `internal/domain/task_state.go`
  - `internal/store/entities.go:121` and the other ~15 getters
  - `internal/httpapi/helpers.go`, `internal/domain/scope.go` (`isNotFoundErr`)

## Locked defaults

| Item | Value |
|------|-------|
| Invalid value | **Fail closed**, name the offending parameter |
| `trace_add` id collision | Refuse unless `replace=true`; response gains `created: bool` |
| `replace=true` scope | Entity row only — **no cascade** to links/changes/evidence |
| `goal_id` | Validated as a UUID **and** resolved via `GetGoal` |
| `work_state` | Closed 8-value vocabulary, indexed error for `work_states[i]` |
| Self-link | Refused **before** any lookup or write, at the domain layer |
| Missing entity | Typed `store.NotFoundError`, never raw `sql.ErrNoRows` |

## The three failure modes

| # | Today | After |
|---|-------|-------|
| #150 | `trace_add` with an existing id overwrites, returns `ok:true` | Refuses; `replace=true` opt-in; `created:false` reported |
| #151 | `id="not-a-uuid"`, `goal_id="nope"` accepted → dangling rows | Both validated before insert |
| #154 | `work_state="TYPO"` matches nothing, returns a confident empty list | Error naming `work_state` (or `work_states[3]`) |
| #162 | `rel=blocks`, `from == to` persists a task blocking itself | Refused in `LinkBlocks` **and** `LinkAPIContract` |

## Where the guard belongs

**Domain, not MCP.** `refuseSelfLink` lives in `internal/domain/scope.go` and is called from `LinkBlocks` (`:239`) and `LinkAPIContract` (`:161`) — so the CLI, HTTP and MCP adapters are all covered by construction. Same for `NormalizeWorkState` and `EntityExists`: `internal/mcp` calls them, but nothing *only* in `internal/mcp` enforces the rule.

The `isNotFoundErr` helpers in `internal/httpapi` and `internal/domain` were updated to check `store.IsNotFound` **first**, with the legacy string sniff retained as a fallback — so a typed error is never misread as "unknown".

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S02-01**

## Next

`P45-S02-01`