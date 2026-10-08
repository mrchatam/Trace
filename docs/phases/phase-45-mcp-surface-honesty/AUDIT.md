# Phase 45 — Post-implementation audit (2026-10-04)

A full-surface audit pass ran after S01–S06 landed: mechanical checks (`go vet`, `gofmt`), an #155 completeness probe (bogus-entity sweep over **every** tool, not just the four named in the issue), an #154-class sweep (authoritative-empty results for nonexistent parents), and a response-shape review. Findings and dispositions:

## Fixed during the audit

| # | Finding | Fix |
|---|---------|-----|
| A1 | `internal/store/capability.go` `GetCapabilityBySlug` was the one unconverted getter — `trace_capability require` leaked raw `sql: no rows in result set` (#155) | Converted to `notFound("capability", slug)` |
| A2 | `trace_capability missing` returned an authoritative empty list for a nonexistent task (#154 class) | `GetTask` precondition + `wrapToolErr` |
| A3 | `trace_impact finding list` returned an authoritative empty list for a nonexistent decision (#154 class) | `GetDecision` precondition + `wrapToolErr` |
| A4 | `trace_tasks goal_id=<nonexistent>` returned an authoritative empty list (#154 class) | `GetGoal` precondition; `TestTraceTasksParity` updated to assert the fail-closed behaviour (its old assertion **enshrined** the dishonest empty result) |
| A5 | `trace_changes list task_id=<nonexistent>` / `trace_regressions list task_id=<nonexistent>` — same class | `GetTask` preconditions + `wrapToolErr` |
| A6 | `trace_review create` with a bogus `task_id` **created the review row, then failed the link** — orphaned entity + error return ⇒ caller retries and duplicates | Task link target validated **before** the write |
| A7 | `toolChanges` opened the store once, then each action handler opened it **again** (the S01 double-open defect, unfixed here) | Single open, store passed to the action handlers |
| A8 | Same double-open in `toolRegressions` | Fixed identically |
| A9 | `trace_changes list` / `trace_regressions list` hardcoded `truncated:false` — a claim the store cannot back | `truncated` is now the page-fill signal (`len(items) == effective limit`), and the effective `limit` is reported |
| A10 | `resolveExploreBounds` errors said "must be 1..64" while the code **clamps** over-max values — the message lied about the behaviour | Negative values still rejected; message no longer claims over-max rejection; clamps mirror `compiler.ExploreOpts` normalization |
| A11 | `descVersion` omitted the `project_root` field it returns (#159 class) | Description regenerated |
| A12 | `descAgents` hardcoded `list|describe|recommend` outside the shared vocabulary (#159 class) | Generated from `agentActions` |
| A13 | `internal/mcp/schema.go` had a duplicated `reviewResultNames` vocabulary — the exact drift this phase exists to remove | Deduplicated |
| A14 | `trace_agents recommend` errors carried the library prefix (`agents recommend: …`) instead of the tool name | Wrapped as `trace_agents recommend: …` |

## Guard added

`TestToolErrorsNeverLeakDriverText` grew from 5 to **14** cases: every tool that can fail on a missing/invalid entity is now swept, including `trace_context`, `trace_capability missing|require`, `trace_plan deep`, `trace_impact finding list|walk`, `trace_loop status`, `trace_regressions list`, `trace_changes compare`. A future unconverted getter or unvalidated filter fails this test instead of reaching a model. (`trace_loop gate` is deliberately asserted via its envelope, not an error — the `allowed:false` gate envelope *is* the honest answer.)

## Accepted residuals

- **`trace_capability list` / `missing` / `trace_impact finding list` keep domain-named keys** (`capabilities`, `missing`, `findings`) rather than `items`. Deliberate: the MCP responses mirror the CLI's output documents byte-for-byte (`cmd/trace/capability.go:144`), and MCP tools mirror CLI commands in this repo. `listEnvelope` is the standard for *generic history* lists; entity-specific collections keep their domain vocabulary. Renaming would need a cross-surface (CLI+MCP+HTTP) decision.
- **`recommend` hint is static.** When `task_id` is given, `EmptyCatalogHint` (about `trace install agents`) is emitted even when the catalog is populated — the empty result came from the task, not the catalog. Distinguishing the two would need a reason field on the recommendation; noted for a future round.
- **`trace_changes compare` returns the raw domain result** (no envelope). It is a diff document, not a list; envelope-wrapping it would add noise, not honesty.
- **`implication for replace=true` unchanged:** it still overwrites the entity row only, no cascade (per the locked D2).

## Verification after fixes

```
gofmt -l internal/ cmd/                     → clean
go vet ./...                                → clean
CGO_ENABLED=1 go build ./...                → ok
CGO_ENABLED=1 go build -o bin/trace-mcp ./cmd/trace-mcp   → ok
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace          → ok
CGO_ENABLED=1 go test ./... -count=1        → ALL PACKAGES PASS (exit 0)
```

Notably the full run includes `evals/perf` and the two port-race tests that had failed under machine load in the prior session — on an idle machine they pass, confirming those were load-induced, not code defects. S01–S06 named tests re-run green; protocol conformance + key-shape guards green.
