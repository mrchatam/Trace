# P23-S01-02 — Review premature implementation gate library

## Metadata
- id: P23-S01-02
- todo_ids: [P23-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [user-codegraph]
- verification: automated

## Objective
Independent review: evaluator reuses `BuildPolicyInputs` + `SelectNext` without policy fork; violations stable for S02 `trace.loop.gate.v1` + S04 status `violations[]`; no CLI in S01; loop keepers green.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S01-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S01-01 deliverable: [01-premature-impl-gate-lib.md](./01-premature-impl-gate-lib.md)

## Session start
Follow agent-loop-protocol. Fresh reviewer context. Board edits: **status + notes only**.

## Keeper tests (must re-run — all green)

```bash
go test ./internal/deliberation/...
go test ./internal/loop/...
go test ./internal/domain/... -run 'Gate|Premature'
go test ./cmd/trace -run 'TestLoopNext|TestLoopApply|TestLoopStatus'
```

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Policy reuse | `grep SelectNext` in `gate.go` — must call deliberation package, not reimplement branches |
| No fork | `internal/deliberation/select.go` unchanged by S01 |
| API shape | `EvaluateGate(ctx, dom, plan, st, taskID, gateFor)` in `internal/loop/gate.go` |
| Violation schema | Struct fields: `code`, `for`, `message`, `recommended_phase`, `reason_code` with JSON tags |
| Domain error | `PrematureImplementation.Code() == "premature_implementation"` |
| GateFor table | Matches 01-premature-impl-gate-lib.md locked table (orient / edit / execute / done / export) |
| Edit rule | Edit allows only when SelectNext → EXECUTE && !stopped |
| Execute rule | Stricter: EXECUTE + execute_pending reason + inputs.ExecutePending |
| Done rule | Blocks verification debt, regression, blocking uncertainty, incomplete deliberation phases; does **not** duplicate review PASS / operator (S03) |
| Named tests | All 17 tests from 01 prompt present and passing |
| No CLI | `cmd/trace` untouched by S01 |
| No SQL | No new migrations |

## Review checklist

- [ ] **Blocker:** Duplicate SelectNext policy table anywhere outside `deliberation/select.go`
- [ ] **Blocker:** GateFor mapping diverges from ENFORCEMENT.md / 01 locked table
- [ ] **Blocker:** Missing named unit tests from 01 prompt
- [ ] **Blocker:** Loop keeper tests regressed
- [ ] **High:** Violation JSON not ready for S02 (`trace.loop.gate.v1` field parity)
- [ ] **High:** Edit gate allows non-EXECUTE phases (premature impl hole)
- [ ] **High:** Done gate silently skips verification debt or open regression
- [ ] **Medium:** `loadDeliberationState` duplicated incorrectly (should match loop semantics)
- [ ] **Medium:** p19 saturation ignored (always false)
- [ ] **Low:** Violation messages opaque or missing recommended_phase when blocked by SelectNext
- [ ] **Nit:** Exported symbols beyond `EvaluateGate`, types needed by S02/S04

## S02 handoff verification

Confirm S02 can wrap without changes:

```go
allowed, violations, err := loop.EvaluateGate(ctx, dom, plan, st, taskID, loop.GateForEdit)
// → JSON envelope with violations[], exit 0/1
```

## Spawn policy

- **blocker/high:** inline fix if ≤10 lines and zero policy change; else spawn `P23-S01-02a` implement + `02b` review immediately below this row
- **medium:** prefer spawn unless trivial typo
- Do not rewrite S01-00/S01-01 `done` prompts

## Exit criteria

- [ ] No open blocker/high without pending forward row
- [ ] Confidence **medium** or **high** with command output in Notes
- [ ] Residual risks listed if medium (e.g. export extensions deferred S03)
- [ ] APPROVE or spawn documented on board

## Minimal todos

- [ ] Re-run keeper tests; paste pass summary in Notes
- [ ] Walk GateFor table against `gate.go` line-by-line
- [ ] Verify all 17 named tests exist in `gate_test.go`
- [ ] Confirm `cmd/trace` diff empty for S01 scope
- [ ] Set row done with confidence + residuals
