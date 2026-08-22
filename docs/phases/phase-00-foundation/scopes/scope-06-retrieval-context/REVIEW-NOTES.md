# P00-S06-02 — Scope review notes (2026-08-15)

Independent review of S06 against `01-retrieval.md` + TODO Notes for `P00-S06-01`. Fresh session; claims verified in-repo.

## Plan (executed)

1. Diff claims vs `internal/retrieval`, `internal/compiler`, store mig `004_fts` + FTS/Sync APIs
2. Re-run `CGO_ENABLED=0` retrieval/compiler/store/domain and full `./...`
3. Severity-tag findings; inline-fix high (FTS Open backfill); no spawns
4. Write these notes; mark board + SCOPE-TODOS; light thicken upcoming S07 Open note

## Verdict

**APPROVE** — one **high** fixed inline (FTS upgrade backfill). Confidence: **high**. Spawns: **none**. Next board row: **P00-S07-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Store mig `004_fts.sql` applies on Open; unicode61; no source BLOBs | Pass (`schema/004_fts.sql`; store Open test lists `fts_docs`; `HasBlobLikeColumns` skips FTS5 shadow tables) |
| `SearchFTS` + `SyncEntityFTS` / `SyncFileFTS` / `RebuildFTS`; Upsert* wired | Pass (`fts.go`, `entities.go`, `entities_causal.go`, `file_graph.go`) |
| Open backfill when FTS empty but content exists (post-004 upgrade) | Pass after review fix (`ensureFTSPopulated` + `TestFTSBackfillOnOpenWhenIndexEmpty`) |
| `ListTasksByGoalID` / `GetFileByID` | Pass (`helpers.go` + tests) |
| `retrieval.Exact` / `Search` / `Expand` (depth 1–2) / `Why` + reason codes; depth >2 rejected | Pass (`retrieval_test.go`) |
| Graph expand: causal links + `goal_id` + structural symbols/imports | Pass (expand neighbors + structural test) |
| `compiler.TaskContext` Layer 0–1 JSON+MD; default depth 1; `ExpandContext` depth 2 | Pass (`compiler_test.go`) |
| Budgets 4096 / 32; `truncated` when drops; hard item cap | Pass (`TestTaskContextAndBudgets`, `TestItemCapNeverExceeded`) |
| Untrusted labeling on retrieved text (JSON + MD callouts) | Pass (`trust=untrusted_data`; MD contains `untrusted_data` + authority callout) |
| No embeddings / vector deps; packages `retrieval` + `compiler` only (no context/contextx) | Pass (`go.mod` grep; package tree) |
| No dump API; no MCP/daemon/HTTP; no new CLI; no analyzers/gitcli imports | Pass (`cmd/trace` help/version only; import grep) |
| `CGO_ENABLED=0` retrieval/compiler/store/domain; `go test ./...` | Pass (re-run below) |

## Re-verification commands (2026-08-15)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... -count=1   # ok
go test ./... -count=1                                                                                                    # ok
```

## Findings

### Blocker

None.

### High (fixed inline)

- **FTS empty after 004 onto existing DB:** Sync on Upsert covered greenfield inserts, but applying `004_fts` left `fts_docs` empty for pre-existing goals/tasks/files/symbols until a manual `RebuildFTS`. Failure mode: silent empty `Search` / weak TaskContext FTS enrichment on upgraded `.trace/` DBs. **Fix:** `store.ensureFTSPopulated()` from `Open` when `fts_docs` is empty but content tables have rows; regression `TestFTSBackfillOnOpenWhenIndexEmpty`. Doc comment updated in `store/doc.go`.

### Medium (residual — no spawn)

- **Silent omit on optional loads:** `compileAtDepth` ignores `GetGoal` errors when `goal_id` is set, and swallows `Why` errors when `IncludeWhy` is true. Layer 0 can omit a broken goal without surfacing failure; optional Why can vanish silently. Acceptable for P0-X if store FKs hold; S07 should treat missing goal as a hard error at the CLI boundary if desired.
- **Byte-sliced excerpts:** `excerpt` / `excerptBody` cut at 240 bytes (not runes). Can yield invalid UTF-8 mid-emoji; current `encoding/json` tolerates via replacement. Prefer rune-safe trim later.

### Low / nit

- `TestNoDumpAPI` only references package constants — API absence is still clear from package surface/docs.
- Layer 0 `work_state` carried as synthetic `task_state` item (`trust=system`) rather than fields on the task item — meets the lock, slightly unusual shape for S07 adapters.
- Spec soft-prefers `domain.Service` reads; locked ctor is `retrieval.New(*store.Store)` — store reads are intentional.

## Spawns

None.

## Residual risks

- Partial FTS drift if rows are mutated outside Upsert*/Replace* paths (no triggers by design).
- Why expands at default depth 1 only (plus events / optional VCS) — sufficient for P0-X #3; deeper chains need explicit Expand.
- Q-INJECTION remains OPEN; untrusted labeling is provisional mitigation (A14).
- Layers 2–3 / embeddings deferred (DR-NOSSEM / DR-PACK).

## Forward edits this review

- `internal/store/open.go`, `fts.go`, `fts_test.go`, `doc.go` — Open FTS backfill + test
- `SCOPE-TODOS.md` — mark S06-01/02 done
- `docs/TODO.md` — `P00-S06-02` → done + notes
- `scopes/scope-07-cli/01-cli.md` — note Open FTS backfill (CLI must not reimplement Rebuild)
