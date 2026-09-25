# Phase 44 / S02 / Planner — Backend links

## Metadata
- id: P44-S02-00
- todo_ids: [P44-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, incremental-implementation]
- mcps: [user-trace]
- verification: automated

## Objective

Plan migration + domain + retrieval + MCP + OpenAPI for scope link rels per **LOCKED** [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md). Thicken `01-implement.md` to run alone. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — D1–D3, §3, §8
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- Live targets:
  - `internal/domain/service.go` (rel consts)
  - `internal/domain/doc.go`
  - `internal/store/links.go` + migrations
  - `internal/retrieval/project_graph.go`, `neighborhood.go`, `graph_neighbors.go`
  - `internal/mcp/tools_write.go`
  - `internal/domain/seed_export.go` / seed import
  - `api/openapi.yaml`

## Session start

Follow agent-loop-protocol. Block if S01 locks not APPROVED.

## Locked defaults

| Item | Value |
|------|-------|
| Store | `entity_links` only (Law 13) |
| MVP rels | `scope_member`, `api_contract`, `implements`, `blocks` |
| Scope records | Thin table or entity per D1 LOCKED |
| Graph API | Optional `scope_id` on nodes; `scope` query param; provenance on edges |
| Caps | `max_nodes` required; hard max 5000 |
| MCP | Extend `trace_link` allowlist — no new write tool |
| Portable | Migration + seed export allowlist + `trace/graph.json` before PR |

## Role work

1. Lock file touch-list + migration strategy in `01-implement.md`.
2. Define acceptance tests (unit: insert MVP rels; ProjectGraph includes scope edges; reject unbounded).
3. Thicken `02-review.md` export round-trip checklist.

## Exit criteria

- [x] `01-implement.md` runnable alone
- [x] Board Notes; next **P44-S02-01**

## Minimal todos

- [x] Thicken 01 + 02 from LOCKED design
- [x] Board done

## Next

`P44-S02-01`
