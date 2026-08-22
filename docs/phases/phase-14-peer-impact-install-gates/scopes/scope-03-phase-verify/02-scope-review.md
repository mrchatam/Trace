# P14 / S03 / 02 — Scope review (Phase 14 VERIFY / phase close)

## Metadata
- id: P14-S03-02
- todo_ids: [P14-S03-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Phase 14 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Reject any claim that Phase 14 passed without **S01–S02 named tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing plan/impact/install MCP as a VERIFY requirement. Reject claiming “every MCP call gated” from Assert library/CLI alone. Spot-check: S01 ImpactWalk named + Gate F + S02 install/decision named + Cursor keepers + ablation + honesty/Gates/compat/p0x/x0/`./...`; Gate C `dry_run:false` intact; Assert≠MCP + optional allowContainsOut Notes explicit.

Do **not** rewrite Phase 13 `done` history. Phase 13 historical `no successor` stays intact as history — Phase 14 was a **forward** human reopen.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-impact-walks/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-install-capability-gates/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Goals sequence: [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — #2 S05 / #3 plan simulate / #4 D21+ stay off-board by default
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 13 VERIFY review [`../../../phase-13-import-resolve-honesty/scopes/scope-04-phase-verify/02-scope-review.md`](../../../phase-13-import-resolve-honesty/scopes/scope-04-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S03-01).

## Review focus
- Did VERIFY independently run **S01** ImpactWalk named tests + Gate F — not copy Notes?
- Did VERIFY re-prove **S02** install named + `TestCapabilityDecision*` + `TestInstallCursor*` + ablation?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat **13** + p0x + x0 + product `./...`?
- Evidence table covers S01–S02 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (Assert≠MCP; optional allowContainsOut; graphify path FAIL; CGO0 analyzers; goals #2–#4) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim a Phase 15 / S05 / plan simulate / D21+ scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 15 artifacts required — confirm absence is intentional; remaining research/dogfood stays parallel / research-only
  - [ ] Forward-only: Phase 13 historical `no successor` left intact as history
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 14 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01+S02 named regressions re-proven (ImpactWalk + install/decision + Cursor keepers + Gate F + ablation) |
| 2 | Carry-forward honesty/E–H/ablation/compat **13**/p0x/x0 green |
| 3 | Gate C `dry_run:false` intact |
| 4 | Dry-run ≠ Gate C/F/G/ablation/H/checklist |
| 5 | DR-HANDOFF closed per FINAL (default `no successor`) |
| 6 | No auto-board of S05 / plan simulate / D21+ / Phase 15 |
| 7 | Phase 13 historical `no successor` left intact |
| 8 | Assert≠MCP honesty Note present; S02 APPROVE ≠ MCP request-path gating |
| 9 | Board + AGENTS/README synced on complete |
| 10 | No silent residuals; product `./...` PASS (graphify space FAIL OK) |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01–S02 path | Named tests per `01-verify.md` locks |
| MCP | **Nine** tools + `trace_version`; install/decide MCP dump **out**; Assert ≠ MCP dispatch |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat **13**; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Compat | Ceiling **13** (mig 013); no 014+ |
| `./...` | Product pkgs PASS; graphify space FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + AGENTS/README focus). Do not rewrite Phase 14 `done` prompts. Do not invent Phase 15 without explicit promotion. On APPROVE: mark Phase 14 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 14 complete; no accidental Phase 15 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; research/dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S02 + suite evidence (+ optional fresh runs)
- [ ] Confirm Assert≠MCP Notes + DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 14 complete)

## Out of scope
- Inventing Phase 15 / S05 / plan simulate / D21+ without promotion
- Weakening prior gates
- Requiring install/decide MCP
- Wiring Assert into MCP dispatch as a VERIFY requirement
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments (method notes)
- Rewriting Phase 00–13 `done` history
