# P22-S01-04 — Review: artifact relationships

## Metadata
- id: P22-S01-04
- todo_ids: [P22-S01-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C01** is fully met (all artifact kinds in the bullet, including tests from S01-01).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- [03-artifact-relationships.md](03-artifact-relationships.md)
- Live: `imports` already models dependencies; `symbols.kind` is free text; no modules table

## Review checklist

| Check | Evidence |
|-------|----------|
| Files, modules, components, functions, types, APIs | `TestArtifactEdgesFunctionsTypesAPIs` + `TestModuleContainsSymbols` |
| Tests still validate | S01-01 keepers green (`kind=test` + `validates`) |
| Incremental | `TestReplaceFileEdgesIsFileLocal` — outgoing batch includes validates **and** contains/exports |
| No `depends_on` duplicate of `imports` | grep `depends_on` writes; C01 deps remain the `imports` table |
| No new mig 023+ | `internal/store/schema/` still max **022**; compat **22** |
| No new ontology table | no `modules` / `components` SQL |

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestArtifactEdgesFunctionsTypesAPIs|TestModuleContainsSymbols|TestIndexDiscoversGoTestFunctions|TestValidatesEdgeExtractedFromImport|TestReplaceFileEdgesIsFileLocal'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Spawn policy

If C01 still unmet: spawn **`P22-S01-04a` + `P22-S01-04b`** immediately below. **Do not** close with residuals for later.

## Exit criteria

- [ ] C01 closed or spawned
- [ ] Confidence **high** with re-run output in Notes
