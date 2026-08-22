# P23-S05-01 — Implement harness install rules

## Metadata
- id: P23-S05-01
- todo_ids: [P23-S05-01]
- role: implementer
- skills: [incremental-implementation, documentation-and-adrs]
- verification: automated

## Objective
Extend `trace install` with enforcement rules, AGENTS.md block, optional cursor-hook per **S05-00 FINAL locks**. Thin install-layer only — **no gate policy changes**.

## References
- S05-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- Live: `internal/install/{cursor,claude,githook}.go`, `cmd/trace/install.go`
- S02 gate CLI: `cmd/trace/loop.go` — exit 0/1/2 (hook contract)
- Cursor hooks: project `.cursor/hooks.json` schema version 1

## Session start
Follow agent-loop-protocol. Board edits: **status + notes only**.

## Locked defaults (from S05-00 — do not re-debate)

| Item | Value |
|------|-------|
| cursor/claude | Extend existing targets — rules + AGENTS.md; MCP paths unchanged |
| cursor-hook | New target `cursor-hook`, tier CONDITIONAL |
| AGENTS.md | `# begin-trace-enforcement` / `# end-trace-enforcement` idempotent merge |
| git-hook | **Untouched** — no edits to `traceHookFragment()` or P22 behavior |
| agents | Bundled SQLite profiles unchanged |
| Config | Do **not** auto-write `.trace/config.json`; document default off |
| Env | `TRACE_TASK_ID`, `TRACE_PROJECT_ROOT` documented in rules + help |

## Files to create/modify

| File | Action |
|------|--------|
| `internal/install/types.go` | Add `TargetCursorHook` |
| `internal/install/registry.go` | Register `cursorHookTarget{}` |
| `internal/install/markers.go` | Shared marker upsert/remove (parameterized markers) |
| `internal/install/enforcement.go` | Rules bodies, AGENTS/CLAUDE merge helpers |
| `internal/install/cursor.go` | Rules file + AGENTS.md; pass-through ProjectRoot |
| `internal/install/claude.go` | Claude rules + AGENTS.md |
| `internal/install/cursorhook.go` | hooks.json + shell script install/uninstall/detect |
| `internal/install/githook.go` | Optional refactor to shared markers only — **zero behavior change** |
| `cmd/trace/install.go` | `ProjectRoot` for cursor; `cursor-hook` case |
| `cmd/trace/help.go` | New targets + env vars |
| `internal/install/enforcement_test.go` | Named unit tests |
| `cmd/trace/install_test.go` | CLI help / integration tests as needed |
| `CONTRIBUTING.md` | Harness vs product enforcement paragraph |

**Do not modify:** `internal/loop/gate.go`, `cmd/trace/loop.go` gate path, `githook.go` fragment content.

## Implementation sketch

### `internal/install/enforcement.go`

```go
const (
    enforcementBeginMarker = "# begin-trace-enforcement"
    enforcementEndMarker   = "# end-trace-enforcement"
)

func EnforcementRulesMarkdown() string { /* locked body from S05-00 */ }
func AgentsEnforcementBlock() string   { /* locked AGENTS block */ }

func UpsertAgentsMD(root string) error { /* merge into AGENTS.md */ }
func StripAgentsMD(root string) error  { /* uninstall */ }

func UpsertClaudeRules(root string) error { /* CLAUDE.md markers or fallback file */ }
func StripClaudeRules(root string) error

func WriteCursorRulesMDC(root string) error { /* .cursor/rules/trace-enforcement.mdc */ }
func RemoveCursorRulesMDC(root string) error
```

Reuse parameterized `upsertMarkedFragment(content, begin, end, fragment)` from `markers.go`.

### `internal/install/cursor.go` — extend `Install`

```go
func (c cursorTarget) Install(opts InstallOpts) error {
    // 1. Existing MCP upsert / print (unchanged stdout JSON on print-only)
    // 2. If opts.ProjectRoot != "":
    //    print-only → rules hint on ErrOut
    //    --write → WriteCursorRulesMDC + UpsertAgentsMD
}
```

`Uninstall`: existing MCP removal + `RemoveCursorRulesMDC` + `StripAgentsMD`.

### `internal/install/cursorhook.go`

```go
type cursorHookTarget struct{}

func (cursorHookTarget) ID() string   { return TargetCursorHook }
func (cursorHookTarget) Tier() string { return TierConditional }

func (cursorHookTarget) Detect(opts InstallOpts) DetectResult {
    // .cursor/ under projectRoot OR hooks.json contains trace-loop-gate.sh
}

func (cursorHookTarget) Install(opts InstallOpts) error {
    // refuse without .cursor/ when --write
    // write script 0755, merge hooks.json, UpsertAgentsMD
}

func (cursorHookTarget) Uninstall(opts InstallOpts) error {
    // remove script, remove matching preToolUse entry, StripAgentsMD
}
```

### `cmd/trace/install.go`

```go
case install.TargetCursorHook:
    return cmdInstallCursorHook(abs, args[1:])

func cmdInstallCursor(args []string) int {
    // add resolveRoot abs → opts.ProjectRoot = abs
}
```

## Named tests (minimum — S05-01 must implement all)

### Cursor rules + AGENTS.md

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestInstallCursorIncludesLoopGateRule` | temp project + `.cursor/`; `cursor --write` | `.cursor/rules/trace-enforcement.mdc` contains `loop gate` and `TRACE_TASK_ID` |
| `TestInstallAgentsMDEnforcementBlock` | temp project; `cursor --write` | `AGENTS.md` contains begin/end enforcement markers + gate command |
| `TestInstallAgentsMDMarkersIdempotent` | run `cursor --write` twice | single block; updated content; no duplicate markers |
| `TestInstallCursorRulesPrintOnlyStderr` | print-only (no `--write`) | stdout valid MCP JSON; stderr mentions loop gate |
| `TestInstallCursorUninstallRemovesRules` | install then uninstall cursor | rules file gone; AGENTS block stripped; MCP trace removed |

### Claude rules

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestInstallClaudeIncludesLoopGateRule` | `.claude/` marker; `claude --write` | `CLAUDE.md` or fallback file references `loop gate` |
| `TestInstallClaudeRulesWithCLAUDEmd` | `CLAUDE.md` present | marker block inside CLAUDE.md |
| `TestInstallClaudeRulesFallbackFile` | no CLAUDE.md | `.claude/trace-enforcement-rules.md` created |

### cursor-hook

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestInstallCursorHookCallsGate` | `.cursor/` + `--write` | script contains `loop gate` and `--for edit` |
| `TestInstallCursorHookPreToolUseMatcher` | after install | hooks.json has preToolUse entry with trace script path |
| `TestInstallDetectIncludesCursorHook` | `install detect` | list includes `cursor-hook` with CONDITIONAL tier |
| `TestInstallDetectCursorHookConditional` | no `.cursor/` | detected=false with reason |
| `TestInstallCursorHookUninstallRemovesScript` | install then uninstall | script removed; hooks.json entry removed |

### Idempotency + regression

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestInstallEnforcementIdempotent` | cursor + cursor-hook `--write` twice | no duplicate hooks.json entries; files overwritten cleanly |
| `TestInstallGitHookUnchanged` | install git-hook `--write` | fragment still `# begin-trace` (not enforcement markers); index/export commands present |
| `TestInstallCursorMCPUnchanged` | existing cursor MCP tests | all green — MCP merge behavior preserved |

### Help

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestHelpIncludesCursorHook` | `trace help` | mentions `cursor-hook`, `TRACE_TASK_ID`, loop gate |

## CONTRIBUTING paragraph (FINAL text)

Add under harness / agent workflow section:

> **Harness vs product enforcement:** Trace product gates (`trace loop gate`, `--enforce` on DONE/export, status `violations[]`) read SQLite evidence and are authoritative. Harness install (`trace install cursor|claude|cursor-hook`) writes Cursor/Claude rules and optional pre-edit hooks that **call** those CLIs — best-effort reminders, not a second policy engine. Default `.trace/config.json` enforce mode is off; teams opt in explicitly. See `docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md`.

## Keeper commands (after implementation)

```bash
go test ./internal/install/... -run 'TestInstall|Enforcement|CursorHook|GitHook'
go test ./cmd/trace -run 'TestInstall|TestHelp'
```

## Exit criteria

- [ ] `--write` installs rules without breaking existing cursor MCP / claude MCP install
- [ ] AGENTS.md merge idempotent; CLAUDE path works with and without CLAUDE.md
- [ ] cursor-hook separate from git-hook; preToolUse calls gate CLI
- [ ] Further rules cover status / DONE `--enforce` / export `--strict` (in rules text)
- [ ] Default enforce remains off (no auto config write)
- [ ] All named tests green
- [ ] No git commit wrap; no daemon; no network

## Next

**P23-S05-02** after this row is `done`.
