# P07 / S03 / 02 — Scope review (Phase 07 VERIFY)

## Metadata
- id: P07-S03-02
- todo_ids: [P07-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 07 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 08 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Gate H passed without `evals/perf` `TestPlantedPerfLadderGateH` (+ schema/metrics + measure-then-threshold derivation). Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, or Gate H. Spot-check: harness was **created** in VERIFY (S01/S02 left none); rungs smoke/~1k/~10k; thresholds derived not invented; S01 T0+isolation; S02 Go golden; carry-forward bars; Gate C `dry_run:false` intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-incremental-indexing/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-language-plugins/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 8 Ecosystem & hardening
- [docs/STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §10
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 06 VERIFY review [`../../../phase-06-environment-capability/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-06-environment-capability/scopes/scope-03-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **create** `evals/perf` (S01/S02 left none) and independently run `TestPlantedPerfLadderGateH` — not copy Notes?
- Did VERIFY prove `schema-gate-h.json` v1 + temp `metrics-gate-h.json` (`dry_run:false`) + measure-then-threshold derivation in VERIFY-NOTES?
- Rungs smoke / ~1k / ~10k present; 100k/1M deferred (not claimed as CI pass)?
- Structural: `t0_skip_ok` + `incremental_isolation_ok` (+ optional `go_adapter_exercised`)?
- Did VERIFY re-run S01 T0/isolation + S02 Go golden + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + `./...`?
- Evidence table covers Gate H path + S01/S02 + Gate C `dry_run:false` + dry-run≠Gate C/≠Gate F/≠Gate G/≠ablation/≠Gate H + law checks + handoff?
- Residuals (DPC-global, GC-03/04, A5, deferred 100k/1M, S01/S02 lows) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened; thresholds not invented as theater?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-08-ecosystem-hardening/README.md` exists (goal = Ecosystem & hardening / A_PROJECT_PLAN Phase 8)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 08 section; first pending row is **`P08-00`** after Phase 07’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 07 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Gate H path | `evals/perf` **`TestPlantedPerfLadderGateH`** + `schema-gate-h.json` / temp `metrics-gate-h.json` |
| Harness ownership | Created in S03-01 VERIFY (not pre-seeded by S01/S02) |
| Thresholds | Measure-then-threshold — derivation must appear in VERIFY-NOTES |
| Rungs | smoke + ~1k + ~10k; 100k/1M deferred |
| S01 / S02 | T0+isolation tests; Go golden + tree-sitter-go v0.25.0 |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 08 path | `docs/phases/phase-08-ecosystem-hardening/` |
| First Phase 08 board row | `P08-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no ecosystem product features on this row |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 08** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 07 `done` prompts. Do not execute Phase 08 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 08 runnable **or** explicit stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 07 complete; next runnable is **`P08-00`** (or stop)
- [ ] Explicit: Gate H = planted `evals/perf` ladder; measure-then-threshold; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H; Gate C artifacts intact; no commercial perf theater

## Minimal todos
- [ ] Compare VERIFY claims vs Gate H harness + S01/S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 08 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 07 complete; next `P08-00`)

## Out of scope
- Executing Phase 08 implement wave
- Inventing Gate H pass without harness evidence
- Weakening prior gates
