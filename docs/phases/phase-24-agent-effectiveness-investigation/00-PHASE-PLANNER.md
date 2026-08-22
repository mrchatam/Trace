# P24-00 — Phase 24 scaffold: agent effectiveness investigation

## Metadata
- id: P24-00
- todo_ids: [P24-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents, research, grilling]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Lock Phase 24 against E01 dogfood + Phase 23 limits. Scaffold S01–S05 investigation scopes. **No product Go on this row.**

Human promoted after: E01 complete, G1 thought process did not work as intended (seed-only tasks, thin graph, loop saturation, orchestrator bypass).

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [phase README](README.md)
- [INVESTIGATION.md](INVESTIGATION.md)
- [experiments/RESULTS.md](../../../experiments/RESULTS.md) — E01 row
- [experiments/ab-incident-tracker/](../../../experiments/ab-incident-tracker/)
- Phase 23 [ENFORCEMENT.md](../phase-23-enforcement-choke-points/ENFORCEMENT.md) — what P23 did **not** fix
- Live: `internal/loop/`, `internal/deliberation/`, `cmd/trace/loop.go`, `cmd/trace-mcp/`, `internal/install/`
- [docs/TODO.md](../../TODO.md), [docs/TODO/phase-24.md](../../TODO/phase-24.md)

## Session start

User request: dedicated investigation phase. **Updated 2026-08-20:** E01 **Session B** (directed gap analysis) shows Trace **can** record discoveries/decisions when asked — but still no new tasks and hop_budget blocks verify. Phase 24 must compare Session A vs B, not treat E01 as uniformly failed.

## Live inventory (investigation targets)

| Surface | Phase | Investigation angle |
|---------|-------|---------------------|
| `trace loop next/apply/status/gate` | P19–P23 | Saturation, STOP, hop budget, agent adherence |
| `SelectNext` / deliberation phases | P20–P22 | When agents are told STOP vs EXECUTE |
| `trace add` / task linking | core | FM-01 seed anchoring |
| MCP tools (15) | P22 | Discovery, trace_add prominence, orchestrator use |
| `trace install` rules/hook | P23 | Orchestrator vs subagent coverage |
| Seed import shape | P17 | Fixed task list anchoring |
| Dogfood protocol | experiments | Multitask, prepare.sh, arm isolation |

## Locked defaults (phase)

| Item | Value |
|------|-------|
| Goal | Diagnose agent effectiveness gaps; ranked interventions |
| Method | Post-mortem + codebase audit + external research + matrix |
| Product Go | **None** in P24-00; investigate rows write docs/evidence only |
| Spikes | Only if a scope 01 row explicitly allows ≤50 lines proof |
| Forbidden | daemon; hosted MCP; rewriting P23 done history |
| Successor | DR-HANDOFF proposes Phase 25 themes — human promotes |

## Scope order (locked)

1. S01 dogfood-postmortem
2. S02 codebase-loop-audit
3. S03 external-research
4. S04 intervention-matrix
5. S05 VERIFY + DR-HANDOFF

## Planner work (this row)

1. Confirm INVESTIGATION.md covers failure modes + questions.
2. Ensure S01–S05 have `00-PLANNER` / `01-*` / `02-review` / `SCOPE-TODOS.md`.
3. Add Phase 24 to `docs/TODO.md` + `docs/TODO/phase-24.md`.
4. Update `AGENTS.md` current focus → P24 next runnable.
5. Open `DR-HANDOFF.md`.

## Exit criteria

- [ ] README + INVESTIGATION.md cover investigation MVP
- [ ] S01–S05 stubs exist
- [ ] Board rows P24-00 … P24-S05-02
- [ ] AGENTS.md next runnable `P24-S01-00` (post P24-00 execution)
- [ ] No product Go

## Next

Orchestrator: **P24-S01-00** after this row is `done`.
