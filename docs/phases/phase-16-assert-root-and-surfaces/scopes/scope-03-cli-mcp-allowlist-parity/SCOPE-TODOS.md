# S03 — CLI vs MCP allowlist parity — scope todos

**Depends-on:** P16-S02-02 APPROVE. Owns **DF-77**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** (2026-08-17) |
| 2 | 01-cli-mcp-allowlist-parity | implement | **done** (P16-S03-01) |
| 3 | 02-scope-review | review | **APPROVE high** (P16-S03-02); next **P16-S04-00** |

## Depends (from S02 — live P16-S02-02 APPROVE)
- S02 CHECK + YOLO heal→PENDING is the shared enum (`014_*`); S03 does **not** add a second decision table or mig **015**
- **`cli:` prefix reserved:** S02 canonicalize does **not** rewrite `cli:*` to `mcp:` (`TestCanonicalizeCustomAndCLISlugsUnchanged`). Dual-slug stays independent
- Unprefixed MCP Names persist as `mcp:`+Name — CLI slugs are `cli:<command>` only; do **not** treat `add` / `trace_why` as CLI slugs
- MCP helper still `mcp:`+Name; nine tools
- Resolve must AUTO_ALLOW `BuiltinCLICapabilitySpecs()` or first `trace add` PENDING-fails

## Phase locks (P16-S03-00 FINAL) — dual-slug
- MCP `mcp:<tool>` and CLI `cli:<command>` are **independent**
- MCP DENIED does **not** deny CLI (CLI-first); CLI DENIED does **not** deny MCP
- Operator Shell lockdown = `decide --slug cli:add DENIED`
- **Gated:** add, link, transition, review, why, context, tasks, seed, impact, plan, index (`reindex`→`cli:index`)
- **Ungated:** help, version, init, **whole** `capability`, migrate, backup, restore, auth, install
- G19: `assertCLICommand` → `AssertToolAllowed("cli:"+cmd)`; domain owns builtin list + `cli:reindex` fold
- No TTY skip; no new MCP tools; no shared slug

## Named tests (S06 import)
`TestCLIAddSucceedsWhenMCPAddDenied`, `TestCLIAddDeniedFailClosed`, `TestCLIWhySucceedsWhenMCPWhyDenied`, `TestCLIWhyDeniedFailClosed`, `TestCLIAddDeniedDoesNotBlockMCPAdd`, `TestUngatedCapabilityDecideWhenCLIAddDenied`, `TestUnprefixedAddDecideDoesNotGateCLI`, `TestCapabilityDecisionAutoAllowBuiltinCLI`, `TestCanonicalizeCLIReindexFoldsToIndex`, `TestCLIIndexAliasDenied` + MCP/S01/S02 keepers

## Reminders
- Next after APPROVE: **P16-S04-00** (`install` stays ungated — live-confirmed P16-S03-02)
