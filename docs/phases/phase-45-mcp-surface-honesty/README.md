# Phase 45 — MCP surface honesty

Human-promoted successor to Phase 44. Source: audit round-5 MCP findings, GitHub issues **#149–#162** (all OPEN at implementation time).

## Goal

Make the MCP tool surface say what it actually does. Fourteen findings, one defect class: **a tool told the agent something untrue, and the agent had no way to find out.** Silent AND-ing of FTS operators, silent overwrite on `trace_add`, silent acceptance of an invalid `work_state` filter, raw `sql: no rows in result set` leaking to the model, tool descriptions that had drifted from the schemas they describe.

The phase makes each surface either honest or **loud**:

| Class | Policy |
|-------|--------|
| Something is silently narrowed | Report the narrowing (`graph_sync_honesty`, `query_note`, `truncated`, `fallback`) |
| Something would be silently lost | **Fail closed** — refuse rather than degrade (`trace_add` on an existing id) |
| Something is invalid | **Fail closed** with a message naming the offending parameter |
| Something is undocumented | One shared vocabulary generates both the JSON-Schema enum *and* the description text |

## Issues in scope

| # | Title (verbatim) | Scope |
|---|------------------|-------|
| [#149](https://github.com/mrchatam/Trace/issues/149) | `trace_search` silently AND-s FTS operators (OR/AND/NOT/NEAR/*) → false-negative results, no warning | S01 |
| [#153](https://github.com/mrchatam/Trace/issues/153) | `trace_search` / `trace_explore` / `trace_impact` walk emit no index-staleness signal (asymmetric with `trace_context`) | S01 |
| [#150](https://github.com/mrchatam/Trace/issues/150) | `trace_add` with an existing id silently overwrites the entity (data loss reported as `ok:true`) | S02 |
| [#151](https://github.com/mrchatam/Trace/issues/151) | `trace_add` accepts non-UUID id and unresolvable `goal_id` → dangling/corrupt graph rows | S02 |
| [#154](https://github.com/mrchatam/Trace/issues/154) | `trace_tasks` silently accepts invalid `work_state` filters and returns an authoritative empty list | S02 |
| [#162](https://github.com/mrchatam/Trace/issues/162) | `trace_link` accepts self-references — `rel=blocks` lets a task block itself, persisted with no warning | S02 |
| [#155](https://github.com/mrchatam/Trace/issues/155) | Raw `sql: no rows in result set` and wrong entity-type names leak from `trace_link` / `trace_why` / `trace_review` / `trace_changes` | S03 |
| [#156](https://github.com/mrchatam/Trace/issues/156) | `trace_loop action=next` returns an opaque domain error where `action=gate` returns a structured gate envelope | S03 |
| [#157](https://github.com/mrchatam/Trace/issues/157) | Inconsistent MCP response envelopes; `trace_changes list` lacks `count`/`truncated`, `trace_agents recommend` drops the empty-catalog hint | S04 |
| [#158](https://github.com/mrchatam/Trace/issues/158) | MCP tool schemas carry no enums/ranges; `trace_explore` skips the validation `trace_context` performs | S04 |
| [#159](https://github.com/mrchatam/Trace/issues/159) | MCP tool descriptions have drifted from the actual schemas (`trace_loop` omits `gate`; capability kind omits `AGENT`) | S04 |
| [#161](https://github.com/mrchatam/Trace/issues/161) | `trace_changes show` returns PascalCase Go field names while `trace_changes list` returns snake_case (third recurrence of this defect class) | S04 |
| [#152](https://github.com/mrchatam/Trace/issues/152) | `trace_version` requires an initialized store, defeating its documented stale-server-detection purpose | S05 |
| [#160](https://github.com/mrchatam/Trace/issues/160) | `trace_tasks` ordering is undocumented and degenerates to UUID order for same-second inserts | S05 |

## Scope sequence

```
S01 FTS honesty        (#149, #153)  opt-in matcher; staleness signal on all three walk tools
 → S02 Write-path      (#150, #151, #154, #162)  fail-closed validation on add/link/tasks
 → S03 Error quality   (#155, #156)  no driver text; structured gate envelope for loop.next
 → S04 Envelope/schema (#157–#159, #161)  uniform envelopes; generated descriptions + enums
 → S05 Store-free      (#152, #160)  trace_version answers from the process alone
 → S06 Guards          (permanent)  key-shape sweep + protocol conformance suite
```

| Scope | Theme | Artifact focus | First row |
|-------|-------|----------------|-----------|
| S01 | FTS honesty | `store.MatchAnd`/`MatchOps`, `matcher` input, `graph_sync_honesty` | P45-S01-00 |
| S02 | Write path | `trace_add` fail-closed, `refuseSelfLink`, `NormalizeWorkState` | P45-S02-00 |
| S03 | Error quality | `wrapToolErr`, `wrapLinkErr`, `trace.loop.gate.v1` for `loop.next` | P45-S03-00 |
| S04 | Envelope/schema | `listEnvelope`, `schemaFor[In]`, generated descriptions, `changeShowRow` | P45-S04-00 |
| S05 | Store-free reads | `store.ProjectInitialized` (stat-only), documented `trace_tasks` order | P45-S05-00 |
| S06 | Guards | lowercase-key sweep over 21 real tool calls; in-memory protocol conformance | P45-S06-00 |

## In scope

- Reportable narrowing (`graph_sync_honesty` on search/explore/impact; `query_note` on the FTS rewrite).
- Fail-closed writes and filters, with the offending parameter named in the message.
- One `listEnvelope` for every list-shaped tool response.
- One shared vocabulary per field, generating the schema enum **and** the description sentence, with a drift guard that fails the build when they diverge.
- Permanent regression guards that replace the ephemeral audit harness.

## Out of scope

- **No FK migration** for `tasks.goal_id` — the write path validates; the schema is unchanged.
- **No nanosecond `created_at`.** #160 is closed by *documenting* the ordering, not by changing the column.
- **No `replace` cascade.** `replace=true` overwrites the entity row only; it does not delete or merge links, changes or evidence attached to it.
- No new MCP tools, no hosted MCP, no `web/` changes, no changes under `cmd/trace`.
- Raising any cap (Laws 6–7 unchanged).

## Moat charter (M-001)

The task loop (Tasks → Loop → Gate → Review) is unchanged. S03 makes `loop.next` *honest about its own preconditions* — it returns the same gate envelope `loop.gate` returns, so an agent that already parses `trace.loop.gate.v1` needs no new parsing path. No scope in this phase adds a capability that competes with the loop; they make the existing loop's failure modes visible.

## Process note (read before trusting the board)

This phase was **implemented in one continuous session by the coding agent**, not by separate planner/implementer/reviewer sub-agent passes. The `02-review.md` file in each scope is therefore a **self-verification record** — the tests that were actually run and the behaviour actually observed — and is labelled as such. It is not an independent approval, and no scope claims one. Treat those rows as evidence, not as a second pair of eyes.

## Board

[`docs/TODO/phase-45.md`](../../TODO/phase-45.md) · Phase planner: [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md)