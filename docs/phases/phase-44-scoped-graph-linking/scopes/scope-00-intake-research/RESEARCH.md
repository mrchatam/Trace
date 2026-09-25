# Phase 44 / S00 — RESEARCH

**Authored:** 2026-08-23 · **Row:** P44-S00-01 · **HEAD anchors:** re-verified against live sources (no drift vs P44-S00-00 table).

**Caps (record only — unchanged):** GUI `PROJECT_MAX_NODES = 500`; API/OpenAPI hard max `5000` (`MaxNeighborhoodNodes`).

---

## 1. Executive summary

The Explore project graph is an honest **goal-centric causal star**: `center` is the first goal, and nearly every task fans out via synthesized `goal_has_task` from `tasks.goal_id`. That topology serves the task loop (Laws / M-001) but does **not** show feature or business cartography—auth FE/BE peers, API contracts, or marketing↔design clusters. Cross-links are sparse because MCP/`trace link` expose only five aliases and domain rels are causal/planning, not scope. Phase 44 closes that gap via light-locked D1–D5 (thin scopes + MVP rels + `scope` filter + CLI inference + keep GUI 500) without raising the browser node budget.

---

## 2. Current link creation paths

| Path | Rel / mechanism | Store | Visible in `/v1/graph`? |
|------|-----------------|-------|-------------------------|
| MCP `trace_link` | Aliases only: `goal-task`, `decision-task`, `discovery-plan-change`, `discovery-mentions-task`, `claim-evidence` (`tools_write.go:183–195`) | `goal-task` → `tasks.goal_id` + event; others → `entity_links` | Yes, if both endpoints in budgeted node set (walk via `graphWalkNeighbors`) |
| CLI `trace link` | Same five aliases (`cmd/trace/link.go:51–64`) | Same as MCP | Same as MCP |
| Domain side-effects | e.g. `review_judges_task` / `review_judges_scope` (`review.go`); `change_implements_decision` (`changes.go`); `uncertainty_blocks_task` (`cognitive.go`); regression/hypothesis/assumption/effect rels (`service.go:33–64`) | `entity_links` (+ events) | Yes when both ends collected as project nodes; **`plan_scope` targets are not project nodes** (see below) |
| `tasks.goal_id` | Synthesized walk edge `goal_has_task` (`graph_neighbors.go:50–86`); not a row in `entity_links` (`doc.go`; `RelGoalHasTaskEvent` is event payload only) | Task column | **Yes** — primary fan-out for every goal-linked task |
| Seed import | Imports allowlisted link rows + plan tree | `entity_links` + `plan_*` tables | Links: yes if in walk; plan tree: **not** as graph edges |
| Seed export `seedExportLinkRels` | Only 4 causal rels: `decision_affects_task`, `discovery_causes_plan_change`, `claim_has_evidence`, `discovery_mentions_task` (`seed_export.go:341–346`) | N/A (read path) | Export omits many live rels (e.g. reviews); plan scopes export as **tree JSON**, not graph edges |
| `plan_scopes` CRUD | Progressive planner hierarchy: phase → scope → deep plan (`plan_hierarchy.go`; `PlanScope` = id/phase_id/title/body/ord/status/…) | `plan_scopes` table | **No** — `collectProjectNodes` never lists them; `grep plan_scope` under `internal/retrieval` → empty. Invisible to project-graph walk |

**Agents cannot create today (via MCP/CLI allowlist):** any D2 MVP scope rel (`scope_member`, `api_contract`, `implements`, `blocks`), bare `depends_on` / `relates_to` / `same_feature_front_back`, or free-form `entity_links` writes outside the five aliases / domain APIs.

---

## 3. Hairball root cause

1. **Center = first goal** — After kind-sorted collect + truncate, `ProjectGraph` sets `center` to the first node with `Kind == "goal"` (`project_graph.go:55–64`). Layout/orient treat that as the hub.
2. **`goal_has_task` fan-out** — For every task with `goal_id`, and for every goal listing its tasks, the walk emits `goal_has_task` (`graph_neighbors.go:50–86`). One goal with many tasks yields a dense star independent of feature boundaries.
3. **Sparse real cross-links** — Walk edges are `entity_links` ∪ synthesized goal edges only (`graph_neighbors.go:9–10, 14–48`). Domain vocabulary is causal/planning (`service.go:33–64`); **no** `scope_member` / `api_contract` / bare `implements` / `blocks`. MCP cannot mint scope cartography (`tools_write.go:183–195`).
4. **GUI reinforces the star** — Overview caps rendered edges at `EDGE_OVERVIEW_MAX = 150` and **prioritizes** `goal_has_task` (and other tokens) when sampling (`graphLayout.ts:37–48`). Tokens `blocks` / `depends_on` / `relates_to` appear in `EDGE_PRIORITY_RELS` but are **aspirational** until D2 rels persist—they do not create clusters today.
5. **No scope filter / membership on nodes** — `GraphNode` is id/kind/title/goal_id only (`neighborhood.go:22–27`); `GraphEdge` is rel/from/to only (`neighborhood.go:30–34`). OpenAPI `/v1/graph` has `mode` / `center` / `max_nodes` / `depth` — **no** `scope` query (`api/openapi.yaml` ~629–671). Caps stay GUI **500** / API **5000**.

**Conclusion:** Topology ≠ feature cartography. The graph correctly expresses the goal/task moat; Phase 44 must add scope membership and typed cross-edges, then filter/layout—not raise GUI 500.

---

## 4. Gap matrix vs INTAKE examples

| Scope example | Desired intra | Desired cross | Today | Gap (S02/S03/S04) |
|---------------|---------------|---------------|-------|-------------------|
| **Frontend auth** | login form ↔ session store ↔ route guard | ↔ backend auth via **`api_contract`** | Tasks (if any) hang under first goal via `goal_has_task`; no `scope_member` cluster; no FE↔BE contract edge; agents cannot `trace_link` scope rels | **S02:** thin scope record + `scope_member` (D1) for FE auth entities; MVP `implements`/`blocks` among tasks; `api_contract` FE→BE. **S03:** path/title rules (`web/` auth) → `INFERRED` membership. **S04:** cluster layout + elevate MVP rels in `EDGE_PRIORITY_RELS` |
| **Backend auth** | handler ↔ middleware ↔ token service | ↔ frontend via OpenAPI/route **`api_contract`** | Same star; Go handlers are not scope-linked; `change_implements_decision` is **not** the D2 task↔task `implements` | **S02:** BE auth `scope_member` + intra `implements`/`blocks`; reciprocal/paired `api_contract` endpoints. **S03:** `internal/` path globs / OpenAPI tag inference. **S04:** second cluster + cross-scope edge styling |
| **Business** | marketing landing ↔ design system ↔ human operator task | ↔ goal/task hierarchy **without** flattening to one goal | Plan/body text may mention areas, but `plan_scopes` never enter the graph; business peers only meet at the goal hub | **S02:** `scope` kind=`business` + `scope_member`; keep `goal_has_task` for loop honesty while scope edges carry cartography. **S03:** map/align titles to `plan_scopes` (no hard FK required). **S04:** ≥2 clusters without collapsing all nodes onto first goal |

Each INTAKE row is closed by **D1 membership** (`scope_member`) plus **D2** MVP rels (`api_contract` for FE↔BE; `implements`/`blocks` for intra peers). Optional D2 (`depends_on`, `relates_to`, `same_feature_front_back`) may enrich but must not block MVP.

---

## 5. Gap matrix vs D1–D5 light locks

| Lock | Locked intent | HEAD state | Build in |
|------|---------------|------------|----------|
| **D1** | Thin `scope` records (id/slug/title/kind) + `scope_member` in `entity_links`; reject tag-column primary; reuse/map `plan_scopes` where titles align | No scope table for feature/business cartography; `plan_scopes` exist for planner only and are **not** in retrieval; no `scope_member` rel | **S02** (records + membership edges + export); **S03** (title/plan_scope alignment inference) |
| **D2** | MVP rels: `scope_member`, `api_contract`, `implements`, `blocks` (optional: `same_feature_front_back`, `depends_on`, `relates_to`) | Domain consts are causal/planning only (`service.go:33–64`); MCP/CLI 5 aliases; name collisions with existing `change_implements_decision` / `uncertainty_blocks_task` (distinct strings — must stay distinct) | **S02** (domain + MCP/CLI allowlist + seed allowlist); **S04** (priority/legend for new rels) |
| **D3** | `GET /v1/graph?mode=project&scope=…&max_nodes=…` + bounded N-hop; `scope_id` on nodes; optional edge provenance; no tiles | No `scope` param; `GraphNode`/`GraphEdge` lack `scope_id` / provenance (`neighborhood.go:22–34`; OpenAPI GraphNode/GraphEdge); `max_nodes` required, reject >5000 — OK | **S02** (API/retrieval/OpenAPI); GUI consumes in **S04** |
| **D4** | Inference **CLI-only**; no index-hook / silent daemon; `source_type=INFERRED`; explicit wins | `entity_links` already has `source_type` + `confidence` (`store/links.go`) — shape ready; **no** `trace graph infer` (or equiv.); no silent path today | **S03** |
| **D5** | Keep GUI **500**; prefer scope filter over raising cap; edge overview 150 stays presentation policy | `PROJECT_MAX_NODES = 500` (`overviewCompose.ts`); API 5000 unchanged | **S04** (filter + cluster; **do not** raise 500 without VERIFY) |

---

## 6. Non-goals confirmed

- **DR-NOSSEM** — no embedding-based auto-clustering / semantic scope fill.
- **D4** — no silent post-`trace index` hook or always-on inference daemon in Phase 44.
- **D5** — no raising GUI default above **500** without VERIFY evidence that scope filter still fails.
- **Law 13** — no parallel relationship store; scope membership and MVP rels live in `entity_links` (+ thin scope records as entities/rows, not a second edge table).
- **Law 19** — GUI/HTTP remain adapters; walk/filter logic stays in `internal/retrieval` (no business-logic fork in `web/`).
- **M-001** — scope cartography **enriches** orientation; it does **not** replace the goal/task loop, board order, or gate rels (`review_judges_task`, etc.).

---

## 7. Open risks for S01 (wording only)

1. **`implements` / `blocks` naming** — Proposed D2 strings must remain distinct from existing `change_implements_decision` and `uncertainty_blocks_task` in docs, MCP aliases, and seed allowlists so agents do not confuse change↔decision or uncertainty gates with task↔task scope edges.
2. **`plan_scopes` vs D1 scope records** — Planner scopes are phase-ordered progressive-plan rows, not feature/layer/business cartography. S01 should word “reuse/map where titles align” without implying 1:1 FK or that every `plan_scope` becomes a graph scope node in S02 MVP.
3. **Seed-export allowlist expansion** — Today only four causal rels export (`seed_export.go:341–346`). D2 + thin scopes need explicit allowlist + portable `trace/graph.json` policy before PRs that change entities (CONTRIBUTING); S01 should state export expectations without reopening D1–D5.
4. **Aspirational GUI tokens** — `EDGE_PRIORITY_RELS` already lists `blocks` / `depends_on` / `relates_to` before persistence; S01/S04 wording should treat them as presentation prep, not evidence that scope semantics exist at HEAD.
5. **Provenance on wire** — Store can hold `INFERRED` via `source_type`, but `GraphEdge` has no provenance field yet; S01 should confirm optional `edge.provenance` mapping (D3) without requiring UI suffix mutation of stored rel names.

**Do not reverse D1–D5.** Wording refinements and collision/export clarifications only; reversals need a REVIEW spawn.
