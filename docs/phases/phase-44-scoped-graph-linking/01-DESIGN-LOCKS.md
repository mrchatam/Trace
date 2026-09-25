# Phase 44 — Design locks

**Status:** LOCKED (S01) — 2026-08-23. Implementers must not contradict locked rows without a forward REVIEW spawn. Wording refinements closed via P44-S01-01; D1–D5 choices unchanged from P44-00 light locks.

## 1. Why the graph looks like “everything → first goal”

| Factor | Mechanism | File (verified P44-00) |
|--------|-----------|------------------------|
| Center selection | `ProjectGraph` sets `center` to **first goal** in kind-sorted list | `internal/retrieval/project_graph.go:55–64` |
| Primary edge | `goal_has_task` synthesized from `tasks.goal_id` for **every** task | `internal/retrieval/graph_neighbors.go:50–86` |
| Sparse cross-links | MCP `trace_link` exposes only 5 rel aliases; most cognitive rels created by domain workflows, not scope | `internal/mcp/tools_write.go:183–195` |
| No scope rels | `entity_links` has causal/planning rels but **no** `scope_member`, `api_contract`, feature-front-back | `internal/domain/service.go:33–64` |
| Plan scopes disconnected | `plan_scopes` table exists for progressive planner; **not** wired into graph walk | `internal/store/plan_hierarchy.go` |
| GUI edge cap | Overview renders ≤150 edges, prioritizing `goal_has_task` | `web/src/lib/graphLayout.ts:36–48` |

**Conclusion:** Topology is honest for **goal-centric causal model**; it is **not** honest for **feature/scope cartography**. Phase 44 adds scope semantics without removing goal/task moat.

---

## 2. Scope as first-class concept

### Options considered

| Option | Model | Pros | Cons |
|--------|-------|------|------|
| **A — Link-only** | `scope_member` (entity → scope slug/id), `scope_related` between scopes | No new entity table; seed-export friendly | Scope “node” may be synthetic in GUI |
| **B — Scope entity** | New `scope` entity type + `scope_member` edges | Visible scope nodes in graph | Migration + export surface |
| **C — Tag column** | `scopes[]` JSON on task/decision/discovery | Simple queries | Weak cross-entity joins; not in graph walk today |

### D1 — LOCKED

**Primary: A + thin scope records (B-lite).**

- Persist lightweight **`scope` records**: `id`, `slug`, `title`, `kind` (`feature` \| `layer` \| `business`).
- Membership via **`scope_member`** in existing `entity_links` (Law 13 — no parallel relationship store).
- Cross-scope via typed rels in §3.
- **`plan_scopes` are planner hierarchy rows**, not feature/layer/business cartography. “Reuse/map where titles align” means **soft title/body alignment for S03 inference** — **not** a hard FK, and **not** “every `plan_scope` becomes a graph scope node” in S02 MVP.
- **Reject C as primary** (orphan tags without graph walk). Pure A (synthetic-only, no record table) only if seed-export cannot host thin records — default remains **A + thin records** (RESEARCH did not prove export blocker).

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

#### Name collisions (D2 vs existing causal)

| Proposed D2 string | Must stay distinct from (existing) | Why |
|--------------------|------------------------------------|-----|
| `implements` | `change_implements_decision` | Task/decision peer edge ≠ change→decision side-effect |
| `blocks` | `uncertainty_blocks_task` | Task↔task scope edge ≠ uncertainty gate |

**Require:** domain consts, MCP/CLI aliases, seed allowlists, and docs use the **exact** D2 strings; do not alias D2 to the longer causal names.

### 3.4 Provenance

| Source | `entity_links.source_type` | Graph `rel` suffix? |
|--------|---------------------------|---------------------|
| Agent MCP/CLI | `USER_ASSERTED` | no |
| Import seed | `IMPORTED` | no |
| Inference job | `INFERRED` | no — optional `edge.provenance` / dashed UI only |

**Lock:** Inference **never deletes** explicit links; lower confidence; GUI dashed edge style. **Do not** require UI to mutate stored `rel` strings with an `inferred_` prefix; dashed styling / provenance field only.

### D2 — LOCKED

**MVP minimum rel set:** `scope_member`, `api_contract`, `implements`, `blocks`.

Optional if cheap in S02 (same allowlist/migration pass): `same_feature_front_back`, `depends_on`, `relates_to`. Do not block S02 MVP on optional rels. Presentation tokens in GUI may already list some of these; that is **not** evidence they persist at HEAD (see §6).

---

## 4. Inference vs explicit

| Path | When | Owner scope |
|------|------|-------------|
| **Explicit** | Agent `trace_link`, domain side-effects (reviews, regressions) | S02 |
| **Inference** | Batch/on-demand: file path rules (`web/` vs `internal/`), task title tokens, shared plan_scope, OpenAPI tag | S03 |
| **Rejected** | LLM-only graph fill without evidence row | — |

### D4 — LOCKED

Inference runs **opt-in CLI only** (`trace graph infer` or equivalent documented command). **No** post-`trace index` hook in P44; **no** silent daemon. Index-hook deferred to residual / future phase.

---

## 5. Graph API & cap policy

### 5.1 Current caps (baseline — verified P44-00)

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

### 5.3 Phase 44 API extensions

| Change | Description |
|--------|-------------|
| `scope_id` on `GraphNode` | Optional — populated from `scope_member` |
| `edge.provenance` | `explicit` \| `inferred` (optional field; map from `source_type`: `USER_ASSERTED`/`IMPORTED` → explicit; `INFERRED` → inferred) |
| `scope` query param | Filter project graph to one scope (+ N-hop) — **bounded** by `max_nodes` |

### D3 — LOCKED

**API shape:** `GET /v1/graph?mode=project&scope=<slug|id>&max_nodes=…` returns members of that scope plus bounded N-hop neighbors (default depth 1–2, hard-capped by `max_nodes`). Reject missing `max_nodes`. **No** pagination/tiles in P44.

Optional `scope_id` on `GraphNode` (from `scope_member`). Optional `edge.provenance` on wire as above — **not** an `inferred_` prefix on stored `rel` strings (align §3.4).

### D5 — LOCKED

Keep **GUI `PROJECT_MAX_NODES = 500`**. API hard max **5000** unchanged. Prefer **scope-filtered** project mode over raising global cap. Edge overview **150** stays presentation policy. Tiles/pagination deferred unless VERIFY proves 500 insufficient *with* scope filter.

---

## 6. UI consequences (S04)

- Force layout: **cluster by scope** (scope centroids + weak inter-scope links).
- Edge styling: causal vs scope vs inferred (width/dash).
- Reduce `goal_has_task` visual dominance when scope edges exist (re-prioritize `EDGE_PRIORITY_RELS` to elevate MVP scope rels).
- `EDGE_PRIORITY_RELS` may already list `blocks` / `depends_on` / `relates_to` — treat as **presentation prep**, not proof that scope semantics exist until S02 persists MVP rels. S04 elevates real MVP rels after they exist.
- Orient panel: document scope legend + inference honesty.
- **Not sufficient alone:** relayout without new edges recreates hairball.
- **Law 19:** GUI/HTTP adapters thin; walk/filter logic in `internal/retrieval`.

---

## 7. Non-goals (phase 1 cut)

- Vector embeddings / automatic semantic scopes (DR-NOSSEM)
- Real-time sync with Codegraph symbol graph
- New MCP write tool beyond extending `trace_link` rel allowlist
- Multi-tenant hosted graph queries
- Replacing `docs/TODO.md` board order with graph scope order
- Auto-infer `review_judges_task` or other **gate** rels
- Silent index-hook inference (D4)
- Raising GUI default above 500 without VERIFY evidence (D5)

---

## 8. Portable graph / migration

- Expanding D2 MVP rels **and** thin `scope` records requires: migration (if new table/entity) + **`seedExportLinkRels` / export allowlist** update + `trace seed export -o trace/graph.json` before PRs that change Trace entities (CONTRIBUTING).
- Today export allowlist is **four causal rels only** (`seed_export.go` ~341–346) — S02 must add MVP scope rels explicitly.
- INFERRED links export with honest `source_type`; re-import idempotent.
- Do **not** reopen D1–D5 while stating export expectations — export policy is an implementation obligation under the locked choices.

---

## 9. Closed decisions (LOCKED)

| ID | Decision | Final choice | Rationale |
|----|----------|--------------|-----------|
| D1 | Scope model | **A + thin scope records**; reject C as primary; plan_scopes soft-map only (not hard FK / not every plan_scope → graph scope in S02) | Graph walk + visible clusters; Law 13; RESEARCH did not prove export blocker for thin records |
| D2 | MVP rel set | `scope_member`, `api_contract`, `implements`, `blocks` (keep distinct from `change_implements_decision` / `uncertainty_blocks_task`) | Covers INTAKE auth FE/BE + intra-scope; collision-safe exact strings |
| D3 | Scope filter API | `scope` query param + bounded N-hop; optional `scope_id` / `edge.provenance`; **no tiles** | Law 6–7; provenance from `source_type`, not `inferred_` rel mutation |
| D4 | Inference trigger | **CLI only** (`trace graph infer` or documented equiv.) | No silent daemon / index-hook; opt-in honesty |
| D5 | GUI cap | **Stay 500**; prefer scope filter; edge overview 150 presentation-only | Perf + Law 6–7; API 5000 unchanged |

Wording refinements closed via P44-S01-01 (RESEARCH §7). Reversing a locked choice requires REVIEW spawn with reason.

---

## 10. Acceptance map (VERIFY)

- [ ] Seeded auth FE/BE scopes → ≥2 clusters + ≥1 cross-scope edge (GUI smoke)
- [ ] `trace_link` / CLI can create `scope_member` + `api_contract`
- [ ] Inferred edges visually distinct; explicit wins
- [ ] Laws 6–7: no unbounded graph; truncation honest
- [ ] Law 19: walk/filter in `internal/retrieval`; GUI adapter only
- [ ] Law 13 / M-001: membership+MVP rels in `entity_links`; goal/task loop preserved
