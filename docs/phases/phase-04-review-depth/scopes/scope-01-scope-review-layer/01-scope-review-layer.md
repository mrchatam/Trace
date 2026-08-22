# P04 / S01 / 01 — Scope review layer

## Metadata
- id: P04-S01-01
- todo_ids: [P04-S01-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement **scope-level review** against live `plan_scopes` and **structured residual tracking** on reviews, extending the Phase 01 task-review surface (reuse `CreateReview` / `SetReviewResult`; no second review stack). Keep task DONE fail-closed and honesty Paths A/B/C green. Provide count/list hooks S02 escape-rate can consume. No daemon/HTTP/embeddings. MCP not required. VerifiedFact out.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- [phase README](../../README.md)
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md) — scope review layer
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 4
- Live: `internal/domain` CreateReview/SetReviewResult/LinkReviewTask; mig 005 `reviews.result`; `internal/store` GetPlanScope / `plan_scopes` (mig 006); thin `trace review create|set`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Package | **`internal/domain`** on `*store.Store` (+ store helpers). **Do not** invent a second review package or put review policy under `internal/planner` |
| Planner | **Do not fork** — validate scopes via `store.GetPlanScope` only; no planner CRUD changes required |
| Reuse | Keep `CreateReview` / `SetReviewResult` / `GetReview` / `LinkReviewTask` (`review_judges_task`) unchanged in behavior |
| Migration | Additive embed **`008_scope_review.sql`** only (do not rewrite `001`–`007`) |
| Scope link | **`review_judges_scope`** — `entity_links` from=`review` to=`plan_scope` (entity type string **`plan_scope`**) |
| Residuals | New table **`review_residuals`** (structured; not free-text body only) |
| Scope close policy | **Recording layer only** — do **not** mutate `plan_scopes.status` or task DONE gates from scope reviews |
| Task DONE | Unchanged: PASS `review_judges_task` **or** explicit `AllowDoneWithoutReview`. Never weaken; honesty never sets escape hatch |
| VerifiedFact | **Out** — residuals are tracking hooks, not a promotion engine |
| CLI | Thin G19: extend `trace review` (`--scope` on create; `residual add\|list`). No business logic in `cmd/trace` |
| MCP | **Not** required. Do not add MCP residual tools this scope |
| CGO | Domain + store APIs must pass `CGO_ENABLED=0` |
| Carry-forward bars | Honesty Paths A/B/C; p0x 7/7; x0; replan Gate E (`TestPlantedDiscoveryReplan`); Gate C artifacts intact |
| Out | Escape-rate harness (S02); phase VERIFY (S03); VerifiedFact; daemon/HTTP/embeddings; forking planner; second review stack |

### Schema (locked)

```sql
-- 008_scope_review.sql (additive)

CREATE TABLE IF NOT EXISTS review_residuals (
    id TEXT PRIMARY KEY,
    review_id TEXT NOT NULL,
    code TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    severity TEXT NOT NULL DEFAULT 'INFO',
    status TEXT NOT NULL DEFAULT 'OPEN',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_review_residuals_review ON review_residuals(review_id);
CREATE INDEX IF NOT EXISTS idx_review_residuals_status ON review_residuals(status);
```

`review_judges_scope` needs **no** new table — use existing `entity_links` (rel locked below).

### Residual vocabulary (locked)

```text
# severity
ResidualSeverityINFO     = "INFO"
ResidualSeverityWARN     = "WARN"
ResidualSeverityBlocking = "BLOCKING"
# reject unknown severities fail-closed; empty → INFO

# status
ResidualStatusOpen     = "OPEN"
ResidualStatusAcked    = "ACKED"
ResidualStatusResolved = "RESOLVED"
# reject unknown; empty on create → OPEN

# recommended codes (string column; reject empty code)
ResidualCodeMissingEvidence  = "MISSING_EVIDENCE"
ResidualCodeOpenGap          = "OPEN_GAP"
ResidualCodeContractGap      = "CONTRACT_GAP"
ResidualCodePolicyException  = "POLICY_EXCEPTION"  # e.g. documented AllowDoneWithoutReview / escape
# Other non-empty codes OK (S02 may count by code); empty code fails validation
```

### Minimum public API (behavior locked; names may vary slightly)

```text
# Consts (domain)
EntityPlanScope        = "plan_scope"
RelReviewJudgesScope   = "review_judges_scope"

# Domain — scope review link
LinkReviewScope(ctx, reviewID, scopeID, meta LinkMeta) error
  // Validates review via GetReview; scope via store.GetPlanScope
  // Inserts entity_links from=review → to=plan_scope, rel=review_judges_scope
  // Appends entity.linked event (same pattern as LinkReviewTask)

# Domain — residuals
AddResidual(ctx, reviewID, ResidualInput) (store.ReviewResidual, error)
  // review must exist; code non-empty; severity/status defaults as above
SetResidualStatus(ctx, residualID, status string, opts ResidualStatusOptions) error
  // Actor+Reason required (mirror SetReviewResult discipline)
ListResidualsByReview(ctx, reviewID) ([]store.ReviewResidual, error)
ListResidualsByScope(ctx, scopeID) ([]store.ReviewResidual, error)
  // Residuals on reviews linked via review_judges_scope to this scope
CountOpenResidualsByScope(ctx, scopeID) (int, error)
  // Count where status=OPEN (S02 escape-rate / residual hooks)

# Store helpers (as needed)
InsertReviewResidual / GetReviewResidual / UpdateReviewResidualStatus
ListReviewResidualsByReviewID
ListReviewResidualsByScopeID  // join entity_links review_judges_scope OR domain does join
```

### Policy (locked)

```text
CreateReview → LinkReviewScope(review, plan_scope) → SetReviewResult(PASS|FAIL|UNCERTAIN)
AddResidual(review, code, …) anytime after review exists (open or closed)

Task DONE gate: UNCHANGED (review_judges_task PASS | AllowDoneWithoutReview)
plan_scopes.status: NOT mutated by this scope
AllowDoneWithoutReview: keep; honesty Paths A/B/C must NOT set it true
```

### Target tree

```text
internal/store/
  schema/008_scope_review.sql
  residuals.go            # ReviewResidual + CRUD/list helpers (or entities_causal.go sibling)
  # GetPlanScope already exists — reuse

internal/domain/
  service.go              # EntityPlanScope, RelReviewJudgesScope (+ residual consts OK here or residual.go)
  review.go               # keep CreateReview/SetReviewResult/LinkReviewTask; add LinkReviewScope
  residual.go             # AddResidual / SetResidualStatus / List* / CountOpen*
  domain_test.go          # policy tests (see Exit criteria)

cmd/trace/
  review.go               # --scope on create; residual add|list
  help.go                 # usage lines
```

### Tests (required)

- `LinkReviewScope` succeeds for existing review + plan_scope; fails on missing scope/review
- Link emits `review_judges_scope` in `entity_links`
- `AddResidual` + `ListResidualsByReview`; severity/status defaults + reject empty code / bad severity
- `ListResidualsByScope` / `CountOpenResidualsByScope` after link + OPEN residual
- `SetResidualStatus` OPEN→ACKED/RESOLVED with actor+reason; reject unknown status
- Task DONE policy regression: EvidenceIDs-alone still rejected; PASS `review_judges_task` still works; escape hatch still explicit
- Store mig 008 applied on Open (table present)
- Do **not** change honesty Paths A/B/C to require residuals

### CLI (thin G19)

```text
trace review create --title <t> [--task <id>] [--scope <plan_scope_id>] [--id <uuid>]
trace review set --id <id> --result PASS|FAIL|UNCERTAIN --reason <text> [--actor <a>]
trace review residual add --review <id> --code <CODE> [--body <text>] [--severity INFO|WARN|BLOCKING]
trace review residual list --review <id> | --scope <plan_scope_id>
```

Stdout machine-friendly JSON lines (match existing review create/set style).

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] Mig `008_scope_review.sql` + store residual helpers live
- [ ] `LinkReviewScope` (`review_judges_scope`) + residual Add/List/Count/SetStatus APIs live
- [ ] Task DONE fail-closed unchanged; VerifiedFact absent; no planner fork / second review stack
- [ ] Thin CLI `--scope` + `residual add|list` wired (G19)
- [ ] Domain tests cover link/residual/DONE regression cases above
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...` green
- [ ] Carry-forward: `CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/...` + `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./...` green
- [ ] Gate C artifacts under `docs/verification/gate-c-x0/` untouched
- [ ] TODO.md status + Notes updated (this row only)

## Minimal todos
- [ ] Store: mig 008 + ReviewResidual CRUD/list helpers
- [ ] Domain: LinkReviewScope + residual APIs + consts; DONE regression tests
- [ ] Thin CLI review `--scope` + residual subcommands + help
- [ ] Full carry-forward bars; board Notes
