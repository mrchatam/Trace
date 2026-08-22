# P16-S05-00 — seed / impact packet + thin MCP (FINAL)

## Metadata
- id: P16-S05-00
- todo_ids: [P16-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `cmd/trace/seed.go`, `internal/compiler`, `cmd/trace/impact.go`, `internal/mcp`, and impact domain APIs, lock **FINAL** defaults for **DF-70, DF-71, DF-72, DF-73, DF-74**. DF-72 is a **thin** G19 `trace_impact` adapter ([`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md); P14 A3 superseded for this tool only). Thicken sibling `01`/`02`/SCOPE-TODOS + light S06 Depends. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition DF-70…74 **fix**; DF-72 **fix** (thin adapter)
- [../../DF-72-FORWARD.md](../../DF-72-FORWARD.md) — SoT over historical P16-00 DF-72 defer
- Live: `cmd/trace/seed.go`; `internal/compiler`; `cmd/trace/impact.go`; `internal/mcp`; `domain.LinkDiscoveryMentionsTask`; `domain.AddImpactFinding` / `AddDecisionAlternative` / `ImpactReport`
- Combo: [`experiments/BATCH-D21-D23.md`](../../../../../experiments/BATCH-D21-D23.md); D22 surfaces [`experiments/ab-combo-context-impact/results/_surfaces/`](../../../../../experiments/ab-combo-context-impact/results/_surfaces/)
- Keeper: `TestImportBoundaryMCPNoPlanImpactIndexTools` — **update** to allow `trace_impact` only
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-70…74
- Quality bar: [P16-S04-00](../scope-04-install-project-root/00-PLANNER.md) FINAL
- S04 live: **P16-S04-02 APPROVE** — DF-68 `-C` ProjectRoot; `install` ungated; nine MCP tools until this scope
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Depends-on: **P16-S04-02 APPROVE** (board, already true). Phase locks below. **Unattended:** no architecture blockers; defaults below are FINAL. **DF-72 is in scope** (thin adapter). Do **not** re-litigate daemon / install / decide MCP. Do **not** reopen S04 install `-C`. **No SwitchMode** (orchestrator). **No product Go.**

## Depends (from S04 — live P16-S04-02 APPROVE)

- **DF-68** is owned by S04 (`cmdInstall` keeps `root`; Abs `-C` → `InstallOpts.ProjectRoot` for detect/claude/uninstall; cursor omits ProjectRoot). Confirm vs [../scope-04-install-project-root/REVIEW-NOTES.md](../scope-04-install-project-root/REVIEW-NOTES.md).
- S05 must **not** reopen `cmdInstall`, Cursor STABLE home, CONDITIONAL markers, or `cli:install`.
- Do **not** add install/decide MCP. Do **not** thread `InstallOpts.ProjectRoot` into seed/impact/compiler.
- `install` stays **ungated**. Board-order only — this scope is seed / compiler packet / `trace_impact`.
- Catalog is still **nine** names including `trace_version` until 01 lands `trace_impact`.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| Seed link switch | Yes — `goal_has_task`/`goal-task`, `decision_affects_task`/`decision-task`, `discovery_causes_plan_change`/`discovery-plan-change`, `claim_has_evidence`/`claim-evidence` | **Omits** `LinkDiscoveryMentionsTask` (DF-70). Unknown rel → `seed: unknown link rel` |
| CLI/MCP `discovery-mentions-task` | Yes — DF-42 `cmd/trace/link.go` + `tools_write.go` → store `discovery_mentions_task` | Seed does not accept underscore **or** hyphen |
| Seed top-level allowlist | Yes — reject unknown keys | No `findings` / `alternatives` (DF-73). D22 planted via `hooks/g1-post-seed.sh` |
| `domain.AddImpactFinding` / `AddDecisionAlternative` / `ImpactReport` | Yes — tables + rollup live (mig ≤14) | Seed never calls them |
| Compiler `Packet` | Yes — items, why_trace, required/missing caps; `SchemaVersion` **0.2** | **No** impact fields (DF-71). D22 MD lists both decisions, omits DESTRUCTIVE / `overall_class` |
| `trace why` | Yes — encodes `retrieval.WhyResult` (`seed_type`, `steps`) | No impact attach (DF-71). Why is **not** a compiler packet today |
| `trace impact report` wrapper | Yes — snake_case `ok`, `decision_id`, `overall_class`, … | Nested `findings` / `alternatives` encode **untagged** store structs → PascalCase `ID`, `ImpactClass`, `IsRecommended` (DF-74; D22 `impact_report.json`) |
| `impact finding add` JSON | Yes — already maps `id` / `impact_class` | List dumps store structs (same PascalCase class) |
| MCP catalog | Yes — **nine**: why, context, add, link, transition, review, tasks, capability, version | **No** `trace_impact` (DF-72). Boundary keeper **forbids** the name |
| MCP Assert | Yes — `assertMCPToolAllowed` → `mcp:<Name>` at each tool | New tool must Assert at entry |
| `BuiltinMCPCapabilitySpecs` | Yes — nine `mcp:` slugs | Must add `mcp:trace_impact` |
| `BuiltinCLICapabilitySpecs` | Yes — already includes `cli:impact` | **Unchanged** (dual-slug; do not add `cli:trace_impact`) |
| Impact walk (P14) | Yes — `retrieval.ImpactWalk`; CLI `impact walk` | **Do not** rewrite R2 / walk semantics. MCP walk is a thin call |
| Compat / mig | Ceiling **14** | No new mig unless a test proves schema need (none live — reuse impact tables) |

**Bug path DF-70 (live):** D21/D23 `gt.json` omitted `discovery_mentions_task` links; `prepare.sh` planted via post-seed `trace link discovery-mentions-task`. Seed switch has no `RelDiscoveryMentionsTask` / `discovery-mentions-task` case.

**Bug path DF-71 (live):** D22 `context.md` items 4–5 are decisions (plain-text + append-only) with **no** finding class. `impact_report.json` has `overall_class: DESTRUCTIVE`. Agent must know to call `trace impact report`.

**Bug path DF-72 (live):** `user-trace` catalog has no impact tool. `RegisteredToolNames` length 9. `TestImportBoundaryMCPNoPlanImpactIndexTools` fatals on `trace_impact`.

**Bug path DF-73 (live):** seed allowlist has no `findings`/`alternatives`; D22 hook calls CLI `impact finding add` + `alternative add --recommended`.

**Bug path DF-74 (live):** `cmdImpactReport` wraps snake_case keys but `rep.Findings` / `rep.Alternatives` are `store.DecisionImpactFinding` / `DecisionAlternative` **without** `json` tags → encoding/json emits Go field names.

### Live MCP catalog (before this scope)

| # | Name |
|--:|------|
| 1 | `trace_why` |
| 2 | `trace_context` |
| 3 | `trace_add` |
| 4 | `trace_link` |
| 5 | `trace_transition` |
| 6 | `trace_review` |
| 7 | `trace_tasks` |
| 8 | `trace_capability` |
| 9 | `trace_version` |

After this scope: **ten** names — insert `trace_impact` **before** `trace_version`. Phrase “Ten tools + `trace_version`” in the phase README means this **ten-name** catalog **including** version (S04: “nine until S05”), **not** eleven.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Home DF-70/73 | Thin `cmd/trace/seed.go` (+ `cli_test.go`). Domain link/impact APIs **already exist** — no new `entity_links` rels, no new tables |
| Home DF-71 | `internal/compiler` Packet attach + Markdown `## Impact`. Why CLI/MCP **inherit** via domain helper (do **not** change Why BFS / P14 walk) |
| Home DF-74 | `json` tags on `store.DecisionImpactFinding` + `store.DecisionAlternative` (canonical snake_case). CLI report/list and MCP inherit |
| Home DF-72 | New thin file `internal/mcp/tools_impact.go` + register in `server.go`. Pattern = `trace_capability` (`action` switch). G19: call existing domain / `retrieval.ImpactWalk` only |
| Seed rel DF-70 | Switch accepts `domain.RelDiscoveryMentionsTask` (`discovery_mentions_task`) **and** hyphen `discovery-mentions-task` → `LinkDiscoveryMentionsTask`. Same dual-token pattern as `decision_affects_task` / `decision-task` |
| Seed keys DF-73 | Top-level **`findings`** and **`alternatives`** (arrays, optional). Add both to the unknown-key allowlist. Historical stub names `impact_findings` / `decision_alternatives` stay **rejected** as unknown |
| Unknown keys | Unknown **other** top-level keys still rejected. Nested extra fields: ignore (same as live `seedEntity`) — do **not** invent a nested allowlist |
| Seed finding object | `decision_id` (required), `impact_class` (required), `kind` (required), `uncertainty` (optional, empty → UNKNOWN), `body`, `related_type`, `related_id`, `id` (all optional) |
| Seed alternative object | `decision_id` (required), `title` (required), `body`, `id` (optional), `recommended` (bool, optional, default false) → `AlternativeInput.Recommended`. Input key is **`recommended`**, not `is_recommended` |
| Seed order | Entities → links → **findings** → **alternatives** → transitions. Findings/alternatives after decisions exist; `related_id` may point at a same-doc task |
| Seed summary | Keep `ok` / `created` / `links` / `transitions`; add counts `findings` and `alternatives` (integers) |
| Packet DF-71 | Additive field `impact` `[]DecisionImpact` `json:"impact,omitempty"`. **Do not** bump `SchemaVersion` (stays **0.2**) |
| `DecisionImpact` | `decision_id`, `overall_class`, `overall_uncertainty` (omitempty), `has_unknown`, `incomplete`, `findings` (store rows after snake_case tags). **No** alternatives on the packet (report remains the alternatives surface) |
| When to attach | Neighborhood = `ListLinksTo(task, taskID)` with `rel=decision_affects_task`. Call `domain.ImpactReport` per decision (no forked rollup). Include a decision **iff** `len(Findings) > 0`. Empty → omit `impact` entirely |
| Attach vs budget | Impact attach is **after** item trim (like capabilities). Do **not** drop DESTRUCTIVE because the decision item was trimmed |
| Markdown | `## Impact` after Items, before Why trace / Capabilities. Must show `overall_class` and each finding `impact_class` so MD-only agents see `DESTRUCTIVE` |
| Why inherit | Do **not** change `retrieval.Why` algorithm. After Why: if seed is **task**, merge `impact` from `ImpactSummariesForTask`; if seed is **decision**, merge that decision’s report when it has findings. CLI `cmd/trace/why.go` + MCP `toolWhy` encode the extra field. JSON key **`impact`** (same shape as packet) |
| Domain helper | `Service.ImpactSummariesForTask(ctx, taskID)` (name may vary; one helper). `compiler` **may import** `domain` solely to call ImpactReport / this helper. Retrieval stays store-only |
| Report JSON DF-74 | Nested findings/alternatives snake_case: `id`, `decision_id`, `impact_class`, `uncertainty`, `kind`, `body`, `related_type`, `related_id`, `created_at`, `updated_at`; alts: `id`, `decision_id`, `title`, `body`, `is_recommended`, `created_at`, `updated_at`. Wrapper keys already snake_case — keep |
| Finding/alternative **list** | Same DTOs (tags). Do **not** leave list PascalCase while report is snake_case |
| Walk JSON | **Unchanged** (P14). MCP walk returns the same library result the CLI already encodes. Do **not** retag walk/blast types in this scope |
| MCP tool | Name `trace_impact`. Slug **`mcp:trace_impact`**. `assertMCPToolAllowed(ctx, st, "trace_impact")` at **every** action entry (after `openStore`, before domain) |
| MCP `action` | Required: `finding` \| `alternative` \| `report` \| `walk`. Unknown action → error (fail-closed) |
| MCP `op` | `finding`: `add` \| `list` (required). `alternative`: `add` \| `list` \| `recommend` (required). `report` / `walk`: omit `op` (ignore if present) |
| MCP fields | See input-schema table below. `project` optional (existing `-C` / cwd). `decision` primary, `decision_id` alias (capability `task`/`task_id` pattern). Finding add: `class` primary, `impact_class` alias |
| MCP annotations | One mixed tool: `ReadOnlyHint=false`, `DestructiveHint=false`, `OpenWorldHint=false` (same class as `trace_capability`; does not mutate plan/tasks) |
| Catalog | **Ten** names including `trace_version`. Order: why, context, add, link, transition, review, tasks, capability, **`trace_impact`**, version |
| Specs | `BuiltinMCPCapabilitySpecs` same ten `mcp:` slugs. `BuiltinCLICapabilitySpecs` **unchanged** (`cli:impact` already exists). Dual-slug: denying `mcp:trace_impact` does **not** block `trace impact` |
| Boundary keeper | Allow `trace_impact` (must be registered). Still **forbid** `trace_plan`, `trace_index`, and any `trace_install` / `trace_decide` |
| G19 | No business logic in MCP/CLI adapters. No new impact rollup. No daemon/HTTP |
| Compat | Ceiling **14**; **no** mig 015+. **Do not edit** done `014_*.sql` nine-Name IN list — live canonicalize uses `BuiltinMCPCapabilitySpecs` (new `trace_impact` AUTO_ALLOWs; no unprefixed rows exist yet) |
| Forbidden | Daemon; new entity_links rels; rewriting P14 impact walk (R2); install/decide/plan/index MCP; reopening S04 `-C`; `cli:install`; bumping packet schema to 0.3; forking OverallClass |

### Seed JSON v1 — locked keys (additive)

```text
allowed top-level += "findings", "alternatives"

findings[]:   decision_id, impact_class, kind  [uncertainty, body, related_type, related_id, id]
alternatives[]: decision_id, title             [body, id, recommended]
```

Call `AddImpactFinding` / `AddDecisionAlternative` (existing validation: class/kind vocab, decision must exist). Do not insert via store bypass.

### Compiler packet — locked fields

```text
Packet.impact []DecisionImpact `json:"impact,omitempty"`
DecisionImpact:
  decision_id, overall_class, overall_uncertainty?, has_unknown, incomplete, findings[]
SchemaVersion remains "0.2"
```

### MCP `trace_impact` input schema (FINAL)

| Field | JSON | Required when | Notes |
|-------|------|----------------|-------|
| project | `project` | never | existing root override |
| action | `action` | **always** | `finding` \| `alternative` \| `report` \| `walk` |
| op | `op` | finding / alternative | `add` \| `list` \| `recommend` |
| decision | `decision` | finding, alternative, report | UUID; alias `decision_id` |
| class | `class` | finding add | alias `impact_class`; vocab SAFE\|CAUTION\|HIGH\|DESTRUCTIVE\|REVERSAL |
| kind | `kind` | finding add | AFFECTED_WORK\|INVALIDATED_ASSUMPTION\|WORK_AT_RISK\|NEW_WORK\|UNRESOLVED |
| uncertainty | `uncertainty` | no | empty → UNKNOWN |
| body | `body` | no | finding add / alternative add |
| related_type / related_id | `related_type`, `related_id` | no | finding add |
| title | `title` | alternative add | |
| recommended | `recommended` | no | alternative add bool |
| id | `id` | alternative recommend; optional on add | finding/alternative UUID |
| seeds | `seeds` | walk | string array `file:<uuid>` \| `symbol:<uuid>` (CLI `--seed` repeatable) |
| depth | `depth` | no | walk 1\|2; default library `DefaultImpactDepth()` (2) |

Report JSON keys **match** CLI `cmdImpactReport` map: `ok`, `decision_id`, `affected_task_ids`, `findings`, `alternatives`, `overall_class`, `overall_uncertainty`, `has_unknown`, `incomplete`.

Walk: parse seeds like CLI (`Cut` on `:`, type file\|symbol); `retrieval.New(st).ImpactWalk` — **no** git open required (CLI walk does not). Do not change depth clamp / blast semantics.

### Catalog after this scope (locked order)

1. `trace_why` 2. `trace_context` 3. `trace_add` 4. `trace_link` 5. `trace_transition` 6. `trace_review` 7. `trace_tasks` 8. `trace_capability` 9. **`trace_impact`** 10. `trace_version`

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| `TestSeedImportDiscoveryMentionsTask` | `cmd/trace` | Seed `rel: discovery_mentions_task` **and** hyphen `discovery-mentions-task` (subtests OK) → store `discovery_mentions_task` via `LinkDiscoveryMentionsTask`. Unknown rel still usage-error |
| `TestContextIncludesImpactOverallClass` | `internal/compiler` | Task + `decision_affects_task` + DESTRUCTIVE finding: Packet JSON `impact` has that `decision_id`, `overall_class=DESTRUCTIVE`, findings with `impact_class` (not `ImpactClass`). Neighbor decision with **zero** findings omitted. Markdown contains `DESTRUCTIVE` |
| `TestWhyIncludesImpactOverallClass` | `cmd/trace` | Same fixture: `trace why task <id>` JSON includes `impact` + `overall_class` DESTRUCTIVE. Decision seed with findings also includes `impact` (subtest OK) |
| `TestSeedImportImpactFindings` | `cmd/trace` | Seed v1 `findings` + `alternatives` (`recommended: true`) → `impact report` shows snake_case finding + `is_recommended`. Summary counts. Top-level `impact_findings` still `unknown top-level key` |
| `TestImpactReportJSONSnakeCase` | `cmd/trace` | After `finding add` + `alternative add --recommended`, `impact report` stdout has `impact_class` / `is_recommended` / `overall_class` and must **not** contain `ImpactClass` / `IsRecommended` / `"ID"` as object keys |
| `TestMCPTraceImpactReport` | `internal/mcp` | CallTool `action=report` after planted finding: JSON `overall_class`, `findings[].impact_class`; Assert not skipped (AUTO_ALLOWED default succeeds) |
| `TestMCPImpactDeniedBlocksCallTool` | `internal/mcp` | `DecideTool` DENIED `mcp:trace_impact` (or unprefixed `trace_impact` per S02 canonicalize) → report CallTool fails closed; no domain write |
| `TestToolNamesRegistered` | `internal/mcp` | **Update** want list length **10**; includes `trace_impact` before `trace_version` |
| `TestBuiltinMCPCapabilitySpecs` | `internal/domain` | **Update** ten `mcp:` slugs including `mcp:trace_impact` |
| `TestImportBoundaryMCPNoPlanImpactIndexTools` | `internal/mcp` | **Update:** `trace_impact` **must** be registered; still fatal on `trace_plan` / `trace_index` / install / decide MCP names |
| Keepers | various | `TestLinkDiscoveryMentionsTask` / CLI; `TestMCPVirginProjectDoesNotMkdir`; `TestCLIAddDeniedFailClosed`; `TestPlantedImpactConflictsGateFPrelim`; `TestOpenCreatesDBAndMigratesIdempotent` (v**14**); `TestInstallClaudeDashCRefuseCitesProjectRoot` (S04 — do not regress) |

TDD: named tests first (red: seed unknown rel / unknown `findings` key; packet has no `impact`; report PascalCase nested; MCP name absent / boundary forbids it), then wire switch + allowlist + tags + Packet attach + MCP adapter (green). Do **not** change ImpactReport rollup or Why BFS to make tests pass.

Implementer may add `TestMCPTraceImpactFindingAdd` / `TestMCPTraceImpactWalk` as extra coverage; they are **not** a substitute for `TestMCPTraceImpactReport`.

## Owns

| Item | Intent |
|------|--------|
| DF-70 | Seed import accepts underscore + hyphen mentions-task rels → existing domain link |
| DF-71 | Compiler context packet (JSON+MD) + why JSON include findings / `overall_class` when present |
| DF-73 | Seed JSON v1 `findings` / `alternatives` → existing impact APIs |
| DF-74 | Impact report (and list) nested JSON snake_case |
| DF-72 | Thin MCP `trace_impact`; slug `mcp:trace_impact`; Assert; catalog ten; boundary allow this name only |

## Explicit deferrals

- Daemon / HTTP / hosted MCP
- MCP `trace_plan` / `trace_index` / install / decide tools
- New `entity_links` relations; new impact tables / mig 015
- Rewriting P14 impact walk / R2 `allowContainsOut`
- Packet `SchemaVersion` 0.3; putting alternatives on the context packet
- Reopening S04 install `-C` / Cursor STABLE / `cli:install`
- Session-global DENY across `project=` roots
- PID-kill Cursor reload (DF-22/37)
- Unprefixed command slugs; changing `cli:impact` gating
- P17 `trace/graph.json` export (Depends-on these seed keys; do not implement here)

## Assumptions (unattended)

1. **Domain is correct:** DF-70/73 are seed switch + allowlist gaps. Do not reopen DF-42 link semantics or impact vocab.
2. **DF-74 is tags, not a new DTO layer:** tagging store structs also fixes finding/alternative list. Walk types stay untagged.
3. **Catalog math:** nine including version today → **ten including version** after adding one tool. Not eleven.
4. **Why is adapter inherit:** retrieval Why algorithm unchanged; extra JSON field is G19 merge from domain summaries.
5. **Attach iff findings exist:** D22 decision a1 (no findings) omitted; a2 DESTRUCTIVE included. Incomplete-because-empty does not spam the packet.
6. **compiler → domain import** is allowed for ImpactReport only (avoids forked rollup). No domain → compiler cycle.
7. **CGO:** named `cmd/trace` tests run **CGO1**. `internal/compiler`, `internal/mcp`, `internal/domain`, `internal/store` named tests stay CGO0.
8. S06 VERIFY imports the named tests in the table + catalog ten + Gate F keeper + compat **14**.

## Effects on later scopes

- **S06:** Import named tests above (seed mentions-task + seed findings + context overall_class + why overall_class + report snake_case + `TestMCPTraceImpactReport` + `TestMCPImpactDeniedBlocksCallTool` + updated `TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs` / boundary keeper). Claim DF-70…74 only when those pass. Catalog **10** including `trace_version`. Still no install/decide/plan/index MCP. Gate F keeper must stay green (walk unrewritten).
- **P17:** seed export of `discovery_mentions_task` + `findings`/`alternatives` Depends-on this lock. Do not implement P17 here; do not rewrite P17-00.

## Planner work

1. [x] Inventory live seed switch / compiler Packet / impact JSON / MCP catalog + D21–D23 dogfood
2. [x] Lock seed keys, compiler fields, MCP input schema, catalog ten
3. [x] Lock named tests per DF including `TestMCPTraceImpactReport`
4. [x] Thicken `01-seed-impact-packet.md` + `02-scope-review.md` + SCOPE-TODOS
5. [x] Light S06 Depends (named tests + ten-tool catalog)
6. [x] Board Notes → next **P16-S05-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL** with seed keys + compiler fields + MCP schema + catalog 10
- [x] DF-72 thin MCP locked (not defer)
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P16-S05-01**

## Minimal todos
- [x] Inventory live seed / compiler / impact / MCP / DF-70…74 dogfood / S04 APPROVE
- [x] FINAL locks + named tests
- [x] Thicken 01/02/SCOPE-TODOS + S06 pointers
- [x] Board sync

## Next
**P16-S05-01** (implement DF-70/71/72/73/74).
