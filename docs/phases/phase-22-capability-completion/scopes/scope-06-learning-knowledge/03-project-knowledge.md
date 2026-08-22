# P22-S06-03 — Implement: project-specific knowledge

## Metadata
- id: P22-S06-03
- todo_ids: [P22-S06-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Learn from previous engineering decisions** (**C10**) by persisting **project-specific engineering knowledge** that accumulates over time (**C21, C26, C27**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- S06-01: `025_engineering_knowledge.sql`, `RefreshChangePatterns`, `change_patterns` populated
- Live: `decision_reconsiderations`, `reflections`, `improvements`, `RecordDecisionReconsideration`
- Seed: mirror `improvements[]` pattern in `seed_export.go` / `seed_import.go`

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| Mig **025** with both tables (S06-01) | `engineering_knowledge` rows, domain CRUD |
| Pattern refresh + list (S06-01) | `SynthesizeKnowledge`, knowledge CLI |
| Seed `improvements[]`, `decision_reconsiderations[]` | `change_patterns[]`, `engineering_knowledge[]` keys |
| Compat **25** | **026+** |

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **None new** — reuse **025** `engineering_knowledge`; compat stays **25** |
| CRUD | **`UpsertEngineeringKnowledge`**, **`GetEngineeringKnowledge`**, **`ListEngineeringKnowledge`** (topic/status filters optional) |
| body_json | Structured object — e.g. `{ "decision_id", "pattern", "summary", "source_entity_type", "source_entity_id" }`; max **8192** bytes (match change text cap) |
| evidence_ids | JSON array, max **32** items; each id validated via `GetEvidence` (mirror improvements) |
| Synthesize | **`SynthesizeKnowledge(ctx)`** — calls **`RefreshChangePatterns`** first; then upsert rows from: (1) **decision_reconsiderations** → topic `decision`; (2) **reflections** → topic `reflection`; (3) **change_patterns** where `count_positive ≥ 2` OR `count_negative ≥ 2` → topic `pattern`; (4) **improvements** → topic `improvement`; idempotent on `(source_type, source_entity_id)` composite or stable derived id |
| C10 | Reconsideration-sourced rows include `decision_id` in `body_json`; **`TestKnowledgeLinksDecision`** |
| C21/C26/C27 | Synthesize creates/updates rows over time; second synthesize updates `updated_at` without duplicating stable ids |
| Seed | Export/import **`change_patterns[]`** + **`engineering_knowledge[]`** stable ids (D-22-19) |
| CLI | `trace knowledge list [--topic <t>] [--limit N]`; `trace knowledge synthesize` JSON `{ok, created, updated}` |
| Capability | `cli:knowledge` AUTO_ALLOW |
| MCP | **No new tools** — catalog **13** |
| No LLM | Zero external model calls (W-27) |
| Checklist | C10, C21, C26, C27 **unboxed** until S06-04 |

## Requirements

1. **`internal/domain/knowledge.go`** — CRUD + synthesize orchestration.
2. **`internal/store/knowledge.go`** — upsert/list/get for `engineering_knowledge`; export helpers for patterns.
3. Extend **`seed_export.go` / `seed_import.go`** — additive keys; import upserts rows (links via JSON sufficient per S04 lock).
4. CLI + help + capability spec.
5. Named tests below.

## Touch files

- `internal/domain/knowledge.go`, `knowledge_test.go` (new)
- `internal/store/knowledge.go` (new)
- `internal/domain/seed_export.go`, `seed_import.go`, `seed_export_test.go`
- `cmd/trace/knowledge.go`, `knowledge_test.go` (new)
- `cmd/trace/root.go`, `help.go`
- `internal/domain/capability.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestUpsertEngineeringKnowledge` | CRUD roundtrip + validation (evidence ids, body size) |
| `TestSynthesizeKnowledgeFromPatterns` | C21/C26/C27 — pattern threshold ≥2 produces knowledge row |
| `TestKnowledgeLinksDecision` | C10 — reconsideration → knowledge with decision_id |
| `TestSeedExportIncludesKnowledge` | export/import round-trip patterns + knowledge |
| `TestSeedExportIncludesImprovements` | keeper (S04) — improvements key still PASS |

```bash
go test ./internal/domain/... -count=1 -run 'TestUpsertEngineeringKnowledge|TestSynthesizeKnowledgeFromPatterns|TestKnowledgeLinksDecision|TestSeedExportIncludesKnowledge|TestSeedExportIncludesImprovements'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestKnowledge'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 25 (no 026+)
```

## Exit criteria

- [ ] C10, C21, C26, C27 true (named tests)
- [ ] Compat **25** unchanged
- [ ] Seed keys live
- [ ] Checklist caps **unboxed** until S06-04
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Knowledge CRUD + synthesize
- [ ] Seed export/import keys
- [ ] CLI list + synthesize
- [ ] Tests
- [ ] Board status + notes

## Residual risks (carry to S06-04)

- **Synthesize idempotency** — second run must not duplicate rows (stable id derivation)
- **Superseded status** — v1 may never set `superseded`; reviewer accepts if only `active` written
- **Pattern-only synthesize** — knowledge from improvements must not require manual Upsert first
