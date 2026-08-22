# P20-S04-00 — Verify / evaluate / gates planner

## Metadata
- id: P20-S04-00
- todo_ids: [P20-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Test vs Verification vs Evaluation as **distinct result kinds**, promotion gates, thin Baseline, verification debt. **Not** a full test runner product. No product Go this row.

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [COVERAGE.md](../../COVERAGE.md) §§11, 12, 13, 20, merge table (Test/Verification/Evaluation, Baseline)
- Laws **2**, **3**, **14**, **15** ([G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md))
- Live: schema max **017** → next migration **`018_outcome_results_baselines.sql`**
- Live S01: `PolicyInputs.VerificationIncomplete` + SelectNext row 7 → `VERIFY` / `verification_incomplete`
- Live S03: `changes`/`change_paths`/`effects` — **no** tests/baseline/score columns (S03-02 APPROVE)
- Live DONE: Review PASS + `AllowOperatorDone`; evidence IDs alone do **not** authorize DONE (`internal/domain/doc.go`)
- §16 Experiments + §18 Risk-adaptive verification → **Future** (not S04 implement)

## Doc map
§11, 12, 13, 20, 29F–H

## Live inventory (2026-08-18)

| Surface | Location | S04 action |
|---------|----------|------------|
| `tasks.work_state` | PENDING…DONE (no new values) | **Do not** add `implemented`/`verified`/`evaluated`/`promotable`. Gates via linked `outcome_results` + deliberation inputs. |
| `goals` | `001_init.sql` | Verification **requires** `goal_id` (Requirement merged into Goal body per COVERAGE — no Requirement table). |
| `evidence` | `001_init.sql` + `claim_has_evidence` | Verification **requires** ≥1 `outcome_supported_by` link → evidence ID. |
| `reviews` + DONE policy | `005_review_promotion.sql`, `task_state.go` | **Unchanged.** DONE = Review PASS + operator flag (or hatch). Verification debt is **not** a silent DONE bypass. |
| `changes` | `017_changes_effects.sql` | Implementation-complete **signal** for debt query only. **No** new columns. |
| `deliberation` SelectNext | `internal/deliberation/select.go` | S04 **does not** edit. S06 wires `VerificationIncomplete` from S04 query helper. |
| Test runner product | **missing** | **Out.** Library records outcomes; does not embed pytest/go test/jest. |
| Baselines / outcomes | **missing** | S04-01 adds `baselines` + `outcome_results`. |
| Next migration | max **017** | S04-01 adds **`018_outcome_results_baselines.sql`** |
| Compat ceiling | **17** after S03 | S04-01 bumps to **18** (forbid `019+`) |

## Paths storage / table naming (FINAL)

| Option | Verdict |
|--------|---------|
| Three tables (`tests`, `verifications`, `evaluations`) | **Rejected** — duplicates task_id/actor/timestamps; harder debt query |
| `verification_results` with kind enum | **Rejected** — name implies verify-only; kind=`test`/`evaluation` misleading |
| **`outcome_results`** with `kind` discriminator + thin **`baselines`** | **Locked** |

S06 apply keys (`test_results`, `verifications`, `evaluations`) are **transport aliases** that all persist into `outcome_results` with the matching `kind`. S04-01 is library-only — S06 owns loop apply wiring.

## FINAL locked defaults (S04-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`018_outcome_results_baselines.sql`** — additive; do not rewrite 001–017; **no ALTER** on `changes` / `tasks.work_state` / S02 tables |
| Compat ceiling | **18** after S04-01 (forbid `019+`); bump `evals/compat`, `production_hardening_test`, `deliberation_test` EmbedExpected, `TestOpenCreatesDBAndMigratesIdempotent` |
| New tables | `outcome_results`, `baselines` only |
| Forbidden | Experiment tables; risk-adaptive matrix; test-runner subprocess product; `score_*` / `tests` / `verification_runs` / `baseline` columns on `changes`; new `work_state` enum values; raw CoT / full log blobs |
| Result kinds | `test` \| `verification` \| `evaluation` — single table, `kind` CHECK |
| Test | Named `test_name` + `test_status` (`pass`\|`fail`\|`skip`\|`error`) + bounded `summary`. **Agent claim alone insufficient** — gate helpers require a stored row (Law 2). |
| Verification | Requires `goal_id` + ≥1 linked evidence ID (`outcome_supported_by`). `test` pass **≠** verification satisfied (Law 2). Status `verified`\|`failed`\|`partial` — not boolean Review PASS. |
| Evaluation | Numeric `scores_json` compared to `baselines.scores_json` → library-computed `comparison_json` with per-dimension deltas + `regression` flags. **Not** a boolean PASS/FAIL outcome (Law 15 for automatable deltas). |
| Baseline | Thin row: `git_commit` OID + `scores_json` only (+ label/metadata). No embedded test runner output blobs. |
| Verification debt | First-class query for S06 packet + S01 `verification_incomplete` input (see below). |
| DONE transition | **Unchanged:** Review PASS + `AllowOperatorDone` (or hatch). Debt blocks via deliberation → VERIFY phase, **not** by silently skipping Review gate. |
| work_state | **No new values.** Optional `BLOCKED` transition remains caller-driven; prefer debt flag over enum explosion. |
| §16 / §18 | **Future** — no experiment runner, no risk-adaptive test-selection engine; optional code comment stub only |
| CLI / MCP / loop / SelectNext | **Library-only** — S06 owns apply keys and packet `verification_debt` |

### `baselines`

```sql
CREATE TABLE IF NOT EXISTS baselines (
    id TEXT PRIMARY KEY,
    git_commit TEXT NOT NULL,
    scores_json TEXT NOT NULL DEFAULT '{}',
    label TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_baselines_git_commit ON baselines(git_commit);
```

| Field | Lock |
|-------|------|
| `git_commit` | Git OID only. Must match `^[0-9a-fA-F]{7,64}$`, store lowercase. Reference evidence — need not exist in `vcs_commits`. |
| `scores_json` | JSON **object** (map string→number or string→string for non-numeric dims). Max **16384** bytes. Fail closed on invalid JSON / array root / empty `{}` at create. |
| `label` | Human/agent label (e.g. `B100`). Optional empty allowed. Max **256** bytes. |
| History | Law 11 — no `DeleteBaseline`. Supersede via new baseline row + pointer from evaluation. |

### `outcome_results`

```sql
CREATE TABLE IF NOT EXISTS outcome_results (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('test', 'verification', 'evaluation')),
    test_name TEXT NOT NULL DEFAULT '',
    test_status TEXT NOT NULL DEFAULT ''
        CHECK (test_status IN ('', 'pass', 'fail', 'skip', 'error')),
    goal_id TEXT NOT NULL DEFAULT '',
    verification_status TEXT NOT NULL DEFAULT ''
        CHECK (verification_status IN ('', 'verified', 'failed', 'partial')),
    baseline_id TEXT NOT NULL DEFAULT '',
    scores_json TEXT NOT NULL DEFAULT '{}',
    comparison_json TEXT NOT NULL DEFAULT '{}',
    summary TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_outcome_results_task_id ON outcome_results(task_id);
CREATE INDEX IF NOT EXISTS idx_outcome_results_kind ON outcome_results(task_id, kind);
CREATE INDEX IF NOT EXISTS idx_outcome_results_goal_id ON outcome_results(goal_id);
CREATE INDEX IF NOT EXISTS idx_outcome_results_baseline_id ON outcome_results(baseline_id);
```

Kind-specific column usage (fail-closed cross-kind pollution):

| kind | Required non-empty | Must be empty / ignored |
|------|-------------------|-------------------------|
| `test` | `test_name`, `test_status` ∈ {pass,fail,skip,error} | `goal_id`, `baseline_id`, `verification_status`, `scores_json`, `comparison_json` (store `{}`) |
| `verification` | `goal_id`, `verification_status` ∈ {verified,failed,partial}, ≥1 `outcome_supported_by` → evidence | `test_name`, `test_status`, `baseline_id`, `scores_json`, `comparison_json` |
| `evaluation` | `baseline_id`, `scores_json` (valid object); library fills `comparison_json` | `test_name`, `test_status`, `goal_id`, `verification_status` |

| Field | Lock |
|-------|------|
| `summary` | Bounded excerpt (command, exit code, one-line result). Max **4096** bytes. **Not** a full test log blob (Law 1 spirit). |
| `scores_json` | Evaluation input scores. Max **16384** bytes. Valid JSON object. |
| `comparison_json` | **Library-computed only** on evaluation create/update. Per-dimension `{baseline, current, delta, regression}` + optional `overall_regression`. Agent must not supply boolean PASS. |
| `task_id` | Required. Task must exist. |
| `goal_id` (verification) | Must match an existing goal. For task-scoped verification, should equal `tasks.goal_id` (fail closed if task has goal_id and they differ). |

### Entity links (new rel)

| rel | from → to | When |
|-----|-----------|------|
| `outcome_supported_by` | outcome_result → evidence | Required for `kind=verification` (≥1). Optional for `kind=test` (e.g. junit artifact as evidence ref, not blob). |

No duplicate `goal_has_outcome` link — use `goal_id` column on verification rows.

Constants on `internal/domain/service.go`: `EntityOutcomeResult`, `EntityBaseline`, `RelOutcomeSupportedBy`, events `outcome.recorded`, `baseline.created`, `evaluation.compared`.

### Promotion gates (library helpers — not new work_state)

Doc §12 flow as **check functions** (fail-closed booleans + reasons). **Not** a subprocess test runner.

```text
HasImplementationSignal(task)  → ∃ change status ∈ {RECORDED, COMPARED}
CheckTestGate(task, testName)   → ∃ kind=test, matching test_name, test_status=pass
CheckVerificationGate(task)     → ∃ kind=verification, verification_status=verified, goal_id set, ≥1 evidence link
CheckEvaluationGate(task, baselineID) → ∃ kind=evaluation for baseline with comparison_json computed
```

**Critical Law 2 lock:** `CheckTestGate` true **does not** imply `CheckVerificationGate` true. Named test must prove this.

**DONE / Review:** Unchanged. Gates inform deliberation and debt; they do **not** replace Review PASS for `TransitionTask → DONE`.

### Evaluation vs baseline (not boolean)

`CompareScoresToBaseline(current, baseline)` returns structured comparison:

```json
{
  "baseline_id": "uuid",
  "git_commit": "abc123…",
  "dimensions": {
    "correctness": {"baseline": 0.98, "current": 0.95, "delta": -0.03, "regression": true},
    "performance_p95_ms": {"baseline": 310, "current": 280, "delta": -30, "regression": false}
  },
  "overall_regression": true
}
```

- Numeric dimensions: regression when `delta` worsens per dimension policy (lower-is-better vs higher-is-better inferred from sign convention documented in domain — default **higher score is better** for 0..1 metrics; latency ms **lower is better** when key contains `_ms` or `latency`).
- Missing baseline dimension: `"regression": false`, `"delta": null`, note in comparison — do not invent PASS.
- S05 consumes `overall_regression` / per-dimension flags — S04 **does not** create `regressions` table.

### Verification debt query (S06 + S01 input)

**Function (domain + store):**

- `HasVerificationDebt(ctx, taskID) (bool, error)`
- `CountTasksWithVerificationDebt(ctx) (int, error)` optional
- `ListVerificationDebtSummary(ctx, taskID) ([]DebtItem, error)` bounded for packets (Law 6)

**Debt when ALL hold:**

1. `HasImplementationSignal(taskID)` — ≥1 `changes` row for task with `status IN ('RECORDED','COMPARED')`.
2. Task has non-empty `goal_id`.
3. **No** satisfactory verification: missing `outcome_results` row with `kind='verification'` AND `verification_status='verified'` AND `goal_id = task.goal_id` AND ≥1 `outcome_supported_by` evidence link.
4. `verification_status='partial'` counts as debt (not satisfied).

**S01 wiring (S06, not S04):** `PolicyInputs.VerificationIncomplete = HasVerificationDebt(taskID)`.

**Packet field (S06):** `verification_debt` = bounded JSON from `ListVerificationDebtSummary` (missing dimensions as strings, e.g. `"performance benchmark"`, `"integration evidence"` — caller-supplied labels in verification body/summary allowed; MVP may return generic `"verification missing for goal"`).

### Domain API (S04-01)

| API | Role |
|-----|------|
| `CreateBaseline(ctx, BaselineInput)` | git_commit + scores_json |
| `GetBaseline(ctx, id)` | |
| `RecordTestOutcome(ctx, TestOutcomeInput)` | kind=test; fail closed without test_name |
| `RecordVerificationOutcome(ctx, VerificationOutcomeInput)` | goal_id + evidenceIDs + status; inserts links |
| `RecordEvaluationOutcome(ctx, EvaluationOutcomeInput)` | baseline_id + scores_json; computes comparison_json |
| `CompareScoresToBaseline(current, baseline)` | pure helper |
| `CheckTestGate` / `CheckVerificationGate` / `CheckEvaluationGate` | gate helpers |
| `HasVerificationDebt` / `ListVerificationDebtSummary` | S06 packet + SelectNext input |
| `HasImplementationSignal` | debt precondition |

Inputs mirror S03 style — see 01-implement prompt.

### Named tests (minimum — exact names)

1. `TestRecordTestOutcomeRequiresNameAndStatus`
2. `TestTestPassAloneCannotSatisfyVerificationGate`
3. `TestVerificationRequiresGoalAndEvidenceIDs`
4. `TestVerificationMissingEvidenceFailClosed`
5. `TestEvaluationComparesScoresToBaselineNotBoolean`
6. `TestEvaluationRegressionFlagInComparisonJSON`
7. `TestBaselineStoresCommitOIDAndScoresJSONOnly`
8. `TestVerificationDebtWhenImplementationWithoutVerification`
9. `TestVerificationDebtClearsWhenVerifiedWithEvidence`
10. `TestPromotionGateRequiresStoredTestNotAgentClaim`
11. `TestOutcomeResultsSchemaNoBlobColumns`
12. `TestUnknownOutcomeKindFailClosed`
13. `TestEvaluationMissingBaselineFailClosed`
14. `TestPartialVerificationCountsAsDebt`

### Files (S04-01)

| Path | Role |
|------|------|
| `internal/store/schema/018_outcome_results_baselines.sql` | Tables + CHECKs + indexes |
| `internal/store/outcomes.go` + `outcomes_test.go` | CRUD + debt SQL |
| `internal/domain/outcomes.go` + `outcomes_test.go` | Gate helpers + named tests |
| `internal/domain/service.go` | Entity/rel/event constants |
| Compat / embed tests | Ceiling **18**, forbid **019+** |

Do **not** edit: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go`, `internal/mcp`, `internal/gitcli`, `changes` schema.

### Later scopes (upcoming notes only)

- **S05:** `comparison_json.regression` flags are inputs; no `regressions` table in S04.
- **S06:** apply keys `test_results`/`verifications`/`evaluations`; packet `verification_debt`; wire `VerificationIncomplete`.
- **S07:** seed export of outcomes/baselines is explicit residual.

## Merge table

COVERAGE §11–13, §20 satisfied by `outcome_results` + `baselines` + debt query. Metric/Observation remain merged into evaluation scores (no Metric table). Test/Verify/Evaluate remain **result kinds**, not runners.

## Planner work

1. [x] Lock table name `outcome_results` + kind enum + `baselines`.
2. [x] Lock verification: goal_id + evidence links; test ≠ verify.
3. [x] Lock evaluation vs baseline (structured comparison, not boolean).
4. [x] Lock verification debt query for S06.
5. [x] Lock DONE/Review policy unchanged; debt via deliberation.
6. [x] Thicken `01-verify-evaluate-gates.md` + `02-scope-review.md`.
7. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] 01/02 thickened with named tests + proof commands
- [x] Explicit: Trace does not claim tests passed without recorded test result evidence
- [x] §16/§18 remain Future
- [x] Compat ceiling **18** locked
- [x] No product Go

## Next

Orchestrator: **P20-S04-01** after this row is `done`.
