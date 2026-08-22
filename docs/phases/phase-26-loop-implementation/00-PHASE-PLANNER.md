# Phase 26 — Loop investigations, planning & implementation

**Phase planner.** Runs as row `P26-00` on the board.

## Metadata
- id: P26-00
- todo_ids: [P26-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: []
- verification: automated
- hooks: []

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended orchestrator runs may execute after prior prompt approval.

## Mission

Address the three open issues surfaced by E02 / Phase 25 DR-HANDOFF:

| Issue | INT-IDs | Scope |
|-------|---------|-------|
| `ParentOrchestratorRule` not wired into installed cursor rules | INT-04 wiring gap | S04 |
| Discoveries produce 0 new tasks (FM-10) | INT-01, INT-06 | S02 (P25-A) |
| Loop saturates early; no deliberation reset after gap pass (FM-03) | INT-02, INT-05, INT-09 | S03 (P25-B) |

## Scope sequence

```
S00 Investigation → S01 Planning → S02 P25-A impl → S03 P25-B impl → S04 Installer fix → S05 VERIFY
```

S04 (installer fix) is small and low-risk — may run in parallel with S03 after S01 planning is done, but default sequence is serial.

## Hard constraints (from project laws)

- No daemon, HTTP server, or hosted service
- No rewriting done board history
- No full-rebuild-on-any-change indexer
- SQLite migration permitted only if schema version bumped properly
- `go test ./...` must stay green on `internal/` packages after each scope

## Deliverables expected by VERIFY (S05)

| ID | Deliverable | Source scope |
|----|-------------|-------------|
| D1 | `ParentOrchestratorRule` text appears in `cursorRulesMDCContent` output | S04 |
| D2 | `loop apply` or `trace add task` promotion path for BLOCKING discoveries | S02 |
| D3 | P19/hop thresholds recalibrated; greenfield build no longer sticky-STOPs verify | S03 |
| D4 | Deliberation reset API clears `Stopped`, resets `hop_count` on verify task | S03 |
| D5 | Unified STOP reason UX (consistent `reason_code` in gate JSON + export) | S03 |
| D6 | `go test ./internal/...` PASS; `score.sh G1 --p25` fully PASS on fresh install | S04/S05 |

## Planner gate (this scope — P26-00)

Verify that all scope scaffold paths below exist before closing this row:

- `docs/phases/phase-26-loop-implementation/` directory ✓ (this file)
- `scopes/scope-00-loop-audit/` — 00-PLANNER.md, 01-loop-audit.md
- `scopes/scope-01-planning/` — 00-PLANNER.md, 01-task-breakdown.md
- `scopes/scope-02-discovery-task-promotion/` — 00-PLANNER.md, 01-implement.md, 02-review.md
- `scopes/scope-03-loop-recalibration/` — 00-PLANNER.md, 01-implement.md, 02-review.md
- `scopes/scope-04-installer-fix/` — 00-PLANNER.md, 01-implement.md, 02-review.md
- `scopes/scope-05-verify/` — 00-PLANNER.md, 01-verify.md, 02-dr-handoff.md
- `DR-HANDOFF.md` with Status: OPEN
- Board row 435 (P26-00) in `docs/TODO/phase-26.md` ✓

If any path is missing: create it before marking P26-00 done.
