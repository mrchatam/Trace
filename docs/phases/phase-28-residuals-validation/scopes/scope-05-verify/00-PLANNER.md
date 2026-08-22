# P28-S05-00 — Scope planner (full regression VERIFY)

## Metadata
- id: P28-S05-00
- todo_ids: [P28-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: mixed
- hooks: []

## Objective

Lock Phase 28 VERIFY floor and DR-HANDOFF successor rules. Finalize `01-verify.md` + `02-dr-handoff.md` after S00–S04 complete.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 28 README](../../README.md)
- [Phase 27 VERIFY-NOTES](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md) — baseline
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify floor

| Block | Command / check |
|-------|-----------------|
| Build | `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` |
| Unit | `go test ./internal/... -count=1` |
| Cmd | `go test ./cmd/trace/... -count=1` |
| Install | `go test ./internal/install/... -count=1` |
| Matrix | `evals/p28-regression/` if S01 added it |
| Score build | `./score.sh G1 --p25 --arm build` |
| Score directed | `./score.sh G1 --p25 --arm directed` (+ attestation if Session-B ran) |
| Hook smoke | strict deny without task ID (test or script ref) |

## Evidence dir

`experiments/runs/YYYY-MM-DD-p28-s05-01-verify/evidence/`

## Successor defaults (S05-02)

| Outcome | Successor |
|---------|-----------|
| R1–R5 closed, regression green | `no successor` |
| Session-B still FAIL after S02 | Optional Phase 29 harness-only (human promote) |
| Hook regression | Spawn repair row; do not close phase |

Successor decision: **never TBD**.

## Planner gate

- [ ] S00–S04 rows done
- [ ] `01-verify.md` + `02-dr-handoff.md` runnable
- [ ] Residual register R1–R8 final disposition template in verify notes

## Exit criteria

- [ ] Verify prompt locked for fresh subagent
- [ ] Board row P28-S05-00 Notes cite verify floor
- [ ] Next runnable **P28-S05-01**

## Todo updates

Status + notes on **P28-S05-00** only.

## Next

`P28-S05-01`
