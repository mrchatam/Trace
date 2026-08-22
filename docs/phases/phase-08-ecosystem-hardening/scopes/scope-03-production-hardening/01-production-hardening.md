# P08 / S03 / 01 — Production hardening

## Metadata
- id: P08-S03-01
- todo_ids: [P08-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **migrations hygiene + backup/restore + local auth/binding** locked by `00-PLANNER.md`. No cloud OAuth. No source BLOBs. Respect S02 path-local `.trace` + `trace.lock`. Do not weaken Gate H / Gate C / honesty / G19.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks finalized 2026-08-16
- [phase README](../../README.md)
- Live: `internal/store/{open,migrate,lock}.go`; `cmd/trace/root.go` + `help.go`; `HasBlobLikeColumns`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (FINAL — P08-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | Prefer **`internal/store`** + thin **`cmd/trace`** (no cobra; no new top-level packages unless Notes) |
| Migrate | Embed pattern unchanged; forward-only; Open applies pending; **`trace migrate status`** |
| `011_*` | **Not mandatory** — prefer none; additive only if Notes force non-secret DB metadata |
| Backup | Consistent **`trace.db`** snapshot while lock held → **`trace backup -o <path>`** |
| Restore | **`trace restore <path> [--force]`** → path-local `.trace/trace.db` + migrate + **rebind** Abs `root_path` |
| Token file | Optional **`.trace/access.token`** (0600); gate on **`store.Open`** via env **`TRACE_ACCESS_TOKEN`** |
| Auth CLI | **`trace auth set` / `clear` / `status`** |
| MCP | No new tools; Open inherits auth |
| Proof | Migrate hygiene + backup round-trip + auth fail-closed + lock respect |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Cloud OAuth; daemon/HTTP primary; source BLOBs; shared-parent DB; swarm; adapter coupling; Gate C rewrite; Gate H invent |

## S01 / S02 carry-in (do not reopen)
- No analyzer API version in SQLite; do not fold `LanguageAdapter` into migrate/backup
- Path-local `<absRoot>/.trace/`; exclusive `trace.lock`; backup must not assume shared parent DB
- Lock conflict / unauthorized → operational fail (CLI exit **2** per live taxonomy)

## Depends-on
- `P08-S02-02` done; `P08-S03-00` done

## Extension points (exact files)

| File / area | Work |
|-------------|------|
| `internal/store/migrate.go` (+ small status helper) | Export migration status (applied versions / max / embed expected); keep embed apply path |
| `internal/store/open.go` (+ `auth.go` or equivalent) | After lock (or with clear order documented): if `.trace/access.token` non-empty, require env match; export `ErrUnauthorized` (name may vary — Notes if renamed) |
| `internal/store/backup.go` (or equivalent) | Backup snapshot + Restore into project `.trace/`; rebind `projects.root_path`; never copy source trees |
| `internal/store/*_test.go` | Hygiene + round-trip + auth fail-closed + lock interaction tests |
| `cmd/trace/root.go` + `help.go` | Dispatch `migrate` / `backup` / `restore` / `auth`; help lines |
| `cmd/trace/migrate.go`, `backup.go`, `restore.go`, `auth.go` | Thin argv adapters (G19 — call store only) |
| `internal/mcp` | **No** new tools; ensure Open auth errors surface clearly (same as lock) |

Do **not** invent WAL-as-primary concurrency or shared `.trace` modes.

## Role work

1. TDD: failing tests for migrate status / Open idempotency (extend if needed), backup→restore round-trip, auth fail-closed, backup fails when locked.
2. Implement store helpers (status, backup, restore+rebind, token gate on Open).
3. Thin CLI commands + help; map errors to exit **1** (usage) / **2** (operational).
4. Prefer **no** `011_*`. If adding one, document why in board Notes and keep additive.
5. Assert backup/restore path does not introduce source BLOB columns (`HasBlobLikeColumns` false).
6. Run locked verify suite; board **status + Notes only**.

### Test requirements (minimum)

**Migrate hygiene** (names may vary; Notes if renamed):

- Re-Open on already-migrated DB is idempotent (existing coverage OK if still green).
- `MigrationStatus` (or CLI-tested equivalent) reports max version matching embed `010` (or latest after any optional `011_*`).
- Optional: DB stopped before latest embed version receives pending mig on Open (temp fixture).

**`TestBackupRestoreRoundTrip`**:

- Open root A; write a distinguishable store/domain row (public API already used in tests).
- `Backup` to temp file; close; restore into root B (or A with `--force` semantics at API level); Open B; row visible; `projects.root_path` == Abs(B); `HasBlobLikeColumns` false.
- Token file not included unless explicitly requested.

**`TestLocalAccessTokenFailClosed`**:

- Write `.trace/access.token`; Open without env → error; wrong env → error; matching `TRACE_ACCESS_TOKEN` → success.
- After `auth clear` (or remove file), Open without env succeeds.

**`TestBackupFailsWhenLocked`** (or restore):

- Hold Open on root; concurrent backup/restore attempt fail-closed (`ErrLocked` or equivalent) without corrupting DB.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional: confirm Gate C artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false` N=3 (do **not** rewrite packs).

## Exit criteria
- [ ] Embed migrate hygiene intact; `trace migrate status` works
- [ ] `trace backup` / `trace restore` round-trip green; no source BLOBs; rebind Abs root
- [ ] Local token gate on Open + `trace auth *`; fail-closed tests green
- [ ] Respects `trace.lock`; path-local `.trace` only
- [ ] No cloud OAuth / daemon / new MCP tools; no LanguageAdapter coupling
- [ ] Prefer no `011_*` (or documented additive only)
- [ ] Carry-forward suite green (incl. Gate H / Gate C intact)
- [ ] Board Notes ready for **P08-S03-02**

## Out of scope
- Plugin APIs (S01); worktree redesign (S02); `evals/compat` ownership (S04)
- Daemon / always-on HTTP / cloud SaaS auth / hosted IdP
- Downgrade migrations; WAL-as-primary concurrency; swarm
- Inventing Gate H thresholds or rewriting Mode-B Gate C packs

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] Add migrate/backup/auth fail-closed tests (fail first)
- [ ] Implement store status + backup/restore+rebind + Open token gate
- [ ] Thin CLI `migrate` / `backup` / `restore` / `auth` + help
- [ ] Run locked verify suite; mark P08-S03-01 done with Notes
