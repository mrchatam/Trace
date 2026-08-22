# P17-S02-02 — commit convention scope review

**Date:** 2026-08-17  
**Reviewer:** independent (fresh session)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P17-S03-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | DF-82 path lock **`trace/graph.json`** | PASS | `CONTRIBUTING.md` § Portable graph (git) item 1; `README.md` clone recipe; `cmd/trace/help.go` lines 48–53; `TestHelpSeedExportPath` |
| 2 | `.gitignore` unchanged — **`.trace/` only** | PASS | `.gitignore` lists `.trace/` only; no `trace/` or `graph.json` ignore |
| 3 | AGENTS export-before-PR (plan tree included) | PASS | `AGENTS.md` Hard boundaries bullet: `trace seed export -o trace/graph.json` + evidence-not-identity |
| 4 | Clone recipe: init → import → index → why/context/plan | PASS | `README.md` fenced block (init/import/index/plan show/why/context); `CONTRIBUTING.md` item 3 prose |
| 5 | **`exported_at_commit` / git SHA = evidence, not identity** | PASS | `CONTRIBUTING.md` item 4; help `Sets exported_at_commit (git SHA evidence, not identity)`; `TestHelpSeedExportPath` |
| 6 | Git **author + SHA** documented as evidence | PASS | `CONTRIBUTING.md` item 4: git author + commit SHA + `exported_at_commit` |
| 7 | **`transition.actor` / review actor / `as_operator` ≠ auth** | PASS | `CONTRIBUTING.md` item 4 (DF-44 cross-ref); help lines 34–35 `flag≠identity` + `Actor string ≠ auth`; `TestAsOperatorFlagIdentityDocs` |
| 8 | Merge: human resolve on `graph.json`; **no merge driver** | PASS | `CONTRIBUTING.md` item 5 |
| 9 | UUID last-import-wins pointer to **S03** | PASS | `CONTRIBUTING.md` item 5 → **Phase 17 S03** |
| 10 | DF-86 hook **not** implemented; docs optional later | PASS | No `install git-hook` in `*.go`; `CONTRIBUTING.md` item 6 optional/deferred |
| 11 | DF-28 help handoff still present | PASS | `TestHelpHandoffSoT` green; help Handoff block unchanged |
| 12 | DF-44 help still present | PASS | `TestAsOperatorFlagIdentityDocs` green |
| 13 | **No** S01 export/seed code changes | PASS | Scope limited to `AGENTS.md`, `CONTRIBUTING.md`, `README.md`, `cmd/trace/help.go`, `cmd/trace/cli_test.go` (`TestHelpSeedExportPath`); `seed.go` / `seed_export.go` unchanged vs S01 (no `graph.json` in export builder) |
| 14 | No hosted MCP / no `trace-mcp` internet pointer | PASS | Changed docs/help do not add internet MCP pointers; help Build note remains local stdio only |

## Verify (independent re-run)

```text
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSeedExportPath|TestSeedExport|TestHelpHandoffSoT|TestAsOperatorFlagIdentityDocs'
→ PASS (0.176s)

CGO_ENABLED=0 go test ./cmd/trace/... -count=1 -run 'TestHelpSeedExportPath|TestSeedExport|TestHelpHandoffSoT|TestAsOperatorFlagIdentityDocs'
→ cmd/trace BUILD FAIL (tree-sitter CGO carry-forward; non-fail per S01 board)

CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedExport|TestSeedImport'
→ PASS (S01 keepers unchanged)
```

Named tests (cmd/trace, CGO=1):

- `TestHelpSeedExportPath` — PASS
- `TestHelpHandoffSoT` — PASS
- `TestAsOperatorFlagIdentityDocs` — PASS
- `TestSeedExport*` (S01 carry-forward) — PASS

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| low | CGO=0 `cmd/trace` build still tree-sitter-blocked (pre-P17 carry-forward); CGO=1 verify is authoritative per board |
| low | Second import of same export may UNIQUE-fail — **S03** idempotent upsert (expected S02 gap per 02-scope-review.md) |
| nit | No `.git/` in workspace — checklist #13 verified by live file inventory + S01 implementer scope, not `git diff` |

## Architecture compliance

- FINAL locks in `00-PLANNER.md` / `01-commit-convention.md` satisfied
- DF-82: path convention + gitignore unchanged + export-before-PR agent docs
- DF-85 docs: git author+SHA + `exported_at_commit` = evidence; actor ≠ identity in contributor docs
- Forbidden items absent: no hook implementation, no seed/export code changes, no `trace/` gitignore

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P17-S03-00**.
