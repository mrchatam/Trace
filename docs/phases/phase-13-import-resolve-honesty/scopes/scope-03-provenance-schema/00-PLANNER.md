# P13-S03-00 — Provenance schema / enum residuals (FINAL)

## Metadata
- id: P13-S03-00
- todo_ids: [P13-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S03 for **DF-64** (enum CHECK / empty / garbage) and product residuals **DF-66**, **DF-67**. **No product Go in this row.** DF-66 = documented wontfix; DF-67 = explicit out-of-bar residual — do not silently drop.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 5 (inferred ≠ extracted)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A5
- DOGFOOD-FINDINGS DF-64/66/67; [`_bughunt/post-p12/POST-P12-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md)
- Repros: `_bughunt/post-p12/{prov,symstale}/`
- Live: `imports.provenance` mig 011; `ReplaceFileImports` / constants in `internal/store/file_graph.go`; analyzers EXTRACTED/AMBIGUOUS; `index_honesty` file-hash only; compat ceiling 11
- Depends: **P13-S02-02 APPROVE** (packet SchemaVersion `0.2`; no inventive packet bump for enum)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Prefer store honesty over inventing call-graph INFERRED.

## Live inventory (2026-08-17, post–S02 APPROVE)

| Area | Present? | Gap vs DF |
|------|----------|-----------|
| Mig **011** `imports.provenance TEXT NOT NULL DEFAULT 'EXTRACTED'` | **Present** | **No CHECK**; raw sqlite accepts `''` / `MADE_UP` (**DF-64**) |
| `ReplaceFileImports` empty→EXTRACTED | **Present** | Garbage strings **passthrough** (no reject) |
| Expand/Why/compiler copy `Import.Provenance` → `edge_provenance` | **Present** | Empty DB value + JSON `omitempty` **hides** field (looks missing) |
| Analyzer EXTRACTED; py wildcard AMBIGUOUS; no analyzer INFERRED | **Present** (P12) | No CLI/analyzer product path to set INFERRED (**DF-66**) — Law 5 via store fixture only |
| Named Law 5 fixtures (`TestImportProvenanceRoundTrip`, Expand/Why/compiler INFERRED) | **Present** | Keep green |
| `buildIndexHonesty` file-hash only | **Present** | Symbol items not in bar (**DF-67** / `symstale/`) |
| Compat embed ceiling **11**, forbid `012_*` | **Present** | Must bump to **12** when mig 012 lands |
| Packet SchemaVersion `0.2` | **Present** (S02) | **No** packet schema bump for enum work |

### Fixture map (bughunt → lock)

| DF | Fixture | Observed defect |
|----|---------|-----------------|
| DF-64 | `prov/` sqlite patch | Empty / `MADE_UP` store; empty omitempty-hides |
| DF-66 | analyzer index + help | INFERRED unreachable via live analyzer/CLI |
| DF-67 | `symstale/` | Symbol rows stay after disk mutate; only file-hash banner |

## Owns (phase lock)

| DF | Intent |
|----|--------|
| DF-64 | Harden provenance enum — reject/normalize empty + garbage; CHECK + write-path validation; no omitempty hide of empty |
| DF-66 | Documented **wontfix** for product analyzer/CLI INFERRED setter; retain Law 5 store-fixture path + thin docs |
| DF-67 | Explicit **out-of-bar** residual (file-hash honesty only); VERIFY note — no symbol-entity staleness in S03 |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | **`internal/store` primary** (`ReplaceFileImports` + read normalize); thin docs in `docs/ANALYZER_CONTRIBUTION.md`; **prefer zero** analyzer/retrieval/compiler edits |
| Migration | **`012_import_provenance_enum.sql`** — rebuild `imports` with `CHECK (provenance IN ('EXTRACTED','INFERRED','AMBIGUOUS'))`; pre-copy heal: `''`→`EXTRACTED`, unknown→`EXTRACTED` (migrate-only) |
| Compat | Bump `evals/compat` ceiling **11→12** (011 present; 012 allowed; **no 013+**) |
| Packet SchemaVersion | Keep **`0.2`** — no inventive bump for enum |
| Law 5 | INFERRED ≠ EXTRACTED; store fixture path retained |
| MCP / CLI | **No** new provenance CLI/MCP tool (G19) |
| Forbidden | Silent coerce on **write** for garbage; analyzer/call-graph INFERRED; symbol honesty impl; daemon; full-rebuild; board spawn by implementer; weakening P12 provenance named tests |
| Carry-forward | P12 provenance named tests; S01 Expand/Why import; S02 honesty residuals; honesty A/B/C+G; E/F/ablation/H/compat; p0x; x0; Gate C `dry_run:false` |

### DF-64 — provenance enum harden

| Lock | Value |
|------|-------|
| Write | Allowed: `EXTRACTED` \| `INFERRED` \| `AMBIGUOUS` only. Empty → `EXTRACTED`. Garbage → **error** (no silent coerce on write) |
| Read | `ListImportsByPath` (and any shared normalize helper): empty → `EXTRACTED` before return (defense for pre-012 DBs) |
| SQL CHECK | Mig 012 table rebuild + CHECK; index `idx_imports_file_id` restored |
| JSON | After normalize, structural hops surface a real enum — no empty `omitempty` hide |
| Named tests | Reject garbage write; empty→EXTRACTED write+read; three-enum round-trip; Expand/Why INFERRED fixture still green |

### DF-66 — INFERRED product path (documented wontfix)

| Lock | Value |
|------|-------|
| Disposition | **wontfix** product analyzer/CLI setter this scope |
| Retain | Law 5 store-fixture path (`ReplaceFileImports` + named Expand/Why/compiler tests) — must stay green |
| Docs | Thin paragraph in `docs/ANALYZER_CONTRIBUTION.md`: analyzers emit EXTRACTED/AMBIGUOUS only; INFERRED is store/fixture until future call-graph work; no `trace` CLI to set import provenance |
| Forbidden | Inventing call-graph INFERRED; new MCP/CLI provenance command |

### DF-67 — symbol-entity staleness (explicit residual)

| Lock | Value |
|------|-------|
| Disposition | **Out of honesty bar** this phase — do **not** implement symbol hashing / drop |
| Keep | File-hash `index_honesty` only (S02 universe) |
| VERIFY | S04 must record residual: symbol items can remain after disk mutate while file banner fires (`symstale/`) |
| Forbidden | Expanding honesty to whole-index or symbol content in S03 |

## Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestImportProvenanceRoundTrip` | **Keep** — three enums + empty default |
| `TestAnalyzerImportProvenanceExtracted` | **Keep** |
| `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` | **Keep** — INFERRED fixture ≠ silent EXTRACTED |
| `TestContextWhyTraceEdgeProvenance` | **Keep** |
| New (name free) | Garbage provenance on `ReplaceFileImports` → error |
| New (name free) | Empty write → stored/read EXTRACTED; Expand surfaces EXTRACTED (not omitempty-hide) |
| Optional | Direct sqlite/migrate heal path if easy without product surface |

## Effects on later scopes

- **S04 VERIFY:** Named DF-64 regressions; DF-66 wontfix docs + fixture evidence; DF-67 residual note in VERIFY-NOTES; compat ceiling 12.
- Do **not** invent packet SchemaVersion `0.3` for enum.

## Out of scope (this planner row)

- Product Go / starting **P13-S03-01**
- Symbol honesty implementation
- Analyzer INFERRED / CLI provenance setter

## Todo updates
Board: this row → `done` with FINAL Notes; next **P13-S03-01**. Status + notes only on own row after thickening 01/02/SCOPE-TODOS.

## Exit criteria
- [x] FINAL locks; thicken 01/02/SCOPE-TODOS; light S04 Depends
- [x] DF-66/67 not silently dropped
- [x] No product Go
- [x] Board next **P13-S03-01**

## Minimal todos
- [x] Inventory live schema vs DF-64/66/67
- [x] Lock FINAL defaults (mig 012 + write reject + DF-66 wontfix + DF-67 residual)
- [x] Thicken 01 + 02 + SCOPE-TODOS + light S04 Depends
- [x] Board + AGENTS/README/STATE pointers
