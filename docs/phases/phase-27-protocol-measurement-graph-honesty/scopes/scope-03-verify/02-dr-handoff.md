# P27-S03-02 — DR-HANDOFF Phase 27 close

## Metadata
- id: P27-S03-02
- todo_ids: [P27-S03-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of S03-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 27 DR-HANDOFF** with explicit successor decision (**never TBD**). Scaffold Phase 28 **only** if successor table selects it. Phase 27 complete when this row is `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks
- [01-verify.md](01-verify.md) — locked verify floor + pass/fail policy
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S03-01
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [docs/TODO.md](../../../../TODO.md)
- [AGENTS.md](../../../../../AGENTS.md)
- Pattern: [P26 S05-02](../../../phase-26-loop-implementation/scopes/scope-05-verify/02-dr-handoff.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S03-01 verifier. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-03-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p27-s03-01-verify/evidence/` |
| Harness | `experiments/ab-p25-gap-pass-validation/score.sh` (enforce upgrade) |
| Phase handoff | `DR-HANDOFF.md` |
| Prior scopes | S00 `AUDIT.md`; S01/S02 implement+review Notes |

## Locked DR-HANDOFF close policy (FINAL — S03-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S03-01** — verify floor + enforce upgrade + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S03-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S03-01 `done`; INT-07/08/10 closure signals PASS per VERIFY-NOTES; `go test ./internal/...` PASS on re-run |
| Default successor | **`no successor`** — INT-07/08/10 delivered; build thin FAIL remains expected measurement |
| Alternative successor | **Phase 28** — only if table row below selects it |
| Regression path | Spawn `P27-S03-02a` implement + `02b` review; **do not** close Phase 27 |
| Must not | Leave `Successor decision: TBD`; rewrite Phase 26/S01–S02 `done` history; implement new product scope in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 27 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S03-01 + independent spot-check) | Decision | Phase 28 scaffold? |
|------------------------------------------------|----------|-------------------|
| Closure signals PASS; only expected thin/enforce residuals; P25-1/2 PASS | **`no successor`** | **No** |
| Product honesty regression (thin still exports under `--strict --enforce`) | **Do not close** — spawn repair | No |
| `score.sh` missing `--enforce` or still WARN-only on thin | **Do not close** — spawn repair | No |
| Unit/honesty tests FAIL on re-run | **Do not close** — spawn repair | No |
| VERIFY PASS but human requests follow-on dogfood / new INT work | **Phase 28** — theme named by human in Notes | **Yes** (minimal runnable) |
| VERIFY PASS but blocking residual needs a new phase (not fixable by spawn) | **Phase 28** — name residual theme explicitly | **Yes** |
| Human named different successor before this row | Document in Notes with evidence | Per human |

**Never** leave successor as `TBD` when marking this row `done`.

### Human-gated criteria (non-blocking)

| Item | Gate | Effect on close |
|------|------|-----------------|
| P25-4 operator attestation | Human / RESULTS.md | Document residual; **do not** block close |
| Session-B directed dogfood (P25-3b PASS) | Human agent session | Optional; thin directed FAIL **OK** for close |
| Promote Phase 28 without regression | Human | Use Phase 28 row in table; scaffold required |

If a human-gated item is **required** by human before close: mark this row **`blocked`** awaiting evidence under `docs/verification/` — do **not** self-claim.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00: Investigation complete (`AUDIT.md`)
- [ ] S01: Protocol v2 implementation + review done (INT-08/10)
- [ ] S02: Graph honesty implementation + review done (INT-07)
- [ ] S03: VERIFY — enforce upgrade + score/rubric evidence; successor documented (**never TBD**)

### Residuals to list on close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| Build-only P25-3a FAIL | Expected RUBRIC baseline |
| Thin G2 FAIL under `--strict --enforce` | Expected harness alignment |
| Directed P25-3b without Session-B | Optional human dogfood |
| P25-4 attestation | Manual protocol |
| S02-02 BLOCKING duplicate orphan message | Low; deferred polish |
| FM-07 warn-only drift | Keep warn-only |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 27 outcome | INT-08/10 protocol v2 (P25-3a/3b + arms) + INT-07 seed export graph honesty; score.sh `--strict --enforce` |
| Verify delta vs Phase 26 | P25-3a thin FAIL remains expected; product+harness now fail closed on dishonest thin export |
| Residuals (non-blocking) | P25-4; optional Session-B dogfood; BLOCKING duplicate msg |
| Forward | Human promotes next phase when ready |
```

If **Phase 28** successor:

```markdown
| Successor decision | **Phase 28** — <explicit theme from decision table> |
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair rows; successor stays pending until pass (**still not TBD** — say `pending repair spawn`).

## Phase handoff scaffold (mandatory only if Phase 28 successor)

Per agent-loop-protocol, before marking `done` with Phase 28 successor:

| Artifact | Path (default) |
|----------|----------------|
| Phase README | `docs/phases/phase-28-*/README.md` |
| Phase planner | `docs/phases/phase-28-*/00-PHASE-PLANNER.md` |
| DR-HANDOFF open | `docs/phases/phase-28-*/DR-HANDOFF.md` (**OPEN**) |
| Scope stub | ≥1 scope with `00-PLANNER.md`, `01-*.md`, `02-*.md`, `SCOPE-TODOS.md` |
| Board file | `docs/TODO/phase-28.md` — **P28-00** first `pending` |
| Index link | `docs/TODO.md` |

If **`no successor`**: no Phase 28 folder required; VERIFY Notes + DR-HANDOFF must say `no successor`.

## Board / index update duties (on APPROVE)

| Update | Duty |
|--------|------|
| `DR-HANDOFF.md` | CLOSED + successor row filled |
| `docs/TODO/phase-27.md` | This row `done` + Notes (verdict, successor) |
| `docs/TODO.md` | Phase 27 status → done (or next = P28-00 if promoted); orchestrator paste |
| `AGENTS.md` | Current focus / orchestrator paste reflects close or Phase 28 |
| Phase 28 board | Only if successor = Phase 28 |

## Locked re-verify commands (minimum — reviewer)

```bash
cd /home/ali/Desktop/Trace

# Evidence exists
EVID=$(ls -d experiments/runs/*-p27-s03-01-verify/evidence 2>/dev/null | tail -1)
test -d "$EVID" && ls -la "$EVID"
test -f docs/phases/phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md

# Enforce upgrade present
rg -n 'strict --enforce|--enforce' experiments/ab-p25-gap-pass-validation/score.sh

# Unit + honesty (do not trust S03-01 alone)
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
go test ./internal/... -count=1
go test ./cmd/trace/... -run 'SeedExport|Enforce|Strict' -count=1

# Harness arms (independent of Notes)
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd experiments/ab-p25-gap-pass-validation
./score.sh G1 --p25 --arm build 2>&1 | tee /tmp/p27-s03-02-reverify-build.txt
./score.sh G1 --p25 --arm directed 2>&1 | tee /tmp/p27-s03-02-reverify-directed.txt
```

Confirm: thin build still shows **P25-3a FAIL** (expected) and G2/enforce **FAIL** (not WARN-only); P25-1/2 PASS.

## Review checklist

### Blockers

- [ ] Missing VERIFY-NOTES or evidence dir
- [ ] S03-01 `failed` without pending spawn
- [ ] Thin graph still passes product `--strict --enforce`
- [ ] `score.sh` still warn-only (no enforce FAIL on thin)
- [ ] `go test ./internal/...` FAIL
- [ ] Successor left **TBD** at close
- [ ] DR-HANDOFF CLOSED while verify bar failed

### High

- [ ] VERIFY-NOTES verdict PASS but spot-checks fail
- [ ] Arms missing P25-3a/3b labels
- [ ] Phase 28 successor chosen but scaffold missing
- [ ] AGENTS.md / TODO index not updated when phase complete

### Medium / low

- [ ] Evidence stale vs git SHA — note in REVIEW-NOTES
- [ ] S02-02 duplicate BLOCKING message not listed in residuals
- [ ] Missing RESULTS.md row

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Spawn `P27-S03-02a` implement + `02b` review immediately below |
| medium | Prefer spawn unless ≤5-line doc fix |
| low / nit | List in REVIEW-NOTES; do not block close |

## Evidence artifacts (reviewer output)

- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update [DR-HANDOFF.md](../../DR-HANDOFF.md) on APPROVE
- Update `docs/TODO.md` + orchestrator paste in `AGENTS.md`
- Register Phase 28 board only if successor promoted

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

## Exit criteria

- [ ] Independent re-verify minimum commands PASS
- [ ] S03-01 evidence reviewed (VERIFY-NOTES + archive + enforce upgrade)
- [ ] Closure signals confirmed; expected thin residuals explicit
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with explicit successor (`no successor` default — **never TBD**)
- [ ] Phase 28 scaffold **runnable** if that successor chosen
- [ ] `docs/TODO.md` (+ AGENTS.md) reflects phase complete / next runnable
- [ ] Confidence **high** (or **medium** with explicit residuals)
- [ ] All Phase 27 board rows `done`

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF without independent re-verify
- Rewriting Phase 26 or S01–S02 `done` history
- Implementing new product features in this review row
- Requiring Session-B dogfood to close when closure signals already PASS

## Minimal todos

- [ ] Read VERIFY-NOTES + evidence manifest
- [ ] Re-run locked re-verify commands
- [ ] Execute review checklist + successor table
- [ ] Close DR-HANDOFF with locked template
- [ ] Scaffold Phase 28 if promoted; else document `no successor`
- [ ] Update TODO index + AGENTS.md; write REVIEW-NOTES.md; set row `done`

## Next

**none** (if `no successor`) — or **P28-00** if Phase 28 promoted
