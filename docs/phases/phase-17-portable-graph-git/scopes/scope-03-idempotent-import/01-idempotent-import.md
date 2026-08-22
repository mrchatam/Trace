# P17 / S03 / 01 — idempotent import

## Metadata
- id: P17-S03-01
- todo_ids: [P17-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-81**, **DF-83**, and **DF-84 import** per sibling **00-PLANNER FINAL** + [DF-84-FORWARD.md](../../DF-84-FORWARD.md). Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling **FINAL:** [00-PLANNER.md](00-PLANNER.md)
- Live: `cmd/trace/seed.go`; `internal/domain/create.go`; `internal/domain/link.go`; `internal/store/links.go`; `internal/store/plan_hierarchy.go`; `internal/store/impact.go`; `CONTRIBUTING.md`

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (from 00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| DF-81 | Re-import same seed succeeds; duplicate links **no-op**; insert-only `entity.created` / `entity.linked` |
| DF-84 | Plan-tree UUID upsert (`plan_phases.id`, `plan_scopes.id`, `scope_deep_plans.id`); `goal_plan_state` last-wins on `goal_id` |
| DF-83 | Same UUID last-import-wins (entities **and** plan rows **and** findings/alternatives); merge docs in CONTRIBUTING |
| Tasks | Upsert **must not** overwrite local `work_state` |
| Transitions | Skip when task already at target `work_state` |
| G19 | Domain **`seed_import.go`** helpers; thin `cmdSeedImport` |
| Forbidden | NDJSON; git merge driver; MCP; HTTP/daemon; S01 export changes |

## Preflight / Plan

1. Read `00-PLANNER FINAL` + live `cmdSeedImport` loop.
2. Red: add named tests in `cmd/trace/cli_test.go` (or `seed_test.go`).
3. Store: `InsertLinkOrIgnore`, `UpsertPlanPhase/Scope/DeepPlan`, `UpsertDecisionImpactFinding/Alternative`.
4. Domain: `seed_import.go` — upsert entities (insert-only events), idempotent links, plan, findings, transition skip.
5. Wire `cmdSeedImport` to domain helpers.
6. CONTRIBUTING: add union-by-id merge paragraph (locked string in 00-PLANNER).
7. Green: locked `-run` filter + full `./cmd/trace/...` + `./internal/...` as needed.

## Role work

### Files (expected touch set)

| File | Change |
|------|--------|
| `internal/store/links.go` | `InsertLinkOrIgnore` |
| `internal/store/plan_hierarchy.go` | `UpsertPlanPhase`, `UpsertPlanScope`, `UpsertScopeDeepPlan` |
| `internal/store/impact.go` | `UpsertDecisionImpactFinding`, `UpsertDecisionAlternative` |
| `internal/store/*_test.go` | Unit tests for upsert/no-op (optional but encouraged) |
| `internal/domain/seed_import.go` | Idempotent import API |
| `internal/domain/seed_import_test.go` | Optional pure tests |
| `cmd/trace/seed.go` | Replace import loops with domain seed-import |
| `cmd/trace/cli_test.go` | Named tests (required) |
| `CONTRIBUTING.md` | Merge union-by-id paragraph |

Do **not** edit `seed_export.go`, help strings, `.gitignore`, or MCP.

### Named tests (required — red then green)

| Test | Intent |
|------|--------|
| **`TestSeedImportIdempotent`** | Same file twice → exit 0; stable counts |
| **`TestSeedImportDuplicateLinksNoOp`** | Duplicate link endpoints → no error; link count stable |
| **`TestSeedImportSameIdLastWins`** | Second file overwrites body/title/plan fields for same UUID |
| **`TestSeedImportPlanTreeIdempotent`** | Plan arrays twice → same ids; no PK errors |

Keepers must stay green: `TestSeedExportRoundTrip`, `TestSeedImportAndWhy`, `TestSeedImportDiscoveryMentionsTask`, `TestSeedImportImpactFindings`, `TestHelpSeedExportPath`.

### Locked verify

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 \
  -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestSeedImport'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Minimal todos
- [ ] Red: four named `TestSeedImport*` tests
- [ ] Store upsert / InsertLinkOrIgnore
- [ ] Domain `seed_import.go` + wire `cmdSeedImport`
- [ ] CONTRIBUTING merge paragraph
- [ ] Green: locked verify + keepers

## Exit criteria
- [ ] DF-81/83/84 named tests green
- [ ] CONTRIBUTING merge union-by-id paragraph present
- [ ] No regression on S01 export / S02 help keepers
- [ ] Board Notes → **P17-S03-02**
