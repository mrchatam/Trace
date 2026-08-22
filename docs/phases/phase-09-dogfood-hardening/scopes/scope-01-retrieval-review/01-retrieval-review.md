# P09 / S01 / 01 — Retrieval review support (DF-01)

## Metadata
- id: P09-S01-01
- todo_ids: [P09-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Fix DF-01 so `trace why` / `trace context` (and library `Why` / `Expand` / `TaskContext`) succeed when a **review** is linked to a task. Add `case "review"` to ExactLookup (`lookupEntity`) and wire review link reason codes. Prefer a full Hit over fail-soft skip. Keep carry-forward gates green. **Do not** change DONE/review promotion policy.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL 2026-08-16
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-01 / DF-09
- Live: `internal/retrieval/{exact,expand,why,types}.go`; `store.GetReview`; honesty plant shape in `evals/honesty`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (FINAL — P09-S01-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Bug | DF-01 `retrieval: unknown entity type "review"` |
| Fix site | `lookupEntity` in `internal/retrieval/exact.go` — add `case "review"` |
| Store API | Existing `(*store.Store).GetReview` (mig 005) — **no** new migration |
| Hit | `EntityType: "review"`; `Title: r.Title`; `Excerpt: r.Result` if non-empty else `excerpt(r.Body)`; `ReasonCode` / `Score` / `Distance` from caller |
| Link reasons | In `hitFromLinkNeighbor` (`expand.go`): map `review_judges_task` + `review_judges_scope` to reason-code consts with those exact strings (mirror Decision/Discovery/Claim mapping) |
| Policy | **Full hit required** — do **not** treat unknown `"review"` as soft-skip / empty neighbor |
| Regression plant | Like D07/honesty: `CreateReview` → `LinkReviewTask` (`review_judges_task`) → `SetReviewResult(..., PASS, …)` then `Why("task", taskID)` **and** context path (`compiler.TaskContext` or Expand depth-1 seeds) must succeed with a review neighbor present |
| Named test | Prefer `TestWhyAndContextWithLinkedReview` (name may vary; must prove why+context after review exists) |
| Packages | Product Go in `internal/retrieval` (+ tests). Compiler test OK if needed to prove TaskContext. No new CLI/MCP tools |
| Carry-forward | honesty A/B/C + Gate G; `evals/p0x`; `evals/x0`; `CGO_ENABLED=1 go test ./...` |
| Forbidden | Fail-soft skip; new `011_*` mig; daemon/HTTP/embeddings; weakening DONE=PASS\|escape; Gate C pack rewrite; board spawn rights |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/retrieval/exact.go` | `case "review":` → `GetReview` → Hit |
| `internal/retrieval/types.go` | Add `ReasonReviewJudgesTask` / `ReasonReviewJudgesScope` consts (`"review_judges_task"` / `"review_judges_scope"`) |
| `internal/retrieval/expand.go` | Map those rels in `hitFromLinkNeighbor` switch |
| `internal/retrieval/*_test.go` (and/or `internal/compiler` test) | Regression: plant review like honesty; assert Why + context/Expand succeed and include review |
| Optional comment | `ExactQuery` EntityType comment may list `review` |

## Role work

1. TDD: add failing regression that plants a linked PASS review then calls `Why` on the task (expect error today).
2. Implement `case "review"` Hit mapping; wire link reason consts.
3. Extend assertion: Expand/TaskContext includes review neighbor; Exact by id for `"review"` works.
4. Run locked verify suite; board **status + Notes only**.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/domain/... ./internal/store/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-check: after plant, library `Why("task", id)` must **not** return `unknown entity type "review"`.

## Exit criteria
- [ ] `lookupEntity` handles `"review"` via `GetReview`; Hit title + result/body excerpt as locked
- [ ] `hitFromLinkNeighbor` maps `review_judges_task` / `review_judges_scope` reason codes
- [ ] Regression proves Why + context/Expand succeed with linked review (no fail-soft skip)
- [ ] No new migration; DONE/review policy unchanged
- [ ] Carry-forward: honesty Gate G + A/B/C, p0x, x0, `./...` green
- [ ] Board Notes ready for **P09-S01-02**

## Out of scope
- S02 discoverability (`trace tasks` / seed path)
- S03 install-wire / MCP json writer
- FTS ranking changes (SyncEntityFTS already indexes reviews)
- Looking up `plan_scope` / capability entity types (only if blocking DF-01 — Note and leave for spawn/review)
- Rewriting Gate C Mode-B packs or Phase 00–08 history

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] Add failing Why/context-with-review regression (honesty/D07 plant shape)
- [ ] Implement `case "review"` in `lookupEntity` + link reason consts/mapping
- [ ] Assert Exact review-by-id + Expand/TaskContext include review neighbor
- [ ] Run locked verify suite; mark P09-S01-01 done with Notes
