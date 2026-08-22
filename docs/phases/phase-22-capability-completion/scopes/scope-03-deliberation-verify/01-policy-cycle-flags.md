# P22-S03-01 — Implement: store-driven cycle flags

## Metadata
- id: P22-S03-01
- todo_ids: [P22-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Wire **`BuildPolicyInputs`** so EXECUTE/TEST/EVALUATE/REFLECT are store-driven. Closes **C09**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/loop/policy.go` (cycle flags **never set** today), `internal/deliberation/select.go` (14-row table — **do not reorder**)
- Domain: `internal/domain/outcomes.go` — `HasImplementationSignal`, `HasVerificationDebt`, gate helpers, `ListOutcomeResultsByTaskKind`
- Store: `internal/store/outcomes.go` — `HasImplementationSignal`, `HasVerificationDebt`
- Reflections: `internal/store/regressions.go` — `ListReflectionsByTaskID`
- Schema max **023**; compat **23** — **no migration this row**

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **None** — reuse existing tables |
| Compat | Stays **23** (forbid 024+) |
| `select.go` | **Do not change** priority order (P21-S03-01 FINAL) |
| `VerificationIncomplete` | Keep existing `dom.HasVerificationDebt` assignment |
| ExecutePending | Plan exists + critiqued + no blocking uncertainty + no open regression + **no** implementation signal |
| TestPending | Implementation signal + no `kind=test` outcome since latest RECORDED/COMPARED change |
| EvaluationPending | Verification gate pass **or** any verification row + no computed evaluation |
| ReflectPending | Computed evaluation exists + no reflection at/after that evaluation timestamp |
| `ReplanNeeded` | Leave **false** (S03-05 owns regression→replan if needed) |
| `OpenDecisionAlternatives` | Leave **0** (not C09; future scope) |

## Locked query semantics

Implement helpers on `domain.Service` (preferred) or `store.Store` if SQL-only:

1. **`HasTestOutcomeSinceLatestChange(taskID)`** — find latest change for task (`status IN (RECORDED, COMPARED)` by `created_at` desc); return true if any `outcome_results` `kind=test` for task has `created_at >=` change `created_at`. No change → false.
2. **`HasComputedEvaluation(taskID)`** — reuse `comparisonComputed` logic from outcomes.go (non-empty `dimensions` in `comparison_json`).
3. **`HasReflectionSinceEvaluation(taskID)`** — latest computed evaluation by `created_at`; true if any reflection for task has `created_at >=` evaluation `created_at`.

Wire all four bools + existing fields in `BuildPolicyInputs` return struct.

## Requirements

1. Live queries in `internal/loop/policy.go` calling domain helpers.
2. **`trace loop next`** and **`trace loop status`** must expose non-stub flags in `deliberation.policy_inputs` and correct `recommended_phase` / `phase` when seeded (not unit-only).
3. `applyWrites` path unchanged — plan_changes still force `plan_critiqued` before read.
4. Named tests below + keepers.

## Touch files

- `internal/loop/policy.go`
- `internal/loop/policy_test.go` (**new**)
- `internal/domain/outcomes.go` (+ `_test.go` for new helpers)
- `internal/store/outcomes.go` (only if SQL helper needed)
- `cmd/trace/loop_test.go` (CLI keeper)

## Named tests

| Test | Proves |
|------|--------|
| `TestBuildPolicyInputsSetsExecutePending` | Plan+critique, no change → `execute_pending=true` |
| `TestBuildPolicyInputsSetsTestPendingAfterChange` | RECORDED change, no test since → `test_pending=true` |
| `TestBuildPolicyInputsSetsEvaluationPending` | Verified + no eval → `evaluation_pending=true` |
| `TestBuildPolicyInputsSetsReflectPending` | Eval computed + no reflection → `reflect_pending=true` |
| `TestLoopNextExecuteWhenPendingLive` | CLI `loop next` → `policy_inputs.execute_pending` + phase EXECUTE |
| `TestSelectNextNeverExecuteOnBlockingUncertainty` | keeper (deliberation package) |
| `TestSelectNextExecuteWhenPending` | keeper |
| `TestVerificationDebtWhenImplementationWithoutVerification` | keeper — debt still wired |

```bash
go test ./internal/loop/... ./internal/deliberation/... ./internal/domain/... -count=1 -run 'TestBuildPolicyInputs|TestSelectNext|TestVerificationDebtWhenImplementationWithoutVerification'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopNextExecute|TestLoopStatusDeliberationFields|TestLoopNextDeliberationSectionPresent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] C09 true (cycle flags store-driven, not zero stubs)
- [ ] Named tests PASS; compat **23**
- [ ] Checklist box **not** checked until S03-02 review closes C09
- [ ] Board Notes: test output

## Minimal todos

- [ ] Domain helpers for test-since-change / eval / reflection
- [ ] Wire `BuildPolicyInputs` four flags
- [ ] Unit + CLI tests
- [ ] Board status + notes
