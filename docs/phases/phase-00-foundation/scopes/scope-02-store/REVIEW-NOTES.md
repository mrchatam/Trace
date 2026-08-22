# P00-S02-02 — Scope review notes (2026-08-15)

Independent review of S02 against `01-store.md` + TODO Notes for `P00-S02-01`. Fresh session; claims verified in-repo.

## Verdict

**APPROVE** — no blocker / high findings. Confidence: **high**.

Go `1.22` → `1.24.0` bump is **acceptable** (required by `modernc.org/sqlite@v1.45.0`, which declares `go 1.24.0`). Documented forward in README + upcoming S03–S07 prompts; no revert/spawn.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `Open` → `<root>/.trace/trace.db`, mkdir `0755`, migrate idempotent | Pass (`open.go`, `TestOpenCreatesDBAndMigratesIdempotent`, `TestTraceDirPermissionsAndPath`) |
| Embedded migration v1 + `schema_migrations` | Pass (`migrate.go` + `schema/001_init.sql`; version `1` recorded) |
| Tables: projects, goals, decisions, assumptions, tasks, discoveries, claims, evidence, reviews, plan_changes, events, files, symbols, imports, schema_migrations | Pass (test enumerates all) |
| Goal + Task + Event roundtrip; provenance `source_type` + `status` | Pass (`TestGoalTaskEventRoundtrip`) |
| File upsert + `ReplaceFileSymbols` / `ReplaceFileImports` isolate per path (DR-INCREMENTAL) | Pass (`TestReplaceFileSymbolsIncrementalIsolation`) |
| No source-content BLOB / body on `files` | Pass (`TestNoSourceContentColumns` + schema audit) |
| Driver `modernc.org/sqlite` v1.45.0; no CGO required | Pass (`go.mod`; `CGO_ENABLED=0 go test ./internal/store/...` ok) |
| `go test ./internal/store/...` and `go test ./...` | Pass (2026-08-15 re-run) |
| No MCP/daemon/HTTP; no `internal/sqlite`; `cmd/trace` not wired to store | Pass |
| Package tree matches locked layout | Pass (`doc.go`, `open.go`, `migrate.go`, `schema/001_init.sql`, `entities.go`, `file_graph.go`, `store_test.go`) |

## Findings

### Medium (docs drift — fixed forward, no spawn)

- **README** still said “Requires Go 1.22+” while `go.mod` is `go 1.24.0` / `toolchain go1.24.2`. Updated README to 1.24+.
- **S04 `01-analyzers.md`** locked Persistence only named `ReplaceFileSymbols`; store also ships `ReplaceFileImports`. Thickened S04 (+ light notes on S03/S05/S07 for Go 1.24 and store surface).

### Low / nit (no spawn)

- `schema_migrations` is created both in `migrate.go` and again in `001_init.sql` (`IF NOT EXISTS`) — harmless redundancy.
- `EnsureRFC3339` in `entities.go` is exported but unused in-package; fine to leave until S05/S07 need it.
- Causal entity `body` columns exist (allowed); only `files` source bodies are forbidden — correctly scoped in tests.

### Blocker / high

None.

## go.mod 1.24 judgment

| Question | Answer |
|----------|--------|
| Required by sqlite dep? | Yes — `modernc.org/sqlite@v1.45.0` module `GoVersion` / `go.mod` = `1.24.0` |
| Acceptable vs S01 “do not bump without need”? | Yes — need established by chosen pure-Go driver |
| Spawn pin/downgrade? | **No** — staying on current modernc is preferable to freezing an older sqlite for 1.22 |
| Action | Document in upcoming prompts + README (done this review) |

## Spawns

None.

## Residual risks

- Cross-entity FKs among causal tables remain stubby (e.g. `tasks.goal_id` not `REFERENCES goals`) — intentional for S02; S05 owns richer domain constraints.
- Path normalization is slash-normalize only (no `..` rejection); acceptable until CLI/analyzer inputs are untrusted at the boundary.
- Event append-only is API-level only (no UPDATE/DELETE helpers); raw SQL could still mutate — fine for local single-user store.

## Forward edits this review

- `README.md` — Go 1.24+
- `scopes/scope-03-vcs/01-vcs.md` — Go 1.24 lock note
- `scopes/scope-04-analyzers/01-analyzers.md` — `ReplaceFileImports` + Go 1.24
- `scopes/scope-05-causal/01-causal.md` — store Goal/Task/Event helpers exist
- `scopes/scope-07-cli/01-cli.md` — `store.Open` + Go 1.24
