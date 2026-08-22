# P17 / S02 / 01 — commit convention (FINAL locks from 00-PLANNER)

## Metadata
- id: P17-S02-01
- todo_ids: [P17-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, writing-for-agents]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-82** + **DF-85 docs** per sibling **00-PLANNER FINAL** + [DF-84-FORWARD.md](../../DF-84-FORWARD.md). Docs + help + one named test only. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- Live: `.gitignore`; `AGENTS.md`; `CONTRIBUTING.md`; `README.md`; `cmd/trace/help.go`
- S01 (depend, do not touch): `cmd/trace/seed.go`; `internal/domain/seed_export.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks. Do **not** implement S03 upsert or DF-86 git-hook. **No board spawn.**

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Path | **`trace/graph.json`** recommended commit path |
| `.gitignore` | **Unchanged** — `.trace/` only; **`trace/` committable** |
| Docs | AGENTS export-before-PR + CONTRIBUTING portable-graph section + README clone recipe |
| Help | Locked `seed export` block in 00-PLANNER (path + evidence not identity) |
| Attribution | Git author+SHA + `exported_at_commit` = evidence; actor / `as_operator` ≠ auth |
| Merge | Human resolve on `graph.json`; no merge driver; last-import-wins → **S03** |
| DF-86 | Hook **not** implemented; docs say optional later, export-before-PR suffices |
| Forbidden | seed/export code changes; gitignore `trace/`; git-hook; wrapping `git commit`; hosted MCP |

Copy exact locked strings from **00-PLANNER** — do not invent alternate paths or merge semantics.

## Extension points / files likely touched

| File | Change |
|------|--------|
| `AGENTS.md` | One **Portable graph** bullet under Hard boundaries |
| `CONTRIBUTING.md` | New **`## Portable graph (git)`** section |
| `README.md` | New **`## Portable graph (clone recipe)`** section |
| `cmd/trace/help.go` | Locked `seed export` help block only |
| `cmd/trace/cli_test.go` | **`TestHelpSeedExportPath`** |
| `.gitignore` | Verify only — **no edit** unless regression |

Do **not** touch `seed.go`, domain export, store, or MCP.

## Named tests (required)

| Test | Intent |
|------|--------|
| **`TestHelpSeedExportPath`** | Help mentions `trace/graph.json`, `exported_at_commit`, evidence-not-identity |
| Keepers | `TestHelpHandoffSoT` (DF-28); `TestAsOperatorFlagIdentityDocs` (DF-44); S01 `TestSeedExport*` unchanged |

TDD: red `TestHelpSeedExportPath` first, then help + docs.

## Role work
1. Red **`TestHelpSeedExportPath`** (help lacks path/evidence wording today).
2. Update **`cmd/trace/help.go`** per locked string.
3. Add AGENTS / CONTRIBUTING / README sections per 00-PLANNER.
4. Confirm `.gitignore` still `.trace/` only.
5. Self-check exit criteria; board **status + Notes only** → **P17-S02-02**.

## Locked verify (minimum)

```text
CGO_ENABLED=0 go test ./cmd/trace/... -count=1 -run 'TestHelp'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
```

Full product suites should stay green (carry-forward gates).

## Exit criteria
- [ ] **`TestHelpSeedExportPath`** green; DF-28/44 help keepers still PASS
- [ ] AGENTS + CONTRIBUTING + README sections match FINAL substance
- [ ] Help `seed export` block matches 00-PLANNER lock
- [ ] `.gitignore` still lists `.trace/` only (no `trace/` ignore)
- [ ] No seed/export/MCP code changes
- [ ] Board Notes → **P17-S02-02**

## Minimal todos
- [ ] Red `TestHelpSeedExportPath`
- [ ] Help string tweak
- [ ] AGENTS + CONTRIBUTING + README docs
- [ ] Verify + board sync
