# Phase 26 — Loop investigations, planning & implementation

Human-promoted successor after E02 dogfood closed Phase 25 (P25-C validated; P25-2 installer check failed). Implements deferred **P25-A**, **P25-B**, and the **P25-2** wiring fix.

## Goal

Close the three open gaps from Phase 25 DR-HANDOFF so discovery→task promotion, loop recalibration/reset, and parent-orchestrator install text all work without re-litigating Phase 25 history.

## Evidence basis

- [Phase 25 DR-HANDOFF](../phase-25-orchestrator-gap-pass/DR-HANDOFF.md) (CLOSED; successor = Phase 26)
- E02: P25-1/P25-3 PASS; **P25-2 FAIL** (`ParentOrchestratorRule` defined but not in `cursorRulesMDCContent`)
- Intervention matrix: [INTERVENTION-MATRIX.md](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) (INT-01/02/05/06/09 + installer INT-04 wiring)

## Scope sequence (locked)

```
S00 audit → S01 plan → S02 P25-A → S03 P25-B → S04 installer → S05 VERIFY
```

| Scope | Theme | INT / gap |
|-------|--------|-----------|
| S00 | Investigation / codebase audit | INT-01/02/05/06/09 + P25-2 |
| S01 | Task breakdown (`PLAN.md`) | all |
| S02 | Discovery→task promotion | INT-01, INT-06 |
| S03 | Loop recalibration + deliberation reset | INT-02, INT-05, INT-09 |
| S04 | Wire `ParentOrchestratorRule` into install MDC | P25-2 / INT-04 surface |
| S05 | VERIFY + DR-HANDOFF close | D1–D6 |

**Parallelism:** S04 is small and *may* run beside S03 after S01, but **default board order is serial** (S00→…→S05). Do not start a higher-order pending row while a lower-order row is open.

## In scope

- Audit + plan artifacts under this phase folder
- Product changes in loop/store/MCP/install/CLI paths named by AUDIT.md / PLAN.md
- Installer wiring of existing `ParentOrchestratorRule` into cursor + Claude fallback rule bodies
- Tests keeping `go test ./internal/...` green

## Out of scope

- Daemon / HTTP / hosted service on Trace core
- Full-rebuild-on-any-change indexer
- Rewriting Phase 25 (or earlier) `done` board/prompt history
- Deep threshold debates in S00 (document options; S01 lists; implementer+tests decide)
- Silent autonomous task spawn without `loop apply` or explicit `trace add task`

## Phase locks

| Item | Value |
|------|-------|
| Product code on P26-00 | **No** (phase planner only) |
| SQLite migration | Only with schema version bump |
| Human gate on promotion | Required (no background spawn) |
| VERIFY signal | P25-2 PASS on fresh `trace install cursor --write` + `go test ./internal/...` |
| Successor | Decided at S05-02 only |

## Board

[`docs/TODO/phase-26.md`](../../TODO/phase-26.md) — first runnable after phase planner: **P26-S00-00**.
