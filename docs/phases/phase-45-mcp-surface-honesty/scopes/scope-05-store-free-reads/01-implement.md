# Phase 45 / S05 / Implement — Store-free reads

## Row `P45-S05-01` — **done** 2026-10-04. Issues [#152](https://github.com/mrchatam/Trace/issues/152), [#160](https://github.com/mrchatam/Trace/issues/160).

## 1. `store.ProjectInitialized` — `internal/store/exists.go` (new)

```go
func ProjectInitialized(projectRoot string) bool
```

**Stat-only**: checks the DB file path under the project root. No `store.Open`, no lock, no mkdir, no migration check. This is the only safe way to answer "is a project here?" from a tool that must work on locked or uninitialized projects (#152's whole scenario).

## 2. `toolVersion` rewritten — `internal/mcp/tools_parity.go:314`

Before: `openStore` (creating/locking the DB) → `assertMCPToolAllowed` (the capability gate) → payload.
After: resolves the project root, answers from the **process alone**:

```go
b, err := json.Marshal(buildmeta.Payload(map[string]any{
    "project_initialized": store.ProjectInitialized(abs),
    "project_root":        abs,
}))
```

No store open ⇒ no mkdir on a virgin project, no hang on a locked one, and **no capability gate** — a locked project is exactly when stale-server detection matters most. `{ok,name,version}` still come from `buildmeta`; the description (S04) documents the full response shape.

## 3. `trace_tasks` ordering documented (#160)

`TasksInput.Limit` jsonschema tag now carries the contract:

> results are ordered by created_at ASC, id ASC

…matching the store's `ORDER BY created_at ASC, id ASC` (`helpers.go:112`), and the tool description explains the same-second behaviour: rows created within the same second tie-break on `id`, so **use the cursor for stable paging** (UUIDs are random, not time-ordered — the tiebreak is deterministic but not chronological within the second). No migration, no nanosecond column; the non-goal stands.

## 4. Acceptance — all PASS

- `TestVersionAnswersWithoutStore` (new): constructs a server on a directory with **no** DB and calls `trace_version` — must answer, with `project_initialized:false`, and must not create anything.
- `TestMCPVirginProjectDoesNotMkdir` (rewritten virgin subtest): a virgin `defaultRoot` + `trace_version` leaves the directory empty.
- `TestMCPAssertBuiltinAutoAllowedSucceeds` updated: AUTO_ALLOW durability is now asserted via `trace_tasks` (trace_version no longer runs the gate, so it can no longer prove gate behaviour).

## Next

`P45-S05-02`