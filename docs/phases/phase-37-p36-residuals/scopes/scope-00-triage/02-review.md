# P37-S00-02 — Review triage

## Metadata
- id: P37-S00-02
- todo_ids: [P37-S00-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Independent review of `RESIDUALS.md` vs INTAKE R1–R11, DESIGN-LOCKS, planner-locked defaults (S00-00), and **live repo**. Every residual must have accept/defer/reject with effort and risk. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — live facts + preliminary triage
- [01-triage.md](01-triage.md) — locked defaults for comparison
- `RESIDUALS.md` (S00-01 deliverable)
- Phase 36 [VERIFY-NOTES.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md) § Residuals

## Session start

Follow agent-loop-protocol Session start. **Fresh review session** — do not assume S00-01 claims; spot-check **≥5** live-code citations (mandatory IDs: R1, R5, R2 or R3, R6, one defer).

## Locked expectations (S00-00 — challenge only with evidence)

| ID | Expected decision | Key check |
|----|-------------------|-----------|
| R1 | Accept advisory-only | §3 signal: plan_changes ≥ 1, !PlanExists, `advisories[]`, no store mutation |
| R2 | Accept | Law 19 handlers in `internal/httpapi/` → `planner.Service`; OpenAPI touch |
| R3 | Accept | `tools_loop.go` new `gate` case → `loop.EvaluateGate` like `cmd/trace/loop.go:122–176` |
| R4 | Accept (S) | Help text gap in `cmd/trace/plan.go` (~lines 55–72, 334–360) |
| R5 | Accept (S) | Partial shipped: `advisory.go` exists; `StatusResult` still no `advisories` (`apply.go:196–209`) |
| R6 | Accept (S) | `enforce.go:43–57` untested |
| R7 | Defer | Default off preserved; re-defer owner + trigger in §6 |
| R8 | Accept minimal | Overview beyond TaskDetail; Law 19 — API consumer only |
| R9 | Defer → S03 | No feet-seller history rewrite |
| R10 | Accept → S03 | Browser verify, not S02 |
| R11 | Accept doc | Not new MCP critique tool; loop apply path documented |

## Checklist vs DESIGN-LOCKS + INTAKE

### Coverage
- [ ] All R1–R11 addressed — no orphan INTAKE rows
- [ ] Each row has **decision + effort (S/M/L) + risk (low/med/high)**
- [ ] Summary counts match table rows

### Phase 36 guarantees (must preserve in rationale)
- [ ] MCP `trace_plan` + bootstrap unchanged as primary recovery path
- [ ] Terminal gate honesty (`goal_plan_gap_terminal_advisory`) not weakened
- [ ] Active non-terminal `plan_missing` still blocks edit (`internal/deliberation/select.go:28–29`)
- [ ] Partial P36 S02 shipment acknowledged (R4/R5/R6/R7/R8 rows)

### R1 — PlanExists bridge
- [ ] **Reject** silent PlanExists satisfy explicitly stated
- [ ] Accept path uses **advisory only** — no `policy.go` bridge, no bootstrap side effects in status
- [ ] R1 signal distinct from R5 (plan_changes vs task-count threshold)
- [ ] Threshold N documented (expected: **1**)

### R2 / R3 — Law 19 adapters
- [ ] R2: HTTP handlers thin — cite target file(s); no business logic fork in `web/`
- [ ] R2: At minimum POST bootstrap (or S01-expanded set) — not unbounded CRUD
- [ ] R3: MCP gate mirrors CLI JSON envelope + exit semantics
- [ ] No duplicate gate evaluation logic outside `internal/loop`

### R5 / R8 — advisories dependency
- [ ] R5 accept includes `advisories[]` on `trace.loop.status.v1`
- [ ] S02 wave order places R5 (and R1) before R8 GUI consumption
- [ ] R8 GUI reads API/status — no planner imports in `web/`

### R7 — product decision
- [ ] Defer (not silent accept of default flip)
- [ ] §6 owner + trigger for re-defer
- [ ] Acknowledges P36 §2.6 nudge already shipped

### R11 — critique path
- [ ] Doc/workflow path preferred over new MCP tool
- [ ] References existing greenfield proof (P36 Block 0 MCP test)
- [ ] CLI verify script gap documented as doc/fixture issue, not product fork

### S02 scope bounds
- [ ] Accept set implementable in one S02 wave (Waves A–D) or explicit split with rationale
- [ ] R9/R10 not incorrectly assigned to S02 implementation
- [ ] Re-defer registry §6 complete for R7, R8-full, R9

### Live-code spot-checks (record in board Notes)

Verify at least these against repo:

| Claim | Verify command / read |
|-------|----------------------|
| No `advisories` on StatusResult | `internal/loop/apply.go` struct fields |
| GoalStructureWarning exists | `internal/planner/advisory.go:13–41` + `advisory_test.go` |
| trace_loop actions | `internal/mcp/tools_loop.go:39–47` — no `gate` |
| HTTP plans read-only | `internal/httpapi/server.go:282` — `GET /v1/plans` only |
| Bootstrap help gap | `cmd/trace/plan.go` help vs PLAN §2.2 risk text |
| WarnIfTraceDirWithoutConfig | `internal/config/enforce.go:43–57`; grep `enforce_test.go` |

## Findings severity

Record in board Notes:

| Severity | Action |
|----------|--------|
| blocker | Fix RESIDUALS.md inline if trivial; else spawn `P37-S00-02a` implement + `02b` review immediately below this row |
| high | Same as blocker |
| medium | Prefer spawn unless one-line RESIDUALS fix |
| low / nit | Note only; may fix inline |

## Exit criteria

- [ ] No open blocker/high without spawn or inline RESIDUALS fix
- [ ] All checklist items above addressed (pass or documented exception in Notes)
- [ ] Spot-check results recorded in board Notes (pass/fail per cite)
- [ ] Confidence **medium+** with evidence; list residual risks if medium
- [ ] Board row `P37-S00-02` → `done`; on PASS next **P37-S01-00**

## Minimal todos

- [ ] Read RESIDUALS.md end-to-end
- [ ] Run checklist vs DESIGN-LOCKS + INTAKE + locked defaults
- [ ] Spot-check ≥5 live-code claims; record in Notes
- [ ] Fix RESIDUALS.md inline if trivial; else spawn below this row
- [ ] Update board status + notes only

## Next

`P37-S01-00` (on PASS)
