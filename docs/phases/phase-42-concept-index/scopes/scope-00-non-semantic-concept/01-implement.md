# P42-S00-01 — Implement (G6 non-semantic concept)

## Metadata
- id: P42-S00-01
- todo_ids: [P42-S00-01]
- role: implementer
- skills: [backend-dev, context-engineering, test-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Implement **G6**: non-semantic concept retrieval via graph labels/summaries (G-004b). M-001: merges into task loop; Laws 6–7 caps preserved; no vector channel.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [LAW-REVIEW-NOTES.md](LAW-REVIEW-NOTES.md) — **required PASS** from S00-00
- [REMEDIATION-PLAN G6](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-004b](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Live anchors (P42-S00-00 re-verified 2026-08-22):
  - `internal/retrieval/doc.go:8–9` — DR-NOSSEM; extend with graph-label channel honesty
  - `internal/retrieval/types.go:8–26` — add `ReasonGraphLabelMatch`
  - `internal/retrieval/search.go:9–44` — pattern for bounded FTS + intent
  - `internal/retrieval/intent.go` — G9 `ExtractIntent` / `IntentInput` feed G6 terms
  - `internal/retrieval/engine.go` — Engine entry point
  - `internal/store/fts.go:25–69` — `SearchFTS` + entity_type in hits
  - `internal/compiler/compiler.go:155–180` — merge point after intent, before file-seed expand
  - `internal/compiler/explore.go:112–130` — explore search path (mirror compile merge)
  - `internal/retrieval/retrieval_test.go`, `p20_test.go` — extend test package

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Preflight

- Confirm [`LAW-REVIEW-NOTES.md`](LAW-REVIEW-NOTES.md) shows **PASS** — else `blocked` with reason.
- Re-read `00-PLANNER.md` locked defaults before coding.
- Re-read live anchors: `doc.go:8–9`, `types.go:8–26`, `search.go:9–44`, `compiler.go:155–180`, `explore.go:110–128`, `fts.go:151–160`.

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-004b |
| Verdict | **Accept — ship** (P42-00 lock) |
| Reason code | `graph_label_match` — new locked string |
| Concept entity types | `discovery`, `assumption`, `decision`, `goal`, `claim` (bounded set; no effect/file) |
| Query source | G9 `ExtractIntent` keywords + entity hints (not raw title-only) |
| Limit | Default 16 when opts.Limit ≤ 0; hard max 64 (match Search cap logic in `search.go:15–20`) |
| Merge | Append concept hits to compile/explore candidates; fail-open on error (DF-87) |
| Dedupe | Same entity key as generic FTS hit → keep `graph_label_match` (drop duplicate `fts_match`) |
| 1-hop graph attach | **Defer** — optional peer-inspired hop not required for G6 closure (LAW-REVIEW PASS) |
| Fallback | If blocked: document in Notes — do not slip vector leg |
| Graph export | If entities/reason_codes change registry: `trace seed export -o trace/graph.json` |

## Touch-list (library → adapters → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `internal/retrieval/types.go` | Add `ReasonGraphLabelMatch = "graph_label_match"` |
| 2 | `internal/retrieval/concept.go` | **New** — `SearchGraphLabels(ctx, intent Intent, opts SearchOptions) ([]Hit, error)` |
| 3 | `internal/retrieval/concept.go` | Filter FTS hits to concept entity types; set reason_code |
| 4 | `internal/retrieval/doc.go` | Document graph-label channel + DR-NOSSEM boundary |
| 5 | `internal/compiler/compiler.go` | After title FTS + query FTS (`~L170–180`), before file-seed expand: `SearchGraphLabels`; merge + dedupe |
| 6 | `internal/compiler/explore.go` | After `eng.Search` (`~L117`), call `SearchGraphLabels`; append to `out.SearchHits`; dedupe; fail-open |
| 7 | `internal/retrieval/concept_test.go` | **New** — G6-C1–C7 tests |
| 8 | `internal/compiler/compiler_test.go` | G6 compile-merge integration test (G6-C3) |
| 9 | `docs/RETRIEVAL_AND_CONTEXT.md` §2 | Mark graph-label channel shipped (not semantic) |

**Explicit non-touch:**

- `web/`, `internal/httpapi/` — no GUI/HTTP concept surface in S00
- G-004a vector / embedding index
- New MCP tool (reuse existing context/explore paths)
- Full-graph scan / unbounded entity iteration
- LLM concept extraction

## SearchGraphLabels sketch (library — Law 19)

```go
// conceptEntityTypes — bounded set; no file/symbol/task in this channel.
var conceptEntityTypes = map[string]struct{}{
    "discovery": {}, "assumption": {}, "decision": {},
    "goal": {}, "claim": {},
}

// SearchGraphLabels: G9 intent terms → store.SearchFTS (multi-OR like search.go)
// → filter hits to conceptEntityTypes → ReasonGraphLabelMatch → cap limit.
func (e *Engine) SearchGraphLabels(ctx context.Context, intent Intent, opts SearchOptions) ([]Hit, error)
```

Reuse `intentSearchTerms(intent)` from `search.go` (same package). Do **not** duplicate G9 extraction.

## Compiler merge insertion (exact)

After `compiler.go:180` (query FTS append), before file-seed expand block:

```go
labels, err := c.retr.SearchGraphLabels(ctx, taskIntent, searchOpts)
if err != nil {
    labels = nil // DF-87 fail-open
}
candidates = mergeConceptHits(candidates, labels) // dedupe: graph_label_match wins
```

## Explore merge insertion (exact)

After `explore.go:124` (`out.SearchHits = hits`), when `searchQ != ""`:

```go
labelHits, lerr := eng.SearchGraphLabels(ctx, retrieval.ExtractIntent(intentIn), searchOpts)
if lerr != nil {
    labelHits = nil
}
out.SearchHits = mergeConceptHits(out.SearchHits, labelHits)
```

## Implementation order

```text
1. ReasonGraphLabelMatch + doc.go honesty
2. SearchGraphLabels in retrieval/concept.go (store FTS + entity filter)
3. mergeConceptHits helper (retrieval or compiler — prefer retrieval if shared)
4. Unit tests G6-C1–C7 in concept_test.go
5. Wire compiler.go + explore.go merge (fail-open)
6. Compiler integration test G6-C3
7. REVISE RETRIEVAL_AND_CONTEXT.md §2 graph-label bullet
8. go test ./internal/retrieval/... ./internal/compiler/... -count=1
9. trace seed export if entity schema changes
```

## Acceptance tests (must pass)

| ID | Suggested name | Assert |
|----|----------------|--------|
| G6-C1 | `TestSearchGraphLabelsDiscovery` | Intent keyword matches discovery body → hit with `graph_label_match` |
| G6-C2 | `TestSearchGraphLabelsEntityFilter` | file/symbol hits excluded from concept channel |
| G6-C3 | `TestContextIncludesGraphLabels` | Compile with task + matching discovery → packet item with `graph_label_match` |
| G6-C4 | `TestSearchGraphLabelsCap` | Limit honored; no unbounded results |
| G6-C5 | `TestSearchGraphLabelsNoSemantic` | No vector deps; no `semantic_match` reason_code |
| G6-C6 | `TestSearchGraphLabelsDeterministic` | Same input → same hits order |
| G6-C7 | `TestSearchGraphLabelsFailOpen` | Inject FTS error (test double or closed store) on concept path → `Context()` still returns packet with task seed |

## Regression tests (must stay green)

- All G8/G9/G1 tests
- `TestNoDumpAPI`
- Default packet shape unchanged when no concept matches

## Role work

1. Implement graph-label channel as **bounded FTS over concept entity types** — not semantic.
2. Wire compile/explore as thin merge (Law 19).
3. Preserve moat: task_id required on compile path.
4. Self-check G6-C1–C7 before marking row done.

## Exit criteria

- [ ] G6-C1–C7 green
- [ ] LAW-REVIEW-NOTES PASS cited in Notes
- [ ] `RETRIEVAL_AND_CONTEXT.md` §2 updated
- [ ] Board row → `done` with files + test command in Notes

## Next

`P42-S00-02`
