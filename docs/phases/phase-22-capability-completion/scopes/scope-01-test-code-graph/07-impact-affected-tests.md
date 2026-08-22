# P22-S01-07 — Implement: impact walk affected tests

## Metadata
- id: P22-S01-07
- todo_ids: [P22-S01-07]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Impact analysis **identifies potentially affected tests** via reverse `validates` edges. Closes **C07**.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- Live walk: `internal/retrieval/impact_walk.go` — `ImpactWalk(ctx, seeds, depth)`; seeds `file\|symbol` only; neighbors in `impactNeighborsFile` / `impactNeighborsSymbol`; `MaxImpactBlast = 64`; `maxImpactDepth = 2`
- Live tests: `internal/retrieval/impact_walk_test.go` — `TestImpactWalkMultiSeedExcludeSeeds`, `TestImpactWalkContainsAsymmetryNoSiblings`, `TestImpactWalkIncomingImportHop`, `TestImpactWalkLoudTruncation`, `TestImpactWalkHopRiskIncreases` (**keep green**)
- Live CLI: `cmd/trace/impact.go` `cmdImpactWalk` encodes `seeds`/`blast`/`blast_total`/`blast_kept`/`truncated`/`depth`
- Live MCP: `internal/mcp/tools_impact.go` `impactWalk` — same JSON keys; catalog **10** tools (`TestToolNamesRegistered`). **No new tool.** Query-by-test (C31) is **S05-03**.

## Locked defaults

| Item | Value |
|------|-------|
| Seeds | Still `file\|symbol` only (do not add `test` as a seed type) |
| Neighbors | Reverse `validates`: from a **symbol** seed, tests with `to_symbol_id = seed`; from a **file** seed, tests with `to_file_id = seed` (and contained symbols) |
| Blast | Tests appear in `blast` as `entity_type=symbol` (test `kind`) with `EdgeProvenance` from the edge |
| Convenience field | `ImpactWalkResult.AffectedTests` JSON `affected_tests` — **filter of kept blast**, not extra hops / not extra cap |
| Depth | Still **1..2**; `TestImpactWalkDepthStillCapped` = depth 3 fail-closed (today `ImpactWalk` already errors; name the test) |
| Cap | `MaxImpactBlast` still 64 on combined blast |
| CLI / MCP | Thin adapters: encode the library DTO (add `affected_tests` to the existing marshal maps in `cmdImpactWalk` and `impactWalk`). **Do not** add MCP tools or a `trace tests` CLI |
| G19 | `internal/retrieval` must not import `cmd/` or `internal/mcp` |
| Architecture stubs | Do **not** walk `architectural_boundary` as an impact hop. Layer identity files at `architecture/<layer>` are TO-only stubs (never IndexFile'd, never FROM). Neighbors stay reverse `validates` plus existing contains/imports. C02 keepers `TestArchitecturalBoundaryEdges` / `TestArchitecturalBoundaryOverlayExtracted` already cover membership; do not reimplement FileLayer in retrieval |

## Requirements

1. Extend `impactNeighborsSymbol` / `impactNeighborsFile` using `ListValidatesForSymbol` (and a file-level list if needed) from S01-01.
2. Named tests + keep existing `TestImpactWalk*` green.
3. If JSON shape grows: keep `TestImpactWalkCLI` (`cmd/trace/cli_test.go`) green.

## Touch files

- `internal/retrieval/impact_walk.go`, `impact_walk_test.go`
- `cmd/trace/impact.go` (thin JSON field)
- `internal/mcp/tools_impact.go` (same DTO keys)
- `internal/mcp/mcp_test.go` keepers (`TestToolNamesRegistered` stays **10**)

## Named tests

| Test | Proves |
|------|--------|
| `TestImpactWalkIncludesAffectedTests` | changing `Foo` lists `TestFoo` (via `validates`) |
| `TestImpactWalkDepthStillCapped` | depth >2 still fails closed |
| Existing `TestImpactWalk*` | stay green |

```bash
go test ./internal/retrieval/... -count=1 -run 'TestImpactWalk'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestImpact'
go test ./internal/mcp/... -count=1 -run 'TestToolNamesRegistered'
```

## Exit criteria

- [ ] C07 true
- [ ] Named tests PASS; MCP catalog still 10
- [ ] Board Notes

## Minimal todos

- [ ] Reverse validates in `impactNeighbors*`
- [ ] `affected_tests` on library DTO + CLI/MCP encode
- [ ] Tests + board notes
