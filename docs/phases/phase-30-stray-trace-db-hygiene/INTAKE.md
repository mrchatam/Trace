# INTAKE — stray root `trace.db`

**Copied 2026-08-21** from live dogfood investigation (do not treat as Trace-repo SoT until Phase 30 S00 re-verifies):

`/home/ali/Desktop/feet seller telegram app/docs/BUG_DUPLICATE_TRACE_DB.md`

## Symptom

After using Trace on an initially empty project (planning phase), the workspace showed **two** `trace.db` files:

| Path | Role (intake claim) |
|------|---------------------|
| `.trace/trace.db` | Live Trace SQLite (~600KB) |
| `<project>/trace.db` | **0-byte stub**, not used by Trace |

## Intake conclusion (must re-verify in S00)

Trace path resolution is **correct** (`internal/store/open.go` → `.trace/trace.db` only). The root stub was created by a **Cursor agent** running `python3 sqlite3.connect('trace.db')` from project cwd (creates empty file on connect-if-missing).

## Intake suggested product work (S01 plans; do not implement in S00)

1. Docs / install rules: never open/create `<project>/trace.db`
2. Optional warn-on-open if stray root `trace.db` exists
3. Recommend gitignore / scaffold ignore for root `trace.db`
4. Agent guidance: use `trace` CLI/MCP, not raw sqlite on `trace.db`

## Workaround (operator)

Delete empty root stub; only `.trace/trace.db` is authoritative.
