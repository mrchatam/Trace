# P11-S01-00 — Index partial-path GC (FINAL)

## Metadata
- id: P11-S01-00
- todo_ids: [P11-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts for **DF-40**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 / DR-INCREMENTAL
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-40 (DF-20 residual)
- Phase 10 S03 FINAL: [../../../../phases/phase-10-integrity-surfaces/scopes/scope-03-index-gc/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-03-index-gc/00-PLANNER.md) — full-tree GC + argv missing-path delete; **not** partial-rename
- Live: `cmd/trace/index.go`; `internal/store/file_graph.go` `ListFilePaths` / `DeleteFileByPath` / `ContentHash`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (A1–A7 + live inventory + P10 S03 locks do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `cmd/trace/index.go` | `fullTree := len(args)==0`: after walk, set-diff GC via `ListFilePaths` + `DeleteFileByPath`. **Partial argv:** index upsert only; missing-on-disk argv → single-path delete. **No** rename / hash orphan GC on partial argv |
| stderr | `indexed N, skipped M, removed K` (P10 S03) |
| `internal/store` | `ListFilePaths`, `DeleteFileByPath` (+ FTS file+symbol cleanup), `GetFileByPath` exposes `ContentHash`; **no** `ListFilePathsByContentHash` yet |
| DF-20 tests | `TestIndexGCAfterPathRename` (full-tree), `TestIndexArgvMissingPathDeletesOnlyThatPath`, `TestIndexIncrementalIsolation` — all green; none cover rename + `index <new-path>` only |
| Dogfood | `experiments/_post_p10/BUGHUNT.md` DF-40 repro; `POST-P10-DOGFOOD.md` — `index <new-path>` leaves `src/old.ts` / `src/a.ts` until bare `index` |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-40 | After rename, `trace index <new-path>` must drop ghost old path/symbols/FTS without requiring full-tree index | **Content-hash orphan GC** on partial argv — not project-wide set-diff, not full-rebuild |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-40** only (DF-20 full-tree + missing-argv stays as shipped in P10 S03 — do not regress) |
| Packages | **`cmd/trace/index.go`** — orchestration after successful partial index; **`internal/store`** — optional helper `ListFilePathsByContentHash(hash string) ([]string, error)` (or equivalent query). Keep analyzers upsert-only; **no** `internal/indexer` package |
| Migration | **None** — reuse `DeleteFileByPath` + existing `files.content_hash` |
| GC trigger — full tree | **Unchanged** (P10): `len(args)==0` → set-diff delete every DB path ∉ live walk set |
| GC trigger — missing argv | **Unchanged** (P10): explicit path missing on disk → delete **that** path only |
| GC trigger — partial argv + rename (NEW) | When `len(args)>0` and a path was **successfully indexed**, after upsert: for that file’s `content_hash` H, find other DB paths with hash H; for each candidate path ≠ indexed path, if **missing on disk** → `DeleteFileByPath`. Count toward `removed` |
| Isolation bar | Partial argv must **not** delete siblings that still exist on disk (even if content hashes collide / duplicate files). `TestIndexIncrementalIsolation` stays green |
| Detection strategy | **Content-hash orphan** preferred (works for uncommitted renames + dogfood fixtures). Git rename detection is **optional sugar only** — not required; must not replace hash orphan logic |
| Residual OK | Rename **plus** content edit so hash ≠ old row: may leave ghost until full-tree `index` or explicit missing-argv delete of old path — document in Notes; not a DF-40 fail |
| Progress stderr | Keep `indexed N, skipped M, removed K`; DF-40 rename+partial must show `removed ≥ 1` (or exact `removed 1` in test) |
| Architecture | **DR-INCREMENTAL** — set-diff / orphan deletes only; **forbidden** wipe/rebuild entire code graph on single-file index |
| Tests (required) | (1) **`TestIndexPartialArgvGCAfterRename`** (or equiv): index `a.js`+`b.js` → rename `a.js`→`c.js` → `trace index c.js` only → `a.js` gone from store+FTS; `c.js` present with symbols; `b.js` untouched; stderr contains `removed`. (2) Keep `TestIndexGCAfterPathRename`, `TestIndexArgvMissingPathDeletesOnlyThatPath`, `TestIndexIncrementalIsolation`. (3) Prefer assert no ghost FTS for old path |
| Suggested store unit | If adding `ListFilePathsByContentHash`: unit test hash match + empty |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-20 behaviors; prior P11 N/A (first product scope) |
| Forbidden | Project-wide set-diff on every argv index; full-rebuild indexer; new mig unless blocked; MCP index tools; daemon/HTTP/embeddings; rewriting Phase 00–10 `done` history; Mode-B Gate C pack rewrite; S02+ product work |

## Effects on later scopes
- **S02** (review/operator): no index coupling — serial after S01 review only.
- **S08 VERIFY:** include DF-40 partial-argv rename regression in evidence table (alongside DF-20 full-tree).

## Exit
- [x] Thicken `01-index-partial-path-gc.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S02)
- [x] Board Notes; next **P11-S01-01**
- [x] Product Go — **not** this row
