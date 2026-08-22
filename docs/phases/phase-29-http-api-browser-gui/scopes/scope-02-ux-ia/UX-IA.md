# UX-IA — Phase 29 browser GUI

**Date:** 2026-08-21  
**Row:** P29-S02-01  
**Status:** Spec for S04/S05 implementers — docs only (no `web/` yet)

**Contract vs ship:** OpenAPI `x-trace-wave` (`p0` \| `p1` \| `defer`) is the **API** contract wave. Columns below use `api_wave` for that tag and `gui_ship` (`S04` \| `S05` \| `defer`) for when the **browser GUI** ships the surface. They are independent: an op can be API `p0` and still `gui_ship: S05` (e.g. rich graph, loop apply).

---

## 1. Product framing

Trace GUI is an **operator / agent tool** for a single bound project root served by opt-in `trace serve` — not a marketing landing, not a multi-tenant cloud console, not a second source of truth.

| Principle | Implication |
|-----------|-------------|
| **Law 19** | Every screen reads/writes only via `/v1` → library. No browser SQLite, no “edit `graph.json` in the SPA as SoT.” |
| **Laws 6–7** | Progressive, budgeted context. No full-graph dump screens; no primary CTA that downloads the entire graph over HTTP. |
| **Local-first** | Default loopback bind; token paste only for non-loopback (ADR). No cloud OAuth screens in Phase 29. |
| **Chrome** | Dense operator UI: left nav (desktop) / bottom tabs (narrow); project title + health/version strip always visible. |
| **Stack hint (S04)** | TypeScript + Vite + React under `web/` (ADR). IA does not invent Vue or Three.js-first viz. Rich 2D graph (xyflow-class) = **GUI S05**. |

Primary jobs: see what is active (goals/tasks/loop gate), move work safely (transitions / loop apply with confirm), inspect bounded neighborhood, record discoveries/decisions, and run seed export/import with honesty about what HTTP returns.

---

## 2. Information architecture (nav + hierarchy)

### Global chrome

| Element | Source | Notes |
|---------|--------|-------|
| Project title / root | `getProject` → `ProjectResponse.root` (+ `store_ready`) | Click → Project / open workspace |
| Health / version strip | `getHealth`, `getVersion` | Shows `ok`, `api_version`, `trace_version` |
| Optional search entry | `search` | S04: optional; may defer full chrome to S05 and use Tasks filters only |
| Auth indicator | Client-only + Settings | Loopback: “local trust”; non-loopback: token present/missing |

### Primary nav (locked)

Desktop: **left nav**. Narrow: **bottom tabs** (same ids).

| Nav id | Label | Default route | Primary screen |
|--------|-------|---------------|----------------|
| overview | Overview | `/` | Overview dashboard |
| tasks | Tasks | `/tasks` | Tasks board (+ `/tasks/:id` detail) |
| loop | Loop | `/loop` | Loop console |
| graph | Graph | `/graph` | Graph explorer |
| discoveries | Discoveries | `/discoveries` | Discoveries & decisions |
| seed | Seed | `/seed` | Seed export/import |
| settings | Settings | `/settings` | Settings |

**Project / open workspace** is not a nav item: it is the **gate / empty shell** when `store_ready` is false or `/v1/project` fails, and reachable from the project title in chrome.

### Hierarchy sketch

```text
[Chrome: project · health/version · search?]
├── Overview          goals summary · active task · loop gate/status strip
├── Tasks             list → detail (context/why drawers optional)
├── Loop              status/gate (S04) → next/apply/reset (S05)
├── Graph             stub (S04) → expand-on-demand canvas (S05)
├── Discoveries       list/detail (S04 read) → create/promote (S05)
├── Seed              status + honesty (S04) → export/import actions (S05)
└── Settings          theme · token · bind/version display
```

Reviews, plans, capability, impact, changes, regressions, agents: **no primary nav** in Phase 29 MVP; enter from task detail / S05 surfaces when OpenAPI `x-trace-wave` + this doc’s `gui_ship` allow (no separate FEATURE-MATRIX file yet).

---

## 3. Screen specs (×8 must-cover)

### 3.1 Project / open workspace

| Field | Content |
|-------|---------|
| **Purpose** | Confirm the bound project is ready before operator work, or block with a clear empty state. |
| **Primary user job** | Know which root is bound and whether `.trace/` is usable; recover when not. |
| **Key data** | `ProjectResponse` (`root`, `store_path`, `store_ready`, `flags`); `HealthResponse`; `VersionResponse`. |
| **Primary actions** | Retry health/project (client refresh → `getHealth` / `getProject`); copy root path (client-only); link to docs for `trace init` / `trace serve` (client-only). **No** multi-root picker. |
| **Empty state** | `store_ready: false` or missing `.trace/`: full-page blocked empty — “No Trace store at this root. Run init in the project, then restart serve.” Do not show fake boards. |
| **Error state** | `INTERNAL_ERROR` / unreachable: “Cannot reach API”; `UNAUTHORIZED` if token required and missing (deep-link Settings). |
| **`gui_ship`** | **S04** (health + readiness). Multi-root / polish → **defer**. |
| **Honesty / Law 6** | N/A beyond stating the server binds one project root (no browser-chosen second SoT). |

### 3.2 Overview dashboard

| Field | Content |
|-------|---------|
| **Purpose** | At-a-glance goals, active task, and loop gate/status violations without dumping the graph. |
| **Primary user job** | Answer “what should I care about right now?” and jump to Tasks or Loop. |
| **Key data** | Goals via `search` and/or `listTasks` + entity kinds; active task from loop/`TaskRow`; `LoopStatusResponse`; `LoopGateResponse` (violations strip). Optional later: `getContext` for active task snippet. |
| **Primary actions** | Open active task → Tasks detail; “Open Loop” → Loop; refresh strip → `getLoopStatus` / `getLoopGate`. |
| **Empty state** | No goals/tasks: “No goals yet — add via CLI/MCP or Discoveries (S05).” Gate clear: quiet “Gate OK” strip. |
| **Error state** | Gate/status `FORBIDDEN` / `INTERNAL_ERROR`: show envelope `message` + code; do not invent pass/fail. |
| **`gui_ship`** | **S04**: goals summary, active task, **loop gate/status violations strip**. Deeper widgets (impact/changes) → **S05** / optional. |
| **Honesty / Law 6** | Brief: overview is summaries only — not a full plan-tree or graph dump. |

### 3.3 Graph explorer

| Field | Content |
|-------|---------|
| **Purpose** | Inspect a **budgeted** neighborhood around a chosen center; expand on demand. |
| **Primary user job** | Understand local relations without loading the whole graph. |
| **Key data** | `BoundedGraph` (`center`, `max_nodes`, `nodes`/`GraphNode`, `edges`/`GraphEdge`, `truncated`); center pick via search/`EntitySummary` / task id. |
| **Primary actions** | **S04:** pick center + optional single budgeted `getGraph` **or** entity list + “Open in graph (S05)” placeholder — never unbounded viz. **S05:** expand neighbor (new center / raise budget within caps), open entity detail. |
| **Empty state** | No center selected: prompt to pick entity/task. Empty neighborhood: “No edges in budget — try another center.” |
| **Error state** | `BUDGET_EXCEEDED` / `VALIDATION_ERROR` (missing `center`/`max_nodes`): explain required params; `NOT_FOUND` for bad center. |
| **`gui_ship`** | **Split:** stub **S04**; rich xyflow-class expand-on-demand **S05**. API remains `api_wave: p0`. |
| **Honesty / Law 6** | **Required.** Banner when `truncated: true`. Copy: “Budgeted neighborhood — not the full project graph.” Forbid “Download entire graph” primary CTA. |

### 3.4 Tasks board + task detail

| Field | Content |
|-------|---------|
| **Purpose** | Browse tasks and inspect one task with safe transition affordances. |
| **Primary user job** | Find work, copy `TRACE_TASK_ID`, understand state, optionally transition. |
| **Key data** | `TaskListResponse` / `TaskRow`; detail `TaskRow` + optional `ContextPacket` / `WhyPacket`; `TransitionResult` on write. |
| **Primary actions** | List/filter → `listTasks`; open → `getTask`; copy `TRACE_TASK_ID=<uuid>` (client-only); **S04 light:** may call `createTransition` for ordinary moves the API accepts — on `FORBIDDEN`/`CONFLICT`/`VALIDATION_ERROR` show the envelope only; **do not** encode DONE/review/cap policy in the SPA (Law 19). **S05:** fuller transition UX + enforce/gate awareness (`getLoopGate`, `TransitionResult.warning`). Optional drawers: `getContext`, `getWhy`. |
| **Empty state** | No tasks: “No tasks in store.” Filtered empty: clear filters CTA. |
| **Error state** | Transition `CONFLICT` / `FORBIDDEN` / `VALIDATION_ERROR`: surface envelope; DONE failures cite review/caps honestly (actor ≠ auth). |
| **`gui_ship`** | **Split:** list + detail read + `TRACE_TASK_ID` + light transition **S04**; full transitions / gate-aware DONE **S05**. |
| **Honesty / Law 6** | Context/why are bounded packets — show depth params; never imply complete closure of all related entities. |

### 3.5 Loop console

| Field | Content |
|-------|---------|
| **Purpose** | Operate the deliberation loop with full CLI-parity five-op surface. |
| **Primary user job** | See status/gate; in S05, fetch next, apply, or reset with confirm. |
| **Key data** | `LoopStatusResponse`; `LoopGateResponse`; `LoopNextResponse`; `LoopApplyEnvelope` / `LoopApplyResult`; `LoopResetResponse`. |
| **Primary actions** | All five ops mapped: `getLoopStatus`, `getLoopGate`, `getLoopNext`, `postLoopApply`, `postLoopReset`. **S04:** read-only status + gate summary (also mirrored on Overview); Loop page may be a shell that deep-links Overview strip. **S05:** interactive next / apply / reset + gate detail; apply/reset require keyboard-confirm destructive dialog. |
| **Empty state** | No active seed/next: “No next packet — gate may be blocking or queue empty” (use gate payload when present). |
| **Error state** | Apply `CONFLICT` / `FORBIDDEN`; reset failures via envelope; never hide gate deny. |
| **`gui_ship`** | **Split:** **S04 read** (status + gate); **S05 write** (next / apply / reset + full console). |
| **Honesty / Law 6** | N/A beyond not inventing loop state client-side (Law 19). |

### 3.6 Discoveries & decisions

| Field | Content |
|-------|---------|
| **Purpose** | Entity-centric list/detail for `discovery` and `decision`; promote into tasks/links. |
| **Primary user job** | Find or record a discovery/decision and connect it to work. |
| **Key data** | `EntitySummary` / `CreateEntityRequest` (`kind: discovery\|decision`); `SearchResponse`; `LinkSummary`; `TaskRow` when promoting. |
| **Primary actions** | Find → `search` (kind filter); detail → `getEntity`; **S05:** create → `createEntity`; relate → `createLink`; promote → `createTransition` and/or create linked `task` via `createEntity` + `createLink`. Capability catalog `listCapability` = **S05 enrichment only** (not a second SoT). |
| **Empty state** | “No discoveries or decisions yet.” |
| **Error state** | Create/link `VALIDATION_ERROR` / `CONFLICT`; search empty ≠ error. |
| **`gui_ship`** | **Split:** list/search read + detail (if cheap) **S04**; create + promote **S05**. |
| **Honesty / Law 6** | Brief: lists are search/filtered, not full entity dump. |

### 3.7 Seed export/import + honesty warnings

| Field | Content |
|-------|---------|
| **Purpose** | Show seed readiness and (in S05) run export/import jobs with path confinement honesty. |
| **Primary user job** | Know last seed status; export/import without mistaking HTTP for a full-graph download API. |
| **Key data** | `SeedStatus` (`ready`, `last_export_at`, `last_import_at`, `notes`); `SeedJobStatus` (`status`, `summary`, `path`, counts, `error`); requests `SeedExportRequest` / `SeedImportRequest` (`output_path` / `input_path`). |
| **Primary actions** | Refresh → `getSeedStatus`; **S05:** export → `postSeedExport`; import → `postSeedImport`. Paths are project-root confined (ADR). |
| **Empty state** | Never exported: “No export recorded — status only until you run export (S05) or CLI.” |
| **Error state** | Path escape / traversal → `VALIDATION_ERROR` / `FORBIDDEN` with clear “path must stay under project root”; job `failed` shows `SeedJobStatus.error`. |
| **`gui_ship`** | **Split:** status + honesty copy **S04**; export/import actions + error honesty **S05**. |
| **Honesty / Law 6** | **Required.** Persistent warning: “HTTP seed responses are **status/summary only** — not a full-graph document. Portable graph remains `trace seed` / `trace/graph.json` workflow; paths are server-local under the bound root.” Never offer “download entire graph as JSON body” as primary UX. |

### 3.8 Settings

| Field | Content |
|-------|---------|
| **Purpose** | Client chrome for theme, bearer token, and display of serve bind/version/project. |
| **Primary user job** | Connect securely when non-loopback; personalize theme; verify what server they hit. |
| **Key data** | Display: `ProjectResponse`, `VersionResponse`, bind addr from client config / serve URL (client-only). Token: local storage/session (client-only). Theme: local (client-only). |
| **Primary actions** | Paste/clear bearer (client-only; sent as `Authorization` on non-loopback); theme toggle (client-only); **S05 optional:** `postAuthToken` if exposed (`x-trace-gui: false` today — treat as API-only / future). No cloud OAuth. |
| **Empty state** | Loopback: “Token not required (loopback-trust).” Non-loopback without token: blocking callout → paste token. |
| **Error state** | `UNAUTHORIZED` on any `/v1` call: banner + focus token field. |
| **`gui_ship`** | **S04:** theme + token paste + show bind/version/project. Token generate via API → **S05** if productized; no cloud auth ever in this phase. |
| **Honesty / Law 6** | N/A. Law 19: do not invent browser SQLite or parallel project roots. |

---

## 4. Empty / error / honesty patterns

### Empty

| Pattern | When | UI |
|---------|------|-----|
| **Blocked project** | `store_ready: false` | Full-page empty; hide nav destinations that need store |
| **Quiet empty list** | Zero tasks / discoveries | One sentence + optional CLI hint |
| **Filtered empty** | Search/filter no hits | “No matches” + clear filters |
| **Gate clear** | Loop gate OK | Compact positive strip, not a celebration modal |

### Error (ADR envelope)

Map `error.code` from `ErrorEnvelope`:

| Code | Operator-facing default |
|------|-------------------------|
| `VALIDATION_ERROR` / `BAD_REQUEST` | Fix inputs (paths, required fields, graph budget params) |
| `UNAUTHORIZED` | Paste token in Settings (non-loopback) |
| `FORBIDDEN` | Gate/capability deny — show message; do not fake success |
| `NOT_FOUND` | Entity/task missing — offer search |
| `CONFLICT` | State machine / apply conflict — refresh status then retry |
| `BUDGET_EXCEEDED` | Lower `max_nodes` or pick tighter center (Graph/Context) |
| `INTERNAL_ERROR` | Retry; do not leak stacks |
| `NOT_IMPLEMENTED` | Hide or disable deferred ops in GUI |

### Honesty (cross-cutting)

1. **Graph:** budget + `truncated` banner; no full-dump CTA.  
2. **Seed:** status/summary only over HTTP; path confinement; CLI/`trace/graph.json` remains portable SoT path.  
3. **Actor / token:** `actor` on transitions is provenance, not auth; bearer is serve access, not operator identity.  
4. **Law 19:** no “edit graph JSON in browser as SoT” pattern on any screen.

---

## 5. GUI ship matrix (S04 vs S05 vs defer)

| Area | S04 MVP (GUI-P0) | S05 rich (GUI-P1) | defer |
|------|------------------|-------------------|-------|
| Project / open | Health + project readiness; blocked empty if `.trace/` missing | Polish | Multi-root |
| Overview | Goals summary, active task, **loop gate/status** strip | Deeper widgets (impact/changes) optional | — |
| Graph | Stub: center picker + one budgeted fetch **or** list + “open in graph (S05)” | Expand-on-demand 2D explorer (budgeted) | Three.js / unbounded viz |
| Tasks | List + detail read; `TRACE_TASK_ID`; light transition if safe | Full transitions + enforce/gate awareness | — |
| Loop | Read-only status/gate on Overview and/or Loop shell | Full console: next / apply / reset + gate detail | — |
| Discoveries | List/search read + detail (if cheap) | Create discovery/decision + promote-to-task | — |
| Seed | Status + honesty copy | Export/import actions + path errors | Full-graph body download |
| Settings | Theme + token paste + bind/version/project | Token generate via API if exposed | Cloud OAuth |
| Reviews | — | Review list/detail (API `p1`) | — |
| Plans / capability / impact / agents / changes / regressions | — | As FEATURE-MATRIX allows | — |
| Deferred API (`x-trace-wave: defer`) | No GUI | No GUI unless later promote | backup/restore/migrate/install/patterns/knowledge/eval/events SSE |

**Locks held:** rich graph = **S05**; reviews **not** in S04; loop **S04 read / S05 write**; seed honesty from S04; Law 19 everywhere.

---

## 6. OpenAPI operation map

| operationId | path | api_wave | gui_ship | screen(s) |
|-------------|------|----------|----------|-----------|
| getHealth | `GET /v1/health` | p0 | S04 | Project, Settings, chrome |
| getVersion | `GET /v1/version` | p0 | S04 | Project, Settings, chrome |
| getProject | `GET /v1/project` | p0 | S04 | Project, chrome |
| listTasks | `GET /v1/tasks` | p0 | S04 | Tasks, Overview |
| getTask | `GET /v1/tasks/{task_id}` | p0 | S04 | Tasks detail |
| getLoopStatus | `GET /v1/loop/status` | p0 | S04 | Overview, Loop |
| getLoopGate | `GET /v1/loop/gate` | p0 | S04 | Overview, Loop |
| getLoopNext | `GET /v1/loop/next` | p0 | S05 | Loop |
| postLoopApply | `POST /v1/loop/apply` | p0 | S05 | Loop |
| postLoopReset | `POST /v1/loop/reset` | p0 | S05 | Loop |
| createEntity | `POST /v1/entities` | p0 | S05 | Discoveries (create/promote); Tasks (create task) |
| getEntity | `GET /v1/entities/{entity_id}` | p0 | S04 | Discoveries detail; Graph center |
| createLink | `POST /v1/links` | p0 | S05 | Discoveries promote/relate |
| createTransition | `POST /v1/transitions` | p0 | S04 light / S05 full | Tasks; Discoveries promote |
| getContext | `GET /v1/context` | p0 | S04 optional / S05 | Tasks detail drawer |
| getWhy | `GET /v1/why` | p0 | S04 optional / S05 | Tasks / entity drawer |
| search | `GET /v1/search` | p0 | S04 (filters) / S05 (chrome) | Discoveries, Graph picker, chrome |
| getSeedStatus | `GET /v1/seed/status` | p0 | S04 | Seed |
| postSeedExport | `POST /v1/seed/export` | p0 | S05 | Seed |
| postSeedImport | `POST /v1/seed/import` | p0 | S05 | Seed |
| getGraph | `GET /v1/graph` | p0 | S04 stub / S05 rich | Graph |
| listReviews | `GET /v1/reviews` | p1 | S05 | Reviews (no S04 nav) |
| createReview | `POST /v1/reviews` | p1 | S05 | Reviews / Tasks |
| getReview | `GET /v1/reviews/{review_id}` | p1 | S05 | Reviews |
| listPlans | `GET /v1/plans` | p1 | S05+ | no must-cover screen / S05+ |
| listCapability | `GET /v1/capability` | p1 | S05 | Discoveries enrichment; Tasks |
| getImpact | `GET /v1/impact` | p1 | S05+ | Overview widget optional / S05+ |
| listChanges | `GET /v1/changes` | p1 | S05+ | no screen / S05+ |
| listRegressions | `GET /v1/regressions` | p1 | S05+ | no screen / S05+ |
| listAgents | `GET /v1/agents` | p1 | S05+ | no screen / S05+ |
| getIndexStatus | `GET /v1/index` | p1 | defer | no GUI (`x-trace-gui: false`) |
| postIndexReindex | `POST /v1/index` | defer | defer | no GUI |
| postAuthToken | `POST /v1/auth/token` | p1 | S05 optional | Settings (API `x-trace-gui: false` today) |
| postBackup | `POST /v1/backup` | defer | defer | no GUI |
| postRestore | `POST /v1/restore` | defer | defer | no GUI |
| postMigrate | `POST /v1/migrate` | defer | defer | no GUI |
| postInstall | `POST /v1/install` | defer | defer | no GUI |
| listPatterns | `GET /v1/patterns` | defer | defer | no GUI |
| listKnowledge | `GET /v1/knowledge` | defer | defer | no GUI |
| getEval | `GET /v1/eval` | defer | defer | no GUI |
| getEventsSSE | `GET /v1/events` | defer | defer | no GUI |

### Feature checklist (must-cover → ops)

| Feature | S04 | S05 | Primary operationIds |
|---------|-----|-----|----------------------|
| Project readiness | ✓ | polish | getHealth, getVersion, getProject |
| Overview + gate strip | ✓ | widgets | listTasks, search, getLoopStatus, getLoopGate |
| Graph stub / rich | stub | rich | getGraph (+ search) |
| Tasks + TRACE_TASK_ID | ✓ | full transition | listTasks, getTask, createTransition |
| Loop five-ops | status+gate | next+apply+reset | getLoopStatus, getLoopGate, getLoopNext, postLoopApply, postLoopReset |
| Discoveries read/write | read | create/promote | search, getEntity, createEntity, createLink, createTransition |
| Seed honesty / actions | status | export/import | getSeedStatus, postSeedExport, postSeedImport |
| Settings chrome | ✓ | token API | client-only + getVersion/getProject; postAuthToken optional |
| Reviews | — | ✓ | listReviews, getReview, createReview |

---

## 7. Accessibility / keyboard (operator UI)

Documented floor for S04+ (do not build in this row):

- Primary nav items are **focusable** with visible focus rings (not color-only).
- Icon-only controls (refresh, copy task id, theme) have **accessible names**.
- **Tasks detail** and **Loop apply/reset**: full keyboard path; destructive apply/reset use an explicit confirm step (Enter on confirm, Escape dismiss).
- Prefer `prefers-reduced-motion: reduce` — no essential info only in motion; graph layout animation optional and disableable.
- Status/gate strips expose text, not icon-only severity.
- Skip link to main content recommended for left-nav layout.

---

## 8. Out of scope / non-goals

- Product UI implementation under `web/` (S04/S05).
- `internal/httpapi` handlers or `trace serve` (S03).
- Rewriting ADR or OpenAPI in this row.
- Marketing landing, PWA install prompts, cloud OAuth, multi-tenant tenancy UI.
- Browser MCP `/rpc` transport; Vue/Three-first stack.
- Full-graph dump UX; editing portable JSON as browser SoT.
- Always-on daemon; open bind without `--allow-remote`.
- Promoting reviews or rich graph into **S04** MVP.
- Multi-root project switcher (defer).

---

## References

- [00-PLANNER.md](00-PLANNER.md) — S02 locks  
- [01-ux-ia.md](01-ux-ia.md) — this implementer prompt  
- [../scope-00-research/RESEARCH.md](../scope-00-research/RESEARCH.md)  
- [docs/adr/ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)  
- [api/openapi.yaml](../../../../../api/openapi.yaml)  
