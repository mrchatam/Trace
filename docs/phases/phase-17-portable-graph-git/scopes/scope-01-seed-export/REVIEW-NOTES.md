# P17-S01-02 — seed export scope review

**Date:** 2026-08-17  
**Reviewer:** independent (fresh session)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P17-S02-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | DF-80 `trace seed export [-o]` + stdout default | PASS | `cmd/trace/seed.go` `cmdSeed` dispatches `export`; stdout when `-o` omitted; `cmd/trace/help.go` lines 48–51 |
| 2 | Round-trip ids + links + plan-tree ids | PASS | `TestSeedExportRoundTrip` green; asserts goals/tasks/decisions/links/findings + plan_phases/scopes/deep plans (ACTIVE+SUPERSEDED)/goal_plan_state |
| 3 | DF-84 all plan rows incl SUPERSEDED deep plans | PASS | Fixture includes SUPERSEDED deep plan; `ListAllScopeDeepPlans` unfiltered; round-trip asserts both deep plan ids |
| 4 | DF-85 `exported_at_commit` evidence not identity | PASS | `TestSeedExportWritesExportedAtCommit` (git HEAD match + non-git omit + re-import with SHA); import allowlist includes key; import path never persists SHA |
| 5 | Exclude denied surfaces | PASS | `TestSeedExportOmitsDeniedSurfaces` green |
| 6 | Tasks without `work_state`; canonical underscore rels | PASS | `SeedTask` omits work_state; export uses `goal_has_task`, `decision_affects_task`, etc.; omit test asserts no `work_state` |
| 7 | G19 domain export + thin CLI; no MCP seed tool | PASS | `internal/domain/seed_export.go` `BuildSeedDocument`; CLI encode only; no `trace_seed` MCP registration |
| 8 | Did not duplicate P16 S05 import | PASS | Findings/alternatives/mentions import unchanged pattern (domain APIs); additive plan import + allowlist only |
| 9 | `.gitignore` still `.trace/` only (no graph.json) | PASS | `.gitignore` lists `.trace/`; `trace/graph.json` not gitignored |
| 10 | P16 seed keepers PASS | PASS | All `TestSeedImport*` keepers green in verify run |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestSeedExport|TestSeedImport'
→ cmd/trace BUILD FAIL (tree-sitter CGO carry-forward); domain/store no matching tests

CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
→ PASS (cmd/trace 0.523s; all packages ok)
```

Named tests (cmd/trace, CGO=1):

- `TestSeedImportAndWhy` — PASS
- `TestSeedImportRelativePathAgainstC` — PASS
- `TestSeedImportFromIDAliases` — PASS
- `TestSeedImportMissingEndpointsMessage` — PASS
- `TestSeedImportDiscoveryMentionsTask` — PASS (underscore/hyphen/unknown)
- `TestSeedImportImpactFindings` — PASS
- `TestSeedExportRoundTrip` — PASS
- `TestSeedExportOmitsDeniedSurfaces` — PASS
- `TestSeedExportWritesExportedAtCommit` — PASS (git + non-git)

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| low | CGO=0 `cmd/trace` build still tree-sitter-blocked (pre-P17 carry-forward); CGO=1 verify is authoritative per board |
| low | Second import of same export may UNIQUE-fail on links/plan rows — **S03** idempotent scope (expected S01 gap) |
| nit | `TestSeedExportRoundTrip` does not assert exported JSON contains synthesized `goal_has_task` link; task `goal_id` + entity ids covered |

## Architecture compliance

- FINAL locks in `00-PLANNER.md` / `01-seed-export.md` satisfied
- G19: export logic in `internal/domain/seed_export.go`; store list-all helpers in `plan_hierarchy.go`, `helpers.go`, `impact.go`
- Import order: entities → links → findings/alternatives → plan tree → transitions
- Forbidden surfaces not exported; `exported_at_commit` omitempty when unknown

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P17-S02-00**.
