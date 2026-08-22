# P21-S08-01 — Verify Phase 21

## Metadata
- id: P21-S08-01
- todo_ids: [P21-S08-01]
- role: verify
- skills: [writing-for-agents, planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Run the **locked verify floor** (complete P20 S07-01 floor + P21 S01–S07 deltas), archive **live §31 mini-eval** evidence under `experiments/runs/YYYY-MM-DD-p21-s08-01-verify/`, and map results to [COVERAGE.md](../../../phase-20-cognitive-deliberation/COVERAGE.md) Must items. **Does not** close DR-HANDOFF (S08-02 owns).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar (this scope)
- [WORK-MAP.md](../../WORK-MAP.md) W-14
- [DECISION-LOG.md](../../DECISION-LOG.md) D-05, D-14 (retired / promoted)
- P20 verify floor: [01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md)
- [TRACE_THOUGHTPROCESS.md §31](../../../../TRACE_THOUGHTPROCESS.md) (Desired Outcome)
- P20 evidence pattern: `experiments/runs/2026-08-18-p20-s07-01-verify/`

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row verifies and records evidence; it does not open new product direction beyond Phase 21 locks.

## Locked defaults (FINAL — S08-00)

| Item | Value |
|------|-------|
| Schema max | **021** (`015`–`021` present; **no 022+**) |
| Compat ceiling | **21** — re-lock at VERIFY time via `TestCompatibilitySecurityChecklist` |
| Loop schemas | **additive v1** unchanged: `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` |
| CLI surface | `trace loop next --task <id>`, `trace loop apply [--in <path>]`, `trace loop status --task <id> [--goal <id>]`, `trace why …`, `trace seed export\|import` |
| Hop budget | **N=12** (`internal/deliberation`) |
| SelectNext table | **14-row FINAL** (S03): EXPLORE row 6; EXECUTE→TEST→VERIFY→EVALUATE→REFLECT→REPLAN rows 8–13 |
| Seed export (P20 entities) | **Include** 10 JSON keys (11 tables; `change_paths` nested in `changes[].paths`) — D-05 **retired** |
| Pre-PR export | Run `trace seed export -o trace/graph.json` once; P20 keys present when `.trace/` populated |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p21-s08-01-verify/evidence/` |
| Notes artifact | Optional `VERIFY-NOTES.md` in this scope folder (recommended) |
| Mini-eval | **Live repo** — real `trace/` CLI against temp DB; **≥5** distinct deliberation phases (not fixture-only residual) |
| Out-of-scope | No product Go; no DR-HANDOFF close; no Phase 22 scaffold |

### Seed export bar (locked — replaces P20 omit policy)

P20 cognition tables **must** round-trip in `BuildSeedDocument` / `ImportSeedDocument` (`internal/domain/seed_export.go`, `seed_import.go`):

- `deliberation_states`
- `uncertainties`, `hypotheses`, `decision_reconsiderations`
- `changes` (with nested `paths`), `effects`
- `outcome_results`, `baselines` (optional `status`, `supersedes_id` after S04)
- `regressions`, `reflections`

**Not in seed (unchanged):** transitions, work_state, reviews, capabilities, tokens, index surfaces.

**VERIFY requires:**

1. P17 keepers still green (no regression on portable graph):
   ```bash
   CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
   ```
2. P21 seed round-trip:
   ```bash
   go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'
   ```
3. One `trace seed export -o trace/graph.json` (or temp path) captured in evidence; grep confirms P20 keys when DB populated.
4. Notes state: portable clone of P20 cognition is **implemented** (D-05 closed); empty local `.trace/` may omit keys — unit tests prove presence when populated.

## Locked verify command floor (FINAL)

S08-01 must run and report PASS/FAIL for **every** block below. `-count=1` recommended. `GOMODCACHE`+`GOPROXY=off` on full product bar if offline.

### Block A — Complete P20 S07-01 floor (must stay green)

Copy verbatim from [P20 01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md) "Locked verify command floor":

```bash
# --- S01 controller + hop budget ---
go test ./internal/deliberation/... -count=1 -run 'TestSelectNext|TestSelectNextNeverExecuteOnBlockingUncertainty|TestApplyTransitionStickyInvestigateConsumesHopBudget|TestApplyTransitionHopBudgetDoesNotIncrementPastN|TestTransitionPayloadJSONRequiredFields'

# --- S02 cognitive artifacts (14 named) ---
go test ./internal/domain/ -count=1 -run 'TestCreateUncertaintyDefaultsOpenInfo|TestBlockingUncertaintyRequiresTaskID|TestBlockingUncertaintyIncrementsCountForTask|TestInfoUncertaintyDoesNotIncrementBlockingCount|TestResolveUncertaintyClearsBlockingCount|TestSupersedeUncertaintyClearsBlockingCount|TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestInvalidateAssumptionSetsStaleAndKeepsRow|TestInvalidateAssumptionSupersededNoDelete|TestInvalidateAssumptionEmitsImpactFindingOnLinkedDecision|TestInvalidateAssumptionOptionalPlanAffectingDiscovery|TestHypothesisLinksEvidenceWithoutDiscoveryTable|TestDecisionReconsiderPreservesDecisionAndAlternatives|TestUnknownUncertaintySeverityFailClosed'

# --- S03 change/effects (14 named) ---
go test ./internal/domain/ -count=1 -run 'TestCreateChangeWithGitSHAAndPathsNoBlob|TestCreateChangeRequiresTaskIDAndPath|TestRecordExpectedThenActualSupported|TestRecordActualRequiresExpectedDimension|TestUnknownEffectComparisonFailClosed|TestRecordActualContradictedLinksHypothesisWithoutDiscoveryFork|TestRecordActualContradictedOptionalPlanAffectingDiscovery|TestContradictedEffectFiresDecisionReconsideration|TestContradictedEffectDoesNotCreateRegressionOrAutoReplan|TestParentChangeChain|TestResolveChangePathViaGitNotSQLite|TestResolveChangePathFailsClosedWithoutCommit|TestOversizedEffectTextFailClosed|TestRecordChangeCommitThenCompared'

# --- S04 verify/evaluate gates (14 named) ---
go test ./internal/domain/ -count=1 -run 'TestRecordTestOutcomeRequiresNameAndStatus|TestTestPassAloneCannotSatisfyVerificationGate|TestVerificationRequiresGoalAndEvidenceIDs|TestVerificationMissingEvidenceFailClosed|TestEvaluationComparesScoresToBaselineNotBoolean|TestEvaluationRegressionFlagInComparisonJSON|TestBaselineStoresCommitOIDAndScoresJSONOnly|TestVerificationDebtWhenImplementationWithoutVerification|TestVerificationDebtClearsWhenVerifiedWithEvidence|TestPromotionGateRequiresStoredTestNotAgentClaim|TestEvaluationMissingBaselineFailClosed|TestPartialVerificationCountsAsDebt'

# --- S05 regression/reflection (14 named) ---
go test ./internal/domain/ -count=1 -run 'TestRecordRegressionFromEvaluationDefaultsCorrelated|TestRecordRegressionFromContradictedEffectDefaultsCorrelated|TestCorrelationAndContradictionNeverAutoSetCaused|TestLinkHypothesisUpgradesToHypothesizedNotCaused|TestSetAttributionCausedFailClosedWithoutEvidence|TestSetAttributionCausedFailClosedFromCorrelated|TestSetAttributionCausedRequiresConfirmedHypothesisAndEvidence|TestHasOpenRegressionFeedsApplyDeliberationTransition|TestResolveRegressionClearsHasOpenRegression|TestReflectionPersistsStructuredFieldsQueryable|TestReflectionEssayOnlyFailClosed|TestObservedRelationshipLinkWithConfidenceNoEvidence|TestCausalRelationshipFailClosedWithoutEvidence|TestUnknownAttributionFailClosed'

# --- S01 domain transition event ---
go test ./internal/domain/ -count=1 -run 'TestApplyDeliberationTransitionPersistsEvent|TestApplyDeliberationTransitionPlanMissing|TestApplyDeliberationTransitionRequiresIDs'

# --- S06 loop integration (14 named — now in internal/loop after S06) ---
go test ./internal/loop/... -count=1 -run 'TestLoopNextDeliberationSectionPresent|TestLoopNextPolicyInputsLiveQueries|TestLoopNextInvestigateEmphasizesUncertainties|TestLoopNextExecuteEmphasizesContextAndRelated|TestLoopNextVerifySurfacesVerificationDebt|TestLoopApplyUnknownWriteKeyFailsClosed|TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopApplyRegressionWriteAffectsPolicyInputs|TestLoopApplyDeliberationTransitionEvent|TestLoopApplyReplaySkipsDuplicateTransition|TestLoopStatusDeliberationFields|TestLoopStatusBlockedWhenBlockingUncertainty|TestLoopApplyNoPartialWritesOnValidationFailure|TestLoopRecentChangesNoFileBytes'

# --- P19 loop keepers (must stay green) ---
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'

# --- Store embed + migration + no-blob ---
go test ./internal/store/... -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestNoSourceContentColumns|TestDeliberationStateTableMigrated|TestChangeStoreRoundtrip|TestRegressionStoreRoundtrip|TestReflectionStoreRoundtrip|TestOutcomeStoreRoundtrip|TestBaselineStoreRoundtrip|TestCreateUncertaintyStoreRoundtrip|TestHypothesisStoreRoundtrip|TestBlockingCountSQL|TestDecisionReconsiderStoreRoundtrip'

# --- P17 seed keepers (portable graph regression) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
```

### Block B — P21 S01–S07 deltas (named keepers)

```bash
# --- S01 portable P20 seed (D-05 retired) ---
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'

# --- S02 retrieval + FTS (D-06, D-07 closed) ---
go test ./internal/retrieval/... -count=1 -run 'TestExactLookupUncertainty|TestExactLookupHypothesis|TestWhyUncertaintyIncludesGraphSteps|TestNormalizeEntityTypeP20Aliases'
go test ./internal/store/... -count=1 -run 'TestSyncEntityFTSUncertainty|TestRebuildFTSIncludesP20Types'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopNextInvestigateNoRetrievalStderr|TestCausalWhyContextRoundTrip|TestLoopNextInvestigateEmphasizesUncertainties'

# --- S03 full SelectNext cycle (D-03, D-04, D-18 closed) ---
go test ./internal/deliberation/... -count=1 -run 'TestSelectNextExecuteWhenPending|TestSelectNextTestWhenExecuteDone|TestSelectNextEvaluateWhenTestPass|TestSelectNextReflectWhenEvaluationDone|TestSelectNextReplanWhenFlagged|TestSelectNextExploreWhenAlternatives|TestSelectNextFullCycleOrdering|TestTransitionPayloadIncludesOptionalScores'

# --- S04 baseline promotion (D-09, D-10 closed) ---
go test ./internal/domain/... -count=1 -run 'TestPromoteBaselineSupersedesPrior|TestPromoteBaselineIdempotent|TestEvalRegressionBlocksPromotionGate|TestEvalRegressionGateClearsAfterResolve|TestPromotionGateIndependentOfTestPassAlone|TestBaselinePromotionRequiresStoredEvaluation'
go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestLoopApplyDeliberationTransitionEvent'

# --- S05 observability + why (D-11, D-12 closed) ---
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestCLIWhyUncertainty|TestCLIWhyRegression|TestWhyTaskIncludesDeliberationTransition|TestLoopNextHistoricalRelationshipsSection|TestHistoricalRelationshipsObservedVsCaused'

# --- S06 apply hardening (D-08, D-13, D-15 closed) ---
go test ./internal/loop/... -count=1 -run 'TestLoopApplyTransactionalRollbackOnFailure|TestLoopApplyGoalIDMismatchFailsClosed|TestApplyDeliberationTransitionRequiresMatchingGoalID|TestValidateApplyEnvelopeSpawnedTaskGoalMismatch|TestLoopApplySuccessPersistsLoopStepEvent'
go test ./internal/domain/... -count=1 -run 'TestApplyDeliberationTransitionRequiresMatchingGoalID|TestApplyDeliberationTransitionRequiresIDs'

# --- S07 promoted MVP cuts (D-01, D-02 closed) ---
go test ./internal/domain/... -count=1 -run 'TestCreateExperimentLinksOutcome|TestExperimentStatusLifecycle|TestNoExperimentRunnerInvoked'
go test ./internal/loop/... -count=1 -run 'TestRiskHintsManyPaths|TestRiskHintsBlockingUncertainty|TestLoopNextRiskHintsBounded'
```

### Block C — Compat + optional product bar

```bash
# --- Compat ceiling 21 (CGO required) ---
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# --- Product bar (optional broad; recommended once) ---
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

If any **named** test in the floor is absent at preflight, mark row **`failed`** with exact name — do not invent replacements.

## Required CLI evidence (archive)

Capture stdout/stderr from **ordinary CLI** (repo root or explicit `-C`). Minimum files under `evidence/`:

| # | Artifact | Proves |
|---|----------|--------|
| 1 | `01-seed-export-p20.json` | `trace seed export` includes P20 keys when DB populated; grep sidecar `01-seed-export-p20-grep.txt` |
| 2 | `02-full-cycle-next.json` | **≥5** distinct `deliberation.phase` values across scripted journey (ORIENT/INVESTIGATE/…/EXECUTE/VERIFY/EVALUATE) |
| 3 | `03-why-uncertainty.json` | `trace why` resolves P20 `uncertainty` type (exit 0 + steps) |
| 4 | `04-baseline-promote.json` | `PromoteBaseline` / eval gate / loop `promotion_blocked` after regression |
| 5 | `05-transactional-apply-fail.json` | apply mid-failure stderr; 0 partial rows (or reference unit test log) |
| 6 | `06-risk-hints.json` | `risk_hints.items` in loop next JSON (S07) |
| 7 | `07-historical-relationships.json` | `historical_relationships` section in loop next (S05) |
| 8 | `99-run-metadata.txt` | git SHA, go version, all command blocks + PASS/FAIL summary |

**P20 baseline artifacts** (still valuable — include or cross-reference P20 run):

| # | Artifact | Proves |
|---|----------|--------|
| — | `02-apply-blocking-uncertainty.json` | blocking uncertainty → INVESTIGATE, not EXECUTE |
| — | `04-apply-test-without-verify.json` | test pass alone; verification debt persists |
| — | `05-loop-status-deliberation.json` | `trace.loop.status.v1` deliberation fields |

Reference setup pattern: [P20 verify evidence](../../../../experiments/runs/2026-08-18-p20-s07-01-verify/evidence/) (goal → task → plan → loop next → apply → status).

## Must coverage checklist (COVERAGE.md → evidence)

S08-01 Notes must cite **PASS** evidence for each **Must** row (test name and/or CLI file). Inherit P20 S07-01 map; add P21 anchors:

| Doc § | Topic | P21 evidence anchor |
|------:|-------|---------------------|
| 1 | Externalized cognition loop | Block A S06 + Block B cycle tests + CLI #2 |
| 5 | Deterministic SelectNext | **14-row** table; `TestSelectNextFullCycleOrdering` |
| 13 | Baseline promotion | `TestPromoteBaselineSupersedesPrior`; CLI #4 |
| 16 | Experiments (thin promote) | `TestCreateExperimentLinksOutcome`; **not** bake-off engine |
| 17 | Historical relationships | CLI #7; `TestLoopNextHistoricalRelationshipsSection` |
| 18 | Risk-adaptive hints (minimal) | CLI #6; `TestRiskHintsManyPaths` |
| 25 | Observability / why | CLI #3; `TestCLIWhyUncertainty` |
| 29O | Failure modes | gates + tx rollback (Block A + B S06) |
| 29Q | Subsystem tests | named floor Blocks A–C |
| 31 | Desired outcome mini-eval | **live** CLI #2 — **≥5 phases** (D-14 promoted) |

**Should** items: include in Notes; §31 live mini-eval is **required** (not fixture-only acceptable).

### §29O failure-mode safeguards (must cite)

| Failure mode | Safeguard evidence |
|--------------|-------------------|
| Hallucinated findings | `TestPromotionGateRequiresStoredTestNotAgentClaim` |
| False verification | `TestVerificationMissingEvidenceFailClosed`; CLI test-without-verify |
| False causal attribution | `TestSetAttributionCausedFailClosed*` |
| Infinite deliberation loops | hop budget N=12; P19 saturation coexist |
| Incomplete evidence cannot promote | blocking uncertainty + verification debt + eval regression → INVESTIGATE/VERIFY/`promotion_blocked` |
| Partial apply writes | `TestLoopApplyTransactionalRollbackOnFailure` (D-08) |
| Wrong goal binding | `TestLoopApplyGoalIDMismatchFailsClosed` (D-13) |

## §31 live mini-eval spec (FINAL — not fixture-only)

Approximate [TRACE_THOUGHTPROCESS §31](../../../../TRACE_THOUGHTPROCESS.md) **through the loop CLI** on the **live Trace repo** with real `trace/` binary and temp `.trace/` (or isolated `TRACE_HOME`):

```text
1. seed import trace/graph.json (with P20 cognition populated — export from test DB or clone)
2. ORIENT / loop next — observe initial phase
3. apply blocking uncertainty → INVESTIGATE (not EXECUTE)
4. resolve/supersede uncertainty → advance toward PLAN/CRITIQUE
5. plan exists + critiqued → EXECUTE eligible when flags set (S03)
6. apply change + expected effect
7. apply test outcome (kind=test) — verification debt still visible
8. apply verification + evidence → debt clears
9. baseline + evaluation → comparison_json; regression → promotion_blocked
10. optional: promote baseline when gate clear
11. trace why uncertainty / task — P20 steps present
12. loop next shows risk_hints + historical_relationships when preconditions met
```

**Minimum mini-eval claims (Notes + evidence):**

1. At least **5** distinct deliberation phases in CLI `loop next` output across the scripted flow (stricter than P20's 3).
2. At least **1** blocking uncertainty prevented EXECUTE-phase recommendation.
3. Seed export/import round-trip includes P20 keys (test or CLI #1).
4. Full cycle phase reachable: EXECUTE and at least one of TEST/VERIFY/EVALUATE observed when pending flags set.
5. Baseline promotion or `promotion_blocked` after eval regression (CLI #4).
6. `trace why` on P20 entity succeeds (CLI #3).
7. `risk_hints` and `historical_relationships` sections present when preconditions met (CLI #6–7).

**Required workspace:** **Live repo** — build `trace` from this checkout; use temp DB directory. **Not acceptable:** fixture-only `.trace/` without documenting live CLI path as primary evidence.

**Residual risk (if any step scripted via domain helpers):** one-line note — long-horizon multi-agent §31 "future agent receives knowledge" still depends on seed clone workflow in the wild.

## Do not

- Close [`DR-HANDOFF.md`](../../DR-HANDOFF.md) — S08-02 owns
- Mark human-gated criteria `done` without evidence files
- Revert to P20 seed omit policy (D-05 retired)
- Add migration 022+ or compat ceiling >21
- Scaffold Phase 22

## Exit criteria

- [ ] Evidence directory populated (CLI + test command log)
- [ ] Block A (P20 floor) + Block B (P21 deltas) + Block C compat reported PASS (or row `failed` with reason)
- [ ] Must checklist mapped in board Notes with evidence pointers
- [ ] §31 **live** mini-eval completed — **≥5 phases** in CLI evidence
- [ ] Seed export: P17 keepers PASS + P20 keys documented (D-05 closed)
- [ ] Compat ceiling **21** confirmed
- [ ] Residual risks listed (experiments not in seed, Trace does not run tests autonomously D-16, hosted MCP out D-19)
- [ ] DR-HANDOFF remains **OPEN**

## Minimal todos

- [ ] Preflight: confirm all named tests exist (live grep)
- [ ] Run Block A + B + C; capture results in `99-run-metadata.txt`
- [ ] Script §31 live CLI flow; archive under `experiments/runs/YYYY-MM-DD-p21-s08-01-verify/`
- [ ] Run `trace seed export -o trace/graph.json`; document P20 keys
- [ ] Write VERIFY-NOTES.md (recommended) + board Notes with Must map
- [ ] Set row `done` or `failed` — **do not** close DR-HANDOFF
