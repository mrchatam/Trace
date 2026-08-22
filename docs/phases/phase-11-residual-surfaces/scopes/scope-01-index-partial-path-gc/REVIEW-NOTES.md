# P11-S01-02 — REVIEW-NOTES (index partial-path GC)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-40** rename + `index <new-path>` removes old path/symbols/FTS | **Pass** — `gcContentHashOrphans` after successful partial argv upsert; `TestIndexPartialArgvGCAfterRename` (ghost `a.js` gone + FTS clean; `c.js` has `alpha`; `b.js` intact; stderr contains `removed`) |
| 2 | Mechanism is content-hash orphan GC — **not** project-wide argv set-diff | **Pass** — partial path uses `ListFilePathsByContentHash` + missing-on-disk delete only; `fullTree` set-diff remains gated on `len(args)==0` |
| 3 | Isolation: on-disk siblings survive partial argv | **Pass** — orphans skipped when `Stat` succeeds; `TestIndexIncrementalIsolation` green |
| 4 | P10 DF-20 retained: full-tree GC + missing-argv single delete | **Pass** — `TestIndexGCAfterPathRename`, `TestIndexArgvMissingPathDeletesOnlyThatPath` green; code paths unchanged in shape |
| 5 | G19 — no domain fork in adapters; analyzers upsert-only | **Pass** — GC only in `cmd/trace/index.go` + store helper; analyzers have no Delete/ListFilePathsByContentHash |
| 6 | No forbidden architecture | **Pass** — no `011_*` mig (schema ends `010_capability_surface.sql`); no MCP index; no daemon/HTTP/embeddings/full-rebuild/`internal/indexer` |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 suites + product `./cmd/... ./internal/... ./evals/...`; Gate C artifacts `dry_run:false` N=3; G1 understanding_accuracy **0.8** > B0 **0.0** (untouched) |
| 8 | Board Notes accurate; planner row no product Go | **Pass** — P11-S01-00 Notes claim no product Go; P11-S01-01 Notes match live APIs/tests |

## Focus answers

- Partial GC deletes only **missing-on-disk** same-hash candidates (not live duplicates).
- `removed` is observable (`removed += n` from `gcContentHashOrphans`; test asserts stderr contains `removed`).
- No accidental full-tree GC when `len(args)>0` (`fullTree` set-diff block separate).

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | `TestIndexPartialArgvGCAfterRename` asserts `Contains(stderr, "removed")` not exact `removed 1` | Residual OK — locks allow `removed ≥ 1` |
| low | No dedicated duplicate-hash on-disk sibling CLI test | Residual — code path explicit; isolation covered for distinct-hash sibling |

## Residuals (explicit)

1. **Rename + content edit** so hash ≠ old row — ghost may remain until full-tree `index` or explicit missing-argv delete (planner Residual OK).
2. **`go test ./...`** FAIL only on pre-existing `similar projects/graphify` path space (non-product); product pkgs PASS.
3. File-local only — no dependent-cascade reindex (Law 4 / DR-INCREMENTAL unchanged).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS
Named: TestListFilePathsByContentHash / TestIndexPartialArgvGCAfterRename / TestIndexGCAfterPathRename / TestIndexArgvMissingPathDeletesOnlyThatPath / TestIndexIncrementalIsolation → PASS
```

## Next

**P11-S02-00**
