# Phase 45 / S04 / Planner — Envelope + schema + descriptions

## Metadata
- id: P45-S04-00
- todo_ids: [P45-S04-00]
- role: planner
- skills: [api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan the uniform list envelope (#157), schema enums/ranges (#158), generated descriptions with a drift guard (#159), and the snake_case `trace_changes show` DTO (#161). **No product code in this row.**

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md) — SDK facts 1–3
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)
- Live targets:
  - `internal/mcp/tools_changes.go` (bare `store.Change` marshal)
  - `internal/mcp/tools_agents.go:60` (missing empty-catalog hint on recommend)
  - `internal/mcp/tools_parity.go` (trace_tasks without `ok`; trace_changes list without limit)
  - `internal/mcp/server.go` (17 tools, inferred schemas)
  - `internal/mcp/mcp_test.go` (3 tests asserting the old bare-array shapes)

## Locked defaults

| Item | Value |
|------|-------|
| List envelope | `{ok, items, count, truncated, limit?, next_cursor?, hint?}` — **one** shape |
| Envelope scope | Generic **history** lists (changes, regressions, agents, tasks). Entity-specific collections (`capabilities`, `missing`, `findings`) keep domain-named keys — they mirror the CLI output documents |
| Schema mechanism | `schemaFor[In](constraints)`: infer via `jsonschema.For[In]`, apply enum/min/max in code, panic at construction on a missing property or unresolvable schema |
| Description mechanism | Generated from the **same** Go vocabularies the enums come from; a drift guard fails the build on divergence |
| #161 | Tagged `changeShowRow` DTO; `RowID` (internal SQLite rowid) deliberately not on the wire |
| Ranges | Only where the code enforces them (explore/context/impact bounds, limit caps). **No decorative min/max** |

## Why panics at construction

A schema that does not resolve, or a constraint naming a property that does not exist, is a programming error — not a runtime condition. Failing at `NewServer` means every test that constructs the server catches it immediately; a lazy failure would surface as a malformed tool listing mid-session for an agent.

## Why generated descriptions

The drift in #159 (`trace_loop` description omitting `gate`; capability kinds omitting `AGENT`) happened because two humans wrote two copies of the same fact. `join(vocabulary)` used in both the schema enum and the `fmt.Sprintf` description makes one copy, and `TestToolDescriptionsCoverSchemaEnums` enforces it.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S04-01**

## Next

`P45-S04-01`