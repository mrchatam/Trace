# P36-S03-02 — DR-HANDOFF Phase 36 close

## Metadata
- id: P36-S03-02
- todo_ids: [P36-S03-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 36 DR-HANDOFF** with explicit successor (**never TBD**). Default successor **`no successor`** unless VERIFY residuals require a thin follow-on phase scaffold. Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 36 complete when this row is `done`. **No product code.** Do **not** implement a successor in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S03-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-36.md](../../../../TODO/phase-36.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p36-s03-01-verify/evidence/` |
| Pinned summaries (optional) | `docs/verification/phase-36-gate-honesty/` |
| Phase handoff | `DR-HANDOFF.md` |
| Design locks | `DESIGN-LOCKS.md` |
| Phase board | `docs/TODO/phase-36.md` |
| Prior scope artifacts | S00 INVESTIGATION, S01 PLAN, S02 implement + review |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; verify blocks 0–7 green per VERIFY-NOTES + independent spot-check; DESIGN-LOCKS must-fix ticked |
| Default successor | **`no successor`** — MCP plan + bootstrap + install + terminal honesty shipped; known deferred items (bridge, HTTP POST plan) are **not** a new phase by default |
| Thin follow-on exception | Only if VERIFY/spot-check leaves **blocking** product residuals needing a new phase theme **or** human promotes follow-on |
| Cloud / hosted SaaS | **Not** a Phase 36 successor — separate product/repo |
| Regression path | Spawn `P36-S03-02a` implement + `02b` review; **do not** close Phase 36 |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S02 `done` history; ship product in this row; start implementing successor |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 36 board rows `done` |
| Portable graph | If VERIFY touched entity semantics post-bootstrap, confirm `trace/graph.json` export per CONTRIBUTING (feet-seller bootstrap is planner state — may warrant Trace repo export on PR boundary) |

### Successor decision table (locked — pick exactly one)

Aligned to VERIFY outcome (blocks 0–7 from [01-verify.md](01-verify.md)):

| Outcome (from S03-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; MCP/bootstrap/install/terminal honesty proven; active PLAN preserved; feet-seller bootstrap OK; residuals only as listed | **`no successor`** | Close DR-HANDOFF; mark Phase 36 **done** in TODO/AGENTS; idle orchestrator paste |
| Block 0 tests FAIL or acceptance tests regressed | **Do not close** — spawn repair | Keep OPEN; insert S03-02a/b; successor = `pending repair spawn` |
| Blocks 2–4 FAIL (terminal honesty or GUI red on DONE) | **Do not close** — spawn repair | Keep OPEN; gate/GUI remediation phase or hotfix row |
| Block 5 FAIL (active work PLAN weaken) | **Do not close** — spawn repair | Keep OPEN — critical invariant regression |
| Block 6 FAIL (bootstrap cannot recover feet-seller) | **Thin follow-on OR repair** | If tool bug → repair spawn; if heuristic insufficient → human promotes **Phase 37 — plan recovery hardening** scaffold |
| VERIFY-NOTES missing blocks, evidence dir absent, or pre-bootstrap feet-seller still shows `plan_missing` block | **Do not close** — send back S03-01 or spawn repair | Keep OPEN |
| VERIFY PASS but agents still skip bootstrap in practice (harness gap) | **`no successor`** default | Document residual; optional future phase for enforce default / bridge advisory — **not** auto-spawn unless human promotes |
| Human wants PlanExists bridge (§2.4) next | **Thin follow-on phase** | Human promotes Phase 37+ with bridge theme; close Phase 36 first with **`no successor`** unless scaffold ready |
| Human wants HTTP POST plan routes / OpenAPI parity | **Residual note only** | Default **`no successor`** — defer to future HTTP/OpenAPI phase |
| Human wants hosted SaaS next | **Not** Trace core successor | Close with **`no successor`**; cloud remains separate product/repo |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| PlanExists bridge (§2.4 defer) | Accept — explicit bootstrap path shipped; revisit if dogfood still skips |
| HTTP POST plan mutation routes | Accept — MCP primary; `GET /v1/plans` exists |
| MCP `trace_loop action=gate` | Accept — CLI/hook/GUI HTTP sufficient |
| Bootstrap help human-refinement note (S02-02 low) | Accept — document in DR-HANDOFF residuals |
| Loop status `advisories[]` for goal-structure warning | Accept — plan show + MCP show field sufficient |
| `WarnIfTraceDirWithoutConfig` unit test gap | Accept — init/install stderr nudge shipped |
| Enforce default `warn` when `.trace/` without config | Accept defer — opt-in model preserved |
| Goal/plan surface UX beyond TaskDetail | Accept defer — terminal advisory + bootstrap hint sufficient for Phase 36 |
| Feet-seller post-bootstrap planner minimal quality | Accept — human refinement via `create-coarse` / `deep` expected per PLAN §2.2 |
| Hosted SaaS / cloud | Separate product — not Phase 36 residual |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# DR-HANDOFF — Phase 36

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 35 CLOSED |
| Theme | Planning model alignment — MCP plan / bootstrap / install + terminal gate honesty |
| Outcome | Agents satisfy PlanExists via trace_plan + bootstrap; terminal DONE honest advisory; feet-seller recovered; active PLAN preserved |
| Successor decision | **no successor** |
| Residuals (non-blocking) | PlanExists bridge defer; HTTP POST plan defer; enforce default defer; bootstrap refinement UX; … |
| Close owner | P36-S03-02 |
| Verify | Cite VERIFY-NOTES + experiments/runs/…-p36-s03-01-verify/evidence/ |

## Scope checklist

- [x] S00 investigate (`INVESTIGATION.md`)
- [x] S01 plan (`PLAN.md`)
- [x] S02 implement + tests + review
- [x] S03 VERIFY + successor documented (**no successor**)
```

If thin follow-on (bootstrap hardening or bridge): set `Successor decision` to **Phase 37 — \<theme\>**; first runnable **P37-00**; confirm scaffold present before close.

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md
EVID=$(ls -d experiments/runs/*-p36-s03-01-verify/evidence 2>/dev/null | tail -1)
test -n "$EVID" && test -d "$EVID"

go build -o /tmp/trace ./cmd/trace

# Block 0 spot-check (subset)
go test ./internal/loop/... ./internal/mcp/... -count=1 \
  -run 'Greenfield|FeetSeller|ActiveWork|TerminalPlanGap' 

FEET="/home/ali/Desktop/feet seller telegram app"
STEP1=33247e2d-aa10-4b25-b194-4b7afb5a6359
GOAL=353b12a4-57dd-4d68-8379-b2024e064733

# If pre-bootstrap evidence captured, jq spot-check archived JSON:
# jq '.allowed, .violations[0].reason_code' "$EVID/02-feet-step1-done-gate-pre-bootstrap.json"

# Live spot (may differ post-bootstrap block 6):
/tmp/trace -C "$FEET" plan show --goal "$GOAL" | jq '{current_scope_id, has_deep: (.current_deep_plan != null)}'

# Active-work regression (temp)
ACT=$(mktemp -d)
/tmp/trace -C "$ACT" init
G=$(/tmp/trace -C "$ACT" add goal --title spot | jq -r '.id')
T=$(/tmp/trace -C "$ACT" add task --title spot --goal-id "$G" | jq -r '.id')
/tmp/trace -C "$ACT" transition --task "$T" --to IN_PROGRESS --reason spot >/dev/null
/tmp/trace -C "$ACT" loop gate --task "$T" --for edit | jq '.allowed, .reason_code'

# MCP tool count
go test ./internal/mcp/... -count=1 -run TestRegisteredToolNames_IncludesTracePlan
```

Confirm VERIFY-NOTES: overall PASS; blocks 0–7 ticked; DESIGN-LOCKS must-fix addressed; DR-HANDOFF still OPEN before this row closes it.

### TODO.md / AGENTS.md updates (on APPROVE)

**If `no successor`:**

1. `docs/TODO.md` orchestrator paste: Phase 00–36 complete; no active phase (idle awaiting human promotion).
2. Phase boards table: Phase 36 status `done`; Next `—`.
3. `AGENTS.md` Current focus: Phase 36 complete; planning model alignment shipped (MCP `trace_plan`, bootstrap, install contract, terminal gate honesty); cloud remains separate product.

**If thin follow-on:** point Active phase / next runnable at scaffolded phase planner; never TBD.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-03-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P36-S03-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** no successor  (or Phase 37 / P37-00 | pending repair spawn)

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| S02 acceptance subset | |
| Feet-seller plan_exists post-bootstrap | |
| Active work plan_missing block | |
| 16-tool MCP lock | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Next
(idle / P37-00 / P36-S03-02a)
```

### On FAIL / repair spawn

Insert immediately below this row on phase board:

| Order | ID | Role |
|------:|----|------|
| 633a | P36-S03-02a | implement repair |
| 633b | P36-S03-02b | review repair |

Keep DR-HANDOFF **OPEN**. Do not mark Phase 36 done.

## Role work

1. Fresh-session re-verify S03-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` (findings + confidence + successor decision).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor **`no successor`** or scaffolded Phase 37 (**never TBD**).
4. Update `docs/TODO.md` + `AGENTS.md` per decision table.
5. If thin follow-on: create runnable next-phase scaffold before close (protocol duties in agent-loop-protocol).
6. Do **not** rewrite S00–S02 `done` history or S03-01 VERIFY-NOTES content except to cite them.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with explicit successor (**no successor** or Phase 37 / P37-00 — never TBD)
- [ ] If follow-on: runnable scaffold present (README + phase planner + stubs + board + TODO index)
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 36 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

Idle (**no successor**) — or successor first row / repair spawn
