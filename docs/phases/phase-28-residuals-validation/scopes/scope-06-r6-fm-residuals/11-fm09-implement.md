# P28-S06-11 — FM-09 / FR-P28-06 implementer

## Metadata
- id: P28-S06-11
- todo_ids: [P28-S06-11]
- role: implementer
- skills: [test-driven-development, writing-for-agents]
- verification: mixed
- hooks: []

## Objective

**FR-P28-06 / FM-09:** Prove mode collapse stays closed beyond single Session-B (build ≠ directed richness). Protocol/dogfood: repeat dual-lane score (**no** `prepare.sh` wipe); optional second directed fixture.

## References

- [00-PLANNER.md](00-PLANNER.md)
- S02 Session-B artifacts: `../scope-02-session-b-dogfood/`
- S05 dual-lane policy: `../scope-05-verify/01-verify.md`, `VERIFY-NOTES.md`
- `experiments/ab-p25-gap-pass-validation/score.sh` `--arm build|directed`
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start.

## Acceptance hint

Dual-lane: thin build baseline documented; directed P25-3b PASS; rich build labeled post-directed — not conflated with Session-A thin FAIL.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json
test -f docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-B-SCORE.txt
grep -n 'arm\|P25-3' experiments/ab-p25-gap-pass-validation/score.sh | head
# Do NOT run prepare.sh
```

## Suggested work

1. Re-run or document dual-lane scoring without prepare wipe.
2. Ensure labels separate thin baseline vs post-directed rich build.
3. Optional second directed fixture if G1 alone is insufficient.
4. Deliver `FM09-NOTES.md` + score snippets under this scope or `experiments/runs/…`.

## Out of scope

- `prepare.sh` wipe of G1; claiming rich build as Session-A thin; other FMs

## Exit criteria

- [ ] Acceptance hint met with labeled evidence
- [ ] Next runnable **P28-S06-12**

## Todo updates

Status + notes on **P28-S06-11** only.

## Next

`P28-S06-12`
