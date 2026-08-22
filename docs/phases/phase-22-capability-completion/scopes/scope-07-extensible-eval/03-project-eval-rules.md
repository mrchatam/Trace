# P22-S07-03 — Implement: project-specific eval rules

## Metadata
- id: P22-S07-03
- todo_ids: [P22-S07-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Allow project-specific evaluation rules** (**C41**) via committed **`trace/eval-rules.json`** driving enabled mechanisms and invariant overrides. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- S07-01: `internal/eval` registry + built-ins + `eval_rule_sets` table
- S03: `CheckArchitecturalInvariants`, `RuleInternalMustNotImportCmd` — extend via rules, do not delete default domain rule

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| `internal/eval` registry + four built-ins (S07-01) | Rules loader / `trace/eval-rules.json` |
| `eval_rule_sets` table (likely empty) | `trace eval rules` CLI |
| Schema **26**; compat **26** | Seed `eval_rules_path` key |

## Locked defaults

| Item | Value |
|------|-------|
| Rules path | **`trace/eval-rules.json`** at repo root (portable; not `.trace/`) |
| Resolver | `eval.RulesPath(root string) string` → `filepath.Join(root, "trace", "eval-rules.json")` |
| Load | `eval.LoadRules(ctx, root, st)` — read file if present; parse JSON; validate `version==1`; upsert cache row `id='default'` |
| Missing file | Return default rules: all four built-in mechanism ids enabled; default invariant `internal_must_not_import_cmd` **enabled** |
| Invalid JSON / wrong version | **Fail-closed** `ErrValidation`; do not partially apply |
| Unknown mechanism id in file | Skip with logged honesty — do not fail load |
| Mechanism filter | `RunAll` / registry uses **`rules.Mechanisms`** list (order preserved in file) intersect registered ids |
| Invariant override | `architectural_invariant` mechanism reads `rules.Invariants` — when `internal_must_not_import_cmd` **disabled**, return **Passed=true** even if domain check would fail (C41 proof) |
| CLI | `trace eval rules [--root]` JSON stdout: `{path, loaded, mechanisms[], invariants[], cached_at}` |
| Capability | `cli:eval` AUTO_ALLOW (new builtin) |
| Seed (W-31) | Export **`eval_rules_path`** string when file exists on disk at export time; import stores pointer only — **no rules body in seed JSON** |
| Testdata | Use `internal/eval/testdata/eval-rules*.json` — **do not** commit repo-root policy file unless neutral example |
| CONTRIBUTING | Document `trace/eval-rules.json` path + schema (one paragraph) |
| Schema | **No new mig** — reuse 026; compat stays **26** |
| MCP | **No new tools** — catalog stays **13** |
| Checklist | C41 **unboxed** until S07-04 |

## Requirements

1. **`rules.go`** — parse/validate/default; types `RulesFile`, `InvariantRule`.
2. **Cache** — on successful load, `UpsertEvalRuleSet` with full body + `source_path` + `updated_at`.
3. **Wire `RunAll`** — accept optional `*RulesFile`; when nil, load from root (for CLI/tests).
4. **`TestProjectEvalRulesOverrideDefaultInvariant`** — fixture: internal→cmd violation exists; rules disable default invariant → `architectural_invariant` **Passed=true**.
5. **`TestMissingEvalRulesUsesBuiltins`** — no file → four mechanisms enabled.
6. **CLI + help + capability** — thin encode of load result.
7. **Seed export** — additive `eval_rules_path` in `seed_export.go` / import ignore body.

## Touch files

- `internal/eval/rules.go`, `rules_test.go`
- `internal/eval/run.go` (extend RunOptions with rules)
- `internal/eval/mechs/architectural_invariant.go` (honor disabled rules)
- `internal/eval/testdata/*.json`
- `cmd/trace/eval.go`, `eval_test.go` (new)
- `cmd/trace/root.go`, `help.go`
- `internal/domain/capability.go`
- `internal/domain/seed_export.go`, `seed_import.go` (pointer only)
- `CONTRIBUTING.md` (path paragraph)

## Named tests

| Test | Proves |
|------|--------|
| `TestProjectEvalRulesLoaded` | C41 — parse testdata file; cache row written |
| `TestProjectEvalRulesOverrideDefaultInvariant` | C41 — disabled invariant skips failure |
| `TestMissingEvalRulesUsesBuiltins` | default four mechanisms when file absent |
| `TestEvalRegistryMultipleMechanisms` | keeper (S07-01) |
| `TestCompatibilitySecurityChecklist` | ceiling still **26** |

```bash
go test ./internal/eval/... -count=1 -run 'TestProjectEvalRules|TestMissingEvalRules|TestEvalRegistry'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestEvalRules'
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesKnowledge|TestInvariant'
ls internal/store/schema/*.sql | wc -l  # still 26
```

## Exit criteria

- [ ] C41 true (evidence via named tests)
- [ ] Compat **26** unchanged
- [ ] Seed `eval_rules_path` export when file present
- [ ] Board Notes: test output summary

## Minimal todos

- [ ] Rules loader + cache upsert
- [ ] Invariant override in architectural_invariant mech
- [ ] CLI eval rules + capability
- [ ] Seed pointer + CONTRIBUTING
- [ ] Board status + notes

## Residual risks (carry to S07-04)

- **Domain rule duplication** — override must not fork `CheckArchitecturalInvariants` SQL; wrap only
- **Repo-root example file** — prefer testdata; committed example must not invent Trace product policy
- **CoordinateVerification still unwired** — rules affect `RunAll` only until future row
