# P32-S00-00 — Scope planner (research)

## Metadata
- id: P32-S00-00
- todo_ids: [P32-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, brainstorming, research]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock research scope so a fresh subagent can produce `RESEARCH.md`: peer matrix (Graphify, Understand-Anything), current `web/` inventory vs explorer bar, depth/visual gap list, `/v1` reuse map, and a short **P32-PORT** recommendation for S02. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 32 README](../../README.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- Live: `web/src/App.tsx`, `web/src/screens/Graph.tsx`, `api/openapi.yaml`, `internal/httpapi/bind.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Do not reopen DESIGN-LOCKS.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-research/RESEARCH.md` |
| Product / serve edits | **No** on S00 |
| Peer bar | Between Graphify and Understand-Anything (Trace keeps plan/task/decision semantics) |
| P32-PORT | Note + recommend options for S02; do **not** implement |
| Sequence | S00 → S01 serial |

## Must answer (handoff to 01)

1. What does current `web/` do vs explorer job (graph-home + rich inspector)?
2. Which peer patterns to borrow vs reject (cite briefly)?
3. Depth gaps vs visual gaps (ordered for S03 vs S04)?
4. Which `/v1` ops already cover inspector needs?
5. P32-PORT: recommend #1 friendly error / #2 auto-port / #3 docs-only multi-`--addr` (min prefer #1)?

## Planner gate

- [ ] `01-research.md` runnable (metadata, exit criteria, RESEARCH template sections)
- [ ] `02-review.md` checklist vs DESIGN-LOCKS + P32-PORT section
- [ ] `SCOPE-TODOS.md` accurate
- [ ] Do **not** write `RESEARCH.md` in this planner row

## Exit criteria

- [ ] Research implementer prompt locked for a fresh subagent
- [ ] Board Notes cite what was thickened
- [ ] Next remains **P32-S00-01**

## Todo updates

Status + notes on **P32-S00-00** only.

## Next

`P32-S00-01`
