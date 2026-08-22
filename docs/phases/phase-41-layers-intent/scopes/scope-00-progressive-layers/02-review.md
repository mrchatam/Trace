# P41-S00-02 — Review (G8 progressive layers)

## Metadata
- id: P41-S00-02
- todo_ids: [P41-S00-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter, security-and-hardening]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Fresh independent review of S00-01 G8 implementation vs REMEDIATION-PLAN G8, GAP-REGISTRY G-003, M-001 moat, Laws 6–7/19.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G8-L1–L7 acceptance map + touch-list
- [00-PLANNER.md](00-PLANNER.md) — locks + live repo gap
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [REMEDIATION-PLAN G8](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-003](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Pre-ship baseline: [h3-layers-designed-vs-shipped.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md)
- P41-S00-00 live re-verify (2026-08-22): no drift — `ContextOptions` lacks `MaxLayer`; `compiler_test.go:60–61` still asserts default layer≤1; CLI/MCP depth-only

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b below this row |
| Ship vs spec-revise | **Ship expected** — spec-revise acceptable only with documented ADR + §4 revise |
| Default behavior | L0–L1 only unless `max_layer` opt-in |
| Cap defaults | 4096/32/64 unchanged unless explicit scoped L2/L3 note |

## Review checklist

### A — G8 gap closure

- [ ] `max_layer` opt-in (default 1) on compiler path
- [ ] Live JSON can show `item.layer` 2 and/or 3 when opted in (G8-L2, G8-L3)
- [ ] L2/L3 content aligns RETRIEVAL_AND_CONTEXT §4 (not arbitrary relabel of L1 hits)
- [ ] Reason codes documented for L2/L3 paths (reuse `graph_neighbor`/`recent_event`/`historical_vcs` where aligned)
- [ ] `compiler_test.go` default-path layer≤1 assertion preserved; G8-L1–L7 cover opt-in
- [ ] `doc.go` defer line updated to shipped honesty

### B — M-001 moat

- [ ] Default `trace context` / MCP context unchanged for moat path (G8-L1)
- [ ] Task UUID + Layer-0 core always present
- [ ] No query-only layer expansion path
- [ ] Layer expansion merges into task packet — not standalone dump

### C — Laws 6–7

- [ ] L2/L3 subject to same budget caps (G8-L4)
- [ ] Trim priority L0 > L1 > L2 > L3 (G8-L5)
- [ ] No full-graph dump (G8-L7, `TestNoDumpAPI`)
- [ ] `--depth` independent of layer (G8-L6)

### D — Law 19

- [ ] Admission/compile logic in `internal/compiler/` (library first)
- [ ] CLI/MCP thin — pass-through `max_layer` only

### E — Tests

- [ ] G8-L1–L7 evidenced green
- [ ] G1/G2 regression green
- [ ] `go test ./internal/compiler/... ./internal/retrieval/... ./internal/mcp/... -count=1` passes

### F — Rejects

- [ ] No auto-load L2/L3 on default context
- [ ] No semantic/vector channel (G-004a)
- [ ] No default cap inflation
- [ ] No conflating `--depth` with progressive layer

### G — Live verification commands

```bash
# Default — layer max 1
trace context <task-id> --format json | jq '.layer, [.items[].layer] | unique'

# Opt-in L2
trace context <task-id> --max-layer 2 --format json | jq '.layer, [.items[].layer] | unique'
```

## Exit criteria

- [ ] APPROVE or REJECT with confidence in board Notes
- [ ] Zero open blocker/high without pending spawn
- [ ] Next **P41-S01-00** on APPROVE

## Next

`P41-S01-00` (on APPROVE)
