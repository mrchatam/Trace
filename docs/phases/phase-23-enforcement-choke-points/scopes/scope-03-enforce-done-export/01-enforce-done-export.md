# P23-S03-01 — Implement enforce on DONE and strict export

## Metadata
- id: P23-S03-01
- todo_ids: [P23-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- verification: automated

## Objective
Wire optional **`--enforce`** on `trace transition … DONE` and **`--strict`** / **`--enforce`** on `trace seed export` per **S03-00 FINAL locks**. Thin adapter over S01 `EvaluateGate` — **no policy fork**, no changes to `internal/loop/gate.go` policy.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S03-00 planner: [00-PLANNER.md](./00-PLANNER.md)
- S01 library: `internal/loop/gate.go` — `EvaluateGate`, `GateForDone`, `GateForExport`
- Live: `cmd/trace/transition.go`, `cmd/trace/seed.go`, `cmd/trace/loop.go` (gate exit pattern)

## Session start
Follow agent-loop-protocol. Board edits: **status + notes only**.

## Locked defaults (from S03-00 — do not re-debate)

| Item | Value |
|------|-------|
| Transition scope | **`--enforce` applies only when `--to DONE`** (case-insensitive match on `DONE`); ignored for other target states |
| Transition gate | `EvaluateGate(..., GateForDone)` **before** `domain.TransitionTask` |
| Without `--enforce` | **Unchanged** — all existing flags (`--allow-done`, `--as-operator`, `--allow-missing-caps`, `--evidence`) behave exactly as today |
| `--enforce` + escape hatches | Gate runs **first**; if gate blocks, transition never runs (even with `--allow-done`). If gate allows, domain escape hatches apply as today |
| Transition block exit | **1** (policy block — same harness contract as `trace loop gate` blocked); stderr = `violations[0].Message`; **no** stdout success JSON |
| Transition internal error | **2** (`exitFail`) — store open, `EvaluateGate` returned `err` |
| Export flags | `--strict` enables validation; `--enforce` requires `--strict` (if `--enforce` without `--strict` → usage error exit **1**) |
| Export gate | `EvaluateGate(..., GateForExport)` — same policy as `done` in S01 (verification debt, open regression, deliberation incomplete) |
| Export task scope | Scan every task where `work_state ∉ {DONE, SKIPPED, STALE}`; optional **`--task <id>`** narrows to one task |
| Export doc keys | After `BuildSeedDocument`: require `version == 1`; `goals` and `tasks` slices non-nil (empty OK) |
| `--strict` alone | Print violations to stderr (one line per violation, prefix `seed export strict: task <id>:`); exit **0** |
| `--strict --enforce` | Exit **1** on any violation; **do not write** file or stdout payload |
| Config `enforce` mode | **No effect** in S03 — S04 may add stderr verbosity; never auto-enable `--enforce` |
| Evaluator | S01 `EvaluateGate` only — **do not** duplicate DONE checks from `domain.TransitionTask` (review PASS, FAIL, caps) inside gate |

### Two-layer DONE model (FINAL — implementers must preserve)

| Layer | What it checks | When |
|-------|----------------|------|
| **Gate (`GateForDone`)** | Deliberation policy: verification debt, open regression, blocking uncertainty, phase not ready | Only when `--enforce` |
| **Domain (`TransitionTask`)** | Review PASS + `--as-operator`, linked FAIL, missing caps, legal edges | Always (unchanged) |

Gate does **not** check review PASS (S01-00 lock). `--enforce` is opt-in **deliberation** enforcement on top of existing domain DONE policy.

## Files to create/modify

| File | Action |
|------|--------|
| `cmd/trace/transition.go` | Add `--enforce` flag; pre-call `EvaluateGate` when `--to DONE` |
| `cmd/trace/seed.go` | Add `--strict`, `--enforce`, optional `--task`; validation before write |
| `cmd/trace/help.go` | Extend `transition` + `seed export` help lines |
| `cmd/trace/cli_test.go` or `transition_test.go` / `seed_test.go` | Named CLI tests below (prefer colocated with existing transition/seed tests) |

**Do not modify:** `internal/loop/gate.go`, `internal/deliberation/select.go`, `internal/domain/task_state.go` DONE policy, MCP tools.

## Implementation sketch — transition

```go
enforce := fs.Bool("enforce", false, "Run deliberation gate before DONE (--to DONE only)")

// after parse, before TransitionTask:
if *enforce && strings.EqualFold(strings.TrimSpace(*to), store.WorkStateDone) {
    plan := planner.New(st)
    allowed, violations, err := loop.EvaluateGate(ctx, svc, plan, st, *taskID, loop.GateForDone)
    if err != nil {
        fmt.Fprintf(os.Stderr, "transition: gate: %v\n", err)
        return exitFail
    }
    if !allowed {
        if len(violations) > 0 {
            fmt.Fprintf(os.Stderr, "transition: %s\n", violations[0].Message)
        } else {
            fmt.Fprintf(os.Stderr, "transition: gate blocked\n")
        }
        return exitGateBlocked // 1 — match loop.go
    }
}
// existing svc.TransitionTask(...) unchanged
```

Update usage string to mention `[--enforce]` on the DONE note line.

Mirror `cmdLoopGate` for: `resolveRoot`, `store.Open`, `planner.New(st)`, defer close.

## Implementation sketch — seed export

```go
strict := fs.Bool("strict", false, "Validate export honesty before write")
enforce := fs.Bool("enforce", false, "Fail closed on strict violations (requires --strict)")
taskFilter := fs.String("task", "", "Evaluate export gate for one task only (optional)")

// usage: reject --enforce without --strict

// after BuildSeedDocument, before json.Marshal/write:
if *strict {
    if doc.Version != 1 {
        // structural violation
    }
    violations := collectExportViolations(ctx, svc, plan, st, *taskFilter)
    for _, v := range violations {
        fmt.Fprintf(os.Stderr, "seed export strict: task %s: %s\n", v.TaskID, v.Message)
    }
    if *enforce && len(violations) > 0 {
        return exitGateBlocked // 1; no write
    }
}
// existing marshal + write path
```

`collectExportViolations`:

1. If `--task` set: evaluate that task only (error if task missing → exit **2**).
2. Else: `st.ListTasks()`, filter non-terminal states, call `EvaluateGate(..., GateForExport)` per task.
3. Collect blocked results as `{TaskID, Violation}` structs (at most one violation per task per S01).
4. Stable order: by task `created_at` then `id` (same as `ListTasks`).

**No write on enforce failure:** when `-o path` set, file must not exist or must remain unchanged (use temp write + rename only after validation passes, or stat before/after in tests).

## Named CLI tests (minimum — S03-01 must implement all)

### Transition

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestTransitionDoneEnforceBlocksVerificationDebt` | init + goal/task/plan critiqued + change without verification + review PASS + `--as-operator` | `--enforce` → exit **1**, task stays non-DONE, stderr contains `verification_incomplete` |
| `TestTransitionDoneWithoutEnforceUnchanged` | same debt setup + review PASS + `--as-operator`, **no** `--enforce` | exit **0**, task → DONE (domain allows; gate not consulted) |
| `TestTransitionDoneEnforceAllowsClean` | `seedFullCycleClear` equivalent via CLI + review PASS + `--as-operator` + `--enforce` | exit **0**, task → DONE |
| `TestTransitionDoneEnforcePreservesAllowDone` | no review; `--allow-done` without `--enforce` | exit **0** + WARNING (existing test pattern) |
| `TestTransitionDoneEnforceBlocksDespiteAllowDone` | verification debt + `--allow-done --enforce` | exit **1**, task not DONE |
| `TestTransitionDoneEnforceIgnoredForNonDone` | `--to IN_PROGRESS --enforce` | exit **0** (enforce no-op) |
| `TestTransitionDoneEnforcePreservesDomainReviewGate` | PASS without `--as-operator`, no `--enforce` | exit **≠ 0** (domain block unchanged) |
| `TestTransitionDoneEnforceStderrHint` | blocked enforce | stderr non-empty, contains violation message substring |
| `TestHelpIncludesTransitionEnforce` | `trace help` | mentions `--enforce` on transition / DONE |

Reuse helpers from `loop_test.go` / `cli_test.go`: `run`, `captureStderr`, `addGoalForLoopTest`, `createCurrentDeepPlanForLoopTest`, review create/set flow from `TestReviewCreateSetDone`.

### Seed export

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestSeedExportStrictEnforceNoWriteOnViolation` | debt task + `-o out.json --strict --enforce` | exit **1**; `out.json` absent or unchanged |
| `TestSeedExportStrictWithoutEnforceExitZero` | debt task + `--strict` (stdout or `-o`) | exit **0**; stderr mentions violation; export payload still written when `-o` |
| `TestSeedExportStrictEnforceBlocksOpenRegression` | open regression fixture | exit **1**, no write with `-o` |
| `TestSeedExportStrictCleanAllowsWrite` | full cycle clear | `--strict --enforce -o` → exit **0**, valid JSON file |
| `TestSeedExportWithoutStrictUnchanged` | any seeded project | plain `seed export` → exit **0** (P17 keepers unchanged) |
| `TestSeedExportEnforceRequiresStrict` | `--enforce` without `--strict` | exit **1** usage |
| `TestSeedExportStrictTaskFilter` | two tasks, one clean one debt; `--strict --task <debt-id> --enforce` | exit **1**; `--task <clean-id> --strict --enforce` → exit **0** |
| `TestHelpIncludesSeedExportStrict` | `trace help` | mentions `--strict` and `--enforce` on seed export |

## Keeper tests (must stay green)

```bash
go test ./cmd/trace -run 'TestTransitionDoneEnforce|TestSeedExportStrict|TestReviewCreateSetDone|TestAllowDoneWarnsOnStderr|TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit|TestLoopGate|TestLoopNext|TestLoopApply|TestLoopStatus|TestHelpIncludesLoop'
go test ./internal/loop/... -run Gate
```

## Help text (FINAL — add to `help.go`)

**Transition** — append to existing DONE note:

```
Optional --enforce runs deliberation gate (GateForDone) before DONE; exit 1 when
blocked. Without --enforce, behavior unchanged. --enforce does not bypass review
PASS/--as-operator domain checks when gate allows.
```

**Seed export** — extend block:

```
  seed export [-o <file>] [--strict] [--enforce] [--task <id>]
                        … existing export description …
                        --strict validates export honesty (GateForExport per active
                        task, or --task only). --enforce requires --strict; exit 1 and
                        no write on violation. --strict alone warns on stderr, exit 0.
```

## Implementation notes

- Share `exitGateBlocked = 1` with `loop.go` (same const in `main` package — already defined in `loop.go`; reuse, do not redefine conflicting values).
- Transition enforce block: **no** `trace.loop.gate.v1` JSON on stdout (S03 scope — stderr only; hooks use `trace loop gate` for JSON).
- Export `--strict` without violations: no stderr noise (or single optional line — prefer silent).
- Terminal task filter constants: `store.WorkStateDone`, `WorkStateSkipped`, `WorkStateStale`.
- Do not read `.trace/config.json` in S03 (S04 scope).
- CONTRIBUTING: one bullet under portable graph noting opt-in `--enforce` on DONE transition and `seed export --strict --enforce` for CI/harness (optional doc touch — only if CONTRIBUTING already mentions export path).

## Exit criteria

- [ ] All named CLI tests green (17 minimum: 9 transition + 8 export)
- [ ] Backward compat: no flags → identical behavior to pre-S03
- [ ] `--enforce` fail-closed; export enforce never writes on violation
- [ ] Gate reuse verified: single `EvaluateGate` call site per command path
- [ ] Help updated
- [ ] P17 seed export keepers green when not strict
- [ ] No changes to `internal/loop/gate.go` policy

## Minimal todos

- [ ] Add `--enforce` to `cmdTransition` + gate pre-check
- [ ] Add `--strict` / `--enforce` / `--task` to `cmdSeedExport` + `collectExportViolations`
- [ ] Implement all named tests
- [ ] Update help strings
- [ ] Run keeper commands; board row status + notes only

## Next

**P23-S03-02** review after this row is `done`.
