# P08 / S04 / 01 — Phase 08 VERIFY (compat+security checklist)

## Metadata
- id: P08-S04-01
- todo_ids: [P08-S04-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 08 — **compatibility + security checklist** planted harness + S01–S03 surfaces + carry-forward honesty/Gates/ablation/p0x/x0/Gate H/Gate C — against live packages.

**Harness ownership (FINAL):** S01–S03 left **no** `evals/compat`. This VERIFY row **creates** the planted checklist harness under `evals/compat` (like Gate H / Gate F / Gate G planted packages), then proves must-pass items. Do **not** mark checklist blocked-until-harness; do **not** invent commercial multi-model security theater.

Do **not** trust S01–S03 Notes alone. Do **not** reopen Gate C, invent VerifiedFact, declare A1 / product thesis, or scaffold Phase 09.

Write durable evidence, then either:

1. **Pass** → declare **Phase 08 VERIFY PASS / compat+security checklist green** + **start DR-HANDOFF** = **`no successor`**, or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features outside `evals/compat` harness creation (+ spawn remediations if needed).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 8 ends roadmap
- Sibling REVIEW-NOTES: [S01](../scope-01-plugin-apis/REVIEW-NOTES.md), [S02](../scope-02-worktrees/REVIEW-NOTES.md), [S03](../scope-03-production-hardening/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 07 VERIFY [`../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P08-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Validation gate | Compatibility + security checklist (`A_PROJECT_PLAN` Phase 8) — **not** commercial multi-model security theater |
| Checklist path / harness | **`evals/compat`** package; named test **`TestCompatibilitySecurityChecklist`** |
| Schema / metrics | Committed **`schema-compat.json`** **v1**; temp **`metrics-compat.json`** under `t.TempDir()` (schema-validated; `dry_run:false`) |
| Harness ownership | **S04-01 creates** `evals/compat` as VERIFY work (S01–S03 seeded none). Pattern: Gate F/G/H planted packages. **Do not** leave checklist blocked awaiting a prior seed. |
| Allowed Go on this row | **`evals/compat/**` only** (doc.go, schema, planted fixtures, named test). No opportunistic `cmd/` / `internal/` product rewrites unless spawn remediation requires them. |
| Checklist pass bar | Named test PASS + schema-valid temp metrics + all must-pass `*_ok` true + `language_adapter_api_version == 1` |
| S01 surface (must stay green) | `LanguageAdapterAPIVersion=1`; `TestLanguageAdapterAPIVersion`; `TestBuiltinLanguageAdaptersContributionPath`; static adapters; **no** `.so` / megastore; **no** `011_*` from S01 |
| S02 surface (must stay green) | Path-local Abs→`<root>/.trace/`; exclusive `trace.lock` / `ErrLocked`; `TestProjectBindPathLocalIsolation`; `TestConcurrentStoreOpenFailClosed`; **no** walk-up / shared parent; **no** swarm |
| S03 surface (must stay green) | `MigrationStatus` / migrate status (embed max applied; **no** `011_*`); `TestBackupRestoreRoundTrip` (+ lock fail-closed backup/restore); `TestLocalAccessTokenFailClosed`; Abs rebind; `HasBlobLikeColumns` false; token excluded by default; **no** new MCP auth/backup tools; **no** cloud OAuth |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` |
| Gate H | **Green** — `TestPlantedPerfLadderGateH` |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go** |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / **this checklist** — Phase 01 dry-run is regression-only |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out; **100k/1M planted CI ladders** deferred |
| Residuals (non-blocking carry) | S03: argv token exposure; restore TOCTOU window; no `--include-token` unit test; S02: lock CLI exit **2** (`exitFail`); DPC-global; A5 SQLite ceiling |
| VerifiedFact | Still **out** |
| Product Go (non-harness) | **Forbidden** — no indexer rewrite / graph DB / language mega-registry / OAuth |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; checklist does not require new MCP tools |
| Successor | **`no successor`** (`A_PROJECT_PLAN` ends at Phase 8) unless user Notes explicitly promote Phase 09 |

### Metrics schema shape (v1 — lock fields)

`schema-compat.json` **required** keys (`additionalProperties: true` OK):

| Field | Notes |
|-------|-------|
| `schema_version` | const `1` |
| `gate` | const `"compat-security"` |
| `suite` | const `"compat"` |
| `dry_run` | boolean — must be **`false`** for pass artifact |
| `named_test` | `"TestCompatibilitySecurityChecklist"` |
| `language_adapter_api_version` | integer — must equal `1` |
| `language_adapter_api_version_ok` | boolean |
| `path_local_bind_ok` | boolean |
| `trace_lock_ok` | boolean |
| `migrate_status_ok` | boolean |
| `backup_restore_ok` | boolean |
| `no_blob_columns_ok` | boolean |
| `local_auth_fail_closed_ok` | boolean |
| `g19_ok` | boolean |
| `no_daemon_http_primary_ok` | boolean |
| `no_011_mig_ok` | boolean |
| `s01_hooks` / `s02_hooks` / `s03_hooks` | string arrays naming consumed surfaces |

### Locked verify commands

```bash
# --- Compat+security checklist: create harness if absent, then primary named test ---
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# Compat package
CGO_ENABLED=1 go test ./evals/compat/... -count=1

# S01 plugin contribution surface
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestLanguageAdapterAPIVersion|TestBuiltinLanguageAdaptersContributionPath|TestIndexFile|TestDetectLanguage'

# S02 path-local + lock
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestProjectBindPathLocalIsolation|TestConcurrentStoreOpenFailClosed'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInitFailClosedWhenStoreLocked|TestMigrateBackupAuthCLI'

# S03 migrate / backup / auth
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestMigrationStatus|TestBackupRestoreRoundTrip|TestLocalAccessTokenFailClosed|TestBackupFailsWhenLocked|TestRestoreFailsWhenLocked'

# Gate H carry-forward
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Supporting surfaces
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... -count=1

# Full regression bar (includes compat + prior evals)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... ./evals/compat/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -v -run TestCompatibilitySecurityChecklist
test -f evals/compat/schema-compat.json
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
# Spot G19: library packages must not import cmd/trace or cmd/trace-mcp
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19) — evals may drive CLI/analyzers as tests do
- [ ] Checklist evidence is **`evals/compat` `TestCompatibilitySecurityChecklist`** + schema/metrics — not vibes / Notes-only
- [ ] `LanguageAdapterAPIVersion=1`; path-local + `trace.lock`; migrate/backup/auth fail-closed; no BLOBs; no `011_*`
- [ ] No new MCP auth/backup tools; Open inherits auth/lock gates
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Mode-B packs not falsified
- [ ] Embeddings / VerifiedFact / `plan simulate` still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES
- [ ] **No Phase 09 scaffold** unless Notes explicitly promote

### DR-HANDOFF duties (this row + S04-02)

Per protocol Phase handoff + **DR-HANDOFF**. On checklist + regression green → record **`no successor`**. Do **not** create Phase 09 folder/board unless user Notes promote.

| Who | Duty |
|-----|------|
| **P08-S04-01 (this VERIFY)** | On pass: write `VERIFY-NOTES.md` with verdict + evidence; explicitly record **DR-HANDOFF = `no successor`** (`A_PROJECT_PLAN` ends at Phase 8). Do **not** invent Phase 09. |
| **P08-S04-02 (final review)** | **Owns completion check** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` (or an explicitly promoted follow-on is fully scaffolded). Marks Phase 08 complete only then. |

**Counterfactual:** If checklist / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** claim Phase 08 complete; do **not** invent a successor to dodge a red VERIFY.

## Board rights
Verify: **status + notes** on `P08-S04-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create** `evals/compat` harness. Do **not** rewrite Phase 08 `done` history. Do **not** mark `P08-S04-02` done. Do **not** scaffold Phase 09 without explicit promotion.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02/S03 REVIEW-NOTES + checklist locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability,perf}`; **no** prior `evals/compat` expected (create it).
3. Plan: create checklist harness → run locked commands → fill evidence → pass→VERIFY-NOTES + `no successor` **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Checklist harness create + re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/compat/` created; `doc.go` names compat+security checklist |
| Named test | `TestCompatibilitySecurityChecklist` PASS under locked command |
| Schema / metrics | Committed `schema-compat.json` v1; test writes schema-valid temp `metrics-compat.json` (`dry_run:false`) |
| Must-pass flags | All `*_ok` true; `language_adapter_api_version == 1` |
| Honesty of claim | Notes map checklist to this harness — **not** S01–S03 Notes alone |

Record pass/fail in VERIFY-NOTES evidence table.

### B. S01 / S02 / S03 surfaces (required)

| Check | Expect |
|-------|--------|
| S01 API version + contribution | `TestLanguageAdapterAPIVersion` + `TestBuiltinLanguageAdaptersContributionPath` PASS; const `= 1` |
| S02 path-local + lock | Isolation + concurrent Open fail-closed PASS |
| S03 migrate/backup/auth | MigrationStatus + backup↔restore + local-auth fail-closed + lock-on-backup/restore PASS; no `011_*` |

### C. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Gate H | `TestPlantedPerfLadderGateH` PASS |
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate G | `TestHonestyEscapeRateGateGPrelim` PASS |
| Gate E | `TestPlantedDiscoveryReplan` PASS |
| Gate F | `TestPlantedImpactConflictsGateFPrelim` PASS |
| Ablation | `TestPlantedCapabilitySelectionAblation` PASS |
| P0-X 7/7 | `evals/p0x` PASS |
| X0 packages | `./evals/x0/...` PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist | Explicit in VERIFY-NOTES |
| Full suite | Locked full command PASS |

### D. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### E. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Checklist harness created (`evals/compat`) | | path + files |
| `TestCompatibilitySecurityChecklist` | | named test PASS |
| `schema-compat.json` v1 + temp `metrics-compat.json` | | file present; `dry_run:false`; validation |
| `language_adapter_api_version == 1` | | metrics + S01 test |
| Path-local bind + `trace.lock` | | S02 tests / checklist flags |
| Migrate status (no `011_*`) | | S03 / checklist |
| Backup↔restore + no BLOBs + Abs rebind | | S03 / checklist |
| Local auth fail-closed | | S03 / checklist |
| G19 + no daemon/HTTP primary | | law checks |
| S01 contribution-path tests | | analyzers |
| S02 isolation / concurrent Open | | store |
| S03 migrate/backup/auth tests | | store + CLI |
| Gate H | | `TestPlantedPerfLadderGateH` |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| Gate G prelim | | `TestHonestyEscapeRateGateGPrelim` |
| Gate E mini-eval | | `TestPlantedDiscoveryReplan` |
| Gate F prelim | | `TestPlantedImpactConflictsGateFPrelim` |
| Capability ablation | | `TestPlantedCapabilitySelectionAblation` |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C / ≠ F / ≠ G / ≠ ablation / ≠ Gate H / ≠ checklist | | explicit |
| `go test ./...` (+ compat) | | PASS/FAIL |
| Law checks | | no daemon/HTTP; no committed `.trace/`; G19; no BLOBs; no full-rebuild; no commercial security theater |
| Residuals (non-blocking) | | S03 argv token; restore window; S02 exit 2; deferred ladders; … |
| DR-HANDOFF | | **`no successor`** recorded? |

Also record: date, go/CGO note, residuals carried forward.

### F. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 08 VERIFY PASS / compat+security checklist green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: checklist = `evals/compat` named planted test; Gate C artifacts intact; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ this checklist.
2. **Start DR-HANDOFF:** record **`no successor`** in VERIFY-NOTES (and board Notes). Do **not** scaffold Phase 09.
3. Board Notes: short “checklist + S01–S03 + honesty/Gates/ablation/Gate H/p0x/x0 PASS; Gate C intact; DR-HANDOFF=`no successor`; see VERIFY-NOTES.md; pending P08-S04-02 handoff close”.
4. Mark `P08-S04-01` **done**. Do **not** mark S04-02 done.

### G. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P08-S04-01a` (implement) + `P08-S04-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent checklist pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming checklist without `TestCompatibilitySecurityChecklist`; claiming Gate C from dry-run; requiring new MCP for checklist; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S03 prompts; promoting GC-03/04 without evidence; inventing Phase 09 to dodge FAIL; inventing VerifiedFact; shipping commercial multi-model security theater / full-rebuild indexer; introducing `011_*` without justification documented in Notes.

### Spawn ID convention

```text
… P08-S04-01  (this VERIFY)
… P08-S04-01a (remediation implement)   ← insert immediately below
… P08-S04-01b (remediation review)
… P08-S04-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P08-S04-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Scaffolding Phase 09 / inventing a successor without explicit promotion
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Expanding checklist into commercial multi-model security theater
- Introducing VerifiedFact promotion engine / `plan simulate` / graph DB / OAuth / daemon
- Product Go outside `evals/compat` (unless spawn remediation)

## Todo updates
Status + Notes on `P08-S04-01`; spawn rows if needed; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] `evals/compat` checklist harness created + `TestCompatibilitySecurityChecklist` + schema/metrics documented
- [ ] S01–S03 surfaces re-proved
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + p0x + x0 + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes checklist path + dry-run≠… + residuals + **`no successor`** handoff
- [ ] Law checks recorded; Mode-B packs not falsified; no commercial security theater
- [ ] Either gates green **and** DR-HANDOFF started as **`no successor`** **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No non-harness product feature Go on this row

## Minimal todos
- [ ] Preflight: S01–S03 REVIEW-NOTES + Gate C metrics + confirm `evals/compat` absent → create
- [ ] Create planted checklist harness + `schema-compat.json` v1
- [ ] Run `TestCompatibilitySecurityChecklist` + locked S01–S03 / carry-forward commands
- [ ] Write `VERIFY-NOTES.md` evidence table (incl. `no successor`)
- [ ] Pass → board Notes + handoff start; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
