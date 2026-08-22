# P23-S04-01 — Implement violations on status + config

## Metadata
- id: P23-S04-01
- todo_ids: [P23-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- verification: automated

## Objective
Add `violations[]` to loop status and load `.trace/config.json` enforce modes. Thin wiring over S01 `EvaluateGate` + new `internal/config` loader — **no gate policy changes**.

## References
- S04-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S01: `internal/loop/gate.go` (read-only)
- Live: `internal/loop/apply.go`, `cmd/trace/loop.go`

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Status schema | Additive `violations[]` only — `trace.loop.status.v1` string unchanged |
| Status gate | `EvaluateGate(..., GateForEdit)` inside `Status()` |
| Parity | `status.violations` == `loop gate --for edit`.violations for same state |
| Config path | `.trace/config.json` |
| Config shape | `{ "enforce": "off" \| "warn" \| "strict" }` |
| Default | Missing/malformed/invalid → `off` |
| Loader | `internal/config.LoadEnforceMode(root)` |
| `off` | JSON violations only; no config stderr hints |
| `warn` / `strict` | Same runtime behavior: stderr hints when violations non-empty; exit **0** |
| Auto-enforce | **None** — config does not enable S03 `--enforce` flags |
| Init | Do not require `trace init` to write config |
| Gate CLI | Unchanged exit contract; config not consulted in S04 for gate |

## Files to touch

| File | Action |
|------|--------|
| `internal/config/enforce.go` | Create `EnforceMode`, `LoadEnforceMode` |
| `internal/config/enforce_test.go` | Loader unit tests |
| `internal/loop/apply.go` | Add `Violations []Violation` to `StatusResult`; populate in `Status()` |
| `cmd/trace/loop.go` | `cmdLoopStatus`: load config; stderr hints for warn/strict |
| `cmd/trace/loop_test.go` | Named CLI tests below |
| `cmd/trace/help.go` | Config + status violations docs |
| `internal/loop/gate.go` | **Do not modify** |

## Implementation sketch

### `internal/config/enforce.go`

```go
package config

type EnforceMode string

const (
    EnforceOff    EnforceMode = "off"
    EnforceWarn   EnforceMode = "warn"
    EnforceStrict EnforceMode = "strict"
)

func LoadEnforceMode(projectRoot string) EnforceMode {
    // read <root>/.trace/config.json
    // json.Unmarshal into struct { Enforce string `json:"enforce"` }
    // validate; any failure → EnforceOff
}
```

### `internal/loop/apply.go` — `Status()`

After building deliberation (both insufficient_history and saturated paths):

```go
dom := domain.New(st)
allowed, violations, err := EvaluateGate(ctx, dom, plan, st, seed.TaskID, GateForEdit)
if err != nil {
    return StatusResult{}, fmt.Errorf("loop status: gate: %w", err)
}
if allowed || violations == nil {
    violations = []Violation{}
}
// attach to StatusResult{ ..., Violations: violations }
```

Add to `StatusResult`:

```go
Violations []Violation `json:"violations"`
```

Ensure `nil` slice marshals as `[]` (initialize empty slice when allowed).

### `cmd/trace/loop.go` — `cmdLoopStatus`

```go
res, err := loop.Status(...)
mode := config.LoadEnforceMode(abs)
if (mode == config.EnforceWarn || mode == config.EnforceStrict) && len(res.Violations) > 0 {
    for _, v := range res.Violations {
        fmt.Fprintf(os.Stderr, "loop status: %s\n", v.Message)
    }
}
// encode res to stdout; return exitOK (0) on success
```

## Named tests (minimum — S04-01 must implement all)

### Status violations

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestLoopStatusIncludesViolationsWhenBlocked` | init + blocking uncertainty on task | `violations` len ≥ 1; `violations[0].code=premature_implementation`; `for=edit` |
| `TestLoopStatusViolationsEmptyWhenClean` | init + full cycle clear / execute-ready | `violations` is empty array |
| `TestLoopStatusViolationsMatchGateEdit` | blocked fixture | parse status JSON + run `loop gate --for edit`; violation arrays deeply equal |
| `TestLoopStatusViolationsAlwaysArray` | any status output | `violations` key present; type array (not null) |
| `TestLoopStatusSchemaVersionUnchanged` | any status output | `schema_version=trace.loop.status.v1` |
| `TestLoopStatusDeliberationFields` | existing keeper | still passes (additive JSON only) |
| `TestLoopStatusBlockedWhenBlockingUncertainty` | existing keeper | `deliberation.blocked=true` unchanged |

### Config enforce modes

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestTraceConfigEnforceDefaultOff` | no config file; blocked task | stdout has violations; stderr empty |
| `TestTraceConfigEnforceMalformedFailClosedOff` | `{not json` in config; blocked | behaves as off (stderr empty) |
| `TestTraceConfigEnforceInvalidValueFailClosedOff` | `{"enforce":"loud"}` | behaves as off |
| `TestTraceConfigEnforceWarnSurfacesStderr` | `{"enforce":"warn"}` + blocked | exit **0**; stderr contains violation message substring |
| `TestTraceConfigEnforceStrictSurfacesStderr` | `{"enforce":"strict"}` + blocked | same as warn (exit **0**, stderr hint) |
| `TestTraceConfigEnforceOffNoStderrOnViolation` | `{"enforce":"off"}` + blocked | violations in stdout; stderr empty |
| `TestHelpIncludesTraceConfig` | `trace help` | mentions `.trace/config.json` and enforce modes |

### Config loader unit tests (`internal/config/enforce_test.go`)

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestLoadEnforceModeMissingFile` | empty temp dir | `EnforceOff` |
| `TestLoadEnforceModeValidValues` | write off/warn/strict files | each parses correctly |
| `TestLoadEnforceModeMalformedJSON` | invalid file | `EnforceOff` |
| `TestLoadEnforceModeUnknownValue` | `"enforce":"yolo"` | `EnforceOff` |

Helper for config tests:

```go
func writeTraceConfig(t *testing.T, root, content string) {
    t.Helper()
    dir := filepath.Join(root, ".trace")
    if err := os.MkdirAll(dir, 0o755); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(content), 0o644); err != nil {
        t.Fatal(err)
    }
}
```

## Keeper tests (must stay green)

```bash
go test ./internal/config/...
go test ./internal/loop/... -run 'Status|Gate'
go test ./cmd/trace -run 'TestLoopStatus|TestLoopGate|TestLoopNext|TestLoopApply|TestHelpIncludesLoop|TestTransitionDoneEnforce|TestSeedExportStrict'
```

## Help text (FINAL)

See S04-00 locks — add loop status + `.trace/config.json` blocks to `help.go`.

## Implementation notes

- Status internal gate error → `loop status: gate: …` → cmd exit **2** (`exitFail`), same as other status failures.
- Do not add `--enforce` to `loop status` — enforcement remains opt-in on transition/export (S03) and gate CLI block exit (S02).
- P19/P20 status keeper tests must pass without asserting absence of `violations` key.
- `.trace/` stays gitignored — config is local-only.
- Optional doc touch: one ENFORCEMENT.md cross-link already sufficient; no CONTRIBUTING change required unless already editing.

## Exit criteria

- [ ] All named tests green (14 CLI + 4 unit minimum)
- [ ] Status violations match gate for same state
- [ ] Config load fail-closed to `off` on malformed JSON
- [ ] Default off; no auto-enforce on transition/export
- [ ] Loop + gate keeper tests green
- [ ] No changes to `internal/loop/gate.go`

## Minimal todos

- [ ] Create `internal/config` loader + tests
- [ ] Add `Violations` to `StatusResult` + `Status()` gate call
- [ ] Wire config stderr hints in `cmdLoopStatus`
- [ ] Implement all named tests
- [ ] Update help strings
- [ ] Run keeper commands; board row status + notes only

## Next

**P23-S04-02** review after implementation.
