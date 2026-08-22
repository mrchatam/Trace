# P28-S06-07 — FM-07 / FR-P28-04 implementer

## Metadata
- id: P28-S06-07
- todo_ids: [P28-S06-07]
- role: implementer
- skills: [documentation-and-adrs, writing-for-agents]
- verification: mixed
- hooks: []

## Objective

**FR-P28-04 / FM-07:** Keep git-sparsity / post-hoc SPEC commits as **warn-only** unless product explicitly adds plan-before-edit mode. Record an explicit decision: remain warn-only (document in VERIFY/harness) **or** ship plan-before-edit gate with tests. Do **not** silently turn FM-07 into hard fail.

## References

- [00-PLANNER.md](00-PLANNER.md) — default: remain warn-only
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-07
- `experiments/ab-p25-gap-pass-validation/score.sh` FM-07 warn block (~L135–157)
- [VERIFY-NOTES.md](../scope-05-verify/VERIFY-NOTES.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start. If upgrading to fail-closed plan-before-edit, clarify with human first (material product decision).

## Acceptance hint

Acceptance = explicit decision recorded: remain warn-only (document) **or** ship plan-before-edit gate with tests.

## Preflight

```bash
cd /home/ali/Desktop/Trace
sed -n '130,165p' experiments/ab-p25-gap-pass-validation/score.sh
grep -n 'FM-07\|git-spars\|warn' experiments/ab-p25-gap-pass-validation/PROTOCOL.md | head
```

## Suggested work (default path)

1. Document warn semantics in harness PROTOCOL + residual-wave note `FM07-DECISION.md`.
2. Cross-link VERIFY / agent harness so operators know warn ≠ fail.
3. Only if human approved: implement plan-before-edit fail-closed with tests — otherwise stay warn-only.

## Out of scope

- Silent hard-fail without decision doc; daemon/HTTP; other FMs

## Exit criteria

- [ ] `FM07-DECISION.md` (or equivalent) states chosen path
- [ ] Next runnable **P28-S06-08**

## Todo updates

Status + notes on **P28-S06-07** only.

## Next

`P28-S06-08`
