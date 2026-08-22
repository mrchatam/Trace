# Phase 30 — Stray root `trace.db` hygiene

**Phase planner.** Runs as row `P30-00` **only after Phase 29 is closed**.

## Metadata
- id: P30-00
- todo_ids: [P30-00]
- role: planner
- skills: [planning-and-task-breakdown, diagnosing-bugs]
- verification: automated

## Gate: do not run while Phase 29 is active

If `docs/TODO.md` still lists Phase 29 as **active**, stop. Finish `P29-S07-02` first.

## Mission

Investigate the dual-`trace.db` operator confusion, **plan** hygiene (docs/warn/gitignore/install rules), **implement**, and verify.

Do **not** change the canonical store path unless S00 proves Trace actually opens `<root>/trace.db`.

## Intake (starting hypothesis)

See [INTAKE.md](INTAKE.md). Likely: agent `sqlite3.connect('trace.db')` from cwd, not a Trace dual-store bug.

## Scope sequence

```
S00 Investigate (independent) → S01 Plan → S02 Implement + review → S03 VERIFY
```

## Hard constraints

- Canonical store remains `<root>/.trace/trace.db` unless S00 overturns
- No silent delete of operator files without a documented flag
- No daemon/HTTP scope creep (Phase 29 owns serve/GUI)
- Tests for warn-on-stray if that ships
- `go test ./internal/...` stays green

## Planner gate (P30-00)

Verify scaffold paths + `DR-HANDOFF.md` OPEN + Phase 29 **done** before closing this row.

## P30-00 outcome (2026-08-21)

Gate **PASS**. Phase 29 closed; scaffold S00–S03 complete; DR-HANDOFF OPEN; README + S00 prompts thickened for independent run; S01–S03 light-thickened. Canonical store assumption unchanged for S00 re-verify. No product code.

## Next

`P30-S00-00`
