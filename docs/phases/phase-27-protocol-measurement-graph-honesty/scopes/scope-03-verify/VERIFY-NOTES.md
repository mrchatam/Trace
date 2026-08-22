# Phase 27 VERIFY notes

**Run:** 2026-08-20  
**Row:** P27-S03-01  
**Evidence dir:** experiments/runs/2026-08-20-p27-s03-01-verify/evidence/  
**Git SHA:** `unknown` (workspace has no `.git`)  
**Harness:** score.sh T02 `--strict --enforce` (both arms)

## Verdict

PASS — confidence **high** (honesty/enforce upgrade + P25-1/2 hold; expected thin residuals only)

## Closure signals (Phase 27)

| Check | Target | This run |
|-------|--------|----------|
| INT-07 product honesty | Thin enforce blocked | **PASS** |
| INT-08 score enforce | G2 FAIL on thin (not WARN-only) | **PASS** |
| INT-10 arms | P25-3a/3b labels + arm isolation | **PASS** |
| Unit + cmd tests | PASS | **PASS** |
| P25-1 / P25-2 | PASS | **PASS** |

## Harness score (--p25)

| ID / arm | Result | Notes |
|----------|--------|-------|
| P25-1 | PASS | GapPassPrompt in installed rules (build + directed) |
| P25-2 | PASS | Parent orchestrator in installed rules (build + directed) |
| P25-3a (build) | FAIL expected | discoveries=0 decisions=0 — build-only thin baseline |
| P25-3b (directed) | FAIL OK | No Session-B dogfood; label `P25-3b` present |
| G2 enforce | FAIL expected | WARN→FAIL after T02 `--strict --enforce`; thin honesty line |
| FM-07 | PASS | export SHA matches HEAD (warn-only path unused) |

## Commands (this run)

```bash
# Block 0 — score.sh T02 upgraded to --strict --enforce (both arms)
# Block 1
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
go test ./internal/... -count=1
go test ./cmd/trace/... -run 'SeedExport|Enforce|Strict' -count=1
# Block 2 — thin demo (P26 snippet)
$TRACE_BIN -C "$THIN_WS" seed import cmd/trace/testdata/p26-export-snippet.json
$TRACE_BIN -C "$THIN_WS" seed export -o /tmp/p27-thin-out.json --strict --enforce  # exit 1, no write
$TRACE_BIN -C "$THIN_WS" seed export -o /tmp/p27-thin-warn.json --strict           # exit 0, file written
# Block 3
./prepare.sh G1
./score.sh G1 --p25 --arm build
./score.sh G1 --p25 --arm directed
./score.sh G1 --p25 --arm foo   # exit 2
```

## Residuals (non-blocking)

| Topic | Disposition |
|-------|-------------|
| Build thin P25-3a FAIL | Expected measurement baseline |
| Thin G2 enforce FAIL | Expected after S03-01 |
| Directed without Session-B | P25-3b FAIL OK |
| P25-4 attestation | Manual / skip |
| S02-02 BLOCKING duplicate msg | Note only |
| Overall VERDICT FAIL (thin) | Acceptable — closure signals PASS |
| Gate lines on fresh seed (`plan_missing` / `plan_uncritiqued`) | Accompany honesty FAIL; do not mask thin residual |

## Gaps / spawn

(none)

## DR-HANDOFF status

**OPEN** — S03-02 closes with successor decision (never TBD).
