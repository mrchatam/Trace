# P12-S03-00 — Phase VERIFY / peer-honesty closeout (FINAL)

## Metadata
- id: P12-S03-00
- todo_ids: [P12-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 12 VERIFY evidence table: **S01 + S02 named regressions** + **carry-forward gates** + product `./...`. Decide **DR-HANDOFF** = **`no successor`** unless Notes promote a further peer-technique slice. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 11 S08 [`../../../phase-11-residual-surfaces/scopes/scope-08-phase-verify/`](../../../phase-11-residual-surfaces/scopes/scope-08-phase-verify/)
- Sibling REVIEW-NOTES: [S01](../scope-01-edge-provenance/REVIEW-NOTES.md), [S02](../scope-02-packet-honesty/REVIEW-NOTES.md)
- Research: [docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — ranks 4+ stay deferred
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — grill only if A1–A7 conflict.

## Depends-on (S01–S02 — landed)

| Scope | Board | Locks imported |
|-------|-------|----------------|
| S01 | **APPROVED** high (P12-S01-02) | Mig `011_import_edge_provenance`; `imports.provenance`; JSON `edge_provenance`; named: `TestImportProvenanceRoundTrip`, `TestAnalyzerImportProvenanceExtracted`, `TestExpandImportEdgeProvenance`, `TestWhySurfacesEdgeProvenance`, `TestContextWhyTraceEdgeProvenance` |
| S02 | **APPROVED** high (P12-S02-02) | Live: `SchemaVersion` `0.2`; Budget totals/cap + loud MD; `index_honesty` false-fresh + **sort-then-cap 8**; Law 18 untouched. Named: `TestBudgetLoudTotals`, `TestCandidateCapSetsTruncated`, `TestIndexStaleBanner` (+ S01 `TestContextWhyTraceEdgeProvenance`). Residual: stale test asserts sorted/`len≤8`/primary path, not exact lex-first-8 set |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 12 | S01 edge provenance + S02 packet honesty | Closed by S01–S02 APPROVE high — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | No enum CHECK on provenance; synthetic context JSON fixture; stale-banner test does not pin exact lex-first-8 set; symbol-entity staleness deferred | Forward notes only — **not** a successor phase |
| Research ranks 4–20 / FUTURE S03–S05 | Impact walks, install gates, supersession, etc. | Stay research-only — **not** Phase 13 unless human promotes |
| Parallel dogfood (not board-blocking) | ladders / Cursor MCP reload / experiments | Stay in `experiments/` — **not** boarded |
| Known `./...` nit | `similar projects/graphify` space-in-path FAIL; CGO0 analyzers FAIL (tree-sitter) | Pre-existing non-product / expected — VERIFY records **product pkgs PASS** |

**DR-HANDOFF = `no successor`.** No Notes or APPROVE residuals justify promoting research ranks 4+. Reopen only with explicit human promotion + scaffold (same posture as Phase 10/11 historical close / Phase 12 forward reopen).

## Planner work
1. Lock VERIFY command set (S01–S02 named tests + carry-forward + product `./...`).
2. Thicken `01-verify.md` evidence table + spawn 01a/b/c + handoff **start**.
3. Thicken `02-scope-review.md` owns DR-HANDOFF **completion** (`no successor`).
4. SCOPE-TODOS + board sync; stamp `DR-HANDOFF.md` ownership.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Phase gate | Phase 12 peer-honesty-surfaces closeout — **not** a new planted eval gate |
| S01 home | Edge provenance on **import** edges: mig 011; `imports.provenance`; JSON/MD `edge_provenance` |
| S02 home | Packet honesty: SchemaVersion `0.2`; Budget totals/cap; `index_honesty` sort-then-cap 8; Law 18 causal STALE untouched |
| Migration | **None** from VERIFY — mig 011 already landed in S01 |
| S01 named | `TestImportProvenanceRoundTrip` (store); `TestAnalyzerImportProvenanceExtracted` (analyzers, CGO1); `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` (retrieval); `TestContextWhyTraceEdgeProvenance` (compiler) |
| S02 named | `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner` (compiler) — keep S01 `TestContextWhyTraceEdgeProvenance` green |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` N=3; Phase 11 DF surfaces stay green via product pkgs / supporting suites |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` — **product pkgs PASS**; known FAIL only `similar projects/graphify` space (non-product); CGO0 analyzers FAIL OK residual |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S03-01 starts** Notes; **S03-02 owns completion**. Do **not** scaffold Phase 13 / research S03–S05 without explicit promotion |
| Forbidden | Scaffolding research impact/install/supersession product scopes without human promotion; Mode-B Gate C rewrite; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–11 `done` history; claiming Phase 11 historical handoff was wrong |

### Locked verify command set (FINAL)

```bash
# --- S01 edge provenance ---
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestImportProvenanceRoundTrip'
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance'
CGO_ENABLED=0 go test ./internal/compiler/... -count=1 -run 'TestContextWhyTraceEdgeProvenance'
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestAnalyzerImportProvenanceExtracted'

# --- S02 packet honesty (+ S01 compiler regression) ---
CGO_ENABLED=0 go test ./internal/compiler/... -count=1 -run 'TestBudgetLoudTotals|TestCandidateCapSetsTruncated|TestIndexStaleBanner|TestContextWhyTraceEdgeProvenance'

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat checklist
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Supporting surfaces (optional strong evidence; Phase 11 DF green via these + full bar)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... ./internal/retrieval/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1

# Full regression bar (product pkgs; graphify space FAIL is known residual)
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# S02 residual OK: TestIndexStaleBanner need not pin exact lex-first-8 membership
# Research ranks 4+: stay off-board unless Notes explicitly promote
```

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked
- [x] SCOPE-TODOS + board Notes; next `P12-S03-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S03-01)
- Product Go / new MCP tools / daemon / mig
- Scaffolding Phase 13 / research S03–S05 without explicit promotion
- Closing parallel dogfood experiments
- Claiming Phase 11 historical handoff was wrong
