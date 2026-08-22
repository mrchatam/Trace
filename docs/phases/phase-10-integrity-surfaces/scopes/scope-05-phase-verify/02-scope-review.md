# P10 / S05 / 02 — Scope review (Phase 10 VERIFY / phase close)

## Metadata
- id: P10-S05-02
- todo_ids: [P10-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S05 (Phase 10 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Phase 10 passed without **S01–S04 named DF tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing plan/impact/index MCP as a VERIFY requirement. Spot-check: DF-19/23/25/27/29 + nine MCP tools + DF-20 GC + DF-17/18/24/26/31 + honesty Path C / Gate G; honesty/Gates/ablation/compat/p0x/x0/`./...`; Gate C `dry_run:false` intact.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-retrieval-why-fidelity/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-mcp-parity-install/REVIEW-NOTES.md)
- [S03 REVIEW-NOTES.md](../scope-03-index-gc/REVIEW-NOTES.md)
- [S04 REVIEW-NOTES.md](../scope-04-operator-capability-gates/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 09 VERIFY review [`../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/02-scope-review.md`](../../../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S05-01).

## Review focus
- Did VERIFY independently run **S01** DF-19/23/25/27/29 named tests — not copy Notes?
- Did VERIFY re-prove **S02** nine MCP tools + DF-21/22/32 (G19; no plan/impact/index MCP)?
- Did VERIFY re-prove **S03** `TestIndexGCAfterPathRename` + argv isolation + `TestIndexIncrementalIsolation`?
- Did VERIFY re-prove **S04** DF-17/18/24/26/31 + honesty Path C operator flag + Gate G hatch?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat + p0x + x0 + product `./...`?
- Evidence table covers S01–S04 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (`plan_scope` Exact out; Mode-B historical; Cursor MCP reload; graphify path FAIL; optional ab-* re-runs) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim a Phase 11 scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 11 artifacts required — confirm absence is intentional; remaining dogfood stays parallel `experiments/`
  - [ ] Forward-only: Phase 09 historical `no successor` left intact as history
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 10 complete; follow spawn trail

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01 path | retrieval/compiler DF-19/23/25/27/29 named tests |
| S02 / S03 / S04 | Nine MCP + GC rename/isolation + operator/capability named tests |
| MCP | **Nine** tools; plan/impact/index MCP **out** |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| `./...` | Product pkgs PASS; graphify space FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF). Do not rewrite Phase 10 `done` prompts. Do not invent Phase 11 without explicit promotion. On APPROVE: mark Phase 10 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 10 complete; no accidental Phase 11 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S04 + suite evidence (+ optional fresh runs)
- [ ] Confirm DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 10 complete)

## Out of scope
- Inventing Phase 11 without promotion
- Weakening prior gates
- Requiring plan/impact/index MCP
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments
- Rewriting Phase 09 `done` history
