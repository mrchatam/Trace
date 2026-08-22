# P38-S07-02 — DR-HANDOFF Phase 38 close

## Metadata
- id: P38-S07-02
- todo_ids: [P38-S07-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown, shipping-and-launch]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed (spot-check + scaffold)
- hooks: []

## Objective

Independent **fresh-session** review of S07-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 38 DR-HANDOFF** with explicit successor (**never TBD**). Deliver **runnable Phase 39 scaffold** per [REMEDIATION-PLAN.md](../scope-06-remediation-plan/REMEDIATION-PLAN.md) §3/§6 (entry co-wave **G1 + G3 + G4**). Update `docs/TODO.md` + `AGENTS.md`. Phase 38 complete when this row is `done`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S07-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S07-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S07-01
- [REMEDIATION-PLAN.md](../scope-06-remediation-plan/REMEDIATION-PLAN.md) — §1 M-001, §2 G1–G9, §3 Phase 39 sketch, §6 successor
- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md)
- [SATURATION-NOTES.md](../scope-05-saturation-gate/SATURATION-NOTES.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Phase 37 [DR-HANDOFF.md](../../../phase-37-p36-residuals/DR-HANDOFF.md) — predecessor close pattern
- Pattern: [P24 S05-02 review](../../../phase-24-agent-effectiveness-investigation/scopes/scope-05-phase-verify/02-scope-review.md) (Phase scaffold)
- Pattern: [P37 S03-02](../../../phase-37-p36-residuals/scopes/scope-03-verify/02-dr-handoff.md) (DR-HANDOFF close)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-38.md](../../../../TODO/phase-38.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S07-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-07-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p38-s07-01-verify/evidence/` |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |
| Remediation plan | [REMEDIATION-PLAN.md](../scope-06-remediation-plan/REMEDIATION-PLAN.md) |
| All S00–S06 deliverables | See [01-verify.md](01-verify.md) manifest |
| Phase board | [docs/TODO/phase-38.md](../../../../TODO/phase-38.md) |

## Locked DR-HANDOFF close policy (FINAL — S07-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S07-01** — archive + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S07-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S07-01 `done`; verify blocks 0–6 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 39 — Context orient & harness** (human promotes **P39-00**) |
| Entry co-wave | **G1 + G3 + G4** — query+task orient merge, MCP/harness orient, dual-stack docs |
| Idle alternative | **`no successor`** — if human explicitly defers Phase 39 (document reason in Notes) |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S06 `done` history; ship product in this row; implement Phase 39 themes here |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 38 board rows `done` |
| Portable graph | Confirm `trace/graph.json` current if entity changes during P38 (expect no-op — investigation only) |

### Successor decision table (locked — pick exactly one)

| Outcome (from S07-01 + spot-check) | Decision | Next action |
|-------------------------------------|----------|-------------|
| VERIFY blocks 0–6 green; 7 artifacts archived; saturation + plan APPROVE | **Phase 39 — Context orient & harness (G1+G3+G4)** | Close DR-HANDOFF; scaffold Phase 39; mark Phase 38 **done** |
| Same as above but human defers implementation | **`no successor`** | Close DR-HANDOFF; Phase 38 done; Phase 39 folder optional stub-only |
| Block 0 FAIL (product code in P38) | **Do not close** — spawn repair | Keep OPEN; insert S07-02a/b |
| Block 1 FAIL (artifact missing / H* gap) | **Do not close** — send back S07-01 or spawn S04/S06 repair | Keep OPEN |
| Block 4 FAIL (MP or other peer cites missing) | **Do not close** — send back S03/S07-01 | Keep OPEN |
| VERIFY-NOTES missing or evidence dir absent | **Do not close** — send back S07-01 | Keep OPEN |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00: `INVESTIGATION-INDEX.md` (H1–H11 register)
- [ ] S01: `TRACE-AUDIT.md` (live Trace audit)
- [ ] S02: `PEER-CG.md` (Codegraph peer)
- [ ] S03: `PEER-UA-GF.md` (UA + Graphify + **Mempalace §3**)
- [ ] S04: `GAP-REGISTRY.md` (matrix Trace \| CG \| UA \| GF \| **MP**)
- [ ] S05: `SATURATION-NOTES.md` (APPROVE saturated)
- [ ] S06: `REMEDIATION-PLAN.md` (G1–G9 ranked)
- [ ] S07: `VERIFY-NOTES.md` + successor documented

### P38 outcome summary (for DR-HANDOFF prose)

One paragraph covering:

1. **Saturated investigation** — 11 gaps (G-001…G-011), M-001 moat non-gap, G-004a defer
2. **Peer coverage** — CG, UA, GF, **MP** mechanism cites; H7 compose-first before unified explore
3. **Plan only** — REMEDIATION-PLAN ranks G1→G9; Phase 39 entry **G1+G3+G4**; H11 doc-only
4. **No implement** — zero product Go/TS in P38 window

### Top 3 G* themes for human promotion (locked default)

| Rank | Theme | GAP ids | Phase 39 scope sketch |
|------|-------|---------|----------------------|
| 1 | **G1** Query+task orient merge | G-001, G-002 | Optional `query` on context/compiler; merge search hits |
| 2 | **G3** MCP/harness orient | G-006, G-010 | Playbook, moat-first bootstrap, 9/16 hygiene |
| 3 | **G4** Dual-stack docs (H11) | G-011 | CONTRIBUTING/AGENTS Trace+CG recipe — **doc-only** |

Secondary queue (document in DR-HANDOFF, do not scaffold implement rows yet): G5 GUI orient, G2 compose-first, G2 unified explore Phase 40+.

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| GAP-REGISTRY §6 H7 "Open" | Forward-only — S05 defer→S06 closed |
| G5 GUI start in Phase 39 | May slip to Phase 40 if capacity-bound (REMEDIATION-PLAN §5 Q1) |
| G2 unified `trace_explore` | Phase 40+ after G1 + law spike |
| G-004a vector | Permanent defer — DR-NOSSEM |
| G9 intent pipeline | Phase 41+ or doc-revise |
| Harness 9/16 Cursor fix | G3 — docs vs code TBD in Phase 39 planner |

### DR-HANDOFF.md update template (on APPROVE — default Phase 39)

```markdown
# DR-HANDOFF — Phase 38

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 37 CLOSED |
| Theme | Retrieval & context peer-gap investigation |
| Outcome | Saturated; GAP-REGISTRY 11 gaps; REMEDIATION-PLAN G1–G9; **plan only — no implement** |
| Successor decision | **Phase 39 — Context orient & harness (G1+G3+G4 entry co-wave)** — human promotes P39-00 |
| Top themes | G1 orient merge, G3 harness, G4 dual-stack docs |
| Moat | M-001 preserved — all remediation merges into task loop, never replaces |
| Residuals (non-blocking) | G5 slip; G2 explore Phase 40+; G-004a vector defer; G9 TBD |
| Close owner | P38-S07-02 |
| Verify | VERIFY-NOTES + experiments/runs/…-p38-s07-01-verify/evidence/ |

## Scope checklist

- [x] S00 INVESTIGATION-INDEX.md
- [x] S01 TRACE-AUDIT.md
- [x] S02 PEER-CG.md
- [x] S03 PEER-UA-GF.md (+ MP §3)
- [x] S04 GAP-REGISTRY.md (+ MP column)
- [x] S05 SATURATION-NOTES.md (APPROVE)
- [x] S06 REMEDIATION-PLAN.md
- [x] S07 VERIFY + successor

## Forward artifacts for Phase 39

- REMEDIATION-PLAN.md §2–§3
- GAP-REGISTRY.md
- SATURATION-NOTES.md
- PEER-CG.md, PEER-UA-GF.md, TRACE-AUDIT.md
```

If **`no successor`**: set `Successor decision` accordingly; minimal Phase 39 stub OK but not required.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

## Phase handoff scaffold (mandatory on APPROVE with Phase 39 successor)

Per agent-loop-protocol, before marking this row `done` with Phase 39 successor:

### Required Phase 39 artifacts

| Artifact | Path (locked default) |
|----------|----------------------|
| Phase README | `docs/phases/phase-39-context-orient-harness/README.md` |
| Phase planner | `docs/phases/phase-39-context-orient-harness/00-PHASE-PLANNER.md` (**runnable** — row P39-00) |
| Design SoT stub | `docs/phases/phase-39-context-orient-harness/INTAKE.md` — goal, G1/G3/G4 in/out, M-001 charter, links to P38 artifacts |
| DR-HANDOFF open | `docs/phases/phase-39-context-orient-harness/DR-HANDOFF.md` (**OPEN**) |
| Scope 00 — G1 | `scopes/scope-00-context-orient-merge/` — `00-PLANNER.md`, `01-implement.md`, `02-review.md`, `SCOPE-TODOS.md` (minimal stubs OK) |
| Scope 01 — G3 | `scopes/scope-01-harness-orient/` — same stub set |
| Scope 02 — G4 | `scopes/scope-02-dual-stack-docs/` — same stub set |
| Scope 03 — VERIFY | `scopes/scope-03-verify/` — `00-PLANNER.md`, `01-verify.md`, `02-dr-handoff.md`, `SCOPE-TODOS.md` |
| Board file | `docs/TODO/phase-39.md` — **P39-00** phase planner first **pending** row |
| Index link | `docs/TODO.md` — Phase 39 row in phase boards table |

### Phase 39 README minimum content

```markdown
# Phase 39 — Context orient & harness

Human-promoted successor to Phase 38 close. **Implement** entry co-wave G1+G3+G4.

## Goal
- G1: Query+task orient merge (G-001, G-002)
- G3: MCP/harness orient playbook (G-006, G-010)
- G4: Dual-stack documentation — doc-only (G-011)

## Evidence basis
- Phase 38 REMEDIATION-PLAN.md, GAP-REGISTRY.md, SATURATION-NOTES.md
- Phase 38 DR-HANDOFF.md

## In scope
(per INTAKE — planner thickens at P39-00)

## Out of scope
- G2 unified trace_explore (Phase 40+)
- G-004a vector / embeddings
- Rewriting Phase 38 investigation history
- Product dual-index default

## Moat charter (M-001)
All work merges peer patterns INTO task loop + gates — never replaces moat.
```

### Phase 39 board minimum rows (`docs/TODO/phase-39.md`)

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 671 | P39-00 | pending | `00-PHASE-PLANNER.md` | First runnable after human promotion |
| 672+ | P39-S00-00 … | pending | scope stubs | Planner fills at P39-00 |

Exact order IDs may shift — use next free order after Phase 38 row 670. **P39-00 must be first pending** after Phase 38 complete.

### 00-PHASE-PLANNER.md minimum (runnable alone)

Must include:

- Metadata id P39-00, role planner
- Objective: lock scopes against REMEDIATION-PLAN G1/G3/G4 + live repo
- References: P38 REMEDIATION-PLAN §2 G1/G3/G4, GAP-REGISTRY, DESIGN-LOCKS, project-rules Law 6–7 / 19
- Locked defaults: entry co-wave G1+G3+G4; M-001 moat; G4 doc-only; no G2 explore ship
- Exit criteria: scope stubs runnable; board points to P39-S00-00
- Next: `P39-S00-00`

### TODO.md / AGENTS.md updates (on APPROVE)

**If Phase 39 promotion (default scaffold path):**

1. `docs/TODO.md`: Phase 00–38 complete; Phase 39 scaffold pending human promotion; update orchestrator paste
2. Phase boards: Phase 38 all rows `done`; add Phase 39 row with link to `phase-39.md`
3. `AGENTS.md` Current focus: Phase 38 complete — investigation saturated, remediation plan shipped; Phase 39 implement scaffold ready; **human promotes P39-00**

**If `no successor`:**

1. Orchestrator paste: Phase 00–38 complete; idle awaiting human promotion
2. Phase 38 board all rows `done`
3. `AGENTS.md`: Phase 38 complete; plan-only delivered

**Never** leave next runnable as TBD.

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
P38=docs/phases/phase-38-retrieval-context-peer-gaps

test -f "$P38/scopes/scope-07-verify/VERIFY-NOTES.md"
EVID=$(ls -d experiments/runs/*-p38-s07-01-verify/evidence 2>/dev/null | tail -1)
test -n "$EVID" && test -d "$EVID"
test -f "$EVID/manifest.sha256"

# Block 0 spot-check
test ! -s "$EVID/00-product-commits-since-promotion.txt" || \
  ! grep -v '^$' "$EVID/00-product-commits-since-promotion.txt" | grep -q .

# Block 1 — artifacts
test -f "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md"
grep -c 'G-00[0-9]\|G-01[01]' "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md" | awk '$1>=11'

# Block 2 — saturation
grep -q 'ready_for_REMEDIATION_PLAN.*true' "$P38/scopes/scope-05-saturation-gate/SATURATION-NOTES.md"

# Block 3 — rank + rejects
grep -q 'G1.*G3.*G4' "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md"
grep -c '^| [0-9]' "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md" | awk '$1>=12'

# Block 4 — MP §3
grep -q '§3.*Mempalace\|Mempalace' "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md"
grep -q 'searcher.py\|layers.py' "$P38/scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md"

# Block 5 — M-001
grep -q 'M-001' "$P38/scopes/scope-04-gap-registry/GAP-REGISTRY.md"

# Block 6 — Phase 39 in plan
grep -q 'Phase 39' "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md"
grep -q 'G1.*G3.*G4\|G1 + G3 + G4' "$P38/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md"
```

Confirm VERIFY-NOTES: overall PASS; blocks 0–6 ticked; DR-HANDOFF still OPEN before this row closes it.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-07-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P38-S07-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** Phase 39 / P39-00 | no successor | pending repair spawn

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir + manifest | |
| Block 0 product boundary | |
| Block 1 artifacts + H* | |
| Block 2 saturation | |
| Block 3 REMEDIATION-PLAN | |
| Block 4 MP peer cites | |
| Block 5 M-001 | |
| Phase 39 scaffold | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Phase 39 scaffold checklist
- [ ] README.md
- [ ] 00-PHASE-PLANNER.md runnable
- [ ] INTAKE.md
- [ ] DR-HANDOFF OPEN
- [ ] Scope stubs S00–S03
- [ ] docs/TODO/phase-39.md
- [ ] docs/TODO.md index link
- [ ] AGENTS.md updated

## Next
(P39-00 / idle / P38-S07-02a)
```

### On FAIL / repair spawn

Insert immediately below this row on phase board:

| Order | ID | Role |
|------:|----|------|
| 670a | P38-S07-02a | implement repair |
| 670b | P38-S07-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 38 done.

## Role work

1. Fresh-session re-verify S07-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick scope checklist; set successor **Phase 39 / no successor / repair** (**never TBD**).
4. **Scaffold Phase 39** per table above (if Phase 39 successor).
5. Update `docs/TODO.md` + `AGENTS.md` per decision table.
6. Mark all Phase 38 board rows `done`.
7. Run `trace seed export -o trace/graph.json` if warranted (expect no-op).
8. Do **not** rewrite S00–S06 `done` history or S07-01 VERIFY-NOTES except to cite them.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**Phase 39 G1+G3+G4 / no successor / repair** — never TBD)
- [ ] Phase 39 scaffold complete (if default successor): README, 00-PHASE-PLANNER, INTAKE, DR-HANDOFF OPEN, scope stubs, `phase-39.md`, TODO.md link
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 38 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

**P39-00** (if human promotes) — idle (**no successor**) — or repair spawn
