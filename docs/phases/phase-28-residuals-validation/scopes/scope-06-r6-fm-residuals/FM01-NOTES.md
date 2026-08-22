# FM-01 / FR-P28-01 — notes (P28-S06-01)

**Date:** 2026-08-20  
**Status:** implemented (human gate preserved)

## What shipped

1. **`domain.ListPromotionCandidates`** — BLOCKING discoveries with no live `discovery_mentions_task` target.
2. **`SeedImportSummary.promotion_candidates` + `promotion_hint`** — filled at end of `ImportSeedDocument` (empty `[]` when none). Import still does **not** spawn tasks.
3. **`loop next`** — `buildPromotionCandidates` reuses the domain list (type alias).
4. **Harness** — `GapPassPrompt` step 4: after `seed import`, read `promotion_candidates` / `loop next`; promote or decline; no invented UUIDs.
5. **Help** — `seed import` documents stdout candidates + no auto-spawn.

## Human gate (FR-P28-D1)

Auto-spawn on import / `trace add discovery` remains **deferred**. Expansion requires explicit:

- `trace add task --from-discovery <discovery_id>`, or
- `trace loop apply` with `spawned_tasks[].discovery_id`

## Evidence

```bash
GOPROXY=direct go test ./internal/domain/ ./internal/loop/ ./internal/install/ -count=1 -run 'Promote|ImportSeedDocumentSurfaces|NextPromotion|GapPass'
GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1
```

Primary regression: `TestImportSeedDocumentSurfacesPromotionCandidates` — roster import → orphan BLOCKING → candidates + hint → promote with discovery UUID → candidates clear; no auto-spawn on import.
