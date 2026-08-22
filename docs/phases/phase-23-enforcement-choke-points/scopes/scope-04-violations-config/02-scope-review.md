# P23-S04-02 — Review violations + config

## Metadata
- id: P23-S04-02
- todo_ids: [P23-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent review: additive status JSON, config defaults fail-closed to `off`, gate parity, and **no accidental auto-enforce** from config `strict`.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S04-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S04-01 deliverable: [01-violations-config.md](./01-violations-config.md)
- S01 library: `internal/loop/gate.go` (must be unchanged by S04)
- S03 enforce: `cmd/trace/transition.go`, `cmd/trace/seed.go` (must not gain config-driven auto-enforce)

## Session start
Follow agent-loop-protocol. Fresh reviewer context. Board edits: **status + notes only**.

## Keeper tests (must re-run — all green)

```bash
go test ./internal/config/...
go test ./internal/loop/... -run 'Status|Gate'
go test ./cmd/trace -run 'TestLoopStatus|TestTraceConfig|TestLoopGate|TestLoopNext|TestLoopApply|TestHelpIncludesLoop|TestTransitionDoneEnforce|TestSeedExportStrict|TestReviewCreateSetDone|TestAllowDoneWarnsOnStderr'
```

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Status additive | `grep Violations` in `apply.go` — field on `StatusResult`; `schema_version` constant unchanged |
| Gate parity | `Status()` calls `EvaluateGate(..., GateForEdit)` only — no SelectNext fork |
| Parity test | `TestLoopStatusViolationsMatchGateEdit` passes |
| Empty array | Clean status has `"violations": []` not null/absent |
| Config loader | `internal/config/enforce.go` — fail-closed to `off` |
| Missing file | `TestTraceConfigEnforceDefaultOff` + `TestLoadEnforceModeMissingFile` |
| Malformed JSON | `TestTraceConfigEnforceMalformedFailClosedOff` |
| Warn stderr | `TestTraceConfigEnforceWarnSurfacesStderr` — exit **0**, stderr non-empty |
| Strict = warn | `TestTraceConfigEnforceStrictSurfacesStderr` — same behavior as warn |
| Off silent | `TestTraceConfigEnforceOffNoStderrOnViolation` |
| No auto-enforce transition | Plain `transition … DONE` without `--enforce` unchanged with config `strict` |
| No auto-enforce export | Plain `seed export` unchanged with config `strict` |
| S01 untouched | `internal/loop/gate.go` diff empty for S04 |
| S03 untouched | transition/seed enforce paths do not import `internal/config` |
| Gate CLI unchanged | `cmdLoopGate` does not read config (exit 0/1/2 contract preserved) |
| Named tests | All tests from 01 prompt present and passing |
| Help | `trace help` documents config + status violations |
| Gitignore | `.trace/config.json` under gitignored `.trace/` |

## Review checklist

- [ ] **Blocker:** Duplicate SelectNext / policy logic outside `gate.go`
- [ ] **Blocker:** Config `strict` auto-enables transition `--enforce` or export `--strict --enforce`
- [ ] **Blocker:** `trace.loop.status.v1` schema string changed
- [ ] **Blocker:** `violations` omitted or `null` on clean status
- [ ] **Blocker:** Status violations use gate other than `GateForEdit` without S04-00 lock
- [ ] **Blocker:** Malformed config crashes or enables strict behavior
- [ ] **Blocker:** Missing named tests from 01 prompt
- [ ] **Blocker:** P19 loop status keeper tests regressed
- [ ] **High:** Status violations disagree with `loop gate --for edit` for same fixture
- [ ] **High:** Config `warn`/`strict` changes status exit code away from 0
- [ ] **High:** Changes to `internal/loop/gate.go` (policy belongs in S01)
- [ ] **High:** S03 transition/export behavior changed when config present
- [ ] **Medium:** Help text missing config default-off note
- [ ] **Medium:** `trace init` unexpectedly writes config (unless explicitly locked otherwise)
- [ ] **Low:** Inconsistent stderr prefix vs gate/transition
- [ ] **Nit:** SQL migration added for config (should be file-only)

## Default-off verification (walk through)

1. Fresh `trace init` — no `.trace/config.json`:
   - `loop status --task <id>` → exit **0**, violations populated when blocked, **no** stderr hints.
2. Write `{"enforce":"warn"}` — blocked task:
   - status exit **0**, violations in JSON, stderr hints present.
3. With config `strict`, verification-debt task:
   - `transition … DONE --as-operator` **without** `--enforce` → still succeeds (domain allows).
   - Same with `--enforce` → exit **1** (unchanged S03).

## S05 handoff verification

Harness install rules will reference:

```bash
trace loop status --task "$TRACE_TASK_ID"
# parse stdout JSON violations[] for agent context

# optional local config:
# .trace/config.json → { "enforce": "warn" }
```

Confirm:

- Status JSON is stable for hook-less agent reads.
- Config pointer in help/ENFORCEMENT aligns with install docs.
- `strict` documented as CI recommendation only — hooks still use explicit `--enforce` on DONE/export.

## Spawn policy

- **blocker/high:** inline fix if ≤10 lines and zero policy change; else spawn `P23-S04-02a` implement + `02b` review immediately below this row
- **medium:** prefer spawn unless trivial typo
- Do not rewrite S04-00/S04-01 `done` prompts

## Exit criteria

- [ ] No open blocker/high without pending forward row
- [ ] Confidence **medium** or **high** with command output in Notes
- [ ] Residual risks listed if medium (e.g. future env override, init --write-config)
- [ ] APPROVE or spawn documented on board

## Minimal todos

- [ ] Re-run keeper tests; paste pass summary in Notes
- [ ] Walk `Status()` gate call + parity test
- [ ] Walk config loader fail-closed paths
- [ ] Confirm transition/export ignore config
- [ ] Verify all named tests exist
- [ ] Confirm `internal/loop/gate.go` diff empty for S04
- [ ] Set row done with confidence + residuals
