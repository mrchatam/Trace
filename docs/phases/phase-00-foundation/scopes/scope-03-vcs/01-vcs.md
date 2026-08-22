# P00 / S03 / 01 — VCS adapter + git CLI

## Metadata
- id: P00-S03-01
- todo_ids: [P00-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Define the **VCS adapter interface** in `internal/vcs`, implement it with the **`git` CLI** in `internal/gitcli`, and support a **thin commit/path index** with **incremental history refresh** (new commits only — no full history rewrite on every call). Content always comes from Git; SQLite stores references (OIDs, paths, metadata) only — **never** source blobs or permanent full diffs.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) (Git adapter + no content duplication)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G1 (Git canonical / no blobs), G12 (incremental)
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — D2, DR-GIT, DR-INCREMENTAL, DR-SURFACE, DR-RISK
- [STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §2–§5
- [ARCHITECTURE.md](../../../../ARCHITECTURE.md) — VCS adapter responsibilities
- Historical T003: [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md)
- Prior: S01 stubs `internal/vcs`, `internal/gitcli`; S02 `internal/store` Open/migrate + `UpsertFile(..., gitOID *string)` (no body columns); `go.mod` **go 1.24.0**

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (sqlite floor; do **not** downgrade) |
| Interface package | `internal/vcs` (package `vcs`) — types + `Repository` (or equivalent) interface + test fake/mock |
| Implementation package | `internal/gitcli` (package `gitcli`) — **only** git CLI backend this phase |
| Git access | `exec` subprocess to host `git` binary (DR-GIT). Prefer `CommandContext` |
| Forbidden backends | **No** libgit2, **no** go-git, **no** other embedded Git library unless a future Decision supersedes DR-GIT |
| Working tree | Bound project root that is (or contains) a Git work tree; resolve via `git -C <root> …` / `GIT_DIR` as needed |
| Content resolution | `git show <rev>:<path>` / `git cat-file` / equivalent — bytes returned to caller; **do not** write file bodies into SQLite |
| Store / blobs | Optional: call `store.UpsertFile` with path + content hash + `git_oid` for metadata only. **Never** add source BLOB/body columns |
| Thin history index | Persist **references only**: commit OID, parent OID(s), committer/author time (RFC3339 or unix), optional subject **summary** (short), and **changed path list** per commit. Do **not** store full patch/diff bodies permanently (T003 out of scope) |
| Index location | Same project DB: `<root>/.trace/trace.db` via new embedded store migration (e.g. `002_vcs_index.sql`) + thin helpers in `internal/store` **or** write helpers colocated but schema still owned by `store` embed.fs — **one SQLite per project** (DR-DB). Do not invent a second DB file |
| Incremental refresh | Keep a durable **watermark** (e.g. last indexed commit OID / HEAD tip processed). `Refresh` walks only commits not yet indexed; a second `Refresh` with no new commits must not rewrite the whole index |
| Abstraction | All library consumers (S04+, tests) depend on `vcs` interface — not `gitcli` concretions — except construction/wiring |
| Surface | Library only this scope — **no** new `cmd/trace` subcommands (S07); **no** MCP/daemon/HTTP |
| Host assumption | `git` available on PATH in tests (create temp repos with `git init` + commits) |
| Depends on | S01–S02 done |
| Repo state at start | `internal/vcs/doc.go` + `internal/gitcli/doc.go` stubs; store has `files.git_oid` but no commit-index tables yet |

### Minimum `vcs` interface surface (names may vary; behavior locked)

```text
// Construction (gitcli):
Open(repoRoot string) (vcs.Repository, error)  // or New(repoRoot) — fail clearly if not a git repo / git missing

// Identity / refs:
Head(ctx) (CommitOID string, error)
IsRepo(ctx) (bool, error)   // optional if Open already validates

// Content (Git canonical — G1):
ShowFile(ctx, rev, path) ([]byte, error)     // equals `git show rev:path` for tracked paths
// optional: ListTree / PathExists — only if needed for tests; keep minimal

// History / changes (index-accelerated after Refresh; may fall back to git):
History(ctx, path string, limit int) ([]CommitMeta, error)           // history(file)
CommitsBetween(ctx, fromRev, toRev string) ([]CommitMeta, error)     // commits_between
Changes(ctx, commitOID string) ([]PathChange, error)                 // changes(commit) — paths (+ status A/M/D if cheap)
LastChanged(ctx, path string) (CommitMeta, error)                    // last_changed(file)

// Incremental index:
Refresh(ctx) (RefreshResult, error)  // index new commits since watermark; no full rebuild when unchanged
```

`CommitMeta` at minimum: `OID`, `ParentOIDs` (or first parent), `CommittedAt`, optional short `Subject`.  
`PathChange` at minimum: `Path`, optional `Status` (`A`/`M`/`D`/`R`…).  
`RefreshResult` at minimum: counts of newly indexed commits (0 when noop).

### Suggested git CLI mapping (implementer guidance)

| Capability | Typical git invocation |
|------------|------------------------|
| Validate repo | `git rev-parse --is-inside-work-tree` |
| HEAD | `git rev-parse HEAD` |
| File at rev | `git show <rev>:<path>` |
| Commits between | `git rev-list --reverse <from>..<to>` or `git log --format=…` |
| Paths in commit | `git diff-tree --no-commit-id --name-status -r <oid>` |
| File history | `git log --format=… -- <path>` |
| Incremental tip walk | `git rev-list --reverse <watermark>..HEAD` (empty → noop) |

Parse porcelain carefully; treat non-zero exit + stderr as errors (no silent empty success on hard failures).

### Target tree

```text
internal/vcs/
  doc.go              # package comment: adapter interface
  types.go            # CommitMeta, PathChange, RefreshResult, errors
  repository.go       # Repository interface
  fake.go / mock_test.go  # ≥1 fake/mock used by unit tests outside gitcli

internal/gitcli/
  doc.go
  open.go             # Open(repoRoot) → vcs.Repository
  exec.go             # git runner helper (-C root, ctx, env)
  history.go          # History / CommitsBetween / Changes / LastChanged
  content.go          # ShowFile
  refresh.go          # incremental index refresh + watermark
  gitcli_test.go      # temp git repos (see Exit criteria)

internal/store/       # only if adding index tables (preferred)
  schema/002_*.sql    # commits + commit_paths (or equivalent) + watermark/meta
  vcs_index.go        # thin UpsertCommit / List… helpers OR watermark get/set
```

Do **not** add `go-git` / libgit2 modules. Do **not** wire CLI commands. Do **not** change S02 File/Symbol/Import semantics beyond optional OID upserts in tests.

### Out of scope (this row)

- go-git / libgit2 implementation
- Storing full diffs / patch text permanently
- Analyzer / causal / retrieval / CLI product commands
- MCP, daemon, HTTP
- Large-repo performance tuning (note unknown; keep API honest)

## Board rights
Implementer: update **status + notes only** on `P00-S03-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `internal/vcs` exports a `Repository` (or equivalent) interface covering History, CommitsBetween, Changes, LastChanged, ShowFile, Refresh (+ Head or equivalent)
- [ ] `internal/gitcli` implements that interface via **git CLI only** (no go-git/libgit2 in `go.mod`)
- [ ] Temp-repo tests: `ShowFile` bytes equal `git show` for a tracked path at a commit
- [ ] Temp-repo tests: History / CommitsBetween / Changes / LastChanged behave correctly on a multi-commit fixture
- [ ] **Incremental refresh test**: index N commits → `Refresh` → add K new commits → `Refresh` indexes only new ones (or watermark advances without rewriting prior rows en masse); third `Refresh` with no new commits is a noop (0 new / no destructive wipe)
- [ ] ≥1 fake/mock of the `vcs` interface used in a unit test (can live under `internal/vcs`)
- [ ] No source blobs / full diff bodies persisted in SQLite; schema audit or test assertion if new tables added
- [ ] `go test ./internal/vcs/... ./internal/gitcli/...` and `go test ./...` pass; `CGO_ENABLED=0` still OK
- [ ] No MCP/daemon/HTTP; no new product CLI commands
- [ ] TODO.md Notes for `P00-S03-01` updated; status `done`

## Minimal todos
- [ ] Define `vcs` types + `Repository` interface + fake/mock
- [ ] Implement `gitcli.Open` + git exec helper (`-C` root, context)
- [ ] Implement ShowFile + Head (content/ref path)
- [ ] Add thin commit/path index schema (store migration) + watermark
- [ ] Implement Refresh (incremental) + History / CommitsBetween / Changes / LastChanged
- [ ] Tests: temp git fixture; content=`git show`; incremental refresh; interface fake
- [ ] Board status + notes (commands, package paths, migration id if any)
