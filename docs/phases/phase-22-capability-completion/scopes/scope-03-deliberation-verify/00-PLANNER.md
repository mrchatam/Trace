# P22-S03-00 — Planner: deliberation policy + verification cycle

## Metadata
- id: P22-S03-00
- todo_ids: [P22-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S03 against live `BuildPolicyInputs`, outcomes gates, and loop status. Owned: **C09, C11, C12, C13, C14, C15, C36, C38-CLI**. **No product Go.**

## Live inventory (2026-08-18, post-S02)

| Surface | Live state |
|---------|------------|
| Schema max | **023** (`023_graph_sync.sql`; 23 embed sql files) |
| Compat ceiling | **23** (`evals/compat/compat_test.go` — no 024+) |
| `internal/loop/policy.go` | Sets `BlockingUncertaintyCount`, `PlanExists`, `PlanCritiqued`, `VerificationIncomplete` (`HasVerificationDebt`), `OpenRegression`, `P19Saturated` only — **never assigns** `ExecutePending`, `TestPending`, `EvaluationPending`, `ReflectPending`, `ReplanNeeded`, `OpenDecisionAlternatives` (zero-value stubs) |
| `internal/deliberation/select.go` | **14-row FINAL** first-match table (P21-S03-01); rows 8–13 honor cycle flags when set |
| `internal/deliberation/types.go` | `PolicyInputs` has all cycle bool fields + JSON tags |
| Domain outcomes | `RecordTestOutcome`, `RecordVerificationOutcome`, `RecordEvaluationOutcome`; gates `CheckTestGate`, `CheckVerificationGate`, `CheckEvaluationGate`; `HasImplementationSignal`, `HasVerificationDebt`, `CompareScoresToBaseline` — package comment still says “no runner” (D-22-03 supersedes for **explicit** invoke only) |
| Store outcomes | `ListOutcomeResultsByTaskKind`, `HasImplementationSignal` (changes `RECORDED`/`COMPARED`), `HasVerificationDebt`, `countSatisfiedVerifications` |
| Changes (S02) | `PromoteVCSCommitToChange`, `GetChangeByGitCommit`, `ListChangesByTaskID` — implementation signal is task-scoped, not VCS-only |
| Test graph (S01) | `validates` edges, `ImpactWalkResult.AffectedTests`, `ListValidatesForSymbol` / `ListValidatesForFile`, `ListCrossLayerImports`, `FileLayer` — **no** `trace test` CLI |
| `internal/testrun` | **Absent** — no package, no `trace/test-runner.json` reader |
| Loop next/status | `BuildPolicyInputs` → `SelectNext` in `next.go` + `apply.go` `Status()`; packet exposes full `policy_inputs`; `statusBlocked` = blocking uncertainty \|\| open regression \|\| `VerificationIncomplete` only |
| Loop status extras | `verification_debt` section + `promotion_blocked` (eval regression advisory); **no** `verification_cycle` bundle |
| Reflections | `reflections` table (task_id only); `RecordReflection` / `ListReflectionsByTaskID`; **no** entity_link from reflection → evaluation outcome today |
| CLI root | `cmd/trace/root.go` — no `test`, `verify`, or `outcomes` subcommands |
| Help | `cmd/trace/help.go` — loop next/status documented; no test/verify lines |
| MCP catalog | **10** tools unchanged (C38 MCP half is S08) |
| `policy_test.go` | **Absent** under `internal/loop/` |

S01 closed **C01–C03, C07**; S02 closed **C04–C06, C25** — do not reopen in S03 prompts.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-03, D-22-04, D-22-20
- [WORK-MAP.md](../../WORK-MAP.md) W-09…W-13
- Coverage: [README.md](../../README.md) C09, C11–C15, C36, C38-CLI rows
- Prior cycle table: P21-S03-01 (`select.go` — **do not reorder**)

## FINAL locked defaults

| Item | Value |
|------|-------|
| SQL | **None** unless implementer proves a hole — reuse `outcome_results`, `changes`, `reflections`, `entity_links`, S01 `code_edges` / `imports` |
| Compat | Stays **23** (after S02); **forbid 024+** entire S03 scope |
| ExecutePending | `PlanExists` && `PlanCritiqued` && `BlockingUncertaintyCount==0` && !`OpenRegression` && !`HasImplementationSignal(task)` |
| TestPending | `HasImplementationSignal(task)` && !`HasTestOutcomeSinceLatestChange(task)` — latest change = max `created_at` among `changes` with `task_id` and `status IN (RECORDED, COMPARED)`; satisfied when ∃ `outcome_results` row `kind=test` for task with `created_at >=` that change timestamp |
| VerificationIncomplete | **Keep** existing `HasVerificationDebt` wiring (implementation + goal + no verified outcome with evidence) — maps to VERIFY row |
| EvaluationPending | (`CheckVerificationGate(task, "")` **or** ∃ any `kind=verification` row for task) && !∃ `kind=evaluation` with `comparisonComputed(comparison_json)` for task |
| ReflectPending | ∃ computed evaluation (`kind=evaluation` + `comparisonComputed`) for task && !∃ `reflections` row for task with `created_at >=` that evaluation’s `created_at` |
| ReplanNeeded | **Out of scope S03-01** — leave stub false unless S03-05 coordinator sets it from regression signal |
| Test invoke | New `internal/testrun`: detect `go.mod` → `go test ./...` with bounded timeout; optional `trace/test-runner.json` `{"command","args","cwd"}` override; always `domain.RecordTestOutcome`; **no daemon** (D-22-03) |
| Relevant tests | Primary: S01 `validates` from changed paths (`ListChangePaths` on latest change) + optional `ImpactWalk` on those paths → `AffectedTests`; fallback: Go package dirs of changed `.go` paths → `go test ./<pkg>` |
| Cycle gate | `loop status` adds `verification_cycle: {execute_pending, test_pending, verification_pending, evaluation_pending, reflect_pending, incomplete_reason}` derived from `PolicyInputs` + debt; `SelectNext` already blocks skip when flags true; **do not** rewrite `TransitionTask` DONE (D-22-20) |
| Coordinate scoring | `domain.CoordinateVerification(ctx, taskID, opts)` runs in order: testrun (if test pending) → ensure verification evidence → evaluation vs active baseline; CLI `trace verify run --task <id> [--force-eval]` |
| Invariants (C14) | `domain.CheckArchitecturalInvariants(taskID)` — for paths in latest change, join `imports` + `ListCrossLayerImports` / layer membership; default forbidden: **`internal/` importer → `cmd/` target** (one rule only; S07 extends via `trace/eval-rules.json`) |
| Iteration compare (C15) | `domain.CompareIterationOutcomes(taskID, kind)` — last two `outcome_results` of same `kind` for task: delta on `test_status` or score dimensions |
| CLI additions | `trace test run`, `trace verify run`, `trace verify invariants`, `trace outcomes compare` — register in `root.go` + `help.go`; capability keys `cli:test`, `cli:verify`, `cli:outcomes` AUTO_ALLOW pattern like `cli:changes` |

## Named tests

| Test | Row |
|------|-----|
| `TestBuildPolicyInputsSetsExecutePending` | S03-01 |
| `TestBuildPolicyInputsSetsTestPendingAfterChange` | S03-01 |
| `TestBuildPolicyInputsSetsEvaluationPending` | S03-01 |
| `TestBuildPolicyInputsSetsReflectPending` | S03-01 |
| `TestLoopNextExecuteWhenPendingLive` | S03-01 |
| `TestTestRunRecordsOutcome` | S03-03 |
| `TestTestRunSelectsValidatingTests` | S03-03 |
| `TestTestRunFailClosedWithoutCommand` | S03-03 |
| `TestVerificationCycleBlocksSkipInStatus` | S03-05 |
| `TestCoordinateVerificationOrder` | S03-05 |
| `TestRegressionDetectedVsPriorPassingTest` | S03-05 |
| `TestInvariantFailOnForbiddenLayerImport` | S03-07 |
| `TestCompareIterationOutcomes` | S03-07 |

## Keeper floor (do not break)

P21 domain gate tests: `TestTestPassAloneCannotSatisfyVerificationGate`, `TestPromotionGateRequiresStoredTestNotAgentClaim`, `TestVerificationDebtWhenImplementationWithoutVerification`, `TestEvaluationComparesScoresToBaselineNotBoolean`, `TestSelectNextNeverExecuteOnBlockingUncertainty`, deliberation ordering tests.

S01/S02 spot-check: `TestImpactWalkIncludesAffectedTests`, `TestPromoteVCSCommitCreatesChangeIdempotent`, `TestCompatibilitySecurityChecklist` (ceiling **23**).

## Residual risks for S03-01

| Risk | Mitigation locked in 01 |
|------|-------------------------|
| `TestPending` vs stale pre-change test rows | `HasTestOutcomeSinceLatestChange` compares outcome `created_at` to latest change only |
| `EvaluationPending` fires before verification debt cleared | `select.go` order: `VerificationIncomplete` (row 10) before `EvaluationPending` (row 11) — keep `HasVerificationDebt` wired |
| `ReflectPending` without evaluation→reflection FK | Use timestamp rule: reflection `created_at >=` latest computed evaluation; document in test |
| Execute never reached when plan critiqued but change exists | By design: implementation signal clears `ExecutePending`; TEST row takes over |
| `OpenDecisionAlternatives` still stubbed | Out of S03-01 scope (EXPLORE row); do not block C09 on it |
| Accidental schema bump | Grep: 23 sql files only; compat stays **23** |
| CLI next shows stub flags | `TestLoopNextExecuteWhenPendingLive` asserts JSON `policy_inputs.execute_pending` + `deliberation.phase` EXECUTE |

## Exit criteria

- [x] 01–08 thickened vs live `policy.go` / outcomes / loop status
- [x] No product Go

## Next

**P22-S03-01**
