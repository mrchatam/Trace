# P21-S08-00 — Planner: phase verify + DR-HANDOFF

## Metadata
- id: P21-S08-00
- todo_ids: [P21-S08-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock Phase 21 verify floor (P20 floor + S01–S07 deltas), live §31 mini-eval, and DR-HANDOFF close policy. **No product Go this row.**

## References
- [WORK-MAP.md](../../WORK-MAP.md) W-14
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — **OPEN** until S08-02
- [DECISION-LOG.md](../../DECISION-LOG.md) D-05, D-14
- P20 verify: [01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md)
- P20 review close: [02-scope-review.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/02-scope-review.md)
- Live mini-eval pattern: `experiments/runs/2026-08-18-p20-s07-01-verify/`

## Live inventory (confirmed 2026-08-18, post S01–S07)

| Surface | Location | Today (live read) | S08 verify proof |
|---------|----------|-------------------|------------------|
| Schema files | `internal/store/schema/` | **21** files; max **021_experiments.sql**; **no 022+** | `TestMigrationStatusReportsEmbedMax` + compat |
| Compat ceiling | `evals/compat/compat_test.go` | **21** (forbid **022+**) | `TestCompatibilitySecurityChecklist` |
| Seed P20 keys | `seed_export.go` / `seed_import.go` | 10 JSON keys; 11 tables | `TestSeedExportIncludesP20Cognition`, `TestSeedImportP20RoundTrip` |
| Retrieval P20 types | `retrieval/exact.go`, `types.go` | 8 types + `outcome` alias | `TestExactLookupUncertainty`, `TestLoopNextInvestigateNoRetrievalStderr` |
| FTS P20 sync | `store/fts.go`, `domain/cognitive.go` | Upsert syncs uncertainty/hypothesis | `TestSyncEntityFTSUncertainty`, `TestRebuildFTSIncludesP20Types` |
| SelectNext table | `deliberation/select.go` | **14-row FINAL** | `TestSelectNextFullCycleOrdering`, cycle row tests |
| Baseline promotion | `domain/outcomes.go` + mig **020** | `PromoteBaseline`, `CheckPromotionGate` | `TestPromoteBaselineSupersedesPrior`, `TestEvalRegressionBlocksPromotionGate` |
| Why + historical | `retrieval/why.go`, `loop/next.go` | P20 why steps; `historical_relationships` cap 8 | `TestCLIWhyUncertainty`, `TestLoopNextHistoricalRelationshipsSection` |
| Tx apply | `store/tx.go`, `loop/apply.go` | `WithTx` wraps mutating apply | `TestLoopApplyTransactionalRollbackOnFailure` |
| Goal guard | `loop/apply.go`, `domain/deliberation.go` | goal_id match fail-closed | `TestLoopApplyGoalIDMismatchFailsClosed` |
| internal/loop tests | `internal/loop/*_test.go` | **≥8** apply tests (D-15 closed) | count + named floor |
| Experiments | mig **021** + `domain/experiments.go` | thin table; no runner | `TestCreateExperimentLinksOutcome`, `TestNoExperimentRunnerInvoked` |
| Risk hints | `loop/risk_hints.go` | 4 codes, cap 4, advisory | `TestRiskHintsManyPaths`, `TestLoopNextRiskHintsBounded` |
| MCP catalog | unchanged | **10** tools | compat grep |
| Loop schema strings | `next.go`, `apply.go`, `status.go` | v1 unchanged | keeper tests |

### P21 delta → verify proof (FINAL)

| Delta | Verify proof |
|-------|--------------|
| Seed includes P20 (D-05 retired) | `TestSeedExportIncludesP20Cognition` + CLI export grep |
| Retrieval clean INVESTIGATE (D-06) | `TestLoopNextInvestigateNoRetrievalStderr` |
| FTS on P20 upserts (D-07) | `TestSyncEntityFTSUncertainty` |
| Full SelectNext cycle (D-03) | `TestSelectNextFullCycleOrdering` |
| Baseline promotion (D-09, D-10) | `TestPromoteBaselineSupersedesPrior` |
| Why P20 + historical (D-11, D-12) | `TestCLIWhyUncertainty`, `TestLoopNextHistoricalRelationshipsSection` |
| Transactional apply (D-08, D-13, D-15) | `TestLoopApplyTransactionalRollbackOnFailure` |
| Experiments + risk hints (D-01, D-02) | S07 6 named tests |
| Live §31 mini-eval (D-14) | CLI evidence **≥5 phases** |
| Compat ceiling | **21** (locked — S07 landed) |

## FINAL verify floor (S08-01)

**Block A:** Inherit **entire** P20 S07-01 command floor ([P20 01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md)).

**Block B:** P21 S01–S07 named keepers (see [01-verify.md](01-verify.md)).

**Block C:** Compat ceiling **21** + optional product bar.

Preflight: all named tests exist (grep confirmed 2026-08-18).

## §31 live mini-eval (FINAL)

| Item | Value |
|------|-------|
| Workspace | **Live repo** — build `trace` from checkout; temp DB — **not fixture-only** |
| Scenario | Import seed with P20 cognition → blocking uncertainty → INVESTIGATE → resolve → CRITIQUE → EXECUTE path → test/verify/eval → baseline promote / promotion_blocked |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p21-s08-01-verify/evidence/` |
| Minimum phases observed | **≥5** distinct `deliberation.phase` in CLI JSON (D-14 promoted from P20 fixture-scale) |
| Seed bar | Export **includes** P20 keys when DB populated (D-05 retired) |

## DR-HANDOFF close (S08-02 only)

| Field | Lock |
|-------|------|
| Who closes | **S08-02 only** |
| Successor | **`no successor`** default unless human names Phase 22 — **never TBD** |
| Residuals | experiments not in seed; D-16 autonomous runner out; D-19 MCP out |
| P20 DR-HANDOFF | Historical `no successor` **unchanged** |

## Named tests (S08-01 keeper summary)

| Scope | Count | Anchor tests |
|-------|------:|--------------|
| P20 floor | ~90+ named | Full block in `01-verify.md` Block A |
| S01 seed | 2 | `TestSeedExportIncludesP20Cognition`, `TestSeedImportP20RoundTrip` |
| S02 retrieval | 7 | `TestExactLookupUncertainty`, `TestLoopNextInvestigateNoRetrievalStderr`, FTS 2 |
| S03 cycle | 8 | `TestSelectNextFullCycleOrdering`, `TestSelectNextExecuteWhenPending`, … |
| S04 promotion | 6+ | `TestPromoteBaselineSupersedesPrior`, `TestEvalRegressionBlocksPromotionGate` |
| S05 observability | 5 | `TestCLIWhyUncertainty`, `TestLoopNextHistoricalRelationshipsSection` |
| S06 apply | 5+ | `TestLoopApplyTransactionalRollbackOnFailure`, goal_id guards |
| S07 MVP cuts | 6 | experiments 3 + risk hints 3 |
| Compat | 1 | `TestCompatibilitySecurityChecklist` (ceiling **21**) |

## Touch files

- `docs/phases/phase-21-thoughtprocess-completion/scopes/scope-08-phase-verify/01-verify.md` — thickened
- `docs/phases/phase-21-thoughtprocess-completion/scopes/scope-08-phase-verify/02-scope-review.md` — thickened
- `experiments/runs/` — S08-01 evidence
- `DR-HANDOFF.md` — S08-02 closes

## Planner work

1. [x] Live inventory: schema max **021**, compat **21**, all S01–S07 surfaces landed.
2. [x] Lock Block A (P20 floor) + Block B (P21 deltas) + Block C (compat).
3. [x] Lock §31 **live** mini-eval — **≥5 phases**, not fixture-only (D-14).
4. [x] Lock seed **include** P20 policy (D-05 retired).
5. [x] Lock DR-HANDOFF close: **`no successor`** default; never TBD.
6. [x] Thicken `01-verify.md` + `02-scope-review.md` with commands, evidence table, Must map, spawn policy.
7. [x] Update `SCOPE-TODOS.md`.
8. [x] No product Go.

## Exit criteria

- [x] Verify floor + live eval locked
- [x] DR-HANDOFF policy locked
- [x] 01/02 thickened enough for S08-01/02 alone
- [x] No product Go

## Next

**P21-S08-01**
