# Phase 44 / S01 / Planner — Design locks

## Metadata
- id: P44-S01-00
- todo_ids: [P44-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, domain-modeling, api-and-interface-design]
- mcps: [user-trace]
- verification: automated

## Objective

Close D1–D5 formally in [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) (status → **LOCKED**). Thicken upcoming S02–S04 implement prompts from S00 RESEARCH.md evidence. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — P44-00 light locks
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) (must exist / APPROVE)
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol. Block if S00 RESEARCH not APPROVED.

## Locked defaults (inherit P44-00 unless RESEARCH forces spawn)

| Item | Value |
|------|-------|
| D1 | A + thin scope records; reject C primary |
| D2 | MVP: `scope_member`, `api_contract`, `implements`, `blocks` |
| D3 | `scope` query param + bounded N-hop; no tiles |
| D4 | Inference CLI only |
| D5 | GUI 500 stays |
| Reversal | Requires REVIEW spawn + reason — do not silently flip |

## Role work

1. Re-read RESEARCH.md gaps.
2. Thicken `01-lock.md` with exact doc edits + acceptance checklist.
3. Thicken `02-review.md` with independent lock rubric.
4. Optionally pre-thicken S02–S04 `01-implement.md` touch-lists (upcoming-only rights).

## Exit criteria

- [ ] `01-lock.md` runnable; D1–D5 close plan explicit
- [ ] Board Notes; next **P44-S01-01**

## Minimal todos

- [ ] Confirm light locks vs RESEARCH
- [ ] Thicken 01 + 02
- [ ] Board done

## Next

`P44-S01-01`
