# P22-S09-06 — Review: phase/task routing + subagent recommendations

## Metadata
- id: P22-S09-06
- todo_ids: [P22-S09-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Independent review of S09-05. Confirm recommendations only (no spawn), correct routing, fresh-subagent hint for reviews per agent-loop-protocol.

## Session start

**Fresh subagent** (not S09-05). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [05-routing-loop-packet.md](05-routing-loop-packet.md)
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-25, D-22-28, D-22-29
- [README.md](../../README.md) — E01, E02, C09/C28/C39 matrix rows
- `docs/rules/agent-loop-protocol.md` — fresh subagent per review row

## Review checklist

### E01 — subagent recommendation

- [ ] `TestRecommendSubagentWhenAvailable` PASS — `use_subagent: true` + prompt_stub when `harness:subagent` AVAILABLE
- [ ] When hook UNKNOWN/UNAVAILABLE: `use_subagent: false`; recommendation still present (honest inline fallback)
- [ ] Prompt stub mentions fresh session / not implementer
- [ ] **No** Task tool, subprocess, or MCP that spawns agents

### E02 — phase/task routing

- [ ] `TestLoopNextIncludesHarnessRecommendations` PASS
- [ ] CRITIQUE → code-reviewer or nested-reviewer (`TestRecommendAgentForPhaseCritique` or packet test)
- [ ] Perf title/tags → performance-reviewer (`TestRecommendPerformanceReviewerForPerfTask`)
- [ ] INVESTIGATE/ORIENT → explore when seeded
- [ ] Cap **4** entries max

### Checklist strengthening (note only — do not box)

- [ ] **C09**: structured thought process — loop next surfaces reviewer profile + subagent hint for CRITIQUE
- [ ] **C28**: redundant work — routing suggests specialized reviewer instead of generic implementer continuation
- [ ] **C39**: workflow — harness recommendations in normal loop next path

### Packet quality

- [ ] Additive field on `trace.loop.next.v1`; does not break existing consumers
- [ ] `TestLoopNextDeliberationSectionPresent`, `TestLoopApply` keepers green
- [ ] Empty catalog → honest empty (`TestNoRecommendationWhenCatalogEmpty`)
- [ ] `missing_capabilities[]` populated when requirements not AVAILABLE

### Hard boundaries

- [ ] Grep loop/agents — no spawn/runner
- [ ] Schema **27**; MCP catalog still **14**

## Spawn policy

Unmet E01/E02 → spawn **`P22-S09-06a`** implement + **`P22-S09-06b`** review below this row.

## Re-run commands

```bash
go test ./internal/loop/... ./internal/agents/... -count=1 -run 'TestLoopNextIncludesHarness|TestRecommendSubagent|TestRecommendPerformance|TestNoRecommendationWhenCatalogEmpty|TestLoopNextDeliberationSectionPresent|TestLoopApply'
rg 'Task\(|exec\.Command' internal/loop/next.go internal/agents/ || true
```

## Exit criteria

- [ ] E01/E02 closed or spawned
- [ ] C09/C28/C39 strengthening noted in board Notes (with test paths)
- [ ] Confidence **high**
