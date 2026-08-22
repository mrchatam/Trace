# P10-S04-02 — REVIEW-NOTES (operator / capability gates)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-17** →DONE = hatch **or** (PASS ∧ `AllowOperatorDone`); Actor alone insufficient | **Pass** — `TransitionTask` (`task_state.go`): hatch bypasses PASS+operator; else requires `findPassReviewID` + `AllowOperatorDone`. `TestOperatorDoneRequiresFlag` rejects `Actor:"operator"` without flag; accepts `Actor:"agent"` + flag |
| 2 | **DF-17** CLI `--as-operator` + MCP `as_operator` (G19) | **Pass** — `cmd/trace/transition.go` → `AllowOperatorDone`; `tools_write.go` `AsOperator` → domain; CLI `TestDoneRequiresReviewPass` path asserts deny-without / allow-with `--as-operator`; MCP transition test deny-without / allow-with `AsOperator` |
| 3 | **DF-17** honesty Path C operator flag; A/B reject | **Pass** — `honesty_test.go` Path C sets `AllowOperatorDone: true`; A/B still reject; reject reason asserts mention `AllowOperatorDone`/`as-operator` |
| 4 | **DF-18** leave DONE invalidates PASS → `UNCERTAIN` | **Pass** — `invalidatePassReviewsOnReopen` via `SetReviewResult` → `UNCERTAIN`; `TestReopenInvalidatesPassReviews` (sticky PASS gone; new PASS + operator OK) |
| 5 | **DF-24** missing caps fail-closed; override flags | **Pass** — gate before mutate; `TestMissingCapabilitiesBlockTransition`; CLI `--allow-missing-caps` / MCP `allow_missing_caps` wired on `TransitionOptions` |
| 6 | **DF-26** hatch loud WARNING; Gate G hatch works | **Pass** — CLI stderr `WARNING` + `allow-done` (`TestAllowDoneWarnsOnStderr`); MCP `"warning"` (`TestTransitionAllowDoneEmitsWarning`); `TestOperatorDoneHatchBypassesOperator` + Gate G Escape-1 still hatch-only success |
| 7 | **DF-31** capability missing without task usable | **Pass** — CLI usage + `trace tasks` hint (`TestCapabilityMissingRequiresTaskHint`); MCP `task/task_id` + `trace_tasks` (`TestCapabilityMissingRequiresTaskParam`) |
| 8 | No new mig; no Actor-name allowlist; nine MCP tools | **Pass** — schema stops at `010_capability_surface.sql` (no `011_*`); DONE auth never branches on Actor string; `server.go` / `TestToolNamesRegistered` nine tools |
| 9 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + S01–S03 + Gate C | **Pass** — locked CGO0/CGO1 suites green; Gate C `docs/verification/gate-c-x0/` `dry_run:false` N=3; G1 mean understanding_accuracy **0.800** > B0 **0.000** intact |
| 10 | Board Notes accurate; planner no product Go | **Pass** — P10-S04-00 Notes claim no product Go; P10-S04-01 Notes cite live tests accurately |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | No dedicated CLI/MCP test asserting `--allow-missing-caps` / `allow_missing_caps` end-to-end (domain override covered; adapters wire flags) | Residual — wiring present; S05 spot-check OK |
| nit | Cursor IDE MCP `trace_transition` schema may lag until reload (`as_operator` / `allow_missing_caps` in source) | Residual OK (S02 carry-forward) |

## Residuals (explicit)

1. **ab-operator-gate** experiment re-run optional (unit/MCP tests cover probe shape).
2. **Cursor MCP reload** still manual — live IDE tool schema may omit new params until reload.
3. **`go test ./...`** FAIL only on pre-existing `similar projects/graphify` path space (non-product).
4. CLI/MCP missing-caps override lacks a dedicated adapter test (low).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./... -count=1  → product pkgs PASS; FAIL only similar projects/graphify (space)
Gate C: dry_run:false N=3; G1 mean understanding_accuracy 0.800 > B0 0.000 intact
```

## Next

**P10-S05-00** (no spawn inserted).
