# P08 / S04 / 02 — Scope review (Phase 08 VERIFY / phase close)

## Metadata
- id: P08-S04-02
- todo_ids: [P08-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (Phase 08 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that the checklist passed without `evals/compat` `TestCompatibilitySecurityChecklist` (+ `schema-compat.json` v1 + temp `metrics-compat.json`). Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or this checklist. Spot-check: harness was **created** in VERIFY (S01–S03 left none); S01 `LanguageAdapterAPIVersion=1`; S02 path-local + `trace.lock`; S03 migrate/backup/auth fail-closed; no BLOBs; no `011_*`; G19; no daemon/HTTP primary; carry-forward bars; Gate C `dry_run:false` intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-plugin-apis/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-worktrees/REVIEW-NOTES.md)
- [S03 REVIEW-NOTES.md](../scope-03-production-hardening/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 8 is last planned phase
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 07 VERIFY review [`../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S04-01).

## Review focus
- Did VERIFY **create** `evals/compat` (S01–S03 left none) and independently run `TestCompatibilitySecurityChecklist` — not copy Notes?
- Did VERIFY prove `schema-compat.json` v1 + temp `metrics-compat.json` (`dry_run:false`) + all must-pass `*_ok` + `language_adapter_api_version == 1`?
- Did VERIFY re-run S01 contribution tests + S02 path-local/lock + S03 migrate/backup/auth + honesty A/B/C + Gate G/E/F + ablation + Gate H + p0x + x0 + `./...`?
- Evidence table covers checklist path + S01–S03 + Gate C `dry_run:false` + dry-run≠Gate C/≠Gate F/≠Gate G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (S03 argv token, restore window, S02 exit 2, deferred 100k/1M, GC-03/04, A5) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened; checklist not invented as theater?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] Board / phase README do **not** claim a Phase 09 scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 09 artifacts required — confirm absence is intentional, not forgotten
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 08 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Checklist path | `evals/compat` **`TestCompatibilitySecurityChecklist`** + `schema-compat.json` / temp `metrics-compat.json` |
| Harness ownership | Created in S04-01 VERIFY (not pre-seeded by S01–S03) |
| Must-pass | API version=1; path-local + lock; migrate status; backup↔restore + no BLOBs; local-auth fail-closed; G19; no daemon/HTTP primary; no `011_*` |
| S01 / S02 / S03 | Contribution-path + isolation/lock + migrate/backup/auth tests stay green |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold fixes; no Phase 09 product features on this row |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF). Do not rewrite Phase 08 `done` prompts. Do not invent Phase 09 without explicit promotion. On APPROVE: mark Phase 08 complete; next runnable = **none** (roadmap closed) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 08 complete; no accidental Phase 09 without promotion
- [ ] Explicit: checklist = planted `evals/compat`; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; no commercial security theater

## Minimal todos
- [ ] Compare VERIFY claims vs checklist harness + S01–S03 + suite evidence (+ optional fresh harness runs)
- [ ] Confirm DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 08 complete)

## Out of scope
- Inventing Phase 09 without promotion
- Inventing checklist pass without harness evidence
- Weakening prior gates
- Executing product feature work outside review/spawn rights
