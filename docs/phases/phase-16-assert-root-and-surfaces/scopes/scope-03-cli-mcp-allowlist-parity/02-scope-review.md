# P16 / S03 / 02 — Scope review (CLI vs MCP allowlist) FINAL checklist

## Metadata
- id: P16-S03-02
- todo_ids: [P16-S03-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S03 **DF-77**. Fresh subagent ≠ implementer. Compare claims + **00-PLANNER FINAL** + `01` to live code/tests. Spawn `P16-S03-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S04-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-cli-mcp-allowlist-parity.md](01-cli-mcp-allowlist-parity.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `cmd/trace/root.go`, gated cmd files, `internal/domain/capability_decision.go`, `internal/mcp/assert.go`
- Hunt: `experiments/_bughunt/post-p15/` `mcp_cli_add.*` / `cap_cli_why.*` (historical; tests are SoT)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open dual-slug isolation, write-set trim, unprefixed `add`→`cli:add`, or gating `capability decide` (FINAL).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Dual-slug: `mcp:trace_add` DENIED does **not** block `trace add` | `TestCLIAddSucceedsWhenMCPAddDenied` |
| 2 | `cli:add` DENIED fail-closes CLI add (no entity) | `TestCLIAddDeniedFailClosed` |
| 3 | Reverse isolation: `cli:add` DENIED does not block MCP `trace_add` | `TestCLIAddDeniedDoesNotBlockMCPAdd` |
| 4 | `why` gated; MCP why DENIED still allows CLI why | `TestCLIWhyDeniedFailClosed`; `TestCLIWhySucceedsWhenMCPWhyDenied` |
| 5 | Operator escape: `capability decide` works when `cli:add` DENIED; `init`/`install`/`help`/`version` never Assert | `TestUngatedCapabilityDecideWhenCLIAddDenied`; grep ungated files for `AssertToolAllowed`/`assertCLICommand` |
| 6 | G19: CLI calls `domain.AssertToolAllowed("cli:"+cmd)`; `cli:reindex`→`cli:index`; unprefixed `add` is not a CLI slug | Helper grep; `TestCanonicalizeCLIReindexFoldsToIndex`; `TestCLIIndexAliasDenied`; `TestUnprefixedAddDecideDoesNotGateCLI` |
| 7 | Builtin CLI AUTO_ALLOWED; MCP specs still nine; helper `mcp:`+Name unchanged | `TestCapabilityDecisionAutoAllowBuiltinCLI`; `TestToolNamesRegistered`; grep `assertMCPToolAllowed` |
| 8 | No shared slug / TTY skip / YOLO flags / mig 015 / new MCP tools / Assert on install | Diff + Notes; compat still **14** |
| 9 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01 virgin + S02 CHECK keepers | Re-run `01` locked verify cmds |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P16-S03-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P16-S04-00** unless spawn.
5. If APPROVE: S06 must import named isolation/DENIED/alias/ungated-decide tests (already pointed on S06 SCOPE-TODOS).

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals listed)
- [x] REVIEW-NOTES.md written
- [x] Board status + Notes; next **P16-S04-00** (unless spawn)
- [x] No rewrite of done P16-S03-00/01 history

## Minimal todos
- [x] Independent verify + checklist
- [x] REVIEW-NOTES.md
- [x] Board sync
