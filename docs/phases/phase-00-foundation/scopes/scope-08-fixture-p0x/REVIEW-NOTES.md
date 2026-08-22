# P00-S08-02 — Review notes (Fixture + P0-X harness)

**Date:** 2026-08-15  
**Verdict:** APPROVE  
**Confidence:** high  
**Spawns:** none

## Claims vs evidence

| Claim (P00-S08-01 Notes / 01-fixture.md) | Evidence |
|------------------------------------------|----------|
| `fixtures/x0` Apache-2.0 synthetic | `fixtures/x0/LICENSE` (Apache 2.0 full text); README marks synthetic / DR-BENCH |
| TS + Python sources, non-vacuous | `src/greeter.ts` (`greet` + `./format` import); `src/math_util.py` (`add`/`hypotenuse` + `math`) |
| Human GT seed v1 stable UUIDs | `seed/gt.json` version 1; IDs `1111…`/`2222…`/`3333…`/`4444…`/`5555…` match README + harness consts |
| Abs seed path (not under `-C`) | CLI `os.ReadFile(path)` with raw argv; harness `filepath.Abs` + `TestSeedPathIsAbsoluteNotUnderC` |
| 7/7 criteria + ≥5 queries | `TestP0XAllCriteria` subtests `criterion-1`…`7` + five named query subtests |
| Incremental #7 sibling isolation | Mutate TS only → `index src/greeter.ts` → Py fingerprint + content_hash unchanged; TS gains `greetAgain` |
| Metrics schema in temp | `metrics-p0x.json` with `ok` + `criteria["1"…"7"]` all true |
| No committed `.trace/` | No `.trace` under `fixtures/`/`evals/`; root + fixture `.gitignore` list `.trace/` |
| No MCP/daemon/`evals/x0` | Tree is `evals/p0x` only; harness is CLI + store asserts |

## Independent verification (2026-08-15)

```text
CGO_ENABLED=1 go test ./evals/p0x/... -count=1   # PASS (incl. all 7 + seed-abs)
CGO_ENABLED=1 go test ./... -count=1             # PASS
```

## Findings

| Severity | Location | Failure mode | Disposition |
|----------|----------|--------------|-------------|
| low | `evals/p0x/p0x_test.go` `decision-constraint` | Soft OR fallback to narrative substring `"typescript greeter"` could pass without DecisionID / `decision_affects_task` if why/context ever regress | Residual — primary path still checks DecisionID; not vacuous today |
| nit | `evals/p0x/p0x_test.go` why step decode | Unchecked `s.(map[string]any)` panics on malformed CLI JSON instead of `t.Fatal` | Residual |
| nit | `fixtures/x0/src/greeter.ts` | Imports `./format` with no `format.ts` (dangling import for index row) | Intentional; documented in README |

No blocker/high. No spawn.

## Cross-scope / upcoming

- S09-01 Depends already points at `CGO_ENABLED=1 go test ./evals/p0x/...` + abs seed path; lightly thickened with live UUID map + metrics schema for VERIFY.
- Next board row: **P00-S09-00** (phase VERIFY planner).

## Residuals (explicit)

1. Soft `decision-constraint` fallback (low).  
2. Panic-on-bad-JSON type asserts in harness (nit).  
3. Workspace has no `.git` here — “no committed `.trace/`” checked via tree + ignore rules, not `git ls-files`.
