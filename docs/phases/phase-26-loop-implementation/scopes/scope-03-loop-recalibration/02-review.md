# P26-S03-02 — P25-B review

## Metadata
- id: P26-S03-02
- todo_ids: [P26-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of P26-S03-01 against PLAN.md **S03-T01–T07** and locked defaults in `01-implement.md`. Spawn forward on HIGH findings. Do not rewrite this review prompt after `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md)
- [00-PLANNER.md](00-PLANNER.md)
- [PLAN.md](../scope-01-planning/PLAN.md) — S03
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) — INT-02 / INT-05 / INT-09

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Fresh context — do not share the implementer session.

## Review checklist

### INT-02 / D1 / T01–T02 / T07

- [ ] Named const `SaturationEmptyThreshold = 2` present; not magic bare `2` scattered without name
- [ ] First pure-empty apply does **not** set sticky `Stopped=true`
- [ ] Second consecutive pure-empty saturates with `stop_reason` / reason `p19_saturated`
- [ ] Discoveries-only apply (discoveries imported, zero plan/spawn) does **not** increment consecutive-empty and does not saturate
- [ ] Non-empty apply (plan_changes or spawned_tasks) zeroes the counter
- [ ] `MaxIterationsReached` still immediate-saturates
- [ ] Column `consecutive_empty_applies` + migration **028**; store/domain wired; no unversioned schema
- [ ] `HopBudget` still **12**
- [ ] Replay / `p19SaturatedFromLastStep` / status / next / gate no longer use single last-step zero-write as sole rule

### INT-05 / D2 / T03–T04

- [ ] Domain reset clears `Stopped`, `StopReason`, `HopCount`, `consecutive_empty_applies`
- [ ] Persisted `CurrentPhase` set to EXECUTE; `PlanCritiqued` preserved
- [ ] CLI `trace loop reset --task <id>` wired + help; no MCP reset
- [ ] After reset: no immediate re-STOP from prior saturation / counter alone
- [ ] Documented: edit gate may still block when `plan_critiqued=false` (SelectNext not EXECUTE) — not a defect if saturation sticky-STOP is gone

### INT-09 / D3 / T05–T06

- [ ] When `Stopped`, `SelectNext` prefers **persisted** `StopReason` (not blanket `hop_budget_exceeded`)
- [ ] `hop_budget_exceeded` only when hop ≥ budget (or empty reason fallback)
- [ ] Status JSON includes `stop_reason` alongside `why_selected`
- [ ] After saturation STOP: gate `reason_code` == export `stop_reason` (`p19_saturated`)
- [ ] Next packet no longer contradicts export without the persisted reason winning

### Laws / blast radius

- [ ] No daemon/HTTP; no MCP reset tool
- [ ] S02 promotion surfaces untouched (`spawned_task_ids`, `discovery_id`, promote helper, GapPassPrompt, MCP description)
- [ ] `ParentOrchestratorRule` still **not** wired (S04)
- [ ] Package blast radius within AUDIT/PLAN paths (`internal/deliberation`, `internal/loop`, `internal/domain`, `internal/store`, `cmd/trace`)
- [ ] Embed expectations updated to **28** (`deliberation_test` + `TestMigrateBackupAuthCLI`)

### Tests

- [ ] D4: first empty / second empty / discoveries-only / reset / reason alignment covered
- [ ] `go test ./internal/...` PASS
- [ ] `go test ./cmd/trace/...` PASS if migrate/CLI tests updated

## Spawn policy

HIGH → insert `P26-S03-02a` implement + `P26-S03-02b` re-review **immediately below** this row, with full protocol prompts.  
No HIGH → close with confidence (**medium** only if residuals listed).

## Exit criteria

- [ ] No open HIGH without pending spawn
- [ ] Confidence medium+ with evidence in Notes (files + tests)
- [ ] Own row `done` / `failed` / spawn recorded
