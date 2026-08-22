# Phase 25 / Scope 01 / 02a-gap-pass-prompt-wiring

## Metadata
- id: P25-S01-02a
- todo_ids: [P25-S01-02a]
- role: implementer
- skills: []
- mcps: [user-trace]
- agents: []
- verification: automated
- hooks: []

## Objective
Fix INT-03 reachability: `internal/install.GapPassPrompt` must be referenced (not orphaned) by an install/agent rule path so the mandatory gap-pass instruction text is surfaced after build sessions.

## References
- `docs/rules/agent-loop-protocol.md`
- `docs/rules/project-rules.md`
- `docs/phases/phase-25-orchestrator-gap-pass/scopes/scope-01-gap-pass-install/02-scope-review.md` (review finding: HIGH — orphaned `GapPassPrompt`)
- `internal/install/gappass.go`
- `internal/install/enforcement.go` (install writes Cursor/Claude rule bodies; also writes marker-delimited `AGENTS.md` blocks)
- `internal/install/enforcement_test.go` (extend with wiring assertions)
- `docs/phases/phase-25-orchestrator-gap-pass/scopes/scope-01-gap-pass-install/01-gap-pass-install.md` (locked defaults for D1)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). This is a fresh implementer session.

## Locked defaults
| Item | Value |
|------|-------|
| Interventions | INT-03 only (wire `GapPassPrompt`) |
| Core boundary | No daemon, no HTTP, no new CLI subcommands |
| Product perimeter | Changes confined to `internal/install/` (code + tests). No SQLite migrations. |
| Trigger semantics | End of build session (not async, not event-driven) |

## Preflight / Plan
1. Run repo search to confirm orphan status:
   - `rg GapPassPrompt` should currently show only the constant file + tests.
2. Identify the install surfaces that are part of “install/agent rule path”:
   - Cursor rule body: `.cursor/rules/trace-enforcement.mdc` (written by install).
   - Claude fallback rules file: `.claude/trace-enforcement-rules.md` (or `CLAUDE.md` marker block).
   - Marker-delimited `AGENTS.md` enforcement block (written when marker present).
3. Decide where the gap-pass text should be inserted so it’s clearly “after build session completes”.
4. Update wiring by referencing `GapPassPrompt` from the selected rule body function(s).
5. Extend tests to assert the wiring (presence of at least the key phrases).

## Role work
1. Wire `GapPassPrompt` into at least one of the install-written agent instruction outputs (Cursor, Claude, and/or `AGENTS.md` block).
   - The wiring must include a direct code reference to the exported `GapPassPrompt` constant (so `rg GapPassPrompt` finds non-test usages).
2. Update/add `_test.go` coverage:
   - Extend `internal/install/enforcement_test.go` with a keeper test that asserts the installed output contains:
     - `mandatory gap pass`
     - `TRACE_TASK_ID`
     - and either `trace gap --task` or `trace loop status`
3. Keep `ParentOrchestratorRule` and `HookDriftNote` unchanged unless tests force a change.

## Todo updates
Update only board row `P25-S01-02a` status + notes in `docs/TODO/phase-25.md`.

## Exit criteria
- [ ] `rg GapPassPrompt` shows at least one non-test code reference (wiring to install rule path).
- [ ] At least one install surface output (Cursor and/or Claude rules and/or `AGENTS.md` marker block) includes the mandatory gap-pass phrases.
- [ ] `go test ./internal/install/...` passes.
- [ ] Boundary checks: no daemon/HTTP/new CLI subcommand/SQLite migration added by this change (confirm with searches restricted to `internal/install/`).
- [ ] `go build ./...` either passes or fails with the same pre-existing fixture-path error (document in board notes if still failing).

## Minimal todos
- [ ] Add wiring reference(s) from install rule output code to `GapPassPrompt`
- [ ] Add/extend `internal/install/enforcement_test.go` to assert surface output includes gap-pass text
- [ ] `go test ./internal/install/...`
- [ ] Update `docs/TODO/phase-25.md` row `P25-S01-02a` with evidence

