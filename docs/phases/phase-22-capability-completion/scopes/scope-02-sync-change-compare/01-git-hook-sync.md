# P22-S02-01 — Implement: `trace install git-hook`

## Metadata
- id: P22-S02-01
- todo_ids: [P22-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Installable **local git hook** so Trace can update as the project changes (**C25**). Promotes **DF-86** / **D-22-02**. Does **not** wrap `git commit` (**D-22-16**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `internal/install/registry.go` — cursor+claude | `git-hook` target, `internal/install/githook.go` |
| `cmd/trace/install.go` — cursor/claude/detect/uninstall | `case install.TargetGitHook` |
| Pattern: marker-gated uninstall in `cursor.go` / `claude.go` | post-commit hook file |
| Schema **022**; compat **22** | **No 023+ this row** |

## Locked defaults

| Item | Value |
|------|-------|
| Target id | **`git-hook`**, tier **CONDITIONAL** (`install.TierConditional`) |
| Constant | `install.TargetGitHook = "git-hook"` in `types.go` |
| CLI | `trace install git-hook [--write]`; `trace install uninstall git-hook` |
| Hooks dir | **`git -C <root> rev-parse --git-path hooks`** → absolute path (honors **`core.hooksPath`**, worktrees) |
| Hook files | **`post-commit`** required; **`pre-push`** optional (same fragment) |
| Markers | `# begin-trace` / `# end-trace` (shell comments) |
| Hook body (between markers) | See **Hook script** below |
| Detect | `git rev-parse --is-inside-work-tree` under `ProjectRoot` |
| Print mode | `--write` false → print hook script to `Out` (JSON or raw shell — match claude print pattern) |
| Help | `cmd/trace/help.go` + usage strings list `git-hook` |
| G19 | Logic in `internal/install/githook.go`; CLI thin |

### Hook script (locked)

```sh
# begin-trace
ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
command -v trace >/dev/null 2>&1 || exit 0
PATHS="$(git diff-tree --no-commit-id --name-only -r HEAD 2>/dev/null || true)"
if [ -n "$PATHS" ]; then
  trace -C "$ROOT" index $PATHS || true
else
  trace -C "$ROOT" index || true
fi
trace -C "$ROOT" seed export -o trace/graph.json 2>/dev/null || true
# end-trace
```

- **`|| true`** on trace invocations — hook must never block git.
- **Never** invoke `git commit`, replace git binary, or `git add`.
- pre-push variant: same block (export optional); may noop when no new commit — acceptable.

## Requirements

1. Register `gitHookTarget{}` in `internal/install/registry.go` (third entry; stable sort by id).
2. `Install` with `--write`: mkdir hooks dir, append or replace marked fragment in `post-commit` (+ optional `pre-push`).
3. `Uninstall`: strip marked fragment only; leave other hook lines; chmod +x preserved.
4. `cmdInstall` switch + usage strings in `cmd/trace/install.go`.
5. Tests in temp git repo (`git init`); set `core.hooksPath` in at least one test.
6. CONTRIBUTING: hook optional; manual `trace index` + `trace seed export` remain valid (DF-86 backup).

## Touch files

- `internal/install/githook.go`, `githook_test.go`, `registry.go`, `types.go`
- `cmd/trace/install.go`, `help.go`
- `internal/install/install_test.go` or dedicated githook tests
- `CONTRIBUTING.md` (portable graph / continuous update paragraph)

## Named tests

| Test | Proves |
|------|--------|
| `TestInstallGitHookWritesPostCommit` | resolved hooks dir contains executable post-commit with `trace -C` + `index` |
| `TestInstallGitHookHonorsCoreHooksPath` | `git config core.hooksPath .husky` → fragment under `.husky/post-commit`, not `.git/hooks` |
| `TestInstallGitHookDoesNotWrapCommit` | installer + hook body contain no `git commit` / no git binary replacement |
| `TestUninstallGitHookRemovesFragment` | markers gone; sibling hook lines kept |
| `TestInstallDetectIncludesGitHook` | `ListTargets` JSON includes id `git-hook`, tier CONDITIONAL |

```bash
go test ./internal/install/... -count=1 -run 'TestInstallGitHook|TestUninstallGitHook|TestInstallDetectIncludesGitHook'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelp|TestInstall'
```

## Exit criteria

- [ ] C25 has installable continuous-update path (hook) + docs for manual `trace index`
- [ ] DF-86: no wrap commit; hooks path via rev-parse
- [ ] Named tests PASS
- [ ] Schema still max **022**; compat still **22**
- [ ] Board Notes → **Next `P22-S02-02`**

## Minimal todos

- [ ] Target + hook write/uninstall + hooksPath test
- [ ] CLI + help
- [ ] Tests + CONTRIBUTING
- [ ] Board notes
