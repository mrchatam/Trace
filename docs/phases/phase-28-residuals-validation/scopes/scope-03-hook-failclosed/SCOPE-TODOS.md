# S03 scope todos — hook failClosed

| ID | Status | Prompt | Notes |
|----|--------|--------|-------|
| P28-S03-00 | done | [00-PLANNER.md](00-PLANNER.md) | Option A locked 2026-08-20 |
| P28-S03-01 | pending | [01-implement.md](01-implement.md) | Implement Option A + INT-11 |
| P28-S03-02 | pending | [02-review.md](02-review.md) | Independent review |

## Locked targets (from P28-S03-00)

| Target | Lock |
|--------|------|
| Behavior | `enforce=strict` + empty `TRACE_TASK_ID` → **deny** |
| Default-off | `off` / `warn` / missing / invalid + empty task → **allow** |
| Option B | Deferred |
| Cursor `failClosed` JSON field | Keep **`false`**; policy failClosed = script |
| Files | `enforcement.go` (script + ParentOrchestratorRule), `cursorhook.go` (HookDriftNote), `hook_drift_test.go` (new), `enforcement_test.go` |
| Residuals | Close R2, R3, R8 |
| Out | Multitask rewrite; daemon/HTTP; S04 polish |

## Implementer acceptance (S03-01)

- [ ] Strict + no task deny (test)
- [ ] Non-strict + no task allow (test)
- [ ] INT-11 drift test PASS
- [ ] `go test ./internal/install/...` PASS
