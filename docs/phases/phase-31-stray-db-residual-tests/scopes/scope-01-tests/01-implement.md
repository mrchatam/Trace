# P31-S01-01 — Implement must-add tests + repro + docs

## Metadata
- id: P31-S01-01
- todo_ids: [P31-S01-01]
- role: implementer
- skills: [test-driven-development, diagnosing-bugs]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship **all three** GAPS must-add items locked by S01-00: **G1**, **G5**, **G6**. No path redesign, no silent delete, no GUI, no warn-suppress flag, no invented `trace serve` startup open.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [GAPS.md](../scope-00-inventory/GAPS.md) — SoT (must-add only)
- Live: `internal/store/stray_trace_db_test.go`, `internal/store/open.go`
- P30 VERIFY Block 3 pattern: `docs/phases/phase-30-stray-trace-db-hygiene/scopes/scope-03-verify/01-verify.md` (temp repro)

## Session start

Follow agent-loop-protocol. Do **not** invent gaps beyond the locked list below. G2 is nice-to-have — **skip** unless G1/G5/G6 are done and Notes explicitly say leftover capacity (default: skip).

## Locked must-add (frozen by P31-S01-00)

| ID | Deliverable | Home |
|----|-------------|------|
| **G1** | Dir-named root `trace.db` → quiet (no warn; open OK; live DB `.trace/trace.db`) | `internal/store/stray_trace_db_test.go` (extend; sibling only if file grows unwieldy) |
| **G5** | Durable dogfood repro script | `scripts/repro-stray-trace-db.sh` (create repo-root `scripts/` if missing) |
| **G6** | Document multi-open warn as intentional | Docs-only: `CONTRIBUTING.md` **and** `AGENTS.md` |

## Locked defaults

| Item | Value |
|------|-------|
| Scope | G1 + G5 + G6 only |
| Canonical store | `<root>/.trace/trace.db` via `store.Open` / `OpenExisting` — unchanged |
| Warn choke | `warnIfStrayRootTraceDB` in `openStore` (`open.go:85`); Stat-only regular-file (`open.go:144–149`) |
| Multi-open | Once per `openStore`; multiple CLI opens may re-emit — **document**, do not add suppress flag |
| Serve | Per-request `store.Open` only — do not invent process-startup open / G3 |
| G2 CLI harness | **Out of mandatory exit** — skip by default |
| G3 / G4 | Deferred — leave alone |
| Package floor | `go test ./internal/...` PASS |
| Delete / GUI / path redesign | **Forbidden** |

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f internal/store/stray_trace_db_test.go
test -f internal/store/open.go
test -f docs/phases/phase-31-stray-db-residual-tests/scopes/scope-00-inventory/GAPS.md
rg -n 'warnIfStrayRootTraceDB|IsRegular|traceDirName|dbFileName' internal/store/open.go
rg -n 'TestOpenWarnsWhenRootStub|TestOpenQuiet|TestOpenLeaves' internal/store/stray_trace_db_test.go
test ! -f scripts/repro-stray-trace-db.sh || echo 'script already present — update in place'
```

## Task cards (implement in order)

### G1 — Dir-stub quiet unit test

**File:** `internal/store/stray_trace_db_test.go`

**Suggested name:** `TestOpenQuietWhenRootStubIsDirectory`

**Steps:**
1. `root := t.TempDir()`
2. `os.Mkdir(filepath.Join(root, dbFileName), 0o755)` — root path named `trace.db` is a **directory**, not a regular file
3. `buf := captureWarnWriter(t)`
4. `s, err := Open(root)` — must succeed; `defer s.Close()`
5. Assert `buf.Len() == 0` (quiet — code path `!IsRegular()` at `open.go:146–147`)
6. Assert `s.DBPath() == filepath.Join(root, traceDirName, dbFileName)` and live file exists under `.trace/`
7. Assert root dir stub still exists (`os.Stat` on the mkdir path) — never removed/renamed

**Reuse:** existing `captureWarnWriter` helper; do **not** use `writeRootStub` (that writes a regular file).

**Hard check:** Existing four tests still PASS: `TestOpenWarnsWhenRootStubPresent`, `TestOpenExistingWarnsWhenRootStubPresent`, `TestOpenQuietWhenNoRootStub`, `TestOpenLeavesRootStubUntouched`.

### G5 — Durable repro script

**File:** `scripts/repro-stray-trace-db.sh` (new; create `scripts/` at repo root)

**Requirements:**
- `#!/usr/bin/env bash` + `set -euo pipefail`
- Resolve Trace binary: prefer `TRACE` env, else `$REPO/bin/trace` if present, else `go build -o "$REPO/bin/trace" ./cmd/trace` from repo root (detect repo as script’s parent-of-`scripts`)
- Temp dir (`mktemp -d`); cleanup on exit
- Flow (mirror P30 VERIFY Block 3):
  1. `"$TRACE" -C "$TMP" init`
  2. Assert **no** `$TMP/trace.db` after init; assert `$TMP/.trace/trace.db` exists
  3. `( cd "$TMP" && python3 -c "import sqlite3; sqlite3.connect('trace.db').close()" )`
  4. Capture stderr from `"$TRACE" -C "$TMP" tasks` (or any CLI that opens the store)
  5. Assert stderr contains: `project-root trace.db exists but is not the Trace store` and `.trace/trace.db` and `agents: use CLI/MCP`
  6. Assert root stub size/mtime unchanged; live store still `$TMP/.trace/trace.db`
- Exit non-zero on any FAIL; print PASS lines on success
- Do **not** delete/rename the root stub as part of “cleanup of the bug” — only remove the temp dir at end

**Self-check:** run the script once from a clean checkout; Notes must cite the path.

### G6 — Multi-open docs

**Files:** `CONTRIBUTING.md` (Portable graph / live store bullet) **and** `AGENTS.md` (SQLite / hard-boundary bullet)

**Content (short — do not essay):**
- Warn fires **once per `openStore`** when `<root>/trace.db` is a regular file
- Multiple CLI/MCP/HTTP opens each call `openStore` → warn **may re-emit** — intentional
- Live DB remains `.trace/trace.db`; root path is never the store
- **No** persistent suppress flag in this product

Do not change product code for G6. Optional: leave `open.go:19–21` comment as-is (already accurate).

## Hard checks (must stay true)

- [ ] `open.go` join remains `.trace` + `trace.db` (`traceDirName` / `dbFileName`)
- [ ] No `os.Remove` / rename of `<root>/trace.db` in product code
- [ ] Existing four stray tests + new G1 test PASS
- [ ] `/trace.db` still in `.gitignore` + `fixtures/x0/.gitignore` (do not touch other ignore scaffolds)
- [ ] No GUI / serve feature work; no suppress flag

## Role work

1. Implement G1 unit test; run `go test ./internal/store/ -count=1 -run 'TestOpen'` (or full package).
2. Add G5 script; run it once; fix until PASS.
3. Add G6 docs notes to CONTRIBUTING + AGENTS.
4. Floor: `go test ./internal/...` PASS.
5. Board Notes: list files + test name + script path + doc anchors.

## Exit criteria

- [ ] G1 shipped (named test + quiet assertions)
- [ ] G5 checked in at `scripts/repro-stray-trace-db.sh` and runs PASS
- [ ] G6 note in both `CONTRIBUTING.md` and `AGENTS.md`
- [ ] `go test ./internal/...` PASS
- [ ] No store-path change; no silent delete; no GUI; no suppress flag
- [ ] Board Notes list files + `TestOpenQuietWhenRootStubIsDirectory` (or actual name) + script path
- [ ] Next: **P31-S01-02**

## Minimal todos

- [ ] G1 dir-stub quiet unit in `stray_trace_db_test.go`
- [ ] G5 `scripts/repro-stray-trace-db.sh`
- [ ] G6 CONTRIBUTING + AGENTS multi-open note
- [ ] `go test ./internal/...` green
- [ ] Board status + Notes on **P31-S01-01** only

## Todo updates

Status + notes on **P31-S01-01** only. Do not edit upcoming review prompts.

## Next

`P31-S01-02`
