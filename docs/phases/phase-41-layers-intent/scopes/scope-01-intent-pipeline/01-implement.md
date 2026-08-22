# P41-S01-01 — Implement (G9 intent pipeline)

## Metadata
- id: P41-S01-01
- todo_ids: [P41-S01-01]
- role: implementer
- skills: [backend-dev, context-engineering, test-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Implement **G9**: bounded rule-based intent extraction in `internal/retrieval/` (G-009). M-001: merges into task loop; Laws 6–7 caps preserved.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G9](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-009](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [RETRIEVAL_AND_CONTEXT.md §3](../../../../RETRIEVAL_AND_CONTEXT.md)
- Live anchors (P41-S01-00 re-verified 2026-08-22 — **no drift vs P41-00**; G8 shipped S00):
  - `internal/retrieval/doc.go:7–16` — DR-NOSSEM semantic forbidden; L2/L3 reason codes shipped (G8)
  - `internal/retrieval/search.go:10–53` — FTS entry; `SearchOptions` today is `{Limit}` only — extend for intent
  - `internal/retrieval/types.go:60–63` — `SearchOptions`; add optional intent input field
  - `internal/retrieval/engine.go:8–26` — Engine API; no intent stage
  - `internal/retrieval/layer_enrich.go` — G8 L3 enrich (orthogonal; do not conflate with intent)
  - `internal/compiler/compiler.go:155–172` — title FTS (`task.Title`) + G1 query FTS (`opts.Query`); wire intent **before** both Search calls
  - `internal/compiler/compiler.go:14–38` — `ContextOptions` includes `MaxLayer` (G8) + `Query` (G1); intent is separate pre-channel stage
  - `internal/compiler/packet.go:95–114` — `Packet` has no intent field; optional `IntentSummary` add
  - `internal/compiler/packet.go:200–201` — "project intent" = Law 9 trust label only (do not repurpose)
  - `internal/compiler/explore.go:110–115` — standalone FTS uses raw `opts.Query`; pass intent when task loaded
  - `internal/store/entities.go:46–58` — `Task{Title, Body}` (no separate Objective field — use Body)
  - Evidence: [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-009 |
| Verdict | **Accept — implement bounded** (P41-00 lock) |
| Algorithm | Rule-based tokenization: `Task.Title` + `Task.Body` + optional query → deduped keywords; regex extract UUIDs, paths, symbols into `EntityHints` |
| No LLM | Deterministic only — no external API, no embeddings |
| G9 vs G1 | G9 extracts structured intent **before** FTS channels; G1 merges query hits **into** packet after retrieval (compiler.go:165–172 unchanged in behavior) |
| G9 vs G8 | G8 layer admission is post-candidate; G9 enriches FTS query strings pre-Search — orthogonal |
| Intent struct | `Intent{Keywords []string, EntityHints []EntityHint, Scope string, Source string}` — `Source` records inputs used (`task`, `query`, etc.) |
| Wiring (primary) | `ExtractIntent(IntentInput) Intent` in `intent.go`; `Search` calls it when `SearchOptions.Intent != nil` to build enriched FTS query string |
| Wiring (compiler) | `compileAtDepth` passes `IntentInput{TaskTitle, TaskBody, Query}` on both title-FTS and query-FTS Search calls (lines 157, 168) |
| Packet field | Optional compact `intent_summary` on `Packet` JSON when intent non-empty — capped (~256 chars total keywords) |
| §3 doc | Update §3: intent stage **shipped**; semantic leg **deferred (DR-NOSSEM)** — not aspirational whole-pipeline |
| Fallback | If blocked: doc-revise §3 as aspirational + ADR supersede — no silent skip |
| Graph export | If entities change: `trace seed export -o trace/graph.json` (intent summary on packet unlikely to trigger) |

## Touch-list (library → compiler wire → doc → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `internal/retrieval/intent.go` (new) | `IntentInput`, `Intent`, `EntityHint`, `ExtractIntent`, `Intent.FTSQuery()` (bounded keyword join for FTS) |
| 2 | `internal/retrieval/intent_test.go` (new) | G9-I1–I6 |
| 3 | `internal/retrieval/types.go` | Add `Intent *IntentInput` to `SearchOptions` (optional; nil = legacy raw-`q` behavior for MCP direct search) |
| 4 | `internal/retrieval/search.go` | When `opts.Intent != nil`, build FTS query via `ExtractIntent(*opts.Intent).FTSQuery()`; else use raw `q` (backward compat) |
| 5 | `internal/retrieval/doc.go` | Document intent pre-stage + note intent does not add new reason_code (still `fts_match`) |
| 6 | `internal/compiler/compiler.go` | Pass `SearchOptions{Limit: 16, Intent: &retrieval.IntentInput{...}}` at title + query Search sites (~157, ~168) |
| 7 | `internal/compiler/packet.go` (optional) | `IntentSummary *IntentSummary` on `Packet` — omit when empty; set from `ExtractIntent` in compile path |
| 8 | `internal/compiler/explore.go` (optional) | When task loaded, enrich explore FTS via same `IntentInput` pattern (~115) |
| 9 | `docs/RETRIEVAL_AND_CONTEXT.md` §3 | Revise pipeline diagram — intent shipped; semantic deferred; align with shipped compiler order |
| 10 | `docs/RETRIEVAL_AND_CONTEXT.md` §2 Semantic | Add DR-NOSSEM defer note cross-ref §3 |

**Explicit non-touch:**

- G-004a semantic/vector channel
- LLM inference packages
- MCP new tool — intent flows through existing search/context
- `web/` — out of scope unless trivial adapter passthrough already exists
- Replacing G1 query merge behavior (hits still merge at compiler.go:165–172)
- G8 `MaxLayer` / `layer_enrich.go` behavior

## Implementation order

```text
1. IntentInput + Intent + ExtractIntent + FTSQuery (rule-based, deterministic)
2. Unit tests G9-I1–I6 in intent_test.go
3. Extend SearchOptions; wire Search to ExtractIntent when Intent set
4. Compiler Search call sites — pass IntentInput (title+body+query)
5. Optional Packet.IntentSummary + explore.go parity
6. Revise RETRIEVAL_AND_CONTEXT §3 (+ §2 semantic defer cross-ref)
7. go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1
8. trace seed export if entity schema changes
```

## G9 vs G1 boundary (doc comment required)

Add package/func doc on `ExtractIntent`:

```text
G9 intent precedes retrieval channels: structured keyword/entity extraction from task+query
feeds FTS query building. G1 (compiler ContextOptions.Query) merges query FTS hits into the
packet after Expand — complementary, not duplicate. Intent never replaces task_id moat.
```

## Acceptance tests (must pass)

| ID | Suggested name | Assert |
|----|----------------|--------|
| G9-I1 | `TestExtractIntentFromTask` | `IntentInput{TaskTitle:"Ship auth", TaskBody:"Add JWT middleware"}` → non-empty `Keywords` containing stem tokens |
| G9-I2 | `TestExtractIntentEntityHints` | Body with UUID (`550e8400-e29b-41d4-a716-446655440000`), path (`internal/foo/bar.go`), symbol → `EntityHints` populated |
| G9-I3 | `TestExtractIntentQueryMerge` | Task + `Query:"extra token"` → query tokens appended; task tokens not duplicated |
| G9-I4 | `TestSearchUsesIntent` | `Search(ctx, "", SearchOptions{Intent: &input})` uses enriched FTS string vs raw empty `q`; or mock/store proves different hit set vs raw title-only |
| G9-I5 | `TestIntentNoSemantic` | No embedding/vector imports; no `semantic_match` reason_code; grep package for forbidden deps |
| G9-I6 | `TestExtractIntentDeterministic` | Same `IntentInput` → byte-identical `Intent` JSON across 10 calls |

## Regression tests (must stay green)

- All G1 tests (`TestG1QueryHitMerged`, `TestG1TitleFTSStillRunsWithQuery`, `TestG1QueryExpandDedupe`, `TestG1QueryCapHonesty`, `TestG1QuerySearchFailOpen`)
- All G8 tests (G8-L1–L7, G8-L2-MCP) — intent must not change default `MaxLayer` behavior
- `TestNoDumpAPI`
- DR-NOSSEM — no semantic channel

## Role work

1. Implement intent as **pre-retrieval stage** — bounded and deterministic.
2. Do not conflate with G1 packet merge or G8 layer admission.
3. Preserve MCP direct `Search(query)` backward compat when `Intent` nil.
4. Revise §3 doc to honest shipped state.
5. Self-check G9-I1–I6 before marking row done.

## Exit criteria

- [ ] G9-I1–I6 green **or** documented doc-revise supersede with ADR
- [ ] §3 updated to reflect shipped/deferred legs
- [ ] Board row → `done` with files + test command in Notes

## Next

`P41-S01-02`
