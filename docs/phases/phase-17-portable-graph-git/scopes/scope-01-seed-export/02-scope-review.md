# P17 / S01 / 02 — seed export review (FINAL checklist)

## Metadata
- id: P17-S01-02
- todo_ids: [P17-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of **DF-80 / DF-84 / DF-85** vs FINAL locks + live `seed export`. Fresh subagent ≠ implementer. Spawn `P17-S01-02a`/`02b` on blocker/high. Prefer `REVIEW-NOTES.md`. Next **P17-S02-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-seed-export.md](01-seed-export.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- Live: `cmd/trace/seed.go`; `internal/domain/seed_export.go` (or equivalent); `cmd/trace/help.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open S03 idempotent import or S02 docs convention (FINAL deferrals).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | DF-80 `trace seed export [-o]` + stdout default | Help line; `cmdSeed` dispatch; manual or test export |
| 2 | Round-trip ids + links + **plan-tree ids** | `TestSeedExportRoundTrip` green; inspect plan key shapes vs mig 006 |
| 3 | DF-84 export includes all plan rows (incl. SUPERSEDED deep plans) | Test fixture or store inspection; not ACTIVE-only export |
| 4 | DF-85 `exported_at_commit` evidence; not identity | `TestSeedExportWritesExportedAtCommit`; import allowlist accepts key; importer does not persist SHA |
| 5 | Exclude list: no transitions, index, token, capabilities, tool decisions, events, reviews | `TestSeedExportOmitsDeniedSurfaces` |
| 6 | Tasks export without `work_state`; links use canonical underscore rels | JSON inspection in round-trip / omit test |
| 7 | G19: library export + thin CLI; **no** MCP seed tool | Domain owns build; grep no new MCP registration |
| 8 | Did **not** duplicate P16 S05 import work | Findings/alternatives/mentions-task import unchanged except allowlist extension |
| 9 | `.gitignore` still `.trace/` only | Read `.gitignore` |
| 10 | P16 seed keepers still PASS | `TestSeedImportDiscoveryMentionsTask`, `TestSeedImportImpactFindings`, `TestSeedImportAndWhy` |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P17-S01-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P17-S02-00** unless spawn.

**Expected S03 gap (non-fail):** second import of same export file may still UNIQUE-fail on links/plan rows — idempotent behavior is S03, not S01.

## Locked verify (re-run)

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestSeedExport|TestSeedImport'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
```

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P17-S02-00** (unless spawn)
- [ ] No rewrite of done P17-S01-00/01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
