# P29-S06-00 — Scope planner (production hardening)

## Metadata
- id: P29-S06-00
- todo_ids: [P29-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, security-and-hardening, shipping-and-launch]
- verification: automated

## Objective

Lock production hardening prompts: security defaults, packaging, docs, AGENTS carve-out, cloud appendix (design only).

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults (final for S06-01)

| Item | Value |
|------|-------|
| Deploy target | Local production use — **not** multi-tenant SaaS |
| Live baseline | `internal/httpapi` + `web/` SPA + `trace serve` (S03–S05 shipped) |
| Bind / auth | Keep ADR: default `127.0.0.1:7432`; non-loopback needs `--allow-remote` + bearer; loopback-trust (token optional; enforce if set) |
| CORS | **Default deny** / never `*`. Optional **exact-origin** reflect only via `--cors-origin URL` (Vite DX). No wildcard, no substring match |
| CSP (static) | Set tight CSP on SPA/static responses (`default-src 'self'`; no CDN; `frame-ancestors 'none'`) |
| Static default | `<root>/web/dist` — **never** document or encourage `--static-dir` = project root |
| Static footgun | Help + quickstart warn; **refuse** `--static-dir` that resolves to project root (would expose `.trace/`) |
| Packaging | **Primary:** documented two-artifact (`trace` + `web/dist`). **Plus:** optional `go:embed` of `web/dist` as **fallback** when disk `index.html` missing (disk wins when present) |
| Seed HTTP | Keep `strict`/`task_id` → **501** `NOT_IMPLEMENTED`; CLI retains honesty export. Do **not** invent HTTP strict without a shared library helper |
| `mapDomainErr` | UUID / “must be UUID” (and plain `domain.ErrValidation`) → **400** `VALIDATION_ERROR` (not 500). Dogfood `rl…` loop ids |
| Promote honesty (low) | Discoveries promote: surface `createTransition` deny envelope (do not silently swallow) |
| Docs | User quickstart for `trace serve` + browser; AGENTS + project-rules carve-out from S00 RESEARCH |
| Cloud | Appendix **design only** (authn/z, tenancy, TLS) — no hosted implement |
| Perf / Law 6–7 | Keep budgeted `/v1/graph`; no unbounded dump endpoints; add/confirm list `limit` caps where missing |
| Out of scope | Multi-tenant SaaS, OAuth, public bind default, MCP `/rpc` browser, Three.js, full-graph export over HTTP |

## Work items (for S06-01)

- Security defaults audit (bind, token, CORS, CSP for static). Optional Vite-dev **exact-origin** CORS reflect (`--cors-origin`) — S03 defaults remain deny / no `*`
- Document `--static-dir` footgun; refuse project-root static dir; keep default `web/dist`
- Packaging: two-artifact quickstart + embed fallback when disk dist absent
- Note: HTTP seed export rejects `strict`/`task_id` with 501 — CLI retains honesty/gate export
- **S04/S05 residual:** `mapDomainErr` UUID → `VALIDATION_ERROR`
- **S05 residual (low):** promote `createTransition` deny honesty in SPA
- User docs quickstart + AGENTS.md / project-rules carve-out
- Cloud path appendix — design only
- Performance: pagination/`limit` defaults; no unbounded graph dump endpoints

## Exit criteria (this planner row)

- [x] Locked defaults table final (above)
- [x] `01-implement.md` + `02-review.md` + `SCOPE-TODOS.md` thickened for unattended S06-01
- [x] No product code; do not start P29-S06-01

## Next

P29-S06-01 → P29-S06-02
