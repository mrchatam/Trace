# P28-S03-02 — Hook failClosed independent review

## Metadata
- id: P28-S03-02
- todo_ids: [P28-S03-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review that **R2 / R3 / R8** are closed by product evidence (not Notes alone): Option A deny under strict+no-task, default-off preserved, INT-11 drift automation present. Fresh subagent — do **not** reuse P28-S03-01 session.

## References

- [01-implement.md](01-implement.md) — locked Option A
- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R2/R3/R8
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) — FM-05
- Live: `internal/install/enforcement.go`, `cursorhook.go`, `hook_drift_test.go`, `enforcement_test.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked review policy

| Item | Rule |
|------|------|
| Option A | Strict + empty `TRACE_TASK_ID` **must** deny; off/warn/missing **must** allow |
| Option B | Must **not** appear as implemented parent-orchestrator detection |
| Cursor `failClosed` | Entry remains `false` unless planner lock was explicitly superseded (it was not) |
| Default-off | Missing config + empty task → allow |
| Out of scope creep | Daemon/HTTP / Multitask rewrite → blocker |
| R2/R3/R8 close | Only with passing tests + live script body matching Option A |

## Verify checklist

### R2 — INT-04 beyond install text

- [ ] `CursorLoopGateHookScript` no longer always-allows on empty `TRACE_TASK_ID`
- [ ] Strict path denies with JSON `permission` (deny) and useful message
- [ ] `ParentOrchestratorRule` text matches strict-gated semantics (not contradicting script)

### R3 — FM-05 / failClosed

- [ ] Unit evidence: strict config + no task → deny
- [ ] Unit evidence: non-strict + no task → allow
- [ ] Cursor hooks.json `failClosed: false` retained; HookDriftNote documents script-level policy

### R8 — INT-11 drift

- [ ] `hook_drift_test.go` (or equivalent) asserts hooks.json shape (`command`, `matcher`, `failClosed`)
- [ ] Allow and deny stdout shapes validated as JSON with `permission`
- [ ] Existing `TestHookDriftNoteNonEmpty` still meaningful or superseded with stronger checks

### Process / blast radius

- [ ] `go test ./internal/install/... -count=1` PASS (reviewer re-runs)
- [ ] No daemon/HTTP; no Multitask product rewrite
- [ ] S04 honesty/attestation not silently implemented here
- [ ] Default-off projects preserved (explicit test or clear script branch)

### Live spot-checks (reviewer runs)

```bash
REPO=/home/ali/Desktop/Trace
cd "$REPO"

# Option A present in script source
sed -n '100,140p' internal/install/enforcement.go
grep -n 'strict\|TRACE_TASK_ID\|permission' internal/install/enforcement.go | head -40

# Cursor failClosed still false at install site
grep -n 'failClosed' internal/install/cursorhook.go

# INT-11 test file
test -f internal/install/hook_drift_test.go || rg -n 'HookDrift|failClosed|permission' internal/install/*_test.go

# Re-verify
go test ./internal/install/... -count=1
```

Confirm empty-task branch is **not** still:

```bash
# Anti-pattern (pre-S03) — must NOT be the only empty-task path:
# if [[ -z "$task_id" ]]; then echo '{"permission":"allow"}'; exit 0; fi
```

## Findings severity

| Level | Action |
|-------|--------|
| blocker | Still allow-on-empty under strict; default-off broken (deny without strict); no INT-11 test; daemon/HTTP added |
| high | Docs/`ParentOrchestratorRule` contradict script; Cursor `failClosed` flipped without note; tests only grep strings without simulating deny |
| medium | HookDriftNote stale; ENFORCEMENT.md omitted when behavior changed |
| low / nit | Wording; test name polish |

## Spawn policy

Insert **immediately below** this row (`P28-S03-02a` implement / `02b` review) if blocker/high remains open:

- Fix Option A script + tests, **or**
- Restore default-off allow path if over-denied

Do not rewrite `done` P28-S03-01 prompt; spawn forward.

## Verdict

| Outcome | Next |
|---------|------|
| Option A + INT-11 PASS, default-off OK | APPROVE; R2/R3/R8 closable in Notes; proceed **P28-S04-00** |
| Gap remains | spawn 02a/02b; do not start S04 until resolved or orchestrator decides |

Confidence must be **medium or high** with evidence. Empty APPROVE without re-run tests is invalid.

## Todo updates

Status + notes on **P28-S03-02**; spawn rows allowed.

## Exit criteria

- [ ] Checklist complete with live evidence
- [ ] No open blocker/high without pending spawn
- [ ] Confidence medium or high
- [ ] Next runnable stated (P28-S04-00 or spawn)
