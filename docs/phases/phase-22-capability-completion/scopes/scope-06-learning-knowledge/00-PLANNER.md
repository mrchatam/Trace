# P22-S06-00 — Planner: learning + project knowledge

## Metadata
- id: P22-S06-00
- todo_ids: [P22-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S06. Owned: **C10, C19, C20, C21, C22, C23, C24, C26, C27**. Deterministic aggregation — **no ML**. **No product Go.**

## Live inventory (2026-08-18, post-S05)

| Surface | Live state |
|---------|------------|
| Schema max | **024** (`024_impact_compare.sql`; 24 embed sql files) |
| Compat ceiling | **24** (`evals/compat/compat_test.go` — no 025+) |
| Changes | `changes` + `change_paths` (path, symbol_id); VCS promote (S02); `ListChangesRecent`, CLI `list\|show\|capture\|compare` |
| Effects | `effects` per `(change_id, dimension)` — `comparison` ∈ `{supported, partially_supported, contradicted}` |
| Outcomes | `outcome_results` kinds test/verification/evaluation; test pass/fail rows (S03) |
| Regressions | `regressions` + `regression_associated_change` links (S04) |
| Improvements | `improvements` table + `RecordImprovement` + seed export (S04) |
| Cognition | `decisions`, `decision_reconsiderations`, `reflections` — seed round-trip live |
| Evidence queries | S05: `ListWorkedApproaches`, `trace outcomes worked`, `trace tests verifying`, etc. |
| Context packet | S05: `evaluations`, `reflections`, `planning_evidence` cap **8**; schema `"0.2"` |
| Loop next | S05: `planning_evidence` section mirrors packet evidence |
| Patterns / knowledge | **Absent** — no `change_patterns`, no `engineering_knowledge`, no `patterns.go` / `knowledge.go` |
| CLI | **No** `trace patterns`, `trace knowledge`; `trace changes` has no `similar` subcommand |
| MCP catalog | **13** tools — no pattern/knowledge tools (G19: CLI-first this scope) |
| Seed export | exports `improvements`, `decision_reconsiderations`, `reflections`; **no** `change_patterns` / `engineering_knowledge` keys |

S01–S05 closed — do not reopen schema 024 tables, FTS sync, or evidence CLI in S06 prompts except as **read dependencies**.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-11 (deterministic, no ML), D-22-19 (seed keys), D-22-21 (mig 025)
- [WORK-MAP.md](../../WORK-MAP.md) W-23…W-27
- Coverage: [README.md](../../README.md) C10, C19–C24, C26, C27 rows
- S05 overlap: C33 `ListWorkedApproaches` / `trace outcomes worked` — S06 **extends** via knowledge + context (do not remove S05)

## FINAL locked defaults

| Item | Value |
|------|-------|
| Mig | **`025_engineering_knowledge.sql`** — **`change_patterns`** + **`engineering_knowledge`** (both tables in one mig; **S06-01 owns mig + patterns**; S06-03 fills knowledge rows) |
| Compat | **25** after **S06-01** (S06-03/S06-05 stay **25**; forbid **026+** entire S06) |
| change_kind | **`InferChangeKind(change)`** — primary path prefix = first segment of lexicographically smallest `change_paths.path`; format **`seg:<segment>`** (e.g. `seg:internal`, `seg:cmd`); no paths → **`seg:unknown`** |
| outcome_kind | **`ClassifyChangeOutcome(change)`** — one bucket per change, priority: **`regression`** > **`effect_contradicted`** > **`improvement`** > **`effect_supported`** > **`test_fail`** > **`test_pass`** > **`neutral`** |
| Pattern refresh | **`RefreshChangePatterns(ctx)`** — full deterministic rebuild (DELETE+INSERT from changes+effects+outcomes+regressions+improvements); **no ML**; not on every index — explicit call + **`SynthesizeKnowledge`** prelude |
| Pattern counts | **`count_positive`**: outcome ∈ {improvement, effect_supported, test_pass}; **`count_negative`**: outcome ∈ {regression, effect_contradicted, test_fail}; **`last_seen`** = max change `created_at` in bucket |
| Similar query | **`QuerySimilarChanges(ctx, opts)`** — filter by **`PathPrefix`** (paths LIKE `prefix%`) **or** **`ChangeKind`**; return prior changes + compact effects/outcome summary; default limit **32**, cap **64**; paths-only (no blobs) |
| Tend | **`ListTendencies(ctx)`** — from `change_patterns` where `count_positive ≥ 2` → `direction=improve`; `count_negative ≥ 2` → `direction=damage`; threshold **2** (both may apply to different kinds) |
| Knowledge table | `engineering_knowledge`: `id`, `title`, `body_json` (structured JSON object), `topic`, `evidence_ids_json`, `confidence`, `status` ∈ {`active`,`superseded`}, `source_type`, `created_at`, `updated_at` |
| Synthesize | **`SynthesizeKnowledge(ctx)`** — upsert rows from: **decision_reconsiderations**, **reflections** (summary), **patterns** (kinds with either count ≥ 2), **improvements**; provenance via `evidence_ids_json` + `source_type`; **no LLM** |
| C10 link | Knowledge from reconsiderations must reference **`decision_id`** in `body_json` or validated evidence link — **`TestKnowledgeLinksDecision`** |
| Seed (D-22-19) | Additive keys **`change_patterns[]`**, **`engineering_knowledge[]`** export/import stable ids (S06-03) |
| Context S06-05 | Packet + loop next: **`tendencies`** (cap **8**), **`successful_approaches`** (cap **8** — reuse/extend `ListWorkedApproaches` + active knowledge rows), optional loop **`similar_changes`** (cap **8**) when task has recent change paths |
| CLI | `trace patterns refresh\|list`; `trace changes similar --path <prefix> \| --kind <seg:…>`; `trace knowledge list\|synthesize\|tendencies` |
| MCP | **No new tools** — catalog stays **13**; `trace_context` inherits compiler fields |
| G19 | Aggregation in domain/store; CLI thin encode |
| Checklist | Implementers **unbox** owned caps; reviewers **box** after review rows |

### Mig 025 DDL (locked shape)

```sql
CREATE TABLE change_patterns (
    change_kind TEXT NOT NULL,
    outcome_kind TEXT NOT NULL,
    count_positive INTEGER NOT NULL DEFAULT 0,
    count_negative INTEGER NOT NULL DEFAULT 0,
    last_seen TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (change_kind, outcome_kind)
);

CREATE TABLE engineering_knowledge (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    body_json TEXT NOT NULL DEFAULT '{}',
    topic TEXT NOT NULL DEFAULT '',
    evidence_ids_json TEXT NOT NULL DEFAULT '[]',
    confidence REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'superseded')),
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX idx_engineering_knowledge_topic ON engineering_knowledge(topic);
CREATE INDEX idx_engineering_knowledge_status ON engineering_knowledge(status);
```

## Named tests

| Test | Row |
|------|-----|
| `TestPatternCountsFromChangesAndOutcomes` | S06-01 |
| `TestQuerySimilarChanges` | S06-01 |
| `TestUpsertEngineeringKnowledge` | S06-03 |
| `TestSynthesizeKnowledgeFromPatterns` | S06-03 |
| `TestKnowledgeLinksDecision` | S06-03 |
| `TestSeedExportIncludesKnowledge` | S06-03 |
| `TestTendHelpHurtInContext` | S06-05 |
| `TestSuccessfulApproachesSurfaced` | S06-05 |
| `TestLoopNextIncludesEvidenceForDecisions` | S06-05 |

## Exit criteria

- [x] 01–06 thickened
- [x] Mig 025 locked for S06-01
- [x] No product Go

## Next

**P22-S06-01**
