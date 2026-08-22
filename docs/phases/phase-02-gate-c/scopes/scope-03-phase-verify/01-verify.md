# P02 / S03 / 01 — Phase 02 VERIFY

## Metadata
- id: P02-S03-01
- todo_ids: [P02-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 02 closeout — Gate C decision artifacts + S02 harden regressions + honesty/p0x/x0 bars — against live packages. Do **not** trust S01/S02 Notes alone. Do **not** treat Phase 01 dry-run as Gate C pass.

Write durable evidence, then either:

1. **Pass** → declare **Phase 02 VERIFY PASS** + **start DR-HANDOFF** Phase 03 (`phase-03-progressive-planner` + board `P03-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features. Gate C verdict remains **Go** unless fresh artifact checks contradict recorded evidence (do not silently flip Go→No-Go without evidence + Notes).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 2 Gate C; Phase 3 progressive planner
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0
- Sibling: [GATE-C-NOTES.md](../scope-01-x0-gate-c/GATE-C-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-slice-hardening/REVIEW-NOTES.md)
- Metrics: [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 01 VERIFY [`../../../phase-01-x0-readiness/scopes/scope-05-phase-verify/01-verify.md`](../../../phase-01-x0-readiness/scopes/scope-05-phase-verify/01-verify.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate C verdict (recorded) | **Go** — mean G1 `understanding_accuracy` 0.800 > B0 0.000; kill not fired |
| Gate C evidence | `GATE-C-NOTES.md` + `docs/verification/gate-c-x0/` (`metrics-b0.json`, `metrics-g1.json`, `pins.md`) with `dry_run:false`, N=3/condition |
| Dry-run ≠ Gate C | Phase 01 `TestX0DryRunMetricsB0AndG1` (`dry_run:true`) is regression-only — **not** Gate C pass |
| Honesty | `evals/honesty` Paths A/B/C; no `AllowDoneWithoutReview` in proof |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 packages | Dry-run harness + Gate C recorded Mode-B packs remain green |
| S02 GC-01 | `TestWhyTaskIncludesDiscoveryPlanChange`, `TestTaskContextIncludesDiscoveryPlanChange` |
| S02 GC-02 | `TestFixtureReadmeHasNoGTUUIDOracle`; fixture tree hash pin below |
| Fixture hash pin | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals | **GC-03/04 still deferred** unless VERIFY Notes promote with evidence |
| Mode-B packs | Historical evidence — do **not** require q3 pack rewrite for VERIFY pass |
| Residual (non-blocking) | Global attach of all `discovery_causes_plan_change` edges on every task Expand — see S02 REVIEW-NOTES |
| Product Go | **Forbidden** — verify + evidence + handoff scaffold + spawn prompts only |
| MCP | Optional; Gate C was CLI Mode-B — do not require MCP for VERIFY pass |
| Daemon / HTTP / embeddings | Still forbidden as primary |

### Locked verify commands

```bash
# Full regression bar
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./... -count=1

# S02 GC-01 surfaces
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... \
  -run 'TestWhyTaskIncludesDiscoveryPlanChange|TestTaskContextIncludesDiscoveryPlanChange' -count=1

# S02 GC-02 guard
CGO_ENABLED=1 go test ./evals/x0/ -run TestFixtureReadmeHasNoGTUUIDOracle -count=1

# Fixture content hash (must match pins)
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
# Inspect Gate C artifacts (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] Gate C evidence is `dry_run:false` artifacts — **not** Phase 01 dry-run alone
- [ ] Mode-B packs not falsified to invent q3 pass
- [ ] Embeddings still absent
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. Gate C outcome is **Go** → scaffold Phase 03. Do **not** leave Phase 03 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P02-S03-01 (this VERIFY)** | On pass: **Create** Phase 03 scaffold under `docs/phases/phase-03-progressive-planner/` (checklist below) **and** append Phase 03 board section to `docs/TODO.md` with `P03-00` as first pending row after Phase 02’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P02-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 02 complete only when handoff + VERIFY evidence agree. |

#### Phase 03 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-03-progressive-planner/
  README.md                 # goal = progressive planner (A_PROJECT_PLAN Phase 3)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 03** section after Phase 02 rows; first pending ID **`P03-00`** → `phase-03-progressive-planner/00-PHASE-PLANNER.md`. Do **not** start executing Phase 03 until Phase 02 S03-02 is `done`.

Deep tasking stays with Phase 03’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If fresh checks overturned Go → No-Go, record explicit stop / `no successor` in VERIFY-NOTES instead of scaffolding (user override may reopen later). That path is **not** expected given S01 APPROVE + S02 harden.

## Board rights
Verify: **status + notes** on `P02-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 03 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 02 `done` history. Do **not** mark `P02-S03-02` or Phase 03 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + GATE-C-NOTES + S02 REVIEW-NOTES residuals.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x}`, Gate C metrics under `docs/verification/gate-c-x0/`.
3. Plan: artifact check → run commands → fill evidence → pass→VERIFY-NOTES + Phase 03 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Gate C artifact re-check (required — not dry-run alone)

Confirm all of:

| Check | Expect |
|-------|--------|
| `GATE-C-NOTES.md` verdict | **Go** |
| Kill criteria | Not fired (G1 mean 0.800 > B0 0.000) |
| Metrics files exist | `docs/verification/gate-c-x0/metrics-b0.json`, `metrics-g1.json` |
| `dry_run` | `false` in both metrics |
| N | 3 runs per condition (or equivalent quality objects) |
| Pins | `docs/verification/gate-c-x0/pins.md` present; fixture hash aligns with GC-02 pin after S02 |
| Honesty of claim | Notes do **not** treat Phase 01 dry-run as Gate C pass |

Record pass/fail in VERIFY-NOTES evidence table.

### B. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Fixture hash must match pin. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### C. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate C verdict + kill | | `GATE-C-NOTES.md` Go; G1 0.800 > B0 0.000 |
| Gate C metrics `dry_run:false` N=3 | | paths under `docs/verification/gate-c-x0/` |
| Dry-run ≠ Gate C | | explicit: Phase 01 dry-run not used as pass |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| S02 GC-01 | | Why + TaskContext discovery↔plan_change tests |
| S02 GC-02 | | README guard + hash `15fe50a1…` |
| `go test ./...` | | PASS/FAIL |
| Deferrals GC-03/04 | | still deferred (or promoted with reason) |
| Mode-B packs | | historical; no rewrite required |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19 |
| Residual DPC-global | | noted non-blocking (S02 REVIEW-NOTES) |
| DR-HANDOFF | | Phase 03 path created? board `P03-00` present? |

Also record: date, go/CGO note, residuals carried forward.

### D. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 02 VERIFY PASS / Gate C closeout green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: Gate C **Go** re-confirmed; Phase 01 dry-run ≠ Gate C.
2. **Start DR-HANDOFF:** create `docs/phases/phase-03-progressive-planner/` per checklist + append Phase 03 section / `P03-00` to `docs/TODO.md`.
3. Board Notes: short “Gate C Go re-check + honesty/p0x/x0/S02 PASS; Phase 03 scaffold started; see VERIFY-NOTES.md; pending P02-S03-02 handoff close”.
4. Mark `P02-S03-01` **done**. Do **not** mark S03-02 done.

### E. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P02-S03-01a` (implement) + `P02-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent Gate C No-Go, or rewrite Mode-B packs to fake q3.
4. **Forbidden “fixes”:** claiming Gate C from dry-run alone; requiring MCP for Gate C; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence.

### Spawn ID convention

```text
… P02-S03-01  (this VERIFY)
… P02-S03-01a (remediation implement)   ← insert immediately below
… P02-S03-01b (remediation review)
… P02-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P02-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing progressive planner product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs for q3
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 03 implement wave before S03-02 closes

## Todo updates
Status + Notes on `P02-S03-01`; spawn rows if needed; Phase 03 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent Gate C artifact check + honesty + p0x + x0 + S02 GC-01/02 + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes dry-run≠Gate C + deferrals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified
- [ ] Either gates green **and** Phase 03 scaffold started (folder + `P03-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: GATE-C-NOTES + metrics + S02 residuals + hash pin
- [ ] Re-check Gate C artifacts (`dry_run:false`, N=3, kill, Go)
- [ ] Run honesty / p0x / x0 / GC-01 / GC-02 / hash / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-03-progressive-planner` + board `P03-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
