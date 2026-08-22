# P30-S00-00 — Scope planner (investigate)

## Metadata
- id: P30-S00-00
- todo_ids: [P30-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, diagnosing-bugs, ask-questions-if-underspecified]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock investigation scope so a fresh subagent produces `INVESTIGATION.md` that independently confirms or overturns [`INTAKE.md`](../../INTAKE.md). **No product code.** Finalize sibling `01-investigate.md` + `02-review.md` if still thin; do not start S01.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [Phase 30 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [INTAKE.md](../../INTAKE.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Prefer resolving from repo evidence over asking unless a verdict hinges on missing dogfood artifacts.

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-investigate/INVESTIGATION.md` |
| Product Go / install edits | **No** on S00-01 / S00-02 |
| Canonical path assumption | `<root>/.trace/trace.db` until overturned with file:line + repro |
| Dogfood project | Do **not** delete files under the external dogfood tree; cite intake only |
| Temp repro | Use a disposable `mktemp -d` under `/tmp` (or similar) |
| HTTP note | Check `internal/httpapi` only for path risk — no serve feature work |
| Sequence | S00 → S01 → S02 → S03 serial |

## Must answer (handoff to 01)

1. Does any Trace code path create `<root>/trace.db`?
2. Can CLI / MCP / `trace serve` accidentally open a store at that path (cwd vs `-C` / `project`)?
3. Reproduce (or refute) the 0-byte stub via `sqlite3.connect('trace.db')` from project cwd after `trace init`.
4. Confirm `internal/store/open.go` join is only `.trace` + `trace.db` (cite lines).
5. Phase 29 HTTP: any new write/open of root `trace.db`?

## Live anchors (verify still true at planner time)

| Topic | Path / note |
|-------|-------------|
| Open join | `internal/store/open.go` — `traceDirName = ".trace"`, `dbFileName = "trace.db"`, `filepath.Join(absRoot, traceDirName, dbFileName)` |
| MCP open | `internal/mcp/project.go` → `store.OpenExisting` |
| HTTP open | `internal/httpapi/server.go` → `store.Open(s.root)` |
| CLI help | `cmd/trace/help.go` documents `.trace/trace.db` |
| This Trace checkout | Live store under `.trace/`; root `trace.db` may be absent (absence ≠ intake false) |

## Planner gate

- [ ] `01-investigate.md` runnable (metadata, preflight, search roots, repro steps, INVESTIGATION template, exit criteria)
- [ ] `02-review.md` has verify checklist + spawn policy (`02a`/`02b` only on blocker/high)
- [ ] `SCOPE-TODOS.md` lists S00 board rows
- [ ] Live anchors still accurate (adjust `01` if renamed)

## Exit criteria

- [ ] Investigation implementer prompt locked enough for a fresh subagent
- [ ] Board row P30-S00-00 Notes cite what was verified/thickened
- [ ] Next runnable remains **P30-S00-01** (do not start S01)

## Todo updates

Status + notes on **P30-S00-00** only.

## Next

`P30-S00-01`
