# P08 / S02 / 01 — Multi-agent worktrees / safe project bind

## Metadata
- id: P08-S02-01
- todo_ids: [P08-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **path-local project bind** documentation/behavior + **concurrent-agent fail-closed** locking locked by `00-PLANNER.md`. Multi-root = isolated abs roots (incl. git linked worktrees), **not** swarm orchestration. Do not weaken Gate H / Gate C / G19. Do not couple to `LanguageAdapter`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks finalized 2026-08-16
- [phase README](../../README.md)
- Live: `cmd/trace/root.go` (`-C`/`--project`); `internal/store/open.go`; `internal/gitcli/open.go`; `internal/mcp/project.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (FINAL — P08-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Bind | **Per-root `.trace`** — `<absRoot>/.trace/trace.db` only |
| Root resolve | `filepath.Abs` only — **no** walk-up to parent `.trace` / git common-dir |
| Multi-root | Independent abs roots → independent DBs; no merged graph |
| Git worktrees | Each worktree root = own Trace project; `gitcli` stays `git -C <that root>` |
| Shared `.trace` | **Forbidden** |
| Concurrency | Exclusive advisory lock **`<absRoot>/.trace/trace.lock`** in `store.Open` → held until `Close`; second Open → clear exported error (e.g. `store.ErrLocked`) |
| In-process | Keep `SetMaxOpenConns(1)` |
| Migration | **No** `011_*` |
| Packages | Prefer **`internal/store`** + thin CLI/MCP error surfacing; no daemon |
| CLI / MCP | Lock conflict → fail closed (CLI exit **1**); clear stderr |
| Help | Brief note: one open store per project root; parallel agents use separate `-C` / worktree roots |
| Proof | Path-local isolation test + concurrent Open fail-closed test (+ optional git worktree dual-bind) |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Swarm; shared-DB mode; daemon/HTTP primary; adapter coupling; Gate C rewrite; Gate H invent |

## S01 carry-in (do not reopen)
- Analyzer contribution = `LanguageAdapterAPIVersion=1` + static adapters in `internal/analyzers`
- Worktree bind must **not** depend on adapter registration; leave DetectLanguage/IndexFile contracts intact

## Depends-on
- `P08-S01-02` done; `P08-S02-00` done

## Extension points (exact files)

| File / area | Work |
|-------------|------|
| `internal/store/open.go` (+ small helper e.g. `lock.go`) | Acquire exclusive `trace.lock` under `.trace/` after mkdir, before/with DB open; release on `Close`; export `ErrLocked` (or equivalent) |
| `internal/store/*_test.go` | Isolation + concurrent-open fail-closed tests (primary proof) |
| `cmd/trace/help.go` (and/or init stderr paths) | One-line bind / concurrency note; map `ErrLocked` → exit fail with clear message where Open is used |
| `internal/mcp` (thin) | Surface lock errors clearly via existing `openStore` — no new tools |
| `internal/gitcli` | **No** bind-policy change required if store Open enforces lock; keep store Open on same abs root |
| Docs comments | Optional `doc.go` note on path-local bind |

Do **not** invent `internal/project` unless store helpers become unwieldy — prefer store-owned lock.

## Role work

1. TDD: failing tests for (a) two abs roots → isolated DBs / project rows; (b) second `store.Open` on same root while first open → `ErrLocked` (or locked error).
2. Implement exclusive `trace.lock` in `Open`/`Close`; keep `MaxOpenConns(1)`; **no** mig; **no** walk-up.
3. Surface lock errors from CLI (and MCP open path) with fail-closed exit/message.
4. Brief help text for `-C` / one-writer-per-root.
5. Optional: if `git` present, test linked worktree dual-root (two `.trace` dirs; Heads may differ) — nice-to-have, not a Gate claim.
6. Run locked verify suite; board **status + Notes only**.

### Test requirements (minimum)

**`TestProjectBindPathLocalIsolation`** (name may vary; Notes if renamed):

- Two temp dirs `A` and `B`; `store.Open(A)` and `store.Open(B)` succeed concurrently (different roots).
- Assert distinct `DBPath()` under `A/.trace/` vs `B/.trace/`.
- Write a distinguishable row in A (e.g. ensureProject / insert via public store API already used in tests); B must not see A's `projects.root_path` / entity.

**`TestConcurrentStoreOpenFailClosed`**:

- `Open(root)` succeeds; second `Open(root)` in another goroutine/process-equivalent **fails** with locked error without corrupting DB.
- After `Close` of first, `Open(root)` succeeds again.

Do **not** require multi-OS flock portability theater beyond what Go stdlib/`golang.org/x/sys` already used in-repo allows — prefer portable approach already acceptable for Linux CI; document residual if any.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional: confirm Gate C artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false` N=3 (do **not** rewrite packs).

## Exit criteria
- [ ] Path-local bind unchanged in spirit (Abs → `<root>/.trace/trace.db`); **no** parent walk-up / shared `.trace`
- [ ] Exclusive `trace.lock` on Open; second Open fail-closed with clear error
- [ ] Isolation + concurrent fail-closed tests green
- [ ] CLI/MCP surface lock failure clearly; help note present
- [ ] **No** `011_*`; no swarm/daemon/HTTP; no LanguageAdapter coupling
- [ ] Carry-forward suite green (incl. Gate H / Gate C intact)
- [ ] Board Notes ready for **P08-S02-02**

## Out of scope
- Plugin API redesign (S01); backup/auth (S03); checklist harness (S04)
- Shared parent `.trace` mode; swarm / agent orchestration frameworks
- Daemon / HTTP / embeddings / new MCP tools
- Inventing Gate H thresholds or rewriting Mode-B Gate C packs
- WAL mode as primary concurrency strategy

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] Add isolation + concurrent Open fail-closed tests (fail first)
- [ ] Implement `.trace/trace.lock` in store Open/Close; export clear ErrLocked
- [ ] Surface lock errors in CLI (+ MCP open); brief help note
- [ ] Run locked verify suite; mark P08-S02-01 done with Notes
