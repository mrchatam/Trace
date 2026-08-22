# P28-S02-02 — Session-B independent review

## Metadata
- id: P28-S02-02
- todo_ids: [P28-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent review of Session-B **evidence** (not a second gap pass). Confirm P25-3b was scored on the **directed** arm, Session-A was not wiped or re-scored as build, and R1 can close or must stay open with a spawn.

**Fresh subagent** — do not reuse P28-S02-01 session.

## References

- [01-run-and-score.md](01-run-and-score.md)
- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R1
- `SESSION-B-NOTES.md`, `SESSION-B-SCORE.txt`, `SESSION-A-GRAPH-SNAPSHOT.json` (this folder)
- [`experiments/RESULTS.md`](../../../../../experiments/RESULTS.md)
- [RUBRIC.md](../../../../../experiments/ab-p25-gap-pass-validation/RUBRIC.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked review policy

| Item | Rule |
|------|------|
| Re-run directed score | **Optional** confirm: `./score.sh G1 --p25 --arm directed` (no `--arm build`) |
| Rebuild G1 | **Forbidden** (`prepare.sh`) |
| Trace product diffs | Implementer must have **zero** `internal/` / `cmd/trace/` changes; if present → blocker spawn |
| M-16 | Must **not** be cited as P25-3b PASS |
| R1 close | Only if P25-3b **PASS** in `SESSION-B-SCORE.txt` **and** honesty G2 PASS (or documented enforce exit 0) |

## Verify checklist

### Deliverables

- [ ] `SESSION-B-NOTES.md` exists (score summary, counts, P25-4 attestation, arm isolation)
- [ ] `SESSION-B-SCORE.txt` exists; command line includes `--arm directed` (not default build)
- [ ] `SESSION-A-GRAPH-SNAPSHOT.json` exists; snapshot disc/dec match thin Session-A (or notes explain pre-existing richness)
- [ ] `experiments/RESULTS.md` has a **separate** `E02-SB` row (P27 E02 row still intact)
- [ ] Live `runs/G1` still exists (not wiped)

### P25 / honesty (from transcript + live graph)

- [ ] P25-1 PASS (mandatory gap pass text)
- [ ] P25-2 PASS (Parent orchestrator text)
- [ ] P25-3b PASS **or** FAIL with spawn (discoveries≥1 OR decisions≥1)
- [ ] G2 `--strict --enforce` not silently skipped
- [ ] Live `runs/G1/trace/graph.json` counts match notes
- [ ] Discoveries have `discovery-mentions-task`; decisions have `decision-task`

### Arm isolation / process

- [ ] No evidence of `./prepare.sh` during S02
- [ ] No post-mutation `--arm build` score claimed as Session-A
- [ ] Gap used `loop status` (not a fake `trace gap` success)
- [ ] Optional promotion documented if used; absence is OK
- [ ] Hook deny / honesty-dup / attestation harness **not** implemented here
- [ ] S01 `TEST-MATRIX.md` not rewritten

### Live spot-checks (reviewer runs)

```bash
REPO=/home/ali/Desktop/Trace
HARNESS="$REPO/experiments/ab-p25-gap-pass-validation"
SCOPE="$REPO/docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood"

test -f "$SCOPE/SESSION-B-NOTES.md"
test -f "$SCOPE/SESSION-B-SCORE.txt"
test -f "$SCOPE/SESSION-A-GRAPH-SNAPSHOT.json"
grep -q 'E02-SB' "$REPO/experiments/RESULTS.md"
grep -q -- '--arm directed' "$SCOPE/SESSION-B-SCORE.txt" || grep -q 'p25 arm: directed' "$SCOPE/SESSION-B-SCORE.txt"
grep 'P25-3b' "$SCOPE/SESSION-B-SCORE.txt"

# Optional re-score directed only — do not use default --arm
# cd "$HARNESS" && ./score.sh G1 --p25 --arm directed
```

If S02-01 is `blocked` with `SESSION-B-BLOCKED.md`: review the block reason; do **not** invent a P25-3b PASS; spawn restore+rerun pair if G1 can be reconstructed without pretending Session-A still exists.

## Findings severity

| Level | Action |
|-------|--------|
| blocker | No notes/score; `prepare.sh` wipe; `--arm build` after mutation presented as Session-A; Trace product code changed; P25-3b PASS claimed without disc/dec evidence |
| high | Missing honesty links with enforce ignored; RESULTS overwrote P27 E02 row; directed score missing `--arm directed` |
| medium | P25-4 attestation missing; snapshot missing but Session-A otherwise intact |
| low / nit | Wording; extra G1 product churn unrelated to gap |

## Spawn policy

Insert **immediately below** this row (`P28-S02-02a` implement / `02b` review) if blocker/high remains:

- Re-run directed gap + score (no prepare) **or**
- Restore protocol (human) if G1 destroyed

P25-3b **FAIL** after a real directed session: document; **S03/S04 may still proceed**; keep R1 open in review notes (S05 VERIFY re-reads). Do not silently close R1.

## Verdict

| Outcome | Next |
|---------|------|
| P25-3b PASS + isolation OK | R1 closed in review Notes; proceed **P28-S03-00** |
| P25-3b FAIL (session ran) | R1 open; S03 still next unless spawn |
| Blocked / missing G1 | spawn or leave S02 failed; do not start S03 until orchestrator decides |

Confidence must be **medium or high** with evidence listed. Empty APPROVE without files is invalid.

## Todo updates

Status + notes on **P28-S02-02**; spawn rows allowed. Do not rewrite S02-01 prompt if `done`.

## Exit criteria

- [ ] Checklist completed against live files
- [ ] Findings listed (or explicitly none)
- [ ] R1 disposition stated (closed / open)
- [ ] Board Notes cite score path + P25-3b line
- [ ] Next runnable named (`P28-S03-00` or spawn)

## Next

**P28-S03-00** (unless spawn)
