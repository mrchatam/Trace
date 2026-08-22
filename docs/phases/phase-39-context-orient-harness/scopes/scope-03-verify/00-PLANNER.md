# P39-S03-00 — Scope planner (VERIFY)

## Metadata
- id: P39-S03-00
- todo_ids: [P39-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, shipping-and-launch, qa-lead]
- verification: automated

## Objective

Lock S03 VERIFY blocks for G1+G3+G4 delivery + DR-HANDOFF close policy. Thicken `01-verify.md` and `02-dr-handoff.md`. **No product code in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Pattern: [P38 S07-01 verify](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-07-verify/01-verify.md)
- Pattern: [P37 S03-01 verify](../../../phase-37-p36-residuals/scopes/scope-03-verify/01-verify.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: DR-HANDOFF secondary queue is authority for successor.

## Locked defaults (FINAL — P39-00)

| Item | Value |
|------|-------|
| Verify scope | G1 product + G3 product/docs + G4 docs-only |
| Precondition | P39-S00-02, P39-S01-02, P39-S02-02 all APPROVE |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p39-s03-01-verify/evidence/` |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required** at S03-01) |
| DR-HANDOFF | Stays **OPEN** until S03-02 |
| Successor | **Phase 40+** — G5 GUI orient + G2 unified explore (see DR-HANDOFF forward note) |
| Close owner | P39-S03-02 |
| Product boundary | S03-00/01 verify only — no feature work |

## Verify blocks (for 01-verify.md)

| Block | Check |
|-------|-------|
| 0 | G1 acceptance T1–T6 green + S00-02 APPROVE |
| 1 | G3 G3-A1–A6 + S01-02 APPROVE; 16 MCP tools |
| 2 | G4 G4-D1–D8 + S02-02 APPROVE; doc-only diff |
| 3 | M-001 moat preserved (task loop, gates, no query-only, no 1-tool facade) |
| 4 | Laws 6–7 caps honest; Law 19 library-first |
| 5 | DR-HANDOFF forward note lists G5/G2 secondary queue; Phase 40+ successor named |
| 6 | `trace seed export` if entities changed during P39 |

## Exit criteria

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable with blocks 0–6
- [ ] Successor never TBD at close template
- [ ] Board row → `done` with Notes

## Next

`P39-S03-01`
