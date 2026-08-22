# P38-S00-00 — Scope planner (investigation index)

## Metadata
- id: P38-S00-00
- todo_ids: [P38-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, research]
- mcps: [user-trace]
- verification: automated

## Objective

Lock S00: author **`INVESTIGATION-INDEX.md`** — hypothesis register (H1–H11 + spawn slots), peer map, investigation methods, spawn rules. **No product code. No GAP-REGISTRY in this scope.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- P24 [EXTERNAL-RESEARCH.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md)

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md` |
| Product edits | **Forbidden** |
| Mode | Investigate only — no implement suggestions as tasks |
| Spawn rule | New H* or peer slice → new board row in S01–S03, not silent backlog |

## Must answer (handoff to 01-investigate)

1. Per H1–H11: investigation method (Trace live / peer read / both), peer cite targets, success criteria for **verified vs rejected**.
2. Which optional tools (Trace MCP, Codegraph MCP, Graphify examples) apply per hypothesis.
3. Evidence path convention under `experiments/runs/`.
4. When to spawn extra investigation rows vs fold into S01–S03.
5. Explicit **non-goals** (no REMEDIATION-PLAN in S00).

## Planner gate

- [ ] `01-investigate.md` runnable alone
- [ ] `02-review.md` checklist vs DESIGN-LOCKS saturation prerequisites
- [ ] `SCOPE-TODOS.md` board IDs 647–649 accurate
- [ ] Do **not** write `INVESTIGATION-INDEX.md` in planner row

## Exit criteria

- [ ] S00-01/02 prompts thickened
- [ ] Board `P38-S00-00` → `done`

## Next

`P38-S00-01`
