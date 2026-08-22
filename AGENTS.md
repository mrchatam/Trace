# Trace

Local-first project knowledge graph + progressive planning for AI coding agents.

## Stack

- **Go** module `github.com/mrchatam/Trace`
- CLI: `trace`
- SQLite under `.trace/trace.db` only (gitignored); never open or create project-root `trace.db` (agents: use CLI/MCP). Stray-root warn fires once per `openStore` when `<root>/trace.db` is a regular file; multiple CLI/MCP/HTTP opens may re-emit (intentional; no suppress flag).
- tree-sitter analyzers (TS/JS, Python)
- Git via CLI behind a VCS interface

## Agent workflow

1. **Source of order:** [`docs/TODO.md`](docs/TODO.md) (index) → active board under [`docs/TODO/`](docs/TODO/) — one subagent per row, strict order.
2. **Protocol:** [`docs/rules/agent-loop-protocol.md`](docs/rules/agent-loop-protocol.md)  
   Session start: **Agent mode → ask if unclear (grill/brainstorm/skills) → Plan mode → execute.**
3. **Laws:** [`docs/rules/project-rules.md`](docs/rules/project-rules.md) + [`docs/init/G_PROJECT_LAWS.md`](docs/init/G_PROJECT_LAWS.md)
4. **Skills map:** [`docs/rules/skills-map.md`](docs/rules/skills-map.md)
5. **Design baseline:** [`docs/`](docs/) + planning registers [`docs/init/`](docs/init/)

### Optional Codegraph (complement)

Symbol-level code exploration is **optional** via separate Codegraph MCP — Trace owns task loop + evidence. See [CONTRIBUTING — Trace + Codegraph](CONTRIBUTING.md#trace--codegraph-optional-dual-stack). Orchestrator/harness may register both MCP servers; user's choice per repo.

## Orchestrator (paste per phase)

```text
Phase 00–43 complete — do not re-run closed rows.
@docs/TODO.md

- Active phase: **none** — Phase 43 GitHub hygiene complete
- Next runnable: **idle** (human promotion only)
- Follow docs/rules/agent-loop-protocol.md
```

## Current focus

**Phase 43 complete** (2026-08-22) — closed at `P43-S01-02`; GitHub triage (0 open issues), extended evals aligned with Phase 29 HTTP policy, perf ceilings re-locked; CI + GHCR green on main. Board: [`docs/TODO/phase-43.md`](docs/TODO/phase-43.md).

**Phase 42 complete** (2026-08-22) — closed at `P42-S02-02`; G6 graph-label concept retrieval, G7 index freshness & langs delivered; M-001 preserved; REMEDIATION-PLAN G1–G9 complete. Board: [`docs/TODO/phase-42.md`](docs/TODO/phase-42.md).

**Phase 41 complete** (2026-08-22) — closed at `P41-S02-02`; G8 opt-in L2/L3 layers, G9 rule-based intent pipeline delivered; M-001 preserved. Board: [`docs/TODO/phase-41.md`](docs/TODO/phase-41.md).

**Phase 40 complete** (2026-08-22) — closed at `P40-S02-02`; G5 GUI graph orient, G2 unified `trace_explore` delivered; M-001 preserved. Board: [`docs/TODO/phase-40.md`](docs/TODO/phase-40.md).

**Phase 39 complete** (2026-08-22) — closed at `P39-S03-02`; G1 query merge, G3 harness orient, G4 dual-stack docs delivered; M-001 preserved. Board: [`docs/TODO/phase-39.md`](docs/TODO/phase-39.md).

**Phase 38 complete** (2026-08-22) — closed at `P38-S07-02`; peer-gap investigation saturated (11 gaps G-001…G-011); `REMEDIATION-PLAN.md` G1–G9 ranked; **plan only — no implement**. Board: [`docs/TODO/phase-38.md`](docs/TODO/phase-38.md).

**Phase 37 complete** (2026-08-22) — closed at `P37-S03-02`; P36 residuals R1–R11 closed (advisories[], HTTP plan routes, MCP loop gate, enforce nudge, plan UX, bootstrap help, tests, live GUI verify). R7/R9/R8-full re-deferred. Board: [`docs/TODO/phase-37.md`](docs/TODO/phase-37.md).

**Phase 36 complete** (2026-08-22) — closed at `P36-S03-02`; MCP `trace_plan`, bootstrap, terminal gate honesty shipped.

**Phase 35 complete** (2026-08-21) — closed at `P35-S03-02`; active-task pick P1→P2→P3a on Overview/Loop.

**Phase 34 complete** — embed SPA + auto free-port; consumer `.trace/` only.

## Hard boundaries

### Phase 29 — Opt-in local HTTP + browser GUI (carve-out)

- **Allowed:** Opt-in `trace serve` on Trace core that exposes a versioned HTTP/JSON API and serves the browser GUI. Default bind **`127.0.0.1`**. Explicit flags required for non-loopback bind; prefer a one-time or configured token when bind is not loopback.
- **Law 19:** HTTP handlers and the browser UI are **adapters only**. They call the canonical Go library/API. No second source of truth, no business-logic fork in `web/`, no parallel SQLite from the browser.
- **Still forbidden:** Always-on network daemon; open bind (`0.0.0.0` / public internet) as default; pointing local product MCP (`trace-mcp`) at the public internet; full-graph dump as default API/GUI behavior (Laws 6–7).
- **Cloud path:** OpenAPI is the shared contract for a **future hosted product** (separate deploy/repo). Phase 29 does not ship multi-tenant hosting, OAuth, billing, or tenancy. Design notes: [`docs/phases/phase-29-http-api-browser-gui/CLOUD-APPENDIX.md`](docs/phases/phase-29-http-api-browser-gui/CLOUD-APPENDIX.md).
- **Historical note:** FR-P28-X1 / older “no HTTP on core” language is superseded **only** for this opt-in local carve-out; it does not authorize silent daemons or public defaults.

### Other

- No full-rebuild-on-any-change indexer architecture
- No rewriting `done` board history — spawn forward
- Implementers: board **status + notes only**; reviewers/planners change upcoming work
- Closing phase VERIFY scaffolds next phase (DR-HANDOFF)
- **Portable graph:** Before a PR that changes Trace entities, export the semantic graph (including plan tree) to `trace/graph.json` via `trace seed export -o trace/graph.json` (`.trace/` stays local; git SHA is evidence, not identity — see CONTRIBUTING). Default export omits reviews/work_state; clone import tasks are PENDING — see CONTRIBUTING.
