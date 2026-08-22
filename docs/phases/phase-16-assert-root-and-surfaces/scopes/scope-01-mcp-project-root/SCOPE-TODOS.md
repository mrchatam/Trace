# S01 — MCP project root / auto-init — scope todos

**Depends-on:** P16-00 FINAL. Owns **DF-76** (high).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** FINAL 2026-08-17 |
| 2 | 01-mcp-project-root | implement | **done** 2026-08-17 |
| 3 | 02-scope-review | review | **done** APPROVE high 2026-08-17 — next **P16-S02-00** |

## Phase locks (P16-S01-00 FINAL)
- API: `store.OpenExisting` + sentinel `ErrNotInitialized` — exists = regular file `<abs>/.trace/trace.db`
- MCP `openStore` → `OpenExisting` (not `Open`); missing store → CallTool error; **no** mkdir / AUTO_ALLOWED DB
- CLI `store.Open` mkdir for `trace init` **stays**
- bind-to-defaultRoot **rejected** (fail-closed, no silent fallback)
- Per-store SoT HOLD (DENIED on initialized A does not apply to initialized B)
- Assert slug `mcp:`+Name **unchanged**; no new tools / mig / daemon

## Named tests (S01-01 must land; S01-02 + S06 import)
- `TestMCPVirginProjectDoesNotMkdir`
- `TestMCPInitializedOtherRootIsolated`
- Keepers: `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered`
- Store: `TestOpenExistingMissingReturnsErrNotInitialized`, `TestOpenExistingEmptyTraceDir`, keeper `TestOpenCreatesDBAndMigratesIdempotent`
- CLI keeper: `TestInitCreatesDB`

## Depends (to S02 — light)
S02 CHECK/slug does not require S01 product, but **board order is sequential**. After S01 APPROVE: **P16-S02-00**.

## Depends (to S06 — light)
VERIFY imports S01 named virgin + isolation tests + P15 Assert keepers (already on S06 SCOPE-TODOS).

## Reminders
- G19; no new MCP tools; no daemon; no session-global DENY
- Forward-only; implementers: status + Notes only
