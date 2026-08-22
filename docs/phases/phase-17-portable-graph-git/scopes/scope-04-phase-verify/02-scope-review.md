# P17 / S04 / 02 — Scope review (Phase 17 VERIFY / phase close)

## Metadata
- id: P17-S04-02
- todo_ids: [P17-S04-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (Phase 17 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** as **`no successor`** before marking this row `done` (unless Notes explicitly promote a follow-on **and** that follow-on is fully scaffolded). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

Reject any claim that Phase 17 passed without **S01–S03 named DF tests** and **two-clone recipe**. Reject treating Phase 01 dry-run as Gate C, Gate F, Gate G, ablation, Gate H, or checklist. Reject failing VERIFY for encryption wontfix / reviews omitted / **DF-86 hook absent** / CGO=0 `cmd/trace` FAIL. Reject inventing hosted MCP / research S05 / plan simulate / D21+ as successors. Spot-check: S01–S03 named + two-clone + carry-forward honesty/Gates/compat/p0x/x0/product pkgs; Gate C `dry_run:false` intact; DF-86 absence recorded as non-blocking.

Do **not** rewrite Phase 16 `done` history. Phase 16 historical `no successor` stays intact as history.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- S01–S03 REVIEW-NOTES (all APPROVE high)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- Goals sequence: [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — #2 S05 / #3 plan simulate / #4 D21+ stay off-board by default
- Gate C: [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md), [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 16 VERIFY review [`../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/02-scope-review.md`](../../../phase-16-assert-root-and-surfaces/scopes/scope-06-phase-verify/02-scope-review.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S04-01). Unattended: no Plan-mode switch.

## Review focus
- Did VERIFY independently run **S01–S03** named tests (not copy Notes)?
- Did VERIFY prove **two-clone recipe** — no shared `.trace/`; why + context + plan; offline?
- Did VERIFY re-run honesty A/B/C + Gate G/E/F + ablation + Gate H + compat + p0x + x0 + product `./cmd|internal|evals`?
- Evidence table covers S01–S03 + two-clone + Gate C `dry_run:false` + dry-run≠Gate C/≠F/≠G/≠ablation/≠Gate H/≠checklist + law checks + handoff?
- DF-86 hook absence noted as **non-blocking** — and **not** used as fail criteria?
- `.gitignore` `.trace/` only; no new MCP seed tool; no HTTP/daemon/hosted MCP?
- On fail: remediations spawned with full prompts; bars not weakened?
- **DR-HANDOFF (this row owns completion):**
  - [ ] `VERIFY-NOTES.md` explicitly records **`no successor`**
  - [ ] [`DR-HANDOFF.md`](../../DR-HANDOFF.md) closed / stamped
  - [ ] Board / phase README / `AGENTS.md` synced on Phase 17 complete
  - [ ] If Notes **did** promote a successor: folder + `00-PHASE-PLANNER` + ≥1 scope stub + board first pending row exist (default path does **not** promote)
  - [ ] Default path: research S05 / plan simulate / D21+ / hosted MCP **not** auto-boarded
  - [ ] Forward-only: Phase 16 historical `no successor` left intact
- If handoff incomplete / ambiguous: **finish documentation here** (reviewer rights) or spawn — do **not** mark `done` until complete
- If VERIFY failed without recovery: do not mark Phase 17 complete; follow spawn trail

## Checklist

| # | Check |
|---|--------|
| 1 | S01–S03 named DF tests re-proven (not Notes-only) |
| 2 | Two-clone recipe re-proven — `TestPortableGraphTwoCloneWhyContextPlan` + VERIFY-NOTES clone detail |
| 3 | Carry-forward honesty/E–H/ablation/compat/p0x/x0/product pkgs green |
| 4 | Gate C `dry_run:false` intact |
| 5 | Dry-run ≠ Gate C/F/G/ablation/H/checklist |
| 6 | DF-86 hook absence **non-fail** — grep evidence in VERIFY-NOTES |
| 7 | `.gitignore` `.trace/` only; `trace/graph.json` not ignored |
| 8 | No new MCP seed tool; no HTTP/daemon/hosted MCP |
| 9 | DR-HANDOFF closed per FINAL (default `no successor`) |
| 10 | No auto-board of research S05 / plan simulate / D21+ / hosted MCP |
| 11 | Phase 16 history intact; P16 DFs not re-claimed |
| 12 | Encryption / reviews omitted / CGO0 cmd/trace **not** claimed fixed; **not** used as fail criteria |
| 13 | Board + AGENTS/README synced on complete; product pkgs PASS (graphify space / CGO0 analyzers FAIL OK) |

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| S01–S03 path | Named tests per `01-verify.md` / `00-PLANNER.md` locks |
| Two-clone | `TestPortableGraphTwoCloneWhyContextPlan` + shell corroboration in VERIFY-NOTES |
| Carry-forward | Honesty A/B/C; Gate G/E/F; ablation; Gate H; compat; p0x; x0 |
| Gate C integrity | `GATE-C-NOTES.md` **Go**; metrics `dry_run:false`; — re-check, no new Go |
| Product bar | `./cmd\|internal\|evals` PASS; R3 graphify / R4 CGO0 FAIL = known residual OK |
| Re-check commands (spot-check ≥1, prefer all) | Locked commands from `01-verify.md` |
| Successor | **`no successor`** unless explicit promotion — hosted MCP is **not** that promotion |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc fixes / spawn remediation |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may complete handoff docs** (DR-HANDOFF + AGENTS/README focus). Do not rewrite Phase 17 `done` prompts. On APPROVE: mark Phase 17 complete; handoff text stays **`no successor`**.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done` (`no successor` **or** explicit promoted successor scaffolded)
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 17 complete; next runnable **none** (unless promoted)
- [ ] Explicit: Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist; Gate C artifacts intact

## Minimal todos
- [ ] Compare VERIFY claims vs S01–S03 + two-clone + suite evidence (+ optional fresh runs)
- [ ] Confirm DF-86 non-blocking + DR-HANDOFF = `no successor`
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; close DR-HANDOFF; board update (Phase 17 complete)

## Out of scope
- Inventing research S05 / plan simulate / D21+ / hosted MCP without promotion
- Weakening prior gates
- Failing VERIFY for DF-86 / encryption / reviews omitted / CGO0 cmd/trace
- Executing product feature work outside review/spawn rights
- Rewriting Phase 00–16 `done` history
