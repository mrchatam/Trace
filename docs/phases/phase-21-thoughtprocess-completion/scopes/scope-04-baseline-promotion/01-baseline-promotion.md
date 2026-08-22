# P21-S04-01 — Implement: baseline promotion

## Metadata
- id: P21-S04-01
- todo_ids: [P21-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Implement baseline promotion chain + eval regression promotion gate per S04-00. Wire optional `promotion_blocked` on loop status. **Advisory gate only** — no auto `work_state` / DONE changes.

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-09, D-10
- [WORK-MAP.md](../../WORK-MAP.md) W-05, W-06
- P20 gates: [scope-04/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-04-verify-evaluate-gates/00-PLANNER.md)
- Live: `internal/domain/outcomes.go`, `internal/store/outcomes.go`, `internal/loop/apply.go`, `internal/domain/seed_export.go`

## Locked defaults (from S04-00 — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`020_baselines_promotion.sql`** — ALTER `baselines` only; additive |
| New columns | `status` (`active`\|`superseded`, default `active`); `supersedes_id` (default `''`) |
| Chain key | One `active` baseline per **`git_commit` + `label`** |
| Promote | `PromoteBaseline` supersedes prior active on same key; sets `supersedes_id`; idempotent |
| Regression gate | Latest evaluation with computed comparison: `overall_regression=true` → blocked / `eval_regression` |
| No eval | Gate blocked / `no_stored_evaluation` — test pass alone insufficient |
| DONE / work_state | **Unchanged** — `TransitionTask` does not consult promotion gate |
| Loop status | Additive `promotion_blocked: {present, reason}` on `trace.loop.status.v1` |
| Compat ceiling | **20** (forbid **021+**) |
| Seed | Optional `status`, `supersedes_id` on `SeedBaseline` — backward compatible |

### Migration SQL (copy verbatim)

```sql
-- Migration v20: baseline promotion chain (Phase 21 S04).
-- Additive only; do not rewrite 001-019.

ALTER TABLE baselines ADD COLUMN status TEXT NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'superseded'));
ALTER TABLE baselines ADD COLUMN supersedes_id TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_baselines_commit_label_status
    ON baselines(git_commit, label, status);
```

## Live inventory (before — confirmed S04-00)

| Surface | Location | Today |
|---------|----------|-------|
| Baselines table | `018_outcome_results_baselines.sql` | 7 columns; no status chain |
| `store.Baseline` | `outcomes.go` L37–45 | No Status / SupersedesID |
| `UpsertBaseline` | `outcomes.go` L81–112 | INSERT/UPDATE 7 columns |
| `CreateBaseline` | `domain/outcomes.go` L285–321 | Insert-only snapshot |
| `CompareScoresToBaseline` | `domain/outcomes.go` L161–226 | Computes `overall_regression` |
| Gate helpers | `domain/outcomes.go` L544–627 | Test / Verification / Evaluation only |
| `comparisonComputed` | `domain/outcomes.go` L528–542 | Private helper — reuse for gate |
| Loop status | `apply.go` L514–576 | No `promotion_blocked` |
| `statusBlocked` | `deliberation_packet.go` L317–319 | Uncertainty / regression / verification debt only |
| Compat | `evals/compat/compat_test.go` L260–268 | Ceiling **19**, forbid **020+** |
| Schema files | `internal/store/schema/` | **19** files, max **019** |

## Requirements

### 1. Migration + store

1. Add `internal/store/schema/020_baselines_promotion.sql` (FINAL shape above).
2. Extend `store.Baseline` with `Status`, `SupersedesID` (`store.BaselineStatusActive`, `store.BaselineStatusSuperseded` constants).
3. Update `UpsertBaseline`, `GetBaseline`, `ListAllBaselines` to read/write new columns.
4. Add store helpers:
   - `GetActiveBaselineByCommitLabel(gitCommit, label string) (Baseline, error)` — returns `ErrNotFound` or equivalent when none
   - `SetBaselinePromotion(id, status, supersedesID string) error`

### 2. Domain — promote / supersede

1. **`PromoteBaseline(ctx, baselineID)`**
   - Validate baseline exists.
   - If already `active` → no-op (idempotent).
   - Find prior **active** baseline with same `git_commit`+`label` (exclude self).
   - If prior exists: `SupersedeBaseline(prior.ID)`.
   - Set target `status=active`, `supersedes_id=prior.ID` (or `''` if first in chain).
   - Append `baseline.promoted` event (`EventBaselinePromoted` on `EntityBaseline`).

2. **`SupersedeBaseline(ctx, baselineID)`** — set `status=superseded`; no delete.

3. **`CreateBaseline`** — ensure new rows persist `status=active`, `supersedes_id=''`.

### 3. Domain — `CheckPromotionGate`

```go
func (s *Service) CheckPromotionGate(ctx context.Context, taskID, baselineID string) (allowed bool, reason string, err error)
```

Logic (locked):

1. Require valid `taskID`.
2. Optional `baselineID`: if set, baseline must exist (else fail closed / `baseline_not_found`).
3. Load evaluations for task; filter `comparisonComputed(comparison_json)`.
4. If none → `(false, "no_stored_evaluation", nil)`.
5. Pick **latest** by `created_at` desc, `id` desc.
6. Parse `comparison_json` → if `overall_regression==true` → `(false, "eval_regression", nil)`.
7. Else → `(true, "", nil)`.

**Do not** call from `TransitionTask`. **Do not** auto-transition work_state.

### 4. Loop status wiring

1. Add `PromotionBlocked` struct + field on `StatusResult` (`json:"promotion_blocked,omitempty"`).
2. In `Status()` (after deliberation build): call `CheckPromotionGate(ctx, seed.TaskID, "")`.
3. Set `promotion_blocked.present = !allowed`; copy `reason` when blocked.
4. **Do not** fold into `statusBlocked()` — keep orthogonal.

### 5. Seed export/import (minimal)

- `SeedBaseline`: add optional `Status`, `SupersedesID` with `omitempty`.
- Export when non-default; import defaults `active` / `''` when absent.
- Re-run S01 seed round-trip tests after changes.

### 6. Compat ceiling **20**

Bump embed expected **19 → 20** in compat + store tests listed in S04-00. Assert **021+** forbidden. Update compat doc string.

### 7. Tests — **6 named + P20 keepers**

| Test | Proves |
|------|--------|
| `TestPromoteBaselineSupersedesPrior` | B100 active → B101 promoted → B100 superseded, chain linked |
| `TestPromoteBaselineIdempotent` | Double promote → no error, single supersede |
| `TestEvalRegressionBlocksPromotionGate` | `overall_regression=true` → blocked / `eval_regression` |
| `TestEvalRegressionGateClearsAfterResolve` | Older regression + newer clean eval → gate open |
| `TestPromotionGateIndependentOfTestPassAlone` | Test pass only → `no_stored_evaluation` |
| `TestBaselinePromotionRequiresStoredEvaluation` | No eval → gate blocked |

Use score fixtures from existing `TestEvaluationRegressionFlagInComparisonJSON` patterns.

## Touch files

- `internal/store/schema/020_baselines_promotion.sql` (**new**)
- `internal/store/outcomes.go` — struct, CRUD, promotion helpers
- `internal/store/outcomes_test.go` — extend roundtrip for status columns
- `internal/domain/outcomes.go` — PromoteBaseline, SupersedeBaseline, CheckPromotionGate
- `internal/domain/outcomes_test.go` — 6 new tests
- `internal/domain/service.go` — `EventBaselinePromoted`
- `internal/domain/seed_export.go`, `seed_import.go` — optional baseline chain fields
- `internal/loop/apply.go` — `PromotionBlocked` + wire in `Status()`
- `evals/compat/compat_test.go`, `doc.go`
- `internal/store/production_hardening_test.go`, `deliberation_test.go`

## Keeper floor

```bash
go test ./internal/domain/... -count=1 -run 'TestPromoteBaselineSupersedesPrior|TestPromoteBaselineIdempotent|TestEvalRegressionBlocksPromotionGate|TestEvalRegressionGateClearsAfterResolve|TestPromotionGateIndependentOfTestPassAlone|TestBaselinePromotionRequiresStoredEvaluation|TestEvaluationComparesScoresToBaselineNotBoolean|TestEvaluationRegressionFlagInComparisonJSON|TestPromotionGateRequiresStoredTestNotAgentClaim'
go test ./internal/store/... -count=1 -run 'TestBaselineStoreRoundtrip|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

Optional loop status assertion (add if cheap):

```bash
go test ./cmd/trace -count=1 -run 'TestLoopStatusPromotionBlockedWhenEvalRegression'
```

## Minimal todos

- [ ] Add mig 020 + store Baseline promotion columns/helpers
- [ ] Implement PromoteBaseline / SupersedeBaseline / CheckPromotionGate
- [ ] Wire promotion_blocked on loop Status()
- [ ] Seed export/import additive fields
- [ ] Bump compat ceiling to 20
- [ ] Add 6 named tests; run keeper floor; paste PASS in board Notes

## Exit criteria

- [ ] 6 new + S04 P20 outcome tests PASS
- [ ] Compat ceiling **20**; no **021+** migration file
- [ ] `TransitionTask` unchanged (no auto DONE block)
- [ ] Board row status + Notes only

## Next

**P21-S04-02**
