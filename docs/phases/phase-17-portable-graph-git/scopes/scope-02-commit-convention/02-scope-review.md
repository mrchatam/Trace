# P17 / S02 / 02 — commit convention review (FINAL checklist)

## Metadata
- id: P17-S02-02
- todo_ids: [P17-S02-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of **DF-82 / DF-85 docs** vs FINAL locks + live repo. Fresh subagent ≠ implementer. Spawn `P17-S02-02a`/`02b` on blocker/high. Prefer `REVIEW-NOTES.md`. Next **P17-S03-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-commit-convention.md](01-commit-convention.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- Live: `.gitignore`; `AGENTS.md`; `CONTRIBUTING.md`; `README.md`; `cmd/trace/help.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open S01 export code or S03 upsert (FINAL deferrals).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | DF-82 path lock **`trace/graph.json`** | CONTRIBUTING + README + help; `TestHelpSeedExportPath` |
| 2 | `.gitignore` unchanged — **`.trace/` only** | Read `.gitignore`; no `trace/` or `graph.json` ignore |
| 3 | AGENTS export-before-PR one-liner (plan tree included) | Read `AGENTS.md` Hard boundaries bullet |
| 4 | Clone recipe: init → import → index → why/context/plan | README fenced block + CONTRIBUTING prose |
| 5 | **`exported_at_commit` / git SHA = evidence, not identity** | CONTRIBUTING + help wording; test substring |
| 6 | Git **author + SHA** documented as evidence | CONTRIBUTING attribution bullet |
| 7 | **`transition.actor` / review actor / `as_operator` ≠ auth** | CONTRIBUTING cross-ref DF-44; help still has Actor≠auth |
| 8 | Merge: human resolve on `graph.json`; **no merge driver** | CONTRIBUTING merge bullet |
| 9 | UUID last-import-wins pointer to **S03** | CONTRIBUTING (not implemented upsert here) |
| 10 | DF-86 hook **not** implemented; docs do not require it | Grep no `install git-hook`; CONTRIBUTING says optional later |
| 11 | DF-28 help handoff still present | `TestHelpHandoffSoT` |
| 12 | DF-44 help still present | `TestAsOperatorFlagIdentityDocs` |
| 13 | **No** S01 export/seed code changes | `git diff` scope limited to docs + help + test |
| 14 | No hosted MCP / no `trace-mcp` internet pointer | Grep changed files |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P17-S02-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P17-S03-00** unless spawn.

**Expected S03 gap (non-fail):** second import of same file may still UNIQUE-fail — idempotent upsert is S03.

## Locked verify (re-run)

```text
CGO_ENABLED=0 go test ./cmd/trace/... -count=1 -run 'TestHelp'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
```

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P17-S03-00** (unless spawn)
- [ ] No rewrite of done P17-S01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
