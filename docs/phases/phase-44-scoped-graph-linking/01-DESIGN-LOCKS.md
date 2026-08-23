# Phase 44 — Design locks (draft)

**Status:** Draft for S01 lock. Implementers must not contradict locked rows without a forward spawn.

## 1. Why the graph looks like “everything → first goal”

| Factor | Mechanism | File |
|--------|-----------|------|
| Center selection | `ProjectGraph` sets `center` to **first goal** in kind-sorted list | `project_graph.go:55–64` |
| Primary edge | `goal_has_task` synthesized from `tasks.goal_id` for **every** task | `graph_neighbors.go:50–86` |
| Sparse cross-links | MCP `trace_link` exposes only 5 rel aliases; most cognitive rels created by domain workflows, not scope | `tools_write.go:183–195` |
| No scope rels | `entity_links` has causal/planning rels (`decision_affects_task`, …) but **no** `scope_member`, `api_contract`, feature-front-back | `domain/service.go:33–64` |
| Plan scopes disconnected | `plan_scopes` table exists for progressive planner; **not** wired into graph walk | `store/plan_hierarchy.go` |
| GUI edge cap | Overview renders ≤150 edges, prioritizing `goal_has_task` | `graphLayout.ts:36–48` |

**Conclusion:** Topology is honest for **goal-centric causal model**; it is **not** honest for **feature/scope cartography**. Phase 44 adds scope semantics without removing goal/task moat.

---

## 2. Scope as first-class concept

**Lock candidate (S01 chooses one primary + one optional):**

| Option | Model | Pros | Cons |
|--------|-------|------|------|
| **A — Link-only** | `scope_member` (entity → scope slug/id), `scope_related` between scopes | No new entity table; seed-export friendly | Scope “node” may be synthetic in GUI |
| **B — Scope entity** | New `scope` entity type + `scope_member` edges | Visible scope nodes in graph | Migration + export surface |
| **C — Tag column** | `scopes[]` JSON on task/decision/discovery | Simple queries | Weak cross-entity joins; not in graph walk today |

**Draft recommendation:** **A + optional B** — introduce lightweight **`scope` records** (id, slug, title, kind: `feature|layer|business`) linked via `scope_member`; cross-scope via typed rels below. Reuse `plan_scopes.id` where titles align (mapping table or inference — S03).

---

## 3. Link taxonomy (phase 1)

### 3.1 Intra-scope (same feature area)

| Rel | From → To | Explicit | Inference hint |
|-----|-----------|----------|----------------|
| `scope_member` | any entity → scope | agent / import | task title prefix, path glob |
| `implements` | task/decision → task | agent | — |
| `blocks` | task → task | agent | dependency text |
| `relates_to` | any → any (same scope) | agent | co-occurrence in plan scope body |

### 3.2 Cross-scope / front-back

| Rel | From → To | Example |
|-----|-----------|---------|
| `api_contract` | task (frontend) → task (backend) | “Login POST /v1/auth” |
| `same_feature_front_back` | scope ↔ scope | auth-FE ↔ auth-BE bundle |
| `depends_on` | scope/task → scope/task | billing depends on auth |

### 3.3 Preserved causal rels (unchanged)

`decision_affects_task`, `discovery_causes_plan_change`, `review_judges_task`, `goal_has_task` (via `goal_id`), etc. — **still authoritative** for loop/gate/review.

### 3.4 Provenance

| Source | `entity_links.source_type` | Graph `rel` suffix? |
|--------|---------------------------|---------------------|
| Agent MCP/CLI | `USER_ASSERTED` | no |
| Import seed | `IMPORTED` | no |
| Inference job | `INFERRED` | optional `inferred_` prefix in UI only — store canonical rel |

**Lock:** Inference **never deletes** explicit links; lower confidence; GUI dashed edge style.

---

## 4. Inference vs explicit

| Path | When | Owner scope |
|------|------|-------------|
| **Explicit** | Agent `trace_link`, domain side-effects (reviews, regressions) | S02 |
| **Inference** | Batch/on-index: file path rules (`web/` vs `internal/`), task title tokens, shared plan_scope, OpenAPI tag | S03 |
| **Rejected** | LLM-only graph fill without evidence row | — |

Inference runs **opt-in** (`trace graph infer` or post-`trace index` hook) — no silent daemon.

---

## 5. Graph API & cap policy

### 5.1 Current caps (baseline)

| Layer | Cap | Rationale |
|-------|-----|-----------|
| OpenAPI `/v1/graph` | `max_nodes` **required**, max **5000** | Law 6 — reject unbounded |
| Go retrieval | `MaxNeighborhoodNodes = 5000` | Hard reject above cap |
| GUI project default | **`PROJECT_MAX_NODES = 500`** | Browser perf + readable overview (Phase 40 G5) |
| GUI neighborhood expand | 50 nodes, depth 2 | Progressive drill-down (Law 7) |
| GUI edge render | 150 overview edges | Anti-hairball presentation |

### 5.2 Why 500 (user question)

1. **Law 6:** No full-project dump by default — explicit budget.
2. **Law 7:** Overview is orient entry; deeper views use `center` + smaller budgets.
3. **OpenAPI:** Allows up to 5000 for power users/CLI; GUI chooses conservative 500.
4. **Perf:** Force layout (`d3-force`) is O(n²) per tick; 500 nodes × simulation ≈ acceptable in browser; 5000 is not.
5. **Honesty:** `truncated: true` + `total_entities` must display when project exceeds budget.

### 5.3 Phase 44 API extensions (draft)

| Change | Description |
|--------|-------------|
| `scope_id` on `GraphNode` | Optional — populated from `scope_member` |
| `edge.provenance` | `explicit` \| `inferred` (optional field) |
| `scope` query param | Filter project graph to one scope (+ N-hop) — **bounded** |
| Pagination/tiles | **Defer** unless VERIFY proves 500 insufficient with scope filter; prefer scope filter over raising global cap |

**Lock candidate:** Keep **500 GUI default**; add **scope-filtered** project mode before raising cap.

---

## 6. UI consequences (S04)

- Force layout: **cluster by scope** (scope centroids + weak inter-scope links).
- Edge styling: causal vs scope vs inferred (width/dash).
- Reduce `goal_has_task` visual dominance when scope edges exist (re-prioritize `EDGE_PRIORITY_RELS`).
- Orient panel: document scope legend + inference honesty.
- **Not sufficient alone:** relayout without new edges recreates hairball.

---

## 7. Non-goals (phase 1 cut)

- Vector embeddings / automatic semantic scopes
- Real-time sync with Codegraph symbol graph
- New MCP write tool beyond extending `trace_link` rel allowlist
- Multi-tenant hosted graph queries
- Replacing `docs/TODO.md` board order with graph scope order
- Auto-infer `review_judges_task` or other **gate** rels

---

## 8. Portable graph / migration

- Any new rel or `scope` table → migration + `seed_export` allowlist + `trace/graph.json` export before PR (CONTRIBUTING).
- Infer links: export with `source_type=INFERRED`; re-import idempotent.

---

## 9. Open decisions (S01 must close)

| ID | Decision |
|----|----------|
| D1 | Scope model: A/B/C final |
| D2 | Minimum rel set for MVP (suggest: `scope_member`, `api_contract`, `implements`, `blocks`) |
| D3 | Scope filter API shape |
| D4 | Inference trigger: CLI only vs index hook |
| D5 | GUI cap: stay 500 vs scoped tiles |

---

## 10. Acceptance map (VERIFY)

- [ ] Project graph with seeded auth FE/BE scopes shows **≥2 clusters** and **≥1 cross-scope** edge in GUI smoke
- [ ] `trace_link` (or CLI) can create `scope_member` + `api_contract`
- [ ] Inferred edges visually distinct; explicit wins on conflict
- [ ] Laws 6–7 tests/grep: no unbounded graph route; truncation honest
- [ ] Law 19: GUI adapter only; graph walk in `internal/retrieval`
