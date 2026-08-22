# P21-S04-02 — Review: baseline promotion

## Metadata
- id: P21-S04-02
- todo_ids: [P21-S04-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: baseline promote/supersede chain (B100→B101) + eval regression promotion gate match S04-00 locks; mig **020** only; **no auto DONE**; loop `promotion_blocked` advisory; compat ceiling **20**.

## Session start
**Fresh subagent** (not S04-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-baseline-promotion.md](01-baseline-promotion.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-09, D-10
- [WORK-MAP.md](../../WORK-MAP.md) W-05, W-06

## Review checklist

| Check | Evidence |
|-------|----------|
| Mig **020** only | `020_baselines_promotion.sql` present; schema dir has **20** files; **no 021+** |
| Column shape | `status` CHECK active/superseded; `supersedes_id` default `''`; index on commit+label+status |
| Baseline chain | At most one `active` per `git_commit`+`label`; B100→B101 supersede link preserved |
| `PromoteBaseline` | Supersedes prior active; sets `supersedes_id`; idempotent on re-promote |
| `SupersedeBaseline` | Status flip only — no DELETE (Law 11) |
| `CreateBaseline` | New rows default `active` |
| `CheckPromotionGate` | Latest evaluation with computed comparison; `eval_regression` when `overall_regression=true` |
| Gate vocabulary | Reasons: `eval_regression`, `no_stored_evaluation`, `baseline_not_found` (when applicable) |
| Test pass ≠ promotion | `TestPromotionGateIndependentOfTestPassAlone` green |
| Regression clears | `TestEvalRegressionGateClearsAfterResolve` green |
| No eval → blocked | `TestBaselinePromotionRequiresStoredEvaluation` green |
| **No auto DONE** | `TransitionTask` / `task_state.go` unchanged — grep no `CheckPromotionGate` in transition path |
| Loop status | `StatusResult.promotion_blocked` populated; **not** merged into `statusBlocked()` |
| Status schema | `trace.loop.status.v1` unchanged (additive field only) |
| Seed compat | Old seeds without status keys import as `active`; export includes chain when set |
| 6 named tests | All exist + PASS |
| P20 keepers | `TestEvaluationComparesScoresToBaselineNotBoolean`, `TestEvaluationRegressionFlagInComparisonJSON`, `TestPromotionGateRequiresStoredTestNotAgentClaim` green |
| Compat ceiling **20** | `evals/compat` + store embed tests assert 20; forbid 21+ |
| Event | `baseline.promoted` (or locked name) on promote |

## D-09 / D-10 closure

- **D-09 promote:** B100→B101 chain via `PromoteBaseline` + `supersedes_id` — not insert-only baselines.
- **D-10 promote:** `overall_regression` wired to `CheckPromotionGate` + loop `promotion_blocked`; advisory only.

## Keeper command floor

```bash
go test ./internal/domain/... -count=1 -run 'TestPromoteBaselineSupersedesPrior|TestPromoteBaselineIdempotent|TestEvalRegressionBlocksPromotionGate|TestEvalRegressionGateClearsAfterResolve|TestPromotionGateIndependentOfTestPassAlone|TestBaselinePromotionRequiresStoredEvaluation|TestEvaluationComparesScoresToBaselineNotBoolean|TestEvaluationRegressionFlagInComparisonJSON|TestPromotionGateRequiresStoredTestNotAgentClaim'
go test ./internal/store/... -count=1 -run 'TestBaselineStoreRoundtrip|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'
go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestLoopApplyDeliberationTransitionEvent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review focus

- **Latest eval wins:** Gate uses most recent evaluation by timestamp, not first regressed row.
- **Idempotent promote:** Double promote must not flip active baseline twice or duplicate supersede.
- **Orthogonal blocked signals:** `deliberation.blocked` (uncertainty/regression/debt) vs `promotion_blocked` (eval regression) — both may be true independently.
- **Blast radius:** S07 owns mig 021 — confirm S04 did not land experiment tables.
- **P20 keeper:** Boolean PASS rejection, verification≠test, agent claim gate unchanged.

## Spawn policy

- **Na (implement):** mig wrong/missing, chain broken, gate logic wrong, auto DONE introduced, compat not 20, named test missing/failing
- **Nb (review):** re-review after Na
- Do **not** spawn for optional dedicated CLI loop-status test if `CheckPromotionGate` unit tests + manual status JSON inspection suffice

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-09 + D-10 closure evidenced
- [ ] Spawn Na/Nb only if promotion/gate gap found

## Next

**P21-S05-00** (unless Na spawned)
