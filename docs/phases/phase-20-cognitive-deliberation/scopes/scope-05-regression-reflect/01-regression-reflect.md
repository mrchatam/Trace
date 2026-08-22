# P20-S05-01 — Implement regression / reflect / history

## Metadata
- id: P20-S05-01
- todo_ids: [P20-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Implement `regressions` + `reflections`, attribution honesty (`correlated` ≠ `caused`), structured reflection writes, and observed vs causal `entity_links` per S05-00 FINAL. **Library-only.**

## Session start
Follow agent-loop-protocol. Unattended: execute after S05-00 is `done`. Board edits: **status + notes only**. Do not re-debate locks. **No CLI / MCP / loop / SelectNext edits.**

## Locked defaults (from S05-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`019_regressions_reflections.sql`** (after 018). Additive. No ALTER on `outcome_results` / `effects` / `hypotheses` / `entity_links` / `tasks` |
| Tables | `regressions`, `reflections` only |
| Derivation | Evaluation `comparison_json` regression flags **or** contradicted effect — nothing else |
| Attribution | Create always `correlated`. Enum `correlated`\|`hypothesized`\|`caused` |
| Caused | **Never auto.** `SetRegressionAttributionCaused` only after hypothesized + CONFIRMED hypothesis + ≥1 evidence |
| Hypothesis | `LinkHypothesisToRegression` → `hypothesized` not `caused`. Do not fork Discovery. Do not hook `ConfirmHypothesis` |
| Reflection | Structured JSON arrays; **no `body`**. Essay-only (summary/note without arrays) fail closed |
| Links | `observed_relationship` + confidence; `caused_by` requires evidence IDs. Causal link ≠ attribution `caused` |
| Open regression | `HasOpenRegression` for S06; do not auto-hop |
| Compat | Ceiling **19**; forbid `020+` |
| Out | CLI, MCP, loop apply keys, seed export, §16/§18 engines, SelectNext table edits |

Full SQL, rels, API signatures, caused evidence policy, and reflection JSON shapes: [00-PLANNER.md](00-PLANNER.md) FINAL section. Copy them; do not invent columns.

## Requirements

1. Store CRUD matching locked SQL (CHECK enums, UNIQUE source, indexes).
2. Domain wrappers with fail-closed `ErrValidation`:
   - `RecordRegressionFromEvaluation` / `RecordRegressionFromContradictedEffect`
   - `LinkHypothesisToRegression` / `SetRegressionAttributionCaused`
   - `ResolveRegression` / `SupersedeRegression`
   - `HasOpenRegression` / `CountOpenRegressionsByTaskID` / `ListOpenRegressions`
   - `CreateReflection` / `GetReflection` / `ListReflectionsByTaskID`
   - `RecordObservedRelationship` / `RecordCausalRelationship`
3. Entity/rel/event constants on `internal/domain/service.go` (`regression`, `reflection`, rels in S05-00, events `regression.recorded`, `regression.attribution_changed`, `reflection.recorded`, …).
4. Evaluation path: fail closed if outcome missing, not `kind=evaluation`, or `comparison_json` has no `overall_regression` and no dimension `regression: true`. Persist `attribution=correlated`. Do **not** insert from `RecordEvaluationOutcome`.
5. Effect path: fail closed if effect missing or `comparison != contradicted`. Persist `correlated`. Do **not** set `caused`. Do **not** auto-copy `hypothesis_explains_effect`.
6. `SetRegressionAttributionCaused`: fail closed from `correlated`; fail closed with empty/missing evidence; fail closed unless linked hypothesis is `CONFIRMED`.
7. Reflection: fail closed if all three arrays empty; fail closed on non-array JSON / `body`-style essay column (column must not exist). Persist assumption links.
8. `observed_relationship`: confidence in [0,1]; no evidence required. `caused_by`: evidence IDs required; must **not** flip regression attribution.
9. Bump embed/compat tests that still expect **18** / forbid **019** → ceiling **19** / forbid **020+**.
10. Named test: `HasOpenRegression==true` fed as complete `PolicyInputs` into `ApplyDeliberationTransition` → `INVESTIGATE` / `open_regression` (not `SelectNext` alone; not auto-called from Record*).

## Input shapes (domain)

See 00-PLANNER `EvaluationRegressionInput` / `EffectRegressionInput` / `ReflectionInput` / `RelInput`. Copy; do not add a `body` or `attribution` override on create.

## Named tests (must exist and pass)

Exact names from 00-PLANNER (14). Minimum proofs:

- Evaluation with `overall_regression` → row `attribution=correlated` not `caused`
- Contradicted effect → `correlated`; contradiction + optional S03 hypothesis link does not set `caused`
- `RecordRegressionFromEvaluation` / FromEffect reject create-time `caused`/`hypothesized`
- `LinkHypothesisToRegression` upgrades `correlated` → `hypothesized`, still not `caused`
- `SetRegressionAttributionCaused` without evidence IDs → fail closed
- `SetRegressionAttributionCaused` from `correlated` (no hypothesized step) → fail closed
- Caused succeeds only with CONFIRMED hypothesis + existing evidence rows + `regression_supported_by`
- Open row → `HasOpenRegression` true; passed into `ApplyDeliberationTransition` → INVESTIGATE / `open_regression`
- `ResolveRegression` → open flag false
- Reflection stores the three JSON arrays and they are readable by Get/List (S06-shaped)
- Reflection with only `summary` / `broaden_tests_note` → fail closed; no `body` column
- `observed_relationship` persists with confidence and zero evidence links
- `caused_by` without evidence → fail closed; with evidence does not set `attribution=caused`
- Unknown attribution string → fail closed

## Likely touch points

- `internal/store/schema/019_regressions_reflections.sql` (**new**)
- `internal/store/regressions.go` + `regressions_test.go` (**new**)
- `internal/domain/regressions.go` + `regressions_test.go` (**new**)
- `internal/domain/service.go` (constants only)
- `evals/compat/compat_test.go` + `evals/compat/doc.go` (18→19, 019 present, forbid 020+)
- `internal/store/production_hardening_test.go`
- `internal/store/deliberation_test.go` (EmbedExpected 19)
- `internal/store/store_test.go` (`TestOpenCreatesDBAndMigratesIdempotent` include 19)

Do **not** touch: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go`, `internal/mcp`, `internal/domain/outcomes.go`, `internal/domain/changes.go`, `internal/domain/cognitive.go`, `internal/store/schema/018_*`.

## Proof commands

```bash
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestRegression|TestReflection|TestNoSourceContentColumns'
go test ./internal/domain/ -count=1 -run 'TestRecordRegression|TestCorrelation|TestLinkHypothesis|TestSetAttribution|TestHasOpenRegression|TestResolveRegression|TestReflection|TestObservedRelationship|TestCausalRelationship|TestUnknownAttribution'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 Loop keepers must stay green (this row must not change loop).

## Todo updates
Status + notes only on `P20-S05-01`. Next after green: `P20-S05-02`.

## Exit criteria

- [ ] Migration 019 applied; EmbedExpected/compat ceiling **19**
- [ ] All 14 named tests green
- [ ] correlated ≠ caused enforced; never auto-`caused` from correlation or contradiction
- [ ] Reflection structured fields persist; essay-only fail closed
- [ ] `observed_relationship` vs `caused_by` evidence policy
- [ ] `HasOpenRegression` helper exists; SelectNext table untouched
- [ ] §16/§18 not implemented (stub note only)
- [ ] No daemon / CoT blobs / hosted MCP / CLI in this row
