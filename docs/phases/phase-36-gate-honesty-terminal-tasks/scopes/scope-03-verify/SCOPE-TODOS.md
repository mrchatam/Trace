# Scope 03 — board map

**S03 VERIFY** — live feet-seller + greenfield agent path + DR-HANDOFF. Serial: **P36-S03-00 → P36-S03-01 → P36-S03-02**.

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 631 | P36-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock verify blocks — **done** |
| 632 | P36-S03-01 | [01-verify.md](01-verify.md) | Implementer | `VERIFY-NOTES.md` + evidence |
| 633 | P36-S03-02 | [02-dr-handoff.md](02-dr-handoff.md) | Closer | DR-HANDOFF CLOSED + successor |

## Key verify dimensions (locked S03-00)

| # | Dimension | Evidence path |
|---|-----------|---------------|
| 0 | S02 scoped `go test` re-check | `experiments/runs/…/evidence/00-s02-scoped-tests.txt` |
| 1 | Greenfield CLI bootstrap → edit gate pass | `…/evidence/01-greenfield/` |
| 2 | Feet-seller Step1 `--for done` terminal advisory (pre-bootstrap) | `…/02-feet-step1-done-gate-pre-bootstrap.json` |
| 3 | Feet-seller Loop112 `--for done` same | `…/03-feet-loop112-done-gate-pre-bootstrap.json` |
| 4 | GUI TaskDetail warn not error on DONE | `…/04-gui/` |
| 5 | Active IN_PROGRESS task still `plan_missing` block | `…/05-active-work/` |
| 6 | Live `plan bootstrap --goal` + PlanExists | `…/06-post-bootstrap-*.json` |
| 7 | Residuals + successor prep for S03-02 | VERIFY-NOTES §residuals |

**Mutation policy:** blocks 2–4 read-only on feet-seller; **block 6 only** mutates fixture (PLAN §2.8).

**Pinned optional:** `docs/verification/phase-36-gate-honesty/`

## Locked fixture

| Item | Value |
|------|-------|
| Path | `/home/ali/Desktop/feet seller telegram app` |
| Goal | `353b12a4-57dd-4d68-8379-b2024e064733` |
| Step 1 | `33247e2d-aa10-4b25-b194-4b7afb5a6359` |
| Loop 112 | `99d8fb92-65ac-462c-82c4-21bcf198c09e` |
| Binary | `go build -o /tmp/trace ./cmd/trace` |

## Expected terminal gate JSON (post-S02)

- `allowed: true`
- `violations[0].reason_code: goal_plan_gap_terminal_advisory`
- CLI exit **0** — not pre-S02 `plan_missing` block

## Successor default (S03-02)

**`no successor`** unless VERIFY exposes blocking gap (bootstrap failure, active PLAN regression, terminal honesty fail).

## Deferred (not VERIFY fail)

- PlanExists bridge (§2.4)
- HTTP POST plan routes
- MCP `trace_loop gate`
- Enforce default flip
- Bootstrap help refinement note (S02-02 low)
