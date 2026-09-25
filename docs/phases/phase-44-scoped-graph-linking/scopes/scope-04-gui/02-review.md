# Phase 44 / S04 / Review — GUI

## Metadata
- id: P44-S04-02
- todo_ids: [P44-S04-02]
- role: reviewer
- skills: [code-review, accessibility-auditing, auditing-performance]
- verification: mixed

## Objective

Independent review of S04 GUI against **LOCKED** [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §6 + **D5**, implement acceptance in [`01-implement.md`](01-implement.md), and **a11y + perf at 500**. Confirm Law 19 (no retrieval/membership walk in `web/`). Fresh subagent — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-implement.md`](01-implement.md) — cluster / priority / dashed / fixture acceptance
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — §3.4, §5–§6, D5
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- Live: `web/src/lib/graphLayout.ts`, `overviewCompose.ts`, `Graph.tsx`, `GraphOrientPanel.tsx`

## Session start

Fresh subagent. Follow agent-loop-protocol. Re-read board Notes on **P44-S04-01** before judging. Prefer browser evidence for clusters + a11y; unit tests for priority/cap.

## Review rubric

### G1 — Clusters & acceptance

| Check | Pass |
|-------|------|
| Layout clusters by **API `scope_id`** (missing → ungrouped); no client membership walk | Yes |
| Smoke / Notes: **≥2** visible scope clusters with seeded data | Yes |
| **≥1** cross-scope edge visible (prefer `api_contract` in priority set) | Yes |
| Zero-`scope_id` payload does not invent fake clusters | Yes |

### G2 — Edges & provenance

| Check | Pass |
|-------|------|
| MVP rels elevated: `scope_member`, `api_contract`, `implements`, `blocks` | Yes |
| When Tier A present, `goal_has_task` does not dominate overview cap | Yes |
| `EDGE_OVERVIEW_MAX = 150` unchanged | Yes |
| `provenance=inferred` → dashed / distinct class; solid for explicit/omitted | Yes |
| Stored `rel` strings **not** prefixed `inferred_` | Yes |

### G3 — Caps & Law 6–7 / D5

| Check | Pass |
|-------|------|
| `PROJECT_MAX_NODES` / `UI_CAP` still **500** | Yes |
| Project fetch still passes explicit `max_nodes`; no unbounded client fetch | Yes |
| Truncation / omitted banner still honest | Yes |
| Optional `scope` query is server filter only (thin client) | Yes |

### G4 — Law 19 / M-001

| Check | Pass |
|-------|------|
| No graph walk / scope filter logic forked into `web/` beyond presentation | Yes |
| No SQLite / second SoT from browser | Yes |
| Orient copy still moat-first (Tasks → Loop → Gate → Review) | Yes |

---

## A11y checklist (mandatory)

| ID | Check | Pass |
|----|-------|------|
| A1 | Orient / legend explains scopes + dashed=inferred in **text** (not color-only) | Yes |
| A2 | Dash pattern (or pattern + label) distinguishes inferred; color alone insufficient | Yes |
| A3 | Orient panel keeps accessible name (`aria-labelledby` / region); dismiss control labelled | Yes |
| A4 | Truncation / budget messaging remains readable to SR / visible text | Yes |
| A5 | Graph nodes remain keyboard-focusable where they were (no regression) | Yes |
| A6 | Contrast of dashed edges + legend text acceptable on default theme | Yes |

Evidence: browser aria snapshot and/or Notes with `data-testid` cites (`graph-orient-panel`, legend copy).

---

## Perf checklist at 500 (mandatory)

| ID | Check | Pass |
|----|-------|------|
| P1 | Default project fetch still `max_nodes=500` — no silent bump toward API 5000 | Yes |
| P2 | Force/scope layout bounded (iterations / strength constants documented or unchanged-order); no O(n²) **extra** client membership inference | Yes |
| P3 | Overview edge cap 150 still applied in project mode | Yes |
| P4 | No full-graph dump UI path; scope filter does not disable `max_nodes` | Yes |
| P5 | Unit tests cover tiered priority + layout helpers without requiring 5000-node browser run | Yes |

Spot-check: `rg PROJECT_MAX_NODES web/` → still 500; no new default >500.

---

## Evidence expectations

- [ ] Re-run or cite implementer web unit tests (`graphLayout.test.ts` / npm script)
- [ ] Browser or DOM Notes for ≥2 clusters + cross-scope edge
- [ ] Rubric G1–G4 + A1–A6 + P1–P5 all PASS, or spawn forward with severity

## Exit criteria

- [ ] **APPROVE** with confidence medium/high + evidence, **or** spawn `P44-S04-02a/02b` for blocker/high
- [ ] SCOPE-TODOS S04-2 / S04-C updated via board Notes as appropriate
- [ ] Next **P44-S05-00**

## Minimal todos

- [ ] Diff vs `01-implement.md` + locks
- [ ] a11y + perf checklists evidenced
- [ ] Board Notes (APPROVE / FAIL + spawn)

## Next

`P44-S05-00`
