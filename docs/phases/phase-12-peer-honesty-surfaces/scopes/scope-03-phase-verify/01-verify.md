# P12 / S03 / 01 — Phase 12 VERIFY (peer-honesty closeout)

## Metadata
- id: P12-S03-01
- todo_ids: [P12-S03-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 12 — S01 edge-provenance + S02 packet-honesty named regressions + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S02 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add plan/impact/index MCP dump, or scaffold Phase 13 / research S03–S05 without promotion.

Write durable evidence, then either:

1. **Pass** → declare **Phase 12 VERIFY PASS / peer-honesty surfaces green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P12-S03-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-edge-provenance/REVIEW-NOTES.md), [S02](../scope-02-packet-honesty/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 11 VERIFY [`../../../phase-11-residual-surfaces/scopes/scope-08-phase-verify/01-verify.md`](../../../phase-11-residual-surfaces/scopes/scope-08-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P12-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Peer-honesty-surfaces closeout (S01+S02) — **not** a new `evals/*` planted gate |
| S01 | Mig 011 `imports.provenance`; JSON/MD `edge_provenance`; named store/analyzer/retrieval/compiler tests |
| S02 | SchemaVersion `0.2`; Budget totals/cap + loud MD; `index_honesty` false-fresh + sort-then-cap 8; Law 18 untouched |
| S01 named | `TestImportProvenanceRoundTrip`; `TestAnalyzerImportProvenanceExtracted`; `TestExpandImportEdgeProvenance`; `TestWhySurfacesEdgeProvenance`; `TestContextWhyTraceEdgeProvenance` |
| S02 named | `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner` (+ S01 compiler regression) |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` (Path C = FAIL→UNCERTAIN then PASS+`AllowOperatorDone`) |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` (hatch Escape-1 retained) |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Residuals (non-blocking) | No provenance enum CHECK; synthetic context JSON fixture; stale test does not pin exact lex-first-8 set; symbol-entity staleness out; graphify space FAIL; CGO0 analyzers FAIL OK; research ranks 4+ deferred; parallel experiments |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; **no** plan/impact/index MCP dump; no new MCP tool menu |
| Mig | No new migration from Phase 12 VERIFY |
| Full bar | Product packages in `CGO_ENABLED=1 go test ./...` PASS; known FAIL only `similar projects/graphify` space-in-path (non-product) |
| Successor | **`no successor`** — research ranks 4+ / dogfood stay off-board unless Notes explicitly promote |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 | Store round-trip; analyzer EXTRACTED (+ AMBIGUOUS as covered); Expand/Why `edge_provenance`; compiler WhyTrace/Item + MD; INFERRED surfaces via store fixture (not analyzer) |
| S02 | Loud Budget totals (`items_total`/`items_kept`/`candidates_capped`); silent-cap ⇒ `truncated`; `index_honesty.stale_paths` + false-fresh + MD banner; sort-then-cap 8; Law 18 causal STALE untouched |
| S01∩S02 | `TestContextWhyTraceEdgeProvenance` still green after packet honesty |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false`; Phase 11 DF surfaces via product pkgs |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no Phase 13 / research S03–S05 without promotion |

### Locked verify commands

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
# Spot-check: Law 18 — index_honesty path never sets Item.Provenance.Status = STALE
# S02 residual OK: TestIndexStaleBanner need not pin exact lex-first-8 membership
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S02 evidence is **named tests** — not Notes-only
- [ ] MCP remains **nine** tools; no plan/impact/index MCP dump; no new tool menu from Phase 12
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / Neo4j SoT still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration from VERIFY; mig 011 already in S01
- [ ] Causal `confidence` / `Item.Provenance` not overloaded by `edge_provenance`
- [ ] Law 18 causal STALE not mutated from index drift
- [ ] **No Phase 13 / research S03–S05 scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 00–11 `done` history; Phase 11 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 13 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P12-S03-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note residuals + research ranks 4+ stay deferred. Do **not** invent Phase 13. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P12-S03-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 12 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 12 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P12-S03-01a` | implement remediation (full prompt) |
| `P12-S03-01b` | review remediation |
| `P12-S03-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns.

## Board rights
Verify: **status + notes** on `P12-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 12 `done` history. Do **not** mark `P12-S03-02` done. Do **not** scaffold Phase 13 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S02 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 edge provenance (required)

| Check | Expect |
|-------|--------|
| Store | `TestImportProvenanceRoundTrip` PASS — EXTRACTED/INFERRED/AMBIGUOUS persist; empty→EXTRACTED |
| Analyzers | `TestAnalyzerImportProvenanceExtracted` PASS (CGO1) — AST EXTRACTED; wildcard AMBIGUOUS; no analyzer INFERRED |
| Retrieval | Expand/Why copy `edge_provenance`; INFERRED store fixture not silent-as-EXTRACTED |
| Compiler | `TestContextWhyTraceEdgeProvenance` PASS — JSON/MD surface; causal confidence untouched |

### B. S02 packet honesty (required)

| Check | Expect |
|-------|--------|
| Loud totals | `TestBudgetLoudTotals` — `items_total`>`items_kept`, `truncated`, MD `items=` kept/total |
| Candidate cap | `TestCandidateCapSetsTruncated` — `candidates_capped` + `truncated` when Layer-1 capped |
| Index honesty | `TestIndexStaleBanner` — hash mismatch → `stale_paths`; false-fresh on errors/missing; MD banner; sorted + `len≤8` |
| Law 18 | Index honesty never sets causal `Provenance.Status` to STALE |
| S01 intact | `TestContextWhyTraceEdgeProvenance` still green |

### C. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | Product pkgs PASS (graphify space FAIL OK residual) |

### D. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals (incl. research ranks 4+ still off-board)
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood / research FUTURE may continue off-board)
6. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)

On FAIL: spawn `P12-S03-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P12-S03-01` only. Do not mark `P12-S03-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P12-S03-01`; next `P12-S03-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01–S02 named regressions
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/`./...`
- [ ] Inspect Gate C `dry_run:false`
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / new mig
- Scaffolding Phase 13 / research S03–S05 without promotion
- Re-scoring Gate C
- Closing parallel dogfood experiments
- Rewriting Phase 00–11 history
- Strengthening `TestIndexStaleBanner` to exact lex-first-8 (optional residual only)
