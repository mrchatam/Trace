# P38-S03-00 — Scope planner (UA + Graphify + Mempalace peers)

## Metadata
- id: P38-S03-00
- todo_ids: [P38-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, research]
- verification: automated

## Objective

Lock S03 UA + Graphify + **Mempalace** investigation. Output **`PEER-UA-GF.md`** (§1 UA · §2 GF · §3 MP). H1 (partial), H4, H8 + MP slices H6/H9. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- S00 `INVESTIGATION-INDEX.md` (after S00-02 APPROVE)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Peer | Path |
|------|------|
| UA | `similar projects/Understand-Anything/` — `context-builder.ts`, SearchEngine, `onboard-builder.ts` |
| Graphify | `similar projects/graphify/` — `validate.py`, `symbol_resolution.py`, `graph.html`, worked examples |
| **Mempalace** | `similar projects/mempalace/` — `searcher.py`, `layers.py`, `mcp_server.py`, `knowledge_graph.py`, `fact_checker.py` |

Artifact: `PEER-UA-GF.md`. Read-only. Sample paths OK — no full monorepo scan.

**Mempalace:** human-added peer 2026-08-22; integrated into existing H* (no H12+ at planner).

## Must answer for 01

1. UA `context-builder.ts`: query+task packet shape vs Trace compiler.
2. UA SearchEngine: query-driven retrieval mechanism (file:line).
3. Graphify EXTRACTED vs INFERRED edges — relevance to H4 (semantic/concept).
4. Graphify `graph.html` / worked examples — orient/onboarding UX for H8.
5. Mempalace hybrid search (vector + BM25) vs Trace FTS-only — H4 + DR-NOSSEM.
6. Mempalace memory stack + MCP surface vs Trace packet/tools — H1/H6/H8 partial.
7. Mempalace fact_checker vs Trace intent docs — H9 contrast partial.
8. Moat seed: what UA/GF/MP lack (task loop, gates, evidence).

## Planner gate

- [x] `01-investigate.md` has ordered investigation todos T0–T11 (multiple)
- [x] `02-review.md` requires file:line cites; no fake cites if peer missing
- [x] SCOPE-TODOS IDs 656–658
- [x] PEER-FIXTURES.md + INVESTIGATION-INDEX updated for Mempalace

## Exit criteria

- [x] S03-01/02 prompts runnable alone
- [x] Board `P38-S03-00` → `done`

## Next

`P38-S03-01`
