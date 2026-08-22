# P20-S05-02 — Review regression / reflect / history

## Metadata
- id: P20-S05-02
- todo_ids: [P20-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: **correlated ≠ caused**; S03 contradiction and S04 evaluation flags are inputs only; reflection updates the graph with structured fields (not raw CoT essays); observed vs causal links; library-only.

## Session start
Follow agent-loop-protocol. Unattended after S05-01 `done`. Reviewer ≠ implementer session. Board: status/notes; spawn forward on blocker/high. Do not rewrite `done` prompts.

## Keeper tests (must re-run)

```bash
go test ./internal/domain/ -count=1 -run 'TestRecordRegression|TestCorrelation|TestLinkHypothesis|TestSetAttribution|TestHasOpenRegression|TestResolveRegression|TestReflection|TestObservedRelationship|TestCausalRelationship|TestUnknownAttribution'
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestRegression|TestReflection|TestNoSourceContentColumns'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review checklist

Compare S05-01 Notes + repo to [00-PLANNER.md](00-PLANNER.md) FINAL locks.

- [ ] `019_regressions_reflections.sql` exists; no 020+; embed/compat ceiling **19**
- [ ] Tables are `regressions` + `reflections` only — **no** Experiment / `relationships` / risk-matrix tables
- [ ] Attribution enum exactly `correlated` \| `hypothesized` \| `caused`; create defaults **`correlated`**
- [ ] Evaluation `comparison_json` flags and contradicted `effects` are **inputs**; S03/S04 schemas **unaltered**
- [ ] S03 contradiction does **not** persist `attribution=caused` (S03-02 residual intact)
- [ ] `RecordEvaluationOutcome` still does **not** auto-insert regressions
- [ ] `LinkHypothesisToRegression` → `hypothesized` **not** `caused`; Discovery is not used as Hypothesis
- [ ] `ConfirmHypothesis` does **not** auto-set `caused`
- [ ] `SetRegressionAttributionCaused` fail-closed without evidence, from `correlated`, or with OPEN (not CONFIRMED) hypothesis
- [ ] `caused_by` link requires evidence and does **not** flip `regressions.attribution`
- [ ] `observed_relationship` allows confidence without evidence
- [ ] Reflections have structured JSON arrays; **no** `body`/essay-only path; `TestReflectionEssayOnlyFailClosed` exists
- [ ] `HasOpenRegression` matches OPEN status; feeds `ApplyDeliberationTransition` (not SelectNext alone); **no** auto-hop on record
- [ ] `internal/deliberation/select.go` **untouched**
- [ ] §16/§18 not implemented (`broaden_tests_note` is a stub string only)
- [ ] No P19 loop/CLI/MCP edits; no CoT/blobs
- [ ] Law 19: library-only this scope

## Spawn rule

blocker/high: small inline fix **or** insert `P20-S05-02a` (implement) + `P20-S05-02b` (review) immediately below this row with full prompts. Medium: prefer spawn unless trivial.

## Exit criteria

- blocker/high fixed or spawned forward
- confidence medium or high with evidence
- residuals listed explicitly if medium (never silent)
- Next runnable after APPROVE: **P20-S06-00**
