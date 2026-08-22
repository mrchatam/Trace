# P23-S02-01 — Implement loop gate CLI

## Metadata
- id: P23-S02-01
- todo_ids: [P23-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- verification: automated

## Objective
Implement `trace loop gate` subcommand + `trace.loop.gate.v1` per **S02-00 FINAL locks**. Thin adapter over S01 `EvaluateGate` — **no policy fork**.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S02-00 planner: [00-PLANNER.md](./00-PLANNER.md)
- S01 library: `internal/loop/gate.go` — `EvaluateGate`, `GateFor`, `Violation`
- Live reuse: `cmd/trace/loop.go` (`cmdLoopNext` store/domain/planner open path)

## Session start
Follow agent-loop-protocol. Board edits: **status + notes only**.

## Locked defaults (from S02-00 — do not re-debate)

| Item | Value |
|------|-------|
| Subcommand | `gate` under existing `trace loop` |
| Flags | `--task` required; `--for` optional default `edit` |
| Valid `--for` | `orient`, `edit`, `execute`, `done`, `export` |
| Exit | **0** allowed, **1** blocked, **2** usage/internal (gate-specific — not global `exitUsage`) |
| Stdout JSON | Only on exit 0/1 (policy paths) |
| Evaluator | S01 `EvaluateGate` only |
| MCP | **No new tool** |

### Exit mapping (implement exactly)

```go
// Policy result from EvaluateGate:
if err != nil {
    // stderr: err.Error(); return exitFail (2); no stdout JSON
}
if allowed {
    // emit envelope; return exitOK (0)
}
// blocked: emit envelope; stderr = violations[0].Message; return 1  ← NOT exitUsage
```

Use literal `return 1` for blocked (or named `exitGateBlocked = 1` local const in `loop.go`).

### Envelope struct (FINAL shape)

```go
type gateEnvelope struct {
    SchemaVersion    string           `json:"schema_version"`
    TaskID           string           `json:"task_id"`
    For              string           `json:"for"`
    Allowed          bool             `json:"allowed"`
    RecommendedPhase string           `json:"recommended_phase,omitempty"`
    ReasonCode       string           `json:"reason_code,omitempty"`
    Violations       []loop.Violation `json:"violations"`
}
```

Build rules:

- `schema_version` = `"trace.loop.gate.v1"`
- `violations` = `nil` slice → marshal as `[]` (never `null`)
- When `allowed=false` and `len(violations)==1`: copy `recommended_phase` + `reason_code` to top level from `violations[0]`
- When `allowed=true`: omit top-level `recommended_phase` / `reason_code`; `violations` = empty array

### Help text (FINAL — add to `printLoopHelp`)

```
  gate --task <id> [--for orient|edit|execute|done|export]
        Check deliberation gate for a task. Emits trace.loop.gate.v1 JSON on stdout.
        Exit 0 when allowed, 1 when blocked, 2 on usage or internal error.
        Default --for is edit (pre-edit harness choke point).
```

Also add `gate` case to `cmdLoop` switch and usage on unknown subcommand path.

## Files to create/modify

| File | Action |
|------|--------|
| `cmd/trace/loop.go` | Add `case "gate"`, `cmdLoopGate`, extend `printLoopHelp` |
| `cmd/trace/loop_test.go` | Add named gate CLI tests below |

**Do not modify:** `internal/loop/gate.go`, `internal/deliberation/select.go`, MCP tools.

## Implementation sketch

```go
func cmdLoopGate(root string, args []string) int {
    fs := flag.NewFlagSet("loop gate", flag.ContinueOnError)
    fs.SetOutput(os.Stderr)
    taskID := fs.String("task", "", "seed task UUID")
    gateFor := fs.String("for", "edit", "gate context: orient|edit|execute|done|export")
    // parse flags → validate task + for enum → resolveRoot → store.Open
    // domain.New(st), planner.New(st)
    // allowed, violations, err := loop.EvaluateGate(ctx, dom, plan, st, *taskID, loop.GateFor(*gateFor))
    // build gateEnvelope, json.Marshal to stdout
    // map exit per locked table
}
```

Mirror `cmdLoopNext` for: `resolveRoot`, `store.Open`, `failCLIDenied`, defer close.

Validate `--for` before calling `EvaluateGate` (invalid value → stderr + exit 2, do not rely on evaluator error alone).

## Named CLI tests (minimum — S02-01 must implement all)

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestLoopGateAllowedExitZero` | init + goal/task/plan critiqued (reuse loop test helpers) | `loop gate --task <id> --for edit` → exit **0**, `allowed=true`, `violations=[]` |
| `TestLoopGateBlockedExitOne` | init + blocking uncertainty linked to task | exit **1** (not `exitUsage`), `allowed=false`, `violations[0].code=premature_implementation` |
| `TestLoopGateJSONSchemaVersion` | any gate invocation (allowed or blocked) | stdout parses; `schema_version == "trace.loop.gate.v1"` |
| `TestLoopGateTopLevelLiftFromViolation` | blocked edit (e.g. plan missing) | top-level `recommended_phase` + `reason_code` match `violations[0]` |
| `TestLoopGateAllowedEmptyViolations` | execute-ready seed | `allowed=true`; `violations` is empty array (not null) |
| `TestLoopGateBlockedStderrHint` | blocked scenario | stderr non-empty, contains violation message substring |
| `TestLoopGateDefaultForEdit` | execute-ready; omit `--for` | same result as explicit `--for edit` |
| `TestLoopGateInvalidForFailClosed` | `--for not-a-gate` | exit **2**; stdout empty or no valid gate JSON |
| `TestLoopGateMissingTaskFlag` | `loop gate` without `--task` | exit **2**; usage hint on stderr |
| `TestLoopGateUnknownTaskOrientBlocked` | random UUID, `--for orient` | exit **1**; JSON with `gate_orient_failed`, `reason_code=task_not_found` |
| `TestLoopGateOrientAllowed` | seeded goal/task/plan | `--for orient` → exit **0**, `allowed=true` |
| `TestLoopGateDoneBlockedVerificationDebt` | change without verification | `--for done` → exit **1**, reason `verification_incomplete` |
| `TestLoopGateExecuteAllowedWhenPending` | execute_pending clear (mirror S01/gate unit test setup) | `--for execute` → exit **0** |
| `TestHelpIncludesLoopGate` | `trace help` | output contains `loop gate --task` and exit-code hint |

Use existing helpers from `loop_test.go`: `addGoalForLoopTest`, `addTaskForLoopTest`, `createCurrentDeepPlanForLoopTest`, `captureStdout`, `run`, temp dir + `init`.

## Implementation notes

- **Harness-agnostic:** no Cursor/Claude imports in `cmd/trace`.
- Blocked exit **1** is intentional for shell hooks: `trace loop gate ... \|\| test $? -eq 1` distinguishes policy block from usage failure (2).
- Do not duplicate SelectNext or policy tables — delegate 100% to `EvaluateGate`.
- JSON to stdout via `fmt.Println` or `json.NewEncoder(os.Stdout)` — match `loop next` style (single line JSON).
- On exit 2, do not emit partial gate JSON.

## Exit criteria

- [ ] All 14 named CLI tests green
- [ ] JSON always valid on stdout for exit 0/1 gate invocations
- [ ] Help includes `loop gate` with usage + exit-code line
- [ ] P19/P20 loop keepers still green: `go test ./cmd/trace -run 'TestLoopNext|TestLoopApply|TestLoopStatus|TestHelpIncludesLoop'`
- [ ] No changes to `internal/loop/gate.go` policy
- [ ] No new MCP tools

## Minimal todos

- [ ] Add `gate` dispatch + `cmdLoopGate` in `loop.go`
- [ ] Implement envelope builder + exit mapping
- [ ] Extend `printLoopHelp` per locked text
- [ ] Add all 14 named tests in `loop_test.go`
- [ ] Run keeper tests; fix regressions only in loop CLI files
- [ ] Board row: status + notes only
