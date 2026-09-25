# Phase 44 / S00 / Research

## Metadata
- id: P44-S00-01
- todo_ids: [P44-S00-01]
- role: implementer
- skills: [research, code-explorer]
- mcps: [user-trace, user-codegraph]
- verification: automated
- agents: []

## Objective

Author [`RESEARCH.md`](RESEARCH.md) in this folder: current link-creation paths, hairball root cause (file/line cites), and gap matrices vs [`INTAKE.md`](../../INTAKE.md) examples **and** D1–D5 light locks. Researcher must not re-debate Phase 44 locks — document gaps only. **No product code** (Notes-only for typos; no refactors).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 13, 19
- [`INTAKE.md`](../../INTAKE.md) — auth FE/BE + business examples
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — light-locked D1–D5 (do not reopen)
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- P44-S00-00 planner Notes on board (anchor re-verify)

## Session start

Follow agent-loop-protocol Session start. Unattended authority: INTAKE + P44-00 light locks. If HEAD contradicts a cited lock line, record the delta in RESEARCH.md §7 risks — do **not** edit `01-DESIGN-LOCKS.md` or P44-00 history.

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/RESEARCH.md` |
| Product code | **Forbidden** |
| Caps | **Record only** — GUI `PROJECT_MAX_NODES=500`, API/OpenAPI hard max **5000**; do not change code |
| D1–D5 | Document gaps vs locks; do not reverse light locks |
| Name collisions | Distinguish proposed `implements` / `blocks` (D2) from existing `change_implements_decision` / `uncertainty_blocks_task` |
| Theme | Meaningful scope clusters (auth FE/BE, business) — not UI-only |

## Preflight — live anchors (re-verified P44-S00-00 @ HEAD 2026-08-23)

Treat these as true unless researcher re-greps and finds drift (then note in §7):

| Claim | Evidence @ HEAD |
|-------|-----------------|
| Center = first goal | `project_graph.go:55–64` |
| `max_nodes` required; reject >5000; truncate | `project_graph.go` + `neighborhood.go:8` `MaxNeighborhoodNodes=5000` |
| Walk = `entity_links` + synthesized `goal_has_task` | `graph_neighbors.go:50–86` |
| `GraphNode` = id/kind/title/goal_id only (no `scope_id`) | `neighborhood.go:22–27` |
| `GraphEdge` = rel/from/to only (no provenance field) | `neighborhood.go:30–34` |
| Domain rels causal/planning; **no** `scope_member` / `api_contract` / bare `implements`/`blocks` | `service.go:33–64` |
| MCP `trace_link` = **5 aliases** only | `tools_write.go:183–195` (`goal-task`, `decision-task`, `discovery-plan-change`, `discovery-mentions-task`, `claim-evidence`) |
| `entity_links` has `source_type` + `confidence` | `store/links.go` |
| `plan_scopes` CRUD exists; **not** in retrieval collect | `store/plan_hierarchy.go`; grep `plan_scope` under `internal/retrieval` → empty |
| Seed export link allowlist narrow (4 causal rels) | `seed_export.go:341–346` (`seedExportLinkRels`) |
| GUI node budget 500 | `overviewCompose.ts` `PROJECT_MAX_NODES = 500` |
| GUI edge cap 150; priority leads with `goal_has_task` | `graphLayout.ts:37–48` (`EDGE_OVERVIEW_MAX`, `EDGE_PRIORITY_RELS`) |
| OpenAPI `/v1/graph` max 5000; `max_nodes` required | `api/openapi.yaml` |

**Note for gap matrix:** `EDGE_PRIORITY_RELS` already lists string tokens `blocks`, `depends_on`, `relates_to` but domain MVP scope rels are **not** persisted yet — treat as aspirational GUI tokens, not live scope semantics.

## File targets (read-only)

| Path | Question to answer in RESEARCH.md |
|------|-----------------------------------|
| `internal/retrieval/project_graph.go` | How are nodes collected/sorted? How is `center` chosen? Truncation honesty? |
| `internal/retrieval/graph_neighbors.go` | Which edges enter the walk? How is `goal_has_task` synthesized? |
| `internal/retrieval/neighborhood.go` | Cap constant; `GraphNode`/`GraphEdge` fields vs D3 (`scope_id`, provenance) |
| `internal/domain/service.go` | Full rel const list vs D2 MVP set |
| `internal/domain/doc.go` | Documented required/optional rels; goal→task via `goal_id` |
| `internal/mcp/tools_write.go` (`toolLink`) | Exact allowlist; what agents cannot create today |
| `internal/store/links.go` | Persistence shape ready for `INFERRED`? |
| `internal/store/plan_hierarchy.go` | What is a `PlanScope`? Why invisible to project graph? |
| `internal/domain/seed_export.go` | Which links export; plan tree vs graph edges |
| `web/src/lib/graphLayout.ts` | Edge priority / overview cap vs hairball |
| `web/src/lib/overviewCompose.ts` | Project graph fetch budget |
| `api/openapi.yaml` `/v1/graph` | Contract vs D3 `scope` filter (absent today) |

Optional (if clarifying create paths): CLI link wrappers under `cmd/` / `internal/cli` — cite only if found; do not invent.

## RESEARCH.md required sections (runnable template)

Author **exactly** these H2 sections (extra H3 OK):

### 1. Executive summary
2–4 sentences: hairball is goal-centric topology (honest for loop); missing scope cartography; Phase 44 fills via D1–D5 without raising GUI 500.

### 2. Current link creation paths
Table columns: **Path** | **Rel / mechanism** | **Store** | **Visible in `/v1/graph`?**

Must cover at least:
- MCP `trace_link` (5 aliases)
- Domain side-effects (reviews → `review_judges_task`, etc.)
- `tasks.goal_id` → synthesized `goal_has_task`
- Seed import / `seedExportLinkRels`
- `plan_scopes` (planner hierarchy — **not** graph walk)

### 3. Hairball root cause
Cite file:line for: center selection, `goal_has_task` fan-out, sparse cross-links, GUI `EDGE_OVERVIEW_MAX` + priority favoring `goal_has_task`. Conclude: topology ≠ feature cartography.

### 4. Gap matrix vs INTAKE examples
One row per INTAKE example. Columns: **Scope example** | **Desired intra** | **Desired cross** | **Today** | **Gap (S02/S03/S04)**

| Scope example | Must address |
|---------------|--------------|
| Frontend auth | login form ↔ session store ↔ route guard; cross via **`api_contract`** to BE |
| Backend auth | handler ↔ middleware ↔ token service; cross via contract to FE |
| Business | marketing landing ↔ design system ↔ human operator task; hierarchy without flattening to one goal |

Each row must name which D2 rels (or D1 membership) close the gap.

### 5. Gap matrix vs D1–D5 light locks
Table: **Lock** | **Locked intent** | **HEAD state** | **Build in** (S02/S03/S04)

Rows: D1 (thin scope + `scope_member`), D2 (MVP rels), D3 (`scope` query + N-hop), D4 (CLI inference only), D5 (GUI 500 + filter first).

### 6. Non-goals confirmed
Bullet list affirming: DR-NOSSEM (no embeddings), no silent index-hook/daemon (D4), no raising GUI 500 without VERIFY (D5), no parallel relationship store (Law 13), no GUI business logic (Law 19), M-001 moat (scope enriches; does not replace task loop).

### 7. Open risks for S01 (wording only)
≤5 bullets. Examples OK: `implements` vs `change_implements_decision` naming; whether `plan_scopes` map 1:1 to D1 scope records; seed-export allowlist expansion. **Do not** propose reversing D1–D5.

## Acceptance checks (self-check before board `done`)

- [ ] All 7 sections present with non-empty substance
- [ ] INTAKE three example rows covered (§4)
- [ ] D1–D5 gap rows covered (§5)
- [ ] Hairball cites ≥3 live file:line anchors
- [ ] Caps stated as 500 GUI / 5000 API (unchanged)
- [ ] No product code diffs in the commit/working tree for this row
- [ ] `SCOPE-TODOS.md` S00-1 → done when artifact exists

## Exit criteria

- [ ] `RESEARCH.md` present with all required sections
- [ ] Gap matrix covers INTAKE three example rows
- [ ] Board row `done` with path to `RESEARCH.md`; next **P44-S00-02**

## Minimal todos

- [ ] Re-skim live anchors if HEAD moved since P44-S00-00
- [ ] Author RESEARCH.md per template
- [ ] Self-check acceptance list
- [ ] Update board Notes only (+ SCOPE-TODOS S00-1 / S00-G)

## Todo updates

Implementer: **status + notes only** on `P44-S00-01`. Do not edit P44-00 or thicken later scopes.

## Next

`P44-S00-02`
