# Scope S01 — Capability surface

**Depends-on:** Phase 05 complete; `P06-00` done.

**Out:** Selection ablation (S02); phase VERIFY (S03); daemon/HTTP/embeddings primary; ontology megastore.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P06-S01-00 | planner | done | 2026-08-16: locked mig `010_capability_surface`; kinds SKILL\|RULE\|MCP\|TOOL\|HOOK; domain+store+compiler; thin `trace capability`; packet required+missing; no new MCP tools / entity_links |
| P06-S01-01 | implement | done | 2026-08-16: mig 010 + domain/store/compiler + CLI; BuiltinMCP no auto-seed; bars green |
| P06-S01-02 | review | done | 2026-08-16: APPROVE high; [REVIEW-NOTES.md](REVIEW-NOTES.md); next P06-S02-00 |

## Checklist

- [x] P06-S01-00 planner
- [x] P06-S01-01 implement
- [x] P06-S01-02 review

## Expected S02 hooks (locked by S01-00; live-confirmed P06-S01-02)

- `UpsertCapability` / `GetCapability` / `GetCapabilityBySlug` / `ListCapabilities` (kinds + status)
- `RequireCapability` / `UnrequireCapability` / `ListRequiredCapabilities`
- `MissingCapabilities(taskID)` — required ∩ (absent \| status ≠ AVAILABLE)
- Compiler packet `required_capabilities` + `missing_capabilities` (no catalog dump)
- `BuiltinMCPCapabilitySpecs()` for six MCP tool slugs (`mcp:trace_*`) — plant via Upsert, no auto-seed
- Mig `010_capability_surface.sql` only — S02 must not fork schema
- CLI optional for humans; **harness must call library APIs** (G19)
