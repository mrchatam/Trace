# P17-S03-02 — idempotent import scope review

**Date:** 2026-08-17  
**Reviewer:** independent (fresh session)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P17-S04-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Second import same seed exits 0; no UNIQUE/PK failures | PASS | `TestSeedImportIdempotent` — round-trip fixture (entities + link + findings + plan tree incl. SUPERSEDED deep plan) → second import exit 0; link/phase/finding counts stable |
| 2 | Duplicate links no-op; count stable | PASS | `InsertLinkOrIgnore` ON CONFLICT DO NOTHING (`internal/store/links.go`); `TestSeedImportDuplicateLinksNoOp` |
| 3 | Entity upsert UUID last-wins; `entity.created` insert-only | PASS | `seedEntityExists` + `upsertEntityCreated` (`internal/domain/seed_import.go`); second import CLI JSON `"created": {}`; `TestSeedImportSameIdLastWins` |
| 4 | Task `work_state` not overwritten on upsert | PASS | `UpsertTaskFromSeed` ON CONFLICT omits `work_state` (`internal/store/entities.go:192–204`) |
| 5 | Plan tree upsert by id; `goal_plan_state` last-wins; SUPERSEDED round-trip | PASS | `UpsertPlanPhase/Scope/DeepPlan` + `UpsertGoalPlanState`; fixture includes SUPERSEDED deep plan; `TestSeedImportPlanTreeIdempotent` |
| 6 | Last-wins entity + plan fields | PASS | `TestSeedImportSameIdLastWins` — goal title/body + plan phase title/body/ord |
| 7 | Findings/alternatives upsert by id | PASS | `UpsertDecisionImpactFinding/Alternative`; round-trip fixture in `TestSeedImportIdempotent` / `TestSeedExportRoundTrip` |
| 8 | Transitions skip when already at target | PASS | `ImportSeedTransition` early return when `task.WorkState == tr.To` (`seed_import.go:619–621`) |
| 9 | CONTRIBUTING union-by-id merge docs | PASS | `CONTRIBUTING.md` § Portable graph item 5 — entities + plan arrays; no merge driver |
| 10 | `exported_at_commit` ignored for identity | PASS | Allowlisted in `seed.go`; `TestSeedExportWritesExportedAtCommit` re-import with SHA still exit 0 |
| 11 | G19 — domain `seed_import.go`; thin CLI; no MCP | PASS | `internal/domain/seed_import.go` owns import; `cmdSeedImport` delegates to `ImportSeedDocument`; no MCP seed tool |
| 12 | S01/S02 keepers green | PASS | `TestSeedExportRoundTrip`, `TestHelpSeedExportPath`, P16 seed import tests PASS |
| 13 | Scope boundary — no export/help/gitignore/hook changes | PASS | S03 scope: domain seed-import + store upserts + `cmdSeedImport` wire + named tests + CONTRIBUTING merge paragraph; export builder / help / `.gitignore` unchanged vs S01/S02 |

## Verify (independent re-run)

```text
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSeedImport|TestSeedExport'
→ PASS (cmd/trace 0.687s; all packages ok)

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestHelpSeedExportPath'
→ PASS (0.177s)

CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ PASS (full suite; cmd/trace 2.439s)

CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestSeedImportIdempotent|...'
→ cmd/trace BUILD FAIL (tree-sitter CGO carry-forward; non-fail per S01/S02 board)
```

Named tests (cmd/trace, CGO=1):

- `TestSeedImportIdempotent` — PASS
- `TestSeedImportDuplicateLinksNoOp` — PASS
- `TestSeedImportSameIdLastWins` — PASS
- `TestSeedImportPlanTreeIdempotent` — PASS
- `TestSeedExportRoundTrip` — PASS
- `TestHelpSeedExportPath` — PASS
- P16 keepers (`TestSeedImportAndWhy`, `TestSeedImportDiscoveryMentionsTask`, `TestSeedImportImpactFindings`, …) — PASS

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| low | CGO=0 `cmd/trace` build still tree-sitter-blocked (pre-P17 carry-forward); CGO=1 verify is authoritative |
| low | No dedicated named test for task `work_state` preservation on re-import — behavior locked in `UpsertTaskFromSeed` SQL (work_state excluded from ON CONFLICT UPDATE) |
| nit | `SeedImportSummary.Links` counts links processed in file, not rows inserted — cosmetic CLI summary only; DB idempotency verified by tests |
| nit | No `.git/` in workspace — scope boundary verified by live file inventory vs FINAL locks |

## Architecture compliance

- FINAL locks in `00-PLANNER.md` / `01-idempotent-import.md` satisfied
- DF-81: idempotent re-import; `InsertLinkOrIgnore`; goal_has_task no-op
- DF-83: UUID last-import-wins; insert-only `entity.created` / `entity.linked`
- DF-84 import: plan-tree upserts; SUPERSEDED deep plans in round-trip fixture
- G19: import logic in domain; CLI thin; MCP unchanged
- Forbidden items absent: no merge driver, no `.gitignore` change, no export shape change, no hook

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P17-S04-00**.
