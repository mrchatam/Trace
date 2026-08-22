# P10-S03-02 — REVIEW-NOTES (index GC)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-20** full-tree index removes DB paths absent from walk set | **Pass** — `cmdIndex` `fullTree` set-diff after walk (`ListFilePaths` ∉ live → `DeleteFileByPath`); `TestIndexGCAfterPathRename` (ghost `a.js` gone; `c.js` has `alpha`; sibling `b.js` intact; stderr `removed 1`) |
| 2 | Delete clears files + symbols + imports + FTS | **Pass** — `DeleteFileByPath` deletes `fts_docs` file-by-id + symbol-by-path (mirrors `SyncFileFTS`), then `DELETE FROM files` (CASCADE); `TestListFilePathsAndDeleteFileByPath` asserts GetFile/symbols/imports + FTS clean + idempotent |
| 3 | Explicit argv index **does not** project-wide GC | **Pass** — GC block gated on `fullTree`; `TestIndexIncrementalIsolation` green |
| 4 | Missing explicit argv path deletes **that** path only | **Pass** — `!fullTree` + `fs.ErrNotExist` → `DeleteFileByPath(rel)` only; `TestIndexArgvMissingPathDeletesOnlyThatPath` |
| 5 | **No** full-rebuild-on-any-change | **Pass** — set-diff deletes only; analyzers stay `IndexFile` upsert/replace per path; no wipe/rebuild graph |
| 6 | **No** new migration `011_*` | **Pass** — schema stops at `010_capability_surface.sql` |
| 7 | **No** MCP index tools; analyzers upsert-only | **Pass** — nine MCP tools; `TestImportBoundaryMCPNoPlanImpactIndexTools`; analyzers have no Delete/ListFilePaths |
| 8 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + S01/S02 + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 suites + product pkgs in `./...`; Gate C artifacts `dry_run:false` N=3; G1 understanding_accuracy **0.8** > B0 **0.0** (untouched) |
| 9 | Board Notes accurate; planner row no product Go | **Pass** — P10-S03-00 Notes claim no product Go; P10-S03-01 Notes cite live tests accurately |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Argv missing path always increments `removed` even if path was never in DB (idempotent delete) | Residual — stderr may over-count; behavior matches locks |
| nit | Experiment ab-index still needs agents to run `trace index` after rename | Residual OK per review prompt — product GC correct when invoked |

## Residuals (explicit)

1. **File-local only** — no dependent-cascade reindex (unchanged Law 4 / DR-INCREMENTAL).
2. **ab-index / dogfood** — operators/agents must still invoke full-tree `trace index` after rename; no filesystem watchers.
3. **`go test ./...`** FAIL only on pre-existing `similar projects/graphify` path space (non-product).
4. Argv missing-path `removed` counter may over-count never-indexed paths (low).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./... -count=1  → product pkgs PASS; FAIL only similar projects/graphify (space)
Gate C: dry_run:false N=3; G1 mean understanding_accuracy 0.800 > B0 0.000 intact
```

## Next

**P10-S04-00** (no spawn inserted).
