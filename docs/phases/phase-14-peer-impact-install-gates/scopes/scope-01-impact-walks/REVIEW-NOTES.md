# P14-S01-02 — Scope review notes (Impact walks)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist evidence

| # | Result | Evidence |
|---|--------|----------|
| 1 | PASS | One `seen` + seedKeys exclusion in `impact_walk.go`; `TestImpactWalkMultiSeedExcludeSeeds` |
| 2 | PASS | `allowContainsOut` false after contains-UP; `TestImpactWalkContainsAsymmetryNoSiblings` |
| 3 | PASS | Neighbors = contains-OUT (gated) + `listIncomingImporters` only; Expand still uses outgoing `ListImportsByPath`; `TestImpactWalkIncomingImportHop` |
| 4 | PASS | Cap 64 + `blast_total`/`blast_kept`/`truncated`; `TestImpactWalkLoudTruncation` |
| 5 | PASS | `hop_risk = float64(hop)`; `TestImpactWalkHopRiskIncreases` |
| 6 | PASS | `TestPlantedImpactConflictsGateFPrelim` green under CGO0/CGO1/product suites; CLI still `finding`/`alternative`/`report`/`walk` |
| 7 | PASS | Expand `case "symbol"` still parent-file; ImpactWalk is separate file/API |
| 8 | PASS | Thin `cmdImpactWalk`; MCP still 9 tools; `TestImportBoundaryMCPNoPlanImpactIndexTools` forbids `trace_impact` |
| 9 | PASS | No `internal/impact`; schema still ends at `012_*` (no new mig); no daemon/HTTP/embeddings |
| 10 | PASS | Fresh verify (below); Gate C `metrics-*.json` still `dry_run: false` |
| 11 | PASS | P14-S01-00 Notes docs-only; S01-01 Notes match live walk + tests |
| 12 | PASS | Walk returns structural JSON only — no finding plant / plan mutation |

## Independent verify (fresh)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/store/... ./evals/impact/... ./evals/honesty/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/impact/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS
```

## Findings

### blocker / high
None.

### medium
1. **`allowContainsOut` late upgrade may skip re-expand** (`impact_walk.go` ~110–125) — If a file is first reached via contains-UP (`allowContainsOut=false`) and later rediscovered via import with `allowContainsOut=true` at a **greater** hop, the flag upgrades but the node is not re-enqueued (`upgraded && existing.hop == d` fails). Same-hop multi-seed upgrade works; primary symbol-seed asymmetry test covers the locked path. Residual incomplete contains-OUT only on that asymmetric multi-path edge. No spawn (not high). Optional future harden / VERIFY spot-check.

### low
2. **`listIncomingImporters` full `ListImportEdges` scan** per file expansion — fine for S01 bar; large graphs may want indexed reverse lookup later (planner already allowed store helper).
3. **No dedicated negative test** that outgoing import targets of a seed are absent — code review confirms; IncomingImportHop only proves incoming.

### nit
4. Duplicate identical seeds are appended twice to `Seeds` (dedupe only on `seen`/frontier).
5. `ImpactWalk` ignores `ctx` (`_ = ctx`) — no cancel propagation.

## Board / next
- Mark **P14-S01-02** done; next runnable **P14-S02-00**.
- No rewrite of P14-S01-00/01 prompt history.
