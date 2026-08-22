# Phase 30 — Stray root `trace.db` hygiene

**Active** after Phase 29 close (`P29-S07-02`, 2026-08-21). Human-promoted after dogfood: a planning workspace grew an extra **0-byte** `trace.db` at project root beside the live store `.trace/trace.db`.

## Goal

Operators and agents must not confuse a stray root `trace.db` with the live store. Trace must not look like split-brain. Prefer docs/warn/gitignore/install-rule hygiene over path redesign.

## Light locks (do not reopen without S00 evidence)

| Lock | Value |
|------|-------|
| Canonical store | `<projectRoot>/.trace/trace.db` via `store.Open` / `OpenExisting` |
| Path change | **Forbidden** unless S00 proves Trace opens/creates `<root>/trace.db` |
| Delete policy | No silent delete of operator files; any remove requires an explicit, documented flag |
| Scope creep | No daemon/HTTP/GUI work (Phase 29 owns `trace serve`) |
| Tests | Any warn-on-stray behavior needs tests; `go test ./internal/...` stays green |

## Live repo baseline (P30-00, 2026-08-21)

| Area | State |
|------|-------|
| Store | `internal/store/open.go` — `filepath.Join(absRoot, ".trace", "trace.db")` only |
| This checkout | `.trace/trace.db` present; **no** root `trace.db` |
| MCP | `internal/mcp/project.go` → `store.OpenExisting`; `-C` / tool `project` override |
| HTTP | `internal/httpapi` → `store.Open(s.root)` (same join; refuse static=project root) |
| Install | `internal/install/*` (agents/rules/hooks) — candidate for “never open root `trace.db`” guidance |
| Intake | [`INTAKE.md`](INTAKE.md) — hypothesis only until S00 re-verifies |

## Scope index (serial)

```
S00 Investigate → S01 Plan → S02 Implement + review → S03 VERIFY + DR-HANDOFF
```

| Scope | Title | Primary artifact |
|-------|-------|------------------|
| S00 | Independent investigation vs INTAKE | `INVESTIGATION.md` |
| S01 | Hygiene plan (docs / warn / gitignore / install) | `PLAN.md` |
| S02 | Implement planned hygiene + review | code + tests |
| S03 | VERIFY + close handoff | `VERIFY-NOTES.md`, successor |

## In scope

- Re-verify whether Trace ever creates/opens `<root>/trace.db`
- Plan + ship docs, optional warn-on-open, gitignore/scaffold, install-rule text
- Tests for any new warn/path behavior
- Operator clarity: only `.trace/trace.db` is authoritative

## Out of scope

- Moving SoT off `.trace/trace.db` (unless S00 overturns)
- Auto-deleting user files without an explicit documented action
- Phase 29 HTTP/GUI feature work; cloud/hosted SaaS
- Rewriting prior phase `done` history

## Intake → S00

[`INTAKE.md`](INTAKE.md) claims the stub came from an agent `sqlite3.connect('trace.db')` from cwd. **Treat as unproven** until S00 writes `INVESTIGATION.md` with a verdict: Trace bug | agent hygiene | both.

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until S03-02.
