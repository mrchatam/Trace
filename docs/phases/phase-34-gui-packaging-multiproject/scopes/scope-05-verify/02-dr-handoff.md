# P34-S05-02 — DR-HANDOFF Phase 34 close

## Metadata
- id: P34-S05-02
- todo_ids: [P34-S05-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S05-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 34 DR-HANDOFF** with explicit successor (**never TBD**). Default successor **`no successor`** unless VERIFY residuals require a thin follow-on phase scaffold (per agent-loop-protocol **Phase handoff**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 34 complete when this row is `done`. **No product code.** Do **not** implement a successor in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (PLAN T9)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S05-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-34.md](../../../../TODO/phase-34.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S05-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-05-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p34-s05-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Design locks | `DESIGN-LOCKS.md` |
| Phase board | `docs/TODO/phase-34.md` |
| Prior scope artifacts | S00 RESEARCH; S01 PLAN; S02–S04 board Notes / reviews |

## Locked DR-HANDOFF close policy (FINAL — S05-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S05-01** — VERIFY floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S05-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S05-01 `done`; verify blocks 0–5 green per VERIFY-NOTES + independent spot-check; L1–L4 + Docs ticked; T10 stub-fail addressed; block 6 live **or** waived |
| Default successor | **`no successor`** — L1–L4 shipped (embed + auto-port + `.trace/` only); known nits (contributor `web/` DX, StaticDir path string, optional CI) are deferred/optional, **not** a new phase by default |
| Thin follow-on exception | Only if VERIFY/spot-check leaves **blocking** product residuals that need a new phase theme **or** human promotes a follow-on — then scaffold per protocol (below) |
| Cloud / hosted SaaS | **Not** a Phase 34 successor — separate product/repo |
| Regression path | Spawn `P34-S05-02a` implement + `02b` review; **do not** close Phase 34 |
| Must not | Leave `Successor decision: TBD`; treat contributor `web/` docs as requiring a phase; rewrite S00–S04 `done` history; ship product in this row; start implementing a successor |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 34 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S05-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; residuals only as listed (contributor DX OK, live smoke waive OK) | **`no successor`** | Close DR-HANDOFF; mark Phase 34 **done** in TODO/AGENTS; idle orchestrator paste |
| Embed stub / auto-port / docs-primary / Law 19 / public default FAIL | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b; successor = `pending repair spawn` |
| VERIFY-NOTES missing blocks, L1–L4 unticked, or evidence dir absent | **Do not close** — spawn repair or send back S05-01 | Keep OPEN |
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
2. **Board:** `docs/TODO/phase-NN.md` with phase planner as **first pending** row after Phase 34’s last `done` row
3. **Index:** link in `docs/TODO.md` phase boards table; orchestrator paste → Active phase NN / next runnable `PNN-00`
4. **AGENTS.md** Current focus → Phase NN
5. Deep tasking belongs to the **next** phase planner — this row delivers a **runnable handoff**, not finished implement prompts

If default **`no successor`**: do **not** invent Phase 35 stubs. Update TODO/AGENTS to Phase 34 **complete** / idle (Phases 00–34 closed; cloud remains separate product).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-34-gui-packaging-multiproject/scopes/scope-05-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p34-s05-01-verify/evidence || ls experiments/runs/ | grep p34-s05-01
go test ./internal/httpapi/ -run 'TestStaticCSPAndEmbedFallback|TestListenAutoPort_' -count=1 -p 1
go test ./cmd/trace/ -run 'TestGuiServeConcurrentDefaultDistinctPorts|TestGuiExplicitDefaultAddrBusyNoHop' -count=1 -p 1
# Stub-fail:
grep -q 'Embedded GUI stub' internal/httpapi/embeddist/index.html && echo FAIL || echo PASS-no-stub
grep -n 'id="root"\|/assets/' internal/httpapi/embeddist/index.html | head -20
go run ./cmd/trace gui --help 2>&1 | head -40
grep -nE '7432|7441|embed|\.trace/|no auto free-port|two-artifact' docs/gui-quickstart.md | head -40
```

Confirm VERIFY-NOTES: overall PASS; L1–L4 + Docs ticked; T10 stub-fail addressed; block 6 live **or** waived; DESIGN-LOCKS addressed; residuals listed; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 research (`RESEARCH.md`)
- [ ] S01 plan (`PLAN.md`)
- [ ] S02 embed / static defaults
- [ ] S03 auto free-port
- [ ] S04 docs + tests
- [ ] S05 VERIFY + successor documented (**never TBD**; default **no successor**)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| Contributor Trace-checkout `web/` DX | Documented + labeled — not consumer requirement |
| Default StaticDir path string `<root>/web/dist` | Resolution disk→embed→placeholder — accept |
| Optional CI embed-gui workflow | Out of phase / deferred |
| Live consumer-temp smoke | Captured in S05-01 **or** waived — not required for close either way if Blocks 1–5 green |
| Explore craft / hosted SaaS | Separate / prior phase — not Phase 34 residual |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# DR-HANDOFF — Phase 34

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-21 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 33 CLOSED (`P33-S06-02`) |
| Theme | Embed real SPA; auto port; consumer `.trace/` only |
| Outcome | L1–L4: real SPA from binary; auto free-port 7432–7441; consumer `.trace/` only; docs T8; VERIFY PASS |
| Successor decision | **no successor** |
| Residuals (non-blocking) | Contributor web DX; StaticDir path string; optional CI; live smoke captured\|waived; … |
| Close owner | P34-S05-02 |
| Verify | Cite VERIFY-NOTES + evidence dir; L1–L4 + Docs ticked |

## Scope checklist

- [x] S00 research (`RESEARCH.md`)
- [x] S01 plan (`PLAN.md`)
- [x] S02 embed / static defaults
- [x] S03 auto free-port
- [x] S04 docs + tests
- [x] S05 VERIFY + successor documented (**no successor**)
```

If thin follow-on: set `Successor decision` to **Phase NN — \<theme\>**; first runnable **PNN-00**; confirm scaffold present before close.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### TODO.md / AGENTS.md updates (on APPROVE)

**If `no successor`:**

1. `docs/TODO.md` orchestrator paste: Phase 00–34 complete; no active phase (or idle awaiting human promotion); do not invent Phase 35.
2. Phase boards table: Phase 34 status `done`; Next `—`.
3. `AGENTS.md` Current focus: Phase 34 complete; L1–L4 shipped (embed SPA + auto-port + `.trace/` only); cloud remains separate product.

**If thin follow-on:** point Active phase / next runnable at the scaffolded phase planner; never TBD.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-05-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P34-S05-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** no successor  (or Phase NN / PNN-00 | pending repair spawn)

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| Static embed + stub-fail | |
| Auto-port + concurrent | |
| Docs T8 / help | |
| L1–L4 ticks | |
| Live consumer-temp | LIVE \| WAIVED |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Next
(idle / PNN-00 / P34-S05-02a)
```

### On FAIL / repair spawn

Insert immediately below this row:

| Order | ID | Role |
|------:|----|------|
| 609a | P34-S05-02a | implement repair |
| 609b | P34-S05-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 34 done.

## Role work

1. Fresh-session re-verify S05-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor **`no successor`** or scaffolded Phase NN (**never TBD**).
4. Update `docs/TODO.md` + `AGENTS.md` per decision table.
5. If thin follow-on: create runnable next-phase scaffold before close (protocol duties above).
6. Do **not** rewrite S00–S04 `done` history or S05-01 VERIFY-NOTES content except to cite them.

## Todo updates

Status + notes on **P34-S05-02**; may update TODO.md / AGENTS.md / DR-HANDOFF; may spawn repair rows below this row if needed; may create **upcoming** next-phase files only when successor is a thin follow-on.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**no successor** or Phase NN / PNN-00 — never TBD)
- [ ] If follow-on: runnable scaffold present (README + phase planner + stubs + board + TODO index)
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 34 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

Idle (**no successor**) — or successor first row / repair spawn
