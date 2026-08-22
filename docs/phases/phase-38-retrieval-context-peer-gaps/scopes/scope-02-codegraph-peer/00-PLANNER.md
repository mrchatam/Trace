# P38-S02-00 — Scope planner (Codegraph peer)

## Metadata
- id: P38-S02-00
- todo_ids: [P38-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, research]
- mcps: [user-codegraph]
- verification: automated

## Objective

Lock S02 Codegraph peer deep-dive. Output **`PEER-CG.md`**. Compare explore tool, watch/index, MCP ergonomics vs Trace (H1 partial, H5, H6, H7). **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- S00 `INVESTIGATION-INDEX.md` (after S00-02 APPROVE)
- P24 [EXTERNAL-RESEARCH.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Peer root | `similar projects/codegraph/` |
| Optional | Codegraph MCP `codegraph_explore` if `.codegraph/` on Trace or sample |
| Artifact | `PEER-CG.md` |
| No | Copying CG into Trace; implement trace_explore |

## Must answer for 01

1. Mechanism of `codegraph_explore` (inputs, outputs, blast radius).
2. Index/watch vs Trace `trace index`.
3. Single-tool UX vs Trace 16-tool — agent discovery implications (observation).
4. P24 transfer items — still deferred? evidence.
5. What Trace should **not** adopt (Law 6/7, daemon).

## Planner gate

- [x] `01-investigate.md` has ordered investigation todos T0–T9 (multiple)
- [x] `02-review.md` requires mechanism cites from peer repo (Checklist B peer mechanism gate)
- [x] SCOPE-TODOS IDs 653–655 accurate

## Exit criteria

- [x] S02-01/02 prompts thickened against live peer repo (spot-check: tools.ts L1163 explore schema, L1275 DEFAULT_MCP_TOOLS, L3017 blast radius, L3193 handleExplore; watcher.ts L1–80; daemon.ts L1–25; no .codegraph on Trace; P24 CG transfer still deferred)
- [x] Board `P38-S02-00` → `done`

## Next

`P38-S02-01`
