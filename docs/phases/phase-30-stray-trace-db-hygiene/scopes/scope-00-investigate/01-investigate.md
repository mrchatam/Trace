# P30-S00-01 — Investigate implementer

## Metadata
- id: P30-S00-01
- todo_ids: [P30-S00-01]
- role: implementer
- skills: [diagnosing-bugs, systematic-debugging]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent forensic pass vs [`INTAKE.md`](../../INTAKE.md). Write `INVESTIGATION.md` with an explicit verdict. **Investigation only — no product code.**

## Session start

Follow agent-loop-protocol Session start. Do not treat INTAKE as proven. Prefer repo + temp-dir evidence over the external dogfood tree (cite intake only; do not delete dogfood files).

## Locked defaults (planner-verified 2026-08-21)

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-investigate/INVESTIGATION.md` (this scope folder) |
| Product Go / install edits | **Forbidden** |
| Canonical path until overturned | `<root>/.trace/trace.db` |
| Temp repro | Disposable `mktemp -d` under `/tmp` (or similar) |
| HTTP | Path-risk check only — no serve feature work |
| Dogfood | Do **not** mutate/delete files under the external dogfood project named in INTAKE |

## Must answer (INVESTIGATION.md must cover all five)

1. Does any Trace code path create `<root>/trace.db`?
2. Can CLI / MCP / `trace serve` accidentally open a store at that path (cwd vs `-C` / MCP `project`)?
3. Reproduce (or refute) the 0-byte stub via `sqlite3.connect('trace.db')` from project cwd after `trace init`.
4. Confirm `internal/store/open.go` join is only `.trace` + `trace.db` (cite lines).
5. Phase 29 HTTP: any new write/open of root `trace.db`?

## Live anchors (re-verify; adjust citations if moved)

Planner gate (2026-08-21) confirmed these still hold — re-Stat / re-grep before citing:

| Topic | Citation |
|-------|----------|
| Constants | `internal/store/open.go` L15–16: `traceDirName = ".trace"`, `dbFileName = "trace.db"` |
| OpenExisting join | `open.go` L57: `filepath.Join(absRoot, traceDirName, dbFileName)` |
| openStore join | `open.go` L77 + L92: `traceDir` then `dbPath := filepath.Join(traceDir, dbFileName)` |
| DBPath helper | `open.go` L154: `filepath.Join(s.projectRoot, traceDirName, dbFileName)` |
| MCP | `internal/mcp/project.go` L17–33 `resolveProject` (override → defaultRoot → cwd); L41 `store.OpenExisting(abs)` |
| HTTP | `internal/httpapi/server.go` L193–198 `openStore` → `store.Open(s.root)`; L97–99 refuse `StaticDir ==` project root (exposes `.trace/`) |
| CLI help | `cmd/trace/help.go` documents `.trace/trace.db` (init/backup/restore strings) |
| This checkout | Live store `.trace/trace.db` present; **no** root `trace.db` (absence ≠ intake false) |

## Preflight

1. Confirm Phase 29 is **done** on `docs/TODO.md` (not active); Phase 30 is active.
2. Confirm `INVESTIGATION.md` does not already claim a final verdict (overwrite only if restarting with board Notes).
3. Skim `internal/store/open.go` constants + `openStore` / `OpenExisting` dbPath joins before searching.
4. Confirm you will **not** edit product Go, install rules, or gitignore.

## Search roots (required)

Run ripgrep (or equivalent) for `trace.db`, path joins to db, and accidental root opens:

| Area | Paths |
|------|-------|
| Store | `internal/store/` |
| CLI | `cmd/trace/` |
| MCP | `internal/mcp/` |
| HTTP (P29) | `internal/httpapi/` |
| Install (guidance only) | `internal/install/` — note whether rules mention store path |

Record every hit that **creates**, **opens**, or **documents** a SQLite path. Distinguish documentation strings from executable joins. Especially flag any `Join(..., "trace.db")` **without** `.trace`.

## Repro protocol (temp dir)

Use a disposable directory (do **not** mutate the external dogfood project):

```bash
TMP=$(mktemp -d)
# From Trace repo: use `go run ./cmd/trace` or an existing `trace` binary on PATH
# 1) Init live store
trace -C "$TMP" init   # or: go run ./cmd/trace -C "$TMP" init
# 2) Confirm only .trace/trace.db exists (no root stub yet)
ls -la "$TMP" "$TMP/.trace"
# 3) Agent-style stub (relative name from project cwd)
(
  cd "$TMP"
  python3 -c "import sqlite3; sqlite3.connect('trace.db').close()"
)
ls -la "$TMP/trace.db"   # expect 0-byte (or sqlite header size — note actual size)
# 4) Trace still uses .trace/ — note DBPath if printed
trace -C "$TMP" tasks    # or version/status that opens store
# Optional: from inside $TMP without -C, confirm cwd-root still opens .trace/ not ./trace.db
# 5) Cleanup
rm -rf "$TMP"
```

Also note MCP: server `-C` / tool arg `project` resolution (`internal/mcp/project.go`, tool schemas). Can cwd alone cause MCP to open `<cwd>/trace.db`? (Expected: no — `OpenExisting` under `.trace/`.)

## HTTP / serve (path risk only)

Confirm `trace serve` / `internal/httpapi` opens via `store.Open(projectRoot)` → `.trace/trace.db`. Note static-dir refusal of project root (exposes `.trace/`) — not a root-`trace.db` creator, but cite if relevant to operator confusion. Answer must-question **#5** explicitly.

## Do not

- Change product code, install rules, or gitignore
- Delete files in the dogfood project named in INTAKE
- Implement warn-on-stray (that is S01/S02)
- Start S01 planning beyond a short “suggested next” section
- Write product fixes “while you are here”

## Deliverable: `INVESTIGATION.md`

Write at `scopes/scope-00-investigate/INVESTIGATION.md` using this template:

```markdown
# INVESTIGATION — stray root trace.db

## Verdict
**Trace bug** | **agent hygiene** | **both** — one sentence why.

## Evidence: store path
- File:line citations for Open / OpenExisting / dbPath join
- Statement: Trace does / does not open `<root>/trace.db`

## Evidence: creators of root stub
- Trace code paths that create `<root>/trace.db` (none | list)
- Repro of python/sqlite stub (commands + observed size)

## CLI / MCP / HTTP
- How project root is chosen (-C, cwd, MCP project)
- Any accidental open of root `trace.db`? (yes/no + proof)
- Phase 29 HTTP: any new root write/open? (yes/no + proof)

## Relation to INTAKE.md
- Confirmed / partially confirmed / overturned (what changed)

## Recommendations for S01 (non-binding)
- Docs / warn / gitignore / install / path fix (only if Trace bug)
```

## Exit criteria

- [ ] `INVESTIGATION.md` present with explicit verdict
- [ ] All five **Must answer** items addressed
- [ ] Store path cited with file:line
- [ ] Stub repro confirmed or refuted with commands
- [ ] CLI/MCP/HTTP path risk addressed
- [ ] No product code diffs in this scope
- [ ] Board row P30-S00-01 Notes summarize verdict + evidence pointers

## Todo updates

Status + notes on **P30-S00-01** only.

## Next

`P30-S00-02`
