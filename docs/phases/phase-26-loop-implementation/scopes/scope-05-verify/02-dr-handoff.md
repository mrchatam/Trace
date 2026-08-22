# P26-S05-02 — DR-HANDOFF Phase 26 close

## Metadata
- id: P26-S05-02
- todo_ids: [P26-S05-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of S05-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 26 DR-HANDOFF** with explicit successor decision (**never TBD**). Scaffold next phase **only** if successor is promoted. Phase 26 complete when this row is `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks
- [01-verify.md](01-verify.md) — locked verify floor
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S05-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [PLAN.md](../scope-01-planning/PLAN.md)
- [docs/TODO.md](../../../../TODO.md)
- [AGENTS.md](../../../../../AGENTS.md)
- Pattern: [P24 S05-02 review](../../../phase-24-agent-effectiveness-investigation/scopes/scope-05-phase-verify/02-scope-review.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S05-01 implementer. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-05-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p26-s05-01-verify/evidence/` |
| Results row | `experiments/RESULTS.md` |
| Phase handoff | `DR-HANDOFF.md` |
| Scope deliverables | S00 `AUDIT.md`, S01 `PLAN.md`, S02–S04 implement+review Notes |

## Locked DR-HANDOFF close policy (FINAL — S05-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S05-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S05-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S05-01 `done`; P25-2 **PASS** on verify run; `go test ./internal/...` PASS |
| Default successor | **`no successor`** — prefer until fresh E03 dogfood if D1–D6 verified |
| Alternative successor | **Phase 27** — protocol/measurement (INT-08/10 + INT-07) if promotion/reset work but graph still thin |
| Regression path | Spawn `P26-S05-02a` implement + `02b` review; **do not** close Phase 26 |
| Must not | Leave `Successor decision: TBD`; rewrite Phase 25 `done` history; implement new product scope in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 26 board rows `done` |

### Successor decision table (locked)

| Outcome (from S05-01 + spot-check) | Decision | Phase 27 scaffold? |
|------------------------------------|----------|-------------------|
| D1–D6 verified; P25-1/2/3 PASS; D4 reset works | **`no successor`** until E03 dogfood | **No** |
| D1/D2/D4 pass; P25-3 FAIL or thin graph; installer OK | **Phase 27** — INT-08/10 protocol + INT-07 graph honesty | **Yes** (minimal runnable) |
| P25-2 FAIL or unit tests FAIL | **Do not close** — spawn repair forward | No |
| Human named different successor before this row | Document in Notes with evidence | Per human |

**E03** here means a **future** full dogfood session — not required to close Phase 26 when verify re-score passes.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00: Loop audit complete (`AUDIT.md`)
- [ ] S01: Planning complete (`PLAN.md`)
- [ ] S02: P25-A implementation + review done
- [ ] S03: P25-B implementation + review done
- [ ] S04: Installer fix implementation + review done
- [ ] S05: VERIFY — `score.sh G1 --p25` or documented partial; P25-2 PASS; successor documented

### Residuals to list on close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| INT-04 hook permissive (`TRACE_TASK_ID` unset) | Document; harness hardening deferred |
| P25-4 operator attestation | Manual protocol; not blocking close |
| Option B partial verify | Note in DR-HANDOFF if S05-01 used fallback |
| Claude fallback rules vs cursor MDC parity | S04 wired both; spot-check optional |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** (until E03 dogfood) |
| Phase 26 outcome | P25-A promotion + P25-B saturation/reset + P25-2 installer wiring; E02 P25-2 gap closed |
| E02 → verify delta | P25-2: FAIL → PASS |
| Residuals (non-blocking) | INT-04 hook enforcement beyond install text; P25-4 attestation |
| Forward | Human promotes next phase when ready |
```

If **Phase 27** successor:

```markdown
| Successor decision | **Phase 27** — protocol/measurement + graph richness (INT-08/10, INT-07) |
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair rows; successor stays pending until pass.

## Phase handoff scaffold (mandatory only if Phase 27 successor)

Per agent-loop-protocol, before marking `done` with Phase 27 successor:

| Artifact | Path (default) |
|----------|----------------|
| Phase README | `docs/phases/phase-27-*/README.md` |
| Phase planner | `docs/phases/phase-27-*/00-PHASE-PLANNER.md` |
| DR-HANDOFF open | `docs/phases/phase-27-*/DR-HANDOFF.md` (**OPEN**) |
| Scope stub | ≥1 scope with `00-PLANNER.md`, `01-*.md`, `02-*.md`, `SCOPE-TODOS.md` |
| Board file | `docs/TODO/phase-27.md` — **P27-00** first `pending` |
| Index link | `docs/TODO.md` |

If **`no successor`**: no Phase 27 folder required; VERIFY Notes + DR-HANDOFF must say `no successor`.

## Locked re-verify commands (minimum — reviewer)

```bash
cd /home/ali/Desktop/Trace

# Evidence exists
EVID=$(ls -d experiments/runs/*-p26-s05-01-verify/evidence 2>/dev/null | tail -1)
test -d "$EVID" && ls -la "$EVID"
test -f docs/phases/phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md

# Unit + install (do not trust S05-01 alone)
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
go test ./internal/... -count=1
go test ./internal/install/... -count=1 -run TestInstallCursorRules

# P25-2 closure (temp dir — independent of G1 state)
tmpdir=$(mktemp -d)
./bin/trace -C "$tmpdir" init
./bin/trace -C "$tmpdir" install cursor --write
grep -qi 'Parent orchestrator' "$tmpdir/.cursor/rules/trace-enforcement.mdc"
rm -rf "$tmpdir"

# If Option A was claimed — re-run score (or confirm score artifact in evidence)
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd experiments/ab-p25-gap-pass-validation
./score.sh G1 --p25
```

Broader walk of VERIFY-NOTES D1–D6 table encouraged; **minimum** commands above mandatory for APPROVE.

## Review checklist

### Blockers

- [ ] Missing VERIFY-NOTES or evidence dir
- [ ] S05-01 `failed` without pending spawn
- [ ] P25-2 still FAIL on independent grep
- [ ] `go test ./internal/...` FAIL
- [ ] Successor left **TBD** at close
- [ ] DR-HANDOFF CLOSED while verify bar failed

### High

- [ ] VERIFY-NOTES verdict PASS but spot-checks fail
- [ ] Option B used without documenting P25-3 SKIP
- [ ] Phase 27 successor chosen but scaffold missing
- [ ] AGENTS.md / TODO index not updated when phase complete

### Medium / low

- [ ] Evidence stale vs git SHA — note in REVIEW-NOTES
- [ ] Missing RESULTS.md row

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Spawn `P26-S05-02a` implement + `02b` review immediately below |
| medium | Prefer spawn unless ≤5-line doc fix |
| low / nit | List in REVIEW-NOTES; do not block close |

## Evidence artifacts (reviewer output)

- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update [DR-HANDOFF.md](../../DR-HANDOFF.md) on APPROVE
- Update `docs/TODO.md` phase-26 status + orchestrator paste in `AGENTS.md` (or cite in board Notes)
- Register Phase 27 board only if successor promoted

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

## Exit criteria

- [ ] Independent re-verify minimum commands PASS
- [ ] S05-01 evidence reviewed (VERIFY-NOTES + archive)
- [ ] P25-2 closure confirmed (FAIL → PASS vs E02)
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with explicit successor (`no successor` default)
- [ ] Phase 27 scaffold **runnable** if that successor chosen
- [ ] `docs/TODO.md` reflects phase complete / next runnable
- [ ] Confidence **high** (or **medium** with explicit residuals)
- [ ] All Phase 26 board rows `done`

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF without independent P25-2 re-grep
- Rewriting Phase 25 or S02–S04 `done` history
- Implementing new product features in this review row

## Minimal todos

- [ ] Read VERIFY-NOTES + evidence manifest
- [ ] Re-run locked re-verify commands
- [ ] Execute review checklist
- [ ] Close DR-HANDOFF with locked template
- [ ] Scaffold Phase 27 if promoted; else document `no successor`
- [ ] Write REVIEW-NOTES.md; set row `done`

## Next

**none** (if `no successor`) — or **P27-00** if Phase 27 promoted
