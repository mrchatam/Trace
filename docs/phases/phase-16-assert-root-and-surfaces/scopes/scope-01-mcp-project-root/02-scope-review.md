# P16 / S01 / 02 — Scope review (MCP project root) FINAL checklist

## Metadata
- id: P16-S01-02
- todo_ids: [P16-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S01 DF-76 (MCP no auto-init). Fresh subagent ≠ implementer. Compare claims + **00-PLANNER FINAL** + `01` to live code/tests. Spawn `P16-S01-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S02-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-mcp-project-root.md](01-mcp-project-root.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `internal/mcp/project.go`, `internal/store/open.go`, named tests in `internal/mcp` + `internal/store`
- Hunt: `experiments/_bughunt/post-p15/{mcp-deny,mcp-noinit,mcp-fresh}/` (historical; tests are SoT)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open bind-to-defaultRoot vs OpenExisting (FINAL).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Virgin `project=` / missing `.trace/trace.db` → CallTool error; **no** auto-mkdir AUTO_ALLOWED DB | `TestMCPVirginProjectDoesNotMkdir`; grep `openStore` uses `OpenExisting` not `Open` |
| 2 | Empty `.trace/` without db also fail-closed (no `trace.db` created) | `TestOpenExistingEmptyTraceDir` + virgin MCP subcase |
| 3 | Bound-root DENIED still fail-closed | `TestMCPAssertDeniedBlocksCallTool` |
| 4 | Initialized other root stays isolated (no session-global DENY) | `TestMCPInitializedOtherRootIsolated` |
| 5 | CLI `store.Open` / `trace init` still creates `.trace/` | `TestOpenCreatesDBAndMigratesIdempotent`; `TestInitCreatesDB` |
| 6 | Sentinel `ErrNotInitialized`; MCP wraps `%w`; `errors.Is` in tests | Code + store named tests |
| 7 | No bind-to-defaultRoot; no new MCP tools; G19; no daemon; Assert slug unchanged | Diff + `TestToolNamesRegistered` |
| 8 | Carry-forward honesty/E–H/ablation/compat **13**/p0x/x0 + product pkgs | Re-run `01` locked verify cmds |
| 9 | R2/R3/R4 / DF-75/77/78 not falsely claimed fixed | Diff + Notes |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P16-S01-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P16-S02-00** unless spawn.
5. If APPROVE: S06 must import named virgin + isolation tests (already stubbed on S06 SCOPE-TODOS) — thicken S02 only if a gap is found (S02 does not need S01 product).

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P16-S02-00** (unless spawn)
- [ ] No rewrite of done P16-S01-00/01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
