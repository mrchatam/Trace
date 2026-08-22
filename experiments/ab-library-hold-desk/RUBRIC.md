# E03 Rubric — full-stack regression

## Product (both arms)

- `go test ./...` PASS (`--test`)
- Library hold features: titles/copies, hold lifecycle, roles, REST, desk UI, public waitlist/availability, audit

## B0

| ID | Pass |
|----|------|
| B0-1 | No `.trace/` |
| B0-2 | Deliverable beyond starter (`cmd/holddeskd` + `internal/` or SPEC) |

## G1 harness

| ID | Pass |
|----|------|
| G2 | Export present + `--strict --enforce` honesty |
| G3 | ≥1 goal, ≥3 tasks |
| E1–E3 | cursor rules, hook, enforce strict |
| FM-07 | Warn-only SHA drift |

## G1 P25 (`--p25`)

| ID | Arm | Pass |
|----|-----|------|
| P25-1 | both | `mandatory gap pass` in rules |
| P25-2 | both | `Parent orchestrator` in rules |
| P25-3a | build | discoveries≥1 OR decisions≥1 (**target PASS** for E03) |
| P25-3b | directed | same threshold — only if Session B run |
| P25-4 | per arm | `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` |

## Verdict → Phase 29

| P25-1/2 | P25-3a | P25-3b | Decision |
|---------|--------|--------|----------|
| PASS | PASS | — | **Stack OK** — no Phase 29 |
| PASS | FAIL | PASS | Default gap weak — optional Phase 29 |
| PASS | FAIL | FAIL / skipped | Investigate — likely Phase 29 |
| FAIL | — | — | Install regression — fix before Phase 29 theme |

## RESULTS attestation

```text
## E03 G1 Session A (build)
- P25-4 build: no human gap prompt — Y
- P25-3a: PASS/FAIL (disc=… dec=…)

## E03 G1 Session B (if run)
- P25-4 directed: Y
- P25-3b: PASS/FAIL
```
