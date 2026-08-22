# FM-08 / FR-P28-05 — notes (P28-S06-09)

**Date:** 2026-08-20  
**Status:** implemented (no auto-spawn; D1 preserved)

## What shipped

1. **`traceAddDescription` (INT-06 reinforce)** — discovery-first kind list retained; explicit prefer task/promotion over discovery-only; after BLOCKING discovery promote via `trace_add kind=task` or `spawned_tasks` + `discovery_id` **before product edits**.
2. **MCP `AddInput.Kind` jsonschema** — aligned to `discovery|task|goal|…` (was goal-first).
3. **`GapPassPrompt` step 3** — **Post-discovery nudge (FM-08):** after CLI/MCP discovery add (esp. BLOCKING), promote with `--from-discovery` or `spawned_tasks[].discovery_id` before product edits; do not discovery-only then edit. Steps 4–7 (import candidates, write-before-export) preserved.

## Human gate (FR-P28-D1)

No auto-spawn on `trace_add` / `trace add discovery`. Promotion remains explicit CLI/apply only.

## Evidence

```bash
GOPROXY=direct go test ./internal/mcp/... ./internal/install/... -count=1 -run 'TraceAddDescription|TraceAddKindSchema|GapPass'
GOPROXY=direct go test ./internal/loop/... -count=1 -run 'PromotesBlockingDiscovery'
GOPROXY=direct go test ./internal/mcp/... ./internal/install/... -count=1
```

- Description regression: `TestTraceAddDescriptionMentionsPromotionPath` (+ prefer / discovery-only / before product edits).
- Schema order: `TestTraceAddKindSchemaListsDiscoveryFirst`.
- GapPass: `TestGapPassPromptNonEmpty` asserts Post-discovery nudge phrases.
- Apply smoke: `TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask` (discovery → spawned task before further work).
