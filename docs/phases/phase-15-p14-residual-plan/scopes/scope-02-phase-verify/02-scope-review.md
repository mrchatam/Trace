# P15 / S02 / 02 — Scope review (Phase 15 VERIFY / phase close)

## Metadata
- id: P15-S02-02
- todo_ids: [P15-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S02 (Phase 15 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Reject any claim that Phase 15 passed without **S01 named MCP Assert tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing install/decide MCP as a VERIFY requirement. Reject failing VERIFY for R2 defer / R3–R4 wontfix. Spot-check: S01 MCP Assert named + tool count + carry-forward honesty/Gates/compat/p0x/x0/product pkgs; Gate C `dry_run:false` intact; R2/R3/R4 Notes explicit as non-blocking.

Do **not** rewrite Phase 14 `done` history. Phase 14 historical `no successor` stays intact as history — Phase 15 was a **forward** human reopen.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-mcp-assert-dispatch/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Goals sequence: [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — #2 S05 / #3 plan simulate / #4 D21+ stay off-board by default
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 14 VERIFY review [`../../../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S02-01).

## Review focus
- Did VERIFY independently run **S01** MCP Assert named tests (`TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered`) — not copy Notes?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat **13** + p0x + x0 + product `./cmd|internal|evals`?
- Evidence table covers S01 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (R2 defer; R3/R4 wontfix; goals #2–#4) noted as non-blocking — and **not** used as fail criteria?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim a Phase 16 / S05 / plan simulate / D21+ scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 16 artifacts required — confirm absence is intentional; remaining research/dogfood stays parallel / research-only
  - [ ] Forward-only: Phase 14 historical `no successor` left intact as history
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 15 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01 MCP Assert named regressions re-proven (DENIED + AUTO_ALLOWED + tool names + specs) |
| 2 | Carry-forward honesty/E–H/ablation/compat **13**/p0x/x0/product pkgs green |
| 3 | Gate C `dry_run:false` intact |
| 4 | Dry-run ≠ Gate C/F/G/ablation/H/checklist |
| 5 | DR-HANDOFF closed per FINAL (default `no successor`) |
| 6 | No auto-board of S05 / plan simulate / D21+ / Phase 16 |
| 7 | Phase 14 historical `no successor` left intact |
| 8 | R2/R3/R4 disposition unchanged; VERIFY did **not** fail for them |
| 9 | Board + AGENTS/README synced on complete |
| 10 | No silent residuals; product pkgs PASS (graphify space / CGO0 analyzers FAIL OK) |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01 path | Named tests per `01-verify.md` locks |
| MCP | **Nine** tools + `trace_version`; install/decide MCP dump **out**; Assert **is** on MCP dispatch (R1) |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat **13**; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Compat | Ceiling **13** (mig 013); no 014+ |
| Product bar | `./cmd\|internal\|evals` PASS; R3 graphify / R4 CGO0 FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + AGENTS/README focus). Do not rewrite Phase 15 `done` prompts. Do not invent Phase 16 without explicit promotion. On APPROVE: mark Phase 15 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 15 complete; no accidental Phase 16 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; research/dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01 + suite evidence (+ optional fresh runs)
- [ ] Confirm R2/R3/R4 non-blocking + DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 15 complete)

## Out of scope
- Inventing Phase 16 / S05 / plan simulate / D21+ without promotion
- Weakening prior gates
- Requiring install/decide MCP
- Fixing R2 / R3 / R4 as a VERIFY requirement
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments (method notes)
- Rewriting Phase 00–14 `done` history
