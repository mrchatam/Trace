# Phase 44 / S02 / Implement backend links

## Metadata
- id: P44-S02-01
- todo_ids: [P44-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd, api-and-interface-design]
- mcps: [user-trace]
- verification: automated

## Objective

Ship scope rels in `entity_links`, thin scope records (per D1), graph walk / project graph fields, `trace_link` aliases, seed export/import, OpenAPI `GraphNode`/`GraphEdge` + `scope` query param. **Prerequisite:** S01 locks APPROVE (P44-S01-02). Inference out of scope (S03). GUI out of scope (S04).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — **LOCKED (S01)**; D1–D3, §3, §8
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) §2, §5, §7
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- CONTRIBUTING — portable `trace seed export -o trace/graph.json` before PRs that change Trace entities

## Session start

Follow agent-loop-protocol Session start. Block if S01 review not APPROVE.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Library first | Domain + store + `internal/retrieval`; HTTP/MCP thin adapters (Law 19) |
| Store | Membership + MVP rels in **`entity_links` only** (Law 13) — no parallel relationship store |
| Scope records | **New** thin table `scopes` (D1 A+B-lite). **Not** `plan_scopes` (planner hierarchy; soft-map only in S03) |
| MVP rels (exact strings) | `scope_member`, `api_contract`, `implements`, `blocks` |
| Collisions | **Distinct** from `change_implements_decision` / `uncertainty_blocks_task` — never alias D2 → causal names |
| Optional rels | `same_feature_front_back`, `depends_on`, `relates_to` — same allowlist pass **if cheap**; do not block MVP |
| Graph API | Optional `scope_id` on nodes; optional `edge.provenance`; `scope` query param; `max_nodes` required; hard max **5000** |
| Caps | Do **not** raise API 5000 or GUI 500 |
| Inference | Out of scope — do **not** invent `INFERRED` rows here (S03) |
| GUI | Out of scope — do **not** edit `web/` |
| MCP | Extend `trace_link` allowlist only — **no** new write tool |
| Portable | Migration + seed allowlist + `trace/graph.json` before PR |

### Preflight anchors @ HEAD (re-verify if drift)

| Claim | Location |
|-------|----------|
| No scope MVP rels | `internal/domain/service.go:33–64` |
| MCP 5 kebab aliases only | `internal/mcp/tools_write.go:183–195` + `cmd/trace/link.go:51–64` |
| Seed export 4 causal rels | `internal/domain/seed_export.go` `seedExportLinkRels` ~341–346 |
| `GraphNode` = id/kind/title/goal_id | `internal/retrieval/neighborhood.go:22–27` |
| `GraphEdge` = rel/from/to | `internal/retrieval/neighborhood.go:30–34` |
| ProjectGraph center = first goal | `internal/retrieval/project_graph.go:55–64` |
| Walk = entity_links + synthesized `goal_has_task` | `internal/retrieval/graph_neighbors.go` |
| HTTP project mode ignores `scope` | `internal/httpapi/handlers_retrieval.go:158–195` |
| Latest schema | `028_*.sql` → next migration **`029_*.sql`** |
| `plan_scopes` exists, unwired to graph walk | `006_plan_hierarchy.sql` / `plan_hierarchy.go` |

---

## Migration strategy

### M1 — Additive schema only

Add **`internal/store/schema/029_graph_scopes.sql`** (name may vary; version **029**):

```sql
-- Migration v29: thin graph scope records (Phase 44 S02).
-- Additive only; do not rewrite 001–028.
-- plan_scopes (v6) remain planner hierarchy — NOT this table.

CREATE TABLE IF NOT EXISTS scopes (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    kind TEXT NOT NULL
        CHECK (kind IN ('feature', 'layer', 'business')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_scopes_kind ON scopes(kind);
```

**Rules:**
- Do **not** ALTER `plan_scopes` for cartography.
- Do **not** add CHECK on `entity_links.rel` (rels stay free-string + domain consts).
- Membership is **`entity_links`**: `rel=scope_member`, `to_type=scope`, `to_id=<scopes.id>` (from any allowed entity type).
- No FK from `entity_links` to `scopes` required (match existing link style); domain validates endpoints exist.

### M2 — Domain entity type

| Const | Value | Notes |
|-------|-------|-------|
| `EntityScope` | `"scope"` | Graph cartography — **≠** `EntityPlanScope` / `"plan_scope"` |

### M3 — Seed / portable surface

| Change | Detail |
|--------|--------|
| `SeedDocument` | New array e.g. `scopes` / `graph_scopes` (prefer **`scopes`**) with `{id, slug, title, kind}` |
| `seedExportLinkRels` | Append four MVP strings (exact) |
| `ImportSeedLink` | Cases for MVP rels (+ optional if shipped); `source_type` honest (`IMPORTED` / existing import path) |
| Import order | Upsert thin scopes **before** links that reference them |
| CONTRIBUTING | After entity/schema change: `trace seed export -o trace/graph.json` |

### M4 — Rollback / compat

- Fresh DBs apply 029 via embed migrate.
- Existing `.trace/trace.db` migrate on open (existing path).
- Old seeds without `scopes` remain valid (empty array).
- Unknown optional rels stay rejected until added.

---

## File touch-list (ordered vertical slices)

Implement bottom-up; keep each slice testable.

### Slice A — Schema + store CRUD

| Path | Do |
|------|----|
| `internal/store/schema/029_graph_scopes.sql` | Create `scopes` table (M1) |
| `internal/store/scopes.go` (new) | `UpsertScope`, `GetScope`, `GetScopeBySlug`, `ListScopes` |
| `internal/store/store_test.go` / scoped test | Migration applies; table listed if schema inventory asserts tables |
| `internal/store/links.go` | No schema change; reuse `InsertLink` / list helpers |

### Slice B — Domain consts + link APIs

| Path | Do |
|------|----|
| `internal/domain/service.go` | Add `EntityScope`; `RelScopeMember`, `RelAPIContract`, `RelImplements`, `RelBlocks` (exact D2). Keep causal consts untouched |
| `internal/domain/doc.go` | Document MVP scope rels + collision note |
| `internal/domain/link.go` or `scope.go` (new) | `CreateScope` / upsert; `LinkScopeMember`, `LinkAPIContract`, `LinkImplements`, `LinkBlocks` (validate endpoints; `InsertLink` + `appendLinked`; default `USER_ASSERTED`) |
| Endpoint rules (locked intent) | `scope_member`: any supported entity → `scope`; `api_contract`: typically task→task; `implements`: task/decision → task; `blocks`: task → task. Reject wrong types fail-closed |

### Slice C — MCP + CLI (adapters)

| Path | Do |
|------|----|
| `internal/mcp/tools_write.go` | Extend `toolLink` switch + `LinkInput` jsonschema string |
| `cmd/trace/link.go` | Mirror same aliases |
| `cmd/trace/help.go` | Mention new rels |
| Alias map (recommend) | Accept **exact** D2 (`scope_member`, …) **and** kebab (`scope-member`, `api-contract`). `implements` / `blocks` are identical kebab/snake. Stored `rel` = exact D2 |
| Forbidden | New MCP tool; aliasing to `change_implements_decision` / `uncertainty_blocks_task` |

### Slice D — Seed export / import

| Path | Do |
|------|----|
| `internal/domain/seed_export.go` | Export `scopes`; extend `seedExportLinkRels` with MVP four |
| `internal/domain/seed_import.go` | Import scopes; `ImportSeedLink` MVP cases; idempotent duplicates |
| `cmd/trace/seed.go` | Allowlist key for `scopes` if keyed import filters exist |
| Tests | Round-trip: create scope + MVP links → export → fresh DB import → same links/scopes |
| Portable | `trace seed export -o trace/graph.json` before PR |

### Slice E — Retrieval + HTTP + OpenAPI

| Path | Do |
|------|----|
| `internal/retrieval/neighborhood.go` | `GraphNode.ScopeID` optional; `GraphEdge.Provenance` optional (`explicit`\|`inferred`) |
| `internal/retrieval/exact.go` (`lookupEntity`) | Resolve `scope` via store GetScope |
| `internal/retrieval/project_graph.go` | Include scope nodes in `collectProjectNodes`; kind order entry for `"scope"`; populate `scope_id` from `scope_member` when cheap; `ProjectGraphOpts.Scope` (slug\|id) |
| Scope filter algorithm | Resolve scope → collect member ids via `scope_member` → seed set = members (+ scope node) → expand N-hop (default depth **1**, allow **2**, hard-capped by `max_nodes`) via existing walk → edges among included nodes. Missing/`max_nodes` invalid → error (Law 6). **No** tiles/pagination |
| Provenance mapping | From link `source_type`: `USER_ASSERTED`/`IMPORTED` → `explicit`; `INFERRED` → `inferred`. S02 may leave provenance empty when unset; do **not** rewrite stored `rel` with `inferred_` prefix |
| `internal/retrieval/graph_neighbors.go` | Prefer unchanged (all `entity_links` already walked); ensure scope endpoints resolve |
| `internal/httpapi/handlers_retrieval.go` | Pass `scope` query into `ProjectGraphOpts`; keep `max_nodes` required |
| `api/openapi.yaml` | Param `scope` on `/v1/graph`; `GraphNode.scope_id`; `GraphEdge.provenance` enum |

### Slice F — Tests (acceptance — must pass)

See **Acceptance tests** below. Prefer `go test` under `internal/domain`, `internal/store`, `internal/retrieval`, plus focused MCP/CLI if existing patterns allow without heavy harness.

### Out of touch-list (forbidden this row)

- `web/**` (S04)
- Inference CLI / `INFERRED` writers (S03)
- Raising caps; pagination tiles
- Hard FK `plan_scopes` ↔ `scopes`
- New MCP write tool

---

## Acceptance tests (implementer checklist)

### A1 — Domain / store MVP rels

- [ ] Insert/link each MVP rel succeeds with exact string (`scope_member`, `api_contract`, `implements`, `blocks`).
- [ ] Const values ≠ `change_implements_decision` / `uncertainty_blocks_task` (compile-time or unit assert).
- [ ] Thin scope CRUD: upsert by id/slug; `kind` ∈ {feature, layer, business}; invalid kind rejected.
- [ ] `scope_member` targets `to_type=scope` and existing scope id.
- [ ] Duplicate endpoints no-op / unique constraint (match existing `entity_links` UNIQUE).

### A2 — ProjectGraph includes scope edges

- [ ] With two tasks + one scope + `scope_member` + `api_contract` (and budget large enough), `ProjectGraph` edges include those MVP rels (both endpoints in node set).
- [ ] Scope nodes appear with `kind=scope` when scopes exist and fit budget.
- [ ] Optional: member nodes expose `scope_id` when membership exists.
- [ ] `plan_scopes` rows alone do **not** appear as `kind=scope` nodes.

### A3 — Scope filter + Law 6

- [ ] `ProjectGraph` / HTTP with `scope=<slug|id>` + `max_nodes=N` returns only members + bounded N-hop (≤ N nodes).
- [ ] Missing `max_nodes` still rejected.
- [ ] `max_nodes > 5000` → `ErrBudgetExceeded` / HTTP `BUDGET_EXCEEDED`.
- [ ] Unrelated entities outside scope+hops absent from nodes.

### A4 — MCP / CLI

- [ ] `trace_link` (or CLI mirror) creates `scope_member` and `api_contract` (exact stored rel).
- [ ] Unknown rel still rejected; existing five aliases unchanged.
- [ ] Help/schema text lists new aliases.

### A5 — Seed export round-trip (see also `02-review.md`)

- [ ] `seedExportLinkRels` contains all four MVP rels.
- [ ] Export includes thin `scopes` array.
- [ ] Import into empty store restores scopes + MVP links idempotently.
- [ ] After schema/entity change: `trace/graph.json` refreshed for PR.

### A6 — Non-goals smoke

- [ ] No `INFERRED` rows created by S02 paths.
- [ ] No edits under `web/`.
- [ ] Causal goal/task / review gate rels still work (M-001 smoke if cheap).

---

## Suggested implement order (Minimal todos)

1. [ ] Migration 029 + store scope CRUD + tests
2. [ ] Domain consts + CreateScope + four link methods + A1 tests
3. [ ] Retrieval: lookup scope, ProjectGraph nodes/edges, scope filter, optional fields + A2/A3
4. [ ] MCP + CLI aliases + A4
5. [ ] Seed export/import + A5; `trace seed export -o trace/graph.json`
6. [ ] OpenAPI + HTTP `scope` param
7. [ ] Board Notes → **P44-S02-02**

## Exit criteria

- [ ] Unit tests for MVP rels + ProjectGraph includes scope edges
- [ ] `scope` filter bounded by `max_nodes`; unbounded rejected
- [ ] `trace_link` creates `scope_member` + `api_contract`
- [ ] Seed export allowlist updated; `trace/graph.json` if entities changed
- [ ] No name collision with `change_implements_decision` / `uncertainty_blocks_task`
- [ ] Board `done`; next **P44-S02-02**

## Next

`P44-S02-02`
