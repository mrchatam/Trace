# P21-S04-00 — Planner: baseline promotion + eval regression gate

## Metadata
- id: P21-S04-00
- todo_ids: [P21-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock baseline promote/supersede chain (B100→B101) and promotion block when `overall_regression` in latest evaluation. **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-09, D-10
- [WORK-MAP.md](../../WORK-MAP.md) W-05, W-06
- P20 baseline: [scope-04/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-04-verify-evaluate-gates/00-PLANNER.md)
- Live: `internal/domain/outcomes.go`, `internal/domain/outcomes_test.go`, `internal/store/outcomes.go`, `internal/store/schema/018_outcome_results_baselines.sql`, `internal/loop/apply.go`, `evals/compat/compat_test.go`

## Live inventory (confirmed 2026-08-18)

| Surface | Today (live read) | S04 action |
|---------|-------------------|------------|
| `CreateBaseline` | `outcomes.go` L285–321 — insert-only via `UpsertBaseline`; emits `baseline.created` | **Add** `PromoteBaseline`, `SupersedeBaseline` |
| `store.Baseline` | `outcomes.go` L37–45 — id, git_commit, scores_json, label, source_type, timestamps only | **Add** `Status`, `SupersedesID` fields |
| Baseline chain | `018_outcome_results_baselines.sql` — no `status` / `supersedes_id` | **Add** columns via mig **020** |
| `CompareScoresToBaseline` | `outcomes.go` L161–226 — computes `overall_regression` | **Keep** — gate reads stored `comparison_json` |
| `RecordEvaluationOutcome` | `outcomes.go` L453–516 — persists library-computed comparison | Unchanged |
| Gate helpers | `CheckTestGate`, `CheckVerificationGate`, `CheckEvaluationGate` only | **Add** `CheckPromotionGate` |
| `TransitionTask` | `task_state.go` — Review PASS + operator; **no** eval regression check | **Do not** auto-block DONE; gate is advisory |
| Loop `statusBlocked` | `deliberation_packet.go` L317–319 — uncertainty / open_regression / verification_debt | **Separate** `promotion_blocked` on status JSON |
| `StatusResult` | `apply.go` L192–202 — no promotion fields | **Add** optional `promotion_blocked` |
| Schema / compat | max mig **019** (19 files); compat ceiling **19** | S04-01 adds **020**; ceiling **20** (forbid **021+**) |
| Seed export | `SeedBaseline` — 7 fields; no status chain | **Additive** optional `status`, `supersedes_id` (omitempty) |
| P20 outcome tests | 11 tests in `outcomes_test.go`; no promote/gate tests | **Add** 6 named S04 tests |

## FINAL locked defaults (S04-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`020_baselines_promotion.sql`** — additive ALTER on `baselines` only; do **not** rewrite 001–019 |
| New columns | `status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','superseded'))`; `supersedes_id TEXT NOT NULL DEFAULT ''` |
| Index | `idx_baselines_commit_label_status ON baselines(git_commit, label, status)` |
| Chain key | **`git_commit` + `label`** — at most one `active` baseline per key |
| `CreateBaseline` | New rows default `status='active'`, `supersedes_id=''` |
| Promote flow | `PromoteBaseline(ctx, baselineID)` → load target; find prior **active** row with same `git_commit`+`label` (if any); mark prior `superseded`; set target `active` with `supersedes_id=prior.ID`; emit `baseline.promoted` event |
| `SupersedeBaseline` | Store/domain helper: set `status=superseded` on given id (no delete — Law 11) |
| Idempotent promote | Second `PromoteBaseline` on already-`active` id → no-op success |
| Regression gate | Latest `kind=evaluation` for task (by `created_at` desc, tie-break `id` desc) with computed `comparison_json`: if `overall_regression=true` → `CheckPromotionGate` returns `(false, "eval_regression", nil)` |
| No evaluation | No stored evaluation with computed comparison → `(false, "no_stored_evaluation", nil)` |
| Test-only | Stored test pass **without** evaluation → promotion gate **blocked** (`no_stored_evaluation`) |
| Gate open | Latest evaluation exists and `overall_regression=false` → `(true, "", nil)` |
| Resolved regression | Newer evaluation without `overall_regression` supersedes older regressed eval for gate purposes |
| Auto work_state | **Forbidden** — gate is advisory; operator decides DONE; `TransitionTask` unchanged |
| Loop status | Additive `promotion_blocked: {present: bool, reason: string}` on `trace.loop.status.v1` — populated from `CheckPromotionGate(seed.task_id, "")` or active baseline for task if wired |
| Compat ceiling | **20** after S04-01 (forbid **021+** until S07) |
| §16 / §18 | Unchanged Future |

### Migration SQL (FINAL shape)

```sql
-- Migration v20: baseline promotion chain (Phase 21 S04).
-- Additive only; do not rewrite 001-019.

ALTER TABLE baselines ADD COLUMN status TEXT NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'superseded'));
ALTER TABLE baselines ADD COLUMN supersedes_id TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_baselines_commit_label_status
    ON baselines(git_commit, label, status);
```

### Domain API (FINAL)

```go
// PromoteBaseline marks baselineID active and supersedes prior active baseline for same git_commit+label.
func (s *Service) PromoteBaseline(ctx context.Context, baselineID string) (store.Baseline, error)

// SupersedeBaseline marks a baseline superseded (Law 11 — no delete).
func (s *Service) SupersedeBaseline(ctx context.Context, baselineID string) error

// CheckPromotionGate is advisory — does not mutate work_state or auto-DONE.
// baselineID optional: when empty, gate uses latest evaluation baseline_id if any.
func (s *Service) CheckPromotionGate(ctx context.Context, taskID, baselineID string) (allowed bool, reason string, err error)
```

**Latest evaluation selection:** `ListOutcomeResultsByTaskKind(taskID, evaluation)` → filter `comparisonComputed(comparison_json)` → sort `created_at` desc, `id` desc → take first → parse `overall_regression` from JSON.

**Promotion gate reasons (locked vocabulary):**

| reason | When |
|--------|------|
| `""` | Allowed |
| `eval_regression` | Latest evaluation has `overall_regression=true` |
| `no_stored_evaluation` | No evaluation with computed comparison for task |
| `baseline_not_found` | `baselineID` arg set but row missing (fail closed) |

### Loop status additive field

```go
type PromotionBlocked struct {
    Present bool   `json:"present"`
    Reason  string `json:"reason,omitempty"`
}
// StatusResult.PromotionBlocked *PromotionBlocked `json:"promotion_blocked,omitempty"`
```

Populate in `Status()` / `buildStatusDeliberation` via domain `CheckPromotionGate`. **Do not** merge into `statusBlocked()` — promotion is orthogonal to deliberation INVESTIGATE triggers.

### Seed v1 additive (backward compatible)

| JSON key | Type | Default on import |
|----------|------|-------------------|
| `status` | string | `active` when absent |
| `supersedes_id` | string | `""` when absent |

Export when non-default; old seeds without keys import as `active`/empty.

### Compat / embed bump (S04-01)

Update hard-coded **19 → 20** in:

- `evals/compat/compat_test.go` — ceiling check + forbid `021+`
- `internal/store/production_hardening_test.go` — `EmbedExpected`
- `internal/store/deliberation_test.go` — `TestMigrationStatusReportsEmbedMax`
- `evals/compat/doc.go` — prose ceiling

**Do not** add `021_experiments.sql` in S04 (S07 lock).

## Named tests (S04-01)

| Test | Proves |
|------|--------|
| `TestPromoteBaselineSupersedesPrior` | B100 active → create B101 same commit+label → promote B101 → B100 `superseded`, B101 `active`, B101.supersedes_id=B100 |
| `TestPromoteBaselineIdempotent` | Double promote same ID → still active, no duplicate supersede events |
| `TestEvalRegressionBlocksPromotionGate` | Evaluation with `overall_regression=true` → gate blocked, reason `eval_regression` |
| `TestEvalRegressionGateClearsAfterResolve` | Regressed eval then newer eval without regression → gate open |
| `TestPromotionGateIndependentOfTestPassAlone` | Test pass only → gate blocked (`no_stored_evaluation`) |
| `TestBaselinePromotionRequiresStoredEvaluation` | No eval row → gate blocked |

**Keep** all P20 outcome tests green (`TestEvaluationComparesScoresToBaselineNotBoolean`, `TestEvaluationRegressionFlagInComparisonJSON`, `TestPromotionGateRequiresStoredTestNotAgentClaim`, etc.).

## Touch files

- `internal/store/schema/020_baselines_promotion.sql` (**new**)
- `internal/store/outcomes.go` — Baseline struct, Upsert/Get/List scan promotion columns; `GetActiveBaselineByCommitLabel`, `SetBaselinePromotion`
- `internal/domain/outcomes.go` — PromoteBaseline, SupersedeBaseline, CheckPromotionGate
- `internal/domain/outcomes_test.go` — 6 named tests
- `internal/domain/seed_export.go` + `seed_import.go` — optional status/supersedes_id
- `internal/domain/service.go` — `EventBaselinePromoted` constant
- `internal/loop/apply.go` — `PromotionBlocked` on `StatusResult`; wire in `Status()`
- `evals/compat/compat_test.go` — ceiling **20**
- `internal/store/production_hardening_test.go`, `deliberation_test.go` — embed max **20**

## Planner work

1. [x] Live inventory `outcomes.go` / store / schema 018 / loop status / compat.
2. [x] Lock W-05 (baseline promote/supersede B100→B101) + W-06 (eval regression promotion gate).
3. [x] Lock mig **020** shape + compat ceiling **20** (forbid 021+).
4. [x] Thicken `01-baseline-promotion.md` + `02-scope-review.md` with before-state, 6 tests, keeper floor.
5. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Mig 020 shape locked
- [x] 6 named tests locked
- [x] 01/02 thickened enough to implement alone
- [x] No product Go

## Next

**P21-S04-01**
