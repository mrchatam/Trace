# P23-S04-00 — Violations + config planner

## Metadata
- id: P23-S04-00
- todo_ids: [P23-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock `violations[]` on `trace.loop.status.v1` and `.trace/config.json` enforce modes (`off`|`warn`|`strict`, default `off`). **No product Go this row.**

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S01 evaluator — **done** (`internal/loop/gate.go`): `Violation`, `EvaluateGate(..., GateForEdit)`
- S02 gate CLI — **done**: `trace.loop.gate.v1` uses same `Violation` shape
- S03 enforce — **done** (S03-01): transition/export `--enforce` opt-in; **does not** read config (S04 scope)
- Live: `internal/loop/apply.go` (`Status`, `StatusResult`), `cmd/trace/loop.go` (`cmdLoopStatus`)

### S03 handoff (S03-02 review pending — not blocking S04 planner)

- Transition `--enforce` and export `--strict --enforce` are flag-opt-in only.
- S04 **must not** wire config `strict` to auto-enable those flags.
- S03-02 evidence checklist includes "No config read" — preserved in S04.

## Live inventory (P23-00 + S02 verified)

**Loop status today** (`internal/loop/apply.go`, `cmd/trace/loop.go`):

- `StatusResult` emits `schema_version: trace.loop.status.v1` with deliberation, promotion_blocked, verification_cycle.
- **No `violations[]`** — S04 adds additive field only; schema string unchanged.
- `deliberation.blocked` already reflects SelectNext policy; `violations[]` mirrors gate evaluator output (machine-checkable, harness-parseable).

**Config today**:

- **No** `.trace/config.json` loader — S04 introduces file-based config.
- **No** `internal/config` package — S04 creates it.
- `trace init` creates `.trace/trace.db` only — does **not** write config (S04-01 optional doc only).

## Live touch points (S04-01)

| File | Change |
|------|--------|
| `internal/config/enforce.go` | **New** — `LoadEnforceMode(root)` → `off`/`warn`/`strict`; fail-closed to `off` |
| `internal/config/enforce_test.go` | **New** — unit tests for loader |
| `internal/loop/apply.go` | `StatusResult.Violations []Violation`; `Status()` calls `EvaluateGate(..., GateForEdit)` |
| `cmd/trace/loop.go` | `cmdLoopStatus`: load config; stderr hints when `warn`/`strict` + non-empty violations |
| `cmd/trace/loop_test.go` | Named status + config CLI tests |
| `cmd/trace/help.go` | Document `.trace/config.json` shape + enforce modes |
| `internal/loop/gate.go` | **No changes** — policy stays in S01 |

## Locked defaults (S04-01 must not re-debate)

| Item | Value |
|------|-------|
| Status schema | Additive `violations[]` on existing `trace.loop.status.v1` — **schema string unchanged** |
| Violation shape | Same objects as S01 `loop.Violation` / gate JSON `violations[]` entries |
| Status gate | `EvaluateGate(..., GateForEdit)` — matches default harness edit choke point |
| Parity rule | For same task/store state, `status.violations` **equals** `loop gate --for edit`.violations |
| Empty violations | Always `"violations": []` when clean — never `null`, never omitted |
| Config path | `<projectRoot>/.trace/config.json` |
| Config shape | `{ "enforce": "off" \| "warn" \| "strict" }` — top-level object only (MVP) |
| Default | Missing file, unreadable file, malformed JSON, unknown `enforce` value → **`off`** (fail-closed) |
| Load order | **File only** — no env override in S04 (defer `TRACE_ENFORCE` to post-MVP) |
| Loader home | **`internal/config`** — reusable; loop/status/gate cmd import it |
| SQL | **None** — file-based config OK for MVP |
| `off` | Status JSON includes `violations[]`; **no** config-driven stderr hints on status; CLI exit **0** |
| `warn` | Status JSON includes `violations[]`; when non-empty, **one stderr line per violation** (prefix `loop status:`); exit **0** |
| `strict` | **Identical behavior to `warn` in S04** — difference is documentation / S05 install recommendation only |
| Auto-enforce | **`strict` does NOT** auto-enable `--enforce` on transition or export — **explicit flags only** (S03 locks preserved) |
| Init | S04-01 does **not** require `trace init` to create config; help documents manual creation |
| Gate CLI | Config modes do **not** change gate exit codes (0/1/2); blocked gate already stderr-hints |
| Transition/export | S04 does **not** read config for enforce behavior — S03 flag semantics unchanged |

### Config load contract (FINAL)

```go
// internal/config/enforce.go
type EnforceMode string

const (
    EnforceOff    EnforceMode = "off"
    EnforceWarn   EnforceMode = "warn"
    EnforceStrict EnforceMode = "strict"
)

// LoadEnforceMode reads <root>/.trace/config.json.
// Returns EnforceOff when file missing or any parse/validation error.
func LoadEnforceMode(projectRoot string) EnforceMode
```

Validation rules:

1. File missing → `off`
2. JSON parse error → `off`
3. `enforce` key missing → `off`
4. `enforce` not one of `off|warn|strict` (case-sensitive) → `off`
5. Extra JSON keys ignored (forward-compatible)

### `trace.loop.status.v1` extension (FINAL)

#### New field on `StatusResult`

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `violations` | array | yes | Always present; `[]` when gate allows edit |

#### Violation element (same as S01 / gate)

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `code` | string | yes | e.g. `premature_implementation`, `gate_orient_failed` |
| `for` | string | yes | Always `"edit"` for status (GateForEdit) |
| `message` | string | yes | Human-readable |
| `recommended_phase` | string | when set | From SelectNext |
| `reason_code` | string | yes | e.g. `blocking_uncertainty` |

**No top-level lift** on status (unlike gate envelope) — consumers read `violations[0]` directly or use existing `deliberation.recommended_phase`.

#### Example — blocked (edit)

```json
{
  "schema_version": "trace.loop.status.v1",
  "seed": { "task_id": "550e8400-e29b-41d4-a716-446655440000", "goal_id": "…" },
  "reason": "insufficient_history",
  "saturated": false,
  "deliberation": {
    "phase": "EXECUTE",
    "recommended_phase": "INVESTIGATE",
    "why_selected": "blocking_uncertainty",
    "hop_count": 0,
    "stopped": false,
    "blocked": true,
    "policy_inputs": { }
  },
  "violations": [
    {
      "code": "premature_implementation",
      "for": "edit",
      "message": "edit blocked: recommended phase INVESTIGATE (blocking_uncertainty)",
      "recommended_phase": "INVESTIGATE",
      "reason_code": "blocking_uncertainty"
    }
  ]
}
```

#### Example — clean (edit allowed)

```json
{
  "schema_version": "trace.loop.status.v1",
  "seed": { "task_id": "550e8400-e29b-41d4-a716-446655440000", "goal_id": "…" },
  "reason": "tasks_and_plan_unchanged",
  "saturated": true,
  "deliberation": {
    "phase": "EXECUTE",
    "recommended_phase": "EXECUTE",
    "blocked": false
  },
  "violations": []
}
```

### Config file example

```json
{
  "enforce": "warn"
}
```

### Help text (FINAL — S04-01 add to `help.go`)

**Loop status** — extend existing line:

```
  loop status --task <id> [--goal <id>]
                        Report trace.loop.status.v1 from persisted loop-step evidence.
                        Includes violations[] (edit gate parity). Optional
                        .trace/config.json enforce mode (off|warn|strict, default off):
                        warn/strict print stderr hints when violations present; exit stays 0.
```

**Config** — new short block under loop or general notes:

```
  .trace/config.json    Local enforce mode: { "enforce": "off"|"warn"|"strict" }.
                        Default off when missing. Does not auto-enable --enforce on
                        transition or seed export (use explicit flags).
```

## Planner work

1. [x] Lock config load order (file only; no env override).
2. [x] Lock loader home (`internal/config`).
3. [x] Lock status gate (`GateForEdit`) + parity with gate CLI.
4. [x] Lock enforce mode semantics (off/warn/strict; no auto-enforce).
5. [x] Thicken `01-violations-config.md` with JSON examples + named tests.
6. [x] Thicken `02-scope-review.md` with evidence table + keeper cmds.
7. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Config + status extensions locked
- [x] 01/02/SCOPE-TODOS runnable alone
- [x] No product Go

## Next

**P23-S04-01** after this row is `done`.
