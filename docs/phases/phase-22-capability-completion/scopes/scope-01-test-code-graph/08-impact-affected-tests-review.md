# P22-S01-08 — Review: impact walk affected tests

## Metadata
- id: P22-S01-08
- todo_ids: [P22-S01-08]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C07** (and that C01–C03 still hold). Scope S01 is complete only if all four owned capabilities are met.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- [07-impact-affected-tests.md](07-impact-affected-tests.md)
- Live: `ImpactWalk` seeds file\|symbol; MCP `trace_impact` walk already exists; catalog 10; C31 query-by-test is **S05-03**

## Review checklist

| Check | Evidence |
|-------|----------|
| C07 | `TestImpactWalkIncludesAffectedTests` — reverse `validates` in blast / `affected_tests` |
| Depth cap | `TestImpactWalkDepthStillCapped` + existing `TestImpactWalk*` |
| Seeds unchanged | still file\|symbol only |
| No new MCP tool | `TestToolNamesRegistered` still 10 names |
| No `trace tests` CLI | grep `cmd/trace` — that is S05-03 |
| G19 | retrieval does not import cmd/mcp |
| C01–C03 regression | analyzer keepers below (include overlay EXTRACTED) |
| No architecture hop | Impact neighbors do not follow `architectural_boundary`; `architecture/<layer>` stubs never appear as FROM / blast via that rel |
| Compat | still **22**; no 023+ |

```bash
go test ./internal/retrieval/... -count=1 -run 'TestImpactWalk'
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestIndexDiscoversGoTestFunctions|TestValidatesEdgeExtractedFromImport|TestArtifactEdgesFunctionsTypesAPIs|TestArchitecturalBoundaryEdges|TestArchitecturalBoundaryOverlayExtracted|TestReplaceFileEdgesIsFileLocal'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run 'TestToolNamesRegistered'
```

## Spawn policy

If C07 (or C01–C03 regression) unmet: spawn **`P22-S01-08a` + `P22-S01-08b`**. Do not close S01 with residuals.

## Exit criteria

- [ ] C01, C02, C03, C07 closed or spawned
- [ ] Confidence **high** with re-run output in Notes
