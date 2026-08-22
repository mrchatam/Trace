# P39-S01-00 — Scope planner (G3 harness orient)

## Metadata
- id: P39-S01-00
- todo_ids: [P39-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering, writing-for-agents]
- mcps: [user-trace]
- verification: automated

## Objective

Lock S01 **G3**: MCP discovery & harness orient (G-006, G-010, 9/16 fold). Thicken implement/review prompts. **No product code in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — Q5 9/16 resolution
- [REMEDIATION-PLAN §2 G3](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [PEER-CG §3](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22):
  - `internal/mcp/server.go:28–36` — `NewServer` no Instructions yet
  - `internal/mcp/server.go:227–235` — 16 locked tools
  - `internal/mcp/server.go:153–161` — `trace_version` for stale detection
  - `internal/install/cursor.go:12–13` — `CursorReloadTip` (DF-22/50)
  - `internal/install/bootstrap_hint.go:12–36` — plan bootstrap hint
  - `internal/install/agents.go:17–57` — harness agent defaults
  - `CONTRIBUTING.md:64–70` — agent workflow + MCP reload note
  - go-sdk v1.4.0 — `mcp.ServerOptions.Instructions` supported

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE locks are authority.

## Locked defaults (FINAL — P39-00)

| Item | Value |
|------|-------|
| GAP ids | G-006, G-010 (+ 9/16 harness fold) |
| Verdict | **Accept** per REMEDIATION-PLAN G3 |
| Peer pattern | CG SERVER_INSTRUCTIONS + ranked orient (**contrast** — do not copy 1-tool-only) |
| Lead tools | Moat-first: `trace_tasks` → `trace_context`(+query) → `trace_loop` → `trace_review` → `trace_plan` |
| Read tools | `trace_search`, `trace_why`, `trace_impact`, `trace_capability` — secondary, capped |
| Write tools | Remain fully registered — **never hide** behind read facade |
| 9/16 fix | **Docs + MCP Instructions + trace_version callout** — stale stdio root cause |
| G2 compose-first | Orient recipe lives in **MCP Instructions** (Phase 39) — not unified explore |
| MCP tool count | **16** unchanged |
| Out | CG 1-tool default; MP 44-tool copy; tool count reduction |

## Accept / reject (G3)

| Decision | Item |
|----------|------|
| **Accept** | MCP server `Instructions` via go-sdk `ServerOptions` in `NewServer` |
| **Accept** | Moat-first orient playbook in Instructions (task loop, gates, plan tree) |
| **Accept** | Ranked read-tool sequence + compose-first G2 recipe in Instructions |
| **Accept** | CONTRIBUTING + install path doc/st stderr reinforcing reload after rebuild |
| **Accept** | Optional: strengthen `PrintBootstrapHintIfNeeded` / install agents messaging for moat |
| **Reject** | Reduce registered tools to fix 9/16 |
| **Reject** | Hide write/task tools from MCP registration |
| **Reject** | Bundled Codegraph MCP |
| **Reject** | `trace_explore` implement |

## Must lock for S01-01 (delivered in thickened 01-implement)

1. Touch-list: MCP Instructions → install/docs → tests.
2. Five acceptance criteria (see `01-implement.md`).
3. Instructions content outline (moat-first + ranked tools + stale-server hygiene).

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] 9/16 scope locked: docs + Instructions (not tool reduction)
- [ ] Board row → `done` with Notes

## Next

`P39-S01-01`
