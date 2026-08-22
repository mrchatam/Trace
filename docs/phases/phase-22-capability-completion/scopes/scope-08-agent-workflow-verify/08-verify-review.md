# P22-S08-08 — Review: phase verify + close DR-HANDOFF

## Metadata
- id: P22-S08-08
- todo_ids: [P22-S08-08]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Independent **fresh-session** review of S08-07 evidence; re-run the locked floor (do **not** trust Notes alone); **close DR-HANDOFF** only if C01–C43 and checklist are fully `[x]` (hard-boundary outs already non-goals).

## Session start

Follow [agent-loop-protocol.md](../../../../../rules/agent-loop-protocol.md). **Fresh session** — not the S08-07 agent. **Unattended:** execute until APPROVE or spawn forward.

## References

- [07-verify.md](07-verify.md)
- [VERIFY.md](../../VERIFY.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [README.md](../../README.md) coverage matrix
- [`docs/CAPABILITIES_CHECKLIST.md`](../../../../CAPABILITIES_CHECKLIST.md)
- [00-PLANNER.md](00-PLANNER.md)

## Locked DR-HANDOFF close policy

| Field | Locked value |
|-------|--------------|
| Who closes | **S08-08 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Successor | **`no successor`** — **only if** all 43 matrix rows + checklist fully `[x]` |
| Must not | Successor `Phase 23` for leftover capabilities; `TBD`; rewrite P21 historical `no successor` |
| If FAIL | Keep DR-HANDOFF **OPEN**; spawn `P22-S08-08a` implement + `08b` review immediately below |

### Residuals allowed at close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| Hosted MCP / daemon / HTTP | D-22-14 **out** |
| ML / embeddings / graph DB | D-22-11 / Law 13 **out** |
| Wrap `git commit` | D-22-16 **out** |
| Code graph omitted from seed | D-22-06 **keep** — clones `index` |
| DONE/Review PASS policy | D-22-20 **keep** |

**Forbidden residual:** any C01–C43 matrix row or checklist bullet still `[ ]`.

### DR-HANDOFF.md update template (on APPROVE)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 22 outcome | All 43 capabilities `[x]`; schema through **027**; local git-hook; CLI+stdio MCP query/loop/agents |
| Residuals (non-blocking) | hosted MCP out; ML out; code graph local+index |
| Forward (human queue) | none required for checklist completeness |
```

## Review checklist

1. **Independently re-run** S08-07 command floor (full copy from 07-verify.md) — session ≠ S08-07.
2. Walk `CAPABILITIES_CHECKLIST.md` — **zero** `[ ]` on APPROVE.
3. Walk README **C01–C43** — matches VERIFY-NOTES; spot-check S08 (C28/C38/C39) and S09 (E01–E04) evidence.
4. Schema **27** sql files; compat **27**; **no 028+**.
5. Confirm S01–S08 reviewers spawned rather than deferring (“later”, “Phase 23”).
6. Confirm `./bin/trace-mcp -h` lists **15** tools post-S09.
7. If APPROVE: update `AGENTS.md` + `docs/TODO.md` Current focus → Phase 22 **complete**; next runnable **none** on phase-22 board.

## Spawn policy

If any capability still `[ ]` or VERIFY-NOTES FAIL without completed spawn chain: spawn **`P22-S08-08a` + `P22-S08-08b`** (or specific Na/Nb) **on this board**. **Do not** mark this review `done` while residuals remain.

## Re-run commands

Same as [07-verify.md](07-verify.md) command floor +:

```bash
grep -c '^\- \[ \]' docs/CAPABILITIES_CHECKLIST.md  # expect 0 on APPROVE
./bin/trace-mcp -h | wc -l  # sanity: lists trace_agents
```

## Exit criteria

- [ ] Independent re-verify **PASS**
- [ ] C01–C43 all `[x]`; checklist **141/141 `[x]`**
- [ ] E01–E04 PASS in VERIFY-NOTES
- [ ] DR-HANDOFF **CLOSED** with **`no successor`** only on full pass
- [ ] Confidence **high**
- [ ] Board Notes: evidence + successor line + index focus update

## Deliverables on APPROVE

- [ ] `DR-HANDOFF.md` updated (template above)
- [ ] `docs/TODO/phase-22.md` — phase row complete if all board rows done
- [ ] `docs/TODO.md` — next runnable **none** (or forward queue per human)
- [ ] Optional: `REVIEW-NOTES.md` in this folder
