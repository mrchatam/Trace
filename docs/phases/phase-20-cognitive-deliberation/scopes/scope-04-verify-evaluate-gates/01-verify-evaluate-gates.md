# P20-S04-01 — Implement verify / evaluate / gates

## Metadata
- id: P20-S04-01
- todo_ids: [P20-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Implement `outcome_results` + `baselines`, promotion gate helpers, and verification debt queries per S04-00 FINAL. **Not** a test runner product. **Library-only.**

## Session start
Follow agent-loop-protocol. Unattended: execute after S04-00 is `done`. Board edits: **status + notes only**. Do not re-debate locks. **No CLI / MCP / loop / SelectNext edits.**

## Locked defaults (from S04-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`018_outcome_results_baselines.sql`** (after 017). Additive. No ALTER on `changes` / `tasks` / S02 tables. |
| Tables | `outcome_results` (kind `test`\|`verification`\|`evaluation`), `baselines` (git SHA + scores JSON) |
| Forbidden on `changes` | Any `tests` / `verification_runs` / `baseline` / `score_*` columns |
| work_state | **No new enum values** |
| Test | `test_name` + `test_status`; stored row required for gate — agent claim insufficient (Law 2) |
| Verification | `goal_id` + ≥1 `outcome_supported_by` → evidence; `test` pass ≠ verified |
| Evaluation | Library-computed `comparison_json` vs baseline — **not** boolean PASS |
| DONE policy | **Unchanged** — Review PASS + operator; debt via deliberation only |
| Compat | Ceiling **18**; forbid `019+` |
| Out | CLI, MCP, loop apply keys, seed export, subprocess test runner, §16/§18 engines |

Full SQL, rels, API signatures, debt definition, and comparison JSON shape: [00-PLANNER.md](00-PLANNER.md) FINAL section. Copy them; do not invent columns.

## Requirements

1. Store CRUD matching locked SQL (CHECK enums, indexes, kind-specific column hygiene).
2. Domain wrappers with fail-closed `ErrValidation`:
   - `CreateBaseline` / `GetBaseline`
   - `RecordTestOutcome` / `RecordVerificationOutcome` / `RecordEvaluationOutcome`
   - `CompareScoresToBaseline` (pure)
   - `CheckTestGate` / `CheckVerificationGate` / `CheckEvaluationGate`
   - `HasImplementationSignal` / `HasVerificationDebt` / `ListVerificationDebtSummary`
3. Entity/rel/event constants on `internal/domain/service.go` (`outcome_result`, `baseline`, `outcome_supported_by`, `outcome.recorded`, `baseline.created`, `evaluation.compared`).
4. Verification: fail closed if `goal_id` empty, zero evidence IDs, or evidence row missing.
5. Evaluation: fail closed if `baseline_id` empty, invalid `scores_json`, or agent tries to pre-set boolean pass in `comparison_json`.
6. Test: fail closed if `test_name` empty or unknown `test_status`.
7. Bump embed/compat tests that still expect **17** / forbid **018** → ceiling **18** / forbid **019+**.
8. Extend no-blob discipline: new tables must not store full log blobs (summary capped 4096; scores/comparison capped 16384).

## Input shapes (domain)

```text
BaselineInput:
  GitCommit   string   // required; OID 7–64 lowercase
  ScoresJSON  string   // required JSON object, non-empty
  Label       string
  SourceType  string

TestOutcomeInput:
  TaskID      string   // required
  TestName    string   // required
  TestStatus  string   // pass|fail|skip|error
  Summary     string   // optional bounded
  EvidenceIDs []string // optional outcome_supported_by
  Actor       string
  SourceType  string

VerificationOutcomeInput:
  TaskID               string   // required
  GoalID               string   // required
  VerificationStatus   string   // verified|failed|partial
  EvidenceIDs          []string // required ≥1
  Summary              string
  Actor, SourceType    string

EvaluationOutcomeInput:
  TaskID      string   // required
  BaselineID  string   // required
  ScoresJSON  string   // required JSON object
  Actor, SourceType string
  // comparison_json computed by domain — not supplied by caller
```

## Named tests (must exist and pass)

See 00-PLANNER list (14 names). Minimum proofs:

- Record test with name + pass; gate helper returns true only when row exists
- Test pass + no verification row → `CheckVerificationGate` false; `CheckTestGate` true does not satisfy verification
- Verification without goal_id or without evidence → fail closed
- Verification with goal + evidence + verified → `CheckVerificationGate` true
- Evaluation produces `comparison_json` with dimension deltas; not `{pass: true}`
- Regression dimension sets `regression: true` in comparison when score drops
- Baseline row stores OID + scores JSON only (no blob columns)
- Change RECORDED + no verification → `HasVerificationDebt` true
- Verified verification with evidence → debt false
- `HasImplementationSignal` false when no changes → no debt
- Unknown kind / cross-kind column pollution → fail closed
- Evaluation without baseline → fail closed
- `partial` verification status → still debt

## Likely touch points

- `internal/store/schema/018_outcome_results_baselines.sql` (**new**)
- `internal/store/outcomes.go` + `outcomes_test.go` (**new**)
- `internal/domain/outcomes.go` + `outcomes_test.go` (**new**)
- `internal/domain/service.go` (constants only)
- `evals/compat/compat_test.go` + `evals/compat/doc.go` (17→18, 018 present, forbid 019+)
- `internal/store/production_hardening_test.go`
- `internal/store/deliberation_test.go` (EmbedExpected 18)
- `internal/store/store_test.go` (`TestOpenCreatesDBAndMigratesIdempotent` include 18)

Do **not** touch: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go`, `internal/mcp`, `internal/store/schema/017_*`, `internal/domain/changes.go`.

## Proof commands

```bash
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestOutcome|TestBaseline|TestNoSourceContentColumns'
go test ./internal/domain/ -count=1 -run 'TestRecordTestOutcome|TestTestPassAlone|TestVerification|TestEvaluation|TestBaseline|TestVerificationDebt|TestPromotionGate|TestPartialVerification'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 Loop keepers must stay green (this row must not change loop).

## Todo updates
Status + notes only on `P20-S04-01`. Next after green: `P20-S04-02`.

## Exit criteria

- [ ] Migration 018 applied; EmbedExpected/compat ceiling **18**
- [ ] All 14 named tests green
- [ ] Law 2: test pass ≠ verification; claims ≠ evidence in gate paths
- [ ] Law 14/15: evaluation is structured comparison, not boolean PASS; deterministic checks
- [ ] No test-runner subprocess product embedded
- [ ] No new `work_state` values; DONE policy unchanged
- [ ] No score/test columns added to `changes`
- [ ] No daemon / CoT blobs / hosted MCP / CLI in this row
