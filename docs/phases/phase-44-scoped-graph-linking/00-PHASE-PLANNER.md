# Phase 44 — Scoped graph linking & scope semantics

**Phase planner.** Row `P44-00`. **Done** 2026-08-23 (human-promoted; planner executed).

## Metadata
- id: P44-00
- todo_ids: [P44-00]
- role: planner
- skills: [planning-and-task-breakdown, domain-modeling, incremental-implementation]
- mcps: [trace_context, trace_why]
- verification: automated

## Mission

Lock Phase 44 scopes against [`INTAKE.md`](INTAKE.md) and [`01-DESIGN-LOCKS.md`](01-DESIGN-LOCKS.md). Thicken scope planners; confirm board row order. **No product code.**

## Gate

Human promotion satisfied 2026-08-23. Phase 43 complete (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | Intake + research | P44-S00-00 → 02 | Current link model + gap matrix |
| S01 | Design locks | P44-S01-00 → 02 | Finalize 01-DESIGN-LOCKS (D1–D5 → LOCKED) |
| S02 | Backend links | P44-S02-00 → 02 | entity_links rels, domain, MCP, OpenAPI, graph walk |
| S03 | Inference | P44-S03-00 → 02 | Opt-in CLI derived links + provenance |
| S04 | GUI | P44-S04-00 → 02 | Scope-aware layout + orient (Law 19) |
| S05 | VERIFY | P44-S05-00 → 02 | VERIFY-NOTES + DR-HANDOFF |

## Locked defaults (FINAL-for-S01 — light-locked P44-00)

| Item | Value |
|------|-------|
| Entry | User intake 2026-08-23 — scoped interconnection, not UI-only |
| M-001 moat | Scope links enrich graph; never replace task loop |
| Law 6–7 | Bounded graph API; **500 GUI default**; API max **5000** |
| Law 13 | Reuse `entity_links`; no parallel relationship store |
| Law 19 | GUI/HTTP adapters thin; walk logic in `internal/retrieval` |
| DR-NOSSEM | No vector/semantic scope clustering in P44 |
| Phase 43 | Do not reopen GitHub hygiene rows |
| D1 | A + thin scope records (`scope` + `scope_member`); reject tag-only primary |
| D2 | MVP rels: `scope_member`, `api_contract`, `implements`, `blocks` |
| D3 | `scope` query param + bounded N-hop; no tiles/pagination |
| D4 | Inference: **CLI only** (no index hook) |
| D5 | GUI cap stays **500**; scope filter before raising cap |

## Live-repo research (P44-00 verified)

| File | Finding |
|------|---------|
| `internal/retrieval/project_graph.go:55–64` | `center` = first `kind==goal` in budgeted node list |
| `internal/retrieval/project_graph.go:22–74` | `max_nodes` required; hard reject >5000; `Truncated` when total exceeds budget |
| `internal/retrieval/graph_neighbors.go:8–88` | Walk = `entity_links` both dirs + synthesized `goal_has_task` from `tasks.goal_id` |
| `internal/retrieval/neighborhood.go:8–34` | `MaxNeighborhoodNodes=5000`; `GraphNode` has id/kind/title/goal_id only (no `scope_id` yet) |
| `internal/domain/service.go:33–64` | Causal/planning rels present; **no** scope taxonomy |
| `internal/domain/doc.go:14–22` | Documented required/optional rels; goal→task via `goal_id` only |
| `internal/mcp/tools_write.go:183–195` | `trace_link` allowlist: 5 aliases only |
| `internal/store/links.go` | `entity_links` + `source_type`/`confidence` — ready for INFERRED |
| `internal/store/plan_hierarchy.go` | `plan_scopes` CRUD/list; **not** in project graph collect |
| `web/src/lib/overviewCompose.ts` | `PROJECT_MAX_NODES = 500` |
| `web/src/lib/graphLayout.ts:36–48` | `EDGE_OVERVIEW_MAX=150`; priority set leads with `goal_has_task` |
| `api/openapi.yaml` `/v1/graph` | `max_nodes` required, maximum 5000 |

## Planner gate (P44-00)

- [x] Re-read live repo files in Research anchors
- [x] Thicken S00–S05 `00-PLANNER` + `01-*` + `02-*` prompts
- [x] Board `docs/TODO/phase-44.md` — scope rows filled / notes thickened
- [x] Close open decisions D1–D5 in 01-DESIGN-LOCKS or defer with reason → **light-locked**; S01 formal APPROVE
- [x] No product code in this row
- [x] `SCOPE-TODOS.md` per scope folder

## Exit criteria

- [x] Scope stubs runnable; board points to **P44-S00-00**
- [x] 01-DESIGN-LOCKS ready for S01 lock pass (status: Light-locked pending S01 APPROVE)
- [x] Successor after P44: **TBD at VERIFY** (default `no successor`)

## Next

`P44-S00-00`
