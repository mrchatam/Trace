# P21-S07-01 — Implement: promoted MVP cuts

## Metadata
- id: P21-S07-01
- todo_ids: [P21-S07-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Thin `experiments` table + deterministic `risk_hints` in loop next per S07-00. Closes D-01 + D-02 promoted MVP cuts. **No bake-off runner, no ML, no subprocess test execution.**

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-01, D-02
- [WORK-MAP.md](../../WORK-MAP.md) W-12, W-13
- Live: `internal/loop/next.go`, `internal/store/changes.go`, `internal/domain/cognitive.go`, `internal/domain/outcomes.go`

## Prerequisites
**S06 done** — transactional apply + compat **20** landed. S07 is first scope to add mig **021**.

## Locked defaults (from S07-00 — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`021_experiments.sql`** only — no ALTER on prior tables |
| Experiment fields | id, task_id, label, hypothesis_summary, status, outcome_result_id, timestamps |
| Status values | `planned`, `running`, `completed` — forward lifecycle only |
| Outcome link | Via `experiments.outcome_result_id` — **not** new `outcome_results.kind` |
| Risk hints | Additive `risk_hints` on `NextPacket`; max **4** items |
| Hint codes | `many_paths`, `blocking_uncertainty`, `high_churn_path`, `missing_verification` |
| Thresholds | many_paths: **>8** paths on latest change; churn: **≥3** changes per path; blocking: count>0; verification: debt present |
| Priority | blocking → missing_verification → many_paths → high_churn_path |
| Loop schema | `trace.loop.next.v1` **unchanged** |
| Compat ceiling | **21** (forbid **022+**) |
| Seed export | **No changes** |
| MCP | No new tools |

## Live inventory (before — confirmed S07-00)

| Surface | Location | Today |
|---------|----------|-------|
| Schema | `internal/store/schema/` | **20** files; max **020**; `019` notes "No experiment tables" |
| Next packet | `loop/next.go` L47–62 | No `risk_hints` field |
| Change paths | `store/changes.go` | `ListChangePaths`; no churn aggregate helper |
| Policy signals | `loop/policy.go` | Blocking + verification debt already queried for SelectNext |
| Compat | `evals/compat/compat_test.go` L260+ | Ceiling **20**; forbids **021+** |
| Experiments domain | — | **Missing** |

## Requirements

### 1. Migration + store (W-12)

1. Add `internal/store/schema/021_experiments.sql` per FINAL SQL in 00-PLANNER.
2. Add `internal/store/experiments.go`:
   - `Experiment` struct matching table columns
   - `UpsertExperiment`, `GetExperiment`, `ListExperimentsByTaskID`
   - Empty `outcome_result_id` = no link (not NULL — match repo convention)
3. Register in embed/migration list (follow `020` pattern).
4. Add `ListHighChurnPaths(taskID, minChanges int)` — SQL GROUP BY path HAVING COUNT(DISTINCT change_id) >= minChanges.

### 2. Domain experiments (W-12)

Add `internal/domain/experiments.go`:

1. `CreateExperiment(ctx, ExperimentInput)` — allocate id, default `planned`, emit optional event only if cheap.
2. `SetExperimentStatus(ctx, id, status)` — validate enum; allow forward transitions; idempotent same status.
3. `LinkExperimentOutcome(ctx, experimentID, outcomeResultID)` — verify outcome exists; set link.
4. `GetExperiment`, `ListExperimentsByTaskID` — pass-through store.

**Forbidden:** Any call to `os/exec`, test runners, or subprocess spawning from these APIs.

### 3. Risk hints in loop next (W-13)

1. Add types to `internal/loop/risk_hints.go` (or extend `deliberation_packet.go`):

```go
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

2. Implement `buildRiskHintsSection` per 00-PLANNER rules:
   - Evaluate four rules in priority order
   - Append hints until cap **4**
   - Empty conditions → `{freshness, items:[]}` (valid packet)
3. Add `RiskHints RiskHintsSection` to `NextPacket` in `next.go`.
4. Wire in `BuildNextPacket` — pass `domain.Service` + task id; reuse existing freshness string.

**Advisory only:** Hints do **not** mutate SelectNext, work_state, or auto-run tests.

### 4. Compat ceiling bump

Update `evals/compat/doc.go` and `evals/compat/compat_test.go`:
- Expect embed/applied **21**
- Require `021_experiments.sql`
- Forbid **022+** (mirror S04→S07 pattern from 020→021)

### 5. Tests (6 named + keepers)

| Test | Assert |
|------|--------|
| `TestCreateExperimentLinksOutcome` | Create experiment; upsert test outcome; link; `GetExperiment` shows `outcome_result_id` |
| `TestExperimentStatusLifecycle` | `planned`→`running`→`completed`; bogus status → validation error |
| `TestRiskHintsManyPaths` | Task + change with **9** paths → builder includes `many_paths` |
| `TestRiskHintsBlockingUncertainty` | OPEN BLOCKING uncertainty on task → `blocking_uncertainty` hint |
| `TestLoopNextRiskHintsBounded` | Seed all four conditions → `len(risk_hints.items) <= 4`; first hint is `blocking_uncertainty` |
| `TestNoExperimentRunnerInvoked` | Experiment create/link path never imports/calls exec runner (test helper or build-tag guard) |

Also keep green: S06 `TestLoopApplyTransactionalRollbackOnFailure`, S05 `TestLoopNextHistoricalRelationshipsSection`, S04 `TestPromoteBaselineSupersedesPrior`, P20 change/outcome keepers.

## Touch files

- `internal/store/schema/021_experiments.sql`
- `internal/store/experiments.go`
- `internal/store/changes.go` — `ListHighChurnPaths` (if not separate file)
- `internal/domain/experiments.go`
- `internal/domain/experiments_test.go`
- `internal/loop/risk_hints.go`
- `internal/loop/risk_hints_test.go`
- `internal/loop/next.go`
- `evals/compat/compat_test.go`
- `evals/compat/doc.go`

## Keeper floor

```bash
go test ./internal/domain/... -count=1 -run 'TestCreateExperimentLinksOutcome|TestExperimentStatusLifecycle|TestNoExperimentRunnerInvoked'
go test ./internal/loop/... -count=1 -run 'TestRiskHintsManyPaths|TestRiskHintsBlockingUncertainty|TestLoopNextRiskHintsBounded'
go test ./internal/store/... -count=1 -run 'TestMigrationStatusReportsEmbedMax|TestOpenCreatesDBAndMigratesIdempotent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] Mig **021** applied; schema dir has **21** files
- [ ] 6 named tests PASS
- [ ] Compat ceiling **21**
- [ ] `trace.loop.next.v1` unchanged; `risk_hints` present in JSON
- [ ] No runner subprocess; no MCP additions
- [ ] Board notes: test output summary

## Next

**P21-S07-02**
