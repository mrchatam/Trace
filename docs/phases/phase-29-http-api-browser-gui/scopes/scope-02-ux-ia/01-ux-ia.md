# P29-S02-01 — UX IA implementer

## Metadata
- id: P29-S02-01
- todo_ids: [P29-S02-01]
- role: implementer
- skills: [frontend-design, planning-and-task-breakdown]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Author **docs-only** `UX-IA.md` at:

`docs/phases/phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md`

Cover: primary nav + IA, all **eight must-cover screens**, empty/error/honesty states, **GUI ship wave** (S04 MVP vs S05 rich) distinct from OpenAPI `x-trace-wave`, and a production feature checklist mapped to OpenAPI `operationId`s / paths.

**No** product UI code (`web/`), **no** `internal/httpapi`, **no** OpenAPI/ADR rewrites.

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults** (honor; do not reopen)
- [../scope-00-research/RESEARCH.md](../scope-00-research/RESEARCH.md) — peer IA + surface map
- [docs/adr/ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../api/openapi.yaml)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §§6–7, §19
- Phase README + S04/S05 planner stubs (ship targets; S02 owns the IA split)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). S02-00 locks below are settled — do not re-grill stack, bind, CORS, or Law 19.

## Locked defaults (S02-00 — do not reopen)

| Item | Value |
|------|-------|
| Output | `scopes/scope-02-ux-ia/UX-IA.md` only |
| Chrome | Operator / agent tool UI — not marketing landing |
| Stack hint (for later S04) | TS + Vite + React under `web/` (ADR); IA must not invent Vue/Three-first |
| Progressive context | No full-graph dump screens or “download entire graph” primary CTA (Laws 6–7) |
| Contract vs GUI ship | OpenAPI `x-trace-wave` = **API** contract wave. Columns in UX-IA: `api_wave` + `gui_ship` (`S04` \| `S05` \| `defer`) |
| Graph | `/v1/graph` is API **p0** (budgeted). **Rich** explorer (xyflow-class, expand-on-demand) = **GUI S05**. S04 may ship a **stub**: entity picker + optional single budgeted fetch / “open in graph (S05)” placeholder — never unbounded viz |
| Reviews | Keep **GUI S05** / API p1 — **do not promote** reviews into S04 MVP |
| Loop console | Spec full console against `GET /v1/loop/status`, `GET /v1/loop/gate`, `GET /v1/loop/next`, `POST /v1/loop/apply`, `POST /v1/loop/reset`. **S04:** status + gate summary on Overview (read). **S05:** full interactive console (next / apply / reset + gate detail) |
| Discoveries & decisions | Entity-centric: types `discovery` \| `decision` via `POST /v1/entities`, detail `GET /v1/entities/{id}`, find via `GET /v1/search`, relate via `POST /v1/links`, promote via `POST /v1/transitions` and/or create linked `task`. Capability catalog `GET /v1/capability` = **S05** enrichment — not a second SoT |
| Seed honesty | Screens use `GET /v1/seed/status` + export/import **job status/summary** only. Warn that HTTP bodies are not full-graph downloads; path fields are project-root confined (ADR). **S04:** status + warnings. **S05:** export/import actions + error honesty |
| Settings | Bind addr / token / theme are **client chrome** over serve config: display health/version/project; token paste for non-loopback; theme local. Do not invent browser SQLite or parallel project roots |
| Auth UX | Loopback-trust (no token gate). Non-loopback: bearer paste in Settings (ADR). No cloud OAuth screens |
| Law 19 | Every screen’s data/actions map to `/v1` ops → library. Forbid “edit graph JSON in browser as SoT” patterns |

### Must-cover screens (all required in UX-IA)

1. Project / open workspace  
2. Overview dashboard (goals, active task, loop violations)  
3. Graph explorer (bounded; expand-on-demand)  
4. Tasks board + task detail (transitions, `TRACE_TASK_ID` hint)  
5. Loop console (status / gate / next / apply / reset)  
6. Discoveries & decisions (+ promote)  
7. Seed export/import + honesty warnings  
8. Settings (bind addr, token, theme)

### Recommended primary nav (lock unless evidence forces change)

Peer lean (agentrq sidebar + UA single-pane): **left nav** (desktop) / **bottom tabs** (narrow):

| Nav id | Label | Default route sketch |
|--------|-------|----------------------|
| overview | Overview | `/` |
| tasks | Tasks | `/tasks` |
| loop | Loop | `/loop` |
| graph | Graph | `/graph` |
| discoveries | Discoveries | `/discoveries` |
| seed | Seed | `/seed` |
| settings | Settings | `/settings` |

Global chrome: project title from `GET /v1/project`, health/version strip, optional search entry (`GET /v1/search`) — search results may deep-link into tasks/entities; full search chrome can be S05 if MVP uses tasks filters only.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-00-research/RESEARCH.md
test -f docs/adr/ADR-HTTP-API-GUI.md
test -f api/openapi.yaml
test ! -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md
# optional sanity: loop gate/reset present
grep -q '/v1/loop/gate' api/openapi.yaml && grep -q '/v1/loop/reset' api/openapi.yaml
```

If `UX-IA.md` already exists from a partial run, **update in place** to meet exit criteria (do not fork a second IA file).

## Role work

### 1. Write `UX-IA.md` with these sections (required)

```markdown
# UX-IA — Phase 29 browser GUI

## 1. Product framing
## 2. Information architecture (nav + hierarchy)
## 3. Screen specs (×8 must-cover)
## 4. Empty / error / honesty patterns
## 5. GUI ship matrix (S04 vs S05 vs defer)
## 6. OpenAPI operation map
## 7. Accessibility / keyboard (operator UI)
## 8. Out of scope / non-goals
```

### 2. Per-screen spec minimum

For each of the eight screens, include:

| Field | Content |
|-------|---------|
| Purpose | One sentence |
| Primary user job | What the operator finishes here |
| Key data | Fields / entities (from OpenAPI schemas by name) |
| Primary actions | Buttons/flows → `operationId` or “client-only” |
| Empty state | What to show when no data |
| Error state | Map to ADR error envelope codes where relevant |
| `gui_ship` | `S04` \| `S05` \| split (read S04 / write S05) |
| Honesty / Law 6 | Required for Graph + Seed; others N/A or brief |

### 3. GUI ship matrix (locked lean — refine wording, do not invert)

| Area | S04 MVP (GUI-P0) | S05 rich (GUI-P1) |
|------|------------------|-------------------|
| Project / open | Health + project readiness; blocked empty if `.trace/` missing | Polish / multi-root **defer** (single bound root) |
| Overview | Goals summary, active task, **loop gate/status** violations strip | Deeper widgets (impact/changes) optional |
| Graph | Stub: pick center + budgeted fetch **or** entity list + defer rich canvas | Expand-on-demand explorer (2D; budgeted) |
| Tasks | List + detail read; show `TRACE_TASK_ID`; light transition if safe | Full transitions + enforce/gate awareness |
| Loop | Read-only status/gate on Overview and/or Loop page shell | Full console: next / apply / reset + gate |
| Discoveries | List/search read + detail (if cheap) | Create discovery/decision + promote-to-task |
| Seed | Status + honesty copy | Export/import actions + path confinement errors |
| Settings | Theme + token paste + show bind/version | Token generate via API if exposed; no cloud auth |
| Reviews | — | Review list/detail (API p1) |
| Plans / capability / impact / agents | — | As FEATURE-MATRIX allows |
| Deferred API (`x-trace-wave: defer`) | No GUI | No GUI unless later promote |

### 4. OpenAPI map table

Produce a table: `operationId` | path | `api_wave` | `gui_ship` | screen(s). Cover at least all `x-trace-gui: true` ops used by must-cover screens; mark unused p1/defer as “no screen / S05+”.

### 5. A11y floor (document, don’t build)

- Focusable primary nav; visible focus
- Labels on icon-only controls
- Prefer keyboard paths for Tasks detail + Loop apply (destructive confirm)
- Respect reduced-motion (note for S04+)

## Exit criteria

- [ ] `UX-IA.md` exists at the locked path with all required sections
- [ ] All eight must-cover screens specified
- [ ] `api_wave` vs `gui_ship` explicit; graph rich = S05; reviews not in S04
- [ ] Loop maps all five ops (status/gate/next/apply/reset) with S04 read / S05 write split
- [ ] Seed honesty + no full-graph dump UX stated
- [ ] OpenAPI operation map present
- [ ] No product UI code; no ADR/OpenAPI edits

## Todo updates

Status + notes on **P29-S02-01** only.

## Next

**P29-S02-02**
