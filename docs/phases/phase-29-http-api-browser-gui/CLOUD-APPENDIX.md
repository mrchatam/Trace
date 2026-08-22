# Cloud appendix (design only) — Phase 29

**Status:** design notes for a future hosted product that reuses the **same OpenAPI** contract as local `trace serve`.  
**Phase 29 ships:** local-first opt-in HTTP + browser GUI only. **No** hosted deploy, multi-tenant SaaS, OAuth, billing, or TLS-termination product code in this phase.

## Reuse

- Contract SoT: [`api/openapi.yaml`](../../../api/openapi.yaml) + [`docs/adr/ADR-HTTP-API-GUI.md`](../../../adr/ADR-HTTP-API-GUI.md)
- Handlers remain Law-19 adapters over the canonical library; cloud is another deploy surface, not a second business-logic SoT

## Authn / authz (future)

| Concern | Direction |
|---------|-----------|
| Authn | Operator / user identity at the edge (OIDC/OAuth or similar) — **not** loopback-trust |
| Authz | Project/workspace scoped permissions; map to Trace project root or tenant DB |
| Tokens | Short-lived bearer or session cookies over TLS; never log secrets |
| Local carve-out | Loopback-trust + optional bearer remains valid for `trace serve` on a developer machine |

## Tenancy

- Reserve a request header (e.g. `X-Trace-Tenant` / workspace id) in product design — **not** required for local serve
- One logical Trace project graph per tenant workspace; no cross-tenant graph reads
- SQLite-per-workspace or equivalent isolation; do not share a single world-writable store

## TLS / edge

- Terminate TLS at reverse proxy / load balancer / platform edge
- Application continues to speak HTTP internally; HSTS and cert lifecycle are edge concerns
- Do not default Trace to open bind (`0.0.0.0`) without auth — same as local `--allow-remote` discipline

## CORS / browser

- Hosted GUI and API should be same-site or explicitly allowlisted exact origins (never `*`)
- CSP remains on static assets; tighten further for public internet (drop `'unsafe-inline'` when possible)

## Explicit non-goals (Phase 29)

- Multi-tenant control plane, billing meters, org admin UI
- Shipping always-on public daemons or changing the default bind away from loopback
- Exposing MCP tool `/rpc` as the browser transport
- Promoting HTTP seed `strict`/`task_id` without a shared library helper (CLI retains honesty export)

## Promotion path

When a hosted product is scheduled, open a new phase that: (1) locks tenancy + auth ADRs, (2) keeps OpenAPI compatibility, (3) does not weaken local-first defaults for `trace serve`.
