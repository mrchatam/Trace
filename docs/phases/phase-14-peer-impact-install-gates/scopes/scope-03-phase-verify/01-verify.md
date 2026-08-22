# P14 / S03 / 01 — Phase 14 VERIFY (peer-impact-install-gates closeout)

## Metadata
- id: P14-S03-01
- todo_ids: [P14-S03-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 14 — S01 ImpactWalk named regressions + S02 install/capability-gate named regressions + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S02 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add plan/impact/install MCP dump, claim “every MCP call gated,” invent Phase 15 / S05 / `plan simulate` / D21+ without promotion.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Write durable evidence, then either:

1. **Pass** → declare **Phase 14 VERIFY PASS / peer-impact-install-gates green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P14-S03-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-impact-walks/REVIEW-NOTES.md), [S02](../scope-02-install-capability-gates/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 13 VERIFY [`../../../phase-13-import-resolve-honesty/scopes/scope-04-phase-verify/01-verify.md`](../../../phase-13-import-resolve-honesty/scopes/scope-04-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md) — must be **FINAL**

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P14-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Peer-impact-install-gates closeout (S01+S02) — **not** a new `evals/*` planted gate |
| S01 | ImpactWalk multi-seed BFS + contains asymmetry + loud truncation; Gate F kept |
| S02 | Install registry + CONDITIONAL markers + mig **013** decisions; Assert fail-closed; CLI-first |
| S01 named | `TestImpactWalkMultiSeedExcludeSeeds`; `TestImpactWalkContainsAsymmetryNoSiblings`; `TestImpactWalkIncomingImportHop`; `TestImpactWalkLoudTruncation`; `TestImpactWalkHopRiskIncreases`; `TestPlantedImpactConflictsGateFPrelim` |
| S02 named | `TestInstallDetectListsCursorStable`; `TestInstallCursorUninstallIdempotent`; `TestInstallConditionalRefusesWithoutMarker`; `TestInstallConditionalWritesWithMarker`; `TestCapabilityDecisionAutoAllowBuiltinMCP`; `TestCapabilityDecisionUnknownPendingFailClosed`; `TestCapabilityDecisionHumanAllowPersists`; `TestCapabilityDecisionDenyBlocks`; keep `TestInstallCursor*` |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` (mig ceiling **13**, no 014+) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Residuals (non-blocking) | Assert ≠ MCP dispatch; optional `allowContainsOut` spot-check; graphify space FAIL; CGO0 analyzers FAIL OK; goals #2–#4 deferred |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; **no** install/decide MCP; nine tools + `trace_version` |
| Mig | No new migration from Phase 14 VERIFY |
| Full bar | Product packages in `CGO_ENABLED=1 go test ./...` PASS; known FAIL only `similar projects/graphify` space-in-path (non-product) |
| Successor | **`no successor`** — S05 / plan simulate / D21+ stay off-board unless Notes explicitly promote |
| Optional strong evidence | `TestImpactWalkCLI`; MCP env `GOMODCACHE`/`GOPROXY=off` if proxy 403 |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 ImpactWalk | Multi-seed + seed exclusion; contains asymmetry no siblings; incoming import hop (+ provenance when set); loud truncation totals; hop_risk monotonic; Gate F planted green |
| S02 install | detect lists Cursor STABLE; uninstall idempotent; CONDITIONAL refuse/write with marker; `TestInstallCursor*` keepers |
| S02 decisions | AUTO_ALLOWED builtin; unknown PENDING fail-closed; human ALLOWED persists; DENIED blocks; ablation green |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat **13**; p0x; x0; Gate C `dry_run:false` |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Honesty residuals | Assert ≠ MCP dispatch recorded; optional allowContainsOut Note |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no Phase 15 / S05 / plan simulate without promotion |

### Locked verify commands

```bash
# --- S01 ImpactWalk ---
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestImpactWalk'
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim

# --- S02 install registry + capability decisions ---
CGO_ENABLED=0 go test ./internal/install/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallCursorUninstallIdempotent|TestInstallConditional'
CGO_ENABLED=0 go test ./internal/domain/... -count=1 -run 'TestCapabilityDecision'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursor|TestInstallUsage|TestImpactWalkCLI'

# Capability ablation keep
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

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
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/retrieval/... ./internal/install/... -count=1
# If mcp proxy 403: GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off …
CGO_ENABLED=0 go test ./internal/mcp/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1

# Full regression bar (product pkgs; graphify space FAIL is known residual)
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Honesty: AssertToolAllowed is library/CLI only — do not claim MCP request-path gating
# Optional non-blocking: spot-check allowContainsOut residual in impact_walk.go
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S02 evidence is **named tests** — not Notes-only
- [ ] MCP remains **nine** tools + `trace_version`; no install/decide MCP tools
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / Neo4j SoT still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration from VERIFY; mig 013 already in S02; compat ceiling **13**
- [ ] No YOLO / AllowAll defaults
- [ ] Assert ≠ MCP dispatch — honesty Note only (not a VERIFY fail)
- [ ] **No Phase 15 / S05 / plan simulate / D21+ scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 00–13 `done` history; Phase 13 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 15 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P14-S03-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note Assert≠MCP + optional allowContainsOut + goals #2–#4 stay deferred. Do **not** invent Phase 15. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P14-S03-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 14 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 14 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P14-S03-01a` | implement remediation (full prompt) |
| `P14-S03-01b` | review remediation |
| `P14-S03-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns.

## Board rights
Verify: **status + notes** on `P14-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 14 `done` history. Do **not** mark `P14-S03-02` done. Do **not** scaffold Phase 15 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S02 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` and `internal/{retrieval,install,domain}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 ImpactWalk (required)

| Check | Expect |
|-------|--------|
| Multi-seed + exclusion | Seeds absent from blast; shared neighbors once at min hop |
| Contains asymmetry | No sibling symbols via contains climb |
| Incoming import | Importer appears; provenance when set |
| Loud truncation | Cap 64 + `blast_total`/`blast_kept`/`truncated` |
| Hop risk | Deeper hop ⇒ ≥ hop_risk |
| Gate F | Planted prelim PASS |

### B. S02 install + capability decisions (required)

| Check | Expect |
|-------|--------|
| Detect | Cursor STABLE + non-empty reason |
| Uninstall | Idempotent; only `mcpServers.trace` removed |
| CONDITIONAL | Refuse without marker; write with marker |
| Cursor keepers | `TestInstallCursor*` PASS |
| Decisions | AUTO_ALLOWED / PENDING fail-closed / ALLOWED persists / DENIED blocks |
| Ablation | `TestPlantedCapabilitySelectionAblation` PASS |

### C. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS (ceiling **13**) |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | Product pkgs PASS (graphify space FAIL OK residual) |

### D. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals (Assert≠MCP; optional allowContainsOut; goals #2–#4 still off-board)
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood / research FUTURE may continue off-board)
6. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)

On FAIL: spawn `P14-S03-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P14-S03-01` only. Do not mark `P14-S03-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] Assert≠MCP + optional allowContainsOut residuals explicit in Notes
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P14-S03-01`; next `P14-S03-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01–S02 named regressions
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/`./...`
- [ ] Inspect Gate C `dry_run:false`
- [ ] Record Assert≠MCP honesty Note
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / new mig
- Scaffolding Phase 15 / S05 / plan simulate / D21+ without promotion
- Re-scoring Gate C
- Wiring Assert into MCP dispatch
- Closing parallel dogfood experiments
- Rewriting Phase 00–13 history
- Claiming S02 APPROVE means MCP request-path gating
