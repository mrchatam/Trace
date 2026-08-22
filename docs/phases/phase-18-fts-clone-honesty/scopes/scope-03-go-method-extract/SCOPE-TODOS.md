# S03 — Go method extract — scope todos

**Depends-on:** P18-S02-02 APPROVE (board order). No S02 code coupling. **S02-00 FINAL:** seed exclude unchanged (docs/help only). **S03-00 is FINAL.**

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** |
| 2 | 01-go-method-extract | implementer | **done** |
| 3 | 02-scope-review | reviewer | **APPROVE** — [REVIEW-NOTES.md](REVIEW-NOTES.md); next **P18-S04-00** |

## Phase locks (00 FINAL)

- DF-89 named golden **`TestIndexFileGoHandlerMethods`** in `internal/analyzers/analyzers_test.go`
- Fixture: `internal/analyzers/testdata/handler_methods.go` (new; 00 locked source). Do **not** mutate `sample.go`
- Exact `kind:name`: `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes` (sorted equality)
- Fail bar: the two `method:*` names must appear
- Keeper: `TestIndexFileGoGolden` (`sample.go` / `method:Run`)
- Live extract already matches (00 overlay). Fix `goSymbolQuery` / extract **only if** named test red; file-local incremental stays
- CGO: **`CGO_ENABLED=1`** required for analyzer tests
- No why-by-name; no indexer architecture; no FTS/seed/MCP; do not reverse DF-88

## Locked verify

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'
```

## Reminders
- SoT: [00-PLANNER.md](00-PLANNER.md) **FINAL**
- S04 VERIFY imports `TestIndexFileGoHandlerMethods` + keeper (live name from [REVIEW-NOTES.md](REVIEW-NOTES.md) after 02)
- S05 rebuild is after VERIFY — not this scope
