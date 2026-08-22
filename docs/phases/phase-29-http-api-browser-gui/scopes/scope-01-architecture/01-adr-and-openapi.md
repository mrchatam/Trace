# P29-S01-01 — ADR + OpenAPI

## Metadata
- id: P29-S01-01
- todo_ids: [P29-S01-01]
- role: implementer
- skills: [documentation-and-adrs, api-and-interface-design]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Author **contract-only** artifacts for Phase 29 HTTP API + GUI architecture:

1. `docs/adr/ADR-HTTP-API-GUI.md`
2. `api/openapi.yaml` (OpenAPI 3.x)

**No** `trace serve` implementation, **no** `internal/httpapi` package, **no** `web/` scaffold in this row (those are S03/S04). Create parent dirs `docs/adr/` and `api/` as needed for the two files only.

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults**
- [../scope-00-research/RESEARCH.md](../scope-00-research/RESEARCH.md) — surface map + families + risks (honor; do not reopen settled locks)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 29 README](../../README.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Human locks and S01-00 locks are settled — do not re-grill bind/auth/CORS/MCP-RPC/paths.

## Locked defaults (copy into ADR “Decision” tables)

| Item | Value |
|------|-------|
| Cmd | Opt-in `trace serve` |
| Bind | `127.0.0.1` default; non-loopback requires `--allow-remote` |
| Auth | Loopback: no token required. Non-loopback: bearer required |
| Package (future) | `internal/httpapi` |
| Prefix | `/v1` |
| OpenAPI file | `api/openapi.yaml` |
| ADR file | `docs/adr/ADR-HTTP-API-GUI.md` |
| CORS | Deny by default; never `*` |
| GUI static | Disk `web/dist` first; embed → S06 |
| Transport | REST/OpenAPI only — **no** browser MCP `/rpc` |
| Graph | Budgeted `/v1/graph` (`max_nodes` + center); no full dump |
| Law 19 | Handlers + UI = adapters only → canonical Go library / `.trace/` SQLite |
| Cloud | Same contract; tenancy/OAuth/hosting **out**; document extension hooks only |

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-00-research/RESEARCH.md
test ! -d internal/httpapi
test ! -d web
! grep -q 'case "serve"' cmd/trace/root.go
```

If `serve` / `httpapi` / `web/` already exist unexpectedly, **stop** and mark `blocked` — do not invent a parallel contract.

## Role work

### 1. ADR (`docs/adr/ADR-HTTP-API-GUI.md`)

Use a standard ADR shape (Status Accepted, Date, Context, Decision, Consequences, Alternatives). Must cover:

| Section | Required content |
|---------|------------------|
| Context | Local-first GUI need; FR-P28-X1 carve-out; Law 19; S00 peer evidence (1–2 sentences each) |
| Decision — serve | Opt-in lifecycle; flags (`--addr` / `--allow-remote` / token flags as named); refuse open bind without flag |
| Decision — package | `internal/httpapi` thin handlers → library; CLI wires serve only |
| Decision — versioning | `/v1` URL prefix; API version vs Trace binary version fields |
| Decision — auth | Loopback-trust vs non-loopback bearer (S01-00 lock); how token is generated/printed (high-level) |
| Decision — CORS | Deny default; same-origin SPA; no `*`; Vite-dev exception if any |
| Decision — errors | One JSON error envelope (`code`, `message`, optional `details`); status mapping sketch |
| Decision — static GUI | Serve `web/dist`; embed deferred S06 |
| Decision — graph / progressive context | Budgeted graph; Laws 6–7; no full dump |
| Decision — MCP boundary | MCP stays stdio-local; GUI never uses `/rpc` |
| Decision — cloud hooks | Reserved OpenAPI extensions / headers; not implemented |
| Consequences | S03 implements server; S04+ GUI; S06 security/embed/law text |
| Alternatives rejected | Always-on daemon; CORS `*`; MCP-RPC browser transport; Vue stack; Three.js-first; full-graph default |

### 2. OpenAPI (`api/openapi.yaml`)

- OpenAPI **3.0 or 3.1**; `info.title` Trace HTTP API; `servers` example `http://127.0.0.1:PORT` (placeholder port OK).
- Tag by resource family; use `operationId`s stable for codegen later.
- Annotate operations with `x-trace-wave: p0 | p1 | defer` (and optionally `x-trace-gui: true`).
- Shared components: error schema, pagination/cursor if used, entity summary, task row, bounded graph node/edge, seed status.
- Security schemes: document optional `bearerAuth`; note in description that loopback may omit.

#### Must include (paths or path stubs with schemas) — P0 families

From RESEARCH surface map / API resource families:

| Family | Minimum operations |
|--------|-------------------|
| `health` | `GET /v1/health` |
| `version` | `GET /v1/version` |
| `project` | `GET /v1/project` (root, store readiness, safe flags) |
| `tasks` | list + get (+ filter query params) |
| `loop` | status, next, apply (and gate/reset if library exposes — align names to CLI/MCP semantics) |
| `entities` | create + get (add parity) |
| `links` | create (link parity) |
| `transitions` | create/apply (transition parity) |
| `context` | get bounded packet (expand params) |
| `why` | get bounded packet |
| `search` | query |
| `seed` | export/import **status** / summary — **not** default full-graph body |
| `graph` | budgeted subgraph (`max_nodes`, center id) |

#### Must appear as P1 or explicit deferred (do not omit silently)

| Family | Wave |
|--------|------|
| `reviews` | p1 (S02 may promote) |
| `plans` | p1 |
| `capability`, `impact`, `changes`, `regressions`, `agents` | p1 |
| `index` | p1 / defer UI |
| `auth` local token issue/rotate | p1 |
| backup/restore/migrate, install, patterns/knowledge deep, eval/outcomes, SSE/push | **deferred** section or `x-trace-wave: defer` |

#### Explicit non-goals in OpenAPI description or ADR

- No `POST /rpc` MCP mirror for the browser
- No default unbounded graph export
- No multi-tenant / OAuth routes implemented

### 3. Board Notes (this row only)

Record: ADR path, OpenAPI path, OpenAPI version, count of path items, any intentional schema TODOs that S03 must not invent silently.

## Exit criteria

- [ ] `docs/adr/ADR-HTTP-API-GUI.md` exists and covers serve/auth/CORS/Law19/cloud hooks/error model
- [ ] `api/openapi.yaml` exists; P0 families covered; P1/defer explicit
- [ ] No full-graph dump default; no CORS `*`; no MCP `/rpc` browser transport
- [ ] No product server/UI package created (`internal/httpapi`, `web/`, `serve` still absent)
- [ ] Board Notes on **P29-S01-01** with paths

## Minimal todos

- [ ] Preflight PASS
- [ ] Write ADR
- [ ] Write OpenAPI (P0 + deferred list)
- [ ] Self-check exit criteria
- [ ] Board status + Notes

## Todo updates

Status + notes on **P29-S01-01** only. Do not start review or S02.

## Next

**P29-S01-02**
