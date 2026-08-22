# P17 / S04 / 01 — Phase 17 VERIFY (portable-graph-git closeout)

## Metadata
- id: P17-S04-01
- todo_ids: [P17-S04-01]
- role: verify
- skills: [systematic-debugging, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 17 — S01–S03 named DF regressions + **two-clone** git-JSON recipe (why/context/plan; no shared `.trace/`; offline) + carry-forward honesty/Gates/ablation/compat/p0x/x0 + product pkgs — against live packages.

Do **not** create a new planted eval gate. Do **not** trust S01–S03 Notes alone. Do **not** fail for encryption-wontfix, omitted reviews, **absent DF-86 git-hook**, CGO=0 `cmd/trace` tree-sitter FAIL, or S03 `work_state` SQL-only preservation. Do **not** invent research S05 / `plan simulate` / D21+ / hosted MCP without promotion. Do **not** claim P16 DFs or rewrite P16 history.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Write durable evidence, then either:

1. **Pass** → declare **Phase 17 VERIFY PASS / portable-graph-git green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (`P17-S04-01a` / `01b` / +`01c`).

No product features on this row **except**: implement **`TestPortableGraphTwoCloneWhyContextPlan`** if absent at preflight (VERIFY remediation — do not weaken bar) + spawn remediations if a bar fails.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Sibling REVIEW-NOTES: [S01](../scope-01-seed-export/REVIEW-NOTES.md), [S02](../scope-02-commit-convention/REVIEW-NOTES.md), [S03](../scope-03-idempotent-import/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- Pattern: Phase 16 VERIFY [`../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/01-verify.md`](../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md) — must be **FINAL**

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify). Unattended: no Plan-mode switch.

## Locked defaults (FINAL — P17-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Portable-graph-git closeout (S01–S03 named DFs + two-clone recipe) — **not** a new `evals/*` planted gate |
| Path lock | **`trace/graph.json`** |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| Compat checklist | **Green** — `TestCompatibilitySecurityChecklist` (mig ceiling as of P16 close) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist — Phase 01 dry-run is regression-only |
| Two-clone recipe | **Fail bar** — `TestPortableGraphTwoCloneWhyContextPlan` (CGO=1). Implement if absent at preflight |
| DF-86 | **Non-fail** — `trace install git-hook` absent OK; grep and record |
| Residuals (non-blocking) | Encryption wontfix; reviews omitted; DF-86 hook absent; CGO=0 `cmd/trace` FAIL; S03 work_state SQL-only |
| VerifiedFact | Still **out** |
| Product Go | **Forbidden** except two-clone named test if absent + spawn remediation |
| MCP / daemon / HTTP | Still forbidden; **no** new MCP seed tool; local stdio only |
| Full bar | Product packages `./cmd\|internal\|evals` PASS; R3 graphify space FAIL OK; R4 CGO0 analyzers / CGO0 `cmd/trace` FAIL OK |
| Successor | **`no successor`** — research S05 / plan simulate / D21+ / hosted MCP stay off-board unless Notes explicitly promote |
| Evidence artifact | **`VERIFY-NOTES.md`** in this folder |

### Evidence table (fill in VERIFY-NOTES.md)

| Bucket | Must prove |
|--------|------------|
| S01 DF-80/84/85 | `TestSeedExportRoundTrip` (plan tree + findings); `TestSeedExportOmitsDeniedSurfaces`; `TestSeedExportWritesExportedAtCommit`; P16 import keepers PASS |
| S02 DF-82/85 | `TestHelpSeedExportPath`; path `trace/graph.json`; `.gitignore` `.trace/` only; actor ≠ auth docs (`TestAsOperatorFlagIdentityDocs`); DF-28 handoff keeper |
| S03 DF-81/83/84 | `TestSeedImportIdempotent`; `TestSeedImportDuplicateLinksNoOp`; `TestSeedImportSameIdLastWins`; `TestSeedImportPlanTreeIdempotent` PASS |
| Two-clone recipe | Two temp dirs, **no shared `.trace/`** with source; each init → import → index → why + context + plan show; **offline** |
| DF-86 | Hook absent — **non-fail**; grep evidence recorded |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false`; product pkgs |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Residuals OK | Encryption / reviews omitted / DF-86 / CGO0 cmd/trace **not** fail criteria |
| Laws | No daemon/HTTP; no new MCP seed tool; no hosted MCP; P16 history intact |

### Locked verify commands

Copy from sibling `00-PLANNER.md` **Locked verify command set (FINAL)** — do not invent a shorter substitute.

```bash
# --- S01 DF-80/84/85 export + P16 import keepers ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit|TestSeedImportAndWhy|TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestSeedImportFromIDAliases|TestSeedImportRelativePathAgainstC|TestSeedImportMissingEndpointsMessage'

# --- S02 DF-82/85 help/path + actor≠auth keepers ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSeedExportPath|TestHelpHandoffSoT|TestAsOperatorFlagIdentityDocs|TestSeedExport'

# --- S03 DF-81/83/84 idempotent import + carry-forward ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestHelpSeedExportPath'

# --- Two-clone git-JSON recipe (fail bar) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run TestPortableGraphTwoCloneWhyContextPlan

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
# DF-86 grep: no install git-hook in *.go
# .gitignore: .trace/ only; trace/graph.json not ignored
# Gate C artifact inspect: dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
```

Two-clone shell recipe (secondary corroboration — copy from `00-PLANNER.md` **Locked two-clone shell recipe**; document ids + clone paths in VERIFY-NOTES).

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] S01–S03 evidence is **named tests** — not Notes-only
- [ ] **No** new MCP seed tool; local stdio MCP unchanged
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] `.gitignore` lists `.trace/` only — `trace/graph.json` **not** ignored
- [ ] **DF-86 hook absence does not fail VERIFY**
- [ ] **No research S05 / plan simulate / D21+ scaffold** unless Notes explicitly promote
- [ ] Phase 16 historical `no successor` left intact — do not rewrite P16 history
- [ ] Forward-only: do **not** rewrite Phase 00–16 `done` history

### DR-HANDOFF duties (this row + S04-02)

Per protocol Phase handoff + [`DR-HANDOFF.md`](../../DR-HANDOFF.md). On green → record **`no successor`**.

| Who | Duty |
|-----|------|
| **P17-S04-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence table + two-clone recipe + DF-86 non-fail; explicitly record **DR-HANDOFF = `no successor`** (start). Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) status toward closed. |
| **P17-S04-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 17 complete only then. |

**Counterfactual:** If primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 17 complete.

**Spawn policy (fail):** insert immediately below this board row:

| ID | Role |
|----|------|
| `P17-S04-01a` | implement remediation (full prompt) |
| `P17-S04-01b` | review remediation |
| `P17-S04-01c` | re-VERIFY (optional if needed after 01b) |

## Board rights
Verify: **status + notes** on `P17-S04-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails. Do **not** rewrite Phase 17 `done` history. Do **not** mark `P17-S04-02` done.

## Preflight / Plan
1. Re-read this prompt + board row + S01–S03 REVIEW-NOTES + locks above.
2. Confirm `00-PLANNER.md` is **FINAL**.
3. Check whether `TestPortableGraphTwoCloneWhyContextPlan` exists — if absent, plan to implement (VERIFY remediation).
4. Plan: run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. S01–S03 named DF regressions (required)

| Check | Expect |
|-------|--------|
| S01 DF-80/84/85 | Export round-trip (plan tree); omit denied; `exported_at_commit`; P16 import keepers PASS |
| S02 DF-82/85 | Help path lock; handoff keeper; actor≠auth; `.gitignore` `.trace/` only |
| S03 DF-81/83/84 | Idempotent import; duplicate links no-op; last-wins entity+plan; plan tree idempotent PASS |

### B. Two-clone recipe (required — fail bar)

| Check | Expect |
|-------|--------|
| Named test | `TestPortableGraphTwoCloneWhyContextPlan` PASS (implement if absent) |
| Isolation | Two dirs; **neither** uses source `.trace/`; each has own `.trace/` |
| Workflow | init → seed import → index → why + context + plan show |
| Offline | No account, no HTTP, no MCP server |

### C. Carry-forward gates (required)

| Check | Expect |
|-------|--------|
| Honesty A/B/C + Gate G | PASS |
| Gate E / F / ablation | PASS |
| Gate H + compat checklist | PASS |
| p0x 7/7 + x0 | PASS |
| Gate C artifacts | `dry_run:false` N=3 intact — inspect only |
| Product pkgs | `./cmd\|internal\|evals` PASS (R3 graphify / R4 CGO0 FAIL OK) |

### D. Residuals (must record, must not fail)

| Residual | Disposition | VERIFY rule |
|----------|-------------|-------------|
| Encryption-as-git | wontfix | Note only — do **not** fail |
| Reviews omitted from export | out | Note only — do **not** fail |
| DF-86 git-hook | deferred | Absent OK — do **not** fail |
| CGO=0 `cmd/trace` | carry-forward | CGO=1 authoritative — do **not** fail |
| S03 work_state preservation | SQL-only test gap | Note only — do **not** fail |

### E. Evidence + handoff

Write `VERIFY-NOTES.md` with:

1. Verdict line (PASS/FAIL)
2. Evidence table (command → result)
3. Two-clone recipe (named test + optional shell corroboration)
4. Law checks + DF-86 grep
5. Residuals / deferrals
6. Explicit **DR-HANDOFF = `no successor`**
7. Stamp [`DR-HANDOFF.md`](../../DR-HANDOFF.md) toward closed (start)

On FAIL: spawn `P17-S04-01a` / `01b` (+ `01c` if needed); do not weaken bars.

## Todo updates
Status + Notes on `P17-S04-01` only. Do not mark `P17-S04-02` done.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn trail)
- [ ] Two-clone recipe proven (named test; shell corroboration in VERIFY-NOTES)
- [ ] DF-86 hook **not** required for PASS
- [ ] `VERIFY-NOTES.md` written with evidence table + law checks
- [ ] DR-HANDOFF **started** = `no successor`
- [ ] Board Notes on `P17-S04-01`; next `P17-S04-02` (or spawn trail)
- [ ] Explicit: dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ H / ≠ checklist
- [ ] P16 history intact; no hosted MCP successor claimed

## Minimal todos
- [ ] Preflight: `TestPortableGraphTwoCloneWhyContextPlan` exists or implement
- [ ] Run S01–S03 named DF regressions (CGO=1 cmd/trace)
- [ ] Run two-clone named test
- [ ] Run carry-forward honesty/Gates/ablation/H/compat/p0x/x0/product pkgs
- [ ] Grep DF-86 absent; `.gitignore` check
- [ ] Write VERIFY-NOTES + start DR-HANDOFF = `no successor`
- [ ] Board update (or spawn on fail)

## Out of scope
- Product features / new MCP tools / DF-86 hook (except two-clone test + spawn)
- Scaffolding research S05 / plan simulate / D21+ / hosted MCP without promotion
- Re-scoring Gate C
- Rewriting Phase 00–16 history
- Using CGO0 `./cmd/trace/...` as a fail bar
