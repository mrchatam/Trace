# P28-S07-02 — Residual-wave DR-HANDOFF close

## Metadata
- id: P28-S07-02
- todo_ids: [P28-S07-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of S07-01 residual-wave verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close the Residual wave section** in `DR-HANDOFF.md` only — leave **S05 CLOSED** history intact. Update `docs/TODO.md` + `AGENTS.md` when wave complete. Successor **never TBD**. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + Phase handoff
- [00-PLANNER.md](00-PLANNER.md) — S07-00 locks
- [01-verify.md](01-verify.md) — locked residual-wave verify floor
- [VERIFY-NOTES-RESIDUAL-WAVE.md](VERIFY-NOTES-RESIDUAL-WAVE.md) — produced by S07-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [docs/TODO.md](../../../../TODO.md)
- [AGENTS.md](../../../../../AGENTS.md)
- [../scope-06-r6-fm-residuals/SCOPE-TODOS.md](../scope-06-r6-fm-residuals/SCOPE-TODOS.md)
- Pattern: [../scope-05-verify/02-dr-handoff.md](../scope-05-verify/02-dr-handoff.md) — S05 closed Residual-wave-style section only here

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S07-01 verifier. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Residual-wave verify notes | `scopes/scope-07-residual-wave-verify/VERIFY-NOTES-RESIDUAL-WAVE.md` |
| Evidence archive | `experiments/runs/…-p28-s07-01-verify/evidence/` |
| FR Notes / decision | `scopes/scope-06-r6-fm-residuals/FM*-NOTES.md`, `FM07-DECISION.md` |
| S05 history | `DR-HANDOFF.md` § S05 close — **immutable** |
| Residual wave section | `DR-HANDOFF.md` § Residual wave — **this row closes** |
| Forward queue | `docs/TODO/forward-p28-residuals.md` — still superseded index only |

## Locked close policy (FINAL — S07-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S07-01** — verify floor + VERIFY-NOTES-RESIDUAL-WAVE; Residual wave stays **OPEN** |
| Who closes Residual wave | **S07-02 only** |
| S05 section | Remains **CLOSED** (immutable) — never rewrite |
| Status on pass | Residual wave section → **CLOSED** |
| Closure prerequisite | S07-01 `done` PASS; FR-P28-01…07 reviews APPROVE; floor green per notes |
| Default successor | **`no successor`** when FR closed + regression green |
| Phase 29 | **Human promote only** — do not invent a phase; scaffold only if human asks |
| Regression path | Spawn `P28-S07-02a` implement + `02b` review; **do not** close Residual wave |
| Must not | Leave successor `TBD`; run `./prepare.sh`; rewrite S05 CLOSED history; edit S05 VERIFY-NOTES; ship product |
| Phase residual wave complete | **Yes** when this row `done` + Residual wave **CLOSED** + S07 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S07-01 + independent spot-check) | Decision | Phase 29 scaffold? |
|------------------------------------------------|----------|-------------------|
| FR-P28-01…07 closed; unit/cmd/install/hook green; dual-lane honest (no prepare) | **`no successor`** | **No** |
| Unit/cmd/install/hook FAIL on re-run | **Do not close** — spawn repair | No |
| Missing / contradictory FR Notes vs board APPROVE | **Do not close** — spawn repair | No |
| VERIFY used `./prepare.sh` or conflated thin/rich | **Do not close** — spawn repair; evidence tainted | No |
| VERIFY PASS but human requests a new theme | **Phase 29** — theme named by human in Notes | **Yes** (minimal runnable) only if human promotes |
| Human named different successor before this row | Document in Notes with evidence | Per human |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
# Do NOT prepare.sh
GOPROXY=direct go test ./internal/install/... -run 'CursorLoopGateFailClosed|HookDrift' -count=1
test -f docs/phases/phase-28-residuals-validation/scopes/scope-07-residual-wave-verify/VERIFY-NOTES-RESIDUAL-WAVE.md
# At least one FR evidence path:
test -f docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals/FM01-NOTES.md
test -f docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals/FM07-DECISION.md
# Forward queue still superseded index:
grep -qi superseded docs/TODO/forward-p28-residuals.md
# Optional directed re-score (no wipe) — only if S07-01 claimed live re-score:
# export TRACE_BIN=$PWD/bin/trace
# cd experiments/ab-p25-gap-pass-validation && P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed --test
```

Confirm VERIFY-NOTES-RESIDUAL-WAVE: FR table closed; D1–D4/X1 listed deferred; dual-lane not conflated; no prepare wipe; S05 VERIFY-NOTES not rewritten.

### Residual-wave scope checklist (tick on APPROVE in DR-HANDOFF)

From [DR-HANDOFF.md](../../DR-HANDOFF.md) Residual wave section:

- [ ] S06 planner scaffold (`P28-S06-00`)
- [ ] S06 FR-P28-01…07 implement + review pairs (`P28-S06-01`…`14`)
- [ ] S07 residual-wave VERIFY + DR-HANDOFF close (`P28-S07-00`…`02`)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| FR-P28-D1 auto-spawn | Deferred — human gate |
| FR-P28-D2 Graphiti | Deferred — human promote |
| FR-P28-D3 RESULTS parser | Deferred — env attestation closed R5 |
| FR-P28-D4 Hook Option B | Deferred — Option A locked |
| FR-P28-X1 Daemon/HTTP | Deferred — project laws |
| FM-07 warn-only | By design (closed as warn-only) |

### DR-HANDOFF Residual wave update template (on APPROVE — default `no successor`)

Update **only** the Residual wave section (keep S05 CLOSED history above intact):

```markdown
## Residual wave (post-close) — CLOSED

**Status:** **CLOSED** (YYYY-MM-DD)

| Field | Value |
|-------|-------|
| Opened | 2026-08-20 |
| Closed | YYYY-MM-DD |
| Theme | R6 / FM-01/02/04/07/08/09/10 (FR-P28-01…07) |
| Outcome | FR-P28-01…07 closed; residual-wave VERIFY PASS; dual-lane honest (no prepare) |
| Successor decision | **no successor** |
| Residuals (non-blocking) | D1–D4/X1 deferred; FM-07 warn-only by design |
| Forward | Human promotes Phase 29 only if needed |
| Close owner | P28-S07-02 |
```

Tick residual-wave scope checklist items `[x]`. Leave S05 section untouched.

If **Phase 29** (human only):

```markdown
| Successor decision | **Phase 29** — <explicit theme from human> |
```

If verify **failed**: keep Residual wave **OPEN**; spawn repair; successor = **`pending repair spawn`**.

## Role work

1. Fresh-session re-verify S07-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES-RESIDUAL-WAVE.md` in this folder (findings + confidence + successor pick).
3. On APPROVE: update `DR-HANDOFF.md` Residual wave → CLOSED; tick checklist; set successor (never TBD).
4. Update `docs/TODO.md`: Phase 28 residual wave → idle / done; orchestrator paste for idle or next human promote.
5. Update `AGENTS.md` current focus (wave closed; no successor default).
6. Do **not** rewrite S05 CLOSED history or S05 VERIFY-NOTES.

## Todo updates

Status + notes on **P28-S07-02**; may update TODO.md / AGENTS.md / DR-HANDOFF Residual wave section; may spawn repair rows below this row if needed.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES-RESIDUAL-WAVE.md`
- [ ] Residual wave CLOSED in DR-HANDOFF with successor **not** TBD
- [ ] S05 CLOSED section unchanged
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All P28 residual-wave board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence

## Minimal todos

- [ ] Spot-check VERIFY-NOTES-RESIDUAL-WAVE + hook smoke + one FR path
- [ ] Write REVIEW-NOTES-RESIDUAL-WAVE.md
- [ ] Close Residual wave in DR-HANDOFF or spawn repair
- [ ] Update TODO.md + AGENTS.md
- [ ] Mark P28-S07-02 `done` / `failed` / `blocked`

## Next

— (phase residual wave complete on APPROVE + successor **no successor** default)
