# Phase 44 — Intake

**Source:** user promotion 2026-08-23 (post Phase 43 idle).

## Problem statement

The Explore **project graph** (`GET /v1/graph?mode=project`) visually collapses into a star: most nodes appear connected only to the **first goal**, with little indication of *what a node is within a feature*, *how peers in the same scope relate*, or *why cross-cutting edges exist*.

## User examples (desired linking)

| Scope | Intra-scope | Cross-scope |
|-------|-------------|-------------|
| Frontend auth | login form ↔ session store ↔ route guard | ↔ backend auth via **API contract** edge |
| Backend auth | handler ↔ middleware ↔ token service | ↔ frontend via OpenAPI/route contract |
| Business | marketing landing ↔ design system ↔ human operator task | ↔ goal/task hierarchy without flattening to one goal |

## Success criteria (human)

1. Graph shows **clusters** aligned with logical scopes (auth, billing, design, …), not one hub goal.
2. Edge **semantics** are visible (`api_contract`, `implements`, `blocks`, `scope_member`, …).
3. Agents can **create** scope links (`trace_link` / decisions / discoveries); system may **infer** weak links with honest provenance.
4. Laws 6–7 preserved: bounded API, progressive expansion, no default full dump.

## Questions resolved at scaffold (planner may refine)

| # | Question | Draft answer |
|---|----------|--------------|
| Q1 | UI-only? | **No** — domain + API first, GUI follows |
| Q2 | New phase? | **Yes — Phase 44** on board |
| Q3 | 500 cap | GUI default (`PROJECT_MAX_NODES`); API hard cap 5000 — see 01-DESIGN-LOCKS |
| Q4 | Scope entity vs tag? | **TBD in S01** — prefer link rels + optional `scope` label table over orphan nodes |

## References (live code)

- `internal/retrieval/project_graph.go` — collect all entities; center = first goal
- `internal/retrieval/graph_neighbors.go` — walks `entity_links` + `goal_id`
- `internal/domain/doc.go` — documented rel vocabulary
- `web/src/lib/overviewCompose.ts` — `PROJECT_MAX_NODES = 500`
- `api/openapi.yaml` — `/v1/graph` max 5000

## Non-goals (intake)

- Embedding-based auto-clustering
- Public/hosted graph API changes beyond OpenAPI parity
- Rewriting Phase 40–43 delivery history
