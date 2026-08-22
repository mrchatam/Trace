# P31-S00-01 — Inventory (write GAPS.md)

## Metadata
- id: P31-S00-01
- todo_ids: [P31-S00-01]
- role: implementer
- skills: [diagnosing-bugs, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Produce [`GAPS.md`](GAPS.md) that inventories missing automated coverage for Phase 30 stray-root hygiene. **Docs only — no Go, no new tests, no scripts.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [Phase 31 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — suggested coverage candidates
- [P30 VERIFY-NOTES](../../../phase-30-stray-trace-db-hygiene/scopes/scope-03-verify/VERIFY-NOTES.md)
- [P30 S02-02 board Notes](../../../../TODO/phase-30.md) — dir-stub quiet nit (row P30-S02-02)
- Live: `internal/store/stray_trace_db_test.go`, `internal/store/open.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Re-verify live anchors below; do not reopen store-path design.

## Locked defaults

| Item | Value |
|------|-------|
| Deliverable | `GAPS.md` in this folder |
| Product / test code | **Forbidden** |
| Path redesign / silent delete / GUI | **Forbidden** |
| Categories | must-add \| nice-to-have \| out-of-scope \| defer-with-reason |
| Canonical path | `<root>/.trace/trace.db` (locked) |
| Preferred test home | `internal/store` unit (`stray_trace_db_test.go` or sibling) before CLI/`serve` harness |
| Script home if required | `scripts/repro-stray-trace-db.sh` (no `scripts/` dir yet) or `experiments/…` |

## Live snapshot (P31-S00-00 verified 2026-08-21 — re-check before writing GAPS)

| Topic | Evidence |
|-------|----------|
| Warn choke | `open.go` L85 `warnIfStrayRootTraceDB(absRoot)` inside `openStore`; covers `Open`, `OpenExisting`→`openStore`, Restore rebind |
| Stat-only regular-file | `open.go` L144–149: `os.Stat` on `<absRoot>/trace.db`; `err != nil \|\| !fi.Mode().IsRegular()` → **quiet return** (no warn, open not failed for stub reason) |
| Join | `traceDirName=".trace"` + `dbFileName="trace.db"` → live DB `filepath.Join(traceDir, dbFileName)` L102 |
| Unit tests present | `TestOpenWarnsWhenRootStubPresent`, `TestOpenExistingWarnsWhenRootStubPresent`, `TestOpenQuietWhenNoRootStub`, `TestOpenLeavesRootStubUntouched` — **no** dir-named stub quiet case |
| Gitignore with `/trace.db` | **Only** repo `.gitignore` + `fixtures/x0/.gitignore`. `web/.gitignore` is frontend-only (no `/trace.db`). Experiment `.gitignore`s lack it. **No** `internal/install` scaffold writes `/trace.db` |
| CLI path | `cmd/trace/tasks.go` → `store.Open(abs)` (stderr warn when stub present) |
| Serve / HTTP | `internal/httpapi/server.go` L193–194 `openStore` → `store.Open(s.root)` — **per handler request**, not a dedicated process-startup open |
| MCP | `OpenExisting` → same `openStore` warn choke (unit already covers OpenExisting + regular stub) |
| Durable repro script | **None** checked in; P30 VERIFY one-shot under `experiments/runs/2026-08-21-p30-s03-01-verify/evidence/` |
| P30 residuals | Agents can still create stubs; delete future-only; multi-open warn (once per `openStore`) intentional — no suppress flag |

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f internal/store/stray_trace_db_test.go
test -f internal/store/open.go
rg -n 'warnIfStrayRootTraceDB|/trace\.db|strayRootDBWarn' --glob '!docs/phases/phase-31*/**' --glob '!similar projects/**'
rg -n 'TestOpenWarnsWhenRootStub|TestOpenExisting|TestOpenQuiet|TestOpenLeaves|IsRegular|Mkdir.*trace\.db' internal/store/
rg -n '/trace\.db' --glob '**/.gitignore' --glob '!similar projects/**'
rg -n 'store\.(Open|OpenExisting)|openStore' internal/httpapi/server.go cmd/trace/tasks.go internal/mcp/
test ! -f docs/phases/phase-31-stray-db-residual-tests/scopes/scope-00-inventory/GAPS.md || echo 'GAPS.md already exists — overwrite only if correcting this row'
```

## Must answer (fill in GAPS.md)

1. Which of the six phase-planner candidates are still **untested** in-repo (no durable automated coverage)?
2. For each: **must-add** vs **nice-to-have** vs **out-of-scope** vs **defer-with-reason**?
3. Preferred home for each must-add: `internal/store` unit vs CLI harness vs script-only vs docs-only?
4. Any additional ignore scaffolds beyond `.gitignore` + `fixtures/x0/.gitignore` that should get `/trace.db`?
5. Confirm join + Stat-only regular-file warn still true (cite `open.go` file:line).

## Role work

1. Re-run preflight; cite live line numbers in GAPS header.
2. Read P30 `VERIFY-NOTES.md` residuals + `stray_trace_db_test.go` + `warnIfStrayRootTraceDB` in `open.go`.
3. Grep warn / `/trace.db` / serve open / install gitignore scaffolds (exclude `similar projects/`).
4. For each phase-planner candidate (and any new find), classify using the guidance table below — **your classifications own GAPS**; planner hints are starting points, not locks.
5. Write `GAPS.md` from the template. Do not implement tests or scripts.

## Planner hints (starting points — S00-01 owns final disposition)

| # | Candidate | Live status (S00-00) | Hint |
|---|-----------|----------------------|------|
| 1 | Dir-named root `trace.db` → quiet | Code quiet via `!IsRegular()`; **no unit test** (P30-S02-02 nit) | Strong **must-add** → `internal/store` unit (`Mkdir` root `trace.db`, assert no warn + open OK) |
| 2 | CLI `trace tasks -C` stderr | P30 VERIFY one-shot PASS; **no durable `go test` / script** | Prefer **nice-to-have** or **must-add** only if you want CLI-surface lock; else cover via store unit + script (#5). If must-add: reuse `cmd/trace` harness patterns (`cli_test.go` / exec) |
| 3 | `trace serve` startup warn | HTTP opens via `store.Open` **per request** (`server.go` L193–194), not a separate startup hook | Likely **nice-to-have** / **defer**: store unit already covers `Open` choke; dedicated serve test adds little unless GAPS finds a serve-only path. Do **not** invent a startup open that does not exist |
| 4 | Extra ignore scaffolds | Only product scaffolds with `/trace.db` are repo root + `fixtures/x0` | Likely **out-of-scope** / **defer**: no install template emits consumer `.gitignore`; `web/` and experiment ignores are not product store scaffolds. Note CONTRIBUTING already recommends `/trace.db` for consumers |
| 5 | Dogfood repro script | VERIFY evidence only; no `scripts/` | Strong **must-add** or **nice-to-have**: durable `scripts/repro-stray-trace-db.sh` (init → python stub → warn → `.trace/` live; stub untouched). Prefer script over duplicating VERIFY prose |
| 6 | Document multi-open warn | Intentional in VERIFY residuals; not spelled in AGENTS/CONTRIBUTING as “once per openStore / multi CLI opens OK” | **must-add** as docs-only (short CONTRIBUTING or AGENTS note) **or** nice-to-have; **no** suppress flag |

**Also classify explicitly in Out of scope / Defer:** silent/flagged delete; store-path redesign; GUI (Phase 32); persistent warn suppress.

## GAPS.md template (required sections)

```markdown
# GAPS — Phase 31 stray residual testing

**Date:** YYYY-MM-DD
**Author row:** P31-S00-01
**Baseline tests:** internal/store/stray_trace_db_test.go (list case names)
**open.go cite:** warnIfStrayRootTraceDB / IsRegular (file:line)

## Must-add (S01 must ship)

| ID | Gap | Preferred home | Notes |
|----|-----|----------------|-------|
| G1 | … | store unit / CLI / script / docs | … |

## Nice-to-have

| ID | Gap | Notes |
|----|-----|-------|

## Out of scope (this phase)

| Item | Reason |
|------|--------|
| Store-path redesign | Locked |
| Silent / flagged delete | Locked |
| GUI | Phase 32 |

## Defer-with-reason

| ID | Gap | Why defer |
|----|-----|-----------|

## Candidate disposition (from 00-PHASE-PLANNER)

1. Dir-named root `trace.db` quiet — …
2. CLI `trace tasks -C` stderr — …
3. `trace serve` startup warn — …
4. Extra ignore scaffolds — …
5. Dogfood repro script — …
6. Document multi-open warn — …

## Ignore-scaffold audit

| Path | Has `/trace.db`? | Action |
|------|------------------|--------|
| `.gitignore` | yes/no | keep / … |
| `fixtures/x0/.gitignore` | yes/no | keep / … |
| other finds | … | must-add / out / defer |

## Handoff to S01

- Must-add count: N
- Script required: yes/no
- Docs-only IDs: …
- S01 must not: path change, delete, GUI, invent serve startup open
```

## Exit criteria

- [ ] `GAPS.md` exists with all template sections filled (including ignore-scaffold audit)
- [ ] Every phase-planner candidate has a disposition
- [ ] `open.go` join + Stat-only regular-file cite present
- [ ] No product/test code changes
- [ ] Board Notes cite path to `GAPS.md` + must-add count
- [ ] Next remains **P31-S01-00** (do not start S01)

## Todo updates

Status + notes on **P31-S00-01** only.

## Next

`P31-S01-00`
