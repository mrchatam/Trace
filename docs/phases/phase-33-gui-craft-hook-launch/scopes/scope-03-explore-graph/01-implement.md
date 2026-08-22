# P33-S03-01 — Implement Explore project graph

## Metadata
- id: P33-S03-01
- todo_ids: [P33-S03-01]
- role: implementer
- skills: [frontend-design, impeccable, ui-ux-pro-max, frontend-ui-engineering]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Make Explore (`/`) feel like a Graphify-class **interactive project graph** on open: pan/zoom/click, budgeted overview per S01 UX-IA, inspector retained. Laws **6–7**, Law **19**. Skills required. Full shell colorize remains **S04** (attach hooks; optional token-class placeholders OK).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Theme B
- [00-PLANNER.md](00-PLANNER.md) — locked defaults + Must-answer (resolved)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md) — model **(D)+(B)+(C)**; reject full dump / Leiden / seed-export-as-graph-body
- [`../scope-01-design-ux/UX-IA.md`](../scope-01-design-ux/UX-IA.md) — open sequence, budgets, empty/error table, §S03 handoff
- [`../scope-01-design-ux/DESIGN.md`](../scope-01-design-ux/DESIGN.md) — `data-kind` / `data-state` + kind labels (do **not** invent palette)
- Live: `web/src/screens/Graph.tsx`, `web/src/api/ops.ts`, `web/src/App.tsx`, `web/src/components/Inspector.tsx`, `web/e2e/s03-depth.spec.ts`, `api/openapi.yaml`

## Session start

Follow agent-loop-protocol Session start. **Before large UI edits**, read/apply: `frontend-design`, `impeccable` (craft-floor / operate restraint — not full colorize), `ui-ux-pro-max` (dense developer-tool / Operate signals; reject purple/cream/broadsheet, HUD glow). Note skill use in board Notes.

## Locked defaults (planner — do not re-debate)

| Item | Value |
|------|-------|
| Home | Keep **`/`** as Explore graph — **S02 contract:** `trace gui` lands `http://{addr}/` (**not** `/overview`). Do not relocate graph home |
| IA | UX-IA **(D)+(B)+(C)**; Explore ≠ Nav Overview |
| Stack | Evolve existing `@xyflow/react` 2D in `Graph.tsx` — **no** second SPA, **no** Three.js |
| API | **reuse** `ops.ts` only: `getProject`, `listTasks`, `search`, `getGraph`. No new overview endpoint; no Leiden; no seed-export-as-graph-body |
| Seed target / cap | Target **6**, hard ≤ **8** |
| Per-seed fetch | `getGraph(center=seed, max_nodes=**40**, depth=**2**)` parallel |
| Merge cap | Honor **`UI_CAP=100`** visible nodes (never >150; no bump this phase) |
| Expand | User-driven re-center: `getGraph(…, max_nodes≤**50**, depth=2)`; **prefer replace** neighborhood around new center; optional seed “pin” class if seed still visible; **no** load-all / expand-all |
| `getProject` | Chrome/metadata only unless response later gains a real graph-addressable id — **today** `ProjectResponse` = `root` + `store_ready` → **do not invent a fake project node** |
| Task seed order | One `listTasks({ limit: 50 })` (or equivalent): prefer **`IN_PROGRESS`**, then non-terminal (`PENDING`, `AWAITING_REVIEW`, `BLOCKED`, …); include `DONE`/`SKIPPED` **only** if needed to reach a useful minimum |
| Search fill | `q` is **required** (empty → 400). After tasks, fill remaining slots with non-empty FTS queries for high-signal kinds — e.g. sequential `search("goal"\|"capability"\|"decision"\|"discovery", limit≈20)` — take hits not already seeded; stop at seed cap. Skip inventing nodes |
| First paint | Loading → merged overview canvas (**not** EmptyState “Pick center” / “No center selected” as default hero) |
| Manual center | Keep search + task pick as **secondary toolbar** (collapse/demote “Pick center” panel; do not delete affordance) |
| Keyboard | Chrome + node-list buttons remain fully operable; canvas selection path required or residual risk explicit (see Must-answer #5) |
| Color / hooks | Attach `data-kind` (+ `data-state` when work_state known on tasks); kind **text** always visible; may wire CSS vars from DESIGN.md lightly — **S04** owns shell-wide colorize |
| Out | Three.js; second SQLite; business logic in `web/`; `/overview` as Explore; Leiden; seed-export graph body; unbounded dump; inventing palette; PATH/`trace gui` changes |

### Must-answer locks (resolved by P33-S03-00)

1. **Open data path** — On Explore mount: (0) show loading chrome → (D) `getProject` (chrome) + prioritize `listTasks` + `search` fill → dedupe ≤8 → (B) `Promise.allSettled`-style parallel `getGraph(seed, 40, 2)` → merge/dedupe nodes+edges → trim to `UI_CAP=100` (keep all seed ids first, then neighbors by seed proximity / degree) → (4) first paint → (C) expand only on user action.
2. **Interaction smoke** — Pan/zoom (xyflow Controls); click → select + inspector; double-click / Expand / “Use as center” → re-center fetch ≤50. Happy path: first paint shows interactive nodes without mandatory empty gate. Truly empty store → Empty copy per UX-IA (not the old “pick center” hero).
3. **Empty / partial / hard-fail** — Per UX-IA table: Loading “Building project overview…”; No seeds → helpful empty + Tasks link + secondary center; Partial → subgraph + non-blocking banner + Retry failed seeds; Hard → ErrorBanner + Retry, empty canvas OK. Never fake success / full dump.
4. **Files** — Primary `web/src/screens/Graph.tsx`; optional pure helper `web/src/lib/overviewCompose.ts` (or similar) for seed+merge (unit-testable); styles only as needed for hooks (`app.css` / existing graph classes). Do **not** change `App.tsx` route away from `/`; do **not** redesign `Overview.tsx`.
5. **Keyboard** — **Required:** tab order through toolbar → node list (`graph-select-node-*`) → inspector; visible focus on selected list/canvas node; no focus trap in ReactFlow. **Canvas:** best-effort (`nodesFocusable` / selected node tabbable, or arrow-key roving if cheap). If canvas keyboard remains limited, **document residual risk in Notes** — chrome+list+inspector must still select/expand without pointer.

## Role work

### 1. Overview-on-open (D+B+C)

Replace center-first empty gate with UX-IA open sequence in `Graph.tsx` (extract helper OK):

```text
mount /
  → set overviewLoading
  → seeds = composeSeeds(getProject, listTasks, search)  // ≤8, target 6
  → if seeds.length === 0 → Empty (no-seeds copy)
  → else parallel getGraph(seed, 40, 2)
  → merge + trim UI_CAP=100
  → set graph + budget chrome (seeds used, nodes/cap, depth)
  → first paint (select nothing or first seed; inspector idle OK)
```

Constants (names flexible; values fixed):

| Constant | Value |
|----------|-------|
| `SEED_TARGET` | 6 |
| `SEED_CAP` | 8 |
| `SEED_MAX_NODES` | 40 |
| `UI_CAP` | 100 |
| `EXPAND_MAX_NODES` | 50 |
| `DEPTH` | 2 |

### 2. Progressive expand + inspector

- Click node → `selectedId` + Inspector (existing adapters).
- Expand / double-click / Use as center → `getGraph(center, ≤50, 2)`; prefer **replace** visible neighborhood; mark overview seeds if still present (`graph-node--center` / seed class + moss accent OK).
- If selection leaves view after expand: **keep last entity** + soft “not in current view” note (UX-IA preference).
- Kind filter (if kept): **client-only** on visible nodes — no unbounded refetch.

### 3. Empty / error chrome

| State | UI |
|-------|-----|
| Loading | Skeleton / status: “Building project overview…” + budget placeholder |
| Happy | Canvas + budget line (seeds, nodes/cap, depth) |
| No seeds | EmptyState copy per UX-IA; CTA `Link` to `/tasks`; secondary manual center |
| Partial | Merged successes + warn banner (failed seed count/ids) + Retry |
| Hard | `ErrorBanner` + Retry; canvas empty OK |

Demote existing “Pick center” panel to secondary (details/collapsed OK); keep `#graph-budget` / testids stable where e2e depends on them, or update e2e in same row.

### 4. `data-kind` / `data-state` hooks

On each canvas node (and list row if cheap): `data-kind={kind}`; for tasks with known `work_state`, `data-state={work_state}`. Keep `.graph-node__kind` text. Do not invent new colors beyond DESIGN.md token names.

### 5. Skills + craft restraint

Apply frontend-design / impeccable / ui-ux-pro-max for hierarchy, focus, empty/error clarity — **not** a full shell recolor (S04). Avoid purple/cream/broadsheet / glow HUD.

### 6. Tests / smoke evidence

| Check | Expect |
|-------|--------|
| Open `/` with seeded store | Canvas (or node list) shows overview nodes **without** clicking pick-task first |
| Budgets | No call with `max_nodes` > UI path caps; overview uses 40; expand ≤50; visible ≤100 |
| Partial | Force one bad seed id in unit/helper test **or** document manual — UI keeps successes |
| E2E | Update `web/e2e/s03-depth.spec.ts` (and `s05-gates` if it assumes pick-first) for overview-first; select ≠ expand still holds |
| Build | `npm run build` in `web/` (or project’s usual web script) green |

Evidence in Notes: screenshot path **or** e2e command + pass; skills named.

## Exit criteria

- [ ] Open Explore shows interactive project graph (screenshot or e2e)
- [ ] Laws 6–7: seed ≤8; per-seed 40/depth 2; UI_CAP=100; expand ≤50; no load-all
- [ ] Empty / partial / hard-fail per UX-IA
- [ ] `data-kind` / `data-state` (or equiv.) present for S04
- [ ] Keyboard path or explicit residual risk in Notes; no trap
- [ ] SPA root still `/` (S02 land intact)
- [ ] Skills noted in board Notes
- [ ] No Law 19 fork; no product changes outside web Explore path

## Minimal todos

- [ ] Seed compose + parallel getGraph + merge under UI_CAP
- [ ] First-paint overview; demote manual center; expand/inspector
- [ ] Empty / partial / hard-fail + hooks
- [ ] E2E/smoke + `npm run build`
- [ ] Board Notes (skills + evidence + any keyboard residual)

## Todo updates

Status + notes on **P33-S03-01** only.

## Next

`P33-S03-02`
