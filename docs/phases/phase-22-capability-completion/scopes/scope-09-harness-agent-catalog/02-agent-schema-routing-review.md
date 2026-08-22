# P22-S09-02 — Review: harness agent schema + routing library

## Metadata
- id: P22-S09-02
- todo_ids: [P22-S09-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Independent review of S09-01. Confirm **E03** foundation (schema, CRUD, seed export) and **E02** routing core (deterministic, no harness execution).

## Session start

**Fresh subagent** (not S09-01). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-agent-schema-routing.md](01-agent-schema-routing.md)
- [README.md](../../README.md) — E01–E04 enhancement matrix
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-26, D-22-28
- P21 keepers: `TestSeedExportOmitsDeniedSurfaces`, `TestUpsertCapabilityGetAndReject`

## Review checklist

### E03 — schema + catalog foundation

- [ ] Exactly **one** new migration: `027_harness_agents.sql`; **no** `028+`
- [ ] `ls internal/store/schema/*.sql | wc -l` → **27**
- [ ] `TestHarnessAgentCatalogMigrate027` PASS; idempotent re-open
- [ ] Tables match locked DDL (slug unique, requirements UNIQUE agent+slug)
- [ ] Nullable `external_url`; `registry_source` enum values documented
- [ ] Seed export/import round-trip for `harness_agents` (additive keys only)

### E02 — routing library (partial until loop wired S09-05)

- [ ] `internal/agents.RecommendAgents` exists; deterministic (`TestRoutingDeterministic`)
- [ ] `TestRecommendAgentForPhaseCritique` PASS
- [ ] `TestRecommendPerformanceReviewerForPerfTask` PASS
- [ ] Empty catalog → empty recommendations (no panic)
- [ ] Result cap **4** enforced or parameterized

### AGENT kind

- [ ] `CapabilityKindAgent` / `AGENT` in store + domain
- [ ] `TestCapabilityKindAgent` PASS; `TestUpsertCapabilityGetAndReject` updated (AGENT accepted)
- [ ] Unknown kinds still fail closed

### Hard boundaries (D-22-25)

- [ ] Grep new code — **no** `Task(`, subprocess agent runner, HTTP fetch for registry
- [ ] **No** CLI `trace agents`, **no** MCP `trace_agents`, **no** loop packet changes this row
- [ ] **No** bundled `default.json` yet — S09-03

### Compat + landmines

- [ ] `CGO_ENABLED=1 go test ./evals/compat/... -run TestCompatibilitySecurityChecklist` — ceiling **27**
- [ ] P21 seed keepers still green
- [ ] MCP catalog still **14** (unchanged)

## Spawn policy

If E03 foundation or E02 routing core unmet: spawn **`P22-S09-02a`** implement + **`P22-S09-02b`** review immediately below. **Do not** mark `done` with “later”.

## Re-run commands

```bash
go test ./internal/agents/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestHarnessAgent|TestCapabilityKindAgent|TestRecommend|TestRoutingDeterministic'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 27
rg -l 'Task\(|exec\.Command.*agent|http\.Get.*registry' internal/agents/ internal/store/harness_agents.go internal/domain/harness_agents.go || true
go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces'
```

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with re-run output in Notes
- [ ] E03 foundation closed; E02 routing library closed at library layer (loop wiring deferred to S09-05)
- [ ] Board Notes: findings + confidence + test citations
