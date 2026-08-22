# P00 / S02 / 01 — SQLite .trace/ store

## Metadata
- id: P00-S02-01
- todo_ids: [P00-S02-01]
- role: implementer
- skills: [incremental-implementation]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement per-project SQLite under **`.trace/`** in `internal/store`: migrations, core causal entity tables, thin append-only events, and File/Symbol/Import stubs that support **per-file incremental** replace (DR-INCREMENTAL / P0-X #7 substrate). No analyzers, retrieval, VCS, or CLI wiring.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) (store + P0-X bar)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) (esp. G1, G5, G12, G18)
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-DB, DR-TRACEDIR, DR-EVT, DR-CLAIM, DR-INCREMENTAL
- [STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md)
- [PROJECT_MODEL.md](../../../../PROJECT_MODEL.md) provenance fields
- Historical T002: [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md)
- Prior: S01 done — module `github.com/mrchatam/Trace`, stub `internal/store/doc.go`, `.gitignore` has `.trace/`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Package path | `internal/store` (package `store`; **not** `internal/sqlite`) |
| Store dir | `<projectRoot>/.trace/` — create if missing (`0755`); already gitignored |
| DB file | `<projectRoot>/.trace/trace.db` (one SQLite per bound project; DR-DB) |
| Driver | `modernc.org/sqlite` (pure Go, **no CGO**) via `database/sql` |
| Migrations | Hand-rolled: table `schema_migrations(version INTEGER PRIMARY KEY)` + embedded SQL (`embed.FS`); apply automatically on `Open` |
| Entity IDs | `TEXT` primary keys (UUID v4 strings) |
| Events | Thin append-only `events` table (DR-EVT); **no** UPDATE/DELETE API for events |
| Timestamps | RFC3339 / ISO-8601 text in UTC (or integer unix seconds — pick one and use consistently; prefer TEXT RFC3339) |
| Provenance columns (semantic entities) | `source_type`, `confidence`, `status` (`ACTIVE` \| `STALE` \| `SUPERSEDED` allowed), `created_at`, `updated_at`, `last_verified_at` (nullable OK where unused) |
| File identity (incremental) | `files.path` UNIQUE (repo-relative, forward slashes); `content_hash` TEXT NOT NULL; optional `git_oid` TEXT; upsert by path |
| Incremental contract | Replacing symbols/imports for one `path` must **not** delete or rewrite other files’ rows (DR-INCREMENTAL substrate for S04 / P0-X #7) |
| Forbidden | Source file BLOBs / body columns; full-rebuild-as-default architecture; MCP/daemon/HTTP; analyzers; FTS5; CLI `init` (S07 calls `store.Open`) |
| Depends on | S01 scaffold done |
| Repo state at start | `internal/store/doc.go` stub only; no SQLite dep yet |

### Tables required (migration v1)

**Meta / project**

- `schema_migrations`
- `projects` — at least `id`, `root_path` (or bind key), `created_at`

**Causal / work (materialized; DR-EVT + DR-CLAIM)**

- `goals`, `decisions`, `assumptions`, `tasks`, `discoveries`, `claims`, `evidence`, `reviews`, `plan_changes`
- Each: `id TEXT PK`, title/body or equivalent TEXT fields as needed for roundtrip, plus provenance columns above
- FKs between entities may be nullable stubs; S05 owns rich domain API — store needs durable rows + enough columns for Goal/Task/Event tests

**Events (append-only)**

- `events`: `id TEXT PK`, `ts`, `type`, `entity_type`, `entity_id`, `payload_json` (TEXT, may be `{}`)
- Insert-only from store API

**Code graph stubs (filled by S04; schema + helpers here)**

- `files`: `id`, `path` UNIQUE, `content_hash`, `git_oid` NULL, `language` NULL, `indexed_at`, status-ish columns as needed
- `symbols`: `id`, `file_id` FK → `files`, `name`, `kind`, `start_line`, `end_line` (or equivalent span)
- `imports`: `id`, `file_id` FK, `imported_path` (or similar), optional `symbol` NULL

**Hard rule:** no column that stores source file contents.

### Minimum public API (package `store`)

```text
Open(projectRoot string) (*Store, error)  // mkdir .trace/, open DB, migrate
Close() error
// Goal + Task: Upsert/Get (or Insert+Get) sufficient for P0-X #1 substrate
// Event: Append + list/get by entity or recent
// File: UpsertFile(path, contentHash, gitOID optional)
//       ReplaceFileSymbols(path, symbols) — delete symbols for that file_id only, then insert
//       (imports: same per-file replace pattern if implemented in S02)
```

Exact method names may vary; behavior above is locked. Keep types in `internal/store` (S05 may wrap later).

### Target tree (under `internal/store/`)

```text
internal/store/
  doc.go                 # keep / refresh package comment
  open.go                # Open / Close / Store type
  migrate.go             # schema_migrations + apply
  schema/
    001_init.sql         # or equivalent embedded migration(s)
  entities.go            # Goal/Task/Event (+ File helpers) — split files OK
  file_graph.go          # UpsertFile / ReplaceFileSymbols (optional split)
  store_test.go          # temp-dir roundtrips + incremental isolation
```

Do **not** create `internal/sqlite`. Do **not** wire `cmd/trace` in this scope.

## Board rights
Implementer: update **status + notes only** on `P00-S02-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `Open` on a temp project root creates `.trace/trace.db` and applies migrations idempotently (second `Open` does not fail / re-apply)
- [ ] Tables exist for: projects, goals, decisions, assumptions, tasks, discoveries, claims, evidence, reviews, plan_changes, events, files, symbols, imports, schema_migrations
- [ ] Goal + Task + Event roundtrip test passes (insert/get; event append visible)
- [ ] Provenance columns present on Goal (and preferably Task); values survive roundtrip for at least `source_type` + `status`
- [ ] File upsert + `ReplaceFileSymbols` (or equivalent) for path A does **not** remove symbols for path B (incremental isolation test)
- [ ] Schema has **no** source-content BLOB / body column (test or explicit assertion)
- [ ] Driver is `modernc.org/sqlite`; `go test ./internal/store/...` passes; `go test ./...` still green
- [ ] No MCP/daemon/HTTP; no analyzer/retrieval code in this scope
- [ ] TODO.md Notes for `P00-S02-01` updated; status `done`

## Minimal todos
- [ ] Add `modernc.org/sqlite`; implement `Open` / `Close` + `.trace/` creation
- [ ] Embedded migration v1: all required tables + provenance + File/Symbol/Import stubs
- [ ] Goal / Task / Event write+read helpers
- [ ] File upsert + per-file symbol replace helpers
- [ ] Tests: open/migrate idempotent; Goal+Task+Event roundtrip; incremental file isolation; no content BLOB
- [ ] Board status + notes (DB path, migration version, test commands)
