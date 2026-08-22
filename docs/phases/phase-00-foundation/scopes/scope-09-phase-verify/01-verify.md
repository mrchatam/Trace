# P00 / S09 / 01 — Phase 00 VERIFY / P0 close

## Metadata
- id: P00-S09-01
- todo_ids: [P00-S09-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove **DR-P0X 7/7** against live `fixtures/x0` + `evals/p0x`, write durable evidence, and either (a) declare **P0 closable** or (b) **spawn forward-only remediations** for any failing criterion. Do **not** trust S08 Notes alone. No product features; no MCP/daemon/HTTP “fixes.”

On full pass: record closeout in board Notes + `VERIFY-NOTES.md` + init registers (`F_QUESTION_LEDGER` / `E_ASSUMPTION_REGISTER` as specified below). Phase 01 stays blocked until this row **and** `P00-S09-02` are `done`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — VERIFY may spawn; forward-only
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — P0-X 7/7 bar
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G12 incremental; no blob dump
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment P0-X
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-P0, DR-P0X, DR-INCREMENTAL
- [docs/init/F_QUESTION_LEDGER.md](../../../../init/F_QUESTION_LEDGER.md) — Q-P0X-PASS / Q-P0-DONE
- [docs/init/E_ASSUMPTION_REGISTER.md](../../../../init/E_ASSUMPTION_REGISTER.md) — A15, A6, A17
- [../scope-08-fixture-p0x/01-fixture.md](../scope-08-fixture-p0x/01-fixture.md) — 7↔assert map + query set
- [../scope-08-fixture-p0x/REVIEW-NOTES.md](../scope-08-fixture-p0x/REVIEW-NOTES.md) — residuals to re-check (not excuses)
- Phase README: [../../README.md](../../README.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Bar | **DR-P0X 7/7** including incremental (#7); full-rebuild-on-any-change = fail |
| Agent Gate C / X0 | **Not** this row (Phase 01+) |
| MCP / daemon / HTTP | Still **forbidden**; must not be added to “fix” a fail |
| Product Go | **Forbidden** on this row — verify + docs/evidence + spawn prompts only |
| Primary command | `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` from module root |
| Regression command | `CGO_ENABLED=1 go test ./... -count=1` |
| Fixture | `fixtures/x0` (Apache-2.0 synthetic; TS `src/greeter.ts` + Py `src/math_util.py`) |
| Seed | `fixtures/x0/seed/gt.json` v1; harness uses **absolute** seed path (`-C` does **not** rewrite it) |
| Live GT UUIDs | goal=`11111111-1111-1111-1111-111111111111` · task=`22222222-2222-2222-2222-222222222222` · decision=`33333333-3333-3333-3333-333333333333` · discovery=`44444444-4444-4444-4444-444444444444` · plan_change=`55555555-5555-5555-5555-555555555555` |
| Harness entry | `evals/p0x/p0x_test.go` → `TestP0XAllCriteria` subtests `criterion-1`…`criterion-7` + nested queries under `criterion-6-queries` |
| Metrics schema | Temp `metrics-p0x.json`: `{"ok": true, "criteria": {"1": true … "7": true}, "timings_ms"?: …}` — written inside test temp (not required in repo) |
| Evidence file | Write **`VERIFY-NOTES.md`** in this scope folder (durable). Board Notes = short summary + pointer |
| Human | Optional spot-check only if you flag ambiguity; do not block on human unless automation cannot run |
| Depends on | S01–S08 all `done` (incl. S08 harness claiming 7/7). Re-run anyway |

### P0-X criteria → what VERIFY must confirm

Use the S08 locked map; independently confirm each via harness PASS **and** a one-line evidence note (subtest name / assertion gist):

| # | Criterion | Live assert locus (S08 harness) |
|---|-----------|----------------------------------|
| 1 | Goal/Task/Decision/Discovery round-trip + ACTIVE provenance; task work_state after transition | `criterion-1-roundtrip` — Get* + statuses + `IN_PROGRESS` + `goal_id` |
| 2 | Files + minimal symbols/imports (TS **and** Py) | `criterion-2-files-symbols` — both paths; ≥1 symbol **or** import each |
| 3 | `trace why` causal chain + reason codes | `criterion-3-why` — `why task <TaskID>`; non-empty `reason_code`; goal/decision neighbor |
| 4 | `trace context` bounded | `criterion-4-context` — items ≤32; token_limit 4096; trust/`untrusted_data` labeling |
| 5 | Human seed matches GT | `criterion-5-gt-match` — all five UUIDs; `decision_affects_task`; `discovery_causes_plan_change` |
| 6 | ≥5 deterministic understanding queries (no LLM) | `criterion-6-queries` — `why-task`, `why-decision`, `decision-constraint`, `import-or-symbol-neighbor`, `context-boundedness` |
| 7 | Incremental one-file update **without** full fixture rebuild | `criterion-7-incremental` — mutate TS only → `index src/greeter.ts` → Py fingerprint + `content_hash` unchanged; TS gains `greetAgain` |

### Architecture / law checks (must all hold)

- [ ] No MCP server / daemon / HTTP listener introduced on the P0-X path
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Criterion #7 proves **localized** update (sibling isolation) — not “reindex everything and hope”
- [ ] Library packages still do not import `cmd/trace` (G19)
- [ ] S08 residuals noted but **not** treated as free passes: soft `decision-constraint` OR; panic-prone JSON asserts — re-confirm primary paths still hold

## Board rights
Verify role: **status + notes** on `P00-S09-01`; **may spawn** implement+review remediation pairs **immediately below this row** if any of 7 fail or law checks fail. Do **not** rewrite `done` history. Do **not** thicken Phase 01 beyond a one-line blocker note if P0 fails. Do **not** mark Phase 01 `P01-00` runnable yourself — leave blocked until S09-02 also closes the phase.

## Preflight / Plan
1. Re-read this prompt + board row + S08 REVIEW-NOTES residuals.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, `evals/p0x`, `fixtures/x0`.
3. Plan: run harness → fill evidence table → pass→closeout docs **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Independent re-run (required)

From module root `/home/ali/Desktop/Trace` (or workspace root):

```bash
CGO_ENABLED=1 go test ./evals/p0x/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Both must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`. If `evals/p0x` skips because CGO/binary build failed, treat as **gate fail** (not skip).

Optional deeper signal (not a substitute for the package test):

```bash
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
```

### B. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| # | Result (pass/fail) | Evidence (subtest / log gist) |
|---|--------------------|-------------------------------|
| 1 | | |
| 2 | | |
| 3 | | |
| 4 | | |
| 5 | | |
| 6 | | (≥5 named queries listed) |
| 7 | | (sibling isolation gist) |

Also record: date, go/CGO note, whether `go test ./...` passed, quick tree check (no `.trace/` committed under fixtures/evals), MCP/daemon absent.

### C. On **all 7 pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **P0-X PASS / P0 closable**, confidence, evidence table, residuals carried forward (S08 soft OR etc. as **non-blocking** if primary asserts still green).
2. Update init registers (**forward edits to living registers — allowed**):
   - `F_QUESTION_LEDGER.md`: record that **Q-P0-DONE / P0 close** is satisfied by independent S09 VERIFY re-run on date; point at `VERIFY-NOTES.md` + harness commands. Keep Q-P0X-PASS criteria text; add evidence line.
   - `E_ASSUMPTION_REGISTER.md`: set **A15** (and optionally **A17** if structural bar held) Validation=`P00-S09-01` / Status=`VALIDATED` (or equivalent clear status). Leave **A1** (Gate C) as `EXPERIMENT_REQUIRED`.
3. Board Notes: short “7/7 PASS; P0 closable pending P00-S09-02; see VERIFY-NOTES.md”.
4. Mark `P00-S09-01` **done**. Do **not** flip `P01-00` yourself.

### D. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which criterion/law, minimal reproduction.
2. Mark this row **`failed`** **or** keep `in_progress` only while inserting spawns — prefer: set Notes, insert remediations, then set this row **`blocked`** on remediations **or** leave `failed` with spawns pending (orchestrator re-runs VERIFY after fixes). **Preferred pattern:** insert `P00-S09-01a` (implement) + `P00-S09-01b` (review) immediately below this row; set this row **`blocked`** with reason “awaiting 01a/01b”; after remediations `done`, orchestrator re-opens a fresh VERIFY pass (new spawn `P00-S09-01c` re-verify **or** re-run this prompt — prefer **new** `01c` verify row if this row was already `failed`/`blocked` to preserve history).
3. Spawn prompts must be **full** protocol skeletons (metadata, locked defaults, exit criteria, skills/MCPs). Scope remediations to the failing layer (store/analyzers/domain/retrieval/compiler/CLI/fixture/harness) — **not** MCP.
4. **Forbidden “fixes”:** adding daemon/HTTP/MCP; switching to full-rebuild indexer; weakening the 7-point bar; rewriting `done` S08 prompts.

### Spawn ID convention

```text
… P00-S09-01  (this VERIFY)
… P00-S09-01a (remediation implement)   ← insert immediately below
… P00-S09-01b (remediation review)
… P00-S09-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P00-S09-02  (phase review — only after VERIFY finally done)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing features “while you’re here”
- Running / creating `evals/x0` agent Gate C
- Unblocking Phase 01 board rows before S09-02
- Weakening budgets (4096 / 32) or dropping query count below 5

## Todo updates
Status + Notes on `P00-S09-01`; spawn rows if needed; living init registers on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` **and** `./...` evidence recorded
- [ ] Evidence table vs all 7 criteria filled in `VERIFY-NOTES.md`
- [ ] Law checks (no MCP/daemon; no committed `.trace/`; #7 localized; G19) recorded
- [ ] Either **P0 closable** (registers + Notes updated) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: tree + S08 residuals + GT UUID map
- [ ] Run `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` and `./...`
- [ ] Write `VERIFY-NOTES.md` evidence table (7/7 + laws)
- [ ] Pass → update F/E registers + board Notes; Fail → spawn 01a/01b (+01c plan) without MCP
- [ ] Board status + SCOPE-TODOS
