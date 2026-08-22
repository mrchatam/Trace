# P16-S04-00 — install `-C` vs cwd (FINAL)

## Metadata
- id: P16-S04-00
- todo_ids: [P16-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `cmd/trace/install.go` + `internal/install`, lock **FINAL** defaults for **DF-68**: CONDITIONAL install/detect honor global `-C`/`--project`. Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition DF-68 **fix**; DF-22/37 **carry-forward**
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — S04: thread `cmdInstall(root)` into `InstallOpts.ProjectRoot`
- Live: `cmd/trace/install.go` — `cmdInstall(_ string, …)` **drops** `root`; `cmdInstallDetect` `ListTargets(InstallOpts{})`; `cmdInstallClaude`/`cmdInstallUninstall` `ProjectRoot: cwd`; `internal/install/claude.go` already `projectRoot(opts)`
- Dogfood: [`POST-P15-DOGFOOD.md`](../../../../../experiments/POST-P15-DOGFOOD.md); [`_post_p15/install_C_cwd.json`](../../../../../experiments/_post_p15/install_C_cwd.json)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-68, DF-22, DF-37
- Quality bar: [P16-S03-00](../scope-03-cli-mcp-allowlist-parity/00-PLANNER.md) FINAL
- S03 live: **P16-S03-02 APPROVE** — `install` ungated (no `failCLIDenied` / `assertCLICommand` / `AssertToolAllowed` in `install.go`)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Depends-on: **P16-S03-02 APPROVE** (board, already true). Phase locks below. **Unattended:** no architecture blockers; defaults below are FINAL. Do **not** add allowlist / `cli:install` work. **No SwitchMode** (orchestrator). **No product Go.**

## Depends (from S03 — live P16-S03-02 APPROVE)

S03 FINAL left `install` **ungated**. Live after S03:

- `cmd/trace/install.go` has **zero** `failCLIDenied` / `assertCLICommand` / `AssertToolAllowed`.
- `cmdInstall` still **drops** `root` (`func cmdInstall(_ string, args []string)`).
- DF-77 dual-slug does **not** Assert `install detect|cursor|claude|uninstall`. Operator escape + this scope’s `-C` stay independent of `cli:add`.
- Do **not** canonicalize install subcommands into `cli:` slugs here. Do **not** add mig 015.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| Global `-C` | Yes — `parseGlobal` → `cmdInstall(root, cmdArgs)` | Root **dropped** at CLI |
| Library `InstallOpts.ProjectRoot` | Yes — `claude.go` `projectRoot(opts)`; empty → `"."` | CLI never passes parsed `-C` |
| `install detect` | Yes — `ListTargets(InstallOpts{})` | Empty opts → claude reason `under .` (dogfood) even with `-C` |
| `install claude` | Yes — `ProjectRoot: os.Getwd()` | Dogfood `/tmp` + `-C scratch` cited **`/tmp`** |
| `install uninstall` | Yes — `ProjectRoot: cwd` for all targets | Claude uninstall is project-relative; same cwd bug |
| `install cursor` | Yes — `MCPJSON` / `$HOME/.cursor/mcp.json` | **No** ProjectRoot (correct STABLE home) |
| Cursor Detect | Yes — `HomeDir` / `cursorMCPPath` | Unchanged by `-C` (dogfood HOLD) |
| CONDITIONAL markers | Yes — `.claude/` dir **or** `CLAUDE.md` file | Do **not** change |
| DF-22/37 tip | Yes — `CursorReloadTip` print + `--write` | Keepers; **no** PID kill |
| CLI Assert on install | **Absent** | Keep ungated |
| MCP install tool | **Absent** | Keep absent |

**Bug path DF-68 (live):** `trace -C <root> install claude` from a different cwd resolves `hasClaudeMarker` against **process cwd**. `install detect` claude `reason` says `under .` because `ListTargets` gets empty `ProjectRoot`. Cursor STABLE home detect is unaffected.

### Live install subcommands (`cmd/trace/install.go`)

| Token | Handler | Live ProjectRoot | After this scope |
|-------|---------|------------------|------------------|
| `detect` | `cmdInstallDetect` | **none** (`InstallOpts{}`) | **Yes** — Abs parsed `-C` / cwd |
| `claude` | `cmdInstallClaude` | `os.Getwd()` | **Yes** — Abs parsed `-C` / cwd |
| `uninstall <target>` | `cmdInstallUninstall` | `os.Getwd()` | **Yes** — Abs parsed `-C` / cwd (claude uses it; cursor ignores) |
| `cursor` | `cmdInstallCursor` | **unset** (home / `--mcp-json`) | **No** — STABLE `$HOME/.cursor/mcp.json`; do not treat `-C` as mcp.json |
| unknown | usage | — | unchanged |

`parseGlobal` already defaults omitted `-C` to `os.Getwd()`. Other commands `resolveRoot` (`filepath.Abs`) before use. Install must **Abs once** in `cmdInstall` (or each ProjectRoot call site) so refuse/detect reasons cite the **absolute** `-C` path, not a relative string that `Stat`s against cwd.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Home | Thin `cmd/trace/install.go` (+ `install_test.go`). Library `internal/install` **already** root-aware — **no** marker-rule / Detect algorithm change |
| Thread | `cmdInstall(root, …)` **keeps** `root`. `resolveRoot(root)` → `InstallOpts.ProjectRoot` for **detect**, **claude**, **uninstall** |
| Cursor install | **Do not** set ProjectRoot as the mcp.json location. STABLE detect/install/uninstall stay `HomeDir` / `--mcp-json` / `$HOME/.cursor/mcp.json` |
| Detect JSON | Claude `reason` contains the Abs project root (not `under .` when `-C` points elsewhere; not process cwd). Cursor `reason` still home/`mcp.json` path |
| No `-C` | Parsed root is cwd (existing `parseGlobal`). Detect then cites **Abs(cwd)** instead of `"."` — intended reason-string fix, not a marker-rule change |
| Markers | Unchanged: `.claude/` directory **or** `CLAUDE.md` file under ProjectRoot |
| DF-22 / DF-37 | Keep `TestInstallCursor*` print/write reload tip (`CursorReloadTip`); **no** PID kill / auto-reload |
| Ungated | `install` stays **ungated** — no `cli:install`, no Assert |
| G19 | CLI adapter only; CONDITIONAL logic stays in `internal/install` |
| Compat | Ceiling **14**; **no** mig 015+ |
| Forbidden | New MCP install tools; changing CONDITIONAL marker rules; daemon/HTTP; `cli:install`; rewriting Phase 00–15 `done` history; using `-C` to relocate Cursor home mcp.json |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat **14**; p0x; x0; product pkgs `./cmd\|internal\|evals`. S01 virgin + S02 CHECK/canonicalize + S03 dual-slug keepers stay green |

### Which install subcommands take ProjectRoot (locked)

| CLI | Sets `InstallOpts.ProjectRoot`? | Why |
|-----|--------------------------------|-----|
| `trace [-C root] install detect` | **Yes** | Claude CONDITIONAL marker + reason cite `<root>`. Cursor Detect **ignores** the field (HomeDir/MCPJSON) |
| `trace [-C root] install claude [--write] [--bin]` | **Yes** | Marker check + write `<root>/.claude/trace-mcp.json` |
| `trace [-C root] install uninstall claude` | **Yes** | Remove `<root>/.claude/trace-mcp.json` |
| `trace [-C root] install uninstall cursor` | **Yes (harmless)** | Cursor Uninstall ignores ProjectRoot; passing Abs root is OK and keeps one uninstall opts builder |
| `trace [-C root] install cursor [--write] …` | **No** | STABLE home / `--mcp-json` only. `-C` must **not** change detect or write path |

Implementer may share a small opts helper for detect/claude/uninstall. Do **not** invent a second `--project` flag on install subcommands — global `-C` only.

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| `TestInstallClaudeDashCRefuseCitesProjectRoot` | `cmd/trace` | Dogfood class: `Chdir(cwdDir)` (no marker); `-C projectDir` (no marker); `install claude` exit fail; stderr cites **Abs(projectDir)** after `under `; must **not** cite cwdDir as the marker root |
| `TestInstallClaudeDashCIgnoresCwdMarker` | `cmd/trace` | cwdDir **has** `.claude/`; projectDir has **no** marker; `-C projectDir install claude` still refuse citing **projectDir** (cwd marker must not authorize) |
| `TestInstallClaudeDashCWriteUsesProjectRoot` | `cmd/trace` | cwdDir no marker; projectDir **has** `.claude/`; `-C projectDir install claude --write` writes `projectDir/.claude/trace-mcp.json`; cwdDir has **no** write |
| `TestInstallDetectDashCClaudeReasonCitesRoot` | `cmd/trace` | Same split dirs; `-C projectDir install detect` JSON: claude `detected: true` (marker on project); `reason` contains Abs(projectDir); **not** `under .` as the root; **not** cwdDir |
| `TestInstallDetectDashCCursorHomeUnchanged` | `cmd/trace` | Same `-C projectDir` detect: cursor `reason` does **not** contain projectDir (STABLE home / mcp.json path unchanged) |
| `TestInstallUninstallClaudeDashCUsesProjectRoot` | `cmd/trace` | After write under `-C projectDir`, `install uninstall claude` with same `-C` removes that file; a decoy under cwdDir is untouched |
| `TestInstallCursorPrintReloadTip` | `cmd/trace` | **Keeper** DF-22/37 print tip |
| `TestInstallCursorWriteMergeBackup` | `cmd/trace` | **Keeper** DF-22/37 write tip + backup |
| `TestInstallDetectListsCursorStable` | `internal/install` | **Keeper** STABLE home detect |
| `TestInstallConditionalRefusesWithoutMarker` | `internal/install` | **Keeper** CONDITIONAL refuse |
| `TestInstallConditionalWritesWithMarker` | `internal/install` | **Keeper** CONDITIONAL write |
| `TestCLIAddDeniedFailClosed` | `cmd/trace` | **S03 keeper** — dual-slug still green |
| `TestToolNamesRegistered` | `internal/mcp` | **Keeper** — still nine until S05 |
| `TestMCPVirginProjectDoesNotMkdir` | `internal/mcp` | **S01 keeper** |
| `TestOpenCreatesDBAndMigratesIdempotent` | `internal/store` | **Keeper** — version **14** |

TDD: named CLI `-C` vs cwd tests first (red: `cmdInstall` drops root / detect empty opts), then thread `resolveRoot` → `ProjectRoot` (green). Do **not** change `hasClaudeMarker` to make tests pass.

Library production Go is **not** required unless a test proves a library bug (none live). Optional extra `internal/install` test that `ListTargets` with `ProjectRoot` cites that path is allowed, not a substitute for the CLI named tests.

## Owns

| Item | Intent |
|------|--------|
| DF-68 | `install detect` / `install claude` / `install uninstall` (when project-relative) honor Abs parsed `-C` |
| DF-22 / DF-37 | Tip keepers only — print/write reload tip; no PID kill |
| Subcommand table | Which handlers set ProjectRoot vs Cursor STABLE home |

## Explicit deferrals

- Relocating Cursor mcp.json via `-C` — **rejected** (STABLE home)
- PID-kill / auto-reload of live Cursor MCP (DF-22/37 residual ops)
- `cli:install` / Assert on `cmdInstall` (S03 ungated)
- Changing CONDITIONAL markers (`.claude/` \| `CLAUDE.md`)
- New MCP install / decide tools; daemon/HTTP
- S05 seed/impact/`trace_impact`; unprefixed command slugs
- R2 `allowContainsOut`; R3 graphify space; R4 CGO0 analyzers
- Session-global DENY across `project=` roots

## Assumptions (unattended)

1. **Library is correct:** DF-68 is the CLI dropping `root` / empty detect opts. Do not reopen marker rules.
2. **Abs at CLI:** refuse/detect reasons must be comparable to dogfood (`under /tmp` vs `under <abs -C>`). `resolveRoot` matches every other command.
3. **Detect without `-C`:** citing Abs(cwd) instead of `"."` is the same root, clearer string — in scope.
4. **One uninstall opts builder** may always set ProjectRoot; cursor ignores it.
5. **chdir in tests:** required to prove cwd ≠ `-C` (dogfood from `/tmp`). Follow `cli_test.go` Chdir+Cleanup.
6. **CGO:** named `cmd/trace` tests run **CGO1** (R4: CGO0 `./cmd/trace/...` tree-sitter). `internal/install` named tests stay CGO0.
7. S06 VERIFY imports the DF-68 CLI named tests + Cursor tip keepers + S01–S03 keepers + compat **14**.

## Effects on later scopes

- **S05:** Do **not** thread install ProjectRoot into seed/impact/MCP. Do **not** add install MCP. DF-68 is this scope; S05 Depends is board-order only.
- **S06:** Import named tests in the table above (DashC refuse/ignore-cwd/write/detect/cursor-home/uninstall + `TestInstallCursor*` + CONDITIONAL library keepers). Claim DF-68 only when those pass. DF-22/37 remain residual **ops** (manual reload) — tip keepers only.

## Planner work

1. [x] Inventory live `cmdInstall` drop + detect empty opts + dogfood `install_C_cwd.json`
2. [x] Lock which install subcommands take ProjectRoot; Cursor STABLE unchanged
3. [x] Lock named CLI `-C` vs cwd tests + DF-22/37 keepers
4. [x] Thicken `01-install-project-root.md` + `02-scope-review.md` + SCOPE-TODOS
5. [x] Light S05 Depends (no install MCP / no ProjectRoot into seed); S06 named-test pointer
6. [x] Board Notes → next **P16-S04-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL** with subcommand × ProjectRoot table
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P16-S04-01**

## Minimal todos
- [x] Inventory live install / DF-68 dogfood / S03 ungated keepers
- [x] FINAL locks + named tests
- [x] Thicken 01/02/SCOPE-TODOS + S05/S06 pointers
- [x] Board sync

## Next
**P16-S04-01** (implement DF-68 thread `-C` → `InstallOpts.ProjectRoot`).
