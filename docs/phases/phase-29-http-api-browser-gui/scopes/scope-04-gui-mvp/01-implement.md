# P29-S04-01 — GUI MVP implementer

## Metadata
- id: P29-S04-01
- todo_ids: [P29-S04-01]
- role: implementer
- skills: [frontend-ui-engineering, incremental-implementation]
- mcps: [cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Scaffold and ship the **S04 MVP SPA** under `web/` against the live HTTP API (`trace serve` / `internal/httpapi`). Honor UX-IA `gui_ship: S04` only. **Law 19:** browser calls `/v1` only — no second SoT, no SQLite, no invented transition policy.

**Do not** start P29-S04-02. **Do not** implement S05 surfaces (reviews, xyflow, loop write, seed export/import, discovery create).

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults**
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [UX-IA.md](../scope-02-ux-ia/UX-IA.md) — screens, empty/error, ship matrix
- [ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../api/openapi.yaml)
- Live static: [`internal/httpapi/static.go`](../../../../../internal/httpapi/static.go) — `web/dist/index.html` replaces placeholder

## Session start

Follow agent-loop-protocol Session start. Human locks + S00–S03 + S04-00 defaults are **settled** — do not re-debate stack, bind, CORS, or IA.

## Locked defaults

| Item | Value |
|------|-------|
| Root | `web/` (private npm package) |
| Stack | TypeScript + Vite + React 18+ |
| Router | `react-router-dom` BrowserRouter |
| CSS | Plain CSS + CSS variables; no Tailwind/MUI/shadcn |
| State | Context (theme, token, baseUrl, project/health) + local screen state; no Redux/Zustand/TanStack Query |
| Types | `npx openapi-typescript ../../api/openapi.yaml -o src/api/schema.d.ts` (adjust relative path) |
| Client | Hand wrappers in `src/api/` by `operationId`; `fetch` to relative `/v1/...`; attach Bearer when token set |
| Build out | `web/dist` → consumed by `trace serve` (default StaticDir `<root>/web/dist`) |
| Vite proxy | Dev only: `/v1` → `http://127.0.0.1:7432` |
| Search JSON | `{items:[{id,kind,title,snippet?}]}` — never assume CLI `hits` |
| Graph S04 | Stub: center input + `getGraph?center=&max_nodes=` (cap e.g. 50) **or** entity/task list + “Open in graph (S05)” — show `truncated` honesty banner when true |
| Agents | **No UI.** If a helper exists, recommend requires `task_id` or `phase` (live API) |

### Package layout (conceptual minimum)

```
web/
  package.json
  vite.config.ts          // outDir: dist; proxy /v1
  index.html
  tsconfig.json
  src/
    main.tsx
    App.tsx               // BrowserRouter + routes
    styles/
      tokens.css          // CSS variables; light/dark via data-theme
      app.css
    api/
      schema.d.ts         // generated
      client.ts           // base fetch + envelope error type
      ops.ts              // getHealth, listTasks, … by operationId
    context/
      AppChrome.tsx       // theme, token, project, health, version
    layout/
      Shell.tsx           // chrome + left nav / bottom tabs
      Nav.tsx
    screens/
      ProjectGate.tsx
      Overview.tsx
      Tasks.tsx
      TaskDetail.tsx
      Loop.tsx            // read-only
      GraphStub.tsx
      Discoveries.tsx
      Seed.tsx            // status + honesty
      Settings.tsx
    components/           // ErrorBanner, EmptyState, GateStrip, CopyButton, …
```

### Routes (UX-IA)

| Path | Screen | Notes |
|------|--------|-------|
| `/` | Overview | Also reachable as Overview nav |
| `/tasks` | Tasks list | |
| `/tasks/:taskId` | Task detail | Copy `TRACE_TASK_ID=<uuid>` |
| `/loop` | Loop read-only | Status + gate; no next/apply/reset controls |
| `/graph` | Graph stub | Budgeted one-shot or placeholder |
| `/discoveries` | Discoveries read | search + getEntity detail |
| `/seed` | Seed | Status only + honesty copy |
| `/settings` | Settings | Theme, token, version/project/bind display |

**Project gate:** when `getProject` fails or `store_ready === false`, render full-page blocked empty (hide store-dependent destinations). Retry → `getHealth` / `getProject`.

### S04 operation map (call these; do not invent)

| Screen | operationIds |
|--------|----------------|
| Chrome / gate | `getHealth`, `getVersion`, `getProject` |
| Overview | `listTasks`, `search` (goals summary), `getLoopStatus`, `getLoopGate` |
| Tasks | `listTasks`, `getTask`, optional `createTransition` (show envelope on deny); optional `getContext` / `getWhy` drawers |
| Loop | `getLoopStatus`, `getLoopGate` only |
| Graph stub | `search` and/or task pick + `getGraph` (require `center` + `max_nodes`) |
| Discoveries | `search` (kind filter if API allows via `q`), `getEntity` |
| Seed | `getSeedStatus` only |
| Settings | client-only theme/token + display from `getVersion` / `getProject` |

**Forbidden in S04 UI:** `getLoopNext`, `postLoopApply`, `postLoopReset`, `postSeedExport`, `postSeedImport`, `createEntity`, `createLink`, `listReviews` / review writes, agents screens, defer ops.

### Error / empty (ADR envelope)

Parse `{error:{code,message,details}}`. Map UX-IA §4: `UNAUTHORIZED` → Settings; `FORBIDDEN`/`CONFLICT`/`VALIDATION_ERROR` on transition → show message only (no SPA policy). Seed honesty banner always visible on Seed screen.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -d internal/httpapi
grep -q 'case "serve"' cmd/trace/root.go
test -f api/openapi.yaml
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md
# web/ may be absent at start — create it this row
```

## Preflight / Plan

Short plan then execute. Vertical slices (each leaves a buildable app):

1. Vite+React+TS scaffold + tokens CSS + Shell/nav + Settings (theme/token) + proxy
2. API client + types; Project gate (health/project/version)
3. Tasks list + detail + TRACE_TASK_ID copy (+ light transition)
4. Overview + Loop read-only (status/gate strip)
5. Discoveries read + Graph stub + Seed status/honesty
6. `npm run build` → confirm `web/dist/index.html`; smoke with `trace serve`

## Role work

1. Implement per locked layout and S04 op map.
2. Empty/error states per UX-IA for each shipped screen.
3. Ensure production build is served by existing `trace serve` static path (no second server).
4. Board **P29-S04-01** status + Notes only (build commands, smoke steps, key paths). Do not edit 02-review substance.

## Todo updates

Status + notes on **P29-S04-01** only.

## Exit criteria

- [ ] `cd web && npm install && npm run build` PASS → `web/dist/index.html` exists
- [ ] Smoke: `trace serve` (loopback) → open UI → **tasks list visible** (or honest empty if store has no tasks)
- [ ] All UX-IA S04 screens present (gate, overview, tasks, loop RO, graph stub, discoveries read, seed status, settings)
- [ ] No S05 CTAs that call forbidden write ops (loop apply/reset, seed export/import, discovery create)
- [ ] No browser SQLite / parallel SoT; all data via `/v1`
- [ ] Search uses `items[]` shape
- [ ] Transition denials show API envelope only
- [ ] Basic a11y: focusable nav, labeled icon buttons, gate/status text not icon-only
- [ ] Board Notes with build + smoke evidence

## Minimal todos

- [ ] Scaffold `web/` (Vite React-TS) + CSS tokens + Shell/nav
- [ ] `openapi-typescript` + `src/api/client.ts` + ops wrappers
- [ ] Project gate + chrome health/version/project
- [ ] Tasks list/detail + TRACE_TASK_ID + light transition
- [ ] Overview + Loop read-only
- [ ] Discoveries read + Graph stub + Seed status/honesty + Settings
- [ ] Production build replaces placeholder; document smoke in Notes

## Next

**P29-S04-02**
