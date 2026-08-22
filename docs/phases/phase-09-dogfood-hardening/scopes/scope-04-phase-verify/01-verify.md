# P09 / S04 / 01 — Phase 09 VERIFY (dogfood hardening closeout)

## Metadata
- id: P09-S04-01
- todo_ids: [P09-S04-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 09 — DF-01 regression + S02/S03 dogfood UX surfaces + carry-forward honesty/Gates/ablation/compat/p0x/x0/Gate H/Gate C — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S03 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, add MCP list-tasks, or scaffold Phase 10.

Write durable evidence, then either:

1. **Pass** → declare **Phase 09 VERIFY PASS / dogfood hardening green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features on this row (except spawn remediations if a bar fails).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-retrieval-review/REVIEW-NOTES.md), [S02](../scope-02-discoverability/REVIEW-NOTES.md), [S03](../scope-03-install-wire/REVIEW-NOTES.md)
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md), [experiments/LADDER.md](../../../../../experiments/LADDER.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 08 VERIFY [`../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/01-verify.md`](../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P09-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Dogfood hardening closeout (DF-01/02/03/04/05 product surfaces) — **not** a new `evals/*` planted gate |
| DF-01 regression | **`TestWhyAndContextWithLinkedReview`** under `./internal/retrieval/...` |
| S02 discoverability | **`TestTasksListAfterSeed`** + **`TestSeedImportRelativePathAgainstC`** under `./cmd/trace/...`; optional `TestListTasks` store |
| S03 install-wire | **`TestInstallCursorPrintSnippet`** (prefer also WriteMergeBackup / WriteInvalidJSON / PrintBin); docs README + `experiments/ab-simple/PROTOCOL.md` DF-05 |
| MCP surface | Still **six** tools (`trace_why`/`trace_context`/`trace_add`/`trace_link`/`trace_transition`/`trace_review`) — **no** list-tasks requirement |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out; 100k/1M planted CI ladders deferred |
| Residuals (non-blocking) | S01: `plan_scope` ExactLookup out; scope-only review expand untested; S03: degenerate non-object `mcpServers` |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** on this row except spawn remediation |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; do **not** add list-tasks MCP |
| Successor | **`no successor`** — ladder gaps stay parallel `experiments/` unless Notes explicitly promote Phase 10 |

### Locked verify commands

```bash
# --- DF-01 primary regression ---
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run TestWhyAndContextWithLinkedReview

# S02 discoverability + DF-04 seed vs -C
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run TestListTasks
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestTasksListAfterSeed|TestSeedImportRelativePathAgainstC'

# S03 install-wire (+ preferred merge/invalid JSON)
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursor'

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
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... -count=1

# Full regression bar
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
# Spot-check install snippet shape (no --write required for VERIFY)
go build -o /tmp/trace-verify ./cmd/trace && /tmp/trace-verify install cursor | grep -E 'mcpServers|"trace"|workspaceFolder'
# Docs DF-05: README + experiments/ab-simple/PROTOCOL.md mention run-folder / workspaceFolder footgun
# MCP still six tools — grep BuiltinMCPCapabilitySpecs / mcp tool registration; no trace_tasks
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] DF-01 evidence is **`TestWhyAndContextWithLinkedReview`** — not Notes-only
- [ ] S02/S03 surfaces proven by named tests + DF-05 docs; **no** MCP list-tasks
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Mode-B packs not falsified
- [ ] Embeddings / VerifiedFact / `plan simulate` still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] No new migration `011_*` from Phase 09
- [ ] **No Phase 10 scaffold** unless Notes explicitly promote

### DR-HANDOFF duties (this row + S04-02)

Per protocol Phase handoff + **DR-HANDOFF**. On green → record **`no successor`**. Do **not** create Phase 10 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P09-S04-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table; explicitly record **DR-HANDOFF = `no successor`**. Note remaining ladder gaps (D08/D09/combos/multi-agent; DF-11/12 tighten) stay on parallel dogfood track. Do **not** invent Phase 10. |
| **P09-S04-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 09 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 09 complete; do **not** invent a successor to dodge a red VERIFY.

## Board rights
Verify: **status + notes** on `P09-S04-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 09 `done` history. Do **not** mark `P09-S04-02` done. Do **not** scaffold Phase 10 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02/S03 REVIEW-NOTES + locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf,compat}` exist.
3. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. DF-01 regression (required)

| Check | Expect |
|-------|--------|
| Named test | `TestWhyAndContextWithLinkedReview` PASS |
| Behavior | Why + Expand succeed with linked PASS review; no `unknown entity type "review"` |
| Honesty of claim | Maps to live `lookupEntity` `case "review"` — not Notes-only |

### B. S02 discoverability (required)

| Check | Expect |
|-------|--------|
| `trace tasks` | `TestTasksListAfterSeed` PASS — JSON id/title/work_state/goal_id |
| Seed vs `-C` | `TestSeedImportRelativePathAgainstC` PASS |
| No MCP list-tasks | MCP still six tools |

### C. S03 install-wire + DF-05 (required)

| Check | Expect |
|-------|--------|
| Install tests | `TestInstallCursor*` PASS (at least PrintSnippet) |
| Snippet shape | `mcpServers.trace` with `-C` `${workspaceFolder}` (stdio; default `trace-mcp`) |
| Docs | README + `experiments/ab-simple/PROTOCOL.md` document install + **DF-05** (open **run folder** as workspace) |

### D. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| `./...` | PASS |

### E. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Law checks
4. Residuals / deferrals
5. Explicit **DR-HANDOFF = `no successor`** (+ one-liner that ladder gaps remain parallel)

On FAIL: spawn `P09-S04-01a` / `01b` (+ `01c` if needed) with full prompts; do not weaken bars.

## Todo updates
Status + Notes on `P09-S04-01` only. Do not mark `P09-S04-02` done.

## Exit criteria
- [ ] Locked commands run independently (or FAIL+spawn trail)
- [ ] `VERIFY-NOTES.md` written with evidence table
- [ ] DF-01 + S02 + S03 + carry-forward + `./...` green **or** remediations spawned
- [ ] DR-HANDOFF started as **`no successor`** (completion owned by S04-02)
- [ ] Board Notes updated; next **`P09-S04-02`** (or spawn rows)

## Minimal todos
- [ ] Run locked VERIFY commands
- [ ] Spot-check install docs + MCP six-tool surface
- [ ] Write VERIFY-NOTES.md
- [ ] Board update (pass → handoff start; fail → spawn)

## Out of scope
- Product features / new MCP tools / daemon
- Scaffolding Phase 10 without promotion
- Re-scoring Gate C / inventing new Go
- Closing the dogfood ladder (experiments continue in parallel)
