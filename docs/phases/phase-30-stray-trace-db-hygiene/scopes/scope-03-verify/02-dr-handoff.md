# P30-S03-02 — DR-HANDOFF Phase 30 close

## Metadata
- id: P30-S03-02
- todo_ids: [P30-S03-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents]
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 30 DR-HANDOFF** with explicit successor (**never TBD**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 30 complete when this row is `done`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [PLAN.md](../scope-01-plan/PLAN.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-30.md](../../../../TODO/phase-30.md)
- [AGENTS.md](../../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p30-s03-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Investigation | `scopes/scope-00-investigate/INVESTIGATION.md` |
| Plan | `scopes/scope-01-plan/PLAN.md` |
| Phase board | `docs/TODO/phase-30.md` |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; verify blocks 0–4 green per VERIFY-NOTES + independent spot-check |
| Default successor | **`no successor`** — S00 found **no** store-path defect (agent hygiene / INTAKE confirmed); hygiene T1–T4 shipped |
| Exception | Only if independent spot-check **overturns** S00 (Trace actually creates/opens root `trace.db`) — **do not close**; spawn `P30-S03-02a`/`02b` repair; do **not** invent a new phase theme here |
| Cloud / hosted SaaS | **Not** a Phase 30 successor — separate product/repo |
| Regression path | Spawn `P30-S03-02a` implement + `02b` review; **do not** close Phase 30 |
| Must not | Leave `Successor decision: TBD`; claim store-path redesign needed without evidence; rewrite S00–S02 `done` history; ship product in this row; start implementing a new phase |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 30 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S03-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; residuals only as listed; S00 agent-hygiene stands | **`no successor`** | Close DR-HANDOFF; mark Phase 30 **done** in TODO/AGENTS; idle orchestrator paste |
| Test FAIL / init creates root db / join changed / warn missing / silent delete | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b |
| VERIFY-NOTES missing blocks or evidence dir absent | **Do not close** — spawn repair or send back S03-01 | Keep OPEN |
| Spot-check proves Trace dual-store / root path open (overturns S00) | **Do not close** — spawn repair; document overturn | Keep OPEN; human promotes any follow-on phase |
| VERIFY PASS but human names different successor | Document human theme in Notes | Scaffold only if human promotes |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-30-stray-trace-db-hygiene/scopes/scope-03-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p30-s03-01-verify/evidence || ls experiments/runs/ | grep p30-s03-01
go test ./internal/store/ \
  -run 'TestOpenWarnsWhenRootStubPresent|TestOpenExistingWarnsWhenRootStubPresent|TestOpenQuietWhenNoRootStub|TestOpenLeavesRootStubUntouched' \
  -count=1
grep -n '/trace.db' .gitignore fixtures/x0/.gitignore
grep -n 'warnIfStrayRootTraceDB\|traceDirName\|dbFileName' internal/store/open.go
# Confirm join still .trace + trace.db (Stat-only warn; no delete of root stub in open.go)
```

Confirm VERIFY-NOTES: overall PASS; residuals listed; does not claim store-path redesign; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 independent investigation
- [ ] S01 plan
- [ ] S02 implement + review
- [ ] S03 VERIFY + successor documented (**never TBD**; default **no successor**)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| Agents can still create root stubs | Mitigated by warn + `/trace.db` gitignore; agent hygiene |
| Optional delete of root stub | Future-only — not shipped |
| Warn once per open in long-lived serve | Acceptable; no suppress flag in Phase 30 |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# DR-HANDOFF — Phase 30

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-21 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 29 (`P29-S07-02` CLOSED → this phase) |
| Theme | Stray root `trace.db` hygiene |
| Outcome | S00 agent hygiene confirmed; T1–T4 shipped; VERIFY PASS |
| Successor decision | **no successor** |
| Residuals (non-blocking) | Agents may still create stubs; optional delete future-only |
| Close owner | P30-S03-02 |

## Scope checklist

- [x] S00 independent investigation
- [x] S01 plan
- [x] S02 implement + review
- [x] S03 VERIFY + successor documented
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

## Role work

1. Fresh-session re-verify S03-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` in this folder (findings + confidence + successor pick).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor (never TBD; default **no successor**).
4. Update `docs/TODO.md`: Phase 30 → **done**; orchestrator paste idle / no next phase unless human promotes.
5. Update `AGENTS.md` current focus (Phase 30 closed; no successor default).
6. Do **not** rewrite S00–S02 `done` history or S03-01 VERIFY-NOTES content except to cite them.

## Todo updates

Status + notes on **P30-S03-02**; may update TODO.md / AGENTS.md / DR-HANDOFF; may spawn repair rows below this row if needed.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with successor **not** TBD (default **no successor**)
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 30 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

Successor per DR-HANDOFF (default **none**).
