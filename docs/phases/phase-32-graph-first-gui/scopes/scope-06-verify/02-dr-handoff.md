# P32-S06-02 — DR-HANDOFF Phase 32 close

## Metadata
- id: P32-S06-02
- todo_ids: [P32-S06-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S06-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 32 DR-HANDOFF** with explicit successor (**never TBD**). Default successor **`no successor`** unless VERIFY residuals require a thin follow-on phase scaffold (per agent-loop-protocol **Phase handoff**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 32 complete when this row is `done`. **No product code.** Do **not** implement a successor in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S06-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S06-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-32.md](../../../../TODO/phase-32.md)
- [AGENTS.md](../../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S06-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-06-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p32-s06-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Design locks | `DESIGN-LOCKS.md` |
| Port open item | `OPEN-PORT-MULTI.md` |
| Phase board | `docs/TODO/phase-32.md` |
| Prior scope artifacts | S00 `RESEARCH.md`, S01 `UX-IA.md`, S02 `NO-GAPS.md` |

## Locked DR-HANDOFF close policy (FINAL — S06-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S06-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S06-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S06-01 `done`; verify blocks 0–5 green per VERIFY-NOTES + independent spot-check; **P32-PORT** ticked (#1 + docs; #2 deferred OK) |
| Default successor | **`no successor`** — explorer bar shipped; known nits (#2 auto-port, chrome box-shadow, screenshot pipeline) are deferred/optional, **not** a new phase by default |
| Thin follow-on exception | Only if VERIFY/spot-check leaves **blocking** product residuals that need a new phase theme **or** human promotes a follow-on — then scaffold per protocol (below) |
| Cloud / hosted SaaS | **Not** a Phase 32 successor — separate product/repo |
| Regression path | Spawn `P32-S06-02a` implement + `02b` review; **do not** close Phase 32 |
| Must not | Leave `Successor decision: TBD`; treat deferred #2 as requiring a phase without human/blocker evidence; rewrite S00–S05 `done` history; ship product in this row; start implementing a successor |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 32 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S06-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; residuals only as listed (PORT #2 defer, craft nits, no screenshots) | **`no successor`** | Close DR-HANDOFF; mark Phase 32 **done** in TODO/AGENTS; idle orchestrator paste |
| Build / e2e / PORT #1 / Laws / graph-home FAIL | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b; successor = `pending repair spawn` |
| VERIFY-NOTES missing blocks, P32-PORT unticked, or evidence dir absent | **Do not close** — spawn repair or send back S06-01 | Keep OPEN |
| VERIFY PASS but **blocker** residual needs a new Trace-core theme (human-promoted or clearly beyond deferred nits) | **Thin follow-on phase** | Close only after **runnable** next-phase scaffold exists (see below); first runnable = that phase’s `00-PHASE-PLANNER` |
| Human wants hosted SaaS next | **Not** Trace core successor | Close with **`no successor`** (or human override naming separate product/repo); do not invent cloud board here |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Successor scaffold expectations (only if thin follow-on chosen)

Per [agent-loop-protocol Phase handoff](../../../../rules/agent-loop-protocol.md), the **closing** phase owns the next-phase scaffold — not a later ad-hoc session. Before marking this row `done` with a follow-on successor:

1. **Folder:** `docs/phases/phase-NN-<slug>/` with at least:
   - `README.md` (goal, scope list, in/out)
   - `00-PHASE-PLANNER.md` (runnable)
   - `DR-HANDOFF.md` (OPEN)
   - Per-scope stubs: `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
2. **Board:** `docs/TODO/phase-NN.md` with phase planner as **first pending** row after Phase 32’s last `done` row
3. **Index:** link in `docs/TODO.md` phase boards table; orchestrator paste → Active phase NN / next runnable `PNN-00`
4. **AGENTS.md** Current focus → Phase NN
5. Deep tasking belongs to the **next** phase planner — this row delivers a **runnable handoff**, not finished implement prompts

If default **`no successor`**: do **not** invent Phase 33 stubs. Update TODO/AGENTS to Phase 32 **complete** / idle (Phases 00–32 closed; cloud remains separate product).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-32-graph-first-gui/scopes/scope-06-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p32-s06-01-verify/evidence || ls experiments/runs/ | grep p32-s06-01
(cd web && npm run build)
go test ./internal/httpapi/ -run 'TestIsAddrInUse|TestFormatAddrInUse' -count=1
go test ./cmd/trace/ -run 'TestServe' -count=1
# Optional if time: (cd web && npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts)
grep -n 'DEFAULT_MAX\|UI_CAP' web/src/screens/Graph.tsx
grep -n 'getImpact\|FormatAddrInUseMessage\|DefaultAddr' web/src/api/ops.ts internal/httpapi/*.go | head -40
grep -n 'Multi-project\|7433\|auto-port' docs/gui-quickstart.md | head -30
```

Confirm VERIFY-NOTES: overall PASS; **P32-PORT** ticked (#1 + docs; #2 deferred); DESIGN-LOCKS/Laws addressed; residuals listed; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 research (`RESEARCH.md` + P32-PORT note)
- [ ] S01 UX IA (`UX-IA.md`)
- [ ] S02 API gaps **and** P32-PORT ship (`NO-GAPS.md` + #1; #2 deferred)
- [ ] S03 depth implement (graph-home + inspector)
- [ ] S04 visual craft
- [ ] S05 polish + multi-project port docs
- [ ] S06 VERIFY + successor documented (**never TBD**; default **no successor**)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| P32-PORT #2 auto free-port / `:0` | Deferred — OPEN-PORT-MULTI; not required for close |
| Serve “listening on” before bind fail | S02 low — non-blocking |
| Chrome `box-shadow` transition unused | S04/S05 nit — deferred |
| Canvas keyboard select via list | Acceptable |
| No screenshots / media pipeline | Deferred — no pipeline |
| Hosted SaaS / cloud | Separate product — not Phase 32 residual |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# DR-HANDOFF — Phase 32

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-21 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 31 CLOSED (`P31-S02-02`) |
| Theme | Graph-first GUI explorer |
| Outcome | Hybrid C graph-home + inspector depth; S04 craft; P32-PORT #1 + docs (#2 deferred); VERIFY PASS |
| Successor decision | **no successor** |
| Residuals (non-blocking) | PORT #2 defer; chrome box-shadow nit; no screenshots; serve listening-on ordering |
| Close owner | P32-S06-02 |
| Verify | Cite VERIFY-NOTES + evidence dir; P32-PORT ticked |

## Scope checklist

- [x] S00 research (`RESEARCH.md` + P32-PORT note)
- [x] S01 UX IA (`UX-IA.md`)
- [x] S02 API gaps **and** P32-PORT ship
- [x] S03 depth implement (graph-home + inspector)
- [x] S04 visual craft
- [x] S05 polish + multi-project port docs
- [x] S06 VERIFY + successor documented (**no successor**)
```

If thin follow-on: set `Successor decision` to **Phase NN — \<theme\>**; first runnable **PNN-00**; confirm scaffold present before close.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### TODO.md / AGENTS.md updates (on APPROVE)

**If `no successor`:**

1. `docs/TODO.md` orchestrator paste: Phase 00–32 complete; no active phase (or idle awaiting human promotion); do not invent Phase 33.
2. Phase boards table: Phase 32 status `done`; Next `—`.
3. `AGENTS.md` Current focus: Phase 32 complete; explorer GUI shipped; P32-PORT #1+#docs; #2 deferred; cloud remains separate product.

**If thin follow-on:** point Active phase / next runnable at the scaffolded phase planner; never TBD.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-06-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P32-S06-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** no successor  (or Phase NN / PNN-00 | pending repair spawn)

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| web build | |
| P32-PORT Go tests | |
| DESIGN-LOCKS / getImpact / budgets | |
| Port docs / loopback | |
| P32-PORT tick (#1+#docs; #2 defer) | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Next
(idle / PNN-00 / P32-S06-02a)
```

### On FAIL / repair spawn

Insert immediately below this row:

| Order | ID | Role |
|------:|----|------|
| 569a | P32-S06-02a | implement repair |
| 569b | P32-S06-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 32 done.

## Role work

1. Fresh-session re-verify S06-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor **`no successor`** or scaffolded Phase NN (**never TBD**).
4. Update `docs/TODO.md` + `AGENTS.md` per decision table.
5. If thin follow-on: create runnable next-phase scaffold before close (protocol duties above).
6. Do **not** rewrite S00–S05 `done` history or S06-01 VERIFY-NOTES content except to cite them.

## Todo updates

Status + notes on **P32-S06-02**; may update TODO.md / AGENTS.md / DR-HANDOFF; may spawn repair rows below this row if needed; may create **upcoming** next-phase files only when successor is a thin follow-on.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**no successor** or Phase NN / PNN-00 — never TBD)
- [ ] If follow-on: runnable scaffold present (README + phase planner + stubs + board + TODO index)
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 32 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

Idle (**no successor**) — or successor first row / repair spawn
