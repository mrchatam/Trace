# P22-S02-00 — Planner: graph sync + change capture + state compare

## Metadata
- id: P22-S02-00
- todo_ids: [P22-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S02 against live install registry, index honesty, and change tables. Owned: **C04, C05, C06, C25**. **No product Go.**

## Live inventory (2026-08-18, post-S01)

| Surface | Live state |
|---------|------------|
| Schema max | **022** (`022_code_relationships.sql`; 22 embed sql files) |
| Compat ceiling | **22** (`evals/compat/compat_test.go` — no 023+) |
| Install registry | `internal/install/registry.go` — **`cursor`** (STABLE), **`claude`** (CONDITIONAL); **no `git-hook`** |
| CLI install | `cmd/trace/install.go` — detect / uninstall / cursor / claude only; usage string omits git-hook |
| Help | `cmd/trace/help.go` — install line lists cursor+claude only |
| Index honesty (disk) | `internal/compiler/index_honesty.go` — `buildIndexHonesty` compares indexed `content_hash` vs disk sha256; wired in `compiler.go` → `Packet.IndexHonesty`; loop `contextSectionFreshness` marks **dirty** when stale |
| Graph sync watermark | **Absent** — no `graph_sync_state`; no HEAD-vs-indexed-commit banner |
| VCS index | `internal/gitcli/refresh.go` — `Refresh()` uses `vcs_meta` key **`vcs_index_watermark`** (`store.MetaVCSWatermark`); upserts `vcs_commits` + `vcs_commit_paths` |
| Index CLI | `cmd/trace/index.go` — calls `repo.Refresh` then file-local `IndexFile`; **no** post-index watermark; **no** `index status` subcommand |
| Changes domain | `internal/domain/changes.go` — `CreateChange` / `RecordChangeCommit` only; **no** `PromoteVCSCommitToChange` |
| Changes store | `internal/store/changes.go` — upsert/get/list; **no** `GetChangeByGitCommit` |
| State compare | **Absent** — no `CompareStates`; `vcs.Repository` has `Changes` / `CommitsBetween` but no two-OID diff helper in domain |
| Loop apply path | `internal/loop/apply.go` — `ApplyChange` still the agent path; keeper `TestLoopApplyDeliberationTransitionEvent` |
| MCP catalog | **10** tools (unchanged; S02 does not add MCP) |

S01 closed **C01, C02, C03, C07** — do not reopen in S02 prompts.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-02, D-22-06, D-22-07, D-22-10, D-22-15, D-22-16
- DF-86: [`DF-84-FORWARD.md`](../../../phase-17-portable-graph-git/DF-84-FORWARD.md)
- Coverage: [`README.md`](../../README.md) C04/C05/C06/C25 rows

## FINAL locked defaults

| Item | Value |
|------|-------|
| Hook target id | **`git-hook`**, tier **CONDITIONAL** |
| Hook install path | Resolve via **`git rev-parse --git-path hooks`** (honors **`core.hooksPath`** + worktrees). **Do not** hardcode `.git/hooks/` only |
| Hook files | **`post-commit`** (required); **`pre-push`** (optional, same marker fragment) |
| Hook kind | **post-commit** and/or **pre-push** — **never** wraps `git commit` (D-22-16) |
| Hook markers | `# begin-trace` / `# end-trace` — append/replace fragment only; uninstall removes fragment; preserve user hook lines |
| Hook body | `trace -C "$ROOT" index` on commit-changed paths (`git diff-tree --no-commit-id --name-only -r HEAD` or equivalent); then **optional** `trace seed export -o trace/graph.json` (best-effort; non-fatal on failure) |
| Binary | `trace` on PATH; hook uses `command -v trace` guard |
| Detect | inside git work tree (`git rev-parse --is-inside-work-tree`) |
| Uninstall | Removes only Trace-owned hook fragments |
| Honesty (commit lag) | When git HEAD ≠ `graph_sync_state.last_indexed_commit` → visible notice on context packet + **`trace index status`** (new subcommand) |
| Honesty (disk lag) | Keep existing `buildIndexHonesty` / `IndexHonesty` — **orthogonal** to commit watermark |
| Mig | **`023_graph_sync.sql`** — table **`graph_sync_state`** single row `id=1`: `last_indexed_commit`, `last_indexed_at`, `hook_installed` (bool/int) |
| Compat after S02-03 | **23** (forbid 024+ until S04) |
| VCS watermark | **Keep** `vcs_meta` / `MetaVCSWatermark` separate — Refresh owns VCS index; graph_sync owns symbol/file index lag |
| Change capture | **`domain.PromoteVCSCommitToChange(oid)`** from `vcs_commits` + paths; idempotent on OID; `source_type=VCS`; `task_id=trace:vcs-capture` (sentinel — no FK on changes.task_id) |
| Meaningfulness | Skip commits with **zero** paths after filter: analyzer-supported ext (`analyzers.DetectLanguage`) **or** path already in `files` table; **`--all`** captures every indexed vcs path |
| State compare | **`domain.CompareStates(fromOID, toOID)`** → added/removed/modified paths + linked change ids when present; **no blobs** |
| CLI compare | `trace changes compare --from <oid> --to <oid>` JSON |
| Daemon | **Forbidden** (D-22-14) |

## Named tests (scope)

| Test | Row |
|------|-----|
| `TestInstallGitHookWritesPostCommit` | S02-01 |
| `TestInstallGitHookDoesNotWrapCommit` | S02-01 |
| `TestUninstallGitHookRemovesFragment` | S02-01 |
| `TestInstallDetectIncludesGitHook` | S02-01 |
| `TestGraphSyncStaleWhenHeadDiffers` | S02-03 |
| `TestHookIndexUpdatesLastIndexedCommit` | S02-03 |
| `TestPromoteVCSCommitCreatesChangeIdempotent` | S02-05 |
| `TestCompareStatesPathDeltaNoBlob` | S02-07 |

## Residual risks for S02-01

| Risk | Mitigation locked in 01 |
|------|-------------------------|
| `core.hooksPath` relative vs absolute | Use `git rev-parse --git-path hooks`; resolve relative to project root |
| Pre-existing user hook content | Marker append; uninstall strips markers only |
| Hook runs before `trace` on PATH | `command -v trace` guard; exit 0 noop |
| Accidental commit wrapper | Grep + `TestInstallGitHookDoesNotWrapCommit`; no `git commit` in installer |
| S02-01 must not bump schema | No 023 in S02-01; watermark lands S02-03 |

## Exit criteria

- [x] 01–08 runnable (thickened below)
- [x] DF-86 locks preserved (no wrap commit; rev-parse hooks path)
- [x] No product Go

## Next

**P22-S02-01**
