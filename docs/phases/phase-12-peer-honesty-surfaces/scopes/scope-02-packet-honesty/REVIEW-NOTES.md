# P12-S02-02 — REVIEW-NOTES (packet honesty)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | SchemaVersion `0.2`; Budget `items_total`/`items_kept`/`candidates_capped`; MD `items=kept/total` | **Pass** — `packet.go` `SchemaVersion = "0.2"`; Budget JSON tags; `RenderMarkdown` budget line; `TestBudgetLoudTotals` |
| 2 | `truncated=true` when kept&lt;total **or** MaxCandidateHits capped — no silent caps | **Pass** — `compiler.go` `compileAtDepth`: `itemsTotal` pre-`trimToBudget`, then `if itemsKept < itemsTotal \|\| candidatesCapped { truncated = true }`; Layer-1 peek sets `candidatesCapped`; `TestCandidateCapSetsTruncated` |
| 3 | `index_honesty.stale_paths` on hash mismatch; false-fresh on errors; MD banner; sort-then-cap 8 | **Pass** — `buildIndexHonesty` unique→`sort.Strings`→`[:8]`; `continue` on missing row/disk/I/O; MD banner before Items; `TestIndexStaleBanner` |
| 4 | Causal `Provenance.Status` / Law 18 STALE not mutated from index drift | **Pass** — honesty path never sets Status; test asserts no item `Provenance.Status == "STALE"` after disk mutate |
| 5 | S01 intact — `edge_provenance` pass-through; not written into causal confidence / `Item.Provenance` | **Pass** — `compiler.go` copies `EdgeProvenance` on Item/WhyTrace; `TestContextWhyTraceEdgeProvenance` PASS |
| 6 | G19 — thin adapters only; no new MCP tool menu | **Pass** — no `index_honesty` / budget honesty forks under `cmd/` or `internal/mcp`; library marshal only |
| 7 | No forbidden architecture | **Pass** — emission-time sha256 vs `files.content_hash`; no size/mtime migration; no daemon watcher / embeddings / Neo4j / full-rebuild |
| 8 | Carry-forward + Gate C `dry_run:false`; named tests PASS; Notes accurate; planner no product Go on retry | **Pass** — suites below; Gate C metrics `dry_run:false` N=3; implementer gap-finish (sort-then-cap + named tests) matches FINAL inventory; P12-S02-00 Notes claim no product Go (types pre-existed) |

## Focus answers

- **Totals:** `items_total = len(items)` after Layer-0/1 assemble, **before** `trimToBudget`; `items_kept = len(kept)` after trim (`compiler.go` `compileAtDepth`).
- **False-fresh:** `index_honesty.go` `continue` when `GetFileByID` fails / empty path or hash / `sha256FileHex` errors; omit `index_honesty` when no stale paths.
- **S01 EdgeProvenance:** still wired on Layer-1 Item assign and WhyTrace copy; MD/JSON fields unchanged.
- **Implementer posture:** gap finish only — did **not** re-ship SchemaVersion/`Budget`/`IndexHonesty` types; product delta is **sort-then-cap** in `buildIndexHonesty` plus three named tests.

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | `TestIndexStaleBanner` asserts sorted + `len≤8` + primary path present, not the exact lex-first-8 set if &gt;8 kept stale files | Residual OK — production path is sort-then-cap; VERIFY may optionally strengthen assert |

## Residuals (explicit)

1. Named stale-banner test does not pin exact lex-first-8 membership (code does).
2. Symbol-entity staleness still out of bar (FINAL deferral).
3. Product `./...` known FAIL only outside product pkgs (`similar projects/graphify` space) — not re-run here; product `./cmd/... ./internal/... ./evals/...` PASS.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1
  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
  → PASS
Named: TestBudgetLoudTotals / TestCandidateCapSetsTruncated / TestIndexStaleBanner /
  TestContextWhyTraceEdgeProvenance → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run=false N=3 each
```

## Next

**P12-S03-00**
