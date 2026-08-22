# Phase 27 — Protocol measurement + graph honesty

**Phase planner.** Runs as row `P27-00` on the board.

## Metadata
- id: P27-00
- todo_ids: [P27-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: []
- verification: automated
- hooks: []

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended orchestrator runs may execute after prior prompt approval.

## Mission

Address deferred **P25-D** and **P25-E** from Phase 24, promoted because Phase 26 verify left build-only G1 with **P25-3 FAIL** (thin graph) while product loop/promotion/install work passed.

| Theme | INT-IDs | Scope |
|-------|---------|-------|
| Experiment protocol v2 | INT-08, INT-10 | S01 |
| Graph export honesty | INT-07 | S02 |

## Scope sequence

```
S00 Investigation → S01 Protocol v2 → S02 Graph honesty → S03 VERIFY
```

## Hard constraints (from project laws)

- No daemon, HTTP server, or hosted service
- No rewriting done board history
- No full-rebuild-on-any-change indexer
- `go test ./internal/...` must stay green after product-touching scopes

## Planner gate (this scope — P27-00)

Verify that all scope scaffold paths below exist before closing this row:

- `docs/phases/phase-27-protocol-measurement-graph-honesty/` directory ✓ (this file)
- `scopes/scope-00-investigation/` — 00-PLANNER.md, 01-investigation.md, 02-review.md
- `scopes/scope-01-protocol-v2/` — 00-PLANNER.md, 01-implement.md, 02-review.md
- `scopes/scope-02-graph-honesty/` — 00-PLANNER.md, 01-implement.md, 02-review.md
- `scopes/scope-03-verify/` — 00-PLANNER.md, 01-verify.md, 02-dr-handoff.md
- `DR-HANDOFF.md` with Status: OPEN
- Board row in `docs/TODO/phase-27.md` ✓

If any path is missing: create it before marking P27-00 done.

## Exit criteria

- [ ] Scope list and run order verified against live repo + Phase 26 verify residuals
- [ ] Scope stubs are runnable and references are valid
- [ ] Board points to next runnable scope planner (`P27-S00-00`)

## Next

`P27-S00-00`
