# P17-S03-00 — idempotent import (FINAL)

## Metadata
- id: P17-S03-00
- todo_ids: [P17-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** defaults for **DF-81** (idempotent re-import), **DF-83** (conflict / merge behavior), and **DF-84 import** (plan-tree UUID upsert). Thicken sibling `01`/`02`/SCOPE-TODOS + light S04 Depends. **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- S01 FINAL: [../scope-01-seed-export/00-PLANNER.md](../scope-01-seed-export/00-PLANNER.md)
- S02 FINAL: [../scope-02-commit-convention/00-PLANNER.md](../scope-02-commit-convention/00-PLANNER.md)
- Live: `cmd/trace/seed.go` (`cmdSeedImport`); `internal/domain/create.go` (`Create*` → `Upsert*` + `appendCreated`); `internal/domain/link.go` (`InsertLink` + `appendLinked`); `internal/store/links.go` (`UNIQUE(from_type, from_id, rel, to_type, to_id)`); `internal/store/plan_hierarchy.go` (`InsertPlanPhase`/`InsertPlanScope`/`InsertScopeDeepPlan` plain INSERT; `UpsertGoalPlanState` already ON CONFLICT); `internal/store/impact.go` (`InsertDecisionImpactFinding`/`InsertDecisionAlternative` plain INSERT); `CONTRIBUTING.md` § Portable graph (git)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: **S01 FINAL + S02 FINAL** (both done). **No product Go.**

## Depends-on S01 FINAL (2026-08-17 — do not re-lock)

- **`trace seed export`** + importable JSON with plan-tree arrays + optional `exported_at_commit`
- First-import plan rows via **`Insert*`** / **`UpsertGoalPlanState`** (duplicate id **may fail** — expected until this scope)
- Named **`TestSeedExportRoundTrip`** (entity + link + plan-tree ids + findings)

S03 owns **idempotent** entity/link/plan/findings upsert (DF-81/83/84 import). Do not re-implement export.

## Depends-on S02 FINAL (2026-08-17 — do not re-lock)

- Commit path **`trace/graph.json`**; `.gitignore` **unchanged** (`.trace/` only)
- Export-before-PR in **AGENTS**; clone recipe in **CONTRIBUTING** / **README**
- Git author+SHA + **`exported_at_commit`** = **evidence**, not identity; **actor ≠ auth**
- Merge conflict on `graph.json` = **human** git resolve; **no** merge driver (S03 implements last-import-wins upsert + thickens merge docs)
- **DF-86** git-hook **deferred**

## Live inventory (confirmed 2026-08-17)

| Area | Present? | S03 gap |
|------|----------|---------|
| Entity import | `Create*` → store `Upsert*` **but always `appendCreated`** | Re-import duplicates `entity.created` events; must gate events to insert-only |
| Link import | `Link*` → `InsertLink` | Second import **UNIQUE-fails** on `(from_type, from_id, rel, to_type, to_id)` |
| `goal_has_task` | `LinkGoalTask` sets `tasks.goal_id` + link event (no `entity_links` row) | Must no-op when task already on goal |
| Plan phases/scopes/deep | `InsertPlanPhase` / `InsertPlanScope` / `InsertScopeDeepPlan` | Second import **PRIMARY KEY fails** |
| `goal_plan_state` | `UpsertGoalPlanState` ON CONFLICT | Already idempotent — keep; last-wins on `current_scope_id` |
| Findings/alternatives | `InsertDecisionImpactFinding` / `InsertDecisionAlternative` | Second import **PRIMARY KEY fails** — include in idempotent upsert (round-trip fixture has findings) |
| Transitions | `TransitionTask` strict graph | Re-import with transitions must **skip** when task already at target `work_state` |
| `exported_at_commit` | Allowlisted; ignored | **No change** — still not identity |
| Merge docs | CONTRIBUTING §5 one-liner | S03 adds **union-by-id** paragraph for entities **and plan arrays** |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF-81 | Second `seed import` of the **same file** exits **0**; no UNIQUE / PK errors on entities, links, plan rows, findings, or alternatives |
| Entity upsert | **UUID primary key last-import-wins** on all seed entity tables (`goals`, `tasks`, `decisions`, `assumptions`, `discoveries`, `plan_changes`, `claims`, `evidence`). Overwrite `title`, `body`, provenance fields, `goal_id` (tasks), `severity` (discoveries). **Preserve `created_at`** on conflict; refresh `updated_at`. **Do not overwrite `work_state`** on task upsert (seed JSON has no `work_state`; keep local machine state) |
| Events (entities) | **`entity.created` only on first insert** (row did not exist). Upsert-no-change and upsert-update: **no** extra `entity.created`. **`entity.linked` only when a new link row is inserted** (see links) |
| Links (DF-81) | Duplicate `(from_type, from_id, rel, to_type, to_id)` → **no-op success** (not UNIQUE error). Implement via store **`InsertLinkOrIgnore`** (`INSERT … ON CONFLICT(…) DO NOTHING`) or equivalent existence check. **`goal_has_task`**: no-op when `tasks.goal_id` already equals target goal |
| Plan tree (DF-84) | **`UpsertPlanPhase`**, **`UpsertPlanScope`**, **`UpsertScopeDeepPlan`** by **`id`** — last-import-wins on all seed fields (`title`, `body`, `ord`, `status`, `auto_replan_count`, `content_json`, FK refs). Preserve `created_at` on conflict. **`goal_plan_state`**: existing `UpsertGoalPlanState` — last-wins on `current_scope_id` keyed by **`goal_id`** |
| Findings/alternatives | **`UpsertDecisionImpactFinding`**, **`UpsertDecisionAlternative`** by **`id`** — last-import-wins (required so full round-trip fixture re-imports cleanly) |
| DF-83 / last-wins | Same UUID, different `body`/`title`/`content_json` in a **later import** → **last-import-wins** (upsert overwrites). Git merge on `trace/graph.json` is **human** resolve; **no** merge driver. After merge, **`trace seed import trace/graph.json`** applies union-by-id (see merge docs) |
| Merge / union-by-id | Human merges JSON arrays in git. **Entities and plan arrays** (`plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state`): keep **one object per UUID** (`id` or `goal_id` for goal_plan_state). Duplicate UUIDs in the merged file → **last array entry wins** on import. Cross-PR additive UUIDs coexist (union). Document in **CONTRIBUTING** (substance locked below) |
| `exported_at_commit` | Accept; **ignore** for identity/merge (unchanged from S01) |
| Transitions | Default export **omits** them. If seed file still contains `transitions[]`: **skip** (success, no new event) when task `work_state` already equals `to`; otherwise run existing `TransitionTask` (may still fail on illegal edge — not S03 scope to relax graph) |
| G19 home | New **`internal/domain/seed_import.go`** (name may vary) owns idempotent import helpers; **`cmd/trace/seed.go`** stays thin — replace direct `Create*` / `Insert*` / `Link*` loops with domain seed-import API. **No** MCP tool |
| Forbidden | NDJSON split; custom git merge driver; exporting reviews this scope; SHA-as-id; HTTP/daemon; changing S01 export shapes; changing `.gitignore`; DF-86 hook |

### Import order (unchanged from S01)

1. Causal entities (goals → tasks → … → evidence) — **via seed-import upsert helpers**
2. Links → findings → alternatives — **idempotent helpers**
3. `plan_phases` → `plan_scopes` → `scope_deep_plans` → `goal_plan_state` — **upsert helpers**
4. Transitions (if present) — **skip-if-already-at-target**

Parse **`exported_at_commit`** but **do not branch** on it.

### Store helpers (implementer — expected)

Add in `internal/store` (names may vary; behavior locked):

| Helper | Behavior |
|--------|----------|
| `InsertLinkOrIgnore` | ON CONFLICT on UNIQUE `(from_type, from_id, rel, to_type, to_id)` DO NOTHING; return `(inserted bool, link, err)` |
| `UpsertPlanPhase` | ON CONFLICT(`id`) DO UPDATE — last-wins fields; preserve `created_at` |
| `UpsertPlanScope` | same pattern |
| `UpsertScopeDeepPlan` | same pattern |
| `UpsertDecisionImpactFinding` | ON CONFLICT(`id`) DO UPDATE |
| `UpsertDecisionAlternative` | ON CONFLICT(`id`) DO UPDATE |

Optional: `EntityExists(type, id)` or `GetGoal` error check — for insert-only `entity.created`.

**Do not** change UNIQUE constraint on `entity_links`.

### Domain seed-import API (implementer — expected)

```text
ImportSeedGoal / ImportSeedTask / … / ImportSeedEvidence   — upsert + insert-only entity.created
ImportSeedLink(rel, from, to)                               — InsertLinkOrIgnore + insert-only entity.linked
ImportSeedFinding / ImportSeedAlternative                   — upsert by id
ImportSeedPlanPhase / ImportSeedPlanScope / ImportSeedDeepPlan / ImportSeedGoalPlanState
ImportSeedTransition                                        — skip if already at target work_state
```

`cmdSeedImport` calls these instead of `Create*` / raw store inserts.

### CONTRIBUTING — merge paragraph (S03 adds under § Portable graph)

After existing item 5 (**Merge:**), add locked substance (wording may tighten):

> When merging parallel PRs, resolve `trace/graph.json` conflicts manually in git (**no** merge driver). Combine arrays by **UUID union**: keep every distinct `id` (or `goal_id` for `goal_plan_state`). If the same UUID appears twice after your edit, **keep one object — the later entry in the array wins** on the next `trace seed import`. Re-import after merge; importer applies **last-import-wins** upsert for entities, links (duplicate no-op), and plan-tree rows.

Do **not** remove S02 evidence/hook bullets.

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| **`TestSeedImportIdempotent`** | `cmd/trace` | Import fixture JSON (reuse `TestSeedExportRoundTrip`-class payload: entities + link + findings + plan tree) → **second import same file** → exit 0; entity/link/plan/finding counts unchanged; no duplicate link rows |
| **`TestSeedImportDuplicateLinksNoOp`** | `cmd/trace` | Import once → second import → `ListLinksByRel` count unchanged; duplicate rel endpoints do not error |
| **`TestSeedImportSameIdLastWins`** | `cmd/trace` | Import file A → import file B same UUIDs different `body`/`title`/`content_json` → store reflects **B** values (entity + at least one plan row) |
| **`TestSeedImportPlanTreeIdempotent`** | `cmd/trace` | Import with `plan_phases`/`plan_scopes`/`scope_deep_plans`/`goal_plan_state` → second import → same ids; phase/scope/deep/status/`current_scope_id` preserved; no PK errors |
| Keepers | `cmd/trace` | `TestSeedExportRoundTrip`; `TestSeedImportAndWhy`; `TestSeedImportDiscoveryMentionsTask`; `TestSeedImportImpactFindings`; `TestHelpSeedExportPath` |

TDD: named tests red first, then store upserts + domain seed-import + wire `cmdSeedImport`. **`TestSeedImportIdempotent`** is the primary DF-81 gate.

Implementer may add `internal/store/*_test.go` / `internal/domain/seed_import_test.go`; **named tests above in `cmd/trace` are required**.

### Locked verify (minimum)

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 \
  -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestSeedImport'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Owns

| DF | Intent |
|----|--------|
| DF-81 | Idempotent re-import: entity UUID upsert (insert-only events) + duplicate-link no-op |
| DF-83 | Last-import-wins upsert + merge/union docs (human git resolve; no driver) |
| DF-84 | Plan-tree idempotent upsert on import |

## Explicit deferrals

| Item | Owner |
|------|-------|
| Export shapes / `exported_at_commit` export | **S01** (done) |
| Path / gitignore / export-before-PR / actor≠auth docs | **S02** (done) |
| Git hook auto-export | **DF-86** deferred |
| Two-clone VERIFY recipe | **S04** |
| MCP seed tool | Forbidden |

## Light Depends — downstream scopes

### S04 phase-verify
- **Depends-on this FINAL:** named tests `TestSeedImportIdempotent`, `TestSeedImportDuplicateLinksNoOp`, `TestSeedImportSameIdLastWins`, `TestSeedImportPlanTreeIdempotent` green; CONTRIBUTING merge union-by-id paragraph; S01 round-trip + S02 help keepers still green. S04 two-clone recipe assumes re-import and merge behavior from this scope.

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable thickened prompts
- [x] SCOPE-TODOS synced
- [x] CONTRIBUTING merge paragraph locked (implementer applies in 01)
- [x] No product Go this row
- [x] Next **P17-S03-01**
