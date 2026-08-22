# P23-S05-02 — Review harness install rules

## Metadata
- id: P23-S05-02
- todo_ids: [P23-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent review: install snippets call gate CLI correctly; **git-hook unchanged**; rules opt-in; further-enforcement scope (status, DONE, export) present; print/`--write` contract preserved.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S05-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S05-01 deliverable: [01-harness-install-rules.md](./01-harness-install-rules.md)
- S02 gate CLI: exit 0/1/2 contract (hook must not parse stdout on exit 2)
- P22 git-hook: `internal/install/githook.go` fragment unchanged

## Session start
Follow agent-loop-protocol. Fresh reviewer context. Board edits: **status + notes only**.

## Keeper tests (must re-run — all green)

```bash
go test ./internal/install/... -run 'TestInstall|Enforcement|CursorHook|GitHook'
go test ./cmd/trace -run 'TestInstall|TestHelp'
go test ./cmd/trace -run 'TestLoopGate|TestLoopNext|TestLoopApply|TestLoopStatus'
```

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Loop gate in cursor rules | `grep loop gate` in `.cursor/rules/trace-enforcement.mdc` fixture or test output |
| Loop gate in claude rules | CLAUDE.md block or fallback file references gate CLI |
| AGENTS.md markers | `# begin-trace-enforcement` / `# end-trace-enforcement` present after install |
| Marker idempotent | Second `--write` does not duplicate markers (test output) |
| cursor-hook script | `.cursor/hooks/trace-loop-gate.sh` calls `trace … loop gate --task … --for edit` |
| hooks.json merge | `preToolUse` entry with script path; sibling hooks preserved |
| TRACE_TASK_ID empty | Hook script allows when unset (fail-open) |
| Gate exit 1 → deny | Hook returns deny JSON / exit 2 |
| git-hook unchanged | `traceHookFragment()` still `# begin-trace`; index + seed export commands present |
| git-hook markers distinct | git `# begin-trace` ≠ AGENTS `# begin-trace-enforcement` |
| MCP cursor unchanged | Existing `TestInstallCursor*` green; mcp.json merge + backup behavior intact |
| cursor MCP + rules | `--write` produces both mcp entry and project rules file |
| Print-only cursor | stdout = MCP JSON only; rules hint on stderr |
| Detect cursor-hook | `install detect` lists `cursor-hook` CONDITIONAL |
| Config not auto-written | No `.trace/config.json` created by install --write |
| Further rules | Rules text mentions status violations, `--enforce` DONE, export `--strict` |
| ENFORCEMENT link | Rules / AGENTS block link phase ENFORCEMENT.md |
| Named tests | All tests from 01 prompt present and passing |
| CONTRIBUTING | Harness vs product paragraph present |

## Review checklist

- [ ] **Blocker:** git-hook fragment or behavior changed (P22 regression)
- [ ] **Blocker:** Hook script calls wrong CLI or wrong `--for` value
- [ ] **Blocker:** Gate policy duplicated in install package (SelectNext fork)
- [ ] **Blocker:** AGENTS.md markers collide with git-hook `# begin-trace`
- [ ] **Blocker:** cursor-hook writes post-commit logic (must be pre-edit only)
- [ ] **Blocker:** Missing named tests from 01 prompt
- [ ] **Blocker:** Existing cursor MCP install tests regressed
- [ ] **High:** Install auto-writes `.trace/config.json` with strict mode
- [ ] **High:** Hook fail-closed when `trace` missing (should fail-open per S05-00)
- [ ] **High:** Uninstall leaves hook script or duplicate hooks.json entries
- [ ] **High:** Rules omit status / DONE / export references (user request)
- [ ] **Medium:** Help missing `cursor-hook` or env var docs
- [ ] **Medium:** Print-only emits rules on stdout (breaks MCP JSON piping)
- [ ] **Medium:** Changes to `internal/loop/gate.go` or gate CLI exit codes
- [ ] **Low:** Network fetch or daemon spawn in install path
- [ ] **Nit:** Hook script not executable (mode != 0755)

## git-hook unchanged verification (walk through)

1. `trace install git-hook --write` in temp git repo.
2. Read `.git/hooks/post-commit` — must contain:
   - `# begin-trace` / `# end-trace` (not `-enforcement`)
   - `trace … index`
   - `trace … seed export`
3. Compare fragment to pre-S05 baseline — no diff in hook body.
4. `TestInstallGitHookUnchanged` (or equivalent) passes.

## Harness gate contract verification

Confirm hook + rules align with S02 exit semantics:

```bash
export TRACE_TASK_ID="<uuid>"
export TRACE_PROJECT_ROOT="/path/to/project"
trace loop gate --task "$TRACE_TASK_ID" --for edit
# exit 0 → hook allows; rules say proceed
# exit 1 → hook denies; rules say read recommended_phase from stdout JSON
# exit 2 → hook allows (fail-open); rules say surface stderr
```

## S06 handoff verification

Phase verify (S06) will smoke:

```bash
trace install cursor --write
trace install cursor-hook --write
# rules present; git-hook still P22; gate CLI green
```

Confirm install detect JSON includes all four targets: `claude`, `cursor`, `cursor-hook`, `git-hook`.

## Spawn policy

- **blocker/high:** inline fix if ≤10 lines and zero policy change; else spawn `P23-S05-02a` implement + `02b` review immediately below this row
- **medium:** prefer spawn unless trivial typo
- Do not rewrite S05-00/S05-01 `done` prompts

## Exit criteria

- [ ] No open blocker/high without pending forward row
- [ ] Confidence **medium** or **high** with command output in Notes
- [ ] Residual risks listed if medium (e.g. Cursor tool name drift in matcher)
- [ ] APPROVE or spawn documented on board

## Minimal todos

- [ ] Re-run keeper tests; paste pass summary in Notes
- [ ] Walk git-hook fragment — unchanged
- [ ] Walk hook script gate call + env contract
- [ ] Walk AGENTS.md idempotent merge
- [ ] Verify all named tests exist
- [ ] Confirm no gate.go / loop gate CLI diffs in S05 scope
- [ ] Set row done with confidence + residuals
