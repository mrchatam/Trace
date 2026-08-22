# P16-S03-00 — Install `-C` vs cwd (stub — thicken vs live)

## Metadata
- id: P16-S03-00
- todo_ids: [P16-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** defaults for **DF-68**: `trace -C <root> install …` must resolve Claude CONDITIONAL markers against `-C`, not process cwd. Preserve **DF-22/37** reload-tip keepers. **No product Go.**

## Live gap (P16-00)
`cmdInstall` is `func cmdInstall(_ string, args []string)` — **discards CLI `root`**. `cmdInstallDetect` uses `InstallOpts{}`; `cmdInstallClaude` sets `ProjectRoot: cwd`. Library `projectRoot` already honors `opts.ProjectRoot`.

## Inherited locks
- Pass CLI `-C` root into `InstallOpts.ProjectRoot` for **detect**, **claude install**, **uninstall**
- Cursor STABLE home detect unchanged
- Keep `TestInstallCursor*` reload tip (DF-22/37 residual — **no** PID kill)
- Named: `TestInstallClaudeHonorsDashC`; `TestInstallDetectClaudeReasonUsesRoot`

## Planner work
1. [ ] Confirm CLI root plumbing vs library
2. [ ] Thicken 01/02; **FINAL**; next **P16-S03-01**
