# P04-S02-02 — Scope review notes (honesty escape-rate / Gate G prelim)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P04-S02-01 against S02-00 / `01-honesty-escape-rate.md` locks. Claims match live repo: named `TestHonestyEscapeRateGateGPrelim` under `evals/honesty` (no new package); committed `schema-gate-g.json` v1; temp `metrics-gate-g.json` written + jsonschema-validated; planted tallies escapes=1 / caught=2 / attempts=3 / escape_rate≈1/3; hatch only on Escape-1; `TestHonestyFailClosedPlantedClaim` Paths A/B/C intact and never set `AllowDoneWithoutReview:true`; S01 hooks exercised (`LinkReviewScope` / `review_judges_scope`, OPEN `POLICY_EXCEPTION`, `CountOpenResidualsByScope==1`); no product Go outside harness; VerifiedFact absent; Gate E + Gate C `dry_run:false` N=3 G1 0.800>B0 0.000 + p0x/x0/`./...` green. S03 VERIFY Depends already name this test + schema path. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Named test under `evals/honesty` (not new package) | **Pass** — `TestHonestyEscapeRateGateGPrelim` in `evals/honesty/honesty_test.go`; tree is `doc.go` + test + schema only |
| `schema-gate-g.json` v1 + temp `metrics-gate-g.json` validated | **Pass** — schema `schema_version` const 1; write+`validateGateGMetricsFile` before return |
| Escape formula escapes=1, caught=2, attempts=3; hatch=escape only | **Pass** — Caught-1 EvidenceIDs-alone; Caught-2 FAIL review; Escape-1 new task + hatch; PASS→DONE excluded from attempts |
| Paths A/B/C intact; never hatch in A/B/C | **Pass** — separate `TestHonestyFailClosedPlantedClaim`; all three transitions `AllowDoneWithoutReview: false` |
| S01 hooks: LinkReviewScope / CountOpen / OPEN POLICY_EXCEPTION | **Pass** — `runS01ResidualTally` + metrics `s01_hooks` array |
| No product Go outside harness; no daemon/HTTP/embeddings; VerifiedFact absent | **Pass** — S02-01 touched only `evals/honesty/*`; VerifiedFact only in docs/comments |
| Gate E / Gate C / p0x / x0 / `./...` | **Pass** — fresh suites below; Gate C artifacts not rewritten |
| S03 Depends name test + schema | **Pass** — `00-PLANNER` / `01-verify` / `SCOPE-TODOS` cite `TestHonestyEscapeRateGateGPrelim` + `schema-gate-g.json` / temp metrics |

## Claims → evidence

| Claim (P04-S02-01 Notes) | Evidence |
|--------------------------|----------|
| `TestHonestyEscapeRateGateGPrelim` | `evals/honesty/honesty_test.go` |
| `schema-gate-g.json` v1 | `evals/honesty/schema-gate-g.json` (`const: 1`) |
| Temp `metrics-gate-g.json` | Written under `t.TempDir()`; validated via `jsonschema` Draft2020 |
| escapes=1 caught=2 attempts=3 | Asserted in test + metrics fields |
| OPEN `POLICY_EXCEPTION` + `CountOpenResidualsByScope` | `runS01ResidualTally` + `ResidualCodePolicyException` / severity WARN |
| A/B/C untouched | `TestHonestyFailClosedPlantedClaim` body unchanged in behavior; hatch only in Gate G Escape-1 |
| Gate C untouched | `docs/verification/gate-c-x0/metrics-{b0,g1}.json` still `dry_run:false` N=3 means 0.000 / 0.800 |
| CGO bars | Fresh re-verify below |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... -count=1
  → PASS (honesty, replan / Gate E)

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./... -count=1
  → PASS (p0x, x0, honesty, replan, domain, store, planner, mcp, analyzers, …)
```

Gate C spot-check: `dry_run:false`; N=3; mean G1 0.800 > B0 0.000 — packs not rewritten.

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._

### low

1. **Schema does not require specific `s01_hooks` strings** — only `minItems: 1`. The planted metrics include the three locked tokens; a future metrics writer could omit them and still validate. Acceptable for prelim; VERIFY can assert contents if desired.

### nit

1. Escape-rate range check nests `(0.33,0.34)` with an exact `1.0/3.0` fallback — redundant but correct for `1/3`.
2. Gate G planted Caught paths do not create the Claim/Evidence link used in A/B/C narrative — intentional mirror of reject behavior only.

## Spawns
_None._

## Next board row
**P04-S03-00** (Phase 04 VERIFY planner). Do not invent full Gate G production without this harness bar.
