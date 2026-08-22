# P39-S03-02 — DR-HANDOFF Phase 39 close

## Metadata
- id: P39-S03-02
- todo_ids: [P39-S03-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown, shipping-and-launch]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed (spot-check + scaffold)
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 39 DR-HANDOFF** with explicit successor (**never TBD**). Deliver **runnable Phase 40+ scaffold** per [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §3 Phase 40+ (entry themes **G5 + G2**). Update `docs/TODO.md` + `AGENTS.md`. Phase 39 complete when this row is `done`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S03-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) — §2 G5/G2, §3 Phase 40+
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Phase 38 [DR-HANDOFF.md](../../../phase-38-retrieval-context-peer-gaps/DR-HANDOFF.md) — predecessor close pattern
- Pattern: [P38 S07-02](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-07-verify/02-dr-handoff.md) (Phase scaffold)
- Pattern: [P37 S03-02](../../../phase-37-p36-residuals/scopes/scope-03-verify/02-dr-handoff.md) (DR-HANDOFF close)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-39.md](../../../../TODO/phase-39.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p39-s03-01-verify/evidence/` |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |
| G1 deliverable | S00 implement + S00-02 APPROVE |
| G3 deliverable | S01 implement + S01-02 APPROVE |
| G4 deliverable | S02 implement + S02-02 APPROVE |
| Phase board | [docs/TODO/phase-39.md](../../../../TODO/phase-39.md) |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — VERIFY-NOTES + evidence dir; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; verify blocks 0–6 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 40+ — Read surface & retrieval depth** (human promotes **P40-00**) |
| Entry themes | **G5** GUI graph orient start + **G2** unified `trace_explore` (after G1 law spike) |
| Secondary queue | G6, G7 per REMEDIATION-PLAN rank; G8/G9 Phase 41+ |
| Idle alternative | **`no successor`** — if human explicitly defers Phase 40 (document reason in Notes) |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S02 `done` history; ship product in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 39 board rows `done` |
| Portable graph | Confirm `trace/graph.json` current if entity changes during P39 |

### Successor decision table (locked — pick exactly one)

| Outcome (from S03-01 + spot-check) | Decision | Next action |
|-------------------------------------|----------|-------------|
| VERIFY blocks 0–6 green; G1+G3+G4 delivered | **Phase 40+ — G5 + G2 entry** | Close DR-HANDOFF; scaffold Phase 40+; mark Phase 39 **done** |
| Same as above but human defers implementation | **`no successor`** | Close DR-HANDOFF; Phase 39 done; Phase 40 folder optional stub-only |
| Block 0 FAIL (G1 tests red) | **Do not close** — spawn repair | Keep OPEN; insert S03-02a/b |
| Block 1 FAIL (G3 / tool count) | **Do not close** — send back S01 or S03-01 | Keep OPEN |
| Block 2 FAIL (G4 product code or checklist) | **Do not close** — send back S02 or S03-01 | Keep OPEN |
| Block 3 FAIL (M-001) | **Do not close** — spawn repair | Keep OPEN |
| VERIFY-NOTES missing or evidence dir absent | **Do not close** — send back S03-01 | Keep OPEN |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] **S00** G1 context orient merge — optional query; merge FTS; task moat preserved
- [ ] **S01** G3 harness orient — MCP Instructions playbook; moat-first; 9/16 hygiene
- [ ] **S02** G4 dual-stack docs — CONTRIBUTING/AGENTS Trace+CG recipe (**doc-only**)
- [ ] **S03** VERIFY + successor documented (VERIFY-NOTES + REVIEW-NOTES)

### P39 outcome summary (for DR-HANDOFF prose)

One paragraph covering:

1. **Entry co-wave delivered** — G1 query+task merge, G3 MCP/harness orient, G4 dual-stack docs
2. **M-001 preserved** — task loop + gates primary; query additive; CG optional complement
3. **Compose-first shipped** — G2 ranked read recipe in Instructions (unified explore deferred)
4. **Forward queue** — G5 GUI orient + G2 `trace_explore` → Phase 40+

### Top themes for human promotion (locked default)

| Rank | Theme | GAP ids | Phase 40+ scope sketch |
|------|-------|---------|------------------------|
| 1 | **G5** Graph-first onboarding UX | G-008 | GUI `/` graph route orient; Law 19 adapter only |
| 2 | **G2** Unified `trace_explore` | G-007 | Task-aware capped read after G1 + law spike |

Secondary queue (document in DR-HANDOFF, do not scaffold implement rows yet): G6 non-semantic retrieval, G7 index freshness, G8 layers, G9 intent.

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| G5 full GUI orient | Phase 40+ entry — may start with static sketch |
| G2 unified explore | Phase 40+ after law spike; compose-first already in S01 |
| G-004a vector | Permanent defer — DR-NOSSEM |
| G9 intent pipeline | Phase 41+ or doc-revise |
| `instructions.go:25` Phase 39 S02 stub | Optional P40 doc hygiene |
| `.codegraph/` gitignore | Typically added by `codegraph init` — not Trace product |

### DR-HANDOFF.md update template (on APPROVE — default Phase 40+)

```markdown
# DR-HANDOFF — Phase 39

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 38 CLOSED |
| Theme | Context orient & harness — G1+G3+G4 implement |
| Outcome | G1 query merge, G3 harness orient, G4 dual-stack docs; M-001 preserved |
| Successor decision | **Phase 40+ — Read surface & retrieval depth (G5 + G2 entry)** — human promotes P40-00 |
| Top themes | G5 GUI orient, G2 unified trace_explore |
| Moat | M-001 preserved — peer patterns merge into task loop |
| Residuals (non-blocking) | G6/G7 secondary; G8/G9 Phase 41+; G-004a vector defer |
| Close owner | P39-S03-02 |
| Verify | VERIFY-NOTES + experiments/runs/…-p39-s03-01-verify/evidence/ |

## Scope checklist

- [x] S00 G1 context orient merge
- [x] S01 G3 harness orient
- [x] S02 G4 dual-stack docs
- [x] S03 VERIFY + successor

## Forward artifacts for Phase 40+

- REMEDIATION-PLAN.md §3 Phase 40+
- Phase 39 DR-HANDOFF secondary queue
- G1 shipped (`ContextOptions.Query`, MCP `query`)
- G3 shipped (`ServerInstructions()` orient playbook)
```

If **`no successor`**: set `Successor decision` accordingly; minimal Phase 40 stub OK but not required.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

## Phase handoff scaffold (mandatory on APPROVE with Phase 40+ successor)

Per agent-loop-protocol, before marking this row `done` with Phase 40+ successor:

### Required Phase 40+ artifacts

| Artifact | Path (locked default) |
|----------|----------------------|
| Phase README | `docs/phases/phase-40-read-surface-retrieval-depth/README.md` |
| Phase planner | `docs/phases/phase-40-read-surface-retrieval-depth/00-PHASE-PLANNER.md` (**runnable** — row P40-00) |
| Design SoT stub | `docs/phases/phase-40-read-surface-retrieval-depth/INTAKE.md` — goal, G5/G2 in/out, M-001 charter, links to P39 + REMEDIATION-PLAN |
| DR-HANDOFF open | `docs/phases/phase-40-read-surface-retrieval-depth/DR-HANDOFF.md` (**OPEN**) |
| Scope 00 — G5 | `scopes/scope-00-gui-graph-orient/` — `00-PLANNER.md`, `01-implement.md`, `02-review.md`, `SCOPE-TODOS.md` (minimal stubs OK) |
| Scope 01 — G2 | `scopes/scope-01-unified-explore/` — same stub set |
| Scope 02 — VERIFY | `scopes/scope-02-verify/` — `00-PLANNER.md`, `01-verify.md`, `02-dr-handoff.md`, `SCOPE-TODOS.md` |
| Board file | `docs/TODO/phase-40.md` — **P40-00** phase planner first **pending** row |
| Index link | `docs/TODO.md` — Phase 40 row in phase boards table |

**Note:** Exact scope count may expand at P40-00 (e.g. G6/G7 as secondary scopes). Minimum scaffold above is mandatory; planner thickens at human promotion.

### Phase 40+ README minimum content

```markdown
# Phase 40+ — Read surface & retrieval depth

Human-promoted successor to Phase 39 close. **Implement** entry themes G5 + G2.

## Goal
- G5: Graph-first onboarding UX — GUI orient adapter (G-008)
- G2: Unified trace_explore — task-aware capped read (G-007)

## Evidence basis
- Phase 38 REMEDIATION-PLAN.md §3 Phase 40+
- Phase 39 DR-HANDOFF.md + VERIFY-NOTES
- G1 query merge shipped (Phase 39 S00)
- G3 compose-first recipe shipped (Phase 39 S01)

## In scope
(per INTAKE — planner thickens at P40-00)

## Out of scope
- G-004a vector / embeddings
- Always-on daemon / public bind defaults
- Product dual-index default
- Rewriting Phase 39 delivery history

## Moat charter (M-001)
Unified explore merges into task loop — never query-only replacement.
GUI orient is Law 19 adapter over canonical library/API.
```

### Phase 40+ board minimum rows (`docs/TODO/phase-40.md`)

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 684 | P40-00 | pending | `00-PHASE-PLANNER.md` | First runnable after human promotion |
| 685+ | P40-S00-00 … | pending | scope stubs | Planner fills at P40-00 |

Exact order IDs may shift — use next free order after Phase 39 row 683. **P40-00 must be first pending** after Phase 39 complete.

### 00-PHASE-PLANNER.md minimum (runnable alone)

Must include:

- Metadata id P40-00, role planner
- Objective: lock scopes against REMEDIATION-PLAN G5/G2 + live repo + G1 dependency
- References: P39 VERIFY-NOTES, REMEDIATION-PLAN §2 G5/G2, GAP-REGISTRY G-007/G-008, project-rules Law 6–7 / 19
- Locked defaults: G5 GUI adapter-only; G2 after law spike; M-001 moat; no G-004a vector
- Exit criteria: scope stubs runnable; board points to P40-S00-00
- Next: `P40-S00-00`

### TODO.md / AGENTS.md updates (on APPROVE)

**If Phase 40+ promotion (default scaffold path):**

1. `docs/TODO.md`: Phase 00–39 complete; Phase 40 scaffold pending human promotion; update orchestrator paste
2. Phase boards: Phase 39 all rows `done`; add Phase 40 row with link to `phase-40.md`
3. `AGENTS.md` Current focus: Phase 39 complete — G1+G3+G4 delivered; Phase 40+ scaffold ready; **human promotes P40-00**

**If `no successor`:**

1. Orchestrator paste: Phase 00–39 complete; idle awaiting human promotion
2. Phase 39 board all rows `done`
3. `AGENTS.md`: Phase 39 complete; forward queue documented

**Never** leave next runnable as TBD.

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
P39=docs/phases/phase-39-context-orient-harness

test -f "$P39/scopes/scope-03-verify/VERIFY-NOTES.md"
EVID=$(ls -d experiments/runs/*-p39-s03-01-verify/evidence 2>/dev/null | tail -1)
test -n "$EVID" && test -d "$EVID"

# Block 0 — G1 tests
go test ./internal/compiler/... -count=1 \
  -run 'TestG1QueryHitMerged|TestG1QuerySearchFailOpen' 2>&1 | tail -3

# Block 1 — G3 + 16 tools
go test ./internal/mcp/... -count=1 \
  -run 'TestServerInstructionsNonEmpty|TestToolNamesRegistered' 2>&1 | tail -3

# Block 2 — G4 docs exist
grep -q 'Trace + Codegraph' CONTRIBUTING.md
grep -q 'Optional Codegraph' AGENTS.md

# Block 3 — no trace_explore product tool
! grep -rq 'AddTool.*trace_explore' internal/mcp/

# Block 5 — successor in DR-HANDOFF
grep -q 'Phase 40' "$P39/DR-HANDOFF.md"
grep -q 'G5\|G2' "$P39/DR-HANDOFF.md"
```

Confirm VERIFY-NOTES: overall PASS; blocks 0–6 ticked; DR-HANDOFF still OPEN before this row closes it.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-03-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P39-S03-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** Phase 40+ / P40-00 | no successor | pending repair spawn

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| Block 0 G1 | |
| Block 1 G3 | |
| Block 2 G4 doc-only | |
| Block 3 M-001 | |
| Block 4 Laws | |
| Block 5 successor named | |
| Phase 40+ scaffold | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Phase 40+ scaffold checklist
- [ ] README.md
- [ ] 00-PHASE-PLANNER.md runnable
- [ ] INTAKE.md
- [ ] DR-HANDOFF OPEN
- [ ] Scope stubs S00–S02 (G5, G2, VERIFY)
- [ ] docs/TODO/phase-40.md
- [ ] docs/TODO.md index link
- [ ] AGENTS.md updated

## Next
(P40-00 / idle / P39-S03-02a)
```

### On FAIL / repair spawn

Insert immediately below this row on phase board:

| Order | ID | Role |
|------:|----|------|
| 683a | P39-S03-02a | implement repair |
| 683b | P39-S03-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 39 done.

## Role work

1. Fresh-session re-verify S03-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick scope checklist; set successor **Phase 40+ G5+G2 / no successor / repair** (**never TBD**).
4. **Scaffold Phase 40+** per table above (if Phase 40+ successor).
5. Update `docs/TODO.md` + `AGENTS.md` per decision table.
6. Mark all Phase 39 board rows `done`.
7. Run `trace seed export -o trace/graph.json` if warranted.
8. Do **not** rewrite S00–S02 `done` history or S03-01 VERIFY-NOTES except to cite them.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**Phase 40+ G5+G2 / no successor / repair** — never TBD)
- [ ] Phase 40+ scaffold complete (if default successor): README, 00-PHASE-PLANNER, INTAKE, DR-HANDOFF OPEN, scope stubs, `phase-40.md`, TODO.md link
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 39 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

**P40-00** (if human promotes) — idle (**no successor**) — or repair spawn
