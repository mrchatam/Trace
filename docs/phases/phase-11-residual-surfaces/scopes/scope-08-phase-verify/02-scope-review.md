# P11 / S08 / 02 — Scope review (Phase 11 VERIFY / phase close)

## Metadata
- id: P11-S08-02
- todo_ids: [P11-S08-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S08 (Phase 11 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

Reject any claim that Phase 11 passed without **S01–S07 named DF tests** and carry-forward bars. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject inventing plan/impact/index MCP as a VERIFY requirement. Spot-check: DF-40 + DF-43/44 + DF-47 + DF-41/51 + DF-49/35/48/42 + DF-50/22/37 + DF-33/30/46/45/28 + honesty/Gates/ablation/compat/p0x/x0/`./...`; Gate C `dry_run:false` intact.

Do **not** rewrite Phase 10 `done` history. Phase 10 historical `no successor` stays intact as history — Phase 11 was a **forward** reopen.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [S01 REVIEW-NOTES.md](../scope-01-index-partial-path-gc/REVIEW-NOTES.md)
- [S02 REVIEW-NOTES.md](../scope-02-review-pass-fail-operator/REVIEW-NOTES.md)
- [S03 REVIEW-NOTES.md](../scope-03-store-lock-concurrency/REVIEW-NOTES.md)
- [S04 REVIEW-NOTES.md](../scope-04-capability-upsert-hatch/REVIEW-NOTES.md)
- [S05 REVIEW-NOTES.md](../scope-05-retrieval-why-depth-trust/REVIEW-NOTES.md)
- [S06 REVIEW-NOTES.md](../scope-06-mcp-install-reload/REVIEW-NOTES.md)
- [S07 REVIEW-NOTES.md](../scope-07-seed-plan-review-polish/REVIEW-NOTES.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Dogfood: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 10 VERIFY review [`../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/02-scope-review.md`](../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S08-01).

## Review focus
- Did VERIFY independently run **S01** DF-40 named tests (+ DF-20 retain) — not copy Notes?
- Did VERIFY re-prove **S02** DF-43 sibling FAIL+PASS + DF-44 conscious-flag + honesty Path C?
- Did VERIFY re-prove **S03** DF-47 retry + exclusivity + serialize guidance?
- Did VERIFY re-prove **S04** DF-41 slug upsert + DF-51 hatch≠missing-caps?
- Did VERIFY re-prove **S05** DF-49/35/48/42 named tests?
- Did VERIFY re-prove **S06** tip parity + nine tools/`trace_version` (no PID kill)?
- Did VERIFY re-prove **S07** DF-33/30/46/45/28 named tests?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat + p0x + x0 + product `./...`?
- Evidence table covers S01–S07 + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- Residuals (rename+edit ghost; graphify path FAIL; CGO0 analyzers; Cursor MCP reload; optional ab-* re-runs) noted as non-blocking?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` do **not** claim a Phase 12 scaffold unless Notes promoted one
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist and are runnable (finish here if incomplete)
  - [ ] Default path: **no** Phase 12 artifacts required — confirm absence is intentional; remaining dogfood stays parallel `experiments/`
  - [ ] Forward-only: Phase 10 historical `no successor` left intact as history
  - [ ] Findings: 18 P11 DFs closed or explicitly residual-listed (if S08-01 deferred findings flip, finish here)
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 11 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01–S07 DF regressions evidenced via **named tests** |
| 2 | Carry-forward gates green; dry-run ≠ claims intact |
| 3 | Product `./...` PASS (known non-product graphify path OK; CGO0 analyzers OK residual) |
| 4 | DR-HANDOFF closed (`no successor` or promoted Phase 12 scaffold) |
| 5 | AGENTS/TODO/findings consistent with closeout |
| 6 | No silent residuals; Phase 10 history not rewritten |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01–S07 path | Named DF tests per `01-verify.md` locks |
| MCP | **Nine** tools; plan/impact/index MCP dump **out** |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; G1 0.800 > B0 0.000 — re-check, no new Go |
| `./...` | Product pkgs PASS; graphify space FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + findings flip). Do not rewrite Phase 11 `done` prompts. Do not invent Phase 12 without explicit promotion. On APPROVE: mark Phase 11 complete; next runnable = **none** (roadmap closed again) unless promoted.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 11 complete; no accidental Phase 12 without promotion
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact; dogfood may continue off-board

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S07 + suite evidence (+ optional fresh runs)
- [ ] Confirm DR-HANDOFF = `no successor` (or finish promoted successor scaffold)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 11 complete)

## Out of scope
- Inventing Phase 12 without promotion
- Weakening prior gates
- Requiring plan/impact/index MCP
- Executing product feature work outside review/spawn rights
- Closing parallel dogfood experiments (method notes)
- Rewriting Phase 00–10 `done` history
