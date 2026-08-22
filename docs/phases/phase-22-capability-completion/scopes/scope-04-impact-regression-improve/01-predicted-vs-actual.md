# P22-S04-01 — Implement: predicted vs actual impact

## Metadata
- id: P22-S04-01
- todo_ids: [P22-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Compare predicted impact with actual impact** after implementation (**C08**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/retrieval/impact_walk.go` (`ImpactWalkResult`, `AffectedTests`), `cmd/trace/impact.go` (walk only), `internal/store/changes.go` (`ListChangePaths`)
- S01: reverse `validates`, `AffectedTests` filter — reuse retrieval engine, do not duplicate BFS

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `ImpactWalk` + CLI/MCP walk | `impact_predictions` table, domain predict/compare |
| `change_paths` with optional `symbol_id` | `RecordPredictedImpact`, `CompareActualImpact` |
| `effects` expected/actual per dimension | blast-set snapshot/compare (C08) |
| Schema **023**; compat **23** | **024+ until this row** |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **`024_impact_compare.sql`** — creates **`impact_predictions`** + **`improvements`** (improvements CRUD is S04-05; **create both tables here**, S04-05 implements domain only) |
| Compat | Bump to **24** (`evals/compat`, store embed max) |
| One row per change | `change_id` PK; upsert on predict |
| Predict JSON | `{seeds, blast_keys[], affected_test_keys[], depth, blast_total, blast_kept, truncated}` — keys `"file:<uuid>"\|"symbol:<uuid>"`; **no paths with source text** |
| Seeds from change | `ListChangePaths(changeID)` → symbol seed if `symbol_id` set; else file id from indexed path |
| Compare | Re-walk stored seeds + stored depth; set diff → `compare_json` `{matched, unexpected, missed}` + `compared_at` |
| Fail-closed | No prediction row → compare error; unindexed path → predict error |
| CLI | `trace impact predict --change <id> [--depth 1\|2]`; `trace impact compare --change <id>` |
| G19 | Logic in `internal/domain/impact_compare.go` (+ store); CLI thin |
| MCP | **No new tools** — catalog stays **10** |

## Requirements

1. **`RecordPredictedImpact(ctx, changeID, walkResult)`** — serialize keys from `ImpactWalkResult`; upsert `impact_predictions`.
2. **`PredictImpactForChange(ctx, changeID, depth)`** — build seeds from paths, run `retrieval.Engine.ImpactWalk`, call RecordPredictedImpact.
3. **`CompareActualImpact(ctx, changeID)`** — load row, re-walk, diff sets, persist `compare_json` + `compared_at`.
4. CLI subcommands under existing `trace impact` (add `predict`, `compare`).
5. **`improvements` table in 024** (empty until S04-05) — do not implement CRUD this row beyond migration.
6. Named tests + compat bump.

## Touch files

- `internal/store/schema/024_impact_compare.sql`
- `internal/domain/impact_compare.go`, `impact_compare_test.go` (new)
- `internal/store/impact_compare.go` (new — or extend `impact.go` if tiny)
- `cmd/trace/impact.go`, `help.go`
- `evals/compat/compat_test.go`, store embed list
- `docs/CAPABILITIES_CHECKLIST.md` — box C08 only after S04-02 review (implementer: **unboxed** / note only)

## Named tests

| Test | Proves |
|------|--------|
| `TestRecordPredictedImpactThenCompareActual` | predict → mutate graph or paths → compare returns delta structure |
| `TestImpactCompareUnexpectedAndMissed` | controlled fixture: one unexpected + one missed key |
| `TestImpactCompareFailClosedWithoutPrediction` | compare without row → ErrValidation |
| `TestOpenCreatesDBAndMigratesIdempotent` | 024 applies cleanly |
| `TestMigrationStatusReportsEmbedMax` | embed max **24** |
| `TestCompatibilitySecurityChecklist` | ceiling **24** |
| `TestImpactWalkIncludesAffectedTests` | keeper (S01) |

```bash
go test ./internal/domain/... ./internal/store/... ./internal/retrieval/... -count=1 -run 'TestRecordPredictedImpact|TestImpactCompare|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestImpactWalkIncludesAffectedTests'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestImpact'
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C08 true (stored prediction vs actual blast keys, not one-off print)
- [ ] Compat **24**; exactly **24** sql files
- [ ] Checklist C08 **not** boxed until S04-02 review
- [ ] Board Notes: test output

## Minimal todos

- [ ] Mig 024 + store APIs
- [ ] Domain predict/compare + tests
- [ ] CLI predict/compare
- [ ] Compat/embed bump
- [ ] Board status + notes
