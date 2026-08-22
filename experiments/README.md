# Trace dogfood experiments

Live **Cursor agent** A/B runs: same vague product ask, **B0** (no Trace) vs **G1** (Trace + optional Phase 23 enforcement).

## Invariant

| Piece | Rule |
|-------|------|
| Product tree | `ab-*/project/` — README must **not** mention Trace |
| Arms | `runs/B0` vs `runs/G1` after `./prepare.sh` |
| Workspace | Agent opens **`runs/B0` or `runs/G1`** as Cursor root — never the Trace monorepo |
| Scoring | `./score.sh` per arm; rows in [RESULTS.md](RESULTS.md) |
| Ladder | [LADDER.md](LADDER.md) — add rungs as gaps appear |

Harness proofs live in `evals/*`. **Harness green ≠ dogfood green.**

## Active experiment

| ID | Path | Status |
|----|------|--------|
| **E03** | [ab-library-hold-desk](ab-library-hold-desk/) | **done** — P25-3a PASS (disc=2 dec=5); Session-B P25-3b PASS (disc=4 dec=6, +1 promoted task) |
| **E02** | [ab-p25-gap-pass-validation](ab-p25-gap-pass-validation/) | **done** — Session-B + Phase 28 verify |
| **E01** | [ab-incident-tracker](ab-incident-tracker/) | **done** — two-mode model; fed Phase 24 |

## Quick start (E03)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
cd experiments/ab-library-hold-desk
./prepare.sh
```

Then **File → Open Folder** → `runs/B0` or `runs/G1` and paste the matching prompt from `prompts/`.

See [E02](ab-p25-gap-pass-validation/) for the checkout-desk P25 validation that closed Phase 25–28 residuals.
