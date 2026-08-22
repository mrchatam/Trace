# P20-S05-00 — Regression / reflect / history planner

## Metadata
- id: P20-S05-00
- todo_ids: [P20-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock thin **Regression** (attribution `correlated` vs `hypothesized` vs `caused`), **Reflection** structured artifacts, and **observed vs causal** `entity_links`. No product Go this row.

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [COVERAGE.md](../../COVERAGE.md) §§15, 17, 19, 29J; merge table (Regression / Reflection / Relationship)
- Laws **1**, **2**, **5**, **11**, **14**, **15**, **16**, **19** ([G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md))
- Live: schema max **018** → next migration **`019_regressions_reflections.sql`**
- Live S04: `outcome_results.comparison_json` library-computed `{dimensions[k].regression, overall_regression}` — **inputs**, no `regressions` table in 018 (S04-02 APPROVE)
- Live S03: `effects.comparison=contradicted` + optional `hypothesis_explains_effect` / PLAN_AFFECTING Discovery / FIRED `contradicted_effect` — **inputs**; never `attribution=caused`
- Live S02: `hypotheses` status `OPEN`\|`CONFIRMED`\|`REJECTED`\|`SUPERSEDED`; `ConfirmHypothesis` / `CreateHypothesis`; do **not** fork Discovery as Hypothesis
- Live S01: `PolicyInputs.OpenRegression` + SelectNext row → `INVESTIGATE` / `open_regression` — S05 **exposes** query; S06 wires; S05 does **not** edit `select.go`
- Live `entity_links`: already has `rel` + `confidence` (`003_causal_domain.sql`) — **no ALTER**
- §16 Experiments + §18 Risk-adaptive verification → **Future** (not S05 implement)

## Doc map
§15, 17, 19, 29J

## Live inventory (2026-08-18)

| Surface | Location | S05 action |
|---------|----------|------------|
| `outcome_results` kind=evaluation | `018_outcome_results_baselines.sql` + `CompareScoresToBaseline` | **Input.** `overall_regression` / per-dimension `regression` flags derive a thin `regressions` row. Do **not** ALTER 018. `RecordEvaluationOutcome` stays as shipped (does not auto-insert regressions). |
| `effects.comparison` | `017_changes_effects.sql` | **Input.** `contradicted` derives a thin `regressions` row at `attribution=correlated`. Do **not** ALTER 017. |
| `hypothesis_explains_effect` | S03 rel | **Input.** Optional later upgrade path via explicit `LinkHypothesisToRegression`. Do not auto-copy on record. |
| `hypotheses` | `016_cognitive_artifacts.sql` + `ConfirmHypothesis` | Reuse. Upgrade `correlated` → `hypothesized` on explicit link. `CONFIRMED` is **necessary but not sufficient** for `caused`. Do **not** hook `ConfirmHypothesis` to auto-set `caused`. |
| `entity_links.confidence` | mig 003 | Reuse. New rels only; **no** extra metadata columns. |
| `PolicyInputs.OpenRegression` | `internal/deliberation/types.go` | Expose `HasOpenRegression`; **do not** call `SelectNext` / auto-hop / auto-replan. |
| Uncertainty auto-spawn | S03 lock: do not auto-spawn Question | **Out.** INVESTIGATE via S01 input is the MVP “question”; caller may create an uncertainty separately. |
| Experiment / risk-adaptive engine | missing | **Out.** §16/§18 Future. Optional reflection `broaden_tests_note` stub only. |
| Next migration | max **018** | S05-01 adds **`019_regressions_reflections.sql`** |
| Compat ceiling | **18** after S04 | S05-01 bumps to **19** (forbid `020+`) |

## Paths / table naming (FINAL)

| Option | Verdict |
|--------|---------|
| Fold regression into `outcome_results` / `effects` (`attribution` column) | **Rejected** — pollutes S03/S04 kinds; cannot represent “open regression” independently of a still-valid evaluation row |
| Essay `reflections.body` plus optional tags | **Rejected** — §19 requires queryable structured learning, not CoT |
| Relationship graph DB / new `relationships` table | **Rejected** — Law 13; COVERAGE reuses `entity_links` |
| **`regressions` + `reflections` tables; observed/causal as `entity_links` rels** | **Locked** |

S06 apply keys (`regressions`, `reflections`, optional relationship writes) are **transport aliases**. S05-01 is library-only — S06 owns loop apply wiring.

## FINAL locked defaults (S05-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`019_regressions_reflections.sql`** — additive; do not rewrite 001–018; **no ALTER** on `outcome_results` / `baselines` / `changes` / `effects` / `hypotheses` / `entity_links` / `tasks.work_state` |
| Compat ceiling | **19** after S05-01 (forbid `020+`); bump `evals/compat`, `production_hardening_test`, `deliberation_test` EmbedExpected, `TestOpenCreatesDBAndMigratesIdempotent` |
| New tables | `regressions`, `reflections` only |
| Forbidden | Experiment tables; risk-adaptive matrix; `relationships` table; `body`/essay column on reflections; `attribution=caused` default; auto-hop/auto-replan; raw CoT / log blobs |
| Attribution enum | `correlated` \| `hypothesized` \| `caused` — CHECK exact |
| Create default | **Always `correlated`.** Fail closed if create input sets `hypothesized` or `caused`. |
| Caused honesty | **Never auto-set `caused` from correlation, contradiction, evaluation flags, or hypothesis link.** Only `SetRegressionAttributionCaused` after evidence policy. |
| Hypothesis upgrade | Explicit `LinkHypothesisToRegression` (OPEN or CONFIRMED hypothesis) → `hypothesized` if currently `correlated`. Does **not** set `caused`. Reuse S02 `CreateHypothesis` / `ConfirmHypothesis` — do not fork Discovery. |
| Reflection | Structured JSON array columns (see below) — **not** essay-only. No `body` column. |
| Links | `observed_relationship` + confidence; `caused_by` requires evidence. Causal **link** ≠ regression `attribution=caused`. |
| Open regression | `HasOpenRegression(taskID)` for S06 packet + S01 `open_regression` input. |
| §16 / §18 | **Future** — no experiment runner; no risk-adaptive test-selection engine; optional `broaden_tests_note` one-liner on reflection only |
| CLI / MCP / loop / SelectNext | **Library-only** — S06 owns apply keys and `PolicyInputs.OpenRegression` wiring |

### `regressions`

```sql
CREATE TABLE IF NOT EXISTS regressions (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    source_kind TEXT NOT NULL
        CHECK (source_kind IN ('evaluation', 'contradicted_effect')),
    source_id TEXT NOT NULL,
    dimension TEXT NOT NULL DEFAULT '',
    attribution TEXT NOT NULL DEFAULT 'correlated'
        CHECK (attribution IN ('correlated', 'hypothesized', 'caused')),
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RESOLVED', 'SUPERSEDED')),
    summary TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(source_kind, source_id, dimension)
);
CREATE INDEX IF NOT EXISTS idx_regressions_task_id ON regressions(task_id);
CREATE INDEX IF NOT EXISTS idx_regressions_status ON regressions(task_id, status);
CREATE INDEX IF NOT EXISTS idx_regressions_attribution ON regressions(attribution);
```

| Field | Lock |
|-------|------|
| `task_id` | Required. Task must exist. Fail closed if empty. |
| `source_kind` | `evaluation` = S04 outcome_result id; `contradicted_effect` = S03 effect id. |
| `source_id` | Must exist as that kind (`GetOutcomeResult` kind=evaluation, or `GetEffect` with `comparison=contradicted`). |
| `dimension` | For evaluation: `'overall'` when `overall_regression` else the first regressing dimension key. For effect: the effect’s `dimension`. Max **64** bytes. |
| `attribution` | Create **must** persist `correlated`. Unknown values fail closed. |
| `status` | Lifecycle (not provenance). OPEN counts for `HasOpenRegression`. |
| `summary` | Bounded excerpt. Max **4096** bytes. **Not** a CoT essay or full comparison dump. |
| History | Law 11 — no `DeleteRegression`. RESOLVED/SUPERSEDED explicit. |
| Uniqueness | One row per `(source_kind, source_id, dimension)`. Re-record same source → return existing (or fail closed duplicate — pick **return existing**, do not clone). |

**Thin derivation (only these two sources):**

| Source | When a row may be created | Default attribution |
|--------|---------------------------|---------------------|
| Evaluation | `kind=evaluation` AND (`overall_regression==true` OR any `dimensions.*.regression==true`) | `correlated` |
| Contradicted effect | `effects.comparison=='contradicted'` | `correlated` |

**Not sources:** `kind=test` fail/error; verification `failed`/`partial`; `supported`/`partially_supported` effects; agent claim without a stored evaluation/effect row (Law 2).

S04 `RecordEvaluationOutcome` **does not** insert regressions (already shipped). S05 adds explicit `RecordRegressionFromEvaluation`. S06 apply may call both.

**Attribution transitions**

| From | To | API | Gate |
|------|----|-----|------|
| (create) | `correlated` | `RecordRegressionFromEvaluation` / `RecordRegressionFromContradictedEffect` | Fail closed if caller passes any other attribution |
| `correlated` | `hypothesized` | `LinkHypothesisToRegression` | Hypothesis row exists; status `OPEN` or `CONFIRMED`; inserts `hypothesis_explains_regression` |
| `hypothesized` | `caused` | `SetRegressionAttributionCaused` | Evidence policy below — **all** clauses |
| `correlated` | `caused` | — | **Fail closed** (must pass hypothesized) |
| `caused` | (other attribution) | — | **Fail closed** — terminal attribution; SUPERSEDE the row to reverse (Law 11) |

`ConfirmHypothesis` **must not** flip linked regressions to `caused`. `hypothesis_explains_effect` **must not** auto-set `caused` or even `hypothesized` until `LinkHypothesisToRegression`.

### Evidence policy for `attribution=caused` (fail-closed)

`SetRegressionAttributionCaused(ctx, regressionID, evidenceIDs)` succeeds only when **all** hold:

1. Explicit API call (no implicit path from Record*/Link*/ConfirmHypothesis/contradiction/evaluation flags).
2. Current `attribution == hypothesized` (not `correlated`).
3. ≥1 `hypothesis_explains_regression` link to a hypothesis with status **`CONFIRMED`** (OPEN is insufficient).
4. `evidenceIDs` non-empty; every id exists as `evidence`; insert `regression_supported_by` (regression → evidence).
5. Agent narrative / confidence number / `observed_relationship` / `caused_by` link alone **do not** satisfy this policy (Law 2, 15).

Inserting `caused_by` on `entity_links` is a **separate** §17 historical edge and **does not** set `regressions.attribution`.

### `reflections`

```sql
CREATE TABLE IF NOT EXISTS reflections (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    invalidated_assumptions_json TEXT NOT NULL DEFAULT '[]',
    new_dependencies_json TEXT NOT NULL DEFAULT '[]',
    useful_tests_json TEXT NOT NULL DEFAULT '[]',
    broaden_tests_note TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_reflections_task_id ON reflections(task_id);
```

| Field | Lock |
|-------|------|
| **No `body` column** | Essay-only storage is a defect. `summary` is a bounded caption, not the payload. |
| `summary` | Optional caption. Max **4096** bytes. |
| `invalidated_assumptions_json` | JSON **array of assumption UUID strings**. Max **32** items. Each id must exist (`GetAssumption`). Domain also inserts `reflection_invalidates_assumption` links. |
| `new_dependencies_json` | JSON **array of objects** `{ "kind": "path"\|"symbol"\|"file", "ref": "<string>" }`. Max **32**. `ref` max **512** bytes; paths via `store.NormalizePath` when kind=path/file. **No** file contents. |
| `useful_tests_json` | JSON **array of test_name strings**. Max **32**. Each max **256** bytes. Names, not outcome_result ids (test may not have been recorded). |
| `broaden_tests_note` | §18 **Future stub only**. Optional one-liner max **256** bytes. **Must not** select tests, expand a matrix, or feed SelectNext. |
| Create fail-closed | At least **one** of the three arrays must be non-empty `[]` is empty. Summary-only or note-only → `ErrValidation`. Invalid JSON / object root / string essay in an array column → fail closed. |
| History | Law 11 — no `DeleteReflection`. |

S06 reads these columns (and the assumption links) for packet slices — they must be queryable without parsing an essay.

§19 mappings **not** given their own columns: “what did we learn” → `summary`; “hypotheses supported/rejected” → reuse S02 hypothesis status + optional links (do not duplicate a fourth essay). Unexpected side effects stay on S03 `effects`.

### Entity links (new rels — no ALTER `entity_links`)

| rel | from → to | When | Evidence |
|-----|-----------|------|----------|
| `regression_from_evaluation` | regression → outcome_result | Always on evaluation-derived create | n/a |
| `regression_from_effect` | regression → effect | Always on contradicted-effect create | n/a |
| `hypothesis_explains_regression` | hypothesis → regression | `LinkHypothesisToRegression` | n/a (upgrade to hypothesized only) |
| `regression_supported_by` | regression → evidence | Required for `caused` | evidence row must exist |
| `reflection_invalidates_assumption` | reflection → assumption | Each id in `invalidated_assumptions_json` | assumption row must exist |
| `observed_relationship` | caller from → to | §17 observed (typically change/path/symbol → effect/outcome/task) | **Not** required. `confidence` in **[0, 1]**. Fail closed if confidence unset/NaN or outside range. |
| `caused_by` | caller from → to | §17 claimed causal | **Required:** ≥1 `relationship_supported_by` (from-entity → evidence) on the **same** call. Fail closed if `evidenceIDs` empty or evidence missing. Does **not** set `regressions.attribution`. |
| `relationship_supported_by` | from-entity of the `caused_by` edge → evidence | Inserted with `RecordCausalRelationship` | evidence must exist |

`observed_relationship` and `caused_by` are distinct rels. Do **not** encode causality as “high confidence observed”. Confidence on observed is correlation strength, not proof (Law 5: inferred ≠ verified).

Constants on `internal/domain/service.go`: `EntityRegression`, `EntityReflection`, rels above, events `regression.recorded`, `regression.attribution_changed`, `regression.resolved`, `reflection.recorded`, `relationship.observed`, `relationship.caused`.

### Open regression query (S06 + S01 input)

**Function (domain + store):**

- `HasOpenRegression(ctx, taskID) (bool, error)`
- `CountOpenRegressionsByTaskID(ctx, taskID) (int, error)`
- `ListOpenRegressions(ctx, taskID) ([]Regression, error)` bounded for packets (Law 6) — max **32** rows

**Open when:** `status='OPEN'` for that `task_id`. RESOLVED/SUPERSEDED do not count.

**S01 wiring (S06, not S05):** `PolicyInputs.OpenRegression = HasOpenRegression(taskID)`.

Named test: `HasOpenRegression==true` passed as complete `PolicyInputs` into `ApplyDeliberationTransition` → `INVESTIGATE` / `open_regression` (does **not** call `SelectNext` alone; does **not** auto-hop on record).

**Do not** auto-spawn Uncertainty, Discovery, PlanChange, or REPLAN (Laws 9, 16). S03 FIRED `contradicted_effect` reconsideration stays S03’s; S05 does not re-fire it.

### Domain API (S05-01)

| API | Role |
|-----|------|
| `RecordRegressionFromEvaluation(ctx, EvaluationRegressionInput)` | Inspect stored evaluation `comparison_json`; fail closed if no regression flags; persist `correlated` |
| `RecordRegressionFromContradictedEffect(ctx, EffectRegressionInput)` | Effect must be `contradicted`; persist `correlated` |
| `LinkHypothesisToRegression(ctx, hypothesisID, regressionID)` | → `hypothesized` from `correlated` |
| `SetRegressionAttributionCaused(ctx, regressionID, evidenceIDs)` | Evidence policy; fail closed otherwise |
| `ResolveRegression` / `SupersedeRegression` | status transition; reason required |
| `HasOpenRegression` / `CountOpenRegressionsByTaskID` / `ListOpenRegressions` | S06 + SelectNext input |
| `CreateReflection(ctx, ReflectionInput)` | Structured arrays; fail closed if essay-only |
| `GetReflection` / `ListReflectionsByTaskID` | |
| `RecordObservedRelationship(ctx, RelInput)` | rel=`observed_relationship` + confidence; no evidence required |
| `RecordCausalRelationship(ctx, RelInput, evidenceIDs)` | rel=`caused_by`; evidence required; does not set attribution |

```text
EvaluationRegressionInput:
  OutcomeID   string   // required; kind=evaluation
  TaskID      string   // required; must match outcome.task_id
  Actor, Summary, SourceType string
  Confidence  float64

EffectRegressionInput:
  EffectID    string   // required; comparison=contradicted
  TaskID      string   // required; must match parent change.task_id
  Actor, Summary, SourceType string
  Confidence  float64

ReflectionInput:
  TaskID                     string   // required
  Summary                    string
  InvalidatedAssumptionIDs   []string // JSON + links
  NewDependencies            []DependencyRef
  UsefulTests                []string
  BroadenTestsNote           string   // §18 stub only
  Actor, SourceType          string
  Confidence                 float64

DependencyRef: Kind (path|symbol|file), Ref (required)

RelInput:
  FromType, FromID, ToType, ToID string  // all required; entities must exist for known types
  Confidence float64                     // required [0,1]
  SourceType string
```

Events: `entity.created` (regression/reflection); `regression.recorded`; `regression.attribution_changed` (payload old/new); `regression.resolved`; `reflection.recorded`; `entity.linked` for new rels.

### Named tests (minimum — exact names)

1. `TestRecordRegressionFromEvaluationDefaultsCorrelated`
2. `TestRecordRegressionFromContradictedEffectDefaultsCorrelated`
3. `TestCorrelationAndContradictionNeverAutoSetCaused`
4. `TestLinkHypothesisUpgradesToHypothesizedNotCaused`
5. `TestSetAttributionCausedFailClosedWithoutEvidence`
6. `TestSetAttributionCausedFailClosedFromCorrelated`
7. `TestSetAttributionCausedRequiresConfirmedHypothesisAndEvidence`
8. `TestHasOpenRegressionFeedsApplyDeliberationTransition`
9. `TestResolveRegressionClearsHasOpenRegression`
10. `TestReflectionPersistsStructuredFieldsQueryable`
11. `TestReflectionEssayOnlyFailClosed`
12. `TestObservedRelationshipLinkWithConfidenceNoEvidence`
13. `TestCausalRelationshipFailClosedWithoutEvidence`
14. `TestUnknownAttributionFailClosed`

### Files (S05-01)

| Path | Role |
|------|------|
| `internal/store/schema/019_regressions_reflections.sql` | Tables + CHECKs + indexes |
| `internal/store/regressions.go` + `regressions_test.go` | CRUD + open-count SQL |
| `internal/domain/regressions.go` + `regressions_test.go` | Attribution policy + named tests |
| `internal/domain/service.go` | Entity/rel/event constants |
| Compat / embed tests | Ceiling **19**, forbid **020+** |

Do **not** edit: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go`, `internal/mcp`, `internal/gitcli`, `internal/domain/outcomes.go` (no auto-insert), `internal/domain/changes.go`, `internal/domain/cognitive.go` (`ConfirmHypothesis` stays attribution-agnostic).

### Later scopes (upcoming notes only)

- **S06:** apply keys `regressions` / `reflections`; packet `open_regressions[]` bounded; wire `PolicyInputs.OpenRegression = HasOpenRegression`; optional relationship writes fail-closed; ceiling **19** after this migration.
- **S07:** seed export of `regressions`/`reflections` is an explicit residual; VERIFY ceiling = live embed max (19 after S05-01); correlated ≠ caused remains a VERIFY fail bar.
- **§18:** `broaden_tests_note` is a stub string — S06 must not treat it as a test-selection policy.

## Merge table

COVERAGE §15, §17, §19, §29J satisfied by `regressions` + `reflections` + `observed_relationship`/`caused_by` on existing `entity_links`. Experiment remains Future. Metric/Observation stay merged into S04 comparison / S03 effects. No Relationship table.

## Planner work

1. [x] Lock `regressions` derivation (evaluation flags **or** contradicted effect) + attribution enum + caused evidence policy.
2. [x] Lock `reflections` structured fields (no essay-only body).
3. [x] Lock `observed_relationship` vs `caused_by` on `entity_links`.
4. [x] Lock named tests for 01; compat ceiling **19** after 01.
5. [x] Thicken `01-regression-reflect.md` + `02-scope-review.md`.
6. [x] Note later scopes; §18 Future; no product Go.

## Exit criteria

- [x] 01/02 thickened with named tests + proof commands
- [x] Causal honesty locked: never auto-`caused` from correlation
- [x] Compat ceiling **19** locked
- [x] §18 remains Future (stub note only)
- [x] Library-only; no product Go

## Next

Orchestrator: **P20-S05-01** after this row is `done`.
