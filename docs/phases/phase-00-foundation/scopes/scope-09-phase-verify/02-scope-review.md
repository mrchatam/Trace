# P00 / S09 / 02 — Scope review (Phase 00 VERIFY / P0 close)

## Metadata
- id: P00-S09-02
- todo_ids: [P00-S09-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S09 VERIFY / P0 close. Confirm `VERIFY-NOTES.md` + board claims match a **fresh** harness re-run (or honest fail+spawn trail). Severity-tag findings; small doc fixes or spawn `a`/`b` pairs with **full** prompts. Forward-only. May thicken **upcoming** Phase 01 stubs only if P0 close surface needs a pointer.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling `01-verify.md` + `VERIFY-NOTES.md` + board Notes
- [docs/init/C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — 7/7 bar
- [../scope-08-fixture-p0x/REVIEW-NOTES.md](../scope-08-fixture-p0x/REVIEW-NOTES.md) — residuals

## Session start
Agent → clarify → Plan → review (fresh subagent ≠ S09-01).

## Review focus
- Did VERIFY **independently** re-run `CGO_ENABLED=1 go test ./evals/p0x/...` + `./...` (not copy S08 Notes)?
- Evidence table covers all 7 criteria with real subtest/log gists?
- Law checks: no MCP/daemon/HTTP; no committed `.trace/`; #7 sibling isolation; G19 library≠CLI?
- On pass: F_QUESTION_LEDGER / E A15 (and A17 if claimed) updated honestly — A1 still EXPERIMENT_REQUIRED?
- On fail: remediations spawned with full prompts; bar not weakened; no MCP “fix”?
- Residuals (soft decision-constraint OR, JSON panics) listed — not silently ignored if they undermine confidence
- Phase 01 still blocked until this review `done`

## Locked expectations
| Item | Value |
|------|-------|
| Gate artifact | `VERIFY-NOTES.md` in this folder |
| Re-check command | `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` (reviewer should spot-check) |
| Confidence bar | Prefer **high**; **medium** only with explicit residuals |

## Exit criteria
- [ ] Findings recorded (REVIEW-NOTES.md recommended)
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated; Phase 00 closable only if VERIFY + this review agree

## Minimal todos
- [ ] Compare VERIFY claims vs evidence + optional fresh harness run
- [ ] Fix docs or spawn
- [ ] Re-verify
- [ ] Board update
