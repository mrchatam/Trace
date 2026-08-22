# E02 Rubric — P25 validation

## Product (both arms)

- `go test ./...` PASS (`--test`)
- Checkout desk features: assets, checkout/return, roles, API, admin UI, public availability, audit

## B0 harness

| ID | Pass |
|----|------|
| B0-1 | No `.trace/` |
| B0-2 | Deliverable `cmd/checkoutd` or equivalent |

## G1 harness (Phase 23 baseline)

| ID | Pass |
|----|------|
| G2 | `trace/graph.json` exported |
| G3 | ≥1 goal, ≥3 tasks in graph |
| E1–E3 | cursor rules, hook, strict config |

## G1 P25-specific (`--p25`)

| ID | Arm | Pass | Maps to |
|----|-----|------|---------|
| P25-1 | both | Installed rules contain `mandatory gap pass` | INT-03 wiring |
| P25-2 | both | Installed rules contain `Parent orchestrator` | INT-04 docs |
| P25-3a | build | Graph: **discoveries ≥ 1** OR **decisions ≥ 1** after build-only session | H-P25-1 baseline |
| P25-3b | directed | Same threshold — **PASS required** after directed gap pass | P25-C validation |
| P25-4 | per arm | Env attestation before score: `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` arm-matched → pass; unset → skip; also record in RESULTS.md | Protocol |
| P25-5 | both | **Optional fail:** new tasks > 5 seed (would indicate P25-A not needed yet) | FM-01 |

Score with `./score.sh G1 --p25 --arm build` (Session A) or `./score.sh G1 --p25 --arm directed` (Session B). `./score.sh G1 --p25` defaults to `--arm build`.

## Expected failures (do not blame P25-C)

| Check | Expected on build-only (P25-3a) |
|-------|-----------------------------------|
| P25-3a graph richness | **FAIL** when `discoveries=0 decisions=0` — documented baseline, not install regression |
| Verify task gate allows edit | May still FAIL (`hop_budget_exceeded`) |
| New tasks from discoveries | May still be 0 |
| Uncertainties resolved | May still be 0 |

## Verdict matrix → Phase 26

| P25-1/2 | P25-3a (build) | P25-3b (directed) | Verdict | Phase 26 |
|---------|----------------|-------------------|---------|----------|
| PASS | FAIL (thin graph) | — | **Install OK; behavior unchanged** | P25-A and/or stronger INT-04 hook |
| PASS | — | PASS | **P25-C validated (directed)** | Optional: P25-A/B by priority, or `no successor` + more dogfood |
| PASS | PASS | — | Build arm rich (unusual) | Investigate agent behavior |
| FAIL | — | — | Install regression | Fix Phase 25 wiring first |
| PASS | PASS but G1≡B0 | — | Harness works; product parity | Protocol INT-08; not blocking P25-C |

## E01 comparison row

Record in RESULTS.md: E01 Session A vs E02 G1 on discoveries/decisions/tasks without human gap prompt.

## RESULTS.md attestation template

```text
## E02 G1 Session A (build arm)
- P25-4 build: no human gap prompt before score — attested Y/N
- P25-3a outcome: PASS/FAIL (discoveries=… decisions=…)

## E02 G1 Session B (directed arm)
- P25-4 directed: gap prompt sent — attested Y/N
- P25-3b outcome: PASS/FAIL (discoveries=… decisions=…)
```
