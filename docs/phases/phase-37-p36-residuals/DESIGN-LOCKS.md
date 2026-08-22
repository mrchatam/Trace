# DESIGN-LOCKS — Phase 37

**Human-promoted 2026-08-22.** Phase 36 residuals closure.

| Lock | Value |
|------|-------|
| Predecessor | Phase 36 CLOSED (`P36-S03-02`) |
| Theme | **Close Phase 36 residuals** — triage then implement or re-defer |
| Must do | S00 triage every R1–R11 with accept/defer/reject + evidence |
| Must preserve | Phase 36 guarantees: MCP `trace_plan`, bootstrap, terminal advisory, active `plan_missing` block |
| Must not | Silent PlanExists bridge; weaken active PLAN enforcement; fork business logic in `web/` (Law 19) |
| Dogfood | `/home/ali/Desktop/feet seller telegram app` — read-only unless R9 recovery scoped |
| Out of scope | Hosted SaaS; full plan reconstruction from 123 tasks |

## Accept criteria sketch (S01 locks per residual)

| Residual | Accept if… |
|----------|------------|
| R1 advisory bridge | Recommends bootstrap; **does not** set `PlanExists=true` |
| R2 HTTP plan routes | Thin handlers over `internal/planner`; OpenAPI updated |
| R3 MCP loop gate | `trace_loop action=gate` mirrors CLI JSON + exit semantics |
| R4 bootstrap help | Help text mentions human refinement via create-coarse/deep |
| R5 advisories[] | `trace.loop.status.v1` includes goal-structure warning when triggered |
| R6 config warn test | Unit test for `WarnIfTraceDirWithoutConfig` |
| R7 enforce warn default | Explicit product decision in S01 — default stays `off` unless locked |
| R8 plan UX | At least one surface beyond TaskDetail OR defer with owner |
| R9 feet refinement | Documented path or tooling; no mandatory history rewrite |
| R10 live GUI | Browser spot-check terminal + plan surfaces OR evidence pin |
| R11 critique path | MCP/doc path for plan_critiqued after bootstrap |

## Success sketch

Post-P37: no orphaned P36 deferrals without owner; loop status carries advisories; bootstrap help honest; optional HTTP/MCP gate parity; GUI plan gap visible beyond TaskDetail (if accepted).
