# P11 / S08 / 01 — Phase 11 VERIFY (residual-surfaces closeout)

## Metadata
- id: P11-S08-01
- todo_ids: [P11-S08-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 11 — S01–S07 DF regressions (**18** canonical DFs) + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S07 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add plan/impact/index MCP dump, or scaffold Phase 12.

Write durable evidence, then either:

1. **Pass** → declare **Phase 11 VERIFY PASS / residual surfaces green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P11-S08-01a` / `01b` / +`01c`).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-index-partial-path-gc/REVIEW-NOTES.md), [S02](../scope-02-review-pass-fail-operator/REVIEW-NOTES.md), [S03](../scope-03-store-lock-concurrency/REVIEW-NOTES.md), [S04](../scope-04-capability-upsert-hatch/REVIEW-NOTES.md), [S05](../scope-05-retrieval-why-depth-trust/REVIEW-NOTES.md), [S06](../scope-06-mcp-install-reload/REVIEW-NOTES.md), [S07](../scope-07-seed-plan-review-polish/REVIEW-NOTES.md)
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 10 VERIFY [`../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/01-verify.md`](../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P11-S08-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Residual-surfaces closeout (18 P11 DFs) — **not** a new `evals/*` planted gate |
| S01 DF-40 | `TestIndexPartialArgvGCAfterRename` (+ `TestListFilePathsByContentHash`); retain DF-20 trio |
| S02 DF-43/44 | Sibling FAIL+PASS named tests + conscious-flag docs + `TestOperatorDoneRequiresFlag` |
| S03 DF-47 | Retry / exclusivity / serialize guidance / init fail-closed named tests |
| S04 DF-41/51 | Slug upsert + hatch≠missing-caps + WARNING wording tests |
| S05 DF-49/35/48/42 | Symbol Exact/Why; depth-2 no sibling body; Law 9+4 MD; discovery→task link |
| S06 DF-50/22/37 | Print+write tip parity; nine tools / `trace_version`; tip/docs only (no PID kill) |
| S07 DF-33/30/46/45/28 | Seed aliases; plan `phases:[]`+`tasks` snake_case; review get/show/list; thin help handoff SoT |
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
| Residuals (non-blocking) | rename+edit ghost until full-tree; graphify space FAIL; CGO0 analyzers FAIL OK; Cursor MCP reload manual; parallel ab-* |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; **no** plan/impact/index MCP dump |
| Mig | No `011_*` from Phase 11 |
| Full bar | Product packages in `CGO_ENABLED=1 go test ./...` PASS; known FAIL only `similar projects/graphify` space-in-path (non-product) |
| Successor | **`no successor`** — parallel dogfood stays off-board unless Notes explicitly promote Phase 12 |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 | DF-40 partial-path GC after rename (`TestIndexPartialArgvGCAfterRename`); DF-20 retain |
| S02 | DF-43 PASS+FAIL sibling block; DF-44 flag≠identity docs; honesty Path C |
| S03 | DF-47 retry-on-brief-contention + exclusivity fail-closed + serialize guidance |
| S04 | DF-41 slug upsert same-id; DF-51 hatch does not bypass missing-caps |
| S05 | DF-49 symbol Exact/Why; DF-35 no sibling body; DF-48 Law 9+4 / `untrusted_data`; DF-42 discovery→task |
| S06 | DF-50/22/37 tip parity print+write; nine tools / `trace_version` |
| S07 | DF-33 aliases; DF-30 `phases:[]`+`tasks`; DF-46 snake_case; DF-45 review CLI; DF-28 help handoff SoT |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false` |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Laws | No daemon/HTTP/embeddings; no full-rebuild indexer; G19; no Phase 12 without promotion |

### Locked verify commands

```bash
# --- S01 index partial-path GC (DF-40) + DF-20 retain ---
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestListFilePathsByContentHash'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestIndexPartialArgvGCAfterRename|TestIndexGCAfterPathRename|TestIndexArgvMissingPathDeletesOnlyThatPath|TestIndexIncrementalIsolation'

# --- S02 review PASS+FAIL / operator identity (DF-43/44) ---
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... -count=1 -run 'TestSiblingFailBlocksDone|TestSiblingPassAloneAllowsDone|TestSiblingPassPlusUncertainAllowsDone|TestHatchBypassesSiblingFail|TestOperatorDoneRequiresFlag|TestAsOperatorSchemaIdentityDocs'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestAsOperatorFlagIdentityDocs'

# --- S03 store lock / concurrency (DF-47) ---
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestOpenRetrySucceedsWhenLockReleasedSoon|TestConcurrentStoreOpenFailClosed|TestErrLockedSerializeGuidance'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSerializeLockGuidance|TestInitFailClosedWhenStoreLocked'

# --- S04 capability upsert + hatch vs caps (DF-41/51) ---
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... -count=1 -run 'TestUpsertCapabilityBySlugUpdatesExisting|TestUpsertCapabilityGetAndReject|TestAllowDoneDoesNotBypassMissingCaps|TestTransitionAllowDoneEmitsWarning'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestAllowDoneWarnsOnStderr'

# --- S05 retrieval why / depth / trust / DPC (DF-49/35/48/42) ---
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... -count=1 -run 'TestWhySymbolExact|TestGetSymbolByID|TestExpandDepth2NoSiblingTaskBody|TestExpandContextDepth2NoSiblingTaskBody|TestDecisionMarkdownTrustLabels|TestLinkDiscoveryMentionsTask|TestWhyTaskDPCMultiGoalNoForeignPollution'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestLinkDiscoveryMentionsTaskCLI'

# --- S06 MCP / install reload (DF-50/22/37) ---
CGO_ENABLED=0 go test ./internal/mcp/... -count=1 -run 'TestToolNamesRegistered|TestTraceVersion'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursorPrintReloadTip|TestInstallCursorWriteMergeBackup|TestInstallCursorWriteCreateMissing'

# --- S07 seed / plan / review polish (DF-33/30/46/45/28) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportFromIDAliases|TestSeedImportMissingEndpointsMessage|TestPlanShowSnakeCaseAndEmptyPhases|TestPlanShowWithPhasesSnakeCase|TestReviewGetShowList|TestHelpHandoffSoT'

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

# Full regression bar (product pkgs; graphify space FAIL is known residual)
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Nine MCP tools — grep registration / BuiltinMCPCapabilitySpecs; no plan/impact/index MCP
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# Actor≠auth: spot-check TestOperatorDoneRequiresFlag rejects Actor:"operator" without AllowOperatorDone
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Findings: on PASS, flip DF-22/28/30/33/35/37/40–51 to closed in DOGFOOD-FINDINGS.md (forward-only)
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S07 evidence is **named tests** — not Notes-only
- [ ] MCP remains **nine** tools; no plan/impact/index MCP dump
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Embeddings / VerifiedFact / `plan simulate` still out
- [ ] No full-rebuild-on-any-change indexer architecture (DF-40 is hash-orphan delete, not full rebuild)
- [ ] No new migration `011_*` from Phase 11
- [ ] **No Phase 12 scaffold** unless Notes explicitly promote
- [ ] Forward-only: do **not** rewrite Phase 00–10 `done` history; Phase 10 historical `no successor` left intact as history

### DR-HANDOFF duties (this row + S08-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**. Do **not** create Phase 12 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P11-S08-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`** (start). Optionally flip the 18 DF rows closed in findings. Note residuals stay parallel `experiments/`. Do **not** invent Phase 12. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P11-S08-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 11 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 11 complete; do **not** invent a successor to dodge a red VERIFY.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P11-S08-01a` | implement remediation (full prompt) |
| `P11-S08-01b` | review remediation |
| `P11-S08-01c` | re-VERIFY (optional if needed after 01b) |

Do not weaken bars to avoid spawns.

## Board rights
Verify: **status + notes** on `P11-S08-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 11 `done` history. Do **not** mark `P11-S08-02` done. Do **not** scaffold Phase 12 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S07 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01 index partial-path GC (required)

| Check | Expect |
|-------|--------|
| DF-40 | `TestIndexPartialArgvGCAfterRename` PASS; old path/FTS gone; sibling intact |
| DF-20 retain | Full-tree + missing-argv + isolation tests PASS |

### B. S02 review / operator (required)

| Check | Expect |
|-------|--------|
| DF-43 | Sibling FAIL+PASS rejects →DONE; PASS alone / PASS+UNCERTAIN OK; hatch bypasses FAIL |
| DF-44 | Conscious flag≠identity docs; Actor≠auth |
| Honesty | Path C supersedes FAIL before DONE; A/B/C + Gate G green |

### C. S03 store lock (required)

| Check | Expect |
|-------|--------|
| Retry | Brief release during wait → second Open succeeds |
| Exclusivity | Held lock past budget → `ErrLocked`; compat `trace_lock_ok` |
| UX | Serialize CLI↔MCP / worktree guidance; exit **2** |

### D. S04 capability / hatch (required)

| Check | Expect |
|-------|--------|
| DF-41 | Empty-ID re-declare same slug → same id + update; different-id clash fails |
| DF-51 | Hatch alone does not bypass missing-caps; WARNING mentions override |

### E. S05 retrieval / trust (required)

| Check | Expect |
|-------|--------|
| DF-49 | `why`/`Exact` symbol by id |
| DF-35 | Depth-2 no sibling task body |
| DF-48 | MD Law 9 honor; JSON `untrusted_data` |
| DF-42 | CLI/MCP `discovery-mentions-task` → store rel |

### F. S06 MCP install reload (required)

| Check | Expect |
|-------|--------|
| DF-50/22 | Print + write same stderr tip; stdout JSON-only |
| DF-37 | Tip/docs only — no PID kill / daemon / new tools |
| Surface | Nine tools + `trace_version` still registered |

### G. S07 seed / plan / review polish (required)

| Check | Expect |
|-------|--------|
| DF-33 | `from_id`/`to_id` aliases + empty-endpoint message |
| DF-30/46 | `phases: []` + `tasks`; snake_case keys |
| DF-45 | `review get\|show\|list [--task]` |
| DF-28 | Help thin handoff SoT; no entity/mig |

### H. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | Product pkgs PASS (graphify space FAIL OK residual) |

### I. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that parallel dogfood may continue)
6. Optional: findings DF status flip for the 18 closed DFs

On FAIL: spawn `P11-S08-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P11-S08-01` only. Do not mark `P11-S08-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DR-HANDOFF **started** = `no successor` (or explicit promotion documented)
- [ ] Board Notes on `P11-S08-01`; next `P11-S08-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist

## Minimal todos
- [ ] Run S01–S07 named regressions
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/`./...`
- [ ] Inspect Gate C `dry_run:false`
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / mig `011_*`
- Scaffolding Phase 12 without promotion
- Re-scoring Gate C
- Closing parallel dogfood experiments
- Rewriting Phase 00–10 history
