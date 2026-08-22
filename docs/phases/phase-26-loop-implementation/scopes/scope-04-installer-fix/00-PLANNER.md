# P26-S04-00 — Scope planner (installer P25-2 fix)

## Metadata
- id: P26-S04-00
- todo_ids: [P26-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Confirm live P25-2 gap and lock implement/review prompts. Small, low-risk. No product code on this row beyond verification reads.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 25 DR-HANDOFF](../../../phase-25-orchestrator-gap-pass/DR-HANDOFF.md)
- `internal/install/enforcement.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| File | `internal/install/enforcement.go` |
| Change | Reference `ParentOrchestratorRule` inside `cursorRulesMDCContent` and `claudeFallbackRulesContent` |
| Test | Assert `"Parent orchestrator"` in cursor rules output |
| Packages | `internal/install/` only |
| Board order | After S03 by default (may parallel S03 only after S01 — prefer serial) |

## Problem

`ParentOrchestratorRule` is defined and non-empty but **not** concatenated into `cursorRulesMDCContent()` / Claude fallback — E02 P25-2 FAIL.

## Verify gap before closing planner

```bash
rg "ParentOrchestratorRule" internal/install/enforcement.go
# Expect: const definition; missing usage inside cursorRulesMDCContent / claudeFallbackRulesContent
```

## Planner gate

- [ ] Gap still present (or Notes explain already fixed + skip implement with evidence)
- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] `SCOPE-TODOS.md` current

## Exit criteria

- [ ] Implementer prompt locked; own Notes cite `rg` evidence
