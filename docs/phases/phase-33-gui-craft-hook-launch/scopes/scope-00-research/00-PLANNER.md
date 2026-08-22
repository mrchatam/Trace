# P33-S00-00 — Scope planner (research)

## Metadata
- id: P33-S00-00
- todo_ids: [P33-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, research]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock research scope so a fresh subagent can produce `RESEARCH.md`: peer launch matrix (Graphify `graph.html`, UA open-browser/`npx` viewer), **Laws 6–7–safe project overview graph** options, and PATH/`trace gui` launch recommendations for S01–S02. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 33 README](../../README.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- Live: `web/src/App.tsx`, `web/src/screens/Graph.tsx`, `cmd/trace/serve.go`, `cmd/trace/root.go`, `cmd/trace/install.go`, `docs/gui-quickstart.md`
- Peers (if present): `similar projects/graphify/`, `similar projects/Understand-Anything/`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Do not reopen DESIGN-LOCKS Themes A–C.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-research/RESEARCH.md` |
| Product / CLI / web edits | **No** on S00 |
| Explore today | Index = Graph (Explore); still **center-first** neighborhood — research the overview-hook upgrade, not “add a Graph route” |
| Laws 6–7 | Budgeted / progressive overview only — reject unbounded full dump as default |
| Launch lean | Prefer subcommand **`trace gui`** (flag `-gui` secondary); reuse serve + P32-PORT |
| PATH | Distinct from `trace install` (agents/MCP) |
| Sequence | S00 → S01 serial |

## Must answer (handoff to 01)

1. Peer launch patterns to borrow vs reject (UA open-browser, Graphify static html, Trace serve today)?
2. How to show a **project overview graph** under Laws 6–7 (clusters / caps / progressive expand / seed centers) — options + recommend one for S01?
3. Does overview need a new `/v1` op, or can existing search + graph + client composition suffice?
4. PATH install options ranked (`go install`, make symlink, package) + how docs should teach them?
5. Confirm CLI shape recommendation for S02 (`trace gui` primary)?

## Planner gate

- [ ] `01-research.md` runnable (metadata, exit criteria, RESEARCH template)
- [ ] `02-review.md` checklist vs DESIGN-LOCKS + INTAKE
- [ ] `SCOPE-TODOS.md` accurate
- [ ] Do **not** write `RESEARCH.md` in this planner row

## Exit criteria

- [ ] Research implementer prompt locked for a fresh subagent
- [ ] Board Notes cite what was thickened
- [ ] Next remains **P33-S00-01**

## Todo updates

Status + notes on **P33-S00-00** only.

## Next

`P33-S00-01`
