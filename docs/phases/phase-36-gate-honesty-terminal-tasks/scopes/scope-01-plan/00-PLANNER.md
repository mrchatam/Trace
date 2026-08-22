# P36-S01-00 — Scope planner (plan)

## Metadata
- id: P36-S01-00
- todo_ids: [P36-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, domain-modeling]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Lock S01 plan scope from `INVESTIGATION.md`. Thicken `01-plan.md` so implementers pick a **fundamental** fix set — not GUI-only. **No product code. No PLAN.md in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) (when exists)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live: `internal/loop/policy.go`, `internal/mcp/server.go`, `cmd/trace/plan.go`, `internal/install/`

## Session start

Follow agent-loop-protocol Session start. Waits for S00-02 PASS.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S01-01 | `scopes/scope-01-plan/PLAN.md` |
| Product edits | **No** on S01-00 |
| Primary fix axis | Trace product gaps (MCP / bootstrap / install) — not "agents should have known" |
| GUI | Secondary honesty only; never sole deliverable |
| Preserve | Active-work PLAN enforcement (`PlanExists` for non-terminal tasks) |

## Must decide (for 01-plan)

Pick from DESIGN-LOCKS candidates — **all first-class options** must appear in thickened `01-plan.md`:

1. **MCP plan tools** — `trace_plan` (create-coarse, set-current, deep, show) mirroring CLI; Law 19 library adapter
2. **Bootstrap command** — `trace plan bootstrap --goal` from plan-changes or seed import
3. **Install contract** — first goal → mandatory create-coarse in AGENTS.md + cursor rules; `trace install` docs
4. **PlanExists bridge** — whether plan-change density can satisfy gate (careful: must not fake plan)
5. **Terminal gate** — library short-circuit for DONE/SKIPPED + goal-level advisory (not global weaken)
6. **Enforce nudge** — document or default `warn` when `.trace/` without config
7. **Goal structure warning** — >N tasks under goal with no coarse plan
8. **Feet-seller recovery** — in S02 vs deferred residual (backfill / acknowledge legacy)

Also lock:

9. Touch-list order — library → MCP → install → HTTP/GUI adapter
10. Acceptance tests — greenfield MCP workflow + feet-seller live
11. Explicit non-goals — GUI-only patch rejected

## Planner gate

- [ ] `01-plan.md` runnable with DESIGN-LOCKS candidates as required sections
- [ ] `SCOPE-TODOS.md` lists MCP/bootstrap/install first-class
- [ ] Do **not** write `PLAN.md` in this planner row

## Next

`P36-S01-01`
