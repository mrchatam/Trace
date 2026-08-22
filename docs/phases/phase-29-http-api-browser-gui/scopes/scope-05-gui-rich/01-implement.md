# P29-S05-01 — Feature-rich GUI implementer

## Metadata
- id: P29-S05-01
- todo_ids: [P29-S05-01]
- role: implementer
- skills: [frontend-ui-engineering, incremental-implementation, webapp-testing]
- mcps: [cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Ship **GUI-P1** on the live S04 SPA under `web/`: rich budgeted graph, loop write console, discoveries create/promote, seed export/import, fuller task transitions, reviews. Track every row in [`FEATURE-MATRIX.md`](FEATURE-MATRIX.md). **Law 19:** `/v1` → library only — no SPA-encoded DONE/review policy, no second SoT.

**Do not** start P29-S05-02. **Do not** rewrite the MVP shell/routes/stack. **Do not** implement S06 hardening (CORS reflect, embed, mapDomainErr).

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults**
- [FEATURE-MATRIX.md](FEATURE-MATRIX.md) — **tracking SoT for this scope**
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §6–7, §19
- [UX-IA.md](../scope-02-ux-ia/UX-IA.md) — `gui_ship: S05`
- [ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../api/openapi.yaml)
- Live baseline: `web/` (S04 MVP), `internal/httpapi`, `trace serve`

## Session start

Follow agent-loop-protocol Session start. Human locks + S00–S04 + S05-00 defaults are **settled** — do not re-debate stack, bind, CORS, Three.js, or IA ship matrix.

## Locked defaults

| Item | Value |
|------|-------|
| Baseline | **Extend** S04 `web/` — keep Shell/Nav/Context/css tokens/api client patterns |
| Stack | TypeScript + Vite + React; BrowserRouter; plain CSS vars; Context + `fetch` (no Redux/TQ) |
| Types | `npm run gen:api` → `src/api/schema.d.ts`; add hand wrappers in `ops.ts` by `operationId` |
| Graph lib | **`@xyflow/react`** — 2D expand-on-demand. **Forbidden:** Three.js, full-project dump UX |
| Graph budget | Default `max_nodes=50`; UI must not offer unbounded; always pass `center` + `max_nodes`; banner when `truncated===true` |
| Loop writes | `getLoopNext` / `postLoopApply` / `postLoopReset` with **confirm dialog** (keyboard: Enter confirm, Escape dismiss) for apply + reset |
| Loop honesty | Keep `Promise.allSettled` (or equivalent) so status failure does **not** blank gate (S04 residual) |
| Discoveries | `createEntity` + `createLink` + `createTransition` for create/promote; preserve initial search |
| Seed | `postSeedExport` / `postSeedImport`; project-relative paths; show path-confinement / 501 `strict`/`task_id` envelopes honestly |
| Reviews | New nav + `/reviews`, `/reviews/:reviewId`; `listReviews` / `getReview`; `createReview` optional from task detail |
| Search JSON | `{items:[…]}` only — never CLI `hits` |
| Optional | `listCapability` enrichment; `listAgents` recommend drawer — mark `deferred` in matrix if skipped |
| Deferred (no GUI) | O03–O06 in FEATURE-MATRIX (impact/plans/token API/index/SSE/…) |
| Tests | Promote + export honesty: Playwright under `web/e2e/` **or** documented browser smoke + cite `go test ./internal/httpapi/...` seed path cases |

### Package / file map (extend)

```
web/
  package.json                 // add @xyflow/react; optional Playwright scripts
  src/
    api/ops.ts                 // + loop next/apply/reset, entity/link, seed export/import, reviews, capability?
    layout/Nav.tsx             // + Reviews
    App.tsx                    // + /reviews routes; GraphStub → Graph
    screens/
      Graph.tsx                // was GraphStub — xyflow canvas
      Loop.tsx                 // write CTAs + confirms
      Discoveries.tsx          // create + promote
      Seed.tsx                 // export/import forms
      Tasks.tsx / TaskDetail.tsx  // fuller transitions + gate-aware DONE UX
      Reviews.tsx / ReviewDetail.tsx  // new
    components/
      ConfirmDialog.tsx        // apply/reset (+ destructive seed import if needed)
  e2e/                         // optional Playwright — promote + seed path honesty
```

### Routes (delta from S04)

| Path | Screen | Change |
|------|--------|--------|
| `/loop` | Loop | Add next / apply / reset + gate detail |
| `/graph` | Graph | xyflow explorer (replace stub list viz) |
| `/discoveries` | Discoveries | Create + promote CTAs |
| `/seed` | Seed | Export/import actions |
| `/tasks/:taskId` | Task detail | Fuller transitions; gate strip / warning display |
| `/reviews` | Reviews list | **New** |
| `/reviews/:reviewId` | Review detail | **New** |

### S05 operation map (call these; do not invent)

| Screen | operationIds |
|--------|----------------|
| Graph | `getGraph`, `search` (center picker) |
| Loop | `getLoopStatus`, `getLoopGate`, `getLoopNext`, `postLoopApply`, `postLoopReset` |
| Discoveries | `search`, `getEntity`, `createEntity`, `createLink`, `createTransition`; optional `listCapability` |
| Seed | `getSeedStatus`, `postSeedExport`, `postSeedImport` |
| Tasks | `listTasks`, `getTask`, `createTransition`, `getLoopGate`; optional `getContext`/`getWhy`/`listAgents` |
| Reviews | `listReviews`, `getReview`; optional `createReview` |

**Forbidden in S05 UI:** OpenAPI `defer` ops; unbounded graph; Three.js; inventing transition policy client-side; treating search `hits` as real.

## Preflight

1. `web/` + `web/dist` exist from S04; `case "serve"` in CLI; `internal/httpapi` present.
2. `FEATURE-MATRIX.md` present (S05-00).
3. `cd web && npm run build` still PASS before large edits (sanity).

## Role work

1. Implement matrix **M01–M07** (and optional O01/O02 or mark deferred with reason).
2. Update `FEATURE-MATRIX.md` Status → `done` / `deferred` + evidence columns.
3. Satisfy **G-promote** + **G-export** evidence gates.
4. Self-check exit criteria; board Notes with build + test commands.

## Minimal todos

- [ ] Add `@xyflow/react`; evolve Graph stub → budgeted expand-on-demand canvas + truncated banner (M01)
- [ ] Loop write CTAs + ConfirmDialog; preserve independent status/gate loads (M02)
- [ ] Discoveries create + promote via createEntity/link/transition (M03)
- [ ] Seed export/import UI + path/501 honesty (M04)
- [ ] Task detail fuller transitions + gate-aware DONE display (M05)
- [ ] Reviews nav + list/detail (M06); search filter chrome as needed (M07)
- [ ] Optional O01/O02 or defer with reason in matrix
- [ ] G-promote + G-export evidence; `npm run build` PASS; update FEATURE-MATRIX rollup
- [ ] Board Notes only on **P29-S05-01**

## Exit criteria

- [ ] Every FEATURE-MATRIX P1 row (`M01–M07`) is `done` or `deferred` with reason; optionals settled
- [ ] G-promote + G-export evidenced in Notes (commands / paths)
- [ ] No unbounded graph / Three.js; Law 19 `/v1` only; SearchResponse `items`
- [ ] Loop status failure does not blank gate
- [ ] `cd web && npm run build` PASS → `web/dist`
- [ ] Board Notes with test/smoke commands

## Todo updates

Status + notes on **P29-S05-01** only.

## Next

**P29-S05-02**
