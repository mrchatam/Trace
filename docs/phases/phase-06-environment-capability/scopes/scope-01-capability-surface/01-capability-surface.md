# P06 / S01 / 01 — Capability surface

## Metadata
- id: P06-S01-01
- todo_ids: [P06-S01-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement the **minimal capability / environment surface**: catalog (skills/rules/MCP/tools/hooks) + task↔required attach + missing-capability warnings on Layer 0–1 packets. Extend `internal/domain` + store mig **`010_capability_surface`** + `internal/compiler` packet fields. Expose APIs S02 ablation can plant/score under `evals/capability`. Keep Gate F / Gate G / Gate E / honesty / p0x / x0 / Gate C bars green. No daemon/HTTP/embeddings. No ontology megastore. No new MCP write tools.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- [phase README](../../README.md)
- [docs/AGENT_ENVIRONMENT.md](../../../../AGENT_ENVIRONMENT.md)
- [docs/PROJECT_MODEL.md](../../../../PROJECT_MODEL.md) §10
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 6 / H7
- Live: MCP six tools; CLI thin argv; compiler Packet items/why/budget only; store migs `001`…`009`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Package | **`internal/domain`** on `*store.Store` + store helpers + **`internal/compiler`** packet attach. **Do not** invent `internal/capability` / ontology megastore / planner fork |
| Migration | Additive embed **`010_capability_surface.sql`** only (do not rewrite `001`–`009`) |
| entity_links | **No new rels** — requirements live in `task_capability_requirements` |
| Kinds | **`SKILL` \| `RULE` \| `MCP` \| `TOOL` \| `HOOK`** (UPPER; reject unknown/empty) |
| Status | **`AVAILABLE` \| `UNAVAILABLE` \| `UNKNOWN`** (UPPER; reject unknown; empty on create → **`UNKNOWN`**) |
| Slug | Unique TEXT; recommend `kind:name` lowercase (e.g. `mcp:trace_why`); uniqueness enforced; format not regex-gated |
| Packet | Additive `required_capabilities` + `missing_capabilities` only — **never** dump full catalog; keep `SchemaVersion` `"0.1"` |
| Missing rule | Required row whose capability is absent **or** `status != AVAILABLE` → missing |
| CLI | Thin G19: **`trace capability`** declare/list/require/unrequire/missing |
| MCP tools | **Not** required. Do **not** add MCP capability write tools; existing six stay G19 |
| Builtin mirror | `BuiltinMCPCapabilitySpecs()` returns six live tool names — **no** auto-seed on Open/init |
| CGO | Domain + store + compiler APIs must pass `CGO_ENABLED=0` where they already do |
| Carry-forward bars | Honesty A/B/C; Gate G (`TestHonestyEscapeRateGateGPrelim`); Gate E (`TestPlantedDiscoveryReplan`); Gate F (`TestPlantedImpactConflictsGateFPrelim`); p0x 7/7; x0; Gate C `dry_run:false` artifacts intact |
| Out | Ablation harness (S02); phase VERIFY (S03); VerifiedFact; `plan simulate`; daemon/HTTP/embeddings; ontology megastore; new MCP tools; GC-03/04 |

### Schema (locked)

```sql
-- 010_capability_surface.sql (additive)

CREATE TABLE IF NOT EXISTS capabilities (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'UNKNOWN',
    body TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_capabilities_kind ON capabilities(kind);
CREATE INDEX IF NOT EXISTS idx_capabilities_status ON capabilities(status);

CREATE TABLE IF NOT EXISTS task_capability_requirements (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    capability_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE(task_id, capability_id)
);
CREATE INDEX IF NOT EXISTS idx_task_capability_requirements_task
    ON task_capability_requirements(task_id);
CREATE INDEX IF NOT EXISTS idx_task_capability_requirements_cap
    ON task_capability_requirements(capability_id);
```

No ALTER on `tasks` — requirements live in the join table.

### Vocabulary (locked)

```text
# kind (reject unknown / empty)
CapabilityKindSkill = "SKILL"
CapabilityKindRule  = "RULE"
CapabilityKindMCP   = "MCP"
CapabilityKindTool  = "TOOL"
CapabilityKindHook  = "HOOK"

# status (empty → UNKNOWN; reject other unknowns)
CapabilityStatusAvailable   = "AVAILABLE"
CapabilityStatusUnavailable = "UNAVAILABLE"
CapabilityStatusUnknown     = "UNKNOWN"
```

### Builtin MCP specs (locked — helper only)

```text
BuiltinMCPCapabilitySpecs() []CapabilitySpec
  // kind=MCP, status=AVAILABLE, slug=mcp:<name>, title=<name>
  // names exactly:
  //   trace_why, trace_context, trace_add, trace_link, trace_transition, trace_review
  // Does NOT write to DB — callers UpsertCapability as needed (tests / S02 plants)
```

### Minimum public API (behavior locked; names may vary slightly)

```text
# Domain — catalog
UpsertCapability(ctx, CapabilityInput) (store.Capability, error)
  // id optional (generate UUID); kind required; slug required+unique;
  // status empty→UNKNOWN; reject bad kind/status
GetCapability(ctx, id) (store.Capability, error)
GetCapabilityBySlug(ctx, slug) (store.Capability, error)
ListCapabilities(ctx, ListCapabilitiesFilter) ([]store.Capability, error)
  // filter: optional Kind, optional Status; empty filter → all (ordered stable by slug)

# Domain — task requirements
RequireCapability(ctx, taskID, capabilityID) (store.TaskCapabilityRequirement, error)
  // task must exist; capability must exist; UNIQUE(task,cap) idempotent OK (return existing)
UnrequireCapability(ctx, taskID, capabilityID) error
  // missing row → no-op or ErrNotFound (pick one; document + test)
ListRequiredCapabilities(ctx, taskID) ([]store.Capability, error)
  // join requirements → capabilities; stable order by slug

# Domain — missing (S02 score hook)
MissingCapabilities(ctx, taskID) ([]store.Capability, error)
  // For each required: if Get fails OR status != AVAILABLE → include
  // Prefer returning Capability rows when present; if row deleted mid-flight treat as missing
  // Never silently drop UNKNOWN/UNAVAILABLE required caps

# Compiler — packet attach (called from TaskContext / ExpandContext)
  Packet.RequiredCapabilities []CapabilityRef  // json:"required_capabilities,omitempty"
  Packet.MissingCapabilities  []CapabilityRef  // json:"missing_capabilities,omitempty"
  CapabilityRef { ID, Kind, Slug, Title, Status }
  // Populate from ListRequiredCapabilities + MissingCapabilities for packet.TaskID
  // Markdown: "## Capabilities" with required + missing subsections
  // Do NOT list unattached catalog entries

# Store helpers (as needed)
UpsertCapability / GetCapability / GetCapabilityBySlug / ListCapabilities
InsertTaskCapabilityRequirement / DeleteTaskCapabilityRequirement /
ListTaskCapabilityRequirementsByTaskID
```

### Policy (locked)

```text
Task CRUD / work_state / DONE policy: UNCHANGED
Capability APIs are inventory + attach + warn only — do NOT auto-replan, mutate planner, or invent entity_links
Fail-closed: reject bad enums; missing warnings surface UNKNOWN/UNAVAILABLE; never claim "all required available" when MissingCapabilities non-empty
Packet: required+missing only — no full-catalog dump (H7 / minimal context)
MCP: no new tools; G19 library-only for existing six
Ontology: kinds above only — no agents/models/permissions tables
plan simulate / VerifiedFact: OUT
```

### Target tree

```text
internal/store/
  schema/010_capability_surface.sql
  capability.go              # Capability + TaskCapabilityRequirement + helpers

internal/domain/
  capability.go              # Upsert/Get/List/Require/Unrequire/Missing + BuiltinMCPCapabilitySpecs
  # existing files unchanged in behavior

internal/compiler/
  packet.go                  # CapabilityRef + Packet fields + Markdown section
  compiler.go                # attach on TaskContext / ExpandContext

cmd/trace/
  capability.go              # thin G19 subcommands
  help.go                    # usage lines
```

### Tests (required)

- Mig `010_capability_surface.sql` applied on Open (both tables present)
- `UpsertCapability` + Get by id/slug; status empty→UNKNOWN; reject bad kind/status; unique slug
- `RequireCapability` + `ListRequiredCapabilities`; UNIQUE idempotent; task/capability must exist
- `UnrequireCapability` behavior as documented
- `MissingCapabilities`: UNAVAILABLE + UNKNOWN required → missing; AVAILABLE → not missing
- `BuiltinMCPCapabilitySpecs` returns exactly the six MCP tool names with `mcp:` slugs
- Compiler: `TaskContext` packet includes required + missing; **no** unattached catalog rows; Markdown section present when any required
- Existing compiler/budget/DPC tests remain green
- Do **not** require `evals/capability` this scope (S02)

### CLI (thin G19)

```text
trace capability declare --kind SKILL|RULE|MCP|TOOL|HOOK --slug <slug> \
  [--title <t>] [--status AVAILABLE|UNAVAILABLE|UNKNOWN] [--body <text>] [--id <uuid>]
trace capability list [--kind <K>] [--status <S>]
trace capability require --task <id> --capability <id|slug>
trace capability unrequire --task <id> --capability <id|slug>
trace capability missing --task <id>
```

Stdout machine-friendly JSON lines (match existing `add`/`impact` style). Exit 0/1/2 per CLI norms.
`--capability` may accept id **or** slug (resolve via GetCapability then GetCapabilityBySlug).

### S02 ablation hooks (must expose — do not implement harness)

S02 planner/implement will plant via domain + compiler APIs (not CLI-only):

- `UpsertCapability` with AVAILABLE vs UNAVAILABLE/UNKNOWN for probes
- `RequireCapability` + `MissingCapabilities` for positive missing-warning probes
- `TaskContext` packet `required_capabilities` / `missing_capabilities` as scoreable signals
- Prefer harness path `evals/capability` / `TestPlantedCapabilitySelectionAblation` + `schema-capability.json` (S02-00 finalizes)

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] Mig `010_capability_surface.sql` + store capability helpers live
- [ ] Domain APIs: Upsert/Get/List; Require/Unrequire/ListRequired; `MissingCapabilities` as locked
- [ ] `BuiltinMCPCapabilitySpecs` returns six MCP slugs; no auto-seed on Open/init
- [ ] Compiler packet + Markdown attach required/missing only; SchemaVersion stays `"0.1"`
- [ ] No new entity_links rels; Task DONE / honesty policies unchanged
- [ ] Thin CLI `trace capability` wired (G19); no new MCP tools; no `internal/capability` package
- [ ] Domain/store/compiler tests cover cases above
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/compiler/...` green
- [ ] Carry-forward: `CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... ./evals/impact/...` + `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./...` green
- [ ] Gate C artifacts under `docs/verification/gate-c-x0/` untouched
- [ ] TODO.md status + Notes updated (this row only)

## Minimal todos
- [ ] Store: mig 010 + capability CRUD/list + task_requirement helpers
- [ ] Domain: capability APIs + BuiltinMCPCapabilitySpecs + MissingCapabilities; consts
- [ ] Compiler: CapabilityRef + packet fields + Markdown; attach in TaskContext/ExpandContext
- [ ] Thin CLI `trace capability` + help
- [ ] Full carry-forward bars; board Notes

## Out of scope
- `evals/capability` / named ablation test (S02)
- Phase VERIFY / Phase 07 handoff (S03)
- Ontology expansion (agents/models/permissions); embeddings; daemon/HTTP; VerifiedFact; `plan simulate`
- New MCP write tools; auto-seed on init; commercial multi-model theater
- GC-03/04 unless Notes promote
