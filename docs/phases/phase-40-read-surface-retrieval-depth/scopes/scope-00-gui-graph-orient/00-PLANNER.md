# P40-S00-00 — Scope planner (G5 GUI graph orient)

## Metadata
- id: P40-S00-00
- todo_ids: [P40-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, ux-designer, frontend-ui-engineering]
- mcps: [user-trace]
- verification: automated

## Objective

Lock S00 **G5** against live repo: graph-first onboarding UX via GUI orient adapter (G-008). Thicken `01-implement.md` + `02-review.md` with file targets, acceptance map, and rejects. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P40-00 resolutions
- [REMEDIATION-PLAN §2 G5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-008](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [Phase 39 VERIFY-NOTES](../../../phase-39-context-orient-harness/scopes/scope-03-verify/VERIFY-NOTES.md) — G5 forward queue
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P40-S00-00):
  - `web/src/App.tsx:21–22` — `/` index → `<Graph />`; `/graph` redirects to `/`
  - `web/src/layout/Nav.tsx:4` — nav label **Explore** for `/` (page `<h1>` still “Graph” — orient copy should say **Explore**)
  - `web/src/screens/Graph.tsx` (~705 lines) — seed-composed overview; ReactFlow canvas
    - Law 6–7 static banner `:470–473` (`center` + `max_nodes` required)
    - Dynamic truncation `:605–611` + budget line `:612–617` (`graphTruncated`, `data-testid="graph-budget-line"`)
    - Mount point for orient panel: after `<h1>` / before `.page-lead` (`:456–461`)
  - `web/src/lib/overviewCompose.ts` — pure seed/merge helpers (`SEED_CAP=8`, `UI_CAP=100`, `DEPTH=2`, `SEED_MAX_NODES=40`, `EXPAND_MAX_NODES=50`)
  - `web/src/api/ops.ts` — `getGraph`, `search`, `listTasks`, `getProject` (thin HTTP client)
  - `web/src/lib/overviewCompose.test.ts` — **node:test** runner (no vitest/`npm test` in `web/package.json`)
  - `internal/httpapi/handlers_retrieval.go:158–201` — `GET /v1/graph` (center + max_nodes required)
  - `internal/httpapi/handlers_retrieval.go:121–155` — `GET /v1/search`
  - `internal/retrieval/neighborhood.go:53+` — canonical bounded graph (library)
  - `internal/install/bootstrap_hint.go` — moat-first install hint pattern (optional one-line graph-first pointer)
  - `CONTRIBUTING.md:64–72` — moat-first MCP/CLI orient (no graph-first GUI subsection yet)
  - **Gap vs G-008:** no `GraphOrientPanel.tsx`; no first-visit orient panel; no install hook narrative for `trace serve` → `/`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: INTAKE + P40-00 locks are authority.

## Locked defaults (FINAL — P40-00)

| Item | Value |
|------|-------|
| GAP ids | G-008 |
| Verdict | **Accept** per REMEDIATION-PLAN G5 |
| Peer pattern | GF graph.html + `/graphify` hook; UA onboard artifact; MP `wake_up()` identity story (**contrast** — adapt, do not port) |
| Route strategy | **Enhance existing `/` Graph route** — graph UI already shipped (P32–P34); not greenfield route |
| Sketch vs full | **Not static-only** — live ReactFlow + seed compose exists; add orient UX layer + install narrative |
| Law 19 | GUI adapter over canonical HTTP/library only — **no business logic fork in `web/`** |
| API surface | Prefer existing `/v1/project`, `/v1/tasks`, `/v1/graph`, `/v1/search`; new HTTP route only if preceded by library function |
| M-001 | Orient narrative leads with **task loop + gates** — graph is entry surface, not product identity |
| Laws 6–7 | Preserve progressive caps UI (budget line, truncation banners) — orient copy must teach capped reads |
| Out | Graphify port; parallel SQLite from browser; graph-only product drift; planner logic in `web/` |
| MCP / G2 | Out of scope — S01 owns `trace_explore` |

## Accept / reject (G5)

| Decision | Item |
|----------|------|
| **Accept** | First-visit orient panel on `/` (dismissible; localStorage or session flag) |
| **Accept** | Moat-first onboarding copy: Explore → pick task → Loop → gate → review |
| **Accept** | Confidence/budget labels tied to existing truncation UI (`graphTruncated`, budget line) |
| **Accept** | Install hook narrative — CONTRIBUTING and/or `internal/install/` pointer (graph-first GUI after `trace serve`) |
| **Accept** | Preserve existing seed-compose graph behavior (no regression) |
| **Reject** | Business logic in `web/` beyond API calls + presentational compose |
| **Reject** | Full static `graph.html` replacement of React SPA |
| **Reject** | Full-graph dump default or cap increase |
| **Reject** | New MCP tools or compiler changes in S00 |

## Must lock for S00-01 (delivered in thickened 01-implement)

1. Touch-list: `GraphOrientPanel.tsx` + `Graph.tsx` + `app.css` + `CONTRIBUTING.md` (+ optional `bootstrap_hint.go`).
2. Seven acceptance criteria G5-A1–A7 (see `01-implement.md`).
3. Law 19 boundary: graph compose stays in `overviewCompose.ts`; orient is presentation only.
4. Tests: **node:test** strip-types pattern (same as `overviewCompose.test.ts`); no vitest in `web/`.

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable with file targets + G5 accept map
- [ ] SCOPE-TODOS updated
- [ ] Board row → `done` with Notes

## Next

`P40-S00-01`
