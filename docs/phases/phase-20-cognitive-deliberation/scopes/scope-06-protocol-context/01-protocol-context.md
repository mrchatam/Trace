# P20-S06-01 — Implement protocol + context selection

## Metadata
- id: P20-S06-01
- todo_ids: [P20-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Extend `internal/loop` + `cmd/trace/loop` so `next` / `apply` / `status` carry deliberation phase, bounded cognitive context, and S02–S05 artifact apply keys per **S06-00 FINAL**. **Product Go this row.**

## Session start
Follow agent-loop-protocol. Unattended: execute after S06-00 is `done`. Board edits: **status + notes only**. Do not re-debate locks. **Do not edit** `internal/deliberation/select.go` or SQL migrations 001–019.

## Locked defaults (from S06-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Schema | **Additive on v1** — `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` unchanged strings |
| SQL migration | **None** — compat ceiling stays **19** |
| Deliberation hop | `domain.ApplyDeliberationTransition` on **every** apply (including zero-delta); hop increments there only |
| PolicyInputs | Complete struct every hop from live queries — see 00-PLANNER assembly table |
| Blocking count | `domain.CountBlockingUncertainties` — never stub 0 when rows exist |
| Open regression | `domain.HasOpenRegression` — never stub false when OPEN rows exist |
| Verification debt | `domain.HasVerificationDebt` + `ListVerificationDebtSummary` for packet |
| Unknown apply keys | Fail closed **before** any domain write |
| recent_changes | SHA + path + comparison only — **no file bytes** (Law 1) |
| MCP | **No new tools** |
| SelectNext | Pure policy — consume via transition only |

Full packet shapes, apply key allowlist, write→domain mapping, phase context table: [00-PLANNER.md](00-PLANNER.md) FINAL section. Copy them; do not invent fields.

## Requirements

1. **`BuildNextPacket`** — after P19 sections built:
   - Load `deliberation_state` (or ORIENT/hop 0 initial)
   - Build `PolicyInputs` via shared helper
   - Run `deliberation.SelectNext` → populate `deliberation` section (`phase`, `why_selected`, `policy_inputs`, `hop_count`, `stopped`, `stop_reason`, `last_phase`)
   - Populate `open_uncertainties`, `verification_debt`, `recent_changes` with bounded store/domain queries
   - Apply phase-specific context emphasis (compiler `ContextOptions` + section caps per 00-PLANNER table)
2. **`ParseApplyEnvelope` / `ValidateApplyEnvelope`** — reject unknown `writes` keys; extend struct for new arrays
3. **`Apply`** — process new write kinds via existing domain APIs; then `ApplyDeliberationTransition`; then P19 `loop.step.applied` event. Replay: idempotent on same `apply_id` (no duplicate transition event)
4. **`Status`** — add `deliberation` block with `phase`, `recommended_phase`, `why_selected`, `hop_count`, `stopped`, `blocked`, `needs_phase`, `policy_inputs`
5. **Shared helper** `BuildPolicyInputs(ctx, dom, planner, taskID, goalID, applyWrites, p19Saturated)` in `internal/loop`
6. Store list helpers if missing: `ListOpenUncertaintiesByTaskID(taskID, limit)`, bounded recent changes packet builder (may live in loop package using existing `ListChangesByTaskID` + paths/effects getters)

## Apply write shapes (minimum)

```text
ApplyUncertainty:  id, title, body?, severity?, status?, kind?, task_id? (required when BLOCKING)
ApplyHypothesis:   id, title, body?, evidence_ids?, uncertainty_id?
ApplyChange:       id, git_commit?, parent_change_id?, actor?, reason?, paths[], expected[]?, decision_id?
ApplyEffect:       change_id, dimension, expected?, actual?, comparison?, evidence_ids?, hypothesis_id?
ApplyTestResult:   id, test_name, test_status, summary?, evidence_ids?
ApplyVerification: id, goal_id, verification_status, evidence_ids (≥1), summary?
ApplyEvaluation:   id, baseline_id, scores_json
ApplyRegression:   source_kind (evaluation|contradicted_effect), source_id, task_id, summary?
ApplyReflection:   summary?, invalidated_assumption_ids?, new_dependencies?, useful_tests?, broaden_tests_note?
```

All IDs UUID. Fail closed on unknown enums (reuse domain validation).

## Named tests (must exist and pass)

See 00-PLANNER list (14 names). Minimum proofs:

- `loop next` JSON includes `deliberation.phase` + `deliberation.why_selected` + full `policy_inputs`
- Blocking uncertainty in store → next recommends INVESTIGATE / `blocking_uncertainty`; `open_uncertainties` non-empty
- EXECUTE recommendation → context packet has file/symbol items; related depth 2 when seeds exist
- VERIFY recommendation + debt → `verification_debt.present=true`
- Apply with extra `writes.unknown_key` → exit fail, zero new rows in store
- Apply uncertainty write → subsequent next shows increased blocking count / INVESTIGATE
- Apply with regression write → `policy_inputs.open_regression` true on next
- Apply emits `deliberation.transition` event with `to_phase` + `reason_code` + `policy_inputs`
- Replay same apply_id → single transition event, `replay:true`
- Status includes deliberation fields; `blocked:true` when blocking uncertainty open
- Validation failure mid-envelope → no partial discoveries/uncertainties persisted
- `recent_changes` entries contain no `content`/`patch`/`blob` keys

## Likely touch points

- `internal/loop/next.go` — deliberation + bounded sections + phase emphasis
- `internal/loop/apply.go` — writes extension, unknown-key guard, transition hook
- `internal/loop/*_test.go` — unit coverage
- `internal/store/cognitive.go` — `ListOpenUncertaintiesByTaskID` if absent
- `cmd/trace/loop.go` — pass `domain.Service` into apply if needed for transition
- `cmd/trace/loop_test.go` — integration named tests

Do **not** touch: `internal/deliberation/select.go`, `internal/mcp`, `internal/store/schema/*`, S02–S05 domain policy files (except constants if unavoidable), seed export (S07).

## Proof commands

```bash
go test ./internal/loop/... -count=1
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopNextDeliberation|TestLoopNextInvestigate|TestLoopNextExecute|TestLoopNextVerify|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyUnknownWriteKey|TestLoopApplyUncertainty|TestLoopApplyRegression|TestLoopApplyDeliberationTransition|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestLoopStatusDeliberation|TestLoopStatusBlocked|TestLoopApplyNoPartialWrites|TestLoopRecentChanges|TestHelpIncludesLoopNext'
go test ./internal/deliberation/...
go test ./internal/domain/ -count=1 -run 'TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestHasOpenRegressionFeedsApplyDeliberationTransition|TestHasVerificationDebt'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 keeper subset must stay green:

```bash
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
```

Compat ceiling **19** must remain (no 020 migration).

## Todo updates
Status + notes only on `P20-S06-01`. Next after green: `P20-S06-02`.

## Exit criteria

- [ ] Additive v1 schemas; P19 keeper tests green
- [ ] All 14 named S06 tests green
- [ ] Unknown apply write keys fail closed with no partial writes
- [ ] `ApplyDeliberationTransition` on apply with complete PolicyInputs
- [ ] Phase-specific context emphasis per 00-PLANNER table
- [ ] No file bytes in packets; no new MCP tools; no mig 020
- [ ] Law 19: extend loop/CLI only — no duplicate domain policy
