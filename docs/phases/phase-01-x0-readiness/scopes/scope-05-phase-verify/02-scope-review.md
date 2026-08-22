# P01 / S05 / 02 — Scope review (Phase 01 VERIFY)

## Metadata
- id: P01-S05-02
- todo_ids: [P01-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S05 (Phase 01 VERIFY). Confirm `VERIFY-NOTES.md` + board claims match a **fresh** suite re-run (or honest fail+spawn trail). **Complete DR-HANDOFF** for Phase 02 before marking this row `done`. Severity-tag findings; small doc/scaffold fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only. Do **not** claim Gate C.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop + Phase handoff)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 2 Gate C
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — X0 vs Gate C
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 00 VERIFY review [`../../../phase-00-foundation/scopes/scope-09-phase-verify/02-scope-review.md`](../../../phase-00-foundation/scopes/scope-09-phase-verify/02-scope-review.md)

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S05-01).

## Review focus
- Did VERIFY **independently** re-run honesty + x0 + p0x + `./...` (not copy S02–S04 Notes)?
- Evidence table covers honesty Paths A/B/C, X0 B0+G1 (`dry_run:true`), p0x 7/7, MCP checklist, law checks?
- **No Gate C** / “G1 beats B0” / A1-validated claim?
- MCP: six tools + stdio `trace-mcp` + G19; X0 still CLI-without-MCP?
- On fail: remediations spawned with full prompts; bars not weakened?
- Residuals listed — not silently ignored if they undermine confidence
- **DR-HANDOFF (this row owns completion):**
  - [ ] `docs/phases/phase-02-gate-c/README.md` exists (goal = Gate C evaluation & slice hardening)
  - [ ] `00-PHASE-PLANNER.md` runnable (session-start + exit criteria)
  - [ ] At least one scope stub with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` (minimal OK)
  - [ ] `docs/TODO.md` has Phase 02 section; first pending row is **`P02-00`** after Phase 01’s last `done` row
  - [ ] Not README-only / blocked-until-noticed (DR-HANDOFF — do not repeat Phase 00→01 gap)
- If handoff incomplete: **finish it here** (reviewer rights on upcoming artifacts) or spawn — do **not** mark `done` until complete

## Locked expectations

| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Re-check commands (spot-check ≥1, prefer all) | `CGO_ENABLED=0 go test ./evals/honesty/... -count=1`; `CGO_ENABLED=1 go test ./evals/x0/... ./evals/p0x/... ./... -count=1` |
| Phase 02 path | `docs/phases/phase-02-gate-c/` |
| First Phase 02 board row | `P02-00` → `00-PHASE-PLANNER.md` |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |
| Product Go | Forbidden except tiny doc/scaffold; no feature work |

## Board rights
Reviewer: status + notes; **may spawn** remediations; **may thicken / complete upcoming Phase 02** scaffold and board rows (DR-HANDOFF). Do not rewrite Phase 01 `done` prompts. Do not execute Phase 02 implement wave.

## Exit criteria
- [ ] Findings recorded (`REVIEW-NOTES.md` preferred)
- [ ] blocker/high fixed or spawned
- [ ] VERIFY evidence re-checked (fresh commands or justified trust of VERIFY-NOTES with spot-check)
- [ ] **DR-HANDOFF complete** before this row `done`
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated — Phase 01 complete; next runnable is **`P02-00`**
- [ ] Explicit: Phase 01 ≠ Gate C pass

## Minimal todos
- [ ] Compare VERIFY claims vs evidence + optional fresh harness runs
- [ ] Finish Phase 02 scaffold / board if incomplete (DR-HANDOFF)
- [ ] Fix docs or spawn
- [ ] Write REVIEW-NOTES.md; board update (Phase 01 complete)
