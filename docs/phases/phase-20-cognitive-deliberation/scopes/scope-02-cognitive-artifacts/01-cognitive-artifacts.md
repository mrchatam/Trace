# P20-S02-01 — Implement cognitive artifacts

## Metadata
- id: P20-S02-01
- todo_ids: [P20-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Implement S02-00 FINAL schema + domain APIs + tests: uncertainty lifecycle, blocking-count query, assumption invalidation, hypothesis, decision reconsideration.

## Session start
Follow agent-loop-protocol. Unattended: execute after S02-00 is `done`. Board edits: **status + notes only**. Do not re-debate locks. **No CLI / MCP / loop / SelectNext table edits.**

## Locked defaults (from S02-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`016_cognitive_artifacts.sql`** (after 015). Additive. No ALTER on `assumptions`/`decisions`. |
| Tables | `uncertainties`, `hypotheses`, `decision_reconsiderations` |
| Blocking query | `CountOpenBlockingUncertaintiesByTaskID` — `severity=BLOCKING` AND `status=OPEN` AND rel `uncertainty_blocks_task` |
| Invalidate | `InvalidateAssumption` → `STALE`/`SUPERSEDED` on **existing** row; event `assumption.invalidated`; no DELETE |
| Reconsider | Child table; triggers `contradicted_effect`\|`new_evidence`\|`invalidated_assumption`; append-only |
| Finding | Reuse Discovery + existing `FindingKindInvalidatedAssumption` — **no Finding table** |
| PolicyInputs | Expose count only; persist hops via existing `ApplyDeliberationTransition` with **complete** inputs (do not call SelectNext alone) |
| Compat | Ceiling **16**; forbid `017+` |
| Out | CLI, MCP, loop apply keys, seed export, FTS for new types, auto-replan, raw CoT |

Full SQL, rels, API signatures, and transition rules: [00-PLANNER.md](00-PLANNER.md) FINAL section. Copy them; do not invent columns.

## Requirements

1. Store CRUD + count matching locked SQL (CHECK enums, indexes).
2. Domain wrappers with fail-closed validation (`ErrValidation`):
   - `CreateUncertainty` / `ResolveUncertainty` / `SupersedeUncertainty` / `CountBlockingUncertainties`
   - `CreateHypothesis` + OPEN→CONFIRMED\|REJECTED\|SUPERSEDED
   - `InvalidateAssumption` per 00-PLANNER input struct
   - `RecordDecisionReconsideration`
3. Entity/rel/event constants on `internal/domain/service.go` (`uncertainty`, `hypothesis`, listed rels).
4. BLOCKING create requires TaskID + `uncertainty_blocks_task` link.
5. Invalidate on `assumption_supports_decision`: impact finding + FIRED reconsideration; optional PLAN_AFFECTING Discovery + `discovery_invalidates_assumption`.
6. Bump embed/compat tests that still expect **15** / forbid **016**.

## Named tests (must exist and pass)

See 00-PLANNER list (14 names). Minimum proofs:

- Create → OPEN; resolve/supersede → count 0
- BLOCKING + task link increments count; INFO does not; missing TaskID on BLOCKING fails
- `TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition`: count=1 passed as `PolicyInputs.BlockingUncertaintyCount` into `ApplyDeliberationTransition` → phase `INVESTIGATE`, reason `blocking_uncertainty` (never EXECUTE)
- Invalidate STALE and SUPERSEDED: row still `GetAssumption`; no delete
- Linked decision gets `INVALIDATED_ASSUMPTION` finding; optional Discovery `PLAN_AFFECTING`
- Hypothesis → evidence via `hypothesis_supported_by` without inserting a Discovery as the hypothesis
- Reconsideration row recorded; Decision + alternatives still present
- Unknown severity fail closed

## Likely touch points

- `internal/store/schema/016_cognitive_artifacts.sql` (**new**)
- `internal/store/cognitive.go` + `cognitive_test.go` (**new**)
- `internal/domain/cognitive.go` + `cognitive_test.go` (**new**)
- `internal/domain/service.go` (constants only)
- `evals/compat/compat_test.go` + `evals/compat/doc.go` (15→16, 016 present, forbid 017+)
- `internal/store/production_hardening_test.go`
- `internal/store/deliberation_test.go` (EmbedExpected 16)
- `internal/store/store_test.go` (`TestOpenCreatesDBAndMigratesIdempotent` include 16)

Do **not** touch: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go` policy table, `internal/mcp`.

## Proof commands

```bash
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestDeliberationStateRoundtrip|TestBlocking|TestCreateUncertainty|TestInvalidate|TestHypothesis|TestDecisionReconsider'
go test ./internal/domain/ -count=1 -run 'TestCreateUncertainty|TestBlocking|TestResolve|TestSupersede|TestCountBlocking|TestInvalidateAssumption|TestHypothesis|TestDecisionReconsider|TestUnknownUncertainty|TestApplyDeliberationTransition'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 Loop keepers must stay green (this row must not change loop).

## Todo updates
Status + notes only on `P20-S02-01`. Next after green: `P20-S02-02`.

## Exit criteria

- [ ] Migration 016 applied; EmbedExpected/compat ceiling **16**
- [ ] All 14 named tests green
- [ ] No Finding/Requirement/Constraint tables
- [ ] Law 11: no silent history delete
- [ ] Blocking count query usable by S06 without further schema work
- [ ] No daemon / CoT blobs / hosted MCP / CLI in this row
