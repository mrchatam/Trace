# P38-S01-00 — Scope planner (Trace live audit)

## Metadata
- id: P38-S01-00
- todo_ids: [P38-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, debugging-and-error-recovery]
- mcps: [user-trace]
- verification: automated

## Objective

Lock S01 **Trace live audit**: investigate shipped retrieval/context/MCP/index/GUI vs INTAKE H2,H3,H5,H6,H9,H10. Output **`TRACE-AUDIT.md`**. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- S00 [`INVESTIGATION-INDEX.md`](../scope-00-investigation-index/INVESTIGATION-INDEX.md) (APPROVED P38-S00-02)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- `internal/compiler/`, `internal/retrieval/`, `internal/mcp/`, `docs/RETRIEVAL_AND_CONTEXT.md`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-01-trace-audit/TRACE-AUDIT.md` |
| Method | Live CLI/MCP + file:line — not docs-only |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s01-651/evidence/` |
| Hypotheses | H2, H3, H5, H6, H9, H10 + H1 partial, H8 partial |
| Verdicts | confirmed gap \| not a gap \| inconclusive |
| Dogfood fixture | Trace repo root (`-C` ok) |
| Spawn | Unbounded audit slice → new S01 row before S02 |
| Product edits | **Forbidden** |

## Must answer for 01-investigate

1. Per hypothesis: **confirmed gap / not a gap / inconclusive** + evidence.
2. Layer 0–1 vs 2–3: what is designed vs shipped?
3. FTS/query: what inputs does compiler actually use?
4. MCP tool surface: list + agent discovery friction notes (observation only).
5. Index: langs, manual vs auto, freshness story.
6. Install/harness: how Trace moat is (or isn't) surfaced vs orient-first peers.

## Planner gate

- [x] `01-investigate.md` has ordered investigation todos T0–T9 (multiple)
- [x] `02-review.md` requires live command evidence for major claims (Checklist B)
- [x] SCOPE-TODOS IDs 650–652 accurate
- [x] Do **not** write `TRACE-AUDIT.md` in planner row

## Exit criteria

- [x] S01-01/02 prompts thickened against live repo (spot-check: compiler L146 title FTS, 16 MCP tools, 5 lang ids, no intent in retrieval/)
- [x] Board `P38-S01-00` → `done`

## Next

`P38-S01-01`
