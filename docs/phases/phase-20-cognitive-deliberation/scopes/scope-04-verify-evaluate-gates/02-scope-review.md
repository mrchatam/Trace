# P20-S04-02 — Review verify / evaluate / gates

## Metadata
- id: P20-S04-02
- todo_ids: [P20-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: Law 2 enforcement (implemented ≠ verified; test-pass ≠ verification; evaluation distinct from boolean review); thin baselines; verification debt query; no work_state explosion; no test-runner product.

## Session start
Follow agent-loop-protocol. Unattended after S04-01 `done`. Reviewer ≠ implementer session. Board: status/notes; spawn forward on blocker/high. Do not rewrite `done` prompts.

## Keeper tests (must re-run)

```bash
go test ./internal/domain/ -count=1 -run 'TestRecordTestOutcome|TestTestPassAlone|TestVerification|TestEvaluation|TestBaseline|TestVerificationDebt|TestPromotionGate|TestPartialVerification'
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestOutcome|TestBaseline|TestNoSourceContentColumns'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review checklist

Compare S04-01 Notes + repo to [00-PLANNER.md](00-PLANNER.md) FINAL locks.

- [ ] `018_outcome_results_baselines.sql` exists; no 019+; embed/compat ceiling **18**
- [ ] Tables are `outcome_results` + `baselines` only — **no** Experiment/risk-matrix/runner tables
- [ ] `kind` enum exactly `test` \| `verification` \| `evaluation` on single table
- [ ] **No** tests/baseline/score columns on `changes` (S03 boundary intact)
- [ ] **No** new `tasks.work_state` enum values
- [ ] Test outcomes require `test_name` + stored row; gate helpers fail without record (Law 2)
- [ ] `CheckTestGate` true does **not** satisfy `CheckVerificationGate`
- [ ] Verification requires `goal_id` + ≥1 `outcome_supported_by` evidence link; fail-closed if missing
- [ ] Evaluation uses library-computed `comparison_json` — not agent boolean PASS (Law 15)
- [ ] Baseline = git OID + `scores_json` only; no log blobs (Law 1 spirit)
- [ ] `HasVerificationDebt` matches locked definition (implementation signal + missing verified outcome)
- [ ] `partial` verification counts as debt
- [ ] DONE policy **unchanged** — Review PASS + operator; debt does not silently authorize DONE
- [ ] No subprocess test runner / pytest / go test wrapper embedded
- [ ] §16/§18 not implemented (Future only)
- [ ] No P19 loop/CLI/MCP/SelectNext edits; no CoT/blobs
- [ ] Law 19: library-only this scope

## Spawn rule

blocker/high: small inline fix **or** insert `P20-S04-02a` (implement) + `P20-S04-02b` (review) immediately below this row with full prompts. Medium: prefer spawn unless trivial.

## Exit criteria

- blocker/high fixed or spawned forward
- confidence medium or high with evidence
- residuals listed explicitly if medium (never silent)
- Next runnable after APPROVE: **P20-S05-00**
