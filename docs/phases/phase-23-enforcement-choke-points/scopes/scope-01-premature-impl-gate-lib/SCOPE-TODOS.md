# S01 — Premature implementation gate library — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P23-S01-00 | scope planner | **done** — GateFor table + violation schema + evaluator API locked in 01/02 |
| 2 | P23-S01-01 | implementer | **next** — `internal/loop/gate.go` + `domain.PrematureImplementation` + 17 named tests |
| 3 | P23-S01-02 | reviewer | pending — policy reuse + keeper tests + S02 handoff |

**Depends on:** P23-00 done. **Blocks:** S02 CLI, S03 enforce, S04 status violations.

## Locked artifacts (S01-00)

- Package: `internal/loop/gate.go` (+ `gate_test.go`)
- Domain: `internal/domain/gate_errors.go` — `PrematureImplementation`, `Code() = premature_implementation`
- API: `EvaluateGate(ctx, dom, plan, st, taskID, gateFor)`
- Violation: `{code, for, message, recommended_phase, reason_code}`
- GateFor: orient | edit | execute | done | export (table in `01-premature-impl-gate-lib.md`)
- No SQL, no CLI
