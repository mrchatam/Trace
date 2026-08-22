# P30-S02-02 — Implement review (T1–T4)

## Metadata
- id: P30-S02-02
- todo_ids: [P30-S02-02]
- role: reviewer
- skills: [code-review-and-quality, debugging-and-error-recovery]
- mcps: [user-trace]
- verification: automated
- hooks: []

## Objective

Independent review of **P30-S02-01** against [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) and phase locks. Fresh subagent — do not share the implementer’s session. S00 baseline: verdict **agent hygiene** — **reject** any “path fix” that redesigns the store or treats root `trace.db` as a second store.

## References

- [`docs/rules/agent-loop-protocol.md`](../../../../../rules/agent-loop-protocol.md)
- [`docs/rules/project-rules.md`](../../../../../rules/project-rules.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- Implementer Notes on board row `P30-S02-01`
- Live anchors: `internal/store/open.go` (`openStore`, `DBPath`); `.gitignore`; `fixtures/x0/.gitignore`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

**Clarification:** none required for locks. If implementer Notes claim a path redesign or delete flag, treat as **blocker** and spawn remediation — do not “accept” scope creep.

## Locked defaults (review against these)

| Item | Value |
|------|-------|
| Canonical DB | `.trace/trace.db` only |
| Path redesign | Forbidden |
| Delete | No silent/default auto-delete of root `trace.db` |
| Warn | stderr (or test writer); once-per-`openStore`; non-fatal; Stat only — never opens/deletes root file |
| Gitignore | `/trace.db` root-only; `.trace/` still ignored |
| Green bar | `go test ./internal/...` PASS |
| Creep | No HTTP/daemon/GUI |

## Preflight

1. Read `P30-S02-01` Notes (claimed files + tests).
2. Diff claimed surfaces vs PLAN §3 file lists.
3. Re-run `go test ./internal/...` yourself (do not trust implementer claim alone).

## Checklist (evidence required)

### Path / store

- [ ] `open.go` join still `filepath.Join(..., ".trace", "trace.db")` (constants `traceDirName`/`dbFileName`) — meaning unchanged
- [ ] No new `sql.Open` / create of `<root>/trace.db`
- [ ] `DBPath()` still under `.trace/`

### T2 warn

- [ ] Single choke in `openStore` (covers `Open`, `OpenExisting`→`openStore`, Restore rebind path)
- [ ] Trigger: regular file at `<absRoot>/trace.db` only; missing/dir/perm on stub check → no warn, open not failed for that reason
- [ ] Channel: stderr or injectable writer — not stdout/MCP JSON
- [ ] Once per `openStore` invocation that sees stray (not a persistent suppress file)
- [ ] Non-fatal: open still succeeds when stub present
- [ ] Never deletes, renames, migrates, or opens the root file
- [ ] Message intent matches PLAN draft (root is not Trace store; use `.trace/trace.db`; agents CLI/MCP)

### T3 gitignore

- [ ] Repo `.gitignore` has `/trace.db` (leading slash)
- [ ] `fixtures/x0/.gitignore` has `/trace.db`
- [ ] `.trace/` still ignored; pattern does not target `.trace/trace.db`

### T1 docs

- [ ] `AGENTS.md` states canonical `.trace/trace.db`; never open/create root stub
- [ ] `docs/rules/project-rules.md` Store row updated accordingly
- [ ] `CONTRIBUTING.md` local-store section + ignore guidance
- [ ] Optional `help.go` one-liner OK if present; **no** Phase 29 HTTP doc rewrites required
- [ ] No dual-store claim; no “path fix” language that contradicts S00

### T4 tests

- [ ] Stub present + `Open` → warn + `.trace/` path + success
- [ ] Stub present + `OpenExisting` → same
- [ ] No stub → no warn + success
- [ ] After warn → root stub still present, size unchanged
- [ ] `go test ./internal/...` PASS (reviewer-run)

### Scope

- [ ] No HTTP/daemon/GUI product changes
- [ ] No silent/default delete of operator root `trace.db`
- [ ] No store-path redesign

## Severity → action

| Severity | Action |
|----------|--------|
| blocker / high | Small inline fix **or** spawn `P30-S02-02a` (implement) + `P30-S02-02b` (review) immediately below this row with full protocol prompts |
| medium | Prefer spawn unless trivial one-liner |
| low / nit | Notes for S03 or inline if trivial |

Re-verify until no open blocker/high without a pending follow-up. Confidence **medium** or **high** with evidence; never silent residuals.

## Spawn policy

If spawning:

- Insert board rows `P30-S02-02a` / `P30-S02-02b` under this review
- Write full prompts under `scopes/scope-02-implement/` (metadata, objective, locked defaults, exit criteria, session-start gate)
- Do **not** rewrite `done` `01-implement.md` history; remediate forward

## Todo updates

Reviewer: status + Notes on `P30-S02-02`; may spawn/thicken **upcoming** only. Mark PASS/FAIL with evidence (cite files, test command, key line refs).

## Exit criteria

- [ ] PASS or FAIL in Notes with evidence (checklist items + `go test ./internal/...`)
- [ ] No open blocker/high without pending 02a/02b
- [ ] Confidence medium or high; residual risks listed if medium
- [ ] Next: **P30-S03-00** (or **P30-S02-02a** if spawned)

## Minimal todos

- [ ] Diff + checklist vs PLAN
- [ ] Re-run `go test ./internal/...`
- [ ] PASS/FAIL Notes; spawn if needed
- [ ] Board: mark `P30-S02-02` done (or leave follow-ups pending)

## Next

`P30-S03-00` (or `P30-S02-02a` if spawned)
