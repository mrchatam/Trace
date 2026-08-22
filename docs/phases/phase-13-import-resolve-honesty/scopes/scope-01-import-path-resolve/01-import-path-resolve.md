# P13 / S01 / 01 — Import path resolve (FINAL)

## Metadata
- id: P13-S01-01
- todo_ids: [P13-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-60** per sibling **00-PLANNER** FINAL locks. Make Expand/Why resolve subdirectory relative imports (importer-dir join + `./` normalize + common extensions / `index` candidates) so live `edge_provenance` appears. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A3 resolve-time
- DF-60 — [DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Repro trees: `experiments/_bughunt/post-p12/{prov,whyfile,rootimp}/`
- Optional: `experiments/ab-import-resolve/`
- Live: `internal/retrieval/expand.go`, `store.NormalizePath`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | `internal/retrieval` only (helper + Expand call site + tests) |
| Strategy | Resolve-time: join `Dir(importer)` + raw relative import; `NormalizePath` / Clean; try extensions + `index.*`; first indexed hit wins |
| Relative | Raw prefix `./` or `../` (slash-normalize backslashes). Bare modules: exact NormalizePath only |
| Extensions | `.js`, `.jsx`, `.mjs`, `.cjs`, `.ts`, `.tsx`, `.go`, `.py` |
| Provenance | Copy existing `Import.Provenance` onto Hit (unchanged P12) |
| Migration / analyzers | **No** mig; **no** analyzer string rewrite; **no** path-align hook as product |
| Positive controls | Root `./util.js`→`util.js`; P12 `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` |
| Forbidden | Full-rebuild; fake calls; new MCP; daemon/HTTP; board spawn; product edits outside retrieval without Notes |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Retrieval | `internal/retrieval/expand.go` | Use resolve helper instead of bare `GetFileByPath(NormalizePath(imp))` |
| Retrieval | `internal/retrieval/` helper (new file OK) | Candidate generation + lookup |
| Tests | `internal/retrieval/retrieval_test.go` (+ optional focused `*_test.go`) | Subdir + extensionless + root control; keep P12 names |
| Analyzers / store / compiler / CLI | — | **Prefer zero edits** |

## Role work
1. TDD first: store-fixture Expand/Why with analyzer-shaped `./util.js` / `./util` under `src/` → neighbor + `EXTRACTED`.
2. Wire Expand file→import loop through helper; keep unresolved → skip.
3. Prove root positive + P12 exact-path fixtures still green.
4. Run locked verify; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [ ] Resolve helper (candidates + first hit) + unit coverage of order/dedupe
- [ ] Expand call site uses importer path + raw `ImportedPath`
- [ ] Named subdir + extensionless + root control tests green
- [ ] P12 `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` green
- [ ] Carry-forward verify cmds below
- [ ] Board row P13-S01-01 Notes; next **P13-S01-02**

## Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestExpandImportEdgeProvenance` | **Keep** — exact-path fixture EXTRACTED/INFERRED |
| `TestWhySurfacesEdgeProvenance` | **Keep** — Why INFERRED fixture |
| New (name free) | Subdir importer `src/app.js` + import `./util.js` → Expand neighbor `src/util.js` + `EXTRACTED` |
| New (name free) | Subdir + `./util` (no ext) resolves to indexed `src/util.ts` (or `.js`) + `EXTRACTED` |
| New (name free) | Root importer `app.js` + `./util.js` still resolves `util.js` + `EXTRACTED` |
| Optional | Why surfaces `edge_provenance` on subdir hop (if not already covered by Expand→Why wiring) |

## Verify commands

```bash
# DF-60 + P12 provenance (CGO0)
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1

# Named focus (adjust -run to match new test names)
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance|TestExpand.*[Ii]mport|TestWhy.*[Ii]mport|TestResolve'

# Carry-forward (CGO1)
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1

# Product packages
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Do **not** treat `go test ./...` space-path FAIL (known graphify) as S01 failure. Gate C `dry_run:false` untouched. Dry-run ≠ Gate C / ≠ H / ≠ checklist.

## Exit criteria
- [ ] FINAL locks met; subdir relative imports Expand→Why with `edge_provenance`
- [ ] Root + P12 exact-path controls green
- [ ] No analyzer rewrite / no mig / no path-align product hook
- [ ] Carry-forward gates green; Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P13-S01-02**
