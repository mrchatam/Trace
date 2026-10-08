# Phase 45 / S03 / Planner — Error quality

## Metadata
- id: P45-S03-00
- todo_ids: [P45-S03-00]
- role: planner
- skills: [api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan the removal of driver text and wrong entity-type names from every MCP error (#155), and give `trace_loop action=next` the same structured envelope `action=gate` already returns (#156). **No product code in this row.**

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md)
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)
- Live targets:
  - `internal/mcp/tools_{link,why,review,transition,changes}.go`
  - `internal/mcp/tools_loop.go` — `mcpGateSchemaVersion = "trace.loop.gate.v1"` (`:30`)
  - `internal/loop/next.go:169` — where the precondition fails today
  - `internal/store/errors.go` (S02) — the typed error this scope renders

## Locked defaults

| Item | Value |
|------|-------|
| Driver text | **Never** leaves the MCP boundary (`sql: no rows…`, driver type names) |
| Entity naming | Name the entity type **actually queried** |
| `trace_link` | Name the `from`/`to` parameter that was wrong + the rel's expectations |
| `loop.next` | Returns the **existing** `trace.loop.gate.v1` envelope — no new schema version |
| New `reason_code`s | `goal_link_missing`, `seed_task_not_found` |
| Scope of the wrapper | One helper, applied at the tool boundary — not a global `Error()` rewrite |

## Why reuse the gate envelope (#156)

`action=gate` and `action=next` are the same tool. An agent that parses `trace.loop.gate.v1` for one action should not need a second parser for the other, and a precondition failure **is** a gate result — the loop is refusing to advance. Returning the gate envelope with `for: "next"` and a `reason_code` is both honest and the cheapest thing an agent can consume.

Adding a `trace.loop.next.v1` schema would have been the alternative. Rejected: two schemas for one tool, and a migration burden for every client, for no new information.

## The `from` vs `to` problem

`trace_link`'s internal lookups run in a fixed order, so the position of the failure does not tell you which parameter the caller got wrong. `wrapLinkErr` therefore compares the missing id against the supplied `from`/`to` and names the parameter that matches — and prints the rel's expected shape so the caller can see the rule they violated.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S03-01**

## Next

`P45-S03-01`