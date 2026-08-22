# P32-S00-01 — Research

## Metadata
- id: P32-S00-01
- todo_ids: [P32-S00-01]
- role: implementer
- skills: [research, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Author `RESEARCH.md` in this scope folder: peer matrix (Graphify, Understand-Anything), current `web/` inventory vs explorer bar, ordered depth/visual gaps, `/v1` reuse map, and a short **P32-PORT** section with an explicit S02 recommendation. **No product UI/serve code.** Do **not** edit DESIGN-LOCKS or reopen locks.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [Phase 32 README](../../README.md)
- [00-PLANNER.md](00-PLANNER.md)
- Live Trace: `web/src/App.tsx`, `web/src/layout/Shell.tsx`, `web/src/layout/Nav.tsx`, `web/src/screens/*`, `web/src/api/ops.ts`, `api/openapi.yaml`, `internal/httpapi/bind.go`, `internal/httpapi/server.go`, `cmd/trace/serve.go`
- Peers (local clones): `similar projects/graphify/`, `similar projects/Understand-Anything/` (viewer README under `understand-anything-plugin/packages/viewer/`)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-00-research/RESEARCH.md` only |
| Code / OpenAPI / serve edits | **None** |
| Peer bar | Between **Graphify** (relation explore / path / explained edges) and **Understand-Anything** (search + graph + rich side inspector). Trace keeps plan/task/decision/review semantics. |
| Shell | Evolve existing `web/` — not a second SPA |
| Graph tech | 2D `@xyflow/react` only; **no** Three.js / 3D recommendation |
| Laws | Law 19 adapters; Laws 6–7 budgeted neighborhood (`center` + `max_nodes`); no full-graph dump as product default |
| Live GUI baseline (2026-08-21 planner check) | Index = **Overview**; Graph = **`/graph` route** (not home). Screens: Overview, Tasks(+detail), Loop, Graph, Discoveries, Reviews(+detail), Seed, Settings (+ ProjectGate). Graph: search→center→budgeted `getGraph` (`DEFAULT_MAX=50`, `UI_CAP=100`); **no** side inspector on Graph. TaskDetail has why/context drawers only. |
| `/v1` present in OpenAPI | `getWhy`, `getContext`, `getImpact`, `getGraph`, `search`, `listReviews` / `getReview` (+ tasks/loop/seed/health). Note: `web/src/api/ops.ts` wraps why/context/graph/search/reviews but **no** `getImpact` helper yet — research must call this out in the API reuse map (gap vs OpenAPI, not invent new API unless library-backed). |
| Serve / P32-PORT | Default `127.0.0.1:7432` (`httpapi.DefaultAddr`); `ListenAndServe` = bare `net.Listen` — fail on conflict; CLI prints `serve: %v`. Prefer recommend **#1** friendly `EADDRINUSE` + `--addr` guidance; optionally note UA viewer auto-increment port as peer pattern for #2. **S02 always ships** P32-PORT even if API is `NO-GAPS.md`. |
| Sequence | S00 → S01 serial; this row only writes RESEARCH |

## Preflight (confirm in Notes / RESEARCH, do not change code)

1. `App.tsx`: index → Overview; path `graph` → Graph.
2. `openapi.yaml`: `/v1/why`, `/v1/context`, `/v1/impact`, `/v1/graph`, `/v1/search`, `/v1/reviews`.
3. `bind.go` + `server.go` ListenAndServe: default addr + fail-on-conflict matches OPEN-PORT-MULTI.
4. Peers readable under `similar projects/` (if a peer tree missing, cite README/docs still available and mark peer rows incomplete — do not block on network).

## Role work

Write `RESEARCH.md` using the **template sections below** (headings required). Keep claims sourced (file path or peer README). Answer the planner must-answer set:

1. What does current `web/` do vs explorer job (graph-home + rich inspector)?
2. Which peer patterns to **borrow** vs **reject** (brief cite)?
3. Depth gaps vs visual gaps (ordered for **S03** vs **S04**)?
4. Which `/v1` ops already cover inspector needs (and which client wrappers are missing)?
5. **P32-PORT:** recommend #1 friendly error / #2 auto-port / #3 docs-only multi-`--addr` (minimum prefer **#1**; S02 owns ship).

### RESEARCH.md template (required headings)

```markdown
# Phase 32 S00 — RESEARCH

## Summary
(3–6 sentences: Trace today vs explorer bar; top depth gap; P32-PORT lean.)

## Peer matrix
| Dimension | Graphify | Understand-Anything | Trace today | Borrow / reject |
|-----------|----------|---------------------|-------------|-----------------|
| Home surface | … | … | Overview nav shell; Graph is a route | … |
| Search → graph | … | … | Graph has search+task pick | … |
| Node inspector | … | … | Graph: none; TaskDetail: why/context drawers | … |
| Path / filter | … | … | Kind filter on search only | … |
| Budget / dump | … | … | getGraph center+max_nodes; UI_CAP 100 | … |
| Local serve / ports | … | UA viewer: `--port`, auto-increment if taken, tokenized URL | Fixed :7432; fail on bind | … |
| Semantics Trace must keep | — | — | plan/task/decision/review/loop | do not clone code-AST-only UX |

### Peer notes (cite)
- Graphify: …
- Understand-Anything: …

## web/ inventory vs explorer job
| Screen / area | Path | Role today | Fit for graph-home explorer |
|---------------|------|------------|-----------------------------|
| Overview | `/` | … | … |
| Graph | `/graph` | … | … |
| TaskDetail | `/tasks/:id` | … | … |
| … | … | … | … |

## Gaps (ordered)
### Depth — S03 candidates (priority order)
1. …
### Visual craft — S04 candidates (only after depth)
1. …
### Explicit non-goals this phase
- 3D / Three.js, second SPA, hosted SaaS, MCP `/rpc` in browser, …

## API reuse map (`/v1` → inspector)
| Inspector need | OpenAPI op | Client (`ops.ts`) | Used in UI today | Notes |
|----------------|------------|-------------------|------------------|-------|
| Why | getWhy | yes | TaskDetail drawer | … |
| Context | getContext | yes | TaskDetail drawer | … |
| Impact | getImpact | **missing wrapper** | no | … |
| Neighborhood | getGraph | yes | Graph | … |
| Search | search | yes | Graph | … |
| Reviews | listReviews / getReview | yes | Reviews screens | … |
| … | … | … | … | … |

Flag any **library-backed** hole for S02 (else expect `NO-GAPS.md` for API) — P32-PORT is separate and still required.

## P32-PORT
### Confirm light review
(Default addr, Listen behavior, multi-project conflict — cite OPEN-PORT-MULTI + live files.)

### Peer patterns
(e.g. UA auto-increment port + tokenized URL — borrow idea vs reject full copy.)

### Recommendation for S02
| Option | Ship? | Rationale |
|--------|-------|-----------|
| #1 Friendly EADDRINUSE + `--addr` examples | **Prefer yes (min)** | … |
| #2 Auto free-port / `:0` | optional / defer? | … |
| #3 Docs-only multi-`--addr` | S05 docs; not sole S02 story | … |

**S02 ownership:** always ships P32-PORT even if API work is `NO-GAPS.md`. Discouraged: `NO-PORT-CHANGE.md` without written reason.

## Handoff to S01
Bullet list of IA constraints RESEARCH implies (graph-home, panels-not-nav-CRUD, Law 6–7 budgets). No UX-IA writing in this row.
```

## Exit criteria

- [ ] `RESEARCH.md` exists with all template headings filled (no empty stubs)
- [ ] Peer borrow/reject and depth-before-visual ordering explicit
- [ ] API reuse map cites real `/v1` ops and notes `getImpact` client gap if still true
- [ ] P32-PORT recommendation explicit for S02 (prefer #1 minimum)
- [ ] No product code / OpenAPI / serve diffs
- [ ] Board Notes cite artifact path `docs/phases/phase-32-graph-first-gui/scopes/scope-00-research/RESEARCH.md`

## Minimal todos

- [ ] Preflight live paths + OpenAPI + bind
- [ ] Skim Graphify + UA viewer READMEs under `similar projects/`
- [ ] Draft peer matrix + web inventory
- [ ] Ordered depth (S03) vs visual (S04) gaps
- [ ] API reuse map + P32-PORT section
- [ ] Update board row **P32-S00-01** status + Notes only

## Todo updates

Status + notes on **P32-S00-01** only.

## Next

`P32-S00-02`
