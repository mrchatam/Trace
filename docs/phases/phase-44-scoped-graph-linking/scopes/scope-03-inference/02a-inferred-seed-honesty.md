# Phase 44 / S03 / 02a — Fix INFERRED seed export/import honesty

## Metadata
- id: P44-S03-02a
- todo_ids: [P44-S03-02a]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- verification: automated

## Objective

Close **H1 high** from **P44-S03-02**: portable seed round-trip must leave `entity_links.source_type=INFERRED` (and confidence `0.4`) honest — **no silent upgrade to `USER_ASSERTED`**.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Re-read P44-S03-02 board Notes before coding.

## References

- [02-review.md](02-review.md) — H1 seed honesty row (FAIL evidence)
- [01-implement.md](01-implement.md) — “Export honesty | INFERRED rows round-trip with `source_type=INFERRED`”
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §3.4 provenance
- Live: `internal/domain/seed_export.go` (`SeedLink`, `BuildSeedDocument` link loop ~468), `seed_import.go` (`ImportSeedLink` → `LinkMeta{}`)

## Root cause (locked)

1. `SeedLink` has **no** `source_type` / `confidence` fields.
2. Export writes `SeedLink{Rel, From, To}` only — drops provenance.
3. Import always uses `meta := LinkMeta{}` → `withDefaults()` → `USER_ASSERTED` / confidence 0.

Proven by independent review repro: after `InferScopes` → `BuildSeedDocument` → `ImportSeedDocument`, `scope_member` lands as `USER_ASSERTED`.

## Locked defaults

| Item | Value |
|------|-------|
| Scope | Seed link provenance only — **do not** change inference rules, CLI, MCP, index hooks, or GUI |
| Wire | Optional JSON fields on `SeedLink`: `source_type`, `confidence` (omitempty OK for legacy seeds) |
| Import | Prefer exported `source_type`/`confidence` when present; empty → existing `LinkMeta{}` defaults (legacy) |
| Rel strings | Still **not** prefixed with `inferred_` |
| Schema | No new migration unless unavoidable (links already store source_type) |
| Caps / D4 | Unchanged — CLI-only infer; no index-hook |

## Requirements

1. Export every `entity_links` row on the seed allowlist **including** `source_type` + `confidence`.
2. Import applies those fields via `LinkMeta` (do not wipe INFERRED → USER_ASSERTED).
3. Named test: infer → export → import → assert `source_type=INFERRED` and `confidence=InferredLinkConfidence`.
4. Keep S02 keepers: `TestSeedExportImportScopesRoundTrip`, `TestNoINFERREDFromScopePaths` (S02 paths still must not *create* INFERRED; round-trip of CLI-created INFERRED is S03a).

## Named tests

| Test | Proves |
|------|--------|
| `TestINFERREDScopeMemberSeedRoundTrip` (or equiv.) | **new** — INFERRED survives export/import |
| `TestSeedExportImportScopesRoundTrip` | S02 scopes CRUD/export still OK |
| `TestInferScopes_*` keepers | Inference behavior unchanged |

```bash
go test ./internal/domain/ -count=1 -run 'TestINFERREDScopeMemberSeedRoundTrip|TestSeedExportImportScopesRoundTrip|TestInferScopes_|TestNoINFERREDFromScopePaths'
```

## Exit criteria

- [ ] SeedLink carries honest provenance; round-trip test PASS
- [ ] Board Notes only (status + notes); next **P44-S03-02b**
- [ ] No `web/` / no new MCP infer tool / no index-hook

## Minimal todos

- [ ] Extend `SeedLink` + export + import
- [ ] Add round-trip honesty test
- [ ] Board notes → **P44-S03-02b**

## Next

`P44-S03-02b`
