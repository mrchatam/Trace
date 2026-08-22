# P22-S03-07 — Implement: invariants + iteration comparison

## Metadata
- id: P22-S03-07
- todo_ids: [P22-S03-07]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Verify **architectural/invariant constraints** where possible (**C14**) and **compare results between iterations** (**C15**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- S01: `internal/analyzers/architecture_graph.go`, `store.ListCrossLayerImports`, `store.FileLayer`
- Changes: latest task change paths via `ListChangePaths`
- S07 will extend rules via `trace/eval-rules.json` — this row ships **one** default forbidden rule only

## Locked defaults

| Item | Value |
|------|-------|
| Invariant check | `domain.CheckArchitecturalInvariants(ctx, taskID)` → `{passed, violations[]}` |
| Default rule | **`internal/` must not import `cmd/`** — detect via cross-layer import where `FromLayer=internal` and `ToLayer=cmd` on paths touched by latest change (or new imports in those files) |
| Scope | Only paths in latest RECORDED/COMPARED change for task; use indexed `imports` + layer membership — **do not** full-repo glob (Law 12) |
| Record | Optional: failed check may append advisory `kind=verification` with summary listing violation paths (evidence optional); fail-closed for “verified” gate if used as hard gate — default **advisory JSON only** from CLI |
| Iterations | `domain.CompareIterationOutcomes(ctx, taskID, kind)` — last two rows of `kind` (`test` or `evaluation`) by `created_at`: emit `{previous, current, delta}` JSON (`test_status` change and/or score dimension deltas) |
| CLI | `trace verify invariants --task <id>` JSON; `trace outcomes compare --task <id> --kind test|evaluation` JSON |
| Schema / compat | **23** — no migration |

## Requirements

1. Invariant library + tests using analyzers testdata (cross-layer fixture from S01-05).
2. Iteration compare pure on stored outcomes (no re-exec tests).
3. CLI wired in `root.go` / `help.go`; capability keys `cli:verify`, `cli:outcomes`.
4. Named tests below.

## Touch files

- `internal/domain/invariants.go`, `invariants_test.go` (**new**)
- `internal/domain/outcomes.go` — `CompareIterationOutcomes` (+ test)
- `cmd/trace/verify.go` — extend with `invariants` subcommand
- `cmd/trace/outcomes.go`, `outcomes_test.go` (**new**)
- `cmd/trace/root.go`, `help.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestInvariantFailOnForbiddenLayerImport` | internal→cmd import in change paths fails |
| `TestInvariantPassWhenNoCrossLayer` | clean change passes |
| `TestCompareIterationOutcomes` | two test outcomes → status delta JSON |
| `TestCompareIterationOutcomesEvaluation` | two evaluations → score delta |
| `TestArchitecturalBoundaryEdges` | keeper (S01) |

```bash
go test ./internal/domain/... ./internal/analyzers/... -count=1 -run 'TestInvariant|TestCompareIterationOutcomes|TestArchitecturalBoundaryEdges'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestVerifyInvariants|TestOutcomesCompare'
```

## Exit criteria

- [ ] C14 and C15 true
- [ ] Named tests PASS; compat **23**
- [ ] Board Notes

## Minimal todos

- [ ] Invariant checker (default rule)
- [ ] Iteration compare API
- [ ] CLI + tests
- [ ] Board notes
