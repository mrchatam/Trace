# P11 / S03 / 01 — Store lock / concurrency

## Metadata
- id: P11-S03-01
- todo_ids: [P11-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-47** per sibling **00-PLANNER** FINAL locks (2026-08-16). Keep exclusive `.trace/trace.lock`; add short bounded Open retry + clearer ErrLocked / help (serialize CLI↔MCP or use worktrees). **No new migration. Do not drop exclusivity.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-47
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md); [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md)
- [phase README](../../README.md)
- Live: `internal/store/{lock,open}.go`; `cmd/trace/help.go`; `internal/mcp/project.go`; `evals/compat`
- Prior: P08 S02 path-local + exclusive lock — keep green
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/store` (+ thin `cmd/trace/help` + optional `internal/mcp` copy) |
| Migration | **None** |
| Exclusivity | Retain `LOCK_EX` for Open→Close; post-budget contention → `ErrLocked` |
| Retry | Bounded ~250–500ms total in store acquire on EWOULDBLOCK/EAGAIN; optional `TRACE_LOCK_WAIT_MS` |
| UX | ErrLocked (+ help/MCP) guide: serialize CLI↔MCP **or** separate `-C`/worktrees; exit **2** unchanged |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P11 DF-40/43/44 |
| Forbidden | Drop flock; multi-writer; indefinite wait; daemon; mig; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/store/lock.go` | Bounded retry loop; thicken `ErrLocked` text (keep `errors.Is`); optional env wait override |
| `internal/store/bind_lock_test.go` (or new `*_test.go`) | `TestOpenRetrySucceedsWhenLockReleasedSoon`; keep `TestConcurrentStoreOpenFailClosed` (ensure wait budget expires under held lock) |
| `cmd/trace/help.go` | Global: serialize CLI↔MCP / one-writer + worktrees |
| `internal/mcp/project.go` / tool descs (optional) | Surface same serialize guidance if errors already wrap `ErrLocked` thinly |
| `cmd/trace/cli_test.go` / MCP tests (optional) | Lock message / exitFail=2 still; init-while-locked stays fail-closed |
| `evals/compat` | Must stay green — exclusivity checklist unchanged in meaning |

## Role work

1. TDD: brief release-during-wait → second Open OK; held lock past budget → `ErrLocked`.
2. Implement retry + clearer sentinel/message in `internal/store` only (G19).
3. Thicken help (+ optional MCP) serialize / worktree guidance.
4. Run locked verify suite; board **status + Notes only** (cite test names + DF-47).

## Algorithm sketch (non-normative — locks above win)

```text
acquireTraceLock:
  open lock file
  deadline = now + waitBudget (~250–500ms)
  loop:
    Flock(LOCK_EX|LOCK_NB)
    if ok → return handle
    if not EWOULDBLOCK/EAGAIN → error
    if now >= deadline → return ErrLocked (actionable text)
    sleep small step
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: retry-on-brief-release; held-lock still ErrLocked; help/ErrLocked mention serialize; exitFail=2; compat `trace_lock_ok`.

## Exit criteria
- [ ] `TestOpenRetrySucceedsWhenLockReleasedSoon` (or equiv) green — DF-47 soft race
- [ ] `TestConcurrentStoreOpenFailClosed` (or equiv) still green — exclusivity
- [ ] ErrLocked and/or help/MCP assert serialize / CLI↔MCP / worktree guidance
- [ ] CLI lock conflict remains exit **2**; no mig; no dropped flock
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S03-02**

## Out of scope
- S04+ Phase 11 product work (DF-41, DF-51, …)
- True multi-writer / shared lock across processes
- Rewriting Mode-B Gate C packs / Phase 00–10 history

## Todo updates
Implementer: **status + notes only**. Record test names + DF-47 evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] Store retry + ErrLocked UX tests + implementation
- [ ] Help / optional MCP serialize guidance
- [ ] Locked verify suite + board Notes
