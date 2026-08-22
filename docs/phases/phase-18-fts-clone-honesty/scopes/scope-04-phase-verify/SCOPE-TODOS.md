# S04 — Phase VERIFY — scope todos

**Depends-on:** P18-S01-02 / S02-02 / S03-02 **APPROVE**. Product evidence for named DF-87/88/89 + keepers + carry-forward. **Does not** close the phase — **P18-S05-00** follows. **P18-S05-02** closes DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **FINAL** (P18-S04-00) |
| 2 | 01-verify | verifier | **PASS** (P18-S04-01) — VERIFY-NOTES; DR-HANDOFF started `no successor` |
| 3 | 02-scope-review | review (product) | **APPROVE** high (P18-S04-02) — next **P18-S05-00**; does **not** own DR-HANDOFF close |

## Phase locks (00 FINAL)

### Named tests (fail bar)

| DF | Named | Keepers |
|----|-------|---------|
| **DF-87** | `TestSanitizeFTSQueryPunctuationClass` · `TestSearchFTSSlashInQuery` · `TestTaskContextSlashTitle` · `TestTaskContextContinuesWhenSearchErrors` | `TestFTSFindsEntityTitleAndPathSymbol` · `TestIncludeWhyFailClosed` |
| **DF-88** (document-only) | `TestHelpCloneTasksImportPending` | `TestSeedExportOmitsDeniedSurfaces` · `TestHelpSeedExportPath` |
| **DF-89** | `TestIndexFileGoHandlerMethods` | `TestIndexFileGoGolden` |

P17 seed keepers (not two-clone): `TestSeedExportRoundTrip` · `TestSeedExportOmitsDeniedSurfaces` · `TestSeedExportWritesExportedAtCommit`.

### CGO matrix

- DF-87 store/compiler/retrieval: **CGO0 authoritative** + **CGO1 corroboration**
- DF-88 `cmd/trace` + DF-89 analyzers: **CGO1 required**
- CGO0 `cmd/trace` / analyzers: **non-fail**
- Product bar: CGO1 + `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off`

### Two-clone

**Not required.** Do not run P17 shell recipe. Do not add a dedicated `-run`. Do not implement `TestPortableGraphTwoCloneWhyContextPlan`.

### DF-88

Document-only re-prove (help named + omit/path keepers). Keep P17 exclude. Clone PENDING expected. No clone-dir hunt.

### Handoff / S05

- Carry-forward honesty/E–H/compat/p0x/x0 + product `./cmd\|internal\|evals`
- Residuals non-fail: DF-86, DF-67, harness rsync/stdio EOF, DF-22/37, stale binaries (**S05**)
- DR-HANDOFF default **`no successor`** — S04-01 starts Notes; S04-02 confirms not closed; **S05-02 owns close**
- Next after S04-02: **P18-S05-00**

## Reminders
- Do **not** auto-scaffold Phase 19 / hosted MCP / research S05
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Phase 17 historical handoff stays intact
- Stale `bin/trace` / `bin/trace-mcp` is **S05**, not a VERIFY fail
- S04-01 does **not** close the phase
- Missing named test → FAIL + spawn (no product Go)
