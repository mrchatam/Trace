# S05 scope todos

| ID | Status | Prompt | Notes |
|----|--------|--------|-------|
| P28-S05-00 | done | [00-PLANNER.md](00-PLANNER.md) | Locked verify floor + dual-lane score + successor defaults |
| P28-S05-01 | pending | [01-verify.md](01-verify.md) | Full regression; write VERIFY-NOTES; DR-HANDOFF stays OPEN |
| P28-S05-02 | pending | [02-dr-handoff.md](02-dr-handoff.md) | Close handoff; successor never TBD |

## Locked defaults (S05-00)

| Item | Value |
|------|-------|
| G1 wipe | **Forbidden** — no `./prepare.sh` |
| Directed (primary) | `P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed` on rich live G1 → P25-3b PASS |
| Build live | `P25_ATTEST_BUILD=Y ./score.sh G1 --p25 --arm build` on rich G1 → P25-3a PASS labeled **post-Session-B** |
| Build thin | Docs via `SESSION-A-GRAPH-SNAPSHOT.json` (disc=0/dec=0); do not recreate |
| Hook smoke | `go test ./internal/install/... -run 'CursorLoopGateFailClosed\|HookDrift\|…'` |
| Matrix | `evals/p28-regression/score_arm_labels_test.sh` + TEST-MATRIX M-01..M-16 |
| Default successor | **`no successor`** when R1–R5 closed + green |
| Phase 29 | Human promote only |
| Never | Successor `TBD` |
