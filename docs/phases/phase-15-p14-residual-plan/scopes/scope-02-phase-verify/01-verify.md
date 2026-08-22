# P15 / S02 / 01 — Phase 15 VERIFY (P14 residual remediation closeout)

## Metadata
- id: P15-S02-01
- todo_ids: [P15-S02-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 15 — S01 MCP Assert named regressions + carry-forward honesty/Gates/ablation/compat/p0x/x0 + product pkgs — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add install/decide MCP dump, claim R2/R3/R4 fixed, invent Phase 16 / S05 / `plan simulate` / D21+ without promotion.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Write durable evidence, then either:

1. **Pass** → declare **Phase 15 VERIFY PASS / P14 residual remediation green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P15-S02-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-mcp-assert-dispatch/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 14 VERIFY [`../../../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md) — must be **FINAL**

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P15-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | P14 residual remediation closeout (S01 MCP Assert) — **not** a new `evals/*` planted gate |
| S01 | Assert on every MCP CallTool path with slug `mcp:<Name>` (incl. `trace_version`) |
| S01 named | `TestMCPAssertDeniedBlocksCallTool`; `TestMCPAssertBuiltinAutoAllowedSucceeds`; `TestToolNamesRegistered`; keep `TestBuiltinMCPCapabilitySpecs` |
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
| Residuals (non-blocking) | **R2 defer** (`allowContainsOut`); **R3 wontfix** (graphify space); **R4 wontfix** (CGO0 analyzers); goals #2–#4 deferred |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary expansion; **no** install/decide MCP; nine tools + `trace_version`; Assert **is** on MCP (R1 closed) |
| Mig | No new migration from Phase 15 VERIFY |
| Full bar | Product packages `./cmd\|internal\|evals` PASS; R3 graphify space FAIL OK; R4 CGO0 analyzers FAIL OK |
| Successor | **`no successor`** — S05 / plan simulate / D21+ stay off-board unless Notes explicitly promote |
| Optional strong evidence | Grep nine `assertMCPToolAllowed` sites; MCP env `GOMODCACHE`/`GOPROXY=off` if proxy 403 |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 MCP Assert | DENIED blocks CallTool; builtin AUTO_ALLOWED succeeds; exactly nine tools registered; slugs ≡ BuiltinMCPCapabilitySpecs |
| S01 wire-up | `assertMCPToolAllowed` at all nine entries incl. `toolVersion` (grep + named tests) |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat **13**; p0x; x0; Gate C `dry_run:false`; product pkgs |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Residuals OK | R2/R3/R4 **not** fail criteria; disposition matrix unchanged |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no Phase 16 / S05 / plan simulate without promotion |

### Locked verify commands

```bash
# --- S01 MCP Assert named ---
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

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

# Product regression bar (prefer over full-module ./... when R3 graphify present)
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Grep: assertMCPToolAllowed at all nine tool entries (incl. toolVersion)
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Spot-check only (non-blocking): R2 allowContainsOut residual still in impact_walk.go
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01 evidence is **named tests** — not Notes-only
- [ ] MCP remains **nine** tools + `trace_version`; no install/decide MCP tools
- [ ] Assert **is** on MCP dispatch (R1) — DENIED + AUTO_ALLOWED named tests green
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / Neo4j SoT still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration from VERIFY; mig 013 already in tree; compat ceiling **13**
- [ ] No YOLO / AllowAll defaults
- [ ] **R2/R3/R4 do not fail VERIFY** (defer / wontfix per disposition)
- [ ] **No Phase 16 / S05 / plan simulate / D21+ scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 00–14 `done` history; Phase 14 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S02-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 16 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P15-S02-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note R2 defer + R3/R4 wontfix + goals #2–#4 stay deferred. Do **not** invent Phase 16. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P15-S02-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 15 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 15 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P15-S02-01a` | implement remediation (full prompt) |
| `P15-S02-01b` | review remediation |
| `P15-S02-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns. Do not weaken bars by treating R2/R3/R4 as fail criteria.

## Board rights
Verify: **status + notes** on `P15-S02-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 15 `done` history. Do **not** mark `P15-S02-02` done. Do **not** scaffold Phase 16 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` and `internal/mcp` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 MCP Assert (required)

| Check | Expect |
|-------|--------|
| DENIED fail-closed | `TestMCPAssertDeniedBlocksCallTool` PASS |
| Builtin AUTO_ALLOWED | `TestMCPAssertBuiltinAutoAllowedSucceeds` PASS |
| Tool registry | `TestToolNamesRegistered` — exactly nine names |
| Specs alignment | `TestBuiltinMCPCapabilitySpecs` PASS; slugs `mcp:`+name |
| Wire-up | Grep/spot: `assertMCPToolAllowed` at all nine entries incl. `toolVersion` |

### B. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS (ceiling **13**) |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| Product pkgs | `./cmd\|internal\|evals` PASS (R3 graphify / R4 CGO0 FAIL OK) |

### C. Residuals (must record, must not fail)

| Residual | Disposition | VERIFY rule |
|----------|-------------|-------------|
| R2 `allowContainsOut` | defer | Note only — do **not** fail |
| R3 graphify space | wontfix | Product bar ≠ full `./...` — do **not** fail |
| R4 CGO0 analyzers | wontfix | Product bar is CGO1 — do **not** fail |

### D. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals (R2 defer; R3/R4 wontfix; goals #2–#4 still off-board)
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood / research FUTURE may continue off-board)
6. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)

On FAIL: spawn `P15-S02-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P15-S02-01` only. Do not mark `P15-S02-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] R2/R3/R4 residuals explicit in Notes as non-blocking
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P15-S02-01`; next `P15-S02-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01 named MCP Assert regressions
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/product pkgs
- [ ] Inspect Gate C `dry_run:false` (optional strong)
- [ ] Record R2/R3/R4 as non-blocking
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / new mig
- Scaffolding Phase 16 / S05 / plan simulate / D21+ without promotion
- Re-scoring Gate C
- Fixing R2 / R3 / R4
- Closing parallel dogfood experiments
- Rewriting Phase 00–14 history
- Claiming S01 APPROVE means R2/R3/R4 are fixed
