# Phase 45 / S01 / Planner — FTS honesty

## Metadata
- id: P45-S01-00
- todo_ids: [P45-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan the opt-in FTS matcher (#149) and the index-staleness signal on the three walk tools (#153). **No product code in this row.**

## References

- [`../00-PHASE-PLANNER.md`](../00-PHASE-PLANNER.md) — locked decision **D1**
- [`../SCOPE-TODOS.md`](../SCOPE-TODOS.md)
- Live targets:
  - `internal/store/fts.go` (`sanitizeFTSQuery:74`)
  - `internal/retrieval/{search,types}.go`
  - `internal/compiler/{explore,graph_sync_honesty}.go`
  - `internal/mcp/tools_{search,explore,impact}.go`
  - `internal/mcp/tools_context.go` — the existing honesty precedent

## Locked defaults

| Item | Value |
|------|-------|
| Mode name | `matcher` input; `and` (default) \| `ops` |
| Default behaviour | `SearchFTS` output **byte-identical** to pre-S01 |
| Operators | **Uppercase only**: `AND OR NOT NEAR NEAR(n)`; `"phrase"`; trailing `*` |
| Operator set | **Opt-in only** — never inferred from the query text |
| Limit cap | **64**, enforced in both modes |
| Fallback | Never silent — always a `query_note` / `fallback` string |
| Staleness | `graph_sync_honesty` on search, explore **and** impact walk |

## Why uppercase-only

FTS5 is uppercase-only, and so is every operator syntax an agent is likely to copy from a manual. Making lowercase `or` an operator would silently change results for anyone searching for the English word "or" — which is the exact defect this phase exists to remove. Uppercase-only means the only way to change behaviour is to pass `matcher=ops` **and** write an operator.

## Role work

1. Lock the tokenizer/compiler contract and its rejection cases in `01-implement.md`.
2. Define the acceptance tests for union-vs-intersection and for the rewrite notice.
3. Thicken `02-review.md` with the honesty invariants that must hold.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board notes; next **P45-S01-01**

## Next

`P45-S01-01`