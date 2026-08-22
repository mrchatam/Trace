# P30-S02-01 — Implement T1–T4 hygiene

## Metadata
- id: P30-S02-01
- todo_ids: [P30-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd, debugging-and-error-recovery]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Ship **only** PLAN.md tasks **T1–T4** for Phase 30 stray root `trace.db` hygiene. Source of truth: [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). S00 verdict is **agent hygiene** (INTAKE confirmed; **no** Trace dual-store bug; **no** store-path change). Fresh implementer should finish alone without re-planning.

## References

- [`docs/rules/agent-loop-protocol.md`](../../../../../rules/agent-loop-protocol.md)
- [`docs/rules/project-rules.md`](../../../../../rules/project-rules.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — SoT
- [`../scope-00-investigate/INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md) — do not re-open
- Live choke: `internal/store/open.go` (`openStore` L71+, `DBPath` L152–155; join is `filepath.Join(absRoot, ".trace", "trace.db")` via `traceDirName`/`dbFileName`)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

**Clarification:** none required — PLAN.md + S00 locks are closed. Do not ask to redesign the store path. Produce a short mental plan matching § Implementation order below, then code.

## Locked defaults (do not renegotiate)

| Item | Value |
|------|-------|
| Canonical DB | `<projectRoot>/.trace/trace.db` only (`store.Open` / `OpenExisting` → `openStore`) |
| Path redesign | **Forbidden** — do not change `dbPath` / `filepath.Join(..., ".trace", "trace.db")` |
| Delete | **No** silent / default / auto-delete of root `trace.db`; optional delete flag is **future-only** (not this row) |
| Warn (T2) | stderr (or injectable `io.Writer`); **once per `openStore` call** that sees stray; **non-fatal**; never opens/deletes/renames root file |
| Gitignore (T3) | Root-only `/trace.db` (leading slash); keep `.trace/` ignored |
| Tests (T4) | Required; green bar `go test ./internal/...` |
| Scope creep | No HTTP / daemon / GUI / Phase 29 doc rewrites; no `internal/install/*` unless you find a single shared agent-facing store-location string — prefer AGENTS + project-rules (PLAN §3 T1) |
| Product Go | Only T2 warn + T4 tests surfaces under `internal/store` |

### T2 warn design (copy from PLAN §4)

| Item | Decision |
|------|----------|
| **Choke** | `openStore` in `internal/store/open.go`, after `filepath.Abs` succeeds and before (or immediately after) `.trace/` mkdir — **one** place so CLI, MCP (`OpenExisting` → `openStore`), and `trace serve` (`store.Open`) all warn |
| **Trigger** | `os.Stat(filepath.Join(absRoot, "trace.db"))` succeeds **and** mode is regular file. Directories / missing / permission errors on the stub check → **no warn** (do not fail open) |
| **Frequency** | Once per `openStore` invocation that sees the stray. Not a persistent on-disk flag |
| **Channel** | stderr only (or package-level `io.Writer`, default `os.Stderr`, for tests). Never stdout / JSON MCP payloads |
| **Severity** | Non-fatal. Open **must still succeed** (subject to existing lock/auth/db errors) |
| **Draft message** | `trace: warning: project-root trace.db exists but is not the Trace store; using .trace/trace.db. Do not open or create a root trace.db (agents: use CLI/MCP).` |
| **Must not** | Delete, rename, migrate, or `sql.Open` the root file; change the `.trace`+`trace.db` join; treat stub as second store |

## Preflight

1. Re-read PLAN.md §3–§6 (tasks, warn design, order, test plan).
2. Confirm live `open.go`: join still `.trace`+`trace.db`; no existing root-stub warn.
3. Confirm `.gitignore` and `fixtures/x0/.gitignore` currently have `.trace/` only (no `/trace.db` yet).
4. Prefer new test file or cases beside existing open tests in `internal/store/store_test.go` / `*_test.go`.

## Implementation order (PLAN §5 — follow exactly)

1. **T2** — stray check + warn writer hook in `openStore` (**no** path join change).
2. **T4** — tests against T2 (stub present / absent / stub untouched). Prefer writing tests immediately after T2 (or TDD: failing tests then warn).
3. **T3** — repo `.gitignore` + `fixtures/x0/.gitignore` → add `/trace.db`.
4. **T1** — docs: `AGENTS.md`, `docs/rules/project-rules.md`, `CONTRIBUTING.md`; optional one-liner in `cmd/trace/help.go` (init/store clarity only).

**Parallel-safe note:** T3 is file-disjoint from T2/T4 Go; T1 docs are disjoint except CONTRIBUTING’s ignore sentence (align with T3). Prefer T2+T4 before docs so acceptance is executable.

## Task cards

### T2 — Warn once in `openStore`

**Files:** `internal/store/open.go` (and tiny helper in same package if needed for writer injection).

**Do:**
- After Abs root resolves, if `<absRoot>/trace.db` is a regular file, write the draft warning once to the warn writer.
- Keep mkdir / lock / token / `sql.Open` on `.trace/trace.db` unchanged.
- Optional: `var warnWriter io.Writer = os.Stderr` (or setter) so tests capture output without hijacking process stderr permanently — restore in tests.

**Acceptance:**
- Trigger after Abs; once per `openStore` that sees stray; non-fatal; never deletes/renames/opens root file; `DBPath()` still under `.trace/`.

### T4 — Tests

**Files:** `internal/store/*_test.go` (new cases beside existing open tests; e.g. `store_test.go` or `stray_trace_db_test.go`).

**Command:** `go test ./internal/...` (also fine to iterate with `go test ./internal/store/ -count=1`).

| Case | Expect |
|------|--------|
| Stub present + `Open` | Warn observed; `DBPath()` / open target is `.trace/trace.db`; open succeeds |
| Stub present + `OpenExisting` (after init/`Open` once) | Same warn + path + success |
| No stub | No warn; open succeeds |
| After warn | Live DB under `.trace/`; root stub still present; size unchanged |

**Do not** assert delete/rename/open of the root file. No HTTP/GUI tests.

### T3 — Gitignore

**Files:**
- `.gitignore` — add `/trace.db` beside existing `.trace/`
- `fixtures/x0/.gitignore` — same pattern
- CONTRIBUTING (in T1) mentions consumers should ignore root `/trace.db` without un-ignoring `.trace/`

**Acceptance:** Pattern is `/trace.db` (leading slash); does not target `.trace/trace.db`; `.trace/` remains ignored.

### T1 — Docs

**Files (required):**
- `AGENTS.md` — stack one-liner: SQLite is `.trace/trace.db` only; never open/create project-root `trace.db` (agents: CLI/MCP)
- `docs/rules/project-rules.md` — Store row: canonical `.trace/trace.db`; root `trace.db` is not Trace
- `CONTRIBUTING.md` — local store section: live store under `.trace/`; warn about accidental root stub; recommend `/trace.db` ignore

**Optional:**
- `cmd/trace/help.go` — one-liner clarifying init/store path only (already mentions `.trace/trace.db`; may add “not project-root `trace.db`”)

**Skip:** Phase 29 HTTP docs; `internal/install/*` rule bodies unless a single shared agent-facing string already documents store location (PLAN prefers AGENTS + project-rules).

**Acceptance:** Readers see: live store is `.trace/trace.db` only; root `trace.db` is not Trace; no dual-store claim.

## Hard checks (every task)

- [ ] `filepath.Join` for the live store still includes `.trace` + `trace.db` under project root
- [ ] No new code path opens or creates `<root>/trace.db` (Stat for warn only)
- [ ] Gitignore pattern is `/trace.db` (root-only), not a pattern that fights `.trace/`
- [ ] No delete of root stub
- [ ] No HTTP/daemon/GUI changes

## Do not

- Expand into HTTP/GUI/daemon
- Change canonical store path
- Auto-delete root `trace.db`
- Re-litigate S00 investigation or invent a “path fix”
- Implement optional documented delete flag

## Role work / loop

```text
LOOP:
  1. Implement next incomplete task in order T2 → T4 → T3 → T1
  2. Self-check against that task’s Acceptance + Hard checks
  3. Fix gaps
  UNTIL all T1–T4 done OR failed/blocked with reason
→ Update own board row status + Notes only (files + test commands)
```

## Todo updates

Implementer: **status + notes only** on `P30-S02-01`. Do not rewrite this prompt or upcoming prompts after `done`.

Board Notes must include:
- Files changed (path list)
- Test command(s) run + PASS
- Explicit: join unchanged; no delete; warn once-per-open non-fatal

## Exit criteria

- [ ] T1–T4 complete per PLAN acceptance (or Notes with explicit defer — should be none for Phase 30)
- [ ] T4 cases above PASS
- [ ] `go test ./internal/...` PASS
- [ ] No path redesign; no silent delete; no HTTP creep
- [ ] Board Notes: files + tests; next pointer **P30-S02-02**

## Minimal todos

- [ ] T2: warn in `openStore` + optional warn writer
- [ ] T4: four named test cases; `go test ./internal/...` green
- [ ] T3: `/trace.db` in repo + `fixtures/x0` gitignores
- [ ] T1: AGENTS + project-rules + CONTRIBUTING (+ optional help)
- [ ] Board: mark `P30-S02-01` done with Notes

## Next

`P30-S02-02`
