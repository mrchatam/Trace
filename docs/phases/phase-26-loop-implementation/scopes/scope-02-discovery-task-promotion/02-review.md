# P26-S02-02 — P25-A review

## Metadata
- id: P26-S02-02
- todo_ids: [P26-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of P26-S02-01 against PLAN.md **S02-T01–T07** and locked defaults in `01-implement.md`. Spawn forward on HIGH findings. Do not rewrite this review prompt after `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md)
- [00-PLANNER.md](00-PLANNER.md)
- [PLAN.md](../scope-01-planning/PLAN.md) — S02
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) — INT-01 / INT-06

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Fresh context — do not share the implementer session.

## Review checklist

### INT-01 / D1–D2 / T01–T03 / T05–T07

- [ ] Domain helper exists; BLOCKING-only; fail-closed on INFO/PLAN_AFFECTING
- [ ] Idempotent: re-promote same discovery → same task ID, no duplicate rows; `ImportSeedTask` work_state preserved
- [ ] Apply `spawned_tasks[].discovery_id` promotes; `ApplyResult.spawned_task_ids` populated
- [ ] CLI `trace loop apply` stdout includes `spawned_task_ids` (JSON encode of result)
- [ ] `trace add task --from-discovery` uses the same helper; help mentions the path
- [ ] Discoveries-only apply / `trace add discovery` / seed import **do not** create tasks
- [ ] `loop next` includes `promotion_candidates` (BLOCKING, unlinked); empty `[]` when none; no spawn
- [ ] Seed import then promote: single task; existing seed task UUIDs unchanged

### INT-06 / D3–D4 / T04

- [ ] `trace_add` Description: discovery first + promotion sentence (task or loop apply)
- [ ] MCP discovery handler still does **not** auto-spawn
- [ ] `GapPassPrompt` has one extra promotion sentence; P25-C gap-pass body not gutted
- [ ] `ParentOrchestratorRule` still **not** wired (S04)

### Human gate + laws

- [ ] No silent background spawn, daemon, HTTP, or new MCP promote tool
- [ ] No `promote_discoveries[]`; apply schema still `trace.loop.apply.v1`
- [ ] No SQLite migration without version bump
- [ ] Saturation / reset / STOP-reason files not “fixed” here (S03)
- [ ] `buildPromotionBlocked` not used as the discovery→task gate
- [ ] Changes stay in AUDIT/PLAN packages (`internal/domain`, `internal/loop`, `cmd/trace`, `internal/mcp`, `internal/install`)

### Tests

- [ ] D5: BLOCKING discovery → linked task (apply and/or `--from-discovery`)
- [ ] Negative: discoveries-only apply does not grow roster
- [ ] MCP description string asserted
- [ ] `go test ./internal/...` PASS (plus `cmd/trace` if CLI tests added)

## Spawn policy

HIGH → insert `P26-S02-02a` implement + `P26-S02-02b` re-review **immediately below** this row, with full protocol prompts.  
No HIGH → close with confidence (**medium** only if residuals listed).

## Exit criteria

- [ ] No open HIGH without pending spawn
- [ ] Confidence medium+ with evidence in Notes (files + tests)
- [ ] Own row `done` / `failed` / spawn recorded
