# Phase 44 / S00 / Planner — Intake + research

## Metadata
- id: P44-S00-00
- todo_ids: [P44-S00-00]
- role: planner
- skills: [research, code-explorer, domain-modeling, planning-and-task-breakdown]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Replan S00 against live repo. Thicken `01-research.md` + `02-review.md` so the researcher can produce `RESEARCH.md` (gap matrix vs [`INTAKE.md`](../../INTAKE.md)) without re-debating Phase 44 locks. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 13, 19
- [`INTAKE.md`](../../INTAKE.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §1 + D1–D5 light locks
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- P44-00 live anchors (do not re-litigate; verify still true):
  - `internal/retrieval/project_graph.go`
  - `internal/retrieval/graph_neighbors.go`
  - `internal/retrieval/neighborhood.go` (`GraphNode` / `GraphEdge`)
  - `internal/domain/service.go`, `internal/domain/doc.go`
  - `internal/mcp/tools_write.go` (`toolLink`)
  - `internal/store/links.go`, `internal/store/plan_hierarchy.go`
  - `web/src/lib/overviewCompose.ts`, `web/src/lib/graphLayout.ts`
  - `api/openapi.yaml` `/v1/graph`

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P44-00 light locks are authority.

## Locked defaults (from P44-00)

| Item | Value |
|------|-------|
| Theme | Meaningful scope clusters (auth FE/BE, business) — not UI-only |
| Output | `RESEARCH.md` in this folder |
| Product code | Forbidden in S00 |
| Cap baseline | GUI 500 / API 5000 (document; do not change) |
| D1–D5 | Light-locked; S00 documents gaps vs locks, does not reopen |

## Role work

1. Confirm P44-00 research table still matches HEAD.
2. Thicken `01-research.md` with: required RESEARCH.md sections, file touch-list (read-only), acceptance checks mapped to INTAKE examples.
3. Thicken `02-review.md` with completeness rubric (auth FE/BE + business row coverage).
4. Update `SCOPE-TODOS.md` checkboxes as planner completes.

## Exit criteria

- [x] `01-research.md` runnable alone (sections + file targets + acceptance)
- [x] `02-review.md` has APPROVE rubric vs INTAKE
- [x] Board Notes; next **P44-S00-01**

## Minimal todos

- [x] Re-verify live anchors vs P44-00 table
- [x] Thicken 01 + 02
- [x] Mark board row done with evidence

## Next

`P44-S00-01`
