# P20-S07-01 VERIFY-NOTES

**Date:** 2026-08-18  
**Evidence dir:** `experiments/runs/2026-08-18-p20-s07-01-verify/evidence/`  
**Row outcome:** **failed** — locked floor `s03-changes` + optional product bar fail on stale S03 keeper (see below). All other layers PASS.

## Pass/fail per layer

| Layer | Result | Evidence |
|-------|--------|----------|
| S01 deliberation (5) | **PASS** | `test-s01-deliberation.log` |
| S02 cognitive (14) | **PASS** | `test-s02-cognitive.log` |
| S03 changes (14) | **FAIL** | `test-s03-changes.log` — `TestContradictedEffectDoesNotCreateRegressionOrAutoReplan` |
| S04 gates (14) | **PASS** | `test-s04-gates.log` |
| S05 regression/reflect (14) | **PASS** | `test-s05-regression.log` |
| S01 domain transition (3) | **PASS** | `test-s01-domain-trans.log` |
| S06 loop locked path (`./internal/loop/...`) | **VACUOUS PASS** | `test-s06-loop-locked-path.log` — no test files in package |
| S06 loop corrected (`./cmd/trace`) | **PASS** | `test-s06-loop-corrected.log` — 14/14 |
| P19 loop keepers (6) | **PASS** | `test-p19-loop.log` |
| Store embed + P20 roundtrips (13) | **PASS** | `test-store.log` — embed max **019** |
| Compat ceiling 19 | **PASS** | `test-compat-rerun.log` |
| P17 seed keepers (3) | **PASS** | `test-p17-seed-rerun.log` |
| Product bar (optional) | **FAIL** | `test-product-bar-rerun.log` — same S03 stale test |

### S03 failure analysis (not product regression)

`TestContradictedEffectDoesNotCreateRegressionOrAutoReplan` asserts `sqlite_master` has **no** `regressions` **table**. Mig **019** (S05) creates that table globally; the test passed at S03 review time. Intended invariant (no auto regression **row** on contradicted effect) remains covered by `TestContradictedEffectFiresDecisionReconsideration` + S05 create-always-`correlated` tests. **Forward fix:** update S03 test to assert zero regression rows / no auto `RecordRegressionFromContradictedEffect`, not table absence.

### S06 path note

Locked floor lists `./internal/loop/...`; named S06 tests live in `cmd/trace/loop_test.go` (library tests import loop). Authoritative run: `./cmd/trace` (14/14 PASS).

## CLI evidence (ordinary)

| # | File | Proves |
|---|------|--------|
| 1 | `01-loop-next-deliberation.json` | `trace.loop.next.v1`; `deliberation.phase=CRITIQUE`, `policy_inputs`, sections |
| 2 | `02-apply-blocking-uncertainty.json` | BLOCKING uncertainty → `INVESTIGATE` / `blocking_uncertainty_count=1` (not EXECUTE) |
| 3 | `03-apply-change-effect-contradiction.json` | change + `comparison=contradicted` in `recent_changes`; verify debt |
| 4 | `04-apply-test-without-verify.json` | test pass; `verification_debt.present=true`, phase VERIFY |
| 5 | `05-loop-status-deliberation.json` | `trace.loop.status.v1`; `blocked=true`, `verification_incomplete` |
| 6 | `06-seed-export-sample.json` | portable export; P20 tables omitted (grep + key audit) |
| 7 | `99-run-metadata.txt` | commands, SHA, summary |

Fixture projects: `project/` (01–02, 05), `project-change/` (03), `project-verify-debt/` (04, §31 mini-eval).

**CLI residual:** `loop next` stderr may log `retrieval: unknown entity type "uncertainty"` while stdout still emits degraded packet when blocking/open regression (S06 documented residual). Zero-delta apply without `spawned_tasks` can saturate P19 step early; evidence scripts use spawned task where needed.

## §31 mini-eval (fixture-scale)

| Claim | Evidence |
|-------|----------|
| ≥3 deliberation phases | CLI: CRITIQUE (`01`), INVESTIGATE (`02`), VERIFY (`04`), INVESTIGATE+open_regression (`mini-eval-after-regression-next.json`) |
| Blocking uncertainty blocks EXECUTE | `02` — phase INVESTIGATE, count=1 |
| Expected vs actual effect | `03` — contradicted in `recent_changes.effects` |
| Test pass ≠ verify clearance | `04` — debt present after test-only apply |
| Regression correlated not caused | `mini-eval-baseline-regression.txt` — `attribution=correlated` |
| Status reflects deliberation/blocked | `05`, `mini-eval-final-status.json` |

**Fixture residual:** live multi-agent variance + long-horizon “future agent receives knowledge” unproven until seed export extends to P20 tables.

## Seed export omit policy + residual

**Policy (locked S07-00):** P20 SQLite tables **not** in `BuildSeedDocument` — omit extension this phase.

**Omitted tables:** `deliberation_state`, `uncertainties`, `hypotheses`, `decision_reconsiderations`, `changes`, `change_paths`, `effects`, `outcome_results`, `baselines`, `regressions`, `reflections`.

**Evidence:** P17 keepers PASS (`test-p17-seed-rerun.log`); `06-seed-export-sample.json` top-level keys lack all P20 entities (`uncertainties`/`deliberation_state` false).

**Forward residual (not Phase 20 fail):** portable clone of P20 cognition artifacts requires human-promoted seed-export extension.

## COVERAGE Must checklist → evidence

| § | Topic | PASS | Anchor |
|---|-------|------|--------|
| 1 | Externalized cognition loop | yes | S06 cmd/trace tests + CLI 01–05 |
| 2 | No raw CoT | yes | `TestNoSourceContentColumns`; mig 015–019 |
| 3A | Thin engineering links | yes | S02 hypothesis tests; S03 change_paths |
| 3B | Deliberation transitions | yes | `TestSelectNext*`, `TestApplyDeliberationTransition*` |
| 4 | Phases ORIENT…REPLAN | yes | CLI `deliberation.phase` 01–04 |
| 5 | Deterministic SelectNext | yes | `TestSelectNext` |
| 6 | Entry/exit + hop budget N=12 | yes | `TestApplyTransitionHopBudget*` |
| 7 | Uncertainty first-class | yes | S02 tests + CLI 02 |
| 8 | Assumption invalidate | yes | `TestInvalidateAssumption*` |
| 9 | Decisions + reconsider | yes | `TestDecisionReconsider*` |
| 10 | Change first-class | yes | S03 tests + CLI 03 |
| 11 | Test ≠ Verify ≠ Evaluate | yes | `TestTestPassAloneCannotSatisfyVerificationGate` |
| 12 | Gates; implemented ≠ verified | yes | `TestPromotionGateRequiresStoredTestNotAgentClaim` |
| 14 | Expected vs actual | yes | S03 contradiction tests + CLI 03 |
| 15 | Regression correlated ≠ caused | yes | `TestCorrelationAndContradictionNeverAutoSetCaused` + mini-eval |
| 19 | Reflection (thin) | yes | `TestReflectionPersistsStructuredFieldsQueryable` + reflection apply |
| 20 | Verification debt | yes | `TestVerificationDebt*` + CLI 04 |
| 21 | Feedback state machine | yes | SelectNext + apply transition |
| 22 | Agent interaction | yes | allowlisted writes; `TestLoopApplyUnknownWriteKeyFailsClosed` |
| 23 | Context by phase | yes | `TestLoopNextInvestigate*` / Execute / Verify |
| 24 | Harness-agnostic | yes | stdout JSON CLI |
| 25 | Observability / why | yes | `why_selected` + transition events (S06 test) |
| 26 | Avoid premature complexity | yes | COVERAGE Future rows not implemented; no mig 020 |
| 28 | Reuse P19 loop | yes | P19 keeper tests PASS |
| 29O | Failure modes | yes | gates fail-closed tests + CLI 04 |
| 29Q | Subsystem tests | **partial** | S03 floor FAIL (stale test); others PASS |
| 30 | Challenge concept | yes | README MVP cuts honored |
| 32 | Incremental MVP | yes | §16/§18 Future; ceiling 19 |

**Should:** §31 mini-eval fixture done; baselines/observed-relationship covered by domain tests.

## Residual risks (listed)

1. **Seed export** — P20 cognition not portable (forward residual).
2. **FTS sync** on apply-created entities (S06/S07 note).
3. **Non-tx upserts** (domain apply paths).
4. **Retrieval/compiler** lacks `uncertainty` entity — degraded context/why on INVESTIGATE (stderr noise).
5. **Stale S03 keeper** after mig 019 — blocks green product bar until forward test fix.

## DR-HANDOFF

**Not closed** — S07-02 owns review + handoff close policy.
