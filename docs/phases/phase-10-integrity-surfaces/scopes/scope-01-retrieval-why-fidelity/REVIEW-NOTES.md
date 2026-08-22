# P10-S01-02 — REVIEW-NOTES (retrieval / why fidelity)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | DF-19 no global all-project DPC | **Pass** — `discoveryPlanChangeHitsForTask` + Expand post-pass only (`expand.go`); global `discoveryPlanChangeHits` gone |
| 2 | DF-19 algorithm vs locks | **Pass** — pair-completion → goal foreign filter → single-goal unattributed; nil `goal_id` → pair-completion only. Tests: `TestWhyTaskDPCGoalScoped`, `TestWhyTaskDPCMultiGoalNoForeignPollution`, `TestTaskContextDPCGoalScoped`, `TestTaskContextMultiGoalOmitsForeignDPC` |
| 3 | DF-23 `plan-change` → `plan_change` | **Pass** — `NormalizeEntityType` in Exact/Why/`lookupEntity`; `TestExactWhyPlanChangeAlias` emits `plan_change` |
| 4 | DF-25 capability Exact/Why | **Pass** — `GetCapability`; Title/status Excerpt; `TestExactWhyCapability` |
| 5 | DF-25 `plan_scope` residual | **Pass** — still unknown type; Notes + this residual list; not implemented |
| 6 | DF-27 MD Law 9/4 | **Pass** — decision/assumption title reworded; JSON `trust` stays `untrusted_data`; `TestDecisionMarkdownTrustLabels`; no TrustSystem elevation for bodies |
| 7 | DF-29 IncludeWhy fail-closed | **Pass** — `compileAtDepth` returns Why err; `TestIncludeWhyFailClosed` (TaskContext + ExpandContext); `Retriever` seam |
| 8 | No new mig / MCP tools / daemon | **Pass** — schema still through `010_*`; MCP still 6 tools (`trace_why`…`trace_review`) |
| 9 | Carry-forward | **Pass** — honesty; E/F/ablation/H; compat; p0x; x0; product pkgs green. Gate C `dry_run:false` N=3 mean G1 0.800 > B0 0.000 intact. `go test ./...` still fails only on pre-existing `similar projects/graphify` path space (non-product) |
| 10 | Board Notes / S02 Depends | **Pass** — P10-S01-01 Notes cite final test names; S02 stubs inherit alias + IncludeWhy |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Pair-completion branch has no dedicated unit test (covered only by algorithm review) | Residual — acceptable; S05 may spot-check if desired |
| low | Assumption MD uses “recorded user decision” phrasing (planner-allowed shape) | Residual — wording OK vs locks |
| nit | `go test ./...` FAIL on `similar projects/graphify` space-in-path | Pre-existing; out of S01 |

## Residuals (explicit)

1. **`plan_scope` Exact/Why** — still unknown (out of S01; intentional).
2. **Mode-B Gate C packs** — historical; not rewritten; scores unchanged.
3. **S02** owns MCP tool parity / install freshness / capability JSON (DF-21/22/32).
4. **Shared multi-goal DPC** (linked to tasks under both A and B) omitted from both goals under strict foreign rule — matches locks; no spawn.
5. **Pair-completion** untested in isolation (low).
6. **`./...`** non-product path FAIL (pre-existing).

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1  → PASS
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... -count=1  → PASS
CGO_ENABLED=1 go test ./... -count=1  → product pkgs PASS; FAIL only similar projects/graphify (space)
```

## Next

**P10-S02-00** (no spawn inserted).
