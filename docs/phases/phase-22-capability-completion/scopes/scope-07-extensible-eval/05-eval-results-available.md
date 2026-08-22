# P22-S07-05 — Implement: evaluation results for future agents

## Metadata
- id: P22-S07-05
- todo_ids: [P22-S07-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Stable library API so **evaluation results are available to future agents** (**C42-library**). S05 already surfaces compact evaluations on context packets; this row makes **`internal/eval`** the query SoT with **`mechanism_id`**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- S05: `internal/compiler/evidence_sections.go` `buildEvaluations` — **do not remove**; optional thin delegate to `eval.ListResults` is OK but not required this row
- S07-01/03: registry + rules loader

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| Context `evaluations[]` cap 8 (S05) — no mechanism_id | `eval.ListResults` API |
| `ListOutcomeResultsByTaskKind(evaluation)` in store | Mechanism id on query rows |
| `RecordEvaluationOutcome` + `comparison_json` | `trace eval results` CLI |

## Locked defaults

| Item | Value |
|------|-------|
| API | **`eval.ListResults(ctx, svc *domain.Service, taskID string) ([]ResultRow, error)`** in `internal/eval/query.go` |
| ResultRow | `{ ID, TaskID, MechanismID, Kind, Passed *bool, Summary, ScoresJSON, ComparisonJSON, CreatedAt }` |
| Mapping | **No new DB columns** — derive `MechanismID` from existing rows: test → `stored_test`; verification → `stored_verification`; evaluation → `stored_evaluation` |
| Evaluations | Include full **`comparison_json`** (not truncated) for kind=evaluation rows |
| Mechanism runs | Include latest **`RunAll`** snapshot optional — **lock v1**: synthesize pass/fail from mapped outcomes + run `architectural_invariant` once when listing if no stored row (honest `RecordedAt` from query time) OR attach last mechanism results only when explicitly requested — **prefer: map stored outcomes only + separate `eval.RunAll` for live checks** |
| List scope | All `outcome_results` for task where kind ∈ {test, verification, evaluation}; newest-first; default limit **32**, cap **64** (match S05 evidence queries) |
| FTS | **No new entity type** — evaluations already sync via `UpsertOutcomeResult` (S05-02a); do not add mechanism-run FTS |
| CLI (optional) | `trace eval results --task <id>` JSON array stdout; reuse `cli:eval` |
| Compiler | **Do not require** compiler refactor this row — C42-library closes on eval package + tests |
| Checklist | C42 full line stays **partially boxed** — S05 surface already `[x]`; this row boxes library half in Notes + test name |
| Schema | stays **26**; compat **26** |
| MCP | **No new tools** — catalog stays **13** |

## Requirements

1. **`ListResults`** — domain/store read only via `svc`; stable sort `created_at DESC, rowid DESC` (match S03/S06 ordering fix).
2. **`MechanismID` on every row** — asserted in `TestEvalResultsIncludeMechanismID`.
3. **`TestListEvaluationResultsForFutureAgents`** — seed task; `RecordEvaluationOutcome`; list; assert id, task_id, comparison_json non-empty, mechanism_id=`stored_evaluation`.
4. Include at least one test + verification row mapping in same test or sibling test.
5. Optional CLI thin encode.
6. **Do not** change S05 packet shape this row (no mechanism_id on `EvaluationItem` required for close).

## Touch files

- `internal/eval/query.go`, `query_test.go` (new)
- `cmd/trace/eval.go` (extend — `results` subcommand optional)
- `cmd/trace/eval_test.go` or extend existing

## Named tests

| Test | Proves |
|------|--------|
| `TestListEvaluationResultsForFutureAgents` | C42-library — evaluation row round-trip |
| `TestEvalResultsIncludeMechanismID` | C40+C42 — every row has mechanism_id |
| `TestContextIncludesEvaluations` | keeper (S05) — surface still PASS if compiler untouched |
| `TestCompatibilitySecurityChecklist` | ceiling **26** |

```bash
go test ./internal/eval/... ./internal/domain/... -count=1 -run 'TestListEvaluationResults|TestEvalResultsIncludeMechanismID|TestRecordEvaluationOutcome'
go test ./internal/compiler/... -count=1 -run TestContextIncludesEvaluations
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestEvalResults|TestEvalRules'
```

## Exit criteria

- [ ] C42-library true (evidence via named tests)
- [ ] Named tests PASS
- [ ] S05 context keeper PASS
- [ ] Board Notes: mention C42 library half closed (checklist line may need full box at S07-06)

## Minimal todos

- [ ] ListResults API + mapping
- [ ] Named tests
- [ ] Optional CLI results
- [ ] Board status + notes

## Residual risks (carry to S07-06)

- **Packet vs library drift** — compiler still reads store directly; reviewer accepts until optional delegate row
- **Live RunAll vs stored** — v1 list is stored-outcomes-first; document in review
- **Cap 32/64** — no overflow keeper required (low)
