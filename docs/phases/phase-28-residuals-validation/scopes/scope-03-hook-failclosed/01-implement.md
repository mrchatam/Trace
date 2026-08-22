# P28-S03-01 — Hook failClosed implementer

## Metadata
- id: P28-S03-01
- todo_ids: [P28-S03-01]
- role: implementer
- skills: [incremental-implementation, security-and-hardening, test-driven-development]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Close **R2 / R3 / R8** (INT-04/11 beyond install text): harden `CursorLoopGateHookScript` so **Option A** fails closed — when `.trace/config.json` has `enforce=strict` and `TRACE_TASK_ID` is absent, preToolUse returns **deny**. Preserve default-off / warn projects (empty task still **allow**). Automate INT-11 hook drift checks.

**No** Cursor Multitask rewrite. **No** daemon/HTTP.

## References

- [00-PLANNER.md](00-PLANNER.md) — locked defaults (this scope)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R2/R3/R8 seeds
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md) — product SoT (brief delta OK)
- [`internal/config/enforce.go`](../../../../../internal/config/enforce.go) — `LoadEnforceMode` fail-open semantics
- Live anchors: `enforcement.go` L106–108 allow path; `cursorhook.go` L124 `failClosed: false`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P28-S03-00 — do not re-debate)

| Item | Value |
|------|-------|
| **Option A (required)** | `enforce=strict` **and** empty `TRACE_TASK_ID` → echo deny JSON + non-zero exit (failClosed) |
| Default-off preserved | `enforce` missing / invalid / `off` / `warn` + empty task → **allow** (today’s L106–108 behavior) |
| Option B | Deferred — do **not** detect Parent orchestrator rule text |
| Config read in hook | `root="${TRACE_PROJECT_ROOT:-$PWD}"`; read `"$root/.trace/config.json"`; only exact `"enforce":"strict"` is strict (mirror `LoadEnforceMode` fail-open) |
| Task ID set | Unchanged: call `trace loop gate --task … --for edit`; missing `trace` binary still allow |
| Cursor `failClosed` field | Keep **`false`** in `.cursor/hooks.json` entry — policy failClosed is **script Option A**; Cursor soft-fail avoids locking out default-off on hook crash |
| INT-11 | Add `internal/install/hook_drift_test.go` (and/or extend `enforcement_test.go`) — golden/assert hooks.json shape + script stdout JSON has `permission` for allow **and** deny |
| Align text | Update `ParentOrchestratorRule` / `HookDriftNote` so docs say strict-gated deny (not “always deny if absent”) |
| Out of scope | Rewriting Cursor Multitask; daemon/HTTP; S04 honesty/attestation; changing HopBudget |

```text
preToolUse
  → empty TRACE_TASK_ID?
       yes → config enforce==strict? → deny : allow
       no  → existing gate path
```

## Preflight

Run from repo root before editing:

```bash
cd /home/ali/Desktop/Trace

test -f docs/phases/phase-28-residuals-validation/scopes/scope-03-hook-failclosed/00-PLANNER.md
test -f docs/phases/phase-28-residuals-validation/scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md

# Live residual anchors (must still show allow-without-task until you change it)
grep -n 'TRACE_TASK_ID' internal/install/enforcement.go | head
sed -n '100,122p' internal/install/enforcement.go
grep -n 'failClosed' internal/install/cursorhook.go
grep -n 'LoadEnforceMode\|EnforceStrict' internal/config/enforce.go

# Baseline
go test ./internal/install/... -count=1
```

Abort if RESIDUAL-AUDIT or `CursorLoopGateHookScript` missing.

## Files to touch

| File | Change |
|------|--------|
| `internal/install/enforcement.go` | Rewrite empty-task branch in `CursorLoopGateHookScript()` for Option A; align `ParentOrchestratorRule` text with strict-gated deny |
| `internal/install/cursorhook.go` | Update `HookDriftNote` (document script-level failClosed + keep schema `failClosed: false`); leave install entry `failClosed: false` |
| `internal/install/hook_drift_test.go` | **New** — INT-11: hooks.json entry keys (`command`, `matcher`, `failClosed`); deny/allow JSON shapes from script |
| `internal/install/enforcement_test.go` | Add/extend: strict+no-task → deny; off/warn/missing+no-task → allow; install still writes gate path |
| `docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md` | Optional brief harness-hook delta (strict + no task → deny) |

Do **not** change `cmd/trace/` product CLI beyond what install already generates. Re-install path is covered by existing cursor-hook Install tests writing the script body.

## Minimal todos

1. **Preflight** — run bash block; record pass in board Notes.
2. **Script Option A** — in `CursorLoopGateHookScript()`, when `task_id` empty: if config at `$root/.trace/config.json` is strict → deny (`permission`/`user_message`/`agent_message` consistent with existing deny shape); else allow.
3. **Docs constants** — `ParentOrchestratorRule` + `HookDriftNote` match Option A + Cursor `failClosed: false` rationale.
4. **Tests** — unit/simulated script:
   - temp root + `.trace/config.json` `{"enforce":"strict"}` + unset `TRACE_TASK_ID` → deny
   - same with `off` / `warn` / missing file → allow
   - `TRACE_TASK_ID` set still contains `loop gate` / `--for edit` in installed script
5. **INT-11** — `hook_drift_test.go`: after Install write, parse hooks.json preToolUse entry for required fields; assert allow/deny stdout is valid JSON with `permission` key.
6. **Regression** — `go test ./internal/install/... -count=1` (and `./internal/config/...` if touched indirectly — usually not).
7. **Board** — **P28-S03-01** status + notes only (behavior delta summary).

## Test commands

```bash
cd /home/ali/Desktop/Trace
go test ./internal/install/... -count=1
# Targeted if useful:
go test ./internal/install/ -count=1 -run 'Hook|CursorLoop|Drift|FailClosed|Strict'
```

Optional manual smoke (not required if unit covers):

```bash
# After writing script to a temp dir with enforce=strict and empty TRACE_TASK_ID
# printf '' | env -u TRACE_TASK_ID TRACE_PROJECT_ROOT="$tmp" bash "$tmp/.cursor/hooks/trace-loop-gate.sh"
# expect permission deny
```

## Exit criteria

- [ ] Strict + no `TRACE_TASK_ID` → deny (unit/simulated preToolUse evidence)
- [ ] Non-strict (`off`/`warn`/missing) + no task → allow
- [ ] Task ID set → existing gate behavior preserved (script still calls `loop gate --for edit`)
- [ ] Cursor hooks.json entry still has `failClosed: false` with HookDriftNote updated
- [ ] INT-11 drift test present and PASS
- [ ] `go test ./internal/install/... -count=1` PASS
- [ ] No daemon/HTTP; default-off projects not broken
- [ ] Board Notes cite Option A + test names

## Todo updates

Per board-rights: status + notes on **P28-S03-01** only. Do not edit `done` planner prompt substance.

## Next

`P28-S03-02` (fresh reviewer subagent)
