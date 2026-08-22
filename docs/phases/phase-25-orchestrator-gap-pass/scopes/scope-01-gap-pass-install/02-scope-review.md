# Phase 25 / Scope 01 / 02-scope-review

## Metadata
- id: P25-S01-02
- todo_ids: [P25-S01-02]
- role: reviewer
- skills: []
- mcps: [user-trace]
- agents: []
- verification: mixed
- hooks: []

## Objective

Independently verify P25-S01-01 against locked defaults, Trace laws, and Phase 25 scope boundaries. Spawn inline fixes or implement+review pairs for any blocker/high findings.

## References

- `docs/rules/agent-loop-protocol.md`
- `docs/rules/project-rules.md`
- `docs/phases/phase-25-orchestrator-gap-pass/GAP-PASS.md`
- `docs/phases/phase-25-orchestrator-gap-pass/scopes/scope-01-gap-pass-install/01-gap-pass-install.md`
- `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md`
- `internal/install/` (all files touched by P25-S01-01)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). This is a **fresh subagent**; do not share context with the P25-S01-01 implementer session.

## Preflight / Plan

1. Read `docs/TODO/phase-25.md` P25-S01-01 Notes to see claimed evidence.
2. Read `01-gap-pass-install.md` exit criteria and locked defaults.
3. List the files actually touched (via git diff or direct read).
4. Plan review pass order: constants → tests → board row → boundary check.

## Verify checklist

### INT-03 — Gap-pass default (D1)

- [ ] `GapPassPrompt` (or equivalent) exported constant exists in `internal/install/`
- [ ] Constant is non-empty and contains mandatory gap-pass trigger language (`gap`, `TRACE_TASK_ID`, `violations`)
- [ ] Constant is reachable from the install/agent rule path (referenced or wired — not orphan)
- [ ] At least one `_test.go` asserts the constant is non-empty

### INT-04 — Parent orchestrator rule (D2)

- [ ] `ParentOrchestratorRule` (or equivalent) exported constant exists
- [ ] Constant references `TRACE_TASK_ID` ownership by parent orchestrator
- [ ] Constant references `preToolUse` deny path and `failClosed` semantics
- [ ] No new shell scripts, binary entrypoints, or daemon added
- [ ] At least one `_test.go` asserts the constant is non-empty

### INT-11 — Hook drift verification (D3)

- [ ] `HookDriftNote` (or equivalent) exported constant exists
- [ ] Constant references Cursor upgrade + `hooks.json` schema check
- [ ] At least one `_test.go` asserts the constant is non-empty

### Boundary checks

- [ ] No daemon, HTTP server, or new CLI subcommand introduced
- [ ] No SQLite migration introduced
- [ ] All changes confined to `internal/install/` package
- [ ] `go build ./...` passes (run and confirm)
- [ ] `go test ./internal/install/...` passes (run and confirm)

### Board hygiene

- [ ] P25-S01-01 row in `docs/TODO/phase-25.md` is marked `done` with evidence notes
- [ ] No `done` row history rewritten

## Spawn policy

Follow agent-loop-protocol reviewer loop (quality-first):

| Finding | Action |
|---------|--------|
| Blocker | Small inline fix (if trivial) **or** spawn `P25-S01-02a` implement + `P25-S01-02b` review immediately below this row |
| High | Spawn implement+review pair if not trivially fixable |
| Medium | Prefer spawn; skip only if absolutely trivial and evidenced |
| Low / nit | Note in board row; no spawn required |

Spawn prompt files must follow the required prompt skeleton from `agent-loop-protocol.md`.

## Exit criteria

- [ ] All INT-03, INT-04, INT-11 verify checklist items pass with evidence
- [ ] All boundary checks pass (build + test confirmed)
- [ ] No open blocker/high findings without a spawned follow-up row
- [ ] Confidence: **high** (or **medium** with explicit residuals listed)
- [ ] Own board row P25-S01-02 updated to `done` with findings summary

## Minimal todos

- [ ] Read P25-S01-01 board notes and touched files
- [ ] Run verify checklist (INT-03 / INT-04 / INT-11 + boundary)
- [ ] Run `go build ./...` and `go test ./internal/install/...`
- [ ] Record findings by severity
- [ ] Spawn if needed; else mark done with confidence level
