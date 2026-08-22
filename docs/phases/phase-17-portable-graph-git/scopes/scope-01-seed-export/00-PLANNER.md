# P17-S01-00 — seed export (FINAL)

## Metadata
- id: P17-S01-00
- todo_ids: [P17-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** defaults for **DF-80** (`trace seed export` + import↔export round-trip), **DF-84** (plan tree in seed), and **DF-85** (`exported_at_commit` on export). Thicken sibling `01`/`02`/SCOPE-TODOS + light S02/S03 Depends. **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md) — SoT for plan-tree keys + `exported_at_commit`
- Research: [docs/research/PORTABLE-GRAPH-GIT-2026-08-17.md](../../../../research/PORTABLE-GRAPH-GIT-2026-08-17.md)
- P16 S05 (depend, do not duplicate): [../../../phase-16-assert-root-and-surfaces/scopes/scope-05-seed-impact-packet/00-PLANNER.md](../../../phase-16-assert-root-and-surfaces/scopes/scope-05-seed-impact-packet/00-PLANNER.md)
- Live: `cmd/trace/seed.go`; `cmd/trace/help.go`; `fixtures/x0/seed/gt.json`; `internal/store/schema/006_plan_hierarchy.sql`; `internal/store/plan_hierarchy.go`; `internal/gitcli/open.go` (`Head` / `rev-parse HEAD`)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Depends-on: **Phase 16 complete** (P16-S06-02 APPROVE — board, already true). Phase locks below. **No product Go.**

## Depends (from P16 S05 — live)

- Seed import already accepts **`findings`** / **`alternatives`** (DF-73) and **`discovery_mentions_task`** + hyphen alias (DF-70). S01 **exports** those keys when present; **must not** re-implement import rels or impact domain.
- P16 keepers stay green: `TestSeedImportDiscoveryMentionsTask`, `TestSeedImportImpactFindings`, `TestSeedImportAndWhy`, etc.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| CLI `seed` | **Import only** — `cmdSeed` requires `args[0] == "import"` | No `export` subcommand (DF-80) |
| Help | `seed import <file>` only | No export line |
| `seedDocument` struct | Causal entities + `findings`/`alternatives` + `transitions` | No plan-tree fields; no `exported_at_commit` |
| Import allowlist | 13 keys (incl. `findings`, `alternatives`) | Missing `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state`, `exported_at_commit` |
| Plan store (mig 006) | `InsertPlanPhase`, `InsertPlanScope`, `InsertScopeDeepPlan`, `UpsertGoalPlanState` | List helpers are **goal/phase scoped** and often **ACTIVE-only** — export needs **all-row** readers |
| Entity lists | `ListGoals`, `ListTasks` | No `ListDecisions` / `ListAssumptions` / … — export needs list-all helpers or one query |
| Links | `ListLinksByRel` | Export filters to seed rel set |
| Git HEAD | `internal/gitcli` `Head()` → `rev-parse HEAD` | Not wired to seed; requires repo (not store) |
| G19 | Import is thin CLI over domain | Export should mirror: **library builds doc**, CLI writes JSON |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| CLI argv | **`trace seed export [-o <file>]`** — stdout when `-o` omitted. **`trace seed import <file>`** unchanged. `cmdSeed` dispatches on `args[0]` (`import` \| `export`). Missing/unknown subcommand → usage (exitUsage) |
| `-o` path | Relative paths resolve under `-C` project root (same rule as import via `resolveSeedPath`). Absolute paths unchanged. Parent dirs for `-o` **must** be created (`MkdirAll`) before write |
| Stdout JSON | Indented two-space JSON + trailing newline (match import summary style). **No** extra wrapper object |
| Recommended path | **`trace/graph.json`** (S02 docs; `-o` is authoritative) |
| Format | Seed JSON **v1** (`version: 1`). Unknown top-level keys still **rejected** on import (except newly allowed keys below) |
| G19 home | New **`internal/domain/seed_export.go`** (or `seed_document.go`) owns `BuildSeedDocument(ctx, *store.Store, ExportOpts)` → JSON-marshalable struct. **`cmd/trace/seed.go`** stays thin: parse flags, open store, call domain, encode stdout/file. **No** MCP tool |
| Git SHA (DF-85) | Top-level **`exported_at_commit`**: full `git rev-parse HEAD` OID when bound `-C` root is inside a git work tree; **omit key** when not a repo / git missing / HEAD unavailable. **Never** use SHA as entity or plan-row id |
| Import `exported_at_commit` | Add to allowlist; **accept and ignore** — do not persist, do not merge, do not fail import |
| Round-trip | **import → export → import** (fresh DB second import) preserves entity ids, link endpoints (`rel`+`from`+`to`), **and** plan-tree row ids (`plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state`) |
| S01 import plan rows | **First insert** via existing `InsertPlanPhase` / `InsertPlanScope` / `InsertScopeDeepPlan` / `UpsertGoalPlanState` (same as today — duplicate id **may fail** until S03). **Idempotent upsert** is **S03** only |
| Forbidden | MCP `trace_seed`; committing `.trace/`; duplicating P16 S05 import; default transition **export**; HTTP/daemon; SHA-as-id; board spawn from implementer |

### Export include (default)

| Bucket | Source | Notes |
|--------|--------|-------|
| `goals`, `tasks`, `decisions`, `assumptions`, `discoveries`, `plan_changes`, `claims`, `evidence` | Store entity tables | Tasks: **`id`, `title`, `body`, `goal_id` only** — no `work_state` |
| `links` | `entity_links` + synthesized goal→task | Canonical **underscore** rels (see below). Endpoints **`from` / `to`** (not `from_id`/`to_id`) |
| `findings`, `alternatives` | Impact tables | When present (P16 S05 shapes) |
| `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` | mig **006** tables | **All rows** (all statuses). Deep plans: **ACTIVE + SUPERSEDED** |
| `exported_at_commit` | `git rev-parse HEAD` | When known (DF-85) |

**Link rels exported** (store truth; import still accepts hyphen aliases):

| Export `rel` | Source |
|--------------|--------|
| `goal_has_task` | Synthesized from `tasks.goal_id` (one link per task with non-empty goal) |
| `decision_affects_task` | `entity_links` |
| `discovery_causes_plan_change` | `entity_links` |
| `claim_has_evidence` | `entity_links` |
| `discovery_mentions_task` | `entity_links` |

Do **not** export review/capability/index rels or event-only `goal_has_task` payloads without a task row.

### Export exclude (default — must not appear in export JSON)

| Surface | Rationale |
|---------|-----------|
| `transitions` | Process replay; clone tasks stay PENDING (DF-28 class) |
| Index (`files`, `symbols`, `imports`, FTS, VCS watermark) | Derived — rebuild with `index` |
| `.trace/access.token`, lock | Local auth / single-writer |
| `capability_tool_decisions`, capabilities, task capability requirements | Per-machine allowlist |
| `events` | Audit log; re-import would duplicate |
| `reviews`, review residuals | Local process |
| Task `work_state` | Omitted from task objects |

**Empty arrays:** emit `[]` for included buckets that are empty (stable round-trip). **Omit** `transitions` key entirely. **Omit** `exported_at_commit` when empty.

### Seed JSON v1 — locked top-level keys

**Import allowlist** (reject any other top-level key):

```text
version
goals, tasks, decisions, assumptions, discoveries, plan_changes, claims, evidence
links
findings, alternatives          # P16 S05 — already live
plan_phases, plan_scopes, scope_deep_plans, goal_plan_state   # DF-84 — S01 additive
exported_at_commit              # DF-85 — accept on import; ignore for identity
transitions                     # import-only; never exported by default
```

**`seedDocument` struct** must add JSON fields for DF-84/85 (implementer). Import path: extend allowlist + unmarshal + import loops.

#### Plan-tree object shapes (match mig 006 / DF-84-FORWARD — do not invent names)

| Key | Object fields (required bold) |
|-----|--------------------------------|
| `plan_phases[]` | **`id`, `goal_id`, `title`, `body`, `ord`, `status`** — timestamps optional on seed v1 |
| `plan_scopes[]` | **`id`, `phase_id`, `title`, `body`, `ord`, `status`, `auto_replan_count`** |
| `scope_deep_plans[]` | **`id`, `scope_id`, `content_json`** (string JSON text), **`status`** |
| `goal_plan_state[]` | **`goal_id`, `current_scope_id`** (nullable string; JSON `null` when unset) |

Export **`content_json`** as stored (string). Import passes through to `InsertScopeDeepPlan`.

#### Import order (S01 — lock)

1. Causal entities (existing order: goals → tasks → … → evidence)
2. Links → findings → alternatives (existing)
3. **`plan_phases` → `plan_scopes` → `scope_deep_plans` → `goal_plan_state`**
4. Transitions (if present — existing behavior)

Parse **`exported_at_commit`** before/after doc unmarshal but **do not branch** on it.

### Store helpers (implementer — expected)

Add narrow **list-all** readers in `internal/store` (names may vary):

- All goals/tasks/decisions/assumptions/discoveries/plan_changes/claims/evidence
- All `entity_links` for export rel set (or filter after `ListLinksByRel` per rel)
- All `plan_phases`, `plan_scopes`, `scope_deep_plans`, all `goal_plan_state` rows
- All findings/alternatives (or per-decision walk if simpler)

**Do not** use ACTIVE-only list helpers for export unless paired with a second query for non-ACTIVE rows.

### Help (S01)

Add line (wording may tighten in S02):

```text
seed export [-o <file>] Export seed JSON v1 (causal entities, links, plan tree,
                        findings/alternatives when present) to stdout or -o.
                        Sets exported_at_commit when -C root is a git repo.
```

Keep existing `seed import` line.

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| **`TestSeedExportRoundTrip`** | `cmd/trace` | Import fixture JSON (entities + links + **plan tree** + optional findings/alternatives) → `seed export` → import into **fresh** DB → assert **same ids** for goals/tasks/decisions/links and plan rows (`plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state.current_scope_id`). Subtest without plan tree optional but **primary** path includes plan tree |
| **`TestSeedExportOmitsDeniedSurfaces`** | `cmd/trace` | DB with index rows, token file on disk, capability decision, review, transition history → export → JSON must **not** contain key `transitions`; must not contain index paths, token material, capability/tool-decision payloads, or reviews. Task objects must not include `work_state` |
| **`TestSeedExportWritesExportedAtCommit`** | `cmd/trace` | **git subtest:** temp repo, `git init`, commit, `trace init` + seed data, `seed export` → `exported_at_commit` non-empty and equals `git rev-parse HEAD`. **non-git subtest:** temp dir without `.git` → key absent or empty; re-import of exported JSON still succeeds |
| Keepers | `cmd/trace` | `TestSeedImportAndWhy`; `TestSeedImportDiscoveryMentionsTask`; `TestSeedImportImpactFindings`; `TestSeedImportFromIDAliases` |

TDD: named tests red first, then export + allowlist + plan import paths (green). Round-trip test may build plan fixture via store/domain APIs in test setup if no committed fixture yet.

Implementer may add `internal/domain/seed_export_test.go` for pure build logic; **named tests above in `cmd/trace` are required** for VERIFY.

## Owns

| DF | Intent |
|----|--------|
| DF-80 | `trace seed export` + library builder + round-trip |
| DF-84 | Export **and** first-import of plan-tree keys |
| DF-85 | Export field `exported_at_commit`; import allow + ignore |

## Explicit deferrals (S02 / S03 / later)

| Item | Owner |
|------|-------|
| Idempotent entity/link/plan upsert | **S03** (DF-81/83/84 import) |
| AGENTS / CONTRIBUTING / merge docs | **S02** (DF-82/85 docs) |
| Git hook auto-export | **DF-86** deferred; not blocking VERIFY |
| MCP seed tool | Forbidden |
| Default export of transitions/reviews/capabilities/index | Forbidden |

## Light Depends — downstream scopes

### S02 commit-convention
- **Depends-on this FINAL:** CLI `seed export [-o]`; recommended path `trace/graph.json`; field name **`exported_at_commit`** (SHA evidence, not identity); export includes **plan tree** keys. S02 documents attribution (git author+SHA vs `actor`); does not re-lock JSON shapes.

### S03 idempotent-import
- **Depends-on S01 export:** round-trip produces importable JSON with plan-tree arrays. S03 owns duplicate-id / duplicate-link behavior; S01 first-import may still UNIQUE-fail on second import of same file (expected until S03).

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable thickened prompts
- [x] No product Go
- [x] Next **P17-S01-01**
