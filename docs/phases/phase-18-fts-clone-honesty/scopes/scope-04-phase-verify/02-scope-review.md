# P18 / S04 / 02 — VERIFY review (product evidence)

## Metadata
- id: P18-S04-02
- todo_ids: [P18-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent re-verify vs `VERIFY-NOTES.md`. Confirm named **DF-87 / DF-88 / DF-89** + keepers + carry-forward. **Confirm** DR-HANDOFF Notes started (`no successor`). Do **not** close DR-HANDOFF. Do **not** mark Phase 18 complete. Next **P18-S05-00**. Two-clone **not required** — reject treating a missing two-clone shell or dedicated `-run` as a VERIFY fail. DF-88 re-prove is **document-only**. Binary rebuild is **S05**, not a VERIFY fail. Do **not** rewrite P17 `done` history.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-verify.md](01-verify.md); [VERIFY-NOTES.md](VERIFY-NOTES.md) (after S04-01)
- S01–S03 REVIEW-NOTES (all APPROVE high)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-88-DECISION.md](../../DF-88-DECISION.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Fresh subagent ≠ S04-01. Unattended: no Plan-mode switch. Re-run locked commands from `00-PLANNER.md`. Stale binaries are **S05**, not a VERIFY fail.

## Review focus
Reject any claim that Phase 18 is complete after S04. Reject closing DR-HANDOFF on this row. Reject requiring the P17 two-clone shell or a dedicated `TestPortableGraphTwoCloneWhyContextPlan` `-run`. Reject reversing DF-88 exclude. Reject treating DF-88 as a clone-dir hunt. Reject treating Phase 01 dry-run as Gate C/F/G/ablation/H/checklist.

## Checklist

| # | Check |
|---|--------|
| 1 | Named DF-87 four tests + keepers; CGO0 **then** CGO1; both slash names |
| 2 | Named DF-88 `TestHelpCloneTasksImportPending` + omit/path keepers — **document-only** |
| 3 | Named DF-89 `TestIndexFileGoHandlerMethods` + keeper `TestIndexFileGoGolden` (CGO1) |
| 4 | P17 seed keepers round-trip + omit + `exported_at_commit` — **not** two-clone |
| 5 | Two-clone **not** required; no dedicated `-run`; shell not a fail bar |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0/product pkgs green |
| 7 | Gate C `dry_run:false` intact; dry-run ≠ C/F/G/ablation/H/checklist |
| 8 | DF-86 / DF-67 / harness / stale binaries **non-fail** |
| 9 | DR-HANDOFF Notes **started** `no successor` — **not closed** |
| 10 | Next **P18-S05-00**; S05 rows still pending; Phase 18 **not** complete |
| 11 | P17 historical `no successor` intact; hosted MCP not boarded |
| 12 | REVIEW-NOTES.md; no product Go invented on VERIFY |

## Locked expectations

| Item | Value |
|------|-------|
| Named path | Per `00-PLANNER.md` FINAL table |
| Two-clone | **Not required** |
| DF-88 | Document-only; exclude stays |
| DR-HANDOFF | Started only; close = **P18-S05-02** |
| Next | **P18-S05-00** unless spawn |
| Product Go | **Forbidden.** Spawn on fail |

## Role work
1. Re-run locked named DF-87/88/89 + keepers + carry-forward (same CGO matrix).
2. Diff vs VERIFY-NOTES. Spot-check DF-88: `SeedTask` still no `work_state`; no include flags.
3. Confirm DR-HANDOFF.md still OPEN; VERIFY-run started; successor TBD until S05-02.
4. Write `REVIEW-NOTES.md`.
5. Board Notes → **P18-S05-00** unless spawn.

## Exit criteria
- [ ] APPROVE high (or medium with residuals listed)
- [ ] Named DF-87/88/89 + keepers independently green
- [ ] Two-clone not required / not used as fail bar
- [ ] DR-HANDOFF Notes started/confirmed — **not** closed
- [ ] Board Notes; next **P18-S05-00** unless spawn
- [ ] Phase 18 **not** marked complete

## Minimal todos
- [ ] Re-verify named DF-87/88/89 + keepers vs VERIFY-NOTES
- [ ] Confirm two-clone not required; DF-88 document-only; S05 still next
- [ ] REVIEW-NOTES.md
- [ ] Board sync → **P18-S05-00**

## Out of scope
- Closing DR-HANDOFF / Phase 18 complete
- Starting S05 rebuild on this row
- Requiring two-clone shell or dedicated two-clone `-run`
- Rewriting P17 `done` history
- Hosted MCP / research S05 / plan simulate / D21+
