# P22-S04-05 — Implement: record improvements

## Metadata
- id: P22-S04-05
- todo_ids: [P22-S04-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Record improvements** as first-class queryable rows (**C18**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — **`improvements` table created in S04-01 mig 024**
- Live: no improvement API; `effects.comparison=supported` is not C18

## Live baseline

| Present | Absent |
|---------|--------|
| `improvements` table (if S04-01 done) | domain CRUD, CLI list, seed key |
| Seed export for regressions/effects | `improvements[]` export/import |
| Schema **024**; compat **24** | **025+** |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **None new** — table from **024**; if S04-01 skipped (should not happen), add table in 024 only |
| Compat | Stays **24** (forbid **025+**) |
| Fields | `id`, `change_id`, `task_id`, `dimension`, `summary`, `evidence_ids_json` (JSON array ≤32 ids), `source_type`, `confidence`, timestamps |
| Summary cap | 4096 bytes (match regression summary) |
| API | **`RecordImprovement`**, **`GetImprovement`**, **`ListImprovementsByChangeID`**, **`ListImprovementsByTaskID`** |
| Evidence | Validate evidence ids exist on record; store ids in JSON column + optional `improvement_supported_by` links — **lock: JSON column sufficient; links optional bonus** |
| CLI | **`trace outcomes improvements --change <id>`** or **`--task <id>`** (extend `cmd/trace/outcomes.go`) |
| Seed (D-22-19) | Add **`improvements[]`** to `SeedDocument` export/import round-trip |
| Auto from effects | **Out of scope** — explicit `RecordImprovement` only (S06 may synthesize knowledge from rows) |
| MCP | S05 query scope — not this row |

## Requirements

1. Store CRUD in `internal/store/improvements.go`.
2. Domain validation + `RecordImprovement` / list APIs.
3. CLI list JSON under `trace outcomes improvements`.
4. Seed export/import with stable ids.
5. Named tests.

## Touch files

- `internal/domain/improvements.go`, `improvements_test.go` (new)
- `internal/store/improvements.go` (new)
- `internal/domain/seed_export.go`, `seed_import.go`, `seed_export_test.go`
- `cmd/trace/outcomes.go`, `help.go`
- `docs/CAPABILITIES_CHECKLIST.md` — C18 note only until S04-06 review

## Named tests

| Test | Proves |
|------|--------|
| `TestRecordImprovementQueryable` | record + list by change and task |
| `TestSeedExportIncludesImprovements` | export/import round-trip when rows exist |
| `TestRecordImprovementFailClosedEmptySummary` | validation |
| `TestSeedExportIncludesP20Cognition` | keeper — existing seed keys still work |

```bash
go test ./internal/domain/... -count=1 -run 'TestRecordImprovement|TestSeedExportIncludesImprovements|TestSeedExportIncludesP20Cognition'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestOutcomes'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # still 24
```

## Exit criteria

- [ ] C18 true
- [ ] Named tests PASS; compat **24**
- [ ] Checklist C18 **not** boxed until S04-06 review
- [ ] Board Notes

## Minimal todos

- [ ] Store + domain CRUD
- [ ] CLI list
- [ ] Seed key
- [ ] Tests
- [ ] Board notes
