# INTAKE — Phase 36 residuals

**Human-promoted 2026-08-22.** Phase 36 closed at `P36-S03-02` with successor **no successor** and documented non-blocking residuals.

## Symptom / trigger

Phase 36 shipped the fundamental planning-model fix (MCP `trace_plan`, bootstrap, install contract, terminal gate honesty). Several items were **deferred** or left as **low/nit** findings. Human requests a dedicated phase to address them.

## Residual inventory (from P36 DR-HANDOFF + VERIFY-NOTES + PLAN.md deferrals)

| ID | Item | P36 disposition | P37 default stance |
|----|------|-----------------|-------------------|
| R1 | **PlanExists bridge** — advisory when plan-changes ≥ N and `!PlanExists` (not silent satisfy) | PLAN §2.4 defer | S00 triage; likely accept advisory-only |
| R2 | **HTTP POST plan routes** — GUI/API parity for plan bootstrap/create | PLAN touch-list defer | S00 triage; Law 19 adapter |
| R3 | **MCP `trace_loop action=gate`** — gate check without shell | defer | S00 triage |
| R4 | **Bootstrap help** — human-refinement note per PLAN §2.2 | S02-02 low | Likely accept (docs/help) |
| R5 | **`loop status advisories[]`** — goal-structure warning (N>15, no plan) wired to status JSON | S02-02 low / §2.7 partial | Likely accept |
| R6 | **`WarnIfTraceDirWithoutConfig` unit test** | nit | Likely accept |
| R7 | **Enforce default `warn`** when `.trace/` exists without config (§2.6 was doc-only) | defer | S00 triage — product decision |
| R8 | **Goal/plan surface UX** beyond TaskDetail (Overview, plan screen) | defer | S00 triage; may split |
| R9 | **Feet-seller planner quality** — post-bootstrap minimal plan; critique/deep refinement path | verify note | S03 dogfood; optional S02 tooling |
| R10 | **Live GUI browser verify** — Block 4 deferred in P36 VERIFY | verify note | S03 accept if cheap |
| R11 | **CLI greenfield critique path** — Block 1 partial (`plan_uncritiqued` after plan chain) | verify note | S00 triage — agent workflow doc or MCP critique seed |

## Not in scope

- Reopening Phase 36 core deliverables (MCP plan, bootstrap, terminal advisory)
- Hosted SaaS
- Silent PlanExists bridge (fake progressive plan)
- Rewriting feet-seller task history

## Desired outcome

Residuals triaged → implemented or explicitly re-deferred with rationale. Agent/harness/GUI paths more complete than post-P36 baseline. VERIFY against feet-seller + greenfield.
