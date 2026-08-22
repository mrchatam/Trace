# P22-S07-00 — Planner: extensible evaluation

## Metadata
- id: P22-S07-00
- todo_ids: [P22-S07-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S07. Owned: **C40, C41, C42-library, C43**. **No product Go.**

## Live inventory (2026-08-18, post-S06)

| Surface | Live state |
|---------|------------|
| Schema max | **025** (`025_engineering_knowledge.sql`; 25 embed sql files) |
| Compat ceiling | **25** (`evals/compat/compat_test.go` — no 026+) |
| Outcomes | `outcome_results` kinds **`test` \| `verification` \| `evaluation`** only (mig 018); `RecordEvaluationOutcome` + `CompareScoresToBaseline` in `internal/domain/outcomes.go` |
| Verification cycle | S03: `CoordinateVerification` test→verify→eval; `BuildPolicyInputs` cycle flags; `HasComputedEvaluation`; `CompareIterationOutcomes` (last-two-of-kind) |
| Invariants | S03: `CheckArchitecturalInvariants` — **one** default rule `internal_must_not_import_cmd` in `internal/domain/invariants.go`; advisory JSON only (no auto verification row); CLI `trace verify invariants` |
| Context / C42-surface | S05: `buildEvaluations` via `ListOutcomeResultsByTaskKind(evaluation)` — cap **8**; **no `mechanism_id`** on packet items |
| Eval package | **Absent** — no `internal/eval/`; no registry, no `Mechanism` interface |
| Project rules | **Absent** — no `trace/eval-rules.json`; no `eval_rule_sets` table |
| Seed | exports `outcome_results`, `change_patterns`, `engineering_knowledge`; **no** `eval_rules_path` key yet (D-22-19 W-31) |
| MCP catalog | **13** tools — no eval-specific MCP (G19: library + CLI this scope) |
| Patterns / knowledge | S06 closed — do not reopen mig 025 tables except as **read dependencies** for eval context |

S01–S06 closed — do not reopen schema 025 tables, FTS sync, patterns/knowledge CLI, or evidence context caps in S07 prompts except as **read dependencies**.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-12 (evaluator contract + rules file), D-22-18 (not a bake-off runner), D-22-19 (seed eval-rules pointer), D-22-21 (mig 026)
- [WORK-MAP.md](../../WORK-MAP.md) W-28…W-31
- Coverage: [README.md](../../README.md) C40–C43 rows (C42 split: S05 surface + S07 library)
- S03 overlap: `CheckArchitecturalInvariants`, `RecordEvaluationOutcome`, `CompareIterationOutcomes` — S07 **wraps** via mechanisms; do not redesign outcomes/changes tables
- S05 overlap: context `evaluations[]` — S07-05 **`ListResults`** is library SoT; compiler may optionally delegate later; **do not remove** S05 surface this scope

## FINAL locked defaults

| Item | Value |
|------|-------|
| Package | `internal/eval` |
| Contract | `type Mechanism interface { ID() string; Run(ctx context.Context, in EvalInput) (EvalResult, error) }` |
| EvalInput | `{ TaskID string; Service *domain.Service }` — mechanisms read store **via domain** (G19); no MCP SQL |
| EvalResult | `{ MechanismID, Passed bool, Summary string, DetailsJSON string, RecordedAt string }` — **no new `outcome_results` kind** |
| Registry | `Register(Mechanism)`, `DefaultRegistry()`, `RunAll(ctx, in, opts)` — runs enabled mechanisms in stable ID order |
| Built-ins | **`stored_test`**, **`stored_verification`**, **`stored_evaluation`**, **`architectural_invariant`** — register in `init()` under `internal/eval/mechs/` |
| Built-in semantics | **stored_test**: pass when `HasTestOutcomeSinceLatestChange`; **stored_verification**: pass when verification gate satisfied (`CheckVerificationGate` or stored verification row); **stored_evaluation**: pass when `HasComputedEvaluation`; **architectural_invariant**: delegate `CheckArchitecturalInvariants`, pass when `Passed` |
| Mig | **`026_eval_rules.sql`** — `eval_rule_sets` cache of committed rules file |
| Compat | **26** after **S07-01** (S07-03/S07-05 stay **26**; forbid **027+** entire S07) |
| Project rules | Committed **`trace/eval-rules.json`** (portable; not `.trace/`). Schema: `{ "version": 1, "mechanisms": ["id", ...], "invariants": [{"id":"...", "enabled": true}] }`. **Missing file → all four built-ins enabled**; unknown mechanism id in file → skip with honesty (do not fail entire load) |
| Rules cache | Single row `id='default'`, `source_path`, `body_json`, `updated_at` — upsert on successful parse of `trace/eval-rules.json` |
| Core model | Do **not** redesign `outcome_results` / `changes` / `CoordinateVerification` flow in S07-01; mechanisms are **additive wrappers** |
| Invariant override (C41) | Rules file `invariants[].enabled=false` disables that rule id for **`architectural_invariant`** mechanism; default rule id **`internal_must_not_import_cmd`** (constant in domain) |
| C42 library | **`eval.ListResults(ctx, svc, taskID)`** — stable API returning evaluation rows + mechanism-shaped rows with **`mechanism_id`**; maps existing `outcome_results` rows to mechanism ids; includes `comparison_json` for evaluations; **no new DB column on outcome_results** |
| CLI | S07-03: `trace eval rules` (show parsed rules + enabled mechanisms); S07-05: optional `trace eval results --task <id>` thin encode of `ListResults` |
| MCP | **No new tools** — catalog stays **13** |
| Seed (W-31) | Additive **`eval_rules_path`** string in seed export (default `"trace/eval-rules.json"` when file exists); **do not** embed rules body in seed JSON (git-canonical file) |
| Checklist | Implementers **unbox** owned caps; reviewers **box** after review rows |

### Mig 026 DDL (locked shape)

```sql
CREATE TABLE eval_rule_sets (
    id TEXT PRIMARY KEY,
    source_path TEXT NOT NULL DEFAULT '',
    body_json TEXT NOT NULL DEFAULT '{}',
    updated_at TEXT NOT NULL
);
```

### trace/eval-rules.json (locked shape)

```json
{
  "version": 1,
  "mechanisms": ["stored_test", "stored_verification", "stored_evaluation", "architectural_invariant"],
  "invariants": [
    {"id": "internal_must_not_import_cmd", "enabled": true}
  ]
}
```

## Named tests

| Test | Row |
|------|-----|
| `TestEvalRegistryMultipleMechanisms` | S07-01 |
| `TestAddMechanismWithoutSchemaChange` | S07-01 |
| `TestProjectEvalRulesLoaded` | S07-03 |
| `TestProjectEvalRulesOverrideDefaultInvariant` | S07-03 |
| `TestMissingEvalRulesUsesBuiltins` | S07-03 |
| `TestListEvaluationResultsForFutureAgents` | S07-05 |
| `TestEvalResultsIncludeMechanismID` | S07-05 |

## Exit criteria

- [x] 01–06 thickened
- [x] Mig 026 locked for S07-01
- [x] No product Go

## Next

**P22-S07-01**
