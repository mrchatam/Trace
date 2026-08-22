# P30-S03-01 — VERIFY

## Metadata
- id: P30-S03-01
- todo_ids: [P30-S03-01]
- role: implementer
- skills: [systematic-debugging, grinding-until-pass]
- verification: automated
- hooks: []

## Objective

Run the Phase 30 VERIFY floor after S00–S02. Capture evidence, write **`VERIFY-NOTES.md`** with per-block PASS/FAIL. Keep **`DR-HANDOFF.md` OPEN** — S03-02 owns successor close. **No product code.** Do **not** start a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor + successor lean (FINAL)
- [PLAN.md](../scope-01-plan/PLAN.md) — T1–T4 SoT
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) — agent hygiene / no store-path change
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- Live: `internal/store/open.go`, `internal/store/stray_trace_db_test.go`, `.gitignore`, `fixtures/x0/.gitignore`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or implement product hygiene.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P30-00 … P30-S02-02 all `done`; PLAN.md T1–T4 SoT; S00 **agent hygiene** / INTAKE confirmed / **no store-path change** |
| Product Go / path redesign | **Forbidden** (evidence + notes only) |
| Binary | Prefer rebuild `bin/trace` from repo HEAD before CLI repro |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p30-s03-01-verify/` (tee under `evidence/` subdir) |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes |
| Successor | **Out of scope** — S03-02 only (default lean **no successor**) |
| Warn message (expect substring) | `project-root trace.db exists but is not the Trace store` + `.trace/trace.db` + `agents: use CLI/MCP` |
| Canonical join | `filepath.Join(absRoot, ".trace", "trace.db")` — unchanged |

### Fail vs residual (locked)

**Fail VERIFY for:** focused store warn tests FAIL; `go test ./internal/...` FAIL; `trace init` creates project-root `trace.db`; after stub, Trace opens/uses root stub instead of `.trace/trace.db`; warn missing on open when regular-file stub present; silent delete/rename of root stub; `open.go` join no longer `.trace`+`trace.db`; `/trace.db` missing from `.gitignore` or `fixtures/x0/.gitignore`.

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| Agents can still `sqlite3.connect('trace.db')` and create stubs | Mitigated by warn + gitignore; cannot prevent all agent mistakes |
| Optional documented delete of root stub | **Future-only** — not Phase 30 |
| Warn once per `openStore` in long-lived `serve` | Acceptable; no persistent suppress flag |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Today's date for dir name: use `date +%Y-%m-%d` (expected **2026-08-21** if run same day as planner).

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p30-s03-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P30-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P30-S02-02 PASS; PLAN T1-T4; no store-path change"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists; metadata cites S02 PASS + no store-path change.

### Block 1 — Focused store warn tests (required)

```bash
go test ./internal/store/ \
  -run 'TestOpenWarnsWhenRootStubPresent|TestOpenExistingWarnsWhenRootStubPresent|TestOpenQuietWhenNoRootStub|TestOpenLeavesRootStubUntouched' \
  -count=1 \
  2>&1 | tee "$EVID/01-store-stray-tests.txt"
```

**Pass:** exit 0; all four cases run and PASS.

### Block 2 — Package bar (required)

```bash
go test ./internal/... -count=1 2>&1 | tee "$EVID/02-internal-test.txt"
```

**Pass:** exit 0.

### Block 3 — Temp repro (required)

```bash
REPO=/home/ali/Desktop/Trace
cd "$REPO"
go build -o bin/trace ./cmd/trace 2>&1 | tee "$EVID/03-go-build.txt"
TMP=$(mktemp -d)
TRACE="$REPO/bin/trace"
{
  echo "tmp=$TMP"
  "$TRACE" -C "$TMP" init
  if test ! -e "$TMP/trace.db"; then echo "PASS: init did not create root trace.db"; else echo "FAIL: root trace.db after init"; fi
  if test -f "$TMP/.trace/trace.db"; then echo "PASS: .trace/trace.db exists after init"; else echo "FAIL: missing .trace store"; fi
  ( cd "$TMP" && python3 -c "import sqlite3; sqlite3.connect('trace.db').close()" )
  ls -la "$TMP/trace.db" "$TMP/.trace/trace.db"
  STAT_BEFORE=$(stat -c '%s %Y' "$TMP/trace.db")
  # Command that opens the store (tasks). Capture stderr for warn.
  "$TRACE" -C "$TMP" tasks >"$EVID/03-tasks-stdout.txt" 2>"$EVID/03-warn-stderr.txt" || true
  STAT_AFTER=$(stat -c '%s %Y' "$TMP/trace.db")
  echo "stub_before=$STAT_BEFORE stub_after=$STAT_AFTER"
  if test "$STAT_BEFORE" = "$STAT_AFTER"; then echo "PASS: stub untouched"; else echo "FAIL: stub mutated"; fi
  # Live store still under .trace/ (exists as regular file)
  if test -f "$TMP/.trace/trace.db"; then echo "PASS: live store still .trace/trace.db"; else echo "FAIL: live store missing"; fi
} 2>&1 | tee "$EVID/03-repro-init.txt"
rm -rf "$TMP"

grep -F 'project-root trace.db exists but is not the Trace store' "$EVID/03-warn-stderr.txt" \
  && grep -F '.trace/trace.db' "$EVID/03-warn-stderr.txt" \
  && grep -F 'agents: use CLI/MCP' "$EVID/03-warn-stderr.txt" \
  && echo "PASS: warn observed" | tee "$EVID/03-warn-verdict.txt" \
  || echo "FAIL: warn missing" | tee "$EVID/03-warn-verdict.txt"
```

**Notes:** Prefer `tasks` (or any CLI that opens the store). Size/mtime of `.trace/trace.db` may change on open — OK. **Root stub** must stay size(+mtime) unchanged. Init must **not** create `$TMP/trace.db`.

**Pass:** no root db after init; python creates root stub; Trace still uses `.trace/`; warn on stderr with locked substrings; stub untouched.

### Block 4 — Docs / gitignore / join spot-check (required)

```bash
{
  echo "=== gitignore ==="
  grep -n '/trace.db' .gitignore fixtures/x0/.gitignore
  echo "=== AGENTS / rules / CONTRIBUTING (canonical path) ==="
  grep -n 'trace.db\|\.trace/' AGENTS.md docs/rules/project-rules.md CONTRIBUTING.md | head -n 80
  echo "=== open.go join + warn choke ==="
  grep -n 'traceDirName\|dbFileName\|warnIfStrayRootTraceDB\|strayRootDBWarn\|filepath.Join' internal/store/open.go
} 2>&1 | tee "$EVID/04-docs-gitignore-join.txt"
```

**Pass:** `/trace.db` present in both gitignores; docs state live store is `.trace/trace.db` / agents must not create root stub; `warnIfStrayRootTraceDB` present; join remains `.trace` + `trace.db` (Stat-only on root name — no delete/open of stub).

### Block 5 — Residuals list (required in notes)

Record non-blocking residuals from the locked table above. Do **not** fail solely for them.

## Do not

- Start a new product feature wave or change `open.go` / gitignore / docs in this row
- Close DR-HANDOFF (that is **P30-S03-02**)
- Claim successor (S03-02 only; default lean **no successor**)
- Treat root stub as a second store
- Silent-delete the stub “to clean evidence”

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-30-stray-trace-db-hygiene/scopes/scope-03-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P30-S03-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p30-s03-01-verify/evidence/

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| 1 Store stray tests | PASS/FAIL | |
| 2 go test ./internal/... | PASS/FAIL | |
| 3 Temp repro (init / stub / warn / untouched) | PASS/FAIL | |
| 4 Docs/gitignore/join | PASS/FAIL | |
| 5 Residuals | listed | |

## Residuals (non-blocking)
- …

## Failures (if any)
- …
```

## Todo updates

Status + notes on **P30-S03-01** only (implementer). Do not edit upcoming S03-02 prompt body beyond what Notes need; do not mark S03-02 done.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0–5
- [ ] Evidence dir present under `experiments/runs/…-p30-s03-01-verify/evidence/`
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P30-S03-02**

## Next

`P30-S03-02`
