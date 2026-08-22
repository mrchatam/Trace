# P22-S02-05 — Implement: meaningful change capture

## Metadata
- id: P22-S02-05
- todo_ids: [P22-S02-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Record every meaningful change** — not only `loop apply` (**C05** / **D-22-07**). Promote VCS commits into `changes` + `change_paths`. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Live baseline

| Present | Absent |
|---------|--------|
| `changes` / `change_paths` mig **017** (`internal/store/schema/017_changes_effects.sql`) | `PromoteVCSCommitToChange` |
| `vcs_commits` / `vcs_commit_paths` mig **002**; `gitcli.Refresh` populates | `GetChangeByGitCommit` |
| `domain.CreateChange` — requires ≥1 path, sets RECORDED when git_commit set | `trace changes` CLI |
| `loop apply` → `ApplyChange` path | auto-capture after Refresh |
| Schema max **23** after S02-03 | **No 024+** (S04 owns 024) |

## Locked defaults

| Item | Value |
|------|-------|
| API | **`func (s *Service) PromoteVCSCommitToChange(ctx, oid string) (store.Change, error)`** in `internal/domain/changes.go` |
| Idempotent | Same OID → return existing change; no duplicate rows (`store.GetChangeByGitCommit`) |
| Paths | From `vcs_commit_paths` for OID; copy path+status; **no blobs** |
| Fields | `git_commit=oid`, `source_type=VCS`, `task_id=trace:vcs-capture`, `reason=vcs_commits.subject`, `status=RECORDED` |
| Meaningful filter | Include commit if **any** path passes: `analyzers.DetectLanguage(path)` **or** `store.GetFileByPath(path)` succeeds |
| `--all` flag | CLI `trace changes capture --all` skips meaningful filter |
| Skip | Zero qualifying paths → no row (not an error) |
| Trigger | **`trace changes capture [--since <oid>] [--all]`**; optional call from end of `cmdIndex` after Refresh (behind small helper — must not break Refresh noop) |
| Loop apply | **Unchanged** — keeper `TestLoopApplyDeliberationTransitionEvent` |
| Schema | Reuse **023** + 017/002 — **no new mig** |

## Requirements

1. `store.GetChangeByGitCommit(oid string) (Change, error)` — index on existing `idx_changes_git_commit`.
2. `PromoteVCSCommitToChange` + unit tests in `internal/domain/changes_test.go`.
3. New `cmd/trace/changes.go`: subcommands **`capture`** (implement now); stub **`compare`** for S02-07 or implement if trivial.
4. Register in `cmd/trace/root.go` + `help.go`.
5. Optional: `PromoteRecentVCSCommits` helper for `--since` / post-index batch.

## Touch files

- `internal/domain/changes.go`, `changes_test.go`
- `internal/store/changes.go`
- `cmd/trace/changes.go`, `root.go`, `help.go`
- `cmd/trace/index.go` (optional capture hook — keep Refresh noop semantics)

## Named tests

| Test | Proves |
|------|--------|
| `TestPromoteVCSCommitCreatesChangeIdempotent` | second promote same OID → same id, one row |
| `TestPromoteVCSCommitRecordsPathsNoBlob` | paths copied; no content columns |
| `TestPromoteVCSCommitSkipsNonMeaningful` | docs-only commit skipped unless `--all` |
| `TestCreateChangeWithGitSHAAndPathsNoBlob` | existing keeper still PASS |
| `TestLoopApplyStillCreatesChange` | `TestLoopApplyDeliberationTransitionEvent` PASS |

```bash
go test ./internal/domain/... -count=1 -run 'TestPromoteVCSCommit|TestCreateChangeWithGitSHAAndPathsNoBlob'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestChangesCapture|TestLoopApplyDeliberationTransitionEvent'
```

## Exit criteria

- [ ] C05 true (VCS-promoted + apply)
- [ ] Law 1 (no blobs)
- [ ] Named tests PASS; compat still **23**
- [ ] Board Notes → **Next `P22-S02-06`**

## Minimal todos

- [ ] Store lookup + Promote API + idempotency
- [ ] CLI capture + optional index wiring
- [ ] Tests
- [ ] Board notes
