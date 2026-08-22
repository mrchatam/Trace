# P04 / S03 / 02 — Scope review (Phase 04 VERIFY)

## Metadata
- id: P04-S03-02
- todo_ids: [P04-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 04 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 05 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Gate G passed without `evals/honesty` `TestHonestyEscapeRateGateGPrelim` (+ schema/metrics evidence). Reject treating Phase 01 dry-run as Gate C or as Gate G. Spot-check planted tallies (escapes=1/caught=2/attempts=3), S01 hooks (`LinkReviewScope` / OPEN `POLICY_EXCEPTION` / `CountOpenResidualsByScope`), Paths A/B/C intact, Gate E green, and Gate C `dry_run:false` artifacts intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-scope-review-layer/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-honesty-escape-rate/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 5 Decision impact & simulation
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate G
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 03 VERIFY review [`../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **independently** re-run Gate G harness (`evals/honesty` `TestHonestyEscapeRateGateGPrelim`) — not copy S02 Notes?
- Did VERIFY prove schema `schema-gate-g.json` v1 + temp `metrics-gate-g.json` + tallies escapes=1/caught=2/attempts=3 + hatch=escape only?
- Did VERIFY prove S01 hooks still green (`LinkReviewScope` / `CountOpenResidualsByScope` / OPEN `POLICY_EXCEPTION`)?
- Did VERIFY re-run honesty A/B/C + Gate E + p0x + x0 + domain/store/planner + `./...`?
- Evidence table covers Gate G path + S01 mig 008 + Gate E + Gate C `dry_run:false` intact + dry-run≠Gate C + law checks + handoff?
- Residuals (DPC-global, non-tx Apply, UNIQUE re-link, MCP no severity, s01_hooks schema looseness) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-05-decision-impact/README.md` exists (goal = Decision impact & simulation / A_PROJECT_PLAN Phase 5)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 05 section; first pending row is **`P05-00`** after Phase 04’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 04 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Gate G path | `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** + `schema-gate-g.json` / temp `metrics-gate-g.json` |
| Tallies | escapes=1 / caught=2 / attempts=3; hatch=escape only in Gate G report |
| S01 / S02 | mig 008 + `LinkReviewScope` / residuals; honesty extend (keep A/B/C) |
| Gate E | `TestPlantedDiscoveryReplan` still green |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 05 path | `docs/phases/phase-05-decision-impact/` |
| First Phase 05 board row | `P05-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no impact-engine / simulate features |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 05** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 04 `done` prompts. Do not execute Phase 05 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 05 runnable **or** explicit stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 04 complete; next runnable is **`P05-00`** (or stop)
- [ ] Explicit: Gate G = planted `evals/honesty` escape-rate report; Phase 01 dry-run ≠ Gate C; Gate C artifacts intact

## Minimal todos
- [ ] Compare VERIFY claims vs Gate G + S01/S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 05 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 04 complete; next `P05-00`)
