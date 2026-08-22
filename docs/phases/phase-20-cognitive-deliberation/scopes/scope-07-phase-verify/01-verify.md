# P20-S07-01 — Verify Phase 20

## Metadata
- id: P20-S07-01
- todo_ids: [P20-S07-01]
- role: verify
- skills: [writing-for-agents, planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Run the **locked verify floor** (named keeper tests + compat ceiling **19** + P17 seed keepers), archive **ordinary CLI** deliberation evidence under `experiments/runs/YYYY-MM-DD-p20-s07-01-verify/`, and map results to **COVERAGE.md Must** items + **§31** mini-eval. **Does not** close DR-HANDOFF (S07-02 owns).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [COVERAGE.md](../../COVERAGE.md)
- [TRACE_THOUGHTPROCESS.md §31](../../../../TRACE_THOUGHTPROCESS.md) (Desired Outcome)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar (this scope)
- Pattern: [P19 S03-01 verify](../../../phase-19-loop-gap-detection/scopes/scope-03-phase-verify/01-verify.md)

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row verifies and records evidence; it does not open new product direction beyond Phase 20 locks.

## Locked defaults (FINAL — S07-00)

| Item | Value |
|------|-------|
| Schema max | **019** (`015`–`019` present; **no 020+**) |
| Compat ceiling | **19** — re-lock at VERIFY time via `TestCompatibilitySecurityChecklist` |
| Loop schemas | **additive v1** unchanged: `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` |
| CLI surface | `trace loop next --task <id>`, `trace loop apply [--in <path>]`, `trace loop status --task <id> [--goal <id>]` |
| Hop budget | **N=12** (`internal/deliberation`) |
| Seed export (P20 entities) | **Omit extension** this phase — P17 keepers must PASS; list omitted tables in Notes as **residual** |
| Pre-PR export | Run `trace seed export -o trace/graph.json` once; note P20 tables absent from JSON (expected) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p20-s07-01-verify/evidence/` |
| Notes artifact | Optional `VERIFY-NOTES.md` in this scope folder (recommended) |
| Out-of-scope | No product Go; no DR-HANDOFF close; no Phase 21 scaffold |

### Seed export bar (locked)

P20 added SQLite tables **not** in `BuildSeedDocument` (`internal/domain/seed_export.go`):

- `deliberation_state` (S01)
- `uncertainties`, `hypotheses`, `decision_reconsiderations` (S02)
- `changes`, `change_paths`, `effects` (S03)
- `outcome_results`, `baselines` (S04)
- `regressions`, `reflections` (S05)

**VERIFY requires:**

1. P17 keepers green (no regression on portable graph):
   ```bash
   CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
   ```
2. One `trace seed export -o trace/graph.json` (or temp path) captured in evidence; grep/inspect confirms **no** keys for tables above.
3. Notes explicitly state: portable clone of P20 cognition artifacts is a **forward residual** (human-promoted follow-on), **not** a Phase 20 fail.

## Locked verify command floor (FINAL)

S07-01 must run and report PASS/FAIL for **every** block below. `-count=1` recommended. `GOMODCACHE`+`GOPROXY=off` on full product bar if offline.

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

# --- S06 loop integration (14 named) ---
go test ./internal/loop/... -count=1 -run 'TestLoopNextDeliberationSectionPresent|TestLoopNextPolicyInputsLiveQueries|TestLoopNextInvestigateEmphasizesUncertainties|TestLoopNextExecuteEmphasizesContextAndRelated|TestLoopNextVerifySurfacesVerificationDebt|TestLoopApplyUnknownWriteKeyFailsClosed|TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopApplyRegressionWriteAffectsPolicyInputs|TestLoopApplyDeliberationTransitionEvent|TestLoopApplyReplaySkipsDuplicateTransition|TestLoopStatusDeliberationFields|TestLoopStatusBlockedWhenBlockingUncertainty|TestLoopApplyNoPartialWritesOnValidationFailure|TestLoopRecentChangesNoFileBytes'

# --- P19 loop keepers (must stay green) ---
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'

# --- Store embed + migration + no-blob ---
go test ./internal/store/... -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestNoSourceContentColumns|TestDeliberationStateTableMigrated|TestChangeStoreRoundtrip|TestRegressionStoreRoundtrip|TestReflectionStoreRoundtrip|TestOutcomeStoreRoundtrip|TestBaselineStoreRoundtrip|TestCreateUncertaintyStoreRoundtrip|TestHypothesisStoreRoundtrip|TestBlockingCountSQL|TestDecisionReconsiderStoreRoundtrip'

# --- Compat ceiling 19 (CGO required) ---
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# --- P17 seed keepers (portable graph regression) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'

# --- Product bar (optional broad; recommended once) ---
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

If any **named** test in the floor is absent at preflight, mark row **`failed`** with exact name — do not invent replacements.

## Required CLI evidence (archive)

Capture stdout/stderr from **ordinary CLI** (repo root or explicit `-C`). Minimum files:

| # | Artifact | Proves |
|---|----------|--------|
| 1 | `01-loop-next-deliberation.json` | `trace.loop.next.v1`; `deliberation.phase`, `why_selected`, `policy_inputs`, bounded sections |
| 2 | `02-apply-blocking-uncertainty.json` | apply `uncertainty` (BLOCKING) → next shows INVESTIGATE / `blocking_uncertainty`; **not** EXECUTE |
| 3 | `03-apply-change-effect-contradiction.json` | apply `change`+`effect` (contradicted) → regression signal or INVESTIGATE/`open_regression` on next |
| 4 | `04-apply-test-without-verify.json` | test pass alone; next still shows `verification_debt` / VERIFY phase when debt open |
| 5 | `05-loop-status-deliberation.json` | `trace.loop.status.v1` with deliberation fields / blocked when blocking uncertainty |
| 6 | `06-seed-export-sample.json` | `trace seed export` output; confirms P20 tables omitted |
| 7 | `99-run-metadata.txt` | git SHA, commands run, pass/fail summary |

Reference setup pattern: [P19 verify evidence](../../../experiments/runs/2026-08-18-p19-s03-01-verify/evidence/) (goal → task → plan → loop next → apply → status).

## Must coverage checklist (COVERAGE.md → evidence)

S07-01 Notes must cite **PASS** evidence for each **Must** row (test name and/or CLI file):

| Doc § | Topic | Evidence anchor |
|------:|-------|-----------------|
| 1 | Externalized cognition loop | S06 loop tests + CLI next/apply/status |
| 2 | No raw CoT | store `TestNoSourceContent`; no blob columns in 015–019 |
| 3A | Thin engineering links | S02 hypothesis/evidence rel; S03 change_paths |
| 3B | Deliberation transitions | `TestSelectNext*` + `TestApplyDeliberationTransition*` |
| 4 | Cognitive phases ORIENT…REPLAN | `deliberation.phase` in CLI #1 |
| 5 | Deterministic SelectNext | `TestSelectNext` table (≥8 rows) |
| 6 | Entry/exit + hop budget | `TestApplyTransitionHopBudget*`; N=12 |
| 7 | Uncertainty first-class | S02 named tests + CLI #2 |
| 8 | Assumption invalidate | `TestInvalidateAssumption*` |
| 9 | Decisions + reconsider | `TestDecisionReconsider*` |
| 10 | Change first-class | S03 named tests + CLI #3 |
| 11 | Test ≠ Verify ≠ Evaluate | `TestTestPassAloneCannotSatisfyVerificationGate` |
| 12 | Gates; implemented ≠ verified | `TestPromotionGateRequiresStoredTestNotAgentClaim` |
| 14 | Expected vs actual | `TestRecordExpectedThenActualSupported`; contradiction tests |
| 15 | Regression correlated ≠ caused | `TestCorrelationAndContradictionNeverAutoSetCaused` |
| 19 | Reflection (thin) | `TestReflectionPersistsStructuredFieldsQueryable` |
| 20 | Verification debt | `TestVerificationDebt*` + `TestLoopNextVerifySurfacesVerificationDebt` |
| 21 | Feedback state machine | SelectNext priority + loop apply transition |
| 22 | Agent interaction | loop apply allowlisted writes |
| 23 | Context by phase | `TestLoopNextInvestigate*` / `Execute*` / `Verify*` |
| 24 | Harness-agnostic | stdout JSON CLI evidence |
| 25 | Observability / why | `why_selected` + `deliberation.transition` event |
| 26 | Avoid premature complexity | COVERAGE Future rows **not** implemented |
| 28 | Reuse P19 loop | P19 keeper tests PASS |
| 29O | Failure modes | gates fail-closed (§29O bullets below) |
| 29Q | Subsystem tests | named floor above |
| 30 | Challenge concept | README MVP cuts honored |
| 32 | Incremental MVP | §16/§18 Future; no mig 020 |

**Should** items (13 baselines, 17 observed relationships, 31 mini-eval): include in Notes; mini-eval required below.

### §29O failure-mode safeguards (must cite)

| Failure mode | Safeguard evidence |
|--------------|-------------------|
| Hallucinated findings | Agent claim without stored row fails gates (`TestPromotionGateRequiresStoredTestNotAgentClaim`) |
| False verification | `TestVerificationMissingEvidenceFailClosed`; CLI #4 |
| False causal attribution | `TestSetAttributionCausedFailClosed*`; create always `correlated` |
| Infinite deliberation loops | hop budget N=12; P19 saturation coexist |
| Incomplete evidence cannot promote | verification debt + blocking uncertainty → INVESTIGATE/VERIFY, not EXECUTE/DONE |

## §31 mini-eval spec (fixture-scale acceptable)

Approximate the §31 story **through the loop CLI** on one seeded goal/task with plan context:

```text
ORIENT/INVESTIGATE → record open uncertainty (apply)
INVESTIGATE → resolve or supersede after “finding” (apply)
PLAN/CRITIQUE → plan exists + critiqued (existing planner or apply plan_change)
EXECUTE path blocked while BLOCKING uncertainty open (next packet)
CHANGE → record change + expected effect (apply)
TEST → record test outcome kind=test (apply)
VERIFY → verification debt visible until verification+evidence (apply + next)
EVALUATE → baseline + evaluation comparison_json (library/apply)
contradiction OR regression → INVESTIGATE/open_regression (apply + next)
REFLECT → structured reflection arrays (apply)
```

**Minimum mini-eval claims (Notes + evidence):**

1. At least **3** distinct deliberation phases observed in CLI `loop next` output across the scripted flow.
2. At least **1** blocking uncertainty prevented EXECUTE-phase recommendation (CLI #2 or equivalent).
3. At least **1** change with expected vs actual effect recorded (supported or contradicted).
4. Test pass **without** verification row did **not** clear verification debt (CLI #4).
5. Contradiction or evaluation regression produced **correlated** attribution (never auto-`caused`).
6. `trace loop status` reflects deliberation/blocked state at least once.

**Acceptable paths:**

- **Preferred:** continuation on live Trace repo with real goal/task/plan (mirror P19 taskboard-style setup).
- **Allowed:** fixture-scale `.trace/` + seeded entities via CLI/domain helpers scripted into evidence dir.

If fixture-scale: one-line **residual risk** — live multi-agent variance and long-horizon §31 “future agent receives knowledge” remain unproven until seed export extends to P20 tables.

## Do not

- Close [`DR-HANDOFF.md`](../../DR-HANDOFF.md) — S07-02 owns
- Mark human-gated criteria `done` without evidence files
- Extend seed export to P20 tables (forward residual only)
- Add migration 020 or compat ceiling >19

## Exit criteria

- [ ] Evidence directory populated (CLI + test command log)
- [ ] All locked verify floor commands reported PASS (or row `failed` with reason)
- [ ] Must checklist mapped in board Notes with evidence pointers
- [ ] §31 mini-eval completed or fixture residual stated
- [ ] Seed export: P17 keepers PASS + P20 omission documented
- [ ] Compat ceiling **19** confirmed
- [ ] Residual risks listed (seed export, FTS sync on apply-created entities, non-tx upserts)

## Minimal todos

- [ ] Preflight: confirm all named tests exist (live grep)
- [ ] Run locked command floor; capture results in `99-run-metadata.txt`
- [ ] Script CLI evidence flow; archive under `experiments/runs/…`
- [ ] Run `trace seed export`; document omitted P20 entities
- [ ] Write VERIFY-NOTES.md (recommended) + board Notes with Must map
- [ ] Set row `done` or `failed` — **do not** close DR-HANDOFF
