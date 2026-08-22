# P09 / S04 / 02 — Scope review (Phase 09 VERIFY / phase close)

## Metadata
- id: P09-S04-02
- todo_ids: [P09-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (Phase 09 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Phase 09 passed without **`TestWhyAndContextWithLinkedReview`** (DF-01) plus S02/S03 spot-checks and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing MCP list-tasks as a VERIFY requirement. Spot-check: DF-01 + `trace tasks` / seed-`-C` + `trace install cursor` + DF-05 docs; honesty/Gates/ablation/compat/p0x/x0/`./...`; Gate C `dry_run:false` intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-retrieval-review/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-discoverability/REVIEW-NOTES.md)
- [S03 REVIEW-NOTES.md](../scope-03-install-wire/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md), [experiments/LADDER.md](../../../../../experiments/LADDER.md)
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 08 VERIFY review [`../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/02-scope-review.md`](../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S04-01).

## Review focus
- Did VERIFY independently run **`TestWhyAndContextWithLinkedReview`** — not copy Notes?
- Did VERIFY re-prove S02 (`TestTasksListAfterSeed` / `TestSeedImportRelativePathAgainstC`) and S03 (`TestInstallCursor*` + DF-05 docs)?
- MCP still six tools — **no** list-tasks invent as a FAIL reason?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat + p0x + x0 + `./...`?
- Evidence table covers DF-01 + S02/S03 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (`plan_scope` lookup out; scope-only review untested; S03 degenerate mcpServers) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] Board / phase README do **not** claim a Phase 10 scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 10 artifacts required — confirm absence is intentional; remaining ladder gaps (D08/D09/combos/multi-agent; DF-11/12) stay parallel `experiments/`
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 09 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| DF-01 path | `internal/retrieval` **`TestWhyAndContextWithLinkedReview`** |
| S02 / S03 | Named CLI/store/install tests + DF-05 docs |
| MCP | Six tools; list-tasks **out** |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF). Do not rewrite Phase 09 `done` prompts. Do not invent Phase 10 without explicit promotion. On APPROVE: mark Phase 09 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 09 complete; no accidental Phase 10 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; dogfood ladder may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs DF-01 + S02/S03 + suite evidence (+ optional fresh runs)
- [ ] Confirm DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 09 complete)

## Out of scope
- Inventing Phase 10 without promotion
- Weakening prior gates
- Requiring MCP list-tasks
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments
