# P09-S01-02 — Scope review notes (2026-08-16)

Independent review of DF-01 retrieval-review fix vs `00-PLANNER.md` / `01-retrieval-review.md` locks + `P09-S01-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `lookupEntity` `case "review"` → `store.GetReview` | Pass (`exact.go` ~162–171) |
| Hit: Title; Excerpt = Result if set else `excerpt(Body)`; ReasonCode from caller | Pass |
| Expand maps `review_judges_task` / `review_judges_scope` | Pass (`types.go` consts + `hitFromLinkNeighbor`) |
| No fail-soft skip of unknown `"review"` (hard error for other unknown types remains) | Pass — default still `unknown entity type`; review is a real case |
| No new migration (`011_*`) | Pass — schema still through `010_capability_surface.sql` |
| Regression `TestWhyAndContextWithLinkedReview` | Pass — Exact + Why + Expand with linked PASS review |
| Why on linked task succeeds; Expand includes review neighbor | Pass |
| DONE/review promotion policy untouched | Pass — no edits in `task_state.go` / DONE gate; domain review APIs unchanged |
| No CLI/MCP scope creep | Pass — product changes confined to `internal/retrieval` (+ test) |
| DF-01 / DF-09 multi-step why/context after review | Pass — plant matches honesty/D07 shape |
| S02 discoverability stubs compatible | Pass — listing tasks does not depend on skipping reviews |
| Carry-forward: honesty A/B/C + Gate G, p0x, x0, `./...` | Pass (independent re-run) |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1  PASS
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1                 PASS
CGO_ENABLED=1 go test ./... -count=1                                                               PASS
```

## Findings

None at blocker / high / medium.

### Low (no spawn)

- `ReasonReviewJudgesScope` is mapped but not exercised by the DF-01 regression (task-link path only). Acceptable: DF-01 bug was task why/context after review; `plan_scope` ExactLookup remains out of S01 scope by planner.
- Empty-Result → body-excerpt branch is implemented but not separately asserted (PASS Result path is covered).

### Nit

- Compiler `TaskContext` not named in the regression; Why uses Expand and TaskContext shares Expand — coverage is sufficient per locks (“TaskContext **or** Expand”).

## Spawns

None.

## Residuals

- Expanding a `review` that only judges a `plan_scope` may still hard-fail on unknown entity type `"plan_scope"` until a later row adds that lookup (explicitly out of S01).
- Pre-existing Expand `isNotFound` soft-continue on missing link targets remains (not DF-01; not a skip of unknown `"review"`).

## Next

**P09-S02-00** (discoverability scope planner). Do not start until orchestrator launches that row.
