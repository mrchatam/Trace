# P04 / S01 / 02 — Scope review (scope review layer)

## Metadata
- id: P04-S01-02
- todo_ids: [P04-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P04-S01-01 against S01-00 locks. APPROVE with evidence or spawn forward remediations. No silent weaken of honesty/p0x/Gate E/Gate C bars.

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- Claims match repo vs S01-00 / `01-scope-review-layer.md` locks:
  - mig **`008_scope_review.sql`** + `review_residuals`
  - `LinkReviewScope` / **`review_judges_scope`** (to=`plan_scope`)
  - residual Add/List/CountOpen/SetStatus vocabulary (INFO|WARN|BLOCKING; OPEN|ACKED|RESOLVED)
  - reuse CreateReview/SetReviewResult; **no** second review stack; **no** planner fork
- Task DONE policy unchanged (PASS `review_judges_task` \| `AllowDoneWithoutReview`); EvidenceIDs alone still rejected
- Honesty Paths A/B/C still fail-closed (**no** `AllowDoneWithoutReview` in honesty proof; residuals not required on A/B/C)
- Gate E (`evals/replan` `TestPlantedDiscoveryReplan`) / Gate C artifacts / p0x / x0 still green
- VerifiedFact not introduced; `plan_scopes.status` not mutated by scope reviews
- Thin CLI G19 (`--scope`, `residual add|list`); no daemon/HTTP/embeddings primary
- Light-check upcoming S02 stubs still list S01 hooks (`CountOpenResidualsByScope` / residual codes / `review_judges_scope`)

## Required re-verify commands
```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
```
Confirm Gate C `docs/verification/gate-c-x0/` metrics still `dry_run:false` (do not rewrite packs).

## Exit criteria
- [x] Verdict + confidence + REVIEW-NOTES.md
- [x] No open blocker/high without spawn
- [x] Board status + Notes; next runnable **P04-S02-00** only after APPROVE

## Minimal todos
- [x] Diff claims vs evidence
- [x] Re-run required tests
- [x] Write REVIEW-NOTES; mark done or spawn
