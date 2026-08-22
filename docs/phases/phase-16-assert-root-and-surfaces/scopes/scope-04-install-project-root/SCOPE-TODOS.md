# S04 — install `-C` vs cwd — scope todos

**Depends-on:** P16-S03-02 APPROVE (board). Owns **DF-68**; **DF-22/37** tip keepers (no PID kill).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** (P16-S04-00) |
| 2 | 01-install-project-root | implement | **done** (P16-S04-01) |
| 3 | 02-scope-review | review | **APPROVE high** (P16-S04-02); next **P16-S05-00** |

## Depends (from S03 — live)
- `install` is **ungated** on the DF-77 dual-slug table. S04 must not add `AssertToolAllowed` / `cli:install`.
- **P16-S03-02 APPROVE:** `cmd/trace/install.go` grep-clean of Assert helpers; `cmdInstall` still dropped `root` at plan time (01 threads it).
- S03 does not change `-C` / `ProjectRoot` (this scope owns DF-68).

## Phase locks (P16-S04-00 FINAL)
- Pass Abs `cmdInstall(root)` into `InstallOpts.ProjectRoot` for **detect**, **claude**, **uninstall**
- Cursor STABLE `$HOME/.cursor/mcp.json` detect/install **unchanged** (`-C` is not mcp.json)
- CONDITIONAL markers unchanged (`.claude/` \| `CLAUDE.md`)
- No new MCP install tools; no PID kill
- Named: `TestInstallClaudeDashCRefuseCitesProjectRoot`, `TestInstallClaudeDashCIgnoresCwdMarker`, `TestInstallClaudeDashCWriteUsesProjectRoot`, `TestInstallDetectDashCClaudeReasonCitesRoot`, `TestInstallDetectDashCCursorHomeUnchanged`, `TestInstallUninstallClaudeDashCUsesProjectRoot` + `TestInstallCursor*` keepers

## Reminders
- **P16-S04-02 APPROVE** — next **P16-S05-00**. Do not spawn 02a/02b.
