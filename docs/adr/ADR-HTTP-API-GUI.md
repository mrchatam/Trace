# ADR-HTTP-API-GUI: Opt-in local HTTP API + browser GUI

## Status

Accepted

## Date

2026-08-21

## Context

Trace needs a **production-ready browser GUI** for local project knowledge (tasks, loop, bounded context/graph) without abandoning local-first defaults. Phase 29 carves this out of the historical FR-P28-X1 “no HTTP on core” deferral: HTTP is **opt-in** on Trace core, library-backed, with an OpenAPI contract reusable by a future hosted product (separate deploy; not this phase’s target).

**Law 19** ([`G_PROJECT_LAWS.md`](../init/G_PROJECT_LAWS.md)): CLI, MCP, and UI must call the canonical library/API — adapters never fork business logic or create a second source of truth. HTTP handlers and the SPA are therefore thin adapters over the existing Go store/library and `.trace/` SQLite.

**S00 peer evidence** ([`RESEARCH.md`](../phases/phase-29-http-api-browser-gui/scopes/scope-00-research/RESEARCH.md)): Understand-Anything and codebase-memory confirm React+Vite local dashboards with **127.0.0.1** bind, no CORS `*`, and **budgeted** graph payloads; agentrq shows `/api/v1` task-board IA. Phase 29 rejects browser MCP `/rpc`, Vue as primary stack, Three.js-first viz, and full-graph dump defaults.

Human locks (2026-08-21) and S01-00 planner locks settle bind, auth, CORS, paths, and transport — this ADR records them; it does not reopen them.

## Decision

### Serve lifecycle

| Item | Decision |
|------|----------|
| Command | Opt-in **`trace serve`** only — never an always-on daemon or install-time auto-start |
| Default bind | **`127.0.0.1`** (loopback) |
| Address flag | `--addr` (host:port); default host is loopback |
| Remote bind | Non-loopback / `0.0.0.0` **refused** unless **`--allow-remote`** is set |
| Token flags | `--token` / `--token-file` (or generate-and-print on first non-loopback serve) — required when bind is non-loopback; optional on loopback |
| Lifecycle | Foreground process; Ctrl-C / SIGTERM stops; no systemd unit shipped by default |

Refuse open bind without `--allow-remote`. Document the refusal in CLI help and S06 security checklist.

### Package and layering (Law 19)

| Item | Decision |
|------|----------|
| Future package | `internal/httpapi` (implemented in S03; named here only) |
| Responsibility | Thin HTTP handlers: parse/validate request → call canonical library → map errors to the JSON envelope |
| CLI | `cmd/trace` wires `serve` only (flags, bind policy, static root); no domain logic in the command |
| UI | `web/` SPA calls `/v1` only; no parallel SQLite, no business-rule fork in the browser |
| SoT | `.trace/` SQLite via existing store/domain packages |

### Versioning

| Item | Decision |
|------|----------|
| URL prefix | **`/v1`** for all JSON API routes |
| API version | Field `api_version` (semver string for the HTTP contract, e.g. `1.0.0`) on `GET /v1/version` |
| Binary version | Field `trace_version` (Trace CLI/binary build version) on the same endpoint |
| Compatibility | Breaking HTTP contract changes require `/v2` (future); additive fields under `/v1` are preferred |

### Authentication

| Mode | Decision |
|------|----------|
| Loopback bind | **Loopback-trust:** no bearer token required for MVP |
| Non-loopback | **Bearer required:** `Authorization: Bearer <token>` |
| Token issuance | On non-loopback serve (or via future `POST /v1/auth/token`), generate a high-entropy token, print once to stdout / write to `--token-file`; operator pastes into clients. Optional loopback token is **not** required for MVP |
| Cloud | Reserved headers / OpenAPI `x-trace-cloud` extensions only — **not implemented** this phase |

### CORS

| Item | Decision |
|------|----------|
| Default | **Deny** cross-origin (no `Access-Control-Allow-Origin` for foreign origins) |
| Same-origin SPA | When the GUI is served by `trace serve` from the same origin, browser same-origin policy applies — no CORS wildcard needed |
| Vite-dev | Optional **localhost reflect** (exact origin match for Vite dev server) may be documented for local DX; **never** `Access-Control-Allow-Origin: *` |
| Ship rule | Production defaults never use `*` |

### Error model

One JSON error envelope for all `/v1` error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable summary",
    "details": {}
  }
}
```

| HTTP status | Typical `code` |
|-------------|----------------|
| 400 | `VALIDATION_ERROR`, `BAD_REQUEST` |
| 401 | `UNAUTHORIZED` (missing/invalid bearer on non-loopback) |
| 403 | `FORBIDDEN` (capability/gate deny mapped from library) |
| 404 | `NOT_FOUND` |
| 409 | `CONFLICT` (state machine / apply conflict) |
| 413 / 400 | `BUDGET_EXCEEDED` (graph/context over budget) |
| 500 | `INTERNAL_ERROR` |
| 501 | `NOT_IMPLEMENTED` (deferred routes if stubbed) |

Handlers must not leak stack traces or filesystem paths outside safe, intentional fields.

### Static GUI

| Item | Decision |
|------|----------|
| Phase 29 | Serve built assets from disk **`web/dist`** (first) |
| Dev | Vite dev server may proxy `/v1` to `trace serve` |
| Embed | Embedding `web/dist` into the binary is **deferred to S06** |
| SPA | TypeScript + Vite + React under `web/` (S00 stack lock) |

### Graph and progressive context (Laws 6–7)

| Item | Decision |
|------|----------|
| Graph | **`GET /v1/graph`** returns a **budgeted** neighborhood/subgraph; require explicit **`max_nodes`** and a **center** entity id (or equivalent seed) |
| Forbidden default | No full-graph dump as default response body (including seed export over HTTP) |
| Seed HTTP paths | `input_path` / `output_path` on seed import/export are **server-local** hints only; handlers **MUST** confine resolution to the bound project root (reject traversal / escape). Seed success bodies stay status/summary — never a default full-graph document |
| Context / why | Bounded packets with explicit expand/depth params — same progressive intent as CLI/MCP |
| Plans | Separate **`/v1/plans`** family (P1) for plan-tree reads; not a substitute for unbounded export |
| Loop parity | HTTP loop family includes **gate** and **reset** (CLI-exposed library surfaces) in addition to status/next/apply — needed for GUI loop console even though MCP `trace_loop` omits them |

### MCP boundary

| Item | Decision |
|------|----------|
| MCP | Remains **local stdio** (`trace-mcp`); not pointed at the public internet |
| Browser transport | **REST/OpenAPI only** |
| Forbidden | No `POST /rpc`, no MCP `tools/call` mirror for the browser |

### Cloud extension hooks (not implemented)

Document only — no Phase 29 runtime behavior:

- OpenAPI vendor extensions: `x-trace-cloud: true` on operations reserved for hosted product
- Reserved request headers: `X-Trace-Tenant`, `X-Trace-Request-Id` (ignored locally)
- Future OAuth / multi-tenant routes are **out of scope**; same `/v1` resource shapes are the intended shared contract

### Artifact locations

| Artifact | Path |
|----------|------|
| This ADR | `docs/adr/ADR-HTTP-API-GUI.md` |
| OpenAPI | `api/openapi.yaml` |

## Consequences

- **S03** implements `trace serve` + `internal/httpapi` against this ADR and `api/openapi.yaml`.
- **S04+** implement the GUI as an OpenAPI client; mutating route ship order is owned by S02/S04 (operations already tagged `x-trace-wave`).
- **S06** hardens security defaults, optional embed, and law/AGENTS carve-out text.
- Future hosted Trace reuses the OpenAPI contract; hosting/tenancy/OAuth remain a separate product.

## Alternatives rejected

| Alternative | Why rejected |
|-------------|--------------|
| Always-on network daemon | Violates local-first / human locks |
| CORS `Access-Control-Allow-Origin: *` | Peer security evidence; over-broad for local serve |
| Browser MCP `POST /rpc` / tools/call | Keeps MCP stdio-local; avoids dual transport |
| Vue as primary SPA stack | No overturning evidence vs React peer majority |
| Three.js-first graph | MVP complexity; prefer 2D budgeted viz later |
| Full-graph dump as default HTTP body | Conflicts with Laws 6–7 and RESEARCH |
| Token required on loopback for MVP | S01-00 chose loopback-trust for lighter local UX |
| Embed GUI in binary in S01/S03 | Deferred to S06; disk `web/dist` first |
