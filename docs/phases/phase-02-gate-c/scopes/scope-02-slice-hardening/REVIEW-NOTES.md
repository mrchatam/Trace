# P02-S02-02 — Scope review notes (slice hardening)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P02-S02-01. Claims map to repo evidence for **GC-01** and **GC-02**. **GC-03/GC-04** remain deferred. Honesty Paths A/B/C, p0x 7/7, x0, and `./...` bars stay green. Gate C **Go** text and Mode-B pack scores were not rewritten. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| GC-01 Why/TaskContext + `discovery_causes_plan_change` | **Pass** — `ListLinksByRel` + `discoveryPlanChangeHits` on task Expand; tests `TestWhyTaskIncludesDiscoveryPlanChange`, `TestTaskContextIncludesDiscoveryPlanChange` |
| GC-01 Expand depth still 1..2; no fake task↔discovery GT edge | **Pass** — `maxExpandDepth=2`; seed links unchanged |
| GC-02 agent README no UUID oracle | **Pass** — `fixtures/x0/README.md` clean; `evals/x0/GT-MAP.md` evaluator-only |
| GC-02 pins hash + Agent brief + guard | **Pass** — hash `15fe50a1…`; `TestFixtureReadmeHasNoGTUUIDOracle` |
| Deferrals GC-03/04 | **Pass** — no live-model pack refresh; no N/variance work |
| Honesty Paths A/B/C | **Pass** — no escape hatch in proof; suite green |
| P0-X 7/7 | **Pass** |
| Scope creep | **Pass** — no daemon/HTTP/embeddings; no progressive planner; Mode-B packs still miss q3 historically |
| Gate C Go not silently altered | **Pass** — GATE-C-NOTES still Go @ 0.800 > 0.000 |

## Claims → evidence

### GC-01

| Claim | Evidence |
|-------|----------|
| `ListLinksByRel` | `internal/store/links.go` + store test |
| Shared attach helper | `internal/retrieval/discovery_plan_change.go`; hooked from task branch in `expand.go` |
| Why surfaces discovery + plan_change | `TestWhyTaskIncludesDiscoveryPlanChange` (+ Expand depth-1 assert) |
| TaskContext / why_trace | `TestTaskContextIncludesDiscoveryPlanChange` |
| Budget prefers causal over FTS | `hitPriority` ranks `ReasonDiscoveryCausesPlanChg` above `ReasonFTSMatch` (`compiler/budget.go`) |
| Depth unchanged / no fake GT edge | Expand rejects outside 1..2; fixture seed still discovery↔plan_change only |

### GC-02

| Claim | Evidence |
|-------|----------|
| Agent README stripped | `fixtures/x0/README.md` — layout/license/harness only |
| Evaluator map | `evals/x0/GT-MAP.md` |
| Oracle policy | `evals/x0/testdata/gate-c/README.md` + `docs/verification/gate-c-x0/pins.md` Agent brief |
| Guard test | `evals/x0/fixture_readme_guard_test.go` |
| Fixture hash | recomputed `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` matches pins |
| Metrics scores unchanged | `metrics-g1.json` still 0.8; g1 packs still empty q3 asserts (historical) |

### Deferred

| ID | Status |
|----|--------|
| GC-03 | Still deferred — model pin `recorded-operator-sim/v1`; packs not refreshed |
| GC-04 | Still deferred — N=3 identical grades; no significance work |

## Findings

### blocker
_None._

### high
_None._

### medium (residual — no spawn)

1. **Global task attach** — every task Expand pulls **all** `discovery_causes_plan_change` edges in the store (locked S02 approach; avoids fake task↔discovery GT edges). Fine for `fixtures/x0` / single-project DBs; multi-goal stores may over-attach. Defer scoping to Phase 03+ unless measurement demands it.

### low (residual)

1. **`seed/gt.json` still on disk** — fairness is policy + agent brief, not mechanical hide. Live operators who ignore the brief can still read UUIDs (known Gate C residual; README oracle removed).
2. **No new Gate C live re-run** — product proof is unit tests; historical Mode-B packs still document q3 miss. Correct for this harden row (do not falsify packs).

### nit
_None material._

## Spawns

None.

## Re-verify

```text
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./... -count=1  → PASS
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
  → 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22 (matches pins)
```

## Residuals for S03 / later

- Re-check GC-01 tests + GC-02 README/hash on VERIFY; do not treat pack rewrite as required.
- GC-03/04 stay deferred unless promoted with Notes.
- Global DPC attach scope (medium residual above).
- Phase 01 dry-run ≠ Gate C; Gate C **Go** already recorded — VERIFY aggregates, does not re-litigate kill.

## Next board row

**P02-S03-00** (Phase 02 VERIFY planner).
