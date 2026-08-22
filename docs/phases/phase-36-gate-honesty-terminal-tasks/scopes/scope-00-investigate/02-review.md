# P36-S00-02 — Review investigation

## Metadata
- id: P36-S00-02
- todo_ids: [P36-S00-02]
- role: reviewer
- skills: [code-review-and-quality, diagnosing-bugs, systematic-debugging]
- mcps: [Shell via trace CLI for re-verify]
- agents: []
- verification: automated
- hooks: []

## Objective

Independent review of `INVESTIGATION.md` vs DESIGN-LOCKS + INTAKE + live repo. Re-run spot-check CLI commands; confirm file:line cites match current code. **No product code.** Spawn `P36-S00-02a`/`02b` only if blocker/high gaps remain after one fix pass.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [01-investigate.md](01-investigate.md)
- Deliverable: `INVESTIGATION.md`
- Live anchors: `internal/loop/policy.go`, `internal/loop/gate.go`, `internal/deliberation/select.go`, `cmd/trace/plan.go`, `internal/mcp/server.go`, `cmd/trace/transition.go`, `internal/install/enforcement.go`, `web/src/screens/TaskDetail.tsx`, `web/src/components/GateStrip.tsx`

## Session start

Follow agent-loop-protocol Session start. Fresh subagent — do not assume implementer session context.

## Locked defaults

| Item | Value |
|------|-------|
| Fixture | `/home/ali/Desktop/feet seller telegram app` (read-only) |
| Re-verify tasks | Step 1 `33247e2d-…`, Loop 112 `99d8fb92-…` |
| PASS bar | High confidence, or medium with explicit residual risks |
| Spawn threshold | Blocker/high without pending follow-up row |

## Review procedure

1. Read INVESTIGATION.md end-to-end.
2. Re-run at least: both `loop gate --for done`, `plan show --goal`, `seed export` count check.
3. Spot-read cited source lines — reject stale line numbers without correction in doc.
4. Score findings: blocker | high | medium | low | nit.
5. Blocker/high → small doc fix in INVESTIGATION.md **or** spawn implement+review pair (prefer doc fix if investigation content only).

## Checklist vs DESIGN-LOCKS + INTAKE

### Verdict and blame (INTAKE + DESIGN-LOCKS theme)

- [ ] **Verdict table** present: Trace product / agent misuse / harness — each with evidence and primary blame assignment
- [ ] Conclusion matches DESIGN-LOCKS theme: **fundamental planning model**, not GUI-only patch
- [ ] Does **not** accept "agents should have known" without MCP/install/product gap analysis

### Two planning systems (DESIGN-LOCKS must-investigate)

- [ ] **plan-change count 11** vs **progressive planner 0** — proved via export (and/or DB/show), not assertion
- [ ] Explains disconnected signals: plan-changes satisfy `PlanCritiqued` only; `PlanExists` reads progressive planner only
- [ ] Live cites: `policy.go:45–49` (PlanExists), `policy.go:60–62` (PlanCritiqued), `deliberation/select.go:28–29` (PLAN/plan_missing)

### Gate behavior (DESIGN-LOCKS must-fix honesty + must-preserve)

- [ ] **Identical gate** on Step 1 and Loop 112 explained (goal-level `PlanExists`, not task terminal state)
- [ ] **`evaluateDone` gap** documented: no terminal `work_state` short-circuit (`gate.go:227–265`)
- [ ] **Must preserve** stated: PLAN gate for **active** non-terminal work — doc does not propose global weaken
- [ ] Terminal DONE + misleading "done blocked" copy acknowledged as secondary (GUI section)

### Transition without gate (INTAKE hypothesis B + C)

- [ ] **123 DONE** path documented: enforce off (missing config), opt-in `--enforce` on transition, MCP `trace_transition` without gate
- [ ] Review PASS + `as_operator` path cited (`internal/mcp/tools_write.go`, `cmd/trace/transition.go`)
- [ ] Explicit: gates surfaced when human opened GUI / ran CLI gate — not during historical agent DONE wave

### MCP + install (DESIGN-LOCKS must-investigate + must-fix product)

- [ ] **MCP surface audit**: lists registered tools; confirms **no** plan create-coarse/set-current/deep/show
- [ ] `trace_add kind=plan-change` documented; `trace_loop` ≠ gate
- [ ] `cmd/trace/plan.go:67` "No MCP plan tools" cited
- [ ] **Install audit**: feet-seller missing AGENTS.md, `.trace/config.json`, cursor rules/hooks — impact on enforce + discoverability
- [ ] Install rules gap: no `create-coarse` bootstrap in `internal/install/enforcement.go`

### Mega-goal + fixes (INTAKE D + DESIGN-LOCKS candidates)

- [ ] **123 tasks / 1 goal** pattern analyzed as amplification factor
- [ ] **Fundamental fix options** ranked for S01 (MCP plan, bootstrap, install, PlanExists bridge, terminal gate, etc.)
- [ ] Options respect **Law 19** (library policy first; MCP/HTTP/GUI adapters)

### GUI secondary (INTAKE + DESIGN-LOCKS must-fix honesty)

- [ ] **TaskDetail** path cited: always loads done gate (`TaskDetail.tsx:76–86`, `:198–210`)
- [ ] **GateStrip** error styling for blocked gate (`GateStrip.tsx:41–56`)
- [ ] GUI framed as adapter/honesty issue, not root cause

### S02 handoff + rejects

- [ ] **Red-capable test sketch** for greenfield bootstrap **and** legacy dogfood (S02 implementer can turn into tests)
- [ ] All **SCOPE-TODOS investigation rejects** addressed (Phase 35 pick, GUI-only, global weaken, delete history)

### Out of scope guardrails (DESIGN-LOCKS)

- [ ] No hosted SaaS / delete feet-seller history / auto-LLM backlog
- [ ] No weakening PLAN gate for active work

## Exit criteria

- [ ] PASS or FAIL with spawn notes on board row **P36-S00-02**
- [ ] On PASS: confidence medium+ with evidence; next **P36-S01-00**
- [ ] On FAIL: spawn rows inserted immediately below with full prompts per agent-loop-protocol

## Minimal todos

- [ ] Read INVESTIGATION.md
- [ ] Re-verify CLI repro + export counts
- [ ] Walk checklist; record findings by severity
- [ ] PASS/FAIL + board Notes

## Next

`P36-S01-00` (on PASS)
