# P40-S02-00 — Scope planner (VERIFY)

## Metadata
- id: P40-S02-00
- todo_ids: [P40-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, qa-lead, shipping-and-launch]
- verification: automated

## Objective

Lock S02 VERIFY blocks for G5+G2 deliverables + Phase 41+ successor policy. Thicken `01-verify.md` + `02-dr-handoff.md`. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Pattern: [P39 S03-00 verify planner](../../../phase-39-context-orient-harness/scopes/scope-03-verify/00-PLANNER.md)
- Pattern: [P39 S03-01 verify](../../../phase-39-context-orient-harness/scopes/scope-03-verify/01-verify.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: DR-HANDOFF secondary queue is authority.

## Locked defaults (FINAL — P40-00)

| Item | Value |
|------|-------|
| Verify scope | G5 product (S00) + G2 product (S01) |
| Precondition | P40-S00-02, P40-S01-02 both **APPROVE** |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p40-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required** at S02-01) |
| DR-HANDOFF | Stays **OPEN** until S02-02 |
| Successor | **Phase 41+** — G8 layers + G9 intent (secondary queue G6/G7 documented, not P40 implement) |
| Close owner | P40-S02-02 |
| Product boundary | S02-00/01 verify only — no feature work |

## Verify blocks (for 01-verify.md)

| Block | Check |
|-------|-------|
| 0 | G5 G5-A1–A7 + S00-02 APPROVE |
| 1 | G2 G2-T1–T7 + S01-02 APPROVE; **17** MCP tools |
| 2 | M-001 moat preserved (task loop primary; explore merges; no query-only; no 1-tool facade) |
| 3 | Laws 6–7 caps honest; Law 19 library-first (GUI + MCP adapters thin) |
| 4 | Secondary queue G6/G7 documented in DR-HANDOFF — **no P40 implement rows** |
| 5 | Phase 41+ successor named (G8/G9) — never TBD |
| 6 | `trace seed export` if entities changed during P40 |

## Exit criteria

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable with blocks 0–6
- [ ] Successor never TBD at close template
- [ ] Board row → `done` with Notes

## Next

`P40-S02-01`
