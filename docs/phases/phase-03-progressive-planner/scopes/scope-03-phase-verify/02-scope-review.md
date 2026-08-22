# P03 / S03 / 02 — Scope review (Phase 03 VERIFY)

## Metadata
- id: P03-S03-02
- todo_ids: [P03-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 03 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 04 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Gate E passed without `evals/replan` `TestPlantedDiscoveryReplan`. Reject treating Phase 01 dry-run as Gate C. Spot-check severity (`PLAN_AFFECTING`+ only) + churn N=5 fail-closed/ack and that Gate C `dry_run:false` artifacts remain intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-coarse-planner/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-discovery-replan/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 4 Review depth & evidence policies
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate E
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 02 VERIFY review [`../../../phase-02-gate-c/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-02-gate-c/scopes/scope-03-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **independently** re-run Gate E harness (`evals/replan` `TestPlantedDiscoveryReplan`) — not copy S02 Notes?
- Did VERIFY prove severity (`PLAN_AFFECTING`+ only; INFO no auto-replan) + churn N=5 fail-closed/ack?
- Did VERIFY re-run honesty + p0x + x0 + planner/store/domain + `./...`?
- Evidence table covers Gate E path + S01 mig 006 + S02 mig 007 + Gate C `dry_run:false` intact + dry-run≠Gate C + law checks + handoff?
- Residuals (DPC-global, non-tx Apply, UNIQUE re-link, MCP no severity) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-04-review-depth/README.md` exists (goal = Review depth & evidence policies / A_PROJECT_PLAN Phase 4)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 04 section; first pending row is **`P04-00`** after Phase 03’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 03 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Gate E path | `evals/replan` **`TestPlantedDiscoveryReplan`** |
| Severity / churn | PLAN_AFFECTING+ only; N=5 fail-closed + ack |
| S01 / S02 | `internal/planner` + mig 006; mig 007 + `ApplyDiscoveryReplan` |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 04 path | `docs/phases/phase-04-review-depth/` |
| First Phase 04 board row | `P04-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no review-depth features |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 04** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 03 `done` prompts. Do not execute Phase 04 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 04 runnable **or** explicit stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 03 complete; next runnable is **`P04-00`** (or stop)
- [ ] Explicit: Gate E = planted `evals/replan` demo; Phase 01 dry-run ≠ Gate C; Gate C artifacts intact

## Minimal todos
- [ ] Compare VERIFY claims vs Gate E + S01/S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 04 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 03 complete; next `P04-00`)
