# P31-S02-02 — DR-HANDOFF Phase 31 close

## Metadata
- id: P31-S02-02
- todo_ids: [P31-S02-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent **fresh-session** review of S02-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 31 DR-HANDOFF** with successor **Phase 32** / first runnable **P32-00** (**never TBD**). Update `docs/TODO.md` + `AGENTS.md` current focus. Phase 31 complete when this row is `done`. **No product code.** Do **not** implement Phase 32 in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S02-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S02-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [GAPS.md](../scope-00-inventory/GAPS.md)
- [REVIEW-NOTES.md](../scope-01-tests/REVIEW-NOTES.md) — S01-02 PASS
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-31.md](../../../../TODO/phase-31.md)
- [docs/TODO/phase-32.md](../../../../TODO/phase-32.md)
- [AGENTS.md](../../../../../../AGENTS.md)
- Phase 32 scaffold: `docs/phases/phase-32-graph-first-gui/` (`00-PHASE-PLANNER.md`, `DESIGN-LOCKS.md`, `DR-HANDOFF.md`)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S02-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-02-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p31-s02-01-verify/evidence/` |
| Phase handoff | `DR-HANDOFF.md` |
| Gaps inventory | `scopes/scope-00-inventory/GAPS.md` |
| S01 review | `scopes/scope-01-tests/REVIEW-NOTES.md` |
| Phase board | `docs/TODO/phase-31.md` |
| Successor board | `docs/TODO/phase-32.md` (first runnable **P32-00**) |

## Locked DR-HANDOFF close policy (FINAL — S02-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S02-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S02-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S02-01 `done`; verify blocks 0–4 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 32** — graph-first GUI; first runnable **P32-00** |
| Regression path | Do **not** close; spawn `P31-S02-02a` implement + `02b` review; successor = **`pending repair spawn`** (still not TBD) |
| Must not | Leave successor `TBD`; redesign store path; start implementing Phase 32; rewrite S00–S01 `done` history; ship product in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 31 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S02-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY floor green; residuals only as listed (G2/G3/G4, multi-open, stubs) | **Phase 32** / **P32-00** | Close DR-HANDOFF; mark Phase 31 **done**; point TODO/AGENTS at Phase 32 |
| Test FAIL / repro FAIL / join changed / warn missing / silent delete / G1 or G5 or G6 missing | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b; successor = `pending repair spawn` |
| VERIFY-NOTES missing blocks or evidence dir absent | **Do not close** — spawn repair or send back S02-01 | Keep OPEN |
| VERIFY PASS but human holds Phase 32 start | Still set successor **Phase 32** / **P32-00**; note human hold in AGENTS | Do **not** write TBD |

**Never** leave successor as `TBD` when marking this row `done`.

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-31-stray-db-residual-tests/scopes/scope-02-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p31-s02-01-verify/evidence || ls experiments/runs/ | grep p31-s02-01
go test ./internal/store/ \
  -run 'TestOpenWarnsWhenRootStubPresent|TestOpenExistingWarnsWhenRootStubPresent|TestOpenQuietWhenNoRootStub|TestOpenLeavesRootStubUntouched|TestOpenQuietWhenRootStubIsDirectory' \
  -count=1
test -x scripts/repro-stray-trace-db.sh && bash scripts/repro-stray-trace-db.sh
grep -n '/trace.db' .gitignore fixtures/x0/.gitignore
grep -n 'warnIfStrayRootTraceDB\|traceDirName\|dbFileName\|IsRegular' internal/store/open.go
grep -n 'openStore\|once per\|suppress' CONTRIBUTING.md AGENTS.md
```

Confirm VERIFY-NOTES: overall PASS; residuals listed; does not claim store-path redesign; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00 inventory (`GAPS.md`)
- [ ] S01 tests + review
- [ ] S02 VERIFY + successor **Phase 32** documented

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| Agents can still create root stubs | Mitigated by warn + `/trace.db` gitignore; agent hygiene |
| Optional delete of root stub | Future-only — not shipped |
| Warn once per `openStore` (multi-open CLI/MCP/HTTP) | Intentional; G6 documented; no suppress flag |
| G2 / G3 / G4 | Nice skip / deferred with reason in GAPS.md — do not reopen path design |

### Optional docs cleanup on close (non-blocking)

S01-02 nit: AGENTS.md orchestrator paste may still say Phase 31 “Next: `P31-S00-00`”. On APPROVE, refresh orchestrator paste + Current focus to Phase 32 / `P32-00` as part of the handoff update (not a product change).

### DR-HANDOFF.md update template (on APPROVE — successor Phase 32)

```markdown
# DR-HANDOFF — Phase 31

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-21 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 30 closed (`P30-S03-02`) |
| Theme | Extra testing for stray root `trace.db` hygiene |
| Outcome | G1+G5+G6 shipped; S01-02 PASS; S02-01 VERIFY PASS; no path redesign |
| Successor decision | **Phase 32** — graph-first GUI; first runnable **P32-00** |
| Residuals (non-blocking) | Agent stubs; optional delete future-only; multi-open once-per-openStore; G2/G3/G4 deferred |
| Close owner | P31-S02-02 |

## Scope checklist

- [x] S00 inventory (`GAPS.md`)
- [x] S01 tests + review
- [x] S02 VERIFY + successor **Phase 32** documented
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-02-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P31-S02-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** Phase 32 / P32-00  (or pending repair spawn)

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| Five store tests | |
| Repro script | |
| gitignore / open.go / G6 | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Next
P32-00 (or P31-S02-02a)
```

## Role work

1. Fresh-session re-verify S02-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` in this folder (findings + confidence + successor Phase 32).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick checklist; set successor **Phase 32** / **P32-00** (never TBD).
4. Update `docs/TODO.md`: Phase 31 → **done**; Next Phase 32 / `P32-00`; orchestrator paste.
5. Update `AGENTS.md` current focus → Phase 32 (refresh stale Phase 31 paste if still present).
6. Confirm Phase 32 board + `00-PHASE-PLANNER` already exist; do **not** deep-plan or implement Phase 32 here.
7. Do **not** rewrite S00–S01 `done` history or S02-01 VERIFY-NOTES content except to cite them.

## Todo updates

Status + notes on **P31-S02-02**; may update TODO.md / AGENTS.md / DR-HANDOFF; may spawn repair rows below this row if needed.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] DR-HANDOFF CLOSED with successor **Phase 32** / **P32-00** (never TBD)
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All Phase 31 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence
- [ ] Board row done with Notes

## Next

`P32-00` (after Phase 31 complete)
