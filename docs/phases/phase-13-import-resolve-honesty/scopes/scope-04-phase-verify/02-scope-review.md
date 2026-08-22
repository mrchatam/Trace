# P13 / S04 / 02 — Scope review (Phase 13 VERIFY / phase close)

## Metadata
- id: P13-S04-02
- todo_ids: [P13-S04-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (Phase 13 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Phase 13 passed without **S01–S03 named DF tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing plan/impact/index MCP as a VERIFY requirement. Reject inventing product INFERRED (DF-66) or symbol honesty (DF-67). Spot-check: S01 DF-60 + S02 DF-61/62/63/65 + S03 DF-64 + P12 keepers + honesty/Gates/ablation/compat/p0x/x0/`./...`; Gate C `dry_run:false` intact; DF-66/67 Notes explicit.

Do **not** rewrite Phase 12 `done` history. Phase 12 historical `no successor` stays intact as history — Phase 13 was a **forward** human reopen.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-import-path-resolve/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-packet-honesty-residuals/REVIEW-NOTES.md)
- [S03 REVIEW-NOTES.md](../scope-03-provenance-schema/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Research deferrals: [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md)
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 12 VERIFY review [`../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/02-scope-review.md`](../../../phase-12-peer-honesty-surfaces/scopes/scope-03-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S04-01).

## Review focus
- Did VERIFY independently run **S01** DF-60 named import-path/Expand/Why tests — not copy Notes?
- Did VERIFY re-prove **S02** DF-61/62/63/65 named tests + P12 Budget/cap/stale/WhyTrace keepers?
- Did VERIFY re-prove **S03** DF-64 named store/Expand/analyzer + compat **12**, and record DF-66 wontfix + DF-67 `symstale/`?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat + p0x + x0 + product `./...`?
- Evidence table covers S01–S03 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (DF-66/67; TaskContext DF-65 shared-path; graphify path FAIL; CGO0 analyzers; research ranks 4+) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- Optional ab-import-resolve probe (if present) treated as dogfood only — **not** Gate C?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim a Phase 14 / research ranks 4+ scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 14 artifacts required — confirm absence is intentional; remaining research/dogfood stays parallel / research-only
  - [ ] Forward-only: Phase 12 historical `no successor` left intact as history
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 13 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01–S03 regressions evidenced via **named tests** (+ P12 keepers) |
| 2 | DF-66 wontfix + DF-67 residual explicit in VERIFY-NOTES |
| 3 | Carry-forward gates green; dry-run ≠ claims intact |
| 4 | Product `./...` PASS (known non-product graphify path OK; CGO0 analyzers OK residual) |
| 5 | DR-HANDOFF closed (`no successor` or promoted successor scaffold) |
| 6 | AGENTS/TODO/phase README consistent with closeout |
| 7 | No silent residuals; Phase 12 history not rewritten; research ranks 4+ not auto-boarded |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01–S03 path | Named tests per `01-verify.md` locks |
| MCP | **Nine** tools; plan/impact/index MCP dump **out**; no provenance MCP/CLI; no Phase 13 tool menu |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Compat | Ceiling **12** (mig 012); no 013+ |
| `./...` | Product pkgs PASS; graphify space FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + AGENTS/README focus). Do not rewrite Phase 13 `done` prompts. Do not invent Phase 14 without explicit promotion. On APPROVE: mark Phase 13 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 13 complete; no accidental Phase 14 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; research/dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S03 + suite evidence (+ optional fresh runs)
- [ ] Confirm DF-66/67 Notes + DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 13 complete)

## Out of scope
- Inventing Phase 14 / research ranks 4+ without promotion
- Weakening prior gates
- Requiring plan/impact/index MCP
- Treating ab-import-resolve as Gate C
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments (method notes)
- Rewriting Phase 00–12 `done` history
