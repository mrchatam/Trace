# P01 / S05 / 01 — Phase 01 VERIFY

## Metadata
- id: P01-S05-01
- todo_ids: [P01-S05-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 01 readiness for Gate C work against live packages — do **not** trust S01–S04 Notes alone. Write durable evidence, either (a) declare **Phase 01 VERIFY green** + start **DR-HANDOFF** Phase 02 scaffold, or (b) **spawn forward-only remediations**. No product features; do **not** claim Gate C / product thesis won.

Independently confirm:

1. **Honesty** — `evals/honesty` fail-closed Paths A/B/C green (S02).
2. **X0 dry-run** — `evals/x0` emits schema-valid temp **B0 + G1** metrics (`dry_run: true`; G1 used CLI `why`/`context`).
3. **MCP** — S04 surface documented (thin adapter; G19; checklist tools present).
4. **P0-X regression** — `evals/p0x` still **7/7 PASS**.
5. Evidence file **`VERIFY-NOTES.md`** in this folder.
6. **DR-HANDOFF started** — Phase 02 folder + planner + stubs + board rows created (or advanced enough that S05-02 can finish the gate).

**Not** Gate C pass. **Not** “G1 beats B0.” Dry-run readiness only.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT, DR-SURFACE
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 1 validation gate; Phase 2 Gate C
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0 (B0 vs G1)
- Sibling REVIEW-NOTES: S02 honesty, S03 x0, S04 mcp (residuals ≠ free passes)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 00 VERIFY [`../../../phase-00-foundation/scopes/scope-09-phase-verify/01-verify.md`](../../../phase-00-foundation/scopes/scope-09-phase-verify/01-verify.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Bar | Honesty green + X0 dry-run B0/G1 metrics + MCP checklist + p0x 7/7 regression |
| Gate C / agent scoring | **Out** — Phase 02 (`phase-02-gate-c`) |
| Product Go | **Forbidden** on this row — verify + evidence + handoff scaffold + spawn prompts only |
| Honesty command | `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` (also re-check under full CGO suite below) |
| X0 command | `CGO_ENABLED=1 go test ./evals/x0/... -count=1` |
| P0-X command | `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` |
| Full regression | `CGO_ENABLED=1 go test ./... -count=1` |
| Honesty entry | `evals/honesty` → `TestHonestyFailClosedPlantedClaim` Paths A/B/C; **no** `AllowDoneWithoutReview` in proof |
| X0 entry | `evals/x0` → `TestX0DryRunMetricsB0AndG1`; schema `evals/x0/schema.json` v1 |
| X0 metrics | Temp only: `metrics-b0.json` + `metrics-g1.json` with `dry_run: true` (not committed) |
| X0 conditions | **B0** = no `trace why`/`context`; **G1** = live CLI `why` + `context` on task `22222222-2222-2222-2222-222222222222` |
| Fixture / seed | `fixtures/x0` + abs `seed/gt.json` v1 (same UUID map as P0-X) |
| MCP checklist | Tools: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`. Binary `cmd/trace-mcp` (`bin/trace-mcp`); official `go-sdk` **v1.4.0** **stdio only**. G19: no domain/retrieval/compiler fork in MCP. **Do not** require re-running X0 via MCP (DR-AGENT) |
| MCP smoke (optional) | `CGO_ENABLED=0 go test ./internal/mcp/... -count=1` — docs checklist is mandatory; package test is strong evidence |
| Evidence file | Write **`VERIFY-NOTES.md`** here (durable). Board Notes = short summary + pointer |
| Depends on | S01–S04 all `done` / APPROVED. Re-run suites anyway |
| Human | Optional spot-check only if automation cannot run; do not invent Gate C scores |

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] X0 remains CLI-path (MCP not required for dry-run)
- [ ] No Gate C “G1 beats B0” / product-thesis claim in VERIFY-NOTES
- [ ] Embeddings still absent

### DR-HANDOFF duties (this row + S05-02)

Per protocol Phase handoff + **DR-HANDOFF**: do **not** leave Phase 02 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P01-S05-01 (this VERIFY)** | **Create** Phase 02 scaffold under `docs/phases/phase-02-gate-c/` (see checklist below) **and** append Phase 02 board section to `docs/TODO.md` with `P02-00` as first pending row after Phase 01’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P01-S05-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 01 complete only when handoff + VERIFY evidence agree. |

#### Phase 02 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-02-gate-c/
  README.md                 # goal = Gate C evaluation & slice hardening (A_PROJECT_PLAN Phase 2)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 02** section after Phase 01 rows; first pending ID **`P02-00`** → `phase-02-gate-c/00-PHASE-PLANNER.md`. Do **not** start executing Phase 02 until Phase 01 S05-02 is `done`.

Deep tasking stays with Phase 02’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

## Board rights
Verify: **status + notes** on `P01-S05-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 02 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 01 `done` history. Do **not** mark `P01-S05-02` or Phase 02 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S02/S03/S04 REVIEW-NOTES residuals.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x}`, `internal/mcp`, `cmd/trace-mcp`.
3. Plan: run commands → fill evidence → pass→VERIFY-NOTES + Phase 02 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`:

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, not a substitute for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=1 go test ./evals/x0/... -count=1 -v -run TestX0DryRunMetricsB0AndG1
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
CGO_ENABLED=0 go test ./internal/mcp/... -count=1
```

All required commands must exit 0 for a green gate. If a package skips because CGO/binary build failed, treat as **gate fail** (not skip). Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### B. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| X0 dry-run B0 metrics | | schema-valid; `dry_run:true`; no why/context in tools |
| X0 dry-run G1 metrics | | schema-valid; `dry_run:true`; why+context used on task `2222…` |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| `go test ./...` | | PASS/FAIL |
| MCP checklist | | six tools + stdio `trace-mcp` + G19 note (or N/A with reason — **not** expected) |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; no Gate C claim |
| DR-HANDOFF | | Phase 02 path created? board `P02-00` present? |

Also record: date, go/CGO note, residuals carried forward (non-blocking if primary asserts green).

### C. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 01 VERIFY PASS / ready for Gate C phase**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: **not Gate C**.
2. **Start DR-HANDOFF:** create `docs/phases/phase-02-gate-c/` per checklist + append Phase 02 section / `P02-00` to `docs/TODO.md`.
3. Board Notes: short “honesty+x0 B0/G1+p0x+MCP PASS; Phase 02 scaffold started; see VERIFY-NOTES.md; pending P01-S05-02 handoff close”.
4. Mark `P01-S05-01` **done**. Do **not** mark S05-02 done.

### D. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P01-S05-01a` (implement) + `P01-S05-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer (domain/CLI/harness/MCP adapter/docs) — **do not** weaken bars or claim Gate C.
4. **Forbidden “fixes”:** declaring Gate C pass; requiring MCP for X0; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S04 prompts.

### Spawn ID convention

```text
… P01-S05-01  (this VERIFY)
… P01-S05-01a (remediation implement)   ← insert immediately below
… P01-S05-01b (remediation review)
… P01-S05-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P01-S05-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing product features “while you’re here”
- Running real multi-model agent Gate C / Go-No-Go scoring
- Declaring A1 / product thesis validated
- Merging honesty into X0 metrics
- Starting Phase 02 implement wave before S05-02 closes

## Todo updates
Status + Notes on `P01-S05-01`; spawn rows if needed; Phase 02 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent honesty + x0 + p0x + `./...` evidence recorded in `VERIFY-NOTES.md`
- [ ] MCP checklist (six tools / stdio / G19) noted with evidence
- [ ] Law checks recorded; **no** Gate C pass claimed
- [ ] Either gates green **and** Phase 02 scaffold started (folder + `P02-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: tree + S02/S03/S04 residuals + UUID/task map
- [ ] Run honesty / x0 / p0x / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-02-gate-c` + board `P02-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
