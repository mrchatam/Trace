# P06 / S03 / 02 — Scope review (Phase 06 VERIFY)

## Metadata
- id: P06-S03-02
- todo_ids: [P06-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 06 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 07 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that capability-selection ablation passed without `evals/capability` `TestPlantedCapabilitySelectionAblation` (+ schema/metrics evidence). Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, or ablation pass. Spot-check planted tallies (TP=3/FN=0/FP=0/TN=1; P/R=1.0), S01 hooks (`UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet required/missing), Paths A/B/C intact, Gate F + Gate G + Gate E green, and Gate C `dry_run:false` artifacts intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-capability-surface/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-capability-selection/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 7 Performance ladder & language plugins
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 05 VERIFY review [`../../../phase-05-decision-impact/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-05-decision-impact/scopes/scope-03-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **independently** re-run ablation harness (`evals/capability` `TestPlantedCapabilitySelectionAblation`) — not copy S02 Notes?
- Did VERIFY prove schema `schema-capability.json` v1 + temp `metrics-capability.json` + tallies TP=3/FN=0/FP=0/TN=1 + P/R=1.0?
- Did VERIFY prove S01 hooks still green (`UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet required/missing)?
- Did VERIFY re-run honesty A/B/C + Gate G + Gate E + Gate F + p0x + x0 + domain/store/planner/compiler + `./...`?
- Evidence table covers ablation path + S01 mig 010 + Gate F + Gate G + Gate E + Gate C `dry_run:false` intact + dry-run≠Gate C/≠Gate F/≠Gate G/≠ablation + law checks + handoff?
- Residuals (DPC-global, non-tx Apply, UNIQUE re-link, MCP no severity, GC-03/04 deferred, `plan simulate` out, S02 lows) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-07-performance-ladder/README.md` exists (goal = Performance ladder & language plugins / A_PROJECT_PLAN Phase 7)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 07 section; first pending row is **`P07-00`** after Phase 06’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 06 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Ablation path | `evals/capability` **`TestPlantedCapabilitySelectionAblation`** + `schema-capability.json` / temp `metrics-capability.json` |
| Tallies | TP=3 / FN=0 / FP=0 / TN=1; precision=1.0; recall=1.0 |
| S01 / S02 | mig 010 + Upsert/Require/Missing + packet; `evals/capability` planted harness |
| Gate F | `TestPlantedImpactConflictsGateFPrelim` still green |
| Gate G | `TestHonestyEscapeRateGateGPrelim` still green |
| Gate E | `TestPlantedDiscoveryReplan` still green |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 07 path | `docs/phases/phase-07-performance-ladder/` |
| First Phase 07 board row | `P07-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no performance-ladder / language-plugin product features on this row |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 07** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 06 `done` prompts. Do not execute Phase 07 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 07 runnable **or** explicit stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 06 complete; next runnable is **`P07-00`** (or stop)
- [ ] Explicit: ablation = planted `evals/capability` P/R harness; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation; Gate C artifacts intact; no commercial capability theater

## Minimal todos
- [ ] Compare VERIFY claims vs ablation + S01/S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 07 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 06 complete; next `P07-00`)
