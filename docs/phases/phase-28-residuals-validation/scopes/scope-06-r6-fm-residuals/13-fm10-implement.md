# P28-S06-13 — FM-10 / FR-P28-07 implementer

## Metadata
- id: P28-S06-13
- todo_ids: [P28-S06-13]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [user-codegraph, user-trace]
- verification: mixed
- hooks: []

## Objective

**FR-P28-07 / FM-10:** Ensure promotion API is used in live loops (API shipped; build-only exports still risk 0 discoveries). Dogfood + optional apply E2E assert (TEST-MATRIX M-01 already covers apply promote); measure live promotion rate.

## References

- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-10; INT-01
- [TEST-MATRIX.md](../scope-01-integration-tests/TEST-MATRIX.md) — M-01
- `internal/domain/promote.go`, `internal/loop/apply_test.go` promote path
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start.

## Acceptance hint

Directed/build run shows ≥1 discovery linked to task **or** spawned task from BLOCKING; auto-spawn without human gate stays out of scope (FR-P28-D1).

## Preflight

```bash
cd /home/ali/Desktop/Trace
GOPROXY=direct go test ./internal/loop/ -count=1 -run 'Promote|Apply.*Spawn|Apply.*Discovery'
grep -n 'M-01\|spawned_tasks' docs/phases/phase-28-residuals-validation/scopes/scope-01-integration-tests/TEST-MATRIX.md | head
```

## Suggested work

1. Measure live or scripted promotion usage (discovery↔task link or spawned_tasks).
2. Strengthen dogfood/apply assert only if M-01 insufficient for live-loop claim.
3. Document promotion rate / residual risk for build-only 0-discovery exports.
4. `FM10-NOTES.md` with evidence; no D1 auto-spawn.

## Out of scope

- FR-P28-D1; daemon/HTTP; rewriting M-01 history without need

## Exit criteria

- [ ] Acceptance hint met with evidence
- [ ] Next runnable **P28-S06-14**

## Todo updates

Status + notes on **P28-S06-13** only.

## Next

`P28-S06-14`
