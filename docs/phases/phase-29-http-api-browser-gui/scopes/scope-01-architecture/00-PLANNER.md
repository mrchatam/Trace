# P29-S01-00 — Scope planner (architecture + OpenAPI)

## Metadata
- id: P29-S01-00
- todo_ids: [P29-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, documentation-and-adrs]
- mcps: []
- verification: automated

## Objective

Lock ADR + OpenAPI v1 for Trace HTTP API from S00 `RESEARCH.md`. Thicken `01-adr-and-openapi.md` + `02-review.md`. **No product server code** beyond contract artifacts.

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults (final — S01-00, 2026-08-21)

Honors RESEARCH Risks / open decisions §1–7. Implementer must not re-debate these.

| Topic | Locked value | RESEARCH § |
|-------|--------------|------------|
| Cmd | `trace serve` (opt-in only) | human lock |
| Bind | Default `127.0.0.1`; refuse `0.0.0.0` / non-loopback without `--allow-remote` | human lock |
| Package | `internal/httpapi` (S03 implements; S01 names only) | §5 |
| API prefix | `/v1` | human lock |
| SoT | SQLite under `.trace/` via existing store/library — handlers are adapters (Law 19) | Law 19 |
| Auth local | **Loopback-trust:** no token required when bound to loopback. **Non-loopback:** require bearer token (`Authorization: Bearer …` or equivalent); document issuance in ADR. Optional loopback token is **not** required for MVP. | §1 |
| Auth cloud | OpenAPI `x-trace-cloud` extensions / reserved headers (`X-Trace-Tenant` placeholder) — **not implemented** this phase | human lock |
| Static GUI | Phase 29: serve `web/dist` from disk (or Vite dev proxy). **Embed into binary deferred to S06.** | §2 |
| SPA root | `web/` (TS + Vite + React) | S00 stack |
| ADR path | `docs/adr/ADR-HTTP-API-GUI.md` | §5 |
| OpenAPI path | `api/openapi.yaml` | §5 |
| CORS | **Default deny.** Same-origin when SPA is served by `trace serve`. No `Access-Control-Allow-Origin: *`. Optional localhost reflect for Vite-dev only — document in ADR; never ship `*` | §6 |
| Browser transport | **REST/OpenAPI only.** Do **not** expose MCP `POST /rpc` / tools/call to the browser | §4 |
| Graph API | Primary: budgeted neighborhood `GET /v1/graph` with required/explicit `max_nodes` (+ center id). Separate `GET /v1/plans` for plan tree. No full-graph dump. S02 picks which is MVP UI; both in OpenAPI | §3 |
| Write surface (contract) | OpenAPI **includes** P0 mutating families (tasks/loop apply, entities, links, transitions, seed status ops). Which mutate routes ship in S04 vs later = **S02/S04** — mark `x-trace-wave: p0\|p1\|defer` on operations | §7 |
| Inputs | `../scope-00-research/RESEARCH.md` required | — |

## Deliverables (authored by P29-S01-01)

- `docs/adr/ADR-HTTP-API-GUI.md`
- `api/openapi.yaml`

## Exit criteria

- [x] Implementer prompt runnable with paths + exit criteria
- [x] Board Notes; next **P29-S01-01**

## Next

P29-S01-01 → P29-S01-02
