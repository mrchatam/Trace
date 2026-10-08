# Phase 45 / S05 / Planner — Store-free reads

## Metadata
- id: P45-S05-00
- todo_ids: [P45-S05-00]
- role: planner
- skills: [api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan the `trace_version` fix (#152) and the `trace_tasks` ordering documentation (#160). **No product code in this row.**

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md)
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)
- Live targets:
  - `internal/mcp/tools_parity.go:314` (`toolVersion` calls `openStore` + the capability gate)
  - `internal/store/open.go:228` (`nowRFC3339` — second granularity)
  - `internal/store/helpers.go:112` (`ListTasksFiltered` ordering)
  - `internal/mcp/mcp_test.go` (`TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestMCPVirginProjectDoesNotMkdir`)

## Locked defaults

| Item | Value |
|------|-------|
| `trace_version` | Answers from the **process alone**: no `openStore`, no capability gate, no mkdir |
| Response | `buildmeta.Payload` + `project_initialized` + `project_root` |
| `project_initialized` | **Stat-only** probe (`store.ProjectInitialized`) — never opens, locks or creates |
| #160 | **Document the ordering** (`created_at ASC, id ASC`); no nanosecond column, no migration |
| Gate behaviour change | `trace_version` no longer runs `assertMCPToolAllowed` — a locked/virgin project still answers (that is the point of #152) |

## Why stat-only for `project_initialized`

The documented purpose of `trace_version` is stale-server detection: an agent calls it to check whether the running stdio process matches the freshly installed binary — typically **after** a rebuild, sometimes on a project that is locked by another process or not yet initialized. Any implementation that opens (or worse, creates) the store defeats that purpose and can hang on a lock or pollute a virgin directory. `os.Stat` on the DB path answers the real question ("is there a project here at all?") with zero side effects.

## Why document-not-migrate for #160

The UUID tiebreak is deterministic and stable — the defect is purely that nothing says so. A nanosecond `created_at` would rewrite history for every existing row (all identical to the second), cost a migration, and change stored formats for a convenience. Documenting the contract in the tool description and schema tag fixes the actual honesty gap.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S05-01**

## Next

`P45-S05-01`