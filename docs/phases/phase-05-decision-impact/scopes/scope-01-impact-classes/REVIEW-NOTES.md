# P05-S01-02 — Scope review notes (2026-08-16)

Independent review of S01 against `01-impact-classes.md` + TODO Notes for `P05-S01-01`. Fresh session; claims verified in-repo.

## Plan (executed)

1. Diff claims vs mig `009_decision_impact.sql`, `internal/store/impact.go`, `internal/domain/impact.go`, thin `cmd/trace/impact.go`
2. Check fail-closed `ImpactReport`, enums, no new entity_links rels, no `internal/impact` / planner fork / MCP impact / `plan simulate`
3. Confirm S02 Gate F hooks usable; Gate C `dry_run:false` intact
4. Re-run honesty / Gate G / Gate E / p0x / x0 / `./...`
5. Severity-tag findings; spawn only for blocker/high (none)
6. Write these notes; mark board + SCOPE-TODOS; light S02 Depends confirm

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**. Next board row: **P05-S02-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Mig **`009_decision_impact.sql`** additive; tables `decision_impact_findings` + `decision_alternatives`; no ALTER on `decisions` | Pass (schema file + Open migrates v9; store test lists both tables; `001`–`008` unchanged) |
| Package **`internal/domain`** + store helpers; **no** `internal/impact`; **no** planner fork | Pass (`internal/impact` absent; no Impact* under `internal/planner`) |
| Enums: bands SAFE\|CAUTION\|HIGH\|DESTRUCTIVE\|REVERSAL; uncertainty KNOWN\|LIKELY\|POSSIBLE\|UNKNOWN (empty→UNKNOWN); kinds AFFECTED_WORK…UNRESOLVED | Pass (`Normalize*` + `TestImpactFindingAddListAndReject`) |
| APIs Add/List findings; Add/List/SetRecommended alternatives (single recommended); `ImpactReport` | Pass (`impact.go` + exclusivity + report tests) |
| Fail-closed: `HasUnknown` / `Incomplete` / OverallClass order SAFE<…<REVERSAL; UNKNOWN findings not omitted | Pass (`TestImpactReportFailClosedAndRollup`) |
| **No new** entity_links rels; keep `decision_affects_task` only | Pass (`TestImpactDoesNotAddEntityLinkRels`; Rel consts unchanged aside from prior phases) |
| Thin CLI `trace impact` (finding/alternative/report) G19; no business logic beyond argv→domain | Pass (`cmd/trace/impact.go` + help + root dispatch) |
| MCP impact tools absent; `plan simulate` absent | Pass (mcp grep empty; no simulate cmd) |
| S02 hooks: `AddImpactFinding`, `LinkDecisionTask`, `ImpactReport` fields match stubs | Pass (S02 Depends already accurate; confirmed live) |
| Honesty A/B/C + Gate G `TestHonestyEscapeRateGateGPrelim` | Pass |
| Gate E `TestPlantedDiscoveryReplan` | Pass |
| p0x 7/7; x0; `CGO_ENABLED=0` domain+store; `CGO_ENABLED=1` `./...` | Pass |
| Gate C `docs/verification/gate-c-x0/` `dry_run:false`; B0 acc 0 / G1 0.8 pattern intact | Pass |

## Re-verification commands (2026-08-16)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... -count=1
# ok

CGO_ENABLED=0 go test ./evals/honesty/ -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim' -count=1
# PASS both

CGO_ENABLED=0 go test ./evals/replan/ -run TestPlantedDiscoveryReplan -count=1
# PASS

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
# ok (p0x + x0 + full tree)
```

Domain impact tests: `TestImpactFindingAddListAndReject`, `TestDecisionAlternativeRecommendExclusivity`, `TestImpactReportFailClosedAndRollup`, `TestImpactDoesNotAddEntityLinkRels`.

## Findings

### Blocker / high

None.

### Medium (residual — no spawn)

- **Non-atomic recommend exclusivity:** `AddDecisionAlternative` (Recommended) and `SetRecommendedAlternative` clear siblings then write without a single DB transaction. A crash between clear and set can leave zero recommended (report surfaces `Incomplete`). Same class as prior non-tx planner/domain residuals; harden later if concurrent writers appear.

### Low / nit

- **OverallClass vs UNKNOWN:** A sole finding with class `SAFE` and uncertainty `UNKNOWN` yields `OverallClass=SAFE` with `HasUnknown`/`Incomplete` true (locked rollup is max severity; fail-closed is via flags). Gate F consumers must score `HasUnknown`/`Incomplete`, not trust `OverallClass` alone — S02 Depends already lists those fields.
- **CLI coverage:** No dedicated `cmd/trace` impact integration test; thin adapter + domain tests cover policy. Optional later.
- **ctx ignored** on impact APIs (`_ = ctx`) — consistent with other domain methods until cancellation is wired.

## Spawns

None.

## S02 note

Upcoming `P05-S02-00` Depends hooks match shipped surface (`AddImpactFinding` / `LinkDecisionTask` / `ImpactReport` + `HasUnknown`/`Incomplete`/`OverallClass`; mig 009). No prompt rewrite required beyond a live-confirm stamp on S02 planner Depends.
