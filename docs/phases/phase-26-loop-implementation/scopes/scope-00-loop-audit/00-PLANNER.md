# P26-S00-00 — Scope planner (loop audit)

## Metadata
- id: P26-S00-00
- todo_ids: [P26-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Finalize the S00 audit implementer prompt against live `internal/` layout and gate **P26-S00-01** so `AUDIT.md` maps INT-01/02/05/06/09 plus the P25-2 installer gap. No product code in this scope.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 26 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [Phase 25 DR-HANDOFF](../../../phase-25-orchestrator-gap-pass/DR-HANDOFF.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-loop-audit/AUDIT.md` |
| Product Go | **No** |
| Threshold numbers | Document options only; do **not** pick final values (S01/S02–S03) |
| Installer gap | Confirm `ParentOrchestratorRule` defined but unused in `cursorRulesMDCContent` / `claudeFallbackRulesContent` |
| Sequence | S00 → S01 → S02 → S03 → S04 → S05 (serial default) |

## Files to audit (at minimum)

| Package / path | What to find |
|----------------|--------------|
| `internal/loop/` | SelectNext, hop_count, saturation, Stopped, apply / spawned tasks |
| `internal/store/` | deliberation state, hop_count, reset surfaces |
| `internal/seed/` | discovery export / task helpers |
| `cmd/trace/` loop cmds | gate, status, apply CLI |
| `internal/mcp/` | `trace_add` description ordering |
| `internal/install/enforcement.go` | `ParentOrchestratorRule` vs `cursorRulesMDCContent` |

## Planner gate

- [ ] `01-loop-audit.md` runnable (exit criteria + no-product-code rule)
- [ ] `SCOPE-TODOS.md` lists S00 board rows
- [ ] Live paths above still exist (adjust `01-loop-audit.md` if renamed)

## Exit criteria

- [ ] Audit implementer prompt locked enough for a fresh subagent
- [ ] Board row P26-S00-00 Notes cite what was verified/thickened
- [ ] Next runnable remains **P26-S00-01** (do not start S01)

## Todo updates

Status + notes on **P26-S00-00** only.
