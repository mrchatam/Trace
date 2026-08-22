// Package retrieval provides hybrid project-graph lookup without embeddings.
//
// Channels (P0-X): exact ID/path/symbol, SQLite FTS5 lexical search, and
// causal/structural graph expand (default depth 1, hard cap depth ≤ 2).
// Optional temporal enrich uses vcs.Repository refs only (no blobs).
//
// Every Hit carries a reason_code from the locked vocabulary. There is no
// DumpGraph / unbounded export API (G6). Semantic/embedding match is forbidden
// (DR-NOSSEM); do not emit reason_code semantic_match.
//
// G9 intent (pre-channel): when SearchOptions.Intent is set, ExtractIntent builds
// a bounded FTS query from task title/body and optional agent query before lexical
// search. Intent does not add a new reason_code — hits remain fts_match.
//
// G6 graph-label channel (non-semantic): SearchGraphLabels runs bounded FTS over
// concept entity types (discovery, assumption, decision, goal, claim) using G9
// intent terms. Hits carry reason_code graph_label_match — lexical label/summary
// match only (DR-NOSSEM; not semantic similarity).
//
// Extended reason codes beyond the locked set must be documented here.
//
// Progressive layer mapping (G8, compiler admission):
//   - graph_neighbor — L2 architectural neighbors (import/symbol hops, impact walk dependents)
//   - recent_event — L2 recent discoveries / events
//   - historical_vcs — L3 temporal refs (VCS LastChanged; commit entity, not blob)
package retrieval
