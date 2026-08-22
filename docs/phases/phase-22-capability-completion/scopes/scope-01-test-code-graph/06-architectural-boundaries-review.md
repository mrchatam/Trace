# P22-S01-06 — Review: architectural boundaries

## Metadata
- id: P22-S01-06
- todo_ids: [P22-S01-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C02** is fully met without a full-rebuild indexer.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- [05-architectural-boundaries.md](05-architectural-boundaries.md)
- Live: `IndexFile` is the incremental unit; `cmd/trace/index.go` already walks files one-by-one

## Review checklist

| Check | Evidence |
|-------|----------|
| C02 edges stored | `TestArchitecturalBoundaryEdges` — `rel=architectural_boundary` |
| Incremental | `TestArchitecturalBoundaryIncremental` — not a full-graph rebuild on one-file change |
| Outgoing merge | `TestArchitecturalBoundarySurvivesFileReindex` — `indexCodeEdges` / `ReplaceFileEdges` batch includes `architectural_boundary` **with** validates+contains+exports. Grep: exactly one `ReplaceFileEdges` per IndexFile path after pre-clear |
| Pre-clear landmine | `IndexFile` still `ReplaceFileEdges(path, nil)` then one rewrite. Boundary-only second replace deletes C01/C03 edges |
| No 023+ | schema dir max **022**; compat **22** |
| No graph DB / LLM layers | grep; overlay `trace/architecture.json` is optional |
| Keepers | `TestReplaceFileEdgesIsFileLocal`; validates/contains still present |

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestArchitecturalBoundaryEdges|TestArchitecturalBoundaryIncremental|TestArchitecturalBoundarySurvivesFileReindex|TestReplaceFileEdgesIsFileLocal'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Spawn policy

If C02 unmet or Law 12 violated: spawn **`P22-S01-06a` + `P22-S01-06b`**. Do not close with residuals.

## Exit criteria

- [ ] C02 closed or spawned
- [ ] Confidence **high** with re-run output in Notes
