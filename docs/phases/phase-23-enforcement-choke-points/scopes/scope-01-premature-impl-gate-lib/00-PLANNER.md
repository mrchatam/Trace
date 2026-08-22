# P23-S01-00 — Premature implementation gate library planner

## Metadata
- id: P23-S01-00
- todo_ids: [P23-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock `domain.PrematureImplementation` + shared gate evaluator reusing `PolicyInputs` / `SelectNext`. **No product Go this row.** Library-only — S02 owns CLI.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/deliberation/select.go`, `internal/loop/{policy,apply,next}.go`, `internal/domain/{regressions,service}.go`

## Live touch points (P23-00 inventory)

| File | Reuse for gate |
|------|----------------|
| `internal/loop/policy.go` | `BuildPolicyInputs` — blocking uncertainty, verify debt, open regression, plan exists |
| `internal/deliberation/select.go` | `SelectNext(dState, inputs)` — phase recommendation |
| `internal/loop/apply.go` | Post-apply SelectNext path — mirror semantics, do not duplicate |
| `internal/domain/regressions.go` | `HasOpenRegression`, verification gate helpers |
| `cmd/trace/transition.go` | DONE checks to align `GateForDone` (no `--enforce` yet) |

## Locked defaults (S01-01 must not re-debate)

| Item | Value |
|------|-------|
| Package | **`internal/loop`** (recommended: new `gate.go`) — single evaluator entrypoint; deliberation stays policy table only |
| Reuse | `BuildPolicyInputs`, `deliberation.SelectNext`, existing gate queries (blocking uncertainty, debt, regression) |
| Domain error | `domain.PrematureImplementation` with stable `Code()` e.g. `premature_implementation` |
| Violation record | Struct reusable by gate JSON + status `violations[]` (S04) |
| `--for` mapping | Table locked in ENFORCEMENT.md; evaluator accepts `GateFor` enum |
| SQL | Prefer **no migration** — pure policy over existing tables; if config needed defer to S04 |
| CLI | **None in S01** |

### GateFor → policy checks (FINAL intent)

| GateFor | Block when |
|---------|------------|
| `orient` | Task/seed not found or not orientable |
| `edit` | Would not allow EXECUTE: blocking uncertainty, open regression, missing plan, or phase ∉ {EXECUTE, …} per S01-00 table |
| `execute` | SelectNext ≠ EXECUTE with execute_pending semantics |
| `done` | DONE promotion gates fail (verification debt, open regression, etc.) — align with existing transition checks |
| `export` | Subset of done + export honesty (S03 may extend) |

### Violation struct (FINAL fields for JSON parity)

| Field | Type | Notes |
|-------|------|-------|
| `code` | string | e.g. `premature_implementation` |
| `for` | string | GateFor value |
| `message` | string | Human-readable |
| `recommended_phase` | string | From SelectNext when blocked |
| `reason_code` | string | e.g. `blocking_uncertainty`, `verification_debt`, `open_regression` |

### Evaluator API (FINAL shape — names may vary)

```go
type GateFor string // orient | edit | execute | done | export

// EvaluateGate(ctx, dom, plan, taskID string, for GateFor) (allowed bool, violations []Violation, err error)
```

- Pure policy path testable with stubbed `PolicyInputs` in unit tests.
- Must not fork SelectNext priority table.
- `GateForEdit`: block when SelectNext would not return EXECUTE.
- `GateForDone`: align with `domain.TransitionTask` DONE prerequisites + verification debt (S03 wires `--enforce`).

## Planner work

1. [ ] Confirm package placement (loop vs deliberation vs domain wrapper).
2. [ ] Lock violation struct fields for JSON parity with `trace.loop.gate.v1`.
3. [ ] Thicken `01-premature-impl-gate-lib.md` + `02-scope-review.md`.
4. [ ] Update `SCOPE-TODOS.md`.

## Exit criteria

- [ ] 01/02/SCOPE-TODOS runnable alone
- [ ] GateFor table + violation schema locked
- [ ] Named unit test rows listed for S01-01
- [ ] No product Go

## Next

**P23-S01-01** after this row is `done`.
