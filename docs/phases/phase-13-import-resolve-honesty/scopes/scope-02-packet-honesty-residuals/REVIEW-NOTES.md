# P13-S02-02 — REVIEW-NOTES (packet honesty residuals)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | DF-61: `stale_total` + `stale_truncated`; cap ≤8; MD total when truncated | **Pass** — `IndexHonesty` fields in `packet.go`; sort-then-cap in `index_honesty.go`; MD `stale_total=N` when truncated; `TestIndexHonestyStaleTotalTruncated` |
| 2 | DF-62: honesty universe = pre-trim file items; trim-dropped stale ≠ false-fresh null | **Pass** — `buildIndexHonesty(c.store, items)` before keep; `TestIndexHonestyPreTrimUniverse` |
| 3 | DF-63: `candidates_capped` → `items_total` ≥ L1 admit universe (not post-cap ≤64); MD `items=k/t` | **Pass** — `admitUniverse` via `layer1AdmitKey` over full candidates; `TestCandidateCapAdmitUniverseTotal` |
| 4 | DF-65: context import-hop `edge_provenance` via Expand on file seeds; no compiler path-join | **Pass** — `compileAtDepth` Expand `fileSeeds` depth 1; resolve stays in `Retriever.Expand`/`resolveImportedFile`; `TestContextImportHopEdgeProvenance` |
| 5 | P12 keepers; Law 18; SchemaVersion `0.2`; false-fresh on I/O miss | **Pass** — named keepers green; honesty never sets causal `Provenance.Status=STALE`; `SchemaVersion = "0.2"`; omit on missing row/disk/I/O |
| 6 | No mig / analyzer rewrite / path-align / new MCP; G19 | **Pass** — delta in `internal/compiler` (+ tests); zero MCP/CLI adapter edits; no path-align / analyzer / mig |
| 7 | S01 Expand/Why import tests still green | **Pass** — `TestExpand.*Import` / `TestWhy.*Import` |
| 8 | Carry-forward + Gate C `dry_run:false`; dry-run ≠ C/H/checklist | **Pass** — suites below; Gate C metrics unchanged `dry_run:false` |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | DF-65 named test exercises `ExpandContext` depth 2 only (not `TaskContext` depth 1) | Residual — both share `compileAtDepth` + file-seed Expand |
| nit | Layer-1 admission is first-wins after `sortHits`; a prior no-provenance file hit could mask a later Expand hop’s `edge_provenance` | Residual — DF-65 fixture util appears only via file Expand; not observed failing |
| nit | DF-61 asserts `stale_total > 8` not exact full count | Residual OK vs lock |

## Residuals (explicit)

1. No dedicated `TaskContext` (depth 1) DF-65 fixture — shared compile path only.
2. First-wins item admission can theoretically drop later `edge_provenance` on duplicate entity IDs.
3. Honesty still bounded to pre-trim / MaxCandidateHits pipeline (whole-index / symbol staleness → S03 DF-67).
4. Optional dogfood / bughunt re-run not required for board close.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1
  → PASS
CGO_ENABLED=0 go test ./internal/compiler/... -count=1 -run 'TestBudgetLoudTotals|TestCandidateCapSetsTruncated|TestIndexStaleBanner|TestContextWhyTraceEdgeProvenance|TestIndexHonesty|TestCandidateCapAdmit|TestContextImportHop'
  → PASS
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpand.*[Ii]mport|TestWhy.*[Ii]mport'
  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
  → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run=false intact
```

## Next

**P13-S03-00** (do not start from this row)
