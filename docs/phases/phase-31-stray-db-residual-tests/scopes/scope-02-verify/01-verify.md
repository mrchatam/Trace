# P31-S02-01 — VERIFY

## Metadata
- id: P31-S02-01
- todo_ids: [P31-S02-01]
- role: verify
- skills: [test-driven-development, grinding-until-pass]
- mcps: []
- verification: automated
- hooks: []

## Objective

Run the locked Phase 31 VERIFY floor after S00–S01. Capture evidence under `experiments/runs/`, write **`VERIFY-NOTES.md`** with per-block PASS/FAIL. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P31-S02-02**. **No product code.** Do **not** start Phase 32.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S02-00)
- [GAPS.md](../scope-00-inventory/GAPS.md) — must-add G1/G5/G6; defer G3/G4; nice G2
- [REVIEW-NOTES.md](../scope-01-tests/REVIEW-NOTES.md) — S01-02 PASS (high)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S02-02
- Live: `internal/store/open.go`, `internal/store/stray_trace_db_test.go`, `scripts/repro-stray-trace-db.sh`, `.gitignore`, `fixtures/x0/.gitignore`, `CONTRIBUTING.md`, `AGENTS.md`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product/test/script bodies.

## Locked defaults (FINAL — S02-00)

| Item | Value |
|------|-------|
| Precondition | P31-00 … P31-S01-02 all `done`; S01-02 **PASS high**; G1+G5+G6 shipped |
| Product Go / path redesign / silent delete / GUI / suppress | **Forbidden** (evidence + notes only) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p31-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S02-02 closes |
| Successor | **Out of scope** — S02-02 only (lean **Phase 32** / **P32-00**) |
| Warn message (expect substring) | `project-root trace.db exists but is not the Trace store` + `.trace/trace.db` + `agents: use CLI/MCP` |
| Canonical join | `filepath.Join(absRoot, ".trace", "trace.db")` — unchanged |
| Repro | Prefer shipped `scripts/repro-stray-trace-db.sh`; missing script → **FAIL** (G5 must-add) |

### Fail vs residual (locked)

**Fail VERIFY for:** focused stray store tests FAIL; `go test ./internal/...` FAIL; `scripts/repro-stray-trace-db.sh` missing or non-zero / not ALL PASS; G1 (`TestOpenQuietWhenRootStubIsDirectory`) missing; G6 multi-open note missing from CONTRIBUTING.md or AGENTS.md; store path redesigned (join no longer `.trace`+`trace.db`); silent delete/rename of root stub in `open.go`; `/trace.db` missing from `.gitignore` or `fixtures/x0/.gitignore`; init creates project-root `trace.db`; after stub, Trace opens/uses root stub instead of `.trace/trace.db`.

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| G2 CLI stderr unit absent | Nice-to-have; G5 script + store units cover CLI path |
| G3 serve “startup” warn harness | Deferred — request-scoped `store.Open` only |
| G4 extra `/trace.db` in web/experiment ignores | Out-of-scope / deferred — no additional product scaffolds |
| Agents can still create root stubs | Mitigated by warn + gitignore; agent hygiene |
| Warn once per `openStore` (multi CLI/MCP/HTTP may re-emit) | Intentional; G6 documents; no suppress flag |
| Optional delete of root stub | Future-only — not this phase |
| G5 script uses Linux `stat -c` | Acceptable for this repo’s Linux CI/dev (S01 nit) |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p31-s02-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P31-S02-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P31-S01-02 PASS high; G1+G5+G6 shipped; no path redesign"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists; metadata cites S01-02 PASS + G1/G5/G6.

### Block 1 — Focused store stray tests (required — all five)

```bash
go test ./internal/store/ \
  -run 'TestOpenWarnsWhenRootStubPresent|TestOpenExistingWarnsWhenRootStubPresent|TestOpenQuietWhenNoRootStub|TestOpenLeavesRootStubUntouched|TestOpenQuietWhenRootStubIsDirectory' \
  -count=1 \
  2>&1 | tee "$EVID/01-store-stray-tests.txt"
```

**Pass:** exit 0; all five cases run and PASS (four Phase 30 + G1 dir-stub quiet).

### Block 2 — Package bar (required)

```bash
go test ./internal/... -count=1 2>&1 | tee "$EVID/02-internal-test.txt"
```

**Pass:** exit 0.

### Block 3 — Repro script (required — G5 shipped)

```bash
test -x scripts/repro-stray-trace-db.sh \
  || { echo "FAIL: scripts/repro-stray-trace-db.sh missing or not executable" | tee "$EVID/03-repro.txt"; exit 1; }
bash scripts/repro-stray-trace-db.sh 2>&1 | tee "$EVID/03-repro.txt"
```

**Pass:** script exit 0; output includes `ALL PASS` (or equivalent per-script PASS lines); confirms warn substrings / stub untouched / live `.trace/trace.db` as the script asserts.

**Fallback (only if script truly absent — still FAIL overall for G5):** do **not** invent a green VERIFY. Record FAIL; optional diagnostic one-shot matching P30 Block 3 shape may be captured under `$EVID/03-fallback-oneshot.txt` for DR-HANDOFF, but overall remains **FAIL** until G5 is restored via repair spawn.

### Block 4 — Docs / gitignore / join / G6 spot-check (required)

```bash
{
  echo "=== gitignore ==="
  grep -n '/trace.db' .gitignore fixtures/x0/.gitignore
  echo "=== G6 multi-open (CONTRIBUTING + AGENTS) ==="
  grep -n 'openStore\|once per\|suppress' CONTRIBUTING.md AGENTS.md
  echo "=== open.go join + warn choke ==="
  grep -n 'traceDirName\|dbFileName\|warnIfStrayRootTraceDB\|strayRootDBWarn\|filepath.Join\|IsRegular\|Remove\|Rename' internal/store/open.go
  echo "=== no Remove/Rename of stub in open.go (expect empty or non-stub cleanup only) ==="
  grep -n 'os\.Remove\|os\.Rename' internal/store/open.go || echo "PASS: no Remove/Rename in open.go"
} 2>&1 | tee "$EVID/04-docs-gitignore-join.txt"
```

**Pass:** `/trace.db` in both gitignores; CONTRIBUTING + AGENTS mention once-per-`openStore` / multi-open re-emit / no suppress; `warnIfStrayRootTraceDB` + `traceDirName`/`.trace` + `dbFileName`/`trace.db` present; Stat/`IsRegular` gate intact; no silent delete/rename of root stub.

### Block 5 — Residuals list (required in notes)

Record non-blocking residuals from the locked table (G2/G3/G4, agent stubs, multi-open, optional delete, Linux `stat -c` nit). Do **not** fail solely for them. Overall PASS only if blocks 0–4 green.

## Do not

- Change `open.go`, tests, `scripts/repro-stray-trace-db.sh`, gitignore, or product docs in this row
- Close DR-HANDOFF (that is **P31-S02-02**)
- Claim/scaffold successor beyond citing lean Phase 32 (S02-02 only)
- Start implementing Phase 32
- Treat root stub as a second store
- Silent-delete the stub “to clean evidence”
- Fail solely because G2 is absent or G3/G4 deferred

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-31-stray-db-residual-tests/scopes/scope-02-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P31-S02-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p31-s02-01-verify/evidence/
**Precondition:** P31-S01-02 PASS high (G1+G5+G6)

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| 1 Store stray tests (5) | PASS/FAIL | names + exit |
| 2 go test ./internal/... | PASS/FAIL | |
| 3 Repro script G5 | PASS/FAIL | scripts/repro-stray-trace-db.sh |
| 4 Docs/gitignore/join/G6 | PASS/FAIL | |
| 5 Residuals | listed | |

## Residuals (non-blocking)
- G2 …
- G3 …
- G4 …
- multi-open once-per-openStore …
- agents can still create stubs …
- optional delete future-only …

## Failures (if any)
- …

## DR-HANDOFF
remains OPEN — close owner **P31-S02-02**

## Next
P31-S02-02
```

## Todo updates

Status + notes on **P31-S02-01** only. Do not mark S02-02 done; do not edit S02-02 beyond what Notes need to cite.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0–5
- [ ] Evidence dir present under `experiments/runs/…-p31-s02-01-verify/evidence/`
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P31-S02-02**

## Next

`P31-S02-02`
