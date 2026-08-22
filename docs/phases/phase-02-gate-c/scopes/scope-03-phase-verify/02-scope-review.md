# P02 / S03 / 02 — Scope review (Phase 02 VERIFY)

## Metadata
- id: P02-S03-02
- todo_ids: [P02-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 02 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 03 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Phase 01 dry-run alone closed Gate C. Spot-check S02 harden (GC-01 tests + GC-02 README/hash) and that Mode-B packs were not falsified for q3.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [GATE-C-NOTES.md](../scope-01-x0-gate-c/GATE-C-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-slice-hardening/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 3 progressive planner
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 01 VERIFY review [`../../../phase-01-x0-readiness/scopes/scope-05-phase-verify/02-scope-review.md`](../../../phase-01-x0-readiness/scopes/scope-05-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY **independently** re-check Gate C artifacts (`dry_run:false`, N=3, kill, Go) — not copy S01 Notes?
- Did VERIFY re-run honesty + p0x + x0 + S02 GC-01/02 + `./...`?
- Evidence table covers Gate C + dry-run≠Gate C + GC-01/02 + deferrals + law checks + handoff?
- Mode-B packs left historical (no q3 rewrite required / no falsification)?
- Residual global DPC-on-every-task noted as non-blocking (S02 REVIEW-NOTES)?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-03-progressive-planner/README.md` exists (goal = progressive planner / A_PROJECT_PLAN Phase 3)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 03 section; first pending row is **`P03-00`** after Phase 02’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete
- If Gate C were No-Go: explicit `no successor` / stop in VERIFY-NOTES instead (not expected)

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Gate C re-check | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 |
| S02 GC-01 | `TestWhyTaskIncludesDiscoveryPlanChange`, `TestTaskContextIncludesDiscoveryPlanChange` |
| S02 GC-02 | `TestFixtureReadmeHasNoGTUUIDOracle`; hash `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Phase 03 path | `docs/phases/phase-03-progressive-planner/` |
| First Phase 03 board row | `P03-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no progressive-planner features |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 03** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 02 `done` prompts. Do not execute Phase 03 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (Phase 03 runnable **or** explicit No-Go stop)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 02 complete; next runnable is **`P03-00`** (or stop)
- [ ] Explicit: Phase 01 dry-run ≠ Gate C pass; Gate C Go re-confirmed

## Minimal todos
- [ ] Compare VERIFY claims vs Gate C + S02 + suite evidence (+ optional fresh harness runs)
- [ ] Finish Phase 03 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 02 complete; next `P03-00`)
