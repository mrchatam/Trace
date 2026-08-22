# P06 / S01 / 00-PLANNER — Capability surface

## Metadata
- id: P06-S01-00
- todo_ids: [P06-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-capability-surface.md` for **skills/rules/MCP/tool capability inventory** against live Trace adapters. Lock package paths, persistence, CLI/packet attach, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 6
- [docs/AGENT_ENVIRONMENT.md](../../../../AGENT_ENVIRONMENT.md)
- [docs/PROJECT_MODEL.md](../../../../PROJECT_MODEL.md) §10
- Live: `internal/mcp` six tools; `cmd/trace` argv surface; `internal/compiler` Packet Layer 0–1; store mig `001`…`009`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (gaps — locked 2026-08-16)
| Item | Today (post–Phase 05 / P06-00) | S01 lock |
|------|--------------------------------|----------|
| MCP tools | `trace_why`/`trace_context`/`trace_add`/`trace_link`/`trace_transition`/`trace_review` (stdio; G19) | Catalog as kind=`MCP` slugs `mcp:<tool>` via `BuiltinMCPCapabilitySpecs()` helper — **no** new MCP write tools |
| CLI surface | Thin argv: init/index/add/link/transition/review/impact/plan/seed/why/context | Thin **`trace capability`** declare/list/require/unrequire/missing |
| Task context | `internal/compiler` Packet — items/why/budget; **no** capabilities fields | Attach `required_capabilities` + `missing_capabilities` only (never full catalog dump) |
| Skills/rules/hooks | Absent as first-class graph | Minimal catalog + task require table; kinds `SKILL`\|`RULE`\|`MCP`\|`TOOL`\|`HOOK` |
| Ontology | N/A | **Reject** megastore — no agents/models/permissions entities this scope |
| entity_links | Existing causal rels only | **No new rels** — dedicated tables (mirror P05 impact pattern) |
| Gate F/G/E/C bars | Green / Go | Must stay green |
| Migration | `001`…`009_decision_impact` | Additive **`010_capability_surface.sql`** |

## Phase defaults already locked (respect — P06-00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Honesty / Gate G / Gate E / Gate F / p0x / x0 / Gate C | Keep green / intact |
| Package hint | Prefer **`internal/domain`** + store + **`internal/compiler`** — avoid ontology megastore |
| Migration hint | Next additive **`010_*`** (name locked here) |
| Ablation path | S02 owns `evals/capability` planted harness |
| Daemon/HTTP/embeddings | Forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted |
| MCP | Capability ids / thin G19 — not new daemon |

## Planner work
- [x] Inventory live MCP/CLI/compiler vs capability-graph needs (confirm table above)
- [x] Lock package / schema / CLI / packet-attach model
- [x] Thicken `01-capability-surface.md` enough to run alone
- [x] Light-update upcoming S02 stubs with selection/ablation hooks from S01 surface
- [x] Sync SCOPE-TODOS + board Notes; mark this row done

## Locked defaults (set 2026-08-16 — do not re-debate in S01-01)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Carry-forward | Honesty A/B/C; Gate G/E/F; p0x; x0; Gate C `dry_run:false` intact; GC-03/04 deferred |
| Daemon/HTTP/embeddings | Forbidden as primary |
| Ontology bloat | Reject — minimal kinds only; no agents/models/permissions catalog |
| Migration | **`010_capability_surface.sql`** — tables `capabilities` + `task_capability_requirements` |
| Package | **`internal/domain`** + store helpers + **`internal/compiler`** packet attach — no `internal/capability` / ontology megastore / planner fork |
| entity_links | **No new rels** this scope |
| Kinds | `SKILL` \| `RULE` \| `MCP` \| `TOOL` \| `HOOK` (UPPER; reject unknown/empty) |
| Availability status | `AVAILABLE` \| `UNAVAILABLE` \| `UNKNOWN` (empty on create → `UNKNOWN`; reject other) |
| Slug | Unique TEXT; recommend `kind:name` lowercase (e.g. `mcp:trace_why`) — uniqueness enforced, format not regex-gated |
| APIs | `UpsertCapability` / `GetCapability` / `GetCapabilityBySlug` / `ListCapabilities`; `RequireCapability` / `UnrequireCapability` / `ListRequiredCapabilities`; `MissingCapabilities` (fail-closed) |
| Builtin MCP mirror | `BuiltinMCPCapabilitySpecs()` returns the six live tool names as kind=`MCP` slug=`mcp:<name>` — **no** auto-seed on `store.Open` / `init` |
| Packet attach | `Packet.RequiredCapabilities` + `Packet.MissingCapabilities` (+ Markdown section); **never** dump full catalog; keep `SchemaVersion` `"0.1"` (additive omitempty) |
| Missing rule | Required cap missing from store **or** status ≠ `AVAILABLE` → listed in `MissingCapabilities` |
| CLI | Thin **`trace capability`** (declare / list / require / unrequire / missing); G19 |
| MCP tools | **Out** this scope — keep existing six; packet fields flow via `trace_context` → compiler automatically |
| S02 hooks | Catalog CRUD + require/unrequire + `MissingCapabilities` + packet fields for `evals/capability` planted ablation |
| Ablation path (S02 owns) | Prefer `evals/capability` / `TestPlantedCapabilitySelectionAblation` + `schema-capability.json` (S02-00 finalizes) |

## Exit criteria
- [x] `01-capability-surface.md` runnable alone
- [x] Paths + model locked
- [x] Light S02 Depends note updated
- [x] No product Go in this row

## Minimal todos
- [x] Inventory MCP/CLI/compiler vs S01 needs
- [x] Thicken 01 + 02 + light S02 Depends
- [x] Mark P06-S01-00 done

## Out of scope
- Selection ablation harness (S02); phase VERIFY (S03); embeddings; daemon; commercial multi-model theater
