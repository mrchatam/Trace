# P10-S03-00 — Index GC (FINAL)

## Metadata
- id: P10-S03-00
- todo_ids: [P10-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S03 implement/review prompts for **DF-20** (ghost files/symbols after rename). **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 incremental; DR-INCREMENTAL
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — index ghosts cluster
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-20 (product); DF-14 residual fuel
- Live: `cmd/trace/index.go` (walk + IndexFile; **no** delete-on-missing); `internal/store/file_graph.go` (Upsert/Replace*; **no** DeleteFile*); `internal/store/fts.go` `SyncFileFTS`; schema `001` symbols/imports `ON DELETE CASCADE`
- S01/S02 inherit — do not re-litigate; S03 owns index GC only
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no user grill required (phase A locks + live inventory).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `cmd/trace/index.go` | Full walk or argv → `indexOne` / `analyzers.IndexFile(AtRev)` upsert only; stderr `indexed N, skipped M`; **never** removes DB rows for paths absent on disk |
| `reindex` | Alias of `index` (same path) |
| `internal/store` | `UpsertFile` / `ReplaceFileSymbols` / `ReplaceFileImports` / `SyncFileFTS`; **no** `ListFilePaths` / `DeleteFileByPath` |
| FK | `symbols`/`imports` → `files(id) ON DELETE CASCADE`; **FTS** (`fts_docs`) is **not** FK-cascaded — must delete explicitly |
| Isolation | `TestIndexIncrementalIsolation` — indexing `a.js` alone must not touch `b.js` |
| T0 / gitignore | Walk already excludes T0 + unsupported lang + best-effort gitignore (P07-S01) |
| Dogfood | `_bughunt/rename_test/`; ab-index / ab-index-stale — rename leaves ghost path/symbol until product GC |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-20 | Incremental `trace index` drops missing paths from files/symbols/imports/FTS | **Delete-on-missing GC** on full-tree index; **not** full-rebuild-on-any-change |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Packages | **`internal/store`** — `ListFilePaths` + `DeleteFileByPath` (name OK if equivalent); **`cmd/trace/index.go`** — GC orchestration after full walk. Keep `internal/analyzers` IndexFile upsert-only (no analyzer rewrite / no `internal/indexer`) |
| Migration | **None** — prefer no `011_*`; DELETE existing rows + FTS |
| GC trigger — full tree | When `len(args)==0` (walkIndexable): after indexing live set, **delete every `files.path` not in the live walk set** |
| GC trigger — explicit argv | **No** project-wide GC (preserve file-local isolation). If an explicit path is **missing on disk**, **delete that path only** from the graph (optional soft progress); do **not** fail the whole command solely because the path is gone if delete succeeds — prefer: missing argv → attempt delete → count as removed/skipped, not hard fail unless store delete errors |
| Live set definition | Same as today’s walk: T0 → DetectLanguage → T0 file → gitignore. Paths that leave the indexable set (deleted, renamed away, newly T0/ignored) are GC candidates on next full index |
| Delete semantics | `DeleteFileByPath`: normalize path; delete `fts_docs` for `entity_type='file'` by file id **and** `entity_type='symbol'` by path (mirror SyncFileFTS cleanup); `DELETE FROM files WHERE path=?` (CASCADE symbols/imports); idempotent OK (missing path = no-op success) |
| List helper | `ListFilePaths() ([]string, error)` — all repo-relative paths currently in `files` (order stable preferred) |
| Progress stderr | Extend to include removals, e.g. `indexed N, skipped M, removed K` (wording flexible; **K** must be observable in tests or Notes) |
| Architecture | **DR-INCREMENTAL** — never wipe/rebuild entire code graph on a single file change; GC is set-diff delete only |
| Tests (required) | (1) **Rename/full-index GC:** index two files → rename/move one path on disk → `trace index` (no args) → old path absent from store + FTS; new path present with current symbols; sibling unchanged. (2) **Isolation:** keep `TestIndexIncrementalIsolation` green (argv single-file must not GC sibling). (3) Prefer assert no ghost symbol name via `ListSymbolsByPath` / FTS after rename |
| Suggested test name | `TestIndexGCAfterPathRename` (or equivalent) in `cmd/trace` |
| Carry-forward | honesty A/B/C + Gate G; Gate E/F; ablation; Gate H; compat; p0x; x0; S01 why/Exact; S02 nine MCP tools; Gate C `dry_run:false` intact |
| Forbidden | Full-rebuild-as-default; new mig unless blocked; MCP index tools; daemon/HTTP/embeddings; rewriting Phase 00–09 / S01 / S02 `done` history; Mode-B Gate C pack rewrite; S04 product work |

## Effects on later scopes
- **S04:** serial after S03 review; no index coupling (operator/capability gates).
- **S05 VERIFY:** include DF-20 rename+full-index regression + isolation bar in evidence table.

## Exit
- [x] Thicken `01-index-gc.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S05 stub if useful)
- [x] Board Notes; next **P10-S03-01**
- [x] Product Go — **not** this row
