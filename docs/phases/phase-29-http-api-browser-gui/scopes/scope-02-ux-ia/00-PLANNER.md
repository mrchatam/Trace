# P29-S02-00 — Scope planner (UX IA)

## Metadata
- id: P29-S02-00
- todo_ids: [P29-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-design]
- verification: automated

## Objective

Produce locked `UX-IA.md`: screens, nav, empty states, production feature checklist mapped to OpenAPI ops. **No product UI code.**

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-02-ux-ia/UX-IA.md` |
| Inputs | S00 RESEARCH + S01 ADR (`docs/adr/ADR-HTTP-API-GUI.md`) + OpenAPI (`api/openapi.yaml`) |
| Progressive context | No full-graph dump screens (Law 6 / progressive retrieval) |
| Product chrome | Operator/agent tool — not marketing landing |
| Contract wave vs GUI ship | OpenAPI `x-trace-wave` = **API** contract wave. GUI ship columns: `S04` \| `S05` \| `defer` (locked in `01-ux-ia.md`) |
| Graph | API `/v1/graph` p0 (budgeted). Rich explorer = **GUI S05**; S04 stub only |
| Reviews | Stay **GUI S05** / API p1 — **not** promoted to S04 MVP |
| Loop console | Map **status / gate / next / apply / reset**. S04: status+gate **read** on Overview/Loop shell. S05: full interactive console |
| Discoveries & decisions | Entities (`discovery`/`decision`) + search + links/transitions promote; `GET /v1/capability` = S05 enrichment |
| Seed honesty | Status/summary + honesty warnings — never imply full-graph download as default HTTP body; path inputs project-root confined |
| Nav lean | Left nav (desktop): Overview, Tasks, Loop, Graph, Discoveries, Seed, Settings |

## Must-cover screens (minimum)

1. Project / open workspace
2. Overview dashboard (goals, active task, loop violations)
3. Graph explorer (bounded; expand-on-demand)
4. Tasks board + task detail (transitions, TRACE_TASK_ID hint)
5. Loop console (status / gate / next / apply / reset)
6. Discoveries & decisions (+ promote)
7. Seed export/import + honesty warnings
8. Settings (bind addr, token, theme)

## Exit criteria

- [x] `01-ux-ia.md` thickened enough for fresh subagent
- [x] `02-review.md` + `SCOPE-TODOS.md` aligned
- [x] Next **P29-S02-01**

## Next

P29-S02-01 → P29-S02-02
