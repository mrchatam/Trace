# P20-00 — Phase 20 scaffold: cognitive deliberation

## Metadata
- id: P20-00
- todo_ids: [P20-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents, grilling]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Lock Phase 20 against live repo + [`docs/TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md) so **every document section is owned** ([COVERAGE.md](COVERAGE.md)). Scaffold S01–S07. **No product Go on this row.**

Do not rewrite Phase 19 `done` history. P19 DR-HANDOFF historical `no successor` stays true at close; this phase is a forward human queue.

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — especially Laws 1–3, 6–8, 11–16, 19
- [phase README](README.md)
- [COVERAGE.md](COVERAGE.md)
- [docs/TRACE_THOUGHTPROCESS.md](../../TRACE_THOUGHTPROCESS.md)
- Live: `internal/loop/{next,apply}.go`, `cmd/trace/loop.go`, `internal/domain/create.go`, `internal/store/{entities,entities_causal,impact}.go`, `internal/planner`
- [docs/TODO.md](../../TODO.md) (index)
- [docs/TODO/phase-20.md](../../TODO/phase-20.md) (Phase 20 row table)

## Session start

Human requested Phase 20 that **comprehensively covers** TRACE_THOUGHTPROCESS. Locks below are the scaffold defaults for P20-00 to confirm/thicken — not a license to implement Go here.

## Live inventory (do not fork)

| Surface | Present | Gap vs thoughtprocess |
|---------|---------|------------------------|
| `trace loop next/apply/status` | **Yes** P19 | No phase controller; apply writes discovery/plan_change/task only |
| Goal/Task/Decision/Assumption/Discovery/PlanChange | **Yes** | No Uncertainty/Question table; Assumption has no first-class invalidate→replan |
| DecisionAlternative + impact findings + uncertainty vocab on findings | **Yes** | No reconsideration triggers on Decision |
| Claim/Evidence/Review + DONE policy | **Yes** | Test/Verify/Evaluate not distinct result kinds; no verification debt |
| planner phases/scopes/deep-plan + discovery replan N=5 | **Yes** | Not the cognitive phase machine |
| Git CLI / file+symbol index / impact walk | **Yes** | No Change entity; no expected/actual effects |
| Events | **Yes** thin (`events` + `AppendEvent`) | Need `deliberation.transition` inspectability (reuse table; S01) |
| Experiments / baselines / regressions | **No** | S04/S05 thin; §16 experiments Future |
| `internal/deliberation` | **No** | S01 new library |
| P19 Loop keeper tests | **Yes** [`cmd/trace/loop_test.go`](../../../cmd/trace/loop_test.go) | Must stay green through S06 unless dual-version plan |

## Locked defaults (phase)

| Item | Value |
|------|-------|
| Goal | State-driven deliberation controller on top of P19 loop |
| Keep | `internal/loop` schemas `trace.loop.{next,apply,status}.v1` — **extend**, do not break P19 tests without dual-version plan |
| New lib | `internal/deliberation` |
| Phases | ORIENT INVESTIGATE EXPLORE PLAN CRITIQUE EXECUTE TEST VERIFY EVALUATE REFLECT REPLAN |
| Policy | Deterministic, table-tested SelectNext; inspectable event |
| Artifacts | See COVERAGE.md merge table |
| Transport | stdout-first; model-agnostic |
| Git | SHA/path refs only (Law 1) |
| Trust | Agent JSON untrusted; gates need evidence (Law 2, 15) |
| Forbidden | daemon; hosted MCP; embeddings; graph DB; raw CoT; replacing P19; implementing §16/§18 product |
| Events | Reuse `store.AppendEvent`; type `deliberation.transition` (S01) — no parallel event store |
| Assumption invalidate | Provenance status transition on existing row (S02) — Law 11 |
| Task gates | Linked result records + verification debt (S04); do not add new `work_state` values without S04-00 proof |
| MCP | No new MCP tools by default (S06); CLI/library inherit G19 |
| Hop budget N | Locked in S01-00 (not this row) |
| Schema evolution | S06 chooses additive P19 fields vs dual-version envelope |

## Scope order (locked)

1. S01 deliberation-controller
2. S02 cognitive-artifacts
3. S03 change-effects
4. S04 verify-evaluate-gates
5. S05 regression-reflect
6. S06 protocol-context
7. S07 VERIFY + DR-HANDOFF

## Planner work (this row)

1. Confirm COVERAGE.md vs TRACE_THOUGHTPROCESS §§1–32 (no silent drops).
2. Ensure S01–S07 have `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
3. Board Phase 20 after Phase 19, before Later developments.
4. Update AGENTS.md current focus.
5. Open DR-HANDOFF.md.

## Exit criteria

- [x] README + COVERAGE.md cover all 32 sections
- [x] S01–S07 stubs exist
- [x] Board rows P20-00 … P20-S07-02
- [x] AGENTS.md next runnable `P20-S01-00` (post P20-00 execution)
- [x] No product Go

## Planner execution (2026-08-18)

Executed against live repo. Confirmed COVERAGE.md owns TRACE_THOUGHTPROCESS §§1–32 with no silent drops (§16/§18 explicitly Future). Reused seams: `internal/loop/{next,apply}.go`, `cmd/trace/loop.go`, `internal/domain/create.go`, `internal/store/{entities,entities_causal,impact}.go`, `internal/planner`, existing `events` table. Applied light phase locks in README + locked-defaults table above. Thickened upcoming S01–S07 `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS` with named live files, P19 keeper tests, fail-closed paths, and do-not-re-debate locks. DR-HANDOFF remains **OPEN**; S07-02 owns successor. No product Go. No new board rows spawned.

## Next

Orchestrator: **P20-S01-00** after this row is `done`.
