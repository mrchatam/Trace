# P39-S02-00 — Scope planner (G4 dual-stack docs)

## Metadata
- id: P39-S02-00
- todo_ids: [P39-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents]
- verification: automated

## Objective

Lock S02 **G4**: dual-stack documentation (G-011). **Doc-only** — CONTRIBUTING/AGENTS Trace+Codegraph recipe. Thicken implement/review prompts. **No product code in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md)
- [REMEDIATION-PLAN §2 G4](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [h11-stack-docs evidence](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [PEER-CG §5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors:
  - `CONTRIBUTING.md` — Trace-only agent workflow today
  - `AGENTS.md` — project orchestrator doc (Trace moat, no CG section)
  - `.trace/` vs `.codegraph/` — separate storage (Law 19)

## Session start

Follow agent-loop-protocol Session start. Unattended: H11 doc-only lock is authority.

## Locked defaults (FINAL — P39-00)

| Item | Value |
|------|-------|
| GAP id | G-011 |
| Verdict | **Accept doc-only**; **Reject** product dual-index default |
| Mode | **Doc-only** — zero product Go/TS changes |
| Primary files | `CONTRIBUTING.md`, `AGENTS.md` |
| Optional | `README.md` one-paragraph pointer if CONTRIBUTING section is long |
| Content | When to `codegraph init` vs `trace index`; orient use-cases; Law 19; storage separation |
| Must not | Mandatory dual-index; bundled MCP; adapter fork in core; `.codegraph/` writes from Trace |
| Cross-ref | S01 MCP Instructions may pointer-only — full recipe here |

## Doc acceptance checklist (for S02-01)

| ID | Requirement |
|----|-------------|
| G4-D1 | Section title identifies Trace + Codegraph as **complementary**, not merged product |
| G4-D2 | Storage: `.trace/` (task/plan/evidence) vs `.codegraph/` (symbol graph) — separate indexes |
| G4-D3 | When to use Trace: task loop, gates, plan tree, evidence, progressive task packet |
| G4-D4 | When to use Codegraph: symbol exploration, call paths, blast radius — **optional** per repo |
| G4-D5 | Law 19: each stack is adapter/MCP over its own store; Trace core does not index into `.codegraph/` |
| G4-D6 | Setup recipe: `trace index` / install path + optional `codegraph init` — neither required for the other |
| G4-D7 | Reject language explicit: no default dual-index, no bundled dual MCP in Trace product |
| G4-D8 | Link to Phase 38 PEER-CG complement note / PEER-FIXTURES for investigation context |

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable with checklist G4-D1–D8
- [ ] Board row → `done` with Notes

## Next

`P39-S02-01`
