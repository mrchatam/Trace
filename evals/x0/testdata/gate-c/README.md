# Gate C answer packs (Mode B)

Recorded operator-simulated packs for Experiment X0 Gate C.

## Protocol

- **Model pin:** `recorded-operator-sim/v1` (not a live LLM; refresh packs via live agent under `TRACE_X0_GATE_C=1` when available).
- **B0:** repo tools on `src/**` only — no `trace why` / `trace context`.
- **G1:** `trace why` + `trace context` on the wiring task (+ repo OK under the same oracle rule).
- **Oracle policy:** B0/G1 agents **must not** treat `fixtures/x0/seed/gt.json`, `evals/x0/GT-MAP.md`, or any UUID/causal answer key as the answer oracle. The harness may still read `seed/gt.json` for import and grading. Agent-facing `fixtures/x0/README.md` deliberately omits the UUID table.
- **Grading:** structured `assert` tokens vs `queries.json` GT keys (free-text is narrative only).

## Files

| File | Condition | Run |
|------|-----------|-----|
| `b0-run1.json` … `b0-run3.json` | B0 | 3 |
| `g1-run1.json` … `g1-run3.json` | G1 | 3 |

Live refresh: re-run agents with the same pins, rewrite packs, then `CGO_ENABLED=1 go test ./evals/x0/ -run TestX0GateCRecordedMetrics`.
