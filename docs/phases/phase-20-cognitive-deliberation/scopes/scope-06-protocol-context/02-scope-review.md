# P20-S06-02 — Review protocol + context

## Metadata
- id: P20-S06-02
- todo_ids: [P20-S06-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: harness-agnostic bounded packets (§24); phase-aware context (§23); structured apply round-trip (§22); P19 backward compatibility; fail-closed apply; no daemon/MCP expansion.

## Session start
Follow agent-loop-protocol. Unattended after S06-01 `done`. Reviewer ≠ implementer session. Board: status/notes; spawn forward on blocker/high. Do not rewrite `done` prompts.

## Keeper tests (must re-run)

```bash
go test ./internal/loop/... -count=1
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopNextDeliberation|TestLoopNextInvestigate|TestLoopNextExecute|TestLoopNextVerify|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyUnknownWriteKey|TestLoopApplyUncertainty|TestLoopApplyRegression|TestLoopApplyDeliberationTransition|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestLoopStatusDeliberation|TestLoopStatusBlocked|TestLoopApplyNoPartialWrites|TestLoopRecentChanges|TestHelpIncludesLoopNext'
go test ./internal/deliberation/...
go test ./internal/domain/ -count=1 -run 'TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestHasOpenRegressionFeedsApplyDeliberationTransition|TestApplyDeliberationTransition'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 subset (must remain green):

```bash
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
```

## Review checklist

Compare S06-01 Notes + repo to [00-PLANNER.md](00-PLANNER.md) FINAL locks.

### Schema + compat

- [ ] `schema_version` remains `trace.loop.next.v1` / `apply.v1` / `status.v1` (additive, not v2)
- [ ] No SQL migration **020**; embed/compat ceiling still **19**
- [ ] P19 keeper tests pass without weakening assertions

### next packet (§22–23)

- [ ] `deliberation` section: recommended phase, `why_selected`, full `policy_inputs`, hop/stop fields
- [ ] `open_uncertainties`, `verification_debt`, `recent_changes` bounded per caps
- [ ] `recent_changes` / context: **no** file bytes, patches, or blobs (Law 1)
- [ ] INVESTIGATE emphasizes uncertainties; EXECUTE emphasizes context/related; VERIFY surfaces debt
- [ ] `PolicyInputs.blocking_uncertainty_count` from live `CountBlockingUncertainties` (not stubbed)
- [ ] `PolicyInputs.open_regression` from live `HasOpenRegression` (not stubbed)
- [ ] `SelectNext` not called in isolation to burn hop budget on next-only path

### apply (§22)

- [ ] Allowlisted `writes` keys only; unknown keys fail closed **before** writes
- [ ] New arrays map to S02–S05 domain APIs (no duplicate SQL in loop)
- [ ] `ApplyDeliberationTransition` runs after writes with **complete** PolicyInputs
- [ ] `deliberation.transition` event payload includes `to_phase`, `reason_code`, `policy_inputs`
- [ ] Replay same `apply_id`: idempotent artifacts + single transition semantics
- [ ] Zero-delta apply still runs transition (hop policy unchanged from S01)
- [ ] Test ≠ verify ≠ evaluate apply keys distinct (Law 2)

### status

- [ ] `deliberation` block: phase, recommended, blocked, needs_phase, hop_count
- [ ] P19 saturation fields unchanged in meaning

### Boundaries

- [ ] `internal/deliberation/select.go` **untouched**
- [ ] **No new MCP tools**
- [ ] No raw CoT storage; no hosted daemon
- [ ] §18 `broaden_tests_note` is stub only — not a test-selection engine
- [ ] S07 seed export still out of scope (note residual only)

## Spawn rule

blocker/high: small inline fix **or** insert `P20-S06-02a` (implement) + `P20-S06-02b` (review) immediately below this row with full prompts. Medium: prefer spawn unless trivial.

## Exit criteria

- blocker/high fixed or spawned forward
- confidence medium or high with evidence
- residuals listed explicitly if medium (never silent)
- Next runnable after APPROVE: **P20-S07-00**
