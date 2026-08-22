# P29-S00-00 — Scope planner (research)

## Metadata
- id: P29-S00-00
- todo_ids: [P29-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified, research]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock investigation scope for Phase 29 research against **live** Trace surfaces and peer GUIs. Finalize `01-peer-and-surface-research.md` + `02-review.md` so a fresh subagent produces `RESEARCH.md` (peer matrix, Trace surface map, stack recommendation, API resource families for S01). **No product HTTP/GUI code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [Phase 29 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [docs/TODO.md](../../../../TODO.md) Later developments (hosted product boundary)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-research/RESEARCH.md` |
| Product Go / web | **No** on S00-01 |
| Peers (minimum) | `similar projects/Understand-Anything` dashboard + ≥2 other modules under `similar projects/` |
| CLI inventory | `cmd/trace/root.go` switch + related cmd files |
| MCP inventory | `internal/mcp/` tool catalog (`mcp_test.go` lists expected names) |
| Law | Handlers/UI → library only; draft carve-out text for AGENTS (applied later in S06) |
| Stack default lean | TypeScript + Vite + React — overturn only with evidence in RESEARCH |
| Sequence | S00 → S01 → … serial |

## Live anchors (verify still true)

| Topic | Path / note |
|-------|-------------|
| No `serve` yet | `cmd/trace/root.go` — no `case "serve"` |
| No httpapi | `internal/httpapi` absent |
| No web GUI | `web/` absent |
| MCP tools | `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`, `trace_search`, `trace_changes`, `trace_regressions`, `trace_loop`, `trace_agents` |
| Peer dashboard | `similar projects/Understand-Anything/understand-anything-plugin/packages/dashboard` |
| Peer skill | `similar projects/Understand-Anything/understand-anything-plugin/skills/understand-dashboard` |
| Peer graph-ui | `similar projects/codebase-memory-mcp/graph-ui` |
| Peer frontend | `similar projects/agentrq/frontend` |

## Planner gate

- [ ] `01-peer-and-surface-research.md` runnable (metadata, preflight, template, exit criteria)
- [ ] `02-review.md` has verify checklist + spawn policy
- [ ] `SCOPE-TODOS.md` lists S00 board rows
- [ ] Live anchors above still accurate (adjust `01` if renamed)

## Exit criteria

- [ ] Research implementer prompt locked enough for a fresh subagent
- [ ] Board row P29-S00-00 Notes cite what was verified/thickened
- [ ] Next runnable remains **P29-S00-01** (do not start S01)

## Todo updates

Status + notes on **P29-S00-00** only.

## Next

`P29-S00-01`
