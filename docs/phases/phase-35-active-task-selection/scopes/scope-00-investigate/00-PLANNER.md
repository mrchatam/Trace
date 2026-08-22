# P35-S00-00 — Scope planner (investigate)

## Metadata
- id: P35-S00-00
- todo_ids: [P35-S00-00]
- role: planner
- skills: [diagnosing-bugs, planning-and-task-breakdown, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock investigation scope so a fresh subagent can produce `INVESTIGATION.md`: feet-seller repro (CLI + HTTP/GUI), root-cause cites for Overview/Loop/listTasks, and a **red-capable** feedback loop description for S01–S02. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [Phase 35 README](../../README.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live: `web/src/screens/Overview.tsx`, `web/src/screens/Loop.tsx`, `web/src/api/ops.ts`, `internal/httpapi/handlers_tasks.go`, `internal/store/helpers.go`, `cmd/trace/AGENTS.md`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Do not reopen DESIGN-LOCKS theme. Fixture path is locked.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-investigate/INVESTIGATION.md` |
| Product / CLI / web edits | **No** on S00 |
| Fixture | `/home/ali/Desktop/feet seller telegram app` — read-only DB |
| Feedback loop | Prefer CLI + HTTP JSON asserts; GUI screenshot/notes secondary |
| Sequence | S00 → S01 serial |

## Must answer (handoff to 01)

1. Exact repro steps + observed bound task id/title on Overview and Loop (feet-seller).
2. Does HTTP `/v1/tasks` truncate at 100 today, or does client `limit: 100` have no effect?
3. Confirm or refute each INTAKE “likely cause” with file:line cites.
4. Is there any library/API “current work” / last-touched / plan-current-scope already?
5. Red-capable command or script outline that fails on “defaults to Step 1” (for S02 TDD).

## Planner gate

- [x] `01-investigate.md` runnable (metadata, exit criteria, INVESTIGATION template)
- [x] `02-review.md` checklist vs DESIGN-LOCKS + INTAKE
- [x] `SCOPE-TODOS.md` accurate
- [x] Do **not** write `INVESTIGATION.md` in this planner row

## Exit criteria

- [x] Investigate implementer prompt locked for a fresh subagent
- [x] Board Notes cite what was thickened
- [x] Next remains **P35-S00-01**

## Todo updates

Status + notes on **P35-S00-00** only.

## Next

`P35-S00-01`
