# Phase 44 / S04 / Implement GUI

## Metadata
- id: P44-S04-01
- todo_ids: [P44-S04-01]
- role: implementer
- skills: [frontend-ui-engineering, verifying-in-browser, incremental-implementation]
- mcps: [user-trace]
- verification: mixed

## Objective

Ship **scope-aware** Explore layout: cluster nodes by API `scope_id`, elevate MVP scope edges in overview priority, style **inferred** edges dashed, keep **`PROJECT_MAX_NODES = 500`**, and document honesty in the orient legend. **Law 19:** `web/` is adapter-only — no second graph walk, no inventing membership from edge crawling when `scope_id` is absent. **Prerequisite:** S02 APPROVE (P44-S02-02) + S03 APPROVE via P44-S03-02b (provenance wire + seed honesty for dashed demo).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §5–§6, **D5**, §3.4 provenance
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) §3 hairball, §7 risk 4
- Phase 40 G5 orient baseline (`GraphOrientPanel.tsx`)
- Live: `web/src/lib/graphLayout.ts`, `web/src/lib/overviewCompose.ts`, `web/src/screens/Graph.tsx`, `web/src/api/ops.ts`, `api/openapi.yaml` (`GraphNode.scope_id`, `GraphEdge.provenance`, `scope` query)
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Block if P44-S02-02 or P44-S03-02b is not APPROVE.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Cap | Keep `PROJECT_MAX_NODES = 500` in `overviewCompose.ts` — **do not raise**; `UI_CAP` stays tied to it |
| Edge overview | `EDGE_OVERVIEW_MAX = 150` stays presentation policy |
| Membership SoT | Cluster key = **`node.scope_id` from API** only. Missing/`undefined` → `__ungrouped__`. **Do not** derive clusters by walking `scope_member` edges in the browser |
| Rel styling | Stored `rel` strings unchanged (no `inferred_` prefix). Dashed style from `edge.provenance === 'inferred'` |
| Priority | Elevate MVP scope rels in `EDGE_PRIORITY_RELS` / tiered cap; reduce `goal_has_task` dominance **when any MVP scope edge is present** in the payload |
| Scope filter | Optional: pass `scope` query to `getProjectGraph` when UI exposes a scope picker — still bounded by `max_nodes=500`. Filter logic stays server-side |
| Types | OpenAPI already has `scope_id` / `provenance`; sync `web/src/api/schema.d.ts` (or thin local widen) so TS consumes them |
| Law 19 | No retrieval/filter rewrite in `web/`; no duplicate SQLite; HTTP client only |
| Out | Raising cap; tiles/pagination; embedding clustering (DR-NOSSEM); forking walk in browser; mutating stored rels |

### Preflight anchors @ HEAD (re-verify if drift)

| Claim | Location |
|-------|----------|
| `PROJECT_MAX_NODES = 500` | `web/src/lib/overviewCompose.ts:8` |
| `EDGE_OVERVIEW_MAX = 150` + aspirational `EDGE_PRIORITY_RELS` (includes `blocks`, **not** `scope_member`/`api_contract`/`implements`) | `web/src/lib/graphLayout.ts:37–48` |
| Force layout soft kind bias only (no scope centroid) | `computeForceLayout` ~301–348 |
| Project fetch ignores `scope` query | `getProjectGraph` in `ops.ts:127–132` |
| Edge style solid only (no provenance dash) | `Graph.tsx` `buildFlowGraph` ~262–281 |
| Orient panel has no scope/inference legend | `GraphOrientPanel.tsx` |
| Wire fields exist server-side | `internal/retrieval/neighborhood.go` `GraphNode.ScopeID`, `GraphEdge.Provenance`; OpenAPI ~1518–1538 |
| Generated TS may lag OpenAPI | `web/src/api/schema.d.ts` GraphNode/GraphEdge **lack** `scope_id`/`provenance` @ planner time — fix in this scope |

---

## Cluster-by-scope layout (spec)

### Inputs (API fields only)

| Field | Use |
|-------|-----|
| `GraphNode.scope_id` | Cluster id; omit → `__ungrouped__` |
| `GraphNode.kind` / `id` | Existing force soft bias + collide (keep) |
| `GraphEdge.{from,to,rel,provenance}` | Link forces + edge style/priority |

### Algorithm (locked shape — implementer may tune constants, not shape)

1. **Partition** nodes by `scope_id` (or `__ungrouped__`).
2. **Scope centroids:** place each non-empty scope cluster center on a ring (or grid) around canvas center. Deterministic order: sort scopes by `scope_id` string (stable). `__ungrouped__` sits at canvas center with weaker pull.
3. **Force:** extend `computeForceLayout` (or sibling `computeScopeForceLayout`) so each sim node gets a **scope attractor** (`forceX`/`forceY` or custom force toward its cluster centroid). Strength: intra-scope cohesion **stronger** than kind-lane bias; keep collide + charge so clusters remain readable at ≤500 nodes.
4. **Link distance:** optional — slightly shorter distance for edges whose both ends share `scope_id`; slightly longer / weaker for cross-scope (`api_contract`, etc.). Do **not** invent edges.
5. **Neighborhood mode** (`layoutMode === 'neighborhood'`): keep radial `computeOverviewPositions`; scope clustering is **project overview** primary. Optional light tint by `scope_id` on nodes is OK if cheap.
6. **Honesty:** if **zero** nodes have `scope_id`, layout may fall back to current kind-biased force (no fake clusters). Smoke fixture **must** populate ≥2 scopes so acceptance is testable.

### Visual cluster cues (adapter)

- Optional: soft fill / ring / label chip per distinct `scope_id` (CSS vars; not a second layout engine).
- Do not require scope **entity nodes** in the RF graph for MVP if API does not emit `kind=scope` nodes — clustering members is enough (D1 membership via `scope_id`).

---

## EDGE_PRIORITY_RELS (spec)

Current set elevates `goal_has_task` first → hairball-friendly overview. Change policy:

### Tier A — MVP scope (always priority when present)

`scope_member`, `api_contract`, `implements`, `blocks`

### Tier B — causal / planning (fill remaining budget)

`decision_affects_task`, `depends_on`, `relates_to`, `mentions`, `supports`, `contradicts`, … (existing non-MVP tokens OK)

### Tier C — demote when scope signal exists

`goal_has_task`: **remove from priority** (or move to last-fill only) **when the edge list contains ≥1 Tier A rel**. If no Tier A edges in payload, keep today’s behavior (goal-centric overview still useful).

### Cap algorithm

Replace flat `Set` filter with **tiered** `capEdgesByPriority`:

1. Take Tier A up to `EDGE_OVERVIEW_MAX`.
2. If room, add Tier B (and Tier C only if no Tier A in payload, or as last fill).
3. Preserve focus-incident override in `filterEdgesForOverview` (unchanged).

Unit-test: with mixed `goal_has_task` + `api_contract` + filler, capped overview **retains** `api_contract` and is **not** flooded by `goal_has_task`.

---

## Inferred edge style (spec)

| `edge.provenance` | Style |
|-------------------|--------|
| `inferred` | Dashed stroke (`strokeDasharray`), muted opacity or distinct CSS class `graph-edge--inferred` |
| `explicit` / omitted | Solid (current default). Treat empty as explicit for styling |

**Do not** change `rel` label text for inferred. Optional subtle label suffix in legend only, not on every edge.

Causal vs scope: optional width/class for Tier A rels (`graph-edge--scope`) vs causal — nice-to-have if cheap; dashed+priority is mandatory.

---

## Orient / legend (spec)

Update `GraphOrientPanel` (and any Graph chrome legend) to state:

1. **Scopes:** nodes cluster by server `scope_id` (feature/layer/business cartography).
2. **Inferred:** dashed edges = `provenance=inferred` (CLI `trace graph infer`); explicit wins / solid.
3. **Budget:** still `mode=project`, cap **500**, truncation banners honest (Laws 6–7).
4. Keep moat-first Tasks → Loop → Gate → Review copy (M-001).

a11y: legend must not rely on color alone (dash pattern + text). Keep `role="region"` / labelled heading.

---

## Types + fetch wiring

| Path | Do |
|------|----|
| `web/src/api/schema.d.ts` and/or ops types | Add optional `scope_id` on nodes; optional `provenance: 'explicit' \| 'inferred'` on edges |
| `GraphNodeMeta` (`overviewCompose.ts` / Graph) | Carry `scope_id` through compose/merge |
| `getProjectGraph(maxNodes, opt)` | Optional `scope?: string` → query `scope` (server filter). Default Explore path may omit (full project ≤500) |
| `Graph.tsx` | Pass `scope_id` into layout; apply edge dash/class from `provenance`; keep truncation banner |

---

## Acceptance fixture (≥2 clusters + cross-scope edge)

Seed (via existing CLI/MCP/`trace_link` / seed import — **do not** invent client membership):

| Artifact | Intent |
|----------|--------|
| Scope `auth-fe` (`kind=feature` or `layer`) | Cluster A |
| Scope `auth-be` | Cluster B |
| ≥2 tasks (or entities) with `scope_member` → each scope so API sets `scope_id` | Visible members in two clusters |
| ≥1 `api_contract` edge FE task → BE task | Cross-scope edge in overview priority |
| Optional: one `INFERRED` `scope_member` (post-`trace graph infer` or seed `source_type`) | Dashed edge visible |

Browser smoke Notes must record: URL/`trace serve`, seed steps, screenshot or DOM evidence of **≥2 clusters** + **≥1 cross-scope** edge.

---

## File touch-list (ordered)

| Slice | Path | Do |
|-------|------|----|
| A | `web/src/api/schema.d.ts`, `ops.ts` | Types + optional `scope` on `getProjectGraph` |
| B | `overviewCompose.ts` (+ tests if merge drops fields) | Preserve `scope_id` on `GraphNodeMeta`; **do not** change `PROJECT_MAX_NODES` |
| C | `graphLayout.ts` + `graphLayout.test.ts` | Scope force layout; tiered `EDGE_PRIORITY_*`; export helpers for tests |
| D | `Graph.tsx` (+ CSS if needed) | Wire layout inputs; dashed inferred; optional scope cue |
| E | `GraphOrientPanel.tsx` | Scope + inference legend honesty |
| F | Board Notes | Smoke evidence + commands |

**Forbidden:** `internal/retrieval` rewrites; raising caps; `web/` membership inference walk.

---

## Exit criteria

- [ ] Cluster layout uses **API `scope_id` only** (Law 19)
- [ ] `EDGE_PRIORITY` elevates MVP scope rels; `goal_has_task` demoted when scope edges exist
- [ ] Inferred edges dashed via `provenance` (no `inferred_` rel mutation)
- [ ] Browser smoke: **≥2** visible scope clusters + **≥1** cross-scope edge with seeded fixture
- [ ] `PROJECT_MAX_NODES` still **500**; `EDGE_OVERVIEW_MAX` still **150**
- [ ] Truncation / omitted banner still honest when truncated
- [ ] `npm test` / web unit tests + build green (repo’s usual web check)
- [ ] Board `done`; next **P44-S04-02**

## Minimal todos

- [ ] Types + preserve `scope_id` through compose
- [ ] Scope cluster force layout
- [ ] Tiered edge priority + inferred dash style
- [ ] Orient legend
- [ ] Unit tests + browser smoke Notes
- [ ] Board Notes

## Next

`P44-S04-02`
