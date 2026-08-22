# P13-S01-02 — REVIEW-NOTES (import path resolve)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | Subdir `./` / `../` → importer-dir join + Expand neighbor + `edge_provenance` | **Pass** — `import_resolve.go` `importPathCandidates`/`resolveImportedFile`; Expand call site `expand.go`; fixtures `TestExpandSubdirRelativeImportJS`, `TestExpandParentRelativeImport` (reviewer inline), Why `TestWhySurfacesSubdirRelativeImportProvenance` |
| 2 | Extensionless `./util` tries locked exts (+ `index.*`) and hits indexed file | **Pass** — candidate order/dedupe `TestImportPathCandidates_extensionlessThenIndex`; Expand `TestExpandSubdirExtensionlessImport` → `src/util.ts` + EXTRACTED |
| 3 | Root `./util.js`→`util.js` + `edge_provenance` | **Pass** — `TestExpandRootRelativeImportPositive` |
| 4 | P12 provenance tests green; causal confidence / `Item.Provenance` untouched | **Pass** — `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance`; Hit only sets `EdgeProvenance` from `imp.Provenance` (no causal Status/confidence mutation) |
| 5 | Retrieval resolve only — no analyzer rewrite / path-align / new mig | **Pass** — product delta confined to `internal/retrieval` (`import_resolve.go` + Expand wire + tests); `NormalizePath` unchanged; no `012_*`; no path-align hooks |
| 6 | Bare modules Expand-skipped (no fake module resolver) | **Pass** — `TestImportPathCandidates_bareModuleExactOnly` (`fmt`, `github.com/…`); non-relative → exact only |
| 7 | Forbidden architecture absent; G19 intact | **Pass** — no full-rebuild indexer, daemon/HTTP, new MCP tools, or fake call edges in S01 delta |
| 8 | Carry-forward + Gate C `dry_run:false`; Notes accurate; no implementer board spawn | **Pass** — suites below; Gate C metrics `dry_run:false`; implementer Notes match live; board had no implementer-spawned rows |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Checklist `../` had no Expand fixture before review | **Fixed inline** — `TestExpandParentRelativeImport` |
| low | Extensionless relatives append exts on **exact** basename before joined path (locked order) — rare false hit if root `util.js` exists beside `src/util.ts` | Residual — matches FINAL candidate order; not renegotiated |
| nit | `index.*` covered in candidate unit test only (no Expand store fixture for `./dir` → `dir/index.js`) | Residual OK |

## Residuals (explicit)

1. Locked candidate order can prefer root `basename+ext` over importer-dir file for extensionless `./util` when both exist.
2. No Expand integration fixture for `index.*` (candidate list asserts presence).
3. Optional dogfood `experiments/ab-import-resolve/` not re-run (A8 — not board blocker).
4. Bare `node:fs` not named in unit test (same exact-only path as `fmt`).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
  → PASS
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpand…|TestWhy…|TestImportPath…'
  → PASS (incl. TestExpandParentRelativeImport after inline)
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
  → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run=false
```

## Next

**P13-S02-00** (do not start from this row)
