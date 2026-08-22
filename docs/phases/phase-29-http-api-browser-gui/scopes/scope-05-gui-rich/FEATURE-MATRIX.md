# Phase 29 / Scope 05 — FEATURE-MATRIX (GUI-P1)

**Authority:** [`../scope-02-ux-ia/UX-IA.md`](../scope-02-ux-ia/UX-IA.md) `gui_ship: S05` + [`api/openapi.yaml`](../../../../../api/openapi.yaml).  
**Baseline:** Live S04 SPA under `web/` — **extend**, do not rewrite MVP.  
**Status vocabulary (this file):** `planned` (S05-01 must ship) · `optional` (ship if cheap; may leave `deferred`) · `deferred` (reason required) · `done` (S05-01 marks after evidence).

S05-00 authored this matrix. **S05-01** updates Status → `done` / keeps `deferred` with reason. **S05-02** verifies every row.

---

## P1 checklist (must-cover)

| ID | Area | Ship target | Primary operationIds | Status | Notes / evidence |
|----|------|-------------|----------------------|--------|------------------|
| M01 | Rich graph explorer | xyflow-class 2D canvas; expand-on-demand; budgeted `getGraph`; truncated honesty banner | `getGraph`, `search` | done | `web/src/screens/Graph.tsx` + `@xyflow/react`. Default `max_nodes=50`, UI cap 100. Click node → re-center. Truncated banner when `truncated===true`. No Three.js. |
| M02 | Loop console write | next / apply / reset + gate detail; confirm dialogs on apply/reset | `getLoopStatus`, `getLoopGate`, `getLoopNext`, `postLoopApply`, `postLoopReset` | done | `Loop.tsx` + `ConfirmDialog` (Enter/Escape). Independent status/gate via `Promise.allSettled`. |
| M03 | Discoveries create + promote | Create discovery/decision; relate; promote → task | `createEntity`, `createLink`, `createTransition`, `search`, `getEntity` | done | Create form + promote = task → `discovery-mentions-task` link → optional transition. Envelope on deny. Initial search preserved. |
| M04 | Seed export / import | Actions + path-confinement error honesty | `getSeedStatus`, `postSeedExport`, `postSeedImport` | done | Project-relative paths; summary-only job panel; strict/`task_id` → 501 envelope; import confirm. |
| M05 | Tasks full transitions | Fuller transition UX + gate-aware DONE | `createTransition`, `getTask`, `getLoopGate`, optional `getContext`/`getWhy` | done | All-states select + suggested; DONE gate strip; `allow_done` checkbox; warning display; createReview from detail. Law 19: no SPA policy. |
| M06 | Reviews list / detail | Primary nav + list/detail (+ create from task if cheap) | `listReviews`, `getReview`, optional `createReview` | done | Nav + `/reviews`, `/reviews/:reviewId`; create from task detail. |
| M07 | Search chrome / filters | Kind/q filters as needed for Discoveries + Graph picker | `search` | done | Kind filter chrome on Discoveries + Graph; response stays `{items:[…]}`. |

### Evidence gates (exit)

| Gate | Requirement | Owner |
|------|-------------|-------|
| G-promote | E2E **or** integration coverage that promote path calls library-backed create/link/transition and surfaces envelope on failure | S05-01 — `cd web && npm run test:e2e` (`e2e/s05-gates.spec.ts` G-promote) |
| G-export | E2E **or** integration + UI honesty for path escape / confinement on seed export or import | S05-01 — Playwright G-export + cite `go test ./internal/httpapi/ -run 'TestSeedPathConfinement|TestSeedExportSummaryOnly'` |

---

## Optional / enrichment (FEATURE-MATRIX allows)

| ID | Area | Ship target | operationIds | Status | Deferral / reason |
|----|------|-------------|----------------|--------|-------------------|
| O01 | Capability catalog | Enrich Discoveries/Tasks detail (read-only catalog) | `listCapability` | deferred | Time-boxed: enrichment only, not SoT; S06+ OK |
| O02 | Agents recommend | Task-detail drawer “recommend harness” | `listAgents` (`action=recommend` + `task_id`\|`phase`) | deferred | No must-cover screen (UX-IA S05+); primary P1 shipped without drawer |
| O03 | Overview widgets | Impact / changes chips | `getImpact`, `listChanges` | deferred | UX-IA S05+ optional; not blocking P1 console/graph/seed |
| O04 | Plans / regressions | Any UI | `listPlans`, `listRegressions` | deferred | No must-cover screen; avoid scope creep |
| O05 | Settings token API | Generate bearer via API | `postAuthToken` | deferred | OpenAPI `x-trace-gui: false`; paste token remains S04 |
| O06 | Index / defer APIs | Any GUI | index, backup, restore, migrate, install, patterns, knowledge, eval, events SSE | deferred | `x-trace-wave: defer` / `x-trace-gui: false` — no GUI |

---

## Cross-cutting locks (do not re-debate)

| Lock | Value |
|------|-------|
| Law 19 | Browser → `/v1` only; hand ops in `web/src/api/ops.ts`; no second SoT |
| Law 6–7 | Budgeted graph; seed status/summary honesty (never full dump default) |
| Stack | Extend S04: TS + Vite + React, BrowserRouter, CSS vars, Context+fetch |
| Graph lib | `@xyflow/react` (xyflow-class). **Not** Three.js |
| Search JSON | `SearchResponse.items` only |
| Loop residual | Dogfood IDs like `rl…` may 500 `INTERNAL_ERROR` on status; SPA keeps gate; S06 prefers `VALIDATION_ERROR` via `mapDomainErr` |
| Nav (S05) | Overview · Tasks · Loop · Graph · Discoveries · Reviews · Seed · Settings |

---

## S05-01 status rollup (fill at implement close)

| Metric | Value |
|--------|-------|
| planned → done | M01–M07 (7/7) |
| deferred with reason | O01, O02 (optional); O03–O06 (locked deferred) |
| G-promote evidence | `cd web && npm run test:e2e` — `e2e/s05-gates.spec.ts` G-promote (API + UI) PASS |
| G-export evidence | Playwright G-export PASS; `go test ./internal/httpapi/ -run 'TestSeedPathConfinement\|TestSeedExportSummaryOnly'` PASS |
