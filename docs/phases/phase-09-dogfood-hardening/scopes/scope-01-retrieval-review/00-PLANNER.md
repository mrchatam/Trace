# P09-S01-00 — Retrieval review entity (DF-01)

## Metadata
- id: P09-S01-00
- todo_ids: [P09-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts to fix `retrieval: unknown entity type "review"` so `trace why`/`context` work after reviews exist. **No product Go in this row.**

## Live confirmation (2026-08-16)
| Surface | Finding |
|---------|---------|
| `internal/retrieval/exact.go` `lookupEntity` | switch has goal…evidence+file; **no** `case "review"` → hard error (not `isNotFound`) |
| Expand path | `ListLinksTo` on task with `review_judges_task` → `hitFromLinkNeighbor("review", …)` → same switch → **why/context fail** |
| Store | `GetReview` + `SyncEntityFTS("review", …)` already exist (mig 005); FTS sync ≠ ExactLookup |
| Domain | `CreateReview` / `SetReviewResult` / `LinkReviewTask` (`review_judges_task`) live (honesty D07 shape) |
| Fail-soft | **Forbidden** — skipping unknown types would hide reviews agents need in the chain |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Bug | DF-01 — `retrieval: unknown entity type "review"` |
| Primary fix | Add `case "review":` in `lookupEntity` (`exact.go`) via `store.GetReview` |
| Hit fields | `EntityType:"review"`; `Title` = review.Title; `Excerpt` = `Result` when non-empty else `excerpt(Body)`; `ReasonCode` from caller arg (unchanged signature) |
| Parallel switches | Grep/fix any other entity-type switches that omit `review` (today: Expand goal_id switch does **not** need a review case — links go through `lookupEntity`). Wire `hitFromLinkNeighbor`: map rels `review_judges_task` / `review_judges_scope` to reason codes equal to those rel strings (add consts alongside existing Decision/Discovery/Claim ones) |
| Prefer | **Full hit** — do **not** fail-soft-skip unknown `"review"` |
| Migration | **None** (reuse mig 005 `reviews` + `GetReview`) |
| Package | `internal/retrieval` (+ tests); no CLI/MCP surface change required |
| Regression | Plant like honesty/D07: CreateReview → LinkReviewTask → SetReviewResult PASS; assert `Why(task)` and `compiler.TaskContext` (or Expand depth 1) succeed and include a review neighbor hit |
| Named test | Prefer `TestWhyAndContextWithLinkedReview` (or split Exact+Why) under `internal/retrieval` and/or `internal/compiler` — implementer may choose one package if both paths are covered |
| Carry-forward | honesty Gate G (`TestHonestyEscapeRateGateGPrelim` + A/B/C) + `evals/p0x` + `evals/x0` + `go test ./...` |
| Forbidden | Fail-soft skip of `"review"`; new mig; daemon/HTTP/embeddings; weakening DONE/review policy; rewriting Phase 00–08 history |

## Exit
- [x] Thicken `01-retrieval-review.md` + `02-scope-review.md`
- [x] SCOPE-TODOS + board Notes; next **P09-S01-01**
- [x] Product Go — **not** this row
