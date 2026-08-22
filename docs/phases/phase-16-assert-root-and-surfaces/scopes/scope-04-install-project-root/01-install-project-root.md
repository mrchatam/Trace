# P16 / S04 / 01 — install `-C` vs cwd (FINAL locks from 00-PLANNER)

## Metadata
- id: P16-S04-01
- todo_ids: [P16-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-68** per sibling **00-PLANNER FINAL**: thread Abs parsed `-C` into `InstallOpts.ProjectRoot` for detect / claude / uninstall-when-project-relative. Cursor STABLE home detect unchanged. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- Live: `cmd/trace/install.go`; `cmd/trace/root.go` `parseGlobal` / `resolveRoot`; `internal/install/claude.go` / `cursor.go`
- Dogfood: [`POST-P15-DOGFOOD.md`](../../../../../experiments/POST-P15-DOGFOOD.md); [`_post_p15/install_C_cwd.json`](../../../../../experiments/_post_p15/install_C_cwd.json)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks (which subcommands take ProjectRoot; Cursor STABLE home; no `cli:install`; no marker-rule change; no PID kill).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | Thin `cmd/trace/install.go` (+ tests). Library already root-aware — **do not** change `hasClaudeMarker` / Cursor `cursorMCPPath` |
| Thread | Keep `cmdInstall(root, args)`. `resolveRoot(root)` → `InstallOpts.ProjectRoot` for **detect**, **claude**, **uninstall** |
| Cursor install | **Do not** set ProjectRoot as mcp.json. STABLE `$HOME/.cursor/mcp.json` / `--mcp-json` / `HomeDir` |
| Detect | `ListTargets` must receive ProjectRoot (not `InstallOpts{}`). Claude reason cites Abs `-C`; cursor reason stays home |
| No `-C` | `parseGlobal` cwd → Abs; detect may cite Abs(cwd) instead of `"."` |
| Markers | Unchanged: `.claude/` **or** `CLAUDE.md` |
| Tip | Keep `installCursorReloadTip` / `CursorReloadTip` print+write; **no** PID kill |
| Ungated | **Never** Assert / `cli:install` |
| Compat | **14**; no **015+** |
| Forbidden | New MCP tools; daemon; changing CONDITIONAL rules; relocating Cursor home via `-C` |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| CLI install | `cmd/trace/install.go` | Keep `root`; Abs; pass ProjectRoot on detect/claude/uninstall; cursor install unchanged |
| CLI tests | `cmd/trace/install_test.go` | Named DashC tests (Chdir+Cleanup like `cli_test.go`; `runCapture` for fail paths) |
| Library production | `internal/install/*` | **Zero** unless a proven library bug (none live) |
| Library tests | `internal/install/install_test.go` | Keepers only; optional extra ListTargets+ProjectRoot test is not a CLI substitute |
| MCP / domain / store | **Zero** | Nine tools; ungated install; compat 14 |

Do **not** add a subcommand `--project` flag. Global `-C` only.

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestInstallClaudeDashCRefuseCitesProjectRoot` | Chdir cwdDir (no marker); `-C projectDir` (no marker); refuse cites Abs(projectDir), not cwdDir |
| `TestInstallClaudeDashCIgnoresCwdMarker` | cwdDir **has** `.claude/`; projectDir none; `-C projectDir` still refuse citing projectDir |
| `TestInstallClaudeDashCWriteUsesProjectRoot` | projectDir has marker; `--write` under `-C` writes there; cwdDir untouched |
| `TestInstallDetectDashCClaudeReasonCitesRoot` | detect JSON claude detected + reason contains Abs(projectDir); not `under .`; not cwdDir |
| `TestInstallDetectDashCCursorHomeUnchanged` | detect cursor reason does **not** contain projectDir |
| `TestInstallUninstallClaudeDashCUsesProjectRoot` | uninstall `-C` removes projectDir file; cwd decoy untouched |
| Keepers | `TestInstallCursorPrintReloadTip`, `TestInstallCursorWriteMergeBackup`, `TestInstallDetectListsCursorStable`, `TestInstallConditionalRefusesWithoutMarker`, `TestInstallConditionalWritesWithMarker`, `TestCLIAddDeniedFailClosed`, `TestToolNamesRegistered` (nine), `TestMCPVirginProjectDoesNotMkdir`, `TestOpenCreatesDBAndMigratesIdempotent` (v14) |

TDD: named DashC tests first (red on live drop-root / empty detect opts), then thread Abs ProjectRoot (green).

## Role work
1. TDD named CLI tests (red: detect `under .` / claude cites cwd).
2. `cmdInstall` keeps `root`; `resolveRoot`; ProjectRoot on detect/claude/uninstall; cursor install stays home/`--mcp-json`.
3. Prove green: named tests + locked verify cmds. Do **not** add MCP tools, Assert, or marker-rule changes.
4. Board **status + Notes only** → next **P16-S04-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/install/... ./internal/mcp/... ./internal/store/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallConditional|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallClaudeDashC|TestInstallDetectDashC|TestInstallUninstallClaudeDashC|TestInstallCursorPrintReloadTip|TestInstallCursorWriteMergeBackup|TestCLIAddDeniedFailClosed'

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Product bar = `./cmd|internal|evals`. Compat ceiling **14**. Named `cmd/trace` tests are **CGO1** (R4: CGO0 `./cmd/trace/...` tree-sitter). Do **not** fail this row for R3 graphify space-in-path on full-module `./...` if present outside product pkgs.

## Exit criteria
- [ ] DF-68: `trace -C <root> install claude` / `install detect` / claude uninstall use Abs `<root>` not process cwd
- [ ] Cursor STABLE home detect unchanged; print/write reload tip keepers green; no PID kill
- [ ] `install` still ungated; nine MCP tools; no 015; CONDITIONAL markers unchanged
- [ ] Named tests pass; locked verify cmds PASS
- [ ] Board Notes → **P16-S04-02**

## Minimal todos
- [ ] Named DashC tests (red → green)
- [ ] Thread `resolveRoot` → ProjectRoot on detect/claude/uninstall
- [ ] Cursor install/detect home path unchanged
- [ ] Locked verify suite
- [ ] Board status + Notes only
