# P31-S01-02 — Scope review (tests)

## Metadata
- id: P31-S01-02
- todo_ids: [P31-S01-02]
- role: reviewer
- skills: [code-review-and-quality, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent fresh-session review of S01-01 against locked must-add **G1, G5, G6**. Small inline fixes OK; spawn `P31-S01-02a` / `02b` for blocker/high. **No GUI / path redesign / suppress flag.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — locked task cards
- [GAPS.md](../scope-00-inventory/GAPS.md)
- Live: `internal/store/open.go`, `internal/store/stray_trace_db_test.go`, `scripts/repro-stray-trace-db.sh`

## Session start

Follow agent-loop-protocol. **Fresh reviewer** — must not be S01-01 implementer. Re-run evidence; do not trust Notes alone.

## Locked must-add under review

| ID | Expectation |
|----|-------------|
| G1 | Unit: dir-named root `trace.db` → empty warnWriter; Open OK; `DBPath` under `.trace/`; stub dir untouched |
| G5 | Checked-in `scripts/repro-stray-trace-db.sh`; reviewer runs it; PASS (warn substrings + stub untouched + live `.trace/`) |
| G6 | CONTRIBUTING **and** AGENTS state once-per-`openStore` / multi-open may re-emit / no suppress flag |

**Out of mandatory bar:** G2 (nice-to-have). **Deferred (do not reopen):** G3 serve startup, G4 extra ignores.

## Evidence checklist

### Per must-add

- [ ] **G1:** Test present (expected name `TestOpenQuietWhenRootStubIsDirectory` or Notes-cited equivalent); uses `os.Mkdir` for root `trace.db`; asserts empty warn + live `.trace/` path
- [ ] **G5:** Script exists; `bash scripts/repro-stray-trace-db.sh` (or Notes-documented invocation) exits 0; stderr/PASS covers warn substrings from `strayRootDBWarn`
- [ ] **G6:** Grep CONTRIBUTING + AGENTS for multi-open / once-per-openStore / no suppress; wording accurate vs `open.go:19–21` + `warnIfStrayRootTraceDB`

### Cross-cutting

- [ ] Join still `.trace` + `trace.db` (cite `open.go` lines for `traceDirName`, `dbFileName`, join)
- [ ] Warn non-fatal; Stat-only regular-file; no delete/rename of root stub in product code
- [ ] Existing four stray tests still present and PASS with G1
- [ ] `go test ./internal/...` PASS (reviewer re-run)
- [ ] `/trace.db` still in `.gitignore` + `fixtures/x0/.gitignore`
- [ ] No product-scope creep (GUI, path redesign, suppress flag, invent serve startup open)
- [ ] G2 absent is OK (nice-to-have); if present, must not break floor

## Findings severity

blocker | high | medium | low | nit

- **blocker/high:** spawn `P31-S01-02a` (implement) + `P31-S01-02b` (review) immediately below this board row (full prompts per agent-loop-protocol); or trivial inline fix + re-verify
- **medium:** prefer spawn unless one-line docs/test assert fix
- Write [`REVIEW-NOTES.md`](REVIEW-NOTES.md) in this folder with findings + confidence

## Exit criteria

- [ ] `REVIEW-NOTES.md` present (findings + confidence medium or high with evidence)
- [ ] No open blocker/high without pending spawn
- [ ] Each of G1/G5/G6 closed or deferred with explicit reason (prefer closed)
- [ ] Board Notes; next **P31-S02-00** on PASS

## Todo updates

Status + notes on **P31-S01-02**; may spawn 02a/02b; may thicken upcoming S02 prompts only.

## Next

`P31-S02-00` (on PASS)
