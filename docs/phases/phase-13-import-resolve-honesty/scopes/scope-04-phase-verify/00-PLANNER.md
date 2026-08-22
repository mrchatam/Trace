# P13-S04-00 — Phase VERIFY / import-resolve-honesty closeout (FINAL)

## Metadata
- id: P13-S04-00
- todo_ids: [P13-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 13 VERIFY evidence table: **S01–S03 named DF-60…67 regressions** + **P12 honesty keepers** + **carry-forward gates** + product `./...`. Decide **DR-HANDOFF** = **`no successor`**. Optional `experiments/ab-import-resolve/` prepare is a dogfood hook only (not Mode-B Gate C). **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 12 VERIFY [`../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/`](../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/)
- Sibling REVIEW-NOTES: [S01](../scope-01-import-path-resolve/REVIEW-NOTES.md), [S02](../scope-02-packet-honesty-residuals/REVIEW-NOTES.md), [S03](../scope-03-provenance-schema/REVIEW-NOTES.md)
- Findings: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Optional dogfood: [experiments/ab-import-resolve/](../../../../../experiments/ab-import-resolve/)
- Research deferrals: [docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — ranks 4+ stay deferred
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — grill only if A1–A7 conflict.

## Depends-on (S01–S03 — landed)

| Scope | Board | Locks imported |
|-------|-------|----------------|
| S01 | **APPROVED** high (P13-S01-02) | DF-60 resolve-time Expand: `import_resolve.go` + Expand wire; named: `TestImportPathCandidates_extensionlessThenIndex`, `TestImportPathCandidates_bareModuleExactOnly`, `TestExpandSubdirRelativeImportJS`, `TestExpandParentRelativeImport`, `TestExpandSubdirExtensionlessImport`, `TestExpandRootRelativeImportPositive`, `TestWhySurfacesSubdirRelativeImportProvenance` (+ P12 Expand/Why keepers) |
| S02 | **APPROVED** high (P13-S02-02) | DF-61/62/63/65: `stale_total`/`stale_truncated` + MD; pre-trim honesty; admit-universe `items_total`; Expand file seeds for `edge_provenance`; SchemaVersion `0.2`. Named: `TestIndexHonestyStaleTotalTruncated`, `TestIndexHonestyPreTrimUniverse`, `TestCandidateCapAdmitUniverseTotal`, `TestContextImportHopEdgeProvenance` (+ P12 keepers) |
| S03 | **APPROVED** high (P13-S03-02) | DF-64 write reject + empty→EXTRACTED + read normalize + mig **012** CHECK; compat ceiling **12**; DF-66 **wontfix** docs+Law 5 fixtures; DF-67 **out-of-bar**. Named: `TestReplaceFileImportsRejectsGarbageProvenance`, `TestImportProvenanceEmptyWriteAndReadNormalize`, `TestExpandEmptyProvenanceSurfacesExtracted` (+ P12 keepers / analyzer EXTRACTED) |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 13 | DF-60…65 + DF-64 (+ DF-66/67 disposition) | Closed by S01–S03 APPROVE high — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | DF-66 documented **wontfix** (confirm docs + Law 5 fixtures); DF-67 symbol-entity staleness out of `index_honesty` (`symstale/`); TaskContext DF-65 shared-path only; first-wins provenance mask theoretical | Forward notes only — **not** a successor phase |
| Research ranks 4–20 / FUTURE | Impact walks, install gates, supersession, etc. | Stay research-only — **not** Phase 14 unless human promotes |
| Parallel dogfood (not board-blocking) | `experiments/ab-import-resolve/` prepare + surface probe; other ladders | Stay in `experiments/` — **not** boarded |
| Known `./...` nit | `similar projects/graphify` space-in-path FAIL; CGO0 analyzers FAIL (tree-sitter) | Pre-existing non-product / expected — VERIFY records **product pkgs PASS** |

**DR-HANDOFF = `no successor`.** No Notes or APPROVE residuals justify a thin Phase 14 scaffold. DF-66/67 are explicit non-product dispositions. Reopen only with explicit human promotion + scaffold (same posture as Phase 10/11/12 historical closes / Phase 13 forward reopen).

## Planner work
1. Lock VERIFY command set (S01–S03 named DF + P12 keepers + carry-forward + product `./...`).
2. Thicken `01-verify.md` evidence table + spawn 01a/b/c + handoff **start**.
3. Thicken `02-scope-review.md` owns DR-HANDOFF **completion** (`no successor`).
4. SCOPE-TODOS + board sync; stamp `DR-HANDOFF.md` ownership.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Phase gate | Phase 13 import-resolve-honesty closeout — named DF-60…67 regressions + P12 honesty keepers — **not** a new planted eval gate |
| S01 home | DF-60 resolve-time import path Expand (`./`/`../`, exts + `index.*`; bare exact-only) |
| S02 home | DF-61/62/63/65 packet/index honesty residuals; SchemaVersion `0.2`; Law 18 causal STALE untouched |
| S03 home | DF-64 mig **012** + write/read normalize; DF-66 **wontfix**; DF-67 residual Note; compat **12** |
| Migration | **None** from VERIFY — mig 012 already landed in S03 |
| S01 named | `TestImportPathCandidates_extensionlessThenIndex`; `TestImportPathCandidates_bareModuleExactOnly`; `TestExpandSubdirRelativeImportJS`; `TestExpandParentRelativeImport`; `TestExpandSubdirExtensionlessImport`; `TestExpandRootRelativeImportPositive`; `TestWhySurfacesSubdirRelativeImportProvenance` |
| S02 named | `TestIndexHonestyStaleTotalTruncated`; `TestIndexHonestyPreTrimUniverse`; `TestCandidateCapAdmitUniverseTotal`; `TestContextImportHopEdgeProvenance` |
| S03 named | `TestReplaceFileImportsRejectsGarbageProvenance`; `TestImportProvenanceEmptyWriteAndReadNormalize`; `TestExpandEmptyProvenanceSurfacesExtracted`; keep `TestImportProvenanceRoundTrip`; `TestAnalyzerImportProvenanceExtracted` |
| P12 keepers | `TestExpandImportEdgeProvenance`; `TestWhySurfacesEdgeProvenance`; `TestContextWhyTraceEdgeProvenance`; `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner` |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` N=3 |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` — **product pkgs PASS**; known FAIL only `similar projects/graphify` space (non-product); CGO0 analyzers FAIL OK residual |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Optional dogfood | `experiments/ab-import-resolve/` `./prepare.sh` + surface probe — **non-blocking** unless Notes elevate |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S04-01 starts** Notes; **S04-02 owns completion**. Do **not** scaffold Phase 14 / research ranks 4+ without explicit promotion |
| Forbidden | Scaffolding research product scopes without human promotion; Mode-B Gate C rewrite; inventing CLI/analyzer INFERRED (DF-66); inventing symbol honesty (DF-67); daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–12 `done` history; claiming Phase 12 historical handoff was wrong |

### Locked verify command set (FINAL)

```bash
# --- S01 DF-60 import path resolve ---
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestImportPathCandidates_|TestExpandSubdir|TestExpandParent|TestExpandRoot|TestWhySurfacesSubdir|TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance'

# --- S02 DF-61/62/63/65 packet honesty residuals (+ P12 keepers) ---
CGO_ENABLED=0 go test ./internal/compiler/... -count=1 -run 'TestIndexHonestyStaleTotalTruncated|TestIndexHonestyPreTrimUniverse|TestCandidateCapAdmitUniverseTotal|TestContextImportHopEdgeProvenance|TestBudgetLoudTotals|TestCandidateCapSetsTruncated|TestIndexStaleBanner|TestContextWhyTraceEdgeProvenance'

# --- S03 DF-64 (+ DF-66 Law 5 fixtures / P12 keepers) ---
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestReplaceFileImportsRejectsGarbageProvenance|TestImportProvenanceEmptyWriteAndReadNormalize|TestImportProvenanceRoundTrip'
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpandEmptyProvenanceSurfacesExtracted|TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance'
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestAnalyzerImportProvenanceExtracted'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat (compat also covers mig 012 ceiling)
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Supporting surfaces (optional strong evidence)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... ./internal/retrieval/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1

# Full regression bar (product pkgs; graphify space FAIL is known residual)
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, **not** substitutes for package PASS / not Mode-B Gate C):

```bash
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# DF-66: docs/ANALYZER_CONTRIBUTION.md § Import edge provenance present; Law 5 fixtures green (no CLI/analyzer INFERRED)
# DF-67: VERIFY-NOTES must name experiments/_bughunt/post-p12/symstale/ residual; file-hash honesty only
# Optional dogfood (non-blocking):
#   cd experiments/ab-import-resolve && ./prepare.sh
#   # surface probe: why/context on project/src/app.js shows EXTRACTED after S01
# Research ranks 4+: stay off-board unless Notes explicitly promote
```

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked (`no successor`)
- [x] SCOPE-TODOS + board Notes; next `P13-S04-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S04-01)
- Product Go / new MCP tools / daemon / mig
- Scaffolding Phase 14 / research ranks 4+ without explicit promotion
- Closing parallel dogfood experiments
- Claiming Phase 12 historical handoff was wrong
