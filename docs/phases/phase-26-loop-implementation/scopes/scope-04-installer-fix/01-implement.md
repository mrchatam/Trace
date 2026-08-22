# P26-S04-01 — Installer fix implementer

## Metadata
- id: P26-S04-01
- todo_ids: [P26-S04-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: []
- verification: automated
- hooks: []

## Objective

Wire `ParentOrchestratorRule` into installed cursor + Claude enforcement rule bodies; add regression test.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md)
- `internal/install/enforcement.go`
- `internal/install/enforcement_test.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Scope | `internal/install/` only |
| Phrase to assert | `"Parent orchestrator"` |
| Order | Keep gate body → GapPassPrompt → ParentOrchestratorRule (or document alternate order in Notes) |
| Edit boundary | Change only `internal/install/enforcement.go` (tests in `internal/install/enforcement_test.go`) |

## Change

In `cursorRulesMDCContent()` and `claudeFallbackRulesContent()`, append `ParentOrchestratorRule` after `GapPassPrompt` (same pattern as `AgentsEnforcementBlock` already uses for GapPassPrompt).

## Test

In `enforcement_test.go`:

```go
if !strings.Contains(cursorOut, "Parent orchestrator") {
    t.Error("cursor rules missing Parent orchestrator rule (INT-04 / P25-2)")
}
```

## Exit criteria

- [ ] `rg` shows `ParentOrchestratorRule` used inside both `cursorRulesMDCContent` and `claudeFallbackRulesContent`
- [ ] `go test ./internal/install/...` PASS
- [ ] No other packages changed

## Minimal todos

- [ ] Wire cursor + Claude content
- [ ] Add/adjust test
- [ ] Own row `done`
