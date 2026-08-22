# P22-S07-01 — Implement: evaluator contract

## Metadata
- id: P22-S07-01
- todo_ids: [P22-S07-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Multiple verification mechanisms** (**C40**) via an additive **`Mechanism`** contract so evaluation can grow **without redesigning the core model** (**C43**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/domain/outcomes.go`, `internal/domain/invariants.go`, `internal/domain/coordinate.go`, `internal/store/outcomes.go`
- S03 closed — do not rewrite `CoordinateVerification`, `RecordEvaluationOutcome`, or outcome kind CHECK

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `outcome_results` kinds test/verification/evaluation | `internal/eval/` package |
| `CheckArchitecturalInvariants` (one default rule) | `Mechanism` interface / registry |
| `HasComputedEvaluation`, verification/test gates | `eval_rule_sets` table |
| Schema **25**; compat **25** | **026+ until this row** |
| MCP catalog **13** | eval MCP tools |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **`026_eval_rules.sql`** — **`eval_rule_sets`** only (planner DDL); **no ALTER** on `outcome_results` |
| Compat | Bump to **26** (`evals/compat`, store embed max, keeper tests) |
| Package layout | `internal/eval/{types,registry,run}.go`; built-ins in `internal/eval/mechs/{stored_test,stored_verification,stored_evaluation,architectural_invariant}.go` |
| Contract | `Mechanism` with `ID() string` + `Run(ctx, EvalInput) (EvalResult, error)` per planner |
| Registry | `Register`, `DefaultRegistry`, `ListMechanismIDs`, `RunAll(ctx, in, RunOptions)` |
| RunOptions | `{ MechanismIDs []string }` — empty → all registered; filter unknown ids out |
| Built-ins | Four ids locked in planner; each `Run` reads **existing store state** via `*domain.Service` — **no subprocess, no daemon** |
| Fake mech test | `TestAddMechanismWithoutSchemaChange` registers **`fake_echo`** in test-only file; proves extensibility |
| Persistence | Mechanism `Run` does **not** insert new outcome kinds; results are returned in-memory (`EvalResult` slice) |
| Rules file | **Not loaded this row** — S07-03; table exists for cache row optional no-op in S07-01 |
| G19 | Mechanism logic in `internal/eval`; domain helpers called, not duplicated |
| MCP | **No new tools** — catalog stays **13** |
| Checklist | C40, C43 **unboxed** until S07-02 |

## Requirements

1. **`026_eval_rules.sql`** — create `eval_rule_sets` per planner; no other DDL.
2. **`EvalInput` / `EvalResult`** types in `internal/eval/types.go` — JSON-friendly result fields.
3. **Registry** — thread-safe or init-only registration; built-ins register from `mechs` package `init()`.
4. **Built-in mechanisms** — map to live domain queries (see planner semantics table); `architectural_invariant` calls `CheckArchitecturalInvariants` unchanged.
5. **`RunAll`** — stable sort by mechanism id; collect results; first hard error fail-closed optional — **lock: continue other mechanisms, set `Passed=false` + Summary on individual failures; return aggregate error only when `TaskID` invalid**.
6. **Store helpers** (optional this row): `UpsertEvalRuleSet` / `GetEvalRuleSet` for cache table — stub OK if unused until S07-03.
7. **Do not** wire `RunAll` into `CoordinateVerification` or CLI verify this row — registry-only land.

## Touch files

- `internal/store/schema/026_eval_rules.sql` (new)
- `internal/eval/types.go`, `registry.go`, `run.go`, `registry_test.go` (new)
- `internal/eval/mechs/*.go` (new — four built-ins)
- `internal/store/eval_rules.go` (new — cache CRUD if needed)
- `evals/compat/compat_test.go`, `evals/compat/doc.go`
- `internal/store/*_test.go` embed-max keepers (if separate from compat)

## Named tests

| Test | Proves |
|------|--------|
| `TestEvalRegistryMultipleMechanisms` | C40 — `DefaultRegistry().ListMechanismIDs()` ≥ 4 including all built-in ids |
| `TestAddMechanismWithoutSchemaChange` | C43 — fake mech runs via `RunAll`; `ls schema/*.sql \| wc -l` still **26**; outcome kind CHECK unchanged |
| `TestOpenCreatesDBAndMigratesIdempotent` | 026 applies cleanly |
| `TestMigrationStatusReportsEmbedMax` | embed max **26** |
| `TestCompatibilitySecurityChecklist` | ceiling **26** |
| `TestRegressionDetectedVsPriorPassingTest` | keeper (S03) — ordering still stable |

```bash
go test ./internal/eval/... ./internal/store/... -count=1 -run 'TestEvalRegistry|TestAddMechanism|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
ls internal/store/schema/*.sql | wc -l  # expect 26
```

## Exit criteria

- [ ] C40, C43 true (evidence via named tests)
- [ ] Compat **26**; exactly **26** sql files
- [ ] Zero ALTER on `outcome_results` / no new outcome kind
- [ ] MCP catalog **13** unchanged
- [ ] Checklist caps **unboxed** until S07-02
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Mig 026 + optional store cache CRUD
- [ ] eval package types + registry + RunAll
- [ ] Four built-in mechanisms + tests
- [ ] Compat/embed bump to 26
- [ ] Board status + notes

## Residual risks (carry to S07-02)

- **Built-in pass semantics** — reviewer must assert each built-in `Run` matches domain gate helpers (not ad-hoc SQL)
- **Import cycle** — `internal/eval` → `domain` only; domain must not import `eval` this scope
- **CoordinateVerification unwired** — intentional; S07-03/S08 may integrate; do not sneak wiring without board row
- **eval_rule_sets empty** — OK until S07-03 loads rules file
