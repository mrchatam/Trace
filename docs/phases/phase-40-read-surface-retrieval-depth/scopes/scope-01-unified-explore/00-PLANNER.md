# P40-S01-00 — Scope planner (G2 unified explore)

## Metadata
- id: P40-S01-00
- todo_ids: [P40-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, context-engineering]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S01 **G2** against live repo: unified task-aware capped `trace_explore` MCP tool (G-007). G1 prerequisite shipped; law spike gate resolved. Thicken `01-implement.md` + `02-review.md`. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P40-00 Q2 resolution
- [REMEDIATION-PLAN §2 G2](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-007](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Phase 39 G1 + G3 shipped:
  - `internal/compiler/compiler.go:158–165` — `ContextOptions.Query` merge
  - `internal/mcp/instructions.go:13–16` — compose-first recipe (ranked read tools)
  - `internal/mcp/server.go:229–237` — **16** tools today; no `trace_explore`
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P40-00):
  - `internal/mcp/tools_search.go` — `trace_search` FTS
  - `internal/mcp/tools_why.go` — `trace_why` causal neighborhood
  - `internal/mcp/tools_context.go` — `trace_context` task packet (+ optional query)
  - `internal/mcp/tools_impact.go` — `trace_impact` (write/read mix)
  - `internal/retrieval/neighborhood.go:53+` — bounded graph
  - `internal/httpapi/handlers_retrieval.go` — context/search/graph HTTP mirrors

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P40-00 locks are authority.

## Locked defaults (FINAL — P40-00)

| Item | Value |
|------|-------|
| GAP ids | G-007 |
| Verdict | **Accept** unified `trace_explore` per REMEDIATION-PLAN G2 secondary path |
| Prerequisite G1 | **Shipped** Phase 39 S00 — query+task merge live |
| Law spike gate | **Waived** — desk-check checklist below; no mandatory live MCP spike before S01-01 (mirrors P39 INTAKE Q2) |
| M-001 | `task_id` **required**; explore merges into task loop — **never query-only** |
| Tool policy | **Add `trace_explore` as 17th MCP tool** — read-only; write tools stay registered |
| Library first | New compose function in `internal/retrieval/` or `internal/compiler/` — MCP/CLI thin adapters |
| Peer pattern | CG `codegraph_explore` single capped call (**contrast** — task-aware + moat, not CG facade) |
| Caps | Progressive defaults aligned with existing search/context limits; explicit `truncated` honesty |
| Instructions | Update `ServerInstructions()` — explore is **convenience after moat path**, not replacement |
| Out | 1-tool CG facade; mega-tool hiding writes; query-only explore; G-004a vector |
| GUI | Out of scope — S00 owns GUI orient |

## Law spike desk-check (S01-00 gate — document, do not block implement)

| Check | Expected |
|-------|----------|
| Task required | Empty `task_id` → validation error |
| Query optional | Merges via G1 path when set |
| Caps honest | Response includes truncation/budget fields |
| Write surface visible | 16 existing tools + explore = 17; loop/review/transition unchanged |
| Not CG equivalent claim | Response shape may differ; task packet section required |
| Compose-first preserved | Instructions still rank manual compose for fine control |

## Accept / reject (G2)

| Decision | Item |
|----------|------|
| **Accept** | Library `Explore` (or equivalent) composing search + task context + neighborhood + why summaries |
| **Accept** | MCP `trace_explore` with `task_id` required, optional `query`, capped opts |
| **Accept** | `RegisteredToolNames()` → **17** tools; update tool-count tests |
| **Accept** | `ServerInstructions()` addendum — explore optional; moat-first lead unchanged |
| **Accept** | CLI mirror optional (`trace explore`) — implement if low cost; not blocking |
| **Reject** | Query-only explore (no task_id) |
| **Reject** | Hiding/reducing write tools to “fix” discovery |
| **Reject** | Full-graph dump default |
| **Reject** | Claim compose ≈ CG explore (h7-compose-desk-check) |

## Must lock for S01-01 (delivered in thickened 01-implement)

1. Touch-list: library explore → MCP tool → instructions → tests.
2. Seven acceptance tests G2-T1–T7 (see `01-implement.md`).
3. Tool count migration 16 → 17 documented in review checklist.

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable with file targets + G2 accept map
- [ ] SCOPE-TODOS updated
- [ ] Board row → `done` with Notes

## Next

`P40-S01-01`
