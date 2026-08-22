# Phase 32 — Graph-first GUI

Raise Trace’s browser UI from Phase 29’s **ops shell** to an **explorer** bar (Graphify ↔ Understand-Anything): **graph as home**, ops as deep panels, **depth before visual craft**.

**Active** after Phase 31 close (2026-08-21). Design SoT: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md). Open work: [`OPEN-PORT-MULTI.md`](OPEN-PORT-MULTI.md) (**P32-PORT**). First runnable after phase planner: **P32-S00-00**.

## Human locks (do not reopen)

| Lock | Value |
|------|-------|
| Primary job | **Hybrid:** graph home; tasks / loop / reviews / seed = deep panels on selection |
| Depth vs craft | **S03 depth first**, **S04 visual second** |
| Graph tech | Keep **2D `@xyflow/react`**; **no Three.js / 3D** default |
| Shell | **Evolve `web/`** — not a second SPA product |
| Laws | **Law 19** adapters only; **Laws 6–7** budgeted neighborhood / `max_nodes` |
| API | Prefer existing `/v1`; S02 only for proven library-backed gaps **plus always P32-PORT** |
| Out | Hosted SaaS, always-on daemon, public bind defaults, MCP `/rpc` in browser, auto-delete root `trace.db` |

Full table: [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Live repo baseline (P32-00, 2026-08-21)

| Area | State |
|------|-------|
| GUI | `web/` ops shell — Overview is index; Graph is a **route** (`App.tsx`), not graph-home |
| Screens | Overview, Tasks(+detail), Loop, Graph, Discoveries, Reviews(+detail), Seed, Settings |
| Graph | `@xyflow/react` + budgeted `getGraph` (`DEFAULT_MAX=50`, `UI_CAP=100`) |
| API | `/v1` already has why, context, impact, graph, search, reviews (`api/openapi.yaml`) |
| Serve | Default `127.0.0.1:7432` (`httpapi.DefaultAddr`); fail on bind conflict — **P32-PORT** |
| Docs | [`docs/gui-quickstart.md`](../../gui-quickstart.md) — multi-project ports still thin |

## Scope index (serial)

```
S00 research (+ P32-PORT note)
  → S01 UX IA
  → S02 API gaps + P32-PORT ship
  → S03 depth
  → S04 visual
  → S05 polish (docs ports)
  → S06 VERIFY + DR-HANDOFF
```

| Scope | Title | Primary artifact | Board |
|-------|-------|------------------|-------|
| S00 | Peer + gap research | `RESEARCH.md` (incl. P32-PORT) | P32-S00-00…02 |
| S01 | UX information architecture | `UX-IA.md` | P32-S01-00…02 |
| S02 | API gaps **and** P32-PORT | OpenAPI diffs and/or `NO-GAPS.md` + port/serve ship | P32-S02-00…02 |
| S03 | Inspector / graph-home depth | Evolve `web/` per UX-IA | P32-S03-00…02 |
| S04 | Visual craft | Typography, density, motion on depth shell | P32-S04-00…02 |
| S05 | Production polish + docs | `gui-quickstart` multi-project ports | P32-S05-00…02 |
| S06 | VERIFY + handoff | `VERIFY-NOTES.md`; close DR-HANDOFF | P32-S06-00…02 |

**P32-PORT ownership:** S00 notes → **S02 ships** (even if API is `NO-GAPS.md`) → S05 docs → S06 VERIFY ticks.

## In scope

- Graph-home shell + rich node inspector (why / context / impact / reviews / links / path-filter)
- Budgeted graph explore (no unbounded dump)
- Address port conflict / multi-project `trace serve` (friendly error and/or auto free-port + docs)
- Visual craft **after** depth lands
- Docs / quickstart updates for explorer + multi-project ports

## Out of scope

- Hosted multi-tenant SaaS / OAuth / billing
- Always-on daemon; `0.0.0.0` as default bind
- Second SQLite or business-logic fork in `web/`
- Pointing local `trace-mcp` at the public internet
- 3D / Three.js as default canvas
- Rewriting Phases 00–31 `done` history
- Silent auto-delete of root `trace.db`

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until `P32-S06-02`. Successor lean: **no successor** unless residuals need a thin follow-on.
