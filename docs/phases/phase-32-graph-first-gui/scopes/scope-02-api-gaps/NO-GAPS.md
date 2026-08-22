# NO-GAPS — Phase 32 S02 library / OpenAPI

**Date:** 2026-08-21  
**Board:** P32-S02-01  
**Verdict:** No new library-backed `/v1` operations are required for the Phase 32 inspector map.

## Evidence

### UX-IA API gap table (`scopes/scope-01-ux-ia/UX-IA.md`)

| Gap | Kind | S02 action |
|-----|------|------------|
| `getImpact` in `ops.ts` | **Client glue** (not library) | Ship wrapper — separate deliverable this row |
| Library-backed HTTP ops for inspector | **None** | Why / context / graph / search / reviews / entity / task / loop already present |
| `/v1/path`, `listChanges`, `listRegressions` | **Not required** | No IA proof of need — omit; do not invent path API |
| **P32-PORT** | Serve UX (not an API gap) | Ship #1 min this row; not an OpenAPI addition |

### Live OpenAPI + handlers

Inspector-relevant paths already exist in `api/openapi.yaml` and `internal/httpapi/`:

| Path | `operationId` | Handler (live) |
|------|---------------|----------------|
| `/v1/why` | getWhy | present |
| `/v1/context` | getContext | present |
| `/v1/impact` | getImpact | `handleGetImpact` |
| `/v1/graph` | (graph neighborhood) | present |
| `/v1/search` | search | present |
| `/v1/reviews` (+ `{review_id}`) | list/get reviews | present |

`/v1/impact` documents optional `task_id`; response schema is `additionalProperties: true`. No OpenAPI change for impact in this row.

## Explicitly out (do not invent)

- `/v1/path`
- Requiring `listChanges` / `listRegressions` wrappers for the inspector map
- New core domain/HTTP ops for graph-home depth

## Same-row non-API deliverables

S02 is **not** API-only. This file covers the **library/OpenAPI** story only. Same implement row also ships:

1. **Client glue:** `getImpact(taskId)` in `web/src/api/ops.ts` → `GET /v1/impact?task_id=`
2. **P32-PORT #1:** friendly address-in-use stderr + `--addr` examples (fail-on-conflict unchanged; default `127.0.0.1:7432`)

**Deferred:** P32-PORT **#2** (auto free-port / `:0`) — not this row; S05 owns multi-project docs polish.
