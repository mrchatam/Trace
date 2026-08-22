# P23-S02-00 — Loop gate CLI planner

## Metadata
- id: P23-S02-00
- todo_ids: [P23-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock `trace loop gate` CLI + `trace.loop.gate.v1` JSON schema. **No product Go this row.**

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S01 evaluator API — **done** (`internal/loop/gate.go`): `loop.EvaluateGate(ctx, dom, plan, st, taskID, gateFor)` → `(allowed, []Violation, err)`; `loop.GateFor{Orient,Edit,Execute,Done,Export}`; `loop.Violation` JSON-tagged for envelope `violations[]`.
- Live: `cmd/trace/loop.go`, `internal/loop/`

### S01 handoff (review P23-S01-02 — do not re-debate)

- CLI is a thin wrapper: open store/domain/planner like `loop next`, call `EvaluateGate`, emit JSON, map exit 0/1.
- When blocked, set top-level `recommended_phase` and `reason_code` from `violations[0]` (S01 returns at most one violation).
- Orient infra failures use violation `code` **`gate_orient_failed`**; policy blocks use **`premature_implementation`** (matches `domain.PrematureImplementation.Code()`).
- Default `--for edit` aligns with harness edit choke point.

## Live touch points (P23-00 inventory)

| File | Reuse for gate CLI |
|------|-------------------|
| `cmd/trace/loop.go` | Add `gate` subcommand dispatch + `cmdLoopGate`; extend `printLoopHelp` |
| `cmd/trace/loop_test.go` | Named gate CLI tests; extend `TestHelpIncludesLoopNext` → `TestHelpIncludesLoop` |
| `cmd/trace/main.go` | Global `exitOK=0`, `exitUsage=1`, `exitFail=2` — **gate uses its own mapping** (see below) |
| `internal/loop/gate.go` | S01 `EvaluateGate` — **no changes in S02** |
| `internal/loop/gate_test.go` | Unit coverage done — CLI tests call through `trace loop gate` only |

## Locked defaults (S02-01 must not re-debate)

| Item | Value |
|------|-------|
| Command | `trace loop gate --task <uuid> [--for orient\|edit\|execute\|done\|export]` |
| Default `--for` | `edit` |
| Exit codes | **0** allowed, **1** blocked, **2** usage/internal error (fail-closed) |
| Stdout | JSON `trace.loop.gate.v1` on **policy paths only** (exit 0 or 1) |
| Stderr | Human-readable one-line hint on block (exit 1); usage/internal errors on exit 2 |
| Schema fields | `schema_version`, `task_id`, `for`, `allowed`, `violations[]`, `recommended_phase`, `reason_code` |
| Implementation | Thin wrapper over S01 `EvaluateGate` |
| MCP | **No new tool** |

### Exit code contract (FINAL — gate-specific)

Gate subcommand **overrides** the general CLI convention (`exitUsage=1`) for harness hook compatibility:

| Exit | Meaning | Stdout JSON | Stderr |
|------|---------|-------------|--------|
| **0** | Allowed (`allowed=true`) | Yes — full envelope | Empty (or debug-only) |
| **1** | Blocked (`allowed=false`, `err==nil`) | Yes — full envelope | One-line hint from `violations[0].message` |
| **2** | Usage or internal error | **No** | Human-readable error (missing `--task`, bad `--for`, store open fail, `EvaluateGate` returned `err`) |

Harness hooks MUST treat exit **1** as "policy blocked — read JSON" and exit **2** as "do not parse stdout JSON."

### `trace.loop.gate.v1` schema (FINAL)

#### Top-level envelope

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `schema_version` | string | yes | Always `"trace.loop.gate.v1"` |
| `task_id` | string | yes | From `--task` flag |
| `for` | string | yes | GateFor value (`orient`, `edit`, `execute`, `done`, `export`) |
| `allowed` | bool | yes | Mirrors `EvaluateGate` return |
| `violations` | array | yes | Empty `[]` when allowed; one element when blocked (S01 guarantee) |
| `recommended_phase` | string | when blocked | Lifted from `violations[0].recommended_phase`; omit when allowed |
| `reason_code` | string | when blocked | Lifted from `violations[0].reason_code`; omit when allowed |

#### Violation element (same as S01 `loop.Violation`)

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `code` | string | yes | `premature_implementation` or `gate_orient_failed` |
| `for` | string | yes | Same as envelope `for` |
| `message` | string | yes | Human-readable; also echoed to stderr on block |
| `recommended_phase` | string | when set | From SelectNext or orient fail |
| `reason_code` | string | yes | e.g. `blocking_uncertainty`, `task_not_found` |

#### Example — blocked (edit)

```json
{
  "schema_version": "trace.loop.gate.v1",
  "task_id": "550e8400-e29b-41d4-a716-446655440000",
  "for": "edit",
  "allowed": false,
  "recommended_phase": "INVESTIGATE",
  "reason_code": "blocking_uncertainty",
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

#### Example — allowed (edit)

```json
{
  "schema_version": "trace.loop.gate.v1",
  "task_id": "550e8400-e29b-41d4-a716-446655440000",
  "for": "edit",
  "allowed": true,
  "violations": []
}
```

### Help text (FINAL — S02-01 must add to `printLoopHelp`)

```
  gate --task <id> [--for orient|edit|execute|done|export]
        Check deliberation gate for a task. Emits trace.loop.gate.v1 JSON on stdout.
        Exit 0 when allowed, 1 when blocked, 2 on usage or internal error.
        Default --for is edit (pre-edit harness choke point).
```

Usage line for flag errors:

```
usage: trace loop gate --task <id> [--for orient|edit|execute|done|export]
```

## Planner work

1. [x] Lock JSON schema + exit code contract for harness hooks.
2. [x] Thicken `01-loop-gate-cli.md` with named CLI tests.
3. [x] Document help text / usage line in 01.

## Exit criteria

- [x] 01/02/SCOPE-TODOS runnable alone
- [x] Schema + exit codes locked
- [x] No product Go

## Next

**P23-S02-01** after this row is `done`.
