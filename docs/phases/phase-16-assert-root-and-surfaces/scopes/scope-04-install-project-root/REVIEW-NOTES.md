# P16-S04-02 REVIEW-NOTES — install `-C` vs cwd / DF-68

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none (`P16-S04-02a` / `02b` not inserted)  
**Next board:** P16-S05-00

Independent review (fresh subagent ≠ implementer). Claims from P16-S04-01 Notes re-verified against live code + locked verify cmds — not trusted alone. Sibling `00-PLANNER.md` is **FINAL**. Marker rules, Cursor STABLE home, `cli:install`, and PID-kill/auto-reload not re-opened.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | `cmdInstall` no longer drops `root`; detect/claude/uninstall get Abs `ProjectRoot` | PASS | `cmd/trace/install.go`: `func cmdInstall(root string, args []string)`; `resolveRoot(root)` then `cmdInstallDetect(abs, …)` / `cmdInstallUninstall(abs, …)` / `cmdInstallClaude(abs, …)`. `ListTargets(install.InstallOpts{ProjectRoot: projectRoot})` — not `InstallOpts{}`. `root.go` `run()` still `cmdInstall(root, cmdArgs)`. No `os.Getwd` in `install.go`. |
| 2 | `-C` refuse/detect reasons cite `<root>` not process cwd | PASS | `TestInstallClaudeDashCRefuseCitesProjectRoot`: Chdir cwdDir (no marker); `-C projectDir` (no marker); exit fail; stderr contains `under `+Abs(projectDir); must not contain Abs(cwdDir). `TestInstallDetectDashCClaudeReasonCitesRoot`: claude `detected: true`; reason contains Abs(projectDir); not `under .`; not cwdDir. Library `claude.go` Detect/Install interpolate `projectRoot(opts)` (Abs from CLI). |
| 3 | Cwd marker does not authorize `-C` other root; write/uninstall land under `-C` | PASS | `TestInstallClaudeDashCIgnoresCwdMarker`: cwdDir has `.claude/`; projectDir none; still refuse citing projectDir. `TestInstallClaudeDashCWriteUsesProjectRoot`: write → `projectDir/.claude/trace-mcp.json`; cwdDir untouched. `TestInstallUninstallClaudeDashCUsesProjectRoot`: same `-C` removes project file; cwd decoy untouched. |
| 4 | Cursor STABLE home detect unchanged (`-C` not mcp.json) | PASS | `cmdInstallCursor(args)` — no `projectRoot` param; `InstallOpts` is Write/Bin/MCPJSON/Out/ErrOut only (no `ProjectRoot`). `cursorMCPPath` still MCPJSON / HomeDir / `$HOME/.cursor/mcp.json`. `TestInstallDetectDashCCursorHomeUnchanged`: cursor reason does **not** contain projectDir. |
| 5 | DF-22/37 tip keepers; **no** PID kill / auto-reload | PASS | `TestInstallCursorPrintReloadTip` + `TestInstallCursorWriteMergeBackup` PASS; `installCursorReloadTip` aliases `install.CursorReloadTip`. Grep `Kill`/`Signal`/`syscall` empty in `cmd/trace` and `internal/install`. |
| 6 | CONDITIONAL markers unchanged; library production untouched unless justified | PASS | `hasClaudeMarker` still `.claude/` directory **or** `CLAUDE.md` file under root. `TestInstallConditionalRefusesWithoutMarker` / `TestInstallConditionalWritesWithMarker` PASS. `projectRoot(opts)` empty → `"."` unused from CLI (always Abs). No library algorithm change required; CLI drop-root was the bug. |
| 7 | `install` still ungated; no new MCP install tools; no 015 | PASS | Grep `failCLIDenied`/`assertCLICommand`/`AssertToolAllowed` absent from `install.go` (present on gated cmds only). `RegisteredToolNames` still nine (`trace_why`…`trace_version`); no install MCP. No `015_*.sql`. Compat EmbedExpected **14**. |
| 8 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01/S03 keepers | PASS | All `01` locked cmds re-run this session (see below). `TestMCPVirginProjectDoesNotMkdir` + `TestCLIAddDeniedFailClosed` + `TestToolNamesRegistered` (nine) + `TestOpenCreates*` (v**14**) in named `-run`. |

### Subcommand × ProjectRoot (live vs FINAL)

| CLI | Sets `InstallOpts.ProjectRoot`? | Live |
|-----|--------------------------------|------|
| `install detect` | **Yes** | Abs `projectRoot` into `ListTargets` |
| `install claude` | **Yes** | Abs into `tgt.Install` |
| `install uninstall <target>` | **Yes** (harmless for cursor) | Abs into `tgt.Uninstall` |
| `install cursor` | **No** | STABLE home / `--mcp-json` only |

## Locked verify (re-run 2026-08-17)

```text
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/install/... ./internal/mcp/... ./internal/store/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallConditional|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
→ ok install, mcp, store

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallClaudeDashC|TestInstallDetectDashC|TestInstallUninstallClaudeDashC|TestInstallCursorPrintReloadTip|TestInstallCursorWriteMergeBackup|TestCLIAddDeniedFailClosed'
→ ok cmd/trace
  PASS: TestInstallClaudeDashCRefuseCitesProjectRoot
  PASS: TestInstallClaudeDashCIgnoresCwdMarker
  PASS: TestInstallClaudeDashCWriteUsesProjectRoot
  PASS: TestInstallDetectDashCClaudeReasonCitesRoot
  PASS: TestInstallDetectDashCCursorHomeUnchanged
  PASS: TestInstallUninstallClaudeDashCUsesProjectRoot
  PASS: TestInstallCursorPrintReloadTip
  PASS: TestInstallCursorWriteMergeBackup
  PASS: TestCLIAddDeniedFailClosed

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 14, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

Named CLI `-run` re-listed with `-v` this session (all nine PASS). Product suite used locked `GOMODCACHE`+`GOPROXY=off` (sandbox 403 class from S03).

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | `TestInstallDetectDashCCursorHomeUnchanged` is negative-only (cursor reason must not contain projectDir). Does not require `.cursor` / `mcp.json` in the reason. Empty reason would still pass. Live `cursor.Detect` always emits a path-bearing reason and ignores `ProjectRoot`. Do not spawn. |
| nit | `TestInstallClaudeDashCIgnoresCwdMarker` asserts `Contains(absProject)` rather than `under `+Abs(projectDir) (the refuse test is stricter). Exit-fail + cwd-cite ban still catch the hunt class. Three `ProjectRoot:` assignments vs optional helper — planner allowed either. |

## Find → refute (not reported as open)

| Proposed | Refute |
|----------|--------|
| `cmdInstall` still drops `root` | Signature keeps `root`; `resolveRoot` Abs; detect/claude/uninstall take Abs. `run()` passes `cmdInstall(root, cmdArgs)`. |
| Detect still `ListTargets(InstallOpts{})` | Live `InstallOpts{ProjectRoot: projectRoot}`. Empty opts would Chdir-detect cwd (no marker) → `detected: false`; named detect test requires `true`. |
| Claude refuse/write still uses process cwd | Named refuse cites Abs(projectDir) not cwdDir; ignore-cwd-marker still exit-fail when only cwd has `.claude/`; write Stat is under projectDir. |
| Uninstall `-C` removes cwd decoy | Named uninstall leaves cwd decoy bytes intact; project file `IsNotExist`. |
| Cursor `-C` relocates mcp.json | `cmdInstallCursor` never receives Abs root; `cursorMCPPath` is MCPJSON/HomeDir only; detect reason must not contain projectDir. |
| `hasClaudeMarker` changed to make tests pass | Same `.claude/` dir **or** `CLAUDE.md` file. CONDITIONAL keepers green. |
| PID kill / auto-reload added | No `Kill`/`Signal` in install/CLI; tip keepers still print `CursorReloadTip`. |
| `install` Asserted / `cli:install` / new MCP install tool / mig 015 | Grep-clean Assert in `install.go`; nine tools; no `015_*.sql`; compat **14**. |
| Missing `install cursor --write -C` write-under-project named test | Locked evidence is grep + detect named test. Write path cannot use `ProjectRoot` because the handler never sets it. Residual coverage only. |
| No DashC test for `CLAUDE.md` (only `.claude/`) | Marker rules unchanged in library; CONDITIONAL keepers cover the pair. CLI `-C` threading is path-agnostic. |
| Detect without `-C` still `under .` | CLI always Abs via `parseGlobal` cwd + `resolveRoot`. Intended reason-string fix (FINAL). Empty-opts path is what the `-C` detect test would fail. |
| `resolveRoot` before cursor means bogus `-C` fails cursor install | `filepath.Abs` does not Stat; nonexistent `-C` still Abs-succeeds. Cursor then ignores the value. |
| Relative `-C` Abs vs Chdir in tests | Tests pass `t.TempDir()` (absolute). Product Abs matches other commands (`parseGlobal` then `resolveRoot`). |

## Residuals for S05 / S06

- **S05:** Do **not** reopen `cmdInstall` / Cursor STABLE home / CONDITIONAL markers. Do **not** add install MCP or `cli:install`. Do **not** thread `InstallOpts.ProjectRoot` into seed/impact/`trace_impact`. DF-68 is closed here.
- **S06:** Import named DashC refuse/ignore-cwd/write/detect/cursor-home/uninstall + `TestInstallCursor*` + library `TestInstallDetectListsCursorStable` / `TestInstallConditional*` + S01/S03 keepers (already on S06 SCOPE-TODOS). Claim DF-68 only when those pass. Runnable `-run` lines copied into S06 `01-verify.md`. DF-22/37 remain residual **ops** (manual reload) — tip keepers only.
- Do **not** fail later rows for R2 defer / R3–R4 wontfix / CGO0 analyzer compile / DF-67.

## Board

- P16-S04-02 → `done`
- Next runnable → **P16-S05-00**
- No `02a`/`02b` spawn
