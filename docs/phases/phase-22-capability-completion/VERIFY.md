# Phase 22 VERIFY (index)

Runnable prompts live in the verify scope:

- Implement/verify: [`scopes/scope-08-agent-workflow-verify/07-verify.md`](scopes/scope-08-agent-workflow-verify/07-verify.md) (`P22-S08-07`)
- Independent review + DR-HANDOFF: [`scopes/scope-08-agent-workflow-verify/08-verify-review.md`](scopes/scope-08-agent-workflow-verify/08-verify-review.md) (`P22-S08-08`)

## Gate (locked)

1. Re-read [`docs/CAPABILITIES_CHECKLIST.md`](../../CAPABILITIES_CHECKLIST.md).
2. Every item is `[x]` with evidence in board Notes / VERIFY-NOTES, **or** a spawned **in-phase** `Na`/`Nb` pair still exists on [`docs/TODO/phase-22.md`](../../TODO/phase-22.md) and is runnable.
3. Fail the gate if any checklist `[ ]` would be “Phase 23 / later / residual”.
4. **S09 enhancements E01–E04** must have evidence (loop `harness_recommendations`, bundled catalog, CLI/MCP `trace agents`) — not post-MVP.
5. DR-HANDOFF successor is **`no successor`** only when the checklist is fully `[x]` and E01–E04 are evidenced.
6. Do not rewrite Phase 00–21 `done` history.

**Board order:** complete **S09** before running S08-07 VERIFY.

Details, command floor, and evidence paths are in the S08-07 prompt.
