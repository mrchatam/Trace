# P29-S04-00 — Scope planner (GUI MVP)

## Metadata
- id: P29-S04-00
- todo_ids: [P29-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-ui-engineering]
- verification: automated

## Objective

Lock MVP GUI implement/review prompts against UX-IA P0 + live API. **No product SPA code this row** — thicken `01-implement.md` + `02-review.md` so a fresh agent can build S04 alone.

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Root | `web/` — Vite app; build `outDir` = `dist` (path `<project>/web/dist`) |
| Stack | TypeScript + Vite + **React 18+**. **No** Next.js, Vue, Svelte, Three.js, or `@xyflow/*` in S04 |
| Router | `react-router-dom` **BrowserRouter**; route paths = UX-IA nav defaults (`/`, `/tasks`, `/tasks/:id`, `/loop`, `/graph`, `/discoveries`, `/seed`, `/settings`) |
| Styling | Plain CSS + CSS variables (dense operator chrome). **No** Tailwind, MUI, shadcn, or CSS-in-JS frameworks |
| State | Thin `fetch` API client + React Context for theme / bearer token / base URL / project+health chrome. **No** Redux, Zustand, or TanStack Query in S04 |
| Client | Generate types with `openapi-typescript` from `api/openapi.yaml` → `web/src/api/schema.d.ts`. Hand-written wrappers in `web/src/api/` keyed by OpenAPI `operationId`. Same-origin relative `/v1` when served by `trace serve` |
| Vite DX | `server.proxy["/v1"]` → `http://127.0.0.1:7432`. Vite-dev CORS origin reflect stays **S06** (S03 CORS remains deny / no `*`) |
| Serve | Static via same `trace serve`. S03 serves placeholder HTML if `web/dist/index.html` missing; S04 build **replaces** that — do not fork a second static server |
| Search | HTTP `SearchResponse` = `{items:[{id,kind,title,snippet?}]}` (OpenAPI SoT) — **not** CLI/MCP `{ok,hits,count}` |
| Agents | Live `GET /v1/agents?action=list\|recommend` — recommend needs `task_id` **or** `phase` (live httpapi; OpenAPI may omit `phase` — call live shape). **No agents screen in S04** (`gui_ship: S05+`) |
| Chrome | Operator/agent product UI — not marketing |
| Layout | Shell: **left nav** ≥768px; **bottom tabs** narrow. Chrome strip: project title + health/version. Full-page **blocked** empty when `store_ready: false` or project unreachable |
| Auth | Loopback-trust (no token required). Settings: paste bearer → `Authorization: Bearer` on all `/v1` calls when set. `UNAUTHORIZED` → banner + Settings |
| IA SoT | [`../scope-02-ux-ia/UX-IA.md`](../scope-02-ux-ia/UX-IA.md) — honor `gui_ship: S04` vs S05; do not promote reviews or rich graph |
| Law 19 | `/v1` only; no browser SQLite / IndexedDB / parallel SoT; transition policy = **library error envelopes only**, not SPA DONE/review rules |

## P0 (MVP) — `gui_ship: S04` from UX-IA

- Project / open workspace (health + `getProject` readiness / blocked empty)
- Overview: goals summary via search/tasks, active task, **loop status + gate** strip
- Tasks list + detail read + `TRACE_TASK_ID`; optional light `createTransition` (envelope-only on deny)
- Loop: **read-only** status/gate (full next/apply/reset → S05)
- Graph: **stub only** (center picker + one budgeted `getGraph` or list + S05 placeholder) — never unbounded viz; **no xyflow**
- Discoveries: **read** list/search + detail (create/promote → S05)
- Seed: **status + honesty** warnings (export/import actions → S05)
- Settings: theme + token paste + bind/version/project display

## Explicitly not S04

Reviews UI; rich graph; loop next/apply/reset; discovery create/promote; seed export/import CTAs; agents nav; cloud OAuth; MCP `/rpc`.

## Exit criteria (scope planner)

- [x] Stack / layout / API client locked above
- [x] `01-implement.md` + `02-review.md` thick enough for a fresh agent
- [x] `SCOPE-TODOS.md` synced to board IDs 515–517
- [ ] Scope product exit (S04-01): `npm run build` succeeds; smoke via `trace serve` shows tasks — **implementer row**

## Next

P29-S04-01 → P29-S04-02
