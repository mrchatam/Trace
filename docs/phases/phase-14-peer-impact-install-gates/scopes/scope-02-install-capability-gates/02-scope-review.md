# P14 / S02 / 02 — Scope review (Install / capability gates) FINAL

## Metadata
- id: P14-S02-02
- todo_ids: [P14-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 install/capability gates. Fresh subagent. Compare claims + **00-PLANNER FINAL** to live code/tests. Spawn `02a`/`02b` for blocker/high. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-install-capability-gates.md](01-install-capability-gates.md) — **FINAL**
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A4
- Research ranks 4–5; Laws 9 / 17
- Live: `internal/install`, `cmd/trace/install.go`, domain/store decisions, `evals/capability`, S01 ImpactWalk
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone.

## Checklist (FINAL)

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Registry + STABLE\|CONDITIONAL\|OPT_IN; Cursor STABLE Detect/Install/Uninstall | Read `internal/install`; `TestInstallDetectListsCursorStable`; `TestInstallCursorUninstallIdempotent` |
| 2 | CONDITIONAL refuses without marker; writes with marker | `TestInstallConditionalRefusesWithoutMarker` + `TestInstallConditionalWritesWithMarker` |
| 3 | detect CLI lists targets; no silent mass-write / install-all | CLI + help; grep forbids install-all YOLO |
| 4 | Uninstall idempotent; Cursor removes only `mcpServers.trace` | Uninstall test + code review |
| 5 | Existing Cursor print/`--write`/reload tip not weakened | Existing `TestInstallCursor*` PASS |
| 6 | Graduated allowlist: builtin exact auto-allow only; no AllowAll/globs | `TestCapabilityDecisionAutoAllowBuiltinMCP`; code review of auto-allow set |
| 7 | Unknown → PENDING; Assert fail-closed; human ALLOWED persists; DENIED blocks | Named decision tests |
| 8 | Durable audit in SQLite mig 013; ≠ chat-as-sole-record | Schema file + store CRUD; decisions list CLI |
| 9 | G19 thin adapters; **no** new MCP tools; nine + `trace_version` retained | `TestToolNamesRegistered` / MCP package; CLI-only decide |
| 10 | S01 ImpactWalk + Gate F **untouched** and green | Diff excludes walk; named ImpactWalk + `TestPlantedImpactConflictsGateFPrelim` PASS |
| 11 | Capability ablation not weakened | `evals/capability` PASS |
| 12 | No daemon/HTTP/embeddings/Neo4j/full-rebuild; no YOLO | Diff + laws skim |
| 13 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` | Re-run 01 verify cmds; Gate C artifacts untouched |
| 14 | Board Notes accurate; planner row had no product Go | P14-S02-00 Notes docs-only; implement scoped to S02-01 |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Diff implement vs inventory gaps; map each named test to code.
3. Fresh verify suite from 01 (include S01 ImpactWalk + Gate F).
4. Decide APPROVE / spawn `P14-S02-02a`/`02b` for blocker/high.
5. Write `REVIEW-NOTES.md` (optional but preferred) + board Notes; next **P14-S03-00** unless spawn.

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P14-S03-00** (unless spawn)
- [x] No rewrite of done P14-S02-00/01 history — spawn forward only
