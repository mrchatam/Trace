# P20-S06-00 — Protocol + context planner

## Metadata
- id: P20-S06-00
- todo_ids: [P20-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock phase-aware `trace loop next` / `apply` / `status` extensions and bounded context selection by deliberation phase. Preserve P19 backward compatibility. **No product Go this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [COVERAGE.md](../../COVERAGE.md) §§22, 23, 24, 29E
- Laws **1**, **2**, **6**, **7**, **19** ([G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md))
- Live P19: `internal/loop/{next,apply}.go`, `cmd/trace/loop.go`, `cmd/trace/loop_test.go`
- Live S01: `internal/deliberation/{select,types}.go`, `internal/domain/deliberation.go` (`ApplyDeliberationTransition`)
- Live S02–S05 queries: `CountBlockingUncertainties`, `HasVerificationDebt`, `HasOpenRegression`
- Live compiler: `internal/compiler` (`TaskContext`, `ContextOptions`)
- Schema max **019** (S05); **no SQL migration 020** for S06 — protocol-only scope

## Doc map
§22 Agent interaction (structured ask/return via loop), §23 Context selection (phase-bounded), §24 Model/harness agnostic (stdout JSON inherit P19), §29E protocol

## Live inventory (2026-08-18)

| Surface | Location | S06 action |
|---------|----------|------------|
| `trace.loop.next.v1` | `internal/loop/next.go` | **Extend** — additive sections; schema_version **unchanged** |
| `trace.loop.apply.v1` | `internal/loop/apply.go` | **Extend** `writes` keys; **fix** unknown-key silent ignore |
| `trace.loop.status.v1` | `internal/loop/apply.go` `Status` | **Extend** — additive `deliberation` block |
| P19 keeper tests | `cmd/trace/loop_test.go` | Must stay green unchanged assertions on v1 + legacy fields |
| `SelectNext` / `ApplyTransition` | `internal/deliberation/select.go` | **Consume only** — do not edit policy table |
| `ApplyDeliberationTransition` | `internal/domain/deliberation.go` | **Call on apply path** after writes + before loop.step event |
| `PolicyInputs` | `internal/deliberation/types.go` | **Wire live queries every hop** — no partial stubs |
| Artifact domain APIs | S02–S05 `internal/domain/*` | **Reuse** for apply write handlers — no duplicate SQL in loop |
| Store list helpers | partial | Add bounded `ListOpenUncertaintiesByTaskID` / `ListRecentChangesPacket` in Go if missing — **no mig 020** |
| Compat ceiling | **19** after S05 | **Unchanged** — S06 is not a schema migration scope |
| MCP | G19 | **No new MCP tools** default |

## Schema evolution decision (FINAL)

| Option | Verdict |
|--------|---------|
| Dual-version `trace.loop.next.v2` / `apply.v2` / `status.v2` | **Rejected** — P19 keepers hard-check `schema_version == *.v1`; dual-version adds harness churn with no breaking need |
| Additive fields on existing **v1** schemas | **Locked** |

**Rationale:** P19 tests assert `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` strings and required legacy sections. Adding optional sibling JSON fields preserves backward compatibility for old harnesses (ignore unknown fields). The only **behavior** change on apply is fail-closed unknown `writes` keys — valid P19 envelopes (exact allowlisted keys) remain accepted.

**No migration 020:** All artifact tables exist (015–019). S06 may add store **query** helpers only.

## FINAL locked defaults (S06-01 must not re-debate)

### Transport + MCP

| Item | Value |
|------|-------|
| Transport | Stdout-first JSON; harness-agnostic (inherit P19) |
| MCP | **No new MCP tools** — G19 library inherit |
| Compat ceiling | **19** (forbid `020+` SQL); S06 does not bump embed max |

### `loop next` packet additions (additive on `trace.loop.next.v1`)

New top-level sections (each with `freshness` like P19 sections):

| Section | Fields | Bounds |
|---------|--------|--------|
| `deliberation` | `phase` (recommended from `SelectNext`), `why_selected` (`reason_code`), `policy_inputs` (full struct), `hop_count`, `stopped`, `stop_reason`, `last_phase` | `phase`/`why_selected` from pure `SelectNext(state, inputs)` — **does not** increment hop (apply only) |
| `open_uncertainties` | `items[]`: `{id, title, severity, status, kind}` | Max **16** OPEN rows task-scoped (BLOCKING first, then INFO) |
| `verification_debt` | `{present: bool, items[]: DebtItem}` | Max **8** items via `ListVerificationDebtSummary` |
| `recent_changes` | `items[]`: `{id, git_commit, status, paths[], effects[]}` | Max **8** changes, newest first; paths max **16** each; effects max **8** each; **comparison enum only** — **no file bytes** (Law 1) |

Existing P19 sections (`seed`, `tasks`, `plan`, `why`, `context`, `related`, `loop_hints`) **unchanged** in shape. Phase-specific **emphasis** adjusts compiler/options and caps (see Context table below) — does not remove sections.

**PolicyInputs assembly (next + apply + status):**

```text
BlockingUncertaintyCount ← domain.CountBlockingUncertainties(ctx, taskID)
PlanExists               ← planner.GetPlan goal has current_scope_id + current_deep_plan (same gate as next preflight)
PlanCritiqued            ← deliberation_state.plan_critiqued OR (apply path: len(writes.plan_changes)>0 before transition)
VerificationIncomplete   ← domain.HasVerificationDebt(ctx, taskID)
OpenRegression           ← domain.HasOpenRegression(ctx, taskID)
P19Saturated             ← apply result only (next/status: false unless last apply said saturated)
```

**Never stub** `BlockingUncertaintyCount=0` when S02 rows exist; **never stub** `OpenRegression=false` when OPEN regressions exist.

### `loop apply` writes extension (additive on `trace.loop.apply.v1`)

**Allowlisted `writes` keys (exact):**

```text
discoveries, plan_changes, spawned_tasks, stop          # P19 legacy
uncertainties, hypotheses, changes, effects             # S02–S03
test_results, verifications, evaluations                # S04 (not outcome_results)
regressions, reflections                                # S05
```

**Fail-closed:** `ParseApplyEnvelope` must inspect raw `writes` object keys. Any key ∉ allowlist → error **before** partial domain writes. Today extra keys silently drop — **S06 fixes this** (P19 valid envelopes unaffected).

**Apply order (locked):**

1. Validate envelope + unknown keys
2. P19 writes: discoveries → plan_changes (+ replan) → spawned_tasks
3. S02–S05 writes: uncertainties → hypotheses → changes → effects → test_results → verifications → evaluations → regressions → reflections
4. Build complete `PolicyInputs` from **post-write** store state (+ `P19Saturated` from step result)
5. `domain.ApplyDeliberationTransition(ctx, seed.task_id, seed.goal_id, inputs)` — persists hop + appends `deliberation.transition`
6. Append `loop.step.applied` event (P19 replay semantics unchanged)

**Do not** call `SelectNext` alone to consume hop budget. **Do not** skip transition when writes are empty (zero-delta apply still advances deliberation per policy).

Each new write array element requires caller-supplied UUID `id` (same as P19 discoveries). Map to domain inputs:

| Apply key | Domain API | Notes |
|-----------|------------|-------|
| `uncertainties[]` | `CreateUncertainty` / resolve via status field | BLOCKING requires `task_id` = seed task |
| `hypotheses[]` | `CreateHypothesis` | Optional `evidence_ids`, `uncertainty_id` |
| `changes[]` | `CreateChange` + `RecordChangeCommit` when `git_commit` set | Paths required; no content blobs |
| `effects[]` | `RecordExpectedEffect` / `RecordActualEffect` on `change_id` | Fail closed unknown comparison |
| `test_results[]` | `RecordTestOutcome` | `test_name` + `test_status` required |
| `verifications[]` | `RecordVerificationOutcome` | `goal_id` + ≥1 `evidence_ids` |
| `evaluations[]` | `RecordEvaluationOutcome` | `baseline_id` + `scores_json`; comparison computed |
| `regressions[]` | `RecordRegressionFromEvaluation` / `FromContradictedEffect` | Create always `correlated`; no override |
| `reflections[]` | `CreateReflection` | Structured arrays only; essay-only fail closed |

Optional `links[]` on items: reuse `ApplyLink` shape; fail closed unsupported rels (same pattern as P19).

### `loop status` additions (additive on `trace.loop.status.v1`)

Add `deliberation` object (always present when task loads):

```text
phase            string   # persisted deliberation_state.current_phase
recommended_phase string  # SelectNext output (same as next packet)
why_selected     string   # reason_code
hop_count        int
stopped          bool
blocked          bool     # blocking_uncertainty OR open_regression OR verification_incomplete
needs_phase      string   # recommended_phase when blocked && !stopped && !p19_saturated; else ""
policy_inputs    object   # full snapshot
```

P19 fields (`saturated`, `reason`, `last_apply_id`, deltas) **unchanged**.

### Phase-specific context slice rules (FINAL)

Phase for emphasis = **recommended** phase from `SelectNext` (not persisted alone). Adjust `BuildNextPacket` compiler call + supplemental section caps:

| Recommended phase | Emphasis | Compiler / packet behavior |
|-------------------|----------|----------------------------|
| **INVESTIGATE** | Questions / uncertainties | `open_uncertainties` at max cap **16**; `context` `MaxItems=24`, `IncludeWhy=true`; `related` depth **1**; trim `plan.lookahead_summary` to **256** chars in packet mirror |
| **EXECUTE** | Files / impact | `context` `IncludeMarkdown=true`, `MaxItems=32`; `related` full walk depth **2**; `open_uncertainties` cap **4** |
| **VERIFY** | Requirements / evidence / debt | `verification_debt` populated; `why` `IncludeWhy=true`; `context` `MaxItems=20`; surface goal-linked evidence refs in `why.impact` (existing path) |
| **TEST** / **EVALUATE** | Recent work signals | `recent_changes` at max cap **8** with effect comparisons; default context |
| **REFLECT** / **REPLAN** | Learning + regressions | Include bounded open regressions summary (max **8**) in new optional `open_regressions` subsection under `deliberation` or top-level `open_regressions` (max **8** `{id, summary, attribution, status}`) |
| **ORIENT** / **EXPLORE** / **PLAN** / **CRITIQUE** | Plan + why default | P19 defaults; `open_uncertainties` cap **8** |

**Law 6–7:** No section may dump full graph. Hard per-packet JSON cap target **512 KiB** — if exceeded, truncate lowest-priority section items (never drop `deliberation` or `seed`).

**§18 stub:** Reflection `broaden_tests_note` may appear in apply writes; **not** a test-selection policy engine.

### P19 keeper tests (must pass after S06-01)

```bash
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
```

### Named tests (minimum — exact names for S06-01)

1. `TestLoopNextDeliberationSectionPresent`
2. `TestLoopNextPolicyInputsLiveQueries`
3. `TestLoopNextInvestigateEmphasizesUncertainties`
4. `TestLoopNextExecuteEmphasizesContextAndRelated`
5. `TestLoopNextVerifySurfacesVerificationDebt`
6. `TestLoopApplyUnknownWriteKeyFailsClosed`
7. `TestLoopApplyUncertaintyWriteAffectsNextSelectNext`
8. `TestLoopApplyRegressionWriteAffectsPolicyInputs`
9. `TestLoopApplyDeliberationTransitionEvent`
10. `TestLoopApplyReplaySkipsDuplicateTransition`
11. `TestLoopStatusDeliberationFields`
12. `TestLoopStatusBlockedWhenBlockingUncertainty`
13. `TestLoopApplyNoPartialWritesOnValidationFailure`
14. `TestLoopRecentChangesNoFileBytes`

### Files (S06-01)

| Path | Role |
|------|------|
| `internal/loop/next.go` | Deliberation sections + phase context emphasis + PolicyInputs builder |
| `internal/loop/apply.go` | Extended writes, unknown-key guard, ApplyDeliberationTransition hook |
| `internal/loop/apply_test.go` / `next_test.go` | Unit tests (new or extend) |
| `internal/store/cognitive.go` or `changes.go` | Bounded list helpers if missing |
| `cmd/trace/loop.go` | Wire domain service for transition (minimal) |
| `cmd/trace/loop_test.go` | Named integration tests above |

Do **not** edit: `internal/deliberation/select.go`, S02–S05 domain policy tables, `internal/mcp`, SQL migrations 001–019.

### Later scopes (upcoming notes only)

- **S07:** seed export of S02–S05 entities; VERIFY §31 mini-eval uses loop next; DR-HANDOFF close
- **FTS:** new artifact types may need FTS sync — residual if apply creates entities without SyncEntityFTS (note in review, do not block S06 on export)

## Planner work

1. [x] Decide additive vs v2 — **additive v1 locked**
2. [x] Lock apply fail-closed unknown-key behavior
3. [x] Lock per-phase context slice rules
4. [x] Thicken `01-protocol-context.md` + `02-scope-review.md`
5. [x] List P19 keeper tests

## Exit criteria

- [x] 01/02 thickened with named tests + proof commands
- [x] Schema evolution decision locked (additive v1, no mig 020)
- [x] P19 Loop tests listed as keepers
- [x] No product Go

## Next

Orchestrator: **P20-S06-01** after this row is `done`.
