# P13-S03-02 — Scope review notes (Provenance schema)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | DF-64 write reject garbage; empty→EXTRACTED; read normalize | **PASS** | `validateImportProvenance` / `normalizeImportProvenance` in `internal/store/file_graph.go`; `TestReplaceFileImportsRejectsGarbageProvenance`; `TestImportProvenanceEmptyWriteAndReadNormalize`; `TestImportProvenanceRoundTrip` |
| 2 | DF-64 mig **012** CHECK; heal; embed ceiling **12** (no 013+) | **PASS** | `internal/store/schema/012_import_provenance_enum.sql` rebuild+CHECK+heal; embed via `schema/*.sql`; `evals/compat` EmbedExpected==12, saw012, forbid 013+ |
| 3 | Empty cannot omitempty-hide after normalize | **PASS** | `TestExpandEmptyProvenanceSurfacesExtracted` — Expand `EdgeProvenance` == `EXTRACTED` |
| 4 | DF-66 **wontfix** + Law 5 fixtures green | **PASS** | `docs/ANALYZER_CONTRIBUTION.md` § Import edge provenance; analyzers EXTRACTED/AMBIGUOUS only; `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` / `TestContextWhyTraceEdgeProvenance` / round-trip INFERRED green |
| 5 | DF-67 **no** symbol honesty; residual explicit | **PASS** | No symbol-hash / symbol-staleness in `index_honesty.go` (file-hash only); residual for S04 VERIFY (`symstale/`) |
| 6 | P12 keepers + SchemaVersion `0.2`; G19 no provenance MCP/CLI | **PASS** | Named P12 tests green; `compiler.SchemaVersion = "0.2"`; no import-provenance CLI/MCP (causal `source_type` flags unrelated) |
| 7 | Prefer zero analyzer/retrieval/compiler churn | **PASS** | Store + mig + compat + docs primary; retrieval = new Expand empty-provenance test only; analyzers/compiler logic unchanged for enum |
| 8 | Carry-forward + Gate C `dry_run:false` | **PASS** | CGO0 store/retrieval/compiler/honesty PASS; CGO1 cmd/trace+analyzers+p0x+x0+honesty+compat PASS; product `./cmd|internal|evals` PASS; Gate C metrics `dry_run:false` N=3 mean G1 **0.800** > B0 **0.000** |

## Independent verify (re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS
```

Named focus: store provenance + Expand/Why/compiler edge provenance + `TestAnalyzerImportProvenanceExtracted` + compat → PASS.

## Findings

None at blocker/high/medium. No spawns.

### Residuals (explicit — never silent)

| ID | Disposition | Carry to |
|----|-------------|----------|
| **DF-66** | Documented **wontfix** — no product analyzer/CLI INFERRED setter; Law 5 via store fixture + docs only | S04 VERIFY Notes (confirm docs + fixtures) |
| **DF-67** | **Out of honesty bar** — file-hash `index_honesty` only; symbol rows can linger after disk mutate (`experiments/_bughunt/post-p12/symstale/`) | S04 VERIFY-NOTES must record |
| low | Garbage reject tested from empty import set; prior-row restore relies on tx rollback (untested restore path) | optional future; not blocking |
| low | `normalizeImportProvenance` only maps empty (not unknown); post-012 CHECK + migrate heal close the gap | none |

## Dry-run ≠ gates

Dry-run ≠ Gate C / ≠ Gate H / ≠ checklist. Gate C artifacts remain Mode-B `dry_run:false`.

## Next

**P13-S04-00** (do not start from this review).
