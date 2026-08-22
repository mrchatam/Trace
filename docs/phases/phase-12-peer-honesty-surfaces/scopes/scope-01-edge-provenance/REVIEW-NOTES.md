# P12-S01-02 — REVIEW-NOTES (edge provenance)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | Mig **011** + `imports.provenance` persists `EXTRACTED\|INFERRED\|AMBIGUOUS`; default EXTRACTED | **Pass** — `internal/store/schema/011_import_edge_provenance.sql` `NOT NULL DEFAULT 'EXTRACTED'`; `ReplaceFileImports` empty→EXTRACTED; `TestImportProvenanceRoundTrip` |
| 2 | Analyzers: AST → EXTRACTED; wildcard → AMBIGUOUS; no fake call edges / no analyzer INFERRED | **Pass** — Go/JS set EXTRACTED; py `*` / empty path → AMBIGUOUS; zero `ImportProvenanceInferred` / calls/usages under `internal/analyzers` |
| 3 | Why + context surface JSON/MD `edge_provenance`; INFERRED fixture not silent-as-EXTRACTED | **Pass** — Expand copies `imp.Provenance`; Why/compiler pass-through; MD ``edge_provenance: `…` ``; named Expand/Why/compiler tests with store-fixture INFERRED |
| 4 | Causal `confidence` / `Item.Provenance` not overloaded | **Pass** — separate `EdgeProvenance` fields; compiler test asserts structural hop leaves `Item.Provenance` unset; mig comment + A4 |
| 5 | G19 — no domain fork in CLI/MCP | **Pass** — no `edge_provenance` / `ImportProvenance` under `cmd/` or `internal/mcp` |
| 6 | No forbidden architecture | **Pass** — column on `imports` only; no daemon/HTTP/embeddings/Neo4j/full-rebuild/call tables |
| 7 | Carry-forward + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 suites + product pkgs; Gate C artifacts `dry_run:false` N=3 each |
| 8 | Board Notes accurate; planner no product Go | **Pass** — P12-S01-00 Notes match FINAL locks / no Go; P12-S01-01 Notes match live APIs/tests |

## Focus answers

- Enum home is **`imports.provenance`**, not causal `confidence`.
- Analyzers never emit `INFERRED`; surfacing proven via store fixture → Expand/Why/compiler.
- MCP/CLI parity is library-owned JSON only (G19).

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | `ReplaceFileImports` does not reject unknown provenance strings | Residual OK — locks require round-trip of three values + empty default, not CHECK |
| low | `TestContextWhyTraceEdgeProvenance` hand-builds Packet for JSON/MD; live `TaskContext` only smoke-checks TaskID | Residual OK — Expand/Why + compiler pass-through covered separately |

## Residuals (explicit)

1. Unknown provenance strings accepted at write (no enum CHECK).
2. No end-to-end TaskContext path that forces a task→file→import hop for `edge_provenance` on live Items (library Expand→Item path still wired).
3. `go test ./...` FAIL only on pre-existing `similar projects/graphify` path space (non-product); product `./cmd/... ./internal/... ./evals/...` PASS.
4. CGO0 cannot build `./internal/analyzers/...` (tree-sitter) — expected; CGO1 analyzers PASS.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
  → PASS (analyzers build-failed under CGO0 — expected)
CGO_ENABLED=1 go test ./internal/analyzers/... ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 → PASS
Named: TestImportProvenanceRoundTrip / TestAnalyzerImportProvenanceExtracted /
  TestExpandImportEdgeProvenance / TestWhySurfacesEdgeProvenance /
  TestContextWhyTraceEdgeProvenance → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run=false N=3
```

## Next

**P12-S02-00**
