# P26-S04-02 — Installer fix review

## Metadata
- id: P26-S04-02
- todo_ids: [P26-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of P26-S04-01.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Checklist

- [ ] `ParentOrchestratorRule` referenced in `cursorRulesMDCContent()`
- [ ] Claude fallback also includes the rule (or Notes justify asymmetry)
- [ ] Generated cursor rules output asserts/includes `"Parent orchestrator"`
- [ ] `go test ./internal/install/...` PASS
- [ ] Package scope stays under `internal/install/`; no daemon/HTTP/schema

## Spawn policy

HIGH → spawn fix+review below. Else close with confidence.

## Exit criteria

- [ ] Notes with confidence + evidence
