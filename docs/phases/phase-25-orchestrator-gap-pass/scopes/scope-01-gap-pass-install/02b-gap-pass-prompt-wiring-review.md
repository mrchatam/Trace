# Phase 25 / Scope 01 / 02b-gap-pass-prompt-wiring-review

## Metadata
- id: P25-S01-02b
- todo_ids: [P25-S01-02b]
- role: reviewer
- skills: []
- mcps: [user-trace]
- agents: []
- verification: mixed
- hooks: []

## Objective
Independently verify the P25-S01-02a wiring fix for INT-03 reachability (GapPassPrompt is no longer orphaned) and rerun the required boundary checks.

## References
- `docs/rules/agent-loop-protocol.md`
- `docs/rules/project-rules.md`
- `docs/TODO/phase-25.md` (rows `P25-S01-02` + `P25-S01-02a`)
- `docs/phases/phase-25-orchestrator-gap-pass/scopes/scope-01-gap-pass-install/02a-gap-pass-prompt-wiring.md`
- `internal/install/gappass.go`
- `internal/install/enforcement.go`
- `internal/install/enforcement_test.go`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). This is a fresh reviewer session.

## Preflight / Plan
1. Read `docs/TODO/phase-25.md` P25-S01-02a Notes (claimed evidence + tests).
2. Identify files touched (prefer `git diff`, else use direct code inspection).
3. Plan review order:
   - confirm wiring is non-test usage
   - confirm surface outputs include mandatory phrases
   - run Go build/test boundaries
   - confirm no daemon/HTTP/new CLI subcommand/SQLite migration added

## Verify checklist
### INT-03 — Gap-pass default (D1)
- [ ] `GapPassPrompt` constant exists and is non-empty
- [ ] `rg GapPassPrompt` shows at least one non-test reference from an install/agent rule output path
- [ ] Installed output contains mandatory gap-pass phrases:
  - `mandatory gap pass`
  - `TRACE_TASK_ID`
  - `trace gap --task` or `trace loop status`

### Boundary checks
- [ ] No daemon, HTTP server, or new CLI subcommand introduced (confirm via restricted searches under `internal/install/` and `cmd/trace/`)
- [ ] No SQLite migration introduced (confirm no `internal/store` schema/migrate code touched; restricted searches)
- [ ] `go build ./...` passes OR fails only with the same pre-existing fixture-path error in `similar projects/graphify` (record evidence if failing)
- [ ] `go test ./internal/install/...` passes (run and confirm)

## Spawn policy
- [ ] Spawn implement+review pair if a blocker/high finding exists that is not trivially fixable.

## Exit criteria
- [ ] INT-03 wiring verified with evidence.
- [ ] `go test ./internal/install/...` PASS.
- [ ] Boundary checks satisfied (or pre-existing unrelated failure documented).
- [ ] Own board row `P25-S01-02b` updated to `done` with findings summary (or `failed/blocked` with reason).

## Minimal todos
- [ ] Confirm non-test references to `GapPassPrompt`
- [ ] Confirm rule surface outputs include the mandatory text phrases
- [ ] Run `go test ./internal/install/...`
- [ ] Run `go build ./...` and capture evidence if it fails
- [ ] Update `docs/TODO/phase-25.md` row `P25-S01-02b` accordingly

