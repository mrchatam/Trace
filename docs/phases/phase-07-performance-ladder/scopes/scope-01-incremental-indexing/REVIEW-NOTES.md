# P07-S01-02 — Scope review notes (2026-08-16)

Independent review of S01 against `00-PLANNER.md` / `01-incremental-indexing.md` + TODO Notes for `P07-S01-01`. Fresh session; claims verified in-repo (no implementer assumptions reused).

## Plan (executed)

1. Diff claims vs `cmd/trace/index.go` (`isT0SkipDir` / `isT0SkipPath` / `walkIndexable` / explicit argv gate) + `cli_test.go` proofs
2. Confirm T0 dirs/suffixes/path-segment locks, walk order T0→DetectLanguage→T0 file/path→gitignore, file-local incremental (no full rebuild / no `internal/indexer`)
3. Confirm **no** `011_*`; no Gate H / `evals/perf` threshold invent; optional seed absent
4. Re-run cmd/trace + analyzers + honesty / Gate G / Gate E / Gate F / capability / p0x / x0 / `./...`
5. Gate C `docs/verification/gate-c-x0/` `dry_run:false` N=3 G1≈0.8 / B0=0.0 intact
6. Severity-tag findings; spawn only for blocker/high (none)
7. Write these notes; mark board; unlock **P07-S02-00**; light-note S02/S03

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**. Next board row: **P07-S02-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| T0 dir basenames (`.git`/`.trace`/`node_modules`/`vendor`/`__pycache__`/`.venv`/`venv`/`dist`/`.next`/`target`/`coverage`); no `build`/`bin` | Pass (`t0SkipDirs` + `isT0SkipDir`) |
| T0 file suffixes `.min.js`/`.min.mjs`/`.min.cjs` | Pass (`isT0SkipPath` basename HasSuffix) |
| Path-segment rule on walk **and** explicit argv | Pass (`walkIndexable` step 3; `cmdIndex` `skipped++` before `indexOne`) |
| Walk order: T0 SkipDir → DetectLanguage → T0 file/path → gitignore | Pass (`walkIndexable` comments + control flow) |
| File-local incremental; no cascade / no `internal/indexer` / no analyzer rewrite | Pass (`indexOne` → `IndexFile`/`IndexFileAtRev`; `internal/indexer` absent; `TestIndexIncrementalIsolation` + analyzers `TestIncrementalIsolation`) |
| Explicit T0 argv → skip not hard fail | Pass (`skipped++`; `TestIndexSkipsExplicitT0Path`) |
| Tests: isolation + `TestWalkIndexableT0AlwaysSkip` + `TestIndexSkipsExplicitT0Path` | Pass |
| **No** `011_*` | Pass (schema through `010_capability_surface.sql` only) |
| No Gate H claim; no `evals/perf` threshold harness | Pass (no `evals/perf/`; only optional `t.Logf` timing in T0 walk test) |
| Honesty A/B/C + Gate G + Gate E + Gate F + capability ablation | Pass |
| p0x 7/7; x0; `CGO_ENABLED=1` `./...` | Pass |
| Gate C `dry_run:false`; N=3; G1 understanding_accuracy 0.8 / B0 0.0 | Pass |
| No daemon/HTTP/embeddings as primary | Pass (CLI walk-only change surface) |

## Re-verification commands (2026-08-16)

```text
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... -count=1
# ok

CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
# ok

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
# ok

CGO_ENABLED=1 go test ./... -count=1
# ok
```

Named proofs: `TestIndexIncrementalIsolation`, `TestWalkIndexableT0AlwaysSkip`, `TestIndexSkipsExplicitT0Path`; analyzers `TestIncrementalIsolation`.

## Findings

### Blocker / high

None.

### Medium

None requiring spawn.

### Low / nit (residual)

- **Explicit `.min.*` argv:** `TestIndexSkipsExplicitT0Path` plants `node_modules/x.js` only; minified suffixes are covered on the walk plant (`foo.min.js` / `.min.mjs` / `.min.cjs`) and share `isT0SkipPath`, but there is no dedicated `trace index foo.min.js` CLI assert.
- **Optional Gate H seed skipped:** No `evals/perf/fixtures` tree; only `t.Logf` wall time in `TestWalkIndexableT0AlwaysSkip`. Allowed by lock; S03 must not assume planted ladders exist yet.
- **Project dir named a T0 basename:** If the project root’s basename were e.g. `coverage`, WalkDir’s root dir entry could `SkipDir` the whole tree — pathological edge; not worth changing this row.

## Spawns

None.

## S02 note (upcoming)

`P07-S02-00` already locks “consume S01 T0 surface; do not regress.” Live surface to preserve: `cmd/trace` `isT0SkipDir` / `isT0SkipPath` / `walkIndexable` order (T0 → lang → T0 file/path → gitignore) + explicit argv skip. Extending `DetectLanguage` must not bypass T0 filters. No substantive prompt rewrite required — live-confirm stamp applied on S02 planner Depends.

## S03 note (upcoming)

S01 did **not** seed `evals/perf` fixtures. Gate H planner should treat ladders as still absent (measurement = CLI T0 + isolation tests + optional later seeds from S02/S03). Do not invent thresholds from S01 `t.Logf` alone.
