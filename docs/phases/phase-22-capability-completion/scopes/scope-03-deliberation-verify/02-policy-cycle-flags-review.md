# P22-S03-02 — Review: cycle flags

## Metadata
- id: P22-S03-02
- todo_ids: [P22-S03-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C09**: live `BuildPolicyInputs` sets all four cycle flags with store-backed queries; `select.go` first-match order unchanged.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. Grep `internal/loop/policy.go` — `ExecutePending`, `TestPending`, `EvaluationPending`, `ReflectPending` must be **assigned from live queries**, not omitted / always false.
2. Grep `internal/deliberation/select.go` — **no priority reorder** vs P21 FINAL table.
3. Confirm locked semantics from [00-PLANNER.md](00-PLANNER.md):
   - Test pending uses **since latest change**, not any historical test row.
   - Reflect pending uses evaluation timestamp vs reflection `created_at`.
4. Schema: **23** sql files, no 024+; compat test PASS.
5. CLI: `trace loop next` JSON shows true flags when seeded (not test-only wiring).
6. Keepers green: `TestSelectNextNeverExecuteOnBlockingUncertainty`, `TestVerificationDebtWhenImplementationWithoutVerification`.

## Spawn policy

If flags still stubbed or semantics wrong: spawn **`P22-S03-02a` + `P22-S03-02b`** immediately below. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/loop/... ./internal/deliberation/... ./internal/domain/... -count=1 -run 'TestBuildPolicyInputs|TestSelectNextNeverExecuteOnBlockingUncertainty|TestVerificationDebtWhenImplementationWithoutVerification'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopNextExecute|TestLoopStatusDeliberationFields'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 23
```

## Exit criteria

- [ ] C09 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: findings + confidence
