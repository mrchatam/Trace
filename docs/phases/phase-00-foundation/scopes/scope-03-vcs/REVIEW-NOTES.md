# P00-S03-02 — Scope review notes (2026-08-15)

Independent review of S03 against `01-vcs.md` + TODO Notes for `P00-S03-01`. Fresh session; claims verified in-repo.

## Verdict

**APPROVE** — one **high** finding fixed inline (History index freshness). No remaining blocker/high. Confidence: **high**. Spawns: **none**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `internal/vcs` `Repository` covers Head, IsRepo, ShowFile, History, CommitsBetween, Changes, LastChanged, Refresh (+ Close) | Pass (`repository.go`, `types.go`) |
| ≥1 fake/mock of interface used in unit test | Pass (`fake.go`, `fake_test.go`) |
| `internal/gitcli` implements via git CLI only (`exec.CommandContext`, `-C` root) | Pass (`exec.go`, `open.go`, …) |
| No go-git / libgit2 in `go.mod` / imports | Pass (`go.mod` = uuid + modernc/sqlite only; grep clean) |
| Temp-repo: `ShowFile` bytes == `git show` | Pass (`TestShowFileMatchesGitShow`) |
| Temp-repo: History / CommitsBetween / Changes / LastChanged multi-commit | Pass (`TestHistoryCommitsBetweenChangesLastChanged`) |
| Incremental Refresh: N → noop → +K → noop; no wipe | Pass (`TestIncrementalRefresh`; watermark + `CountIndexedCommits`) |
| Thin index schema; no source blobs / full diffs in SQLite | Pass (`002_vcs_index.sql`; `TestNoBlobColumnsInVCSIndex` / `HasBlobLikeColumns`) |
| Watermark durable in `vcs_meta` (`vcs_index_watermark`) | Pass (`refresh.go` + store helpers) |
| `CGO_ENABLED=0 go test ./internal/vcs/... ./internal/gitcli/... ./...` | Pass (2026-08-15 re-run after fix) |
| No MCP/daemon/HTTP; no new `cmd/trace` VCS wiring | Pass |
| Cross-scope: S04 can consume `vcs.Repository` path/OID; `files.git_oid` optional | Pass (iface + store `UpsertFile`; S04 Depends note thickened lightly) |

## Findings

### High (fixed inline — no spawn)

- **Stale index preferred for History/LastChanged** (`history.go`): any non-empty `ListIndexedHistory` short-circuited git, so after new commits **without** `Refresh`, History/LastChanged returned outdated tips (and a mid-Refresh partial index could truncate history).
  - **Fix:** prefer index only when durable watermark == current `HEAD`.
  - **Regression:** `TestHistoryFallsBackWhenWatermarkBehindHEAD`.

### Medium / low (no spawn)

- **Branch rewind / force-push:** watermark advances with HEAD; commits no longer on the tip remain in the thin index (by design — no full wipe). Callers should treat index as tip-accelerated, not a full reflog. Residual / acceptable for P0.
- **Rename status:** `parseNameStatus` keeps status letter + new path only; history of the old path may miss rename commits via index. Low for P0.
- **`IsRepo`:** broad `fatal:` / “not a git repository” string matching returns `false` rather than error for some git failures — acceptable Open gate; low.

### Nit

- `UpsertIndexedCommit` ON CONFLICT does not update `seq` (preserves first-index order) — correct for re-upsert.
- `Fake.CommitsBetween` is linear-list semantics, not git rev-list — fine for iface tests only.

### Blocker

None.

## Spawns

None.

## Residual risks

- Large-repo Refresh cost (N commits × `show` + `diff-tree`) unmeasured — API honest; tuning out of scope.
- Concurrent `store.Open` from `gitcli.Open` plus a second store handle (as in tests) relies on SQLite multi-connection behavior — fine for local single-writer P0.
- Index-accelerated `Changes` still uses per-commit rows when present even if watermark lags; safe because each upserted commit’s path list is complete, and unknown OIDs fall back to git.

## Forward edits this review

- `internal/gitcli/history.go` — watermark==HEAD gate before index History
- `internal/gitcli/gitcli_test.go` — `TestHistoryFallsBackWhenWatermarkBehindHEAD`
- `scopes/scope-04-analyzers/01-analyzers.md` — light Depends note: consume `vcs.Repository` (+ `vcs.Fake` in tests); `gitcli.Open` at wiring only; optional `git_oid` on `UpsertFile`
- `SCOPE-TODOS.md` — mark S03-01/02 done
