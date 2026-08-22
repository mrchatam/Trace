# P35-S03-02 — DR-HANDOFF Phase 35 close

## Metadata
- id: P35-S03-02
- todo_ids: [P35-S03-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 35 DR-HANDOFF** with explicit successor (**never TBD**). Default successor **`no successor`** unless VERIFY residuals require a thin follow-on phase scaffold (per agent-loop-protocol **Phase handoff**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 35 complete when this row is `done`. **No product code.** Do **not** implement a successor in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-35.md](../../../../TODO/phase-35.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p35-s03-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Design locks | `DESIGN-LOCKS.md` |
| Phase board | `docs/TODO/phase-35.md` |
| Prior scope artifacts | S00 INVESTIGATION; S01 PLAN; S02 implement + review Notes |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — VERIFY floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; verify blocks 0–5 green per VERIFY-NOTES + independent spot-check; live binds ≠ Step1; URL override evidenced; residuals listed |
| Default successor | **`no successor`** — P1→P2→P3a shipped in GUI; known nits (display vs pick truncation, `listTasksForPick` max-pages, HTTP pagination deferred, Placement A deferred) are **not** a new phase by default |
| Thin follow-on exception | Only if VERIFY/spot-check leaves **blocking** product residuals that need a new phase theme **or** human promotes a follow-on — then scaffold per protocol (below) |
| Hosted SaaS / cloud | **Not** a Phase 35 successor — separate product/repo |
| Regression path | Spawn `P35-S03-02a` implement + `02b` review; **do not** close Phase 35 |
| Must not | Leave `Successor decision: TBD`; invent HTTP pagination phase solely from residual notes; rewrite S00–S02 `done` history; ship product in this row; start implementing a successor |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 35 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S03-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; residuals only as listed (truncation / max-pages / HTTP pagination deferred / Placement A) | **`no successor`** | Close DR-HANDOFF; mark Phase 35 **done** in TODO/AGENTS; idle orchestrator paste |
| Live still binds Step1; unit fail; override overwritten; dogfood mutated | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b; successor = `pending repair spawn` |
| VERIFY-NOTES missing blocks, live binds unticked, or evidence dir absent | **Do not close** — spawn repair or send back S03-01 | Keep OPEN |
| VERIFY PASS but **blocker** residual needs a new Trace-core theme (e.g. human wants OpenAPI pagination + server current-work **now**) | **Thin follow-on phase** | Close only after **runnable** next-phase scaffold exists (see below); first runnable = that phase’s `00-PHASE-PLANNER` |
| Human wants hosted SaaS next | **Not** Trace core successor | Close with **`no successor`** (or human override naming separate product/repo); do not invent cloud board here |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Successor scaffold expectations (only if thin follow-on chosen)

Per [agent-loop-protocol Phase handoff](../../../../rules/agent-loop-protocol.md), the **closing** phase owns the next-phase scaffold — not a later ad-hoc session. Before marking this row `done` with a follow-on successor:

1. **Folder:** `docs/phases/phase-NN-<slug>/` with at least:
   - `README.md` (goal, scope list, in/out)
   - `00-PHASE-PLANNER.md` (runnable)
   - `DR-HANDOFF.md` (OPEN)
   - Per-scope stubs: `00-PLANNER` / `01-*` / `02-*` / `SCOPE-TODOS.md` (minimal OK)
2. **Board:** `docs/TODO/phase-NN.md` with phase planner as **first pending** row after Phase 35’s last `done` row
3. **Index:** link in `docs/TODO.md` phase boards table; orchestrator paste → Active phase NN / next runnable `PNN-00`
4. **AGENTS.md** Current focus → Phase NN
5. Deep tasking belongs to the **next** phase planner — this row delivers a **runnable handoff**, not finished implement prompts

If default **`no successor`**: do **not** invent Phase 36 stubs. Update TODO/AGENTS to Phase 35 **complete** / idle (Phases 00–35 closed; cloud remains separate product).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-35-active-task-selection/scopes/scope-03-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p35-s03-01-verify/evidence || ls experiments/runs/ | grep p35-s03-01
cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts
# Spot-read VERIFY-NOTES live bind table: Overview + Loop ≠ 33247e2d-…; prefer 99d8fb92-…
# Spot-read override evidence file under evidence/
# Optional live re-open (read-only):
#   trace gui -C "/home/ali/Desktop/feet seller telegram app"
```

Confirm VERIFY-NOTES: overall PASS; Blocks 0–5; live binds table filled; override evidenced; residuals listed; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 investigate (`INVESTIGATION.md`)
- [ ] S01 plan (`PLAN.md`)
- [ ] S02 implement + tests + review
- [ ] S03 VERIFY (`VERIFY-NOTES.md`) + successor / CLOSED

### Residuals that do **not** force a successor (default)

| Residual | Why not a phase by default |
|----------|----------------------------|
| Display `limit: 100` vs pick completeness | Bind path correct; UI option list honesty is follow-on polish |
| `listTasksForPick` no max-pages | Low risk until pathological cursors |
| HTTP `limit`/`cursor` still ignored / unimplemented in OpenAPI | Explicitly deferred in PLAN; client already page-ready |
| Placement A Go current-work API | Law-19 upgrade path; not Phase 35 ship gate |
| Optional TRACE_TASK_ID docs paragraph | Non-blocking |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# DR-HANDOFF — Phase 35

**Status:** CLOSED

| Field | Value |
|-------|-------|
| Opened | 2026-08-21 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 34 CLOSED |
| Theme | Active/current task selection + dogfood test |
| Successor decision | **no successor** |
| Close owner | `P35-S03-02` |
| Verify | Cite VERIFY-NOTES + evidence dir; live ≠ Step1; override evidenced |

## Scope checklist (board SoT)

- [x] S00 investigate (`INVESTIGATION.md`)
- [x] S01 plan (`PLAN.md`)
- [x] S02 implement + tests + review
- [x] S03 VERIFY (`VERIFY-NOTES.md`) + successor / CLOSED

## Evidence pointers

- VERIFY-NOTES: `scopes/scope-03-verify/VERIFY-NOTES.md`
- Evidence: `experiments/runs/YYYY-MM-DD-p35-s03-01-verify/evidence/`
- Unit: `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts`
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### TODO.md + AGENTS.md updates (on APPROVE + `no successor`)

- `docs/TODO.md` — Phase 35 marked complete / idle paste (Phases 00–35 closed).
- `AGENTS.md` Current focus — Phase 35 complete; no active phase unless human promotes next.
- Phase board: all rows `done` after this row.

## Exit criteria

- [ ] Independent spot-check of VERIFY-NOTES + unit re-run (minimum)
- [ ] `DR-HANDOFF.md` **CLOSED** with date + evidence pointers + successor **never TBD**
- [ ] Checklist boxes match reality
- [ ] `docs/TODO.md` + `AGENTS.md` current-focus updated
- [ ] Phase 35 board complete after this row
- [ ] No product code

## Minimal todos

1. Confirm S03-01 `done` and artifacts exist.
2. Re-run unit spot-check; read live bind + override evidence.
3. Decide successor per table (`no successor` default).
4. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor.
5. Update `docs/TODO.md` + `AGENTS.md`; mark this board row `done` with Notes.
6. On FAIL: keep OPEN; spawn 02a/02b; do not mark Phase 35 done.

## Board Notes template (fill on completion)

| Field | Value |
|-------|-------|
| VERIFY-NOTES overall | |
| Spot-check unit exit | |
| Live ≠ Step1 confirmed | |
| Override evidenced | |
| Successor | `no successor` \| `pending repair spawn` \| `phase-NN-…` |
| DR-HANDOFF | CLOSED \| OPEN |
| TODO/AGENTS updated | yes/no |

## Next

—(phase complete when APPROVE + CLOSED)
