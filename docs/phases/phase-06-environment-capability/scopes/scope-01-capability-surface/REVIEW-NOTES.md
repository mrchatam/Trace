# P06-S01-02 — Scope review notes (2026-08-16)

Independent review of S01 against `00-PLANNER.md` / `01-capability-surface.md` + TODO Notes for `P06-S01-01`. Fresh session; claims verified in-repo (no implementer assumptions reused).

## Plan (executed)

1. Diff claims vs mig `010_capability_surface.sql`, `internal/store/capability.go`, `internal/domain/capability.go`, compiler packet attach, thin `cmd/trace/capability.go`
2. Check kinds/status enums, APIs, BuiltinMCP six specs (no auto-seed), packet required+missing only / SchemaVersion `0.1`, no `internal/capability` / megastore / planner fork / new MCP tools / new entity_links rels
3. Confirm S02 ablation hooks usable; Gate C `dry_run:false` intact; GC-03/04 still deferred
4. Re-run honesty / Gate G / Gate E / Gate F / p0x / x0 / `./...`
5. Severity-tag findings; spawn only for blocker/high (none)
6. Write these notes; mark board; light S02 Depends confirm

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**. Next board row: **P06-S02-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Mig **`010_capability_surface.sql`** additive; tables `capabilities` + `task_capability_requirements`; no ALTER on `tasks` | Pass (schema matches lock; store Open applies v10; `store_test` lists both tables; `001`–`009` untouched) |
| Package **`internal/domain`** + store + **`internal/compiler`**; **no** `internal/capability` / megastore / planner fork | Pass (`internal/capability` absent; no capability APIs under planner) |
| Kinds `SKILL`\|`RULE`\|`MCP`\|`TOOL`\|`HOOK`; status `AVAILABLE`\|`UNAVAILABLE`\|`UNKNOWN` (empty→UNKNOWN; reject other) | Pass (`Normalize*` + `TestUpsertCapabilityGetAndReject`) |
| APIs Upsert/Get/List + Require/Unrequire/ListRequired + `MissingCapabilities` | Pass (`capability.go` + `TestRequireUnrequireAndMissing`) |
| `BuiltinMCPCapabilitySpecs` six `mcp:trace_*`; **no** auto-seed on Open/init | Pass (`TestBuiltinMCPCapabilitySpecs`; Open catalog empty; init/open have no Upsert) |
| Packet `required_capabilities` + `missing_capabilities` only; SchemaVersion `"0.1"`; MD `## Capabilities`; no catalog dump | Pass (`TestTaskContextIncludesRequiredAndMissingCapabilities`; attach on shared TaskContext/ExpandContext path) |
| Thin CLI `trace capability` declare\|list\|require\|unrequire\|missing; G19 | Pass (`cmd/trace/capability.go` + help + root) |
| **No** new MCP tools; existing six only | Pass (`registerTools` still six names; no capability MCP tools) |
| **No** new entity_links rels | Pass (requirements in join table; Rel consts unchanged) |
| S02 hooks: catalog + require + `MissingCapabilities` + packet fields + BuiltinMCP + mig 010 | Pass (match S02-00 Depends table) |
| Honesty A/B/C + Gate G + Gate E + Gate F | Pass |
| p0x 7/7; x0; `CGO_ENABLED=0` domain+store+compiler; `CGO_ENABLED=1` `./...` | Pass |
| Gate C `docs/verification/gate-c-x0/` `dry_run:false`; N=3 G1 0.8 / B0 0.0 intact | Pass |
| GC-03/04 | Still deferred (no promotion in Notes) |

## Re-verification commands (2026-08-16)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/compiler/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... -count=1
# ok

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
# ok (p0x + x0 + full tree)
```

Domain/compiler capability tests: `TestUpsertCapabilityGetAndReject`, `TestListCapabilitiesFilter`, `TestRequireUnrequireAndMissing`, `TestBuiltinMCPCapabilitySpecs`, `TestResolveCapabilityIDOrSlug`, `TestTaskContextIncludesRequiredAndMissingCapabilities`.

## Findings

### Blocker / high

None.

### Medium

None requiring spawn.

### Low / nit (residual)

- **Packet missing vs domain `MissingCapabilities`:** Compiler `attachTaskCapabilities` uses `ListCapabilitiesRequiredByTaskID` (JOIN) + inline `status != AVAILABLE`. Equivalent for all product-API paths (no `DeleteCapability`). Domain `MissingCapabilities` also walks requirements and surfaces orphan placeholders if a capability row were deleted out-of-band; packet would omit those orphans. S02 plants via Upsert+Require — no gap for ablation.
- **Silent list failure:** attach returns empty on store list error (same class as optional Why attach swallow). Fail-open on I/O rare; harden later if needed.
- **CLI coverage:** No dedicated `cmd/trace` capability integration test; thin adapter + domain/compiler tests cover policy.
- **`ctx` ignored** on capability APIs (`_ = ctx`) — consistent with other domain methods.

## Spawns

None.

## S02 note

Upcoming `P06-S02-00` Depends hooks match the shipped surface (`UpsertCapability` / `ListCapabilities` / `RequireCapability` / `ListRequiredCapabilities` / `MissingCapabilities` / packet `required_capabilities`+`missing_capabilities` / `BuiltinMCPCapabilitySpecs` / mig `010_capability_surface`; CLI optional; harness must call library APIs). No substantive prompt rewrite required — live-confirm stamp applied on S02 planner Depends.
