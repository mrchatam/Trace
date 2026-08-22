# P22-S09-04 — Review: bundled default agent catalog + install

## Metadata
- id: P22-S09-04
- todo_ids: [P22-S09-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Independent review of S09-03. Confirm bundled catalog, install path, `harness:subagent` honesty, and E04 extension docs. **No network fetch** in P22.

## Session start

**Fresh subagent** (not S09-03). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [03-bundled-defaults-install.md](03-bundled-defaults-install.md)
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-27, D-22-30
- S09-01 keepers: `TestHarnessAgentCatalogMigrate027`, routing tests

## Review checklist

### E03 — bundled catalog

- [ ] `trace/agents/default.json` committed; `schema_version: 1`
- [ ] All **6** minimum profiles present with correct slugs (`agent:*`)
- [ ] `TestDefaultCatalogValidJSON`, `TestBundledProfilesIncludeRequirements` PASS
- [ ] Requirements reference valid slug prefixes (`skill:`, `mcp:`, `hook:`)

### E04 — extensible registry hook

- [ ] `registry_source`, `registry_version`, `external_url` fields in JSON + DB
- [ ] `trace/agents/README.md` documents future hosted catalog; explicitly **no network fetch in P22**
- [ ] No `http.Get`, no registry sync daemon

### Install path

- [ ] `trace install agents` works on fresh init fixture
- [ ] `TestInstallAgentsSeedsDefaults`, `TestInstallAgentsIdempotent` PASS
- [ ] `trace init --with-agent-defaults` calls same loader
- [ ] Help mentions `install agents`

### harness:subagent (D-22-29 prep)

- [ ] `hook:harness:subagent` upserted on install
- [ ] Status AVAILABLE only with heuristic or env; else UNKNOWN (not fake AVAILABLE)
- [ ] No subprocess / Task tool invocation

### Hard boundaries

- [ ] No agent runner, no spawn
- [ ] Schema still **27**; compat **27**
- [ ] MCP catalog still **14** (S09-07 adds tool)
- [ ] S09-01 routing tests still green with seeded catalog

## Spawn policy

Unmet E03/E04 install path → spawn **`P22-S09-04a`** implement + **`P22-S09-04b`** review below this row.

## Re-run commands

```bash
go test ./internal/install/... ./internal/agents/... ./cmd/trace/... -count=1 -run 'TestInstallAgents|TestDefaultCatalog|TestSubagentHook'
go test ./internal/agents/... -count=1 -run 'TestRecommend'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
rg 'http\.(Get|Post)|Task\(' internal/install/agents.go internal/agents/builtin.go cmd/trace/install.go || true
```

## Exit criteria

- [ ] default.json committed and valid
- [ ] Install idempotent; no agent runner
- [ ] E03 bundled + E04 extension hook closed at install layer
- [ ] Confidence **high** in Notes
