# INTAKE — Retrieval & context peer-gap investigation

**Human-promoted 2026-08-22.** Phase 37 addresses **P36 residuals** (implement). Phase 38 is **separate**: deep investigation of Trace vs peer projects to identify gaps and potential improvements — **investigate and plan only; no implement.**

## Trigger

Conversation + P24 external research surfaced likely gaps vs **Codegraph**, **Understand Anything (UA)**, **Graphify**, and related peers: query-driven context, code exploration UX, semantic retrieval, index freshness, MCP ergonomics, progressive context layers, graph onboarding. These are **hypotheses** until Phase 38 investigation proves them with live evidence.

## Human locks

| Lock | Value |
|------|-------|
| Phase 37 | **CLOSED** (2026-08-22) — P38 promoted same day; do not reopen P37 rows |
| Phase 38 mode | **Investigation loops → saturation gate → remediation PLAN** |
| Implement | **Forbidden** in P38 — successor phase (human-promoted) owns build |
| Peer use | Encouraged — local clones under `similar projects/` + optional tools (Trace CLI/MCP, Codegraph MCP, Graphify/UA docs) **read-only** |
| Exit investigation | Only when saturation reviewer **confident** no high-value investigation room remains (document rejects) |

## Hypothesis backlog (H* — verify, do not treat as implement list)

| ID | Hypothesis | Peers | Trace touch areas |
|----|------------|-------|-------------------|
| H1 | No unified **query + task** context packet (auth-style orient) | UA `context-builder`, CG `codegraph_explore` | `compiler`, MCP surface |
| H2 | Context compiler FTS uses **task title** only, not agent query | UA SearchEngine(query) | `internal/compiler/compiler.go` |
| H3 | **Layer 2–3** context designed but not shipped (P0-X) | Aider repo map depth | `compiler`, docs |
| H4 | **No semantic** retrieval channel (DR-NOSSEM) limits concept discovery | Graphify docs path | `retrieval` |
| H5 | Code index: **3 langs** vs CG 20+; manual `trace index` vs watcher | CG auto-sync | `analyzers`, `cmd/trace/index` |
| H6 | MCP **16 tools** vs CG **1 explore** → discovery paralysis | CG, CM | `internal/mcp` |
| H7 | **`trace_explore` unified read** missing (P24 transfer deferred) | CG, P24 EXTERNAL-RESEARCH | MCP + library |
| H8 | Explore GUI vs Graphify/UA **onboarding graph** (human hook) | Graphify `graph.html`, UA viewer | `web/`, Phase 32–33 |
| H9 | Intent extraction pipeline in docs **not implemented** | RETRIEVAL_AND_CONTEXT.md | `retrieval`, `compiler` |
| H10 | Trace moat (tasks/gates/evidence) **under-promoted** in install vs peer orient-first | OH, SWE, CG | `install`, harness |
| H11 | Combining Trace + Codegraph **undocumented** as recommended stack | User dogfood | docs only? |

Investigation may **spawn new H*** rows; S03 matrix is SoT.

## Desired outcomes

1. **`GAP-REGISTRY.md`** — evidence-backed gaps + non-gaps + deferred ideas with severity and Trace-law fit  
2. **`SATURATION-NOTES.md`** — explicit “no more investigation room” rationale or spawn list  
3. **`REMEDIATION-PLAN.md`** — prioritized **future phases** / scopes (not P38 code)  
4. **No product commits** except optional read-only evidence under `experiments/runs/`

## Not in scope (P38)

- Implementing H* fixes  
- Reopening Phase 36/37 deliverables  
- Hosted SaaS / daemon defaults  
- Re-running closed P24 rows verbatim (cite P24; extend with live re-verify)
