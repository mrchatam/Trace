# P22-S02-06 — Review: change capture

## Metadata
- id: P22-S02-06
- todo_ids: [P22-S02-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C05** is not leftover as “apply-only”.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

| # | Check |
|---|-------|
| 1 | VCS promote path creates `changes` with `source_type=VCS` |
| 2 | Idempotent on git OID — grep `GetChangeByGitCommit` |
| 3 | Paths from `vcs_commit_paths` only — no blob/patch columns |
| 4 | Meaningful filter + `--all` honored |
| 5 | `loop apply` change path still works |
| 6 | No mig **024+**; compat still **23** |
| 7 | `gitcli.Refresh` noop semantics preserved (second Refresh, no new commits) |

```bash
go test ./internal/domain/... -count=1 -run 'TestPromoteVCSCommit|TestCreateChangeWithGitSHAAndPathsNoBlob'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestChangesCapture|TestLoopApplyDeliberationTransitionEvent'
go test ./internal/gitcli/... -count=1 -run 'TestRefresh' 
```

## Spawn policy

If C05 unmet: spawn **`P22-S02-06a` + `P22-S02-06b`**. Do not close with residuals.

## Exit criteria

- [ ] C05 closed or spawned
- [ ] Confidence **high**
- [ ] Board Notes → **Next `P22-S02-07`**
