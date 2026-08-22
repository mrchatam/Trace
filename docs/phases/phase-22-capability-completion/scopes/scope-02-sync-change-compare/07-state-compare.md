# P22-S02-07 — Implement: compare project states

## Metadata
- id: P22-S02-07
- todo_ids: [P22-S02-07]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Allow **comparison between project states** (**C06** / **D-22-10**). Library + CLI. No blobs.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Live baseline

| Present | Absent |
|---------|--------|
| `vcs.Repository` — `Changes(ctx, oid)`, `CommitsBetween` | `CompareStates` |
| `gitcli` — `diffTreePaths`, `parseNameStatus` (`history.go`) | `trace changes compare` |
| `changes.git_commit` indexed | domain compare helper |
| Schema **23** | **No 024+** |

## Locked defaults

| Item | Value |
|------|-------|
| API | **`func (s *Service) CompareStates(ctx, fromOID, toOID string) (StateCompareResult, error)`** in **`internal/domain/compare.go`** (new) |
| Result shape | `{ from, to, added[], removed[], modified[], change_ids[] }` — paths repo-relative; status letters from git name-status |
| Diff source | **`git diff --name-status from..to`** via `vcs.Repository` (add **`DiffNameStatus(from, to)`** on gitcli if needed — **do not** persist patches) |
| Change linkage | For each OID in range (or endpoints), attach `changes.id` when `GetChangeByGitCommit` hits |
| Validation | Invalid OID → fail closed (`TestCompareStatesUnknownOIDFailClosed`) |
| CLI | **`trace changes compare --from <oid> --to <oid>`** — JSON to stdout |
| Blobs | **Forbidden** — Law 1; no `ShowFile` in compare path unless tests need fixture setup |

## Requirements

1. Implement compare in domain + tests with temp git repo (two commits, known path delta).
2. Extend `internal/gitcli` or `internal/vcs` with name-status range diff if not already available.
3. Wire CLI in `cmd/trace/changes.go` (subcommand **`compare`**).
4. CLI test `TestChangesCompare` in `cmd/trace/changes_test.go` or `cli_test.go`.

## Touch files

- `internal/domain/compare.go`, `compare_test.go`
- `internal/gitcli/history.go` or new `diff.go` (range name-status)
- `internal/vcs/repository.go` (interface method if required)
- `cmd/trace/changes.go`, `changes_test.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestCompareStatesPathDeltaNoBlob` | two commits → correct added/removed/modified |
| `TestCompareStatesLinksChangeWhenPresent` | promoted change id appears in result |
| `TestCompareStatesUnknownOIDFailClosed` | bad OID errors |

```bash
go test ./internal/domain/... ./internal/gitcli/... -count=1 -run 'TestCompareStates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestChangesCompare'
```

## Exit criteria

- [ ] C06 true
- [ ] Named tests PASS; compat still **23**
- [ ] Board Notes → **Next `P22-S02-08`**

## Minimal todos

- [ ] Compare API + git name-status
- [ ] CLI compare
- [ ] Tests
- [ ] Board notes
