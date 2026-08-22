# P22-S06-01 — Implement: patterns + similar changes

## Metadata
- id: P22-S06-01
- todo_ids: [P22-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Identify **recurring patterns** between change types and outcomes (**C19**) and allow **querying historical evidence before similar changes** (**C20**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/domain/changes.go`, `internal/store/changes.go`, `effects`, `outcome_results`, `regressions`, `improvements`
- S05: `ListChangesRecent`, `trace changes list|show` — similar query is **domain API**, not duplicate of list

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| Changes + paths + effects + outcomes + regressions + improvements | `change_patterns` table, pattern refresh |
| `trace changes list\|show\|capture\|compare` | `trace changes similar`, `trace patterns` |
| Schema **24**; compat **24** | **025+ until this row** |
| MCP catalog **13** | pattern/knowledge MCP tools |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **`025_engineering_knowledge.sql`** — creates **`change_patterns`** + **`engineering_knowledge`** (knowledge CRUD is S06-03; **create both tables here**) |
| Compat | Bump to **25** (`evals/compat`, store embed max, keeper tests) |
| change_kind | **`InferChangeKind`** — `seg:<first-path-segment>` from lexicographically smallest `change_paths.path`; no paths → `seg:unknown` |
| outcome_kind | **`ClassifyChangeOutcome`** — priority: regression > effect_contradicted > improvement > effect_supported > test_fail > test_pass > neutral |
| Pattern counts | positive/negative per [00-PLANNER.md](00-PLANNER.md); aggregate into `(change_kind, outcome_kind)` rows |
| Refresh | **`RefreshChangePatterns(ctx)`** — deterministic full rebuild; **no subprocess, no ML** |
| Similar | **`QuerySimilarChanges(ctx, SimilarChangesOpts)`** — `PathPrefix` **or** `ChangeKind` (mutually exclusive; both empty → ErrValidation); returns `{changes[], patterns[]}` compact JSON |
| Similar rows | Each change: id, task_id, change_kind, reason, created_at, paths[] (path only), effects[] (dimension, comparison), outcome_kind |
| Limits | default **32**, cap **64** (match S05 evidence queries) |
| CLI patterns | `trace patterns refresh` (rebuild + JSON summary `{ok,patterns_updated}`); `trace patterns list [--limit N]` (read stored rows) |
| CLI similar | Extend `trace changes` → **`similar --path <prefix> \| --kind <kind> [--limit N]`** JSON stdout |
| Capability | Add `cli:patterns` AUTO_ALLOW; extend `cli:changes` for `similar` |
| G19 | Logic in `internal/domain/patterns.go` + `internal/store/patterns.go`; CLI thin |
| MCP | **No new tools** — catalog stays **13** |
| Checklist | C19, C20 **unboxed** until S06-02 |

## Requirements

1. **`025_engineering_knowledge.sql`** — both tables per planner DDL; no ALTER on 024 tables.
2. **`InferChangeKind` / `ClassifyChangeOutcome`** — pure functions with unit tests; document priority in code comments only if non-obvious.
3. **`RefreshChangePatterns`** — SQL/group in store; domain orchestrates classification per change then upsert pattern rows.
4. **`QuerySimilarChanges`** — query `changes` joined with paths/effects; filter by prefix or precomputed kind; exclude OPEN changes optional — **lock: include all statuses except SUPERSEDED**.
5. **`ListChangePatterns`** store read for list CLI.
6. CLI + help + capability specs.
7. **`engineering_knowledge` table exists but empty** until S06-03 — no CRUD this row beyond migration.

## Touch files

- `internal/store/schema/025_engineering_knowledge.sql` (new)
- `internal/domain/patterns.go`, `patterns_test.go` (new)
- `internal/store/patterns.go`, `patterns_test.go` (new)
- `cmd/trace/patterns.go`, `patterns_test.go` (new)
- `cmd/trace/changes.go`, `changes_test.go` (extend — `similar`)
- `cmd/trace/root.go`, `help.go`
- `internal/domain/capability.go`
- `evals/compat/compat_test.go`, `evals/compat/doc.go`
- `internal/store/*_test.go` embed-max keepers (if separate from compat)

## Named tests

| Test | Proves |
|------|--------|
| `TestPatternCountsFromChangesAndOutcomes` | C19 — seed changes with known kinds/outcomes; refresh; assert pattern row counts |
| `TestQuerySimilarChanges` | C20 — same `seg:internal` kind returns prior changes + effects summary; path prefix filter works |
| `TestQuerySimilarChangesFailClosed` | empty path+kind → ErrValidation |
| `TestOpenCreatesDBAndMigratesIdempotent` | 025 applies cleanly |
| `TestMigrationStatusReportsEmbedMax` | embed max **25** |
| `TestCompatibilitySecurityChecklist` | ceiling **25** |
| `TestChangesCompare` | keeper (S02) — compare still PASS |

```bash
go test ./internal/domain/... ./internal/store/... -count=1 -run 'TestPatternCountsFromChangesAndOutcomes|TestQuerySimilarChanges|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestPatterns|TestChangesSimilar|TestChangesCompare'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 25
```

## Exit criteria

- [ ] C19, C20 true (evidence via named tests)
- [ ] Compat **25**; exactly **25** sql files
- [ ] MCP catalog **13** unchanged
- [ ] Checklist caps **unboxed** until S06-02
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Mig 025 + store pattern APIs
- [ ] Domain infer/classify/refresh/query + tests
- [ ] CLI patterns + changes similar
- [ ] Compat/embed bump to 25
- [ ] Board status + notes

## Residual risks (carry to S06-02)

- **`ClassifyChangeOutcome` priority** — reviewer must assert single bucket per change with fixture covering regression+improvement on same change
- **Path-prefix vs kind consistency** — `QuerySimilarChanges --kind` must use same `InferChangeKind` as refresh, not ad-hoc SQL
- **VCS-capture changes** (`trace:vcs-capture`) — include in patterns if paths present; document in test if excluded
- **Empty pattern table** — `patterns list` before refresh returns `[]`, not error
