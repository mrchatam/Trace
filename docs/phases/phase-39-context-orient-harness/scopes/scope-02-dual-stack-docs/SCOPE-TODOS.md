# Scope 02 — board map (G4)

**S02 dual-stack docs** — G-011 doc-only. Serial: **P39-S02-00 → P39-S02-01 → P39-S02-02**.

| Order | Board ID | Prompt | Role | Status |
|------:|----------|--------|------|--------|
| 678 | P39-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | **done** |
| 679 | P39-S02-01 | [01-implement.md](01-implement.md) | Implementer (docs) | pending |
| 680 | P39-S02-02 | [02-review.md](02-review.md) | Reviewer | pending |

## Planner locks (P39-S02-00 — verified 2026-08-22)

| Lock | Value |
|------|-------|
| Theme | G4 dual-stack documentation |
| Verdict | **Accept doc-only**; reject product integration |
| Files | `CONTRIBUTING.md`, `AGENTS.md` (+ optional README pointer) |
| Checklist | G4-D1–D8 (locked in implement/review prompts) |
| S01 boundary | Moat-first + reload + 9/16 in Agent workflow (`:64–72`); `instructions.go` pointer-only — S02 adds full recipe section |
| Out | Product dual-index default; bundled MCP; Trace core `.codegraph/` access |

## Live-repo spot-check (P39-S02-00)

| Anchor | Status |
|--------|--------|
| CONTRIBUTING moat-first orient (`:68–69`) | confirmed — S01 shipped; **do not duplicate** |
| CONTRIBUTING MCP reload + 9/16 (`:72`) | confirmed — includes dual-stack **placeholder** pointer |
| CONTRIBUTING dual-stack section | **missing** — S02-01 creates `## Trace + Codegraph (optional dual-stack)` |
| AGENTS.md Codegraph subsection | **missing** — S02-01 adds under Agent workflow |
| `instructions.go` Codegraph complement (`:23–26`) | confirmed — pointer-only to CONTRIBUTING; **S02 non-touch** |
| `.trace/` vs `.codegraph/` on Trace repo | no `.codegraph/` (PEER-CG T0 evidence) — doc describes optional per-repo CG |
| PEER-CG §5 + PEER-FIXTURES link targets | paths exist under `docs/phases/phase-38-retrieval-context-peer-gaps/` |

## G4-D1–D8 checklist (authority)

| ID | Requirement |
|----|-------------|
| G4-D1 | Section title: Trace + Codegraph **complementary**, not merged |
| G4-D2 | Storage: `.trace/` vs `.codegraph/` — separate indexes |
| G4-D3 | When Trace: task loop, gates, plan, evidence, task packet |
| G4-D4 | When Codegraph: symbols, call paths, blast — **optional** |
| G4-D5 | Law 19: adapter per store; Trace core no `.codegraph/` index |
| G4-D6 | Setup: independent `trace index` + optional `codegraph init` |
| G4-D7 | Reject: no default dual-index, no bundled dual MCP |
| G4-D8 | Link PEER-CG §5 + PEER-FIXTURES |

## Next

**P39-S02-01** — implement G4 docs per thickened [01-implement.md](01-implement.md).
