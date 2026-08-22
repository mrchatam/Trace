# Peer fixtures — investigation aids (read-only)

Agents **may** use these to understand repos faster during Phase 38. **Do not** modify peer trees unless a row explicitly scopes a throwaway experiment under `experiments/runs/`.

## Local clones (Trace monorepo)

| Peer | Path | Investigate for |
|------|------|-----------------|
| **Codegraph** | [`similar projects/codegraph/`](../../similar%20projects/codegraph/) | `codegraph_explore`, index/watch, MCP single-tool UX, benchmarks README |
| **Understand Anything** | [`similar projects/Understand-Anything/`](../../similar%20projects/Understand-Anything/) | [`context-builder.ts`](../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts), SearchEngine, graph build pipeline |
| **Graphify** | [`similar projects/graphify/`](../../similar%20projects/graphify/) | AST graph, EXTRACTED/INFERRED edges, `graph.html`, worked examples |
| **Mempalace** | [`similar projects/mempalace/`](../../similar%20projects/mempalace/) | Hybrid vector+BM25 search ([`searcher.py`](../../similar%20projects/mempalace/mempalace/searcher.py)), 4-layer memory stack ([`layers.py`](../../similar%20projects/mempalace/mempalace/layers.py)), MCP tool surface ([`mcp_server.py`](../../similar%20projects/mempalace/mempalace/mcp_server.py), [`service.py`](../../similar%20projects/mempalace/mempalace/service.py)), temporal KG ([`knowledge_graph.py`](../../similar%20projects/mempalace/mempalace/knowledge_graph.py)), fact/contradiction check ([`fact_checker.py`](../../similar%20projects/mempalace/mempalace/fact_checker.py)) — contrast Trace task packet, DR-NOSSEM FTS-only retrieval, MCP discovery |
| **Graphiti** | [`similar projects/graphiti/`](../../similar%20projects/graphiti/) | Temporal episodes (contrast only — no daemon transfer) |
| **Codebase Memory MCP** | [`similar projects/codebase-memory-mcp/`](../../similar%20projects/codebase-memory-mcp/) | MCP tool surface density |
| **AgentRQ** | [`similar projects/agentrq/`](../../similar%20projects/agentrq/) | Task dequeue state machine (contrast) |

## Trace internal baselines (do not re-litigate without live re-verify)

| Doc | Use |
|-----|-----|
| [P24 EXTERNAL-RESEARCH.md](../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md) | Prior peer table — **extend**, don’t copy blindly |
| [RETRIEVAL_AND_CONTEXT.md](../../RETRIEVAL_AND_CONTEXT.md) | Design intent vs shipped |
| [internal/compiler/doc.go](../../../internal/compiler/doc.go) | Layer 0–1 shipped |
| [internal/retrieval/doc.go](../../../internal/retrieval/doc.go) | FTS + expand; no semantic |

## Optional tooling (investigator choice)

| Tool | When |
|------|------|
| `trace context`, `trace search`, `trace why`, `trace loop next` | Live Trace behavior on Trace repo or dogfood fixture |
| MCP `trace_context`, `trace_search` | Same via harness |
| **Codegraph MCP** (`codegraph_explore`) | If `.codegraph/` exists on a peer or sample repo — compare orient latency |
| **Graphify** | Read worked example outputs under `similar projects/graphify/worked/` — no full monorepo scan required in P38 |
| **UA** | Read `context-builder.ts` + tests — optional local graph JSON if already built |
| **Mempalace** | Read `searcher.py`, `layers.py`, `mcp_server.py` — optional live MCP if palace initialized; not required in P38 |

## Evidence convention

Store under `experiments/runs/YYYY-MM-DD-p38-<scope>-<row>/evidence/`:

- CLI stdout JSON snippets  
- Side-by-side tables (Trace command vs peer mechanism)  
- File:line cites for both sides  

No secrets; no mutating consumer `.trace/` without row permission.
