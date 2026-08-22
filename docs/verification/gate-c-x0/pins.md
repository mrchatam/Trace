# Gate C X0 pins

| Pin | Value |
|-----|--------|
| Date (UTC) | 2026-08-16 |
| Operator / session | P02-S01-01 implementer (Mode B recorded packs); seed pin stabilized P02-S01-02; fixture README oracle strip P02-S02-01 |
| Model | `recorded-operator-sim/v1` |
| Trace version | `0.0.0-dev` |
| Fixture | `fixtures/x0` |
| Seed pin (metrics) | `fixtures/x0/seed/gt.json` (repo-relative; not temp worktree abs) |
| Fixture content hash (`find fixtures/x0 -type f ! -path '*/.git/*' \| sort \| xargs sha256sum \| sha256sum`) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Query bank | `evals/x0/queries.json` (5 understanding questions, same order B0/G1) |
| Answer packs | `evals/x0/testdata/gate-c/{b0,g1}-run{1,2,3}.json` |
| Metrics | `metrics-b0.json`, `metrics-g1.json` (`dry_run: false`, N=3) |
| RNG / seed | none (deterministic recorded packs) |

## Agent brief (fairness)

- **B0:** read/search under `src/**` only; **must not** call `trace why` / `trace context`; **must not** treat `seed/gt.json`, `evals/x0/GT-MAP.md`, or fixture README as the answer oracle (agent-facing README has no UUID table).
- **G1:** may call `trace why` / `trace context` (CLI) and read `src/**` under the same oracle rule.
- Grading uses structured `assert` tokens vs GT keys (not naive free-text substring). Harness may read seed for import/grading only.

## Live refresh

Optional: set `TRACE_X0_GATE_C=1` and replace packs from a live agent with the same pins, then re-run:

```bash
CGO_ENABLED=1 go test ./evals/x0/ -run TestX0GateCRecordedMetrics -count=1
```
