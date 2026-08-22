# P33-S03-02 — Explore graph review

**Verdict:** **PASS**  
**Confidence:** **high**  
**Date:** 2026-08-21  
**Scope:** Explore overview-on-open (Theme B / D+B+C); Laws 6–7 / 19; S02 land `/`  
**Skills loaded:** frontend-design · impeccable (Operate / craft-floor restraint) · ui-ux-pro-max (dense dev-tool a11y) · code-review-and-quality  
**Evidence:** `node --experimental-strip-types --test src/lib/overviewCompose.test.ts` 7/7; `npm run build` ok; `npm run test:e2e -- e2e/s03-depth.spec.ts` pass (overview-first, select≠expand, `data-kind`)

## Checklist evidence

### Theme B / open experience
- [x] Open `/` runs `loadOverview` on mount — not EmptyState “Pick center” / “No center selected” as happy path
- [x] Loading copy: “Building project overview…” (`graph-overview-loading`)
- [x] Manual center demoted to collapsed `<details data-testid="graph-manual-center">` — secondary, not hero gate

### Laws 6–7 + progressive expand
- [x] Seed pipeline in `Graph.tsx` + `overviewCompose.ts`: `getProject` chrome-only (no fake node) → `listTasks` IN_PROGRESS→non-terminal → `search` fill (`goal|capability|decision|discovery`) → dedupe; **target 6 / ≤8**
- [x] Parallel `Promise.allSettled` → `getGraph(seed, SEED_MAX_NODES=40, DEPTH=2)`
- [x] `mergeOverviewGraphs` trims to **`UI_CAP=100`**, seeds first then proximity/degree
- [x] Expand via `loadGraph` caps `max_nodes≤50` (`EXPAND_MAX_NODES`); user-driven only; no load-all / expand-all / Leiden / seed-export body
- [x] API **reuse** only (`ops.ts`: `getProject`, `listTasks`, `search`, `getGraph`)

### Interaction + inspector
- [x] Pan/zoom (ReactFlow Controls); click → select + Inspector; canvas + list share `onSelect`
- [x] Select ≠ re-center (e2e asserts center stable until “Use as center” / expand)
- [x] Kind filter is client-side on `visibleNodes` / search list only

### Empty / error
- [x] No seeds → EmptyState + Tasks/Overview links; partial → banner + Retry + keeps subgraph; hard → ErrorBanner + Retry

### Keyboard / a11y
- [x] Toolbar → node list buttons (`graph-select-node-*`) → inspector path; kind/state text visible (color-not-only)
- [x] Canvas: `nodesFocusable` + selected `tabIndex={0}` best-effort — **residual accepted** (full canvas arrow-roving not shipped; chrome+list+inspector sufficient)

### Routing / S02 contract
- [x] `App.tsx`: `index` → `<Graph />`; `/overview` remains ops Overview; `/graph` redirects to `/`
- [x] Explore stays at `/` so `trace gui` → `http://{addr}/` still lands Explore

### S04 hooks + craft boundary
- [x] `data-kind` / `data-state` on list buttons + canvas nodes; kind label text always shown
- [x] No invented shell palette; S04 owns colorize
- [x] Implementer Notes cite skills

### Law 19
- [x] `web/` adapters + pure compose helpers only; no parallel SQLite / business-logic fork

### Evidence
- [x] Unit 7/7 + build + e2e cited above

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| low | No-seeds EmptyState CTA is page-footer `Link to="/tasks"` rather than inline on EmptyState | S04/S05 polish OK; not a Theme B blocker |
| residual | Canvas arrow-key roving not shipped | Accepted per Must-answer #5 / implementer Notes; list+inspector cover AC |

**No blocker / high.** No spawn rows (`P33-S03-02a`/`02b`).

## Upcoming thickenings (reviewer rights)

- **S04** — Hooks already on Explore nodes (`data-kind`/`data-state`, `.graph-node__kind` / `__state`, seed/center classes). Land DESIGN tokens onto those selectors; do not rework seed compose or relocate `/`. Preserve kind text labels (color-not-only). Residual canvas keyboard stays out of S04 unless cheap with craft.
- **S05** — Docs primary flip already planned; optional note that Explore open path is overview-composed (not center-gate).

## Next runnable

**P33-S04-00**
