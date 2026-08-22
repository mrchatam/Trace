# P07 / S03 / 00-PLANNER — Phase 07 VERIFY / Gate H

## Metadata
- id: P07-S03-00
- todo_ids: [P07-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 07 VERIFY commands, evidence table, Gate H measurement bar (after first ladders exist), spawn 01a/b/c on fail, and DR-HANDOFF Phase 08 = `phase-08-ecosystem-hardening`. No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 7 / 8
- [docs/STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §10 metrics
- [docs/init/E_ASSUMPTION_REGISTER.md](../../../../init/E_ASSUMPTION_REGISTER.md) A5
- Pattern: Phase 06 VERIFY [`../../../phase-06-environment-capability/scopes/scope-03-phase-verify/`](../../../phase-06-environment-capability/scopes/scope-03-phase-verify/)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Phase defaults already locked (respect — P07-00)
| Item | Value |
|------|-------|
| Gate H harness preference | **`evals/perf`** planted/synthetic ladders |
| Preferred names | `TestPlantedPerfLadderGateH` + `schema-gate-h.json` v1 + temp `metrics-gate-h.json` (**finalize here**) |
| Thresholds | **After first measurements** — do not invent pass numbers in S01/S02; S01 may only seed optional `evals/perf` fixtures/timing logs (no pass claim) |
| S01 Depends | Re-prove T0 skip + sibling isolation; consume any optional `evals/perf` seeds from S01 — **S01-02 note:** no `evals/perf` fixtures seeded (only `t.Logf` in `TestWalkIndexableT0AlwaysSkip`); do not assume planted ladders exist yet |
| S02 Depends | Language adapter **Go** via `tree-sitter-go` v0.25.0 — **S02-02 APPROVE (live):** `LangGo`+`extract_go.go`+golden; **no** `evals/perf` / thresholds from S02 — Gate H may optionally include tiny `.go` fixtures; **thresholds still after measurements** |
| Phase 08 folder | **`phase-08-ecosystem-hardening`** |
| Carry-forward | Honesty A/B/C; Gate G/E/F; capability ablation; p0x; x0; Gate C `dry_run:false` |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H |

## Planner work
1. Lock Gate H harness path + thresholds (or explicit “measure then set” protocol if ladders first appear in S01/S02).
2. Thickened `01-verify.md` evidence table + spawn convention + Phase 08 checklist.
3. Thickened `02-scope-review.md` owns DR-HANDOFF completion.
4. SCOPE-TODOS sync.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Gate H path | **`evals/perf`** / **`TestPlantedPerfLadderGateH`** |
| Schema / metrics | **`schema-gate-h.json` v1** + temp **`metrics-gate-h.json`** (`dry_run:false`) |
| Harness ownership | **S03-01 creates** planted harness as VERIFY work (S01/S02 left none) — do **not** block Gate H awaiting a seed |
| Size ladder | Synthetic **smoke** (~50–200 LOC) + **~1k** + **~10k**; toward 10k–1M tables; **100k/1M CI deferred** |
| Thresholds | **Measure-then-threshold** — derive ceilings after first measurements; no commercial theater invent |
| Optional Go fixtures | Tiny `.go` OK on smoke (optional higher); thresholds still after measure |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` if needed) immediately below |
| DR-HANDOFF | S03-01 starts Phase 08 = **`phase-08-ecosystem-hardening`**; S03-02 owns completion |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; p0x; x0; Gate C `dry_run:false`; S01 T0+isolation; S02 Go golden |

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] Phase 08 folder name locked = `phase-08-ecosystem-hardening`
- [x] Gate H names + measure-then-threshold + S03-01 harness ownership locked
- [x] Board Notes; next `P07-S03-01`

## Out of scope
- Running VERIFY (S03-01)
- Implementing Phase 08
- Inventing Gate H pass numbers without measurements
