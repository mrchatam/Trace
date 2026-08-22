# Phase 33 S01 — UX-IA (Explore overview)

## Job

Open **Explore (`/`)** → interactive **project overview** graph hook (Graphify energy: pan / zoom / click) under Laws **6–7**. Inspector and panels remain. Success: first paint shows a budgeted, seed-composed subgraph — **not** an empty “Pick center” gate as the default hook.

Binding S00 model: **(D) + (B) + (C)** — derive seeds from ops → parallel budgeted `getGraph` → merge client-side → **progressive expand** on user action. API preference: **`reuse`** (`getProject` / `listTasks` / `search` / `getGraph`). Never unbounded dump; never seed-export-as-graph-body.

## Route clarity

| Route | Role | Phase 33 target? |
|-------|------|------------------|
| **`/`** (Graph) | **Explore** — interactive project graph | **Yes** — this IA |
| `/graph` | Redirects to `/` | Same |
| **`/overview`** | Nav **Overview** — ops summary screen | **No** — do not conflate; do not redesign as Explore |

Nav labels stay: Explore → `/`; Overview → `/overview`.

## Open sequence (D+B+C)

### Step 0 — Enter Explore

User navigates to `/` (or `trace gui` lands here). Show **loading** chrome immediately (canvas skeleton + “Building overview…” / budget placeholder). Do **not** block on a center picker as the hero.

### Step 1 — (D) Seed composition

Build an ordered seed id list (hard stop ≤ **8**; target **6**):

1. **`getProject`** — if the project/root is graph-addressable as a node id, include it; else keep project metadata for chrome only (**do not invent fake nodes**).
2. **`listTasks`** — prefer `IN_PROGRESS`, then non-terminal (`PENDING` / other active). Skip `DONE` / `SKIPPED` unless needed to reach a minimum useful set (e.g. empty active list).
3. **`search`** — fill remaining slots with high-signal entities (goals, capabilities, decisions, discoveries) under existing search caps.
4. **Dedupe** by id; stop at seed cap.

If zero seeds after this: go to **Empty** (below). Manual center search remains **secondary** (toolbar), not the first-paint hero.

### Step 2 — (B) Parallel budgeted neighborhoods

For each seed id, call **`getGraph(center=seed, max_nodes=40, depth=2)`** in parallel (reuse `ops.ts`). On individual failure: record error; continue with successes (**partial fail**).

### Step 3 — Merge under UI cap

Client-side merge/dedupe nodes+edges. Enforce visible-node cap **`UI_CAP=100`** (trim policy: keep all seed nodes first, then neighbors by proximity to seeds / edge degree; drop extras with a quiet “+N omitted” budget note if needed). Layout: multi-seed overview (force or deterministic radial clusters by seed) — visual grouping by **kind/state** only (RESEARCH **(A)** inspiration); **no** Leiden/community API.

### Step 4 — First paint

Render merged graph on canvas; select nothing or first seed; inspector idle or shows project summary. Budget chrome: seeds used, nodes shown / cap, depth.

### Step 5 — (C) Progressive expand

User actions only (no expand-all default):

| Action | Behavior |
|--------|----------|
| Click node | Select → inspector |
| Double-click / explicit Expand | Re-center: `getGraph(center=node, max_nodes≤50, depth=2)`; replace or union per implementer choice — **prefer replace neighborhood around new center** while keeping overview seeds in a “pin” set if cheap; never load-all |
| Pan / zoom | Canvas navigation (xyflow Controls) |
| Kind filter (if present) | Client filter of visible nodes; does not fetch unbounded data |

## Budgets

| Lean | Locked default | Hard bound | Rationale |
|------|----------------|------------|-----------|
| Seed count | **Target 6** | ≤ **8** (acceptable 4–8) | Enough for project shape; avoids N×max blowup |
| Per-seed `max_nodes` | **40** | ≤ **50** (`DEFAULT_MAX` today) | 6×40=240 raw → merge trims to UI cap |
| Merged UI visible | Honor **`UI_CAP=100`** | Argue ≤120 only with written rationale; **never >150** | Keep Laws 6–7; **no bump this phase** — stay at 100 |
| Depth | **2** | OpenAPI default | Neighborhood, not whole graph |
| Expand | User-driven re-center / expand only | No “load all” / expand-all default | Law 7 |

## Progressive expand

- **Re-center** means: chosen node becomes graph `center`; fetch a new budgeted neighborhood; update layout; selection stays on that node unless fetch fails.
- Overview “seed pins” (optional): seeds remain visually marked (`--accent` center treatment) when still in the visible set after expand.
- **No** “Load entire project graph” control in Phase 33.
- Clusters: visual only (kind color / state chip); no community endpoint.

## Inspector

- **Retained** (Phase 32 inspector). Overview does not remove it.
- Selection on overview: same as today — click node → inspector loads entity / task context via existing adapters (Law 19).
- After expand/re-center: inspector follows selection; if selection drops out of graph, clear or keep last entity with “not in current view” note (S03 chooses; prefer keep last + soft note).

## Empty / loading / error

| State | Behavior | Copy intent (not CSS) |
|-------|----------|------------------------|
| **Loading** | Skeleton canvas + non-blocking budget line | “Building project overview…” |
| **Happy** | Seed-composed interactive graph | (no empty hero) |
| **No seeds** | Helpful empty; secondary manual center | “Nothing to seed yet. Add tasks or entities, or pick a center to explore a neighborhood.” CTAs: link to Tasks; optional center search. Manual center is **secondary**, not the hero. |
| **Partial API fail** | Show loaded subgraph + non-blocking banner | Cause + Retry for failed seeds; do **not** blank canvas if any seed succeeded |
| **Hard fail** | Error banner + Retry; empty canvas OK | No fake full dump; no silent success |

## Kind / state on canvas (IA note)

Encoding details and tokens: [`DESIGN.md`](DESIGN.md). IA requirement: every visible node exposes **kind label**; task nodes expose **work_state** text when known. Color is not the sole encoder.

## S03 vs S04 ownership

| Owner | Owns | Does not own |
|-------|------|--------------|
| **S03** | Seed → `getGraph` → merge → progressive expand; Graph UX (pan/zoom/click); keep inspector; budget chrome; empty/error per this doc | Full shell recolor; redesigning Nav Overview (`/overview`); new SPA; Three.js; unbounded graph API |
| **S04** | Land `DESIGN.md` tokens into `tokens.css` / `app.css`; colorize chrome, nav, pills, Explore nodes (kind/state), empty/error surfaces | Re-opening this IA model; PATH / `trace gui` |

## S03 handoff checklist

- [ ] On Explore mount: run seed pipeline (getProject → listTasks priority → search fill → dedupe ≤8).
- [ ] Parallel `getGraph` per seed (`max_nodes=40`, `depth=2`); merge/dedupe; enforce `UI_CAP=100`.
- [ ] First paint = merged overview (not EmptyState “Pick center” as default).
- [ ] Keep manual center / search as secondary toolbar affordance.
- [ ] Click → select + inspector; double-click or Expand → user-driven re-center fetch; no load-all.
- [ ] Empty / partial / hard-fail states per table above.
- [ ] Attach `data-kind` / `data-state` (or equivalent) for S04 styling; do not invent palette.
- [ ] Do not change `/overview` ops screen as Explore; do not add Leiden API; do not use seed-export as graph body.
- [ ] Prefer reuse of `web/src/api/ops.ts`; Law 19 adapters only.
