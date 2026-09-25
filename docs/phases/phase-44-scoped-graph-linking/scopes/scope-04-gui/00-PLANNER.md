# Phase 44 / S04 / Planner — GUI

## Metadata
- id: P44-S04-00
- todo_ids: [P44-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-ui-engineering]
- mcps: [user-trace]
- verification: automated

## Objective

Plan scope-aware layout in Graph / graphLayout (Law 19 adapter only). Thicken `01-implement.md` with edge priority + cluster spec. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §6, D5
- Phase 40 G5 orient baseline
- `web/src/screens/Graph.tsx`
- `web/src/lib/graphLayout.ts`
- `web/src/lib/overviewCompose.ts`
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol. Requires S02 (and preferably S03) APPROVE for edge data.

## Locked defaults

| Item | Value |
|------|-------|
| Law 19 | GUI consumes API fields only — no second graph walk in browser |
| Cap | `PROJECT_MAX_NODES = 500` unchanged |
| Priority | Elevate MVP scope rels in `EDGE_PRIORITY_RELS`; reduce `goal_has_task` dominance when scope edges exist |
| Inferred | Dashed / distinct style from `edge.provenance` |
| Orient | Legend: scope + inference honesty |

## Role work

1. Spec cluster-by-scope layout algorithm inputs (scope_id from API).
2. Thicken implement acceptance: ≥2 clusters + ≥1 cross-scope edge with seeded fixture.
3. a11y + perf notes at 500 cap.

## Exit criteria

- [x] `01-implement.md` with edge priority + cluster spec
- [x] Next **P44-S04-01**

## Minimal todos

- [x] Thicken 01 + 02
- [x] Board done

## Next

`P44-S04-01`
