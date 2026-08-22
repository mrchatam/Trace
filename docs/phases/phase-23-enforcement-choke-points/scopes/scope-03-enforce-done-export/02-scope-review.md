# P23-S03-02 — Review enforce DONE + strict export

## Metadata
- id: P23-S03-02
- todo_ids: [P23-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent review: opt-in **`--enforce`** only; escape hatches intact when flag absent; export enforce never writes on failure; thin adapter over S01 — no duplicated DONE/deliberation policy in `cmd/trace`.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S03-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S03-01 deliverable: [01-enforce-done-export.md](./01-enforce-done-export.md)
- S01 library: `internal/loop/gate.go` (must be unchanged by S03)
- Domain DONE policy: `internal/domain/task_state.go` (must be unchanged by S03)

## Session start
Follow agent-loop-protocol. Fresh reviewer context. Board edits: **status + notes only**.

## Keeper tests (must re-run — all green)

```bash
go test ./cmd/trace -run 'TestTransitionDoneEnforce|TestSeedExportStrict|TestReviewCreateSetDone|TestAllowDoneWarnsOnStderr|TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit|TestLoopGate|TestLoopNext|TestLoopApply|TestLoopStatus|TestHelpIncludesLoop|TestHelpIncludesTransitionEnforce|TestHelpIncludesSeedExportStrict'
go test ./internal/loop/... -run Gate
```

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Thin adapter — transition | `grep EvaluateGate` in `cmd/trace/transition.go` — gate call only when `--enforce` && `--to DONE` |
| Thin adapter — export | `grep EvaluateGate` in `cmd/trace/seed.go` — only in strict validation path |
| No policy fork | No `SelectNext` / `BuildPolicyInputs` in transition or seed export paths |
| Default unchanged | Plain `transition … DONE` + plain `seed export` behave as pre-S03 (review/caps/allow-done tests still pass) |
| Enforce scope | `--enforce` on `--to IN_PROGRESS` (etc.) is no-op, not usage error |
| Transition block exit | `--enforce` + gate block → exit **1**, not `exitFail` (2) |
| Transition allows | Clean cycle + review PASS + `--as-operator` + `--enforce` → exit **0** |
| Escape hatch without enforce | `--allow-done` without `--enforce` still works + WARNING |
| Enforce overrides allow-done on policy | Debt + `--allow-done --enforce` → exit **1** |
| Domain review gate preserved | PASS without `--as-operator` still fails without `--enforce` |
| Export enforce no write | `-o path --strict --enforce` with violations: file absent or byte-identical pre-run |
| Export strict warn-only | `--strict` without `--enforce` → exit **0**, stderr has violation lines |
| Export enforce requires strict | `--enforce` alone → usage exit **1** |
| Task filter | `--task` narrows export gate scan |
| Terminal filter | DONE/SKIPPED/STALE tasks skipped in full scan |
| Doc version check | Export strict rejects `version != 1` if ever possible |
| S01 untouched | `internal/loop/gate.go` diff empty for S03 |
| Domain untouched | `task_state.go` DONE policy diff empty for S03 |
| Named tests | All 17 tests from 01 prompt present and passing |
| Help | `trace help` documents `--enforce` / `--strict` |
| No config read | S03 does not load `.trace/config.json` (S04 scope) |

## Review checklist

- [ ] **Blocker:** Duplicate SelectNext / policy logic in `cmd/trace/transition.go` or `seed.go`
- [ ] **Blocker:** `--enforce` changes default behavior when flag absent
- [ ] **Blocker:** Export `--enforce` writes partial/corrupt file on violation
- [ ] **Blocker:** Transition enforce bypasses domain review PASS / FAIL / caps checks
- [ ] **Blocker:** Missing named CLI tests from 01 prompt
- [ ] **Blocker:** P17 seed export keepers regressed (`TestSeedExportRoundTrip`, omit surfaces, exported_at_commit)
- [ ] **High:** Enforce block returns exit 2 instead of 1 (breaks harness parity with loop gate)
- [ ] **High:** `--allow-done` bypasses gate when `--enforce` set
- [ ] **High:** Gate checks review PASS (policy belongs in domain only)
- [ ] **High:** Changes to `internal/loop/gate.go` (policy belongs in S01)
- [ ] **High:** Changes to `domain.TransitionTask` DONE rules
- [ ] **Medium:** `--enforce` without `--strict` on export silently ignored (must be usage error)
- [ ] **Medium:** Help text missing two-layer DONE explanation
- [ ] **Medium:** S03 reads config enforce mode and auto-blocks
- [ ] **Low:** Inconsistent stderr prefix between transition and export
- [ ] **Nit:** Redefines `exitGateBlocked` conflicting with `loop.go`

## Two-layer DONE verification (walk through)

1. Task with verification debt, review PASS, `--as-operator`:
   - Without `--enforce` → DONE succeeds (domain only).
   - With `--enforce` → blocked exit 1 (gate only).
2. Task clean, no review:
   - `--allow-done` without `--enforce` → succeeds + WARNING.
   - `--allow-done --enforce` → if gate clean, succeeds; if debt, gate blocks first.

## S05 handoff verification

Harness install rules will reference:

```bash
trace transition --task "$TRACE_TASK_ID" --to DONE --reason "…" --as-operator --enforce
trace seed export -o trace/graph.json --strict --enforce
```

Confirm exit **1** on policy block is hook-friendly (distinct from usage **1** / internal **2** — document that transition usage remains **1**, internal **2**).

## Spawn policy

- **blocker/high:** inline fix if ≤10 lines and zero policy change; else spawn `P23-S03-02a` implement + `02b` review immediately below this row
- **medium:** prefer spawn unless trivial typo
- Do not rewrite S03-00/S03-01 `done` prompts

## Exit criteria

- [ ] No open blocker/high without pending forward row
- [ ] Confidence **medium** or **high** with command output in Notes
- [ ] Residual risks listed if medium (e.g. multi-task export scan perf, export-honesty extensions beyond GateForExport)
- [ ] APPROVE or spawn documented on board

## Minimal todos

- [ ] Re-run keeper tests; paste pass summary in Notes
- [ ] Walk transition enforce path line-by-line against 01 sketch
- [ ] Walk export strict/enforce path; confirm no write on enforce failure
- [ ] Verify all 17 named tests exist
- [ ] Confirm `internal/loop/gate.go` + `task_state.go` diffs empty for S03
- [ ] Set row done with confidence + residuals
