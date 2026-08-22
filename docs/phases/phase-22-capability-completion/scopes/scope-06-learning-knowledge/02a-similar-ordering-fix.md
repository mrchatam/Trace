# P22-S06-02a — Implement: similar-changes ordering fix

## Metadata
- id: P22-S06-02a
- todo_ids: [P22-S06-02a]
- role: implementer
- skills: [incremental-implementation, test-driven-development, debugging-and-error-recovery]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Fix **C20** similar-change ordering so `TestQuerySimilarChanges` is stable when multiple changes share the same second-precision `created_at`. Spawned from **P22-S06-02** review: keeper **FAIL** on every `-count=1` run.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [02-patterns-similar-review.md](02-patterns-similar-review.md) — spawn rationale
- [01-patterns-similar.md](01-patterns-similar.md) — original deliverable
- `internal/domain/patterns.go` — `listChangesByKind`, `QuerySimilarChanges`
- `internal/store/patterns.go` — `ListChangesByPathPrefix`

## Root cause (locked)

`created_at` is **second-precision** RFC3339. Rapid `CreateChange` calls in tests (and real use) share a timestamp. Newest-first tie-break uses **`id DESC` (UUID lex)**, not insertion order → `TestQuerySimilarChanges` expects newer row first but gets older row when its UUID sorts higher.

Same class of bug fixed in **P22-S03-06a** (`ORDER BY created_at ASC, rowid ASC` for outcomes).

## Locked defaults

| Item | Value |
|------|-------|
| Scope | Ordering fix only — no behavior change to C19 pattern aggregation or CLI/MCP surface |
| Preferred fix | `ORDER BY created_at DESC, rowid DESC` in `ListChangesByPathPrefix`; mirror in `listChangesByKind` (load rowid or push kind filter to store SQL) |
| Schema / compat | **25** — no migration unless unavoidable |
| Limits | default **32**, cap **64** unchanged |
| MCP | catalog stays **13** |

## Requirements

1. `TestQuerySimilarChanges` PASS `-count=30` (non-flaky).
2. Add **`TestQuerySimilarChangesSameSecondTimestamps`** — force equal `created_at`, assert insertion/newest-first order.
3. Path-prefix and kind filters still mutually exclusive fail-closed.

## Named tests

| Test | Proves |
|------|--------|
| `TestQuerySimilarChanges` | C20 — kind + path filters, effects summary, newest-first |
| `TestQuerySimilarChangesSameSecondTimestamps` | **new** — equal `created_at` tie-break stable |
| Keepers | `TestPatternCountsFromChangesAndOutcomes`, `TestQuerySimilarChangesFailClosed`, S05 `TestCLIChanges*`, `TestChangesCompare` |

```bash
go test ./internal/domain/... -count=30 -run TestQuerySimilarChanges$
go test ./internal/domain/... ./internal/store/... -count=1 -run 'TestQuerySimilarChanges|TestPatternCountsFromChangesAndOutcomes|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestPatterns|TestChangesSimilar|TestChangesCompare|TestCLIChanges'
```

## Exit criteria

- [ ] C20 ordering fix landed
- [ ] Keeper non-flaky (`-count=30`)
- [ ] Board status + notes only
- [ ] Checklist **not** boxed (S06-02b closes C19+C20)

## Minimal todos

- [ ] Fix store + domain newest-first tie-break (rowid)
- [ ] Add same-second similar-changes test
- [ ] Board notes
