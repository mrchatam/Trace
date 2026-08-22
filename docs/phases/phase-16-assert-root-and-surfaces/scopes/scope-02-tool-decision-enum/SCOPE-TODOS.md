# S02 — Tool-decision enum + slug prefix — scope todos

**Depends-on:** P16-S01-02 APPROVE (board). Owns **DF-75**, **DF-78**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** FINAL 2026-08-17 |
| 2 | 01-tool-decision-enum | implement | **done** 2026-08-17 |
| 3 | 02-scope-review | review | **done** APPROVE high 2026-08-17 — no spawn |

## Phase locks (P16-S02-00 FINAL)

- Mig **`014_capability_tool_decision_enum.sql`**: rebuild + CHECK `decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED')`; restore decision index; do not rewrite 001–013
- YOLO/garbage **write:** reject (CHECK + domain + store Upsert)
- YOLO/garbage **migrate:** heal → **PENDING** (never AUTO_ALLOWED, never DROP)
- Resolve unknown persisted status → **fail-closed** (no builtin AUTO_ALLOWED upsert)
- Canonicalize Decide+Resolve: exact registered MCP **Name** or `mcp:`+Name → `mcp:`+Name (`BuiltinMCPCapabilitySpecs`); no globs
- Custom slugs unchanged; **`cli:` reserved** (not rewritten, not Asserted here)
- Footgun fold: unprefixed DENIED wins over prefixed AUTO_ALLOWED
- MCP `assertMCPToolAllowed` **unchanged**; prefer zero CLI edits
- Compat ceiling **14** (no 015+)

## Named tests (S02-01 must land; S02-02 + S06 import)

- `TestCapabilityToolDecisionCheckRejectsYOLO`
- `TestCapabilityToolDecisionMigrateHealsYOLOToPending`
- `TestResolveYOLOBuiltinDoesNotAutoAllow`
- `TestDecideUnprefixedMCPNameCanonicalizes`
- `TestCanonicalizeCustomAndCLISlugsUnchanged`
- `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed`
- `TestMCPUnprefixedDecideGatesCallTool`
- Keepers: `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered`, `TestCapabilityDecisionAutoAllowBuiltinMCP`, `TestCapabilityDecisionUnknownPendingFailClosed`, `TestOpenCreatesDBAndMigratesIdempotent` (v14), `TestMCPVirginProjectDoesNotMkdir`

## Depends (to S03 — light)

S03 CLI `cli:` Assert needs this fail-closed enum. **`cli:` prefix is reserved:** S02 canonicalize must not map `cli:add` (or any `cli:`) to `mcp:`. Do **not** implement CLI Assert in S02. MCP DENIED still must not imply CLI DENIED (S03 dual-slug).

## Depends (to S06 — light)

VERIFY imports S02 named CHECK/heal/canonicalize + unprefixed-decide tests + compat **14**. Do not claim DF-77 fixed.

## Reminders

- Next after APPROVE: **P16-S03-00**
- G19; no new MCP tools; no daemon; no YOLO/AllowAll
- Forward-only; implementers: status + Notes only
