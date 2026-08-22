# P08 / S02 / 00-PLANNER — Multi-agent worktrees

## Metadata
- id: P08-S02-00
- todo_ids: [P08-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling prompts for **multi-agent worktrees / safe project bind** after S01 plugin APIs land. Not a swarm framework. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 8
- Live: `cmd/trace` `-C`/`--project`; `store.Open` → `.trace/trace.db`; `internal/gitcli`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (2026-08-16 — S02-00)

| Surface | Behavior today | Gap vs S02 |
|---------|----------------|------------|
| CLI root | `parseGlobal`: `-C`/`--project` or cwd; `resolveRoot` = `filepath.Abs` only | No walk-up; no worktree-aware bind doc/tests |
| MCP root | `resolveProject`: tool `project` → server default → cwd; then Abs | Same Abs-only story |
| `store.Open` | `Abs(root)` → mkdir `<root>/.trace` → open `trace.db`; `SetMaxOpenConns(1)`; `projects.root_path` = abs | Path-local DB already; **no** cross-process lock; no `busy_timeout`/WAL |
| `gitcli.Open` | `git -C <abs>`; `rev-parse --is-inside-work-tree`; store Open on **same** abs (not git common-dir) | Linked worktrees work for git HEAD; store stays per-path (good) but undocumented |
| Concurrent | Unspecified — second process may hit opaque SQLite lock / undefined UX | Need **fail-closed** clear error |
| Shared `.trace` | Not implemented (no parent walk) | Must **lock** “never share / never walk-up” |
| Swarm | Absent | Keep out |

### Git worktree layout vs Trace bind

- **Main worktree** `R/`: Trace bind `Abs(R)` → `R/.trace/trace.db`; gitcli indexes `R` HEAD.
- **Linked worktree** `W/` (`.git` file → common dir): Trace bind `Abs(W)` → `W/.trace/trace.db` (separate); gitcli uses `W` HEAD. **Do not** open `R/.trace` from `W`.
- **Forbidden bind**: walking up from `W` to reuse `R/.trace`, or pointing multiple agents at one abs root without lock.

## Phase defaults already locked (respect — P08-00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Theme | Multi-agent worktrees — safe bind, not orchestration megastore |
| Depends | S01 plugin APIs reviewed (`P08-S01-02` done) before implement |
| Carry-forward | Gate H + honesty + E/F/G + ablation + p0x + x0 + Gate C |
| Forbidden | Daemon/HTTP primary; swarm frameworks |

## Depends note from S01-00 (2026-08-16)
S01 locks `LanguageAdapter` + static table in `internal/analyzers` only — **orthogonal** to worktree / `-C` bind. Do **not** couple project-root resolution or store Open to adapter registration. Keep `DetectLanguage` call sites; IndexFile orchestration stays file-local.

## Locked defaults (FINAL — P08-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Theme | Safe **path-local** project bind + concurrent fail-closed — **not** swarm |
| Bind policy | **Per-root `.trace`** — `store.Open(absRoot)` always uses `<absRoot>/.trace/trace.db` |
| Root resolution | Keep **`filepath.Abs` only** (`-C`/`--project`/cwd; MCP same). **No** walk-up for `.trace` or git common-dir |
| Multi-root | Multiple independent abs roots = multiple isolated DBs; **no** merged multi-root graph |
| Git worktrees | Each worktree checkout root is its own Trace project; gitcli stays `git -C <that root>` |
| Shared parent `.trace` | **Forbidden** this scope (and do not invent symlink tricks to share) |
| Concurrency | **Fail-closed**: exclusive advisory lock file **`<absRoot>/.trace/trace.lock`** acquired in `store.Open`, held until `Close`; second Open on same DB → exported clear error (e.g. `store.ErrLocked`) |
| In-process | Keep `SetMaxOpenConns(1)` |
| SQLite pragmas | Prefer **no** WAL invent; optional `busy_timeout=0` only as belt-and-suspenders — lockfile is primary |
| Migration | **No** `011_*` — lockfile is FS-only under `.trace/` |
| Packages | Prefer **`internal/store`** (Open/Close lock) + thin error surfacing in **`cmd/trace`** / **`internal/mcp`** (G19). No new daemon package |
| CLI / MCP | Same bind rules; lock conflict → non-zero fail (CLI exit **1**); stderr clear message |
| Help / docs | Brief bind note in `trace` help (and/or `store`/`cmd` doc comments): one writer per root; parallel agents → separate worktree roots / `-C` |
| Proof | Automated: path-local isolation + concurrent Open fail-closed; optional git worktree dual-root if `git` available |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Swarm/orchestration; shared-DB worktree mode; daemon/HTTP primary; coupling to `LanguageAdapter`; full-rebuild indexer; Gate C pack rewrite; Gate H threshold invent |

## Planner work (this row)
1. Inventory live root resolution + store Open vs git worktree layouts — **done** (table above).
2. Lock bind policy + concurrency fail-closed — **done** (FINAL table).
3. Thicken `01-worktrees.md` + `02-scope-review.md`; light S03 Depends — **this row**.
4. SCOPE-TODOS sync; board Notes — **this row**.

## Effects on later scopes
| Scope | Note |
|-------|------|
| S03 production | Backup/auth target = **path-local** `<root>/.trace/`; respect `trace.lock` (do not assume shared parent DB). S02 adds **no** mig — S03 may still add `011_*` for backup/auth only |
| S04 VERIFY | Re-prove worktree bind + concurrent fail-closed in checklist evidence |

## Exit criteria
- [x] Live inventory recorded
- [x] Bind + concurrency policy locked
- [x] `01-*.md` + `02-*.md` runnable alone
- [x] Light S03 Depends notes
- [x] No product Go; board Notes; next **P08-S02-01**

## Minimal todos
- [x] Inventory root/`Open`/worktree behavior
- [x] Lock bind + concurrency policy; thicken 01/02; board
