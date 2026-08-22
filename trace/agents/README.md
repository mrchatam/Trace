# Trace harness agent catalog

Trace maintains a **catalog of harness agents it understands** — profiles like `performance-reviewer`, `code-reviewer`, and `explore` — mapped to environment capabilities (skills, MCPs, hooks, tools).

## Important

- Trace **recommends** agents and subagents; it **does not spawn or run** them.
- The harness (Cursor, Claude Code, custom orchestrator) performs delegation.
- When `harness:subagent` is available, loop next may recommend a **fresh subagent** for independent review (CRITIQUE / post-EXECUTE).

## Default catalog

Bundled profiles live in `default.json` (schema v1, six agents). Install into your project store with:

```bash
trace install agents
trace init --with-agent-defaults   # optional one-shot during init
```

Re-run `trace install agents` after upgrading Trace to refresh bundled rows (idempotent upsert by slug).

### Bundled profiles

| slug | subagent_type | phases | requirements |
|------|---------------|--------|--------------|
| `agent:code-reviewer` | code-reviewer | CRITIQUE, VERIFY | `skill:code-review-and-quality` |
| `agent:performance-reviewer` | performance-reviewer | VERIFY, EVALUATE | `skill:performance-optimization` |
| `agent:security-reviewer` | security-reviewer | CRITIQUE, VERIFY | `skill:security-and-hardening` |
| `agent:nested-reviewer` | nested-reviewer | CRITIQUE | `skill:code-review-and-quality` |
| `agent:explore` | explore | INVESTIGATE, ORIENT | `mcp:codegraph_explore` |
| `agent:generalPurpose` | generalPurpose | (fallback) | — |

Install also declares `hook:harness:subagent` as **AVAILABLE** when a multitask harness is detected (`.cursor/` in the project, Cursor rules under `$HOME/.cursor/rules`, or `TRACE_HARNESS_SUBAGENT=1`); otherwise **UNKNOWN** (honest, not silent UNAVAILABLE). Referenced skill/MCP slugs are upserted as **UNKNOWN** stubs when absent.

## User extensions (E04)

Add JSON files under `trace/agents/` following the same v1 shape, or declare profiles via `trace_capability` / domain APIs, then re-run `trace install agents` or upsert directly. Extension fields on each agent row:

- `registry_source` — e.g. `bundled`, `local`, future `hosted`
- `registry_version` — catalog semver or date stamp from the file header
- `external_url` — optional link to upstream profile docs (not fetched in Phase 22)

Profiles reference capability slugs (`skill:…`, `mcp:…`, `hook:harness:subagent`). Missing capabilities surface as `missing_capabilities` in recommendations; Trace does not download skills or MCP servers.

## Future trace host

Schema supports `registry_source`, `registry_version`, and `external_url` for catalogs published by a hosted Trace product. Phase 22 does **not** fetch from the network — local-first only.

See [`docs/phases/phase-22-capability-completion/scopes/scope-09-harness-agent-catalog/`](../../docs/phases/phase-22-capability-completion/scopes/scope-09-harness-agent-catalog/).
