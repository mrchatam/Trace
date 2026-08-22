# Phase 34 — GUI packaging + multi-project ports

Fix post–Phase 33 dogfood: **real Explore SPA embedded** in the Trace binary (consumers never need `web/`), **auto free-port** for concurrent `trace gui`, and **consumer Trace footprint = `.trace/` only**.

Design SoT: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) · [`INTAKE.md`](INTAKE.md). Board: [`docs/TODO/phase-34.md`](../../TODO/phase-34.md).

## Human locks (do not reopen)

| Lock | Value |
|------|-------|
| **L1 — Consumer layout** | User project: Trace may create/use **only** `.trace/` (db, lock, token, …). **Forbidden** as required layout: `web/`, project-root `trace.db`, copying SPA into the repo |
| **L2 — GUI asset source** | Default: serve **real** Explore SPA from Trace binary (`go:embed` of built `web/dist` at build/release). Disk `<project>/web/dist` is **Trace-checkout / contributor DX only**, never the consumer story. Stub unacceptable as shipped default when full dist was built into release |
| **L3 — Auto port** | `trace gui` / `serve`: default bind in use → **automatically** try next free loopback port; print chosen URL; open browser to that URL. Explicit `--addr` stays strict (fail if busy). **Overturns** P32/P33 “no auto-port” for the happy path |
| **L4 — One process = one project** | One store root per process (`-C`/cwd). Multi-project = multiple processes on distinct ports (auto) |
| Law 19 / loopback | Unchanged; no public bind defaults |

Full table: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Live repo baseline (P34-00, 2026-08-21)

| Area | State |
|------|-------|
| Static resolution | `internal/httpapi/static.go`: prefer disk `StaticDir` if `index.html`; else `go:embed` `embeddist/`; else inline placeholder HTML |
| Default StaticDir | `<projectRoot>/web/dist` — **consumer pain** (no `web/` → embed stub) |
| Embed today | **Stub only** (`internal/httpapi/embeddist/index.html`); README still teaches optional manual `cp dist → embeddist` + two-artifact everyday path |
| CLI | `trace gui` exists (`cmd/trace/local_http.go`); `--no-open`; help says **no auto free-port** |
| Port conflict | P32-PORT: `IsAddrInUse` + `FormatAddrInUseMessage` — **fail** + hint pick `--addr` (`internal/httpapi/addr_in_use.go`) |
| Default bind | `127.0.0.1:7432` (`httpapi.DefaultAddr`) |
| Docs | `docs/gui-quickstart.md` still documents disk `web/dist` wins / consumer-facing two-artifact path |
| Supersession | P33 RESEARCH **rejected** UA auto-increment port; **L3 overturns** that reject for default-bind happy path (still loopback; `--addr` pin remains strict) |

## Scope index (serial)

```
S00 research (embed vs sidecar; StaticDir; auto-port peers)
  → S01 plan (release embed + bind policy + test matrix)
  → S02 embed real SPA + static defaults (no consumer web/)
  → S03 auto free-port for gui/serve
  → S04 docs + residual tests
  → S05 VERIFY + DR-HANDOFF
```

| Scope | Title | Primary artifact | Board |
|-------|-------|------------------|-------|
| S00 | Research embed + auto-port | `RESEARCH.md` | P34-S00-00…02 |
| S01 | Plan pipeline + policy | `PLAN.md` | P34-S01-00…01 |
| S02 | Embed SPA + StaticDir defaults | embeddist + resolution + tests | P34-S02-00…02 |
| S03 | Auto free-port | gui/serve bind + tests | P34-S03-00…02 |
| S04 | Docs + residual tests | quickstart / help / AGENTS | P34-S04-00…02 |
| S05 | VERIFY + handoff | `VERIFY-NOTES.md`; close DR-HANDOFF | P34-S05-00…02 |

**Ownership notes**

- **Embed / StaticDir:** S00 recommends → S01 locks `PLAN.md` → **S02 ships** pipeline + resolution (consumer without `web/` gets real SPA).
- **Auto-port:** S00 recommends algorithm → S01 locks flags → **S03 ships** (after S02 on board; do not start S03 until S02-02 PASS).
- **Docs story:** S04 owns consumer-facing flip (no `web/` required; embed primary; auto-port Just Works). S02/S03 may update help minimally.
- **VERIFY:** S05 proves consumer-like temp (`.trace/` only) + real SPA + concurrent distinct ports.

## In scope

- Build/release embed of real `web/dist` into Trace binary
- Default static resolution that does **not** require consumer `web/`
- Auto free-port on default bind for `trace gui` / `serve`; print + open chosen URL
- Docs/help/quickstart/AGENTS aligned with L1–L3
- Tests for consumer-without-web SPA + concurrent ports

## Out of scope

- Hosted multi-tenant SaaS / OAuth / billing
- Always-on daemon; `0.0.0.0` as default bind
- Putting SPA under consumer `.trace/` as a copy of `web/dist` (prefer binary embed)
- Rewriting Phases 00–33 `done` history
- Second SQLite / business-logic fork in `web/` (Law 19)
- Explore/craft UI redesign (Phase 33 closed)

## Cross-scope blast radius

| If … | Then … |
|------|--------|
| S00 prefers install-sidecar over `go:embed` | Flag in RESEARCH; S01 must still meet L2 (SPA from Trace product, not consumer tree) — default lean remains embed |
| Release embed is stub | S05 VERIFY **FAIL**; S02 must make release path embed real dist |
| Auto-port only on `gui` not `serve` | Document in PLAN; prefer **both** share path (same local HTTP) unless PLAN justifies split |
| Docs still require consumer `web/` after S04 | S05 FAIL |
| Explicit `--addr` starts auto-hopping | Violates L3 — keep `--addr` fail-if-busy |

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until `P34-S05-02`. Successor lean: **no successor** unless VERIFY residuals need a thin follow-on.
