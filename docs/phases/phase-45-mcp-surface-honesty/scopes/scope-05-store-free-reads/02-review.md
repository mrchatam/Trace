# Phase 45 / S05 / Review — Store-free reads

## Row `P45-S05-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** Implemented in one session by the coding agent. Evidence, not a second opinion.

## Honesty invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| H1 | `trace_version` answers on a project that does not exist, is locked, or denies the capability (#152) | **yes** | `TestVersionAnswersWithoutStore`; no `openStore`/gate call remains in the handler |
| H2 | `trace_version` has zero side effects | **yes** | Virgin-project subtest: no mkdir, no DB file after the call (`store.ProjectInitialized` is stat-only) |
| H3 | `project_initialized` reflects the filesystem, not an open | **yes** | `exists.go` is a single `os.Stat` on the DB path |
| H4 | The response documents everything it returns (#159 discipline) | **yes** | `descVersion` names `{ok,name,version,project_initialized,project_root}` |
| H5 | `trace_tasks` ordering is stated, and the same-second tiebreak is called out with the remedy (#160) | **yes** | Schema tag + description; cursor advice present |
| H6 | No migration, no column change | **yes** | `git status` — no new migration file; store schema untouched |

## Commands run

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1
CGO_ENABLED=1 go test ./... -count=1
```

All green (full suite passing end-to-end on the final audit pass, including `evals/perf`).

## Residuals (accepted, stated)

- **The same-second UUID tiebreak remains** — deterministic but not chronological within a second. Documented rather than "fixed" per the locked non-goal; if two actors ever need within-second ordering, that is a schema decision, not a doc patch.
- **`project_initialized` is a heuristic** (DB file present ⇒ initialized). A corrupt-but-present DB reads as initialized; the version tool's contract is identity, not health.

## Verdict

**PASS.** #152 and #160 closed by behaviour + permanent tests.