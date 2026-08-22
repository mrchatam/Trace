# P41-S00-01 — Implement (G8 progressive layers)

## Metadata
- id: P41-S00-01
- todo_ids: [P41-S00-01]
- role: implementer
- skills: [backend-dev, context-engineering, test-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Implement **G8**: ship progressive layers L2–L3 in compiler (G-003). M-001: merges into task loop; Laws 6–7 caps preserved; L2–L3 opt-in only.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G8](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-003](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [RETRIEVAL_AND_CONTEXT.md §4](../../../../RETRIEVAL_AND_CONTEXT.md)
- Live anchors (P41-S00-00 re-verified 2026-08-22 — **no drift vs P41-00**):
  - `internal/compiler/doc.go:7` — L2–L3 defer; update when shipped
  - `internal/compiler/packet.go:44,97` — `Item.Layer` / `Packet.Layer` capped 0–1 today
  - `internal/compiler/compiler.go:14–21` — `ContextOptions` (no `MaxLayer` yet)
  - `internal/compiler/compiler.go:79–329` — `compileAtDepth`; all admits layer 0–1
  - `internal/compiler/compiler.go:376–391` — `layer1AdmitKey` only
  - `internal/compiler/budget.go:28–59` — trim prefers L0 > L1 (extend to L2/L3)
  - `internal/retrieval/types.go:10–24` — reuse `graph_neighbor`, `recent_event`, `historical_vcs` for L2/L3 where aligned; document new codes in `retrieval/doc.go`
  - `internal/retrieval/expand.go`, `impact_walk.go`, `impact_bridge.go` — L2/L3 candidate sources
  - `cmd/trace/context.go:18–77` — `--depth` only; add `--max-layer`; extend `flagsFirst` map
  - `internal/mcp/tools_context.go:14–94` — add `max_layer` on `ContextInput` + pass-through
  - `internal/compiler/compiler_test.go:60–61` — default path asserts `pkt.Layer <= 1`; keep for default, add G8-L1–L7 for opt-in

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-003 |
| Verdict | **Accept — ship** (P41-00 lock) |
| Default | `max_layer=1` — L0–L1 unchanged for default `trace context` |
| Opt-in | `max_layer=2` or `3` via CLI/MCP; requires explicit caller request |
| Depth | Graph expand `--depth` 1..2 unchanged; independent of `max_layer` |
| L2 admission | Dependents (impact walk depth 1), discoveries, sibling/future tasks under goal, arch neighbors (import/symbol 1-hop beyond L1) |
| L3 admission | Historical decisions (temporal/ref enrich), cross-module impact (depth-2 walk capped), older evidence |
| Reason codes | Prefer existing: `graph_neighbor` (L2 arch neighbors), `recent_event` (L2 discoveries), `historical_vcs` (L3 temporal); add new only if §4 needs — document in `retrieval/doc.go` |
| Caps | Honor `DefaultTokenBudget=4096`, `DefaultMaxItems=32`, `MaxCandidateHits=64`; L2/L3 items compete in same budget |
| Trim priority | L0 > L1 > L2 > L3; then distance/score |
| Fallback | If ship blocked: spec-revise RETRIEVAL_AND_CONTEXT §4 + ADR — document alternative, no silent defer |
| Graph export | If entities change: `trace seed export -o trace/graph.json` |

## Touch-list (library → adapters → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `internal/compiler/compiler.go` | Add `MaxLayer` to `ContextOptions`; compile L2/L3 admission branches |
| 2 | `internal/compiler/compiler.go` | Add `layer2AdmitKey`, `layer3AdmitKey` (or unified admit with layer return) |
| 3 | `internal/compiler/packet.go` | Allow `Item.Layer` 0–3; `Packet.Layer` = highest included |
| 4 | `internal/compiler/budget.go` | Extend `trimToBudget` layer priority L0→L3 |
| 5 | `internal/compiler/doc.go` | Replace L2–L3 defer with shipped honesty |
| 6 | `internal/retrieval/` (as needed) | Helpers for L2/L3 candidate fetch — no parallel index |
| 7 | `internal/compiler/compiler_test.go` | Keep `TestTaskContextAndBudgets` default layer≤1; add G8-L1–L7 |
| 8 | `cmd/trace/context.go` | `--max-layer 1|2|3` (default 1); wire `opts.MaxLayer`; extend `flagsFirst` |
| 9 | `internal/mcp/tools_context.go` | `max_layer` optional on `ContextInput`; default 1 when omitted |
| 10 | `internal/mcp/mcp_test.go` | `TestMCPContextMaxLayer2` (G8-L2-MCP) |
| 11 | `internal/retrieval/doc.go` | Document any new L2/L3 reason codes |

**Explicit non-touch:**

- `web/` — no GUI layer controls in S00
- `internal/httpapi/` — HTTP mirror optional; not blocking
- G-004a semantic/vector channel
- Default behavior change (must stay L0–L1 unless opted in)
- `trace_explore` — out of scope (G2 shipped Phase 40)

## Implementation order

```text
1. ContextOptions.MaxLayer + admission rules L2/L3
2. Compile path promotes items to layer 2/3 with reason_codes
3. Budget trim layer priority
4. Unit tests G8-L1–L7
5. CLI --max-layer + MCP max_layer thin adapters
6. Update doc.go honesty
7. go test ./internal/compiler/... ./internal/retrieval/... ./internal/mcp/... ./cmd/trace/... -count=1
8. trace seed export if entity schema changes
```

## Acceptance tests (must pass)

| ID | Suggested name | Assert |
|----|----------------|--------|
| G8-L1 | `TestContextDefaultLayer1` | Default compile → `packet.layer` ≤ 1; no L2/L3 items |
| G8-L2 | `TestContextMaxLayer2` | `max_layer=2` → at least one `item.layer==2` with valid reason_code |
| G8-L3 | `TestContextMaxLayer3` | `max_layer=3` → L3 items present when graph has depth; honest empty OK |
| G8-L4 | `TestContextLayerBudgetCap` | L2/L3 subject to same token/item caps; `truncated` when over budget |
| G8-L5 | `TestContextLayerTrimPriority` | L0 kept over L3 when budget tight |
| G8-L6 | `TestContextDepthIndependentOfLayer` | `--depth 2` + `max_layer=1` → no L2 items (depth ≠ layer) |
| G8-L7 | `TestContextNoDump` | No unbounded layer expansion; `MaxCandidateHits` honored |

MCP mirror (if field added):

| ID | Suggested name | Assert |
|----|----------------|--------|
| G8-L2-MCP | `TestMCPContextMaxLayer2` | MCP `max_layer=2` returns layer-2 items |

## Regression tests (must stay green)

- All G1/G2 tests
- `TestNoDumpAPI`
- Default context packet shape unchanged (layer max 1)

## Role work

1. Implement L2/L3 as **opt-in progressive expansion** — not auto-load.
2. Wire CLI/MCP as thin adapters (Law 19).
3. Preserve default L0–L1 behavior for moat path.
4. Self-check G8-L1–L7 before marking row done.

## Exit criteria

- [ ] G8-L1–L7 (+ MCP mirror if applicable) green
- [ ] Default context unchanged; L2/L3 opt-in documented
- [ ] Board row → `done` with files + test command in Notes

## Next

`P41-S00-02`
