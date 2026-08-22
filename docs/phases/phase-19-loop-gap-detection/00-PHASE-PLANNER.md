# P19-00 — Phase 19 scaffold: loop gap detection

## Metadata
- id: P19-00
- todo_ids: [P19-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents, grilling]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Against live repo state and D42 taskboard findings, lock a **thin** Phase 19 that adds a harness-agnostic loop surface to Trace. Scaffold minimal runnable scopes only. **Do not** implement product Go on this row.

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](README.md)
- [experiments/RESULTS.md](../../../experiments/RESULTS.md) — D41, D42
- [AGENTS.md](../../../AGENTS.md)
- [`similar projects/Understand-Anything/README.md`](../../../similar%20projects/Understand-Anything/README.md)
- [docs/TODO.md](../../TODO.md)

## Session start

Agent -> clarify if needed -> Plan -> execute (planner). Human already chose the key Phase 19 lock:

- core loop interface = **stdout-first**
- MVP stop rule = **tasks saturated**

## Locked defaults (phase)

| Item | Value |
|------|-------|
| Goal | Thin CLI loop MVP for progressive gap detection |
| Core commands | `trace loop next`, `trace loop apply`, `trace loop status` |
| Transport | **stdout-first** packet protocol; optional wrappers later |
| Stop rule | No new tasks and no new plan changes, or max iterations |
| Context packet | tasks + why + plan + context + freshness + related files/symbols |
| Freshness | Must be explicit in packet |
| Product surface | CLI/library only; no daemon/hosted service |
| Harness support | Must work in shell and be wrapper-friendly for MCP/IDE/higher-level harnesses |
| Non-goals | full autonomous coding loop, hosted orchestration, embeddings |

## Scope order (locked)

1. **S01 loop-next-packet** — packet builder + freshness + related-files neighborhood
2. **S02 loop-apply-saturation** — structured apply and saturation status
3. **S03 phase-verify** — CLI proof + mini-eval + DR-HANDOFF

## Planner work

1. Inventory existing `tasks` / `why` / `context` / `impact` / planner APIs versus loop MVP needs.
2. Write minimal scope stubs (`00` / `01` / `02` / `SCOPE-TODOS`) for S01-S03.
3. Add Phase 19 rows to `docs/TODO.md` as the next runnable work after Phase 18 and before Later developments.
4. Update `AGENTS.md` current focus to mention Phase 19 scaffold + next runnable.
5. Create `DR-HANDOFF.md` in OPEN state for the phase.

## Exit criteria

- [x] Phase README explains why/how/locks
- [x] S01-S03 each have minimal runnable stubs
- [x] `docs/TODO.md` has Phase 19 rows and next runnable = `P19-00`
- [x] `AGENTS.md` current focus updated
- [x] No product Go in this row

## Minimal todos

- [x] Lock commands, stop rule, and non-goals
- [x] Spawn S01-S03 stubs
- [x] Board + AGENTS sync

## Next

Orchestrator: **P19-S01-00** after this row is marked `done`.
