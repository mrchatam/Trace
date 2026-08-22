# P21-S07-00 — Planner: promoted MVP cuts (§16 + §18)

## Metadata
- id: P21-S07-00
- todo_ids: [P21-S07-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock **thin** §16 experiment records and **minimal** §18 risk-adaptive verification hints — promoted from P20 Future, not full bake-off engine. **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-01, D-02
- [WORK-MAP.md](../../WORK-MAP.md) W-12, W-13
- [COVERAGE.md](../../../phase-20-cognitive-deliberation/COVERAGE.md) §16/§18 Future → promote thin
- P20 deferral: [scope-05/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-05-regression-reflect/00-PLANNER.md) ("§16/§18 Future — no experiment runner")
- Live: `internal/loop/{next.go,deliberation_packet.go,policy.go}`, `internal/store/schema/`, `internal/domain/changes.go`, `evals/compat/compat_test.go`

## Live inventory (confirmed 2026-08-18)

| Surface | Location | Today (live read) | S07 action |
|---------|----------|-------------------|------------|
| Experiments | — | **No** `experiments` table; `019` comment "No experiment tables" | **Add** mig **021_experiments.sql** thin table |
| Outcome kinds | `018_outcome_results_baselines.sql` L19 | CHECK `test`\|`verification`\|`evaluation` only; `kind=experiment` fail-closed in store test | **Keep** — link via `experiments.outcome_result_id`, not new outcome kind |
| Change paths | `017_changes_effects.sql` + `store/changes.go` | `change_paths(change_id, path, …)`; `ListChangePaths` per change | **Read** for risk hints — no schema change |
| Recent changes packet | `deliberation_packet.go` L243–298 | Caps paths at `maxChangePathsCap=16` per change in display | Risk hint uses **full** path count from store (not display cap) |
| Blocking uncertainties | `domain/cognitive.go` L295+ | `CountBlockingUncertainties` wired in `BuildPolicyInputs` | **Reuse** for `blocking_uncertainty` hint |
| Verification debt | `domain/outcomes.go` L763+ | `HasVerificationDebt` wired in loop sections | **Reuse** for `missing_verification` hint |
| Next packet | `next.go` L47–62 | Sections through `historical_relationships`; **no** `risk_hints` | **Add** `RiskHintsSection` (additive JSON) |
| Loop schema string | `next.go` L17 | `trace.loop.next.v1` | **Unchanged** — additive field only |
| Schema / compat | `internal/store/schema/` | **20** files, max **020**; compat ceiling **20** (forbid **021+**) | S07-01 adds **021**; ceiling **21** (forbid **022+**) |
| Seed export | `seed_export.go` | No `experiments` key | **No seed changes** — operational records; S08 verify uses domain/CLI |
| Test runner | repo-wide | Agents/harness run tests; Trace records outcomes only (D-16) | **Forbidden** — `TestNoExperimentRunnerInvoked` proves no subprocess runner in experiment path |
| MCP catalog | unchanged since S05 | **10** tools | **No new tools** |

### P20 Future evidence (why promote thin)

- **D-01:** P20-00 forbidden §16 product; COVERAGE row §16 Future ("first-class Experiment objects and multi-candidate bake-offs").
- **D-02:** P20-00 forbidden §18 product; COVERAGE §18 Future ("Risk-adaptive test-selection policy engine").
- TRACE §16 asks for "smallest practical version" comparing baseline vs candidates — S07 records **metadata + optional outcome link**, not runners.
- TRACE §29F lists "risk-adaptive test selection" — S07 emits **deterministic hints** in loop next, not ML matrix.

## W-12 / W-13 locks

| Work ID | TRACE § | Lock |
|---------|---------|------|
| **W-12** | §16 Experiments | Thin **`experiments`** table (mig **021**); status lifecycle; optional `outcome_result_id` link — **no** bake-off runner, **no** multi-agent orchestration, **no** new `outcome_results.kind` |
| **W-13** | §18 Risk-adaptive hints | **Deterministic** `risk_hints[]` on loop next from change-path count + blocking uncertainties + path churn + verification debt — max **4** hints; advisory only |

## Hard boundaries (unchanged)

- No multi-agent bake-off runner
- No ML risk matrix
- No autonomous test execution (D-16)
- No subprocess test runner invoked from Trace experiment APIs
- No `trace.loop.next.v2` / apply / status version bump

## FINAL locked defaults (S07-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`021_experiments.sql`** — additive CREATE only; do **not** rewrite 001–020 or ALTER `outcome_results` |
| Experiment row | `id`, `task_id`, `label`, `hypothesis_summary`, `status` (`planned`\|`running`\|`completed`), `outcome_result_id` (empty = none), `created_at`, `updated_at` |
| Create default | `status=planned`, `outcome_result_id=''` |
| Status lifecycle | Forward-only: `planned`→`running`→`completed`; idempotent re-set same status OK; **no** auto-transition |
| Outcome link | `LinkExperimentOutcome(experimentID, outcomeResultID)` sets `outcome_result_id` on existing experiment; outcome must exist; **no** runner side effects |
| Domain API | `CreateExperiment`, `SetExperimentStatus`, `LinkExperimentOutcome`, `GetExperiment`, `ListExperimentsByTaskID` |
| Store API | `UpsertExperiment`, `GetExperiment`, `ListExperimentsByTaskID` |
| Events | Optional `experiment.created` / `experiment.status_changed` — **omit** if time-boxed; table rows sufficient for MVP |
| §18 Risk hints | Top-level `risk_hints` on `NextPacket`; shape `{code, severity, detail}` |
| Hint codes | `many_paths`, `blocking_uncertainty`, `high_churn_path`, `missing_verification` — **fixed set only** |
| Hint cap | Max **4** items; deterministic priority order (below) |
| Priority order | 1 `blocking_uncertainty` → 2 `missing_verification` → 3 `many_paths` → 4 `high_churn_path` (emit first N that apply) |
| `many_paths` rule | Most recent change for task (by `created_at` desc, `id` desc) has **>8** `change_paths` rows |
| `blocking_uncertainty` | `CountBlockingUncertainties(taskID) > 0` |
| `high_churn_path` | Any path appears on **≥3** distinct changes for same `task_id` |
| `missing_verification` | `HasVerificationDebt(taskID) == true` |
| Severities | `blocking_uncertainty`→`high`; `missing_verification`→`medium`; `many_paths`→`medium`; `high_churn_path`→`low` |
| Detail strings | Human-readable counts/path snippet — **no** ML scores; e.g. `"9 paths on latest change"`, `"path internal/loop/apply.go touched 3 times"` |
| Loop schema | `trace.loop.next.v1` string **unchanged** (additive JSON field only) |
| Compat ceiling | **21** after S07-01 (forbid **022+**) |
| Seed export | **Unchanged** — no `experiments` JSON key in S07 |
| MCP | No new tools |

### Migration SQL (FINAL shape)

```sql
-- Migration v21: experiments (Phase 21 S07).
-- Additive only; do not rewrite 001-020. Thin §16 record — no bake-off runner.

CREATE TABLE IF NOT EXISTS experiments (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    hypothesis_summary TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'planned'
        CHECK (status IN ('planned', 'running', 'completed')),
    outcome_result_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_experiments_task_id ON experiments(task_id);
CREATE INDEX IF NOT EXISTS idx_experiments_status ON experiments(task_id, status);
```

### Domain API (FINAL)

```go
type ExperimentInput struct {
    TaskID            string
    Label             string
    HypothesisSummary string
}

func (s *Service) CreateExperiment(ctx context.Context, in ExperimentInput) (store.Experiment, error)
func (s *Service) SetExperimentStatus(ctx context.Context, id, status string) error
func (s *Service) LinkExperimentOutcome(ctx context.Context, experimentID, outcomeResultID string) error
```

### Risk hints builder (FINAL)

Implement `buildRiskHintsSection(ctx context.Context, dom *domain.Service, st *store.Store, taskID, freshness string) (RiskHintsSection, error)` in `internal/loop/` (new `risk_hints.go` or `deliberation_packet.go`):

```go
const maxRiskHintsCap = 4
const manyPathsThreshold = 8
const highChurnMinChanges = 3

type RiskHint struct {
    Code     string `json:"code"`
    Severity string `json:"severity"`
    Detail   string `json:"detail"`
}

type RiskHintsSection struct {
    Freshness string     `json:"freshness"`
    Items     []RiskHint `json:"items"`
}
```

Wire in `BuildNextPacket` after `recent_changes` / before return — same task freshness as seed.

Store helper for churn (FINAL):

```go
// ListHighChurnPaths returns paths touched on >= minChanges distinct changes for taskID.
func (s *Store) ListHighChurnPaths(taskID string, minChanges int) ([]string, error)
```

### Compat updates (FINAL)

- `evals/compat/doc.go` — schema through **021_experiments** (no **022+**)
- `evals/compat/compat_test.go` — expect **21** embed/applied files; forbid **022+**; assert `021_experiments.sql` present

## Named tests (S07-01)

| # | Test | Location | Proves |
|---|------|----------|--------|
| 1 | `TestCreateExperimentLinksOutcome` | `internal/domain/experiments_test.go` | Create experiment; link existing `outcome_result_id`; read back |
| 2 | `TestExperimentStatusLifecycle` | `internal/domain/experiments_test.go` | `planned`→`running`→`completed`; invalid status fail-closed |
| 3 | `TestRiskHintsManyPaths` | `internal/loop/risk_hints_test.go` | Change with **9** paths → `many_paths` hint in builder output |
| 4 | `TestRiskHintsBlockingUncertainty` | `internal/loop/risk_hints_test.go` | Blocking uncertainty → `blocking_uncertainty` hint |
| 5 | `TestLoopNextRiskHintsBounded` | `internal/loop/risk_hints_test.go` or `cmd/trace/loop_test.go` | All four rules true → packet has **≤4** hints; priority order preserved |
| 6 | `TestNoExperimentRunnerInvoked` | `internal/domain/experiments_test.go` | Create/link experiment does **not** call `exec.Command` / `os/exec` / spawn test processes (static helper or runtime guard) |

**Keep green (no regressions):** P20 outcome/change keepers, S06 apply tx tests, S05 historical section, compat after ceiling bump, MCP **10** tools.

## Touch files

- `internal/store/schema/021_experiments.sql` — **new**
- `internal/store/experiments.go` — **new** (CRUD)
- `internal/domain/experiments.go` — **new** (thin domain)
- `internal/domain/experiments_test.go` — **new** (tests 1, 2, 6)
- `internal/loop/risk_hints.go` — **new** (builder + types)
- `internal/loop/risk_hints_test.go` — **new** (tests 3–5)
- `internal/loop/next.go` — add `RiskHints` field + wire `buildRiskHintsSection`
- `evals/compat/compat_test.go` — ceiling **21**
- `evals/compat/doc.go` — comment ceiling **21**

**Do not touch:** `select.go`, `apply.go`, `seed_export.go`, MCP handlers, `outcome_results` CHECK kinds.

## Planner work

1. [x] Live inventory: schema max **020**, no experiments table, next packet lacks `risk_hints`, compat **20**.
2. [x] Lock **W-12** (thin experiments table + lifecycle + outcome link — no runner).
3. [x] Lock **W-13** (deterministic risk hints — 4 codes, cap 4, priority order, thresholds).
4. [x] Lock mig **021** + compat **21** (forbid **022+**).
5. [x] Thicken `01-promoted-mvp-cuts.md` + `02-scope-review.md` with before-state, 6 named tests, keeper floor.
6. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Thin experiment + risk hint locks explicit
- [x] W-12/W-13 explicit
- [x] Mig **021** + compat **21** locked
- [x] 6 named tests + touch files locked
- [x] 01/02 thickened enough to implement alone
- [x] No product Go

## Next

**P21-S07-01**
