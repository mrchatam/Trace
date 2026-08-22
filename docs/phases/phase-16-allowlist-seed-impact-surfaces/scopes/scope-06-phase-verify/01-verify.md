# P16 / S06 / 01 — Phase 16 VERIFY

## Metadata
- id: P16-S06-01
- todo_ids: [P16-S06-01]
- role: verify
- skills: [systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 16 — S01–S05 named DF tests + carry-forward bars. Write `VERIFY-NOTES.md`. On pass, **start DR-HANDOFF = `no successor`**. On fail, spawn `P16-S06-01a`/`01b` (+`01c`).

**Stop if sibling `00-PLANNER.md` is not FINAL.**

## References
- [00-PLANNER.md](00-PLANNER.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- S01–S05 REVIEW-NOTES
- Pattern: [../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/01-verify.md](../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/01-verify.md)

## Evidence table (fill in VERIFY-NOTES.md)
| Bucket | Must prove |
|--------|------------|
| S01 DF-76 | Virgin `project=` does not auto-init / escape DENIED |
| S02 DF-75/77/78 | CHECK + no YOLO fail-open; CLI DENIED; unprefixed slug |
| S03 DF-68 | install `-C` honors project root; DF-22 tip keepers |
| S04 DF-70/73 | seed mentions-task + impact findings |
| S05 DF-71/72/74 | packet impact + snake_case + `trace_impact` MCP |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false`; product pkgs |
| Residuals OK | DF-67 defer; R2 defer; R3/R4 wontfix — **not** fail criteria |
| Laws | No daemon/HTTP/embeddings; G19; no install/decide MCP; no full-rebuild indexer |

## Locked verify commands (00-PLANNER FINAL may replace)
```bash
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... ./internal/install/... ./internal/compiler/... ./cmd/trace/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## DR-HANDOFF
On pass: VERIFY-NOTES records **`no successor`** (start). Do **not** invent Phase 17 / S05 supersession / plan simulate / D21+. P15 historical `no successor` left intact.

## Exit criteria
- [ ] Named S01–S05 tests + carry-forward PASS (or honest fail+spawn)
- [ ] VERIFY-NOTES written; DR-HANDOFF started
- [ ] Board Notes → **P16-S06-02**
