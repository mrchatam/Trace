# P10 / S05 / 01 — Phase 10 VERIFY (integrity-surfaces closeout)

## Metadata
- id: P10-S05-01
- todo_ids: [P10-S05-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 10 — S01–S04 DF regressions (retrieval-why, MCP parity, index GC, operator/capability gates) + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S04 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add plan/impact/index MCP, or scaffold Phase 11.

Write durable evidence, then either:

1. **Pass** → declare **Phase 10 VERIFY PASS / integrity surfaces green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P10-S05-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-retrieval-why-fidelity/REVIEW-NOTES.md), [S02](../scope-02-mcp-parity-install/REVIEW-NOTES.md), [S03](../scope-03-index-gc/REVIEW-NOTES.md), [S04](../scope-04-operator-capability-gates/REVIEW-NOTES.md)
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 09 VERIFY [`../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/01-verify.md`](../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P10-S05-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Integrity-surfaces closeout (DF-17…32 product surfaces) — **not** a new `evals/*` planted gate |
| S01 DF-19 | `TestWhyTaskDPCGoalScoped` + `TestWhyTaskDPCMultiGoalNoForeignPollution` (+ `TestTaskContextDPCGoalScoped` / `TestTaskContextMultiGoalOmitsForeignDPC`) — **supersedes** historical GC-01 global-attach names |
| S01 DF-23/25/27/29 | `TestExactWhyPlanChangeAlias` / `TestExactWhyCapability` / `TestDecisionMarkdownTrustLabels` / `TestIncludeWhyFailClosed` |
| S02 MCP surface | **Nine** tools: `trace_why`/`trace_context`/`trace_add`/`trace_link`/`trace_transition`/`trace_review` + `trace_tasks`/`trace_capability`/`trace_version` |
| S02 named | `TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs` / `TestTraceTasksParity` / `TestTraceCapabilityActions` / `TestTraceVersion` / `TestCapabilityListMissingSnakeCase` / `TestImportBoundaryMCPNoPlanImpactIndexTools` (+ install tip asserts) |
| S03 DF-20 | `TestIndexGCAfterPathRename` / `TestIndexArgvMissingPathDeletesOnlyThatPath` / `TestIndexIncrementalIsolation` (+ `TestListFilePathsAndDeleteFileByPath`) |
| S04 DF-17/18/24/26/31 | `TestOperatorDoneRequiresFlag` / `TestOperatorDoneHatchBypassesOperator` / `TestReopenInvalidatesPassReviews` / `TestMissingCapabilitiesBlockTransition` / `TestAllowDoneWarnsOnStderr` / `TestCapabilityMissingRequiresTaskHint` + MCP `TestTransitionAllowDoneEmitsWarning` / `TestCapabilityMissingRequiresTaskParam` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` (Path C = `AllowOperatorDone`; Actor≠auth) |
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
| Fixture hash pin (carry) | Historical pin may drift after S02 README edits — note residual; do **not** fail VERIFY solely on pin drift |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out; 100k/1M planted CI ladders deferred; DF-28/30/33/34…36 out |
| Residuals (non-blocking) | `plan_scope` ExactLookup out; Mode-B packs historical; Cursor MCP reload manual; ab-index / ab-operator-gate optional; adapter missing-caps e2e test low |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; **no** plan/impact/index MCP tools |
| Mig | No `011_*` from Phase 10 |
| Full bar | Product packages in `CGO_ENABLED=1 go test ./...` PASS; known FAIL only `similar projects/graphify` space-in-path (non-product) |
| Successor | **`no successor`** — parallel dogfood stays off-board unless Notes explicitly promote Phase 11 |

### Locked verify commands

```bash
# --- S01 retrieval / why fidelity ---
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... -count=1 -run 'TestWhyTaskDPCGoalScoped|TestWhyTaskDPCMultiGoalNoForeignPollution|TestTaskContextDPCGoalScoped|TestTaskContextMultiGoalOmitsForeignDPC|TestExactWhyPlanChangeAlias|TestExactWhyCapability|TestDecisionMarkdownTrustLabels|TestIncludeWhyFailClosed'

# --- S02 MCP parity + install + snake_case ---
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... -count=1 -run 'TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs|TestTraceTasksParity|TestTraceCapabilityActions|TestTraceVersion|TestImportBoundaryMCPNoPlanImpactIndexTools'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCapabilityListMissingSnakeCase|TestInstallCursor'

# --- S03 index GC ---
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run TestListFilePathsAndDeleteFileByPath
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestIndexGCAfterPathRename|TestIndexArgvMissingPathDeletesOnlyThatPath|TestIndexIncrementalIsolation'

# --- S04 operator / capability gates (+ MCP transition) ---
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... -count=1 -run 'TestOperatorDoneRequiresFlag|TestOperatorDoneHatchBypassesOperator|TestReopenInvalidatesPassReviews|TestMissingCapabilitiesBlockTransition|TestTransitionAllowDoneEmitsWarning|TestCapabilityMissingRequiresTaskParam'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestAllowDoneWarnsOnStderr|TestCapabilityMissingRequiresTaskHint'

# Honesty: Paths A/B/C + Gate G (CGO-free) — Path C operator-flag supersession
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

# Full regression bar (product pkgs; graphify space FAIL is known residual)
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Nine MCP tools — grep registration / BuiltinMCPCapabilitySpecs; no plan/impact/index MCP
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# Actor≠auth: spot-check TestOperatorDoneRequiresFlag rejects Actor:"operator" without AllowOperatorDone
# G19: library packages do not import cmd/trace or cmd/trace-mcp
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S04 evidence is **named tests** — not Notes-only
- [ ] MCP is **nine** tools; no plan/impact/index MCP
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Mode-B packs not falsified
- [ ] Embeddings / VerifiedFact / `plan simulate` still out
- [ ] No full-rebuild-on-any-change indexer architecture (DF-20 is delete-on-missing, not full rebuild)
- [ ] No new migration `011_*` from Phase 10
- [ ] **No Phase 11 scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 09 `done` history

### DR-HANDOFF duties (this row + S05-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 11 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P10-S05-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Note residuals stay parallel `experiments/`. Do **not** invent Phase 11. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P10-S05-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 10 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 10 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P10-S05-01a` | implement remediation (full prompt) |
| `P10-S05-01b` | review remediation |
| `P10-S05-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns.

## Board rights
Verify: **status + notes** on `P10-S05-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 10 `done` history. Do **not** mark `P10-S05-02` done. Do **not** scaffold Phase 11 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S04 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 retrieval / why fidelity (required)

| Check | Expect |
|-------|--------|
| DF-19 | Goal-scoped DPC tests PASS; no foreign multi-goal pollution |
| DF-23/25/27/29 | Alias / capability Exact / MD trust labels / IncludeWhy fail-closed PASS |
| Honesty of claim | Maps to live retrieval/compiler — not Notes-only |

### B. S02 MCP parity + install (required)

| Check | Expect |
|-------|--------|
| Nine tools | `TestToolNamesRegistered` + BuiltinMCP×9 PASS |
| DF-21/22/32 | tasks/capability/version + snake_case + install tip PASS |
| Boundary | No plan/impact/index MCP (`TestImportBoundaryMCPNoPlanImpactIndexTools`) |

### C. S03 index GC (required)

| Check | Expect |
|-------|--------|
| Rename GC | `TestIndexGCAfterPathRename` PASS |
| Argv isolation | Missing argv deletes **only** that path — not project-wide GC |
| Isolation | `TestIndexIncrementalIsolation` PASS |

### D. S04 operator / capability gates (required)

| Check | Expect |
|-------|--------|
| DF-17 | Operator DONE requires flag; Actor alone insufficient; MCP `as_operator` |
| DF-18 | Leave DONE → PASS→UNCERTAIN |
| DF-24/26/31 | Missing-cap fail-closed; loud hatch WARNING; missing-task UX |
| Honesty | Path C operator-flag; Gate G hatch retained |

### E. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | Product pkgs PASS (graphify space FAIL OK residual) |

### F. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood may continue)

On FAIL: spawn `P10-S05-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P10-S05-01` only. Do not mark `P10-S05-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P10-S05-01`; next `P10-S05-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01–S04 named regressions
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/`./...`
- [ ] Inspect Gate C `dry_run:false`
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / mig `011_*`
- Scaffolding Phase 11 without promotion
- Re-scoring Gate C
- Closing parallel dogfood experiments
- Rewriting Phase 09 history
