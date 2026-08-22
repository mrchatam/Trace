# Phase 29 — HTTP API + browser GUI

Human-promoted 2026-08-21 after E03 stack validation. **Goal:** ship a **production-ready, feature-rich browser GUI** on Trace via an **opt-in local HTTP API** on Trace core (library-backed, Law 19), with an **OpenAPI contract** reusable by a future hosted/cloud product (separate deploy — not this phase’s target).

## Locked product direction

```
[Browser GUI] --HTTP/JSON--> [trace serve / internal/httpapi] --calls--> [internal/* library + SQLite .trace/]
                                                                      ^
                                                                      CLI + MCP also call here (Law 19)
```

Later cloud: same OpenAPI behind auth/tenancy gateway in a **separate hosted product** (`docs/TODO.md` Later developments). Phase 29 does **not** ship multi-tenant SaaS.

## Human locks (do not reopen)

| Lock | Value |
|------|-------|
| Delivery | Browser-based GUI |
| Transport | HTTP API on Trace core (local-first) |
| Cmd | Opt-in `trace serve`; default bind `127.0.0.1` |
| Cloud | Same API contract prepared; hosting/tenancy **out** of deploy scope |
| Adapter | UI + HTTP handlers → canonical library only (no second SoT) |
| FR-P28-X1 | Superseded for this carve-out only |

## Live repo baseline (P29-00, 2026-08-21)

| Area | State |
|------|-------|
| CLI | `cmd/trace/root.go` — no `serve` yet; surfaces include tasks, loop, seed, add, transition, review, context, why, search, … |
| MCP | `internal/mcp/` — `trace_tasks`, `trace_loop`, `trace_add`, `trace_why`, `trace_context`, … (parity inventory for S00) |
| HTTP package | **Absent** — expect `internal/httpapi` (or ADR name) in S03 |
| Browser GUI | **Absent** — expect `web/` (or ADR path) in S04 |
| Peers | `similar projects/Understand-Anything` (dashboard package + skills) |

## Scope index (serial)

```
S00 Research → S01 Architecture/OpenAPI → S02 UX IA
  → S03 HTTP API → S04 GUI MVP → S05 GUI rich → S06 Production → S07 VERIFY
```

| Scope | Title | Primary artifact |
|-------|-------|------------------|
| S00 | Peer + surface research | `RESEARCH.md` |
| S01 | Architecture ADR + OpenAPI | `ADR-HTTP-API-GUI.md`, `openapi.yaml` |
| S02 | UX information architecture | `UX-IA.md` |
| S03 | HTTP API + `trace serve` | `internal/httpapi` + CLI |
| S04 | Browser GUI MVP | `web/` P0 screens |
| S05 | Feature-rich GUI waves | P1 + `FEATURE-MATRIX.md` |
| S06 | Production hardening | security, packaging, docs |
| S07 | VERIFY + DR-HANDOFF | `VERIFY-NOTES.md`, successor |

## In scope

- Opt-in local HTTP server + OpenAPI-shaped JSON API
- Browser SPA served by (or alongside) `trace serve`
- Loopback-default security; optional token for non-loopback
- Tests for API; GUI smoke/e2e for critical paths
- Design hooks / appendix for future cloud (authn/z, tenancy) — docs only

## Out of scope

- Always-on network daemon; open bind `0.0.0.0` by default
- Multi-tenant hosted SaaS / OAuth / billing deploy
- Second SQLite or business-logic fork in the browser
- Pointing local `trace-mcp` at the public internet
- Rewriting Phases 00–28 `done` history

## Law note

Historical “no daemon/HTTP on Trace core” (FR-P28-X1 / prior AGENTS boundary) is **human-superseded for Phase 29** for opt-in local HTTP + browser GUI + cloud-ready API shape. Still forbidden without a later promote: silent always-on daemon, open bind by default, forking SoT into the UI.

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until S07-02.
