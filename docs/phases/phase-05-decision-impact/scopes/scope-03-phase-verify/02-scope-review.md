# P05 / S03 / 02 — Scope review (Phase 05 VERIFY)

## Metadata
- id: P05-S03-02
- todo_ids: [P05-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 05 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 06 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Gate F passed without `evals/impact` `TestPlantedImpactConflictsGateFPrelim` (+ schema/metrics evidence). Reject treating Phase 01 dry-run as Gate C or as Gate F. Spot-check planted tallies (TP=3/FN=0/FP=0/TN=1; P/R=1.0), S01 hooks (`AddImpactFinding` / `LinkDecisionTask` / `ImpactReport`), Paths A/B/C intact, Gate G + Gate E green, and Gate C `dry_run:false` artifacts intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-impact-classes/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-gate-f-prelim/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-NOIMP
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 6 Environment/capability graph
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate F
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 04 VERIFY review [`../../../phase-04-review-depth/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-04-review-depth/scopes/scope-03-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **independently** re-run Gate F harness (`evals/impact` `TestPlantedImpactConflictsGateFPrelim`) — not copy S02 Notes?
- Did VERIFY prove schema `schema-gate-f.json` v1 + temp `metrics-gate-f.json` + tallies TP=3/FN=0/FP=0/TN=1 + P/R=1.0?
- Did VERIFY prove S01 hooks still green (`AddImpactFinding` / `LinkDecisionTask` / `ImpactReport` / `decision_affects_task`)?
- Did VERIFY re-run honesty A/B/C + Gate G + Gate E + p0x + x0 + domain/store/planner + `./...`?
- Evidence table covers Gate F path + S01 mig 009 + Gate G + Gate E + Gate C `dry_run:false` intact + dry-run≠Gate C/≠Gate F + law checks + handoff?
- Residuals (DPC-global, non-tx Apply, UNIQUE re-link, MCP no severity, s01_hooks, GC-03/04 deferred, `plan simulate` out) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-06-environment-capability/README.md` exists (goal = Environment/capability graph / A_PROJECT_PLAN Phase 6)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 06 section; first pending row is **`P06-00`** after Phase 05’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 05 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Gate F path | `evals/impact` **`TestPlantedImpactConflictsGateFPrelim`** + `schema-gate-f.json` / temp `metrics-gate-f.json` |
| Tallies | TP=3 / FN=0 / FP=0 / TN=1; precision=1.0; recall=1.0 |
| S01 / S02 | mig 009 + findings/alts/`ImpactReport`; `evals/impact` planted harness |
| Gate G | `TestHonestyEscapeRateGateGPrelim` still green |
| Gate E | `TestPlantedDiscoveryReplan` still green |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 06 path | `docs/phases/phase-06-environment-capability/` |
| First Phase 06 board row | `P06-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no capability-graph / MCP-ontology product features on this row |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 06** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 05 `done` prompts. Do not execute Phase 06 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 06 runnable **or** explicit stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 05 complete; next runnable is **`P06-00`** (or stop)
- [ ] Explicit: Gate F = planted `evals/impact` P/R harness; Phase 01 dry-run ≠ Gate C / ≠ Gate F; Gate C artifacts intact; DR-NOIMP respected

## Minimal todos
- [ ] Compare VERIFY claims vs Gate F + S01/S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 06 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 05 complete; next `P06-00`)
