# P03-S02-02 — Scope review notes (discovery replan)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P03-S02-01. Claims match live repo: mig `007_discovery_severity.sql`, store `IncrementAutoReplanCount`/`AckAutoReplan`, `planner.ApplyDiscoveryReplan` + `AckReplan` with INFO no-supersede / PLAN_AFFECTING+ via `SupersedeDeepPlan` + `LinkDiscoveryPlanChange`, churn N=5 fail-closed + ack reset, `evals/replan` `TestPlantedDiscoveryReplan`, thin CLI G19. Honesty / p0x / x0 / `./...` green. Mode-B Gate C packs untouched (scores still G1 0.800 > B0 ~0; `dry_run:false` N=3). DPC attach still global. No second planner package. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Severity: only PLAN_AFFECTING/BLOCKING auto-replan; INFO no supersede/increment | **Pass** — `replan.go` switch; unit + demo assert count/revision unchanged on INFO |
| Churn N=5 fail-closed; Increment + Ack (reset 0); no unbounded loops | **Pass** — `DefaultMaxAutoReplans=5`; `ErrReplanBudgetExceeded` when `count >= N`; `AckAutoReplan` → 0 |
| S01 consume: SupersedeDeepPlan / GetCurrent\|ListScopes\|GetPlan; LinkDiscoveryPlanChange; no second planner | **Pass** — Apply calls `SupersedeDeepPlan` + domain `LinkDiscoveryPlanChange`; S01 APIs retained in `service.go`; only `internal/planner` |
| Demo `evals/replan` `TestPlantedDiscoveryReplan` | **Pass** — planted Goal→coarse→deep→PA/INFO/budget/ack path |
| CLI G19 thin: `--severity`, `plan apply-discovery`, `plan ack-replan` | **Pass** — `cmd/trace/add.go` + `plan.go` adapters only |
| Out of scope held: DPC global; honesty/p0x/Gate C; Mode-B packs | **Pass** — `discoveryPlanChangeHits` still global; Mode-B `testdata/gate-c/*-run*.json` unchanged; metrics still Go |
| S03 VERIFY stubs list Gate E inputs | **Pass** — `scope-03` `00-PLANNER` + `SCOPE-TODOS` list replan demo + severity + churn |

## Claims → evidence

| Claim (P03-S02-01 Notes) | Evidence |
|--------------------------|----------|
| Mig 007 severity | `internal/store/schema/007_discovery_severity.sql`; embed via `schema/*.sql` |
| Store Increment/Ack | `plan_hierarchy.go` `IncrementAutoReplanCount` / `AckAutoReplan` |
| ApplyDiscoveryReplan + AckReplan | `internal/planner/replan.go` |
| INFO no auto-replan | `applyInfoDiscovery`; tests `TestApplyDiscoveryReplanINFONoSupersede`, demo INFO block |
| PLAN_AFFECTING+/BLOCKING supersede + link + count | `applyAutoReplanDiscovery`; unit + demo |
| N=5 + ack | `TestApplyDiscoveryReplanBudgetAndAck`, demo exhaust→fail→ack→succeed |
| CLI | `add discovery --severity`; `plan apply-discovery` / `ack-replan` |
| No MCP replan tools | MCP still six write/read tools; discovery add defaults severity INFO (no severity field) |

### Implementer residuals (verified)

| Note | Verified | Disposition |
|------|----------|-------------|
| `ApplyDiscoveryReplan` not single transaction — partial failure risk | Yes — order link → supersede → increment; mid-fail can leave link without supersede, or superseded plan without count bump | **Residual medium** — same class as S01 non-tx CreateCoarsePlan; no spawn |
| Re-apply same discovery→plan_change may hit UNIQUE | Yes — `InsertLink` UNIQUE; re-link identical endpoints errors | **Residual low** — happy path creates new PlanChange when ID omitted |
| MCP `trace_add` discovery has no severity | Yes — `tools_write.go` CreateDiscovery without Severity → INFO | **OK** — MCP not in S02 scope; CLI-only severity |
| Gate C packs untouched | Yes — Mode-B packs mtime 01:08; metrics `understanding_accuracy` 0.8 / B0 miss pattern; `dry_run:false` | **Pass** |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./internal/planner/... ./internal/store/... ./internal/domain/... ./evals/replan/...  PASS
CGO_ENABLED=0 go test ./evals/honesty/...                                                               PASS
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/...                                                    PASS
CGO_ENABLED=1 go test ./...                                                                             PASS
```

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._ Residual: non-transactional Apply (see above) — acceptable for CLI/library single-writer; transactional wrap can wait for a measured need.

### low

1. **Duplicate discovery→plan_change link** — re-`Apply` with same `PlanChangeID` fails UNIQUE before supersede; callers should omit ID (auto-create) or use a fresh PlanChange.
2. **Budget TOCTOU** — check-then-increment not serialized across concurrent callers; fine for local-first single-writer.
3. **MCP severity absent** — MCP-created discoveries stay INFO until a later surface adds severity (intentional for this scope).

### nit

1. `AckReplan` ignores AppendEvent errors (`_, _ =`).
2. Severity constants duplicated in `store` / `domain` / `planner` (re-export OK per lock).

## Spawns
_None._

## Next board row
**P03-S03-00** (VERIFY planner — Gate E path + Phase 04 handoff). Do not start from this review session beyond board update.
