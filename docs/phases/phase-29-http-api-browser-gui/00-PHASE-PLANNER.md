# Phase 29 — Trace HTTP API + browser GUI

**Phase planner.** Runs as row `P29-00`.

## Metadata
- id: P29-00
- todo_ids: [P29-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified, research]
- mcps: []
- verification: automated
- hooks: []

## Human locks (2026-08-21)

| Lock | Value |
|------|-------|
| Delivery | **Browser-based GUI** |
| Transport | **HTTP API on Trace core** (local-first now) |
| Cloud | **Same API contract** prepared for later hosted product; multi-tenant cloud **not** deployed in this phase |
| Adapter law | UI + HTTP handlers call **canonical library** only (G_PROJECT_LAWS §19) — no second SoT |
| FR-P28-X1 | **Human-promoted / superseded** for local HTTP + cloud-ready contract; hosted SaaS remains later-developments separate product |

## Mission

Ship a **production-ready, feature-rich browser GUI** on top of Trace by:

1. Researching peer GUIs + Trace surfaces
2. Locking architecture + OpenAPI
3. Implementing local HTTP server (`trace serve` or equivalent)
4. Implementing browser GUI (MVP → rich)
5. Hardening for production local use + documenting cloud extension path

## Scope sequence

```
S00 Research → S01 Architecture/OpenAPI → S02 UX IA
  → S03 HTTP API → S04 GUI MVP → S05 GUI rich → S06 Production → S07 VERIFY
```

Serial default. Do not start S0N until S0(N−1) review is `done` (or explicit parallel note from a later planner).

## Hard constraints

- **Local default:** bind `127.0.0.1` only unless explicitly configured; no open-internet default
- **No second SQLite / business logic fork** in the UI
- **CLI + MCP remain first-class**; HTTP is an additional adapter
- Hosted OAuth/tenancy/billing = **out of Phase 29 deploy scope** (design hooks OK)
- No rewriting Phases 00–28 `done` history
- `go test ./internal/...` (and new packages) stay green

## Light locks (phase-level; scope planners may refine with evidence)

| Item | Default |
|------|---------|
| Package | `internal/httpapi` (ADR may rename with rationale) |
| Cmd | `trace serve` |
| API prefix | `/v1` |
| GUI root | `web/` SPA → static assets served by `trace serve` |
| Stack | S00 recommends; default lean **TypeScript + Vite + React** unless research overturns |
| Auth local | none on loopback; optional bearer when `--allow-remote` |
| OpenAPI path | `api/openapi.yaml` or phase `scopes/scope-01-architecture/openapi.yaml` (S01 locks) |

## End-state definition (production-ready)

| Area | Must have by S07 |
|------|------------------|
| API | Stable OpenAPI; health; read graph/tasks/loop; write transitions/add entities; seed import/export hooks |
| GUI | Project open, graph explorer, task/loop board, discoveries/decisions, gap/status, export honesty signals |
| Security | Loopback default; CORS locked; optional token for non-loopback; no secrets in browser storage beyond session |
| DX | `trace serve` (or documented cmd); README; `trace install` note if needed |
| Quality | Tests for API; GUI smoke/e2e for critical paths; VERIFY evidence |

## Planner gate (P29-00)

Verify before closing this row:

- [ ] `docs/phases/phase-29-http-api-browser-gui/` — README, this file, `DR-HANDOFF.md` **OPEN**
- [ ] `scopes/scope-00-research/` — `00-PLANNER`, `01-peer-and-surface-research`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-01-architecture/` — `00-PLANNER`, `01-adr-and-openapi`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-02-ux-ia/` — `00-PLANNER`, `01-ux-ia`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-03-http-api/` — `00-PLANNER`, `01-implement`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-04-gui-mvp/` — `00-PLANNER`, `01-implement`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-05-gui-rich/` — `00-PLANNER`, `01-implement`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-06-production/` — `00-PLANNER`, `01-implement`, `02-review`, `SCOPE-TODOS`
- [ ] `scopes/scope-07-verify/` — `00-PLANNER`, `01-verify`, `02-dr-handoff`, `SCOPE-TODOS`
- [ ] Board row **502 / P29-00** in `docs/TODO/phase-29.md`
- [ ] Upcoming prompts have protocol stubs (metadata + exit criteria); S00 thickened for first implement wave
- [ ] Live baseline noted: no `internal/httpapi`, no `web/`, no `serve` in `cmd/trace/root.go` yet

## Session start

Follow [agent-loop-protocol.md](../../rules/agent-loop-protocol.md) Session start (Agent → clarify → Plan → execute). Human locks above are settled — do not re-grill them.

## Todo updates

Status + notes on **P29-00** only. May thicken **upcoming** scope folders. No product code.

## Exit criteria

- [ ] Planner gate checklist satisfied
- [ ] README + light locks reflect live repo
- [ ] Each scope has adequate stub/thickened planner + implement + review prompts
- [ ] Board Notes cite what was verified/changed; next runnable **P29-S00-00**
- [ ] `docs/TODO.md` next pointer updated if needed

## Next

`P29-S00-00`
