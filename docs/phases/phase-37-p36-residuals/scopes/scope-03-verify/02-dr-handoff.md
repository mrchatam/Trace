# P37-S03-02 — DR-HANDOFF Phase 37 close

## Metadata
- id: P37-S03-02
- todo_ids: [P37-S03-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown, shipping-and-launch]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 37 DR-HANDOFF** with explicit successor (**never TBD**). Update Phase 36 DR-HANDOFF residuals paragraph if P37 consumed deferrals. Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 37 complete when this row is `done`. **No product code.** Do **not** implement Phase 38 in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S03-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [PLAN.md](../scope-01-plan/PLAN.md) — §3 re-defer registry
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Phase 36 [DR-HANDOFF.md](../../../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md)
- Phase 38 scaffold: [`docs/TODO/phase-38.md`](../../../../TODO/phase-38.md), [`phase-38-retrieval-context-peer-gaps/`](../../../phase-38-retrieval-context-peer-gaps/)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-37.md](../../../../TODO/phase-37.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p37-s03-01-verify/evidence/` |
| Pinned summaries (optional) | `docs/verification/phase-37-p36-residuals/` |
| Phase handoff | `DR-HANDOFF.md` |
| Design locks | `DESIGN-LOCKS.md` |
| Phase board | `docs/TODO/phase-37.md` |
| Prior scope artifacts | S00 RESIDUALS, S01 PLAN, S02 implement + review |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; verify blocks 0–5 green per VERIFY-NOTES + independent spot-check; all S02 accepts evidenced |
| Default successor | **Phase 38 scaffold ready — human promotes `P38-00`** (investigation only; no implement in P38 by default) |
| Idle alternative | **`no successor`** — if human chooses not to promote Phase 38 immediately |
| Cloud / hosted SaaS | **Not** a Phase 37 successor — separate product/repo |
| Regression path | Spawn `P37-S03-02a` implement + `02b` review; **do not** close Phase 37 |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S02 `done` history; ship product in this row; implement Phase 38 |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 37 board rows `done` |
| Portable graph | If S02 entity export warranted, confirm `trace/graph.json` per CONTRIBUTING |

### Successor decision table (locked — pick exactly one)

Aligned to VERIFY outcome (blocks 0–5 from [01-verify.md](01-verify.md)):

| Outcome (from S03-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; R1–R6,R8,R11 shipped; R10 evidenced or explicitly blocked with file; re-defer R7/R9/R8-full documented | **Phase 38 — retrieval/context peer-gap investigation** (scaffold exists) | Close DR-HANDOFF; mark Phase 37 **done**; set TODO/AGENTS to idle **or** point next runnable at **P38-00** if human promotes |
| Same as above but human declines Phase 38 promotion | **`no successor`** | Close DR-HANDOFF; Phase 37 done; orchestrator idle |
| Block 0 P36 regression subset FAIL | **Do not close** — spawn repair | Keep OPEN; insert S03-02a/b; successor = `pending repair spawn` |
| Block 1 FAIL (S02 accept missing evidence) | **Do not close** — spawn repair or send back S03-01 | Keep OPEN |
| Block 3 FAIL (greenfield MCP or gate regression) | **Do not close** — spawn repair | Keep OPEN — critical agent path |
| R1 guard FAIL (`plan_exists` flipped without rows) | **Do not close** — spawn repair | Keep OPEN — DESIGN-LOCKS violation |
| Active-work `plan_missing` weakened | **Do not close** — spawn repair | Keep OPEN — P36 invariant regression |
| VERIFY-NOTES missing blocks or evidence dir absent | **Do not close** — send back S03-01 | Keep OPEN |
| R10 accepted but no browser/screenshot evidence when GUI was reachable | **Do not close** or mark R10 PARTIAL with explicit human gate | Keep OPEN until pinned file or documented blocker |
| Human wants hosted SaaS next | **Not** Trace core successor | Close with Phase 38 or no successor; cloud separate |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### Block 5 — Successor table (must appear in VERIFY-NOTES before close)

S03-01 prepares; S03-02 **confirms** one row:

| VERIFY outcome | Successor | First runnable |
|----------------|-----------|----------------|
| All blocks green | Phase 38 investigation (default) | `P38-00` after human promotion |
| All blocks green, human idle | `no successor` | — |
| Repair needed | `pending repair spawn` | `P37-S03-02a` |

Phase 38 scaffold checklist (already present — verify, do not recreate unless missing):

- [ ] `docs/phases/phase-38-retrieval-context-peer-gaps/README.md`
- [ ] `00-PHASE-PLANNER.md`
- [ ] Scope stubs S00–S07
- [ ] `docs/TODO/phase-38.md` + index link in `docs/TODO.md`

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| R7 enforce default `warn` | Re-defer — R6 test + stderr nudge sufficient |
| R9 feet-seller refinement quality | Re-defer — Block 2 doc path; human dogfood |
| R8-full plan tree GUI | Re-defer — Phase 38 investigation may address |
| R10 blocked (GUI unreachable) | Document blocker; optional human follow-up |
| Feet-seller `plan_uncritiqued` post-bootstrap | Accept — deliberation separate from P37 |
| PlanExists bridge | Permanent reject — advisory-only R1 shipped |
| MCP critique-seed tool | Permanent reject — R11 doc path |
| Hosted SaaS / cloud | Separate product |

### P36 DR-HANDOFF cross-link (on close)

Update Phase 36 DR-HANDOFF residuals paragraph or Phase 37 DR-HANDOFF to note **consumed**:

- HTTP POST plan routes (R2) — shipped
- MCP `trace_loop action=gate` (R3) — shipped
- Loop status `advisories[]` (R5) — shipped
- Bootstrap help refinement (R4), enforce nudge test (R6), Overview surface (R8), critique doc (R11) — shipped
- Remaining P36 deferrals: PlanExists bridge, enforce default (R7), full plan GUI (R8-full)

### DR-HANDOFF.md update template (on APPROVE — default Phase 38 scaffold)

```markdown
# DR-HANDOFF — Phase 37

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 36 CLOSED |
| Theme | Phase 36 residuals closure (R1–R11 triage) |
| Outcome | advisories[] channel; MCP gate; HTTP bootstrap POST; help/test/doc/GUI adapters; R7/R9/R8-full re-deferred |
| Successor decision | **Phase 38 — retrieval/context peer-gap investigation** (human promotes P38-00) |
| Residuals (non-blocking) | R7, R9, R8-full; R10 if partial |
| Close owner | P37-S03-02 |
| Verify | Cite VERIFY-NOTES + experiments/runs/…-p37-s03-01-verify/evidence/ |

## Scope checklist

- [x] S00 triage (`RESIDUALS.md`)
- [x] S01 plan (`PLAN.md`)
- [x] S02 implement + tests + review
- [x] S03 VERIFY + successor documented
```

If **`no successor`**: set `Successor decision` accordingly; Phase 38 remains scaffold-only until human promotes.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-37-p36-residuals/scopes/scope-03-verify/VERIFY-NOTES.md
EVID=$(ls -d experiments/runs/*-p37-s03-01-verify/evidence 2>/dev/null | tail -1)
test -n "$EVID" && test -d "$EVID"

go build -o /tmp/trace ./cmd/trace

# Block 0 spot-check (P36 regression subset)
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... ./internal/config/... -count=1 \
  -run 'Greenfield_MCPPlanBootstrap|FeetSellerExport_GateHonesty|ActiveWork_PlanMissing|TerminalPlanGapAdvisory|PlanBootstrap_Idempotent|GoalStructureWarning_OverThreshold|RegisteredToolNames'

# Block 1 spot-check (S02 accepts)
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./internal/httpapi/... ./cmd/trace/... ./internal/config/... -count=1 \
  -run 'LoopStatus_.*Advisory|BootstrapAdvisoryNeverSetsPlanExists|MCPLoopGate|HTTPPlanBootstrap|PlanHelp_MentionsRefinement|WarnIfTraceDirWithoutConfig'

FEET="/home/ali/Desktop/feet seller telegram app"
STEP1=33247e2d-aa10-4b25-b194-4b7afb5a6359

# Status advisories live sample
/tmp/trace -C "$FEET" loop status --task "$STEP1" | jq '{advisories, plan_exists: .deliberation.policy_inputs.plan_exists}'

# R11 doc cite
grep -c 'trace loop apply' docs/rules/agent-loop-protocol.md

# R10 pinned evidence (if claimed PASS)
test -d docs/verification/phase-37-p36-residuals || test -d "$EVID/04-browser"
```

Confirm VERIFY-NOTES: overall PASS; blocks 0–5 ticked; re-defer registry present; DR-HANDOFF still OPEN before this row closes it.

### TODO.md / AGENTS.md updates (on APPROVE)

**If Phase 38 promotion (default scaffold path):**

1. `docs/TODO.md`: Phase 00–37 complete; Phase 38 scaffold pending human promotion; or set Active phase **38** if promoting now.
2. Phase boards: Phase 37 status `done`; Next `P38-00` or `—`.
3. `AGENTS.md` Current focus: Phase 37 complete — P36 residuals closed; Phase 38 investigation scaffold ready; cloud separate product.

**If `no successor`:**

1. Orchestrator paste: Phase 00–37 complete; idle awaiting human promotion.
2. Phase 37 board all rows `done`.
3. `AGENTS.md`: Phase 37 complete; residuals closure shipped.

**Never** leave next runnable as TBD.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-03-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P37-S03-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** Phase 38 / P38-00 | no successor | pending repair spawn

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| P36 regression subset | |
| S02 acceptance tests | |
| R11 doc cite | |
| R10 browser evidence | |
| Re-defer registry | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## P36 residuals consumed
…

## Next
(P38-00 / idle / P37-S03-02a)
```

### On FAIL / repair spawn

Insert immediately below this row on phase board:

| Order | ID | Role |
|------:|----|------|
| 645a | P37-S03-02a | implement repair |
| 645b | P37-S03-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 37 done.

## Role work

1. Fresh-session re-verify S03-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor **Phase 38 / no successor / repair** (**never TBD**).
4. Cross-link Phase 36 DR-HANDOFF if residuals consumed.
5. Update `docs/TODO.md` + `AGENTS.md` per decision table.
6. Confirm Phase 38 scaffold runnable (do not deep-plan P38 — phase planner owns that).
7. Run `trace seed export -o trace/graph.json` if S02 entity changes not yet exported (CONTRIBUTING).
8. Do **not** rewrite S00–S02 `done` history or S03-01 VERIFY-NOTES content except to cite them.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**Phase 38 / no successor / repair** — never TBD)
- [ ] Phase 36 DR-HANDOFF cross-link updated if applicable
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 37 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

**P38-00** (if promoted) — idle (**no successor**) — or repair spawn
