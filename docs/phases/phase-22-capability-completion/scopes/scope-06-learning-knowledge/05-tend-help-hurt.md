# P22-S06-05 — Implement: tend help/hurt + successful approaches

## Metadata
- id: P22-S06-05
- todo_ids: [P22-S06-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Help agents see **what tends to improve or damage** this project (**C22**), **surface successful approaches** (**C23**), and **decide from accumulated evidence** (**C24**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- S05: `BuildEvidenceSections`, `ListWorkedApproaches`, `planning_evidence` cap **8**
- S06-01: `ListTendencies` from `change_patterns` threshold **2**
- S06-03: `engineering_knowledge` active rows

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| S05 packet: evaluations, reflections, planning_evidence | `tendencies`, `successful_approaches`, `similar_changes` |
| S05 `trace outcomes worked` (C33 query) | C23 context surfacing (packet/loop) |
| Pattern + knowledge tables (S06-01/03) | compiler/loop aggregation for tendencies |
| Compat **25** | **026+** |

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **None**; compat stays **25** |
| Tendencies | **`ListTendencies(ctx)`** — from `change_patterns`: `count_positive ≥ 2` → `direction=improve`; `count_negative ≥ 2` → `direction=damage`; fields: `change_kind`, `outcome_kind`, `direction`, `count_positive`, `count_negative`, `last_seen` |
| Successful approaches | **`ListSuccessfulApproaches(ctx, opts)`** — merge: (a) S05 **`ListWorkedApproaches`** (improvements + test_pass), (b) active **`engineering_knowledge`** where `topic ∈ {improvement, pattern}` and pattern direction improve; dedupe by id; newest first; cap **8** in packet |
| C23 vs C33 | S05 CLI **`trace outcomes worked`** unchanged; S06 adds **context + loop** surfacing — extend, do not remove S05 query |
| Similar (C20 surface) | Optional loop section **`similar_changes`** — when task has latest change paths, call **`QuerySimilarChanges`** with dominant prefix/kind; cap **8** |
| Packet | Add **`tendencies[]`**, **`successful_approaches[]`** on `compiler.Packet` (additive; schema_version stays **`"0.2"`**) |
| Loop next | Add sections **`tendencies`**, **`successful_approaches`**, optional **`similar_changes`** with freshness; **`trace.loop.next.v1`** unchanged |
| CLI | `trace knowledge tendencies [--limit N]` JSON (default 32, cap 64) |
| Markdown | Short labeled sections when `IncludeMarkdown` |
| MCP | **`trace_context`** inherits via compiler — no new tools; catalog **13** |
| C24 | Loop next includes tendencies + successful_approaches + existing planning_evidence — **`TestLoopNextIncludesEvidenceForDecisions`** |
| Checklist | C22, C23, C24 **unboxed** until S06-06 |

## Requirements

1. **`ListTendencies`** in domain (may live in `patterns.go` or `knowledge.go`).
2. **`ListSuccessfulApproaches`** — domain merge helper; reuse `ListWorkedApproaches` internally.
3. Extend **`internal/compiler/evidence_sections.go`** (or sibling) — attach tendencies + successful_approaches; keep cap **8**.
4. Extend **`internal/loop/next.go`** — mirror sections + optional `similar_changes`.
5. CLI `trace knowledge tendencies`.
6. Named tests in compiler, loop, domain, cmd/trace.

## Touch files

- `internal/domain/patterns.go` or `knowledge.go` (extend)
- `internal/compiler/evidence_sections.go`, `packet.go`, `compiler.go`
- `internal/loop/next.go`, `next_test.go`
- `cmd/trace/knowledge.go`, `knowledge_test.go` (extend)
- `docs/CAPABILITIES_CHECKLIST.md` — note only until S06-06 review boxes

## Named tests

| Test | Proves |
|------|--------|
| `TestTendHelpHurtInContext` | C22 — TaskContext JSON includes tendencies when patterns seeded |
| `TestSuccessfulApproachesSurfaced` | C23 — packet includes worked + knowledge rows |
| `TestLoopNextIncludesEvidenceForDecisions` | C24 — loop next has tendencies + successful_approaches + planning_evidence |
| `TestContextIncludesEvaluationsAndReflections` | keeper (S05) |
| `TestLoopNextPlanningEvidenceSection` | keeper (S05) |
| `TestOutcomesWorkedCLI` | keeper (S05 C33 CLI) |

```bash
go test ./internal/compiler/... ./internal/loop/... ./internal/domain/... -count=1 -run 'TestTendHelpHurtInContext|TestSuccessfulApproachesSurfaced|TestLoopNextIncludesEvidenceForDecisions|TestContextIncludesEvaluations|TestLoopNextPlanningEvidence'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestContext|TestLoopNext|TestKnowledgeTendencies|TestOutcomesWorked'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 25
```

## Exit criteria

- [ ] C22, C23, C24 true (named tests)
- [ ] Compat **25** unchanged; MCP **13**
- [ ] S05 evidence sections unchanged behavior
- [ ] Checklist caps **unboxed** until S06-06
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Tendency + successful-approaches domain helpers
- [ ] Compiler + loop sections
- [ ] CLI tendencies
- [ ] Tests
- [ ] Board status + notes

## Residual risks (carry to S06-06)

- **C23 overlap with C33** — reviewer confirms packet surfacing distinct from CLI query
- **Cap-8 overflow** — no named overflow test required (planner lock) but ordering must be newest-first
- **similar_changes** optional — absence OK when task has no changes; must not error loop next
