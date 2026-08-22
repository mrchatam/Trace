# P22-S09-05 — Implement: phase/task routing + subagent recommendations in loop packet

## Metadata
- id: P22-S09-05
- todo_ids: [P22-S09-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Wire **harness agent recommendations** into agent-facing packets. Trace suggests the right profile for the situation (e.g. performance-reviewer for perf tasks) and recommends a **fresh subagent** when independent review is appropriate and `harness:subagent` is available. **E01, E02** — strengthens **C09, C28, C39**. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-agent-schema-routing.md](01-agent-schema-routing.md) — `RecommendAgents` API
- [03-bundled-defaults-install.md](03-bundled-defaults-install.md) — bundled catalog (must be **done**)
- Live: `internal/loop/next.go` (`NextPacket`), `internal/loop/next_test.go`, `docs/rules/agent-loop-protocol.md` (fresh subagent rule)
- Pattern: `WorkConflictsSection` — additive packet section with freshness

## Live baseline (after S09-03)

| Present | Absent |
|---------|--------|
| `RecommendAgents` + seeded bundled catalog | `harness_recommendations[]` in loop next |
| `hook:harness:subagent` on install | Subagent hint in packet |
| Loop next v1 with deliberation, work_conflicts, … | `trace context --agents` |

## Locked defaults

| Item | Value |
|------|-------|
| Packet section | `harness_recommendations` on `NextPacket` — cap **4** entries |
| Schema | Stay `trace.loop.next.v1` — additive JSON field only (no version bump unless compiler requires) |
| Trust | Recommendations are **UNTRUSTED_DATA** / suggestions; document in section or compiler markdown |
| Subagent gate | `use_subagent: true` only when matched agent `recommend_subagent` **and** `harness:subagent` status **AVAILABLE** |
| Missing hook | When UNAVAILABLE/UNKNOWN: still recommend profile; `use_subagent: false`; optional `prompt_stub` suggests inline review |
| Trace role | **Never invoke** subagents — recommendation text only (D-22-25) |

## `harness_recommendations[]` entry shape

```json
{
  "agent_slug": "agent:code-reviewer",
  "subagent_type": "code-reviewer",
  "reason": "phase CRITIQUE matches profile deliberation_phases",
  "confidence": "high",
  "use_subagent": true,
  "prompt_stub": "Fresh subagent for independent review — not the implementer session.",
  "missing_capabilities": ["skill:code-review-and-quality"]
}
```

## Section wrapper (mirror other sections)

```go
type HarnessRecommendationsSection struct {
    Freshness string                      `json:"freshness"`
    Items     []agents.Recommendation     `json:"items"`
}
```

- `Freshness`: `fresh` when catalog non-empty and caps loaded; `unknown` when catalog empty; never claim `fresh` with zero agents.

## Wiring requirements

1. **`BuildNextPacket`**: after deliberation phase resolved, gather:
   - current phase string from deliberation state
   - task title, tags (if stored), goal keywords from plan snapshot
   - harness capabilities map: load all capabilities with kind HOOK/SKILL/MCP where slug matches agent requirements or `harness:subagent`
2. Call `agents.RecommendAgents`; attach section to `NextPacket`.
3. Deterministic rules (re-use library from S09-01 — **do not duplicate** routing logic in loop package):
   - **CRITIQUE** or post-EXECUTE review context → code-reviewer / nested-reviewer
   - perf keywords → performance-reviewer
   - security keywords → security-reviewer
   - **INVESTIGATE** / **ORIENT** → explore
4. **`trace context --agents`** (optional): slim list of top 2 recommendations when flag set — mirror CLI pattern from `--evaluations`.
5. Compiler markdown render (optional): bullet list under `## Harness recommendations` in human packet.
6. Empty catalog → `items: []`, freshness `unknown` — `TestNoRecommendationWhenCatalogEmpty`.

## Touch files

- `internal/loop/next.go` — `HarnessRecommendationsSection`, wire in `BuildNextPacket`
- `internal/loop/next_test.go` — named tests
- `internal/agents/routing.go` — extend only if loop needs helpers (prefer not)
- `internal/compiler/packet.go` — optional markdown section
- `cmd/trace/context.go` — optional `--agents` flag

## Named tests

| Test | Proves |
|------|--------|
| `TestLoopNextIncludesHarnessRecommendations` | section present in JSON |
| `TestRecommendSubagentWhenAvailable` | `use_subagent: true` when hook AVAILABLE |
| `TestRecommendSubagentHonestWhenUnavailable` | `use_subagent: false` when hook UNKNOWN |
| `TestRecommendPerformanceReviewerForPerfTask` | keyword routing in packet |
| `TestNoRecommendationWhenCatalogEmpty` | honest empty |
| `TestLoopNextDeliberationSectionPresent` | keeper — must stay green |

```bash
go test ./internal/loop/... ./internal/agents/... ./cmd/trace/... -count=1 -run 'TestLoopNextIncludesHarness|TestRecommendSubagent|TestRecommendPerformance|TestNoRecommendationWhenCatalogEmpty|TestLoopNextDeliberationSectionPresent'
```

## Exit criteria

- [ ] **E01** fresh subagent recommendation when `harness:subagent` AVAILABLE
- [ ] **E02** phase + keyword routing visible in loop next
- [ ] Notes cite which **C09 / C28 / C39** bullets strengthened (do not checkbox checklist — review/VERIFY owns)
- [ ] Named tests PASS; P21 loop keepers green
- [ ] Board Notes

## Minimal todos

- [ ] `harness_recommendations` section on NextPacket
- [ ] Wire RecommendAgents in BuildNextPacket
- [ ] Subagent gate + prompt_stub
- [ ] Tests + optional context --agents
- [ ] Board notes

## Residual risks (carry to S09-06)

- Task tags not populated on all fixtures — keyword routing may only hit title
- Duplicate recommendations if routing returns overlapping profiles — dedupe by slug
- Markdown compiler section omitted — JSON path is required; markdown optional
- C28/C39 not fully closed until CLI/MCP agents surface (S09-07) + VERIFY
