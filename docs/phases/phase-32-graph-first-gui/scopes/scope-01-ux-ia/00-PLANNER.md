# P32-S01-00 — Scope planner (UX IA)

## Metadata
- id: P32-S01-00
- todo_ids: [P32-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock UX IA so S01-01 authors `UX-IA.md`: graph-home shell, inspector sections, which ops screens become panels, Law 6–7 budgets. Flag any missing `/v1` for S02. **No SPA rewrite in this scope.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [Phase 32 README](../../README.md)
- S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) (**PASS** after P32-S00-02) — especially **Handoff to S01**, **Gaps → IA — S01**, API reuse map, Laws 6–7 / 19
- Live: `web/src/App.tsx`, `web/src/layout/`, `web/src/screens/Graph.tsx`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output of S01-01 | `scopes/scope-01-ux-ia/UX-IA.md` |
| Shell strategy | Evolve `web/` into graph-home (**hybrid C** per DESIGN-LOCKS + RESEARCH) |
| Graph tech | 2D `@xyflow/react` only |
| Budgets | Keep `center` + `max_nodes` / `UI_CAP` 100 — no full-graph dump default (Laws 6–7) |
| Inspector depth map | Prefer why / context / impact / reviews (+ entity summary); reuse `/v1`; note `getImpact` client gap for S02 |
| Port work | **Out of S01** — P32-PORT ships in S02 (#1 min per RESEARCH) |
| Product code | **No** in S01 |

## Must answer (handoff to 01)

1. Default route / home = graph canvas?
2. Inspector section list (why, context, impact, reviews, links, …)?
3. Which current nav screens become secondary panels vs remain routes?
4. Selection / search → center → expand UX sketch?
5. Any missing library-backed `/v1` ops for S02 (else expect `NO-GAPS.md` later; client `getImpact` wrapper is glue, not a new core API)?

## Planner gate

- [x] `01-ux-ia.md` runnable
- [x] `02-review.md` checklist ready
- [x] Do **not** write `UX-IA.md` here

## Exit criteria

- [x] Implementer prompt locked; next **P32-S01-01**

## Todo updates

Status + notes on **P32-S01-00** only.

## Next

`P32-S01-01`
