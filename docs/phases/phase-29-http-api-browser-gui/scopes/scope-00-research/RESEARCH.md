# RESEARCH — Phase 29 HTTP API + browser GUI

**Date:** 2026-08-21  
**Row:** P29-S00-01

## Executive summary

Phase 29 should ship an opt-in local `trace serve` that exposes a versioned JSON/OpenAPI surface and serves a TypeScript + Vite + React SPA under `web/`, with handlers calling only the canonical Go library (Law 19). Peer survey of Understand-Anything dashboard, codebase-memory-mcp graph-ui, and agentrq frontend confirms React+Vite as the dominant local-graph pattern; Vue is viable for task boards but does not overturn the locked lean. Security defaults worth copying: bind `127.0.0.1`, no CORS wildcard, optional one-time token for data endpoints, and **bounded** graph payloads (never full-graph dump as default). Trace CLI/MCP inventories (preflight PASS; no `serve` yet) map cleanly into resource families for S01 OpenAPI, with P0 centered on health, meta, tasks, loop, entities, bounded context/why, and seed status.

## Peer matrix

| Peer | Nav / IA | Graph | Tasks / loop | Backend transport | Takeaways for Trace |
|------|----------|-------|--------------|-------------------|---------------------|
| **Understand-Anything** `packages/dashboard` (+ skill `understand-dashboard`) | Single-pane dashboard: search, filters, persona, file explorer, side info; mobile bottom nav; URL token gate before UI | React + **@xyflow/react** (+ dagre/elk/d3-force/graphology); knowledge + domain graph views; path finder; staleness banner | No first-class task/kanban board — learn/onboarding panels only | Vite middleware serves static JSON graph files from `.ua` / `.understand-anything`; **host `127.0.0.1`**; one-shot `?token=` on data routes; skill launches viewer/`npx vite --host 127.0.0.1` and prints tokenized URL | Copy: loopback + token UX for local serve; graph viz with xyflow for Trace entity/plan graph. **Avoid:** shipping full `knowledge-graph.json` as default HTTP response (conflicts with Trace Law 6 progressive/bounded context). Agent skill pattern for “open dashboard” is useful for later docs/install, not Phase 29 MCP-over-internet. |
| **codebase-memory-mcp** `graph-ui` (+ `src/ui`, `scripts/security-ui.sh`, `scripts/embed-frontend.sh`) | Header tabs: **graph / stats (projects) / control**; query-string routing (`?tab=&project=`); graph tab disabled until project selected | React + **Three.js** (`@react-three/*`) 3D layout; node budget (`max_nodes`, default 5k) with progressive download; detail-by-center-node | Control tab = process/CPU/RAM/logs (daemon ops), not plan tasks | Dual: REST `GET /api/layout?…` (bounded) + JSON-RPC `POST /rpc` mirroring MCP `tools/call`; frontend **embedded** into Go binary (`embed-frontend.sh`); security audit enforces **127.0.0.1 bind**, no `0.0.0.0`, **no CORS `*`**, no external fetch domains in UI source | Copy: hard security gate scripts/checklist for S06; bounded graph API with explicit node budgets; optional embed-of-dist into `trace serve`. Prefer OpenAPI REST over exposing MCP-RPC to the browser (keeps MCP stdio-local). Avoid 3D/Three as MVP complexity unless S02 IA demands it. |
| **agentrq** `frontend` | Sidebar: Overview, workspaces, task filters, events; Vue Router nested routes; login gate; kanban under workspace | No code/knowledge graph — charts/sparklines/heatmap for analytics | First-class **tasks**: list, detail, form, kanban board, scheduled instances, permission verdicts, event stream | Versioned REST **`/api/v1/*`** + SSE event stream; auth (`/auth/*`); PWA; not loopback-only (multi-user product posture) | Copy: task/board IA and `/api/v1` versioning shape for Trace tasks/loop/reviews. Contrast: auth/PWA/push are cloud-product concerns — out of Phase 29 deploy; do not adopt always-on multi-user defaults. Vue stack works but adds a second UI ecosystem vs Trace’s React peer majority. |

**Peer count:** 3 locked peers surveyed (plus Understand skill + codebase-memory security/embed scripts). Stack lean **not overturned**.

## Trace surface map

Inventories re-checked 2026-08-21: CLI switch in `cmd/trace/root.go` matches prompt list; MCP catalog in `internal/mcp/mcp_test.go` lists 15 tools (matches prompt); **`serve` absent** (expected until S03).

| Surface | CLI | MCP | HTTP candidate family | GUI page | Priority (P0/P1/defer) |
|---------|-----|-----|----------------------|----------|------------------------|
| Help / version | `help`, `version` | `trace_version` | `/v1/health`, `/v1/version` | Status strip / About | P0 |
| Project bind / meta | `init` (setup) | — | `/v1/project` (root, store path, readiness) | Home / project header | P0 |
| Tasks board | `tasks` | `trace_tasks` | `/v1/tasks` | Tasks list + detail | P0 |
| Loop / next work | `loop` | `trace_loop` | `/v1/loop` | Loop board / next | P0 |
| Entities write | `add`, `link`, `transition` | `trace_add`, `trace_link`, `trace_transition` | `/v1/entities`, `/v1/links`, `/v1/transitions` | Entity detail + actions | P0 |
| Reviews | `review` | `trace_review` | `/v1/reviews` | Review pane | P1 (S02 may promote to P0 for MVP) |
| Progressive context | `why`, `context` | `trace_why`, `trace_context` | `/v1/why`, `/v1/context` (bounded packets) | Context / Why drawer | P0 |
| Search | `search` | `trace_search` | `/v1/search` | Search | P0 |
| Seed portable graph | `seed` | — | `/v1/seed` (export/import status; not full dump default) | Seed / sync status | P0 |
| Bounded graph viz | (via context/impact/plan tree) | — | `/v1/graph` (neighborhood / budgeted subgraph) | Graph view | P1 (MVP may stub; rich in S05) |
| Plan tree | `plan` | — (MCP intentionally omits `trace_plan`) | `/v1/plans` | Plan tree | P1 |
| Capability | `capability` | `trace_capability` | `/v1/capability` | Capability decisions | P1 |
| Impact | `impact` | `trace_impact` | `/v1/impact` | Impact report | P1 |
| Changes / regressions | `changes`, `regressions` | `trace_changes`, `trace_regressions` | `/v1/changes`, `/v1/regressions` | Changes / regressions | P1 |
| Agents recommend | `agents` | `trace_agents` | `/v1/agents` | Agents | P1 |
| Index / reindex | `index`, `reindex` | — | `/v1/index` (status/trigger) | Ops / index status | P1 / defer UI |
| Patterns / knowledge | `patterns`, `knowledge` | — | `/v1/patterns`, `/v1/knowledge` | Knowledge browse | defer |
| Tests / verify / eval | `test`, `tests`, `verify`, `eval`, `outcomes` | — | optional `/v1/verify` later | Verification | defer |
| Backup / restore / migrate | `backup`, `restore`, `migrate` | — | — (CLI-only ops) | — | defer |
| Auth / install | `auth`, `install` | — | local serve token only (`/v1/auth/token` or header); install stays CLI | Token gate (if non-loopback) | P1 for token; install defer |
| Serve (new) | **absent** → `serve` | — | HTTP server itself | SPA host | P0 (S03) |

**Parity intent:** HTTP/GUI cover MCP write/read cores and human task/loop/context flows; do **not** 1:1 mirror every CLI ops command or invent browser forks of library logic.

## API resource families (for S01 OpenAPI)

Recommended `/v1` families (S01 locks paths/schemas):

1. **`health`** — liveness; process up.
2. **`version`** — Trace version / API version.
3. **`project` / meta** — bound repo, `.trace/` readiness, feature flags safe for UI.
4. **`tasks`** — list, get, filter by status/board row intent.
5. **`loop`** — status, next, apply (library-backed; same semantics as CLI/MCP).
6. **`entities`** — create/get semantic entities (`add` parity).
7. **`links`** — create relations (`link` parity).
8. **`transitions`** — state transitions (`transition` parity).
9. **`reviews`** — review records / depth surfaces.
10. **`context`** — progressive context packets (Law 7); explicit expand params.
11. **`why`** — causal/explanation packets (bounded).
12. **`search`** — FTS / semantic search entry.
13. **`seed`** — export/import job status and summaries; **no** default full-graph body.
14. **`graph`** — budgeted subgraph / neighborhood for GUI (explicit `max_nodes` / center id).
15. **`plans`** — plan tree read (+ limited writes if library supports).
16. **`capability`** — tool decision rows / gates.
17. **`impact`** — impact reports.
18. **`changes`** — change listing.
19. **`regressions`** — regression listing.
20. **`agents`** — list / recommend.
21. **`index`** — index status / optional reindex trigger (ops).
22. **`auth` (local)** — optional bearer/token for non-loopback bind only; loopback may allow unauthenticated local use per S01 ADR.

**Out of first OpenAPI cut (document as deferred):** backup/restore/migrate, install, patterns/knowledge deep browse, eval/outcomes, push/SSE product features.

**Package/path recommendations for S01 (not locked here):** Go package `internal/httpapi`; OpenAPI artifact e.g. `docs/phases/phase-29-http-api-browser-gui/openapi.yaml` or `api/openapi.yaml`; SPA root `web/`.

## Stack options

| Option | Pros | Cons | Verdict |
|--------|------|------|---------|
| **TypeScript + Vite + React** | Matches 2/3 peers; xyflow ecosystem for graphs; strong Vite embed/serve story; team lean already locked | SPA complexity vs templates | **Recommended** |
| TypeScript + Vite + Vue | agentrq proves solid task/kanban IA; Pinia/router mature | Diverges from Understand + codebase-memory; weaker xyflow defaults; second mental model | Reject for Phase 29 (no overturning evidence) |
| Svelte + Vite | Small bundles | No peer in-tree; thinner graph ecosystem for agents | Reject |
| Go `html/template` only | Zero Node toolchain; simplest packaging | Poor interactive graph/board UX; fights production GUI goal | Reject as primary GUI |
| React + Three.js (as default viz) | Rich spatial graphs (codebase-memory) | Heavy deps; overkill for Trace plan/entity MVP | Optional P2 viz; not default stack |

**Recommended stack:** TypeScript + Vite + React SPA in `web/`, served by (or alongside) `trace serve`, talking OpenAPI JSON under `/v1`. Graph MVP: 2D library (e.g. `@xyflow/react`) with **budgeted** fetches — not full dump, not Three.js-first.

## Law carve-out draft

Paste-ready bullets for **AGENTS.md** / **project-rules** (apply in **S06**; draft only):

```markdown
### Phase 29 — Opt-in local HTTP + browser GUI (carve-out)

- **Allowed:** Opt-in `trace serve` on Trace core that exposes a versioned HTTP/JSON API and serves the browser GUI. Default bind **`127.0.0.1`**. Explicit flags required for non-loopback bind; prefer a one-time or configured token when bind is not loopback.
- **Law 19:** HTTP handlers and the browser UI are **adapters only**. They call the canonical Go library/API. No second source of truth, no business-logic fork in `web/`, no parallel SQLite from the browser.
- **Still forbidden:** Always-on network daemon; open bind (`0.0.0.0` / public internet) as default; pointing local product MCP (`trace-mcp`) at the public internet; full-graph dump as default API/GUI behavior (Laws 6–7).
- **Cloud path:** OpenAPI is the shared contract for a **future hosted product** (separate deploy/repo per `docs/TODO.md` Later developments). Phase 29 does not ship multi-tenant hosting, OAuth, billing, or tenancy.
- **Historical note:** FR-P28-X1 / older “no HTTP on core” language is superseded **only** for this opt-in local carve-out; it does not authorize silent daemons or public defaults.
```

Suggested `project-rules.md` settled-stack footnote (S06):

```markdown
| Surface (post–Phase 29) | Library + CLI + MCP (stdio) + **opt-in** `trace serve` HTTP/GUI (loopback default) |
```

(Keep the historical “Surface (P0)” row intact as P0-X history; add a new row rather than rewriting P0 law.)

## Risks / open decisions for S01–S02

1. **Auth on loopback:** Require token always (Understand style) vs loopback-trust + token only for non-loopback (lighter UX). S01 ADR must pick one.
2. **Embed vs separate static dir:** codebase-memory embeds `dist/` into the binary; Trace may start with `trace serve` serving `web/dist` from disk and defer embed to S06.
3. **Graph API shape:** Neighborhood-by-id vs typed plan-tree endpoint first; S02 IA should drive which ships in S04 MVP vs S05.
4. **MCP-RPC in browser:** Do **not** expose `/rpc` tools/call as the GUI transport; keep REST/OpenAPI. MCP remains local stdio.
5. **OpenAPI file location** and exact Go package name — recommend `internal/httpapi` + phase or `api/` OpenAPI; **S01 locks**.
6. **CORS:** Default deny / reflect localhost only; never `*` (peer security-ui).
7. **Write surface in MVP:** How many mutating routes land in S04 vs read-mostly + tasks/loop only — S02/S04 planners decide using this family list.
