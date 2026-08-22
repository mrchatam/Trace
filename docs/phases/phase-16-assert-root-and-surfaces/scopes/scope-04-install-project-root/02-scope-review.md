# P16 / S04 / 02 — Scope review (install `-C`) FINAL checklist

## Metadata
- id: P16-S04-02
- todo_ids: [P16-S04-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S04 **DF-68**. Fresh subagent ≠ implementer. Compare claims + **00-PLANNER FINAL** + `01` to live code/tests. Spawn `P16-S04-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S05-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-install-project-root.md](01-install-project-root.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `cmd/trace/install.go`, `cmd/trace/root.go` `resolveRoot`, `internal/install/claude.go`, `internal/install/cursor.go`
- Dogfood: [`_post_p15/install_C_cwd.json`](../../../../../experiments/_post_p15/install_C_cwd.json) (historical; tests are SoT)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open Cursor STABLE home via `-C`, `cli:install`, CONDITIONAL marker rules, or PID-kill reload (FINAL).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | `cmdInstall` no longer drops `root`; detect/claude/uninstall get Abs `ProjectRoot` | Grep `cmdInstall(`; `resolveRoot`; `ListTargets` not `InstallOpts{}` |
| 2 | `-C` refuse/detect reasons cite `<root>` not process cwd | `TestInstallClaudeDashCRefuseCitesProjectRoot`; `TestInstallDetectDashCClaudeReasonCitesRoot` |
| 3 | Cwd marker does not authorize `-C` other root; write/uninstall land under `-C` | `TestInstallClaudeDashCIgnoresCwdMarker`; `TestInstallClaudeDashCWriteUsesProjectRoot`; `TestInstallUninstallClaudeDashCUsesProjectRoot` |
| 4 | Cursor STABLE home detect unchanged (`-C` not mcp.json) | `TestInstallDetectDashCCursorHomeUnchanged`; grep cursor Install still omits ProjectRoot-as-path |
| 5 | DF-22/37 tip keepers; **no** PID kill / auto-reload | `TestInstallCursorPrintReloadTip`; `TestInstallCursorWriteMergeBackup`; grep no `Kill`/`Signal` in install |
| 6 | CONDITIONAL markers unchanged; library production untouched unless justified | `TestInstallConditional*`; diff `hasClaudeMarker` |
| 7 | `install` still ungated; no new MCP install tools; no 015 | Grep `install.go` for Assert helpers; `TestToolNamesRegistered` (nine); compat **14** |
| 8 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01/S03 keepers | Re-run `01` locked verify cmds |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P16-S04-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P16-S05-00** unless spawn.
5. If APPROVE: S06 must import named DashC + tip keepers (already pointed on S06 SCOPE-TODOS).

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P16-S05-00** (unless spawn)
- [ ] No rewrite of done P16-S04-00/01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
