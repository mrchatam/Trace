# P16 / S06 / 01 — Phase 16 VERIFY (assert-root-and-surfaces closeout)

## Metadata
- id: P16-S06-01
- todo_ids: [P16-S06-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 16 — S01–S05 named DF regressions + carry-forward honesty/Gates/ablation/compat/p0x/x0 + product pkgs — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S05 Notes alone. Do **not** fail for DF-67 / P14 R2 / P15 R3/R4 / S05-02 `attachTaskImpact` swallow / 014 nine-Name list. Re-prove thin `trace_impact` (DF-72) as an S05 named test — do not treat historical P16-00 defer as the live lock ([`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md)). Do **not** invent research S05 / `plan simulate` / D21+ without promotion. **Phase 17** is independently queued after this phase’s board block — VERIFY still **starts** DR-HANDOFF = **`no successor`**; do not rewrite P17 or claim it as P16 successor.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Write durable evidence, then either:

1. **Pass** → declare **Phase 16 VERIFY PASS / assert-root-and-surfaces green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P16-S06-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-mcp-project-root/REVIEW-NOTES.md), [S02](../scope-02-tool-decision-enum/REVIEW-NOTES.md), [S03](../scope-03-cli-mcp-allowlist-parity/REVIEW-NOTES.md), [S04](../scope-04-install-project-root/REVIEW-NOTES.md), [S05](../scope-05-seed-impact-packet/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-72-FORWARD.md](../../DF-72-FORWARD.md)
- Pattern: Phase 15 VERIFY [`../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/01-verify.md`](../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md) — must be **FINAL**

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify). Unattended: no Plan-mode switch.

## Locked defaults (FINAL — P16-S06-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Assert-root-and-surfaces closeout (S01–S05 named DFs) — **not** a new `evals/*` planted gate |
| Catalog | **Ten** MCP names **including** `trace_version` (`trace_impact` before version); slug `mcp:trace_impact`; **no** install/decide/plan/index MCP |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` (S05 rollup unforked keeper) |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` (mig ceiling **14**, no 015+) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Residuals (non-blocking) | **DF-67 defer**; **P14 R2 defer**; **P15 R3/R4 wontfix**; S05-02 `attachTaskImpact` swallow; 014 nine-Name `IN` list; DF-22/37 tip-only |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary expansion; **no** install/decide/plan/index MCP; ten tools including `trace_version`; `trace_impact` **is** in catalog |
| Mig | No new migration from Phase 16 VERIFY |
| Full bar | Product packages `./cmd\|internal\|evals` PASS; R3 graphify space FAIL OK; R4 CGO0 analyzers / CGO0 `cmd/trace` FAIL OK |
| Successor | **`no successor`** — Phase 17 independently queued (do not rewrite; do not claim as this successor); research S05 / plan simulate / D21+ stay off-board unless Notes explicitly promote |
| Optional strong evidence | Grep catalog 10; MCP env `GOMODCACHE`/`GOPROXY=off` if proxy 403 |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 DF-76 | Virgin/`project=`/empty `.trace/` CallTool error, no auto-mkdir; isolation HOLD; `OpenExisting` sentinels; P15 Assert keepers |
| S02 DF-75/78 | Mig 014 CHECK rejects YOLO; heal→PENDING; Resolve no AUTO_ALLOWED overwrite; unprefixed Name → `mcp:` gates CallTool |
| S03 DF-77 | Dual-slug: MCP DENIED ≠ CLI DENIED; `cli:add`/`cli:why` fail-closed; ungated `capability decide`; `cli:reindex`→`cli:index` |
| S04 DF-68 | `-C` ProjectRoot for detect/claude/uninstall; cwd marker ignored; Cursor STABLE home unchanged; DF-22/37 tip keepers |
| S05 DF-70…74 | Seed mentions-task + findings/alternatives; packet `overall_class`; snake_case report; thin `trace_impact` + DENIED; catalog **10**; boundary allows impact only |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat **14**; p0x; x0; Gate C `dry_run:false`; product pkgs |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Residuals OK | DF-67/R2/R3/R4/`attachTaskImpact`/014 nine-Name **not** fail criteria |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no research S05 / plan simulate without promotion; P17 not this successor |

### Locked verify commands

Copy from sibling `00-PLANNER.md` **Locked verify command set (FINAL)** — do not invent a shorter substitute. Summary:

```bash
# --- S01 DF-76 named + P15 keepers ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPVirgin|TestMCPInitialized|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

# --- S02 DF-75/78 ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCapabilityToolDecisionCheckRejectsYOLO|TestCapabilityToolDecisionMigrateHealsYOLOToPending|TestResolveYOLOBuiltinDoesNotAutoAllow|TestDecideUnprefixedMCPNameCanonicalizes|TestCanonicalizeCustomAndCLISlugsUnchanged|TestMigrateUnprefixedDeniedFoldsOverAutoAllowed|TestMCPUnprefixedDecideGatesCallTool|TestMCPAssert|TestToolNamesRegistered|TestOpenCreates|TestMCPVirgin'

# --- S03 DF-77 ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCLIAddDeniedDoesNotBlockMCPAdd|TestUnprefixedAddDecideDoesNotGateCLI|TestCapabilityDecisionAutoAllowBuiltinCLI|TestCanonicalizeCLIReindexFoldsToIndex|TestCanonicalizeCustomAndCLISlugsUnchanged|TestCapabilityDecision|TestMCPAssert|TestMCPUnprefixed|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCLIAddSucceedsWhenMCPAddDenied|TestCLIAddDeniedFailClosed|TestCLIWhySucceedsWhenMCPWhyDenied|TestCLIWhyDeniedFailClosed|TestUngatedCapabilityDecideWhenCLIAddDenied|TestCLIIndexAliasDenied|TestUnprefixedAddDecideDoesNotGateCLI'

# --- S04 DF-68 ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/install/... ./internal/mcp/... ./internal/store/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallConditional|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallClaudeDashC|TestInstallDetectDashC|TestInstallUninstallClaudeDashC|TestInstallCursorPrintReloadTip|TestInstallCursorWriteMergeBackup|TestCLIAddDeniedFailClosed'

# --- S05 DF-70…74 + catalog 10 ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/compiler/... ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestContextIncludesImpactOverallClass|TestMCPTraceImpactReport|TestMCPImpactDeniedBlocksCallTool|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs|TestImportBoundaryMCPNoPlanImpactIndexTools|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestImpactReportJSONSnakeCase|TestWhyIncludesImpactOverallClass|TestCLIAddDeniedFailClosed|TestLinkDiscoveryMentionsTaskCLI|TestInstallClaudeDashCRefuseCitesProjectRoot'

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat checklist
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Product regression bar
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Grep: RegisteredToolNames length 10; trace_impact before trace_version; no plan/index/install/decide MCP
# Gate C artifact inspect: dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Spot-check only (non-blocking): DF-67 / R2 allowContainsOut / 014 nine-Name IN list
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S05 evidence is **named tests** — not Notes-only
- [ ] MCP remains **ten** tools including `trace_version`; `trace_impact` **in** catalog; no install/decide/plan/index MCP
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat **14** still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / Neo4j SoT still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration from VERIFY; mig 014 already in tree; compat ceiling **14**
- [ ] No YOLO / AllowAll defaults
- [ ] **DF-67 / R2 / R3 / R4 / attachTaskImpact swallow / 014 nine-Name do not fail VERIFY**
- [ ] **No research S05 / plan simulate / D21+ scaffold** unless Notes explicitly promote
- [ ] **Phase 17 board rows left intact** — not claimed as this successor
- [ ] Forward-only: do **not** rewrite Phase 00–15 `done` history; Phase 15 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S06-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** rewrite Phase 17. Do **not** auto-board research S05 / plan simulate / D21+.

| Who | Duty |
|-----|------|
| **P16-S06-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note DF-67 defer + R2 defer + R3/R4 wontfix + S05-02 residuals as non-blocking. Do **not** rewrite P17. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P16-S06-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 16 complete only then. Phase 17 already queued independently — confirm it is **not** this successor. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 16 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P16-S06-01a` | implement remediation (full prompt) |
| `P16-S06-01b` | review remediation |
| `P16-S06-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns. Do not weaken bars by treating DF-67/R2/R3/R4 as fail criteria.

## Board rights
Verify: **status + notes** on `P16-S06-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 16 `done` history. Do **not** mark `P16-S06-02` done. Do **not** rewrite Phase 17 rows.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S05 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` and `internal/mcp` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01–S05 named DF regressions (required)

| Check | Expect |
|-------|--------|
| S01 DF-76 | `TestMCPVirginProjectDoesNotMkdir`; `TestMCPInitializedOtherRootIsolated`; `TestOpenExistingMissingReturnsErrNotInitialized`; `TestOpenExistingEmptyTraceDir` PASS |
| S01 P15 keepers | `TestMCPAssertDeniedBlocksCallTool`; `TestMCPAssertBuiltinAutoAllowedSucceeds`; `TestToolNamesRegistered` (exactly **10**) PASS |
| S02 DF-75/78 | CHECK/heal/Resolve/canonicalize/fold/unprefixed-decide named tests PASS; `TestOpenCreates*` includes v**14** |
| S03 DF-77 | Isolation + `cli:` DENIED + alias + ungated decide named tests PASS (CLI under **CGO1**) |
| S04 DF-68 | DashC refuse/ignore-cwd/write/detect/cursor-home/uninstall + tip keepers PASS |
| S05 DF-70…74 | Seed/context/why/findings/snake_case + `TestMCPTraceImpactReport` + `TestMCPImpactDeniedBlocksCallTool` PASS |
| Catalog / boundary | `TestToolNamesRegistered` + `TestBuiltinMCPCapabilitySpecs` = **10** including `trace_version`; `TestImportBoundaryMCPNoPlanImpactIndexTools` allows `trace_impact` only |

### B. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS (ceiling **14**) |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| Product pkgs | `./cmd\|internal\|evals` PASS (R3 graphify / R4 CGO0 FAIL OK) |

### C. Residuals (must record, must not fail)

| Residual | Disposition | VERIFY rule |
|----------|-------------|-------------|
| DF-67 symbol-entity honesty | defer | Note only — do **not** fail |
| P14 R2 `allowContainsOut` | defer | Note only — do **not** fail |
| P15 R3 graphify space | wontfix | Product bar ≠ full `./...` — do **not** fail |
| P15 R4 CGO0 analyzers | wontfix | Product bar is CGO1; CGO0 `cmd/trace` FAIL OK — do **not** fail |
| S05-02 `attachTaskImpact` swallow | residual | capability-style — do **not** fail |
| 014 nine-Name `IN` list | historical | do **not** edit 014; do **not** fail |
| DF-22/37 Cursor reload | ops | tip keepers only — do **not** fail |

### D. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals (DF-67; R2; R3/R4; S05-02 swallow; 014 list; DF-22/37)
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that Phase 17 remains independently queued and is **not** this successor; research/dogfood may continue off-board)
6. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)

On FAIL: spawn `P16-S06-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P16-S06-01` only. Do not mark `P16-S06-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DF-67/R2/R3/R4 residuals explicit in Notes as non-blocking
- [ ] DF-72 named test **is** a fail bar (catalog 10 + `TestMCPTraceImpactReport` / DENIED)
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P16-S06-01`; next `P16-S06-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist
- [ ] Phase 17 rows **not** rewritten; **not** claimed as this successor

## Minimal todos
- [ ] Run S01–S05 named DF regressions (library CGO0 + CLI CGO1)
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/product pkgs
- [ ] Inspect Gate C `dry_run:false` (optional strong)
- [ ] Record DF-67/R2/R3/R4/S05-02 residuals as non-blocking
- [ ] Write VERIFY-NOTES + start DR-HANDOFF = `no successor`
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / new mig
- Rewriting Phase 17 / claiming P17 as this successor
- Scaffolding research S05 / plan simulate / D21+ without promotion
- Re-scoring Gate C
- Claiming DF-67 / R2 / R3 / R4 fixed
- Closing parallel dogfood experiments
- Rewriting Phase 00–15 history
- Using CGO0 `./cmd/trace/...` as a fail bar
