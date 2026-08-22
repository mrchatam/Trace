# Phase 38 — Retrieval & context peer-gap investigation

**Phase planner.** Row `P38-00`. Human promoted **2026-08-22** after P37 CLOSED.

## Metadata
- id: P38-00
- role: planner
- skills: [research, planning-and-task-breakdown, diagnosing-bugs]

## Mission

Scaffold **investigation-heavy** phase: Trace vs Codegraph / UA / Graphify (+ peers). **No implement.** Read [`INTAKE.md`](INTAKE.md), [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md), [`PEER-FIXTURES.md`](PEER-FIXTURES.md).

## Gate

**Do not run P38-00 while Phase 37 is active** unless human explicitly parallel-promotes. Default: P37 CLOSED first.

## Scope sequence (7 scopes — see README)

Investigation loops S00–S05 → plan S06 → verify S07.

## Planner gate (P38-00)

- [x] Phase folder + INTAKE H1–H11 + PEER-FIXTURES + DESIGN-LOCKS
- [x] S00–S07 prompts + SCOPE-TODOS each scope (board 647–670)
- [x] S05 saturation gate documented — blocks S06
- [x] S06 explicitly **plan-only** (REMEDIATION-PLAN.md)
- [x] S01–S03 investigate rows have **multiple todos** each; spawn rules in S00/S05
- [x] Board `docs/TODO/phase-38.md` — all rows pending except P38-00 when done
- [x] TODO.md + AGENTS.md wired (active phase 38)
- [x] No product code

## P38-00 outcome (2026-08-22)

Gate **PASS**. Verified INTAKE H1–H11, PEER-FIXTURES peer paths (CG, UA, Graphify clones present under `similar projects/`), DESIGN-LOCKS saturation/planning scope IDs (fixed S05/S06 numbering drift). All 8 scopes have `00-PLANNER` / `01-*` / `02-*` + `SCOPE-TODOS.md`; board rows 647–670 match prompts. S01/S02/S03 investigate prompts each carry 6–7 ordered todos; S05-02 explicitly blocks S06 until APPROVE; S06 locked plan-only. Spot-checked live repo for investigator anchors: 16 MCP tools (`RegisteredToolNames`), Layers 0–1 shipped / 2–3 not auto-loaded (`internal/compiler/doc.go`), FTS-only retrieval + DR-NOSSEM (`internal/retrieval/doc.go`), four analyzer langs (Go/JS/TS/Python — INTAKE H5 says “3 langs”; S01 must reconcile). Protocol-thickened thin S02–S04 planners + S06-02 review. No product code. Next: **P38-S00-00**.

## Scope checklist (for P38-00 verify)

| Scope | Artifact |
|-------|----------|
| S00 | INVESTIGATION-INDEX.md |
| S01 | TRACE-AUDIT.md |
| S02 | PEER-CG.md |
| S03 | PEER-UA-GF.md |
| S04 | GAP-REGISTRY.md |
| S05 | SATURATION-NOTES.md |
| S06 | REMEDIATION-PLAN.md |
| S07 | VERIFY-NOTES.md + DR-HANDOFF CLOSED |

## Next (when promoted)

`P38-S00-00`
