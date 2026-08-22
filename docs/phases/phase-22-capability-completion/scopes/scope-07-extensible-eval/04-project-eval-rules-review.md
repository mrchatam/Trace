# P22-S07-04 — Review: project eval rules

## Metadata
- id: P22-S07-04
- todo_ids: [P22-S07-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C41** — project-specific rules load from committed file; defaults hold when absent; invariant override works.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [03-project-eval-rules.md](03-project-eval-rules.md)
- S07-02 must be **done** (registry land)
- S03: `TestArchitecturalBoundaryEdges`, `TestInvariant`

## Review checklist

### C41 — project-specific evaluation rules

- [x] Rules path locked: **`trace/eval-rules.json`** (not `.trace/`)
- [x] `TestProjectEvalRulesLoaded` PASS — valid file parsed + cache upsert
- [x] `TestMissingEvalRulesUsesBuiltins` PASS — all four built-ins when file missing
- [x] `TestProjectEvalRulesOverrideDefaultInvariant` PASS — disable `internal_must_not_import_cmd` → mechanism passes despite violation
- [x] Invalid JSON / bad version fail-closed (grep validation paths)
- [x] CLI `trace eval rules` JSON (`TestEvalRules`) PASS

### W-31 — seed pointer

- [x] Seed export includes **`eval_rules_path`** when file exists
- [x] Seed JSON does **not** embed rules body (D-22-19)
- [x] CONTRIBUTING documents path

### Holds (S07-01 / S03)

- [x] C40/C43 still hold — registry + no outcome schema change
- [x] Schema **26**; compat **26**; no 027+
- [x] Domain `CheckArchitecturalInvariants` still default implementation — eval wraps, not replaces
- [x] MCP catalog **13**

## Spawn policy

If unmet: spawn **`P22-S07-04a`** + **`P22-S07-04b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/eval/... -count=1 -run 'TestProjectEvalRules|TestMissingEvalRules|TestEvalRegistry|TestAddMechanism'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestEvalRules'
go test ./internal/domain/... -count=1 -run 'TestInvariant|TestSeedExportIncludesKnowledge'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [x] C41 closed or spawned
- [x] Confidence **high**
- [x] Checklist C41 **boxed** when closed
