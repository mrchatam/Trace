# P11-S03-00 — Store lock / concurrency (FINAL)

## Metadata
- id: P11-S03-00
- todo_ids: [P11-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S03 implement/review prompts for **DF-47**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A4 DF-47 posture
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-47
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md) — 3/4 parallel `context`; MCP↔CLI race
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md) — `exits 2 0 2 2`
- Phase 08 S02 FINAL: exclusive `trace.lock` / `ErrLocked` / path-local `.trace`
- Live: `internal/store/{lock,open}.go`; `cmd/trace/help.go`; `internal/mcp/project.go`; `evals/compat` lock checklist
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (phase A4 + live inventory + P08 exclusivity do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `acquireTraceLock` | `unix.Flock(LOCK_EX\|LOCK_NB)` — **immediate** fail → `ErrLocked` (no retry) |
| `ErrLocked` text | Mentions worktrees / separate `-C`; **does not** say serialize CLI↔MCP or that single-writer is intentional |
| MCP | Per-tool `openStore` + `defer Close` — races CLI when overlapping; not a long-held daemon lock |
| CLI | `ErrLocked` → `exitFail` (**2**); help Global already says one open store / parallel agents use worktrees |
| Tests | `TestConcurrentStoreOpenFailClosed` + compat `trace_lock_ok` — exclusivity green; **no** brief-contention retry test |
| Dogfood | 4× parallel `trace context` → `exits 2 0 2 2`; MCP `trace_tasks` vs CLI same `ErrLocked` |
| Migration | None — lock file is advisory flock, not schema |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-47 | Concurrent CLI↔CLI / CLI↔MCP on same root is brittle / opaque | Keep exclusive lock; add **short bounded retry** + **clearer ErrLocked / docs** so agents know to serialize or use worktrees — **not** multi-writer |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-47** only (P08 path-local bind + exclusivity stay as shipped — do not regress) |
| Packages | **`internal/store`** (`lock.go` / Open path) for retry + `ErrLocked` text; thin **`cmd/trace/help.go`** (+ optional usage stderr already via wrapped err); thin **`internal/mcp`** copy in tool/server descriptions if needed. **G19** — no domain fork; adapters do not reimplement lock policy |
| Migration | **None** |
| Exclusivity | **Retain** exclusive `.trace/trace.lock` for the life of `Open`→`Close`. Second holder after retry budget exhausted → `errors.Is(..., ErrLocked)`. Compat checklist `trace_lock_ok` / `TestConcurrentStoreOpenFailClosed` stay green |
| Retry | Inside store acquire (not CLI-only): on `EWOULDBLOCK`/`EAGAIN`, **short bounded retry** with small sleeps before returning `ErrLocked`. Default total wait budget **≈250–500ms** (pick one constant; document in Notes). Optional `TRACE_LOCK_WAIT_MS` override OK if cheap; default must work with no env. **Forbidden:** indefinite block; multi-second default that looks hung |
| Fail-closed bar | Holder that never releases during the wait window must still yield `ErrLocked` (existing concurrent test stays valid — extend timing if needed so test exceeds budget) |
| Soft-race bar | If first Open releases within the budget, second Open **succeeds** without caller retry loops |
| `ErrLocked` UX | Assertable message must keep `errors.Is` identity (sentinel or wrap with `%w`). Text must guide agents: **serialize** CLI↔MCP / parallel Trace on one root **or** use separate `-C` / worktree roots; single-writer is intentional. Prefer phrases: `serialize`, `CLI`/`MCP`, and existing worktree/`-C` guidance |
| Exit code | CLI remains **`exitFail` (2)** on lock conflict — do not invent a new exit for DF-47 |
| Docs | Thicken `help` Global (and optional MCP open-store / tool desc) with serialize + worktree guidance. Do **not** claim concurrent same-root writers are supported |
| Non-goals / forbidden | Drop exclusive flock; multi-writer SQLite; shared parent lock; daemon / long-lived MCP store across tools; silent swallow of lock errors; full-rebuild indexer; new mig; rewriting Phase 00–10 / P11-S01–S02 `done` history |
| Tests (required) | (1) Keep **`TestConcurrentStoreOpenFailClosed`** (held lock → `ErrLocked`). (2) **`TestOpenRetrySucceedsWhenLockReleasedSoon`** (or equiv): second Open waits briefly while first Close → succeeds. (3) Assert `ErrLocked` (and/or help/MCP) mentions serialize and/or CLI↔MCP / worktree guidance. (4) Keep compat exclusivity / `TestInitFailClosedWhenStoreLocked` green |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P11-S01 DF-40; P11-S02 DF-43/44 |
| Residual OK | Sustained parallel N writers on one root still fail-closed after budget (dogfood `2 0 2 2` shape may persist under true overlap) — success criterion is **actionable UX + brief-race recovery**, not unlimited concurrency |

## Effects on later scopes
- **S04** (capability/hatch): no store-lock coupling — serial after S03 review only. Light Depends note on S04 stubs.
- **S08 VERIFY:** include DF-47 retry-on-brief-contention + exclusivity-still-fail-closed + ErrLocked/docs serialize guidance in evidence table.

## Exit
- [x] Thicken `01-store-lock-concurrency.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S04)
- [x] Board Notes; next **P11-S03-01**
- [x] Product Go — **not** this row
