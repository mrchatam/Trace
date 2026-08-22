# P22-S02-02 — Review: git-hook

## Metadata
- id: P22-S02-02
- todo_ids: [P22-S02-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents, security-and-hardening]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C25** path exists and DF-86 / D-22-16 (no wrap `git commit`, no daemon) holds.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

| # | Check |
|---|-------|
| 1 | `trace install detect` JSON includes **`git-hook`**, tier **CONDITIONAL** |
| 2 | Hooks dir resolved via **`git rev-parse --git-path hooks`** (grep installer; read test `TestInstallGitHookHonorsCoreHooksPath`) |
| 3 | post-commit fragment between `# begin-trace` / `# end-trace`; calls **`trace -C … index`** |
| 4 | Optional **`trace seed export -o trace/graph.json`** present; failures non-fatal |
| 5 | Uninstall removes fragment only — user hook lines preserved |
| 6 | **No** `exec.Command("git", "commit"` / no commit wrapper / no daemon / no HTTP |
| 7 | Schema max still **022**; no 023+ landed in S02-01 |

```bash
go test ./internal/install/... -count=1 -run 'TestInstallGitHook|TestUninstallGitHook|TestInstallDetectIncludesGitHook'
grep -R 'git commit' internal/install/githook.go cmd/trace/install.go || true
grep -R 'ListenAndServe\|http\.Handle' internal/install/ || true
```

## Spawn policy

If C25 unmet or wrap-commit: spawn **`P22-S02-02a` + `P22-S02-02b`**. Do not close with residuals.

## Exit criteria

- [ ] C25 closed or spawned
- [ ] Confidence **high**
- [ ] Board Notes → **Next `P22-S02-03`**
