# P17 / S01 / 01 — seed export (FINAL locks from 00-PLANNER)

## Metadata
- id: P17-S01-01
- todo_ids: [P17-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-80**, **DF-84** (plan tree keys), and **DF-85** (`exported_at_commit`) per sibling **00-PLANNER FINAL** + [DF-84-FORWARD.md](../../DF-84-FORWARD.md). Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- Live: `cmd/trace/seed.go`; `cmd/trace/help.go`; `internal/store/plan_hierarchy.go`; `internal/gitcli`; `fixtures/x0/seed/gt.json`
- P16 S05 (depend): do not re-implement `findings`/`alternatives`/`discovery_mentions_task` import
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks. Do **not** implement S03 idempotent upsert or S02 docs-only convention here. **No board spawn.**

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| CLI | **`trace seed export [-o <file>]`** stdout if no `-o`; **`trace seed import <file>`** unchanged |
| G19 | **`internal/domain/seed_export.go`** (name may vary) builds document; **`cmd/trace/seed.go`** thin encode + file I/O |
| DF-84 keys | `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` — shapes in 00-PLANNER |
| DF-85 | Export sets `exported_at_commit` from `git rev-parse HEAD` when repo; omit when unknown; import **allows + ignores** |
| Export include | Causal entities, seed link rels, findings/alternatives when present, full plan tree (all statuses / all deep-plan revisions) |
| Export exclude | **No** `transitions`; no index/token/lock/capabilities/tool decisions/events/reviews; tasks without `work_state` |
| Links | Canonical underscore rels; synthesize `goal_has_task` from `tasks.goal_id` |
| S01 plan import | `Insert*` / `UpsertGoalPlanState` — duplicate ids may fail until S03 |
| Forbidden | MCP seed tool; `.trace/` commit; duplicating P16 S05; HTTP/daemon; SHA-as-id |

Copy exact key lists, object fields, and import order from **00-PLANNER**. Do not invent alternate plan key names.

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Domain export | **`internal/domain/seed_export.go`** (+ `_test.go` optional) | `BuildSeedDocument`, `ExportOpts`, git HEAD helper |
| Store | `internal/store/helpers.go`, `plan_hierarchy.go`, `links.go`, `impact.go` | List-all readers for export surfaces |
| Seed CLI | `cmd/trace/seed.go` | `export` subcommand; extend `seedDocument` + allowlist + plan import loops |
| Help | `cmd/trace/help.go` | `seed export` line |
| Tests | `cmd/trace/cli_test.go` | Named DF-80/84/85 tests |

Do **not** add MCP tools. Do **not** implement idempotent link/plan upsert (S03).

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestSeedExportRoundTrip` | import → export → fresh import: same entity ids + links + **plan-tree ids** |
| `TestSeedExportOmitsDeniedSurfaces` | No `transitions` key; no index/token/capability/review payloads; no task `work_state` |
| `TestSeedExportWritesExportedAtCommit` | Git repo → non-empty SHA == HEAD; non-git → omit/empty; import still OK |
| Keepers | `TestSeedImportAndWhy`; `TestSeedImportDiscoveryMentionsTask`; `TestSeedImportImpactFindings`; `TestSeedImportFromIDAliases` |

TDD: red named tests first, then implement. Plan-tree round-trip is the **primary** `TestSeedExportRoundTrip` path.

## Role work
1. TDD named tests (red: no export subcommand; plan keys unknown on import; no `exported_at_commit` behavior).
2. Store list-all helpers + domain `BuildSeedDocument`.
3. Wire `seed export` + extend import allowlist/loops for plan keys + `exported_at_commit`.
4. Help line for export.
5. Self-check exit criteria; board **status + Notes only** → **P17-S01-02**.

## Locked verify (minimum)

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestSeedExport|TestSeedImport'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
```

Full product suites should stay green (carry-forward gates).

## Exit criteria
- [ ] DF-80 / DF-84 / DF-85 named tests green
- [ ] Help lists `seed export`; import accepts new keys (incl. `exported_at_commit` ignore)
- [ ] G19: export logic in domain, not duplicated in MCP
- [ ] Board Notes → **P17-S01-02**

## Minimal todos
- [ ] Red named tests
- [ ] Domain export + store list helpers
- [ ] CLI export + import extensions
- [ ] Help + verify + board sync
