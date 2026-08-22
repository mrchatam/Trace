# P22-S04-00 — Planner: predicted vs actual + regression + improvements

## Metadata
- id: P22-S04-00
- todo_ids: [P22-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S04 against live impact walk, effects, regressions, and outcomes. Owned: **C08, C16, C18**. Query of regressions is **S05-03** (CLI/MCP only — domain helpers land here). **No product Go.**

## Live inventory (2026-08-18, post-S03)

| Surface | Live state |
|---------|------------|
| Schema max | **023** (`023_graph_sync.sql`; 23 embed sql files) |
| Compat ceiling | **23** (`evals/compat/compat_test.go` — no 024+) |
| Impact walk | `internal/retrieval/impact_walk.go` — BFS depth 1..2; seeds `file\|symbol`; blast + `AffectedTests` (reverse `validates`); no `architectural_boundary` hop; **not snapshotted on a change** |
| CLI impact | `cmd/trace/impact.go` — `finding\|alternative\|report\|walk` only; **no `predict` / `compare`** |
| MCP impact | `internal/mcp/tools_impact.go` — walk mirror; catalog **10** |
| Decision impact | `internal/domain/impact.go` — findings/alternatives/report (decision-scoped); **not** code-graph predicted vs actual |
| Changes | `changes` + `change_paths` (path, symbol_id); `PromoteVCSCommitToChange` (S02); statuses OPEN/RECORDED/COMPARED/SUPERSEDED |
| Effects | `effects` table — `expected` / `actual` / `comparison` per `(change_id, dimension)`; `RecordExpectedEffect` / `RecordActualEffect`; **not** blast-set compare |
| Regressions | `regressions` table (019); create always **`correlated`**; `SetRegressionAttributionCaused` fail-closed (hypothesized → confirmed hypothesis + evidence); links `regression_from_evaluation`, `regression_from_effect`; **no `change_id` column or regression↔change link** |
| Test regression (C13) | `domain.DetectTestRegression` / `DetectAnyTestRegression` — outcome row ordering; **does not write `regressions` rows or associate changes** |
| Improvements | **Absent** — no table, no domain API, no seed key; positive `effects.comparison=supported` exists but is not first-class C18 |
| Seed export | exports `changes`, `effects`, `regressions`, `reflections`; **no `impact_predictions` or `improvements`** |
| S01 deps | `validates` edges, `AffectedTests`, `ListValidatesForFile` — predict/compare must reuse retrieval walk |
| S03 deps | `trace test run`, `trace verify run`, stored outcomes — compare runs **after** index on change paths, not subprocess re-test |

S01 closed **C01–C03, C07**; S02 closed **C04–C06, C25**; S03 closed **C09, C11–C15, C36, C38-CLI** — do not reopen in S04 prompts.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-08, D-22-09, D-22-19, D-22-21
- [WORK-MAP.md](../../WORK-MAP.md) W-14…W-16
- Coverage: [README.md](../../README.md) C08, C16, C18 rows
- Law 5 (P20): never auto-`caused` from correlation alone

## FINAL locked defaults

| Item | Value |
|------|-------|
| Mig | **`024_impact_compare.sql`** only — **`impact_predictions`** + **`improvements`**; **no ALTER on `regressions`** (change association via `entity_links`) |
| Compat | **24** after S04-01; S04-03/S04-05 stay **24** (forbid **025+** entire S04) |
| Predict timing | **Before implement**: snapshot walk keyed by **`change_id`** (one row per change; upsert) |
| Predict payload | JSON of entity keys + metadata — **no source blobs**: `{seeds, blast_keys[], affected_test_keys[], depth, blast_total, blast_kept, truncated}` where keys are `"file:<uuid>"\|"symbol:<uuid>"` |
| Predict seeds | From `change_paths`: prefer `symbol_id` when set; else resolve `path` → file id via store (fail-closed if path not indexed) |
| Actual | **After** index: re-walk same seeds + depth from stored row; compare path/symbol/test **key sets** |
| Compare delta | Persist + return `{matched[], unexpected[], missed[]}` (sorted keys); optional `impact_compare_json` column on `impact_predictions` **or** compute-only — **lock: store `compare_json` on row at compare time** |
| CLI C08 | `trace impact predict --change <id> [--depth 1\|2]`; `trace impact compare --change <id>` |
| Regression↔change | New link rel **`regression_associated_change`**: `from_type=regression`, `to_type=change`, `to_id=change_id`; query **`ListRegressionsByChangeID`** (+ alias **`RegressionsForChange`**) |
| Auto-associate | `RecordRegressionFromContradictedEffect` **must** insert associated_change link to `effect.change_id` (still **`correlated`**, not caused) |
| Caused + change | `SetRegressionAttributionCaused` unchanged policy; S04-03 adds **`AssociateRegressionWithChange`**; caused path test must show link + attribution |
| Law 5 | Correlation / contradiction / confirmed hypothesis alone **never** set `caused`; keepers **`TestCorrelationAndContradictionNeverAutoSetCaused`**, **`TestSetAttributionCausedFailClosedWithoutEvidence`** |
| Improvements | Dedicated **`improvements`** table (not outcome kind) — C18 queryable without overloading `effects` |
| Improvement fields | `id`, `change_id`, `task_id`, `dimension`, `summary`, `evidence_ids_json` (JSON array, max 32), `source_type`, `confidence`, `created_at`, `updated_at` |
| Improvement API | `RecordImprovement`, `ListImprovementsByChangeID`, `ListImprovementsByTaskID`; summary max 4096 (match regression cap) |
| Improvement CLI | `trace outcomes improvements --change <id>` **or** `--task <id>` (extend existing `outcomes` subcommand — no new root verb) |
| Seed (D-22-19) | Additive `improvements[]` in export/import; export `impact_predictions` **optional** (may omit from seed v1 — lock: **export improvements only**; predictions stay local `.trace/` unless implementer adds both with tests) |
| MCP | **No new MCP tools this scope** — S05/S08 parity later; catalog stays **10** |

### Mig 024 DDL (locked shape)

```sql
-- impact_predictions: one snapshot per change
CREATE TABLE impact_predictions (
    change_id TEXT PRIMARY KEY,
    predicted_json TEXT NOT NULL,
    compare_json TEXT NOT NULL DEFAULT '',
    depth INTEGER NOT NULL DEFAULT 2,
    created_at TEXT NOT NULL,
    compared_at TEXT NOT NULL DEFAULT ''
);

-- improvements: first-class C18 rows
CREATE TABLE improvements (
    id TEXT PRIMARY KEY,
    change_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    dimension TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    evidence_ids_json TEXT NOT NULL DEFAULT '[]',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX idx_improvements_change_id ON improvements(change_id);
CREATE INDEX idx_improvements_task_id ON improvements(task_id);
```

## Named tests

| Test | Row |
|------|-----|
| `TestRecordPredictedImpactThenCompareActual` | S04-01 |
| `TestImpactCompareUnexpectedAndMissed` | S04-01 |
| `TestRegressionLinkedToChangeCausedWithEvidence` | S04-03 |
| `TestSetAttributionCausedFailClosedWithoutEvidence` | keeper |
| `TestCorrelationAndContradictionNeverAutoSetCaused` | keeper |
| `TestRecordImprovementQueryable` | S04-05 |
| `TestSeedExportIncludesImprovements` | S04-05 |
| `TestCompatibilitySecurityChecklist` | S04-01 (ceiling **24**) |

## Keeper floor (do not break)

Regression attribution: `TestSetAttributionCausedRequiresConfirmedHypothesisAndEvidence`, `TestLinkHypothesisUpgradesToHypothesizedNotCaused`, `TestRecordRegressionFromEvaluationDefaultsCorrelated`.

S01/S03 spot-check: `TestImpactWalkIncludesAffectedTests`, `TestRegressionDetectedVsPriorPassingTest` (30/30 stable), `TestPromoteVCSCommitCreatesChangeIdempotent`.

Compat: exactly **23** sql files until S04-01 lands **024**; then **24** only.

## Residual risks for S04-01

| Risk | Mitigation locked in 01 |
|------|-------------------------|
| Change path not indexed → no file/symbol seed | Fail-closed with clear error; test uses indexed fixture from analyzers testdata |
| Blast truncation skews compare | Persist `truncated` in predicted_json; compare uses **kept** keys only; test asserts unexpected/missed on controlled small graph |
| Depth mismatch predict vs compare | Store `depth` on row; compare re-reads stored depth (ignore CLI default) |
| Double predict overwrites silently | Upsert by `change_id`; test asserts second predict replaces predicted_json |
| Accidental 025+ / compat drift | S04-01 only bumps embed + compat; grep 24 sql files; S04-03/05 forbid 025+ |
| Conflating effects compare with C08 | C08 is blast-key sets from ImpactWalk; do not repurpose `effects.expected/actual` |
| `compare` without prior predict | Fail-closed `ErrValidation` |

## Exit criteria

- [x] 01–06 thickened vs live impact/regression/outcome code
- [x] Mig **024** locked (`impact_predictions` + `improvements`; regression link via `entity_links`)
- [x] No product Go

## Next

**P22-S04-01**
