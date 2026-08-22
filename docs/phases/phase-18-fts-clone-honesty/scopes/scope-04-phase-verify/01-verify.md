# P18 / S04 / 01 — Phase 18 VERIFY (fts-clone-honesty product evidence)

## Metadata
- id: P18-S04-01
- todo_ids: [P18-S04-01]
- role: verify
- skills: [systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent re-prove Phase 18 named DFs (**DF-87 / DF-88 / DF-89**) + keepers + carry-forward. Write `VERIFY-NOTES.md`. **Start** DR-HANDOFF Notes (`no successor` default). Do **not** close DR-HANDOFF (**S05-02** owns close). Do **not** mark Phase 18 complete. **Two-clone recipe is not required** (no shell; no dedicated `-run`; do not implement `TestPortableGraphTwoCloneWhyContextPlan`). DF-88 is **document-only** (help + omit keepers — no clone dirs). Stale `bin/trace` / `bin/trace-mcp` is **non-fail** (**S05** still follows). **No product Go.**

**Stop if sibling `00-PLANNER.md` is still DRAFT.** (It is **FINAL** as of P18-S04-00.)

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — required **FINAL** (command set + CGO matrix SoT)
- S01 REVIEW-NOTES: [../scope-01-fts-query-sanitize/REVIEW-NOTES.md](../scope-01-fts-query-sanitize/REVIEW-NOTES.md)
- S02 REVIEW-NOTES: [../scope-02-clone-pending-honesty/REVIEW-NOTES.md](../scope-02-clone-pending-honesty/REVIEW-NOTES.md)
- S03 REVIEW-NOTES: [../scope-03-go-method-extract/REVIEW-NOTES.md](../scope-03-go-method-extract/REVIEW-NOTES.md)
- [../../DF-88-DECISION.md](../../DF-88-DECISION.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Gate C: [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/)
- [docs/TODO.md](../../../../TODO.md)
- Protocol: [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start
Follow agent-loop-protocol. Unattended: no Plan-mode switch. Do **not** trust S01–S03 Notes alone — re-run locked commands. Spawn `P18-S04-01a`/`01b` (+`01c`) on fail. Do **not** start S05. Do **not** rebuild binaries.

## Locked defaults (FINAL — P18-S04-00)

Copy CGO matrix and command set from sibling `00-PLANNER.md`. Do **not** invent a shorter substitute. Do **not** add `TestPortableGraphTwoCloneWhyContextPlan` as a dedicated fail-bar line.

| Item | Value |
|------|-------|
| Named DF-87 | `TestSanitizeFTSQueryPunctuationClass`, `TestSearchFTSSlashInQuery`, `TestTaskContextSlashTitle`, `TestTaskContextContinuesWhenSearchErrors` — keep **both** slash names |
| Named DF-88 | `TestHelpCloneTasksImportPending` — document-only |
| Named DF-89 | `TestIndexFileGoHandlerMethods` (`method:Search` + `method:SearchCursor` + `type:Memory` + `type:Notes`) |
| Keepers | `TestFTSFindsEntityTitleAndPathSymbol`, `TestIncludeWhyFailClosed`, `TestSeedExportOmitsDeniedSurfaces`, `TestHelpSeedExportPath`, `TestIndexFileGoGolden` |
| P17 seed keepers | `TestSeedExportRoundTrip`, omit, `TestSeedExportWritesExportedAtCommit` — **not** two-clone |
| Two-clone | **Not required.** No shell. No dedicated `-run`. Do not implement |
| CGO | DF-87 CGO0 authoritative + CGO1 corroboration; DF-88/89 CGO1; CGO0 cmd/trace / analyzers **non-fail** |
| Product Go | **Forbidden.** Missing named test → FAIL + spawn |
| DR-HANDOFF | **Start** `no successor`. Do **not** close |
| Next | **P18-S04-02** then **P18-S05-00** |
| Evidence | `VERIFY-NOTES.md` in this folder |
| Stale binaries | **Non-fail** — S05 |

## Locked verify commands

Copy from `00-PLANNER.md` **Locked verify command set (FINAL)**. If sandbox module cache 403s on `segmentio/encoding`, retry product bar with `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off`.

## Role work

### A. Named DF regressions (required)

| Check | Expect |
|-------|--------|
| DF-87 | Four named + two S01 keepers PASS on **CGO0** then **CGO1**. Keep both slash names |
| DF-88 | `TestHelpCloneTasksImportPending` + omit + path keepers PASS on **CGO1**. Document-only — no two-clone dirs |
| DF-89 | `TestIndexFileGoHandlerMethods` + `TestIndexFileGoGolden` PASS on **CGO1** |
| P17 seed | Round-trip + omit + `exported_at_commit` PASS. **Not** two-clone recipe |

### B. Two-clone (explicitly out)

Do **not** run the P17 shell recipe. Do **not** require a dedicated `-run`. Do **not** implement `TestPortableGraphTwoCloneWhyContextPlan`. If that test fails inside the product `./cmd/...` bar, that is a product-bar fail (spawn) — still not a reason to add a two-clone shell.

### C. Carry-forward (required)

Honesty A/B/C + G; E / F / ablation; H + compat; p0x 7/7 + x0; Gate C `dry_run:false` N=3 inspect-only; product `./cmd\|internal\|evals` (R3 graphify / R4 CGO0 FAIL OK).

### D. Residuals (record, do not fail)

Stale binaries → S05. DF-86 absent. DF-67 / DF-22 / DF-37. Harness rsync / stdio EOF. CGO0 cmd/trace / analyzers. Clone PENDING **expected** (DF-88).

### E. Evidence + handoff

Write `VERIFY-NOTES.md` with: verdict (product only, **not** phase complete); evidence table; CGO matrix; DF-88 document-only note; two-clone **not** used as fail bar; law checks; residuals; **DR-HANDOFF started = `no successor`** (not closed); stamp DR-HANDOFF VERIFY-run field; next **P18-S04-02** then **P18-S05-00**.

On FAIL: spawn `P18-S04-01a` / `01b` (+ `01c`); do not weaken bars; do not invent tests.

## Architecture / law checks

- [ ] No daemon / HTTP; G19; named tests not Notes-only
- [ ] DF-88 exclude holds; no include flags; document-only re-prove
- [ ] Two-clone shell **not** a fail bar; no dedicated two-clone `-run`
- [ ] Carry-forward honesty/E–H/compat/p0x/x0 green; Gate C `dry_run:false`
- [ ] DF-86 / stale binaries **non-fail**
- [ ] No research S05 / hosted MCP scaffold; P17 history intact
- [ ] DR-HANDOFF **started**, **not** closed; S05 rows still pending

## Preflight
1. Confirm `00-PLANNER.md` is **FINAL**.
2. Grep live `func Test*` names. If a named import is missing → FAIL + spawn; **do not** implement.
3. Run locked commands → VERIFY-NOTES + start `no successor` **or** fail→spawn.

## Exit criteria
- [ ] Locked commands run independently (or fail+spawn)
- [ ] Named DF-87 (CGO0+CGO1) + DF-88 document-only + DF-89 + keepers PASS
- [ ] Two-clone **not** required and **not** used as fail bar
- [ ] `VERIFY-NOTES.md` + DR-HANDOFF **started** (not closed)
- [ ] Board Notes; next **P18-S04-02** (then S05)
- [ ] Dry-run ≠ Gate C / F / G / ablation / H / checklist
- [ ] Phase 18 **not** complete; binaries **not** rebuilt

## Minimal todos
- [ ] Preflight: 00 FINAL; live named tests exist
- [ ] Run DF-87 CGO0 then CGO1
- [ ] Run DF-88 document-only CGO1
- [ ] Run DF-89 CGO1
- [ ] Run P17 seed keepers (not two-clone) + carry-forward
- [ ] Write VERIFY-NOTES + start DR-HANDOFF
- [ ] Board → **P18-S04-02** (or spawn)

## Out of scope
- Product Go / implementing missing tests (spawn instead)
- Two-clone shell / dedicated two-clone `-run` / implementing `TestPortableGraphTwoCloneWhyContextPlan`
- Rebuilding binaries (S05); closing DR-HANDOFF; starting **P18-S05-00**
- Scaffolding research S05 / hosted MCP; rewriting history
- Using CGO0 cmd/trace or CGO0 analyzers as a fail bar
