# P13 / S04 / 01 — Phase 13 VERIFY (import-resolve-honesty closeout)

## Metadata
- id: P13-S04-01
- todo_ids: [P13-S04-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 13 — S01 DF-60 + S02 DF-61/62/63/65 + S03 DF-64/66/67 named regressions + P12 honesty keepers + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S03 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add plan/impact/index MCP dump, invent CLI/analyzer INFERRED (DF-66), invent symbol honesty (DF-67), or scaffold Phase 14 / research ranks 4+ without promotion.

Write durable evidence, then either:

1. **Pass** → declare **Phase 13 VERIFY PASS / import-resolve-honesty green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P13-S04-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-import-path-resolve/REVIEW-NOTES.md), [S02](../scope-02-packet-honesty-residuals/REVIEW-NOTES.md), [S03](../scope-03-provenance-schema/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Optional dogfood: [experiments/ab-import-resolve/](../../../../../experiments/ab-import-resolve/)
- Pattern: Phase 12 VERIFY [`../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P13-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Import-resolve-honesty closeout (S01+S02+S03) — **not** a new `evals/*` planted gate |
| S01 | DF-60 resolve-time Expand; named import-path + Expand/Why subdir/root tests |
| S02 | DF-61/62/63/65; SchemaVersion `0.2`; Law 18 untouched; P12 Budget/cap/stale keepers |
| S03 | DF-64 mig **012** + write/read normalize; DF-66 **wontfix** docs+Law 5; DF-67 Note; compat **12** |
| S01 named | `TestImportPathCandidates_extensionlessThenIndex`; `TestImportPathCandidates_bareModuleExactOnly`; `TestExpandSubdirRelativeImportJS`; `TestExpandParentRelativeImport`; `TestExpandSubdirExtensionlessImport`; `TestExpandRootRelativeImportPositive`; `TestWhySurfacesSubdirRelativeImportProvenance` |
| S02 named | `TestIndexHonestyStaleTotalTruncated`; `TestIndexHonestyPreTrimUniverse`; `TestCandidateCapAdmitUniverseTotal`; `TestContextImportHopEdgeProvenance` |
| S03 named | `TestReplaceFileImportsRejectsGarbageProvenance`; `TestImportProvenanceEmptyWriteAndReadNormalize`; `TestExpandEmptyProvenanceSurfacesExtracted`; `TestImportProvenanceRoundTrip`; `TestAnalyzerImportProvenanceExtracted` |
| P12 keepers | `TestExpandImportEdgeProvenance`; `TestWhySurfacesEdgeProvenance`; `TestContextWhyTraceEdgeProvenance`; `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` (mig ceiling **12**, no 013+) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Residuals (non-blocking) | DF-66 wontfix; DF-67 `symstale/`; TaskContext DF-65 shared-path; graphify space FAIL; CGO0 analyzers FAIL OK; research ranks 4+ deferred; parallel experiments |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; **no** plan/impact/index MCP dump; no new MCP tool menu; **no** provenance MCP/CLI |
| Mig | No new migration from Phase 13 VERIFY |
| Full bar | Product packages in `CGO_ENABLED=1 go test ./...` PASS; known FAIL only `similar projects/graphify` space-in-path (non-product) |
| Successor | **`no successor`** — research ranks 4+ / dogfood stay off-board unless Notes explicitly promote |
| Optional dogfood | `experiments/ab-import-resolve/` prepare + probe — **non-blocking** |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 DF-60 | Subdir `./`/`../` + extensionless + root Expand/Why emit `edge_provenance`; bare modules exact-only; P12 Expand/Why keepers green |
| S02 DF-61/62/63/65 | `stale_total`/`stale_truncated` + MD; pre-trim honesty universe; admit-universe `items_total`; context import-hop `edge_provenance`; SchemaVersion `0.2`; Law 18 untouched; P12 Budget/cap/stale keepers |
| S03 DF-64 | Garbage write reject; empty→EXTRACTED; read normalize; Expand empty surfaces EXTRACTED; mig **012** / compat **12** |
| DF-66 | Docs § Import edge provenance present; Law 5 store-fixture Expand/Why/compiler + round-trip INFERRED green; **no** product CLI/analyzer INFERRED |
| DF-67 | VERIFY-NOTES explicitly records out-of-bar residual (`experiments/_bughunt/post-p12/symstale/`); file-hash `index_honesty` only |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false` |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no Phase 14 / research ranks 4+ without promotion |

### Locked verify commands

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

# Gate H + compat checklist
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

Optional (strong evidence, not substitutes for package PASS):

```bash
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# DF-66: confirm docs/ANALYZER_CONTRIBUTION.md § Import edge provenance; Law 5 fixtures green
# DF-67: record symstale/ residual in VERIFY-NOTES (file-hash only)
# Optional dogfood (non-blocking — not Mode-B Gate C):
#   cd experiments/ab-import-resolve && ./prepare.sh
#   # surface probe: why/context on project/src/app.js should show EXTRACTED after S01
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S03 evidence is **named tests** — not Notes-only
- [ ] MCP remains **nine** tools; no plan/impact/index MCP dump; no provenance MCP/CLI; no new tool menu from Phase 13
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / Neo4j SoT still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration from VERIFY; mig 012 already in S03; compat ceiling **12**
- [ ] Causal `confidence` / `Item.Provenance` not overloaded by `edge_provenance`
- [ ] Law 18 causal STALE not mutated from index drift
- [ ] DF-66: no invented product INFERRED path; DF-67: no invented symbol honesty
- [ ] **No Phase 14 / research ranks 4+ scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 00–12 `done` history; Phase 12 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S04-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 14 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P13-S04-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note DF-66 wontfix + DF-67 `symstale/` + research ranks 4+ stay deferred. Do **not** invent Phase 14. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. Optional: record ab-import-resolve probe result if run (non-blocking). |
| **P13-S04-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 13 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 13 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P13-S04-01a` | implement remediation (full prompt) |
| `P13-S04-01b` | review remediation |
| `P13-S04-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns.

## Board rights
Verify: **status + notes** on `P13-S04-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 13 `done` history. Do **not** mark `P13-S04-02` done. Do **not** scaffold Phase 14 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S03 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 DF-60 import path resolve (required)

| Check | Expect |
|-------|--------|
| Candidates | Extensionless + `index.*` order; bare modules exact-only |
| Expand/Why | Subdir `./`/`../`, extensionless, root → neighbor + `EXTRACTED` / `edge_provenance` |
| P12 keepers | Expand/Why edge provenance still green |

### B. S02 DF-61/62/63/65 packet honesty (required)

| Check | Expect |
|-------|--------|
| DF-61 | `stale_total` + `stale_truncated` + MD when truncated |
| DF-62 | Pre-trim honesty universe (trim-dropped stale ≠ false-fresh null) |
| DF-63 | Admit-universe `items_total` when `candidates_capped` |
| DF-65 | Context import-hop `edge_provenance` via Expand file seeds |
| P12 keepers | Budget totals / candidate cap / index stale banner / WhyTrace provenance |
| Law 18 | Index honesty never sets causal `Provenance.Status` to STALE |
| SchemaVersion | Remains `0.2` |

### C. S03 DF-64/66/67 (required)

| Check | Expect |
|-------|--------|
| DF-64 store | Garbage reject; empty→EXTRACTED; read normalize; round-trip |
| DF-64 Expand | Empty provenance surfaces EXTRACTED |
| DF-64 analyzer | EXTRACTED (+ AMBIGUOUS as covered); no analyzer INFERRED |
| Compat | Ceiling **12**; saw 012; forbid 013+ |
| DF-66 | Docs present; Law 5 fixtures green; **wontfix** disposition confirmed |
| DF-67 | Explicit VERIFY-NOTES residual (`symstale/`); no symbol honesty invented |

### D. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | Product pkgs PASS (graphify space FAIL OK residual) |

### E. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals (DF-66 wontfix; DF-67 `symstale/`; research ranks 4+ still off-board)
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood / research FUTURE may continue off-board)
6. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)
7. Optional: ab-import-resolve probe note (non-blocking)

On FAIL: spawn `P13-S04-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P13-S04-01` only. Do not mark `P13-S04-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DF-66 + DF-67 residuals explicit in Notes
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P13-S04-01`; next `P13-S04-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01–S03 named regressions + P12 keepers
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/`./...`
- [ ] Inspect Gate C `dry_run:false`
- [ ] Confirm DF-66 docs + DF-67 Note
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / new mig
- Scaffolding Phase 14 / research ranks 4+ without promotion
- Re-scoring Gate C
- Treating ab-import-resolve as Gate C
- Closing parallel dogfood experiments
- Rewriting Phase 00–12 history
- Inventing DF-66 product INFERRED or DF-67 symbol honesty
