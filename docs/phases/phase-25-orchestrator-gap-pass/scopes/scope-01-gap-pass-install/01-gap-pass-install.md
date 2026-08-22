# Phase 25 / Scope 01 / 01-gap-pass-install

## Metadata
- id: P25-S01-01
- todo_ids: [P25-S01-01]
- role: implementer
- skills: []
- mcps: [user-trace]
- agents: []
- verification: mixed
- hooks: []

## Objective

Deliver the Phase 25 / P25-C install bundle: three bounded harness/docs changes that make default build sessions end with a mandatory gap-pass review and tighten parent orchestrator task ownership.

### Deliverables

| # | Deliverable | Intervention | Files in scope |
|---|-------------|--------------|----------------|
| D1 | Gap-pass prompt text in install bundle | INT-03 | `internal/install/agents.go` (or new `gappass.go` in same package) |
| D2 | Parent orchestrator `TRACE_TASK_ID` ownership doc + hook note | INT-04 | `internal/install/enforcement.go` (add `ParentOrchestratorRule` constant) |
| D3 | Hook drift verification notes | INT-11 | `internal/install/cursorhook.go` (add `HookDriftNote` constant + doc comment) |

## Locked defaults

| Item | Value |
|------|-------|
| Interventions | INT-03, INT-04, INT-11 only |
| Core boundary | No daemon, no HTTP, no new CLI subcommands |
| Product perimeter | `internal/install/` package only (constants + exported strings) |
| Gap-pass trigger | End of build session (not async, not event-driven) |
| Parent enforcement | Document-and-deny via existing `preToolUse` hook path — no new Go binary entrypoints |
| Schema change | None — no SQLite migration in this scope |
| Board rights | Status + notes only |

## References

- `docs/rules/agent-loop-protocol.md`
- `docs/rules/project-rules.md`
- `docs/phases/phase-25-orchestrator-gap-pass/GAP-PASS.md`
- `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md` (INT-03, INT-04, INT-11 rows)
- `docs/phases/phase-24-agent-effectiveness-investigation/DR-HANDOFF.md`
- `internal/install/enforcement.go` (existing hook constants and rule text)
- `internal/install/cursor.go` (existing install target)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Preflight / Plan

Before coding:

1. Read `internal/install/enforcement.go` — understand existing `CursorEnforcementRule`, `AgentEnforcementRule`, hook script constants.
2. Read `internal/install/agents.go` — understand existing agent rule structure and exported strings.
3. Read `internal/install/cursorhook.go` — understand hook schema write path (INT-11 hook drift context).
4. Confirm no existing gap-pass constant (`grep -r "gap.pass" internal/install/`).
5. Produce a 3-item plan (D1/D2/D3) and state files touched before writing any code.

## Role work

### D1 — Gap-pass prompt install (INT-03)

Add a new exported constant (or small set of constants) to `internal/install/` that contains the gap-pass prompt text agents receive after a build session completes.

Minimum content the constant must include:

```
After completing build work, run a mandatory gap pass:
1. Call `trace gap --task "$TRACE_TASK_ID"` (or `trace loop status`) to surface open gaps.
2. For each BLOCKING gap: either fix inline or spawn a new task via `trace add task`.
3. Do not mark the task DONE until `violations[]` is empty.
```

The constant must be:
- Exported (e.g. `GapPassPrompt`)
- Referenced in the relevant `AgentRule` or `Install` path so reviewers can verify wiring
- Tested: add or extend an existing `_test.go` that asserts `GapPassPrompt` is non-empty and contains the word "gap"

### D2 — Parent orchestrator rule (INT-04)

Extend `internal/install/enforcement.go` (or add a sibling constant file) with a `ParentOrchestratorRule` exported string. It must:

- State that the **parent orchestrator** must set `TRACE_TASK_ID` before delegating to subagents
- State that `preToolUse` deny fires when a Write is attempted without an active `TRACE_TASK_ID`
- Reference the existing `CursorLoopGateHookScript` for implementation path
- Include a `failClosed` note: "When `TRACE_TASK_ID` is absent, deny the edit rather than allowing untracked work."

No new shell scripts or binary entrypoints. Document-level constant only; implementation is the existing hook.

### D3 — Hook drift verification note (INT-11)

Add a `HookDriftNote` exported constant (or extend `CursorLoopGateHookScript` doc comment) in `internal/install/cursorhook.go` that states:

- On each Cursor upgrade, run `trace install --check` (or equivalent) to verify `hooks.json` schema compatibility
- The schema field to watch: `preToolUse` deny array format in `.cursor/hooks.json`
- Escalation: if schema drift detected, re-run `trace install` to regenerate

## Todo updates

Update own board row `P25-S01-01` in `docs/TODO/phase-25.md`: set status `done` with evidence notes listing files changed and constants added.

## Exit criteria

- [ ] `GapPassPrompt` constant exists in `internal/install/`, is exported, non-empty, and contains gap-pass trigger instructions
- [ ] `ParentOrchestratorRule` constant exists, references `preToolUse` deny path and `failClosed` semantics
- [ ] `HookDriftNote` constant exists referencing Cursor upgrade verification steps
- [ ] All three constants appear in at least one `_test.go` assertion (non-empty check minimum)
- [ ] No daemon, HTTP, or new CLI subcommand added
- [ ] No SQLite schema migration introduced
- [ ] `go build ./...` passes

## Minimal todos

- [ ] Read `internal/install/` files (enforcement.go, agents.go, cursorhook.go)
- [ ] Add `GapPassPrompt` constant (INT-03) + test assertion
- [ ] Add `ParentOrchestratorRule` constant (INT-04) + test assertion
- [ ] Add `HookDriftNote` constant (INT-11) + test assertion
- [ ] Run `go build ./...` and `go test ./internal/install/...`
- [ ] Update board row P25-S01-01 to `done` with evidence
