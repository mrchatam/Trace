# P20-S02-02 — Review cognitive artifacts

## Metadata
- id: P20-S02-02
- todo_ids: [P20-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: merge discipline holds; blocking uncertainty is queryable for S01 SelectNext; assumption invalidation is an explicit state transition (Law 11); decision reconsideration is append-only; no raw CoT; no Finding/Requirement/Constraint forks.

## Session start
Follow agent-loop-protocol. Unattended after S02-01 `done`. Reviewer ≠ implementer session. Board: status/notes; spawn forward on blocker/high. Do not rewrite `done` prompts.

## Keeper tests (must re-run)

```bash
go test ./internal/domain/ -count=1 -run 'TestCreateUncertainty|TestBlocking|TestResolve|TestSupersede|TestCountBlocking|TestInvalidateAssumption|TestHypothesis|TestDecisionReconsider|TestUnknownUncertainty'
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review checklist

Compare S02-01 Notes + repo to [00-PLANNER.md](00-PLANNER.md) FINAL locks.

- [ ] COVERAGE.md merge table still accurate (Requirement/Constraint/Finding/Risk/Option **not** new tables)
- [ ] `016_cognitive_artifacts.sql` exists; no 017+; embed/compat ceiling **16**
- [ ] `uncertainties` distinct from Discovery (no reuse of `discoveries` as questions; severity INFO\|BLOCKING only)
- [ ] Blocking count SQL matches lock (`BLOCKING`+`OPEN`+`uncertainty_blocks_task`); INFO does not count; resolve/supersede clears count
- [ ] `TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition` passes complete `PolicyInputs` into `ApplyDeliberationTransition` (does not call `SelectNext` alone)
- [ ] Assumption invalidate: row remains; status STALE/SUPERSEDED only; event `assumption.invalidated`; no DELETE
- [ ] Linked decision: `INVALIDATED_ASSUMPTION` finding + FIRED reconsideration; no auto-replan
- [ ] Optional Discovery is `PLAN_AFFECTING`, not a second Finding type
- [ ] Decision + DecisionAlternative rows survive reconsideration
- [ ] Hypothesis evidence via `hypothesis_supported_by` → evidence (not a Discovery standing in as hypothesis)
- [ ] No P19 loop/CLI/MCP edits; no CoT/blobs
- [ ] Law 19: library-only this scope
- [ ] Residuals listed if seed/FTS omitted (expected per S02-00)

## Spawn rule

blocker/high: small inline fix **or** insert `P20-S02-02a` (implement) + `P20-S02-02b` (review) immediately below this row with full prompts. Medium: prefer spawn unless trivial.

## Exit criteria

- blocker/high fixed or spawned forward
- confidence medium or high with evidence
- residuals listed explicitly if medium (never silent)
- Next runnable after APPROVE: **P20-S03-00**
