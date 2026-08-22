# INVESTIGATION — stray root trace.db

**Scope:** P30-S00-01 · **Date:** 2026-08-21 · **Role:** investigate implementer (no product code)

## Verdict

**agent hygiene** — Trace never creates or opens `<projectRoot>/trace.db`; every executable store join is `<root>/.trace/trace.db`, and a disposable repro shows a **0-byte** root stub from `python3` `sqlite3.connect('trace.db')` after `trace init`, while CLI continues to use the live `.trace/` store.

## Must-answer summary

| # | Question | Answer |
|---|----------|--------|
| 1 | Does any Trace code path create `<root>/trace.db`? | **No** |
| 2 | Can CLI / MCP / `trace serve` accidentally open a store at that path? | **No** |
| 3 | Reproduce (or refute) 0-byte stub via `sqlite3.connect('trace.db')`? | **Reproduced** (0 bytes) |
| 4 | Is `open.go` join only `.trace` + `trace.db`? | **Yes** (cited below) |
| 5 | Phase 29 HTTP: any new write/open of root `trace.db`? | **No** |

## Evidence: store path

Constants and joins in `internal/store/open.go` (re-verified 2026-08-21):

```15:16:internal/store/open.go
	traceDirName = ".trace"
	dbFileName   = "trace.db"
```

```57:57:internal/store/open.go
	dbPath := filepath.Join(absRoot, traceDirName, dbFileName)
```

```77:92:internal/store/open.go
	traceDir := filepath.Join(absRoot, traceDirName)
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return nil, fmt.Errorf("store: mkdir .trace: %w", err)
	}
	// ...
	dbPath := filepath.Join(traceDir, dbFileName)
```

```152:155:internal/store/open.go
// DBPath returns the absolute path to trace.db (for tests/diagnostics).
func (s *Store) DBPath() string {
	return filepath.Join(s.projectRoot, traceDirName, dbFileName)
}
```

Related executable joins (also under `.trace/`):

- `internal/store/backup.go` L98: `dbPath := filepath.Join(traceDir, dbFileName)` after `resolveTraceDir`
- CLI `cmd/trace/init.go` L33–48: `store.Open(abs)` then prints `st.DBPath()` → `…/.trace/trace.db`
- All production CLI commands open via `store.Open(abs)` / Abs project root from `-C` or cwd (`cmd/trace/root.go` L103–138: default root = cwd; Abs only; no walk-up)

**Statement:** Trace **does not** open `<root>/trace.db`. The only SQLite open in product store code is `sql.Open("sqlite", dbPath)` where `dbPath` is always under `.trace/`.

Ripgrep across `internal/store/`, `cmd/trace/`, `internal/mcp/`, `internal/httpapi/`: no `Join(..., "trace.db")` without a preceding `.trace` directory component. Docs/help strings (`cmd/trace/help.go`) document `.trace/trace.db` only.

`internal/install/`: mentions `.trace/config.json` and agent defaults under `.trace/`; **no** guidance string for store path `trace.db` (gap for S01 docs, not a creator).

## Evidence: creators of root stub

### Trace code paths that create `<root>/trace.db`

**None.** `Open` / `openStore` mkdir `.trace/` then open `.trace/trace.db`. `OpenExisting` stats `.trace/trace.db` only and never creates a root-level file. `Restore` installs into `.trace/trace.db`.

This Trace checkout baseline (live): `.trace/trace.db` present (~1.3MB); **no** root `trace.db` (absence ≠ intake false).

### Repro of python/sqlite stub

Disposable tree under `/tmp` (external dogfood **not** touched):

```bash
TMP=$(mktemp -d)
go build -o /tmp/trace-p30 ./cmd/trace   # from Trace repo
/tmp/trace-p30 -C "$TMP" init
# → prints $TMP/.trace/trace.db
# → ls: only .trace/ (trace.db ~733184 bytes); no $TMP/trace.db

(
  cd "$TMP"
  python3 -c "import sqlite3; sqlite3.connect('trace.db').close()"
)
# → $TMP/trace.db exists, size 0 bytes

/tmp/trace-p30 -C "$TMP" tasks   # → []
# from cwd without -C, after stub:
( cd "$TMP" && /tmp/trace-p30 tasks )
# live .trace/trace.db size+mtime unchanged; root stub remains 0 bytes
```

Observed: root stub **0 bytes** (matches intake “empty stub”; not a SQLite header page until further writes).

## CLI / MCP / HTTP

### How project root is chosen

| Surface | Root resolution | Store open |
|---------|-----------------|------------|
| CLI | Global `-C` / `--project`, else **cwd**; `filepath.Abs` only (`root.go`) | `store.Open(abs)` → `.trace/trace.db` |
| MCP | Tool `project` override → server `defaultRoot` → **cwd**; Abs only (`internal/mcp/project.go` L17–33) | `store.OpenExisting(abs)` L41 → stats/opens `.trace/trace.db` only; no auto-init |
| `trace serve` / HTTP | Serve resolves project Abs root; `Server.openStore` → `store.Open(s.root)` (`server.go` L193–198) | Same `.trace/trace.db` |

### Accidental open of root `trace.db`?

**No.** Proof:

1. Path construction always includes `traceDirName` (`.trace`).
2. Temp repro: with root stub present, CLI from that cwd still serves tasks from `.trace/trace.db` (live size unchanged).
3. MCP cannot open `<cwd>/trace.db` via cwd alone — `OpenExisting` looks only at `<abs>/.trace/trace.db`.

### Phase 29 HTTP: any new root write/open?

**No.**

- Open path: `openStore` → `store.Open(s.root)` only (`server.go` L193–198).
- Readiness check: `handlers_meta.go` L23–26 joins `s.root` + `.trace` then `trace.db` (Stat only; does not create):

```23:26:internal/httpapi/handlers_meta.go
	storePath := filepath.Join(s.root, ".trace")
	storeReady := false
	if fi, err := os.Stat(filepath.Join(storePath, "trace.db")); err == nil && fi.Mode().IsRegular() {
		storeReady = true
```

- Static-dir refusal of project root (`server.go` L97–99) prevents exposing `.trace/`; unrelated to creating root `trace.db`.

## Relation to INTAKE.md

**Confirmed.**

- Intake claim that Trace path resolution is correct → **confirmed**.
- Intake claim that root stub comes from agent-style `sqlite3.connect('trace.db')` from project cwd → **confirmed** by independent temp-dir repro (0-byte file).
- Intake suggested product hygiene (docs / warn / gitignore / agent guidance) remains appropriate for S01; **no store-path change** required.

Nothing in S00 overturns the locked canonical path `<root>/.trace/trace.db`.

## Recommendations for S01 (non-binding)

Keep human lock: **no store-path change**.

1. **Docs / install rules / agent guidance:** never open or create `<project>/trace.db`; use `trace` CLI/MCP; cite `.trace/trace.db`.
2. **Optional warn-on-open** if stray root `trace.db` exists (operator confusion), without changing the open path.
3. **gitignore / scaffold:** ignore root `trace.db` (not `.trace/trace.db`).
4. **Do not** “fix” path resolution — there is no Trace dual-store bug to fix.
