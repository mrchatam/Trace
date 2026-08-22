# P08 / S03 / 00-PLANNER — Production hardening

## Metadata
- id: P08-S03-00
- todo_ids: [P08-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling prompts for **production concerns**: migrations hygiene, backup/restore, and **local** auth/binding. No cloud SaaS OAuth. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 8
- Live: `internal/store` embed `schema/001`…`010`; `migrate.go`; Open→`.trace/trace.db` + `trace.lock`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (2026-08-16 — S03-00)

| Surface | Behavior today | Gap vs S03 |
|---------|----------------|------------|
| Migrations | Embed `schema/*.sql` (`001`…`010`); `schema_migrations`; `migrate()` on every `Open`; idempotent re-Open | No CLI status; no documented upgrade/hygiene tests beyond Open idempotency; no downgrade (keep absent) |
| Open | Abs root → mkdir `.trace` → exclusive `trace.lock` → `trace.db` → migrate → `ensureProject` → FTS backfill | Must keep lock; auth gate missing |
| Backup | **Absent** (no export CLI; no archive helpers) | Need consistent `trace.db` snapshot + restore into path-local `.trace/` |
| Auth | Filesystem trust only; MCP/CLI open with no token | Optional **local** token file + env fail-closed — **not** OAuth |
| Bind | `projects.root_path` = Abs at Open; S02 path-local only | Restore must **rebind** root_path to current Abs |
| MCP | Six stdio tools; G19; openStore → store.Open | Inherit auth via Open; **no** new MCP backup/auth tools |
| `011_*` | None (S01/S02 added none) | **Not required** for FS token + file backup; additive only if Notes force DB audit |

### Depends notes (carry-in — do not reopen)
- **S01:** `LanguageAdapterAPIVersion` is a **code const** — do not persist in SQLite; do not couple migrate/backup to adapters.
- **S02:** path-local `<absRoot>/.trace/trace.db` + exclusive `trace.lock` on Open→Close; concurrent same-root → `store.ErrLocked`. Backup/restore/auth **must** target that tree and respect an active lock (no shared-parent DB / swarm).

## Phase defaults already locked (respect — P08-00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Theme | Migrations + backup + local auth |
| Depends | S02 worktrees reviewed before implement (`P08-S02-02` done) |
| Carry-forward | Gate H + honesty + E/F/G + ablation + p0x + x0 + Gate C |
| Forbidden | Cloud SaaS auth theater; daemon as primary |

## Locked defaults (FINAL — P08-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Theme | Migrate hygiene + backup/restore + **local** auth — **not** cloud OAuth |
| Packages | Prefer **`internal/store`** helpers + thin **`cmd/trace`** (stdlib argv; no cobra). No `internal/backup` / `internal/auth` packages unless store helpers become unwieldy |
| Migration policy | Keep embed `schema/NNN_*.sql` + `schema_migrations`; **additive forward-only**; Open applies pending; **no** downgrade CLI; **no** rewrite of applied SQL |
| Migrate API | Export read helpers (e.g. `MigrationStatus` / applied versions + max) used by CLI; Open remains the apply path |
| Migrate CLI | **`trace migrate status`** — machine-friendly stdout (version/max/applied); requires successful Open (lock + auth) |
| `011_*` | **Not mandatory.** Prefer FS token + file backup with **no** new mig. Additive `011_*` only if implementer must store non-secret metadata in SQLite (document in Notes); never for analyzer API version |
| Backup unit | Consistent snapshot of **`trace.db` only** (not whole repo; not source trees). Prefer SQLite backup / `VACUUM INTO` (or equivalent) **while** exclusive `trace.lock` is held |
| Backup exclude | Always exclude runtime `trace.lock`. Exclude `.trace/access.token` **by default**; opt-in `--include-token` |
| Backup CLI | **`trace backup -o <path>`** (or `--output`) → writes snapshot file; clear stderr progress; fail if locked / unauthorized |
| Restore CLI | **`trace restore <path> [--force]`** → install into `<absRoot>/.trace/trace.db`; fail if target exists without `--force`; fail if lock held; then Open → migrate → **rebind** `projects.root_path` to current Abs |
| Restore rebind | After restore Open, `projects.root_path` **must** equal current Abs root (UPDATE if needed). Fail-closed if rebind impossible |
| Auth storage | Optional file **`<absRoot>/.trace/access.token`** (mode **0600**); content = secret string; absent = current trust model |
| Auth gate | On **`store.Open`**: if token file exists and is non-empty, require `TRACE_ACCESS_TOKEN` env equal (constant-time compare). Else exported error (e.g. `store.ErrUnauthorized`) |
| Auth CLI | **`trace auth set <token>`** / **`auth clear`** / **`auth status`** (status = enabled\|disabled only — never print secret) |
| MCP | **No** new tools. Open path inherits token gate (G19). Backup/auth remain CLI-primary |
| Proof | Automated: (1) migrate status + idempotent Open / pending-apply hygiene; (2) backup→restore round-trip (entity survives; no BLOB columns via `HasBlobLikeColumns`); (3) auth fail-closed when token set + env missing/wrong; (4) lock respected on backup/restore |
| Exit codes | Match live CLI: **1** usage; **2** operational (locked / unauthorized / IO) — same taxonomy as S02 lock residual |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Cloud OAuth / hosted IdP; daemon/HTTP primary; source BLOBs in DB or backup; shared-parent `.trace`; swarm; LanguageAdapter coupling; Gate C pack rewrite; Gate H threshold invent; new MCP tools |

## Planner work (this row)
1. Inventory migrate/Open vs backup/auth — **done** (table above).
2. Lock CLI surfaces + local auth + migration policy — **done** (FINAL table).
3. Thicken `01-production-hardening.md` + `02-scope-review.md`; light S04 Depends — **this row**.
4. SCOPE-TODOS sync; board Notes — **this row**.

## Effects on later scopes
| Scope | Note |
|-------|------|
| S04 VERIFY | Checklist must include migrate status, backup round-trip, local-auth fail-closed, no source BLOBs, G19, S02 lock still green. Prefer `evals/compat` ownership stays S04 |

## Exit criteria
- [x] Live inventory recorded
- [x] CLI / auth / migrate policy locked
- [x] `01-*.md` + `02-*.md` runnable alone
- [x] Light S04 Depends notes
- [x] No product Go; board Notes; next **P08-S03-01**

## Minimal todos
- [x] Inventory migrate/Open/backup/auth gaps
- [x] Lock surfaces + exit criteria; thicken 01/02; board
