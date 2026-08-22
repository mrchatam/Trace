# P16 / S02 / 02 — Scope review (tool-decision enum) FINAL checklist

## Metadata
- id: P16-S02-02
- todo_ids: [P16-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S02 **DF-75** + **DF-78**. Fresh subagent ≠ implementer. Compare claims + **00-PLANNER FINAL** + `01` to live code/tests. Spawn `P16-S02-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S03-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-tool-decision-enum.md](01-tool-decision-enum.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `internal/store/schema/014_capability_tool_decision_enum.sql`, `internal/domain/capability_decision.go`, `internal/mcp/assert.go`
- Hunt: `experiments/_bughunt/post-p15/{cap-decisions,mcp-yolo,mcp-footgun}/` (historical; tests are SoT)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open heal→PENDING vs DENIED, glob matching, or CLI Assert (FINAL / S03).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Mig **014** CHECK rejects YOLO/garbage; four enums remain valid | `TestCapabilityToolDecisionCheckRejectsYOLO`; read `014_*.sql` CHECK |
| 2 | Existing YOLO heals to **PENDING** (not AUTO_ALLOWED, not dropped) | `TestCapabilityToolDecisionMigrateHealsYOLOToPending` |
| 3 | Resolve does **not** upsert AUTO_ALLOWED over unknown/healed YOLO | `TestResolveYOLOBuiltinDoesNotAutoAllow`; grep Resolve `default` / no fall-through |
| 4 | Unprefixed registered Name canonicalizes to `mcp:`+Name and gates MCP | `TestDecideUnprefixedMCPNameCanonicalizes`; `TestMCPUnprefixedDecideGatesCallTool` |
| 5 | Exact match only (no globs); custom slugs unchanged; `cli:` not rewritten | `TestCanonicalizeCustomAndCLISlugsUnchanged` |
| 6 | Footgun fold: unprefixed DENIED wins over prefixed AUTO_ALLOWED | `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed` |
| 7 | Compat ceiling **14**; no 015+; P15 Assert helper still `mcp:`+Name | `evals/compat`; grep `assertMCPToolAllowed`; `TestToolNamesRegistered` |
| 8 | DF-77 **not** implemented; no new MCP tools; YOLO/AllowAll absent | Diff + Notes; CLI still no `cli:` Assert |
| 9 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01 virgin keeper | Re-run `01` locked verify cmds |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P16-S02-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P16-S03-00** unless spawn.
5. If APPROVE: S06 must import named CHECK/heal/canonicalize + unprefixed-decide tests (already pointed on S06 SCOPE-TODOS).

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P16-S03-00** (unless spawn)
- [ ] No rewrite of done P16-S02-00/01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
