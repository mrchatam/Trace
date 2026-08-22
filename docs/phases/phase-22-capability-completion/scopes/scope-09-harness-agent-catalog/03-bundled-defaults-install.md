# P22-S09-03 — Implement: bundled default agent catalog + install

## Metadata
- id: P22-S09-03
- todo_ids: [P22-S09-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Ship the **default harness agent catalog** installed with Trace. Maps known agent profiles to skills, MCPs, hooks, and tools Trace already models. **E03** (bundled data) + **E04** (extension fields + docs; no network). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-agent-schema-routing.md](01-agent-schema-routing.md) — schema + CRUD (must be **done**)
- Live: `internal/install/registry.go`, `cmd/trace/install.go`, `trace/agents/README.md`
- Pattern: `internal/install/githook.go` (target registry + idempotent install)

## Live baseline (after S09-01)

| Present | Absent |
|---------|--------|
| mig **027** + harness agent CRUD | `trace/agents/default.json` |
| `trace/agents/README.md` stub | `trace install agents` subcommand |
| `internal/install` target registry (cursor, claude, git-hook) | `TargetAgents` install target |
| Routing library (empty DB → no recs) | Bundled profiles in DB |

## Locked defaults

| Item | Value |
|------|-------|
| Catalog file | `trace/agents/default.json` — **`schema_version`: 1** |
| Docs | `trace/agents/README.md` — user extensions + future host (`registry_source`, `external_url`); **no network fetch in P22** |
| CLI | `trace install agents` — upsert bundled profiles + requirements into `.trace/` |
| Init flag | `trace init --with-agent-defaults` — optional; calls same loader |
| Idempotent | Re-run install → no duplicate requirements rows |
| `harness:subagent` | Upsert HOOK capability `hook:harness:subagent` — status **AVAILABLE** when install heuristics detect multitask harness (Cursor rules dir, `.cursor/` present, or env `TRACE_HARNESS_SUBAGENT=1`); else **UNKNOWN** (honest, not silent UNAVAILABLE) |
| SKILL stubs | Upsert referenced skill slugs as UNKNOWN if absent (metadata only) |
| Registry | **Do not** add `agents` to `install detect` JSON unless trivial — CLI subcommand is sufficient for P22 |

## `default.json` schema (v1)

```json
{
  "schema_version": 1,
  "registry_version": "2026-08-18",
  "agents": [
    {
      "slug": "agent:code-reviewer",
      "title": "Code Reviewer",
      "description": "Standards/maintainability review of a diff.",
      "subagent_type": "code-reviewer",
      "deliberation_phases": ["CRITIQUE", "VERIFY"],
      "task_keywords": [],
      "recommend_subagent": true,
      "registry_source": "bundled",
      "requirements": ["skill:code-review-and-quality"]
    }
  ]
}
```

## Minimum bundled profiles (locked)

| slug | subagent_type | phases | keywords | recommend_subagent | requirements |
|------|---------------|--------|----------|-------------------|--------------|
| agent:code-reviewer | code-reviewer | CRITIQUE, VERIFY | — | true | skill:code-review-and-quality |
| agent:performance-reviewer | performance-reviewer | VERIFY, EVALUATE | perf, performance, latency, benchmark | false | skill:performance-optimization |
| agent:security-reviewer | security-reviewer | CRITIQUE, VERIFY | security, auth, injection, owasp | true | skill:security-and-hardening |
| agent:nested-reviewer | nested-reviewer | CRITIQUE | — | true | skill:code-review-and-quality |
| agent:explore | explore | INVESTIGATE, ORIENT | — | false | mcp:codegraph_explore (optional) |
| agent:generalPurpose | generalPurpose | — (fallback) | — | false | — |

Each agent row: `registry_source: "bundled"`, `registry_version` from file header.

## Requirements

1. **`internal/agents/builtin.go`**: `LoadDefaultCatalog(path string) ([]HarnessAgentBundle, error)` — parse + validate JSON.
2. **`internal/install/agents.go`**: `InstallAgentDefaults(opts InstallOpts) error` — open store, upsert each agent + requirements + capability stubs.
3. **`cmd/trace/install.go`**: add `agents` case → `cmdInstallAgents`.
4. **`cmd/trace/init.go`**: `--with-agent-defaults` flag → call same installer after init.
5. Resolve bundled path: `filepath.Join(moduleRoot, "trace/agents/default.json")` or embed via `//go:embed` — follow existing trace asset patterns.
6. Update `trace/agents/README.md` with install examples + E04 extension point prose.
7. **No network fetch**, no hosted registry sync, no agent runner.

## Touch files

- `trace/agents/default.json` (new)
- `trace/agents/README.md` (extend)
- `internal/agents/builtin.go`, `internal/agents/builtin_test.go` (new)
- `internal/install/agents.go`, `internal/install/agents_test.go` (new)
- `cmd/trace/install.go`, `cmd/trace/init.go`
- `cmd/trace/help.go` — mention `install agents`

## Named tests

| Test | Proves |
|------|--------|
| `TestInstallAgentsSeedsDefaults` | install loads JSON into DB |
| `TestDefaultCatalogValidJSON` | schema validation |
| `TestBundledProfilesIncludeRequirements` | requirements rows for all 6 profiles |
| `TestInstallAgentsIdempotent` | second run no dupes |
| `TestSubagentHookDeclaredOnInstall` | `hook:harness:subagent` upserted |

```bash
go test ./internal/install/... ./internal/agents/... ./cmd/trace/... -count=1 -run 'TestInstallAgents|TestDefaultCatalog|TestSubagentHook'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] **E03** bundled catalog committed + loadable
- [ ] **E04** extension columns populated from JSON; README documents future host; no network code
- [ ] Named tests PASS; compat still **27**
- [ ] Board Notes: profile count + install path

## Minimal todos

- [ ] default.json + README
- [ ] builtin loader + install path
- [ ] CLI `install agents` + init flag
- [ ] Tests + board notes

## Residual risks (carry to S09-04)

- Module root resolution in tests vs installed binary
- Skill slug names must match real bundled skills or stay UNKNOWN honestly
- `mcp:codegraph_explore` may be absent on fresh clone — explore agent should still recommend with missing_capabilities
- Init flag discoverability in help text
