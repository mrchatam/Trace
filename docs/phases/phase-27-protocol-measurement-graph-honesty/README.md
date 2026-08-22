# Phase 27 — Protocol measurement + graph honesty (P25-D + P25-E)

Human-promoted successor after Phase 26 close. Phase 26 shipped P25-A/B and closed the P25-2 installer gap; verify left **P25-3 FAIL** on build-only G1 (thin graph) — measurement and export-honesty work deferred from Phase 24.

## Goal

Strengthen experiment protocol and seed export gates so build-only vs directed-gap sessions are scored distinctly and thin graphs are caught at export/score time — without changing live agent behavior alone.

## Evidence basis

- [Phase 26 DR-HANDOFF](../phase-26-loop-implementation/DR-HANDOFF.md) (CLOSED; successor = Phase 27)
- Phase 26 verify: P25-2 PASS (closure); P25-3 FAIL (`discoveries=0 decisions=0` on build-only G1)
- Phase 24 deferred themes: [INTERVENTION-MATRIX.md](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) — **P25-D** (INT-08/10), **P25-E** (INT-07)

## Scope sequence (initial)

```
S00 investigation → S01 protocol v2 → S02 graph honesty → S03 VERIFY
```

| Scope | Theme | INT |
|-------|--------|-----|
| S00 | Audit protocol + export paths vs Phase 26 outcome | INT-07/08/10 |
| S01 | Experiment protocol v2 (two-session rubric, score.sh fixes) | INT-08, INT-10 |
| S02 | `seed export --strict` graph honesty enforcement | INT-07 |
| S03 | VERIFY + DR-HANDOFF close | D1–D6 TBD at S00 |

**Default board order is serial** (S00→…→S03).

## In scope

- Investigation + plan artifacts under this phase folder
- Experiment harness / protocol docs under `experiments/`
- Product changes to `seed export`, scoring scripts, and rubric as named by audit/plan
- Tests keeping `go test ./internal/...` green

## Out of scope

- Daemon / HTTP / hosted service on Trace core
- Rewriting Phase 25/26 `done` board history
- Re-scoring E01 historical results
- Replacing SQLite with export-as-SoT

## Board

[`docs/TODO/phase-27.md`](../../TODO/phase-27.md) — first runnable after phase planner: **P27-S00-00**.
