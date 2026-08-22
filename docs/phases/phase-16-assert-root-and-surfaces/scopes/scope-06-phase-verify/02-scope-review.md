# P16 / S06 / 02 — Scope review (Phase 16 VERIFY / phase close)

## Metadata
- id: P16-S06-02
- todo_ids: [P16-S06-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S06 (Phase 16 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Reject any claim that Phase 16 passed without **S01–S05 named DF tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing install/decide/plan/index MCP as a VERIFY requirement. Reject failing VERIFY for DF-67 defer / P14 R2 / P15 R3–R4 / S05-02 `attachTaskImpact` swallow / 014 nine-Name list. Reject claiming DF-72 still deferred (live lock is thin `trace_impact` — [`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md)). Spot-check: S01–S05 named + catalog **10** including `trace_version` + carry-forward honesty/Gates/compat **14**/p0x/x0/product pkgs; Gate C `dry_run:false` intact; residuals Notes explicit as non-blocking.

Do **not** rewrite Phase 15 `done` history. Phase 15 historical `no successor` stays intact as history — Phase 16 was a **forward** human reopen. **Phase 17** is independently queued on the board — do **not** rewrite those rows; do **not** claim P17 as this VERIFY successor.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- S01–S05 REVIEW-NOTES (all APPROVE high)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-72-FORWARD.md](../../DF-72-FORWARD.md)
- Goals sequence: [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — #2 S05 / #3 plan simulate / #4 D21+ stay off-board by default
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 15 VERIFY review [`../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/02-scope-review.md`](../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S06-01). Unattended: no Plan-mode switch.

## Review focus
- Did VERIFY independently run **S01–S05** named tests (not copy Notes)? DF-72 `TestMCPTraceImpactReport` / `TestMCPImpactDeniedBlocksCallTool` **is** a fail bar.
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat **14** + p0x + x0 + product `./cmd|internal|evals`?
- Catalog **10** including `trace_version`; `trace_impact` in catalog; **no** install/decide/plan/index MCP?
- Evidence table covers S01–S05 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (DF-67; R2; R3/R4; `attachTaskImpact`; 014 nine-Name) noted as non-blocking — and **not** used as fail criteria?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim Phase 17 as P16 VERIFY successor
  - [ ] Phase 17 rows **232–244 left intact** (independently queued; first P17 implement after every P16 row `done`)
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete) — default path does **not** promote
  - [ ] Default path: remaining research/dogfood stays parallel / research-only; **no** auto-board of research S05 / plan simulate / D21+
  - [ ] Forward-only: Phase 15 historical `no successor` left intact as history
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 16 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01–S05 named DF tests re-proven (not Notes-only); DF-72 named test is a fail bar |
| 2 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0/product pkgs green |
| 3 | Gate C `dry_run:false` intact |
| 4 | Dry-run ≠ Gate C/F/G/ablation/H/checklist |
| 5 | Ten MCP tools + `trace_version`; thin `trace_impact` present; **no** install/decide/plan/index MCP |
| 6 | DR-HANDOFF closed per FINAL (default `no successor`) |
| 7 | No auto-board of research S05 / plan simulate / D21+; Phase 17 **not** rewritten or claimed as this successor |
| 8 | Phase 15 historical `no successor` left intact |
| 9 | DF-67 / R2 / R3 / R4 / S05-02 swallow / 014 nine-Name **not** claimed fixed; **not** used as fail criteria |
| 10 | Board + AGENTS/README synced on complete; product pkgs PASS (graphify space / CGO0 analyzers FAIL OK) |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01–S05 path | Named tests per `01-verify.md` / `00-PLANNER.md` locks |
| MCP | **Ten** tools including `trace_version`; `trace_impact` **in**; install/decide/plan/index MCP dump **out** |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat **14**; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| Compat | Ceiling **14** (mig 014); no 015+ |
| Product bar | `./cmd\|internal\|evals` PASS; R3 graphify / R4 CGO0 FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion — Phase 17 is **not** that promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + AGENTS/README focus). Do not rewrite Phase 16 `done` prompts. Do not rewrite Phase 17 rows. On APPROVE: mark Phase 16 complete; next runnable = **P17-S01-00** only because P17 was independently queued **before** this VERIFY — that is **not** DR-HANDOFF promotion. Handoff text stays **`no successor`**.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 16 complete; Phase 17 rows unchanged except “P16 in progress” → complete so P17-S01-00 becomes runnable
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; research/dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S05 + suite evidence (+ optional fresh runs)
- [ ] Confirm residuals non-blocking + DR-HANDOFF = `no successor` (Phase 17 not this successor)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 16 complete)

## Out of scope
- Inventing research S05 / plan simulate / D21+ without promotion
- Rewriting Phase 17 prompt bodies or claiming P17 as P16 VERIFY successor
- Weakening prior gates
- Requiring install/decide/plan/index MCP
- Failing VERIFY for DF-67 / R2 / R3 / R4 / attachTaskImpact swallow
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments (method notes)
- Rewriting Phase 00–15 `done` history
